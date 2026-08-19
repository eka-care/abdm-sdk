"""Verifying and decoding Eka webhooks. Mirrors
go/services/abdm/carecontext/webhook.go.

``carecontext`` is a library, not a server: you own your webhook endpoint and
register its URL with Eka. Inside your handler, call parse_webhook and then
the responder for the event you received.

ABDM retries a webhook it considers unanswered, so callbacks must be
idempotent — deduplicate on the event's transaction/request id before doing
anything with side effects. A non-2xx response makes ABDM retry.
"""

import hashlib
import hmac
import json
import time

__all__ = [
    "EVENT_DATA_FETCH",
    "EVENT_LINK_STATUS",
    "EVENT_DISCOVER",
    "EVENT_LINK_INIT",
    "EVENT_LINK_CONFIRM",
    "BadSignatureError",
    "StaleTimestampError",
    "UnknownEventError",
    "DataFetchEvent",
    "LinkStatusEvent",
    "Identifier",
    "DiscoverEvent",
    "LinkInitEvent",
    "LinkConfirmEvent",
    "parse_webhook",
]

# Webhook event names, verbatim from the ABDM Connect docs. These are
# inconsistent upstream — note that the link-confirm event has no care_ prefix.
EVENT_DATA_FETCH = "abha.hip_data_fetch"
EVENT_LINK_STATUS = "abha.link_care_context"
EVENT_DISCOVER = "abha.care_context_discover"
EVENT_LINK_INIT = "abha.care_context_discover_link_init"
EVENT_LINK_CONFIRM = "abha.context_discover_link_confirm"

# How far a webhook timestamp may drift, in either direction, before rejection.
_SIGNATURE_TOLERANCE_SECONDS = 3 * 60

# Swapped in tests.
_time_now = time.time


class BadSignatureError(Exception):
    """The payload was not signed by Eka, or was altered, or the signature
    header/secret was malformed/missing. Treat as 401."""


class StaleTimestampError(Exception):
    """The signature is valid but its timestamp is outside the tolerance
    window. Treat as 401."""


class UnknownEventError(Exception):
    """The payload is authentic but its event is not one this package
    handles. Distinguishable from the two errors above on purpose: answer
    200 for this one (a webhook endpoint must not fail on events it does
    not yet care about) and 401 for an inauthentic payload."""


class DataFetchEvent:
    """abha.hip_data_fetch: an HIU wants records. Answer with
    Service.respond_to_fetch. The HIU's key material is captured but kept
    private (``_key_info``) — respond_to_fetch uses it so you never handle
    key material yourself."""

    def __init__(
        self,
        transaction_id="",
        abha_address="",
        oid="",
        partner_patient_id="",
        hip_id="",
        care_contexts=None,
        hi_types=None,
        key_info=None,
    ):
        self.transaction_id = transaction_id
        self.abha_address = abha_address
        self.oid = oid
        self.partner_patient_id = partner_patient_id
        self.hip_id = hip_id
        self.care_contexts = care_contexts or []
        self.hi_types = hi_types or []
        self._key_info = key_info or {}

    def event_name(self):
        return EVENT_DATA_FETCH


class LinkStatusEvent:
    """abha.link_care_context: the async outcome of Service.link."""

    def __init__(
        self,
        abha_address="",
        care_context_id="",
        status="",
        error="",
        retry_count=0,
        oid="",
        partner_patient_id="",
        hip_id="",
    ):
        self.abha_address = abha_address
        self.care_context_id = care_context_id
        self.status = status
        self.error = error
        self.retry_count = retry_count
        self.oid = oid
        self.partner_patient_id = partner_patient_id
        self.hip_id = hip_id

    def event_name(self):
        return EVENT_LINK_STATUS


class Identifier:
    """A verified patient identifier supplied during discovery."""

    def __init__(self, type="", value=""):
        self.type = type
        self.value = value


class DiscoverEvent:
    """abha.care_context_discover: a patient is searching your facility for
    their records. Match them in your own system, conservatively — a loose
    match leaks another patient's records — and answer with on_discover."""

    def __init__(
        self,
        abha_address="",
        patient_name="",
        gender="",
        year_of_birth=0,
        identifiers=None,
        request_id="",
        txn_id="",
        oid="",
        partner_patient_id="",
        hip_id="",
    ):
        self.abha_address = abha_address
        self.patient_name = patient_name
        self.gender = gender
        self.year_of_birth = year_of_birth
        self.identifiers = identifiers or []
        self.request_id = request_id
        self.txn_id = txn_id
        self.oid = oid
        self.partner_patient_id = partner_patient_id
        self.hip_id = hip_id

    def event_name(self):
        return EVENT_DISCOVER


class LinkInitEvent:
    """abha.care_context_discover_link_init: the patient chose care contexts
    to link. Generate an OTP, send it, store it against the ref_num you
    return, then answer with on_link_init."""

    def __init__(
        self,
        abha_address="",
        patient=None,
        request_id="",
        txn_id="",
        oid="",
        partner_patient_id="",
        hip_id="",
    ):
        self.abha_address = abha_address
        self.patient = patient
        self.request_id = request_id
        self.txn_id = txn_id
        self.oid = oid
        self.partner_patient_id = partner_patient_id
        self.hip_id = hip_id

    def event_name(self):
        return EVENT_LINK_INIT


class LinkConfirmEvent:
    """abha.context_discover_link_confirm: the patient submitted the OTP.
    Validate token against what you stored for link_ref_number, then answer
    with on_link_confirm. This event carries no txn_id — link_ref_number is
    the correlation key, and matches the ref_num you sent from
    on_link_init."""

    def __init__(
        self,
        link_ref_number="",
        token="",
        request_id="",
        oid="",
        partner_patient_id="",
        hip_id="",
    ):
        self.link_ref_number = link_ref_number
        self.token = token
        self.request_id = request_id
        self.oid = oid
        self.partner_patient_id = partner_patient_id
        self.hip_id = hip_id

    def event_name(self):
        return EVENT_LINK_CONFIRM


def parse_webhook(body: bytes, signature_header: str, secret: str):
    """Verifies an Eka webhook and decodes it into a typed event.

    ``signature_header`` is the raw ``Eka-Webhook-Signature`` header;
    ``secret`` is the signing key from your webhook subscription; ``body``
    must be the exact bytes received, because the signature covers them
    verbatim.

    Raises BadSignatureError, StaleTimestampError or UnknownEventError.
    Treat the first two as 401 and UnknownEventError as 200.
    """
    _verify_signature(body, signature_header, secret)

    try:
        env = json.loads(body)
    except ValueError as e:
        raise BadSignatureError("carecontext: decode webhook envelope: %s" % e)

    if not isinstance(env, dict):
        raise BadSignatureError("carecontext: decode webhook envelope: not an object")

    event = env.get("event", "")
    data = env.get("data") or {}

    if event == EVENT_DATA_FETCH:
        key_info = data.get("key_information") or {}
        return DataFetchEvent(
            transaction_id=data.get("transaction_id") or env.get("transaction_id", ""),
            abha_address=data.get("abha_address", ""),
            oid=data.get("oid", ""),
            partner_patient_id=data.get("partner_patient_id", ""),
            hip_id=data.get("hip_id", ""),
            care_contexts=data.get("care_contexts") or [],
            hi_types=data.get("hi_types") or [],
            key_info=key_info,
        )
    if event == EVENT_LINK_STATUS:
        return LinkStatusEvent(
            abha_address=data.get("abha_address", ""),
            care_context_id=data.get("care_context_id", ""),
            status=data.get("status", ""),
            error=data.get("error", ""),
            retry_count=data.get("retry_count", 0),
            oid=data.get("oid", ""),
            partner_patient_id=data.get("partner_patient_id", ""),
            hip_id=data.get("hip_id", ""),
        )
    if event == EVENT_DISCOVER:
        identifiers = [
            Identifier(type=i.get("type", ""), value=i.get("value", ""))
            for i in (data.get("identifiers") or [])
        ]
        return DiscoverEvent(
            abha_address=data.get("abha_address", ""),
            patient_name=data.get("patient_name", ""),
            gender=data.get("gender", ""),
            year_of_birth=data.get("year_of_birth", 0),
            identifiers=identifiers,
            request_id=data.get("request_id", ""),
            txn_id=data.get("txn_id", ""),
            oid=data.get("oid", ""),
            partner_patient_id=data.get("partner_patient_id", ""),
            hip_id=data.get("hip_id", ""),
        )
    if event == EVENT_LINK_INIT:
        return LinkInitEvent(
            abha_address=data.get("abha_address", ""),
            patient=data.get("patient"),
            request_id=data.get("request_id", ""),
            txn_id=data.get("txn_id", ""),
            oid=data.get("oid", ""),
            partner_patient_id=data.get("partner_patient_id", ""),
            hip_id=data.get("hip_id", ""),
        )
    if event == EVENT_LINK_CONFIRM:
        return LinkConfirmEvent(
            link_ref_number=data.get("linkRefNumber", ""),
            token=data.get("token", ""),
            request_id=data.get("request_id", ""),
            oid=data.get("oid", ""),
            partner_patient_id=data.get("partner_patient_id", ""),
            hip_id=data.get("hip_id", ""),
        )

    raise UnknownEventError("carecontext: unrecognised webhook event: %r" % event)


def _verify_signature(body: bytes, header: str, secret: str) -> None:
    ts, v1 = None, None
    for part in (header or "").split(","):
        if "=" not in part:
            continue
        k, v = part.strip().split("=", 1)
        if k == "t":
            ts = v
        elif k == "v1":
            v1 = v
    if not ts or not v1:
        raise BadSignatureError("carecontext: malformed webhook signature header")

    # An empty secret would key the HMAC with nothing, which anyone can
    # compute. Refuse rather than authenticate every forgery.
    if not secret:
        raise BadSignatureError("carecontext: no webhook secret configured")

    if isinstance(body, str):
        body = body.encode("utf-8")

    mac = hmac.new(secret.encode("utf-8"), digestmod=hashlib.sha256)
    mac.update(ts.encode("utf-8"))
    mac.update(b".")
    mac.update(body)
    expected = mac.hexdigest()

    if not hmac.compare_digest(expected, v1):
        raise BadSignatureError("carecontext: signature mismatch")

    # Timestamp is checked only after the signature, so an attacker cannot
    # use timing here to learn anything about the key.
    try:
        sec = int(ts)
    except ValueError:
        raise BadSignatureError("carecontext: unparseable timestamp")

    drift = _time_now() - sec
    if drift > _SIGNATURE_TOLERANCE_SECONDS or drift < -_SIGNATURE_TOLERANCE_SECONDS:
        raise StaleTimestampError("carecontext: webhook timestamp outside tolerance")

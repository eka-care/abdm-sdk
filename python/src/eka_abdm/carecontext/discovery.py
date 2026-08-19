"""Answering patient-initiated discovery. Mirrors
go/services/abdm/carecontext/discovery.go.
"""

from ..http import HTTPClient, Headers


class DiscoverResult:
    """Answers a discovery request: the patients you matched, or an error
    explaining why you could not."""

    def __init__(self, patients=None, error=None):
        self.patients = patients or []
        self.error = error


class LinkInitResult:
    """Reports that you sent an OTP. ref_num is your reference for this
    linking attempt; it returns as LinkConfirmEvent.link_ref_number, so
    store your OTP against it. otp_expiry is an ISO-8601 timestamp."""

    def __init__(self, ref_num="", otp_expiry="", error=None):
        self.ref_num = ref_num
        self.otp_expiry = otp_expiry
        self.error = error


class LinkConfirmResult:
    """Reports the outcome of OTP validation: the care contexts to link, or
    an error if the OTP was wrong or expired."""

    def __init__(self, patients=None, error=None):
        self.patients = patients or []
        self.error = error


def _post_discovery(config, path, oid, partner_patient_id, hip_id, body):
    http = HTTPClient(config, resolve_token=True)
    http.do(
        "POST",
        path,
        headers=Headers(patient_id=oid, partner_user_id=partner_patient_id, hip_id=hip_id),
        body=body,
    )


def on_discover(config, event, result: DiscoverResult) -> None:
    """Answers an abha.care_context_discover webhook with the unlinked care
    contexts belonging to the patient described in the event.

    Matching is yours: query your own patient records using the event's
    name, gender, year of birth and identifiers. Match conservatively — a
    loose match discloses another patient's records to the requester. Set
    result.error when nothing matched.
    """
    if event is None:
        raise ValueError("carecontext: nil discover event")

    body = {"request_id": event.request_id, "txn_id": event.txn_id}
    if result.error is not None:
        body["error"] = result.error.to_dict()
    else:
        body["patients"] = [p.to_dict() for p in result.patients]

    _post_discovery(
        config,
        "/abdm/v1/care-contexts/on-discover",
        event.oid,
        event.partner_patient_id,
        event.hip_id,
        body,
    )


def on_link_init(config, event, result: LinkInitResult) -> None:
    """Answers an abha.care_context_discover_link_init webhook, telling
    ABDM you have dispatched an OTP.

    Generating the OTP, delivering it, and storing it against
    result.ref_num are all yours — this package holds no state. The OTP
    must survive until the matching link-confirm webhook arrives.
    """
    if event is None:
        raise ValueError("carecontext: nil link init event")

    body = {"request_id": event.request_id, "txn_id": event.txn_id}
    if result.error is not None:
        body["error"] = result.error.to_dict()
    else:
        if result.ref_num:
            body["ref_num"] = result.ref_num
        if result.otp_expiry:
            body["otp_expiry"] = result.otp_expiry

    _post_discovery(
        config,
        "/abdm/v1/care-contexts/discover/link/on-init",
        event.oid,
        event.partner_patient_id,
        event.hip_id,
        body,
    )


def on_link_confirm(config, event, result: LinkConfirmResult) -> None:
    """Answers an abha.context_discover_link_confirm webhook after you have
    validated the OTP.

    Validate event.token against whatever you stored for
    event.link_ref_number. On success send the care contexts to link; on
    failure set result.error. This endpoint correlates by request_id
    alone — the webhook carries no txn_id, and this call sends none.
    """
    if event is None:
        raise ValueError("carecontext: nil link confirm event")

    body = {"request_id": event.request_id}
    if result.error is not None:
        body["error"] = result.error.to_dict()
    else:
        body["patients"] = [p.to_dict() for p in result.patients]

    _post_discovery(
        config,
        "/abdm/v1/care-contexts/discover/link/on-confirm",
        event.oid,
        event.partner_patient_id,
        event.hip_id,
        body,
    )

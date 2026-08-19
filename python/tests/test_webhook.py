import hashlib
import hmac
import json
import time
import unittest

from eka_abdm.carecontext import webhook
from eka_abdm.carecontext.webhook import (
    BadSignatureError,
    DataFetchEvent,
    DiscoverEvent,
    LinkConfirmEvent,
    LinkInitEvent,
    LinkStatusEvent,
    StaleTimestampError,
    UnknownEventError,
    parse_webhook,
)

SECRET = "whsec"


def sign(body: str, secret: str, ts: int) -> str:
    mac = hmac.new(secret.encode("utf-8"), digestmod=hashlib.sha256)
    mac.update(("%d" % ts).encode("utf-8"))
    mac.update(b".")
    mac.update(body.encode("utf-8"))
    return "t=%d,v1=%s" % (ts, mac.hexdigest())


class TestParseWebhookSignature(unittest.TestCase):
    def setUp(self):
        self.body = (
            '{"service":"abdm","event":"abha.link_care_context",'
            '"data":{"care_context_id":"cc-1","status":"LINKED"}}'
        )

    def test_valid(self):
        e = parse_webhook(
            self.body.encode("utf-8"), sign(self.body, SECRET, int(time.time())), SECRET
        )
        self.assertIsInstance(e, LinkStatusEvent)
        self.assertEqual(e.care_context_id, "cc-1")
        self.assertEqual(e.status, "LINKED")

    def test_tampered_body(self):
        sig = sign(self.body, SECRET, int(time.time()))
        with self.assertRaises(BadSignatureError):
            parse_webhook((self.body + " ").encode("utf-8"), sig, SECRET)

    def test_wrong_secret(self):
        sig = sign(self.body, "other", int(time.time()))
        with self.assertRaises(BadSignatureError):
            parse_webhook(self.body.encode("utf-8"), sig, SECRET)

    def test_empty_secret_rejected(self):
        # Signature is genuinely valid for the empty key — still must be rejected.
        sig = sign(self.body, "", int(time.time()))
        with self.assertRaises(BadSignatureError):
            parse_webhook(self.body.encode("utf-8"), sig, "")

    def test_stale_timestamp(self):
        old = int(time.time()) - 10 * 60
        sig = sign(self.body, SECRET, old)
        with self.assertRaises(StaleTimestampError):
            parse_webhook(self.body.encode("utf-8"), sig, SECRET)

    def test_future_timestamp_also_stale(self):
        future = int(time.time()) + 10 * 60
        sig = sign(self.body, SECRET, future)
        with self.assertRaises(StaleTimestampError):
            parse_webhook(self.body.encode("utf-8"), sig, SECRET)

    def test_malformed_header(self):
        with self.assertRaises(BadSignatureError):
            parse_webhook(self.body.encode("utf-8"), "garbage", SECRET)


class TestParseWebhookEventNames(unittest.TestCase):
    # The upstream event names are inconsistent (note the third lacks the
    # care_ prefix), so each is asserted verbatim.
    def test_each_event_decodes_to_the_right_type(self):
        cases = [
            (webhook.EVENT_DATA_FETCH, DataFetchEvent),
            (webhook.EVENT_LINK_STATUS, LinkStatusEvent),
            (webhook.EVENT_DISCOVER, DiscoverEvent),
            (webhook.EVENT_LINK_INIT, LinkInitEvent),
            (webhook.EVENT_LINK_CONFIRM, LinkConfirmEvent),
        ]
        for event_name, want_type in cases:
            body = json.dumps({"service": "abdm", "event": event_name, "data": {}})
            e = parse_webhook(
                body.encode("utf-8"), sign(body, SECRET, int(time.time())), SECRET
            )
            self.assertIsInstance(e, want_type, msg=event_name)
            self.assertEqual(e.event_name(), event_name)

    def test_event_name_constants_are_exact(self):
        self.assertEqual(webhook.EVENT_DATA_FETCH, "abha.hip_data_fetch")
        self.assertEqual(webhook.EVENT_LINK_STATUS, "abha.link_care_context")
        self.assertEqual(webhook.EVENT_DISCOVER, "abha.care_context_discover")
        self.assertEqual(webhook.EVENT_LINK_INIT, "abha.care_context_discover_link_init")
        # No care_ prefix on this one — that is not a typo.
        self.assertEqual(webhook.EVENT_LINK_CONFIRM, "abha.context_discover_link_confirm")

    def test_unknown_event(self):
        body = json.dumps({"service": "abdm", "event": "abha.something_new", "data": {}})
        with self.assertRaises(UnknownEventError):
            parse_webhook(body.encode("utf-8"), sign(body, SECRET, int(time.time())), SECRET)


class TestParseWebhookDataFetch(unittest.TestCase):
    def test_captures_key_material(self):
        body = json.dumps(
            {
                "service": "abdm",
                "event": "abha.hip_data_fetch",
                "data": {
                    "transaction_id": "txn-1",
                    "abha_address": "p@sbx",
                    "oid": "o-1",
                    "partner_patient_id": "pp-1",
                    "hip_id": "hip-1",
                    "care_contexts": ["cc-1"],
                    "hi_types": ["OPConsultation"],
                    "key_information": {
                        "crypto_alg": "ECDH",
                        "curve": "Curve25519",
                        "dh_public_key": {"key_value": "k", "parameters": "p", "expiry": "e"},
                        "nonce": "n",
                    },
                },
            }
        )
        e = parse_webhook(
            body.encode("utf-8"), sign(body, SECRET, int(time.time())), SECRET
        )
        self.assertIsInstance(e, DataFetchEvent)
        self.assertEqual(e.transaction_id, "txn-1")
        self.assertEqual(e.hip_id, "hip-1")
        self.assertEqual(e.care_contexts, ["cc-1"])
        self.assertEqual(e._key_info["nonce"], "n")
        self.assertEqual(e._key_info["dh_public_key"]["key_value"], "k")

    def test_empty_data_falls_back_to_envelope_transaction_id(self):
        body = json.dumps(
            {"service": "abdm", "event": "abha.hip_data_fetch", "transaction_id": "txn-9"}
        )
        e = parse_webhook(
            body.encode("utf-8"), sign(body, SECRET, int(time.time())), SECRET
        )
        self.assertEqual(e.transaction_id, "txn-9")


class TestParseWebhookLinkConfirmFields(unittest.TestCase):
    def test_camel_case_link_ref_number_field(self):
        body = json.dumps(
            {
                "service": "abdm",
                "event": "abha.context_discover_link_confirm",
                "data": {
                    "linkRefNumber": "ref-123",
                    "token": "111111",
                    "request_id": "req-3",
                    "oid": "o-1",
                    "partner_patient_id": "pp-1",
                    "hip_id": "hip-1",
                },
            }
        )
        e = parse_webhook(
            body.encode("utf-8"), sign(body, SECRET, int(time.time())), SECRET
        )
        self.assertIsInstance(e, LinkConfirmEvent)
        self.assertEqual(e.link_ref_number, "ref-123")
        self.assertEqual(e.token, "111111")


if __name__ == "__main__":
    unittest.main()

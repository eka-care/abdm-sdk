import hashlib
import unittest

from abdm_ecdh import decrypt, generate_key_material

from eka_abdm.carecontext.service import CareContextService
from eka_abdm.carecontext.types import Entry
from eka_abdm.carecontext.webhook import DataFetchEvent
from eka_abdm.config import Config
from stub_server import JSONHandler, StubServer


class OnFetchHandler(JSONHandler):
    received = None

    def do_POST(self):
        OnFetchHandler.received = {
            "path": self.path,
            "headers": {k.lower(): v for k, v in self.headers.items()},
            "body": self.read_json(),
        }
        self.write_json(202, {})


def make_event(hiu, **overrides):
    kwargs = dict(
        transaction_id="txn-1",
        oid="oid-1",
        partner_patient_id="pp-1",
        hip_id="hip-1",
        key_info={
            "crypto_alg": "ECDH",
            "curve": "Curve25519",
            "nonce": hiu.nonce,
            "dh_public_key": {"key_value": hiu.x509_public_key, "parameters": "p"},
        },
    )
    kwargs.update(overrides)
    return DataFetchEvent(**kwargs)


class TestRespondToFetch(unittest.TestCase):
    def test_hiu_can_decrypt_using_only_the_published_key_material(self):
        # A second key pair, acting as the requesting HIU, entirely separate
        # from whatever the responder generates internally.
        hiu = generate_key_material()

        with StubServer(OnFetchHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            bundle = b'{"resourceType":"Bundle","id":"b1"}'
            evt = make_event(hiu)
            service.respond_to_fetch(evt, [Entry(care_context_id="cc-1", bundle=bundle)])

        received = OnFetchHandler.received
        self.assertEqual(received["path"], "/abdm/v1/hip/care-context/data/on-fetch")
        body = received["body"]
        self.assertEqual(body["transaction_id"], "txn-1")

        entry = body["entries"][0]
        self.assertEqual(entry["care_context_id"], "cc-1")
        self.assertEqual(entry["media"], "application/fhir+json")
        self.assertEqual(entry["checksum"], hashlib.md5(bundle).hexdigest())

        # Decrypt as the HIU would, using ONLY the key material the
        # responder published in its own request body.
        ki = body["key_information"]
        dec = decrypt(
            encrypted_data=entry["content"],
            sender_nonce=ki["nonce"],
            requester_nonce=hiu.nonce,
            requester_private_key=hiu.private_key,
            sender_public_key=ki["dh_public_key"]["key_value"],
        )
        self.assertEqual(dec.decrypted_data, bundle.decode("utf-8"))

        # Three distinct ABDM headers, so a transposition cannot pass by
        # coincidence.
        hdr = received["headers"]
        self.assertEqual(hdr["x-pt-id"], "oid-1")
        self.assertEqual(hdr["x-partner-pt-id"], "pp-1")
        self.assertEqual(hdr["x-hip-id"], "hip-1")
        self.assertEqual(len({hdr["x-pt-id"], hdr["x-partner-pt-id"], hdr["x-hip-id"]}), 3)

    def test_rejects_empty_entries(self):
        hiu = generate_key_material()
        # No StubServer needed: the empty-entries check raises before any
        # HTTP call is attempted.
        cfg = Config(base_url="http://unused")
        service = CareContextService(cfg)
        with self.assertRaises(ValueError):
            service.respond_to_fetch(make_event(hiu), [])

    def test_rejects_nil_event(self):
        cfg = Config(base_url="http://unused")
        service = CareContextService(cfg)
        with self.assertRaises(ValueError):
            service.respond_to_fetch(None, [Entry(care_context_id="cc-1", bundle=b"{}")])

    def test_rejects_empty_bundle(self):
        hiu = generate_key_material()
        cfg = Config(base_url="http://unused")
        service = CareContextService(cfg)
        with self.assertRaises(ValueError):
            service.respond_to_fetch(
                make_event(hiu), [Entry(care_context_id="cc-1", bundle=b"")]
            )


if __name__ == "__main__":
    unittest.main()

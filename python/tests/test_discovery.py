import unittest

from eka_abdm.carecontext.discovery import DiscoverResult, LinkConfirmResult, LinkInitResult
from eka_abdm.carecontext.service import CareContextService
from eka_abdm.carecontext.types import DiscoveredCareContext, ErrorDetail, Patient
from eka_abdm.carecontext.webhook import DiscoverEvent, LinkConfirmEvent, LinkInitEvent
from eka_abdm.config import Config
from stub_server import JSONHandler, StubServer


class CaptureHandler(JSONHandler):
    received = None

    def do_POST(self):
        CaptureHandler.received = {
            "path": self.path,
            "headers": {k.lower(): v for k, v in self.headers.items()},
            "body": self.read_json(),
        }
        self.write_json(204, {})


class TestOnDiscover(unittest.TestCase):
    def test_sends_matched_patients(self):
        with StubServer(CaptureHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            evt = DiscoverEvent(
                request_id="req-1", txn_id="txn-1", oid="oid-1",
                partner_patient_id="pp-1", hip_id="hip-1",
            )
            result = DiscoverResult(
                patients=[
                    Patient(
                        ref_num="P1", display="Gajendra", hi_type="OPConsultation",
                        care_contexts=[DiscoveredCareContext(ref_num="CC101", display="OP Consult")],
                    )
                ]
            )
            service.on_discover(evt, result)

        received = CaptureHandler.received
        self.assertEqual(received["path"], "/abdm/v1/care-contexts/on-discover")
        body = received["body"]
        self.assertEqual(body["request_id"], "req-1")
        self.assertEqual(body["txn_id"], "txn-1")
        patient = body["patients"][0]
        self.assertEqual(patient["ref_num"], "P1")
        self.assertEqual(patient["hi_type"], "OPConsultation")

        hdr = received["headers"]
        self.assertEqual(hdr["x-pt-id"], "oid-1")
        self.assertEqual(hdr["x-partner-pt-id"], "pp-1")
        self.assertEqual(hdr["x-hip-id"], "hip-1")
        self.assertEqual(len({hdr["x-pt-id"], hdr["x-partner-pt-id"], hdr["x-hip-id"]}), 3)

    def test_error_result_omits_patients(self):
        with StubServer(CaptureHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            evt = DiscoverEvent(
                request_id="r", txn_id="t", oid="oid-err",
                partner_patient_id="pp-err", hip_id="hip-err",
            )
            service.on_discover(evt, DiscoverResult(error=ErrorDetail(code=1000, message="no match")))

        body = CaptureHandler.received["body"]
        self.assertEqual(body["error"]["message"], "no match")
        self.assertNotIn("patients", body)

        hdr = CaptureHandler.received["headers"]
        self.assertEqual(hdr["x-pt-id"], "oid-err")
        self.assertEqual(hdr["x-partner-pt-id"], "pp-err")
        self.assertEqual(hdr["x-hip-id"], "hip-err")


class TestOnLinkInit(unittest.TestCase):
    def test_sends_ref_num_and_otp_expiry(self):
        with StubServer(CaptureHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            evt = LinkInitEvent(
                request_id="req-2", txn_id="txn-2", oid="oid-2",
                partner_patient_id="pp-2", hip_id="hip-2",
            )
            service.on_link_init(
                evt, LinkInitResult(ref_num="temp", otp_expiry="2026-08-19T10:00:00Z")
            )

        received = CaptureHandler.received
        self.assertEqual(received["path"], "/abdm/v1/care-contexts/discover/link/on-init")
        body = received["body"]
        self.assertEqual(body["ref_num"], "temp")
        self.assertEqual(body["otp_expiry"], "2026-08-19T10:00:00Z")

        hdr = received["headers"]
        self.assertEqual(hdr["x-pt-id"], "oid-2")
        self.assertEqual(hdr["x-partner-pt-id"], "pp-2")
        self.assertEqual(hdr["x-hip-id"], "hip-2")
        self.assertEqual(len({hdr["x-pt-id"], hdr["x-partner-pt-id"], hdr["x-hip-id"]}), 3)

    def test_error_result_omits_ref_num_and_otp_expiry(self):
        with StubServer(CaptureHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            evt = LinkInitEvent(request_id="req-1", txn_id="txn-1")
            service.on_link_init(
                evt,
                LinkInitResult(
                    ref_num="should-not-be-sent",
                    otp_expiry="2026-01-01T00:00:00Z",
                    error=ErrorDetail(code=1000, message="no match"),
                ),
            )

        body = CaptureHandler.received["body"]
        self.assertNotIn("ref_num", body)
        self.assertNotIn("otp_expiry", body)
        self.assertEqual(body["error"]["message"], "no match")
        self.assertEqual(body["request_id"], "req-1")


class TestOnLinkConfirm(unittest.TestCase):
    def test_correlates_by_request_id_only_no_txn_id_sent(self):
        with StubServer(CaptureHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            evt = LinkConfirmEvent(
                request_id="req-3", link_ref_number="temp", token="111111",
                oid="oid-3", partner_patient_id="pp-3", hip_id="hip-3",
            )
            service.on_link_confirm(
                evt, LinkConfirmResult(patients=[Patient(ref_num="P1", display="")])
            )

        received = CaptureHandler.received
        self.assertEqual(received["path"], "/abdm/v1/care-contexts/discover/link/on-confirm")
        body = received["body"]
        self.assertEqual(body["request_id"], "req-3")
        self.assertNotIn("txn_id", body)

        hdr = received["headers"]
        self.assertEqual(hdr["x-pt-id"], "oid-3")
        self.assertEqual(hdr["x-partner-pt-id"], "pp-3")
        self.assertEqual(hdr["x-hip-id"], "hip-3")
        self.assertEqual(len({hdr["x-pt-id"], hdr["x-partner-pt-id"], hdr["x-hip-id"]}), 3)

    def test_error_result_omits_patients(self):
        with StubServer(CaptureHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            evt = LinkConfirmEvent(request_id="req-4", link_ref_number="temp", token="000000")
            service.on_link_confirm(
                evt, LinkConfirmResult(error=ErrorDetail(code=1001, message="bad otp"))
            )

        body = CaptureHandler.received["body"]
        self.assertEqual(body["error"]["message"], "bad otp")
        self.assertNotIn("patients", body)


if __name__ == "__main__":
    unittest.main()

import unittest

from eka_abdm.carecontext import CareContext, CareContextService, Headers, LinkRequest
from eka_abdm.config import Config
from stub_server import JSONHandler, StubServer


class LinkHandler(JSONHandler):
    received = None

    def do_POST(self):
        LinkHandler.received = {
            "path": self.path,
            "headers": {k.lower(): v for k, v in self.headers.items()},
            "body": self.read_json(),
        }
        self.write_json(202, {})


class TestCareContextLink(unittest.TestCase):
    def test_fills_oid_and_partner_user_id_from_headers_when_empty(self):
        with StubServer(LinkHandler) as server:
            cfg = Config(base_url=server.url)
            cfg.set_authorization_token("tok")
            service = CareContextService(cfg)

            headers = Headers(
                patient_id="oid-from-header",
                partner_user_id="partner-from-header",
                hip_id="hip-1",
            )
            request = LinkRequest(
                abha_address="test@sbx",
                care_contexts=[CareContext(care_context_id="cc-1", display="Visit 1")],
            )

            service.link(headers, request)

        body = LinkHandler.received["body"]
        self.assertEqual(LinkHandler.received["path"], "/abdm/v1/care-contexts/link")
        self.assertEqual(body["oid"], "oid-from-header")
        self.assertEqual(body["partner_user_id"], "partner-from-header")
        self.assertEqual(body["abha_address"], "test@sbx")
        self.assertEqual(
            body["care_contexts"], [{"care_context_id": "cc-1", "display": "Visit 1"}]
        )

    def test_does_not_override_explicit_oid_and_partner_user_id(self):
        with StubServer(LinkHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            headers = Headers(
                patient_id="oid-from-header",
                partner_user_id="partner-from-header",
                hip_id="hip-1",
            )
            request = LinkRequest(
                abha_address="test@sbx",
                care_contexts=[],
                oid="explicit-oid",
                partner_user_id="explicit-partner",
            )

            service.link(headers, request)

        body = LinkHandler.received["body"]
        self.assertEqual(body["oid"], "explicit-oid")
        self.assertEqual(body["partner_user_id"], "explicit-partner")

    def test_does_not_mutate_callers_request_object(self):
        with StubServer(LinkHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            headers = Headers(patient_id="oid-x", partner_user_id="partner-x", hip_id="hip-x")
            request = LinkRequest(abha_address="test@sbx", care_contexts=[])

            service.link(headers, request)

        self.assertEqual(request.oid, "")
        self.assertEqual(request.partner_user_id, "")

    def test_three_abdm_headers_sent_with_distinct_values(self):
        with StubServer(LinkHandler) as server:
            cfg = Config(base_url=server.url)
            service = CareContextService(cfg)

            headers = Headers(
                patient_id="value-pt", partner_user_id="value-partner", hip_id="value-hip"
            )
            request = LinkRequest(abha_address="test@sbx", care_contexts=[])

            service.link(headers, request)

        received = LinkHandler.received["headers"]
        self.assertEqual(received["x-pt-id"], "value-pt")
        self.assertEqual(received["x-partner-pt-id"], "value-partner")
        self.assertEqual(received["x-hip-id"], "value-hip")
        self.assertEqual(
            len({received["x-pt-id"], received["x-partner-pt-id"], received["x-hip-id"]}), 3
        )


if __name__ == "__main__":
    unittest.main()

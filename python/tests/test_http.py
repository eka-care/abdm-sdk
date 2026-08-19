import json
import unittest

from eka_abdm.config import Config
from eka_abdm.http import APIError, Headers, HTTPClient
from stub_server import JSONHandler, StubServer


class RecordingHandler(JSONHandler):
    received = None

    def do_POST(self):
        body = self.read_json()
        # HTTP header names are case-insensitive on the wire; lower-case the
        # keys here so assertions don't depend on urllib's exact casing.
        headers = {k.lower(): v for k, v in self.headers.items()}
        RecordingHandler.received = {"headers": headers, "body": body}
        self.write_json(200, {"ok": True})


class TestHTTPClientHeaders(unittest.TestCase):
    def test_bearer_token_and_abdm_headers_are_sent(self):
        with StubServer(RecordingHandler) as server:
            cfg = Config(base_url=server.url)
            cfg.set_authorization_token("static-tok")
            client = HTTPClient(cfg, resolve_token=True)

            headers = Headers(
                patient_id="patient-1", partner_user_id="partner-2", hip_id="hip-3"
            )
            client.do("POST", "/x", headers=headers, body={"a": 1})

        received = RecordingHandler.received["headers"]
        self.assertEqual(received["authorization"], "Bearer static-tok")
        self.assertEqual(received["x-pt-id"], "patient-1")
        self.assertEqual(received["x-partner-pt-id"], "partner-2")
        self.assertEqual(received["x-hip-id"], "hip-3")
        # Three distinct values, so a transposition would fail this assertion.
        self.assertEqual(
            len({received["x-pt-id"], received["x-partner-pt-id"], received["x-hip-id"]}), 3
        )

    def test_token_is_resolved_per_request_not_at_construction(self):
        with StubServer(RecordingHandler) as server:
            cfg = Config(base_url=server.url)
            client = HTTPClient(cfg, resolve_token=True)
            # No token set yet when the client is built.
            calls = {"n": 0}

            def resolver():
                calls["n"] += 1
                return "token-%d" % calls["n"]

            cfg.set_token_func(resolver)

            client.do("POST", "/x", body={})
            first = RecordingHandler.received["headers"]["authorization"]
            client.do("POST", "/x", body={})
            second = RecordingHandler.received["headers"]["authorization"]

        self.assertEqual(first, "Bearer token-1")
        self.assertEqual(second, "Bearer token-2")

    def test_unauthenticated_client_never_resolves_token(self):
        with StubServer(RecordingHandler) as server:
            cfg = Config(base_url=server.url)

            def resolver():
                raise AssertionError("token resolver must not be called")

            cfg.set_token_func(resolver)
            client = HTTPClient(cfg, resolve_token=False)

            client.do("POST", "/x", body={})

        self.assertEqual(RecordingHandler.received["headers"]["authorization"], "Bearer ")


class ErrorHandler(JSONHandler):
    def do_POST(self):
        self.write_json(
            400,
            {"code": 400, "error": "bad request", "source_error": {"code": "E1", "message": "nope"}},
        )


class TestHTTPClientErrors(unittest.TestCase):
    def test_error_response_raises_api_error_with_status_and_body(self):
        with StubServer(ErrorHandler) as server:
            cfg = Config(base_url=server.url)
            client = HTTPClient(cfg, resolve_token=True)

            with self.assertRaises(APIError) as ctx:
                client.do("POST", "/x", body={})

        err = ctx.exception
        self.assertEqual(err.status_code, 400)
        self.assertIn("bad request", err.message)
        self.assertIn("E1", err.message)


if __name__ == "__main__":
    unittest.main()

"""Regression test for the reentrant self-deadlock the Go SDK hit:

    request -> resolve token -> credentials provider takes its lock ->
    refresh call -> request -> resolve token -> SAME lock, SAME thread ->
    hang forever.

Fixed by never resolving a bearer token for the auth/refresh client (see
eka_abdm.http.HTTPClient(resolve_token=False) and
eka_abdm.credentials.AuthService). This test seeds an already-expired access
token with a still-usable refresh token, then makes a request that must
complete: it is run in a background thread and joined with a timeout, so a
reintroduced deadlock fails the test instead of hanging the whole suite.

Manually confirmed to fail (thread still alive after the timeout) when the
auth service's HTTPClient was built with resolve_token=True instead of
False — see python/tests/DEADLOCK_EVIDENCE in the Task 1 report.
"""

import threading
import time
import unittest

from eka_abdm.carecontext import CareContext, CareContextService, Headers, LinkRequest
from eka_abdm.config import Config
from eka_abdm.credentials import AuthService, ClientCredentialsProvider, Credentials
from stub_server import JSONHandler, StubServer


class AuthAndLinkHandler(JSONHandler):
    def do_POST(self):
        if self.path == "/connect-auth/v1/account/refresh":
            self.write_json(
                200,
                {
                    "access_token": "fresh-access",
                    "refresh_token": "fresh-refresh",
                    "expires_in": 3600,
                    "refresh_expires_in": 7200,
                },
            )
        elif self.path == "/abdm/v1/care-contexts/link":
            auth_header = self.headers.get("Authorization", "")
            if auth_header != "Bearer fresh-access":
                self.write_json(500, {"error": "unexpected token: %s" % auth_header})
                return
            self.write_json(202, {})
        else:
            self.write_json(404, {"error": "unexpected path %s" % self.path})


class TestRequestWithExpiredTokenDoesNotDeadlock(unittest.TestCase):
    def test_completes_within_timeout(self):
        with StubServer(AuthAndLinkHandler) as server:
            cfg = Config(base_url=server.url)
            auth = AuthService(cfg)
            provider = ClientCredentialsProvider(auth, "id", "secret")
            # Access token already expired; refresh token still good, so
            # retrieve() takes the refresh branch, not the login branch.
            provider._cache = Credentials(
                "stale-access", "good-refresh", time.time() - 1, time.time() + 7200
            )
            cfg.set_token_func(lambda: provider.retrieve().access_token)

            service = CareContextService(cfg)
            headers = Headers(patient_id="x", partner_user_id="y", hip_id="z")
            request = LinkRequest(abha_address="x@sbx", care_contexts=[])

            result = {}

            def call():
                try:
                    service.link(headers, request)
                    result["ok"] = True
                except Exception as e:  # pragma: no cover - failure path
                    result["error"] = e

            # daemon=True: if this deadlocks, the thread never finishes: without
            # daemon=True a live non-daemon thread would keep the interpreter
            # (and the whole test suite) from exiting even after join() times
            # out and the assertion below fails it.
            t = threading.Thread(target=call, daemon=True)
            t.start()
            t.join(timeout=5)

            if t.is_alive():
                self.fail("request did not complete within 5s: token refresh deadlocked")

        if "error" in result:
            self.fail("request failed: %r" % (result["error"],))
        self.assertTrue(result.get("ok"))


if __name__ == "__main__":
    unittest.main()

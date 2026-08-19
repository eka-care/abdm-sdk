import time
import unittest

from eka_abdm.credentials import ClientCredentialsProvider, Credentials


class FakeAuthService:
    """Duck-types AuthService without touching the network, so the provider's
    retrieve() ladder can be tested in isolation from HTTP."""

    def __init__(self):
        self.login_calls = 0
        self.refresh_calls = 0
        self.next_login_response = {
            "access_token": "login-access",
            "refresh_token": "login-refresh",
            "expires_in": 3600,
            "refresh_expires_in": 7200,
        }
        self.next_refresh_response = {
            "access_token": "refreshed-access",
            "refresh_token": "refreshed-refresh",
            "expires_in": 3600,
            "refresh_expires_in": 7200,
        }
        self.refresh_should_fail = False

    def client_login(self, client_id, client_secret):
        self.login_calls += 1
        return self.next_login_response

    def refresh_token(self, access_token, refresh_token):
        self.refresh_calls += 1
        if self.refresh_should_fail:
            raise RuntimeError("refresh failed")
        return self.next_refresh_response


class TestCredentialsExpiry(unittest.TestCase):
    def test_expired_true_within_buffer(self):
        # Expires in 4 minutes: inside the 5-minute buffer, so already "expired".
        creds = Credentials("a", "r", time.time() + 240, time.time() + 7200)
        self.assertTrue(creds.expired())

    def test_expired_false_outside_buffer(self):
        creds = Credentials("a", "r", time.time() + 600, time.time() + 7200)
        self.assertFalse(creds.expired())

    def test_can_refresh_requires_token_and_unexpired_refresh_window(self):
        creds = Credentials("a", "r", time.time() + 600, time.time() + 3600)
        self.assertTrue(creds.can_refresh())

        expired_refresh = Credentials("a", "r", time.time() + 600, time.time() - 1)
        self.assertFalse(expired_refresh.can_refresh())

        no_refresh_token = Credentials("a", "", time.time() + 600, time.time() + 3600)
        self.assertFalse(no_refresh_token.can_refresh())


class TestClientCredentialsProviderLadder(unittest.TestCase):
    def test_first_call_logs_in(self):
        auth = FakeAuthService()
        provider = ClientCredentialsProvider(auth, "id", "secret")

        creds = provider.retrieve()

        self.assertEqual(creds.access_token, "login-access")
        self.assertEqual(auth.login_calls, 1)
        self.assertEqual(auth.refresh_calls, 0)

    def test_cached_valid_token_is_reused_without_a_call(self):
        auth = FakeAuthService()
        provider = ClientCredentialsProvider(auth, "id", "secret")
        provider.retrieve()

        creds = provider.retrieve()

        self.assertEqual(creds.access_token, "login-access")
        self.assertEqual(auth.login_calls, 1)
        self.assertEqual(auth.refresh_calls, 0)

    def test_expired_with_usable_refresh_token_refreshes(self):
        auth = FakeAuthService()
        provider = ClientCredentialsProvider(auth, "id", "secret")
        # Pre-seed an expired cache with a still-usable refresh token.
        provider._cache = Credentials(
            "stale-access", "good-refresh", time.time() - 1, time.time() + 7200
        )

        creds = provider.retrieve()

        self.assertEqual(creds.access_token, "refreshed-access")
        self.assertEqual(auth.refresh_calls, 1)
        self.assertEqual(auth.login_calls, 0)

    def test_expired_refresh_token_falls_back_to_login(self):
        auth = FakeAuthService()
        provider = ClientCredentialsProvider(auth, "id", "secret")
        provider._cache = Credentials(
            "stale-access", "stale-refresh", time.time() - 1, time.time() - 1
        )

        creds = provider.retrieve()

        self.assertEqual(creds.access_token, "login-access")
        self.assertEqual(auth.refresh_calls, 0)
        self.assertEqual(auth.login_calls, 1)

    def test_refresh_failure_falls_back_to_login(self):
        auth = FakeAuthService()
        auth.refresh_should_fail = True
        provider = ClientCredentialsProvider(auth, "id", "secret")
        provider._cache = Credentials(
            "stale-access", "good-refresh", time.time() - 1, time.time() + 7200
        )

        creds = provider.retrieve()

        self.assertEqual(creds.access_token, "login-access")
        self.assertEqual(auth.refresh_calls, 1)
        self.assertEqual(auth.login_calls, 1)


if __name__ == "__main__":
    unittest.main()

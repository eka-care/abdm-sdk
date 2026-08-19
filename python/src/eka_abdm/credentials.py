"""Token lifecycle: login/refresh calls plus the caching provider that decides
when to reuse, refresh, or re-login. Mirrors go/auth/service.go and
go/auth/credentials.go.

AuthService.http is built with resolve_token=False (see http.py) — this is
what breaks the reentrant self-deadlock: login and refresh authenticate from
the request body and must never resolve a bearer token, because that
resolution is exactly what calls back into ClientCredentialsProvider.retrieve
while it still holds its own lock.
"""

import json
import threading
import time
from dataclasses import dataclass

from .http import HTTPClient

_EXPIRY_BUFFER_SECONDS = 5 * 60  # matches Go's Credentials.Expired 5-minute buffer


@dataclass
class Credentials:
    access_token: str
    refresh_token: str
    expires_at: float  # unix timestamp (time.time())
    refresh_expires_at: float
    source: str = ""

    def expired(self) -> bool:
        return time.time() > (self.expires_at - _EXPIRY_BUFFER_SECONDS)

    def can_refresh(self) -> bool:
        return bool(self.refresh_token) and time.time() < self.refresh_expires_at


def _credentials_from_response(resp: dict, source: str) -> Credentials:
    now = time.time()
    return Credentials(
        access_token=resp.get("access_token", ""),
        refresh_token=resp.get("refresh_token", ""),
        expires_at=now + resp.get("expires_in", 0),
        refresh_expires_at=now + resp.get("refresh_expires_in", 0),
        source=source,
    )


class AuthService:
    """Calls the two token endpoints. Never resolves a bearer token itself."""

    def __init__(self, config):
        self._http = HTTPClient(config, resolve_token=False)

    def client_login(self, client_id: str, client_secret: str) -> dict:
        raw = self._http.do(
            "POST",
            "/connect-auth/v1/account/login",
            body={"client_id": client_id, "client_secret": client_secret},
        )
        return json.loads(raw) if raw else {}

    def refresh_token(self, access_token: str, refresh_token: str) -> dict:
        raw = self._http.do(
            "POST",
            "/connect-auth/v1/account/refresh",
            body={"access_token": access_token, "refresh_token": refresh_token},
        )
        return json.loads(raw) if raw else {}


class ClientCredentialsProvider:
    """Client-credential login with automatic refresh.

    Retrieve ladder, exactly Go's ClientCredentialsProvider.Retrieve:
    1. Cached token still valid -> return it.
    2. Cached token expired but its refresh token is still usable -> refresh.
    3. Otherwise (no cache, or refresh failed) -> full re-login.
    """

    def __init__(self, auth_service: AuthService, client_id: str, client_secret: str):
        self._auth = auth_service
        self._client_id = client_id
        self._client_secret = client_secret
        self._cache = None  # type: Credentials
        self._lock = threading.Lock()

    def retrieve(self) -> Credentials:
        with self._lock:
            if self._cache is not None and not self._cache.expired():
                return self._cache

            if self._cache is not None and self._cache.can_refresh():
                try:
                    resp = self._auth.refresh_token(
                        self._cache.access_token, self._cache.refresh_token
                    )
                    self._cache = _credentials_from_response(
                        resp, "ClientCredentialsProvider(refresh)"
                    )
                    return self._cache
                except Exception:
                    pass  # refresh failed; fall through to a full re-login

            resp = self._auth.client_login(self._client_id, self._client_secret)
            self._cache = _credentials_from_response(
                resp, "ClientCredentialsProvider(login)"
            )
            return self._cache

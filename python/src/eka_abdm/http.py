"""Minimal HTTP transport, stdlib-only (urllib), mirroring go/internal/http/http.go.

Two behaviours matter here more than anywhere else in the port:

1. The bearer token is resolved PER REQUEST (``config.token()`` inside
   ``HTTPClient.do``), never cached at construction. That is what lets a
   long-running process pick up a refreshed token.

2. ``resolve_token=False`` (used only by the auth/login/refresh client) skips
   that resolution entirely and sends a fixed empty bearer. This breaks the
   reentrant deadlock Go hit: request -> resolve token -> credentials
   provider takes its lock -> refresh call -> request -> resolve token ->
   same lock, same thread -> hang forever. Login and refresh authenticate
   from the request body, not a bearer header, so they never need a token
   and must never call back into the token resolver. See
   eka_abdm/credentials.py and tests/test_deadlock.py.
"""

import json
import urllib.error
import urllib.request
from dataclasses import dataclass
from typing import Any, Optional


@dataclass
class Headers:
    """The ABDM request identifiers sent with every care-context call."""

    patient_id: str = ""
    partner_user_id: str = ""
    hip_id: str = ""


class APIError(Exception):
    """Raised for any 4xx/5xx response. Carries the HTTP status and the
    API's error body (parsed when possible, raw text otherwise)."""

    def __init__(self, status_code: int, message: str):
        self.status_code = status_code
        self.message = message
        super().__init__("HTTP %d: %s" % (status_code, message))


def _format_error_body(raw: bytes, status_code: int) -> str:
    if not raw:
        return "HTTP %d" % status_code
    try:
        parsed = json.loads(raw.decode("utf-8"))
    except (ValueError, UnicodeDecodeError):
        return raw.decode("utf-8", errors="replace")

    if not isinstance(parsed, dict):
        return raw.decode("utf-8", errors="replace")

    code = parsed.get("code", status_code)
    error = parsed.get("error", "")
    source_error = parsed.get("source_error")
    if isinstance(source_error, dict):
        return "Error %s: %s (Source: %s - %s)" % (
            code,
            error,
            source_error.get("code", ""),
            source_error.get("message", ""),
        )
    return "Error %s: %s" % (code, error)


class HTTPClient:
    """Performs JSON-in/JSON-out requests against ``config.base_url``.

    Pass ``resolve_token=False`` to build the unauthenticated variant used by
    the auth service (see module docstring).
    """

    def __init__(self, config, resolve_token: bool = True):
        self._config = config
        self._resolve_token = resolve_token

    def do(
        self,
        method: str,
        path: str,
        headers: Optional[Headers] = None,
        body: Any = None,
    ) -> bytes:
        url = self._config.base_url + path
        data = None
        if body is not None:
            data = json.dumps(body).encode("utf-8")

        req = urllib.request.Request(url, data=data, method=method)

        token = self._config.token() if self._resolve_token else ""
        req.add_header("Authorization", "Bearer %s" % token)
        req.add_header("Content-Type", "application/json")
        req.add_header("User-Agent", self._config.user_agent)

        if headers is not None:
            if headers.patient_id:
                req.add_header("X-Pt-Id", headers.patient_id)
            if headers.partner_user_id:
                req.add_header("X-Partner-Pt-Id", headers.partner_user_id)
            if headers.hip_id:
                req.add_header("X-Hip-Id", headers.hip_id)

        try:
            with urllib.request.urlopen(req, timeout=self._config.timeout) as resp:
                return resp.read()
        except urllib.error.HTTPError as e:
            raw = e.read()
            raise APIError(e.code, _format_error_body(raw, e.code))

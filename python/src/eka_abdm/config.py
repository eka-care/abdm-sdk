"""Client configuration: base URL/environment and the per-request token resolver.

Mirrors go/internal/config/config.go. The important behaviour to preserve is
``Config.token()``: it is consulted on every request (see http.py), not once
at client construction, so a token refreshed mid-flight by a credentials
provider is picked up by a long-running process (e.g. a webhook server)
without rebuilding the client.
"""

import threading
from typing import Callable, Optional

PRODUCTION = "production"
DEVELOPMENT = "development"

_BASE_URLS = {
    PRODUCTION: "https://api.eka.care",
    DEVELOPMENT: "https://api.dev.eka.care",
}


class Environment:
    """Namespace for the two known deployment environments."""

    PRODUCTION = PRODUCTION
    DEVELOPMENT = DEVELOPMENT

    @staticmethod
    def base_url(environment: str) -> str:
        """Returns the base URL for an environment, defaulting to production
        for anything unrecognised (matches Go's Environment.GetBaseURL)."""
        return _BASE_URLS.get(environment, _BASE_URLS[PRODUCTION])


class Config:
    """Holds connection settings and the current bearer token / resolver.

    ``token()`` is thread-safe and cheap: it does not itself perform network
    I/O. A token resolver (set via ``set_token_func``) may do so — that is
    the credentials provider's ``retrieve()`` call.
    """

    def __init__(
        self,
        environment: str = PRODUCTION,
        base_url: Optional[str] = None,
        timeout: float = 30.0,
        user_agent: str = "eka-sdk-python/1.0.0",
    ):
        self.environment = environment
        self.base_url = base_url or Environment.base_url(environment)
        self.timeout = timeout
        self.user_agent = user_agent

        self._lock = threading.Lock()
        self._static_token = ""
        self._token_fn = None  # type: Optional[Callable[[], str]]

    def set_authorization_token(self, token: str) -> None:
        """Sets a static bearer token used for every request until changed."""
        with self._lock:
            self._static_token = token

    def set_token_func(self, fn: Optional[Callable[[], str]]) -> None:
        """Installs a resolver consulted before every request. Use this
        instead of set_authorization_token for long-running processes so
        expired tokens refresh."""
        with self._lock:
            self._token_fn = fn

    def token(self) -> str:
        """Returns the current bearer token, calling the resolver if one is
        installed. Safe to call once per request."""
        with self._lock:
            fn, static = self._token_fn, self._static_token
        if fn is None:
            return static
        return fn()

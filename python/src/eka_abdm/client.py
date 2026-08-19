"""Top-level SDK client. Mirrors go/client.go's New()/Login() surface, scoped
to what Task 1 needs: config, auth, and the care-contexts service.
"""

from typing import Optional

from .carecontext.service import CareContextService
from .config import Config, Environment
from .credentials import AuthService, ClientCredentialsProvider


class Client:
    def __init__(
        self,
        environment: str = Environment.PRODUCTION,
        client_id: str = "",
        client_secret: str = "",
        base_url: Optional[str] = None,
        timeout: float = 30.0,
        user_agent: str = "eka-sdk-python/1.0.0",
    ):
        self.config = Config(
            environment=environment,
            base_url=base_url,
            timeout=timeout,
            user_agent=user_agent,
        )
        self._client_id = client_id
        self._client_secret = client_secret

        self.auth = AuthService(self.config)
        self._credentials_provider = None

        self.care_contexts = CareContextService(self.config)

    def login(self) -> None:
        """Authenticates with the configured client_id/client_secret and
        installs a token resolver so every subsequent request refreshes
        automatically. Fails fast on bad credentials rather than on the
        first API call."""
        if not self._client_id or not self._client_secret:
            raise ValueError("client_id and client_secret are required to log in")

        provider = ClientCredentialsProvider(
            self.auth, self._client_id, self._client_secret
        )
        provider.retrieve()
        self.set_credentials_provider(provider)

    def set_credentials_provider(self, provider) -> None:
        """Installs any object exposing retrieve() -> Credentials as the
        token source, resolved on every request."""
        self._credentials_provider = provider
        self.config.set_token_func(lambda: provider.retrieve().access_token)

    def set_access_token(self, token: str) -> None:
        """Supplies a static token instead of client-credential login. Use
        set_credentials_provider if the token also needs to refresh."""
        self.config.set_authorization_token(token)

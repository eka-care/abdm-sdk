package care.eka.abdm;

import care.eka.abdm.carecontext.CareContexts;
import java.io.IOException;
import java.time.Duration;

/**
 * Top-level SDK client. Mirrors {@code go/client.go}'s {@code New()}/{@code
 * Login()} surface and Python's {@code Client}, scoped to what Task 1 needs:
 * config, auth, and the care-contexts service.
 */
public final class EkaClient {

    private final Config config;
    private final AuthService auth;
    private final CareContexts careContexts;
    private final String clientId;
    private final String clientSecret;

    public EkaClient(String environment, String clientId, String clientSecret) {
        this(environment, clientId, clientSecret, Duration.ofSeconds(30), "eka-sdk-java/1.0.0");
    }

    public EkaClient(String environment, String clientId, String clientSecret, Duration timeout, String userAgent) {
        this.config = new Config(environment, Config.baseUrlFor(environment), timeout, userAgent);
        this.clientId = clientId;
        this.clientSecret = clientSecret;
        this.auth = new AuthService(config);
        this.careContexts = new CareContexts(config);
    }

    /**
     * Authenticates with the configured client id/secret and installs a
     * token resolver so every subsequent request refreshes automatically.
     * Fails fast on bad credentials rather than on the first API call.
     */
    public void login() throws IOException, InterruptedException {
        if (clientId == null || clientId.isEmpty()) {
            throw new IllegalStateException("client id is required for authentication");
        }
        if (clientSecret == null || clientSecret.isEmpty()) {
            throw new IllegalStateException("client secret is required for authentication");
        }
        ClientCredentialsProvider provider = new ClientCredentialsProvider(auth, clientId, clientSecret);
        provider.retrieve();
        setCredentialsProvider(provider);
    }

    /** Installs any {@link CredentialsProvider} as the token source, resolved on every request. */
    public void setCredentialsProvider(CredentialsProvider provider) {
        config.setTokenResolver(() -> provider.retrieve().accessToken());
    }

    /** Supplies a static token instead of client-credential login. Use
     * {@link #setCredentialsProvider} if the token also needs to refresh. */
    public void setAccessToken(String token) {
        config.setAuthorizationToken(token);
    }

    public Config config() {
        return config;
    }

    public AuthService auth() {
        return auth;
    }

    public CareContexts careContexts() {
        return careContexts;
    }
}

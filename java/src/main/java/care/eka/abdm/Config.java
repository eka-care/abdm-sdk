package care.eka.abdm;

import java.io.IOException;
import java.time.Duration;

/**
 * Client configuration: base URL/environment and the per-request token
 * resolver. Mirrors {@code go/internal/config/config.go}.
 *
 * <p>The important behaviour to preserve is {@link #token()}: it is consulted
 * on every request (see {@link HttpTransport}), not once at client
 * construction, so a token refreshed mid-flight by a credentials provider is
 * picked up by a long-running process (e.g. a webhook server) without
 * rebuilding the client.
 */
public final class Config {

    public static final String PRODUCTION = "production";
    public static final String DEVELOPMENT = "development";

    /** Returns the base URL for an environment, defaulting to production for
     * anything unrecognised (matches Go's Environment.GetBaseURL). */
    public static String baseUrlFor(String environment) {
        if (DEVELOPMENT.equals(environment)) {
            return "https://api.dev.eka.care";
        }
        return "https://api.eka.care";
    }

    /** Resolves a bearer token, possibly performing network I/O (a refresh or login call). */
    @FunctionalInterface
    public interface TokenResolver {
        String resolve() throws IOException, InterruptedException;
    }

    private final String environment;
    private final String baseUrl;
    private final Duration timeout;
    private final String userAgent;

    private final Object lock = new Object();
    private String staticToken = "";
    private TokenResolver tokenResolver;

    public Config(String environment, String baseUrl, Duration timeout, String userAgent) {
        this.environment = environment;
        this.baseUrl = baseUrl;
        this.timeout = timeout;
        this.userAgent = userAgent;
    }

    public Config(String environment) {
        this(environment, baseUrlFor(environment), Duration.ofSeconds(30), "eka-sdk-java/1.0.0");
    }

    /** Sets a static bearer token used for every request until changed. */
    public void setAuthorizationToken(String token) {
        synchronized (lock) {
            this.staticToken = token == null ? "" : token;
        }
    }

    /** Installs a resolver consulted before every request. Use this instead
     * of {@link #setAuthorizationToken} for long-running processes so expired
     * tokens refresh. */
    public void setTokenResolver(TokenResolver resolver) {
        synchronized (lock) {
            this.tokenResolver = resolver;
        }
    }

    /** Returns the current bearer token, calling the resolver if one is
     * installed. Safe to call once per request. */
    public String token() throws IOException, InterruptedException {
        TokenResolver fn;
        String stat;
        synchronized (lock) {
            fn = this.tokenResolver;
            stat = this.staticToken;
        }
        if (fn == null) {
            return stat;
        }
        return fn.resolve();
    }

    public String environment() {
        return environment;
    }

    public String baseUrl() {
        return baseUrl;
    }

    public Duration timeout() {
        return timeout;
    }

    public String userAgent() {
        return userAgent;
    }
}

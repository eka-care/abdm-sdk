package care.eka.abdm;

import java.io.IOException;

/**
 * Client-credential login with automatic refresh. Mirrors Go's
 * {@code ClientCredentialsProvider.Retrieve} exactly.
 *
 * <p>Retrieve ladder:
 *
 * <ol>
 *   <li>Cached token still valid -&gt; return it.
 *   <li>Cached token expired but its refresh token is still usable -&gt; refresh.
 *   <li>Otherwise (no cache, or refresh failed) -&gt; full re-login.
 * </ol>
 *
 * <p>{@link #retrieve()} is {@code synchronized}. Java monitors are reentrant,
 * unlike Go's {@code sync.Mutex} — so routing {@link AuthService} through a
 * token-resolving {@link HttpTransport} pointed back at this provider would
 * not deadlock the way Go's did. It would instead recurse: this method calls
 * {@code auth.refreshToken}, which (if misconfigured) resolves a token by
 * calling back into this same synchronized method on the same thread. Java
 * lets the reentry through, so the inner call re-reads the same
 * not-yet-updated, still-expired cache and refreshes again, and so on,
 * until the stack overflows or the refresh endpoint is hit an unbounded
 * number of times. See {@code AuthService}'s docs and
 * {@code DeadlockRegressionTest}.
 */
public final class ClientCredentialsProvider implements CredentialsProvider {

    private final AuthService auth;
    private final String clientId;
    private final String clientSecret;
    private Credentials cache;

    public ClientCredentialsProvider(AuthService auth, String clientId, String clientSecret) {
        this.auth = auth;
        this.clientId = clientId;
        this.clientSecret = clientSecret;
    }

    @Override
    public synchronized Credentials retrieve() throws IOException, InterruptedException {
        if (cache != null && !cache.expired()) {
            return cache;
        }

        if (cache != null && cache.canRefresh()) {
            try {
                cache = auth.refreshToken(cache.accessToken(), cache.refreshToken());
                return cache;
            } catch (Exception e) {
                // Refresh failed; fall through to a full re-login.
            }
        }

        cache = auth.clientLogin(clientId, clientSecret);
        return cache;
    }
}

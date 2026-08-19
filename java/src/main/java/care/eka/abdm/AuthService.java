package care.eka.abdm;

import java.io.IOException;
import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Calls the two token endpoints. Mirrors {@code go/auth/service.go} and
 * Python's {@code AuthService}.
 *
 * <p>Built on an {@link HttpTransport} with {@code resolveToken = false} —
 * this is what breaks the reentrant self-deadlock/recursion: login and
 * refresh authenticate from the request body and must never resolve a bearer
 * token, because that resolution is exactly what calls back into
 * {@link ClientCredentialsProvider#retrieve()} while it still holds its own
 * lock.
 */
public final class AuthService {

    private final HttpTransport http;

    public AuthService(Config config) {
        this.http = new HttpTransport(config, false);
    }

    public Credentials clientLogin(String clientId, String clientSecret) throws IOException, InterruptedException {
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("client_id", clientId);
        body.put("client_secret", clientSecret);
        Map<String, Object> resp = http.doJson("POST", "/connect-auth/v1/account/login", null, body);
        return credentialsFromResponse(resp, "ClientCredentialsProvider(login)");
    }

    public Credentials refreshToken(String accessToken, String refreshToken) throws IOException, InterruptedException {
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("access_token", accessToken);
        body.put("refresh_token", refreshToken);
        Map<String, Object> resp = http.doJson("POST", "/connect-auth/v1/account/refresh", null, body);
        return credentialsFromResponse(resp, "ClientCredentialsProvider(refresh)");
    }

    private static Credentials credentialsFromResponse(Map<String, Object> resp, String source) {
        Instant now = Instant.now();
        String accessToken = str(resp.get("access_token"));
        String refreshToken = str(resp.get("refresh_token"));
        long expiresIn = num(resp.get("expires_in"));
        long refreshExpiresIn = num(resp.get("refresh_expires_in"));
        return new Credentials(accessToken, refreshToken, now.plusSeconds(expiresIn), now.plusSeconds(refreshExpiresIn), source);
    }

    private static String str(Object v) {
        return v == null ? "" : v.toString();
    }

    private static long num(Object v) {
        if (v == null) {
            return 0;
        }
        if (v instanceof Number n) {
            return n.longValue();
        }
        return Long.parseLong(v.toString());
    }
}

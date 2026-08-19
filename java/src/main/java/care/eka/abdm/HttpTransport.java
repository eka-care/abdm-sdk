package care.eka.abdm;

import care.eka.abdm.internal.Json;
import com.fasterxml.jackson.core.type.TypeReference;
import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.Map;

/**
 * Minimal HTTP transport, stdlib-only ({@code java.net.http.HttpClient}),
 * mirroring {@code go/internal/http/http.go} and Python's {@code http.py}.
 *
 * <p>Two behaviours matter here more than anywhere else in the port:
 *
 * <ol>
 *   <li>The bearer token is resolved PER REQUEST ({@link Config#token()} inside
 *       {@link #doRaw}), never cached at construction. That is what lets a
 *       long-running process pick up a refreshed token.
 *   <li>{@code resolveToken = false} (used only by {@link AuthService}, the
 *       login/refresh client) skips that resolution entirely and sends a
 *       fixed empty bearer. This breaks the reentrant self-deadlock/recursion
 *       Go and Java each hit in their own way: request -&gt; resolve token -&gt;
 *       credentials provider takes its lock -&gt; refresh call -&gt; request -&gt;
 *       resolve token -&gt; same lock, same thread. Go's {@code sync.Mutex} is
 *       not reentrant, so that shape hangs forever; Java's {@code synchronized}
 *       IS reentrant, so the same shape does not hang — it recurses, calling
 *       the refresh endpoint repeatedly until the stack overflows. Login and
 *       refresh authenticate from the request body, not a bearer header, so
 *       they never need a token and must never call back into the token
 *       resolver. See {@link ClientCredentialsProvider} and
 *       {@code DeadlockRegressionTest}.
 * </ol>
 */
public final class HttpTransport {

    private final Config config;
    private final boolean resolveToken;
    private final HttpClient client;

    public HttpTransport(Config config, boolean resolveToken) {
        this.config = config;
        this.resolveToken = resolveToken;
        this.client = HttpClient.newBuilder()
                .connectTimeout(config.timeout())
                .build();
    }

    /** Performs a JSON-in/bytes-out request. Throws {@link AbdmApiException} for any 4xx/5xx. */
    public byte[] doRaw(String method, String path, Headers headers, Object body) throws IOException, InterruptedException {
        URI uri = URI.create(config.baseUrl() + path);
        HttpRequest.BodyPublisher publisher = HttpRequest.BodyPublishers.noBody();
        if (body != null) {
            byte[] json = Json.MAPPER.writeValueAsBytes(body);
            publisher = HttpRequest.BodyPublishers.ofByteArray(json);
        }

        HttpRequest.Builder builder = HttpRequest.newBuilder(uri)
                .timeout(config.timeout())
                .method(method, publisher);

        String token = resolveToken ? config.token() : "";
        builder.header("Authorization", "Bearer " + token);
        builder.header("Content-Type", "application/json");
        builder.header("User-Agent", config.userAgent());

        if (headers != null) {
            if (!headers.patientId().isEmpty()) {
                builder.header("X-Pt-Id", headers.patientId());
            }
            if (!headers.partnerUserId().isEmpty()) {
                builder.header("X-Partner-Pt-Id", headers.partnerUserId());
            }
            if (!headers.hipId().isEmpty()) {
                builder.header("X-Hip-Id", headers.hipId());
            }
        }

        HttpResponse<byte[]> response = client.send(builder.build(), HttpResponse.BodyHandlers.ofByteArray());
        if (response.statusCode() >= 400) {
            throw AbdmApiException.fromResponse(response.statusCode(), response.body());
        }
        return response.body();
    }

    /** Convenience wrapper that decodes a JSON object response into a {@code Map}. */
    public Map<String, Object> doJson(String method, String path, Headers headers, Object body) throws IOException, InterruptedException {
        byte[] raw = doRaw(method, path, headers, body);
        if (raw == null || raw.length == 0) {
            return Map.of();
        }
        return Json.MAPPER.readValue(raw, new TypeReference<Map<String, Object>>() {
        });
    }
}

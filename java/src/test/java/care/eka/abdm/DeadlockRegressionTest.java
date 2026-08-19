package care.eka.abdm;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTimeoutPreemptively;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.concurrent.atomic.AtomicInteger;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

/**
 * Regression test for gotcha #8 in the M2 language-port design doc: never
 * hold a lock across a token-refresh network call.
 *
 * <p>Go hit a reentrant self-deadlock in this exact shape: request -&gt;
 * resolve token -&gt; {@code ClientCredentialsProvider} takes its lock -&gt;
 * refresh call -&gt; request -&gt; resolve token -&gt; same lock, same thread
 * -&gt; permanent hang. Go's fix was making the auth/refresh client
 * ({@link AuthService}) NOT resolve tokens — it authenticates from the
 * request body and needs no bearer (see {@link HttpTransport}'s
 * {@code resolveToken} flag).
 *
 * <p>Java's {@code synchronized} (used by {@link ClientCredentialsProvider#retrieve})
 * is reentrant, unlike Go's {@code sync.Mutex}, so the same bug shape would
 * NOT deadlock here — it would silently recurse instead: the reentrant call
 * would see the same not-yet-updated, still-expired cache and issue another
 * refresh call, and so on, either overflowing the stack or hammering the
 * refresh endpoint an unbounded number of times. That is easier to miss than
 * a hang, not safer, so this test asserts both:
 *
 * <ol>
 *   <li>the request completes within 5 seconds (catches a hang or a
 *       stack-overflow-then-retry loop), and
 *   <li>the refresh endpoint was called exactly ONCE (catches the recursion
 *       variant specifically — a timeout alone would not, since bounded
 *       recursion that calls refresh 40 times could still finish in time).
 * </ol>
 *
 * <p>This test is only meaningful against the wiring under test:
 * {@link AuthService} building its internal {@link HttpTransport} with
 * {@code resolveToken = false}. That was verified manually during
 * development by temporarily changing {@code AuthService}'s constructor to
 * pass {@code true} instead: the request hung past the 5-second timeout
 * (JUnit reported a TestAbortedException from assertTimeoutPreemptively)
 * before it could even reach the call-count assertion, then the change was
 * reverted. There is no way to reconstruct that broken wiring through
 * {@code AuthService}'s public API — its constructor always builds the safe,
 * non-resolving transport — which is the point: the fix removes the unsafe
 * shape entirely rather than making it merely optional.
 */
class DeadlockRegressionTest {

    private HttpServer server;

    @AfterEach
    void stop() {
        if (server != null) {
            server.stop(0);
        }
    }

    @Test
    void refreshDuringATokenResolveCompletesAndCallsRefreshExactlyOnce() throws Exception {
        AtomicInteger loginCalls = new AtomicInteger();
        AtomicInteger refreshCalls = new AtomicInteger();
        AtomicInteger protectedCalls = new AtomicInteger();

        server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext("/connect-auth/v1/account/login", exchange ->
                respond(exchange, 200, "{\"access_token\":\"initial-token\",\"refresh_token\":\"valid-refresh\","
                        + "\"expires_in\":0,\"refresh_expires_in\":7200}", loginCalls));
        server.createContext("/connect-auth/v1/account/refresh", exchange ->
                respond(exchange, 200, "{\"access_token\":\"refreshed-token\",\"refresh_token\":\"refreshed-refresh\","
                        + "\"expires_in\":3600,\"refresh_expires_in\":7200}", refreshCalls));
        server.createContext("/abdm/v1/care-contexts/link", exchange ->
                respond(exchange, 202, "", protectedCalls));
        server.setExecutor(null);
        server.start();

        String baseUrl = "http://127.0.0.1:" + server.getAddress().getPort();
        Config config = new Config(Config.PRODUCTION, baseUrl, Duration.ofSeconds(5), "test-agent");
        AuthService auth = new AuthService(config); // production wiring: resolveToken = false
        ClientCredentialsProvider provider = new ClientCredentialsProvider(auth, "client-id", "client-secret");

        // Set up the starting state described in the design doc: an expired
        // access token plus a valid refresh token. This first retrieve()
        // call (expires_in=0) is plain setup, not the assertion target.
        Credentials initial = provider.retrieve();
        assertEquals("initial-token", initial.accessToken());
        assertEquals(1, loginCalls.get());
        assertEquals(0, refreshCalls.get());

        // Wire the token resolver exactly as EkaClient.setCredentialsProvider
        // does, then make an actual request through a token-resolving
        // transport. Resolving the token mid-request is what triggers the
        // refresh call inside provider.retrieve() while its lock is held.
        config.setTokenResolver(() -> provider.retrieve().accessToken());
        HttpTransport requestTransport = new HttpTransport(config, true);

        assertTimeoutPreemptively(Duration.ofSeconds(5), () ->
                requestTransport.doRaw("POST", "/abdm/v1/care-contexts/link",
                        new Headers("pt-1", "partner-1", "hip-1"), null));

        assertEquals(1, refreshCalls.get(), "refresh must be called exactly once, not recursively");
        assertEquals(1, protectedCalls.get());
        assertEquals(1, loginCalls.get(), "a valid refresh token must never trigger a re-login");
    }

    private static void respond(HttpExchange exchange, int status, String body, AtomicInteger counter) throws IOException {
        counter.incrementAndGet();
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(status, bytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }
}

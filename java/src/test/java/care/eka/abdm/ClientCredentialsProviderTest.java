package care.eka.abdm;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertSame;

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
 * Exercises the retrieve() ladder against a hermetic stub server: cached ->
 * refresh -> re-login, exactly Go's ClientCredentialsProvider.Retrieve.
 */
class ClientCredentialsProviderTest {

    private HttpServer server;

    @AfterEach
    void stop() {
        if (server != null) {
            server.stop(0);
        }
    }

    private HttpServer startServer(AtomicInteger loginCalls, AtomicInteger refreshCalls, long refreshExpiresIn) throws IOException {
        HttpServer s = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        s.createContext("/connect-auth/v1/account/login", exchange -> {
            loginCalls.incrementAndGet();
            String json = "{\"access_token\":\"login-token\",\"refresh_token\":\"login-refresh\","
                    + "\"expires_in\":3600,\"refresh_expires_in\":" + refreshExpiresIn + "}";
            respond(exchange, 200, json);
        });
        s.createContext("/connect-auth/v1/account/refresh", exchange -> {
            refreshCalls.incrementAndGet();
            String json = "{\"access_token\":\"refreshed-token\",\"refresh_token\":\"refreshed-refresh\","
                    + "\"expires_in\":3600,\"refresh_expires_in\":" + refreshExpiresIn + "}";
            respond(exchange, 200, json);
        });
        s.setExecutor(null);
        s.start();
        return s;
    }

    private static void respond(com.sun.net.httpserver.HttpExchange exchange, int status, String body) throws IOException {
        byte[] bytes = body.getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json");
        exchange.sendResponseHeaders(status, bytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }

    @Test
    void firstRetrieveLogsIn() throws Exception {
        AtomicInteger loginCalls = new AtomicInteger();
        AtomicInteger refreshCalls = new AtomicInteger();
        server = startServer(loginCalls, refreshCalls, 7200);

        Config config = new Config(Config.PRODUCTION, "http://127.0.0.1:" + server.getAddress().getPort(),
                Duration.ofSeconds(5), "test-agent");
        AuthService auth = new AuthService(config);
        ClientCredentialsProvider provider = new ClientCredentialsProvider(auth, "id", "secret");

        Credentials creds = provider.retrieve();
        assertEquals("login-token", creds.accessToken());
        assertEquals(1, loginCalls.get());
        assertEquals(0, refreshCalls.get());
    }

    @Test
    void secondRetrieveReturnsCachedTokenWithoutNetworkCall() throws Exception {
        AtomicInteger loginCalls = new AtomicInteger();
        AtomicInteger refreshCalls = new AtomicInteger();
        server = startServer(loginCalls, refreshCalls, 7200);

        Config config = new Config(Config.PRODUCTION, "http://127.0.0.1:" + server.getAddress().getPort(),
                Duration.ofSeconds(5), "test-agent");
        AuthService auth = new AuthService(config);
        ClientCredentialsProvider provider = new ClientCredentialsProvider(auth, "id", "secret");

        Credentials first = provider.retrieve();
        Credentials second = provider.retrieve();
        assertSame(first, second);
        assertEquals(1, loginCalls.get());
        assertEquals(0, refreshCalls.get());
    }

    @Test
    void expiredCacheWithUsableRefreshTokenRefreshesInsteadOfReLogin() throws Exception {
        AtomicInteger loginCalls = new AtomicInteger();
        AtomicInteger refreshCalls = new AtomicInteger();
        server = startServer(loginCalls, refreshCalls, 7200);

        Config config = new Config(Config.PRODUCTION, "http://127.0.0.1:" + server.getAddress().getPort(),
                Duration.ofSeconds(5), "test-agent");
        AuthService auth = new AuthService(config);
        ClientCredentialsProvider provider = new ClientCredentialsProvider(auth, "id", "secret");

        // Seed an already-expired cache with a still-valid refresh token by
        // logging in once (expires_in=3600 is not expired), then forcing
        // expiry via reflection-free approach: log in with a server that
        // reports expires_in=0 so it is immediately expired but refreshable.
        server.stop(0);
        server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext("/connect-auth/v1/account/login", exchange -> {
            loginCalls.incrementAndGet();
            respond(exchange, 200, "{\"access_token\":\"login-token\",\"refresh_token\":\"login-refresh\","
                    + "\"expires_in\":0,\"refresh_expires_in\":7200}");
        });
        server.createContext("/connect-auth/v1/account/refresh", exchange -> {
            refreshCalls.incrementAndGet();
            respond(exchange, 200, "{\"access_token\":\"refreshed-token\",\"refresh_token\":\"refreshed-refresh\","
                    + "\"expires_in\":3600,\"refresh_expires_in\":7200}");
        });
        server.setExecutor(null);
        server.start();
        config = new Config(Config.PRODUCTION, "http://127.0.0.1:" + server.getAddress().getPort(),
                Duration.ofSeconds(5), "test-agent");
        auth = new AuthService(config);
        provider = new ClientCredentialsProvider(auth, "id", "secret");

        Credentials first = provider.retrieve();
        assertEquals("login-token", first.accessToken());
        assertEquals(1, loginCalls.get());

        Credentials second = provider.retrieve();
        assertEquals("refreshed-token", second.accessToken());
        assertEquals(1, loginCalls.get());
        assertEquals(1, refreshCalls.get());
    }
}

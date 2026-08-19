package care.eka.abdm.carecontext;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import care.eka.abdm.AbdmApiException;
import care.eka.abdm.Config;
import care.eka.abdm.Headers;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

class CareContextsLinkTest {

    private static final ObjectMapper MAPPER = new ObjectMapper();

    private HttpServer server;
    private volatile String capturedMethod;
    private volatile String capturedPath;
    private volatile Map<String, String> capturedHeaders;
    private volatile String capturedBody;
    private volatile int responseStatus = 202;

    @AfterEach
    void stop() {
        if (server != null) {
            server.stop(0);
        }
    }

    private Config startServer() throws IOException {
        server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext("/abdm/v1/care-contexts/link", this::capture);
        server.setExecutor(null);
        server.start();
        return new Config(Config.PRODUCTION, "http://127.0.0.1:" + server.getAddress().getPort(),
                Duration.ofSeconds(5), "test-agent/1.0");
    }

    private void capture(HttpExchange exchange) throws IOException {
        capturedMethod = exchange.getRequestMethod();
        capturedPath = exchange.getRequestURI().getPath();
        capturedHeaders = Map.of(
                "Authorization", firstOrEmpty(exchange, "Authorization"),
                "X-Pt-Id", firstOrEmpty(exchange, "X-Pt-Id"),
                "X-Partner-Pt-Id", firstOrEmpty(exchange, "X-Partner-Pt-Id"),
                "X-Hip-Id", firstOrEmpty(exchange, "X-Hip-Id"),
                "Content-Type", firstOrEmpty(exchange, "Content-Type"));
        capturedBody = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);

        byte[] respBytes = responseStatus >= 400
                ? "{\"code\":400,\"error\":\"bad request\"}".getBytes(StandardCharsets.UTF_8)
                : new byte[0];
        exchange.sendResponseHeaders(responseStatus, respBytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(respBytes);
        }
    }

    private static String firstOrEmpty(HttpExchange exchange, String header) {
        List<String> values = exchange.getRequestHeaders().get(header);
        return values == null || values.isEmpty() ? "" : values.get(0);
    }

    @Test
    void linkSendsThreeDistinctHeadersAndFillsBodyDefaultsWithoutMutatingCaller() throws Exception {
        Config config = startServer();
        config.setAuthorizationToken("bearer-token");
        CareContexts careContexts = new CareContexts(config);

        Headers headers = new Headers("patient-1", "partner-2", "hip-3");
        LinkRequest original = new LinkRequest("test@sbx",
                List.of(new CareContext("cc-1", "Visit on 2026-08-19")));

        careContexts.link(headers, original);

        // Three distinct header values, asserted individually so a
        // transposition cannot pass by coincidence.
        assertEquals("patient-1", capturedHeaders.get("X-Pt-Id"));
        assertEquals("partner-2", capturedHeaders.get("X-Partner-Pt-Id"));
        assertEquals("hip-3", capturedHeaders.get("X-Hip-Id"));
        assertEquals("Bearer bearer-token", capturedHeaders.get("Authorization"));
        assertEquals("application/json", capturedHeaders.get("Content-Type"));

        assertEquals("POST", capturedMethod);
        assertEquals("/abdm/v1/care-contexts/link", capturedPath);

        Map<?, ?> body = MAPPER.readValue(capturedBody, Map.class);
        assertEquals("test@sbx", body.get("abha_address"));
        assertEquals("patient-1", body.get("oid"), "oid must be filled from X-Pt-Id when left empty");
        assertEquals("partner-2", body.get("partner_user_id"), "partner_user_id must be filled from X-Partner-Pt-Id when left empty");

        // The caller's original request object must not be mutated.
        assertEquals("", original.oid());
        assertEquals("", original.partnerUserId());
    }

    @Test
    void linkDoesNotOverrideExplicitOidOrPartnerUserId() throws Exception {
        Config config = startServer();
        CareContexts careContexts = new CareContexts(config);

        Headers headers = new Headers("patient-1", "partner-2", "hip-3");
        LinkRequest original = new LinkRequest("test@sbx",
                List.of(new CareContext("cc-1", "Visit")), "explicit-oid", "explicit-partner");

        careContexts.link(headers, original);

        Map<?, ?> body = MAPPER.readValue(capturedBody, Map.class);
        assertEquals("explicit-oid", body.get("oid"));
        assertEquals("explicit-partner", body.get("partner_user_id"));
    }

    @Test
    void errorResponseThrowsAbdmApiExceptionWithStatusAndBody() throws Exception {
        Config config = startServer();
        responseStatus = 400;
        CareContexts careContexts = new CareContexts(config);

        Headers headers = new Headers("patient-1", "partner-2", "hip-3");
        LinkRequest request = new LinkRequest("test@sbx", List.of(new CareContext("cc-1", "Visit")));

        AbdmApiException ex = assertThrows(AbdmApiException.class, () -> careContexts.link(headers, request));
        assertEquals(400, ex.statusCode());
        assertTrue(ex.errorBody().contains("bad request"));
    }
}

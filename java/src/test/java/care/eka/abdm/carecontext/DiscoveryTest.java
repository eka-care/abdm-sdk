package care.eka.abdm.carecontext;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

import care.eka.abdm.Config;
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
import java.util.Set;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

class DiscoveryTest {

    private static final ObjectMapper MAPPER = new ObjectMapper();

    private HttpServer server;
    private volatile String capturedPath;
    private volatile Map<String, String> capturedHeaders;
    private volatile String capturedBody;

    @AfterEach
    void stop() {
        if (server != null) {
            server.stop(0);
        }
    }

    private Config startServer(String path) throws IOException {
        server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext(path, this::capture);
        server.setExecutor(null);
        server.start();
        return new Config(Config.PRODUCTION, "http://127.0.0.1:" + server.getAddress().getPort(),
                Duration.ofSeconds(5), "test-agent/1.0");
    }

    private void capture(HttpExchange exchange) throws IOException {
        capturedPath = exchange.getRequestURI().getPath();
        capturedHeaders = Map.of(
                "X-Pt-Id", firstOrEmpty(exchange, "X-Pt-Id"),
                "X-Partner-Pt-Id", firstOrEmpty(exchange, "X-Partner-Pt-Id"),
                "X-Hip-Id", firstOrEmpty(exchange, "X-Hip-Id"));
        capturedBody = new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
        byte[] resp = new byte[0];
        exchange.sendResponseHeaders(204, resp.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(resp);
        }
    }

    private static String firstOrEmpty(HttpExchange exchange, String header) {
        List<String> values = exchange.getRequestHeaders().get(header);
        return values == null || values.isEmpty() ? "" : values.get(0);
    }

    private void assertThreeDistinctHeaders(String pt, String partner, String hip) {
        assertEquals(pt, capturedHeaders.get("X-Pt-Id"));
        assertEquals(partner, capturedHeaders.get("X-Partner-Pt-Id"));
        assertEquals(hip, capturedHeaders.get("X-Hip-Id"));
        assertEquals(3, Set.of(pt, partner, hip).size());
    }

    // --- onDiscover --------------------------------------------------------------

    @Test
    void onDiscoverSendsMatchedPatients() throws Exception {
        Config config = startServer("/abdm/v1/care-contexts/on-discover");
        CareContexts careContexts = new CareContexts(config);

        DiscoverEvent event = new DiscoverEvent("", "", "", 0, List.of(), "req-1", "txn-1", "oid-1", "pp-1", "hip-1");
        Patient patient = new Patient("P1", "Gajendra", "OPConsultation",
                List.of(new DiscoveredCareContext("CC101", "OP Consult")));
        careContexts.onDiscover(event, new DiscoverResult(List.of(patient)));

        assertEquals("/abdm/v1/care-contexts/on-discover", capturedPath);
        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        assertEquals("req-1", body.get("request_id"));
        assertEquals("txn-1", body.get("txn_id"));
        @SuppressWarnings("unchecked")
        List<Map<String, Object>> patients = (List<Map<String, Object>>) body.get("patients");
        assertEquals("P1", patients.get(0).get("ref_num"));
        assertEquals("OPConsultation", patients.get(0).get("hi_type"));

        assertThreeDistinctHeaders("oid-1", "pp-1", "hip-1");
    }

    @Test
    void onDiscoverErrorResultOmitsPatients() throws Exception {
        Config config = startServer("/abdm/v1/care-contexts/on-discover");
        CareContexts careContexts = new CareContexts(config);

        DiscoverEvent event = new DiscoverEvent("", "", "", 0, List.of(), "r", "t", "oid-err", "pp-err", "hip-err");
        careContexts.onDiscover(event, DiscoverResult.ofError(new ErrorDetail(1000, "no match")));

        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        @SuppressWarnings("unchecked")
        Map<String, Object> error = (Map<String, Object>) body.get("error");
        assertEquals("no match", error.get("message"));
        assertFalse(body.containsKey("patients"));

        assertThreeDistinctHeaders("oid-err", "pp-err", "hip-err");
    }

    // --- onLinkInit ----------------------------------------------------------------

    @Test
    void onLinkInitSendsRefNumAndOtpExpiry() throws Exception {
        Config config = startServer("/abdm/v1/care-contexts/discover/link/on-init");
        CareContexts careContexts = new CareContexts(config);

        LinkInitEvent event = new LinkInitEvent("", null, "req-2", "txn-2", "oid-2", "pp-2", "hip-2");
        careContexts.onLinkInit(event, new LinkInitResult("temp", "2026-08-19T10:00:00Z"));

        assertEquals("/abdm/v1/care-contexts/discover/link/on-init", capturedPath);
        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        assertEquals("temp", body.get("ref_num"));
        assertEquals("2026-08-19T10:00:00Z", body.get("otp_expiry"));

        assertThreeDistinctHeaders("oid-2", "pp-2", "hip-2");
    }

    @Test
    void onLinkInitErrorResultOmitsRefNumAndOtpExpiry() throws Exception {
        Config config = startServer("/abdm/v1/care-contexts/discover/link/on-init");
        CareContexts careContexts = new CareContexts(config);

        LinkInitEvent event = new LinkInitEvent("", null, "req-1", "txn-1", "", "", "");
        careContexts.onLinkInit(event, LinkInitResult.ofError(new ErrorDetail(1000, "no match")));

        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        assertFalse(body.containsKey("ref_num"));
        assertFalse(body.containsKey("otp_expiry"));
        @SuppressWarnings("unchecked")
        Map<String, Object> error = (Map<String, Object>) body.get("error");
        assertEquals("no match", error.get("message"));
        assertEquals("req-1", body.get("request_id"));
    }

    // --- onLinkConfirm --------------------------------------------------------------

    @Test
    void onLinkConfirmCorrelatesByRequestIdOnlyNoTxnIdSent() throws Exception {
        Config config = startServer("/abdm/v1/care-contexts/discover/link/on-confirm");
        CareContexts careContexts = new CareContexts(config);

        LinkConfirmEvent event = new LinkConfirmEvent("temp", "111111", "req-3", "oid-3", "pp-3", "hip-3");
        careContexts.onLinkConfirm(event, new LinkConfirmResult(List.of(new Patient("P1", ""))));

        assertEquals("/abdm/v1/care-contexts/discover/link/on-confirm", capturedPath);
        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        assertEquals("req-3", body.get("request_id"));
        assertFalse(body.containsKey("txn_id"));

        assertThreeDistinctHeaders("oid-3", "pp-3", "hip-3");
    }

    @Test
    void onLinkConfirmErrorResultOmitsPatients() throws Exception {
        Config config = startServer("/abdm/v1/care-contexts/discover/link/on-confirm");
        CareContexts careContexts = new CareContexts(config);

        LinkConfirmEvent event = new LinkConfirmEvent("temp", "000000", "req-4", "", "", "");
        careContexts.onLinkConfirm(event, LinkConfirmResult.ofError(new ErrorDetail(1001, "bad otp")));

        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        @SuppressWarnings("unchecked")
        Map<String, Object> error = (Map<String, Object>) body.get("error");
        assertEquals("bad otp", error.get("message"));
        assertFalse(body.containsKey("patients"));
    }
}

package care.eka.abdm.carecontext;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

import care.eka.abdm.Config;
import care.eka.abdmecdh.AbdmEcdh;
import care.eka.abdmecdh.DecryptionResponse;
import care.eka.abdmecdh.KeyMaterial;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;
import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.time.Duration;
import java.util.List;
import java.util.Map;
import java.util.Set;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.Test;

class RespondToFetchTest {

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

    private Config startServer() throws IOException {
        server = HttpServer.create(new InetSocketAddress("127.0.0.1", 0), 0);
        server.createContext("/abdm/v1/hip/care-context/data/on-fetch", this::capture);
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
        exchange.sendResponseHeaders(202, resp.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(resp);
        }
    }

    private static String firstOrEmpty(HttpExchange exchange, String header) {
        List<String> values = exchange.getRequestHeaders().get(header);
        return values == null || values.isEmpty() ? "" : values.get(0);
    }

    private static DataFetchEvent makeEvent(KeyMaterial hiu) {
        Map<String, Object> dhPublicKey = Map.of("key_value", hiu.x509PublicKey(), "parameters", "p");
        Map<String, Object> keyInfo = Map.of(
                "crypto_alg", "ECDH", "curve", "Curve25519", "nonce", hiu.nonce(), "dh_public_key", dhPublicKey);
        return new DataFetchEvent("txn-1", "", "oid-1", "pp-1", "hip-1", List.of(), List.of(), keyInfo);
    }

    private static String hexMd5(byte[] plaintext) throws Exception {
        MessageDigest md5 = MessageDigest.getInstance("MD5");
        byte[] sum = md5.digest(plaintext);
        StringBuilder sb = new StringBuilder();
        for (byte b : sum) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    @Test
    void checksumMatchesTheConfirmedTestVector() throws Exception {
        byte[] plaintext = "{\"resourceType\":\"Bundle\",\"id\":\"b1\"}".getBytes(StandardCharsets.UTF_8);
        assertEquals("d8d404a602afd8230d69e45439d91b7f", hexMd5(plaintext));
    }

    @Test
    void hiuCanDecryptUsingOnlyThePublishedKeyMaterialAndThreeHeadersAreDistinct() throws Exception {
        // A second key pair, acting as the requesting HIU, entirely separate
        // from whatever the responder generates internally.
        KeyMaterial hiu = AbdmEcdh.generateKeyMaterial();

        Config config = startServer();
        CareContexts careContexts = new CareContexts(config);

        byte[] bundle = "{\"resourceType\":\"Bundle\",\"id\":\"b1\"}".getBytes(StandardCharsets.UTF_8);
        careContexts.respondToFetch(makeEvent(hiu), List.of(new Entry("cc-1", bundle)));

        assertEquals("/abdm/v1/hip/care-context/data/on-fetch", capturedPath);
        @SuppressWarnings("unchecked")
        Map<String, Object> body = MAPPER.readValue(capturedBody, Map.class);
        assertEquals("txn-1", body.get("transaction_id"));

        @SuppressWarnings("unchecked")
        List<Map<String, Object>> entries = (List<Map<String, Object>>) body.get("entries");
        Map<String, Object> entry = entries.get(0);
        assertEquals("cc-1", entry.get("care_context_id"));
        assertEquals("application/fhir+json", entry.get("media"));
        assertEquals(hexMd5(bundle), entry.get("checksum"));

        // Decrypt as the HIU would, using ONLY the key material the
        // responder published in its own request body.
        @SuppressWarnings("unchecked")
        Map<String, Object> keyInformation = (Map<String, Object>) body.get("key_information");
        @SuppressWarnings("unchecked")
        Map<String, Object> dhPublicKey = (Map<String, Object>) keyInformation.get("dh_public_key");
        DecryptionResponse dec = AbdmEcdh.decrypt(
                (String) entry.get("content"),
                (String) keyInformation.get("nonce"),
                hiu.nonce(),
                hiu.privateKey(),
                (String) dhPublicKey.get("key_value"));
        assertEquals(new String(bundle, StandardCharsets.UTF_8), dec.decryptedData());

        // Three distinct ABDM headers, so a transposition cannot pass by coincidence.
        assertEquals("oid-1", capturedHeaders.get("X-Pt-Id"));
        assertEquals("pp-1", capturedHeaders.get("X-Partner-Pt-Id"));
        assertEquals("hip-1", capturedHeaders.get("X-Hip-Id"));
        assertEquals(3, Set.of(capturedHeaders.get("X-Pt-Id"), capturedHeaders.get("X-Partner-Pt-Id"),
                capturedHeaders.get("X-Hip-Id")).size());
    }

    @Test
    void rejectsEmptyEntries() {
        KeyMaterial hiu = AbdmEcdh.generateKeyMaterial();
        Config config = new Config(Config.PRODUCTION, "http://unused", Duration.ofSeconds(1), "test-agent");
        CareContexts careContexts = new CareContexts(config);
        assertThrows(IllegalArgumentException.class,
                () -> careContexts.respondToFetch(makeEvent(hiu), List.of()));
    }

    @Test
    void rejectsNullEvent() {
        Config config = new Config(Config.PRODUCTION, "http://unused", Duration.ofSeconds(1), "test-agent");
        CareContexts careContexts = new CareContexts(config);
        assertThrows(IllegalArgumentException.class,
                () -> careContexts.respondToFetch(null, List.of(new Entry("cc-1", "{}".getBytes(StandardCharsets.UTF_8)))));
    }

    @Test
    void rejectsEmptyBundle() {
        KeyMaterial hiu = AbdmEcdh.generateKeyMaterial();
        Config config = new Config(Config.PRODUCTION, "http://unused", Duration.ofSeconds(1), "test-agent");
        CareContexts careContexts = new CareContexts(config);
        assertThrows(IllegalArgumentException.class, () -> careContexts.respondToFetch(
                makeEvent(hiu), List.of(new Entry("cc-1", new byte[0]))));
    }
}

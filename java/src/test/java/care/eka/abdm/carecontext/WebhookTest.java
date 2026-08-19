package care.eka.abdm.carecontext;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import org.junit.jupiter.api.Test;

class WebhookTest {

    private static final String SECRET = "whsec";

    private static String sign(String body, String secret, long ts) {
        try {
            Mac mac = Mac.getInstance("HmacSHA256");
            mac.init(new SecretKeySpec(secret.getBytes(StandardCharsets.UTF_8), "HmacSHA256"));
            mac.update(Long.toString(ts).getBytes(StandardCharsets.UTF_8));
            mac.update((byte) '.');
            mac.update(body.getBytes(StandardCharsets.UTF_8));
            byte[] sum = mac.doFinal();
            StringBuilder sb = new StringBuilder();
            for (byte b : sum) {
                sb.append(String.format("%02x", b));
            }
            return "t=" + ts + ",v1=" + sb;
        } catch (NoSuchAlgorithmException | InvalidKeyException e) {
            throw new RuntimeException(e);
        }
    }

    private static long now() {
        return System.currentTimeMillis() / 1000;
    }

    // --- Signature verification -------------------------------------------------

    @Test
    void validSignatureDecodesLinkStatusEvent() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\","
                + "\"data\":{\"care_context_id\":\"cc-1\",\"status\":\"LINKED\"}}";
        Event e = Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sign(body, SECRET, now()), SECRET);
        LinkStatusEvent event = assertInstanceOf(LinkStatusEvent.class, e);
        assertEquals("cc-1", event.careContextId());
        assertEquals("LINKED", event.status());
    }

    @Test
    void tamperedBodyRejected() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\",\"data\":{}}";
        String sig = sign(body, SECRET, now());
        byte[] tampered = (body + " ").getBytes(StandardCharsets.UTF_8);
        assertThrows(BadSignatureException.class, () -> Webhook.parseWebhook(tampered, sig, SECRET));
    }

    @Test
    void wrongSecretRejected() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\",\"data\":{}}";
        String sig = sign(body, "other-secret", now());
        assertThrows(BadSignatureException.class,
                () -> Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sig, SECRET));
    }

    @Test
    void emptySecretRejected() {
        // javax.crypto refuses to even construct a zero-length HMAC key, so
        // there is no "genuinely valid signature for the empty key" to
        // build here (unlike Go/Python, whose HMAC implementations accept
        // one) — the empty-secret guard must reject before any HMAC is
        // computed, which this exercises with an arbitrary v1 value.
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\",\"data\":{}}";
        String header = "t=" + now() + ",v1=" + "0".repeat(64);
        assertThrows(BadSignatureException.class,
                () -> Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), header, ""));
    }

    @Test
    void staleTimestampInThePastRejected() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\",\"data\":{}}";
        long old = now() - 10 * 60;
        String sig = sign(body, SECRET, old);
        assertThrows(StaleTimestampException.class,
                () -> Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sig, SECRET));
    }

    @Test
    void staleTimestampInTheFutureAlsoRejected() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\",\"data\":{}}";
        long future = now() + 10 * 60;
        String sig = sign(body, SECRET, future);
        assertThrows(StaleTimestampException.class,
                () -> Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sig, SECRET));
    }

    @Test
    void malformedSignatureHeaderRejected() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.link_care_context\",\"data\":{}}";
        assertThrows(BadSignatureException.class,
                () -> Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), "garbage", SECRET));
    }

    // --- Event names --------------------------------------------------------------

    @Test
    void eventNameConstantsAreExactIncludingTheMissingCarePrefix() {
        assertEquals("abha.hip_data_fetch", Webhook.EVENT_DATA_FETCH);
        assertEquals("abha.link_care_context", Webhook.EVENT_LINK_STATUS);
        assertEquals("abha.care_context_discover", Webhook.EVENT_DISCOVER);
        assertEquals("abha.care_context_discover_link_init", Webhook.EVENT_LINK_INIT);
        // No care_ prefix on this one — that is not a typo.
        assertEquals("abha.context_discover_link_confirm", Webhook.EVENT_LINK_CONFIRM);
    }

    @Test
    void eachEventDecodesToTheRightType() {
        Object[][] cases = {
                {Webhook.EVENT_DATA_FETCH, DataFetchEvent.class},
                {Webhook.EVENT_LINK_STATUS, LinkStatusEvent.class},
                {Webhook.EVENT_DISCOVER, DiscoverEvent.class},
                {Webhook.EVENT_LINK_INIT, LinkInitEvent.class},
                {Webhook.EVENT_LINK_CONFIRM, LinkConfirmEvent.class},
        };
        for (Object[] c : cases) {
            String eventName = (String) c[0];
            Class<?> wantType = (Class<?>) c[1];
            String body = "{\"service\":\"abdm\",\"event\":\"" + eventName + "\",\"data\":{}}";
            Event e = Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sign(body, SECRET, now()), SECRET);
            assertInstanceOf(wantType, e, eventName);
            assertEquals(eventName, e.eventName(), eventName);
        }
    }

    @Test
    void unknownEventDistinguishableFromInauthentic() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.something_new\",\"data\":{}}";
        assertThrows(UnknownEventException.class,
                () -> Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sign(body, SECRET, now()), SECRET));
    }

    // --- Field-level details ---------------------------------------------------

    @Test
    void dataFetchCapturesKeyMaterialAndFallsBackToEnvelopeTransactionId() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.hip_data_fetch\","
                + "\"data\":{\"transaction_id\":\"txn-1\",\"abha_address\":\"p@sbx\",\"oid\":\"o-1\","
                + "\"partner_patient_id\":\"pp-1\",\"hip_id\":\"hip-1\",\"care_contexts\":[\"cc-1\"],"
                + "\"hi_types\":[\"OPConsultation\"],\"key_information\":{\"crypto_alg\":\"ECDH\","
                + "\"curve\":\"Curve25519\",\"dh_public_key\":{\"key_value\":\"k\",\"parameters\":\"p\","
                + "\"expiry\":\"e\"},\"nonce\":\"n\"}}}";
        Event e = Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sign(body, SECRET, now()), SECRET);
        DataFetchEvent event = assertInstanceOf(DataFetchEvent.class, e);
        assertEquals("txn-1", event.transactionId());
        assertEquals("hip-1", event.hipId());
        assertEquals("cc-1", event.careContexts().get(0));

        String fallbackBody = "{\"service\":\"abdm\",\"event\":\"abha.hip_data_fetch\",\"transaction_id\":\"txn-9\"}";
        Event fallback = Webhook.parseWebhook(
                fallbackBody.getBytes(StandardCharsets.UTF_8), sign(fallbackBody, SECRET, now()), SECRET);
        assertEquals("txn-9", ((DataFetchEvent) fallback).transactionId());
    }

    @Test
    void linkConfirmDecodesCamelCaseLinkRefNumber() {
        String body = "{\"service\":\"abdm\",\"event\":\"abha.context_discover_link_confirm\","
                + "\"data\":{\"linkRefNumber\":\"ref-123\",\"token\":\"111111\",\"request_id\":\"req-3\","
                + "\"oid\":\"o-1\",\"partner_patient_id\":\"pp-1\",\"hip_id\":\"hip-1\"}}";
        Event e = Webhook.parseWebhook(body.getBytes(StandardCharsets.UTF_8), sign(body, SECRET, now()), SECRET);
        LinkConfirmEvent event = assertInstanceOf(LinkConfirmEvent.class, e);
        assertEquals("ref-123", event.linkRefNumber());
        assertEquals("111111", event.token());
    }
}

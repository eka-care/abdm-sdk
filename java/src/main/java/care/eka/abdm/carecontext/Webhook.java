package care.eka.abdm.carecontext;

import care.eka.abdm.internal.Json;
import com.fasterxml.jackson.core.type.TypeReference;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

/**
 * Verifying and decoding Eka webhooks. Mirrors
 * {@code go/services/abdm/carecontext/webhook.go} and Python's
 * {@code webhook.py}.
 *
 * <p>{@code carecontext} is a library, not a server: you own your webhook
 * endpoint and register its URL with Eka. Inside your handler, call
 * {@link #parseWebhook} and then the responder for the event you received.
 *
 * <p>ABDM retries a webhook it considers unanswered, so callbacks must be
 * idempotent — deduplicate on the event's transaction/request id before
 * doing anything with side effects. A non-2xx response makes ABDM retry.
 */
public final class Webhook {

    private Webhook() {
    }

    // Webhook event names, verbatim from the ABDM Connect docs. These are
    // inconsistent upstream — note that the link-confirm event has no
    // care_ prefix.
    public static final String EVENT_DATA_FETCH = "abha.hip_data_fetch";
    public static final String EVENT_LINK_STATUS = "abha.link_care_context";
    public static final String EVENT_DISCOVER = "abha.care_context_discover";
    public static final String EVENT_LINK_INIT = "abha.care_context_discover_link_init";
    public static final String EVENT_LINK_CONFIRM = "abha.context_discover_link_confirm";

    /** How far a webhook timestamp may drift, in either direction, before rejection. */
    private static final Duration SIGNATURE_TOLERANCE = Duration.ofMinutes(3);

    /**
     * Verifies an Eka webhook and decodes it into a typed event.
     *
     * <p>{@code signatureHeader} is the raw {@code Eka-Webhook-Signature}
     * header; {@code secret} is the signing key from your webhook
     * subscription; {@code body} must be the exact bytes received, because
     * the signature covers them verbatim.
     *
     * <p>Throws {@link BadSignatureException}, {@link StaleTimestampException}
     * or {@link UnknownEventException}. Treat the first two as 401 and
     * {@link UnknownEventException} as 200.
     */
    public static Event parseWebhook(byte[] body, String signatureHeader, String secret) {
        verifySignature(body, signatureHeader, secret);

        Map<String, Object> env;
        try {
            env = Json.MAPPER.readValue(body, new TypeReference<Map<String, Object>>() {
            });
        } catch (Exception e) {
            throw new BadSignatureException("carecontext: decode webhook envelope: " + e.getMessage());
        }

        String event = str(env.get("event"));
        Map<String, Object> data = asMap(env.get("data"));

        switch (event) {
            case EVENT_DATA_FETCH -> {
                String txn = str(data.get("transaction_id"));
                if (txn.isEmpty()) {
                    // The envelope's transaction_id is the fallback when data omits it.
                    txn = str(env.get("transaction_id"));
                }
                return new DataFetchEvent(
                        txn,
                        str(data.get("abha_address")),
                        str(data.get("oid")),
                        str(data.get("partner_patient_id")),
                        str(data.get("hip_id")),
                        strList(data.get("care_contexts")),
                        strList(data.get("hi_types")),
                        asMap(data.get("key_information")));
            }
            case EVENT_LINK_STATUS -> {
                return new LinkStatusEvent(
                        str(data.get("abha_address")),
                        str(data.get("care_context_id")),
                        str(data.get("status")),
                        str(data.get("error")),
                        intOf(data.get("retry_count")),
                        str(data.get("oid")),
                        str(data.get("partner_patient_id")),
                        str(data.get("hip_id")));
            }
            case EVENT_DISCOVER -> {
                List<Identifier> identifiers = new ArrayList<>();
                Object rawIdentifiers = data.get("identifiers");
                if (rawIdentifiers instanceof List<?> list) {
                    for (Object o : list) {
                        Map<String, Object> m = asMap(o);
                        identifiers.add(new Identifier(str(m.get("type")), str(m.get("value"))));
                    }
                }
                return new DiscoverEvent(
                        str(data.get("abha_address")),
                        str(data.get("patient_name")),
                        str(data.get("gender")),
                        intOf(data.get("year_of_birth")),
                        identifiers,
                        str(data.get("request_id")),
                        str(data.get("txn_id")),
                        str(data.get("oid")),
                        str(data.get("partner_patient_id")),
                        str(data.get("hip_id")));
            }
            case EVENT_LINK_INIT -> {
                return new LinkInitEvent(
                        str(data.get("abha_address")),
                        data.get("patient"),
                        str(data.get("request_id")),
                        str(data.get("txn_id")),
                        str(data.get("oid")),
                        str(data.get("partner_patient_id")),
                        str(data.get("hip_id")));
            }
            case EVENT_LINK_CONFIRM -> {
                // linkRefNumber is camelCase on the wire, unlike its
                // snake_case neighbours — not a typo, copy verbatim.
                return new LinkConfirmEvent(
                        str(data.get("linkRefNumber")),
                        str(data.get("token")),
                        str(data.get("request_id")),
                        str(data.get("oid")),
                        str(data.get("partner_patient_id")),
                        str(data.get("hip_id")));
            }
            default -> throw new UnknownEventException("carecontext: unrecognised webhook event: \"" + event + "\"");
        }
    }

    private static void verifySignature(byte[] body, String header, String secret) {
        String ts = null;
        String v1 = null;
        for (String part : (header == null ? "" : header).split(",")) {
            String trimmed = part.trim();
            int eq = trimmed.indexOf('=');
            if (eq < 0) {
                continue;
            }
            String k = trimmed.substring(0, eq);
            String v = trimmed.substring(eq + 1);
            if ("t".equals(k)) {
                ts = v;
            } else if ("v1".equals(k)) {
                v1 = v;
            }
        }
        if (ts == null || ts.isEmpty() || v1 == null || v1.isEmpty()) {
            throw new BadSignatureException("carecontext: malformed webhook signature header");
        }

        // An empty secret would key the HMAC with nothing, which anyone can
        // compute. Refuse rather than authenticate every forgery.
        if (secret == null || secret.isEmpty()) {
            throw new BadSignatureException("carecontext: no webhook secret configured");
        }

        String expected;
        try {
            Mac mac = Mac.getInstance("HmacSHA256");
            mac.init(new SecretKeySpec(secret.getBytes(StandardCharsets.UTF_8), "HmacSHA256"));
            mac.update(ts.getBytes(StandardCharsets.UTF_8));
            mac.update((byte) '.');
            mac.update(body == null ? new byte[0] : body);
            expected = hex(mac.doFinal());
        } catch (NoSuchAlgorithmException | InvalidKeyException e) {
            throw new BadSignatureException("carecontext: signature verification failed: " + e.getMessage());
        }

        // Constant-time comparison: a data-dependent-time comparison here
        // would leak the correct signature one byte at a time.
        if (!MessageDigest.isEqual(
                expected.getBytes(StandardCharsets.UTF_8), v1.getBytes(StandardCharsets.UTF_8))) {
            throw new BadSignatureException("carecontext: signature mismatch");
        }

        // Timestamp is checked only after the signature, so an attacker
        // cannot use timing here to learn anything about the key.
        long sec;
        try {
            sec = Long.parseLong(ts);
        } catch (NumberFormatException e) {
            throw new BadSignatureException("carecontext: unparseable timestamp");
        }
        long driftSeconds = System.currentTimeMillis() / 1000 - sec;
        if (Math.abs(driftSeconds) > SIGNATURE_TOLERANCE.toSeconds()) {
            throw new StaleTimestampException("carecontext: webhook timestamp outside tolerance");
        }
    }

    private static String hex(byte[] bytes) {
        StringBuilder sb = new StringBuilder(bytes.length * 2);
        for (byte b : bytes) {
            sb.append(Character.forDigit((b >> 4) & 0xF, 16));
            sb.append(Character.forDigit(b & 0xF, 16));
        }
        return sb.toString();
    }

    private static String str(Object v) {
        return v == null ? "" : v.toString();
    }

    private static int intOf(Object v) {
        if (v instanceof Number n) {
            return n.intValue();
        }
        if (v == null) {
            return 0;
        }
        try {
            return Integer.parseInt(v.toString());
        } catch (NumberFormatException e) {
            return 0;
        }
    }

    @SuppressWarnings("unchecked")
    private static Map<String, Object> asMap(Object v) {
        if (v instanceof Map<?, ?> m) {
            return (Map<String, Object>) m;
        }
        return Map.of();
    }

    private static List<String> strList(Object v) {
        List<String> out = new ArrayList<>();
        if (v instanceof List<?> list) {
            for (Object o : list) {
                out.add(str(o));
            }
        }
        return out;
    }
}

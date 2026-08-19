package care.eka.abdm.carecontext;

import care.eka.abdm.Config;
import care.eka.abdm.Headers;
import care.eka.abdm.HttpTransport;
import care.eka.abdmecdh.AbdmEcdh;
import care.eka.abdmecdh.EncryptionResponse;
import care.eka.abdmecdh.KeyMaterial;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.time.Instant;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * Calls the ABDM care-context APIs. Mirrors
 * {@code go/services/abdm/carecontext/service.go}'s {@code Service}
 * (plus {@code datafetch.go} and {@code discovery.go} for the three
 * responders).
 *
 * <p>The package is boilerplate only. It carries bytes and correlation
 * identifiers; it never decides what your data means. Matching patients,
 * generating OTPs and building FHIR bundles stay in your code.
 *
 * <p>It is a library, not a server: you own your webhook endpoint and
 * register its URL with Eka. Inside your handler, call
 * {@link Webhook#parseWebhook} and then the responder for the event you
 * received.
 *
 * <p>One asymmetry to know about: a bundle whose checksum does not match
 * what the HIU computes is discarded at the far end, and nothing reports
 * that back. Eka answers 202 and {@link #respondToFetch} returns normally
 * either way, so a checksum mismatch is invisible from the HIP side —
 * records appear shared but never arrive.
 */
public final class CareContexts {

    private static final String FHIR_MEDIA = "application/fhir+json";

    private final HttpTransport http;

    public CareContexts(Config config) {
        this.http = new HttpTransport(config, true);
    }

    /**
     * Links care contexts to a patient's ABHA address.
     *
     * <p>The API is asynchronous: a 202 means accepted, not linked. The
     * outcome arrives later as an {@code abha.link_care_context} webhook,
     * which {@link Webhook#parseWebhook} decodes into a
     * {@link LinkStatusEvent}.
     *
     * <p>Does not mutate {@code request} — oid/partnerUserId are filled from
     * headers on a copy when left empty.
     */
    public void link(Headers headers, LinkRequest request) throws IOException, InterruptedException {
        LinkRequest body = request.withHeaderDefaults(headers);
        http.doRaw("POST", "/abdm/v1/care-contexts/link", headers, body.toMap());
    }

    /**
     * Encrypts the given FHIR bundles for the requesting HIU and pushes
     * them, answering an {@code abha.hip_data_fetch} webhook.
     *
     * <p>Pass {@code event} exactly as {@link Webhook#parseWebhook} returned
     * it: it carries the transaction id, the HIU's key material, and the
     * identifiers used as request headers. A fresh ephemeral key pair is
     * generated per call and discarded, so there is no key store to manage
     * and the private key never leaves this process.
     *
     * <p>ABDM retries a fetch it considers unanswered, so this may be called
     * more than once for one transaction. Make your own bundle lookup
     * idempotent.
     */
    public void respondToFetch(DataFetchEvent event, List<Entry> entries) throws IOException, InterruptedException {
        if (event == null) {
            throw new IllegalArgumentException("carecontext: nil data fetch event");
        }
        if (entries == null || entries.isEmpty()) {
            throw new IllegalArgumentException(
                    "carecontext: no entries to send for transaction \"" + (event.transactionId()) + "\"");
        }

        Map<String, Object> theirKeyInfo = event.keyInfo();
        // WE are the sender; the requesting HIU is the requester. The
        // encryption target is the HIU's X.509 public key, taken from
        // theirKeyInfo below — never our own.
        KeyMaterial ours = AbdmEcdh.generateKeyMaterial();

        Map<String, Object> theirDh = asMap(theirKeyInfo.get("dh_public_key"));
        String theirNonce = str(theirKeyInfo.get("nonce"));
        String theirPublicKey = str(theirDh.get("key_value"));

        List<Map<String, Object>> outEntries = new ArrayList<>(entries.size());
        for (Entry entry : entries) {
            if (entry.bundle() == null || entry.bundle().length == 0) {
                throw new IllegalArgumentException(
                        "carecontext: empty bundle for care context \"" + entry.careContextId() + "\"");
            }
            String plaintext = new String(entry.bundle(), StandardCharsets.UTF_8);
            EncryptionResponse enc = AbdmEcdh.encrypt(
                    plaintext, ours.nonce(), theirNonce, ours.privateKey(), theirPublicKey);

            Map<String, Object> outEntry = new LinkedHashMap<>();
            outEntry.put("care_context_id", entry.careContextId());
            outEntry.put("content", enc.encryptedData());
            outEntry.put("checksum", checksum(entry.bundle()));
            outEntry.put("media", FHIR_MEDIA);
            outEntries.add(outEntry);
        }

        Headers headers = new Headers(event.oid(), event.partnerPatientId(), event.hipId());
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("transaction_id", event.transactionId());
        // ponytail: single page. Add chunking if bundles ever exceed the payload limit.
        body.put("page_number", 1);
        body.put("page_count", 1);
        body.put("key_information", ourKeyInformation(ours, theirKeyInfo));
        body.put("entries", outEntries);

        http.doRaw("POST", "/abdm/v1/hip/care-context/data/on-fetch", headers, body);
    }

    /**
     * Returns the hex-encoded MD5 of the plaintext bundle, hashed before
     * encryption as the ABDM contract requires.
     *
     * <p>MD5 and hex encoding are fixed by the wire protocol, confirmed with
     * Eka. This is an interoperability requirement, not a security choice:
     * the HIU recomputes this digest over the decrypted bundle and discards
     * the records if it disagrees. Do not "upgrade" it to SHA-256 — the
     * receiver would reject every bundle, and it would fail silently,
     * because Eka accepts the push with 202 regardless and
     * {@link #respondToFetch} returns normally. Confidentiality comes from
     * the AES-256-GCM encryption of the content, not from this field.
     */
    private static String checksum(byte[] plaintext) {
        try {
            MessageDigest md5 = MessageDigest.getInstance("MD5"); // NOSONAR -- protocol-mandated digest, not a security primitive
            byte[] sum = md5.digest(plaintext);
            StringBuilder sb = new StringBuilder(sum.length * 2);
            for (byte b : sum) {
                sb.append(Character.forDigit((b >> 4) & 0xF, 16));
                sb.append(Character.forDigit(b & 0xF, 16));
            }
            return sb.toString();
        } catch (NoSuchAlgorithmException e) {
            // MD5 is mandated by the JDK's standard algorithm set; this cannot happen.
            throw new IllegalStateException(e);
        }
    }

    /**
     * Describes the ephemeral key material we generated for one response,
     * in the shape the HIU expects.
     *
     * <p>CONFIRM: the docs do not state what a HIP should send for
     * {@code dh_public_key} parameters and expiry. We echo the HIU's
     * parameters and set a 24h expiry.
     */
    private static Map<String, Object> ourKeyInformation(KeyMaterial ours, Map<String, Object> theirs) {
        Map<String, Object> theirDh = asMap(theirs.get("dh_public_key"));
        String params = str(theirDh.get("parameters"));
        if (params.isEmpty()) {
            params = "Curve25519/32byte random key";
        }
        String curve = str(theirs.get("curve"));
        if (curve.isEmpty()) {
            curve = "Curve25519";
        }

        Map<String, Object> dhPublicKey = new LinkedHashMap<>();
        dhPublicKey.put("key_value", ours.x509PublicKey());
        dhPublicKey.put("parameters", params);
        dhPublicKey.put("expiry", DateTimeFormatter.ISO_INSTANT.format(Instant.now().plusSeconds(24 * 60 * 60)));

        Map<String, Object> keyInformation = new LinkedHashMap<>();
        keyInformation.put("crypto_alg", "ECDH");
        keyInformation.put("curve", curve);
        keyInformation.put("nonce", ours.nonce());
        keyInformation.put("dh_public_key", dhPublicKey);
        return keyInformation;
    }

    /**
     * Answers an {@code abha.care_context_discover} webhook with the
     * unlinked care contexts belonging to the patient described in the
     * event.
     *
     * <p>Matching is yours: query your own patient records using the
     * event's name, gender, year of birth and identifiers. Match
     * conservatively — a loose match discloses another patient's records to
     * the requester. Set {@code result.error} when nothing matched.
     */
    public void onDiscover(DiscoverEvent event, DiscoverResult result) throws IOException, InterruptedException {
        if (event == null) {
            throw new IllegalArgumentException("carecontext: nil discover event");
        }
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("request_id", event.requestId());
        body.put("txn_id", event.txnId());
        if (result.error() != null) {
            body.put("error", result.error().toMap());
        } else {
            List<Map<String, Object>> patients = new ArrayList<>();
            for (Patient p : result.patients()) {
                patients.add(p.toMap());
            }
            body.put("patients", patients);
        }
        postDiscovery("/abdm/v1/care-contexts/on-discover", event.oid(), event.partnerPatientId(), event.hipId(), body);
    }

    /**
     * Answers an {@code abha.care_context_discover_link_init} webhook,
     * telling ABDM you have dispatched an OTP.
     *
     * <p>Generating the OTP, delivering it, and storing it against
     * {@code result.refNum} are all yours — this package holds no state.
     * The OTP must survive until the matching link-confirm webhook arrives.
     */
    public void onLinkInit(LinkInitEvent event, LinkInitResult result) throws IOException, InterruptedException {
        if (event == null) {
            throw new IllegalArgumentException("carecontext: nil link init event");
        }
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("request_id", event.requestId());
        body.put("txn_id", event.txnId());
        if (result.error() != null) {
            body.put("error", result.error().toMap());
        } else {
            if (!result.refNum().isEmpty()) {
                body.put("ref_num", result.refNum());
            }
            if (!result.otpExpiry().isEmpty()) {
                body.put("otp_expiry", result.otpExpiry());
            }
        }
        postDiscovery("/abdm/v1/care-contexts/discover/link/on-init", event.oid(), event.partnerPatientId(), event.hipId(), body);
    }

    /**
     * Answers an {@code abha.context_discover_link_confirm} webhook after
     * you have validated the OTP.
     *
     * <p>Validate {@code event.token()} against whatever you stored for
     * {@code event.linkRefNumber()}. On success send the care contexts to
     * link; on failure set {@code result.error}. This endpoint correlates
     * by request_id alone — the webhook carries no txn_id, and this call
     * sends none either.
     */
    public void onLinkConfirm(LinkConfirmEvent event, LinkConfirmResult result) throws IOException, InterruptedException {
        if (event == null) {
            throw new IllegalArgumentException("carecontext: nil link confirm event");
        }
        Map<String, Object> body = new LinkedHashMap<>();
        body.put("request_id", event.requestId());
        if (result.error() != null) {
            body.put("error", result.error().toMap());
        } else {
            List<Map<String, Object>> patients = new ArrayList<>();
            for (Patient p : result.patients()) {
                patients.add(p.toMap());
            }
            body.put("patients", patients);
        }
        postDiscovery("/abdm/v1/care-contexts/discover/link/on-confirm", event.oid(), event.partnerPatientId(), event.hipId(), body);
    }

    private void postDiscovery(String path, String oid, String partnerPatientId, String hipId, Map<String, Object> body)
            throws IOException, InterruptedException {
        Headers headers = new Headers(oid, partnerPatientId, hipId);
        http.doRaw("POST", path, headers, body);
    }

    private static String str(Object v) {
        return v == null ? "" : v.toString();
    }

    @SuppressWarnings("unchecked")
    private static Map<String, Object> asMap(Object v) {
        if (v instanceof Map<?, ?> m) {
            return (Map<String, Object>) m;
        }
        return Map.of();
    }
}

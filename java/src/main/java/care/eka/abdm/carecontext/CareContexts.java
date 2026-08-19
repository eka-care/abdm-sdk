package care.eka.abdm.carecontext;

import care.eka.abdm.Config;
import care.eka.abdm.Headers;
import care.eka.abdm.HttpTransport;
import java.io.IOException;

/**
 * Calls the ABDM care-context APIs. Mirrors
 * {@code go/services/abdm/carecontext/service.go}'s {@code Service}.
 *
 * <p>The package is boilerplate only. It carries bytes and correlation
 * identifiers; it never decides what your data means. Matching patients,
 * generating OTPs and building FHIR bundles stay in your code.
 */
public final class CareContexts {

    private final HttpTransport http;

    public CareContexts(Config config) {
        this.http = new HttpTransport(config, true);
    }

    /**
     * Links care contexts to a patient's ABHA address.
     *
     * <p>The API is asynchronous: a 202 means accepted, not linked. The
     * outcome arrives later as an {@code abha.link_care_context} webhook
     * (webhook parsing is Task 2 scope).
     *
     * <p>Does not mutate {@code request} — oid/partnerUserId are filled from
     * headers on a copy when left empty.
     */
    public void link(Headers headers, LinkRequest request) throws IOException, InterruptedException {
        LinkRequest body = request.withHeaderDefaults(headers);
        http.doRaw("POST", "/abdm/v1/care-contexts/link", headers, body.toMap());
    }
}

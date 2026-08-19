package care.eka.abdm.carecontext;

import care.eka.abdm.Headers;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

/**
 * The body of {@code POST /abdm/v1/care-contexts/link}. Mirrors Go's
 * {@code LinkRequest} / Python's {@code LinkRequest}.
 *
 * <p>{@code oid} and {@code partnerUserId} are filled from the request
 * headers by {@link CareContexts#link} when left empty — the contract
 * requires them in both the headers and the body. This type is immutable, so
 * that fill-in never mutates a caller's original request.
 */
public record LinkRequest(String abhaAddress, List<CareContext> careContexts, String oid, String partnerUserId) {

    public LinkRequest(String abhaAddress, List<CareContext> careContexts) {
        this(abhaAddress, careContexts, "", "");
    }

    public LinkRequest {
        careContexts = careContexts == null ? List.of() : List.copyOf(careContexts);
        oid = oid == null ? "" : oid;
        partnerUserId = partnerUserId == null ? "" : partnerUserId;
    }

    /** Returns a copy with oid/partnerUserId filled from headers where left empty. */
    LinkRequest withHeaderDefaults(Headers headers) {
        String newOid = oid.isEmpty() ? headers.patientId() : oid;
        String newPartnerUserId = partnerUserId.isEmpty() ? headers.partnerUserId() : partnerUserId;
        return new LinkRequest(abhaAddress, careContexts, newOid, newPartnerUserId);
    }

    Map<String, Object> toMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("abha_address", abhaAddress);
        m.put("care_contexts", careContexts.stream().map(CareContext::toMap).collect(Collectors.toCollection(ArrayList::new)));
        if (!oid.isEmpty()) {
            m.put("oid", oid);
        }
        if (!partnerUserId.isEmpty()) {
            m.put("partner_user_id", partnerUserId);
        }
        return m;
    }
}

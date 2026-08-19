package care.eka.abdm;

/**
 * The ABDM request identifiers sent with every call: the Eka patient id
 * (X-Pt-Id), your own patient id (X-Partner-Pt-Id) and your HIP id
 * (X-Hip-Id). Mirrors {@code interfaces.Headers} in Go and {@code Headers}
 * in Python's http module.
 */
public record Headers(String patientId, String partnerUserId, String hipId) {

    public Headers {
        patientId = patientId == null ? "" : patientId;
        partnerUserId = partnerUserId == null ? "" : partnerUserId;
        hipId = hipId == null ? "" : hipId;
    }

    public static Headers empty() {
        return new Headers("", "", "");
    }
}

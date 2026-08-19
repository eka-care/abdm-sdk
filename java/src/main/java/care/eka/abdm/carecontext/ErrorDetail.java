package care.eka.abdm.carecontext;

import java.util.LinkedHashMap;
import java.util.Map;

/**
 * Reports a failure back to ABDM through an {@code on-*} responder — no
 * patient matched, OTP invalid, and so on.
 */
public record ErrorDetail(int code, String message) {

    public ErrorDetail {
        message = message == null ? "" : message;
    }

    Map<String, Object> toMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("code", code);
        m.put("message", message);
        return m;
    }
}

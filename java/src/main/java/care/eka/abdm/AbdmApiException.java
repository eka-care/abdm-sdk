package care.eka.abdm;

import care.eka.abdm.internal.Json;
import java.nio.charset.StandardCharsets;
import java.util.Map;

/**
 * Thrown for any 4xx/5xx response from the ABDM API. Carries the HTTP status
 * code and the raw error body, and formats a message mirroring Go's
 * {@code ErrorResponse.String()} / Python's {@code _format_error_body}.
 */
public final class AbdmApiException extends RuntimeException {

    private final int statusCode;
    private final String errorBody;

    public AbdmApiException(int statusCode, String message, String errorBody) {
        super(message);
        this.statusCode = statusCode;
        this.errorBody = errorBody;
    }

    public int statusCode() {
        return statusCode;
    }

    /** The raw response body, or null if the response had none. */
    public String errorBody() {
        return errorBody;
    }

    static AbdmApiException fromResponse(int statusCode, byte[] raw) {
        String body = raw == null || raw.length == 0 ? null : new String(raw, StandardCharsets.UTF_8);
        return new AbdmApiException(statusCode, formatMessage(raw, statusCode), body);
    }

    private static String formatMessage(byte[] raw, int statusCode) {
        if (raw == null || raw.length == 0) {
            return "HTTP " + statusCode;
        }
        try {
            @SuppressWarnings("unchecked")
            Map<String, Object> parsed = Json.MAPPER.readValue(raw, Map.class);
            Object code = parsed.getOrDefault("code", statusCode);
            Object error = parsed.getOrDefault("error", "");
            Object sourceError = parsed.get("source_error");
            if (sourceError instanceof Map<?, ?> rawSe) {
                @SuppressWarnings("unchecked")
                Map<String, Object> se = (Map<String, Object>) rawSe;
                return "Error %s: %s (Source: %s - %s)".formatted(
                        code, error, se.getOrDefault("code", ""), se.getOrDefault("message", ""));
            }
            return "Error %s: %s".formatted(code, error);
        } catch (Exception e) {
            return new String(raw, StandardCharsets.UTF_8);
        }
    }
}

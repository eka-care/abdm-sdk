package care.eka.abdm.carecontext;

/**
 * The signature is valid but its timestamp is outside the &plusmn;3 minute
 * tolerance window. Mirrors Go's {@code ErrStaleTimestamp} / Python's
 * {@code StaleTimestampError}. Treat as HTTP 401.
 */
public final class StaleTimestampException extends RuntimeException {
    public StaleTimestampException(String message) {
        super(message);
    }
}

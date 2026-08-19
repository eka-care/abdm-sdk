package care.eka.abdm.carecontext;

/**
 * The payload is authentic but its event is not one this package handles.
 * Mirrors Go's {@code ErrUnknownEvent} / Python's {@code UnknownEventError}.
 * Distinguishable from {@link BadSignatureException} and
 * {@link StaleTimestampException} on purpose: answer 200 for this one (a
 * webhook endpoint must not fail on events it does not yet care about), and
 * 401 for an inauthentic payload.
 */
public final class UnknownEventException extends RuntimeException {
    public UnknownEventException(String message) {
        super(message);
    }
}

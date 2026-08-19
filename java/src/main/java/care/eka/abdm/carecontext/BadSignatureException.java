package care.eka.abdm.carecontext;

/**
 * The payload was not signed by Eka, or was altered, or the signature
 * header/secret was malformed/missing. Mirrors Go's {@code ErrBadSignature} /
 * Python's {@code BadSignatureError}. Treat as HTTP 401.
 *
 * <p>Kept distinct from {@link UnknownEventException} on purpose: an
 * unrecognised-but-authentic event must answer 200 (a webhook endpoint must
 * not fail on events it does not yet care about), while an inauthentic one
 * must answer 401.
 */
public final class BadSignatureException extends RuntimeException {
    public BadSignatureException(String message) {
        super(message);
    }
}

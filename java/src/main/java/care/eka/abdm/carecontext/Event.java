package care.eka.abdm.carecontext;

/**
 * Any decoded webhook payload. Mirrors Go's {@code Event} interface / Python's
 * duck-typed {@code event_name()} convention.
 */
public interface Event {
    String eventName();
}

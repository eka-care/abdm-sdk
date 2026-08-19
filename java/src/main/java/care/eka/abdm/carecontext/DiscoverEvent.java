package care.eka.abdm.carecontext;

import java.util.List;

/**
 * {@code abha.care_context_discover}: a patient is searching your facility
 * for their records. Match them in your own system, conservatively — a loose
 * match leaks another patient's records — and answer with
 * {@link CareContexts#onDiscover}.
 */
public record DiscoverEvent(
        String abhaAddress,
        String patientName,
        String gender,
        int yearOfBirth,
        List<Identifier> identifiers,
        String requestId,
        String txnId,
        String oid,
        String partnerPatientId,
        String hipId) implements Event {

    public DiscoverEvent {
        identifiers = identifiers == null ? List.of() : List.copyOf(identifiers);
    }

    @Override
    public String eventName() {
        return Webhook.EVENT_DISCOVER;
    }
}

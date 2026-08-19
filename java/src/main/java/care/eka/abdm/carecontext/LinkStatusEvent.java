package care.eka.abdm.carecontext;

/** {@code abha.link_care_context}: the async outcome of {@link CareContexts#link}. */
public record LinkStatusEvent(
        String abhaAddress,
        String careContextId,
        String status,
        String error,
        int retryCount,
        String oid,
        String partnerPatientId,
        String hipId) implements Event {

    @Override
    public String eventName() {
        return Webhook.EVENT_LINK_STATUS;
    }
}

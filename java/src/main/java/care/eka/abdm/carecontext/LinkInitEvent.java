package care.eka.abdm.carecontext;

/**
 * {@code abha.care_context_discover_link_init}: the patient chose care
 * contexts to link. Generate an OTP, send it, store it against the ref_num
 * you return, then answer with {@link CareContexts#onLinkInit}.
 *
 * <p>{@code patient} is the raw decoded {@code data.patient} JSON value
 * (typically a {@code List<Map<String, Object>>}), left uninterpreted since
 * this package does not decide what your data means.
 */
public record LinkInitEvent(
        String abhaAddress,
        Object patient,
        String requestId,
        String txnId,
        String oid,
        String partnerPatientId,
        String hipId) implements Event {

    @Override
    public String eventName() {
        return Webhook.EVENT_LINK_INIT;
    }
}

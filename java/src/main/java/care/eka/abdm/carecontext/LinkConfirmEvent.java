package care.eka.abdm.carecontext;

/**
 * {@code abha.context_discover_link_confirm}: the patient submitted the OTP.
 * Validate {@code token} against what you stored for {@code linkRefNumber},
 * then answer with {@link CareContexts#onLinkConfirm}.
 *
 * <p>This event carries no {@code txn_id} — {@code linkRefNumber} is the
 * correlation key, and matches the {@code refNum} you sent from
 * {@link CareContexts#onLinkInit}. Note the field itself is camelCase on the
 * wire ({@code linkRefNumber}) where every neighbouring field is
 * snake_case — that is not a typo, copy it verbatim.
 */
public record LinkConfirmEvent(
        String linkRefNumber,
        String token,
        String requestId,
        String oid,
        String partnerPatientId,
        String hipId) implements Event {

    @Override
    public String eventName() {
        return Webhook.EVENT_LINK_CONFIRM;
    }
}

package care.eka.abdm.carecontext;

/**
 * One care context's FHIR bundle, supplied by you and encrypted by
 * {@link CareContexts#respondToFetch}. {@code bundle} is ABDM-compliant FHIR
 * R4 JSON bytes.
 */
public record Entry(String careContextId, byte[] bundle) {
}

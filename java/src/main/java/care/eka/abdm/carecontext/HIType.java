package care.eka.abdm.carecontext;

/**
 * ABDM health information type constants. Mirrors Go's {@code HIType} consts.
 * Kept as plain strings rather than a closed Java {@code enum}: the wire
 * value is an open string set (Go's underlying type is {@code string}, and
 * any value ABDM defines in the future must still round-trip), so callers may
 * always pass a raw string here too.
 */
public final class HIType {

    private HIType() {
    }

    public static final String OP_CONSULTATION = "OPConsultation";
    public static final String PRESCRIPTION = "Prescription";
    public static final String DISCHARGE_SUMMARY = "DischargeSummary";
    public static final String DIAGNOSTIC_REPORT = "DiagnosticReport";
    public static final String IMMUNIZATION_RECORD = "ImmunizationRecord";
    public static final String HEALTH_DOCUMENT_RECORD = "HealthDocumentRecord";
    public static final String WELLNESS_RECORD = "WellnessRecord";
}

package care.eka.abdm.carecontext;

import java.util.List;

/**
 * Reports the outcome of OTP validation: the care contexts to link, or an
 * error if the OTP was wrong or expired.
 */
public record LinkConfirmResult(List<Patient> patients, ErrorDetail error) {

    public LinkConfirmResult(List<Patient> patients) {
        this(patients, null);
    }

    public static LinkConfirmResult ofError(ErrorDetail error) {
        return new LinkConfirmResult(List.of(), error);
    }

    public LinkConfirmResult {
        patients = patients == null ? List.of() : List.copyOf(patients);
    }
}

package care.eka.abdm.carecontext;

import java.util.List;

/**
 * Answers a discovery request: the patients you matched, or an error
 * explaining why you could not.
 */
public record DiscoverResult(List<Patient> patients, ErrorDetail error) {

    public DiscoverResult(List<Patient> patients) {
        this(patients, null);
    }

    public static DiscoverResult ofError(ErrorDetail error) {
        return new DiscoverResult(List.of(), error);
    }

    public DiscoverResult {
        patients = patients == null ? List.of() : List.copyOf(patients);
    }
}

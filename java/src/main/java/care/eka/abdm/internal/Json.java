package care.eka.abdm.internal;

import com.fasterxml.jackson.databind.ObjectMapper;

/** The single shared Jackson mapper used across the SDK's HTTP layer. */
public final class Json {

    private Json() {
    }

    public static final ObjectMapper MAPPER = new ObjectMapper();
}

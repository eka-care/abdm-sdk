package care.eka.abdm.carecontext;

/** A verified patient identifier supplied during discovery. */
public record Identifier(String type, String value) {

    public Identifier {
        type = type == null ? "" : type;
        value = value == null ? "" : value;
    }
}

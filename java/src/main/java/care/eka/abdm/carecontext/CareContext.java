package care.eka.abdm.carecontext;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

/**
 * One logical group of records — a visit, a test, a document. Mirrors Go's
 * {@code CareContext} / Python's {@code CareContext}.
 */
public record CareContext(String careContextId, String display, String hiType, List<String> hiTypes) {

    public CareContext(String careContextId, String display) {
        this(careContextId, display, "", null);
    }

    public CareContext {
        hiType = hiType == null ? "" : hiType;
    }

    Map<String, Object> toMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("care_context_id", careContextId);
        m.put("display", display);
        if (!hiType.isEmpty()) {
            m.put("hi_type", hiType);
        }
        if (hiTypes != null && !hiTypes.isEmpty()) {
            m.put("hi_types", hiTypes);
        }
        return m;
    }
}

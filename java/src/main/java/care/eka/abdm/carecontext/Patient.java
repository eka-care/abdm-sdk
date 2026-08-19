package care.eka.abdm.carecontext;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

/** A matched patient and their unlinked care contexts. */
public record Patient(String refNum, String display, String hiType, List<DiscoveredCareContext> careContexts) {

    public Patient(String refNum, String display) {
        this(refNum, display, "", null);
    }

    public Patient {
        hiType = hiType == null ? "" : hiType;
        careContexts = careContexts == null ? List.of() : List.copyOf(careContexts);
    }

    Map<String, Object> toMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("ref_num", refNum);
        m.put("display", display);
        if (!hiType.isEmpty()) {
            m.put("hi_type", hiType);
        }
        m.put("care_contexts",
                careContexts.stream().map(DiscoveredCareContext::toMap).collect(Collectors.toCollection(ArrayList::new)));
        return m;
    }
}

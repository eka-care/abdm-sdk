package care.eka.abdm.carecontext;

import java.util.LinkedHashMap;
import java.util.Map;

/** One unlinked care context offered during discovery. */
public record DiscoveredCareContext(String refNum, String display) {

    public DiscoveredCareContext {
        refNum = refNum == null ? "" : refNum;
        display = display == null ? "" : display;
    }

    Map<String, Object> toMap() {
        Map<String, Object> m = new LinkedHashMap<>();
        m.put("ref_num", refNum);
        m.put("display", display);
        return m;
    }
}

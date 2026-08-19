package care.eka.abdm.carecontext;

import java.util.List;
import java.util.Map;

/**
 * {@code abha.hip_data_fetch}: an HIU wants records. Answer with
 * {@link CareContexts#respondToFetch}. The HIU's key material is captured
 * but not exposed publicly ({@link #keyInfo()} is package-private) —
 * {@code respondToFetch} uses it so you never handle key material yourself.
 *
 * <p>Not a record: the canonical accessor for {@code keyInformation} must
 * stay package-private, which a record's generated accessor cannot do.
 */
public final class DataFetchEvent implements Event {

    private final String transactionId;
    private final String abhaAddress;
    private final String oid;
    private final String partnerPatientId;
    private final String hipId;
    private final List<String> careContexts;
    private final List<String> hiTypes;
    private final Map<String, Object> keyInfo;

    public DataFetchEvent(
            String transactionId,
            String abhaAddress,
            String oid,
            String partnerPatientId,
            String hipId,
            List<String> careContexts,
            List<String> hiTypes,
            Map<String, Object> keyInfo) {
        this.transactionId = transactionId == null ? "" : transactionId;
        this.abhaAddress = abhaAddress == null ? "" : abhaAddress;
        this.oid = oid == null ? "" : oid;
        this.partnerPatientId = partnerPatientId == null ? "" : partnerPatientId;
        this.hipId = hipId == null ? "" : hipId;
        this.careContexts = careContexts == null ? List.of() : List.copyOf(careContexts);
        this.hiTypes = hiTypes == null ? List.of() : List.copyOf(hiTypes);
        this.keyInfo = keyInfo == null ? Map.of() : Map.copyOf(keyInfo);
    }

    @Override
    public String eventName() {
        return Webhook.EVENT_DATA_FETCH;
    }

    public String transactionId() {
        return transactionId;
    }

    public String abhaAddress() {
        return abhaAddress;
    }

    public String oid() {
        return oid;
    }

    public String partnerPatientId() {
        return partnerPatientId;
    }

    public String hipId() {
        return hipId;
    }

    public List<String> careContexts() {
        return careContexts;
    }

    public List<String> hiTypes() {
        return hiTypes;
    }

    /** The HIU's ECDH key material, as a raw decoded JSON map. Package-private:
     * only {@link CareContexts#respondToFetch} reads this, so no caller ever
     * has to handle key material directly. */
    Map<String, Object> keyInfo() {
        return keyInfo;
    }
}

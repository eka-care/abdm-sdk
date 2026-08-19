package care.eka.abdm.carecontext;

/**
 * Reports that you sent an OTP. {@code refNum} is your reference for this
 * linking attempt; it returns as {@link LinkConfirmEvent#linkRefNumber()},
 * so store your OTP against it. {@code otpExpiry} is an ISO-8601 timestamp.
 */
public record LinkInitResult(String refNum, String otpExpiry, ErrorDetail error) {

    public LinkInitResult(String refNum, String otpExpiry) {
        this(refNum, otpExpiry, null);
    }

    public static LinkInitResult ofError(ErrorDetail error) {
        return new LinkInitResult("", "", error);
    }

    public LinkInitResult {
        refNum = refNum == null ? "" : refNum;
        otpExpiry = otpExpiry == null ? "" : otpExpiry;
    }
}

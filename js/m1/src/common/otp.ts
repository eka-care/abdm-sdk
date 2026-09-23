import type { OtpSystem } from "./types.js";

/** The second scope NHA expects for an OTP sent by the given system. */
export const verifyScope = (system: OtpSystem) => (system === "aadhaar" ? "aadhaar-verify" : "mobile-verify");

export const otpAuthData = (txnId: string, encryptedOtp: string) => ({
  authMethods: ["otp"],
  otp: { txnId, otpValue: encryptedOtp },
});

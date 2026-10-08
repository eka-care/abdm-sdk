import type { AuthResult, Tokens } from "../common/types.js";
import type { AbhaProfile } from "../profile/profile.types.js";

export interface VerifyAadhaarOtpRequest {
  txnId: string;
  otp: string;
  /** The number to put on the ABHA. If Aadhaar has a different one, verify it with `requestMobileOtp`. */
  mobile: string;
}

export interface EnrolResponse {
  txnId: string;
  message: string;
  tokens: Tokens;
  ABHAProfile: AbhaProfile;
  /** false when an ABHA already existed for this Aadhaar. */
  isNew: boolean;
}

export interface MobileOtpRequest {
  txnId: string;
  mobile: string;
}

export interface VerifyMobileOtpRequest {
  txnId: string;
  otp: string;
}

export interface MobileVerifyResponse {
  txnId: string;
  authResult: AuthResult;
  message: string;
  accounts: { ABHANumber: string }[];
}

export interface CreateAbhaAddressRequest {
  txnId: string;
  abhaAddress: string;
  /** Defaults to true. */
  preferred?: boolean;
}

export interface CreateAbhaAddressResponse {
  txnId: string;
  healthIdNumber: string;
  preferredAbhaAddress: string;
}

import type { AuthResult, OtpResponse, OtpSystem, Tokens } from "../common/types.js";

export type LoginHint = "aadhaar" | "mobile" | "abha-number";

export interface LoginOtpRequest {
  loginHint: LoginHint;
  /** Aadhaar, mobile or ABHA number in plain text. */
  value: string;
  /** Defaults to `aadhaar` for an Aadhaar identifier and `abdm` otherwise. */
  otpSystem?: OtpSystem;
}

export interface LoginOtpResponse extends OtpResponse {
  /** Pass this back to `verifyOtp`. */
  otpSystem: OtpSystem;
}

export interface VerifyLoginOtpRequest {
  txnId: string;
  otp: string;
  otpSystem: OtpSystem;
}

export interface LoginAccount {
  ABHANumber: string;
  name: string;
  preferredAbhaAddress: string;
  status: string;
  gender?: string;
  dob?: string;
  profilePhoto?: string;
  verificationType?: string;
  verifiedStatus?: string;
}

export interface LoginVerifyResponse extends Tokens {
  txnId: string;
  authResult: AuthResult;
  message: string;
  accounts: LoginAccount[];
}

export interface VerifyUserRequest {
  txnId: string;
  abhaNumber: string;
  /** The `token` from `verifyOtp` when several accounts matched. */
  tToken: string;
}

export type VerifyUserResponse = Required<Tokens>;

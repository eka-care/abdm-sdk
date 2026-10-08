import type { AddressUser, AuthResult, OtpSystem, Tokens } from "../common/types.js";

export interface AddressSearchResponse {
  abhaAddress: string;
  healthIdNumber?: string;
  fullName: string;
  mobile?: string;
  status: string;
  authMethods: string[];
  blockedAuthMethods: string[];
}

export interface AddressOtpRequest {
  abhaAddress: string;
  otpSystem: OtpSystem;
}

export interface VerifyAddressOtpRequest {
  txnId: string;
  otp: string;
  otpSystem: OtpSystem;
}

export interface AddressVerifyResponse {
  txnId?: string;
  authResult: AuthResult;
  message: string;
  /** `tokens.token` is the X-token for `getProfile` and `getCard`. */
  tokens?: Tokens;
  users: AddressUser[];
}

import type { AddressUser, AuthResult, Tokens } from "../common/types.js";

export interface VerifyPhrOtpRequest {
  txnId: string;
  otp: string;
}

export interface PhrVerifyResponse {
  txnId: string;
  authResult: AuthResult;
  message: string;
  tokens?: Tokens;
  /** ABHA addresses already on this mobile. */
  users?: AddressUser[];
}

export interface PhrSuggestionRequest {
  txnId: string;
  firstName: string;
  lastName?: string;
  dayOfBirth: string;
  monthOfBirth: string;
  yearOfBirth: string;
  email?: string;
}

export interface PhrDetails {
  abhaAddress: string;
  firstName: string;
  middleName?: string;
  lastName?: string;
  gender: "M" | "F" | "O";
  dayOfBirth: string;
  monthOfBirth: string;
  yearOfBirth: string;
  email?: string;
  address?: string;
  pinCode?: string;
  stateCode?: string;
  stateName?: string;
  districtCode?: string;
  districtName?: string;
  profilePhoto?: string;
  /** Plain text; the client encrypts it. */
  password?: string;
}

export interface CreatePhrAddressRequest extends PhrDetails {
  txnId: string;
  /** Plain text; the client encrypts it. */
  mobile: string;
}

export interface PhrEnrolResponse {
  txnId: string;
  message: string;
  tokens: Tokens;
  phrDetails: Omit<PhrDetails, "abhaAddress" | "password"> & {
    abhaAddress: string[];
    fullName: string;
    mobile: string;
  };
}

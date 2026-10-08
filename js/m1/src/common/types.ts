export type OtpSystem = "aadhaar" | "abdm";

export type AuthResult = "success" | "failed";

export interface OtpResponse {
  txnId: string;
  message: string;
}

export interface Tokens {
  token: string;
  expiresIn: number;
  refreshToken?: string;
  refreshExpiresIn?: number;
}

export interface AddressSuggestions {
  txnId: string;
  abhaAddressList: string[];
}

/** An ABHA address as NHA lists it after an OTP is verified. */
export interface AddressUser {
  abhaAddress: string;
  abhaNumber?: string;
  fullName: string;
  kycStatus?: string;
  status: string;
  profilePhoto?: string;
}

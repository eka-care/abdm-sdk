import type { CallOptions } from "../common/config.js";
import type { Encryptor } from "../common/encryptor.js";
import type { HttpClient } from "../common/http.js";
import { otpAuthData } from "../common/otp.js";
import type { AddressSuggestions, OtpResponse } from "../common/types.js";
import type {
  CreateAbhaAddressRequest,
  CreateAbhaAddressResponse,
  EnrolResponse,
  MobileOtpRequest,
  MobileVerifyResponse,
  VerifyAadhaarOtpRequest,
  VerifyMobileOtpRequest,
} from "./aadhaar.types.js";

const ENROL_CONSENT = { code: "abha-enrollment", version: "1.4" };

/** Creates an ABHA number with Aadhaar OTP, then an ABHA address on it. */
export class AadhaarService {
  constructor(
    private readonly http: HttpClient,
    private readonly encryptor: Encryptor,
  ) {}

  requestOtp(aadhaar: string, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/enrollment/request/otp", {
      body: {
        scope: ["abha-enrol"],
        loginHint: "aadhaar",
        loginId: this.encryptor.encrypt(aadhaar),
        otpSystem: "aadhaar",
      },
    }, opts);
  }

  verifyOtp(req: VerifyAadhaarOtpRequest, opts?: CallOptions): Promise<EnrolResponse> {
    const otpValue = this.encryptor.encrypt(req.otp);
    return this.http.request("POST", "/enrollment/enrol/byAadhaar", {
      body: {
        authData: { authMethods: ["otp"], otp: { txnId: req.txnId, otpValue, mobile: req.mobile } },
        consent: ENROL_CONSENT,
      },
    }, opts);
  }

  requestMobileOtp(req: MobileOtpRequest, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/enrollment/request/otp", {
      body: {
        txnId: req.txnId,
        scope: ["abha-enrol", "mobile-verify"],
        loginHint: "mobile",
        loginId: this.encryptor.encrypt(req.mobile),
        otpSystem: "abdm",
      },
    }, opts);
  }

  verifyMobileOtp(req: VerifyMobileOtpRequest, opts?: CallOptions): Promise<MobileVerifyResponse> {
    const otp = this.encryptor.encrypt(req.otp);
    return this.http.request("POST", "/enrollment/auth/byAbdm", {
      body: { scope: ["abha-enrol", "mobile-verify"], authData: otpAuthData(req.txnId, otp) },
    }, opts);
  }

  suggestAddresses(txnId: string, opts?: CallOptions): Promise<AddressSuggestions> {
    return this.http.request("GET", "/enrollment/enrol/suggestion", { headers: { TRANSACTION_ID: txnId } }, opts);
  }

  createAddress(req: CreateAbhaAddressRequest, opts?: CallOptions): Promise<CreateAbhaAddressResponse> {
    return this.http.request("POST", "/enrollment/enrol/abha-address", {
      body: { txnId: req.txnId, abhaAddress: req.abhaAddress, preferred: req.preferred === false ? 0 : 1 },
    }, opts);
  }
}

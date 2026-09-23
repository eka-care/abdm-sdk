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

  async requestOtp(aadhaar: string, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/v3/enrollment/request/otp", {
      body: {
        scope: ["abha-enrol"],
        loginHint: "aadhaar",
        loginId: await this.encryptor.encrypt(aadhaar, opts),
        otpSystem: "aadhaar",
      },
    }, opts);
  }

  async verifyOtp(req: VerifyAadhaarOtpRequest, opts?: CallOptions): Promise<EnrolResponse> {
    const otpValue = await this.encryptor.encrypt(req.otp, opts);
    return this.http.request("POST", "/v3/enrollment/enrol/byAadhaar", {
      body: {
        authData: { authMethods: ["otp"], otp: { txnId: req.txnId, otpValue, mobile: req.mobile } },
        consent: ENROL_CONSENT,
      },
    }, opts);
  }

  async requestMobileOtp(req: MobileOtpRequest, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/v3/enrollment/request/otp", {
      body: {
        txnId: req.txnId,
        scope: ["abha-enrol", "mobile-verify"],
        loginHint: "mobile",
        loginId: await this.encryptor.encrypt(req.mobile, opts),
        otpSystem: "abdm",
      },
    }, opts);
  }

  async verifyMobileOtp(req: VerifyMobileOtpRequest, opts?: CallOptions): Promise<MobileVerifyResponse> {
    const otp = await this.encryptor.encrypt(req.otp, opts);
    return this.http.request("POST", "/v3/enrollment/auth/byAbdm", {
      body: { scope: ["abha-enrol", "mobile-verify"], authData: otpAuthData(req.txnId, otp) },
    }, opts);
  }

  suggestAddresses(txnId: string, opts?: CallOptions): Promise<AddressSuggestions> {
    return this.http.request("GET", "/v3/enrollment/enrol/suggestion", { headers: { TRANSACTION_ID: txnId } }, opts);
  }

  createAddress(req: CreateAbhaAddressRequest, opts?: CallOptions): Promise<CreateAbhaAddressResponse> {
    return this.http.request("POST", "/v3/enrollment/enrol/abha-address", {
      body: { txnId: req.txnId, abhaAddress: req.abhaAddress, preferred: req.preferred === false ? 0 : 1 },
    }, opts);
  }
}

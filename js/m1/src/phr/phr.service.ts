import type { CallOptions } from "../common/config.js";
import type { Encryptor } from "../common/encryptor.js";
import type { HttpClient } from "../common/http.js";
import { otpAuthData } from "../common/otp.js";
import type { AddressSuggestions, OtpResponse } from "../common/types.js";
import type {
  CreatePhrAddressRequest,
  PhrEnrolResponse,
  PhrSuggestionRequest,
  PhrVerifyResponse,
  VerifyPhrOtpRequest,
} from "./phr.types.js";

const SCOPE = ["abha-address-enroll", "mobile-verify"];

/** Creates an ABHA address (PHR) from a mobile number, without an ABHA number. */
export class PhrService {
  constructor(
    private readonly http: HttpClient,
    private readonly encryptor: Encryptor,
  ) {}

  requestOtp(mobile: string, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/enrollment/request/otp", {
      body: {
        scope: SCOPE,
        loginHint: "mobile-number",
        loginId: this.encryptor.encrypt(mobile),
        otpSystem: "abdm",
      },
    }, opts);
  }

  verifyOtp(req: VerifyPhrOtpRequest, opts?: CallOptions): Promise<PhrVerifyResponse> {
    const otp = this.encryptor.encrypt(req.otp);
    return this.http.request("POST", "/enrollment/verify", {
      body: { scope: SCOPE, authData: otpAuthData(req.txnId, otp) },
    }, opts);
  }

  suggestAddresses(req: PhrSuggestionRequest, opts?: CallOptions): Promise<AddressSuggestions> {
    return this.http.request("POST", "/enrollment/suggestion", {
      body: { lastName: "", email: "", ...req },
    }, opts);
  }

  addressExists(abhaAddress: string, opts?: CallOptions): Promise<boolean> {
    return this.http.request("GET", "/enrollment/isExists", { query: { abhaAddress } }, opts);
  }

  createAddress(req: CreatePhrAddressRequest, opts?: CallOptions): Promise<PhrEnrolResponse> {
    const { txnId, mobile, password, ...details } = req;
    const phrDetails = {
      ...details,
      mobile: this.encryptor.encrypt(mobile),
      ...(password && { password: this.encryptor.encrypt(password) }),
    };
    return this.http.request("POST", "/enrollment/enrol", { body: { txnId, phrDetails } }, opts);
  }
}

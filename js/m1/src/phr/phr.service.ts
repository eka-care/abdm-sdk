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

  async requestOtp(mobile: string, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/v3/phr/app/enrollment/request/otp", {
      body: {
        scope: SCOPE,
        loginHint: "mobile-number",
        loginId: await this.encryptor.encrypt(mobile, opts),
        otpSystem: "abdm",
      },
    }, opts);
  }

  async verifyOtp(req: VerifyPhrOtpRequest, opts?: CallOptions): Promise<PhrVerifyResponse> {
    const otp = await this.encryptor.encrypt(req.otp, opts);
    return this.http.request("POST", "/v3/phr/app/enrollment/verify", {
      body: { scope: SCOPE, authData: otpAuthData(req.txnId, otp) },
    }, opts);
  }

  suggestAddresses(req: PhrSuggestionRequest, opts?: CallOptions): Promise<AddressSuggestions> {
    return this.http.request("POST", "/v3/phr/app/enrollment/suggestion", {
      body: { lastName: "", email: "", ...req },
    }, opts);
  }

  addressExists(abhaAddress: string, opts?: CallOptions): Promise<boolean> {
    return this.http.request("GET", "/v3/phr/app/enrollment/isExists", { query: { abhaAddress } }, opts);
  }

  async createAddress(req: CreatePhrAddressRequest, opts?: CallOptions): Promise<PhrEnrolResponse> {
    const { txnId, mobile, password, ...details } = req;
    const phrDetails = {
      ...details,
      mobile: await this.encryptor.encrypt(mobile, opts),
      ...(password && { password: await this.encryptor.encrypt(password, opts) }),
    };
    return this.http.request("POST", "/v3/phr/app/enrollment/enrol", { body: { txnId, phrDetails } }, opts);
  }
}

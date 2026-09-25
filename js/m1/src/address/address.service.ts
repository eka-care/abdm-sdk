import type { CallOptions } from "../common/config.js";
import type { Encryptor } from "../common/encryptor.js";
import { bearer, type HttpClient } from "../common/http.js";
import { otpAuthData, verifyScope } from "../common/otp.js";
import type { OtpResponse } from "../common/types.js";
import type {
  AddressOtpRequest,
  AddressSearchResponse,
  AddressVerifyResponse,
  VerifyAddressOtpRequest,
} from "./address.types.js";

/** Verifies an ABHA address and lists the addresses behind it. */
export class AddressService {
  constructor(
    private readonly http: HttpClient,
    private readonly encryptor: Encryptor,
  ) {}

  search(abhaAddress: string, opts?: CallOptions): Promise<AddressSearchResponse> {
    return this.http.request("POST", "/login/abha/search", { body: { abhaAddress } }, opts);
  }

  requestOtp(req: AddressOtpRequest, opts?: CallOptions): Promise<OtpResponse> {
    return this.http.request("POST", "/login/abha/request/otp", {
      body: {
        scope: ["abha-address-login", verifyScope(req.otpSystem)],
        loginHint: "abha-address",
        loginId: this.encryptor.encrypt(req.abhaAddress),
        otpSystem: req.otpSystem,
      },
    }, opts);
  }

  verifyOtp(req: VerifyAddressOtpRequest, opts?: CallOptions): Promise<AddressVerifyResponse> {
    const otp = this.encryptor.encrypt(req.otp);
    return this.http.request("POST", "/login/abha/verify", {
      body: { scope: ["abha-address-login", verifyScope(req.otpSystem)], authData: otpAuthData(req.txnId, otp) },
    }, opts);
  }

  getProfile(xToken: string, opts?: CallOptions): Promise<Record<string, unknown>> {
    return this.http.request("GET", "/login/profile/abha-profile", { headers: { "X-token": bearer(xToken) } }, opts);
  }

  getCard(xToken: string, opts?: CallOptions): Promise<ArrayBuffer> {
    return this.http.request("GET", "/login/profile/abha/phr-card", { headers: { "X-token": bearer(xToken) } }, opts);
  }
}

import type { CallOptions } from "../common/config.js";
import type { Encryptor } from "../common/encryptor.js";
import { bearer, type HttpClient } from "../common/http.js";
import { otpAuthData, verifyScope } from "../common/otp.js";
import type { OtpResponse } from "../common/types.js";
import type {
  LoginOtpRequest,
  LoginOtpResponse,
  LoginVerifyResponse,
  VerifyLoginOtpRequest,
  VerifyUserRequest,
  VerifyUserResponse,
} from "./login.types.js";

/** Logs in to an ABHA number with an Aadhaar, mobile or ABHA number identifier. */
export class LoginService {
  constructor(
    private readonly http: HttpClient,
    private readonly encryptor: Encryptor,
  ) {}

  async requestOtp(req: LoginOtpRequest, opts?: CallOptions): Promise<LoginOtpResponse> {
    const otpSystem = req.otpSystem ?? (req.loginHint === "aadhaar" ? "aadhaar" : "abdm");
    const res = await this.http.request<OtpResponse>("POST", "/profile/login/request/otp", {
      body: {
        scope: ["abha-login", verifyScope(otpSystem)],
        loginHint: req.loginHint,
        loginId: this.encryptor.encrypt(req.value),
        otpSystem,
      },
    }, opts);
    return { ...res, otpSystem };
  }

  /** Lists matching `accounts`. With several, pick one and call `verifyUser` with `token`. */
  verifyOtp(req: VerifyLoginOtpRequest, opts?: CallOptions): Promise<LoginVerifyResponse> {
    const otp = this.encryptor.encrypt(req.otp);
    return this.http.request("POST", "/profile/login/verify", {
      body: { scope: ["abha-login", verifyScope(req.otpSystem)], authData: otpAuthData(req.txnId, otp) },
    }, opts);
  }

  verifyUser(req: VerifyUserRequest, opts?: CallOptions): Promise<VerifyUserResponse> {
    return this.http.request("POST", "/profile/login/verify/user", {
      body: { txnId: req.txnId, ABHANumber: req.abhaNumber },
      headers: { "T-token": bearer(req.tToken) },
    }, opts);
  }
}

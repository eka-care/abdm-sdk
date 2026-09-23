import type { CallOptions } from "../common/config.js";
import { bearer, type HttpClient } from "../common/http.js";
import type { AbhaProfile } from "./profile.types.js";

/** ABHA number profile, using the X-token from enrolment or login. */
export class ProfileService {
  constructor(private readonly http: HttpClient) {}

  /** `phrAddress` lists every ABHA address on the account. */
  get(xToken: string, opts?: CallOptions): Promise<AbhaProfile> {
    return this.http.request("GET", "/v3/profile/account", { headers: { "X-token": bearer(xToken) } }, opts);
  }

  getCard(xToken: string, opts?: CallOptions): Promise<ArrayBuffer> {
    return this.http.request("GET", "/v3/profile/account/abha-card", { headers: { "X-token": bearer(xToken) } }, opts);
  }
}

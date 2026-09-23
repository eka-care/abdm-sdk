import { AadhaarService } from "./aadhaar/aadhaar.service.js";
import { AddressService } from "./address/address.service.js";
import type { CallOptions, M1Config } from "./common/config.js";
import { AbhaEncryptor, type Encryptor } from "./common/encryptor.js";
import { AbhaHttpClient, type HttpClient } from "./common/http.js";
import { LoginService } from "./login/login.service.js";
import { PhrService } from "./phr/phr.service.js";
import { ProfileService } from "./profile/profile.service.js";

/** Entry point for ABDM M1. Each flow is a separate service sharing one HTTP client and encryptor. */
export class M1Client {
  readonly aadhaar: AadhaarService;
  readonly mobile: PhrService;
  readonly login: LoginService;
  readonly address: AddressService;
  readonly profile: ProfileService;

  private readonly encryptor: Encryptor;

  constructor(config: M1Config) {
    const http: HttpClient = new AbhaHttpClient(config);
    this.encryptor = new AbhaEncryptor(http);

    this.aadhaar = new AadhaarService(http, this.encryptor);
    this.mobile = new PhrService(http, this.encryptor);
    this.login = new LoginService(http, this.encryptor);
    this.address = new AddressService(http, this.encryptor);
    this.profile = new ProfileService(http);
  }

  /** Encrypt a value with the ABHA public key, for calls this client does not wrap. */
  encrypt(value: string, opts?: CallOptions): Promise<string> {
    return this.encryptor.encrypt(value, opts);
  }
}

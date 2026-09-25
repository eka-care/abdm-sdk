import { AadhaarService } from "./aadhaar/aadhaar.service.js";
import { AddressService } from "./address/address.service.js";
import { BASE_URLS, type M1Config } from "./common/config.js";
import { RsaOaepEncryptor } from "./common/encryptor.js";
import { NhaHttpClient } from "./common/http.js";
import { ABHA_PUBLIC_KEY, PHR_PUBLIC_KEY } from "./common/keys.js";
import { LoginService } from "./login/login.service.js";
import { PhrService } from "./phr/phr.service.js";
import { ProfileService } from "./profile/profile.service.js";

/** Entry point for ABDM M1. Wires each service to its NHA API family and encryption key. */
export class M1Client {
  readonly aadhaar: AadhaarService;
  readonly mobile: PhrService;
  readonly login: LoginService;
  readonly address: AddressService;
  readonly profile: ProfileService;

  private readonly abhaEncryptor: RsaOaepEncryptor;

  constructor(config: M1Config) {
    if (!config.accessToken) throw new Error("accessToken is required");

    const urls = { ...BASE_URLS[config.environment ?? "sandbox"], ...config.baseUrls };
    const http = (baseUrl: string) => new NhaHttpClient(baseUrl, config.accessToken, config.fetch);
    const abhaHttp = http(urls.abha);

    this.abhaEncryptor = new RsaOaepEncryptor(config.publicKeys?.abha ?? ABHA_PUBLIC_KEY);
    const phrEncryptor = new RsaOaepEncryptor(config.publicKeys?.phr ?? PHR_PUBLIC_KEY);

    this.aadhaar = new AadhaarService(abhaHttp, this.abhaEncryptor);
    this.login = new LoginService(abhaHttp, this.abhaEncryptor);
    this.profile = new ProfileService(abhaHttp);
    this.mobile = new PhrService(http(urls.phrApp), phrEncryptor);
    this.address = new AddressService(http(urls.phrWeb), this.abhaEncryptor);
  }

  /** Encrypt a value with the ABHA public key, for calls this client does not wrap. */
  encrypt(value: string): string {
    return this.abhaEncryptor.encrypt(value);
  }
}

import { constants, createPublicKey, publicEncrypt, type KeyObject } from "node:crypto";
import type { CallOptions } from "./config.js";
import type { HttpClient } from "./http.js";

export interface Encryptor {
  encrypt(value: string, opts?: CallOptions): Promise<string>;
}

/** RSA/ECB/OAEPWithSHA-1AndMGF1Padding with the ABHA public certificate, as NHA requires. */
export class AbhaEncryptor implements Encryptor {
  private publicKey?: Promise<KeyObject>;

  constructor(private readonly http: HttpClient) {}

  async encrypt(value: string, opts: CallOptions = {}): Promise<string> {
    return encryptWith(await this.loadPublicKey(opts.accessToken), value);
  }

  private loadPublicKey(accessToken?: string): Promise<KeyObject> {
    this.publicKey ??= this.http
      .request<{ publicKey: string }>("GET", "/v3/profile/public/certificate", {}, { accessToken })
      .then(({ publicKey }) => createPublicKey({ key: Buffer.from(publicKey, "base64"), format: "der", type: "spki" }))
      .catch((err) => {
        this.publicKey = undefined;
        throw err;
      });
    return this.publicKey;
  }
}

export function encryptWith(key: KeyObject, value: string): string {
  return publicEncrypt({ key, padding: constants.RSA_PKCS1_OAEP_PADDING, oaepHash: "sha1" }, Buffer.from(value)).toString("base64");
}

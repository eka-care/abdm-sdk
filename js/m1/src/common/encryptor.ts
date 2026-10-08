import { constants, createPublicKey, publicEncrypt, type KeyObject } from "node:crypto";

export interface Encryptor {
  encrypt(value: string): string;
}

/** RSA/ECB/OAEPWithSHA-1AndMGF1Padding, as NHA requires for identifiers, OTPs and passwords. */
export class RsaOaepEncryptor implements Encryptor {
  private readonly key: KeyObject;

  constructor(publicKeyPem: string) {
    this.key = createPublicKey(publicKeyPem);
  }

  encrypt(value: string): string {
    return publicEncrypt({ key: this.key, padding: constants.RSA_PKCS1_OAEP_PADDING, oaepHash: "sha1" }, Buffer.from(value))
      .toString("base64");
  }
}

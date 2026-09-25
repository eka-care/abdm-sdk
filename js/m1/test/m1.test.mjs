import { test } from "node:test";
import assert from "node:assert/strict";
import { constants, generateKeyPairSync, privateDecrypt } from "node:crypto";
import { AbdmError, M1Client } from "../dist/index.js";

const keyPair = () => generateKeyPairSync("rsa", { modulusLength: 2048, publicKeyEncoding: { type: "spki", format: "pem" } });
const abhaKeys = keyPair();
const phrKeys = keyPair();
const decrypt = ({ privateKey }, b64) =>
  privateDecrypt({ key: privateKey, padding: constants.RSA_PKCS1_OAEP_PADDING, oaepHash: "sha1" }, Buffer.from(b64, "base64")).toString();

function fakeNha() {
  const calls = [];
  const json = (body, status = 200) => new Response(JSON.stringify(body), { status, headers: { "content-type": "application/json" } });
  const fetch = async (url, init) => {
    calls.push({ url: String(url), headers: init.headers, body: init.body && JSON.parse(init.body) });
    if (String(url).endsWith("/profile/login/verify")) return json([{ code: "ABDM-1006", message: "Invalid OTP Request" }], 400);
    return json({ txnId: "t1", message: "sent" });
  };
  return { calls, fetch };
}

const client = (opts) => new M1Client({ publicKeys: { abha: abhaKeys.publicKey, phr: phrKeys.publicKey }, ...opts });

test("sends NHA headers and encrypts ABHA calls with the ABHA key", async () => {
  const { calls, fetch } = fakeNha();
  let n = 0;
  const res = await client({ accessToken: () => `token-${++n}`, fetch }).aadhaar.requestOtp("999941057058");

  assert.equal(res.txnId, "t1");
  const [call] = calls;
  assert.equal(call.url, "https://abhasbx.abdm.gov.in/abha/api/v3/enrollment/request/otp");
  assert.equal(call.headers.Authorization, "Bearer token-1");
  assert.match(call.headers["REQUEST-ID"], /^[0-9a-f-]{36}$/);
  assert.ok(!Number.isNaN(Date.parse(call.headers.TIMESTAMP)));
  assert.equal(decrypt(abhaKeys, call.body.loginId), "999941057058");
});

test("mobile PHR calls use the PHR key and production PHR host", async () => {
  const { calls, fetch } = fakeNha();
  await client({ environment: "production", accessToken: "t", fetch }).mobile.requestOtp("9999999999");

  const [call] = calls;
  assert.equal(call.url, "https://phr.abdm.gov.in/api/phr/app/v3/enrollment/request/otp");
  assert.equal(decrypt(phrKeys, call.body.loginId), "9999999999");
});

test("per-call accessToken wins and errors carry the NHA body", async () => {
  const { calls, fetch } = fakeNha();
  const m1 = client({ environment: "production", accessToken: "client-token", fetch });

  const err = await m1.login.verifyOtp({ txnId: "t1", otp: "123456", otpSystem: "abdm" }, { accessToken: "mine" }).catch((e) => e);
  assert.ok(err instanceof AbdmError);
  assert.equal(err.status, 400);
  assert.deepEqual(err.body, [{ code: "ABDM-1006", message: "Invalid OTP Request" }]);

  const [call] = calls;
  assert.equal(call.url, "https://abha.abdm.gov.in/api/abha/v3/profile/login/verify");
  assert.equal(call.headers.Authorization, "Bearer mine");
  assert.deepEqual(call.body.scope, ["abha-login", "mobile-verify"]);
});

test("built-in keys load", () => {
  const m1 = new M1Client({ accessToken: "t" });
  assert.ok(m1.encrypt("1234").length > 100);
});

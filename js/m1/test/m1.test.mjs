import { test } from "node:test";
import assert from "node:assert/strict";
import { constants, generateKeyPairSync, privateDecrypt } from "node:crypto";
import { AbdmError, M1Client } from "../dist/index.js";

const { publicKey, privateKey } = generateKeyPairSync("rsa", { modulusLength: 2048 });
const decrypt = (b64) =>
  privateDecrypt({ key: privateKey, padding: constants.RSA_PKCS1_OAEP_PADDING, oaepHash: "sha1" }, Buffer.from(b64, "base64")).toString();

function fakeAbdm() {
  const calls = [];
  const json = (body, status = 200) => new Response(JSON.stringify(body), { status, headers: { "content-type": "application/json" } });
  const fetch = async (url, init) => {
    const path = new URL(url).pathname;
    calls.push({ path, headers: init.headers, body: init.body && JSON.parse(init.body) });
    if (path.endsWith("/profile/public/certificate")) return json({ publicKey: publicKey.export({ type: "spki", format: "der" }).toString("base64") });
    if (path.endsWith("/enrollment/request/otp")) return json({ txnId: "t1", message: "sent" });
    if (path.endsWith("/profile/login/verify")) return json([{ code: "ABDM-1006", message: "Invalid OTP Request" }], 400);
    return json({});
  };
  return { calls, fetch };
}

test("encrypts identifiers, fetches the key once and sends NHA headers", async () => {
  const { calls, fetch } = fakeAbdm();
  let n = 0;
  const m1 = new M1Client({ accessToken: () => `token-${++n}`, fetch });

  const [a, b] = await Promise.all([m1.aadhaar.requestOtp("999941057058"), m1.aadhaar.requestOtp("999941057058")]);
  assert.equal(a.txnId, "t1");
  assert.equal(b.txnId, "t1");
  assert.equal(calls.filter((c) => c.path.endsWith("/certificate")).length, 1);

  const otp = calls.find((c) => c.path.endsWith("/enrollment/request/otp"));
  assert.match(otp.headers.Authorization, /^Bearer token-\d$/);
  assert.match(otp.headers["REQUEST-ID"], /^[0-9a-f-]{36}$/);
  assert.ok(!Number.isNaN(Date.parse(otp.headers.TIMESTAMP)));
  assert.equal(decrypt(otp.body.loginId), "999941057058");
});

test("per-call accessToken wins and errors carry the NHA body", async () => {
  const { calls, fetch } = fakeAbdm();
  const m1 = new M1Client({ accessToken: "client-token", fetch });

  const err = await m1.login.verifyOtp({ txnId: "t1", otp: "123456", otpSystem: "abdm" }, { accessToken: "mine" }).catch((e) => e);
  assert.ok(err instanceof AbdmError);
  assert.equal(err.status, 400);
  assert.deepEqual(err.body, [{ code: "ABDM-1006", message: "Invalid OTP Request" }]);

  const verify = calls.find((c) => c.path.endsWith("/profile/login/verify"));
  assert.equal(verify.headers.Authorization, "Bearer mine");
  assert.deepEqual(verify.body.scope, ["abha-login", "mobile-verify"]);
  assert.equal(calls.find((c) => c.path.endsWith("/certificate")).headers.Authorization, "Bearer mine");
});

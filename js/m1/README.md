# @eka-care/abdm

ABDM M1 (ABHA) for Node.js 18+. It calls the NHA v3 APIs directly and has no runtime dependencies.

The client adds the `Authorization`, `REQUEST-ID` and `TIMESTAMP` headers and handles RSA encryption of Aadhaar numbers, mobile numbers, OTPs and passwords. Pass everything in plain text.

## Setup

```ts
import { M1Client } from "@eka-care/abdm";

const m1 = new M1Client({
  environment: "sandbox", // or "production"
  accessToken: () => getAbdmSessionToken(), // or a plain string
});
```

### Access token

The client never sees your client ID or secret. You get the ABDM gateway session token (`POST /api/hiecm/gateway/v3/sessions`) and pass it in:

- **A function** (recommended): it's called before every request, so your code can cache the token and renew it when it expires.
- **A string**: fine for scripts. Create a new client when the token expires.
- **Per call**: pass `{ accessToken }` as the last argument to override the token for that call. That argument also takes `requestId`, `signal` and `headers`.

### Keys and hosts

NHA's public keys are built in, the same ones Eka's ndhm service uses. There's one key for ABHA calls and a separate one for the PHR app calls (the mobile PHR address flow). The client makes no network call to fetch a key. If NHA rotates a key, pass the new PEM in `publicKeys: { abha, phr }`.

In production, NHA serves ABHA, PHR app and PHR web calls from different hosts. The client picks the right one for each call. To change any of them, for example to go through a proxy, use `baseUrls: { abha, phrApp, phrWeb }`.

User tokens (the X-token and T-token) are separate. They come back from verify calls, and you pass them to the methods that need them.

## Create an ABHA number with Aadhaar

```ts
const { txnId } = await m1.aadhaar.requestOtp("999941057058");
const enrol = await m1.aadhaar.verifyOtp({ txnId, otp: "123456", mobile: "9999999999" });
// enrol.isNew === false means an ABHA already existed for this Aadhaar.

// If enrol.ABHAProfile.mobile is not the number you passed, verify that number:
await m1.aadhaar.requestMobileOtp({ txnId, mobile: "9999999999" });
await m1.aadhaar.verifyMobileOtp({ txnId, otp: "654321" });

const { abhaAddressList } = await m1.aadhaar.suggestAddresses(txnId);
await m1.aadhaar.createAddress({ txnId, abhaAddress: "john_doe" });

const xToken = enrol.tokens.token;
const profile = await m1.profile.get(xToken); // profile.phrAddress lists the ABHA addresses
```

## Create an ABHA address with a mobile number (PHR)

```ts
const { txnId } = await m1.mobile.requestOtp("9999999999");
const { users } = await m1.mobile.verifyOtp({ txnId, otp: "123456" }); // addresses already on this mobile

const dob = { dayOfBirth: "01", monthOfBirth: "01", yearOfBirth: "1990" };
const { abhaAddressList } = await m1.mobile.suggestAddresses({ txnId, firstName: "John", lastName: "Doe", ...dob });
if (!(await m1.mobile.addressExists("john_doe@sbx"))) {
  await m1.mobile.createAddress({
    txnId, mobile: "9999999999", abhaAddress: "john_doe", firstName: "John", lastName: "Doe", gender: "M", ...dob,
  });
}
```

These calls use the `/v3/phr/app/enrollment` APIs, and NHA has to enable them for your client ID.

## Log in with Aadhaar, a mobile number or an ABHA number

```ts
const otp = await m1.login.requestOtp({ loginHint: "mobile", value: "9999999999" });
const res = await m1.login.verifyOtp({ txnId: otp.txnId, otp: "123456", otpSystem: otp.otpSystem });

// If several accounts share the mobile, pick one:
const { token: xToken } = await m1.login.verifyUser({
  txnId: res.txnId, abhaNumber: res.accounts[0].ABHANumber, tToken: res.token,
});
```

`loginHint` is `aadhaar`, `mobile` or `abha-number`. For an ABHA number, pass `otpSystem: "aadhaar"` to send the OTP to the mobile linked to Aadhaar.

## Verify an ABHA address and list its addresses

```ts
const { authMethods } = await m1.address.search("john_doe@sbx");
const { txnId } = await m1.address.requestOtp({ abhaAddress: "john_doe@sbx", otpSystem: "abdm" });
const { users, tokens } = await m1.address.verifyOtp({ txnId, otp: "123456", otpSystem: "abdm" });
const profile = await m1.address.getProfile(tokens!.token);
```

## Errors

A response that isn't 2xx throws `AbdmError` with `status`, the NHA response `body` and the `requestId` that was sent.

## Layout

```
src/
  client.ts     M1Client: connects each service to its NHA host and key
  common/       shared code: config and hosts, NHA keys, HTTP client, encryption, errors, OTP helpers, shared models
  aadhaar/      create an ABHA number with Aadhaar          (aadhaar.service.ts, aadhaar.types.ts)
  phr/          create an ABHA address from a mobile number (phr.service.ts, phr.types.ts)
  login/        log in with Aadhaar, mobile or ABHA number  (login.service.ts, login.types.ts)
  address/      verify an ABHA address, list its addresses  (address.service.ts, address.types.ts)
  profile/      ABHA profile and card                       (profile.service.ts, profile.types.ts)
```

Services depend on the `HttpClient` and `Encryptor` interfaces in `common/`, not on the concrete classes. Each service can be tested with a fake and used on its own.

## Develop

```sh
npm install
npm test
```

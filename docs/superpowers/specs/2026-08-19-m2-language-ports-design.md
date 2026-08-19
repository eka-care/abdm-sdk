# M2 care-context SDK — Python and Java ports

Date: 2026-08-19
Status: Approved for planning
Reference implementation: `go/` (shipped in PR #3)
API contract: `docs/abdm-api-reference.md`
Parent spec: `docs/superpowers/specs/2026-08-19-abdm-m2-carecontext-design.md`

## Goal

Port the M2 (HIP) care-context surface to Python and Java at full parity with Go: the same
six functions, the same behavioural tests, packaging, and a README.

## Scope

In scope, per language:

- The client core each SDK needs to stand alone: configuration, credential handling with
  runtime refresh, and an HTTP layer that resolves the bearer token per request.
- The six care-context functions.
- Tests, including the two that matter most (below).
- Packaging: `pyproject.toml` / `pom.xml`.
- A README following `go/services/abdm/carecontext/README.md`.

Out of scope:

- **M1 (ABHA identity)** — not ported. Go's `abha/login`, `registration`, `profile` have no
  Python or Java equivalent yet. Each README must say so plainly rather than implying parity.
- **JavaScript** — deliberately not built. `abdm-ecdh` publishes Go, Python and Java only, so a
  JS `respondToFetch` would require hand-rolling ECDH on Curve25519 in Weierstrass form,
  byte-compatible with Java/BouncyCastle's X.509 encoding, plus HKDF-SHA256 and AES-256-GCM.
  Getting that subtly wrong yields ciphertext Eka accepts with a 202 and the recipient cannot
  decrypt — silent, exactly like the checksum risk. Cross-language crypto belongs in
  `abdm-ecdh`, ported once and validated against its existing test vectors. `js/` stays a stub
  recording this reason.
- FHIR construction, persistence, dedupe, Providers/Records listing, Option A upload.

## Dependencies: keep them at one

The Go SDK's only dependency is `abdm-ecdh`. Hold that line.

- **Python**: `abdm-ecdh>=1.0.0` only. Use stdlib `urllib.request` for HTTP, `hmac`/`hashlib`
  for signatures, `json`, `unittest` for tests. No `requests`, no `httpx`, no `pytest`
  requirement (tests must run under plain `python -m unittest`).
- **Java**: `abdm-ecdh` (JitPack) plus Jackson, because Java has no standard JSON binding.
  HTTP is stdlib `java.net.http.HttpClient` (Java 11+). Tests use JUnit 5.

## Crypto APIs to call — never reimplement

Python:
```python
from abdm_ecdh import generate_key_material, encrypt
ours = generate_key_material()           # .private_key .public_key .x509_public_key .nonce
enc  = encrypt(string_to_encrypt=..., sender_nonce=ours.nonce, requester_nonce=hiu_nonce,
               sender_private_key=ours.private_key, requester_public_key=hiu_x509_public_key)
enc.encrypted_data
```

Java:
```java
KeyMaterial ours = AbdmEcdh.generateKeyMaterial();   // privateKey() publicKey() x509PublicKey() nonce()
EncryptionResponse enc = AbdmEcdh.encrypt(plaintext, ours.nonce(), hiuNonce,
                                          ours.privateKey(), hiuX509PublicKey);
enc.encryptedData();
```

Roles are easy to transpose and the failure is silent: **ours is the sender, the HIU is the
requester, and the encryption target is the HIU's X.509 public key.**

## Gotchas the Go implementation paid for

Every one of these caused a real defect or was caught in review. Carry them across verbatim.

1. **Event names are inconsistent upstream.** `abha.care_context_discover`,
   `abha.care_context_discover_link_init`, but `abha.context_discover_link_confirm` — no
   `care_` prefix on the third. Copy verbatim; never derive.
2. **`linkRefNumber` is camelCase** where every neighbouring field is snake_case.
3. **`on-link-confirm` sends `request_id` but no `txn_id`** — its webhook carries none.
4. **Reject an empty webhook secret.** An empty HMAC key is one anybody can compute, so a
   forgotten secret turns the endpoint into an open door for forged `hip_data_fetch` events.
5. **Verify the signature before parsing anything.** HMAC-SHA256 over the literal
   `"{t}.{rawBody}"` using the raw received bytes, constant-time comparison, and a ±3 minute
   window enforced in both directions.
6. **The checksum is hex-encoded MD5 over the PLAINTEXT bundle**, hashed before encryption.
   MD5 and hex are fixed by the wire protocol (confirmed with Eka), not a security choice —
   the comment must say so, or someone will "upgrade" it to SHA-256 and break every bundle.
   Isolate it in one named function and note that a wrong algorithm **fails silently**: Eka
   returns 202, the call returns success, and only the HIU discards the data. Mirror Go's
   `checksum()` in `go/services/abdm/carecontext/datafetch.go` exactly.
7. **Ephemeral key material per call.** Generate inside the responder; never cache it on the
   client or reuse it across calls. The private key must never leave the process.
8. **Never hold a lock across a token-refresh network call.** In Go this produced a reentrant
   self-deadlock: request → resolve token → provider takes its lock → refresh → request →
   resolve token → same lock, same thread. The auth/refresh client must NOT resolve tokens.
   Both ports must avoid the same shape; in Java note that `synchronized` IS reentrant so it
   deadlocks differently — on the nested HTTP call's own lock or not at all — which makes it
   easier to miss, not safer.
9. **When a result carries an error, omit the payload**, sending only correlation ids and the
   error. All three responders behave identically.
10. **An unrecognised event must be distinguishable from an unauthentic one**, because the
    caller answers 200 for the first and 401 for the second.
11. **Document that callbacks must be idempotent** and that a non-2xx makes ABDM retry.

## The two tests that matter

Everything else is ordinary coverage. These two catch silent failures:

1. **The HIU round-trip.** Generate a second key pair acting as the HIU, call the responder,
   then decrypt the pushed ciphertext using ONLY the key material the responder published in
   its request body. If key or nonce roles are transposed this fails; without it, the SDK ships
   bytes that look fine and cannot be read.
2. **Header assertions with distinct values.** Assert `X-Pt-Id`, `X-Partner-Pt-Id` and
   `X-Hip-Id` using three different values, so a transposition cannot pass by coincidence. Go
   shipped this gap and it was caught only in review.

Also required: signature cases (valid, tampered body, wrong secret, empty secret, stale
timestamp, malformed header) and one case per event name asserting it decodes to the right type.

## Naming

Idiomatic per language, same semantics: Python `link`, `parse_webhook`, `respond_to_fetch`,
`on_discover`, `on_link_init`, `on_link_confirm`; Java `link`, `parseWebhook`,
`respondToFetch`, `onDiscover`, `onLinkInit`, `onLinkConfirm`.

## Layout

```
python/
  pyproject.toml
  README.md
  src/eka_abdm/{__init__,client,config,credentials,http,carecontext/{__init__,types,service,webhook,datafetch,discovery}}.py
  tests/
java/
  pom.xml
  README.md
  src/main/java/care/eka/abdm/{EkaClient,Config,Credentials,HttpTransport}.java
  src/main/java/care/eka/abdm/carecontext/{CareContexts,Events,Types}.java
  src/test/java/care/eka/abdm/carecontext/
```

## Open questions

None. The checksum algorithm was confirmed with Eka on 2026-08-19 as hex-encoded MD5 and is
already correct in Go; both ports must match it.

# Eka Care ABDM SDKs

Client libraries for [ABDM](https://abdm.gov.in) (Ayushman Bharat Digital Mission)
via Eka Care's ABDM Connect APIs — one SDK per language, built against a shared,
generated API contract.

These are **libraries you import**, not servers you run. You own your HTTP
endpoints; the SDKs remove the boilerplate — request shapes, webhook signature
verification, ECDH encryption, and correlation — so you call functions instead of
hand-rolling REST calls.

| Language | Status | Install | Location |
|---|---|---|---|
| Go | M1 + M2 (HIP) | `go get github.com/eka-care/abdm-sdk/go` | [`go/`](go/) |
| Python | M2 (HIP) | `pip install eka-abdm` | [`python/`](python/) |
| Java | M2 (HIP) | Maven, via JitPack | [`java/`](java/) |
| JavaScript | blocked — see below | — | [`js/`](js/) |

**M1 (ABHA identity) is Go-only.** The Python and Java SDKs implement M2 and the client core
they need to stand alone; they have no ABHA login, registration or profile.

**JavaScript is deliberately not built.** [`abdm-ecdh`](https://github.com/eka-care/abdm-ecdh)
publishes Go, Python and Java only. A JS `respondToFetch` would mean hand-rolling ECDH on
Curve25519 in Weierstrass form, byte-compatible with Java/BouncyCastle's X.509 encoding, plus
HKDF-SHA256 and AES-256-GCM. Getting that subtly wrong produces ciphertext Eka accepts with a
202 and the recipient cannot decrypt — a silent failure. Cross-language crypto belongs in
`abdm-ecdh`, ported once and validated against its existing test vectors.

## The API contract

[`docs/abdm-api-reference.md`](docs/abdm-api-reference.md) is the contract every
SDK is built against: all 92 ABDM Connect endpoints and 18 webhook events, with
their headers, request/response fields and nested schemas, plus a per-section
coverage count.

It is **generated**, not hand-written. The published docs at
[developer.eka.care](https://developer.eka.care/api-reference/user-app/abdm-connect/overview)
are the source of truth:

```bash
python3 scripts/extract-api-docs.py
```

Re-run it and diff when the published docs change. Never hand-edit the generated file.

## Milestones

ABDM certification is split into four milestones. Each SDK documents which it covers.

- **M1 — ABHA identity**: creation (Aadhaar / mobile / face auth), login, profile, KYC, cards.
- **M2 — Care contexts (HIP)**: linking records to a patient's ABHA address, serving
  encrypted health data to requesting HIUs, and patient-initiated discovery.
- **M3 — Consent (HIU)**: consent lifecycle and fetching records from other providers.
- **M4 — HPR / HFR**: registering practitioners and facilities.

## Encryption

All health-data encryption uses [`eka-care/abdm-ecdh`](https://github.com/eka-care/abdm-ecdh)
(ECDH on Curve25519, HKDF-SHA256, AES-256-GCM), published for Go, Python and Java.
No SDK in this repo implements cryptography itself, and neither should your integration.

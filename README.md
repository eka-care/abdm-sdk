# Eka Care ABDM SDKs

Client libraries for [ABDM](https://abdm.gov.in) (Ayushman Bharat Digital Mission)
via Eka Care's ABDM Connect APIs — one SDK per language, built against a shared,
generated API contract.

These are **libraries you import**, not servers you run. You own your HTTP
endpoints; the SDKs remove the boilerplate — request shapes, webhook signature
verification, ECDH encryption, and correlation — so you call functions instead of
hand-rolling REST calls.

| Language | Status | Location |
|---|---|---|
| Go | M1 + M2 (HIP) | [`go/`](go/) |
| Python | planned | [`python/`](python/) |
| JavaScript | planned | [`js/`](js/) |
| Java | planned | [`java/`](java/) |

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

# eka-abdm — Python SDK

Eka Care ABDM Connect SDK for Python. Currently covers **Milestone 2 (HIP)
care-context**: linking care contexts, serving health data to HIUs, and
answering patient-initiated discovery, plus the client core (config,
credentials with runtime refresh, HTTP transport) it stands on.

**M1 (ABHA identity) is not ported.** Go's `abha/login`, `registration`, and
`profile` have no Python equivalent yet.

Install:

```bash
pip install eka-abdm
```

Runtime dependency: [`abdm-ecdh`](https://pypi.org/project/abdm-ecdh/) only.
Everything else — HTTP, HMAC, JSON — is Python stdlib. Tests run under plain
`python -m unittest`, no `pytest` required.

## Client setup

```python
from eka_abdm import Client, Environment

client = Client(
    environment=Environment.PRODUCTION,
    client_id="...",
    client_secret="...",
)
client.login()  # fails fast on bad credentials; installs a refreshing token resolver

cc = client.care_contexts
```

## Six functions; you write no HTTP calls and touch no cryptography

This is a library, not a server. You own your webhook endpoint and register
its URL with Eka; inside your handler you call `parse_webhook` and one
responder.

```python
from eka_abdm.carecontext import (
    CareContext, DataFetchEvent, DiscoverEvent, DiscoverResult,
    Headers, LinkConfirmEvent, LinkConfirmResult, LinkInitEvent,
    LinkInitResult, LinkRequest, LinkStatusEvent, parse_webhook,
)

# Outbound, when a visit closes:
cc.link(headers, LinkRequest(abha_address="...", care_contexts=[...]))

# Inbound, in your own route:
event = parse_webhook(body, request.headers["Eka-Webhook-Signature"], secret)

if isinstance(event, DataFetchEvent):
    cc.respond_to_fetch(event, entries)
elif isinstance(event, LinkStatusEvent):
    ...  # record the outcome
elif isinstance(event, DiscoverEvent):
    cc.on_discover(event, result)
elif isinstance(event, LinkInitEvent):
    cc.on_link_init(event, result)
elif isinstance(event, LinkConfirmEvent):
    cc.on_link_confirm(event, result)
```

`headers` is an `eka_abdm.carecontext.Headers` (alias of `eka_abdm.http.Headers`)
— the ABDM request identifiers (`X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id`) sent
with every call:

```python
headers = Headers(patient_id=oid, partner_user_id=your_patient_id, hip_id=hip_id)
```

## Verifying webhooks

`parse_webhook(body, signature_header, secret)` verifies before it parses
anything, then decodes into one of five typed events. It raises three
distinguishable exceptions, because your handler must answer differently:

| Exception | Meaning | Respond |
|---|---|---|
| `BadSignatureError` | not signed by Eka, altered, malformed header, or an empty secret | 401 |
| `StaleTimestampError` | valid signature, timestamp outside the ±3 minute window | 401 |
| `UnknownEventError` | authentic payload, event this package does not decode | 200 |

An empty webhook secret is rejected outright — an empty HMAC key is one
anyone can compute, so a forgotten secret would otherwise turn your endpoint
into an open door for forged events.

## A checksum mismatch is silent

`respond_to_fetch` returns normally as soon as Eka answers 202. If a bundle's
checksum does not match what the receiving HIU computes, the HIU discards it
and nothing reports that back — the records look shared but never arrive.
The checksum is **hex-encoded MD5 of the plaintext bundle**, hashed before
encryption. That is fixed by the wire protocol (confirmed with Eka), not a
security choice — do not "upgrade" it to SHA-256.

## Two rules you must follow

1. **Return non-2xx to make ABDM retry.** Retry is the recovery mechanism;
   this package has no queue and no dead-letter handling.
2. **Make your callbacks idempotent.** ABDM re-delivers. Deduplicate on the
   event's transaction or request id before doing anything with side
   effects.

## What stays yours

The package is boilerplate. It never decides what your data means:

- **FHIR bundles.** You supply bytes; `respond_to_fetch` encrypts and pushes
  them.
- **Patient matching** during discovery. Match conservatively — a loose
  match discloses another patient's records to whoever asked.
- **OTP generation, delivery and validation**, and the storage that keeps an
  OTP alive between the link-init and link-confirm webhooks. Store it
  against the `ref_num` you return from `on_link_init`; it comes back as
  `LinkConfirmEvent.link_ref_number`.
- **Link state.** `link` returns after a 202, not "linked". The outcome
  arrives later as a `LinkStatusEvent`.

## Encryption

Delegated entirely to [`abdm-ecdh`](https://github.com/eka-care/abdm-ecdh)
(Curve25519 ECDH, HKDF-SHA256, AES-256-GCM) — never reimplemented here.
`respond_to_fetch` generates an ephemeral key pair per call and discards it,
so there is no key store to manage, nothing to rotate, and the private key
never leaves the process.

## Event names, verbatim

The upstream event names are inconsistent — copy them exactly, never derive:

- `abha.hip_data_fetch`
- `abha.link_care_context`
- `abha.care_context_discover`
- `abha.care_context_discover_link_init`
- `abha.context_discover_link_confirm` — **no** `care_` prefix, unlike its
  neighbours.

## Development

```bash
python3 -m venv .venv
.venv/bin/pip install -e .
.venv/bin/python -m unittest discover -s tests -v
```

Design notes that apply to every language port live in
[`../docs/superpowers/specs/`](../docs/superpowers/specs/). The Go
implementation ([`../go/services/abdm/carecontext/`](../go/services/abdm/carecontext/))
is the reference this package mirrors in behaviour and error semantics.

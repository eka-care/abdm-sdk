# ABDM M2 (HIP) — `carecontext` package design

Date: 2026-08-19
Status: Approved for planning
Repo: `github.com/eka-care/eka-sdk-go` (Go first; spec is language-neutral for the pip/Maven ports)

## Problem

M1, M3 and M4 ship UI SDKs, so integrators get a drop-in flow. M2 is backend-to-backend:
there is no UI to own, so today every integrator hand-rolls the whole thing against raw REST —
webhook signature verification, ECDH key exchange, AES-GCM encryption, checksums, and the
Data-On-Fetch correlation.

Wrapping the REST calls in typed methods is not the deliverable; that is the easy part. The
deliverable is removing the crypto and correlation work from the integrator entirely.

## Scope

In scope (v1):

- Outbound care-context linking (`POST /abdm/v1/care-contexts/link`).
- Inbound `abha.hip_data_fetch`: parse, encrypt, respond via Data-On-Fetch.
- Inbound `abha.link_care_context`: parse into a typed status event.
- Per-request credential resolution so long-running integrator processes do not 401.

Out of scope (v1):

- Discovery (`abha.discover_care_context`) and the link-init/confirm OTP handshake.
  **Full M2 certification requires these**; v1 covers the HIP-initiated path end to end.
- FHIR bundle construction — the integrator supplies bundle bytes (see `abdm-fhir`).
- Persistence, link-state tables, dedupe stores, retry queues.
- Providers/Records listing endpoints.
- Option A (`data` field upload, Eka-hosted storage) — not enabled for most accounts.

## Constraints that shaped the design

1. **No routers or handlers.** This is a library integrators import (`go get`, `pip install`).
   They own their HTTP endpoint and register the webhook URL with Eka themselves. Our code is
   called from inside their handler. This is also what makes the design portable — an
   `http.Handler` has no Python analogue; `respond_to_fetch(event, entries)` does.
2. **Crypto is a solved dependency.** Use `github.com/eka-care/abdm-ecdh/go@v1.0.0`
   (also published for Python and Java). Do not reimplement ECDH/HKDF/AES-GCM.
3. **Credentials belong to the integrator.** They store them; we accept them and refresh
   in memory at runtime.

## Key protocol facts (verified against developer.eka.care)

- Webhook signature: header `Eka-Webhook-Signature`, value `t=<unix>,v1=<hex>`.
  HMAC-SHA256 over the literal string `"{t}.{rawBody}"` using the subscription's shared
  secret. Constant-time compare. Reject when `|now - t| > 3m`.
- `abha.hip_data_fetch` payload carries the **HIU's** `key_information`
  (`crypto_alg`, `curve`, `dh_public_key{key_value, parameters, expiry}`, `nonce`),
  plus `care_contexts[]`, `hi_types[]`, `abha_address`, `oid`, `partner_patient_id`, `hip_id`.
- Data-On-Fetch: `POST /abdm/v1/hip/care-context/data/on-fetch`, correlated by
  `transaction_id`, carrying `entries[]` (`care_context_id`, `content`, `checksum`, `media`),
  `page_number`, `page_count`, and the **HIP's own** `key_information`.

Consequence: because both parties send their own key material in-band, there is no key
lifecycle. Generate an ephemeral keypair per `RespondToFetch` call and discard it. No key
store, no rotation policy, no registration step.

## Package layout

```
services/abdm/carecontext/
  types.go      LinkRequest/Response, CareContext, HIType constants, DataFetchEvent, LinkStatusEvent, Entry
  service.go    Link(), RespondToFetch()
  webhook.go    ParseWebhook()
  webhook_test.go
examples/m2-hip/main.go
```

Mirrors the existing `services/abdm/abha/*` convention (`Service` struct built from
`interfaces.Config`, thin methods over `internal/http`). Registered as
`client.ABDM.CareContexts()` in `services/abdm/client.go`, alongside `Login()`,
`Registration()`, `Profile()`.

## Exported surface — three functions

```go
cc := client.ABDM.CareContexts()

// 1. Outbound, on the integrator's "visit closed" / "report finalized" event.
resp, err := cc.Link(ctx,
    interfaces.Headers{PatientID: oid, PartnerUserID: ptID, HipID: hipID},
    &carecontext.LinkRequest{
        ABHAAddress: "patient@sbx",
        CareContexts: []carecontext.CareContext{{
            CareContextID: "visit-123",
            Display:       "OP Consult 19 Aug",
            HIType:        carecontext.HITypeOPConsultation,
        }},
    })
// 202 Accepted; the real result arrives later as abha.link_care_context.

// 2 + 3. Inside the integrator's own route — their framework, their status codes.
evt, err := carecontext.ParseWebhook(rawBody, r.Header.Get("Eka-Webhook-Signature"), secret)
switch e := evt.(type) {
case *carecontext.DataFetchEvent:
    err = cc.RespondToFetch(ctx, e, []carecontext.Entry{
        {CareContextID: "visit-123", Bundle: fhirJSON},
    })
case *carecontext.LinkStatusEvent:
    err = store.MarkLinked(e.CareContextID, e.Status, e.Error)
}
```

`Entry` is `{CareContextID string; Bundle []byte}` — ABDM-compliant FHIR R4 JSON.

`ParseWebhook` returns a typed event, or `ErrBadSignature` / `ErrStaleTimestamp` /
`ErrUnknownEvent`. Unknown events are a distinct sentinel, not a failure: the integrator
should answer 200 so their endpoint does not break when new event types ship.

`RespondToFetch` takes the event value straight back. That event carries `transaction_id`,
the HIU key material, and the `oid` / `partner_patient_id` / `hip_id` needed for the response
headers, so the integrator never rebuilds a correlation or touches key material — the HIU's
`key_information` is captured on an unexported field and is not reachable from outside.

### What `RespondToFetch` does internally

1. `abdmecdh.New().GenerateKeyMaterial()` — ephemeral, per call.
2. For each entry: `Encrypt(EncryptionRequest{StringToEncrypt: string(Bundle),
   SenderNonce: ours, RequesterNonce: hiu.Nonce, SenderPrivateKey: ours,
   RequesterPublicKey: hiu.DHPublicKey.KeyValue})`.
3. SHA-256 checksum of the **plaintext** bundle (per API contract).
4. `POST /abdm/v1/hip/care-context/data/on-fetch` with `transaction_id`, `entries[]`,
   our `key_information`, and `X-Pt-Id` / `X-Partner-Pt-Id` / `X-Hip-Id` filled from the event.

Single page (`page_number: 1`, `page_count: 1`) in v1.
`ponytail: single-page response; add chunking if bundles exceed the payload limit.`

## Errors, retries, idempotency

We return errors; the integrator chooses the HTTP status, because it is their endpoint.
Two rules go in the README in plain language:

- Return non-2xx to make ABDM retry.
- Dedupe on `TransactionID`, because ABDM **will** retry and re-deliver.

No queue, no dedupe store, no dead-letter handling in the package. Deliberate: an in-package
in-memory dedupe is free for single-instance deploys and silently wrong for multi-replica
ones, which is the more dangerous failure.

## Credentials

The integrator stores credentials; we hold them in memory and refresh at runtime.

Refresh ladder (already implemented in `auth.ClientCredentialsProvider.Retrieve`):
cached and unexpired → return; expired with a valid refresh token → `RefreshToken`;
refresh fails or the refresh token is dead → `ClientLogin` for a fresh session.

Two entry points:

- `WithClientCredentials(id, secret)` — full ladder, can always re-mint.
- `WithAccessToken(tok)` — long-lived or integrator-supplied token. Same ladder where a
  refresh token is present. When both tokens are dead and no client credentials exist there
  is nothing to re-mint from, so fail with an explicit
  "credentials expired; supply client credentials or a fresh token" rather than retrying.

We do **not** call back on refresh. With client_id/secret we can always re-mint, so persisting
refreshed tokens across restarts buys nothing. An `OnTokenRefresh` hook is deferred until an
integrator asks for it.

### Blocking prerequisite: per-request token resolution

`internal/http.Client` snapshots `apiKey` at construction (`internal/http/http.go`, via
`config.GetAPIKey()`), and `Client.Login` sets the token once then rebuilds the ABDM client
(`client.go`). Correct enough for M1's request-scoped use; broken for M2, where a process
runs for days and calls Data-On-Fetch long after login — the token ages out and every push 401s.

Fix, at the one place all callers route through: `internal/http.Client` takes
`tokenFn func(context.Context) (string, error)` and calls it inside `Do()`. `Client.Login`
wires it to `provider.Retrieve`. Falls back to the existing config snapshot when unset, so
M1 services need no changes.

Two consequences worth having: it fixes long-running M1 usage by the same edit, and because
we stop writing `cfg.AuthorizationToken` on every refresh, the unsynchronised write to a
shared `Config` — a real data race once a concurrent webhook server exists — disappears
rather than needing a mutex.

Secondary fix: `StaticCredentialsProvider.Retrieve` currently returns its credentials
unconditionally, expired or not. It must apply the ladder above.

## Testing

`abdm-ecdh` ships its own crypto test suite, so we do not re-test encryption.

The only non-trivial logic we own is `ParseWebhook`. One table test: valid signature,
tampered body, stale `t`, unknown event type. No mocks, no fixtures, no HTTP test server.

`examples/m2-hip/main.go` is the runnable end-to-end reference: link a care context, then a
`net/http` route showing parse-and-respond.

## Language ports

Go first. The three-function surface (`link`, `parse_webhook`, `respond_to_fetch`) maps
directly onto Python and Java, and `abdm-ecdh` already publishes for all three. Each port is
its own project; this spec is the contract they implement.

## Open questions

None blocking. Deferred by explicit decision: discovery events, FHIR builder,
`OnTokenRefresh` hook, response pagination.

# carecontext — ABDM M2 (HIP)

Link care contexts, serve health data to HIUs, and answer patient-initiated
discovery. Six functions; you write no HTTP calls and touch no cryptography.

This is a library, not a server. You own your webhook endpoint and register its
URL with Eka; inside your handler you call `ParseWebhook` and one responder.

```go
cc := client.ABDM.CareContexts()

// Outbound, when a visit closes:
cc.Link(ctx, headers, &carecontext.LinkRequest{ /* ... */ })

// Inbound, in your own route:
evt, err := carecontext.ParseWebhook(body, r.Header.Get("Eka-Webhook-Signature"), secret)
switch e := evt.(type) {
case *carecontext.DataFetchEvent:   cc.RespondToFetch(ctx, e, entries)
case *carecontext.LinkStatusEvent:  /* record the outcome */
case *carecontext.DiscoverEvent:    cc.OnDiscover(ctx, e, result)
case *carecontext.LinkInitEvent:    cc.OnLinkInit(ctx, e, result)
case *carecontext.LinkConfirmEvent: cc.OnLinkConfirm(ctx, e, result)
}
```

See [`examples/m2-hip`](../../../examples/m2-hip) for a complete server.

`headers` is a `carecontext.Headers` — the ABDM request identifiers (`X-Pt-Id`,
`X-Partner-Pt-Id`, `X-Hip-Id`) sent with every call:

```go
headers := carecontext.Headers{PatientID: oid, PartnerUserID: yourPatientID, HipID: hipID}
```

## A checksum mismatch is silent

`RespondToFetch` returns `nil` as soon as Eka answers 202. If the bundle's
checksum does not match what the receiving HIU computes, the HIU discards it and
nothing reports that back — the records look shared but never arrive. The
algorithm is hex-encoded MD5 over the plaintext bundle, fixed by the wire
protocol and confirmed with Eka. It is an interoperability requirement, not a
security choice: do not "upgrade" it to SHA-256, or every bundle you send will
be silently rejected.

## Two rules you must follow

1. **Return non-2xx to make ABDM retry.** Retry is the recovery mechanism; this
   package has no queue and no dead-letter handling.
2. **Make your callbacks idempotent.** ABDM re-delivers. Deduplicate on the
   event's transaction or request id before doing anything with side effects.

## What stays yours

The package is boilerplate. It never decides what your data means:

- **FHIR bundles.** You supply bytes; we encrypt and push them.
- **Patient matching** during discovery. Match conservatively — a loose match
  discloses another patient's records to whoever asked.
- **OTP generation, delivery and validation**, and the storage that keeps an OTP
  alive between the link-init and link-confirm webhooks. Store it against the
  `RefNum` you return from `OnLinkInit`; it comes back as
  `LinkConfirmEvent.LinkRefNumber`.
- **Link state.** `Link` returns 202, not "linked". The outcome arrives later as
  a `LinkStatusEvent`.

## Encryption

Delegated entirely to [`abdm-ecdh`](https://github.com/eka-care/abdm-ecdh)
(Curve25519 ECDH, HKDF-SHA256, AES-256-GCM). `RespondToFetch` generates an
ephemeral key pair per call and discards it, so there is no key store to manage
and nothing to rotate.

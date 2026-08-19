# ABDM SDK — Java

Link care contexts, serve health data to HIUs, and answer patient-initiated
discovery. Six functions; you write no HTTP calls and touch no cryptography.

This is a library, not a server. You own your webhook endpoint and register its
URL with Eka; inside your handler you call `Webhook.parseWebhook` and one
responder.

```java
EkaClient client = new EkaClient(Config.PRODUCTION, clientId, clientSecret);
client.login();
CareContexts cc = client.careContexts();

// Outbound, when a visit closes:
cc.link(headers, new LinkRequest("patient@sbx", List.of(new CareContext("cc-1", "Visit"))));

// Inbound, in your own route:
Event event = Webhook.parseWebhook(rawBody, request.getHeader("Eka-Webhook-Signature"), secret);
if (event instanceof DataFetchEvent e) {
    cc.respondToFetch(e, entries);
} else if (event instanceof LinkStatusEvent e) {
    // record the outcome
} else if (event instanceof DiscoverEvent e) {
    cc.onDiscover(e, result);
} else if (event instanceof LinkInitEvent e) {
    cc.onLinkInit(e, result);
} else if (event instanceof LinkConfirmEvent e) {
    cc.onLinkConfirm(e, result);
}
```

`headers` is a `care.eka.abdm.Headers` — the ABDM request identifiers
(`X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id`) sent with every call:

```java
Headers headers = new Headers(oid, partnerUserId, hipId);
```

## Not ported: M1 (ABHA identity)

This SDK covers M2 (care-context linking, HIP data sharing and discovery)
only. Go's `abha/login`, `registration` and `profile` packages — ABHA
creation, login and profile management — have no Java equivalent here. Do
not assume parity with the Go SDK beyond what is listed above.

## A checksum mismatch is silent

`respondToFetch` returns normally as soon as Eka answers 202. If the
bundle's checksum does not match what the receiving HIU computes, the HIU
discards it and nothing reports that back — the records look shared but
never arrive. The algorithm is hex-encoded MD5 over the plaintext bundle,
fixed by the wire protocol and confirmed with Eka. It is an
interoperability requirement, not a security choice: do not "upgrade" it to
SHA-256, or every bundle you send will be silently rejected.

## Two rules you must follow

1. **Return non-2xx to make ABDM retry.** Retry is the recovery mechanism;
   this package has no queue and no dead-letter handling.
2. **Make your callbacks idempotent.** ABDM re-delivers. Deduplicate on the
   event's transaction or request id before doing anything with side
   effects.

## What stays yours

The package is boilerplate. It never decides what your data means:

- **FHIR bundles.** You supply bytes; we encrypt and push them.
- **Patient matching** during discovery. Match conservatively — a loose
  match discloses another patient's records to whoever asked.
- **OTP generation, delivery and validation**, and the storage that keeps an
  OTP alive between the link-init and link-confirm webhooks. Store it
  against the `refNum` you return from `onLinkInit`; it comes back as
  `LinkConfirmEvent.linkRefNumber()`.
- **Link state.** `link` returns once Eka answers 202, not "linked". The
  outcome arrives later as a `LinkStatusEvent`.

## Encryption

Delegated entirely to [`abdm-ecdh`](https://github.com/eka-care/abdm-ecdh)
(Curve25519 ECDH, HKDF-SHA256, AES-256-GCM), pulled in via JitPack as
`com.github.eka-care:abdm-ecdh`. `respondToFetch` generates an ephemeral key
pair per call and discards it, so there is no key store to manage and
nothing to rotate — the private key never leaves the process.

## Building and testing

```
cd java
mvn -q -B test
```

Tests are hermetic: they start a local `com.sun.net.httpserver.HttpServer`
on an ephemeral port and never touch the network.

Design notes that apply to every language port live in
[`../docs/superpowers/specs/`](../docs/superpowers/specs/). The API contract
both this SDK and Go are built against is
[`../docs/abdm-api-reference.md`](../docs/abdm-api-reference.md).

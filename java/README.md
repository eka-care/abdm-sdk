# ABDM SDK — Java

Not yet implemented.

Planned to mirror the Go SDK's surface, which is documented in [`../go/`](../go/).
The API contract both are built against is [`../docs/abdm-api-reference.md`](../docs/abdm-api-reference.md).

When built, encryption must delegate to [`abdm-ecdh`](https://github.com/eka-care/abdm-ecdh)
(the JitPack artifact `com.github.eka-care:abdm-ecdh`) rather than being reimplemented.

Design notes that apply to every language port live in
[`../docs/superpowers/specs/`](../docs/superpowers/specs/).

"""Answering abha.hip_data_fetch. Mirrors
go/services/abdm/carecontext/datafetch.go.

Encryption is delegated entirely to abdm-ecdh (Curve25519 ECDH, HKDF-SHA256,
AES-256-GCM). We never implement crypto ourselves.
"""

import hashlib
import time

from abdm_ecdh import encrypt, generate_key_material

from ..http import HTTPClient, Headers

_FHIR_MEDIA = "application/fhir+json"


def _checksum(plaintext: bytes) -> str:
    """Returns the hex-encoded MD5 of the plaintext bundle, hashed before
    encryption as the ABDM contract requires.

    MD5 and hex encoding are fixed by the wire protocol, confirmed with Eka.
    This is an interoperability requirement, not a security choice: the HIU
    recomputes this digest over the decrypted bundle and discards the
    records if it disagrees. Do not "upgrade" it to SHA-256 — the receiver
    would reject every bundle, and it would fail silently, because Eka
    accepts the push with 202 regardless and respond_to_fetch returns
    normally. Confidentiality comes from the AES-256-GCM encryption of the
    content, not from this field.
    """
    return hashlib.md5(plaintext).hexdigest()  # noqa: S303 -- protocol-mandated digest, not a security primitive


def _our_key_information(ours, theirs: dict) -> dict:
    """Describes the ephemeral key material we generated for one response,
    in the shape the HIU expects.

    CONFIRM: the docs do not state what a HIP should send for
    dh_public_key parameters and expiry. We echo the HIU's parameters and
    set a 24h expiry.
    """
    their_dh = theirs.get("dh_public_key") or {}
    params = their_dh.get("parameters") or "Curve25519/32byte random key"
    curve = theirs.get("curve") or "Curve25519"
    expiry = time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(time.time() + 24 * 60 * 60))
    return {
        "crypto_alg": "ECDH",
        "curve": curve,
        "nonce": ours.nonce,
        "dh_public_key": {
            "key_value": ours.x509_public_key,
            "parameters": params,
            "expiry": expiry,
        },
    }


def respond_to_fetch(config, event, entries) -> None:
    """Encrypts the given FHIR bundles for the requesting HIU and pushes
    them, answering an abha.hip_data_fetch webhook.

    Pass the event exactly as parse_webhook returned it: it carries the
    transaction id, the HIU's key material, and the identifiers used as
    request headers. A fresh ephemeral key pair is generated per call and
    discarded, so there is no key store to manage and the private key never
    leaves this process.

    ABDM retries a fetch it considers unanswered, so this may be called
    more than once for one transaction. Make your own bundle lookup
    idempotent.
    """
    if event is None:
        raise ValueError("carecontext: nil data fetch event")
    if not entries:
        raise ValueError(
            "carecontext: no entries to send for transaction %r" % event.transaction_id
        )

    ours = generate_key_material()
    key_info = event._key_info  # noqa: SLF001 -- same module family as webhook.py

    out = []
    for entry in entries:
        bundle = entry.bundle
        if not bundle:
            raise ValueError(
                "carecontext: empty bundle for care context %r" % entry.care_context_id
            )
        if isinstance(bundle, str):
            bundle = bundle.encode("utf-8")

        their_dh = key_info.get("dh_public_key") or {}
        enc = encrypt(
            string_to_encrypt=bundle.decode("utf-8"),
            sender_nonce=ours.nonce,
            requester_nonce=key_info.get("nonce", ""),
            sender_private_key=ours.private_key,
            requester_public_key=their_dh.get("key_value", ""),
        )
        out.append(
            {
                "care_context_id": entry.care_context_id,
                "content": enc.encrypted_data,
                "checksum": _checksum(bundle),
                "media": _FHIR_MEDIA,
            }
        )

    # ponytail: single page. Add chunking if bundles ever exceed the payload limit.
    http = HTTPClient(config, resolve_token=True)
    http.do(
        "POST",
        "/abdm/v1/hip/care-context/data/on-fetch",
        headers=Headers(
            patient_id=event.oid,
            partner_user_id=event.partner_patient_id,
            hip_id=event.hip_id,
        ),
        body={
            "transaction_id": event.transaction_id,
            "page_number": 1,
            "page_count": 1,
            "key_information": _our_key_information(ours, key_info),
            "entries": out,
        },
    )

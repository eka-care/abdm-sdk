"""ABDM Milestone 2 for a HIP: linking care contexts, serving health data on
demand, and answering patient-initiated discovery.

The package is boilerplate only. It carries bytes and correlation
identifiers; it never decides what your data means. Matching patients,
generating OTPs and building FHIR bundles stay in your code.

It is a library, not a server: you own your webhook endpoint and register
its URL with Eka. Inside your handler, call parse_webhook and then the
responder for the event you received.

One asymmetry to know about: a bundle whose checksum does not match what the
HIU computes is discarded at the far end, and nothing reports that back. Eka
answers 202 and respond_to_fetch returns normally either way, so a checksum
mismatch is invisible from the HIP side — records appear shared but never
arrive.
"""

from .discovery import DiscoverResult, LinkConfirmResult, LinkInitResult
from .service import CareContextService
from .types import (
    CareContext,
    DiscoveredCareContext,
    Entry,
    ErrorDetail,
    Headers,
    HIType,
    LinkRequest,
    Patient,
)
from .webhook import (
    EVENT_DATA_FETCH,
    EVENT_DISCOVER,
    EVENT_LINK_CONFIRM,
    EVENT_LINK_INIT,
    EVENT_LINK_STATUS,
    BadSignatureError,
    DataFetchEvent,
    DiscoverEvent,
    Identifier,
    LinkConfirmEvent,
    LinkInitEvent,
    LinkStatusEvent,
    StaleTimestampError,
    UnknownEventError,
    parse_webhook,
)

__all__ = [
    "CareContextService",
    "LinkRequest",
    "CareContext",
    "HIType",
    "Headers",
    "Entry",
    "ErrorDetail",
    "DiscoveredCareContext",
    "Patient",
    "DiscoverResult",
    "LinkInitResult",
    "LinkConfirmResult",
    "EVENT_DATA_FETCH",
    "EVENT_LINK_STATUS",
    "EVENT_DISCOVER",
    "EVENT_LINK_INIT",
    "EVENT_LINK_CONFIRM",
    "BadSignatureError",
    "StaleTimestampError",
    "UnknownEventError",
    "DataFetchEvent",
    "LinkStatusEvent",
    "Identifier",
    "DiscoverEvent",
    "LinkInitEvent",
    "LinkConfirmEvent",
    "parse_webhook",
]

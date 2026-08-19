"""Shared care-context types. Mirrors go/services/abdm/carecontext/types.go.

The package is boilerplate only. It carries bytes and correlation
identifiers; it never decides what your data means. Matching patients,
generating OTPs and building FHIR bundles stay in your code.
"""

from dataclasses import dataclass, field
from typing import List, Optional

from ..http import Headers

__all__ = [
    "Headers",
    "HIType",
    "CareContext",
    "LinkRequest",
    "Entry",
    "ErrorDetail",
    "DiscoveredCareContext",
    "Patient",
]


class HIType:
    """An ABDM health information type."""

    OP_CONSULTATION = "OPConsultation"
    PRESCRIPTION = "Prescription"
    DISCHARGE_SUMMARY = "DischargeSummary"
    DIAGNOSTIC_REPORT = "DiagnosticReport"
    IMMUNIZATION_RECORD = "ImmunizationRecord"
    HEALTH_DOCUMENT_RECORD = "HealthDocumentRecord"
    WELLNESS_RECORD = "WellnessRecord"


@dataclass
class CareContext:
    """One logical group of records — a visit, a test, a document."""

    care_context_id: str
    display: str
    hi_type: str = ""
    hi_types: Optional[List[str]] = None

    def to_dict(self) -> dict:
        d = {"care_context_id": self.care_context_id, "display": self.display}
        if self.hi_type:
            d["hi_type"] = self.hi_type
        if self.hi_types:
            d["hi_types"] = list(self.hi_types)
        return d


@dataclass
class LinkRequest:
    """The body of POST /abdm/v1/care-contexts/link.

    oid and partner_user_id are filled from the request headers by
    CareContextService.link when left empty (the API requires them in both
    places) — see service.py.
    """

    abha_address: str
    care_contexts: List[CareContext] = field(default_factory=list)
    oid: str = ""
    partner_user_id: str = ""

    def to_dict(self) -> dict:
        d = {
            "abha_address": self.abha_address,
            "care_contexts": [cc.to_dict() for cc in self.care_contexts],
        }
        if self.oid:
            d["oid"] = self.oid
        if self.partner_user_id:
            d["partner_user_id"] = self.partner_user_id
        return d


@dataclass
class Entry:
    """One care context's FHIR bundle, supplied by you and encrypted by
    Service.respond_to_fetch. bundle is ABDM-compliant FHIR R4 JSON bytes."""

    care_context_id: str
    bundle: bytes


@dataclass
class ErrorDetail:
    """Reports a failure back to ABDM through an on-* responder — no
    patient matched, OTP invalid, and so on."""

    code: int
    message: str

    def to_dict(self) -> dict:
        return {"code": self.code, "message": self.message}


@dataclass
class DiscoveredCareContext:
    """One unlinked care context offered during discovery."""

    ref_num: str
    display: str

    def to_dict(self) -> dict:
        return {"ref_num": self.ref_num, "display": self.display}


@dataclass
class Patient:
    """A matched patient and their unlinked care contexts."""

    ref_num: str
    display: str
    hi_type: str = ""
    care_contexts: List[DiscoveredCareContext] = field(default_factory=list)

    def to_dict(self) -> dict:
        d = {
            "ref_num": self.ref_num,
            "display": self.display,
            "care_contexts": [cc.to_dict() for cc in self.care_contexts],
        }
        if self.hi_type:
            d["hi_type"] = self.hi_type
        return d

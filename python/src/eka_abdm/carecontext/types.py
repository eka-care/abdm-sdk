"""Request types for the care-context link call. Mirrors
go/services/abdm/carecontext/types.go, scoped to what Task 1 needs (link
only — webhook/discovery types are Task 2).
"""

from dataclasses import dataclass, field
from typing import List, Optional

from ..http import Headers

__all__ = ["Headers", "HIType", "CareContext", "LinkRequest"]


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

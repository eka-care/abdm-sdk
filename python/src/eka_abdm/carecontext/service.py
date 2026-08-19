"""Care-context API calls. Mirrors go/services/abdm/carecontext/service.go,
scoped to link only for Task 1.
"""

import dataclasses

from ..http import HTTPClient
from .types import Headers, LinkRequest


class CareContextService:
    def __init__(self, config):
        self._http = HTTPClient(config, resolve_token=True)

    def link(self, headers: Headers, request: LinkRequest) -> None:
        """Links care contexts to a patient's ABHA address.

        The API is asynchronous: a 202 means accepted, not linked. Does not
        mutate the caller's request object — oid/partner_user_id are filled
        from headers on a copy when left empty, because the contract
        requires them in both the headers and the body.
        """
        body = dataclasses.replace(request)
        if not body.oid:
            body.oid = headers.patient_id
        if not body.partner_user_id:
            body.partner_user_id = headers.partner_user_id

        self._http.do(
            "POST",
            "/abdm/v1/care-contexts/link",
            headers=headers,
            body=body.to_dict(),
        )

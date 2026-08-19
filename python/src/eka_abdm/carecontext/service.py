"""Care-context API calls. Mirrors go/services/abdm/carecontext/service.go.

The outbound call (link) and the three responders (respond_to_fetch,
on_discover, on_link_init, on_link_confirm) all live on this one Service, the
same as Go. Parsing an inbound webhook is a free function — see webhook.py —
because it needs no configuration or transport.
"""

import dataclasses

from ..http import HTTPClient
from . import datafetch, discovery
from .types import Headers, LinkRequest


class CareContextService:
    def __init__(self, config):
        self._config = config
        self._http = HTTPClient(config, resolve_token=True)

    def link(self, headers: Headers, request: LinkRequest) -> None:
        """Links care contexts to a patient's ABHA address.

        The API is asynchronous: a 202 means accepted, not linked. Does not
        mutate the caller's request object — oid/partner_user_id are filled
        from headers on a copy when left empty, because the contract
        requires them in both the headers and the body. The outcome arrives
        later as an abha.link_care_context webhook, which parse_webhook
        decodes into a LinkStatusEvent.
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

    def respond_to_fetch(self, event, entries) -> None:
        """Encrypts the given FHIR bundles for the requesting HIU and pushes
        them, answering an abha.hip_data_fetch webhook. See datafetch.py."""
        datafetch.respond_to_fetch(self._config, event, entries)

    def on_discover(self, event, result: discovery.DiscoverResult) -> None:
        """Answers an abha.care_context_discover webhook. See discovery.py."""
        discovery.on_discover(self._config, event, result)

    def on_link_init(self, event, result: discovery.LinkInitResult) -> None:
        """Answers an abha.care_context_discover_link_init webhook. See
        discovery.py."""
        discovery.on_link_init(self._config, event, result)

    def on_link_confirm(self, event, result: discovery.LinkConfirmResult) -> None:
        """Answers an abha.context_discover_link_confirm webhook. See
        discovery.py."""
        discovery.on_link_confirm(self._config, event, result)

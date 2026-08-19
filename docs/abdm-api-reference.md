# ABDM Connect API reference (extracted)

Generated 2026-08-19 by `scripts/extract-api-docs.py` from all
92 endpoints across 141 pages under
`developer.eka.care/api-reference/user-app/abdm-connect`, enumerated via `llms.txt`.

The published docs are the source of truth. Regenerate this file rather than
hand-editing it.

`*` = required. `?` = nullable. `✅` = already implemented in this SDK.

Coverage: **27 of 92** endpoints implemented.

| Section | Endpoints | Implemented |
|---|---:|---:|
| blood-bank | 2 | 0 |
| care-contexts | 3 | 1 |
| care-contexts/discover | 6 | 0 |
| care-contexts/link | 1 | 1 |
| care-contexts/providers | 1 | 0 |
| care-contexts/records | 2 | 0 |
| commons | 3 | 2 |
| consents | 6 | 0 |
| consents/auto-approval | 2 | 0 |
| enrollment/aadhaar | 7 | 6 |
| enrollment/face-auth | 3 | 0 |
| enrollment/mobile | 4 | 4 |
| login | 4 | 3 |
| nhpr-abdm | 1 | 0 |
| nhpr-abdm/hfr | 3 | 0 |
| nhpr-abdm/hpr | 9 | 0 |
| patient-requests | 2 | 0 |
| phys-cons | 11 | 0 |
| profile/cards | 2 | 2 |
| profile/details | 3 | 3 |
| profile/kyc | 3 | 3 |
| profile/search | 1 | 0 |
| providers | 1 | 0 |
| providers/search | 1 | 0 |
| scan-and-pay | 5 | 0 |
| scan-and-share | 3 | 0 |
| session | 3 | 2 |

## blood-bank

### `POST /abdm/uhi/v1/blood-bank/search`

Search Blood Banks

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `blood_component*` string, `blood_group*` string, `district` string, `location` NdhmeventLocation, `radius` ModelRadiusAround, `search_type*` string, `state` string
- response: `request_id` string, `status` string

<details><summary>schemas</summary>

- **ModelRadiusAround** — `unit` string, `value` integer
- **ModelSearchBloodBankPayload** — `blood_component*` string, `blood_group*` string, `district` string, `location` NdhmeventLocation, `radius` ModelRadiusAround, `search_type*` string, `state` string
- **NdhmeventLocation** — `lat` number, `long` number

</details>

### `GET /abdm/uhi/v1/blood-bank/search/{request_id}`

Get Blood Banks

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `request_id` · success: `200`
- request: _none_
- response: `data` ModelBloodBankEkaResponse, `request_id` string, `status` string

<details><summary>schemas</summary>

- **ModelBloodBankEkaResponse** — `context` ModelContext, `error_code` ModelErrorCode, `message` ModelUHIMessage
- **ModelBloodBankStateWiseResponse** — `data` ModelBloodBankEkaResponse, `request_id` string, `status` string
- **ModelCatalogDescriptor** — `descriptor` ModelUHIDescriptor, `providers` []ModelProviders?
- **ModelContact** — `email` string, `phone` string
- **ModelContext** — `action` string, `city` string, `consumer_id` string, `consumer_uri` string, `core_version` string, `country` string, `domain` string, `message_id` string, `provider_id` string, `provider_uri` string, `timestamp` string, `transaction_id` string
- **ModelDescriptor** — `code` string, `name` string
- **ModelErrorCode** — `code` string, `message` string
- **ModelMeasure** — `computed_value` integer, `eliminated_value` integer, `unit` string, `value` integer
- **ModelProviders** — `categories` []ModelUHICategory?, `contact` ModelContact, `descriptor` ModelUHIDescriptor, `fulfillments` []ModelUHIFulfillments?, `id` string, `items` []ModelUHIItem?, `location` ModelUHILocation
- **ModelQuantity** — `count` integer, `measure` ModelMeasure
- **ModelTimeSeries** — `time` ModelTimeStamp
- **ModelTimeStamp** — `timestamp` string
- **ModelUHICategory** — `descriptor` ModelDescriptor, `id` string
- **ModelUHIDescriptor** — `code` string, `images` string, `long_desc` string, `name` string, `short_desc` string
- **ModelUHIFulfillments** — `id` string, `start` ModelTimeSeries, `type` string
- **ModelUHIItem** — `category_id` string, `descriptor` ModelDescriptor, `fulfillment_id` string, `id` string, `quantity` ModelQuantity
- **ModelUHILocation** — `address` string, `city` ModelDescriptor, `country` ModelDescriptor, `descriptor` ModelUHIDescriptor, `gps` string, `id` string
- **ModelUHIMessage** — `catalog` ModelCatalogDescriptor

</details>

## care-contexts

### `POST /abdm/v1/hip/care-context/data/on-fetch` ✅

Push care context data to HIU

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `202`
- request: `entries` []RecordsEntry?, `key_information` RecordsKeyMaterial, `page_count` integer, `page_number` integer, `transaction_id` string
- response: _none_

<details><summary>schemas</summary>

- **RecordsDHPubKey** — `expiry` string, `key_value` string, `parameters` string
- **RecordsEntry** — `care_context_id` string, `checksum` string, `content` string, `media` string
- **RecordsHipDataPushRequest** — `entries` []RecordsEntry?, `key_information` RecordsKeyMaterial, `page_count` integer, `page_number` integer, `transaction_id` string
- **RecordsKeyMaterial** — `crypto_alg` string, `curve` string, `dh_public_key` RecordsDHPubKey, `nonce` string

</details>

### `POST /abdm/v1/hiu/care-context/data/on-push`

Push FHIR parsing status

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `202`
- request: `consent_id` string, `status_response` []RecordsStatusResponse?, `transaction_id` string
- response: _none_

<details><summary>schemas</summary>

- **RecordsHiuDataPushRequest** — `consent_id` string, `status_response` []RecordsStatusResponse?, `transaction_id` string
- **RecordsStatusResponse** — `care_context_id` string, `description` string, `success` boolean

</details>

### `PATCH /abdm/v1/hiu/keyset`

Update public keyset

- success: `200`
- request: `nonce*` string, `public_key*` string
- response: _none_

<details><summary>schemas</summary>

- **RecordsUpdateKeysetRequest** — `nonce*` string, `public_key*` string

</details>

## care-contexts/discover

### `POST /abdm/v1/care-contexts/discover`

Discover unlinked care-contexts

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `hip_id*` string, `ref_id` string
- response: `patient` []DiscoverPatient?, `txn_id` string

<details><summary>schemas</summary>

- **DiscoverCareContext** — `display` string, `id` string
- **DiscoverPatient** — `care_contexts` []DiscoverCareContext?, `display` string, `id` string
- **DiscoverRequest** — `hip_id*` string, `ref_id` string
- **DiscoverResponse** — `patient` []DiscoverPatient?, `txn_id` string

</details>

### `POST /abdm/v1/care-contexts/discover/link/confirm`

Confirm Linking

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `otp*` string, `txn_id*` string
- response: `cc_ref_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ConfirmRequest** — `otp*` string, `txn_id*` string
- **ConfirmResponse** — `cc_ref_id` string, `txn_id` string

</details>

### `POST /abdm/v1/care-contexts/discover/link/init`

Initialise Linking

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `cc_ref_id*` string, `patient_ref_id*` string, `txn_id*` string
- response: `txn_id` string

<details><summary>schemas</summary>

- **InitRequestType3** — `cc_ref_id*` string, `patient_ref_id*` string, `txn_id*` string
- **InitResponseType4** — `txn_id` string

</details>

### `POST /abdm/v1/care-contexts/discover/link/on-confirm`

Discover Link On Confirm

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `204`
- request: `error` ModelsError, `patients` []ModelsPatient?, `request_id*` string
- response: _none_

<details><summary>schemas</summary>

- **ModelsCareContext** — `display` string, `ref_num` string
- **ModelsLinkOnConfirmHmisConnect** — `error` ModelsError, `patients` []ModelsPatient?, `request_id*` string
- **ModelsPatient** — `care_contexts` []ModelsCareContext?, `display` string, `hi_type` string, `ref_num` string

</details>

### `POST /abdm/v1/care-contexts/discover/link/on-init`

Discover Link On Init

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `204`
- request: `error` ModelsError, `otp_expiry` string, `ref_num` string, `request_id*` string, `txn_id*` string
- response: _none_

<details><summary>schemas</summary>

- **ModelsLinkOnInitHmisConnect** — `error` ModelsError, `otp_expiry` string, `ref_num` string, `request_id*` string, `txn_id*` string

</details>

### `POST /abdm/v1/care-contexts/on-discover`

Discover On Care Contexts

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `204`
- request: `error` ModelsError, `patients` []ModelsPatient?, `request_id*` string, `txn_id*` string
- response: _none_

<details><summary>schemas</summary>

- **ModelsCareContext** — `display` string, `ref_num` string
- **ModelsDiscoverCareContextHmisConnect** — `error` ModelsError, `patients` []ModelsPatient?, `request_id*` string, `txn_id*` string
- **ModelsPatient** — `care_contexts` []ModelsCareContext?, `display` string, `hi_type` string, `ref_num` string

</details>

## care-contexts/link

### `POST /abdm/v1/care-contexts/link` ✅

Care context linking

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `202`
- request: `abha_address` string, `care_contexts` []LinkCareContextMeta?, `oid` string, `partner_user_id` string
- response: _none_

<details><summary>schemas</summary>

- **LinkCareContextLinkingRequest** — `abha_address` string, `care_contexts` []LinkCareContextMeta?, `oid` string, `partner_user_id` string
- **LinkCareContextMeta** — `care_context_id` string, `data` string, `display` string, `hi_type` string, `hi_types` []string?

</details>

## care-contexts/providers

### `GET /abdm/v1/care-contexts/providers`

Get Linked Providers

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: _none_
- response: `providers` []ProviderHipDetail?

<details><summary>schemas</summary>

- **ProviderHipDetail** — `hip_id` string, `hip_name` string?
- **ProviderResponse** — `providers` []ProviderHipDetail?

</details>

## care-contexts/records

### `GET /abdm/v1/care-contexts/linked`

Get Linked CareContexts

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `hip_id`, `oid` · success: `200`
- request: _none_
- response: `care_contexts*` []RecordsCareContext?

<details><summary>schemas</summary>

- **RecordsCareContext** — `care_context_id` string, `created_at` string, `display` string, `status` string
- **RecordsResponse** — `care_contexts*` []RecordsCareContext?

</details>

### `POST /abdm/v1/care-contexts/sms/notify`

Notify unlinked care contexts

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `202`
- request: `mobile_number` string
- response: _none_

<details><summary>schemas</summary>

- **LinkUnlinkedRecordsNotificationRequest** — `mobile_number` string

</details>

## commons

### `POST /abdm/na/v1/registration/phr/check` ✅

Does Abha Address Exist

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `abha_address*` string
- response: `exists` boolean

<details><summary>schemas</summary>

- **ExistRequest** — `abha_address*` string
- **ModelsDoesHealthIdExistResponse** — `exists` boolean

</details>

### `GET /abdm/na/v1/registration/suggest` ✅

Suggest Abha Address

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `mn`, `ln`, `fn`, `dob`, `transactionId` · success: `200`
- request: _none_
- response: `suggestions` []string?

<details><summary>schemas</summary>

- **ModelsSuggestHealthIdResponse** — `suggestions` []string?

</details>

### `GET /abdm/v1/registration/pincode/{pincode}`

Pincode details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `pincode` · success: `200`
- request: _none_
- response: `dist_code` string, `dist_name` string, `pincode` string, `state_code` string, `state_name` string

<details><summary>schemas</summary>

- **PincodeResolvedPincodeData** — `dist_code` string, `dist_name` string, `pincode` string, `state_code` string, `state_name` string

</details>

## consents

### `POST /abdm/v1/consents/approve`

Consent approve

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `204`
- request: `access_mode` string, `consent_artefacts*` []ApproveConsent?, `duration` CommonsDuration, `erase_at` string, `hi_types` []string?, `id*` string
- response: _none_

<details><summary>schemas</summary>

- **ApproveConsent** — `access_mode` string, `care_contexts` []CommonsCareContext?, `duration` CommonsDuration, `erase_at` string, `hi_types` []string?, `hip_id` string
- **ApproveRequest** — `access_mode` string, `consent_artefacts*` []ApproveConsent?, `duration` CommonsDuration, `erase_at` string, `hi_types` []string?, `id*` string
- **CommonsCareContext** — `display` string, `id` string
- **CommonsDuration** — `from` string, `to` string

</details>

### `POST /abdm/v1/consents/create`

Consent create

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `204`
- request: `appointment_id` string, `care_contexts` []CommonsCareContextIdentifier, `hip_identifier` CommonsHipIdentifier, `hiu` ConsentEkaHiuIdentifier, `patient` CommonsEkaPatientReference, `period` CommonsPeriodWithExpiry, `purpose` string, `record_types` []string?
- response: `consent_init_id` string

<details><summary>schemas</summary>

- **CommonsCareContextIdentifier** — `cc_ref` string, `patient_ref` string
- **CommonsEkaPatientReference** — `health_id` string, `oid` string
- **CommonsHipIdentifier** — `id` string, `name` string
- **CommonsPeriodWithExpiry** — `expiry` string, `from` string, `to` string
- **ConsentCreateConsentRequest** — `appointment_id` string, `care_contexts` []CommonsCareContextIdentifier, `hip_identifier` CommonsHipIdentifier, `hiu` ConsentEkaHiuIdentifier, `patient` CommonsEkaPatientReference, `period` CommonsPeriodWithExpiry, `purpose` string, `record_types` []string?
- **ConsentCreateConsentResponse** — `consent_init_id` string
- **ConsentEkaHiuIdentifier** — `clinic_id` string, `d_oid` string, `requester` ConsentRequester
- **ConsentRequester** — `identifier` ConsentRequesterIden, `name` string
- **ConsentRequesterIden** — `system` string, `type` string, `value` string

</details>

### `POST /abdm/v1/consents/deny`

Consent deny

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `204`
- request: `id*` string, `reason` string
- response: _none_

<details><summary>schemas</summary>

- **DenyRequest** — `id*` string, `reason` string

</details>

### `POST /abdm/v1/consents/details`

Consent details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `204`
- request: `consent_id` string, `health_id` string
- response: `data` []ConsentConsentDetailData?

<details><summary>schemas</summary>

- **ConsentCareContextDetail** — `created_at` string, `display` string, `documents` []ConsentDocumentDetail, `id` string, `status` string
- **ConsentConsentDetailData** — `care_contexts` []ConsentCareContextDetail?, `hip` ConsentHipDetail
- **ConsentConsentDetailReq** — `consent_id` string, `health_id` string
- **ConsentConsentDetailResp** — `data` []ConsentConsentDetailData?
- **ConsentDocumentDetail** — `fhir_url` string, `id` string
- **ConsentHipDetail** — `id` string, `name` string

</details>

### `POST /abdm/v1/consents/list`

Consent list

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `204`
- request: `hiu` CommonsEkaDocReference, `patient` CommonsEkaPatientReference
- response: `consents` []ConsentConsentDetail?

<details><summary>schemas</summary>

- **CommonsEkaDocReference** — `clinic_id` string, `d_oid` string
- **CommonsEkaPatientReference** — `health_id` string, `oid` string
- **CommonsPeriodWithExpiry** — `expiry` string, `from` string, `to` string
- **ConsentConsentDetail** — `c_at` string, `consent_id` string, `consent_init_id` string, `hi_types` []string?, `period` CommonsPeriodWithExpiry, `status` string, `u_at` string
- **ConsentListConsentRequest** — `hiu` CommonsEkaDocReference, `patient` CommonsEkaPatientReference
- **ConsentListConsentResponse** — `consents` []ConsentConsentDetail?

</details>

### `POST /abdm/v1/consents/revoke`

Consent revoke

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `204`
- request: `consent_artefacts` []string?
- response: _none_

<details><summary>schemas</summary>

- **RevokeRequest** — `consent_artefacts` []string?

</details>

## consents/auto-approval

### `GET /abdm/v1/consents/auto-approval`

Get status

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: _none_
- response: `status` string

<details><summary>schemas</summary>

- **SettingsResponse** — `status` string

</details>

### `PATCH /abdm/v1/consents/auto-approval`

Update auto approval policy

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: _none_
- response: `status` string

<details><summary>schemas</summary>

- **SettingsResponse** — `status` string

</details>

## enrollment/aadhaar

### `POST /abdm/na/v1/registration/aadhaar/auto-login`

Auto Login

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `abha_address*` string, `txn_id*` string
- response: `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

<details><summary>schemas</summary>

- **CreateRequest** — `abha_address*` string, `txn_id*` string
- **CreateResponse** — `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string
- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?

</details>

### `POST /abdm/na/v1/registration/aadhaar/create-phr` ✅

Create

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `abha_address*` string, `txn_id*` string
- response: `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

<details><summary>schemas</summary>

- **CreateRequest** — `abha_address*` string, `txn_id*` string
- **CreateResponse** — `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string
- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?

</details>

### `POST /abdm/na/v1/registration/aadhaar/init` ✅

Generate OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `aadhaar_number*` string
- response: `hint` string?, `txn_id` string

<details><summary>schemas</summary>

- **InitRequest** — `aadhaar_number*` string
- **InitResponse** — `hint` string?, `txn_id` string

</details>

### `POST /abdm/na/v1/registration/aadhaar/mobile/resend` ✅

Mobile Resend OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `txn_id*` string
- response: `hint` string?, `txn_id` string

<details><summary>schemas</summary>

- **ResendRequestType2** — `txn_id*` string
- **ResendResponseType2** — `hint` string?, `txn_id` string

</details>

### `POST /abdm/na/v1/registration/aadhaar/mobile/verify` ✅

Mobile Verify OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `otp*` string, `txn_id*` string
- response: `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string?, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

<details><summary>schemas</summary>

- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?
- **VerifyAbhaProfile** — `abha_address` string, `kyc_verified` string, `name` string
- **VerifyRequestType2** — `otp*` string, `txn_id*` string
- **VerifyResponseType2** — `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string?, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

</details>

### `POST /abdm/na/v1/registration/aadhaar/resend` ✅

Resend OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `txn_id*` string
- response: `hint` string?, `txn_id` string

<details><summary>schemas</summary>

- **ResendRequest** — `txn_id*` string
- **ResendResponse** — `hint` string?, `txn_id` string

</details>

### `POST /abdm/na/v1/registration/aadhaar/verify` ✅

Verify OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `mobile*` string, `otp*` string, `txn_id*` string
- response: `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string?, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

<details><summary>schemas</summary>

- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?
- **VerifyAbhaProfile** — `abha_address` string, `kyc_verified` string, `name` string
- **VerifyRequest** — `mobile*` string, `otp*` string, `txn_id*` string
- **VerifyResponse** — `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string?, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

</details>

## enrollment/face-auth

### `POST /abdm/na/v1/face-auth/capture`

Capture Face

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `txn_id` string
- response: `message` string, `status` string, `txn_id` string

<details><summary>schemas</summary>

- **FaceAuthCapturePIDRequest** — `txn_id` string
- **FaceAuthCapturePIDResponse** — `message` string, `status` string, `txn_id` string

</details>

### `POST /abdm/na/v1/face-auth/init`

Face Auth Init

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: _none_
- response: `hint` string?, `txn_id` string, `url` string

<details><summary>schemas</summary>

- **FaceAuthInitFaceAuthResponse** — `hint` string?, `txn_id` string, `url` string

</details>

### `POST /abdm/na/v1/face-auth/verify`

Verify Face

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `aadhaar` string, `txn_id` string
- response: `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string?, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

<details><summary>schemas</summary>

- **FaceAuthVerifyFaceAuthRequest** — `aadhaar` string, `txn_id` string
- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?
- **VerifyAbhaProfile** — `abha_address` string, `kyc_verified` string, `name` string
- **VerifyResponse** — `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string?, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `token` string?, `txn_id` string

</details>

## enrollment/mobile

### `POST /abdm/na/v1/registration/mobile/create-phr` ✅

Create

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `abha_address*` string, `profile*` DetailsRequestType2, `txn_id*` string
- response: `eka` RegistrationEkaIds, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `success` boolean, `token` string?

<details><summary>schemas</summary>

- **CreateRequestType2** — `abha_address*` string, `profile*` DetailsRequestType2, `txn_id*` string
- **CreateResponseType2** — `eka` RegistrationEkaIds, `profile` ProfileResponse, `refresh_token` string?, `skip_state` string, `success` boolean, `token` string?
- **DetailsRequestType2** — `address` string?, `day_of_birth*` integer, `first_name*` string, `gender*` string, `last_name` string?, `middle_name` string?, `month_of_birth*` integer, `pincode*` string, `year_of_birth*` integer
- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?

</details>

### `POST /abdm/na/v1/registration/mobile/init` ✅

Generate OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `mobile_number*` string
- response: `hint` string?, `txn_id` string

<details><summary>schemas</summary>

- **InitRequestType2** — `mobile_number*` string
- **InitResponseType2** — `hint` string?, `txn_id` string

</details>

### `POST /abdm/na/v1/registration/mobile/resend` ✅

Resend OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `txn_id*` string
- response: `hint` string, `txn_id` string

<details><summary>schemas</summary>

- **ResendRequestType3** — `txn_id*` string
- **ResendResponseType3** — `hint` string, `txn_id` string

</details>

### `POST /abdm/na/v1/registration/mobile/verify` ✅

Verify OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `otp*` string, `txn_id*` string
- response: `abha_profiles` []VerifyAbhaProfile?, `eka` RegistrationEkaIds, `skip_state` string, `txn_id` string

<details><summary>schemas</summary>

- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?
- **VerifyAbhaProfile** — `abha_address` string, `kyc_verified` string, `name` string
- **VerifyRequestType3** — `otp*` string, `txn_id*` string
- **VerifyResponseType3** — `abha_profiles` []VerifyAbhaProfile?, `eka` RegistrationEkaIds, `skip_state` string, `txn_id` string

</details>

## login

### `POST /abdm/na/v1/profile/login/init` ✅

Generate OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `identifier` string, `method` string, `otp_system` string, `txn_id` string
- response: `hint` string, `txn_id` string

<details><summary>schemas</summary>

- **PhrLoginOTPRequest** — `identifier` string, `method` string, `otp_system` string, `txn_id` string
- **PhrLoginOTPResponse** — `hint` string, `txn_id` string

</details>

### `POST /abdm/na/v1/profile/login/methods`

Search Auth Methods

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `abha_address` string
- response: `abha_address` string, `abha_number` string, `auth_methods` []string?, `blocked_methods` []string?, `fullname` string, `mobile` string, `status` string

<details><summary>schemas</summary>

- **PhrSearchAuthMethodRequest** — `abha_address` string
- **PhrSearchAuthMethodResponse** — `abha_address` string, `abha_number` string, `auth_methods` []string?, `blocked_methods` []string?, `fullname` string, `mobile` string, `status` string

</details>

### `POST /abdm/na/v1/profile/login/phr` ✅

Login

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `phr_address` string, `txn_id` string
- response: `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `skip_state` string, `txn_id` string

<details><summary>schemas</summary>

- **PhrPhrAddressLoginRequest** — `phr_address` string, `txn_id` string
- **PhrVerifyLoginOTPResponse** — `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `skip_state` string, `txn_id` string
- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?
- **VerifyAbhaProfile** — `abha_address` string, `kyc_verified` string, `name` string

</details>

### `POST /abdm/na/v1/profile/login/verify` ✅

Verify OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `otp` string, `txn_id` string
- response: `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `skip_state` string, `txn_id` string

<details><summary>schemas</summary>

- **PhrVerifyLoginOTPRequest** — `otp` string, `txn_id` string
- **PhrVerifyLoginOTPResponse** — `abha_profiles` []VerifyAbhaProfile, `eka` RegistrationEkaIds, `hint` string, `profile` ProfileResponse, `skip_state` string, `txn_id` string
- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?
- **RegistrationEkaIds** — `min_token` string, `oid` string?, `uuid` string?
- **VerifyAbhaProfile** — `abha_address` string, `kyc_verified` string, `name` string

</details>

## nhpr-abdm

### `POST /abdm/v1/hip/onboard`

Onboard facility

- success: `200`
- request: `clinic_id` string, `hip_id*` string, `name*` string
- response: `hip_code` string, `hip_id` string, `hip_name` string, `scan_share_url` string

<details><summary>schemas</summary>

- **ModelsAddFacilityRequest** — `clinic_id` string, `hip_id*` string, `name*` string
- **ModelsAddFacilityResponse** — `hip_code` string, `hip_id` string, `hip_name` string, `scan_share_url` string

</details>

## nhpr-abdm/hfr

### `POST /abdm/nhpr/v1/hfr/link`

Link HFR

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `clinic_id` string, `hfr_id` string
- response: `data` NhprUpdateBridgeRequest, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprUpdateBridgeRequest** — `data` NhprUpdateBridgeRequest, `error`, `success` boolean
- **NhprUpdateBridgeRequest** — `clinic_id` string, `hfr_id` string

</details>

### `POST /abdm/nhpr/v1/hfr/onboard`

Onboard HFR

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `clinic_id` string, `hfr_id` string
- response: `data` NhprOnboardFacilityResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprOnboardFacilityResponse** — `data` NhprOnboardFacilityResponse, `error`, `success` boolean
- **NhprOnboardFacilityResponse** — `facility_id` string
- **NhprUpdateBridgeRequest** — `clinic_id` string, `hfr_id` string

</details>

### `GET /abdm/nhpr/v1/hfr/search/{hfrId}`

Search HFR

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `hfrId` · success: `200`
- request: _none_
- response: `data` ModelsFacility, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmInternalNhprRepoModelsFacility** — `data` ModelsFacility, `error`, `success` boolean
- **ModelsFacility** — `address` string, `districtLGDCode` string, `districtName` string, `facilityId` string, `facilityName` string, `facilityStatus` string, `facilityType` string, `facilityTypeCode` string, `latitude` string, `longitude` string, `ownership` string, `ownershipCode` string, `pincode` string, `stateLGDCode` string, `stateName` string, `subDistrictLGDCode` string, `subDistrictName` string, `systemOfMedicine` string, `systemOfMedicineCode` string, `villageCityTownLGDCode`, `villageCityTownName`

</details>

## nhpr-abdm/hpr

### `POST /abdm/nhpr/v1/hpr/aadhaar/check/available`

Check If HPR is Available

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `hprId` string
- response: `data` NhprDoesHprExistResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprDoesHprExistResponse** — `data` NhprDoesHprExistResponse, `error`, `success` boolean
- **NhprDoesHprExistRequest** — `hprId` string
- **NhprDoesHprExistResponse** — `available` boolean

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/check/exist`

Get HPR Info

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `mobile_num` string, `txn_id` string
- response: `data` NhprEKAIsHPRIDExistsResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprEKAIsHPRIDExistsResponse** — `data` NhprEKAIsHPRIDExistsResponse, `error`, `success` boolean
- **NhprEKAIsHPRIDExistsResponse** — `address` string, `category_id` integer, `day_of_birth` string, `district_code` string, `district_name` string, `first_name` string, `gender` string, `hpr_id` string, `hpr_id_number` string, `last_name` string, `middle_name` string, `mobile` string?, `month_of_birth` string, `name` string, `new` boolean, `pincode` string, `profile_photo` string, `state_code` string, `state_name` string, `sub_category_id` integer, `token` string, `txn_id` string, `year_of_birth` string
- **NhprTransactionId** — `mobile_num` string, `txn_id` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/create`

Create HPR

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `council` boolean, `domain_name` string, `email` string, `first_name` string, `hp_category_code` integer, `hp_sub_category_code` integer, `hprId` string, `id_type` string, `last_name` string, `middle_name` string, `notify_user` boolean, `password` string, `pincode` string, `profile_photo` string, `role` integer, `source_type` string, `txnId` string
- response: `data` NhprEKACreateHPRIDResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprEKACreateHPRIDResponse** — `data` NhprEKACreateHPRIDResponse, `error`, `success` boolean
- **NhprEKACreateHPRID** — `council` boolean, `domain_name` string, `email` string, `first_name` string, `hp_category_code` integer, `hp_sub_category_code` integer, `hprId` string, `id_type` string, `last_name` string, `middle_name` string, `notify_user` boolean, `password` string, `pincode` string, `profile_photo` string, `role` integer, `source_type` string, `txnId` string
- **NhprEKACreateHPRIDResponse** — `auth_methods` []string?, `categories`, `category_id` integer, `day_of_birth` string, `district_code` string, `district_name` string, `email` string, `first_name` string, `gender` string, `hpr_id_number` string, `hprId` string, `kyc_photo` string, `last_name` string, `middle_name` string, `mobile` string, `month_of_birth` string, `name` string, `new` boolean, `state_code` string, `state_name` string, `sub_category_id` integer, `year_of_birth` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/demo/verify`

Verify Mobile

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `mobile_num` string, `txn_id` string
- response: `data` NhprEkaVerifiedResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprEkaVerifiedResponse** — `data` NhprEkaVerifiedResponse, `error`, `success` boolean
- **NhprEkaVerifiedResponse** — `error_code` string?, `hpr` object, `reason` string?, `suggested_hpr_ids` []string?, `txnId` string, `uidai_token` string?, `verified` boolean
- **NhprMobileOTPRequest** — `mobile_num` string, `txn_id` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/init`

Get Aadhaar OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `aadhaar` string
- response: `data` NhprEkaAadhaarResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprEkaAadhaarResponse** — `data` NhprEkaAadhaarResponse, `error`, `success` boolean
- **ModelsHPRIDAadhaarInitRequest** — `aadhaar` string
- **NhprEkaAadhaarResponse** — `hint` string, `mobile_num` string, `txn_id` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/mobile/init`

Get Mobile OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `mobile_num` string, `txn_id` string
- response: `hint` string, `mobile_num` string, `txn_id` string

<details><summary>schemas</summary>

- **NhprEkaAadhaarResponse** — `hint` string, `mobile_num` string, `txn_id` string
- **NhprMobileOTPRequest** — `mobile_num` string, `txn_id` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/mobile/verify`

Verify Mobile OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `otp` string, `txnId` string
- response: `data` NhprEkaAadhaarResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprEkaAadhaarResponse** — `data` NhprEkaAadhaarResponse, `error`, `success` boolean
- **NhprEkaAadhaarResponse** — `hint` string, `mobile_num` string, `txn_id` string
- **NhprMobileVerifyRequest** — `otp` string, `txnId` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/suggest`

Suggest IDs

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `mobile_num` string, `txn_id` string
- response: `data` []string

<details><summary>schemas</summary>

- **NhprTransactionId** — `mobile_num` string, `txn_id` string

</details>

### `POST /abdm/nhpr/v1/hpr/aadhaar/verify`

Verify Aadhaar OTP

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `otp` string, `txn_id` string
- response: `data` NhprEkaAadhaarResponse, `error`, `success` boolean

<details><summary>schemas</summary>

- **HandlerTypedAPIResponseGithubComEkaCareNdhmPkgNhprNhprEkaAadhaarResponse** — `data` NhprEkaAadhaarResponse, `error`, `success` boolean
- **NhprEKAAadhaarVerifyRequest** — `otp` string, `txn_id` string
- **NhprEkaAadhaarResponse** — `hint` string, `mobile_num` string, `txn_id` string

</details>

## patient-requests

### `GET /abdm/v1/request`

Requests get details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `id`, `type`, `status`, `oid` · success: `200`
- request: _none_
- response: `access_mode` string, `created_at` string, `duration` CommonsDuration, `erase_at` string, `hi_types` []CommonsHiType?, `hiu` CommonsHiu, `id` string, `providers` []DetailsProviders?, `purpose` CommonsPurpose, `requester` CommonsRequester, `status` string, `updated_at` string

<details><summary>schemas</summary>

- **CommonsDuration** — `from` string, `to` string
- **CommonsHiType** — `display` string, `enabled` boolean, `id` string
- **CommonsHip** — `id` string, `name` string
- **CommonsHiu** — `id` string, `name` string
- **CommonsPurpose** — `code` string, `text` string
- **CommonsRequester** — `name` string
- **DetailsCareContext** — `display` string, `id` string
- **DetailsProviders** — `care_contexts` []DetailsCareContext?, `consent_artefact_id` string?, `hip` CommonsHip
- **DetailsResponse** — `access_mode` string, `created_at` string, `duration` CommonsDuration, `erase_at` string, `hi_types` []CommonsHiType?, `hiu` CommonsHiu, `id` string, `providers` []DetailsProviders?, `purpose` CommonsPurpose, `requester` CommonsRequester, `status` string, `updated_at` string

</details>

### `GET /abdm/v1/requests`

List requests

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `status`, `type`, `oid` · success: `200`
- request: _none_
- response: `authorizations` []ListsRequestDetail?, `consents` []ListsRequestDetail?, `subscriptions` []ListsRequestDetail?

<details><summary>schemas</summary>

- **CommonsDuration** — `from` string, `to` string
- **CommonsHiu** — `id` string, `name` string
- **CommonsPurpose** — `code` string, `text` string
- **CommonsRequester** — `name` string
- **ListsRequestDetail** — `created_at` string, `duration` CommonsDuration, `hiu` CommonsHiu, `id` string, `purpose` CommonsPurpose, `requester` CommonsRequester
- **ListsResponse** — `authorizations` []ListsRequestDetail?, `consents` []ListsRequestDetail?, `subscriptions` []ListsRequestDetail?

</details>

## phys-cons

### `GET /abdm/uhi/v1/physical-consultation/booking/cancel`

Get cancel

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `txn_id`, `provider_id`, `order_id` · success: `200`
- request: _none_
- response: `status` string

<details><summary>schemas</summary>

- **ModelsCancelStatusResponse** — `status` string

</details>

### `POST /abdm/uhi/v1/physical-consultation/booking/cancel`

Cancel

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `order_id` string, `provider_id` string, `transaction_id` string
- response: `order_id` string, `provider_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ModelsCancelOrderRequest** — `order_id` string, `provider_id` string, `transaction_id` string
- **ModelsGetPhysConsCommonResponse** — `order_id` string, `provider_id` string, `txn_id` string

</details>

### `GET /abdm/uhi/v1/physical-consultation/booking/confirm`

Get confirm

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `txn_id`, `provider_id`, `order_id` · success: `200`
- request: _none_
- response: `doctor_info` ModelsDoctorInfo, `order_id` string, `pin` string, `slot` ModelsSlot, `state` string

<details><summary>schemas</summary>

- **ModelsDoctorInfo** — `education` string, `experience` string, `gender` string, `id` string, `languages` []string, `name` string
- **ModelsGetConfirmResponse** — `doctor_info` ModelsDoctorInfo, `order_id` string, `pin` string, `slot` ModelsSlot, `state` string
- **ModelsSlot** — `end` string, `slot_id` string, `start` string

</details>

### `POST /abdm/uhi/v1/physical-consultation/booking/confirm`

Confirm

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `provider_id` string, `transaction_id` string
- response: `order_id` string, `provider_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ModelsConfirmRequest** — `provider_id` string, `transaction_id` string
- **ModelsGetPhysConsCommonResponse** — `order_id` string, `provider_id` string, `txn_id` string

</details>

### `GET /abdm/uhi/v1/physical-consultation/booking/init`

Get init

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `txn_id`, `provider_id`, `order_id` · success: `200`
- request: _none_
- response: `billing` ModelsBilling, `consultation` ModelsConsultation, `doctor` ModelsDoctor, `facility_id` string, `order_id` string, `patient` ModelsPatientType2, `payment` ModelsPayment, `quote` ModelsQuote, `terms` []ModelsTermsSummary?

<details><summary>schemas</summary>

- **ModelsAppointmentSummary** — `billing` ModelsBilling, `consultation` ModelsConsultation, `doctor` ModelsDoctor, `facility_id` string, `order_id` string, `patient` ModelsPatientType2, `payment` ModelsPayment, `quote` ModelsQuote, `terms` []ModelsTermsSummary?
- **ModelsBilling** — `address` string, `email` string, `name` string, `phone` string
- **ModelsConsultation** — `end_time` string, `slot_id` string, `start_time` string, `type` string
- **ModelsDoctor** — `doctor_gender` string, `doctor_image` string, `id` string, `name` string, `tags` []string
- **ModelsPatientType2** — `dob` string, `gender` string, `id` string
- **ModelsPayment** — `status` string, `type` string
- **ModelsQuote** — `breakup` []ModelsQuoteItem?, `total` string
- **ModelsQuoteItem** — `title` string, `value` string
- **ModelsTermsSummary** — `long_desc` string, `short_desc` string, `terms_state` string, `type` string

</details>

### `POST /abdm/uhi/v1/physical-consultation/booking/init`

Init

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `billing` ModelsInitBilling, `doctor_id` string, `end_time` string, `facility_id` string, `price` ModelsPrice, `provider_id` string, `slot_id` string, `start_time` string, `transaction_id` string
- response: `order_id` string, `provider_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ModelsAddress** — `city` string, `country` string, `door` string, `locality` string, `name` string, `state` string
- **ModelsGetPhysConsCommonResponse** — `order_id` string, `provider_id` string, `txn_id` string
- **ModelsInitBilling** — `address` ModelsAddress, `name` string, `phone` string
- **ModelsInitRequest** — `billing` ModelsInitBilling, `doctor_id` string, `end_time` string, `facility_id` string, `price` ModelsPrice, `provider_id` string, `slot_id` string, `start_time` string, `transaction_id` string
- **ModelsPrice** — `currency` string, `value` string

</details>

### `POST /abdm/uhi/v1/physical-consultation/booking/update`

Update

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `end_time` string, `facility_id` string, `order_id` string, `provider_id` string, `slot_id` string, `start_time` string, `transaction_id` string
- response: `order_id` string, `provider_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ModelsGetPhysConsCommonResponse** — `order_id` string, `provider_id` string, `txn_id` string
- **ModelsUpdateOrderRequest** — `end_time` string, `facility_id` string, `order_id` string, `provider_id` string, `slot_id` string, `start_time` string, `transaction_id` string

</details>

### `GET /abdm/uhi/v1/physical-consultation/search/doctors`

Get doctors

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: _none_
- response: `data` []ModelsSearchResponse?

<details><summary>schemas</summary>

- **ModelsDoctorsSearchResponse** — `data` []ModelsSearchResponse?
- **ModelsSearchResponse** — `doctor_name` string, `education` string, `facility_id` string, `gender` string, `hpr_id` string, `image` string, `languages` []string, `provider_address` string, `provider_city` string, `provider_country` string, `provider_gps` string, `provider_name` string, `specialities` []string, `transaction_id` string, `year_of_experience` integer

</details>

### `POST /abdm/uhi/v1/physical-consultation/search/doctors`

Search doctors

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `district_code` string?, `district_name` string?, `doctor_name` string?, `end_time` string, `facility_name` string?, `hpr_id` string?, `latitude` string?, `longitude` string?, `pincode` string?, `radius_km` string?, `speciality` string?, `start_time` string, `state_code` string?, `state_name` string?
- response: `order_id` string, `provider_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ModelsDoctorsSearchRequest** — `district_code` string?, `district_name` string?, `doctor_name` string?, `end_time` string, `facility_name` string?, `hpr_id` string?, `latitude` string?, `longitude` string?, `pincode` string?, `radius_km` string?, `speciality` string?, `start_time` string, `state_code` string?, `state_name` string?
- **ModelsGetPhysConsCommonResponse** — `order_id` string, `provider_id` string, `txn_id` string

</details>

### `GET /abdm/uhi/v1/physical-consultation/search/doctors/slots`

Get slots

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `txn_id`, `provider_id`, `order_id` · success: `200`
- request: _none_
- response: `doctor` ModelsDoctorInfo, `provider_info` ModelsProviderInfo, `slots` []ModelsSlotDetail

<details><summary>schemas</summary>

- **ModelsDoctorInfo** — `education` string, `experience` string, `gender` string, `id` string, `languages` []string, `name` string
- **ModelsProviderInfo** — `address` string, `city` string, `id` string, `name` string, `phone` string
- **ModelsSlotDetail** — `currency` string, `date` string, `end_time` string, `id` string, `price` string, `start_time` string, `type` string
- **ModelsSlotsData** — `doctor` ModelsDoctorInfo, `provider_info` ModelsProviderInfo, `slots` []ModelsSlotDetail

</details>

### `POST /abdm/uhi/v1/physical-consultation/search/doctors/slots`

Search slots

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `end_time*` string, `facility_id*` string, `hpr_id*` string, `provider_id*` string, `start_time*` string, `transaction_id*` string
- response: `order_id` string, `provider_id` string, `txn_id` string

<details><summary>schemas</summary>

- **ModelsGetPhysConsCommonResponse** — `order_id` string, `provider_id` string, `txn_id` string
- **ModelsSlotSearchRequest** — `end_time*` string, `facility_id*` string, `hpr_id*` string, `provider_id*` string, `start_time*` string, `transaction_id*` string

</details>

## profile/cards

### `GET /abdm/v1/profile/asset/card` ✅

ABHA card

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: _none_
- response: _none_

### `GET /abdm/v1/profile/asset/qr` ✅

ABHA QR Code

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `format`, `oid` · success: `200`
- request: _none_
- response: `address` string, `dist name` string, `distlgd` string, `district_name` string, `dob` string, `gender` string, `hid` string, `hidn` string, `mobile` string, `name` string, `phr` string, `state name` string, `statelgd` string

<details><summary>schemas</summary>

- **AssetQRData** — `address` string, `dist name` string, `distlgd` string, `district_name` string, `dob` string, `gender` string, `hid` string, `hidn` string, `mobile` string, `name` string, `phr` string, `state name` string, `statelgd` string

</details>

## profile/details

### `DELETE /abdm/v1/profile` ✅

Delete Profile

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `204`
- request: _none_
- response: _none_

### `GET /abdm/v1/profile` ✅

Profile Details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: _none_
- response: `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?

<details><summary>schemas</summary>

- **ProfileResponse** — `abha_address` string, `abha_addresses` []string, `abha_number` string?, `address` string?, `day_of_birth` integer?, `district_code` string, `district_name` string, `first_name` string?, `full_name` string, `gender` string, `kyc_verified` boolean?, `last_name` string?, `middle_name` string?, `mobile` string?, `month_of_birth` integer?, `pincode` string?, `profile_photo` string, `state_code` string, `state_name` string, `sub_district_code` string, `sub_district_name` string, `town_code` string, `town_name` string, `village_code` string, `village_name` string, `ward_code` string, `ward_name` string, `year_of_birth` integer?

</details>

### `PATCH /abdm/v1/profile` ✅

Update Profile

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `204`
- request: `address` string?, `day_of_birth*` integer, `first_name*` string, `gender*` string, `last_name` string?, `middle_name` string?, `month_of_birth*` integer, `pincode*` string, `year_of_birth*` integer
- response: _none_

<details><summary>schemas</summary>

- **DetailsRequestType2** — `address` string?, `day_of_birth*` integer, `first_name*` string, `gender*` string, `last_name` string?, `middle_name` string?, `month_of_birth*` integer, `pincode*` string, `year_of_birth*` integer

</details>

## profile/kyc

### `POST /abdm/v1/profile/kyc/init` ✅

Kyc init

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `identifier*` string, `method*` string, `user_x_token` string
- response: `txn_id` string

<details><summary>schemas</summary>

- **KycGenerateKycOtpRequest** — `identifier*` string, `method*` string, `user_x_token` string
- **KycGenerateKycOtpResponse** — `txn_id` string

</details>

### `POST /abdm/v1/profile/kyc/resend` ✅

Kyc resend

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `txn_id*` string
- response: `txn_id` string

<details><summary>schemas</summary>

- **KycGenerateKycOtpResponse** — `txn_id` string
- **KycResendKycOtpRequest** — `txn_id*` string

</details>

### `POST /abdm/v1/profile/kyc/verify` ✅

Kyc verify

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `otp*` string, `txn_id*` string, `user_x_token` string
- response: `txn_id` string

<details><summary>schemas</summary>

- **KycVerifyKycOtpRequest** — `otp*` string, `txn_id*` string, `user_x_token` string
- **KycVerifyKycOtpResponse** — `txn_id` string

</details>

## profile/search

### `POST /abdm/v1/profile/search`

Search ABHA

- success: `200`
- request: `mobile*` string
- response: `profiles` []ProfileAbha?, `txn_id` string

<details><summary>schemas</summary>

- **ProfileAbha** — `abha_address` string, `abha_number` string, `age` integer, `gender` string, `index` integer, `internal` boolean, `kyc_verified` string, `name` string
- **ProfileDiscoverAbhaRequest** — `mobile*` string
- **ProfileSearchKycAbhaResponse** — `profiles` []ProfileAbha?, `txn_id` string

</details>

## providers

### `GET /abdm/v1/provider/{hip_id}`

Get Provider Details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `hip_id` · query: `oid` · success: `200`
- request: _none_
- response: `city` string?, `hip_id` string, `hip_name` string?, `state` string?

<details><summary>schemas</summary>

- **ProviderResponseType2** — `city` string?, `hip_id` string, `hip_name` string?, `state` string?

</details>

## providers/search

### `GET /abdm/v1/providers`

Search Providers

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `name`, `oid` · success: `200`
- request: _none_
- response: _none_

<details><summary>schemas</summary>

- **ProviderResponseType2** — `city` string?, `hip_id` string, `hip_name` string?, `state` string?

</details>

## scan-and-pay

### `POST /abdm/v1/profile/scan-and-pay/open-order`

Open Orders

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `counter_id` string, `hip_id` string
- response: `request_id` string, `status` string

### `GET /abdm/v1/profile/scan-and-pay/open-order/{request_id}`

Get Open Orders

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `request_id` · success: `200`
- request: _none_
- response: `data` object, `request_id` string, `status` string?

<details><summary>schemas</summary>

- **ScanAndPayErrorCodeMessage** — `code` string, `message` string
- **ScanAndPayProcedures** — `category` string, `services` []ScanAndPayService?
- **ScanAndPayService** — `amount` number, `description` string, `name` string, `serviceId` string

</details>

### `POST /abdm/v1/profile/scan-and-pay/patient-select`

Patient Select

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: `intent*` string, `procedures*` []ScanAndPayProcedures?, `request_id*` string
- response: `request_id` string, `status` string?

<details><summary>schemas</summary>

- **ScanAndPayPatientSelectEkaRequest** — `intent*` string, `procedures*` []ScanAndPayProcedures?, `request_id*` string
- **ScanAndPayProcedures** — `category` string, `services` []ScanAndPayService?
- **ScanAndPayService** — `amount` number, `description` string, `name` string, `serviceId` string

</details>

### `GET /abdm/v1/profile/scan-and-pay/patient-select/{request_id}`

Get Patient Select Details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `request_id` · success: `200`
- request: _none_
- response: `data` object, `request_id` string, `status` string

<details><summary>schemas</summary>

- **ScanAndPayErrorCodeMessage** — `code` string, `message` string
- **ScanAndPayPaymentEka** — `amount` number, `description` string, `merchant_id` string, `order_number` string, `payment_mode` string, `payment_url` string
- **ScanAndPayProcedures** — `category` string, `services` []ScanAndPayService?
- **ScanAndPayService** — `amount` number, `description` string, `name` string, `serviceId` string

</details>

### `GET /abdm/v1/profile/scan-and-pay/transaction-details/{request_id}`

Get Transaction Details

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `request_id` · success: `200`
- request: _none_
- response: `data` ScanAndPayAcknowledgement, `request_id` string, `status` string?

<details><summary>schemas</summary>

- **ScanAndPayAcknowledgement** — `abhaAddress` string, `openOrderRequestId` string, `orderNumber` string, `paymentDate` string, `paymentReceiptLink` string, `status` string, `transactionId` string

</details>

## scan-and-share

### `POST /abdm/v1/hip/patient/profile/on-share`

On-Share

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `202`
- request: `abha_address*` string, `error` CommonsError, `hip_code` string, `request_id*` string, `token` ScanandshareToken
- response: _none_

<details><summary>schemas</summary>

- **CommonsError** — `code` integer, `message` string
- **ScanandshareHipScanAndShareOnRequest** — `abha_address*` string, `error` CommonsError, `hip_code` string, `request_id*` string, `token` ScanandshareToken
- **ScanandshareToken** — `expires_at` integer, `token_number*` string

</details>

### `POST /abdm/v2/profile/share`

Request to generate token

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `202`
- request: `counter_id*` string, `hip_id*` string, `location` HipLocation, `patient` ScanandsharePatient, `request_id` string, `user_x_token` string
- response: `request_id` string, `status` string

<details><summary>schemas</summary>

- **HipLocation** — `latitude` string?, `longitude` string?
- **ScanandshareAddress** — `district` string, `line` string, `pin_code*` string, `state` string
- **ScanandsharePatient** — `abha_address*` string, `address` ScanandshareAddress, `day_of_birth*` string, `gender*` string, `mobile*` string, `month_of_birth*` string, `name*` string, `year_of_birth*` string
- **ScanandshareProfileShareRequest** — `counter_id*` string, `hip_id*` string, `location` HipLocation, `patient` ScanandsharePatient, `request_id` string, `user_x_token` string

</details>

### `GET /abdm/v2/profile/share/{request_id}`

Get appointment token

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · path: `request_id` · success: `200`
- request: _none_
- response: `expiry` string?, `request_id` string, `status` string, `token_number` string

<details><summary>schemas</summary>

- **ScanandshareResponse** — `expiry` string?, `request_id` string, `status` string, `token_number` string

</details>

## session

### `POST /abdm/v1/session/init` ✅

Init

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `abha_address*` string
- response: `txn_id` string

<details><summary>schemas</summary>

- **InitRequestType4** — `abha_address*` string
- **InitResponseType5** — `txn_id` string

</details>

### `GET /abdm/v1/session/status`

Status

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · success: `200`
- request: _none_
- response: `hpr_id` string, `kyc_verified` boolean, `logged_in` boolean

<details><summary>schemas</summary>

- **AuthLoggedInStatusResponse** — `hpr_id` string, `kyc_verified` boolean, `logged_in` boolean

</details>

### `POST /abdm/v1/session/verify` ✅

Verify

- headers: `X-Pt-Id`, `X-Partner-Pt-Id`, `X-Hip-Id` · query: `oid` · success: `200`
- request: `otp*` string, `txn_id*` string
- response: `refresh_token` string?, `token` string

<details><summary>schemas</summary>

- **VerifyRequestType6** — `otp*` string, `txn_id*` string
- **VerifyResponseType4** — `refresh_token` string?, `token` string

</details>

## webhooks

All webhooks POST to the integrator's registered endpoint with header
`Eka-Webhook-Signature: t=<unix>,v1=<hex>` — HMAC-SHA256 over `"{t}.{rawBody}"` using the
subscription's signing key. Reject when `|now - t|` exceeds ~3 minutes.

Common envelope: `service`, `event`, `event_time`, `transaction_id`, `timestamp`,
`business_id`, `client_id`, `data`. Fields below are those of `data`.

| Event | Description | `data` fields |
|---|---|---|
| `abha.care_context_discover` | Discovery Request Received | `abha_address`, `gender`, `hip_id`, `identifiers`, `oid`, `partner_patient_id`, `patient_name`, `request_id`, `txn_id`, `year_of_birth` |
| `abha.care_context_discover_link_init` | Link Initiation - Generate OTP | `abha_address`, `hip_id`, `oid`, `partner_patient_id`, `patient`, `request_id`, `txn_id` |
| `abha.consent_update` | Consent update | `notification`, `status`, `timestamp` |
| `abha.context_discover_link_confirm` | Link Confirmation - Verify OTP | `hip_id`, `linkRefNumber`, `oid`, `partner_patient_id`, `request_id`, `token` |
| `abha.created` | ABHA Address Created | `abha_address`, `abha_number`, `hip_id`, `oid`, `partner_patient_id` |
| `abha.hip_data_fetch` | Health Data Requested by HIU | `abha_address`, `care_contexts`, `hi_types`, `hip_id`, `key_information`, `oid`, `partner_patient_id`, `transaction_id` |
| `abha.hip_profile_share` | Scan and Share Token Received | `abha_address`, `abha_number`, `hip_code`, `hip_id`, `oid`, `partner_patient_id`, `req_id` |
| `abha.hiu_data_push` | Health Data Received from HIP | `abha_address`, `abha_number`, `consent_id`, `entries`, `hip_id`, `key_information`, `oid`, `partner_patient_id`, `transaction_id` |
| `abha.link_care_context` | Care Context Linking Status | `abha_address`, `care_context_id`, `error`, `hip_id`, `oid`, `partner_patient_id`, `retry_count`, `status` |
| `abha.locker_created` | Locker Created | `abha_address`, `auto_approval_id`, `hip_id`, `oid`, `partner_patient_id`, `subscription_id` |
| `abha.subscription_modified` | Subscription Updated | `abha_address`, `hip_id`, `oid`, `partner_patient_id`, `subscription_details` |
| `abha.subscription_notify` | New Care Context Linked (Subscription) | `abha_address`, `abha_number`, `hip_id`, `oid`, `partner_patient_id`, `subscription_meta` |
| `booking.cancel` | Booking cancel | `status` |
| `booking.confirm` | Booking confirm | `doctor_info`, `order_id`, `pin`, `slot`, `state` |
| `booking.init` | Booking init | `billing`, `consultation`, `doctor`, `facility_id`, `order_id`, `patient`, `payment`, `quote`, `terms` |
| `booking.status` | Booking status | `billing`, `consultation`, `doctor`, `facility_id`, `order_id`, `patient`, `payment`, `quote`, `terms` |
| `doctor.search.results` | Search doctors | `data` |
| `doctor.slot.search.results` | Search slots | `doctor`, `provider_info`, `slots` |

### `abha.care_context_discover`

Discovery Request Received

- `abha_address` string
- `gender` string
- `hip_id` string
- `identifiers` []{type, value}
- `oid` string
- `partner_patient_id` string
- `patient_name` string
- `request_id` string
- `txn_id` string
- `year_of_birth` integer

### `abha.care_context_discover_link_init`

Link Initiation - Generate OTP

- `abha_address` string
- `hip_id` string
- `oid` string
- `partner_patient_id` string
- `patient` []{care_contexts, count, hi_type, id, ref_num}
- `request_id` string
- `txn_id` string

### `abha.consent_update`

Consent update

- `notification` {consentArtefacts, consentRequestId, status}
- `status` string
- `timestamp` string

### `abha.context_discover_link_confirm`

Link Confirmation - Verify OTP

- `hip_id` string
- `linkRefNumber` string
- `oid` string
- `partner_patient_id` string
- `request_id` string
- `token` string

### `abha.created`

ABHA Address Created

- `abha_address` string
- `abha_number` string
- `hip_id` string
- `oid` string
- `partner_patient_id` string

### `abha.hip_data_fetch`

Health Data Requested by HIU

- `abha_address` string
- `care_contexts` []string
- `hi_types` []string
- `hip_id` string
- `key_information` {crypto_alg, curve, dh_public_key, nonce}
- `oid` string
- `partner_patient_id` string
- `transaction_id` string

### `abha.hip_profile_share`

Scan and Share Token Received

- `abha_address` string
- `abha_number` string
- `hip_code` string
- `hip_id` string
- `oid` string
- `partner_patient_id` string
- `req_id` string

### `abha.hiu_data_push`

Health Data Received from HIP

- `abha_address` string
- `abha_number` string
- `consent_id` string
- `entries` []{care_context_id, checksum, content, media}
- `hip_id` string
- `key_information` {crypto_alg, curve, dh_public_key, nonce}
- `oid` string
- `partner_patient_id` string
- `transaction_id` string

### `abha.link_care_context`

Care Context Linking Status

- `abha_address` string
- `care_context_id` string
- `error` NoneType
- `hip_id` string
- `oid` string
- `partner_patient_id` string
- `retry_count` integer
- `status` string

### `abha.locker_created`

Locker Created

- `abha_address` string
- `auto_approval_id` string
- `hip_id` string
- `oid` string
- `partner_patient_id` string
- `subscription_id` string

### `abha.subscription_modified`

Subscription Updated

- `abha_address` string
- `hip_id` string
- `oid` string
- `partner_patient_id` string
- `subscription_details` {categories, period, status, subscription_id}

### `abha.subscription_notify`

New Care Context Linked (Subscription)

- `abha_address` string
- `abha_number` string
- `hip_id` string
- `oid` string
- `partner_patient_id` string
- `subscription_meta` {care_contexts, event_id, hip_id}

### `booking.cancel`

Booking cancel

- `status` string

### `booking.confirm`

Booking confirm

- `doctor_info` {education, experience, gender, id, languages, name}
- `order_id` string
- `pin` string
- `slot` {end, slot_id, start}
- `state` string

### `booking.init`

Booking init

- `billing` {address, email, name, phone}
- `consultation` {end_time, slot_id, start_time, type}
- `doctor` {doctor_gender, doctor_image, id, name, tags}
- `facility_id` string
- `order_id` string
- `patient` {dob, gender, id}
- `payment` {status, type}
- `quote` {breakup, total}
- `terms` []{long_desc, short_desc, terms_state, type}

### `booking.status`

Booking status

- `billing` {address, email, name, phone}
- `consultation` {end_time, slot_id, start_time, type}
- `doctor` {doctor_gender, doctor_image, id, name, tags}
- `facility_id` string
- `order_id` string
- `patient` {dob, gender, id}
- `payment` {status, type}
- `quote` {breakup, total}
- `terms` []{long_desc, short_desc, terms_state, type}

### `doctor.search.results`

Search doctors

- `data` []{doctor_name, education, facility_id, gender, hpr_id, image, languages, provider_address, provider_city, provider_country, provider_gps, provider_name, specialities, transaction_id, year_of_experience}

### `doctor.slot.search.results`

Search slots

- `doctor` {education, experience, gender, id, languages, name}
- `provider_info` {address, city, id, name, phone}
- `slots` []{currency, date, end_time, id, price, start_time, type}

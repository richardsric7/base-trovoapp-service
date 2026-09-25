# Stakeholder Portal Payloads In Swagger Order

Base URL for the local run was `http://localhost:8082/api/v1`.

These payloads use asset/org values from local AdminDB/TrovoWalletDB. Endpoints with `No JSON body` should be called with no request body. Values wrapped in `<...>` must come from a prior response in the same sequence.

Swagger order is display order, not always runnable order. For step-up workflows, use these runnable sequences:

- Fund release approve/execute: row 6 -> row 27 -> row 30 -> row 51 -> row 28 -> row 30 -> row 22 -> row 23.
- Fund release reject: row 6b -> row 52.
- Distribution authorize: row 11 -> row 12 -> row 29 -> row 30 -> row 43.
- Distribution reject: row 12b -> row 44.
- Due diligence approve: row 45 -> row 47 -> row 46.
- Due diligence reject: row 45 using the second asset -> row 48.

| # | Method | Endpoint | Auth | Depends on | Payload file | Exact JSON body | Purpose |
|---:|---|---|---|---|---|---|---|
| 1 | POST | `/stakeholder/admin/asset-assignments` | Admin dashboard token | Active stakeholder orgs and Wallet asset `NVFA` | `02-assignment-main.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","asset_code":"NVFA","trustee_org_id":"sp-verify-trustee-org","asset_manager_org_id":"sp-verify-manager-org","custodian_org_id":"sp-verify-custodian-org","status":"active"}` | Creates AdminDB asset visibility for trustee, manager, and custodian without mutating Wallet tokenization state. |
| 1b | POST | `/stakeholder/admin/asset-assignments` | Admin dashboard token | Active stakeholder orgs and Wallet asset `WSP` | `02-assignment-second.json` | `{"asset_id":"b57c9a37-050b-4431-a9be-09b0780ad5b9","asset_code":"WSP","trustee_org_id":"sp-verify-trustee-org","asset_manager_org_id":"sp-verify-manager-org","custodian_org_id":"sp-verify-custodian-org","status":"active"}` | Second assignment for reject-path tests. |
| 2 | GET | `/stakeholder/asset-manager/assets` | Manager token | Row 1 | None | No JSON body | Lists assets visible to the manager. |
| 3 | PUT | `/stakeholder/asset-manager/assets/{asset_id}` | Manager token | Row 1, `asset_id=87ed319d-dce7-4419-afd7-eea3fc1e4864` | `16-asset-operation.json` | `{"operational_status":"active","notes":"Operational status verified in staging","metadata":{"source":"stakeholder_portal_test_flow"}}` | Updates local portal operational state only. |
| 4 | GET | `/stakeholder/asset-manager/dashboard` | Manager token | Prior manager workflow data | None | No JSON body | Reads manager dashboard aggregates. |
| 5 | GET | `/stakeholder/asset-manager/fund-releases` | Manager token | Row 6 | None | No JSON body | Lists manager-created fund release requests. |
| 6 | POST | `/stakeholder/asset-manager/fund-releases` | Manager token | Row 1 | `05-fund-release-approve-path.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","amount":"1000.50000000","currency":"NGN","purpose":"Verification release - approve path","supporting_document_ids":[]}` | Creates a fund release request for trustee approval and custodian execution flow. |
| 6b | POST | `/stakeholder/asset-manager/fund-releases` | Manager token | Row 1b | `06-fund-release-reject-path.json` | `{"asset_id":"b57c9a37-050b-4431-a9be-09b0780ad5b9","amount":"750.25000000","currency":"NGN","purpose":"Verification release - reject path","supporting_document_ids":[]}` | Creates a separate fund release request for rejection testing. |
| 7 | GET | `/stakeholder/asset-manager/reports` | Manager token | Row 8 | None | No JSON body | Lists manager reports. |
| 8 | POST | `/stakeholder/asset-manager/reports` | Manager token | Row 1 | `17-report.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","asset_code":"NVFA","report_type":"operational","title":"Verification operational report","file_url":"https://example.com/staging/operational-report.pdf"}` | Creates URL-only report metadata. |
| 9 | POST | `/stakeholder/asset-manager/reports/{report_id}/submit` | Manager token | Row 8, use returned `data.id` | None | No JSON body | Submits the generated report to trustee workflow. |
| 10 | GET | `/stakeholder/asset-manager/revenue` | Manager token | Row 1 | None | No JSON body | Lists manager revenue records. |
| 11 | POST | `/stakeholder/asset-manager/revenue` | Manager token | Row 1 | `09-revenue.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","period_start":"2026-01-01","period_end":"2026-03-31","source":"rent","amount":"2500.00000000","currency":"NGN","collected_at":"2026-03-31"}` | Records asset revenue with decimal-safe amount. |
| 12 | POST | `/stakeholder/asset-manager/revenue/submit-distribution` | Manager token | Row 11, use returned revenue ID | `35-distribution-authorize-path-with-revenue.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","revenue_record_ids":["<REVENUE_RECORD_ID_FROM_ROW_11>"],"amount":"2500.00000000","currency":"NGN","source":"rent","scheduled_date":"2026-07-15"}` | Submits recorded revenue for trustee distribution authorization. |
| 12b | POST | `/stakeholder/asset-manager/revenue/submit-distribution` | Manager token | Row 1b | `11-distribution-reject-path.json` | `{"asset_id":"b57c9a37-050b-4431-a9be-09b0780ad5b9","revenue_record_ids":[],"amount":"1000.00000000","currency":"NGN","source":"rent","scheduled_date":"2026-07-20"}` | Creates separate distribution proposal for rejection testing. |
| 13 | POST | `/stakeholder/asset-manager/valuations` | Manager token | Row 1 | `12-valuation.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","valuation":"12500000.00000000","currency":"NGN","methodology":"Independent comparable market analysis","valuation_date":"2026-06-30"}` | Submits asset valuation with decimal-safe amount. |
| 14 | GET | `/stakeholder/asset-manager/valuations/{asset_id}` | Manager token | Row 1 | None | No JSON body | Lists valuation history for a visible manager asset. |
| 15 | POST | `/stakeholder/asset-manager/valuations/{id}/request-independent` | Manager token | Row 13, use returned `data.id` | None | No JSON body | Requests independent valuation review. |
| 16 | GET | `/stakeholder/custodian/accounts` | Custodian token | Seeded `segregated_accounts` row | None | No JSON body | Lists custodian-owned segregated accounts. |
| 17 | GET | `/stakeholder/custodian/accounts/{account_id}` | Custodian token | `account_id=sp-verify-account-main` | None | No JSON body | Reads one custodian-owned segregated account. |
| 18 | GET | `/stakeholder/custodian/compliance` | Custodian token | Compliance item created through `POST /stakeholder/admin/custodian-compliance` | None | No JSON body | Lists custodian compliance items. |
| 19 | PUT | `/stakeholder/custodian/compliance/{item_id}` | Custodian token | `item_id=sp-verify-compliance-main` | `13-compliance-complete.json` | `{"status":"complete"}` | Marks a compliance item complete with audit. |
| 20 | GET | `/stakeholder/custodian/dashboard` | Custodian token | Prior custodian data | None | No JSON body | Reads custodian dashboard aggregates. |
| 21 | GET | `/stakeholder/custodian/fund-releases` | Custodian token | Fund release is trustee approved / execution pending | None | No JSON body | Lists custodian execution queue. |
| 22 | POST | `/stakeholder/custodian/fund-releases/{request_id}/execute` | Custodian token | Verified challenge from rows 28-29 for action `fund_release.execute` | `26-fr-execute.json` | `{"challenge_id":"<FR_EXECUTE_CHALLENGE_ID_FROM_RESPONSE>","execution_reference":"staging-exec-001","status":"processing"}` | Executes an approved fund release; requires step-up. |
| 23 | PUT | `/stakeholder/custodian/fund-releases/{request_id}/status` | Custodian token | Row 22 already moved request into processing | `27-fr-status-complete.json` | `{"execution_reference":"staging-exec-001","status":"completed"}` | Updates an already-executed release to terminal status. |
| 24 | GET | `/stakeholder/shared/assets` | Any stakeholder token | Row 1 | None | No JSON body | Lists visible assets for the caller role. |
| 25 | GET | `/stakeholder/shared/assets/{asset_id}` | Any stakeholder token | Row 1 and visible asset | None | No JSON body | Gets visible asset detail. |
| 26 | GET | `/stakeholder/shared/audit-trail` | Any stakeholder token | Prior audited actions | None | No JSON body | Lists org-scoped audit evidence. |
| 27 | POST | `/stakeholder/shared/authorizations` | Trustee token | Fund release request from row 6 | `20-fr-approve-challenge.json` | `{"action":"fund_release.approve","entity_type":"fund_release_request","entity_id":"<FUND_RELEASE_APPROVE_ID_FROM_CREATE_RESPONSE>"}` | Creates trustee step-up challenge for fund release approval. |
| 28 | POST | `/stakeholder/shared/authorizations` | Custodian token | Fund release is execution pending after row 51 succeeds | `24-fr-execute-challenge.json` | `{"action":"fund_release.execute","entity_type":"fund_release_request","entity_id":"<FUND_RELEASE_APPROVE_ID_FROM_CREATE_RESPONSE>"}` | Creates custodian step-up challenge for execution. |
| 29 | POST | `/stakeholder/shared/authorizations` | Trustee token | Distribution proposal from row 12 | `39-distribution-challenge.json` | `{"action":"distribution.authorize","entity_type":"distribution","entity_id":"<DISTRIBUTION_AUTHORIZE_ID_FROM_SUBMIT_DISTRIBUTION_RESPONSE>"}` | Creates trustee step-up challenge for distribution authorization. |
| 30 | POST | `/stakeholder/shared/authorizations/{challenge_id}/verify` | Same token that created challenge | Use `challenge_id` and `auth_id` from row 27, 28, or 29 | `21-verify-authorization.json` | `{"auth_id":"<AUTH_ID_FROM_CREATE_CHALLENGE_RESPONSE>"}` | Verifies wallet approval for the challenge. |
| 31 | GET | `/stakeholder/shared/documents` | Any stakeholder token | Row 32 | None | No JSON body | Lists documents visible by role and asset. |
| 32 | POST | `/stakeholder/shared/documents` | Manager token | Row 1 | `14-document-asset-scoped.json` | `{"asset_id":"87ed319d-dce7-4419-afd7-eea3fc1e4864","asset_code":"NVFA","category":"legal","title":"Verification document","file_url":"https://example.com/staging/stakeholder-verification.pdf","mime_type":"application/pdf","version":1,"access_roles":["trustee","asset_manager"]}` | Creates URL-only document metadata with role/asset visibility. |
| 32b | POST | `/stakeholder/shared/documents` | Manager token | Regression-only check | `15-document-unscoped-regression.json` | `{"category":"internal","title":"Unscoped document leak regression","file_url":"https://example.com/staging/unscoped.pdf","mime_type":"application/pdf","version":1,"access_roles":["asset_manager"]}` | Regression payload proving unscoped docs do not leak across orgs. |
| 33 | GET | `/stakeholder/shared/documents/{doc_id}` | Any visible stakeholder token | Row 32, use returned `data.id` | None | No JSON body | Reads one visible document. |
| 34 | GET | `/stakeholder/shared/health` | None | None | None | No JSON body | Checks stakeholder portal route availability. |
| 35 | GET | `/stakeholder/shared/notifications` | Any stakeholder token | Prior workflow notification | None | No JSON body | Lists recipient-scoped notifications. |
| 36 | PUT | `/stakeholder/shared/notifications/{id}/read` | Notification recipient token | Row 35, use returned notification ID | None | No JSON body | Marks one owned notification as read. |
| 37 | GET | `/stakeholder/shared/profile` | Any stakeholder token | Active org/member JWT | None | No JSON body | Reads stakeholder profile and role mapping. |
| 38 | PUT | `/stakeholder/shared/profile` | Any stakeholder token | Active org/member JWT | `03-profile-update.json` | `{"first_name":"Portal","last_name":"Verifier"}` | Updates safe member profile fields only. |
| 39 | PUT | `/stakeholder/shared/profile/notifications` | Any stakeholder token | Active org/member JWT | `04-notification-preferences.json` | `{"email_enabled":false,"in_app_enabled":true}` | Updates notification preferences. |
| 40 | GET | `/stakeholder/trustee/dashboard` | Trustee token | Prior trustee workflow data | None | No JSON body | Reads trustee dashboard aggregates. |
| 41 | GET | `/stakeholder/trustee/distributions` | Trustee token | Row 12 | None | No JSON body | Lists trustee distribution queue. |
| 42 | GET | `/stakeholder/trustee/distributions/history` | Trustee token | Rejected or authorized distribution | None | No JSON body | Lists terminal distribution history. |
| 43 | POST | `/stakeholder/trustee/distributions/{dist_id}/authorize` | Trustee token | Row 30 verified distribution challenge | `41-distribution-authorize.json` | `{"challenge_id":"<DISTRIBUTION_CHALLENGE_ID_FROM_RESPONSE>"}` | Authorizes a distribution with consumed step-up. |
| 44 | POST | `/stakeholder/trustee/distributions/{dist_id}/reject` | Trustee token | Row 12b, use returned `data.id` | `07-reject.json` | `{"reason":"Verification rejection reason from test flow"}` | Rejects a proposed distribution with reason/audit. |
| 45 | GET | `/stakeholder/trustee/due-diligence/{asset_id}` | Trustee token | Row 1 | None | No JSON body | Gets or initializes trustee checklist. |
| 46 | POST | `/stakeholder/trustee/due-diligence/{asset_id}/approve` | Trustee token | Row 47 should complete at least one item first | None | No JSON body | Approves local AdminDB due diligence. |
| 47 | PUT | `/stakeholder/trustee/due-diligence/{asset_id}/items/{item_id}` | Trustee token | Row 45, use returned item ID | `08-due-diligence-item-complete.json` | `{"status":"complete","notes":"Verified against source documents in staging"}` | Marks a checklist item complete. |
| 48 | POST | `/stakeholder/trustee/due-diligence/{asset_id}/reject` | Trustee token | Row 45 for second asset | `07-reject.json` | `{"reason":"Verification rejection reason from test flow"}` | Rejects local due diligence with reason. |
| 49 | GET | `/stakeholder/trustee/fund-releases` | Trustee token | Row 6 | None | No JSON body | Lists trustee fund release approval queue. |
| 50 | GET | `/stakeholder/trustee/fund-releases/{request_id}` | Trustee token | Row 6, use returned request ID | None | No JSON body | Reads fund release detail for trustee decision. |
| 51 | POST | `/stakeholder/trustee/fund-releases/{request_id}/approve` | Trustee token | Row 30 verified approval challenge | `22-fr-approve.json` | `{"challenge_id":"<FR_APPROVE_CHALLENGE_ID_FROM_RESPONSE>"}` | Approves a fund release with consumed step-up. |
| 52 | POST | `/stakeholder/trustee/fund-releases/{request_id}/reject` | Trustee token | Row 6b, use returned request ID | `07-reject.json` | `{"reason":"Verification rejection reason from test flow"}` | Rejects a fund release with reason/audit. |

## Negative/Regression Payloads

| Payload file | Exact JSON body | Use |
|---|---|---|
| `N5-invalid-challenge.json` | `{"challenge_id":"not-a-valid-challenge"}` | Proves trustee fund release approval and distribution authorization reject missing/invalid step-up. |
| `N8-status-bypass.json` | `{"status":"completed","execution_reference":"bypass-test"}` | Proves custodian status endpoint cannot bypass the `/execute` step-up path from `execution_pending`. |

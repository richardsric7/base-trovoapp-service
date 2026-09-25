# Stakeholder Portal Test Verification Flow

Date: 2026-06-30

Purpose: verify the Stakeholder Portal implementation with DB-derived payloads, ordered endpoint calls, saved responses, and SQL proof. Run this against local or staging databases only.

Do not mark the portal production-ready unless the two P1 regression checks in this document pass:

- Custodian status update must not execute/complete a release without wallet step-up.
- Unscoped documents must not leak across organizations with the same role.

Route count note: the current code registers 50 unique stakeholder portal endpoints. This runbook has 61 positive verification calls because some endpoints must be exercised more than once for separate roles, happy paths, rejection paths, and terminal state updates.

## Read This First: Where The JSON Is

The payload path is relative to this repo:

```text
/Users/toluwase/GolandProjects/admin-dashboard-api/payloads/stakeholder-portal
```

That folder is created when you run:

```bash
mkdir -p evidence/stakeholder-portal payloads/stakeholder-portal
```

The JSON files are generated from DB-derived environment variables in Step 2. They are not random sample files. The exact JSON bodies are also listed in the payload catalog below, so you do not need to hunt for them.

## Tools

Required local tools:

```bash
psql --version
jq --version
curl --version
```

Set base variables:

```bash
mkdir -p evidence/stakeholder-portal payloads/stakeholder-portal

export BASE_URL="http://localhost:8082/api/v1"
export ADMIN_CONNECTION_STRING="postgres://..."
export WALLET_DB_CONNECTION_STRING="postgres://..."

export TROVO_ADMIN_TOKEN="..."
export TRUSTEE_TOKEN="..."
export MANAGER_TOKEN="..."
export CUSTODIAN_TOKEN="..."

# Optional negative-test tokens from different orgs with same role.
export OTHER_TRUSTEE_TOKEN="..."
export OTHER_MANAGER_TOKEN="..."
export OTHER_CUSTODIAN_TOKEN="..."
```

Prove response payloads after each call:

```bash
cat evidence/stakeholder-portal/03-01-health.status
jq . evidence/stakeholder-portal/03-01-health.json
```

Use this route-count proof before API testing:

```bash
rg -n 'stakeholder.*(GET|POST|PUT)|\.(GET|POST|PUT)\(' internal/components/stakeholder/controllers/main.go \
  | tee evidence/stakeholder-portal/00-route-registration-proof.txt
```

Use this helper for every API call. It stores headers, body, and status separately so IDs can be extracted from JSON reliably.

```bash
api_json() {
  local method="$1"
  local path="$2"
  local token="$3"
  local payload="$4"
  local out="$5"

  if [ "$payload" = "-" ]; then
    curl -sS -X "$method" "$BASE_URL$path" \
      -H "Authorization: Bearer $token" \
      -D "$out.headers" \
      -o "$out.json" \
      -w "%{http_code}" | tee "$out.status"
  else
    curl -sS -X "$method" "$BASE_URL$path" \
      -H "Authorization: Bearer $token" \
      -H "Content-Type: application/json" \
      --data @"$payload" \
      -D "$out.headers" \
      -o "$out.json" \
      -w "%{http_code}" | tee "$out.status"
  fi
}

api_public() {
  local method="$1"
  local path="$2"
  local out="$3"

  curl -sS -X "$method" "$BASE_URL$path" \
    -D "$out.headers" \
    -o "$out.json" \
    -w "%{http_code}" | tee "$out.status"
}
```

## Step 0: DB-Derived Fixture Discovery

The payloads below must be generated from DB values, not guessed IDs. These commands select active, linked organizations and active members.

```bash
export TRUSTEE_ORG_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organizations where status='ACTIVE' and stakeholder_type='trustees' and stakeholder_id is not null order by created_at desc limit 1")"
export MANAGER_ORG_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organizations where status='ACTIVE' and stakeholder_type='asset_manager' and stakeholder_id is not null order by created_at desc limit 1")"
export CUSTODIAN_ORG_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organizations where status='ACTIVE' and stakeholder_type='approved_asset_custodian' and stakeholder_id is not null order by created_at desc limit 1")"

export TRUSTEE_STAKEHOLDER_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select stakeholder_id from organizations where id='$TRUSTEE_ORG_ID'")"
export MANAGER_STAKEHOLDER_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select stakeholder_id from organizations where id='$MANAGER_ORG_ID'")"
export CUSTODIAN_STAKEHOLDER_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select stakeholder_id from organizations where id='$CUSTODIAN_ORG_ID'")"

export TRUSTEE_MEMBER_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organization_members where organization_id='$TRUSTEE_ORG_ID' and status='ACTIVE' order by created_at desc limit 1")"
export MANAGER_MEMBER_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organization_members where organization_id='$MANAGER_ORG_ID' and status='ACTIVE' order by created_at desc limit 1")"
export CUSTODIAN_MEMBER_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organization_members where organization_id='$CUSTODIAN_ORG_ID' and status='ACTIVE' order by created_at desc limit 1")"

export ASSET_ID="$(psql "$WALLET_DB_CONNECTION_STRING" -Atc "select id from tokenized_assets where asset_manager_id=$MANAGER_STAKEHOLDER_ID and approved_asset_custodian_id=$CUSTODIAN_STAKEHOLDER_ID order by created_at desc limit 1")"
export ASSET_CODE="$(psql "$WALLET_DB_CONNECTION_STRING" -Atc "select asset_code from tokenized_assets where id='$ASSET_ID'")"

# Second asset is needed for positive reject-path tests so the approved happy path does not consume the same entity.
export SECOND_ASSET_ID="$(psql "$WALLET_DB_CONNECTION_STRING" -Atc "select id from tokenized_assets where id <> '$ASSET_ID' and asset_manager_id=$MANAGER_STAKEHOLDER_ID and approved_asset_custodian_id=$CUSTODIAN_STAKEHOLDER_ID order by created_at desc limit 1")"
export SECOND_ASSET_CODE="$(psql "$WALLET_DB_CONNECTION_STRING" -Atc "select asset_code from tokenized_assets where id='$SECOND_ASSET_ID'")"

# Required for negative leakage tests. Tokens for these orgs must be exported above.
export OTHER_TRUSTEE_ORG_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organizations where status='ACTIVE' and stakeholder_type='trustees' and stakeholder_id is not null and id <> '$TRUSTEE_ORG_ID' order by created_at desc limit 1")"
export OTHER_MANAGER_ORG_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organizations where status='ACTIVE' and stakeholder_type='asset_manager' and stakeholder_id is not null and id <> '$MANAGER_ORG_ID' order by created_at desc limit 1")"
export OTHER_CUSTODIAN_ORG_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from organizations where status='ACTIVE' and stakeholder_type='approved_asset_custodian' and stakeholder_id is not null and id <> '$CUSTODIAN_ORG_ID' order by created_at desc limit 1")"
```

If `SECOND_ASSET_ID` is empty, create or select another staging tokenized asset with the same manager/custodian linkage before running positive reject tests. Do not fabricate an asset ID.

If any `OTHER_*_ORG_ID` is empty, skip only the matching same-role leakage test or create a real staging organization/member/token for that role. Do not reuse the primary org token for negative leakage tests.

Save fixture proof:

```bash
psql "$ADMIN_CONNECTION_STRING" -c "
select id,name,status,stakeholder_id,stakeholder_type,type
from organizations
where id in (
  '$TRUSTEE_ORG_ID','$MANAGER_ORG_ID','$CUSTODIAN_ORG_ID',
  nullif('$OTHER_TRUSTEE_ORG_ID',''),nullif('$OTHER_MANAGER_ORG_ID',''),nullif('$OTHER_CUSTODIAN_ORG_ID','')
);
" | tee evidence/stakeholder-portal/00-admin-org-fixtures.txt

psql "$ADMIN_CONNECTION_STRING" -c "
select id,email,status,role,organization_id,is_wallet_linked,trovo_wallet_username
from organization_members
where id in ('$TRUSTEE_MEMBER_ID','$MANAGER_MEMBER_ID','$CUSTODIAN_MEMBER_ID');
" | tee evidence/stakeholder-portal/00-admin-member-fixtures.txt

psql "$WALLET_DB_CONNECTION_STRING" -c "
select id,asset_code,asset_name,asset_manager_id,approved_asset_custodian_id,asset_tokenization_status
from tokenized_assets
where id in ('$ASSET_ID','$SECOND_ASSET_ID');
" | tee evidence/stakeholder-portal/00-wallet-asset-fixtures.txt
```

Expected proof:

- Organizations are `ACTIVE`.
- Members are `ACTIVE`.
- Stakeholder IDs are not null.
- Wallet asset `asset_manager_id` equals `$MANAGER_STAKEHOLDER_ID`.
- Wallet asset `approved_asset_custodian_id` equals `$CUSTODIAN_STAKEHOLDER_ID`.
- `organizations.type` is captured only as evidence; portal auth must be based on `stakeholder_type`.

The profile endpoint later proves each token maps back to these DB rows because `.data.member.id`, `.data.organization.id`, and `.data.stakeholder.dashboard_role` must match the fixture.

## Step 1: Migration and Boundary Proof

Apply to local/staging AdminDB only:

```bash
psql "$ADMIN_CONNECTION_STRING" -f migrations/20260629_create_stakeholder_portal_tables.sql \
  | tee evidence/stakeholder-portal/01-migration-output.txt
```

| Proof | Command | Expected |
|---|---|---|
| All portal tables exist in AdminDB | `psql "$ADMIN_CONNECTION_STRING" -c "select table_name from information_schema.tables where table_schema='public' and table_name in ('stakeholder_asset_assignments','stakeholder_audit_logs','stakeholder_notifications','stakeholder_authorization_challenges','fund_release_requests','due_diligence_checklists','due_diligence_items','revenue_records','distributions','asset_valuations','segregated_accounts','compliance_items','stakeholder_documents','stakeholder_notification_preferences','stakeholder_asset_operations','stakeholder_reports') order by table_name;" \| tee evidence/stakeholder-portal/01-admin-tables.txt` | 16 table names returned. |
| Money precision | `psql "$ADMIN_CONNECTION_STRING" -c "select table_name,column_name,numeric_precision,numeric_scale from information_schema.columns where table_name in ('fund_release_requests','revenue_records','distributions','asset_valuations','segregated_accounts') and column_name in ('amount','valuation','balance') order by table_name,column_name;" \| tee evidence/stakeholder-portal/01-money-precision.txt` | Every row has precision `30`, scale `8`. |
| No portal workflow tables in Wallet DB | `psql "$WALLET_DB_CONNECTION_STRING" -c "select table_name from information_schema.tables where table_schema='public' and (table_name like 'stakeholder_%' or table_name in ('fund_release_requests','due_diligence_checklists','distributions','revenue_records','asset_valuations','segregated_accounts','compliance_items'));" \| tee evidence/stakeholder-portal/01-wallet-no-portal-tables.txt` | Zero rows. |

## Step 2: Generate Correct JSON Payloads From DB Fixtures

These payloads are valid for the current request structs in `internal/components/stakeholder/models/requests.go`.

```bash
jq -n \
  --arg asset_id "$ASSET_ID" \
  --arg asset_code "$ASSET_CODE" \
  --arg trustee_org_id "$TRUSTEE_ORG_ID" \
  --arg manager_org_id "$MANAGER_ORG_ID" \
  --arg custodian_org_id "$CUSTODIAN_ORG_ID" \
  '{asset_id:$asset_id,asset_code:$asset_code,trustee_org_id:$trustee_org_id,asset_manager_org_id:$manager_org_id,custodian_org_id:$custodian_org_id,status:"active"}' \
  > payloads/stakeholder-portal/02-assignment-main.json

jq -n \
  --arg asset_id "$SECOND_ASSET_ID" \
  --arg asset_code "$SECOND_ASSET_CODE" \
  --arg trustee_org_id "$TRUSTEE_ORG_ID" \
  --arg manager_org_id "$MANAGER_ORG_ID" \
  --arg custodian_org_id "$CUSTODIAN_ORG_ID" \
  '{asset_id:$asset_id,asset_code:$asset_code,trustee_org_id:$trustee_org_id,asset_manager_org_id:$manager_org_id,custodian_org_id:$custodian_org_id,status:"active"}' \
  > payloads/stakeholder-portal/02-assignment-second.json

jq -n '{first_name:"Portal",last_name:"Verifier"}' \
  > payloads/stakeholder-portal/03-profile-update.json

jq -n '{email_enabled:false,in_app_enabled:true}' \
  > payloads/stakeholder-portal/04-notification-preferences.json

jq -n \
  --arg asset_id "$ASSET_ID" \
  '{asset_id:$asset_id,amount:"1000.50000000",currency:"NGN",purpose:"Verification release - approve path",supporting_document_ids:[]}' \
  > payloads/stakeholder-portal/05-fund-release-approve-path.json

jq -n \
  --arg asset_id "$SECOND_ASSET_ID" \
  '{asset_id:$asset_id,amount:"750.25000000",currency:"NGN",purpose:"Verification release - reject path",supporting_document_ids:[]}' \
  > payloads/stakeholder-portal/06-fund-release-reject-path.json

jq -n '{reason:"Verification rejection reason from test flow"}' \
  > payloads/stakeholder-portal/07-reject.json

jq -n '{status:"complete",notes:"Verified against source documents in staging"}' \
  > payloads/stakeholder-portal/08-due-diligence-item-complete.json

jq -n \
  --arg asset_id "$ASSET_ID" \
  '{asset_id:$asset_id,period_start:"2026-01-01",period_end:"2026-03-31",source:"rent",amount:"2500.00000000",currency:"NGN",collected_at:"2026-03-31"}' \
  > payloads/stakeholder-portal/09-revenue.json

jq -n \
  --arg asset_id "$ASSET_ID" \
  '{asset_id:$asset_id,revenue_record_ids:[],amount:"2500.00000000",currency:"NGN",source:"rent",scheduled_date:"2026-07-15"}' \
  > payloads/stakeholder-portal/10-distribution-authorize-path.json

jq -n \
  --arg asset_id "$SECOND_ASSET_ID" \
  '{asset_id:$asset_id,revenue_record_ids:[],amount:"1000.00000000",currency:"NGN",source:"rent",scheduled_date:"2026-07-20"}' \
  > payloads/stakeholder-portal/11-distribution-reject-path.json

jq -n \
  --arg asset_id "$ASSET_ID" \
  '{asset_id:$asset_id,valuation:"12500000.00000000",currency:"NGN",methodology:"Independent comparable market analysis",valuation_date:"2026-06-30"}' \
  > payloads/stakeholder-portal/12-valuation.json

jq -n '{status:"complete"}' \
  > payloads/stakeholder-portal/13-compliance-complete.json

jq -n \
  --arg asset_id "$ASSET_ID" \
  --arg asset_code "$ASSET_CODE" \
  '{asset_id:$asset_id,asset_code:$asset_code,category:"legal",title:"Verification document",file_url:"https://example.com/staging/stakeholder-verification.pdf",mime_type:"application/pdf",version:1,access_roles:["trustee","asset_manager"]}' \
  > payloads/stakeholder-portal/14-document-asset-scoped.json

jq -n \
  '{category:"internal",title:"Unscoped document leak regression",file_url:"https://example.com/staging/unscoped.pdf",mime_type:"application/pdf",version:1,access_roles:["asset_manager"]}' \
  > payloads/stakeholder-portal/15-document-unscoped-regression.json

jq -n '{operational_status:"active",notes:"Operational status verified in staging",metadata:{source:"stakeholder_portal_test_flow"}}' \
  > payloads/stakeholder-portal/16-asset-operation.json

jq -n \
  --arg asset_id "$ASSET_ID" \
  --arg asset_code "$ASSET_CODE" \
  '{asset_id:$asset_id,asset_code:$asset_code,report_type:"operational",title:"Verification operational report",file_url:"https://example.com/staging/operational-report.pdf"}' \
  > payloads/stakeholder-portal/17-report.json

for f in payloads/stakeholder-portal/*.json; do
  echo "===== $f"
  jq . "$f"
done | tee evidence/stakeholder-portal/02-payload-proof.txt
```

Payload proof requirements:

- Every `asset_id` and `asset_code` must equal the values selected from Wallet DB in Step 0.
- Every organization ID in assignment payloads must equal an active AdminDB organization selected in Step 0.
- Amount strings must keep decimal precision and must be valid for `numeric(30,8)`.
- No payload uses `organizations.type` for authorization.

## Step 2A: Logical Sequence Breakpoints

Use these breakpoints while running the matrix:

| Phase | Rows | What to run before continuing |
|---|---:|---|
| Setup and visibility | 1-13 | Run Step 0 through Step 2 first. Rows 2 and 3 create the AdminDB asset assignments required by later asset-scoped endpoints. |
| Fund-release create/reject | 14-19 | After rows 14 and 15, export `FUND_RELEASE_APPROVE_ID` and `FUND_RELEASE_REJECT_ID` from Step 4. |
| Fund-release approval | 20-22 | Generate `20-fr-approve-challenge.json`, create and verify the wallet challenge, then generate `22-fr-approve.json` from the returned challenge ID. |
| Fund-release execution | 23-27 | Generate `24-fr-execute-challenge.json`, create and verify the wallet challenge, then generate `26-fr-execute.json`. Row 27 is valid only after row 26 has moved the release to `processing`. |
| Due diligence | 28-32 | Export `DD_ITEM_ID` after row 28 and `SECOND_DD_ITEM_ID` after row 31. |
| Revenue and distributions | 33-42 | Export `REVENUE_RECORD_ID`, generate `35-distribution-authorize-path-with-revenue.json`, then export distribution IDs and challenge IDs as shown in Step 4. |
| Valuation, custodian, documents, notifications, dashboards, reports | 43-61 | Seed account/compliance rows from Step 5 before rows 46 and 48. Export `VALUATION_ID`, `DOCUMENT_ID`, `NOTIFICATION_ID`, and `REPORT_ID` as each create/list response is produced. |

## Step 2B: Exact JSON Payload Catalog

Use this table as the JSON column for the endpoint matrix. Values like `$ASSET_ID` are not placeholders to invent; they must come from Step 0 SQL or from a previous endpoint response.

| Payload file | Exact JSON body | Used by | Value proof |
|---|---|---|---|
| `02-assignment-main.json` | `{"asset_id":"$ASSET_ID","asset_code":"$ASSET_CODE","trustee_org_id":"$TRUSTEE_ORG_ID","asset_manager_org_id":"$MANAGER_ORG_ID","custodian_org_id":"$CUSTODIAN_ORG_ID","status":"active"}` | Row 2 | `$ASSET_ID/$ASSET_CODE` from Wallet DB; org IDs from AdminDB active stakeholder orgs. |
| `02-assignment-second.json` | `{"asset_id":"$SECOND_ASSET_ID","asset_code":"$SECOND_ASSET_CODE","trustee_org_id":"$TRUSTEE_ORG_ID","asset_manager_org_id":"$MANAGER_ORG_ID","custodian_org_id":"$CUSTODIAN_ORG_ID","status":"active"}` | Row 3 | Second Wallet DB asset; same active AdminDB orgs. |
| `03-profile-update.json` | `{"first_name":"Portal","last_name":"Verifier"}` | Row 7 | Safe member profile fields only. |
| `04-notification-preferences.json` | `{"email_enabled":false,"in_app_enabled":true}` | Row 8 | Member preference booleans only. |
| `05-fund-release-approve-path.json` | `{"asset_id":"$ASSET_ID","amount":"1000.50000000","currency":"NGN","purpose":"Verification release - approve path","supporting_document_ids":[]}` | Rows 14, N5, N8 | `$ASSET_ID` from Wallet DB and assigned in AdminDB. |
| `06-fund-release-reject-path.json` | `{"asset_id":"$SECOND_ASSET_ID","amount":"750.25000000","currency":"NGN","purpose":"Verification release - reject path","supporting_document_ids":[]}` | Row 15 | `$SECOND_ASSET_ID` from Wallet DB and assigned in AdminDB. |
| `07-reject.json` | `{"reason":"Verification rejection reason from test flow"}` | Rows 19, 32, 38 | Required rejection reason. |
| `08-due-diligence-item-complete.json` | `{"status":"complete","notes":"Verified against source documents in staging"}` | Row 29 | Status allowed by due diligence item state machine. |
| `09-revenue.json` | `{"asset_id":"$ASSET_ID","period_start":"2026-01-01","period_end":"2026-03-31","source":"rent","amount":"2500.00000000","currency":"NGN","collected_at":"2026-03-31"}` | Row 34 | Manager-visible assigned asset. |
| `10-distribution-authorize-path.json` | `{"asset_id":"$ASSET_ID","revenue_record_ids":[],"amount":"2500.00000000","currency":"NGN","source":"rent","scheduled_date":"2026-07-15"}` | Source for row 35 | `$ASSET_ID` from assigned Wallet DB asset. |
| `35-distribution-authorize-path-with-revenue.json` | `{"asset_id":"$ASSET_ID","revenue_record_ids":["$REVENUE_RECORD_ID"],"amount":"2500.00000000","currency":"NGN","source":"rent","scheduled_date":"2026-07-15"}` | Row 35 | `$REVENUE_RECORD_ID` from row 34 response. |
| `11-distribution-reject-path.json` | `{"asset_id":"$SECOND_ASSET_ID","revenue_record_ids":[],"amount":"1000.00000000","currency":"NGN","source":"rent","scheduled_date":"2026-07-20"}` | Rows 36, N6 | Second assigned Wallet DB asset. |
| `12-valuation.json` | `{"asset_id":"$ASSET_ID","valuation":"12500000.00000000","currency":"NGN","methodology":"Independent comparable market analysis","valuation_date":"2026-06-30"}` | Row 44 | Manager-visible assigned asset. |
| `13-compliance-complete.json` | `{"status":"complete"}` | Row 49 | Status allowed by compliance workflow. |
| `14-document-asset-scoped.json` | `{"asset_id":"$ASSET_ID","asset_code":"$ASSET_CODE","category":"legal","title":"Verification document","file_url":"https://example.com/staging/stakeholder-verification.pdf","mime_type":"application/pdf","version":1,"access_roles":["trustee","asset_manager"]}` | Row 50 | Asset visible to uploader; access roles supported by portal. |
| `15-document-unscoped-regression.json` | `{"category":"internal","title":"Unscoped document leak regression","file_url":"https://example.com/staging/unscoped.pdf","mime_type":"application/pdf","version":1,"access_roles":["asset_manager"]}` | N9 | No `asset_id`; must remain scoped to uploader org only. |
| `16-asset-operation.json` | `{"operational_status":"active","notes":"Operational status verified in staging","metadata":{"source":"stakeholder_portal_test_flow"}}` | Row 13 | Manager-visible assigned asset. |
| `17-report.json` | `{"asset_id":"$ASSET_ID","asset_code":"$ASSET_CODE","report_type":"operational","title":"Verification operational report","file_url":"https://example.com/staging/operational-report.pdf"}` | Row 59 | Manager-visible assigned asset. |
| `20-fr-approve-challenge.json` | `{"action":"fund_release.approve","entity_type":"fund_release_request","entity_id":"$FUND_RELEASE_APPROVE_ID"}` | Row 20 | `$FUND_RELEASE_APPROVE_ID` from row 14 response. |
| `22-fr-approve.json` | `{"challenge_id":"$FR_APPROVE_CHALLENGE_ID"}` | Row 22 | `$FR_APPROVE_CHALLENGE_ID` from row 20 response after wallet verification in row 21. |
| `24-fr-execute-challenge.json` | `{"action":"fund_release.execute","entity_type":"fund_release_request","entity_id":"$FUND_RELEASE_APPROVE_ID"}` | Row 24 | Approved fund release from row 22. |
| `26-fr-execute.json` | `{"challenge_id":"$FR_EXECUTE_CHALLENGE_ID","status":"processing","execution_reference":"staging-exec-001"}` | Row 26 | `$FR_EXECUTE_CHALLENGE_ID` from row 24 response after wallet verification in row 25. |
| `27-fr-status-complete.json` | `{"status":"completed","execution_reference":"staging-exec-001"}` | Row 27 | Only valid after row 26 has moved release to `processing`. |
| `39-distribution-challenge.json` | `{"action":"distribution.authorize","entity_type":"distribution","entity_id":"$DISTRIBUTION_AUTHORIZE_ID"}` | Row 39 | `$DISTRIBUTION_AUTHORIZE_ID` from row 35 response. |
| `41-distribution-authorize.json` | `{"challenge_id":"$DISTRIBUTION_CHALLENGE_ID"}` | Row 41 | `$DISTRIBUTION_CHALLENGE_ID` from row 39 response after wallet verification in row 40. |
| `N5-invalid-challenge.json` | `{"challenge_id":"not-a-valid-challenge"}` | N5 | Must fail; proves approval cannot proceed without a real verified challenge. |
| `N8-status-bypass.json` | `{"status":"completed","execution_reference":"bypass-test"}` | N8 | Must fail before custodian `/execute`; proves status endpoint cannot bypass step-up. |
| `N6-wrong-action-challenge.json` | `{"action":"fund_release.approve","entity_type":"fund_release_request","entity_id":"$N5_FUND_RELEASE_ID"}` | N6 setup | Intentionally wrong action/entity for distribution authorization. |
| `N6-wrong-action-authorize.json` | `{"challenge_id":"$N6_WRONG_ACTION_CHALLENGE_ID"}` | N6 | Must fail when used against a distribution. |

## Step 3: Endpoint Verification Matrix

Run the rows in order. If a row depends on a prior row, do not skip the prior row unless the dependent ID already exists and is proven by SQL.

| Seq | Endpoint | Payload file and command | Depends On | Endpoint purpose and task goal | Admin-dashboard product fit | Expected response proof | DB proof |
|---:|---|---|---|---|---|---|---|
| 1 | `GET /stakeholder/shared/health` | `api_public GET /stakeholder/shared/health evidence/stakeholder-portal/03-01-health` | None | Proves stakeholder component is registered. | Gives operators a quick portal availability check. | `200`, body has `status=ok`. | None. |
| 2 | `POST /stakeholder/admin/asset-assignments` | `api_json POST /stakeholder/admin/asset-assignments "$TROVO_ADMIN_TOKEN" payloads/stakeholder-portal/02-assignment-main.json evidence/stakeholder-portal/03-02-assignment-main` | Step 0 | Creates trustee/manager/custodian visibility source. | Lets Trovo admins configure portal asset access without mutating tokenization. | `200`, `data.asset_id=$ASSET_ID`. | Query `stakeholder_asset_assignments` for `$ASSET_ID`. |
| 3 | `POST /stakeholder/admin/asset-assignments` | `api_json POST /stakeholder/admin/asset-assignments "$TROVO_ADMIN_TOKEN" payloads/stakeholder-portal/02-assignment-second.json evidence/stakeholder-portal/03-03-assignment-second` | `SECOND_ASSET_ID` | Creates second asset assignment for reject-path tests. | Keeps approval and rejection workflows independently testable. | `200`, `data.asset_id=$SECOND_ASSET_ID`. | Query assignment row for `$SECOND_ASSET_ID`. |
| 4 | `GET /stakeholder/shared/profile` | `api_json GET /stakeholder/shared/profile "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-04-profile-trustee` | Active trustee token | Proves strict org-member auth and `trustees -> trustee` mapping. | Portal frontend can render the trustee dashboard safely. | `200`, `.data.stakeholder.dashboard_role=trustee`. | Org row has `stakeholder_type='trustees'`. |
| 5 | `GET /stakeholder/shared/profile` | `api_json GET /stakeholder/shared/profile "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-05-profile-manager` | Active manager token | Proves `asset_manager -> asset_manager` mapping. | Portal frontend can render manager dashboard safely. | `200`, `.data.stakeholder.dashboard_role=asset_manager`. | Org row has `stakeholder_type='asset_manager'`; `organizations.type` not used. |
| 6 | `GET /stakeholder/shared/profile` | `api_json GET /stakeholder/shared/profile "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-06-profile-custodian` | Active custodian token | Proves `approved_asset_custodian -> asset_custodian` mapping. | Portal frontend can render custodian dashboard safely. | `200`, `.data.stakeholder.dashboard_role=asset_custodian`. | Org row has `stakeholder_type='approved_asset_custodian'`. |
| 7 | `PUT /stakeholder/shared/profile` | `api_json PUT /stakeholder/shared/profile "$MANAGER_TOKEN" payloads/stakeholder-portal/03-profile-update.json evidence/stakeholder-portal/03-07-profile-update` | Seq 5 | Verifies safe profile update fields only. | Allows organization users to maintain display profile without changing auth/linkage. | `200`, first/last name updated. | Member row first/last name changed; org linkage unchanged. |
| 8 | `PUT /stakeholder/shared/profile/notifications` | `api_json PUT /stakeholder/shared/profile/notifications "$MANAGER_TOKEN" payloads/stakeholder-portal/04-notification-preferences.json evidence/stakeholder-portal/03-08-profile-notifications` | Seq 5 | Verifies notification preference persistence. | Supports dashboard user settings. | `200`, preferences returned. | Query `stakeholder_notification_preferences` by `$MANAGER_MEMBER_ID`. |
| 9 | `GET /stakeholder/shared/assets` | `api_json GET "/stakeholder/shared/assets?search=$ASSET_CODE" "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-09-assets-manager` | Seq 2 | Manager sees asset through wallet manager ID or assignment. | Gives manager only their operational assets. | `200`, records include `$ASSET_ID`; `canRequestFundRelease=true`. | Wallet manager ID or assignment row matches manager org. |
| 10 | `GET /stakeholder/shared/assets` | `api_json GET "/stakeholder/shared/assets?search=$ASSET_CODE" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-10-assets-trustee` | Seq 2 | Trustee sees only explicit assignment. | Prevents trustee cross-asset leakage. | `200`, records include `$ASSET_ID`; trustee action flags true. | Assignment row matches trustee org. |
| 11 | `GET /stakeholder/shared/assets` | `api_json GET "/stakeholder/shared/assets?search=$ASSET_CODE" "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-11-assets-custodian` | Seq 2 | Custodian sees custodied/assigned asset. | Gives custodian only assets they safeguard. | `200`, records include `$ASSET_ID`; `canExecuteFundRelease=true`. | Wallet custodian ID or assignment row matches custodian org. |
| 12 | `GET /stakeholder/shared/assets/:asset_id` | `api_json GET "/stakeholder/shared/assets/$ASSET_ID" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-12-asset-detail` | Seq 2 | Verifies asset detail visibility enforcement. | Supports shared asset detail page without leaking unrelated assets. | `200`, `data.id=$ASSET_ID`. | Assignment row matches. |
| 13 | `PUT /stakeholder/asset-manager/assets/:asset_id` | `api_json PUT "/stakeholder/asset-manager/assets/$ASSET_ID" "$MANAGER_TOKEN" payloads/stakeholder-portal/16-asset-operation.json evidence/stakeholder-portal/03-13-asset-operation` | Seq 9 | Verifies manager operational update path. | Supports asset operations/admin reporting. | `200`, operation record returned. | Query `stakeholder_asset_operations` by asset/org; audit action `asset.operation_update`. |
| 14 | `POST /stakeholder/asset-manager/fund-releases` | `api_json POST /stakeholder/asset-manager/fund-releases "$MANAGER_TOKEN" payloads/stakeholder-portal/05-fund-release-approve-path.json evidence/stakeholder-portal/03-14-fr-create-approve` | Seq 2 | Creates fund release request for happy path. | Core manager-to-trustee-to-custodian workflow. | `201`, `data.status=submitted`. | Query `fund_release_requests`; export `FUND_RELEASE_APPROVE_ID`. |
| 15 | `POST /stakeholder/asset-manager/fund-releases` | `api_json POST /stakeholder/asset-manager/fund-releases "$MANAGER_TOKEN" payloads/stakeholder-portal/06-fund-release-reject-path.json evidence/stakeholder-portal/03-15-fr-create-reject` | Seq 3 | Creates separate fund release for rejection path. | Verifies rejection without corrupting happy path state. | `201`, `data.status=submitted`. | Query row; export `FUND_RELEASE_REJECT_ID`. |
| 16 | `GET /stakeholder/asset-manager/fund-releases` | `api_json GET /stakeholder/asset-manager/fund-releases "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-16-fr-list-manager` | Seq 14 | Manager lists own requests. | Gives manager workflow tracking. | `200`, includes manager org requests only. | DB requester_org_id is `$MANAGER_ORG_ID`. |
| 17 | `GET /stakeholder/trustee/fund-releases` | `api_json GET /stakeholder/trustee/fund-releases "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-17-fr-list-trustee` | Seq 14 | Trustee sees assigned requests. | Gives trustee approval queue. | `200`, includes `$FUND_RELEASE_APPROVE_ID`. | DB trustee_org_id is `$TRUSTEE_ORG_ID`. |
| 18 | `GET /stakeholder/trustee/fund-releases/:request_id` | `api_json GET "/stakeholder/trustee/fund-releases/$FUND_RELEASE_APPROVE_ID" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-18-fr-detail-trustee` | Seq 14 | Trustee reads request detail. | Enables approval decision with request context. | `200`, request ID matches. | DB row matches trustee org. |
| 19 | `POST /stakeholder/trustee/fund-releases/:request_id/reject` | `api_json POST "/stakeholder/trustee/fund-releases/$FUND_RELEASE_REJECT_ID/reject" "$TRUSTEE_TOKEN" payloads/stakeholder-portal/07-reject.json evidence/stakeholder-portal/03-19-fr-reject` | Seq 15 | Verifies trustee rejection path. | Supports controlled denial with audit reason. | `200`, `status=trustee_rejected`. | Row has rejection reason; audit action `fund_release.reject`; notification for manager. |
| 20 | `POST /stakeholder/shared/authorizations` | `api_json POST /stakeholder/shared/authorizations "$TRUSTEE_TOKEN" payloads/stakeholder-portal/20-fr-approve-challenge.json evidence/stakeholder-portal/03-20-fr-approve-challenge` | Seq 14 and Step 4 payload generation | Creates wallet step-up challenge bound to trustee/member/org/action/entity. | Protects high-risk approval action. | `201`, `data.id`, `data.auth_id`, expiry. | Query `stakeholder_authorization_challenges`; export `FR_APPROVE_CHALLENGE_ID`. |
| 21 | `POST /stakeholder/shared/authorizations/:challenge_id/verify` | After approving in wallet app: `api_json POST "/stakeholder/shared/authorizations/$FR_APPROVE_CHALLENGE_ID/verify" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-21-fr-approve-challenge-verify` | Seq 20 and wallet approval | Verifies wallet authorization. | Aligns with existing ServiceLink 2FA product flow. | `200`, `status=verified`. | Challenge row `verified_at` not null. |
| 22 | `POST /stakeholder/trustee/fund-releases/:request_id/approve` | `jq -n --arg id "$FR_APPROVE_CHALLENGE_ID" '{challenge_id:$id}' > payloads/stakeholder-portal/22-fr-approve.json`; call approve. | Seq 21 | Consumes verified challenge and approves release. | Core trustee approval workflow. | `200`, `status=execution_pending`. | Request reviewed fields set; challenge `used`; audit/notifications exist. |
| 23 | `GET /stakeholder/custodian/fund-releases` | `api_json GET /stakeholder/custodian/fund-releases "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-23-fr-list-custodian` | Seq 22 | Custodian sees executable approved requests. | Gives custodian execution queue. | `200`, includes `$FUND_RELEASE_APPROVE_ID`. | DB custodian_org_id is `$CUSTODIAN_ORG_ID`. |
| 24 | `POST /stakeholder/shared/authorizations` | `api_json POST /stakeholder/shared/authorizations "$CUSTODIAN_TOKEN" payloads/stakeholder-portal/24-fr-execute-challenge.json evidence/stakeholder-portal/03-24-fr-execute-challenge` | Seq 22 and Step 4 payload generation | Creates execution step-up challenge. | Protects money movement execution. | `201`, challenge returned. | Challenge row bound to custodian member/org. |
| 25 | `POST /stakeholder/shared/authorizations/:challenge_id/verify` | `api_json POST "/stakeholder/shared/authorizations/$FR_EXECUTE_CHALLENGE_ID/verify" "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-25-fr-execute-challenge-verify` | Seq 24 and wallet approval | Verifies execution challenge. | Uses existing wallet authorization flow. | `200`, `status=verified`. | `verified_at` set. |
| 26 | `POST /stakeholder/custodian/fund-releases/:request_id/execute` | `jq -n --arg id "$FR_EXECUTE_CHALLENGE_ID" '{challenge_id:$id,status:"processing",execution_reference:"staging-exec-001"}' > payloads/stakeholder-portal/26-fr-execute.json`; call execute. | Seq 25 | Consumes challenge and moves request to execution status. | Completes custodian leg of fund-release workflow. | `200`, `status=processing`. | `executed_by_member_id`, `executed_at`, challenge `used`, audit/notifications. |
| 27 | `PUT /stakeholder/custodian/fund-releases/:request_id/status` | `jq -n '{status:"completed",execution_reference:"staging-exec-001"}' > payloads/stakeholder-portal/27-fr-status-complete.json`; call status update. | Seq 26 | Updates already-executed release to terminal status. | Lets custodian report processing outcome. | Secure expected result: allowed only after prior step-up execution path. | Row moves to completed; audit action `fund_release.status_update`. |
| 28 | `GET /stakeholder/trustee/due-diligence/:asset_id` | `api_json GET "/stakeholder/trustee/due-diligence/$ASSET_ID" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-28-dd-get` | Seq 2 | Creates/loads trustee checklist. | Gives trustee compliance checklist. | `200`, `data.items` non-empty. | Export `DD_ITEM_ID`; query checklist/items. |
| 29 | `PUT /stakeholder/trustee/due-diligence/:asset_id/items/:item_id` | `api_json PUT "/stakeholder/trustee/due-diligence/$ASSET_ID/items/$DD_ITEM_ID" "$TRUSTEE_TOKEN" payloads/stakeholder-portal/08-due-diligence-item-complete.json evidence/stakeholder-portal/03-29-dd-item-update` | Seq 28 | Marks a due diligence item complete. | Records trustee verification evidence. | `200`, item status `complete`. | Item verified fields set; audit action exists. |
| 30 | `POST /stakeholder/trustee/due-diligence/:asset_id/approve` | `api_json POST "/stakeholder/trustee/due-diligence/$ASSET_ID/approve" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-30-dd-approve` | Seq 29 | Approves local due diligence state. | Supports trustee approval without hidden tokenization mutation. | `200`, checklist `approved`. | AdminDB checklist approved; Wallet asset unchanged before/after. |
| 31 | `GET /stakeholder/trustee/due-diligence/:asset_id` | `api_json GET "/stakeholder/trustee/due-diligence/$SECOND_ASSET_ID" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-31-dd-get-second` | Seq 3 | Creates separate checklist for rejection. | Keeps approve/reject evidence independent. | `200`. | Export `SECOND_DD_ITEM_ID`; checklist exists for second asset. |
| 32 | `POST /stakeholder/trustee/due-diligence/:asset_id/reject` | `api_json POST "/stakeholder/trustee/due-diligence/$SECOND_ASSET_ID/reject" "$TRUSTEE_TOKEN" payloads/stakeholder-portal/07-reject.json evidence/stakeholder-portal/03-32-dd-reject` | Seq 31 | Verifies due diligence rejection path. | Supports trustee denial with reason. | `200`, checklist `rejected`. | Rejection reason, audit, notifications. |
| 33 | `GET /stakeholder/asset-manager/revenue` | `api_json GET /stakeholder/asset-manager/revenue "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-33-revenue-list-before` | Seq 9 | Baseline manager revenue listing. | Revenue dashboard read path. | `200`. | Records scoped to manager org. |
| 34 | `POST /stakeholder/asset-manager/revenue` | `api_json POST /stakeholder/asset-manager/revenue "$MANAGER_TOKEN" payloads/stakeholder-portal/09-revenue.json evidence/stakeholder-portal/03-34-revenue-create` | Seq 9 | Records asset revenue. | Supports manager revenue reporting. | `201`, `status=recorded`. | Query `revenue_records`; export `REVENUE_RECORD_ID`. |
| 35 | `POST /stakeholder/asset-manager/revenue/submit-distribution` | `api_json POST /stakeholder/asset-manager/revenue/submit-distribution "$MANAGER_TOKEN" payloads/stakeholder-portal/35-distribution-authorize-path-with-revenue.json evidence/stakeholder-portal/03-35-distribution-create-authorize` | Seq 34 and Step 4 payload generation | Creates trustee distribution proposal. | Connects manager revenue to trustee authorization. | `201`, distribution `proposed`. | Export `DISTRIBUTION_AUTHORIZE_ID`; revenue status submitted if ID provided. |
| 36 | `POST /stakeholder/asset-manager/revenue/submit-distribution` | `api_json POST /stakeholder/asset-manager/revenue/submit-distribution "$MANAGER_TOKEN" payloads/stakeholder-portal/11-distribution-reject-path.json evidence/stakeholder-portal/03-36-distribution-create-reject` | Seq 3 | Creates separate distribution for rejection. | Tests negative distribution decision independently. | `201`, distribution `proposed`. | Export `DISTRIBUTION_REJECT_ID`. |
| 37 | `GET /stakeholder/trustee/distributions` | `api_json GET /stakeholder/trustee/distributions "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-37-distribution-list` | Seq 35 | Trustee sees assigned distribution proposals. | Provides trustee distribution queue. | `200`, includes proposed distributions. | DB trustee_org_id matches. |
| 38 | `POST /stakeholder/trustee/distributions/:dist_id/reject` | `api_json POST "/stakeholder/trustee/distributions/$DISTRIBUTION_REJECT_ID/reject" "$TRUSTEE_TOKEN" payloads/stakeholder-portal/07-reject.json evidence/stakeholder-portal/03-38-distribution-reject` | Seq 36 | Verifies distribution rejection. | Trustee can block invalid payout proposals. | `200`, `status=rejected`. | Rejection reason, audit, notification. |
| 39 | `POST /stakeholder/shared/authorizations` | `api_json POST /stakeholder/shared/authorizations "$TRUSTEE_TOKEN" payloads/stakeholder-portal/39-distribution-challenge.json evidence/stakeholder-portal/03-39-distribution-challenge` | Seq 35 and Step 4 payload generation | Creates step-up for distribution authorization. | Protects payout authorization. | `201`, challenge returned. | Challenge scoped to trustee/action/entity. |
| 40 | `POST /stakeholder/shared/authorizations/:challenge_id/verify` | `api_json POST "/stakeholder/shared/authorizations/$DISTRIBUTION_CHALLENGE_ID/verify" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-40-distribution-challenge-verify` | Seq 39 and wallet approval | Verifies distribution challenge. | Uses ServiceLink authorization. | `200`, `status=verified`. | `verified_at` set. |
| 41 | `POST /stakeholder/trustee/distributions/:dist_id/authorize` | `api_json POST "/stakeholder/trustee/distributions/$DISTRIBUTION_AUTHORIZE_ID/authorize" "$TRUSTEE_TOKEN" payloads/stakeholder-portal/41-distribution-authorize.json evidence/stakeholder-portal/03-41-distribution-authorize` | Seq 40 | Authorizes distribution with consumed step-up. | Completes trustee payout authorization. | `200`, `status=authorized`. | Distribution authorized fields set; challenge used. |
| 42 | `GET /stakeholder/trustee/distributions/history` | `api_json GET /stakeholder/trustee/distributions/history "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-42-distribution-history` | Seq 38 or 41 | Reads terminal distribution history. | Compliance/audit view for trustee. | `200`, includes rejected/authorized rows. | DB rows terminal. |
| 43 | `GET /stakeholder/asset-manager/valuations/:asset_id` | `api_json GET "/stakeholder/asset-manager/valuations/$ASSET_ID" "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-43-valuations-list-before` | Seq 9 | Baseline valuation history. | Manager asset valuation page. | `200`. | Rows scoped to manager org. |
| 44 | `POST /stakeholder/asset-manager/valuations` | `api_json POST /stakeholder/asset-manager/valuations "$MANAGER_TOKEN" payloads/stakeholder-portal/12-valuation.json evidence/stakeholder-portal/03-44-valuation-create` | Seq 9 | Submits asset valuation. | Supports valuation updates for tokenized assets. | `201`, `status=submitted`. | Query `asset_valuations`; export `VALUATION_ID`. |
| 45 | `POST /stakeholder/asset-manager/valuations/:id/request-independent` | `api_json POST "/stakeholder/asset-manager/valuations/$VALUATION_ID/request-independent" "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-45-valuation-independent` | Seq 44 | Requests independent valuation. | Supports governance escalation. | `200`, `status=independent_requested`. | Audit action `valuation.request_independent`. |
| 46 | `GET /stakeholder/custodian/accounts` | Seed account from DB first, then `api_json GET /stakeholder/custodian/accounts "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-46-accounts-list` | Seed account SQL below | Lists custodian-owned accounts. | Custodian account dashboard. | `200`, own accounts only. | Query `segregated_accounts`; export `ACCOUNT_ID`. |
| 47 | `GET /stakeholder/custodian/accounts/:account_id` | `api_json GET "/stakeholder/custodian/accounts/$ACCOUNT_ID" "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-47-account-detail` | Seq 46 | Reads account detail with ownership check. | Custodian account detail page. | `200`, account org matches. | DB custodian_org_id matches. |
| 48 | `GET /stakeholder/custodian/compliance` | Create compliance with the admin endpoint first, then `api_json GET /stakeholder/custodian/compliance "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-48-compliance-list` | `POST /stakeholder/admin/custodian-compliance` | Lists custodian compliance items. | Compliance checklist page. | `200`, own org only. | Export `COMPLIANCE_ITEM_ID`. |
| 49 | `PUT /stakeholder/custodian/compliance/:item_id` | `api_json PUT "/stakeholder/custodian/compliance/$COMPLIANCE_ITEM_ID" "$CUSTODIAN_TOKEN" payloads/stakeholder-portal/13-compliance-complete.json evidence/stakeholder-portal/03-49-compliance-update` | Seq 48 | Updates compliance item with audit. | Operational compliance tracking. | `200`, status `complete`. | Completed fields set; audit action `compliance.update`. |
| 50 | `POST /stakeholder/shared/documents` | `api_json POST /stakeholder/shared/documents "$MANAGER_TOKEN" payloads/stakeholder-portal/14-document-asset-scoped.json evidence/stakeholder-portal/03-50-document-create` | Seq 9 | Creates URL-only document metadata scoped to visible asset. | Shared document repository. | `201`, document `active`. | Export `DOCUMENT_ID`; query `stakeholder_documents`. |
| 51 | `GET /stakeholder/shared/documents` | `api_json GET "/stakeholder/shared/documents?asset_id=$ASSET_ID" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-51-document-list` | Seq 50 | Lists visible documents by role and asset. | Shared document page for authorized roles. | `200`, includes `DOCUMENT_ID`. | access_roles include trustee. |
| 52 | `GET /stakeholder/shared/documents/:doc_id` | `api_json GET "/stakeholder/shared/documents/$DOCUMENT_ID" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-52-document-detail` | Seq 50 | Reads document metadata with role/asset checks. | Supports document viewing/download link. | `200`, document ID matches. | DB row active. |
| 53 | `GET /stakeholder/shared/notifications` | `api_json GET /stakeholder/shared/notifications "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-53-notifications-list` | Prior workflow notification | Lists notifications scoped to recipient org. | Notification center. | `200`, records for trustee org. | Export `NOTIFICATION_ID`; DB recipient is trustee org. |
| 54 | `PUT /stakeholder/shared/notifications/:id/read` | `api_json PUT "/stakeholder/shared/notifications/$NOTIFICATION_ID/read" "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-54-notification-read` | Seq 53 | Marks own notification read. | Notification center state. | `200`, `read_at` not null. | DB row read_at set. |
| 55 | `GET /stakeholder/shared/audit-trail` | `api_json GET /stakeholder/shared/audit-trail "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-55-audit-trail` | Prior manager state changes | Lists org-scoped audit trail. | Compliance evidence for admin dashboard. | `200`, records for manager org. | Query `stakeholder_audit_logs` by actor org. |
| 56 | `GET /stakeholder/trustee/dashboard` | `api_json GET /stakeholder/trustee/dashboard "$TRUSTEE_TOKEN" - evidence/stakeholder-portal/03-56-dashboard-trustee` | Prior workflows | Aggregates trustee counts. | Trustee dashboard home. | `200`, role `trustee`. | Counts match workflow SQL. |
| 57 | `GET /stakeholder/custodian/dashboard` | `api_json GET /stakeholder/custodian/dashboard "$CUSTODIAN_TOKEN" - evidence/stakeholder-portal/03-57-dashboard-custodian` | Prior workflows/account seed | Aggregates custody counts/balances. | Custodian dashboard home. | `200`, role `asset_custodian`. | Counts match account/release/compliance SQL. |
| 58 | `GET /stakeholder/asset-manager/dashboard` | `api_json GET /stakeholder/asset-manager/dashboard "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-58-dashboard-manager` | Prior workflows | Aggregates manager assets/revenue/approvals. | Manager dashboard home. | `200`, role `asset_manager`. | Counts match asset/revenue/release SQL. |
| 59 | `POST /stakeholder/asset-manager/reports` | `api_json POST /stakeholder/asset-manager/reports "$MANAGER_TOKEN" payloads/stakeholder-portal/17-report.json evidence/stakeholder-portal/03-59-report-create` | Seq 9 | Creates report metadata. | Manager reporting workflow. | `201`, report generated. | Export `REPORT_ID`; query `stakeholder_reports`. |
| 60 | `GET /stakeholder/asset-manager/reports` | `api_json GET /stakeholder/asset-manager/reports "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-60-report-list` | Seq 59 | Lists generated reports. | Manager reporting page. | `200`, includes report. | DB manager_org_id matches. |
| 61 | `POST /stakeholder/asset-manager/reports/:report_id/submit` | `api_json POST "/stakeholder/asset-manager/reports/$REPORT_ID/submit" "$MANAGER_TOKEN" - evidence/stakeholder-portal/03-61-report-submit` | Seq 59 | Submits report to trustee. | Manager-to-trustee reporting workflow. | `200`, status submitted. | submitted_at/member set; notification to trustee if trustee_org_id exists. |

## Step 4: Extract IDs Between Steps

Run each export after the response file named in the command exists. Do not run this whole block before the matching API row has completed.

```bash
export FUND_RELEASE_APPROVE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-14-fr-create-approve.json)"
export FUND_RELEASE_REJECT_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-15-fr-create-reject.json)"

export FR_APPROVE_CHALLENGE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-20-fr-approve-challenge.json)"
export FR_EXECUTE_CHALLENGE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-24-fr-execute-challenge.json)"

export DD_ITEM_ID="$(jq -r '.data.items[0].id' evidence/stakeholder-portal/03-28-dd-get.json)"
export SECOND_DD_ITEM_ID="$(jq -r '.data.items[0].id' evidence/stakeholder-portal/03-31-dd-get-second.json)"

export REVENUE_RECORD_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-34-revenue-create.json)"
export DISTRIBUTION_AUTHORIZE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-35-distribution-create-authorize.json)"
export DISTRIBUTION_REJECT_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-36-distribution-create-reject.json)"
export DISTRIBUTION_CHALLENGE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-39-distribution-challenge.json)"

export VALUATION_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-44-valuation-create.json)"
export DOCUMENT_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-50-document-create.json)"
export NOTIFICATION_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "select id from stakeholder_notifications where recipient_org_id='$TRUSTEE_ORG_ID' order by created_at desc limit 1")"
export REPORT_ID="$(jq -r '.data.id' evidence/stakeholder-portal/03-59-report-create.json)"
```

Generate challenge/action payloads after IDs exist:

```bash
jq -n --arg id "$FUND_RELEASE_APPROVE_ID" \
  '{action:"fund_release.approve",entity_type:"fund_release_request",entity_id:$id}' \
  > payloads/stakeholder-portal/20-fr-approve-challenge.json

jq -n --arg id "$FR_APPROVE_CHALLENGE_ID" '{challenge_id:$id}' \
  > payloads/stakeholder-portal/22-fr-approve.json

jq -n --arg id "$FUND_RELEASE_APPROVE_ID" \
  '{action:"fund_release.execute",entity_type:"fund_release_request",entity_id:$id}' \
  > payloads/stakeholder-portal/24-fr-execute-challenge.json

jq -n --arg id "$FR_EXECUTE_CHALLENGE_ID" \
  '{challenge_id:$id,status:"processing",execution_reference:"staging-exec-001"}' \
  > payloads/stakeholder-portal/26-fr-execute.json

jq --arg id "$REVENUE_RECORD_ID" '.revenue_record_ids=[$id]' \
  payloads/stakeholder-portal/10-distribution-authorize-path.json \
  > payloads/stakeholder-portal/35-distribution-authorize-path-with-revenue.json

jq -n --arg id "$DISTRIBUTION_AUTHORIZE_ID" \
  '{action:"distribution.authorize",entity_type:"distribution",entity_id:$id}' \
  > payloads/stakeholder-portal/39-distribution-challenge.json

jq -n --arg id "$DISTRIBUTION_CHALLENGE_ID" '{challenge_id:$id}' \
  > payloads/stakeholder-portal/41-distribution-authorize.json
```

## Step 5: Provision Data for Custodian Endpoints

Segregated accounts still have no API create endpoint, so seed the account in local/staging AdminDB. Create compliance requirements through the Trovo-admin endpoint so the custodian can list and update them.

```bash
export ACCOUNT_ID="$(psql "$ADMIN_CONNECTION_STRING" -Atc "
insert into segregated_accounts (
  id,asset_id,asset_code,custodian_org_id,account_type,account_name,balance,currency,status,bank_details,created_at,updated_at
)
values (
  gen_random_uuid()::text,'$ASSET_ID','$ASSET_CODE','$CUSTODIAN_ORG_ID','development_funds',
  'Verification custody account',10000.00000000,'NGN','active','{}'::jsonb,now(),now()
)
returning id;
")"

curl -sS -X POST "$BASE_URL/stakeholder/admin/custodian-compliance" \
  -H "Authorization: $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"custodian_org_id\":\"$CUSTODIAN_ORG_ID\",\"category\":\"regulatory\",\"requirement\":\"Monthly custody report\",\"due_date\":\"2026-09-30\"}" \
  | tee evidence/stakeholder-portal/05-compliance-create.json

export COMPLIANCE_ITEM_ID="$(jq -r '.data.id' evidence/stakeholder-portal/05-compliance-create.json)"
```

Proof:

```bash
psql "$ADMIN_CONNECTION_STRING" -c "select id,asset_id,custodian_org_id,balance,currency,status from segregated_accounts where id='$ACCOUNT_ID';" \
  | tee evidence/stakeholder-portal/05-account-seed-proof.txt

psql "$ADMIN_CONNECTION_STRING" -c "select id,org_id,category,requirement,status from compliance_items where id='$COMPLIANCE_ITEM_ID';" \
  | tee evidence/stakeholder-portal/05-compliance-seed-proof.txt
```

## Step 6: Negative and Security Tests

These tests prove authorization boundaries. They should be run after the positive workflow rows have created data.

Generate exact negative-test payloads as needed:

```bash
jq -n '{challenge_id:"not-a-valid-challenge"}' \
  > payloads/stakeholder-portal/N5-invalid-challenge.json

jq -n '{status:"completed",execution_reference:"bypass-test"}' \
  > payloads/stakeholder-portal/N8-status-bypass.json
```

For N5, create a fresh submitted release so the proof is not affected by rows 14-27:

```bash
api_json POST /stakeholder/asset-manager/fund-releases "$MANAGER_TOKEN" \
  payloads/stakeholder-portal/05-fund-release-approve-path.json \
  evidence/stakeholder-portal/N5-fr-create
export N5_FUND_RELEASE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/N5-fr-create.json)"
api_json POST "/stakeholder/trustee/fund-releases/$N5_FUND_RELEASE_ID/approve" "$TRUSTEE_TOKEN" \
  payloads/stakeholder-portal/N5-invalid-challenge.json \
  evidence/stakeholder-portal/N5-approve-without-step-up
```

For N8, create and trustee-approve a fresh release through the real step-up flow, then call the status endpoint without an execution challenge:

```bash
api_json POST /stakeholder/asset-manager/fund-releases "$MANAGER_TOKEN" \
  payloads/stakeholder-portal/05-fund-release-approve-path.json \
  evidence/stakeholder-portal/N8-fr-create
export N8_FUND_RELEASE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/N8-fr-create.json)"

jq -n --arg id "$N8_FUND_RELEASE_ID" \
  '{action:"fund_release.approve",entity_type:"fund_release_request",entity_id:$id}' \
  > payloads/stakeholder-portal/N8-approve-challenge.json
api_json POST /stakeholder/shared/authorizations "$TRUSTEE_TOKEN" \
  payloads/stakeholder-portal/N8-approve-challenge.json \
  evidence/stakeholder-portal/N8-approve-challenge
export N8_APPROVE_CHALLENGE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/N8-approve-challenge.json)"

# Complete wallet approval for this challenge before the verify call.
api_json POST "/stakeholder/shared/authorizations/$N8_APPROVE_CHALLENGE_ID/verify" "$TRUSTEE_TOKEN" - \
  evidence/stakeholder-portal/N8-approve-challenge-verify
jq -n --arg id "$N8_APPROVE_CHALLENGE_ID" '{challenge_id:$id}' \
  > payloads/stakeholder-portal/N8-approve.json
api_json POST "/stakeholder/trustee/fund-releases/$N8_FUND_RELEASE_ID/approve" "$TRUSTEE_TOKEN" \
  payloads/stakeholder-portal/N8-approve.json \
  evidence/stakeholder-portal/N8-approve

api_json PUT "/stakeholder/custodian/fund-releases/$N8_FUND_RELEASE_ID/status" "$CUSTODIAN_TOKEN" \
  payloads/stakeholder-portal/N8-status-bypass.json \
  evidence/stakeholder-portal/N8-status-bypass
```

For N6, create a fresh distribution and intentionally use a challenge bound to the wrong action/entity:

```bash
api_json POST /stakeholder/asset-manager/revenue/submit-distribution "$MANAGER_TOKEN" \
  payloads/stakeholder-portal/11-distribution-reject-path.json \
  evidence/stakeholder-portal/N6-distribution-create
export N6_DISTRIBUTION_ID="$(jq -r '.data.id' evidence/stakeholder-portal/N6-distribution-create.json)"

jq -n --arg id "$N5_FUND_RELEASE_ID" \
  '{action:"fund_release.approve",entity_type:"fund_release_request",entity_id:$id}' \
  > payloads/stakeholder-portal/N6-wrong-action-challenge.json
api_json POST /stakeholder/shared/authorizations "$TRUSTEE_TOKEN" \
  payloads/stakeholder-portal/N6-wrong-action-challenge.json \
  evidence/stakeholder-portal/N6-wrong-action-challenge
export N6_WRONG_ACTION_CHALLENGE_ID="$(jq -r '.data.id' evidence/stakeholder-portal/N6-wrong-action-challenge.json)"

jq -n --arg id "$N6_WRONG_ACTION_CHALLENGE_ID" '{challenge_id:$id}' \
  > payloads/stakeholder-portal/N6-wrong-action-authorize.json
api_json POST "/stakeholder/trustee/distributions/$N6_DISTRIBUTION_ID/authorize" "$TRUSTEE_TOKEN" \
  payloads/stakeholder-portal/N6-wrong-action-authorize.json \
  evidence/stakeholder-portal/N6-wrong-action-authorize
```

For N9, create the unscoped regression document with manager org A and list documents with manager org B:

```bash
api_json POST /stakeholder/shared/documents "$MANAGER_TOKEN" \
  payloads/stakeholder-portal/15-document-unscoped-regression.json \
  evidence/stakeholder-portal/N9-unscoped-doc-create
export N9_UNSCOPED_DOC_ID="$(jq -r '.data.id' evidence/stakeholder-portal/N9-unscoped-doc-create.json)"

api_json GET /stakeholder/shared/documents "$OTHER_MANAGER_TOKEN" - \
  evidence/stakeholder-portal/N9-other-manager-doc-list
jq --arg id "$N9_UNSCOPED_DOC_ID" '[.data.records[] | select(.id == $id)]' \
  evidence/stakeholder-portal/N9-other-manager-doc-list.json \
  | tee evidence/stakeholder-portal/N9-unscoped-doc-leak-proof.json
```

| Seq | Negative test | Payload / Command | Depends On | Endpoint purpose and task goal | Admin-dashboard product fit | Required result | DB proof |
|---:|---|---|---|---|---|---|---|
| N1 | Trovo admin token rejected from stakeholder profile | `api_json GET /stakeholder/shared/profile "$TROVO_ADMIN_TOKEN" - evidence/stakeholder-portal/N1-admin-rejected-profile` | Admin token | Ensures portal routes require org-member auth. | Separates admin support from stakeholder user access. | `401` or `403`; no profile. | None. |
| N2 | Wrong role blocked | `api_json GET /stakeholder/trustee/dashboard "$MANAGER_TOKEN" - evidence/stakeholder-portal/N2-wrong-role` | Manager token | Ensures role middleware blocks route misuse. | Prevents frontend/backend role confusion. | `403`. | None. |
| N3 | Other trustee cannot read asset | `api_json GET "/stakeholder/shared/assets/$ASSET_ID" "$OTHER_TRUSTEE_TOKEN" - evidence/stakeholder-portal/N3-other-trustee-asset` | Other trustee token | Ensures trustee explicit assignment only. | Prevents asset leakage across trustees. | `404` or `403`. | Assignment row does not match other trustee org. |
| N4 | Other org cannot mark notification read | `api_json PUT "/stakeholder/shared/notifications/$NOTIFICATION_ID/read" "$OTHER_TRUSTEE_TOKEN" - evidence/stakeholder-portal/N4-other-notification-read` | Seq 53 | Ensures notification recipient scoping. | Protects org-specific workflow notifications. | `404` or `403`. | Original row `read_at` unchanged by other org. |
| N5 | Approval without step-up fails | Use `payloads/stakeholder-portal/N5-invalid-challenge.json` against `POST /stakeholder/trustee/fund-releases/$N5_FUND_RELEASE_ID/approve`. | Fresh submitted fund release from N5 setup | Ensures trustee approval requires wallet step-up. | Protects high-risk fund movement approval. | `403`; status remains `submitted`. | Fund release row unchanged. |
| N6 | Wrong-action challenge cannot be consumed | Use `payloads/stakeholder-portal/N6-wrong-action-authorize.json` against `POST /stakeholder/trustee/distributions/$N6_DISTRIBUTION_ID/authorize`. | Fresh proposed distribution and wrong-action challenge from N6 setup | Ensures challenge is bound to action/entity. | Prevents replaying wallet approvals across workflows. | `403`; challenge not `used`. | Challenge row status unchanged. |
| N7 | One-time challenge consumption | Repeat Seq 22 or Seq 26 with same challenge. | Successful challenge consumption | Ensures challenge consumed exactly once. | Prevents replay of high-risk approval/execution. | `403` or `409`. | Challenge remains `used`; no duplicate transition. |
| N8 | P1: status endpoint cannot bypass step-up | Use `payloads/stakeholder-portal/N8-status-bypass.json` against `PUT /stakeholder/custodian/fund-releases/$N8_FUND_RELEASE_ID/status` before custodian `/execute`. | Fresh `execution_pending` release from N8 setup | Detects custodian execution bypass. | Protects money movement execution. | `403`; if `200`, the P1 regression returned. | Fund release remains `execution_pending`; execution fields remain null. |
| N9 | P1: unscoped document cannot leak to same-role org | Create `15-document-unscoped-regression.json` as manager org A, list documents as manager org B, then inspect `N9-unscoped-doc-leak-proof.json`. | Other manager token and `OTHER_MANAGER_ORG_ID` | Detects unscoped document leakage. | Protects org-private document metadata/URLs. | Leak proof JSON is `[]`; if it contains the doc, the P1 regression returned. | Document `uploaded_by_org_id` is manager org A. |

## Step 7: SQL Proof Queries After Workflow

Run these after the sequence.

```bash
psql "$ADMIN_CONNECTION_STRING" -c "
select id,status,amount,currency,requester_org_id,trustee_org_id,custodian_org_id,reviewed_by_member_id,executed_by_member_id
from fund_release_requests
where id in ('$FUND_RELEASE_APPROVE_ID','$FUND_RELEASE_REJECT_ID');
" | tee evidence/stakeholder-portal/07-fund-release-proof.txt

psql "$ADMIN_CONNECTION_STRING" -c "
select id,member_id,organization_id,action,entity_type,entity_id,status,verified_at,used_at
from stakeholder_authorization_challenges
where entity_id in ('$FUND_RELEASE_APPROVE_ID','$DISTRIBUTION_AUTHORIZE_ID');
" | tee evidence/stakeholder-portal/07-challenge-proof.txt

psql "$ADMIN_CONNECTION_STRING" -c "
select action,entity_type,entity_id,actor_org_id,created_at
from stakeholder_audit_logs
where entity_id in ('$FUND_RELEASE_APPROVE_ID','$FUND_RELEASE_REJECT_ID','$DISTRIBUTION_AUTHORIZE_ID','$DOCUMENT_ID','$REPORT_ID')
order by created_at;
" | tee evidence/stakeholder-portal/07-audit-proof.txt

psql "$ADMIN_CONNECTION_STRING" -c "
select recipient_org_id,type,related_entity_type,related_entity_id,read_at,created_at
from stakeholder_notifications
where related_entity_id in ('$FUND_RELEASE_APPROVE_ID','$FUND_RELEASE_REJECT_ID','$DISTRIBUTION_AUTHORIZE_ID','$REPORT_ID')
order by created_at;
" | tee evidence/stakeholder-portal/07-notification-proof.txt

psql "$ADMIN_CONNECTION_STRING" -c "
select id,asset_id,uploaded_by_org_id,access_roles,status,file_url
from stakeholder_documents
where id='$DOCUMENT_ID';
" | tee evidence/stakeholder-portal/07-document-proof.txt

psql "$WALLET_DB_CONNECTION_STRING" -c "
select id,asset_code,asset_tokenization_status,asset_manager_id,approved_asset_custodian_id
from tokenized_assets
where id in ('$ASSET_ID','$SECOND_ASSET_ID');
" | tee evidence/stakeholder-portal/07-wallet-unchanged-proof.txt
```

Expected proof:

- Workflow state is in AdminDB.
- Challenges are bound to member/org/action/entity and used once.
- State-changing actions have audit logs.
- Notifications are scoped to recipient orgs.
- Wallet tokenization rows are not silently mutated by portal workflows.

## Step 8: Final Code-Level Verification

| Area | Command | Required result |
|---|---|---|
| Routes registered | `rg -n 'stakeholder\\.Init|Group\\(\"/stakeholder|\\.(GET|POST|PUT)\\(' main.go internal/components/stakeholder/controllers/main.go \| tee evidence/stakeholder-portal/08-routes.txt` | All stakeholder route groups and endpoints present. |
| AdminDB/TrovoWalletDB boundary | `rg -n 'TrovoWalletDB|TokenizationReadClient|Create\\(|Save\\(|Updates\\(' internal/components/stakeholder \| tee evidence/stakeholder-portal/08-db-boundary.txt` | TrovoWalletDB appears only as tokenization read-client source; portal writes use AdminDB services/models. |
| No hidden tokenization mutation | `rg -n 'makeRequest|/tokenization|Update.*token|GetTokenizedAssetAuthorizationRequest|Send.*Tokenized' internal/components/stakeholder \| tee evidence/stakeholder-portal/08-no-tokenization-mutation.txt` | No stakeholder component tokenization mutation calls. |
| Gofmt | `gofmt -w main.go internal/middleware/stakeholder_auth_middleware.go internal/middleware/stakeholder_auth_middleware_test.go internal/components/stakeholder` | No formatting diff afterward. |
| Tests | `go test ./... \| tee evidence/stakeholder-portal/08-go-test-all.txt` | All packages pass. |
| Swagger | `go run github.com/swaggo/swag/cmd/swag@v1.16.2 init -g main.go \| tee evidence/stakeholder-portal/08-swagger-generate.txt` | Command exits 0; stakeholder paths remain in `docs/swagger.json` and `docs/swagger.yaml`. Existing unrelated warnings should be recorded. |

## Acceptance Rule

SP-001 through SP-048 are verified only when:

- Every positive row in Step 3 returns the expected 2xx response.
- Every negative row in Step 6 returns the expected 401/403/404/409 response.
- SQL proof in Step 7 matches the expected workflow state.
- Wallet DB proof shows no portal workflow tables and no hidden tokenization mutation.
- `go test ./...`, `gofmt`, and Swagger generation pass.
- The two P1 regression tests pass.

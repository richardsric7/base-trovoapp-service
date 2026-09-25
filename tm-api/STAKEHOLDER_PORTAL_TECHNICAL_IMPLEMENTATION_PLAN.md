# Stakeholder Portal Technical Implementation Plan

## Purpose

Build stakeholder-facing API routes for Trustees, Asset Custodians, and Asset Managers on top of this admin-dashboard API codebase.

This plan replaces the earlier high-level implementation plan with codebase-specific architecture decisions, module boundaries, data ownership, state machines, and build order.

## Current Codebase Facts

- The API is a Go 1.21 Gin service.
- `main.go` wires controller packages directly and calls each package `Init(router, s)`.
- `serverModels.Server` exposes three database handles:
  - `AdminDB`: admin dashboard data, organizations, organization members, invites.
  - `TrovoWalletDB`: wallet/tokenization/partner data.
  - `P2P`: P2P data.
- `AdminDB` uses GORM `AutoMigrate` in `internal/db/main.go`.
- Organization records already have `stakeholder_id` and `stakeholder_type`.
- Live AdminDB was inspected on 2026-06-29. It currently has admin/config tables, `organizations`, `organization_members`, `organization_invites`, and `career_roles`; it does not yet have stakeholder portal workflow tables.
- Existing `organizations.stakeholder_id` and `organizations.stakeholder_type` are the current organization-to-stakeholder linkage and should be reused instead of creating a parallel stakeholder organization registry.
- Existing AdminDB constraints matter for portal design:
  - `organizations.email` is globally unique.
  - `organization_members.email` is globally unique, not unique per organization.
  - `organization_members.organization_id` has a foreign key to `organizations.id`.
  - `organization_invites.email` is globally unique.
- Current AdminDB data needs cleanup or seed updates before strict stakeholder auth can pass:
  - Existing organization statuses are `PENDING` and `INACTIVE`, with no `ACTIVE` organizations in the inspected data.
  - Existing organization type casing is inconsistent, for example `ASSET_MANAGER` and `asset_manager`.
  - Existing stakeholder types include `approved_asset_custodian`, `asset_manager`, `legal_and_professionals`, and null values in the inspected data; no trustee organization was present in the inspected data.
- Existing stakeholder types are:
  - `trustees`
  - `approved_asset_custodian`
  - `asset_manager`
  - plus non-portal partner types.
- Organization auth exists, but the active mixed auth path is `AllowOrgOrTrovoAdminNormalized`.
- Current organization middleware does not hydrate `stakeholder_id`, `stakeholder_type`, or a canonical dashboard role from the database.
- Existing tokenization endpoints are registered in `internal/components/usermetrics/controllers/main.go` and proxy to Trovo Wallet `/v1/trovo-manager/...`.
- `trovosdk.GetTokenizedAssetData(assetCode)` only fetches one tokenized asset. It is not suitable as a list source.
- Tokenized asset data currently exposes `assetManagerId` and `approvedAssetCustodianId`. No trustee assignment field is visible in the local code.

## Architecture Decisions

### 1. Add a Dedicated Stakeholder Component

Create a new component instead of extending `organizations` or `usermetrics`.

New package layout:

```text
internal/components/stakeholder/
  controllers/main.go
  handlers/
    assets_handler.go
    trustee_handler.go
    custodian_handler.go
    asset_manager_handler.go
    documents_handler.go
    notifications_handler.go
    profile_handler.go
    authorizations_handler.go
  services/
    asset_service.go
    authz_service.go
    due_diligence_service.go
    fund_release_service.go
    distribution_service.go
    document_service.go
    notification_service.go
    audit_service.go
    dashboard_service.go
  db/
    asset_repository.go
    workflow_repository.go
    document_repository.go
    notification_repository.go
    audit_repository.go
  models/
    requests.go
    responses.go
    constants.go
```

Register it in `main.go`:

```go
stakeholder "admin-panel-dashboard/internal/components/stakeholder/controllers"

stakeholder.Init(router, s)
```

### 2. Keep Portal Workflow State in AdminDB

`AdminDB` should own stakeholder portal workflow data because it already owns organizations and org members.

`TrovoWalletDB` remains the source for tokenization assets and partner records. The stakeholder portal should not mutate tokenization state silently. Any mutation of tokenization state must be an explicit integration boundary.

Do not introduce a duplicate organization-stakeholder mapping table for the MVP. Use:

- `organizations.id` as the portal organization ID.
- `organizations.stakeholder_id` as the linked stakeholder/partner ID.
- `organizations.stakeholder_type` as the raw external stakeholder type.
- a canonical portal role derived from `stakeholder_type`.

Local workflow tables should store stable references:

- `asset_id` or `tokenized_asset_id`
- `asset_code` when available
- `organization_id`
- `stakeholder_id`
- `stakeholder_type`
- actor member IDs
- status
- timestamps

### 3. Use Canonical Dashboard Roles

Do not use route names, organization types, and partner table names interchangeably.

Canonical portal roles:

| Dashboard role | Existing stakeholder_type | Route segment |
|---|---|---|
| `trustee` | `trustees` | `trustee` |
| `asset_custodian` | `approved_asset_custodian` | `custodian` |
| `asset_manager` | `asset_manager` | `asset-manager` |

Add constants in `internal/components/stakeholder/models/constants.go` or shared `internal/models`.

Normalize only for routing/authorization. Do not rewrite external partner type strings in business records unless a separate data migration explicitly owns that cleanup.

`organizations.type` is not reliable as the portal role source because live data contains mixed casing and values. Use `stakeholder_type` for portal role mapping.

### 4. Stakeholder Routes Require Organization Members

MVP stakeholder portal routes should require an organization member token, not a Trovo admin token.

Trovo admin support/impersonation can be added later with explicit audit logging. Do not include it in the MVP.

### 5. Tokenization Read Access Needs a Client Boundary

Do not import or call `usermetrics` handlers from the stakeholder module.

Create a small tokenization read client interface:

```go
type TokenizationReadClient interface {
    ListAssets(ctx context.Context, filters AssetFilters) (AssetPage, error)
    GetAsset(ctx context.Context, idOrCode string) (*TokenizedAsset, error)
}
```

Recommended MVP implementation:

- Use `s.TrovoWalletDB` with read-only GORM models for tokenized assets.
- Keep the interface small so it can later be replaced by a Trovo Wallet service API client.

Avoid using `trovosdk.GetTokenizedAssetData()` for asset lists. It can be used only for single-asset enrichment if its response has fields missing from the DB query.

### 6. Add Trustee Asset Assignment

The current visible tokenization model has manager and custodian IDs, but no trustee assignment. Add a portal-owned assignment table so trustee asset filtering is deterministic.

Table: `stakeholder_asset_assignments`

Fields:

- `id`
- `asset_id`
- `asset_code`
- `trustee_org_id`
- `trustee_stakeholder_id`
- `asset_manager_org_id`
- `asset_manager_stakeholder_id`
- `custodian_org_id`
- `custodian_stakeholder_id`
- `status`
- `assigned_by`
- `created_at`
- `updated_at`

MVP source of assignments:

- Trovo admin or migration script creates trustee assignments.
- Manager and custodian can be synced from tokenized asset fields when possible.

## Security Model

### Stakeholder Auth Middleware

Add `internal/middleware/stakeholder_auth_middleware.go`.

Responsibilities:

1. Parse organization member JWT.
2. Accept both raw token and `Bearer <token>` consistently.
3. Verify member exists and is active in `AdminDB`.
4. Verify organization exists and is active in `AdminDB`.
5. Verify token claims match member and organization.
6. Load organization `stakeholder_id` and `stakeholder_type`.
7. Map `stakeholder_type` to canonical `dashboard_role`.
8. Reject organizations without a supported portal stakeholder type.
9. Set context:
   - `member_id`
   - `member_email`
   - `member_role`
   - `organization_id`
   - `organization_name`
   - `organization_type`
   - `stakeholder_id`
   - `stakeholder_type`
   - `dashboard_role`
   - `trovo_wallet_username`
   - `is_wallet_linked`

Rollout dependency: before enabling strict `OrganizationStatusActive` checks in non-local environments, seed or migrate intended stakeholder organizations to `ACTIVE`. The inspected AdminDB data currently has no active organizations, so strict auth would reject all existing orgs until that data is fixed.

Middleware helpers:

```go
RequireStakeholderRole(roles ...string) gin.HandlerFunc
RequireOrganizationAdmin() gin.HandlerFunc
```

### Step-Up Authorization

High-risk actions must require a fresh Trovo Wallet authorization:

- Trustee fund release approval
- Trustee distribution authorization
- Custodian fund release execution
- Rejection can optionally require step-up based on product policy; MVP can require it for approvals/execution only.

Use existing ServiceLink methods:

- `SendAuthorizationRequest`
- `VerifyAuthorizationRequest`

Add durable challenge table: `stakeholder_authorization_challenges`

Fields:

- `id`
- `member_id`
- `organization_id`
- `stakeholder_type`
- `dashboard_role`
- `action`
- `entity_type`
- `entity_id`
- `auth_id`
- `status`: `pending`, `verified`, `expired`, `failed`, `used`
- `expires_at`
- `verified_at`
- `used_at`
- `created_at`
- `updated_at`

Rules:

- Challenge is bound to member, organization, action, entity type, and entity ID.
- Verified challenge can be used once.
- Challenge expires after a short TTL, for example 10 minutes.
- State-changing handlers must consume the verified challenge inside the same DB transaction as the state change.

Routes:

```text
POST /api/v1/stakeholder/shared/authorizations
POST /api/v1/stakeholder/shared/authorizations/:challenge_id/verify
```

Request:

```json
{
  "action": "fund_release.approve",
  "entity_type": "fund_release_request",
  "entity_id": "..."
}
```

## Data Model

Add GORM models to `internal/models/stakeholder_portal.go` or component-local models if the team wants to keep domain ownership inside the component.

Use explicit SQL migrations for stakeholder workflow tables. The repo already has `migrations/` for `career_roles`, and the portal tables need production-grade indexes, uniqueness constraints, status checks, and foreign keys that should not depend only on GORM `AutoMigrate`.

GORM models should still be added for repository usage. If the team keeps AutoMigrate enabled locally, include the models only as a developer convenience after the SQL migration strategy is defined.

### Core Tables

#### stakeholder_asset_assignments

Purpose: maps assets to stakeholder organizations, especially trustees.

Uniqueness:

- Unique on `asset_id`
- Indexes on `asset_code`, `trustee_org_id`, `asset_manager_org_id`, `custodian_org_id`
- Foreign keys to `organizations(id)` for org ID columns where possible.
- Store `*_stakeholder_id` as the linked external stakeholder ID copied from `organizations.stakeholder_id` at assignment time.

#### due_diligence_checklists

Purpose: trustee verification checklist per asset.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `trustee_org_id`
- `status`: `draft`, `in_review`, `approved`, `rejected`
- `approved_by_member_id`
- `approved_at`
- `rejected_by_member_id`
- `rejected_at`
- `rejection_reason`
- `created_at`
- `updated_at`

#### due_diligence_items

Purpose: checklist items instead of storing all items on one row.

Fields:

- `id`
- `checklist_id`
- `category`
- `item`
- `status`: `pending`, `complete`, `failed`, `not_applicable`
- `notes`
- `verified_by_member_id`
- `verified_at`
- `created_at`
- `updated_at`

#### fund_release_requests

Purpose: Asset Manager requests funds, Trustee approves, Custodian executes.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `requester_org_id`
- `requester_member_id`
- `trustee_org_id`
- `custodian_org_id`
- `amount`
- `currency`
- `purpose`
- `status`
- `supporting_document_ids`
- `reviewed_by_member_id`
- `reviewed_at`
- `rejection_reason`
- `executed_by_member_id`
- `executed_at`
- `execution_reference`
- `failure_reason`
- `created_at`
- `updated_at`

Status machine:

```text
draft -> submitted -> trustee_approved -> execution_pending -> processing -> completed
submitted -> trustee_rejected
execution_pending -> failed
processing -> failed
```

Invariants:

- Only Asset Manager org can create/submit for assigned assets.
- Only Trustee org assigned to the asset can approve/reject.
- Only Custodian org assigned to the asset can execute/update execution status.
- Approval and execution require verified step-up challenge.
- Terminal statuses cannot transition.

#### distributions

Purpose: revenue distribution proposal and trustee authorization.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `proposed_by_org_id`
- `proposed_by_member_id`
- `trustee_org_id`
- `amount`
- `currency`
- `source`
- `scheduled_date`
- `status`: `proposed`, `authorized`, `rejected`, `processing`, `completed`, `failed`
- `authorized_by_member_id`
- `authorized_at`
- `rejection_reason`
- `created_at`
- `updated_at`

#### revenue_records

Purpose: Asset Manager records asset revenue.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `manager_org_id`
- `period_start`
- `period_end`
- `source`
- `amount`
- `currency`
- `status`: `draft`, `recorded`, `submitted_for_distribution`
- `collected_at`
- `created_by_member_id`
- `created_at`
- `updated_at`

#### asset_valuations

Purpose: Asset Manager records valuation history.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `manager_org_id`
- `valuation`
- `currency`
- `methodology`
- `valuation_date`
- `report_document_id`
- `status`: `submitted`, `independent_requested`, `accepted`, `rejected`
- `created_by_member_id`
- `created_at`
- `updated_at`

#### segregated_accounts

Purpose: Custodian account records for an assigned asset.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `custodian_org_id`
- `account_type`: `sale_proceeds`, `development_funds`, `revenue`, `reserve`
- `account_name`
- `balance`
- `currency`
- `status`: `active`, `inactive`, `closed`
- `bank_details`
- `created_by_member_id`
- `created_at`
- `updated_at`

#### compliance_items

Purpose: Custodian compliance checklist.

Fields:

- `id`
- `org_id`
- `category`
- `requirement`
- `status`: `pending`, `complete`, `overdue`, `waived`
- `due_date`
- `completed_at`
- `completed_by_member_id`
- `created_at`
- `updated_at`

#### stakeholder_documents

Purpose: Portal document repository. Store metadata and file URLs, not file bytes.

Fields:

- `id`
- `asset_id`
- `asset_code`
- `uploaded_by_org_id`
- `uploaded_by_member_id`
- `category`
- `title`
- `file_url`
- `mime_type`
- `version`
- `access_roles`
- `status`: `active`, `superseded`, `deleted`
- `created_at`
- `updated_at`

#### stakeholder_notifications

Purpose: In-app notifications.

Fields:

- `id`
- `recipient_org_id`
- `sender_org_id`
- `type`
- `title`
- `message`
- `related_entity_type`
- `related_entity_id`
- `read_at`
- `created_at`

#### stakeholder_audit_logs

Purpose: immutable audit trail for portal actions.

Fields:

- `id`
- `request_id`
- `actor_member_id`
- `actor_org_id`
- `actor_role`
- `action`
- `entity_type`
- `entity_id`
- `before_state`
- `after_state`
- `metadata`
- `created_at`

No update/delete routes should exist for audit logs.

## API Routes

Base group:

```text
/api/v1/stakeholder
```

Controller setup:

```go
apiV1 := router.Group("/api/v1")
stakeholder := apiV1.Group("/stakeholder", middleware.StakeholderAuthMiddleware(s.AdminDB))
```

### Shared Routes

```text
GET  /shared/profile
PUT  /shared/profile

GET  /shared/assets
GET  /shared/assets/:asset_id

POST /shared/authorizations
POST /shared/authorizations/:challenge_id/verify

POST /shared/documents
GET  /shared/documents
GET  /shared/documents/:doc_id

GET  /shared/notifications
PUT  /shared/notifications/:id/read

GET  /shared/audit-trail
```

### Trustee Routes

```text
GET  /trustee/dashboard

GET  /trustee/due-diligence/:asset_id
PUT  /trustee/due-diligence/:asset_id/items/:item_id
POST /trustee/due-diligence/:asset_id/approve
POST /trustee/due-diligence/:asset_id/reject

GET  /trustee/fund-releases
GET  /trustee/fund-releases/:request_id
POST /trustee/fund-releases/:request_id/approve
POST /trustee/fund-releases/:request_id/reject

GET  /trustee/distributions
POST /trustee/distributions/:dist_id/authorize
POST /trustee/distributions/:dist_id/reject
GET  /trustee/distributions/history
```

### Asset Custodian Routes

```text
GET  /custodian/dashboard

GET  /custodian/accounts
GET  /custodian/accounts/:account_id

GET  /custodian/fund-releases
POST /custodian/fund-releases/:request_id/execute
PUT  /custodian/fund-releases/:request_id/status

GET  /custodian/compliance
PUT  /custodian/compliance/:item_id
```

Trovo admins assign custodian compliance requirements through:

```text
POST /stakeholder/admin/custodian-compliance
```

### Asset Manager Routes

```text
GET  /asset-manager/dashboard

GET  /asset-manager/assets
PUT  /asset-manager/assets/:asset_id

GET  /asset-manager/revenue
POST /asset-manager/revenue
POST /asset-manager/revenue/submit-distribution

GET  /asset-manager/valuations/:asset_id
POST /asset-manager/valuations
POST /asset-manager/valuations/:id/request-independent

POST /asset-manager/fund-releases
GET  /asset-manager/fund-releases

POST /asset-manager/reports
GET  /asset-manager/reports
POST /asset-manager/reports/:report_id/submit
```

## Asset Filtering Rules

`AssetService` must apply server-side filtering after resolving context.

### Asset Manager

An asset is visible if:

- tokenized asset `assetManagerId == organization.stakeholder_id`, or
- `stakeholder_asset_assignments.asset_manager_org_id == organization_id`.

### Asset Custodian

An asset is visible if:

- tokenized asset `approvedAssetCustodianId == organization.stakeholder_id`, or
- `stakeholder_asset_assignments.custodian_org_id == organization_id`.

### Trustee

An asset is visible if:

- `stakeholder_asset_assignments.trustee_org_id == organization_id`, or
- `stakeholder_asset_assignments.trustee_stakeholder_id == organization.stakeholder_id`.

If trustee assignment is missing, the asset must not appear in a trustee view.

## Service-Level Patterns

### Handlers

Handlers should:

1. Bind and validate request body/query.
2. Read auth context.
3. Call a service method.
4. Return with `serverResponse.JSON` where possible.

### Services

Services own:

- Permission checks beyond route role.
- State transitions.
- Transactions.
- Audit calls.
- Notification creation.
- Step-up authorization consumption.

### Repositories

Repositories own GORM queries. Keep them small and explicit.

### Audit Helper

Every state-changing service method should call:

```go
AuditService.Record(ctx, AuditEvent{
    ActorMemberID: ...,
    ActorOrgID: ...,
    ActorRole: ...,
    Action: ...,
    EntityType: ...,
    EntityID: ...,
    BeforeState: ...,
    AfterState: ...,
    Metadata: ...,
})
```

Prefer writing audit rows in the same transaction as the state change.

## Build Sequence

### Phase 0: Foundation Refactor

Goal: prepare safe extension points.

Tasks:

- Add stakeholder component skeleton.
- Register stakeholder controller in `main.go`.
- Add role constants and mapping helpers.
- Add `StakeholderAuthMiddleware`.
- Add `RequireStakeholderRole`.
- Add unit tests for role mapping and unsupported stakeholder types.
- Add route groups with placeholder health/profile endpoint.

Acceptance criteria:

- Org member with linked `stakeholder_type=trustees` can access trustee profile.
- Same member cannot access asset-manager route.
- Organization with unsupported/missing stakeholder type receives 403.
- Trovo admin token cannot access stakeholder routes in MVP.

### Phase 1: Asset Read Model

Goal: shared asset list/detail with correct role filtering.

Tasks:

- Add tokenization read models or client interface.
- Implement `TokenizationReadClient` using `TrovoWalletDB`.
- Add `stakeholder_asset_assignments`.
- Add admin/manual seed path for trustee assignments if needed.
- Implement:
  - `GET /shared/assets`
  - `GET /shared/assets/:asset_id`
- Include role action flags in response metadata:
  - Trustee: `canApproveDueDiligence`, `canApproveFundRelease`, `canAuthorizeDistribution`
  - Custodian: `canExecuteFundRelease`, `canUpdateAccount`
  - Asset Manager: `canRequestFundRelease`, `canRecordRevenue`, `canSubmitValuation`

Acceptance criteria:

- Asset Manager sees only assigned/managed assets.
- Custodian sees only assigned/custodied assets.
- Trustee sees only explicitly assigned assets.
- No endpoint returns cross-role assets.

### Phase 2: Audit, Notifications, and Step-Up Authorization

Goal: reusable primitives before business workflows.

Tasks:

- Add GORM models and explicit SQL migrations:
  - `stakeholder_audit_logs`
  - `stakeholder_notifications`
  - `stakeholder_authorization_challenges`
- Implement `AuditService`.
- Implement `NotificationService`.
- Implement authorization challenge create/verify endpoints.
- Add challenge consumption helper for business services.

Acceptance criteria:

- Verified challenge can be consumed once.
- Expired challenge cannot be used.
- Audit record is created for every successful state-changing test endpoint.

### Phase 3: Fund Release Workflow

Goal: implement the first end-to-end workflow across all three roles.

Reason: it exercises Asset Manager creation, Trustee approval, Custodian execution, role filtering, notifications, audit, and 2FA.

Tasks:

- Add `fund_release_requests`.
- Implement Asset Manager routes:
  - `POST /asset-manager/fund-releases`
  - `GET /asset-manager/fund-releases`
- Implement Trustee routes:
  - `GET /trustee/fund-releases`
  - `GET /trustee/fund-releases/:request_id`
  - `POST /trustee/fund-releases/:request_id/approve`
  - `POST /trustee/fund-releases/:request_id/reject`
- Implement Custodian routes:
  - `GET /custodian/fund-releases`
  - `POST /custodian/fund-releases/:request_id/execute`
  - `PUT /custodian/fund-releases/:request_id/status`
- Add notifications:
  - Manager request submitted -> Trustee notified.
  - Trustee approved -> Custodian and Manager notified.
  - Trustee rejected -> Manager notified.
  - Custodian completed/failed -> Trustee and Manager notified.

Acceptance criteria:

- Invalid transitions fail with 409.
- Approval and execution require fresh verified step-up challenge.
- Approve/reject/execute writes audit logs.
- Non-assigned organizations cannot read or mutate requests.

### Phase 4: Trustee Due Diligence

Goal: support trustee verification without corrupting tokenization state.

Tasks:

- Add `due_diligence_checklists` and `due_diligence_items`.
- Implement checklist get/update.
- Implement approve/reject.
- Decide whether approval should call Trovo Wallet tokenization state API.
- If yes, add explicit integration service with idempotency and error handling.

Acceptance criteria:

- Trustee can approve only assigned assets.
- Approval writes audit log and notification.
- If external tokenization update fails, local state remains consistent and user gets a clear error.

### Phase 5: Revenue, Distribution, and Valuation

Goal: implement Asset Manager reporting and Trustee distribution authorization.

Tasks:

- Add `revenue_records`, `distributions`, `asset_valuations`.
- Implement revenue record/create/list.
- Implement submit distribution.
- Implement trustee authorize/reject distribution with step-up.
- Implement valuation list/create/request-independent.

Acceptance criteria:

- Manager can create records only for managed assets.
- Trustee can authorize only assigned asset distributions.
- Authorization requires step-up.
- All state changes are audited.

### Phase 6: Custodian Accounts and Compliance

Goal: support custody/account management after fund release workflow exists.

Tasks:

- Add `segregated_accounts` and `compliance_items`.
- Implement account list/detail.
- Implement compliance list/update.
- Link account views to assigned custodian assets.

Acceptance criteria:

- Custodian sees only own assigned account records.
- Compliance updates are audited.

### Phase 7: Documents and Profile

Goal: shared supporting features.

Tasks:

- Add `stakeholder_documents`.
- Implement document metadata create/list/get.
- Keep file storage out of scope unless storage provider is confirmed.
- Implement profile get/update for allowed fields only.
- Implement notification preference persistence if frontend requires it.

Acceptance criteria:

- Documents are role-filtered.
- Deleted/superseded docs are not returned by default.
- Profile update cannot change organization identity, stakeholder linkage, role, or wallet linkage.

### Phase 8: Dashboard Aggregates and Caching

Goal: role dashboards after workflows exist.

Tasks:

- Implement dashboard services using existing workflow tables and asset read service.
- Add Redis caching only after query shapes stabilize.
- Cache keys must include organization ID and dashboard role.
- Cache invalidation should happen on workflow state changes.

Acceptance criteria:

- Dashboards return correct counts for role and organization.
- Cache never leaks one org's dashboard to another org.

## Error Handling

Use status codes consistently:

- `400`: invalid request body/query.
- `401`: missing/invalid token.
- `403`: valid token but wrong role/org/asset ownership.
- `404`: entity not found within caller's authorized scope.
- `409`: invalid state transition or duplicate idempotent action.
- `422`: semantically invalid input, such as amount <= 0.
- `500`: unexpected server error.
- `502/503`: external tokenization or ServiceLink dependency failed.

## Idempotency and Concurrency

Add idempotency for high-risk transitions:

- Fund release approval
- Fund release execution
- Distribution authorization

Recommended approach:

- Optional `Idempotency-Key` header for POST state transitions.
- Persist key, member ID, action, entity ID, and response status.
- Use DB transaction and row locking for status transitions.

At minimum:

- Check current status before update.
- Update with `WHERE id = ? AND status = ?`.
- Return 409 when no row changes.

## Observability

MVP observability:

- Audit logs for all state changes.
- `log.Printf` around external calls and failed transitions, matching current repo style.
- Include `request_id` in audit metadata if request ID middleware is added.

Future:

- Structured logs.
- Metrics for workflow transitions.
- Metrics for external tokenization call latency/failures.

## Testing Plan

Add tests around the highest-risk contracts.

Unit tests:

- Stakeholder type to dashboard role mapping.
- Unsupported stakeholder types.
- State machine transitions.
- Step-up challenge expiry and one-time use.

Handler/service tests:

- Wrong role receives 403.
- Unassigned org receives 404/403 for asset and workflow entity.
- Fund release lifecycle happy path.
- Duplicate approval returns 409 or idempotent response.
- Approval without challenge returns 403.

Integration tests:

- Tokenization read client filters manager/custodian/trustee assets correctly.
- Audit rows are created in the same transaction as state changes.

## Swagger and Docs

Add Swagger annotations as each handler is implemented.

Do not generate Swagger only at the end; this repo already keeps `docs/` checked in, so each phase should update generated Swagger output.

## Non-Goals for MVP

- Full bank transfer execution.
- Token-holder payout execution.
- Reconciliation engine.
- Background workers/queues.
- Redis caching before dashboard query shape is stable.
- Trovo admin impersonation of stakeholder users.
- File storage provider implementation unless storage requirements are confirmed.
- Replacing existing admin tokenization endpoints.

## Resolved Implementation Decisions

Resolved on 2026-06-29 for this implementation:

1. Trovo Wallet tokenization data visible in this repo does not expose a trustee assignment field. Trustee visibility is therefore owned by AdminDB through `stakeholder_asset_assignments`. Missing trustee assignment means the asset is not visible to trustees.
2. Due diligence approval is local-only for MVP. It updates `due_diligence_checklists` in AdminDB, writes audit/notifications, and does not mutate Trovo Wallet tokenization state.
3. Stakeholder asset reads use direct read-only `s.TrovoWalletDB` access behind `TokenizationReadClient`. The implementation filters by tokenized asset manager/custodian IDs where available and by AdminDB assignment rows.
4. Document storage is URL-only metadata in `stakeholder_documents`. File bytes, object storage credentials, max file size enforcement, and virus scanning are out of scope until a storage provider is selected.
5. Step-up authorization is required for approvals/execution: trustee fund release approval, custodian fund release execution, and trustee distribution authorization. Rejections do not require step-up in MVP.
6. Trovo admins do not have stakeholder view impersonation in MVP. The only admin portal route is the explicit asset assignment creation endpoint.
7. Money fields use Go `github.com/shopspring/decimal` and DB `numeric(30,8)`. Currency values are uppercased and constrained to 3-12 characters.

## First Build Slice

Build this first:

1. `StakeholderAuthMiddleware`
2. role mapping constants
3. stakeholder route groups
4. `stakeholder_asset_assignments`
5. tokenization read client
6. `GET /api/v1/stakeholder/shared/profile`
7. `GET /api/v1/stakeholder/shared/assets`
8. `GET /api/v1/stakeholder/shared/assets/:asset_id`

This slice proves auth, stakeholder linkage, role mapping, asset data access, and cross-role data isolation before implementing workflow tables.

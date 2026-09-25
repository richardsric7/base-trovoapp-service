# Stakeholder Portal Ticket Backlog

## MVP Goal

Deliver the first usable stakeholder portal API slice for organization members: authenticate as a linked stakeholder organization, see the correct assets for that role, and complete one audited fund release workflow across Asset Manager, Trustee, and Asset Custodian.

## Live AdminDB Baseline

AdminDB was inspected on 2026-06-29 before revising this backlog.

- Existing tables include `organizations`, `organization_members`, `organization_invites`, admin/config tables, and `career_roles`.
- No stakeholder portal workflow tables exist yet.
- `organizations.stakeholder_id` and `organizations.stakeholder_type` already exist and are the MVP source for stakeholder linkage.
- `organization_members.organization_id` has a foreign key to `organizations.id`.
- `organizations.email`, `organization_members.email`, and `organization_invites.email` are globally unique.
- Current organization data has `PENDING` and `INACTIVE` statuses only; no inspected organization is `ACTIVE`.
- Organization type casing is inconsistent, for example `ASSET_MANAGER` and `asset_manager`; portal role mapping must use `stakeholder_type`, not `organizations.type`.

## Ticket Format

- **Type:** decision, backend, database, test, docs
- **Depends on:** tickets that must land first
- **Acceptance:** concrete checks before the ticket is done

## Phase 0: Blocking Decisions

### SP-001: Decide Tokenization Asset Read Source

**Type:** decision

**Depends on:** none

**Scope:**

- Decide whether stakeholder asset reads use direct `TrovoWalletDB` queries, Trovo Wallet service API, or a service credential/client.
- Document the selected path in `STAKEHOLDER_PORTAL_TECHNICAL_IMPLEMENTATION_PLAN.md`.
- Define how pagination, filters, and asset detail look for the selected source.

**Acceptance:**

- The plan no longer recommends one asset source while listing it as an open question.
- The first asset list implementation has a known data source and auth model.
- Data leakage risks for the selected source are documented.

### SP-002: Define Trustee Asset Assignment Source

**Type:** decision

**Depends on:** SP-001

**Scope:**

- Use existing `organizations.stakeholder_id` and `organizations.stakeholder_type` as the org-to-stakeholder linkage for the MVP.
- Confirm whether Trovo Wallet tokenization data already has trustee assignment.
- If not, confirm that `stakeholder_asset_assignments` is the system of record for trustee visibility.
- Define how initial trustee assignments are created: admin endpoint, SQL seed, migration backfill, or manual operational script.
- Do not create a parallel stakeholder organization registry unless the existing organization linkage cannot satisfy the workflow.

**Acceptance:**

- Trustee asset visibility has a concrete source.
- The source of organization stakeholder linkage is `organizations.stakeholder_id/stakeholder_type`.
- First build slice includes a way to create or seed assignments.
- Missing trustee assignment behavior is explicitly "not visible".

### SP-003: Define Money Precision Rules

**Type:** decision

**Depends on:** none

**Scope:**

- Choose database type and Go type for `amount`, `balance`, and `valuation`.
- Define currency validation rules.
- Define min/max and decimal places for fund release, revenue, distribution, valuation, and account balances.

**Acceptance:**

- Plan names precise DB column types, for example `numeric(30, 8)`.
- Request validation expectations are documented.
- Tickets that create financial tables reference this decision.

### SP-004: Decide Workflow Migration Strategy and AdminDB Baseline

**Type:** decision

**Depends on:** none

**Scope:**

- Decide whether stakeholder workflow tables are created by explicit SQL migrations, GORM AutoMigrate, or both.
- For production workflow tables, prefer explicit SQL migrations for indexes and constraints.
- Define naming convention for migration files.
- Capture the current AdminDB schema baseline in the technical plan.
- Define whether strict stakeholder auth requires `organizations.status = 'ACTIVE'` on day one, and document the data migration or seed process needed because inspected AdminDB data currently has no active organizations.
- Document the current global uniqueness of `organization_members.email` and `organization_invites.email`; do not design portal invite/member flows that assume one email can belong to multiple organizations unless a separate migration changes those constraints.
- Document organization type casing cleanup or normalization policy. Portal authorization must derive role from `stakeholder_type`, not `organizations.type`.

**Acceptance:**

- The plan states the migration approach.
- Workflow tables have index and constraint strategy.
- No workflow ticket relies only on implicit AutoMigrate for production constraints.
- AdminDB baseline facts are recorded in the plan.
- There is a clear precondition for which stakeholder organizations can authenticate after strict auth is enabled.

## Phase 1: Auth and Route Foundation

### SP-005: Add Stakeholder Component Skeleton

**Type:** backend

**Depends on:** SP-001, SP-002, SP-004

**Scope:**

- Create `internal/components/stakeholder/`.
- Add `controllers/main.go`.
- Add `handlers`, `services`, `db`, and `models` packages.
- Register `stakeholder.Init(router, s)` in `main.go`.
- Add placeholder route under `/api/v1/stakeholder/shared/health`.

**Acceptance:**

- Service compiles.
- `GET /api/v1/stakeholder/shared/health` returns 200 before auth is applied or under a temporary internal route.
- No existing route behavior changes.

### SP-006: Add Portal Role Constants and Mapping

**Type:** backend

**Depends on:** SP-005

**Scope:**

- Add canonical dashboard roles:
  - `trustee`
  - `asset_custodian`
  - `asset_manager`
- Add mapping from existing `stakeholder_type` values:
  - `trustees`
  - `approved_asset_custodian`
  - `asset_manager`
- Add helper that rejects unsupported stakeholder types.
- Treat `organizations.type` as display/classification data only; live values have inconsistent casing and must not drive portal access.
- Normalize input for comparison, but preserve raw `stakeholder_type` in persisted records and responses.

**Acceptance:**

- Unit tests cover supported mappings.
- Unit tests cover unsupported stakeholder types.
- Unit tests prove `organizations.type` does not grant stakeholder route access.
- Constants are reused by middleware and route role checks.

### SP-007: Implement Stakeholder Auth Middleware

**Type:** backend

**Depends on:** SP-006

**Scope:**

- Add `internal/middleware/stakeholder_auth_middleware.go`.
- Accept both raw JWT and `Bearer <token>` formats if existing clients require both.
- Verify JWT signing method.
- Load member from `AdminDB` and require `ACTIVE`.
- Load organization from `AdminDB` and require `ACTIVE`.
- Verify token claims match member and organization.
- Load `stakeholder_id` and `stakeholder_type`.
- Reject missing or unsupported stakeholder linkage.
- Set context values used by stakeholder handlers.

**Acceptance:**

- Org member with supported stakeholder type reaches a protected test route.
- Trovo admin token is rejected for stakeholder routes.
- Missing stakeholder linkage returns 403.
- Suspended/inactive member or organization returns 401/403 consistently.
- Test fixtures or seed data include an active stakeholder organization; otherwise every current live organization would be rejected by design.
- Tests cover invalid token, wrong signing method, and claim mismatch.

### SP-008: Add Stakeholder Role Middleware

**Type:** backend

**Depends on:** SP-007

**Scope:**

- Implement `RequireStakeholderRole(roles ...string)`.
- Apply role middleware to route groups:
  - `/trustee`
  - `/custodian`
  - `/asset-manager`
- Keep `/shared` available to all supported roles.

**Acceptance:**

- Trustee token cannot access asset-manager route.
- Asset Manager token cannot access custodian route.
- Shared profile route works for all three supported roles.

### SP-009: Implement Shared Profile Endpoint

**Type:** backend

**Depends on:** SP-008

**Scope:**

- Add `GET /api/v1/stakeholder/shared/profile`.
- Return member, organization, stakeholder, dashboard role, and wallet linkage fields.
- Include raw `organization.type` and raw `stakeholder_type` for transparency, but use `dashboard_role` for frontend authorization decisions.
- Do not expose password, reset tokens, OTPs, or sensitive internal auth fields.

**Acceptance:**

- Profile response matches authenticated organization member.
- Response includes `dashboard_role`, `stakeholder_id`, and `stakeholder_type`.
- Sensitive fields are absent.

## Phase 2: Asset Read Model and Assignments

### SP-010: Add Asset Assignment Migration

**Type:** database

**Depends on:** SP-002, SP-004

**Scope:**

- Create explicit SQL migration for `stakeholder_asset_assignments`.
- Use AdminDB as the owner of this table.
- Add foreign keys to `organizations(id)` for organization ID columns where possible.
- Include indexes for:
  - `asset_id`
  - `asset_code`
  - `trustee_org_id`
  - `trustee_stakeholder_id`
  - `asset_manager_org_id`
  - `asset_manager_stakeholder_id`
  - `custodian_org_id`
  - `custodian_stakeholder_id`
- Add matching GORM model.

**Acceptance:**

- Migration can be applied on a clean database.
- GORM model maps to the table.
- Table references existing AdminDB organization IDs instead of creating a new organization registry.
- Tests or local verification confirm indexes exist.

### SP-011: Add Asset Assignment Seed/Admin Path

**Type:** backend

**Depends on:** SP-010

**Scope:**

- Implement the assignment creation path chosen in SP-002.
- MVP option: Trovo admin-only endpoint under `/api/v1/stakeholder/admin/asset-assignments`.
- Validate asset identifiers and organization IDs.
- Validate assigned organizations by reading `organizations.stakeholder_id` and `organizations.stakeholder_type`.
- Ensure only supported stakeholder organizations can be assigned.

**Acceptance:**

- A trustee assignment can be created for an asset.
- Duplicate assignment handling is deterministic.
- Invalid org/stakeholder role combinations are rejected.
- Assignment creation is audited if audit exists; otherwise include audit in follow-up SP-018.

### SP-012: Implement Tokenization Read Client

**Type:** backend

**Depends on:** SP-001

**Scope:**

- Add `TokenizationReadClient` interface.
- Implement selected read source from SP-001.
- Support list filters:
  - page
  - limit
  - status
  - asset class/sector if available
  - date range
  - asset code/name search if available
- Support single asset detail by ID or code.

**Acceptance:**

- Client returns paginated assets.
- Client returns one asset detail.
- Errors from external/API/DB source are mapped predictably.
- Unit tests cover query construction or repository filtering.

### SP-013: Implement Role-Based Asset Service

**Type:** backend

**Depends on:** SP-010, SP-012

**Scope:**

- Implement `AssetService.ListAssets`.
- Implement `AssetService.GetAsset`.
- Apply server-side visibility rules:
  - Asset Manager: stakeholder ID or assignment match.
  - Custodian: stakeholder ID or assignment match.
  - Trustee: explicit assignment only.
- Add role-specific action flags.

**Acceptance:**

- Asset Manager sees only managed/assigned assets.
- Custodian sees only custodied/assigned assets.
- Trustee sees only explicitly assigned assets.
- Unassigned asset detail returns 404 or 403 within caller scope.
- Tests cover cross-role leakage cases.

### SP-014: Add Shared Asset Routes

**Type:** backend

**Depends on:** SP-013

**Scope:**

- Add `GET /api/v1/stakeholder/shared/assets`.
- Add `GET /api/v1/stakeholder/shared/assets/:asset_id`.
- Bind and validate pagination/filter query params.
- Return consistent response envelope.

**Acceptance:**

- List endpoint supports pagination.
- Detail endpoint enforces role visibility.
- Response includes role action flags.
- Swagger annotations are added.

## Phase 3: Shared Audit, Notification, and Step-Up Primitives

### SP-015: Add Audit Log Migration and Model

**Type:** database

**Depends on:** SP-004

**Scope:**

- Create `stakeholder_audit_logs`.
- Include immutable event fields:
  - actor member/org/role
  - action
  - entity type/id
  - before state
  - after state
  - metadata
  - created_at
- Add indexes on actor org, entity, action, and created_at.

**Acceptance:**

- Migration applies cleanly.
- No update/delete repository methods are exposed for audit logs.
- Model supports JSON metadata fields.

### SP-016: Implement Audit Service

**Type:** backend

**Depends on:** SP-015

**Scope:**

- Implement `AuditService.Record`.
- Support writing within an existing DB transaction.
- Add helper for extracting actor info from stakeholder auth context.

**Acceptance:**

- Audit event can be written inside a transaction.
- Tests verify audit row is rolled back when transaction fails.
- Audit service is used by assignment creation if SP-011 exists.

### SP-017: Add Notifications Migration and Service

**Type:** backend

**Depends on:** SP-004

**Scope:**

- Create `stakeholder_notifications`.
- Implement `NotificationService.Create`.
- Implement `GET /shared/notifications`.
- Implement `PUT /shared/notifications/:id/read`.

**Acceptance:**

- Notifications are scoped to recipient organization.
- Mark-read cannot mark another org's notification.
- Notification list is paginated.

### SP-018: Add Authorization Challenge Migration

**Type:** database

**Depends on:** SP-003, SP-004

**Scope:**

- Create `stakeholder_authorization_challenges`.
- Add status values:
  - `pending`
  - `verified`
  - `expired`
  - `failed`
  - `used`
- Index by member/org/action/entity/status.

**Acceptance:**

- Migration applies cleanly.
- Table supports one-time challenge consumption queries.
- Expiry fields are indexed or efficiently queryable.

### SP-019: Implement Step-Up Authorization Service

**Type:** backend

**Depends on:** SP-018

**Scope:**

- Implement create challenge using linked Trovo wallet.
- Implement verify challenge using ServiceLink `VerifyAuthorizationRequest`.
- Implement consume verified challenge inside workflow transaction.
- Bind challenge to action and entity.

**Acceptance:**

- Challenge create returns auth ID/deeplink/QR data.
- Challenge verify moves `pending` to `verified`.
- Verified challenge can be consumed exactly once.
- Expired or wrong-action challenge cannot be consumed.
- Tests cover expiry and one-time use.

### SP-020: Add Authorization Challenge Routes

**Type:** backend

**Depends on:** SP-019

**Scope:**

- Add `POST /shared/authorizations`.
- Add `POST /shared/authorizations/:challenge_id/verify`.
- Validate action/entity payloads.

**Acceptance:**

- Unlinked wallet returns a clear 403 or 422.
- Invalid action/entity combination returns 400.
- Verified challenge response includes status and expiry.
- Swagger annotations are added.

## Phase 4: Fund Release Workflow MVP

### SP-021: Add Fund Release Migration and Model

**Type:** database

**Depends on:** SP-003, SP-004, SP-010

**Scope:**

- Create `fund_release_requests`.
- Use money precision rules from SP-003.
- Add status fields and indexes.
- Add constraints needed for state transitions.
- Decide whether `trustee_approved` and `execution_pending` are one or two states; update plan accordingly.

**Acceptance:**

- Migration applies cleanly.
- Status values and financial columns are constrained.
- Indexes support lists by requester, trustee, custodian, asset, status.

### SP-022: Implement Fund Release State Machine

**Type:** backend

**Depends on:** SP-021

**Scope:**

- Implement status transition helper.
- Enforce allowed transitions.
- Use `WHERE id = ? AND status = ?` update pattern or row locking.
- Return 409 for invalid or stale transitions.

**Acceptance:**

- Unit tests cover valid transitions.
- Unit tests cover invalid transitions.
- Duplicate approval/execution is rejected or idempotent by design.

### SP-023: Asset Manager Creates Fund Release Request

**Type:** backend

**Depends on:** SP-013, SP-016, SP-017, SP-022

**Scope:**

- Add `POST /asset-manager/fund-releases`.
- Validate asset visibility and assignment.
- Validate amount/currency/purpose.
- Create request in initial submitted status.
- Notify assigned trustee.
- Audit creation.

**Acceptance:**

- Asset Manager can create request for managed asset.
- Asset Manager cannot create request for unassigned asset.
- Trustee org is resolved from assignment.
- Audit and notification records are created.

### SP-024: Asset Manager Lists Own Fund Release Requests

**Type:** backend

**Depends on:** SP-023

**Scope:**

- Add `GET /asset-manager/fund-releases`.
- Filter to current organization.
- Support pagination and status filter.

**Acceptance:**

- Manager sees only own org requests.
- Pagination works.
- Request status is visible.

### SP-025: Trustee Reviews Fund Release Requests

**Type:** backend

**Depends on:** SP-023

**Scope:**

- Add `GET /trustee/fund-releases`.
- Add `GET /trustee/fund-releases/:request_id`.
- Filter to assigned trustee organization.

**Acceptance:**

- Trustee sees only requests assigned to its org.
- Detail endpoint rejects unrelated requests.
- Supporting document placeholders or URLs are handled according to current document decision.

### SP-026: Trustee Approves or Rejects Fund Release

**Type:** backend

**Depends on:** SP-019, SP-022, SP-025

**Scope:**

- Add `POST /trustee/fund-releases/:request_id/approve`.
- Add `POST /trustee/fund-releases/:request_id/reject`.
- Approval requires verified step-up challenge.
- Rejection follows SP-003/plan policy for step-up if required.
- Notify manager and custodian.
- Audit transition.

**Acceptance:**

- Approval without valid challenge fails.
- Approval moves request to execution-ready status.
- Rejection records reason.
- Invalid repeated approval/rejection returns 409 or idempotent response.

### SP-027: Custodian Lists Executable Fund Releases

**Type:** backend

**Depends on:** SP-026

**Scope:**

- Add `GET /custodian/fund-releases`.
- Show only approved/execution-ready requests assigned to custodian org.
- Support pagination and status filter.

**Acceptance:**

- Custodian sees only assigned requests.
- Non-approved requests are excluded unless explicitly filtered for history.

### SP-028: Custodian Executes Fund Release

**Type:** backend

**Depends on:** SP-019, SP-022, SP-027

**Scope:**

- Add `POST /custodian/fund-releases/:request_id/execute`.
- Add `PUT /custodian/fund-releases/:request_id/status`.
- Execution requires verified step-up challenge.
- Capture execution reference or failure reason.
- Notify trustee and manager.
- Audit transition.

**Acceptance:**

- Execution without valid challenge fails.
- Request can move to processing/completed/failed according to state machine.
- Repeated execution is rejected or idempotent.
- All transitions are audited.

### SP-029: Add Fund Release Workflow Tests

**Type:** test

**Depends on:** SP-023, SP-026, SP-028

**Scope:**

- Test full happy path:
  - manager creates
  - trustee approves
  - custodian executes
- Test unauthorized org access.
- Test invalid state transitions.
- Test missing/used/expired challenge.

**Acceptance:**

- Tests run with `go test`.
- Coverage includes service-level state transition and permission paths.

## Phase 5: Due Diligence

### SP-030: Decide Due Diligence External Mutation

**Type:** decision

**Depends on:** SP-014

**Scope:**

- Decide whether trustee due diligence approval updates Trovo Wallet tokenization state.
- If yes, identify exact endpoint/client and failure handling.
- If no, document local-only semantics.

**Acceptance:**

- Due diligence implementation has a single source-of-truth policy.
- External mutation behavior is testable.

### SP-031: Add Due Diligence Migrations and Models

**Type:** database

**Depends on:** SP-004, SP-030

**Scope:**

- Create `due_diligence_checklists`.
- Create `due_diligence_items`.
- Add indexes by asset, trustee org, status.

**Acceptance:**

- Migration applies cleanly.
- Items are normalized under checklist.
- Status values are constrained.

### SP-032: Implement Trustee Due Diligence APIs

**Type:** backend

**Depends on:** SP-013, SP-016, SP-017, SP-031

**Scope:**

- Add `GET /trustee/due-diligence/:asset_id`.
- Add `PUT /trustee/due-diligence/:asset_id/items/:item_id`.
- Add approve/reject routes.
- Enforce trustee assignment.
- Audit and notify on final decision.

**Acceptance:**

- Trustee can manage only assigned asset checklist.
- Checklist item updates are audited.
- Approval/rejection creates notifications.
- External mutation behavior follows SP-030.

## Phase 6: Revenue, Distribution, and Valuation

### SP-033: Add Revenue, Distribution, and Valuation Migrations

**Type:** database

**Depends on:** SP-003, SP-004

**Scope:**

- Create `revenue_records`.
- Create `distributions`.
- Create `asset_valuations`.
- Add indexes by asset, org, status, and date.

**Acceptance:**

- Migrations apply cleanly.
- Financial fields use SP-003 precision.
- Status values are constrained.

### SP-034: Implement Asset Manager Revenue APIs

**Type:** backend

**Depends on:** SP-013, SP-016, SP-017, SP-033

**Scope:**

- Add `GET /asset-manager/revenue`.
- Add `POST /asset-manager/revenue`.
- Add `POST /asset-manager/revenue/submit-distribution`.
- Enforce managed asset visibility.

**Acceptance:**

- Manager records revenue only for managed assets.
- Submit distribution creates distribution proposal for assigned trustee.
- Audit and notification records are created.

### SP-035: Implement Trustee Distribution APIs

**Type:** backend

**Depends on:** SP-019, SP-034

**Scope:**

- Add `GET /trustee/distributions`.
- Add `POST /trustee/distributions/:dist_id/authorize`.
- Add `POST /trustee/distributions/:dist_id/reject`.
- Add `GET /trustee/distributions/history`.
- Authorization requires step-up challenge.

**Acceptance:**

- Trustee sees only assigned distributions.
- Authorization without valid challenge fails.
- Authorization/rejection audited and notified.

### SP-036: Implement Asset Valuation APIs

**Type:** backend

**Depends on:** SP-013, SP-016, SP-033

**Scope:**

- Add `GET /asset-manager/valuations/:asset_id`.
- Add `POST /asset-manager/valuations`.
- Add `POST /asset-manager/valuations/:id/request-independent`.
- Enforce managed asset visibility.

**Acceptance:**

- Manager can submit valuation only for managed assets.
- Valuation history is asset-scoped and org-authorized.
- Independent request changes status and writes audit.

## Phase 7: Custodian Accounts and Compliance

### SP-037: Add Custodian Account and Compliance Migrations

**Type:** database

**Depends on:** SP-003, SP-004

**Scope:**

- Create `segregated_accounts`.
- Create `compliance_items`.
- Add indexes by custodian org, asset, status, due date.

**Acceptance:**

- Migrations apply cleanly.
- Financial balance fields use SP-003 precision.
- Status values are constrained.

### SP-038: Implement Custodian Account APIs

**Type:** backend

**Depends on:** SP-013, SP-016, SP-037

**Scope:**

- Add `GET /custodian/accounts`.
- Add `GET /custodian/accounts/:account_id`.
- Enforce custodian asset assignment.

**Acceptance:**

- Custodian sees only own assigned asset accounts.
- Account detail rejects unrelated org access.

### SP-039: Implement Custodian Compliance APIs

**Type:** backend

**Depends on:** SP-016, SP-037

**Scope:**

- Add `GET /custodian/compliance`.
- Add `PUT /custodian/compliance/:item_id`.
- Enforce org ownership.
- Audit updates.

**Acceptance:**

- Custodian sees own compliance checklist.
- Updates cannot affect another org's compliance items.
- Updates are audited.

## Phase 8: Documents and Profile Updates

### SP-040: Decide Document Storage Scope

**Type:** decision

**Depends on:** none

**Scope:**

- Decide whether MVP stores file URLs only or uploads to a storage provider.
- If upload provider is needed, define provider, credentials, max file size, MIME allowlist, and virus scanning policy.

**Acceptance:**

- Document implementation ticket has a clear storage path.
- Security constraints for uploads are documented if uploads are in scope.

### SP-041: Add Documents Migration and Model

**Type:** database

**Depends on:** SP-004, SP-040

**Scope:**

- Create `stakeholder_documents`.
- Add indexes by asset, uploader org, category, status.
- Define `access_roles` representation.

**Acceptance:**

- Migration applies cleanly.
- Deleted/superseded documents can be filtered.

### SP-042: Implement Shared Document APIs

**Type:** backend

**Depends on:** SP-013, SP-016, SP-041

**Scope:**

- Add `POST /shared/documents`.
- Add `GET /shared/documents`.
- Add `GET /shared/documents/:doc_id`.
- Enforce asset visibility and document role access.

**Acceptance:**

- User can list only documents for visible assets.
- Role access is enforced.
- Document creation is audited.

### SP-043: Implement Shared Profile Update API

**Type:** backend

**Depends on:** SP-009

**Scope:**

- Add `PUT /shared/profile`.
- Allow only safe member profile fields.
- Do not allow organization identity, stakeholder linkage, role, status, or wallet linkage changes.

**Acceptance:**

- Safe profile fields can update.
- Protected fields are ignored or rejected.
- Update is audited.

## Phase 9: Dashboards, Swagger, and Readiness

### SP-044: Implement Role Dashboard Services

**Type:** backend

**Depends on:** SP-014, SP-029, SP-032, SP-035, SP-039

**Scope:**

- Add:
  - `GET /trustee/dashboard`
  - `GET /custodian/dashboard`
  - `GET /asset-manager/dashboard`
- Aggregate counts from asset service and workflow tables.
- Keep Redis caching out until query behavior stabilizes.

**Acceptance:**

- Dashboard counts are scoped by org and role.
- No cross-org data appears in aggregates.
- Tests cover at least one dashboard per role.

### SP-045: Add Dashboard Cache After Query Shape Stabilizes

**Type:** backend

**Depends on:** SP-044

**Scope:**

- Add Redis caching for dashboard aggregates if performance requires it.
- Include organization ID and role in cache keys.
- Invalidate on workflow state changes.

**Acceptance:**

- Cache key cannot leak across organizations or roles.
- State changes invalidate affected dashboard keys.
- Caching can be disabled through existing `ENABLE_CACHING`.

### SP-046: Complete Swagger Coverage

**Type:** docs

**Depends on:** all implemented endpoint tickets

**Scope:**

- Add Swagger annotations to all stakeholder handlers.
- Regenerate `docs/`.
- Verify auth schemes are documented for organization member tokens.

**Acceptance:**

- Swagger includes all stakeholder routes.
- Request and response schemas are present.
- Existing Swagger routes still render.

### SP-047: End-to-End Security Review

**Type:** test

**Depends on:** SP-044, SP-046

**Scope:**

- Review every stakeholder route for:
  - auth required
  - role required
  - org ownership checked
  - asset visibility checked
  - audit for state changes
  - step-up for high-risk transitions
- Add missing tests for gaps.

**Acceptance:**

- No shared or role route can be accessed without auth.
- No role can read another role/org's workflow data.
- High-risk actions require step-up.
- State-changing routes have audit coverage.

### SP-048: Release Readiness Checklist

**Type:** docs

**Depends on:** SP-047

**Scope:**

- Document environment variables.
- Document migration order.
- Document seed/assignment procedure.
- Document the AdminDB data readiness check:
  - stakeholder organizations have supported `stakeholder_type`
  - stakeholder organizations intended to log in are `ACTIVE`
  - organization type casing is not used for portal authorization
  - global member/invite email uniqueness is accepted or migrated intentionally
- Document rollback plan.
- Document known non-goals and manual operations.

**Acceptance:**

- Another engineer can deploy and seed the stakeholder portal from the checklist.
- Known manual steps are explicit.
- Open questions have been resolved or deferred with owner and date.

## Suggested Milestones

### Milestone 1: Secure Asset Visibility

Tickets: SP-001 through SP-014

Outcome: stakeholder org members can authenticate and see only their assets.

### Milestone 2: Workflow Primitives

Tickets: SP-015 through SP-020

Outcome: audit, notification, and step-up authorization primitives exist.

### Milestone 3: Fund Release MVP

Tickets: SP-021 through SP-029

Outcome: Asset Manager -> Trustee -> Custodian fund release workflow works end to end.

### Milestone 4: Remaining Role Workflows

Tickets: SP-030 through SP-039

Outcome: due diligence, revenue/distribution, valuation, custodian accounts, and compliance are implemented.

### Milestone 5: Supporting Features and Launch

Tickets: SP-040 through SP-048

Outcome: documents, profile update, dashboards, Swagger, security review, and release checklist are complete.

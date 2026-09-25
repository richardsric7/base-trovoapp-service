# Stakeholder Portal Implementation Status

Date: 2026-06-30

## Completed Scope

Implemented code paths for tickets SP-001 through SP-048 in code, migrations, Swagger docs, and release documentation. Test coverage was added for the highest-risk primitives, but it is not exhaustive endpoint-by-endpoint coverage for every acceptance bullet.

Completed ticket groups:

- SP-001 to SP-004: decisions recorded in `STAKEHOLDER_PORTAL_TECHNICAL_IMPLEMENTATION_PLAN.md`.
- SP-005 to SP-009: stakeholder component, route groups, role mapping, strict stakeholder auth middleware, and shared profile endpoints.
- SP-010 to SP-014: explicit assignment migration/model, admin assignment endpoint, `TokenizationReadClient`, role-filtered asset service, and shared asset routes.
- SP-015 to SP-020: audit logs, notification service/routes, authorization challenge table/service/routes, one-time step-up consumption.
- SP-021 to SP-029: fund release table/model, state machine, manager create/list, trustee review/approve/reject, custodian execute/status, and tests.
- SP-030 to SP-032: due diligence local-only decision, checklist/item migrations/models, trustee APIs.
- SP-033 to SP-036: revenue, distribution, and valuation migrations/models/APIs.
- SP-037 to SP-039: segregated account and compliance migrations/models/APIs.
- SP-040 to SP-043: URL-only document decision, document migration/model/APIs, safe profile updates, notification preferences.
- SP-044 to SP-046: role dashboard services/routes, Swagger annotations, regenerated `docs/`.
- SP-047 to SP-048: security review coverage, release readiness checklist, migration/seed/rollback documentation.

## Key Files

- Component: `internal/components/stakeholder/`
- Middleware: `internal/middleware/stakeholder_auth_middleware.go`
- Migration: `migrations/20260629_create_stakeholder_portal_tables.sql`
- Swagger: `internal/components/stakeholder/handlers/swagger_docs.go`, `docs/`
- Env sample: `.env-sample`
- DB/API verification runbook: `STAKEHOLDER_PORTAL_TEST_VERIFICATION_FLOW.md`

## Security Review Fixes

Fixed after review:

- Custodian fund-release status updates can no longer move a request out of `execution_pending`; the request must first go through `/stakeholder/custodian/fund-releases/:request_id/execute`, which consumes a verified `fund_release.execute` wallet challenge.
- Unscoped stakeholder documents are now visible only to the organization that uploaded them, even when another organization has the same dashboard role.

Regression tests added:

- `TestUpdateExecutionStatusRequiresExecuteStepUpBeforeLeavingExecutionPending`
- `TestUpdateExecutionStatusAllowsPostExecutionProgression`
- `TestDocumentServiceUnscopedDocumentsAreScopedToUploaderOrg`
- `TestAuthorizationChallengeWrongActionCannotBeConsumed`

## Release Readiness Checklist

Required environment:

- `ADMIN_CONNECTION_STRING`
- `WALLET_DB_CONNECTION_STRING`
- `JWT_SECRET`
- `TROVO_WALLET_BASE_URL`
- `SERVICE_LINK_USERNAME`
- `SERVICE_LINK_API_KEY`
- `LOGIN_CALLBACK_URL`
- `ENABLE_CACHING=0` for MVP dashboards unless a separate cache validation is completed.

Migration order:

1. Confirm AdminDB has existing `organizations` and `organization_members` tables.
2. Apply `migrations/20260629_create_stakeholder_portal_tables.sql` to AdminDB.
3. Do not apply migrations to production from the app process; run through the normal DB migration/release procedure.

Data readiness:

- Stakeholder organizations that should log in must have `organizations.status = 'ACTIVE'`.
- Organization members that should log in must have `organization_members.status = 'ACTIVE'`.
- `organizations.stakeholder_type` must be one of `trustees`, `approved_asset_custodian`, or `asset_manager`.
- `organizations.stakeholder_id` must be populated for portal organizations.
- `organizations.type` is not used for portal authorization.
- Global uniqueness of organization member/invite email remains accepted for MVP.

Seed/assignment procedure:

1. Create or verify stakeholder organizations through the existing organization/stakeholder linkage flow.
2. Create trustee visibility with `POST /api/v1/stakeholder/admin/asset-assignments` using a Trovo admin token.
3. Include `asset_id`, `trustee_org_id`, and optionally `asset_code`, `asset_manager_org_id`, `custodian_org_id`.
4. Assets without trustee assignment are intentionally invisible to trustee views.

Rollback plan:

- Disable stakeholder routes by rolling back the application deploy.
- If the migration must be rolled back before production data exists, drop the new tables from `migrations/20260629_create_stakeholder_portal_tables.sql` in reverse dependency order.
- If production data exists, preserve table dumps before dropping or run a forward migration that disables route access while retaining audit/workflow records.

## Test Results

Baseline before edits:

- `go test ./...` passed.

After implementation:

- `go test ./...` passed after the first full slice.
- Final `go test ./...` passed on 2026-06-30 after the review fixes and runbook update.
- New tests cover role mapping, document role access, unscoped document org scoping, raw/Bearer token extraction, fund release status transitions, status-step-up bypass prevention, audit rollback, authorization challenge one-time consumption, authorization challenge expiry, wrong-action challenge rejection, and role-based asset visibility filtering.
- `STAKEHOLDER_PORTAL_TEST_VERIFICATION_FLOW.md` now contains DB-derived fixture discovery, exact JSON payload generation, ordered endpoint calls, dependency breakpoints, negative security checks, response evidence files, and SQL proof queries.

Known test gaps:

- There is no full HTTP integration test suite that exercises every registered stakeholder route.
- Full fund release happy path is covered by state-machine and primitive tests, but not yet by a single manager -> trustee -> custodian end-to-end test.
- Notification org scoping is implemented in service queries, but does not yet have a dedicated test.
- Due diligence, revenue/distribution/valuation, custodian account/compliance ownership, and dashboard aggregates need broader service/handler tests before claiming complete acceptance-test coverage.
- The DB-backed endpoint runbook has not been executed in this session because no local/staging DB credentials or real stakeholder/member tokens were provided.

Swagger generation:

- `go run github.com/swaggo/swag/cmd/swag@v1.16.2 init -g main.go` completed on 2026-06-30 and regenerated `docs/`.
- Existing non-stakeholder warnings remained: duplicate legacy tokenization/organization route annotations and an existing career-role example parsing warning.

## Known MVP Boundaries

- Migrations were created but not applied to any live database.
- Due diligence approval is local-only and does not mutate Trovo Wallet tokenization state.
- Document APIs store file URLs and metadata only; no file upload/storage provider is implemented.
- Dashboard Redis caching is not enabled because query shapes are new and MVP-safe without cache.
- Trovo admin stakeholder impersonation is not implemented; only the explicit admin asset-assignment endpoint exists.
- No silent tokenization state mutation is performed.

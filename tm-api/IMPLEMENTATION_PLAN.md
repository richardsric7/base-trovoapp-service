# Stakeholder Portal — Implementation Plan

## Overview

Build out the admin dashboard API to support role-based stakeholder views for Trustees, Asset Custodians, and Asset Managers. The existing codebase already handles organization management, member auth, wallet linkage, and stakeholder type linking — this plan builds on top of that.

**Suggested priority:** Phase 1 -> 2 -> 3 -> 5 -> 4 -> 6 -> 7
(Trustee first because they're the approval bottleneck; Asset Manager before Custodian because they initiate most workflows.)

---

## Phase 1: Foundation — RBAC & Stakeholder Type Routing

**Goal:** Gate endpoints by stakeholder type so each role only accesses their own pages.

### Tasks
- [ ] **1.1** Extend `OrganizationAuthMiddleware` to inject `stakeholder_type` into request context (read from org's linked stakeholder)
- [ ] **1.2** Create `RequireStakeholderType(types ...string)` middleware that checks context and returns 403 if role doesn't match
- [ ] **1.3** Add route groups per stakeholder type in the organizations controller:
  - `/api/v1/stakeholder/trustee/...`
  - `/api/v1/stakeholder/custodian/...`
  - `/api/v1/stakeholder/asset-manager/...`
  - `/api/v1/stakeholder/shared/...` (cross-functional)
- [ ] **1.4** Map existing stakeholder type constants to the three dashboard roles:
  - `trustees` -> Trustee dashboard
  - `approved_asset_custodian` -> Custodian dashboard
  - `asset_manager` -> Asset Manager dashboard

### Files to modify
- `internal/middleware/` — new role-checking middleware
- `internal/components/organizations/controllers/main.go` — new route groups

---

## Phase 2: Shared Tokenized Assets Module

**Goal:** Unified asset listing and detail view accessible to all three roles, with role-appropriate actions.

### Tasks
- [ ] **2.1** `GET /api/v1/stakeholder/shared/assets` — list tokenized assets
  - Paginated, filterable by asset class, status, date range, assigned stakeholder
  - Uses existing `trovosdk.GetTokenizedAssetData()` as data source
  - Response includes role-specific action flags (`canApprove`, `canProcessRelease`, `canUpdateStatus`)
- [ ] **2.2** `GET /api/v1/stakeholder/shared/assets/:asset_id` — asset detail
  - Full asset info: description, financials, operational status, token holders, documents, activity history
  - Role-specific action buttons in response metadata
- [ ] **2.3** Create `stakeholder_assets` service layer to fetch and enrich asset data per role

### New files
- `internal/components/stakeholder/` — new component directory
- `internal/components/stakeholder/handlers/asset_handler.go`
- `internal/components/stakeholder/services/asset_service.go`

---

## Phase 3: Trustee-Specific Endpoints

**Goal:** Dashboard, due diligence, fund management, and distribution authorization for Trustees.

### Tasks

#### 3A — Dashboard Home
- [ ] **3.1** `GET /api/v1/stakeholder/trustee/dashboard` — aggregate stats
  - Assets under trust count & total value
  - Pending milestone verifications
  - Pending fund release approvals
  - Pending distribution authorizations
  - Recent activity feed

#### 3B — Due Diligence & Verification
- [ ] **3.2** DB migration: `due_diligence_checklists` table (asset_id, category, item, status, notes, verified_by, verified_at)
- [ ] **3.3** `GET /api/v1/stakeholder/trustee/due-diligence/:asset_id` — get checklist for asset
- [ ] **3.4** `PUT /api/v1/stakeholder/trustee/due-diligence/:asset_id/item/:item_id` — mark item complete, add notes
- [ ] **3.5** `POST /api/v1/stakeholder/trustee/due-diligence/:asset_id/approve` — approve asset for custody
- [ ] **3.6** `POST /api/v1/stakeholder/trustee/due-diligence/:asset_id/reject` — reject with reasons

#### 3C — Fund Management & Release Approvals
- [ ] **3.7** DB migration: `fund_release_requests` table (id, asset_id, requester_org_id, amount, purpose, status, supporting_docs, reviewed_by, reviewed_at, notes)
- [ ] **3.8** `GET /api/v1/stakeholder/trustee/fund-releases` — list pending/all fund release requests
- [ ] **3.9** `GET /api/v1/stakeholder/trustee/fund-releases/:request_id` — request detail with docs
- [ ] **3.10** `POST /api/v1/stakeholder/trustee/fund-releases/:request_id/approve` — approve release
- [ ] **3.11** `POST /api/v1/stakeholder/trustee/fund-releases/:request_id/reject` — reject with reason

#### 3D — Income Distribution Management
- [ ] **3.12** DB migration: `distributions` table (id, asset_id, proposed_by_org_id, amount, source, scheduled_date, status, authorized_by, authorized_at)
- [ ] **3.13** `GET /api/v1/stakeholder/trustee/distributions` — list distribution proposals
- [ ] **3.14** `POST /api/v1/stakeholder/trustee/distributions/:dist_id/authorize` — authorize distribution
- [ ] **3.15** `POST /api/v1/stakeholder/trustee/distributions/:dist_id/reject` — reject with reason
- [ ] **3.16** `GET /api/v1/stakeholder/trustee/distributions/history` — past distributions

---

## Phase 4: Asset Custodian-Specific Endpoints

**Goal:** Account management, fund release execution, and compliance for Custodians.

### Tasks

#### 4A — Dashboard Home
- [ ] **4.1** `GET /api/v1/stakeholder/custodian/dashboard` — aggregate stats
  - Segregated accounts count & balances
  - Pending fund release requests
  - Reconciliation status
  - Compliance status

#### 4B — Segregated Accounts
- [ ] **4.2** DB migration: `segregated_accounts` table (id, asset_id, custodian_org_id, account_type, account_name, balance, status, bank_details, created_at)
- [ ] **4.3** `GET /api/v1/stakeholder/custodian/accounts` — list accounts with balances
- [ ] **4.4** `GET /api/v1/stakeholder/custodian/accounts/:account_id` — account detail with transaction history

#### 4C — Fund Release Processing
- [ ] **4.5** `GET /api/v1/stakeholder/custodian/fund-releases` — list Trustee-authorized releases pending execution
- [ ] **4.6** `POST /api/v1/stakeholder/custodian/fund-releases/:request_id/execute` — execute the release
- [ ] **4.7** `PUT /api/v1/stakeholder/custodian/fund-releases/:request_id/status` — update release status (processing, completed, failed)

#### 4D — Compliance
- [ ] **4.8** DB migration: `compliance_items` table (id, org_id, category, requirement, status, due_date, completed_at)
- [ ] **4.9** `GET /api/v1/stakeholder/custodian/compliance` — compliance checklist
- [ ] **4.10** `PUT /api/v1/stakeholder/custodian/compliance/:item_id` — update compliance item status

---

## Phase 5: Asset Manager-Specific Endpoints

**Goal:** Asset operations, revenue management, valuations, and reporting for Asset Managers.

### Tasks

#### 5A — Dashboard Home
- [ ] **5.1** `GET /api/v1/stakeholder/asset-manager/dashboard` — aggregate stats
  - Assets managed count & total value
  - Revenue collection status (current period, YTD)
  - Pending milestone completions
  - Pending Trustee approvals

#### 5B — Assets Under Management
- [ ] **5.2** `GET /api/v1/stakeholder/asset-manager/assets` — list managed assets with operational details
- [ ] **5.3** `PUT /api/v1/stakeholder/asset-manager/assets/:asset_id` — update asset operational info

#### 5C — Revenue Management
- [ ] **5.4** DB migration: `revenue_records` table (id, asset_id, manager_org_id, period, source, amount, status, collected_at)
- [ ] **5.5** `GET /api/v1/stakeholder/asset-manager/revenue` — revenue dashboard (current, YTD, by source)
- [ ] **5.6** `POST /api/v1/stakeholder/asset-manager/revenue` — record revenue collection
- [ ] **5.7** `POST /api/v1/stakeholder/asset-manager/revenue/submit-distribution` — submit revenue for distribution (creates distribution proposal for Trustee)

#### 5D — Valuations
- [ ] **5.8** DB migration: `asset_valuations` table (id, asset_id, manager_org_id, valuation, methodology, valuation_date, report_url, status)
- [ ] **5.9** `GET /api/v1/stakeholder/asset-manager/valuations/:asset_id` — valuation history
- [ ] **5.10** `POST /api/v1/stakeholder/asset-manager/valuations` — submit new valuation
- [ ] **5.11** `POST /api/v1/stakeholder/asset-manager/valuations/:id/request-independent` — request independent valuation

#### 5E — Fund Release Requests (from AM side)
- [ ] **5.12** `POST /api/v1/stakeholder/asset-manager/fund-releases` — create fund release request (goes to Trustee for approval)
- [ ] **5.13** `GET /api/v1/stakeholder/asset-manager/fund-releases` — view own requests and their status

#### 5F — Reporting
- [ ] **5.14** `POST /api/v1/stakeholder/asset-manager/reports` — generate report (financial, operational)
- [ ] **5.15** `GET /api/v1/stakeholder/asset-manager/reports` — list generated reports
- [ ] **5.16** `POST /api/v1/stakeholder/asset-manager/reports/:report_id/submit` — submit report to Trustee

---

## Phase 6: Cross-Functional Features

**Goal:** Document management, audit trail, communications, and profile settings.

### Tasks

#### 6A — Document Repository
- [ ] **6.1** DB migration: `documents` table (id, asset_id, uploaded_by_org_id, category, title, file_url, version, access_roles, created_at)
- [ ] **6.2** `POST /api/v1/stakeholder/shared/documents` — upload document
- [ ] **6.3** `GET /api/v1/stakeholder/shared/documents` — list documents (filtered by role access)
- [ ] **6.4** `GET /api/v1/stakeholder/shared/documents/:doc_id` — download/view document

#### 6B — Activity & Audit Trail
- [ ] **6.5** DB migration: `audit_logs` table (id, user_id, org_id, role, action, entity_type, entity_id, details, created_at)
- [ ] **6.6** Create audit logging helper — call from all state-changing handlers
- [ ] **6.7** `GET /api/v1/stakeholder/shared/audit-trail` — list activities (filterable by date, role, action type, asset)

#### 6C — Communications & Notifications
- [ ] **6.8** DB migration: `notifications` table (id, recipient_org_id, sender_org_id, type, title, message, related_entity_type, related_entity_id, read, created_at)
- [ ] **6.9** `GET /api/v1/stakeholder/shared/notifications` — list notifications
- [ ] **6.10** `PUT /api/v1/stakeholder/shared/notifications/:id/read` — mark as read
- [ ] **6.11** Auto-generate notifications on key actions (fund release approved, distribution authorized, milestone verified, etc.)

#### 6D — User Profile
- [ ] **6.12** `GET /api/v1/stakeholder/shared/profile` — get current user profile with org and stakeholder info
- [ ] **6.13** `PUT /api/v1/stakeholder/shared/profile` — update profile info
- [ ] **6.14** `PUT /api/v1/stakeholder/shared/profile/notifications` — update notification preferences

---

## Phase 7: Technical Hardening

- [ ] **7.1** Run all DB migrations in order
- [ ] **7.2** Add database indexes on frequently queried columns (asset_id, org_id, status, created_at)
- [ ] **7.3** Enforce 2FA (via existing wallet linkage auth flow) on: fund release approvals, distribution authorizations, fund release execution
- [ ] **7.4** Add Redis caching for dashboard aggregate queries
- [ ] **7.5** Ensure all list endpoints support pagination (page, limit, sort, filters)
- [ ] **7.6** Add Swagger annotations for all new endpoints
- [ ] **7.7** Add request validation on all POST/PUT endpoints
- [ ] **7.8** Security review: ensure no cross-role data leakage, validate org ownership on all entity access

---

## New Database Tables Summary

| Table | Phase | Purpose |
|---|---|---|
| `due_diligence_checklists` | 3 | Trustee verification items per asset |
| `fund_release_requests` | 3 | Fund release request workflow |
| `distributions` | 3 | Income distribution proposals and authorizations |
| `segregated_accounts` | 4 | Custodian-managed bank accounts |
| `compliance_items` | 4 | Compliance tracking per org |
| `revenue_records` | 5 | Revenue collection records |
| `asset_valuations` | 5 | Valuation history per asset |
| `documents` | 6 | Document repository |
| `audit_logs` | 6 | Activity and audit trail |
| `notifications` | 6 | Inter-stakeholder notifications |

---

## Estimated Endpoint Count

| Area | Endpoints |
|---|---|
| Shared (assets, docs, audit, profile) | ~10 |
| Trustee | ~12 |
| Asset Custodian | ~8 |
| Asset Manager | ~12 |
| **Total new endpoints** | **~42** |

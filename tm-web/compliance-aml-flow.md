# Compliance & AML Flow

*Internal reference — frontend + API surface. Scope: admin console + organisation portal. Backend: `stakeholder` service. As of 2026-09-04.*

How custodian compliance requirements move from creation in the admin console through resolution in the custodian portal — what's actually wired to a backend today, and what's still a static prototype.

## 1. The flow, in order

Compliance is a two-portal handoff: an admin assigns a requirement to a custodian org, and a member of that org resolves it from their own dashboard. The two sides only meet at the API — there's no shared screen.

- **Admin** — admin console, `/complaince-and-aml`
- **Custodian** — organisation portal, `/organisation/compliance`

| Step | Actor | What happens |
|---|---|---|
| 1 | Admin | Opens Compliance & AML from the sidebar. Visible to any admin user — the sidebar entry carries no role or permission check of its own. |
| 2 | Admin *(mock data)* | Lands on **Access Records**: 100 generated login rows (username, IP, auth method, status). Search, pagination and the filter drawer render, but nothing is querying live records. |
| 3 | Admin *(mock data)* | Switches to **Transaction Records**: also 100 generated rows. Clicking one opens a detail route whose Cleared / Flagged / Under Review status, comments and "Cleared By" line are all hardcoded — the row `id` in the URL is never used to fetch anything. |
| 4 | Admin *(wired)* | Switches to **Compliance Requirements** and picks a custodian. The only functional tab. The custodian dropdown loads every organisation and keeps the ones whose `type` matches an asset-custodian variant. |
| 5 | Admin | Fills in category, requirement, optional due date & asset. Formik + Yup validation: custodian, category and requirement text are required; due date and asset ID are optional. |
| 6 | Admin | Submits the form. On success the new item is prepended to a "Requirements Created This Session" table — held in local component state only, so it's gone on refresh.<br>`POST /stakeholder/admin/custodian-compliance` |
| 7 | Custodian | Sees **Compliance** in their org navigation. Only rendered for stakeholder types `asset_custodian` and `approved_asset_custodian` — every other org type (rating agency, issuing house, manager, trustee, legal/financial adviser) never gets this nav item. |
| 8 | Custodian | Opens the Records tab and loads their queue: category, requirement, due date, status and completed-at, paginated.<br>`GET /stakeholder/custodian/compliance?page=&limit=` |
| 9 | Custodian | Changes an item's status to resolve it. Status is an inline dropdown on each row. Updating it invalidates the shared cache tag, so the list and the dashboard's pending count refresh together.<br>`PUT /stakeholder/custodian/compliance/{itemId}` |

## 2. Where each side lives

**Admin console** — route `/complaince-and-aml`, three tabs, one of which talks to the backend.

| Tab | Status |
|---|---|
| Access Records | mock |
| Transaction Records | mock |
| Compliance Requirements | wired |

**Organisation portal** — route `/organisation/compliance`, gated by stakeholder type, not a permission flag.

| Tab | Status |
|---|---|
| Documents tab | n/a to this flow |
| Records tab | wired |

## 3. Create-requirement form

`CreateComplianceRequirement.tsx` — the form the admin fills in at step 5.

| Field | Input | Required? | Notes |
|---|---|---|---|
| `custodianOrgId` | select | required | Sourced from the organisations list, filtered to custodian-type orgs |
| `category` | text | required | e.g. "KYC", "AML Screening", "Custody Audit" |
| `requirement` | textarea | required | Free-text description of what the custodian must do |
| `dueDate` | date | optional | Sent as `due_date` if set |
| `assetId` | text | optional | Links the requirement to a specific tokenized asset |

## 4. API surface

Two independent RTK Query slices — the admin side on `baseApi`, the custodian side on `orgApi` — both backed by the same `stakeholder` service.

**Admin**

| Method | Path | Purpose |
|---|---|---|
| POST | `/stakeholder/admin/custodian-compliance` | Create a requirement and assign it to a custodian org |

No admin-side GET, list, or delete exists — creation is the only operation available from this side.

**Custodian**

| Method | Path | Purpose |
|---|---|---|
| GET | `/stakeholder/custodian/dashboard` | Dashboard summary, incl. pending-compliance count |
| GET | `/stakeholder/custodian/compliance` | Paginated list of the org's assigned requirements |
| PUT | `/stakeholder/custodian/compliance/{itemId}` | Update one item's status |

## 5. Compliance item & status

`ComplianceItem` mirrors the backend model (`stakeholder/models/db_models.go`) and is the shape returned by both the create call and the custodian list.

| Field | Type | Notes |
|---|---|---|
| `id` | string | — |
| `org_id` | string | The assigned custodian org |
| `category` | string | — |
| `requirement` | string | — |
| `status` | string | See enum below — typed loosely on the admin side |
| `due_date` | string? | — |
| `asset_id` | string? | Present only on the admin request type |
| `completed_at` | string? | — |
| `completed_by_member_id` | string? | — |

**Status values:** `pending` → `complete` · `overdue` · `waived`

Defined as `ComplianceStatus` on the custodian side; the admin's `status` field is a plain string with no shared enum import, and the create form never sets it — new items presumably default to `pending` on the backend.

## 6. What to know before building on this

**Access & Transaction Records are prototypes.** All three of the mock-data screens — Access Records, Transaction Records, and its detail page — render hardcoded arrays. Filters, exports and the "Generate Report" button have no handlers. Treat them as scaffolding for future work, not a working AML monitoring surface.

**Admin has no view into existing requirements.** There's no list/get endpoint on the admin side — only create. Requirements assigned in a past session are invisible to the admin unless they log in as the custodian.

**Status enum isn't shared.** The custodian side has a typed `ComplianceStatus` enum; the admin side's item type leaves `status` as a bare string. A drift between the two is easy to introduce silently.

**Custodian-type matching is a string list.** The create form filters organisations by checking `type` against four hardcoded casing variants of "asset custodian." A new casing or naming convention from the orgs API would silently drop valid custodians from the dropdown.

---
*Compiled from `frontend/src/app/(dashboard)/complaince-and-aml`, `frontend/src/app/organisation/compliance`, and the admin / asset-custodian redux API slices.*

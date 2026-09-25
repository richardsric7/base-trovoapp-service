# Trovotech Tokenization Stakeholder Portal — Summary

## What Is This?

A role-based admin dashboard for Trovotech's Asset Tokenization Platform (RATOP). Three types of stakeholders — **Trustee**, **Asset Custodian**, and **Asset Manager** — each get tailored dashboard pages to manage their responsibilities throughout the tokenization lifecycle of real-world assets (real estate, infrastructure, etc.).

---

## The 3 Roles

| Role | Purpose | Key Actions |
|---|---|---|
| **Trustee** | Protects investor interests. Holds assets in trust, controls fund releases, authorizes income distributions to token holders. The "approver" in every workflow. | Approve/reject fund releases, verify milestones, authorize distributions, conduct due diligence |
| **Asset Custodian** | Safeguards money. Maintains segregated bank accounts (sale proceeds, development funds, revenue). Only moves funds when Trustee approves. | Process fund releases, manage segregated accounts, reconciliation, compliance reporting |
| **Asset Manager** | Runs the assets day-to-day. Manages operations, maintenance, tenants, revenue collection. Reports upward to the Trustee. | Update asset status, report revenue, submit milestones, request fund releases, manage valuations |

---

## How They Interact

```
Asset Manager                    Trustee                     Asset Custodian
     |                              |                              |
     |-- requests fund release ---->|                              |
     |-- submits milestone -------->|                              |
     |                              |-- approves release --------->|
     |                              |                              |-- executes fund transfer
     |                              |                              |
     |-- reports revenue ---------->|                              |
     |                              |-- authorizes distribution -->|
     |                              |                              |-- distributes to token holders
```

---

## Dashboard Pages Per Role

### Trustee (8 pages)
1. **Dashboard Home** — stats overview (assets under trust, pending approvals, pending distributions)
2. **Assets Under Trust** — all assets held in trust (existing Tokenized Asset page)
3. **Due Diligence & Verification** — checklist-based verification before accepting assets into custody
4. **Fund Management & Release Approvals** — review and approve/reject fund release requests
5. **Income Distribution Management** — authorize distributions to token holders
6. ~~Compliance & Reporting~~ *(NOT IN SCOPE)*
7. ~~Asset Custody Verification~~ *(NOT IN SCOPE)*
8. **Communications & Notifications** — notification center and inter-stakeholder messaging

### Asset Custodian (7 pages)
1. **Dashboard Home** — account summaries, pending releases, reconciliation status
2. **Assets Under Trust** — assets view (same base data, custodian perspective)
3. **Segregated Accounts Management** — manage accounts by type (sale proceeds, development, revenue)
4. **Fund Release Processing** — execute Trustee-authorized fund releases
5. ~~Account Reconciliation~~ *(NOT IN SCOPE)*
6. **Compliance & Regulatory Reporting** — custody compliance tracking
7. **Communications** — notifications from Trustee and Asset Manager

### Asset Manager (8 pages)
1. **Dashboard Home** — assets managed, revenue status, pending milestones
2. **Assets Under Management** — operational view of managed assets
3. ~~Operational Management~~ *(NOT IN SCOPE)*
4. **Revenue Management & Collection** — track and report revenue
5. **Asset Valuation & Updates** — manage periodic valuations
6. ~~Milestone Management~~ *(NOT IN SCOPE)*
7. **Financial & Operational Reporting** — generate and submit reports to Trustee
8. **Communications** — notifications for approvals, rejections, authorizations

### Cross-Functional Pages (All Roles)
1. **Shared Tokenized Assets Module** — unified asset view with role-specific action buttons
2. **Document Management & Repository** — centralized document storage with role-based access
3. **Activity & Audit Trail** — full activity log for compliance
4. **User Profile & Settings** — profile, 2FA, notification preferences

---

## Key Technical Requirements

- **RBAC** with three distinct roles: `trustee`, `asset_custodian`, `asset_manager`
- **Role-specific data filtering** — users only see data relevant to their role
- **2FA** for high-privilege operations (fund approvals, distributions)
- **Audit logging** on all user actions
- **Pagination & caching** for large data sets
- **Secure session management** with timeouts

---

## What Already Exists in the Codebase

- Organization CRUD with stakeholder type linking
- Member authentication (JWT), invitation & onboarding flow
- Wallet linkage via Trovo SDK (2FA authorization flow)
- Trovo SDK integration (15+ methods: login, auth, user info, tokenized asset data, payments)
- Middleware for org member auth and Trovo admin auth
- Stakeholder types defined: `asset_manager`, `asset_issuing_house`, `approved_asset_custodian`, `legal_and_professionals`, `rating_agency`, `trustees`

# Integration: how tm-api fits into the monorepo

This describes how `tm-api` relates to the rest of `src-monorepo` — its
client (`tm-web`), its sibling backend (`app-backend`), and the
authorization model that gates every route. For the two-database
architecture at a glance, see [README.md](./README.md); for every route's
exact request/response shape, see Swagger UI (linked at the bottom).

## 1. tm-web is tm-api's only client

`tm-web` is the Next.js admin dashboard frontend. It is the only intended
caller of tm-api's `/api/v1/*` routes — tm-api is not a public API.

**The admin login flow (high level):**

1. An operator opens tm-web and submits their Trovo username/email. tm-web calls `POST /login`.
2. `Login` (`internal/components/auth/services/auth.go`) checks the user is both a known `AdminUser` (in `AdminDB`) and a known Trovo Wallet user (in the shared `TrovoWalletDB`), then calls out to **app-backend** via `trovosdk.ServiceLink.SendLoginRequest` — this is a *push-approval* flow: app-backend sends a login-approval push notification to the user's Trovo mobile app, not a password check.
3. The operator approves the login on their phone. app-backend then calls tm-api back at `POST /callbacks/login/{serviceName}` (see `LoginCallback`) with the approval.
4. Meanwhile, the browser is (or can be) subscribed to `GET /login/stream/{loginID}` (`LoginNotificationStream`), an SSE stream that `LoginCallback` broadcasts the approval onto — so the browser learns the login succeeded without polling.
5. Once approved, the browser holds a JWT (issued by app-backend as part of this flow) and sends it as `Authorization: Bearer <token>` on every subsequent request.

**How that JWT is verified on every protected route:** tm-api does **not**
decode/verify the admin JWT locally with a shared secret. Instead,
`middleware.JwtTokenAuthMiddleware` and `middleware.AuthenticateSuperAdmin`
(`internal/middleware/authentication_middleware.go`) call
`trovosdk.ServiceLink.JwtTokenVerify`, which makes an HTTP call to
**app-backend**'s `/v1/servicelinks/token/verify` endpoint on every request.
app-backend returns the token's claims (including the wallet user ID) if
valid; tm-api then looks that wallet user ID up in its own `AdminUser` table
(`AdminDB`) to confirm the caller is a known admin, and to resolve their role.

A **separate** token type — the *organization-member* JWT (for white-label
partner org members using the Organization Portal, not Trovo admins) — is
minted and verified locally by tm-api itself, signed with the shared
`JWT_SECRET` (see `internal/middleware/organization_auth_middleware.go`,
`stakeholder_auth_middleware.go`). This is a distinct auth path from the
Trovo-admin one above; some routes (`middleware.AllowOrgOrTrovoAdminNormalized`)
accept either.

## 2. Relationship to app-backend

`app-backend` is the primary Trovo Wallet / P2P backend and owns the shared
Postgres schema (`TrovoWalletDB`/`P2P`). tm-api relates to it in **two
different ways**, depending on the kind of data involved:

### 2a. Direct shared-DB writes — pure admin/config tables

For a small set of tables that are genuinely just **admin-configured
reference data**, with no app-backend business logic of its own to
preserve, tm-api reads *and writes* the shared database directly via GORM,
the same as it does for its own `AdminDB` tables. Concretely:

- **`curated_assets`** (`internal/models/curated_assets.go`) — the platform's asset catalog, including the `p2p_enabled` flag that gates whether an asset can be used to create a P2P offer. Managed via `PUT /assets/curated`, `PUT /assets/curated/{id}/p2p-enabled`, etc.
- **`service_links`** (`internal/models/service_links.go`) — white-label partner integration accounts (the credentials third parties use to call app-backend's own `/v1/servicelinks` API). app-backend has never exposed a create/edit endpoint for these — every row has historically been provisioned by a direct DB insert — so tm-api owns their lifecycle. Managed via `/service-links`.
- **Fee configs** (`internal/components/usermetrics/services/fees_handler.go`) — per-service-link fee configuration, managed via `/fee/configs`.
- **`kyc_configs`, `faucet_configs`, `doja_widgets`, `kyc_levels`** (`internal/models/configs.go`) — third-party KYC/faucet provider configuration and KYC-level definitions, managed via `/kyc/configs`, `/faucet/configs`, `/doja/widgets`, `/kyc/levels`.

These writes go straight through `s.TrovoWalletDB` with plain GORM calls —
there is no HTTP hop to app-backend for them.

### 2b. app-backend's own HTTP API — everything transactional

For anything that **is** real P2P/wallet business logic — payments,
authorization approvals, tokenized-asset actions, push notifications, or
anything that needs app-backend's own validation/side-effects — tm-api calls
app-backend's HTTP API instead of writing the tables itself. This goes
through `internal/trovosdk` (`ServiceLink`), which wraps app-backend's
`/v1/servicelinks/*` endpoints:

| trovosdk method | app-backend endpoint | Used for |
|---|---|---|
| `JwtTokenVerify` | `POST /v1/servicelinks/token/verify` | Verifying every incoming Trovo-admin JWT (see §1) |
| `JwtTokenRefresh` | `POST /v1/servicelinks/token/refresh` | `POST /login/token/refresh` |
| `JwtTokenDelete` | `DELETE /v1/servicelinks/token` | `POST /logout` |
| `SendLoginRequest` | `POST /v1/servicelinks/login/request/{trovoUser}` | Kicking off the push-approval login flow (§1) |
| `VerifyLoginRequest` | `GET /v1/servicelinks/login/verify/{serviceUsername}/{trovoUser}/{loginID}` | Polling/confirming a login approval |
| `SendAuthorizationRequest` / `VerifyAuthorizationRequest` | `POST/GET /v1/servicelinks/authorize/...` | Step-up authorization approvals (e.g. beneficiary/token-release authorization, used by the stakeholder portal's authorization challenges) |
| `GetTokenizedAssetAuthorizationRequest` | `POST /v1/servicelinks/authorize/tokenized-asset` | Tokenized-asset authorization flow |
| `SendEventRequest` | `POST /v1/servicelinks/events/request` | Generic event/notification requests |
| `SendPushNotification` | `POST /v1/servicelinks/{serviceUsername}/{trovoUser}/push` | Push notifications to a user's device |
| `GetPaymentData` | `GET /v1/servicelinks/payment/request/{trovoUser}` | Building a payment request (actual P2P payment execution) |
| `GetUserInfo` | `GET /v1/servicelinks/{serviceUsername}/{trovoUser}/userinfo` | Fetching a user's profile from app-backend's own view |
| `GetTokenizedAssetData` | `GET /v1/servicelinks/tokenized-asset/{assetCode}` | Reading tokenized-asset detail from app-backend |

Every one of these calls authenticates as the `service_links` row identified
by `SERVICE_LINK_USERNAME`/`SERVICE_LINK_API_KEY` (see
[CONFIGURATION.md](./CONFIGURATION.md)) — i.e. tm-api is, from app-backend's
point of view, just another service-link partner integration, with whatever
permission flags that row has been granted.

### 2c. `payment_histories` — read-only, direct-DB read

`GET /payment/history` (`internal/components/general/services/payment_history.go`)
reads the `payment_histories` table directly via `s.TrovoWalletDB`, the same
physical database app-backend owns and migrates — tm-api never writes or
migrates this table, only reads it. Each row is split into a **source**
side (what left `fromAddress`) and a **destination** side (what arrived at
`toAddress`), replacing the old flat `assetCode`/`contractAddress`/`amount`
fields:

```json
{
  "sourceNetwork": "base",
  "sourceContractAddress": "0x...",
  "sourceAssetCode": "USDC",
  "sourceAmount": "100.0000000",
  "destinationNetwork": "base",
  "destinationContractAddress": "0x...",
  "destinationAssetCode": "USDC",
  "destinationAmount": "100.0000000"
}
```

For a plain payment (everything today) the two sides are identical; a row
where they differ is a swap. The query/search filters follow the same
split (`destinationAssetCode`/`destinationContractAddress` and
`sourceAssetCode`/`sourceContractAddress`, replacing the old flat
`assetCode`/`contractAddress` filter names) — see
[app-backend's INTEGRATION.md](../app-backend/INTEGRATION.md#payment-history-source-vs-destination)
for the full rationale, since app-backend's payment-history-engine is the
table's only writer.

This endpoint, along with every other tm-api route, sits behind the global
rate limiter described in
[CONFIGURATION.md](./CONFIGURATION.md#8-rate-limiting).

## 3. The RBAC model

Three tables in `AdminDB` (`internal/models/admin_permissions.go`) define
tm-api's own authorization model, independent of app-backend:

- **`AdminUser`** — one row per admin-dashboard user, with a `Role` (`SUPER_ADMIN`, `EDIT_LEVEL_ADMIN`, `VIEW_ONLY_ADMIN`, or `SUSPENDED`) and a `Status`.
- **`AdminPermission`** — the fixed catalog of fine-grained permission names (e.g. `VIEW_ADMIN_LIST`, `MANAGE_TROVO_WALLET`, `ACCESS_REPORTS`, `VERIFY_DOCUMENTS`...), seeded on startup by `seedPermissions` in `internal/db/main.go`.
- **`RolePermission`** — the many-to-many mapping of which `Role` has which `AdminPermission`, also seeded by `seedPermissions`.

A route is gated in one of three ways:

1. **Any authenticated admin** — `middleware.JwtTokenAuthMiddleware(s.AdminDB)`. Confirms the caller is *some* known `AdminUser`, without checking role or a specific permission. This is the most common gate (most of `usermetrics`, `users`, etc).
2. **Super-admin only** — `middleware.AuthenticateSuperAdmin(s.AdminDB)`. Same JWT verification, plus an explicit `Role == SuperAdmin` check. Used for admin-management routes (`/admin/invite`, `/admin/remove`, `/admin/list`, the stakeholder portal's `/stakeholder/admin/*` routes, `/audit-trail`, etc).
3. **A specific named permission** — a handler calls `models.HasPermission(db, &adminUser, "PERMISSION_NAME")` itself after the base JWT check, rather than relying on role alone. The clearest example is the P2P Market Reports endpoints (`/p2p/reports/*`): `checkReportsAccess` (`internal/components/usermetrics/services/reports_handler.go`) checks the caller holds `ACCESS_REPORTS` — which `EDIT_LEVEL_ADMIN` and `SUPER_ADMIN` are seeded with, but `VIEW_ONLY_ADMIN` is not.

Separately, **organization members** (white-label partner org staff using
the Organization Portal, not Trovo admins) are gated by
`middleware.OrganizationAuthMiddleware` / `RequireActiveOrganization` /
`RequireStakeholderRole`, checked against their own JWT claims and the
`Organization`/`OrganizationMember` tables — a parallel, simpler model (an
org member has a single `member_role` string, not the
permission-table system above).

### Audit logging

Nearly every mutating admin action is wrapped in
`accesslog.Audit(db, event, category, targetFn)`
(`internal/components/accesslog/middleware.go`), inserted as route
middleware **after** the auth middleware and **before** the handler, e.g.:

```go
apiV1.PATCH("/admin/users/suspend", middleware.AuthenticateSuperAdmin(s.AdminDB),
    accesslog.Audit(s.AdminDB, models.EventUserSuspend, m, accesslog.BodyField("email")),
    adminServices.SuspendUser(s))
```

`Audit` records who made the call (resolved by the auth middleware via
`accesslog.SetActor`, from the caller's IP/user-agent/method/path), what kind
of event it was, and a "target" value extracted from the request (a body
field, a path param, or a custom function) — all written to the
`AdminAccessLog` table in `AdminDB`. The `/audit-trail` and
`/audit-trail/{id}` routes (super-admin only) expose this as a paginated,
filterable security log in the dashboard.

## 4. Swagger UI — the canonical per-endpoint reference

This document describes the *shape* of the integration; for the exact
request/response schema, required permission, and status codes of any
specific route, use the generated Swagger UI rather than reading the
handler source:

- **http://localhost:8082/swagger/index.html** locally (also mounted at `/api/v1/swagger/index.html`), or the equivalent path on a deployed instance.

See [DEPLOYMENT.md](./DEPLOYMENT.md) for how it's regenerated, and
[README.md](./README.md) for how to run tm-api locally to reach it.

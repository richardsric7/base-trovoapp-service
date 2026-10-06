# Integration

Everything `tm-api` connects to, and how. tm-api is the backend of Trovo
Manager; it works in app-backend's database and calls app-backend's API.

| Connects to | Direction | Through | Needed? |
|---|---|---|---|
| [tm-web](#1-tm-web) | tm-web → tm-api | HTTPS, `/api/v1/*` | yes |
| [app-backend's database](#2-app-backends-database) | tm-api reads and writes | Postgres | yes |
| [app-backend's API](#3-app-backends-api) | tm-api → app-backend, and app-backend → tm-api (login callbacks) | HTTPS with a service-link API key | yes |
| [payout-engine and the Public Markets engine](#4-payout-engine-and-the-public-markets-engine) | tm-api → engines | the database, plus Redis wake-ups | for payouts and Public Markets |
| [Redis](#5-redis) | tm-api reads and writes | Redis | recommended |
| [Vault](#6-vault) | tm-api reads and writes | Vault HTTP API | for the vault manager |
| [Base blockchain](#7-base-blockchain) | tm-api reads and sends | JSON-RPC | for the vault manager and Public Markets contracts |
| [Mailgun](#8-mailgun) | tm-api → Mailgun | HTTPS API | for email |
| [Monitoring tools](#9-monitoring-tools) | tm-api → Prometheus, Loki, GlitchTip | HTTP | optional |
| [Discord and IP lookup](#10-discord-and-ip-lookup) | tm-api → services | webhooks, HTTPS | optional |

---

## 1. tm-web

- **What it is and why:** the Trovo Manager website (Next.js). It is the
  only intended caller of tm-api's `/api/v1/*` routes; tm-api is not a
  public API.
- **Direction:** tm-web → tm-api.
- **How they connect:** HTTPS. Every request after login carries
  `Authorization: Bearer <token>`.
- **Settings on this side:** [`CORS_ALLOWED_ORIGINS`](CONFIGURATION.md#cors_allowed_origins)
  (example `https://manager.trovo.example.com`; unset allows any origin).
- **Settings on the other side:** tm-web's `NEXT_PUBLIC_API_BASE_URL`, set
  to tm-api's address + `/api/v1` (see [tm-web/CONFIGURATION.md](../tm-web/CONFIGURATION.md)).
- **How to check it works:** log in to tm-web; the dashboard loads data.
- **When tm-api is down:** tm-web shows errors; nothing else is affected.

### How admins log in

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

**How the admin token is checked on every request:** tm-api does not check
the admin token itself. `middleware.JwtTokenAuthMiddleware` and
`middleware.AuthenticateSuperAdmin` (`internal/middleware/authentication_middleware.go`)
ask app-backend (`POST /v1/servicelinks/token/verify`) on every request.
app-backend returns the token's claims; tm-api then looks the user up in its
`AdminUser` table to confirm they are an admin and find their role.

A **separate** token type — the *organization-member* JWT (for white-label
partner org members using the Organization Portal, not Trovo admins) — is
minted and verified locally by tm-api itself, signed with the shared
`JWT_SECRET` (see `internal/middleware/organization_auth_middleware.go`,
`stakeholder_auth_middleware.go`). This is a distinct auth path from the
Trovo-admin one above; some routes (`middleware.AllowOrgOrTrovoAdminNormalized`)
accept either.

### Who may do what (roles and permissions)

Three tables in `AdminDB` (`internal/models/admin_permissions.go`) define
tm-api's own authorization model, independent of app-backend:

- **`AdminUser`** — one row per admin-dashboard user, with a `Role` (`SUPER_ADMIN`, `EDIT_LEVEL_ADMIN`, `VIEW_ONLY_ADMIN`, or `SUSPENDED`) and a `Status`.
- **`AdminPermission`** — the fixed catalog of fine-grained permission names (e.g. `VIEW_ADMIN_LIST`, `MANAGE_TROVO_WALLET`, `ACCESS_REPORTS`, `VERIFY_DOCUMENTS`...), seeded on startup by `seedPermissions` in `internal/db/main.go`.
- **`RolePermission`** — the many-to-many mapping of which `Role` has which `AdminPermission`, also seeded by `seedPermissions`.

A route is gated in one of three ways:

1. **Any authenticated admin** — `middleware.JwtTokenAuthMiddleware(s.AdminDB)`. Confirms the caller is *some* known `AdminUser`, without checking role or a specific permission. This is the most common gate (most of `usermetrics`, `users`, etc).
2. **Super-admin only** — `middleware.AuthenticateSuperAdmin(s.AdminDB)`. Same JWT verification, plus an explicit `Role == SuperAdmin` check. Used for admin-management routes (`/admin/invite`, `/admin/remove`, `/admin/list`, the stakeholder portal's `/stakeholder/admin/*` routes, `/audit-trail`, etc).
3. **A specific named permission** — a handler calls `models.HasPermission(db, &adminUser, "PERMISSION_NAME")` itself after the base JWT check, rather than relying on role alone. The clearest example is the P2P Market Reports endpoints (`/p2p/reports/*`): `checkReportsAccess` (`internal/components/usermetrics/services/reports_handler.go`) checks the caller holds `ACCESS_REPORTS` — which `EDIT_LEVEL_ADMIN` and `SUPER_ADMIN` are seeded with, but `VIEW_ONLY_ADMIN` is not.

Public Markets (`/public-markets`) uses this third pattern with its own
permissions:
- `MANAGE_PUBLIC_MARKETS` (SUPER_ADMIN, EDIT_LEVEL_ADMIN) to change anything;
- `MANAGE_SETTINGS` for its settings and partners;
- `VIEW_PUBLIC_MARKETS_PII` (SUPER_ADMIN) to see exchange customers' legal names and tax IDs.

On top of these, batch and dividend approvals are limited to the admins
listed in the settings' approver lists.

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

## 2. app-backend's database

- **What it is and why:** app-backend's Postgres database. Trovo Manager
  shows and changes app-backend's data there, and keeps its own admin
  tables (admins, roles, permissions, access logs) in the same database.
- **Direction:** tm-api reads and writes.
- **How they connect:** one direct database connection.
- **Settings on this side:** [`ADMIN_CONNECTION_STRING`](CONFIGURATION.md#admin_connection_string)
  (example `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`).
- **Settings on the other side:** the same database as app-backend's
  `DB_CONNECTION_STRING` ([app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)).
- **How to check it works:** `/ready` reports the database as up.
- **When it is down:** Trovo Manager cannot load or change anything.

**Which data tm-api writes directly.** For data that is purely admin
configuration (no app-backend business logic to go through), tm-api writes
the tables itself:

- **`curated_assets`** (`internal/models/curated_assets.go`) — the platform's asset catalog, including the `p2p_enabled` flag that gates whether an asset can be used to create a P2P offer and the `gas_fee_eligible` flag that lets app users pay network fees in it (through the paymaster; the asset must also be enabled on the paymaster contract and priced by its quote service). Managed via `PUT /assets/curated`, `PUT /assets/curated/{id}/p2p-enabled`, `PUT /assets/curated/{id}/gas-fee-eligible`, etc.
- **`service_links`** (`internal/models/service_links.go`) — white-label partner integration accounts (the credentials third parties use to call app-backend's own `/v1/servicelinks` API). app-backend has never exposed a create/edit endpoint for these — every row has historically been provisioned by a direct DB insert — so tm-api owns their lifecycle. Managed via `/service-links`. Includes `rateLimitPerMinute` (0 = no override): a per-partner rate-limit budget enforced by app-backend's `RateLimitMiddleware` on every request that service link makes, overriding that route's shared default — set this when one partner's legitimate traffic needs a different budget than everyone else's. See app-backend's CONFIGURATION.md "Rate limiting" section for the enforcement side.
- **Fee configs** (`internal/components/usermetrics/services/fees_handler.go`) — per-service-link fee configuration, managed via `/fee/configs`.
- **`fee_exempt_users`** (`internal/models/fee_exempt.go`) — accounts that pay no platform service fees (swap, payment, patron, account recovery, sub-wallet creation, tokenization application, closed group): the platform's own trading or operations accounts. app-backend reads the table whenever it charges a fee, so a change applies to the account's next request. Managed via `GET/POST /fee/exempt-users` and `DELETE /fee/exempt-users/:username` (Trovo admins only: organization members get 403; every change is audited as `fee_exemption.add` / `fee_exemption.remove`, with the admin's email stored as `addedBy`). The username must be an existing Trovo account. The tokenization issuing profile is always exempt and does not need to be listed.
- **Proceeds payouts** (`internal/components/proceedpayouts/`): `proceed_payouts`, their schedules (`tokenized_asset_payout_schedules`), approvals, batches, `payout_engine_states` and the `PROCEED_PAYOUT_FEE` row of `service_fees`. A trustee authorizing a stakeholder distribution registers its payout (`proceedpayouts.DistributionClient`, the stakeholder portal's `DistributionPayoutClient`), and the distribution then follows the payout's status. Trovo admins drive the payout through `/proceed-payouts`:
  - prepare / re-prepare;
  - set the payout's fee (FIXED, or PERCENT with a cap; VAT at the asset country's rate is charged on it);
  - approve, by `PROCEED_PAYOUT_APPROVALS_REQUIRED` distinct admins (neither the preparer nor the fee setter);
  - reject, confirm funding, pause, resume, cancel, retry failed;
  - exclude / include / mark paid a holder;
  - the engine's kill switch and sweep;
  - the fee configuration;
  - the payouts and payout fee / VAT reports.

  All of these are Trovo admins only and audited (`payout.*`, `payout_item.*`, `payout_engine.*`). Each change is a conditional status update, after which tm-api publishes a wake-up on Redis `payout-engine:commands`. `payout-engine` (no HTTP) does the work; see its [INTEGRATION.md](../payout-engine/INTEGRATION.md).
- **Public Markets** (`internal/components/publicmarkets/`): the `public_market_*` tables and their companions (`approved_dealing_members`, `custodian_positions`, `beneficial_ownership_*`, `price_oracle_snapshots`, `wallet_provisioning_requests`). app-backend runs the engine and owns the tables; the models mirror app-backend's `publicmarkets/models/models.go` and are never migrated here (see app-backend's [PUBLIC_MARKETS.md](../app-backend/PUBLIC_MARKETS.md)). Trovo admins operate it through `/public-markets`:
  - **Assets:** add (in SETUP), edit, register the token contract, take live, halt, resume, set a manual price, record the Custodian's position. A contract is verified on Base before it is accepted: the issuing Safe must be its owner and nothing may be minted yet. This needs `BASE_RPC_URL`.
  - **Orders and batches:** approve or reject a net batch above the threshold (Net Creation Approvers only). Roll an unexecuted batch to the next session.
  - **Escalations:** retry an instruction, mark it handled, or record the partner's outcome (a fill or a settlement) for a MANUAL partner.
  - **Reconciliation:** see the latest runs, run reconciliation now, ask a mock Custodian for its position feed.
  - **Corporate actions:** declare, approve a distribution (Dividend Approvers only; the approval binds to the snapshot's checksum), cancel, and follow exchanges' `dividend.paid` confirmations.
  - **Exchange partners:**
    - onboard a verified service link (the signing secret is shown once);
    - edit it — the rate-limit tier sets the service link's `rate_limit_per_minute`;
    - rotate the secret (the old one is valid for 24 hours);
    - suspend or activate;
    - pay a balance back;
    - replay dead-lettered webhooks.
  - **Wallet provisioning:** legal names and tax IDs are masked without `VIEW_PUBLIC_MARKETS_PII`.
  - **Prices, system health and settings:** settings include thresholds, approvers, fees, withholding tax, market hours and holidays, rate-limit tiers, Dealing Members and Custodian integrations.

  tm-api changes rows with conditional updates. Work for the engine goes through two tables:
  - `public_market_job_requests`: RECONCILE, RESUME, POSITION_FEED and EXCHANGE_WITHDRAWAL. The engine checks every minute.
  - unprocessed `public_market_partner_events` with source `MANUAL:<admin email>`: recorded outcomes and declared corporate actions. The engine processes these within seconds.

  Reads need any Trovo admin. Changes need `MANAGE_PUBLIC_MARKETS`, and settings and partners need `MANAGE_SETTINGS`. Every change is audited as `public_markets.*`.
- **`kyc_configs`, `faucet_configs`, `doja_widgets`, `kyc_levels`** (`internal/models/configs.go`) — third-party KYC/faucet provider configuration and KYC-level definitions, managed via `/kyc/configs`, `/faucet/configs`, `/doja/widgets`, `/kyc/levels`.

**Payment history** is read only:

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

## 3. app-backend's API

- **What it is and why:** for anything that is real wallet business logic
  (admin logins, payments, authorizations, push notifications, tokenized
  asset actions), tm-api asks app-backend instead of writing tables, so
  app-backend's checks and side effects apply.
- **Direction:** tm-api → app-backend, through `internal/trovosdk`
  (`ServiceLink`), and app-backend → tm-api for login and authorization
  callbacks.
- **How they connect:** HTTPS. tm-api identifies itself as a **service
  link**: the service link's owner username in the address and its API key
  in the `X-TW-SERVICE-LINK-API-KEY` header. From app-backend's point of view
  tm-api is one more service-link partner, with whatever permissions its
  `service_links` row has.
- **Settings on this side:**

  | Parameter | Example |
  |---|---|
  | [`TROVO_WALLET_BASE_URL`](CONFIGURATION.md#trovo_wallet_base_url) | `https://api.trovo.example.com` |
  | [`SERVICE_LINK_USERNAME`](CONFIGURATION.md#service_link_username) | `trovomanager` |
  | [`SERVICE_LINK_API_KEY`](CONFIGURATION.md#service_link_api_key) | `openssl rand -hex 24` |
  | [`LOGIN_CALLBACK_URL`](CONFIGURATION.md#login_callback_url) | `https://manager-api.trovo.example.com/api/v1/callbacks/login` |
  | [`NATIVE_ASSET_CODE`](CONFIGURATION.md#native_asset_code) | `ETH` |

- **Settings on the other side:** a `service_links` row in app-backend's
  database with that owner username and API key, and the login,
  authorization, payment, push notification, user info, token info and
  tokenized-asset authorization permissions (created as in
  [DEPLOYMENT.md step 3](DEPLOYMENT.md#3-create-trovo-managers-service-link-once)).
  app-backend must be able to reach `LOGIN_CALLBACK_URL`.
- **How to check it works:** an admin login request reaches the admin's
  phone, and approving it logs them in within a second.
- **When it is down:** admins cannot log in, and already logged-in admins
  are refused (their tokens cannot be checked).

The calls tm-api makes:

| trovosdk method | app-backend endpoint | Used for |
|---|---|---|
| `JwtTokenVerify` | `POST /v1/servicelinks/token/verify` | Verifying every incoming Trovo-admin JWT (see [How admins log in](#how-admins-log-in)) |
| `JwtTokenRefresh` | `POST /v1/servicelinks/token/refresh` | `POST /login/token/refresh` |
| `JwtTokenDelete` | `DELETE /v1/servicelinks/token` | `POST /logout` |
| `SendLoginRequest` | `POST /v1/servicelinks/login/request/{trovoUser}` | Kicking off the push-approval login flow ([How admins log in](#how-admins-log-in)) |
| `VerifyLoginRequest` | `GET /v1/servicelinks/login/verify/{serviceUsername}/{trovoUser}/{loginID}` | Polling/confirming a login approval |
| `SendAuthorizationRequest` / `VerifyAuthorizationRequest` | `POST/GET /v1/servicelinks/authorize/...` | Step-up authorization approvals (e.g. beneficiary/token-release authorization, used by the stakeholder portal's authorization challenges) |
| `GetTokenizedAssetAuthorizationRequest` | `POST /v1/servicelinks/authorize/tokenized-asset` | Tokenized-asset authorization flow |
| `SendEventRequest` | `POST /v1/servicelinks/events/request` | Generic event/notification requests |
| `SendPushNotification` | `POST /v1/servicelinks/{serviceUsername}/{trovoUser}/push` | Push notifications to a user's device |
| `GetPaymentData` | `GET /v1/servicelinks/payment/request/{trovoUser}` | Building a payment request (actual P2P payment execution) |
| `GetUserInfo` | `GET /v1/servicelinks/{serviceUsername}/{trovoUser}/userinfo` | Fetching a user's profile from app-backend's own view |
| `GetTokenizedAssetData` | `GET /v1/servicelinks/tokenized-asset/{assetCode}` | Reading tokenized-asset detail from app-backend |

## 4. payout-engine and the Public Markets engine

- **What it is and why:** background engines that do the work admins
  request: payout-engine pays proceeds; app-backend's Public Markets engine
  runs reconciliation, jobs and partner outcomes.
- **Direction:** tm-api → engines. tm-api changes rows (conditional status
  updates, see the proceeds payouts and Public Markets entries in
  [section 2](#2-app-backends-database)); for payouts it also publishes a
  wake-up on Redis channel `payout-engine:commands`.
- **Settings on this side:** the [Redis settings](CONFIGURATION.md#7-redis-cache-wake-ups-and-rate-limiting)
  (`ENABLE_CACHING=1`), [`PROCEED_PAYOUT_APPROVALS_REQUIRED`](CONFIGURATION.md#proceed_payout_approvals_required).
- **Settings on the other side:** payout-engine's Redis settings must point
  to the same Redis ([payout-engine/CONFIGURATION.md](../payout-engine/CONFIGURATION.md)).
- **How to check it works:** pressing **Prepare** on a payout changes its
  status within a second (within 10 seconds without Redis).
- **When an engine is down:** admin actions are recorded and done when it
  comes back.

## 5. Redis

- **What it is and why:** a shared in-memory store. tm-api uses it for the
  response cache, the rate limit, payout-engine wake-ups, passing login
  notifications between copies of tm-api, and a lock that makes vault
  manager key changes run one at a time.
- **Settings on this side:** [`ENABLE_CACHING`](CONFIGURATION.md#enable_caching),
  [`REDIS_HOST`](CONFIGURATION.md#redis_host),
  [`REDIS_PORT`](CONFIGURATION.md#redis_port),
  [`REDIS_PASSWORD`](CONFIGURATION.md#redis_password).
- **How to check it works:** `/ready` reports the cache as up.
- **When it is down or off:** caching and the rate limit stop, payout-engine
  notices actions on its next check, and with more than one copy of tm-api
  logins fall back to polling. Run only one copy without Redis.

## 6. Vault

- **What it is and why:** HashiCorp Vault holds platform signing keys
  managed by the vault manager (Vault Signer) and users' personal secrets.
- **Direction:** tm-api reads and writes secrets.
- **Settings on this side:** [`VAULT_ADDR`](CONFIGURATION.md#vault_addr),
  [`VAULT_TOKEN`](CONFIGURATION.md#vault_token),
  [`PERSONAL_ENVS`](CONFIGURATION.md#personal_envs),
  [`PERSONAL_ENV_VAULT_MOUNT`](CONFIGURATION.md#personal_env_vault_mount),
  [`PERSONAL_ENV_VAULT_PATH_PREFIX`](CONFIGURATION.md#personal_env_vault_path_prefix).
- **Settings on the other side:** a Vault policy for tm-api's token; services
  that read the managed secrets (payout-engine, app-backend) need their own
  read-only tokens.
- **When it is down:** the vault manager's pages answer 503; everything else
  works.

## 7. Base blockchain

- **What it is and why:** tm-api changes Safe owners when the vault manager
  rotates a key, and checks token contracts when Public Markets assets are
  registered.
- **Settings on this side:** [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url),
  [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id),
  [`RPC_TIMEOUT`](CONFIGURATION.md#rpc_timeout).
- **When it is down:** key rotations and contract registrations fail with an
  error; nothing else is affected.

## 8. Mailgun

- **What it is and why:** sends Trovo Manager's emails (admin and
  organization invites, verification codes, notifications).
- **Settings on this side:** [`MAILGUN_DOMAIN`](CONFIGURATION.md#mailgun_domain),
  [`MAILGUN_PRIVATE_API_KEY`](CONFIGURATION.md#mailgun_private_api_key),
  [`MAILGUN_API_KEY`](CONFIGURATION.md#mailgun_api_key),
  [`MAIL_SENDER`](CONFIGURATION.md#mail_sender) and the link settings
  [`TROVO_MANAGER_BASE_URL`](CONFIGURATION.md#trovo_manager_base_url),
  [`ADMIN_DASHBOARD_URL`](CONFIGURATION.md#admin_dashboard_url).
- **Settings on the other side:** a verified sending domain in Mailgun.
- **When it is down:** emails are not sent; the actions themselves still
  happen.

## 9. Monitoring tools

- **What it is and why:** Trovo Manager's system health screens read from
  Prometheus (metrics), Loki (logs) and GlitchTip (crash reports), and poll
  each backend's health. tm-api also reports its own crashes to GlitchTip
  (or Sentry).
- **Direction:** tm-api → each tool, inside your server network (browsers
  never call them).
- **Settings on this side:** [section 10 of CONFIGURATION.md](CONFIGURATION.md#10-system-health-screens-observability).
- **When unset:** each screen says "not configured".

## 10. Discord and IP lookup

- **Discord:** error and warning alerts, each with a built-in fallback
  webhook ([section 13 of CONFIGURATION.md](CONFIGURATION.md#13-discord-alerts)).
- **IP lookup:** adds a country to user records
  ([`IPAPI_HOST`](CONFIGURATION.md#ipapi_host), [`IPAPI_KEY`](CONFIGURATION.md#ipapi_key));
  skipped when unset.

## API reference (Swagger)

For the exact request and response of every route, required permission and
status codes, use the generated reference:
<http://localhost:8082/swagger/index.html> locally (also at
`/api/v1/swagger/index.html`), or the same path on a deployed instance. See
[DEPLOYMENT.md](DEPLOYMENT.md#swagger-api-reference) for how it is
regenerated.

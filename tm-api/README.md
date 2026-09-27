# tm-api

`tm-api` is the admin backend behind **Trovo Wallet's admin dashboard**. It is
a Go REST API (Gin + GORM) that the dashboard frontend, **`tm-web`**, talks to
for everything an operator or support agent does: managing admin accounts and
permissions, suspending/reactivating users, curating the platform's asset
catalog, managing fee configs, reviewing tokenization requests, running the
stakeholder-portal workflows, and more.

Go module name: `admin-panel-dashboard` (the historical name predates the
`tm-api` repo name — you'll see it in import paths throughout the codebase).

If you are new to this codebase, read this file, then:

- **[CONFIGURATION.md](./CONFIGURATION.md)** — every environment variable, what it does, and how to get a real value for it.
- **[DEPLOYMENT.md](./DEPLOYMENT.md)** — prerequisites, building, running locally, Docker, and CI/CD.
- **[INTEGRATION.md](./INTEGRATION.md)** — how tm-api fits with `tm-web` and `app-backend`, the auth flow, and the RBAC/audit-log model.
- **Swagger UI** — the live, per-endpoint API reference (see below).

## Where tm-api sits in the monorepo

```
src-monorepo/
├── tm-api/          <- you are here (this admin backend)
├── tm-web/          <- the admin dashboard frontend (tm-api's only client)
├── app-backend/      <- the main Trovo Wallet / P2P backend
├── app-mobile/       <- the Trovo Wallet mobile app
├── app-web/           <- the end-user wallet web app
├── payment-history-engine/
└── wallet-core/
```

tm-api is **not** public-facing — it is only ever called by tm-web (an
internal admin tool), and it in turn calls out to `app-backend`'s own HTTP
API for certain operations (see [INTEGRATION.md](./INTEGRATION.md)).

## The two-database architecture

tm-api talks to two logically distinct databases, exposed on the `Server`
struct (`internal/server/models/server.go`) as three `*gorm.DB` handles:

| Field | Owns | Notes |
|---|---|---|
| `s.TrovoWalletDB` / `s.P2P` | Shared with **app-backend** | The Trovo Wallet / P2P Postgres schema that `app-backend` is the primary owner of (`users`, P2P orders/offers, tokenized assets, payment history, etc). tm-api **reads** most of this data, and **writes directly** to a specific, deliberately small set of pure admin/config tables that live in this same schema but that app-backend has never exposed a management endpoint for — e.g. `curated_assets` (the asset catalog and its `p2p_enabled` flag), `service_links` (white-label partner integration accounts), `fee configs`, `kyc_configs`, `faucet_configs`, `doja_widgets`, and `kyc_levels`. For everything else — anything that is actual P2P/wallet *business logic* (payments, order matching, KYC verification, etc) — tm-api calls app-backend's own HTTP API instead of touching the tables. See [INTEGRATION.md](./INTEGRATION.md) for concrete examples. |
| `s.AdminDB` | **tm-api's own** database | Admin-only concerns that have no reason to live in app-backend's schema: admin user accounts (`AdminUser`), roles/permissions (`AdminPermission`, `RolePermission`), suspension history, the audit/access log (`AdminAccessLog`), organization/stakeholder-portal data, vault-signer secrets, career listings, and configuration seed data. tm-api owns this schema outright and runs GORM `AutoMigrate` against it on every startup (see `internal/db/main.go`). |

**Current deployment note:** as of this snapshot of the code, both handles
are wired to the *same* physical Postgres connection string
(`ADMIN_CONNECTION_STRING`) — see the comment at the top of
`internal/db/main.go` and in `main.go`'s `AdminDB()` setup. `TrovoWalletDB`,
`P2P` and `AdminDB` are kept as separate variable names (rather than
threading one name through every call site) purely so the distinction above
stays visible in the code, and so a future split back into two real
connections is a small, localized change rather than a rewrite. Don't let
this collapse mislead you: the *logical* ownership split above (which table
belongs to which side, and which side you write to directly vs. call
app-backend for) still applies and is what matters when you're deciding
where a new admin feature's data should live.

## Tech stack

- **Language:** Go 1.23 (toolchain 1.24.3 per `go.mod`)
- **Web framework:** [Gin](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io) with the Postgres driver (`gorm.io/driver/postgres`); `gorm.io/driver/sqlite` is also wired in for local dev/testing without Postgres
- **Auth:** JWT (`golang-jwt/jwt`), verified for Trovo-admin tokens by calling out to app-backend rather than decoding locally — see [INTEGRATION.md](./INTEGRATION.md)
- **API docs:** [swaggo/swag](https://github.com/swaggo/swag) + [gin-swagger](https://github.com/swaggo/gin-swagger)
- **Cache:** Redis (`go-redis/redis/v8`), optional (`ENABLE_CACHING`)
- **Secrets:** HashiCorp Vault (`hashicorp/vault/api`), for the vault-signer module
- **Error/crash reporting:** Sentry-compatible GlitchTip (`getsentry/sentry-go`), plus Discord webhooks for legacy error alerts
- **Email:** Mailgun (`mailgun-go`)

## Directory structure

```
tm-api/
├── main.go                      # Entry point: env/DB setup, route wiring, server start
├── config/                      # envconfig-based config loader (config.Load())
├── docs/                        # Generated Swagger spec (docs.go, swagger.json, swagger.yaml) — do not hand-edit
├── migrations/                  # Historical, hand-run SQL migrations (schema is otherwise managed by GORM AutoMigrate)
├── payloads/                    # Sample/reference request payloads
├── scripts/                     # Build/CI helper scripts (e.g. the swag YAML-structure fix)
└── internal/
    ├── components/               # One package per feature area (see below)
    ├── middleware/                # Auth middleware: JWT verification, org-member auth, stakeholder auth, CORS
    ├── models/                    # Shared GORM models (User, AdminUser, CuratedAsset, ServiceLink, configs, ...)
    ├── db/                        # AdminDB() connection + AutoMigrate + permission seeding
    ├── server/                    # Server struct (the three DB handles + globals), HTTP server start/shutdown, response helpers
    ├── trovosdk/                  # HTTP client (ServiceLink) for calling app-backend's /v1/servicelinks/* API
    ├── mail/, cache/, network/, observe/, errors/, evmkeypair/, gnosissafe/, utils/
    └── ...
```

### `internal/components/*` — feature areas

Each has its own `controllers/` (route registration in `main.go`) and usually
`services/`, `handlers/`, `db/` and/or `models/` subpackages.

| Component | What it does |
|---|---|
| `auth` | Admin login (push-approval flow via the Trovo mobile app), login/auth callbacks, token refresh, logout |
| `general` | Admin invite/suspend/remove, user suspend/lift-suspension, system configuration CRUD, wires in `observability` |
| `observability` | Backend-for-frontend for the "Trovo Manager" observability screens (Prometheus/Loki/GlitchTip/Grafana), queried server-side so the browser never holds monitoring credentials |
| `accesslog` | The audit/access-log: `Audit()` middleware that logs every privileged mutation, plus the `/audit-trail` read endpoints |
| `admin_invite` | Admin-list DB helper used by `general`'s `ListAdminsHandler` |
| `users` | Direct wallet-user detail/update/KYC-update/online-toggle endpoints (talk straight to the shared TrovoWalletDB/P2P schema) |
| `usermetrics` | The largest component: platform metrics/dashboards, P2P reports, curated-asset & service-link & fee-config CRUD, tokenization review/mint workflow, faucet/KYC/doja-widget/KYC-level configs, minting-user grants, partner & country management |
| `organizations` | Organization (white-label partner org) onboarding, member invites, member login/logout, wallet-linking |
| `stakeholder` | The Stakeholder Portal: a large, role-gated sub-app (trustee, custodian, asset manager, legal/financial adviser, issuing house, rating agency) for managing tokenized real-world assets end to end — compliance, due diligence, fund releases, distributions, revenue, valuations, reporting |
| `careers` | Public + admin CRUD for the careers/job-listings page |
| `vaultsigner` | Self-service and admin management of Vault-backed signer secrets and "personal env" secrets |
| `health` | `/health` (liveness) and `/ready` (readiness, with dependency checks) |
| `swagger` | Mounts the Swagger UI and a `/ping` smoke-check |
| `root` | The `/` service-info endpoint |

## Running locally

```bash
# 1. Clone and enter the repo
git clone <repo-url> && cd tm-api

# 2. Install Go 1.23+ (see DEPLOYMENT.md if you need install links)

# 3. Copy the env template and fill in real values (see CONFIGURATION.md)
cp .env-sample .env

# 4. Make sure ADMIN_CONNECTION_STRING in .env points at a reachable Postgres
#    database (or set DB_TYPE=sqlite and point it at a local file for a
#    quick, walletDB-less smoke test)

# 5. Download dependencies
go mod download

# 6. (Optional) Regenerate Swagger docs after changing any @Summary/@Router annotations
make swagger

# 7. Run the server
go run main.go
# or: make run

# 8. Confirm it's up
curl http://localhost:8082/health
```

The default port is `8082` (see `PORT` in [CONFIGURATION.md](./CONFIGURATION.md)).

## API documentation (Swagger UI)

Once running, the full per-endpoint API reference is live at:

- **http://localhost:8082/swagger/index.html** (also mounted at `/api/v1/swagger/index.html`)

This is generated from `@Summary`/`@Router`/... annotations on the handler
functions via `swag init` (see the `swagger` target in `Makefile`) — it is
the canonical source of truth for request/response shapes, required
permissions and status codes. Regenerate it with `make swagger` whenever you
add or change a route or its annotations.

## Tests

```bash
go test ./... -race -coverprofile=coverage.out -covermode=atomic
# or: make test
```

`make ci` runs the same checks CI runs: `tidy-check`, `build`, `vet`, `lint`, `test`.

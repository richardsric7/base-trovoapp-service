# Trovo Wallet API (`app-backend`)

`app-backend` is the core backend for the Trovo Wallet platform. It is a REST
API, written in Go, that owns user accounts, crypto wallets, deposits and
withdrawals, P2P (peer-to-peer) trading, tokenized real-world assets, KYC
verification, and white-label partner ("service link") integrations. If you
are new to this codebase, this file is the map — start here, then follow the
links at the bottom to the deeper docs.

## What this service does

- **User accounts & wallets** — registration, login (see
  [INTEGRATION.md](INTEGRATION.md) for how clients authenticate), wallet
  creation, sub-wallets, shared/approver access to a wallet, account
  recovery.
- **Payments & swaps** — sending assets between wallets, currency swaps,
  crypto deposits/withdrawals.
- **P2P marketplace** — merchants post buy/sell offers for assets against
  fiat, customers place orders, funds move through an on-chain escrow, and
  disputes can be raised and resolved.
- **Tokenized assets** — issuing, subscribing to, and exiting tokenized
  real-world assets (see `TOKENIZATION_PLAN.md` and
  `TOKENIZED_ASSET_PURCHASE_BY_FIAT.md` in this directory for the original
  design notes).
- **KYC** — identity verification via third-party providers (Sumsub, Doja).
- **Service links** — API-key-based integrations that let white-label
  partners (and Trovo's own admin backend) act on behalf of users, mint
  tokens, and read balances/history.

## Where this fits in the monorepo

This repository is one project inside a larger monorepo. The other three
projects that matter for understanding it:

| Project | Relationship to `app-backend` |
| --- | --- |
| **`tm-api`** | Trovo's internal admin/catalog backend. It **shares the same physical database** as `app-backend` and writes admin/catalog tables (curated assets, fee configuration, service link records) directly. For P2P transactional writes (creating orders, resolving disputes), it calls **`app-backend`'s own HTTP API** instead of writing to the DB directly — see [INTEGRATION.md](INTEGRATION.md) for why. |
| **`app-web`** | The end-user web wallet (Vite/React). A REST client of this API. |
| **`app-mobile`** | The end-user mobile wallet (Flutter). Also a REST client of this API. |

Both client apps authenticate with a request-signing scheme, not a bearer
token — see [INTEGRATION.md](INTEGRATION.md#authentication) for the details.

## Tech stack

- **Language**: Go (see `go.mod` — currently `go 1.23`, built with toolchain
  `go1.24.3`)
- **Web framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/), with schema managed by GORM
  `AutoMigrate` (see [DEPLOYMENT.md](DEPLOYMENT.md) for how/when this runs)
- **Database**: PostgreSQL in production; SQLite is supported for local
  development (`DB_TYPE=sqlite`)
- **Docs**: [swaggo/swag](https://github.com/swaggo/swag)-generated OpenAPI
  2.0 (Swagger) docs, served at `/swagger/index.html` once the app is
  running (see [Swagger UI](#swagger-ui) below)
- Also present: Redis (optional response caching), Firebase (push
  notifications + file storage), a Base/EVM blockchain client
  (`go-ethereum`), and a separate CockroachDB-compatible Postgres database
  ("RoachDB") used only by the payment-history tracking tables.

## Directory structure

```
app-backend/
├── main.go              # entrypoint: env checks, DB open + migrate, router wiring, background jobs
├── main_test.go          # integration test suite (needs a live server; excluded from `make test`)
├── docs/                 # swag-generated OpenAPI output (docs.go, swagger.json, swagger.yaml) — do not hand-edit
├── Dockerfile            # multi-stage build → small alpine runtime image
├── Makefile               # CI-parity targets: build, vet, lint, test
└── internal/
    ├── components/        # one folder per business domain — see below
    ├── db/                 # DB connection setup + GORM AutoMigrate model list
    ├── middleware/          # authentication (signature, API key, JWT), CORS
    ├── network/              # Base/EVM chain client & helpers
    ├── basetxn/               # transaction building/signing for the Base chain
    ├── evmkeypair/             # EVM key parsing/generation
    ├── cache/                   # Redis-backed response cache
    ├── mail/, sms/, pns/         # email (Mailgun), SMS (Termii/Infobip), push notifications (Firebase)
    ├── dynamiclinks/               # Firebase Dynamic Links (referral/deep links)
    ├── sharedconfig/                # the GlobalConfig struct threaded through every handler
    └── errors/, validators/          # shared error types and input validators
```

Each folder under `internal/components/` is a self-contained business
domain, and each follows the same internal shape:
`controllers/` (HTTP handlers + route registration), `models/` (GORM
models), `services/` (business logic), and sometimes `db/` (queries) or
`blockchain/`/`safesigner/` (on-chain interactions specific to that domain).

| Component | Owns |
| --- | --- |
| `root` | Health check, Apple/Android app-link well-known files |
| `users` | Accounts, wallets, KYC, shared access, tokenized assets, subwallets, account recovery — the largest component |
| `payments` | Sending assets between wallets |
| `swaps` | Currency/asset swaps |
| `assets` | Curated asset catalog (read side — the catalog itself is largely maintained by `tm-api`) |
| `rates` | Exchange rates |
| `p2p` | The P2P marketplace: offers, orders, escrow, disputes, refunds, merchant status |
| `servicelinks` | White-label partner API-key integrations, plus the `tm-api` server-to-server endpoints |
| `callbacks` | Webhook receivers from external providers (1Liquidity, Doja, Flutterwave) |
| `announcements` | In-app announcement banners and minimum app-version checks |

## Running it locally

You'll need Go installed (see `go.mod` for the version) and, for a full
setup, a Postgres instance — but the fastest path uses the bundled SQLite
mode, which needs nothing else installed. From the `app-backend/` directory:

```bash
# 1. Clone the monorepo and cd into this project
git clone <monorepo-url>
cd src-monorepo/app-backend

# 2. Copy the example environment file
cp .env.example .env

# 3. Fill in .env — see CONFIGURATION.md for what every variable does and
#    how to get a real value for it. The .env.example ships with
#    DB_TYPE=sqlite, so a local Postgres is optional to get started.

# 4. Download Go module dependencies
go mod download

# 5. Run the test suite to confirm everything compiles and passes
go test ./internal/... -vet=off

# 6. Start the server (defaults to 0.0.0.0:8080, or PORT if set)
go run main.go
```

That's 6 commands. The server logs each subsystem as it initializes
(`##users services initialized##`, etc.) and finally `##service started##`
once it's accepting requests.

### Swagger UI

Once the server is running, the full per-endpoint API reference is at:

```
http://localhost:8080/swagger/index.html
```

This is generated from `@Summary`/`@Router`/etc. comments above each
handler function (see the `docs/` folder) — it is always the most
up-to-date and precise source for what a given endpoint accepts and
returns. See [INTEGRATION.md](INTEGRATION.md) for how to authenticate
requests you try from the Swagger UI or elsewhere.

## Further reading

- **[DEPLOYMENT.md](DEPLOYMENT.md)** — how to build, containerize, and
  deploy this service, and what the migration/`DB_AUTOMIGRATE` behavior
  means for you.
- **[CONFIGURATION.md](CONFIGURATION.md)** — every environment variable
  this service reads, what it does, and how to get a real value for it.
- **[INTEGRATION.md](INTEGRATION.md)** — how `app-web`, `app-mobile`,
  `tm-api`, and white-label partners each integrate with this API.

# Trovo Wallet Monorepo — Architecture

This file is the map of the whole monorepo: what each project is, what it depends
on, and how a request actually flows between them. Read this first, then follow
the link at the end of each project's row into that project's own docs for the
detail (how to run it, every config value, every API endpoint).

If you are new here, the one sentence version is: **`app-backend` is the single
source of truth.** Everything else is either a client of it (a UI a person uses),
an admin surface that manages its data (`tm-api` / `tm-web`), or a supporting
library it or its clients depend on (`wallet-core`). `payment-history-engine` is
a separate, narrower service with its own database.

## The seven projects

| Project | What it is | Talks to | Docs |
|---|---|---|---|
| `app-backend` | Go REST API. The core Trovo Wallet backend: accounts, wallets, payments, swaps, P2P marketplace, tokenized assets, KYC, white-label service-link integrations. Owns the primary Postgres database. | Nothing else in this repo (it's the foundation) | [app-backend/README.md](app-backend/README.md) |
| `tm-api` | Go REST API. Backend for the internal admin dashboard (`tm-web`). | Shares `app-backend`'s database directly for admin/catalog tables; calls `app-backend`'s HTTP API for transactional writes. Has its own second database for admin-only data (admin accounts, roles, access logs). | [tm-api/README.md](tm-api/README.md) |
| `payment-history-engine` | Go service. Indexes/tracks payment history against its own separate database. | See its own [INTEGRATION.md](payment-history-engine/INTEGRATION.md) — verify there rather than assuming a relationship to the other services. | [payment-history-engine/README.md](payment-history-engine/README.md) |
| `app-web` | Vite + React + TypeScript. The end-user web wallet. | Calls `app-backend`'s REST API. May use `wallet-core` (verify in its own docs). | [app-web/README.md](app-web/README.md) |
| `app-mobile` | Flutter (Dart). The end-user mobile wallet (iOS + Android). | Calls `app-backend`'s REST API. Uses `wallet-core` as a native library via Dart FFI. | [app-mobile/README.md](app-mobile/README.md) |
| `tm-web` | Next.js + TypeScript. The internal admin dashboard staff use to manage the platform. | Calls `tm-api`'s REST API only — never talks to `app-backend` directly. | [tm-web/README.md](tm-web/README.md) |
| `wallet-core` | Rust crate. Shared cryptographic/wallet primitives (key derivation, signing) compiled two ways: to WebAssembly for web frontends, and to a native library (cdylib/staticlib) that `app-mobile` links via Dart FFI. | Consumed by `app-web`/`tm-web` (as wasm) and `app-mobile` (as a native lib) — verify actual current usage in its own docs, since a library can exist without every consumer having wired it up yet. | [wallet-core/README.md](wallet-core/README.md) |

## How a request actually flows

**A customer using the mobile or web app:**

```
app-mobile / app-web  →  app-backend REST API  →  app-backend's database
      ↑
      └── wallet-core (local key generation / transaction signing,
          never sent to any server)
```

**Staff using the admin dashboard:**

```
tm-web  →  tm-api REST API  →  ┬─ app-backend's database directly
                                │  (curated assets, fee configs, service
                                │  links, KYC/faucet configs — pure admin
                                │  config data, no app-backend business
                                │  logic to preserve)
                                │
                                └─ app-backend's REST API
                                   (P2P offers/orders/dispute resolution —
                                   anything with real business logic stays
                                   behind app-backend's own validation)
                                │
                                └─ tm-api's own database
                                   (admin accounts, roles/permissions,
                                   audit/access logs, suspension history —
                                   admin-only data app-backend has no
                                   concept of)
```

**Why tm-api has two different ways of touching data** — this is the single
most important thing to understand about this monorepo's design, and it is not
obvious from the code alone: for pure catalog/config tables that app-backend
merely *reads* but never enforces business rules on, tm-api writes to them
directly (same physical database, no network hop, no duplicate validation
logic to maintain in two places). For anything with real business logic —
creating a P2P offer, accepting an order, resolving a dispute — tm-api calls
app-backend's own HTTP API instead of writing the rows itself, so that logic
(fee calculation, state machine transitions, balance checks, audit trails)
only ever lives in one place. **When you add a new admin feature, this is the
first design decision to make**: does the table you're touching have any
app-backend business logic guarding it? If yes, add or call an app-backend
endpoint. If no (it's pure config/catalog data), tm-api can write to it
directly the same way it already does for curated assets and service links —
see `tm-api/internal/components/usermetrics/db/curated_assets.go` for the
established pattern to copy (explicit `gorm:"column:..."` tags and a
`TableName()` override on the mirrored struct, map-based `Updates()` instead
of struct-based `Save()`/`Updates()` to avoid GORM silently dropping a field
being set back to a zero value, and a dedicated single-field toggle endpoint
for anything with an active/inactive-style flag).

## A known gap: the CI workflows don't match this repo's layout

`.github/workflows/deploy.yml` and `.github/workflows/pr-checks.yml` reference
directories named `backend/`, `web/`, `mobile/`, `trovotech-io/`, and
`trovo-app-website/`. **None of those directories exist in this repo** — the
actual project directories are `app-backend/`, `app-web/`, `app-mobile/`,
`tm-api/`, `tm-web/`, `payment-history-engine/`, `wallet-core/`. These
workflow files appear to have been inherited from a differently-laid-out
sibling repository and, as written, their path filters will never match a
change made in this repo, so the jobs they gate (build/push images, trigger
Portainer redeploys, PR build/lint/test checks) will not actually run here.
Each project's own `DEPLOYMENT.md` describes what can be verified from that
project's own `Dockerfile`/`Makefile`/CI config instead of trusting these
workflow files. If/when this repo's actual CI is wired up (either by fixing
these paths or replacing them), update this note and each project's
`DEPLOYMENT.md` accordingly.

## The documentation standard going forward

See [CONTRIBUTING.md](CONTRIBUTING.md) for what's required whenever a new
project is added to this monorepo, or an existing one is modified — including
which docs must be kept in sync and how to regenerate Swagger.

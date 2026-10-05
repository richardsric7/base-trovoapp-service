# Trovo Wallet Monorepo — Architecture

This file is the map of the whole monorepo: what each project is, what it depends
on, and how a request actually flows between them. Read this first, then follow
the link at the end of each project's row into that project's own docs for the
detail (how to run it, every config value, every API endpoint).

If you are new here, the one sentence version is: **`app-backend` is the single
source of truth.** Everything else is either a client of it (a UI a person uses),
an admin surface that manages its data (`tm-api` / `tm-web`), or a supporting
library it or its clients depend on (`wallet-core`). `payment-history-engine` is
a separate, narrower service with its own database; `payout-engine` is a worker
on app-backend's database that pays tokenized-asset proceeds to token holders.

## The projects

| Project | What it is | Talks to | Docs |
|---|---|---|---|
| `app-backend` | Go REST API. The core Trovo Wallet backend: accounts, wallets, payments, swaps, P2P marketplace, tokenized assets, Public Markets (tokenized NGX stocks and FMDQ bonds backed by CSCS custody, with exchange-partner endpoints under `/v1/trovo-api/` and its own dividend engine, separate from `payout-engine`; see [app-backend/PUBLIC_MARKETS.md](app-backend/PUBLIC_MARKETS.md)), KYC, white-label service-link integrations. Owns the primary Postgres database. | Nothing else in this repo (it's the foundation) | [app-backend/README.md](app-backend/README.md) |
| `tm-api` | Go REST API. Backend for the internal admin dashboard (`tm-web`). | Shares `app-backend`'s database directly for admin/catalog tables; calls `app-backend`'s HTTP API for transactional writes. Has its own second database for admin-only data (admin accounts, roles, access logs). | [tm-api/README.md](tm-api/README.md) |
| `payment-history-engine` | Go service. Indexes/tracks payment history against its own separate database. | See its own [INTEGRATION.md](payment-history-engine/INTEGRATION.md) — verify there rather than assuming a relationship to the other services. | [payment-history-engine/README.md](payment-history-engine/README.md) |
| `payout-engine` | Go worker, no HTTP. Pays a tokenized asset's proceeds (dividends, interest) to its token holders from a dedicated payout Safe: snapshots holders by replaying the token's transfers, locks a schedule (with the payout's processing fee and VAT) for admin approval, checks the Safe is funded, pays in batches with crash-safe recovery and notifies paid users. One active instance (heartbeat claim). | Works from `app-backend`'s database (payout tables), driven by `tm-api` admin actions, which wake it over Redis; reads and writes Base; signers from the vault manager's managed secret. | [payout-engine/README.md](payout-engine/README.md) |
| `app-web` | Vite + React + TypeScript. The end-user web wallet. | Calls `app-backend`'s REST API. May use `wallet-core` (verify in its own docs). | [app-web/README.md](app-web/README.md) |
| `app-mobile` | Flutter (Dart). The end-user mobile wallet (iOS + Android). | Calls `app-backend`'s REST API. Uses `wallet-core` as a native library via Dart FFI. | [app-mobile/README.md](app-mobile/README.md) |
| `tm-web` | Next.js + TypeScript. The internal admin dashboard staff use to manage the platform. | Calls `tm-api`'s REST API only — never talks to `app-backend` directly. | [tm-web/README.md](tm-web/README.md) |
| `wallet-core` | Rust crate. Shared cryptographic/wallet primitives (key derivation, signing) compiled two ways: to WebAssembly for web frontends, and to a native library (cdylib/staticlib) that `app-mobile` links via Dart FFI. | Consumed by `app-web`/`tm-web` (as wasm) and `app-mobile` (as a native lib) — verify actual current usage in its own docs, since a library can exist without every consumer having wired it up yet. | [wallet-core/README.md](wallet-core/README.md) |
| `paymaster` | Two parts: `contracts/` — `TrovoTokenPaymaster`, an ERC-4337 (EntryPoint v0.7) paymaster that lets wallets pay gas in curated stablecoins; `quote-service/` — Go service that discovers exchange rates from pluggable sources (Chainlink, DEX pool TWAP for cNGN, JSON APIs), adds Trovo's spread and signs per-operation quotes. Configured from Vault. | `app-backend` calls the quote service when building a stablecoin-paid UserOperation; the quote service reads Base (paymaster, EntryPoint deposit, feeds, pools) and external rate APIs; the bundler submits operations to the EntryPoint, which calls the paymaster. | [paymaster/README.md](paymaster/README.md) |
| `market` | `contracts/` — `TrovoOfferBook`, fixed-price offers between curated tokens where every fill needs a platform-signed authorization (how tokenized assets are sold), and `TokenizedAsset`, the per-asset token owned by its issuing Safe. Hardhat project with tests and a Vault-configured deploy script. | `app-backend` builds the mint (issuing Safe opens the sale offer), signs fill authorizations for checked purchases and builds buyer operations; the fiat path has the internal balance minting Safe fill for the buyer. | [market/README.md](market/README.md) |
| `recovery` | `contracts/` — Candide's Social Recovery Module v0.2.0, vendored unmodified and deployed by Trovo, for opt-in account recovery: the platform's recovery guardian can replace a covered wallet's key after a waiting period the user can cancel, never move funds. Hardhat project with tests against real Safe wallets and a Vault-configured deploy script. | `app-backend` builds the users' operations that turn recovery on/off or cancel a recovery, has the guardian start and finalize recoveries, and watches the module for recoveries it did not start. | [recovery/README.md](recovery/README.md) |

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

## CI/CD

There is **no CI/CD pipeline** in this repository. The GitHub Actions
workflows and iOS fastlane setup that came in with the copied projects
targeted other repositories' layouts and infrastructure, so they were
removed. Each project's `DEPLOYMENT.md` describes how to build, test and
run it by hand (Dockerfiles, `Makefile` targets such as `make ci`). When a
pipeline is set up for this repo, document it here and in each project's
`DEPLOYMENT.md`.

## Running more than one instance of a service

`app-backend` and `tm-api` are both safe to scale to multiple instances
behind a load balancer. This wasn't automatic — a few things any single
in-process instance could get away with (an in-memory map/channel for
coordination, a background loop with no idea another copy of itself might
be running elsewhere) needed a shared, cross-instance mechanism instead.
The pattern used everywhere this came up:

- **Cross-instance mutual exclusion** (who's allowed to do X right now)
  is a plain database table claimed with a single atomic `UPDATE ...
  WHERE` statement — never a Postgres advisory lock or `SELECT ... FOR
  UPDATE`, since neither has an equivalent on SQLite, which every one of
  these Go projects also needs to run correctly against for local
  dev/tests. See `app-backend/internal/sharedconfig/distributed_lock.go`
  (the schema-migration lock and the singleton-locked background loops)
  and `app-backend/internal/sharedconfig/channel_accounts.go` (the
  channel-account pool, same idea applied to claiming one of several
  interchangeable rows instead of one named lock).
- **Cross-instance real-time delivery** (an event that happened on
  instance A needs to reach a websocket/SSE connection open on instance
  B) is Redis Pub/Sub, reusing whatever Redis connection the service
  already has for caching — fire-and-forget by design (`PUBLISH` is
  at-most-once), which is an acceptable tradeoff everywhere it's used
  because there's always a durable fallback (REST polling, a push
  notification, the database itself) that isn't affected if the live
  nudge is missed. See `app-backend/internal/sharedconfig/realtime.go`
  (P2P order/offer/escrow/dispute updates over the user websocket) and
  `tm-api/internal/models/streams.go` (the admin QR-login notification
  stream).
- **Cross-instance rate limiting** (a per-caller request budget that must
  hold across every instance, not per-process) is a Redis fixed-window
  counter (`INCR`+`EXPIRE`), reusing the same Redis connection as caching
  and pub/sub above — an in-memory counter would silently multiply the
  effective limit by however many replicas are running, which defeats the
  point of having a limit at all. Fails open (no-ops) when Redis is
  disabled/unreachable, same graceful-degradation posture as caching and
  pub/sub. See `app-backend/internal/middleware/rate_limit_middleware.go`
  and `tm-api/internal/middleware/rate_limit_middleware.go`.
- **Platform signing keys** (a platform Safe, or a key that sends
  transactions) are used by one instance at a time: a named lock held from
  reading the nonce until the transaction is mined (a Safe) or broadcast (a
  key); a transaction signed now and sent later reserves its key's nonce
  under the same lock instead (`sharedconfig.ReserveNonce`). app-backend uses the same `distributed_locks` table, renewed while
  held (`sharedconfig.WithKeyLock`, wired into `internal/gnosissafe`);
  tm-api, which has no lock table, uses a Redis lock (`cache.WithLock`) for
  its vault-signer Safe changes, so running tm-api on more than one
  instance needs Redis (`ENABLE_CACHING=1`).
- **Background jobs** hold their lock for as long as they run (it is
  renewed), and anything that must survive an instance stopping is in the
  database rather than memory (partner callback deliveries, sale-start
  notifications). app-backend shuts down gracefully on SIGTERM; tm-api and
  the paymaster quote service already did.

The paymaster quote service is stateless and scales freely (each instance
does send its own "deposit low" alert). `payment-history-engine` has not been reviewed for
running on more than one instance.

**Going forward:** if you're adding a new project to this monorepo, or a
new piece of shared, cross-instance coordination to an existing one,
follow the same rules — a portable claim table instead of a
database-specific locking primitive, Redis Pub/Sub (with a
non-Redis-dependent fallback) instead of an in-process map/channel, and a
Redis counter instead of an in-process one for anything resembling a rate
limit — so the new code works the same whether it's one instance or ten,
and works the same on SQLite as it does on Postgres.

## The documentation standard going forward

See [CONTRIBUTING.md](CONTRIBUTING.md) for what's required whenever a new
project is added to this monorepo, or an existing one is modified — including
which docs must be kept in sync and how to regenerate Swagger.

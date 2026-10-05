# Integration

Everything `payment-history-engine` connects to, and how. The engine makes
**no calls to app-backend or tm-api, and they make none to it**: they are
connected only through two shared databases. It is closer to a data
pipeline that shares tables with them than to an API.

| Connects to | Direction | Through | Needed? |
|---|---|---|---|
| [app-backend's main database](#1-app-backends-main-database) | engine reads wallets, writes payment history | Postgres | yes |
| [The tracking database ("RoachDB")](#2-the-tracking-database-roachdb) | engine and app-backend both write tracked wallets; engine keeps its scan position | CockroachDB / Postgres | yes |
| [app-backend and tm-api (readers)](#3-app-backend-and-tm-api-the-readers) | they read what the engine writes | the main database | — |
| [Base blockchain](#4-base-blockchain) | engine reads | JSON-RPC over HTTPS | yes |
| [Redis](#5-redis) | engine reads and writes | Redis | optional |
| [Discord](#6-discord-alerts) | engine posts alerts | webhooks | optional |
| [Your hosting platform](#7-health-checks) | platform checks the engine | HTTP `/health`, `/ready` | recommended |

---

## 1. app-backend's main database

- **What it is and why:** the same Postgres database app-backend uses for
  accounts and wallets. The engine finds the wallets to watch there and
  writes the payment history users see.
- **Direction:** the engine reads and writes.
  - **`user_wallets`** (owned by app-backend): the engine reads it to find
    wallets not yet tracked, registers them for tracking, and sets their
    `Tracked` flag.
  - **`payment_history`** (created by app-backend; written by the engine):
    every transfer the engine finds touching a tracked wallet is written
    here by `SavePaymentHistory`
    (`internal/components/payments/services/procedure.go`). The engine does
    not create this table: app-backend must have started against the
    database first.
- **How they connect:** a direct database connection.
- **Settings on this side:** [`DB_CONNECTION_STRING`](CONFIGURATION.md#db_connection_string)
  (example `postgres://trovo:change-me@db.internal:5432/trovo?sslmode=require`),
  [`DB_TYPE`](CONFIGURATION.md#db_type).
- **Settings on the other side:** the same value as app-backend's
  `DB_CONNECTION_STRING` ([app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)).
- **How to check it works:** after a test transfer to a Trovo wallet, a new
  row appears in `payment_history`.
- **When it is down:** the engine cannot record transfers; it resumes from
  its saved block when the database is back.

**Each row** has a source side (`SourceNetwork`, `SourceAssetCode`,
`SourceContractAddress`, `SourceAmount`: what left `FromAddress`) and a
destination side (the same `Destination*` fields: what arrived at
`ToAddress`). The engine sees one leg of each transfer, so it writes the
same values on both sides; a row where they differ would be a swap, which
the engine does not produce today. `SourceNetwork` is `base`.

**If you change `PaymentHistory`** (`internal/components/payments/models/payment_history.go`),
change the matching structs in app-backend and tm-api too: three projects
map Go structs onto the same table, with nothing enforcing they agree.

## 2. The tracking database ("RoachDB")

- **What it is and why:** a separate CockroachDB (or Postgres) database that
  holds the list of wallets and addresses to watch, and how far the engine
  has scanned.
- **Direction:**
  - **`tracked_wallets`, `tracked_addresses`:** written by app-backend when a
    user creates or gets approved for a wallet (`gc.RoachDB.Create(...)` in
    app-backend's `registration.go`, `subwallets.go`, `approvals.go`) and by
    the engine itself; read by the engine to know what to watch. Both
    projects create these tables' structure (`AutoMigrate`).
  - **`monitored_cursors`, `monitored_account_cursors`:** the engine's own
    record of the last block scanned, overall and per address. Nothing else
    touches them.
- **How they connect:** a direct database connection (Postgres protocol).
- **Settings on this side:** [`CDB_CONNECTION_STRING`](CONFIGURATION.md#cdb_connection_string)
  (example `postgresql://root@roach.internal:26257/payment_history?sslmode=require`),
  [`ROACH_DB_TYPE`](CONFIGURATION.md#roach_db_type).
- **Settings on the other side:** app-backend's `CDB_CONNECTION_STRING` must
  point to the same database.
- **How to check it works:** `/ready` lists the tracking database as `up`,
  and `tracked_wallets` gains a row when a test user signs up.
- **When it is down:** the engine cannot start or save its position.

## 3. app-backend and tm-api (the readers)

- **What it is and why:** they show the history the engine writes.
  - app-backend: `GET /v1/trovo-api/users/payment-history/:walletAddress`
    (and the apps' history screens), reading `payment_history` through its
    own database connection.
  - tm-api: `GET /payment/history` for Trovo Manager, reading the same table
    through its connection to that database.
- **Direction:** they read the main database; nobody calls the engine.
- **Settings:** nothing to match beyond using the same database.
- **When the engine is down:** their history endpoints keep working and
  show what has been recorded so far; new transfers appear once the engine
  catches up.

For those endpoints' details see app-backend's and tm-api's Swagger docs
(`app-backend/docs/`, `tm-api/docs/`).

## 4. Base blockchain

- **What it is and why:** the network where transfers happen.
- **Direction:** the engine reads only.
- **What it does:**
  1. **Fills in history** for each tracked address: replays its past token
     `Transfer` logs (`eth_getLogs`, 5,000 blocks at a time) and, for Trovo
     wallets (Safes), the ETH they received (`SafeReceived` events) and sent
     (their user operations, found by the EntryPoint's
     `UserOperationEvent`), resuming from the address's saved cursor.
  2. **Follows new blocks** one at a time. For each block it records ETH
     transfers (a transaction's own value, ETH Trovo wallets send inside
     their user operations, decoded from the EntryPoint's `handleOps`, and
     ETH a Safe receives from contracts) and token `Transfer` logs (marking
     transfers from or to the zero address as `MINT TOKEN` / `BURN TOKEN`).
     Only transfers involving a tracked wallet are kept.
- **How they connect:** JSON-RPC over HTTPS, with a per-request timeout and
  a short pause after repeated failures.
- **Settings on this side:** [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url),
  [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id),
  [`ENTRYPOINT_ADDRESS`](CONFIGURATION.md#entrypoint_address) (must equal
  app-backend's), [`RPC_TIMEOUT`](CONFIGURATION.md#rpc_timeout),
  [`TRACK_ADDRESS_CONCURRENCY`](CONFIGURATION.md#track_address_concurrency).
- **How to check it works:** `/ready` lists the node as `up`, and its scan
  progress moves.
- **When it is down:** scanning pauses; nothing is skipped, and it catches
  up when the node is back.

## 5. Redis

- **What it is and why:** an optional cache.
- **Settings:** [`ENABLE_CACHING`](CONFIGURATION.md#enable_caching),
  [`REDIS_HOST`](CONFIGURATION.md#redis_host),
  [`REDIS_PORT`](CONFIGURATION.md#redis_port),
  [`REDIS_PASSWORD`](CONFIGURATION.md#redis_password),
  [`CACHING_PARAMETER`](CONFIGURATION.md#caching_parameter).
- **When it is down:** with caching on, the engine refuses to start if Redis
  does not answer at start.

## 6. Discord alerts

- **What it is and why:** chat alerts for a full database connection pool,
  node errors and failed payment records.
- **Settings:** [`CONNECTION_WARNING_WEBHOOK`](CONFIGURATION.md#connection_warning_webhook),
  [`EXPANSION_NETWORK_ERROR_WEBHOOK`](CONFIGURATION.md#expansion_network_error_webhook),
  [`FAILED_PAYMENT_ERROR_WEBHOOK`](CONFIGURATION.md#failed_payment_error_webhook).
  Each falls back to a built-in webhook when unset.

## 7. Health checks

- **What it is and why:** the engine's only HTTP server, for your hosting
  platform.
  - `GET /health`: answers while the process runs (liveness).
  - `GET /ready`: checks each dependency and reports the scan's progress
    (readiness). With no wallets tracked yet the scan reports `idle`, which
    is not an error.
  - `GET /swagger/index.html`: the reference for these two.
- **Settings:** [`HEALTH_PORT`](CONFIGURATION.md#health_port) (default
  `8080`).

There is no message queue or event bus between the projects: the engine's
polling loops are the only thing moving data from "app-backend added a
wallet" to "the engine records its transfers".

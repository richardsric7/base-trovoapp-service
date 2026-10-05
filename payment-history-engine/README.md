# Payment History Engine

## What this project does

When money moves in or out of a Trovo wallet, users expect to see it in
their transaction history. Transfers happen on the Base blockchain, often
without going through Trovo at all (someone sending from an exchange, for
example). `payment-history-engine` watches the blockchain for every
transfer to or from a Trovo user's wallet and records it in the payment
history that the Trovo apps and Trovo Manager show.

It runs in the background with no user-facing API: it reads the list of
wallets from app-backend's database and writes the history back there.
Without it, transfers that did not go through Trovo would never appear in
users' history.

- How to deploy: [DEPLOYMENT.md](DEPLOYMENT.md)
- Every setting: [CONFIGURATION.md](CONFIGURATION.md)
- What it connects to: [INTEGRATION.md](INTEGRATION.md)

## What it actually does

On startup, `main.go` launches four long-running loops (goroutines) and never
returns:

1. **Track new user wallets** — polls the primary wallet database for
   `user_wallets` rows that aren't yet tracked, and registers each one in
   RoachDB (`tracked_wallets`, and `tracked_addresses` when address-level
   tracking is on).
2. **Backfill tracked addresses** — for every address in RoachDB's
   `tracked_addresses`, replays its historical Base logs (`eth_getLogs`,
   chunked in ranges of 5,000 blocks) to catch up on B20/ERC-20 `Transfer`
   events involving it and, for Trovo wallets (Safes), their ETH: what the
   Safe received (its `SafeReceived` events) and sent (its user operations,
   found by the EntryPoint's `UserOperationEvent`), resuming from a saved
   per-address cursor.
3. **Monitor the payment stream** — walks new Base blocks one at a time
   (there is no account-agnostic push feed on Base, unlike Stellar Horizon's
   `StreamPayments`), and for each block:
   - records ETH transfers (`nativeTransfersOfTx` in `userops.go`): a
     transaction's own value, the ETH Trovo wallets send inside their
     ERC-4337 user operations (decoded from the EntryPoint's `handleOps`,
     batches included, only for operations that succeeded), and ETH a Safe
     receives from other contracts (`SafeReceived`), and
   - decodes B20/ERC-20 `Transfer` event logs (`processB20TransferLog`),
     classifying `MINT TOKEN` / `BURN TOKEN` transfers to/from the zero
     address.
   Only transfers where the sender or receiver is a tracked wallet are kept;
   everything else is skipped. Each processed block's number is saved as a
   resume cursor.

Every recorded transfer is written as a `PaymentHistory` row via
`SavePaymentHistory` (`internal/components/payments/services/procedure.go`).

## Role in the monorepo

This engine does **not** call `app-backend` or `tm-api` over HTTP, and they
don't call it either — there is no outbound HTTP client to another Trovo
service anywhere in this codebase. Instead, all three services are coupled
through two **shared databases**:

| Database | Env var | Shared table(s) | Who writes | Who reads |
|---|---|---|---|---|
| Primary "wallet" DB (Postgres/SQLite) | `DB_CONNECTION_STRING` | `user_wallets` | app-backend | **this engine** (discovers wallets to track) |
| Primary "wallet" DB | `DB_CONNECTION_STRING` | `payment_history` | **this engine** (writes every transfer it detects) | app-backend (`GET /v1/trovo-api/users/payment-history/:walletAddress`), tm-api (`GET /payment/history`) |
| RoachDB (CockroachDB) | `CDB_CONNECTION_STRING` | `tracked_wallets`, `tracked_addresses` | app-backend (when a user registers/approves a wallet) | **this engine** (what to monitor) |
| RoachDB | `CDB_CONNECTION_STRING` | `monitored_cursors`, `monitored_account_cursors` | **this engine** only | **this engine** only (its own resume state) |

Notably, this engine writes `PaymentHistory` rows into the **primary wallet
DB**, not RoachDB — RoachDB here is used purely for this engine's own
tracking/coordination bookkeeping ("its own database", separate from the
main app's business data). The `payment_history` table's schema is owned and
migrated by app-backend (`app-backend/internal/db/main.go`); this engine only
`AutoMigrate`s the RoachDB tables it owns
(`TrackedWallet`, `TrackedAddress`, `MonitoredCursor`, `MonitoredAccountCursor`).

See [INTEGRATION.md](./INTEGRATION.md) for the full picture, including what
was checked and ruled out.

## Tech stack

- **Language**: Go 1.26 (the version in `go.mod`)
- **HTTP**: plain `net/http` + `http.ServeMux` — **not** Gin, despite Gin
  being a dependency (`github.com/gin-gonic/gin` is pulled in only because
  `internal/errors`, `internal/components/payments/errors` and
  `internal/middleware` were copied from `app-backend` and use `gin.H`/
  `gin.Context`/`gin.HandlerFunc` in their signatures; none of that code is
  wired into a router anywhere in this repo — see [INTEGRATION.md](./INTEGRATION.md))
- **ORM**: [GORM](https://gorm.io) (`gorm.io/gorm`), with
  `gorm.io/driver/postgres` and `gorm.io/driver/sqlite`
- **Databases**: Postgres (primary "wallet" DB) and CockroachDB, a
  Postgres-wire-compatible distributed SQL database, referred to in this repo
  as "RoachDB" (both use the same `gorm.io/driver/postgres` driver — CRDB's
  Postgres-compatibility is what makes that work); SQLite is supported as a
  local/CI escape hatch for both (`DB_TYPE=sqlite`, `ROACH_DB_TYPE=sqlite`)
- **Blockchain client**: [go-ethereum](https://github.com/ethereum/go-ethereum)
  (`ethclient`) against a Base (Ethereum L2) JSON-RPC endpoint
- **Cache**: Redis (`github.com/go-redis/redis/v8`), optional
- **API docs**: [swaggo/swag](https://github.com/swaggo/swag) +
  [swaggo/http-swagger](https://github.com/swaggo/http-swagger) (see below)

## Directory structure

```
payment-history-engine/
├── main.go                     Entry point: DB setup, the 4 worker loops, block/log processing
├── main_test.go                 Manual scripts against a live server (build tag livetests; not run by `go test`)
├── userops.go                   ETH sent/received by Safe wallets: user operation and SafeReceived decoding
├── docs/                         Generated Swagger/OpenAPI spec (swag init output — do not hand-edit)
├── internal/
│   ├── basetxn/                  Base transaction building/signing helpers
│   ├── cache/                    Redis cache wrapper
│   ├── components/
│   │   ├── health/                Liveness/readiness probes + the HTTP server that exposes them
│   │   ├── payments/
│   │   │   ├── db/                 IP-geolocation helper (vestigial, see CONFIGURATION.md)
│   │   │   ├── errors/             gin.H-shaped error types (vestigial — see above)
│   │   │   ├── models/             PaymentHistory, TrackedWallet, TrackedAddress, cursors
│   │   │   └── services/           TrackUserWallet, SavePaymentHistory
│   │   └── users/                 UserWallet / user models read from the primary DB
│   ├── db/                       DB connection setup (primary DB + RoachDB)
│   ├── errors/                   More gin.H-shaped error types (vestigial)
│   ├── evmkeypair/               EVM keypair parsing (address/private-key handling)
│   ├── middleware/                Gin auth/CORS middleware (vestigial — unused, no router mounts it)
│   ├── network/                  Base RPC client, chain ID, mnemonic-derived accounts
│   └── validators/                Input validators (email, address, asset-code format)
├── Dockerfile
├── .env.example
└── LICENSE
```

## Running it locally

You need Go 1.26+ and access to a Postgres-compatible primary DB and a
CockroachDB (or Postgres, via the driver) RoachDB — or use the SQLite escape
hatches to avoid standing up either. See
[CONFIGURATION.md](./CONFIGURATION.md) for what every variable means and how
to get a real value, and [DEPLOYMENT.md](./DEPLOYMENT.md) for Docker/build
details.

```bash
git clone git@github.com:trovotech-technologies/trovo-wallet-payment-history-engine.git payment-history-engine
cd payment-history-engine
cp .env.example .env
```

Then edit `.env` and, at minimum, add the SQLite escape hatches so you don't
need real databases to boot:

```bash
echo 'DB_TYPE=sqlite'                          >> .env
echo 'DB_CONNECTION_STRING=./local-wallet.db'  >> .env
echo 'ROACH_DB_TYPE=sqlite'                    >> .env
echo 'CDB_CONNECTION_STRING=./local-roach.db'  >> .env
```

Fill in the other required variables (see
[CONFIGURATION.md](./CONFIGURATION.md) for the full list and example values),
then:

```bash
go mod download
go build ./...
go run .
```

The engine starts its worker loops (they'll mostly idle, since a fresh
SQLite DB has no `user_wallets`/`tracked_addresses` rows) and its monitoring
HTTP server on `:8080` (override with `HEALTH_PORT`).

## Swagger / API docs

Once running, open:

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Raw OpenAPI spec**: http://localhost:8080/swagger/doc.json

This only documents the monitoring endpoints (`GET /health`, `GET /ready`) —
this service has no other API surface. See
[INTEGRATION.md](./INTEGRATION.md) for how its actual output
(`payment_history` rows) is consumed by other services.

To regenerate the docs after changing a handler's `@`-annotations:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

## Further reading

- [DEPLOYMENT.md](./DEPLOYMENT.md) — prerequisites, building, running via Docker (there is no CI/CD pipeline)
- [CONFIGURATION.md](./CONFIGURATION.md) — every environment variable, explained
- [INTEGRATION.md](./INTEGRATION.md) — how this service fits with app-backend and tm-api

## Error message shape (legacy, mostly unused here)

The `internal/errors` and `internal/components/payments/errors` packages
(copied from `app-backend`) define errors shaped like this for a Gin JSON
response, but as noted above, no route in this repo actually returns one —
only two internal auth-middleware functions call `.JSONError()`, and neither
is wired into a router in this codebase:

| name | type | description |
|---|---|---|
| error | string | Type of error |
| message | string | Human readable description of the error |
| data | string | Extra developer data related to the error |

## License

MIT — see [LICENSE](./LICENSE).

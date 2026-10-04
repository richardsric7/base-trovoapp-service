# Integration with the rest of the monorepo

## Summary

`payment-history-engine` makes **no outbound HTTP calls to `app-backend` or
`tm-api`**, and neither of them calls it over HTTP either — grepped this
repo and both of theirs for `http.NewRequest`, `http.Get`, `sling.New` /
`dghubble/sling` usage, and any hardcoded base URL pointing at another Trovo
service, and found none. All coupling between the three is **indirect**,
through two shared SQL databases. This is worth being explicit about because
it's easy to assume a "payment history engine" calls or is called by the
main wallet API — it doesn't; it's closer to a fully decoupled ETL/indexer
that happens to share tables with two other services.

## What was checked

- `grep -rn "http.NewRequest\|http.Get\|http.Client{" payment-history-engine/` —
  only hits are `internal/components/payments/db/geo.go` (an IP-geolocation
  lookup, unrelated to other Trovo services) and
  `internal/components/health/health.go`'s own outbound ping to
  `BASE_RPC_URL` (the blockchain node, not another Trovo service).
- `grep -rln "payment-history" app-backend/ tm-api/` — the only substantive
  hits are:
  - `app-backend/main.go`: a comment noting the CI migration step
    deliberately doesn't configure `CDB_CONNECTION_STRING` because "RoachDB
    belongs to the payment-history service".
  - `app-backend/internal/components/servicelinks/controllers/handlers_impl.go`
    and `main.go`: the `GET /v1/trovo-api/users/payment-history/:walletAddress`
    endpoint, which reads the `payment_history` table this engine writes
    into (details below) — **via its own database connection, not by calling
    this engine**.
- `grep -rln "gin.New(\|gin.Default(\|gin.Engine" payment-history-engine/` —
  no hits. This confirms this service, despite importing
  `github.com/gin-gonic/gin` transitively (for vestigial error/middleware
  types copied from `app-backend`), never actually starts a Gin server. Its
  only HTTP server is the plain `net/http` one in
  `internal/components/health/server.go`.

## The actual coupling: two shared databases

### 1. Primary "wallet" database (`DB_CONNECTION_STRING`)

This is the same Postgres database `app-backend` uses for user accounts and
wallets.

- **`user_wallets`** — owned and written by `app-backend`. This engine only
  reads it (`main.go`'s "track new user wallets" loop) to discover wallets
  that need to be registered for monitoring, and writes back a `Tracked`
  flag once it has processed one.
- **`payment_history`** — **written by this engine** (`SavePaymentHistory`
  in `internal/components/payments/services/procedure.go`), every time it
  detects a native transfer or B20/ERC-20 `Transfer` event touching a
  tracked wallet. Its schema is migrated by `app-backend`
  (`app-backend/internal/db/main.go` calls
  `gormDB.AutoMigrate(&paymentModels.PaymentHistory{})`) — this engine's own
  `main.go` does **not** migrate this table, it assumes app-backend already
  created it. Read by:
  - `app-backend`'s `GET /v1/trovo-api/users/payment-history/:walletAddress`
    (`app-backend/internal/components/payments/services/history.go`'s
    `GetPaymentHistory`, queried against `gc.DB`, i.e. the same primary DB).
  - `tm-api`'s (the admin dashboard) `GET /payment/history`
    (`tm-api/internal/components/general/services/payment_history.go`,
    queried against its own `walletDB` handle — same physical database).
- **`market_offers`** — owned/migrated by `app-backend`
  (`gormDB.AutoMigrate(&users.MarketOffer{})`). This engine's trade-stream
  loop (`MonitorTradeStream`) only reads it, to decide whether there's any
  active market-making offer worth watching — and today, even when one
  exists, the loop is an idle stub (Base has no DEX/order-book event source
  wired up yet).

### 2. RoachDB / CockroachDB (`CDB_CONNECTION_STRING`)

A separate database, referred to throughout both repos as "RoachDB" — this
engine's own coordination store, kept apart from the main app's business
data.

- **`tracked_wallets`**, **`tracked_addresses`** — written by `app-backend`
  when a user's wallet is created/approved for monitoring
  (`app-backend/internal/components/users/services/subwallets.go`,
  `approvals.go`, `registration.go` all call `gc.RoachDB.Create(...)` on
  these), and also written by **this engine itself** in its "track new user
  wallets" loop. Both services migrate these two tables' schema
  independently (`app-backend/main.go` and this repo's `main.go` both call
  `AutoMigrate` on them) — a bit of duplicated ownership worth knowing about
  if the schema ever needs to change. Read continuously by this engine to
  know what to poll.
- **`monitored_cursors`**, **`monitored_account_cursors`** — written and
  read **only by this engine**. Pure internal bookkeeping: how far into the
  Base chain the global stream and each per-address backfill have gotten.
  Nothing else in the monorepo touches these two tables.

## What this means in practice

- If you change the shape of `PaymentHistory` in
  `internal/components/payments/models/payment_history.go`, you must also
  update the matching struct in `app-backend` (and check `tm-api`'s reads),
  since three codebases independently define Go structs mapped onto the same
  table — there is no shared Go module or schema registry enforcing this.
- `PaymentHistory` splits each row into a **source** side (`SourceNetwork`,
  `SourceAssetCode`, `SourceContractAddress`, `SourceAmount` — what left
  `FromAddress`) and a **destination** side (the equivalent `Destination*`
  fields — what arrived at `ToAddress`). `SavePaymentHistory`
  (`internal/components/payments/services/procedure.go`) takes a `network`
  parameter and, since every current caller in `main.go` only ever observes
  one leg of a transfer, writes it to both sides identically — a plain
  payment genuinely has the same asset/network on both ends. `NetworkBase =
  "base"` is the only network this engine watches today; it's a real column
  (not left implicit) so a future second network or a Base-native bridge
  is a new value here, not another schema change. A row where source and
  destination differ would be a swap — this engine doesn't produce one
  today (Base has no live DEX/AMM integration; `swapTransactionType()` in
  `main.go` has zero callers), but the schema is ready for whoever builds
  that correlation later.
- If this engine is down, `app-backend`'s and `tm-api`'s payment-history
  endpoints keep serving whatever rows already exist — they degrade to
  "stale data", not "erroring out", because they don't depend on this
  engine being reachable, only on the database being reachable.
- If `app-backend` stops writing to `tracked_wallets`/`tracked_addresses`
  (e.g. it's down), this engine simply has nothing new to track — its
  `GET /ready` probe reports `stream.status: "idle"` in that case (see
  `internal/components/health/health.go`'s `hasWork` handling), not an
  error.
- There is no message queue, event bus, or webhook connecting these
  services either — the polling loops in `main.go` (`time.Sleep` between
  passes) are the only mechanism moving data from "app-backend wrote a row"
  to "this engine notices it".

## Canonical per-endpoint reference

This engine's own HTTP surface (`GET /health`, `GET /ready`) is documented
via Swagger — see [README.md](./README.md#swagger--api-docs). Once running:

- Swagger UI: http://localhost:8080/swagger/index.html
- Raw spec: http://localhost:8080/swagger/doc.json

For the endpoints that actually serve this engine's *output* (payment
history) to end users, see `app-backend`'s and `tm-api`'s own Swagger docs
(`app-backend/docs/`, `tm-api/docs/`) — this repo does not duplicate their
annotations here, since it isn't the service that owns those routes.

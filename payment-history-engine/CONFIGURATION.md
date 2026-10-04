# Configuration

Every environment variable read by this service, found via:

```bash
grep -rn "os.Getenv\|os.LookupEnv" --include="*.go" . | grep -v _test.go
```

Loaded from a `.env` file at startup via `godotenv.Load()` (see `main.go`) if
one exists; otherwise real process environment variables are used. A sample
file is at `.env.example` — copy it to `.env` and fill it in.

Variables marked **Required** must be set (mostly checked explicitly in
`main.go`'s startup validation) or the process logs an error and exits
immediately without serving anything.

## Database

### `DB_TYPE`
- **Required**: no (defaults to `postgres`)
- **Example**: `postgres` or `sqlite`
- **What it is**: which GORM driver to use for the primary "wallet" database
  — the same database `app-backend` uses for `user_wallets` and
  `payment_history`.
- **How to get a real value**: use `postgres` against a real Postgres
  instance for anything beyond local development; use `sqlite` locally/in CI
  to avoid needing Postgres at all.

### `DB_CONNECTION_STRING`
- **Required**: yes
- **Example (postgres)**: `postgres://wallet_user:wallet_pass@localhost:5432/trovo_wallet?sslmode=disable`
- **Example (sqlite)**: `./local-wallet.db`
- **What it is**: the connection string (or, for SQLite, a file path) for the
  primary wallet database.
- **How to get a real value**: for Postgres, the shape is
  `postgres://<user>:<password>@<host>:<port>/<database>?sslmode=disable`.
  Stand up a local instance with Docker:
  ```bash
  docker run --name wallet-db -e POSTGRES_USER=wallet_user -e POSTGRES_PASSWORD=wallet_pass -e POSTGRES_DB=trovo_wallet -p 5432:5432 -d postgres:16
  ```
  then use `postgres://wallet_user:wallet_pass@localhost:5432/trovo_wallet?sslmode=disable`.
  For local/CI use, skip Postgres entirely: set `DB_TYPE=sqlite` and point
  this at any writable file path, e.g. `./local-wallet.db`.

### `DB_MAX_OPEN_CONNECTIONS`
- **Required**: no (defaults to `50`)
- **Example**: `50`
- **What it is**: max open connections in the primary DB's connection pool
  (Postgres mode only).
- **How to get a real value**: tune based on your Postgres instance's
  `max_connections` and how many other services share it; the default of 50
  is a reasonable starting point.

### `DB_MAX_IDLE_CONNECTIONS`
- **Required**: no (defaults to `50`)
- **Example**: `50`
- **What it is**: max idle connections kept open in the pool (Postgres mode
  only).
- **How to get a real value**: usually matched to `DB_MAX_OPEN_CONNECTIONS`;
  lower it if you see connection churn become a bottleneck.

### `ROACH_DB_TYPE`
- **Required**: no
- **Example**: `sqlite` (unset/anything else means "use Postgres wire
  protocol against CockroachDB", via `gorm.io/driver/postgres`)
- **What it is**: escape hatch to point RoachDB at a local SQLite file
  instead of a real CockroachDB cluster.
- **How to get a real value**: set to `sqlite` for local dev/CI; leave unset
  in any real environment.

### `CDB_CONNECTION_STRING`
- **Required**: yes (the process calls `log.Fatal` inside `OpenRoachDB` if
  this is unusable)
- **Example (CockroachDB/Postgres)**: `postgresql://root@localhost:26257/payment_history?sslmode=disable`
- **Example (sqlite mode)**: `./local-roach.db`
- **What it is**: connection string for "RoachDB" — this engine's own
  CockroachDB database, used for `tracked_wallets`, `tracked_addresses`,
  `monitored_cursors` and `monitored_account_cursors` (see
  [INTEGRATION.md](./INTEGRATION.md)). Despite the name, it's read with the
  Postgres driver (`gorm.io/driver/postgres`), since CockroachDB speaks the
  Postgres wire protocol.
- **How to get a real value**: run a local single-node CockroachDB (see
  [DEPLOYMENT.md](./DEPLOYMENT.md#standing-up-a-local-cockroachdb-if-you-want-the-real-thing-instead-of-sqlite)):
  ```bash
  cockroach start-single-node --insecure --listen-addr=localhost:26257 --background
  cockroach sql --insecure -e "CREATE DATABASE payment_history;"
  ```
  then use `postgresql://root@localhost:26257/payment_history?sslmode=disable`.
  Or set `ROACH_DB_TYPE=sqlite` and point this at a file path instead.

### `PURGE_TABLES`
- **Required**: no
- **Example**: `1`
- **What it is**: a destructive one-shot switch. If set to exactly `1`, the
  process drops `tracked_wallets`, `tracked_addresses` and
  `monitored_cursors` from RoachDB and exits immediately — it does **not**
  start the worker loops. Use with care; there is no confirmation prompt.
- **How to get a real value**: leave unset. Only set to `1` deliberately,
  e.g. to reset a broken local dev RoachDB.

## Base (blockchain) connection

### `BASE_RPC_URL`
- **Required**: yes
- **Example**: `https://sepolia.base.org` (testnet) or your own/provider
  mainnet endpoint
- **What it is**: the JSON-RPC endpoint this engine polls for blocks and
  `eth_getLogs` queries. Also used by the `GET /ready` probe to confirm the
  node is reachable.
- **How to get a real value**: for testing, Base Sepolia's public endpoint
  (`https://sepolia.base.org`) needs no signup. For production, get an
  API key from a provider like Alchemy (https://www.alchemy.com/) or
  Infura (https://www.infura.io/), or run your own `base-node`
  (https://docs.base.org/tools/node-providers).

### `RPC_TIMEOUT`
- **Required**: no (defaults to `30s`)
- **Example**: `30s`
- **What it is**: the longest any single request to `BASE_RPC_URL` may
  take. If the RPC stops answering, requests give up after this long, and
  after 3 failures in a row the client refuses RPC calls at once for 10
  seconds at a time until a request gets through again. A block whose
  data could not be read is retried on the next poll, never skipped.
- **How to get a real value**: leave unset unless your provider is slow
  for large `eth_getLogs` ranges.

### `TRACK_ADDRESS_CONCURRENCY`
- **Required**: no (defaults to `8`)
- **Example**: `8`
- **What it is**: how many tracked addresses are back-filled from the
  chain at the same time. Each back-fill holds RPC requests and memory, so
  this caps both when many addresses are tracked at once.
- **How to get a real value**: leave unset; raise it only if your RPC
  plan allows more parallel requests.

### `BASE_CHAIN_ID`
- **Required**: no (defaults to `84532`, Base Sepolia)
- **Example**: `8453` (Base mainnet) or `84532` (Base Sepolia)
- **What it is**: the chain ID used to build the EIP-1559 transaction signer
  (`types.LatestSignerForChainID`) when recovering a transaction's sender.
- **How to get a real value**: use `8453` for Base mainnet, `84532` for Base
  Sepolia testnet — see https://chainlist.org/?search=base.

### `NATIVE_ASSET_CODE`
- **Required**: yes
- **Example**: `GAS` (per `.env.example`) or `ETH`
- **What it is**: the display code used for native-currency transfers (Base's
  gas token) when recording payment history, and as the fallback asset code
  in `SavePaymentHistory` when none is given.
- **How to get a real value**: pick whatever label your product wants shown
  for native transfers — there's no on-chain source for this, it's purely
  cosmetic.

### `BLOCKCHAIN_NETWORK_PASSPHRASE`, `BLOCKCHAIN_BASE_RESERVE`
- **Required**: no
- **Example**: leave empty (see `.env.example`)
- **What they are**: vestigial. These map to Stellar-era concepts (signature
  domain separation and per-subentry account reserve) that have no EVM/Base
  equivalent. The code reads them but handles them being unset gracefully
  (`network.GetBlockchainNetworkPassPhrase`/`GetBlockchainBaseReserve`'s doc
  comments say so explicitly), and `main.go` no longer requires them.
- **How to get a real value**: don't set them.

### `BLOCKCHAIN_SWAP_DESTINATION_MIN`
- **Required**: no
- **Example**: `0`
- **What it is**: a decimal minimum used by `network.GetBlockchainSwapDestinationMin`;
  parsed with `decimal.NewFromString`, falling back to `0` if unset/invalid.
- **How to get a real value**: only relevant if/when swap support is wired
  up; leave unset otherwise.

## Wallet key derivation

### `MNEMONIC_TEMP_ACCOUNTS`
- **Required**: yes
- **Example**: `24 words sha` (a placeholder in `.env.example` — replace with
  a real BIP-39 mnemonic)
- **What it is**: the seed phrase this engine derives temporary/intermediate
  accounts from (`network.main.go`).
- **How to get a real value**: generate a real BIP-39 mnemonic with a wallet
  tool you trust (e.g. `openssl rand` piped through a BIP-39 wordlist
  generator, or a hardware wallet's setup flow) — **never reuse a
  mainnet-funded mnemonic in a dev/test environment**, and never commit a
  real one to `.env`.

### `MARKET_MAKING_SALT`, `MNEMONIC_MARKET_MAKING`
- **Required**: yes (both)
- **Example**: any string for the salt; a real BIP-39 mnemonic for the other
- **What they are**: used by `internal/blockchainalgofuncs` to
  deterministically derive the market-making account from a mnemonic + salt.
- **How to get a real value**: same guidance as `MNEMONIC_TEMP_ACCOUNTS` for
  the mnemonic; the salt can be any secret string you control (treat it like
  a password — generate with `openssl rand -hex 16`).

### `MM_FEE_COLLECTION_CHANNEL_ACCOUNT`, `MARKET_MAKING_FEE_WALLET`
- **Required**: yes (both, checked at boot)
- **Example**: a Base address, e.g. `0x0000000000000000000000000000000000000000`
- **What they are**: required by the startup check in `main.go`, but **not
  read anywhere else in this codebase** (grepped) — they appear to be
  reserved for the market-making/trade-stream feature, which is currently an
  idle stub (see `MonitorTradeStream`'s doc comment: "Base has no trade
  stream to watch yet").
- **How to get a real value**: set to any syntactically valid non-empty
  value to satisfy the startup check until this feature is implemented; a
  real value will matter once trade-stream processing is built out.

### `BULK_PAYMENT_SALT`, `MNEMONIC_BULK_PAYMENT`
- **Required**: no (not in the hard-required list, but read by
  `internal/blockchainalgofuncs`)
- **Example**: same shape as the market-making pair above
- **What they are**: derive a bulk-payment account, analogous to the
  market-making derivation.
- **How to get a real value**: same guidance as `MNEMONIC_TEMP_ACCOUNTS`.

### `ENCODER_SALT`
- **Required**: no (but read unconditionally by two functions in
  `internal/blockchainalgofuncs/algofuncs.go`, so effectively required if
  those code paths are exercised)
- **Example**: a random 16+ byte hex string, e.g. `openssl rand -hex 16`
- **What it is**: a salt used in encoding/decoding helper functions.
- **How to get a real value**: `openssl rand -hex 16`; store it securely and
  keep it stable once set (changing it invalidates anything previously
  encoded with it).

## Caching (Redis)

### `ENABLE_CACHING`
- **Required**: yes (must be set, even if to a value other than `1`)
- **Example**: `1` (on) or `0` (off)
- **What it is**: toggles Redis caching. When `1`, `REDIS_HOST`,
  `REDIS_PORT` and `CACHING_PARAMETER` become required and the process
  `log.Fatal`s if the initial `PING`/test-write fails.
- **How to get a real value**: `0` for local dev unless you specifically want
  to test the caching path.

### `REDIS_HOST`, `REDIS_PORT`
- **Required**: only if `ENABLE_CACHING=1`
- **Example**: `localhost`, `6379`
- **What they are**: Redis connection host/port.
- **How to get a real value**: run Redis locally with
  `docker run --name redis -p 6379:6379 -d redis:7`, then use
  `REDIS_HOST=localhost`, `REDIS_PORT=6379`.

### `REDIS_PASSWORD`
- **Required**: no
- **Example**: empty for a local unauthenticated Redis
- **What it is**: Redis `AUTH` password.
- **How to get a real value**: whatever your Redis instance is configured
  with; leave empty for the local Docker command above.

### `CACHING_PARAMETER`
- **Required**: only if `ENABLE_CACHING=1`
- **Example**: `v1`
- **What it is**: a string mixed into every cache key (`internal/cache/main.go`),
  letting you invalidate all cached entries by changing it.
- **How to get a real value**: any short string; bump it whenever you want to
  invalidate the whole cache (e.g. after a schema change to cached payloads).

## Monitoring / health server

### `HEALTH_PORT`
- **Required**: no (defaults to `8080`)
- **Example**: `8080`
- **What it is**: the port the `GET /health`, `GET /ready` and
  `GET /swagger/*` endpoints are served on.
- **How to get a real value**: leave at the default unless it conflicts with
  something else on the host.

### `APP_VERSION`
- **Required**: no
- **Example**: a short git SHA, e.g. `a1b2c3d`
- **What it is**: reported as `version` in the health/readiness responses.
  Normally stamped at build time via the Dockerfile's `-ldflags -X` (see
  [DEPLOYMENT.md](./DEPLOYMENT.md)); this env var is a fallback read only if
  that build-time stamp is empty.
- **How to get a real value**: `git rev-parse --short HEAD`.

## Cursors (manual override / resume)

### `LAST_CURSOR`
- **Required**: no
- **Example**: a Base block number as a string, e.g. `12345678`
- **What it is**: an override for the payment-stream resume cursor,
  compared against the DB-stored cursor and whichever is more recent wins
  (`GetLastCursor` in `main.go`).
- **How to get a real value**: normally leave unset and let the DB-stored
  cursor (`monitored_cursors`) drive resumption; only set this to force a
  specific starting block (e.g. after a manual intervention).

### `TRADE_RESUME_CURSOR`
- **Required**: no (defaults to `"0"`)
- **What it is**: read by `GetTradeResumeCursor`, but the trade stream itself
  is currently an idle stub — has no present effect.
- **How to get a real value**: leave unset.

## Email validation

### `ENABLE_EMAIL_VALIDATION`
- **Required**: no
- **Example**: `1` or `0`
- **What it is**: if `1`, requires `MAILGUN_VALIDATOR_API_KEY` to be set.
- **How to get a real value**: `0` unless you specifically need Mailgun email
  validation exercised.

### `MAILGUN_VALIDATOR_API_KEY`
- **Required**: only if `ENABLE_EMAIL_VALIDATION=1`
- **Example**: `key-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`
- **What it is**: Mailgun's email-validation API key.
- **How to get a real value**: create a free Mailgun account
  (https://www.mailgun.com/) and copy the API key from its dashboard.

## Payments

### `MIN_SENDABLE_AMOUNT`
- **Required**: no
- **Example**: `1` (any value other than empty or `"0"` enables the check)
- **What it is**: used by
  `internal/components/payments/errors/payment_amount_below_min_allowed.go`
  to report a minimum sendable amount in error messages — part of the
  vestigial Gin-error package (see [README.md](./README.md)), not currently
  reachable through this service's live HTTP surface.
- **How to get a real value**: leave unset unless this code path becomes
  reachable.

## Third-party / misc

### `IPAPI_HOST`, `IPAPI_KEY`
- **Required**: no
- **Example**: empty (per `.env.example`)
- **What they are**: read by `internal/components/payments/db/geo.go`'s
  `GetGeoInfo`, an IP-geolocation lookup helper. Not called from anywhere in
  `main.go`'s worker loops — appears to be leftover/unused in this service.
- **How to get a real value**: leave unset unless you wire this helper into
  something. If needed, ip-api.com (https://ip-api.com/) is a common free
  option for `IPAPI_HOST`.

### `CONNECTION_WARNING_WEBHOOK`
- **Required**: no
- **Example**: a Discord webhook URL
- **What it is**: overrides the hardcoded fallback Discord webhook URL in
  `internal/db/main.go` used to post a warning when the DB connection pool
  fills up (`PrintDBStats`). **The hardcoded fallback in source has been
  committed to git history and should be treated as leaked/rotated** — set
  this explicitly rather than relying on the fallback.
- **How to get a real value**: create a Discord webhook on the channel you
  want alerts in (Server Settings → Integrations → Webhooks → New Webhook),
  copy its URL.

### `EXPANSION_NETWORK_ERROR_WEBHOOK`, `FAILED_PAYMENT_ERROR_WEBHOOK`
- **Required**: no
- **Example**: a Discord webhook URL (must be longer than 50 characters to
  override the built-in default — see `internal/network/main.go`)
- **What they are**: Discord webhooks for network-expansion and
  failed-payment error alerts respectively.
- **How to get a real value**: same as `CONNECTION_WARNING_WEBHOOK` above —
  create a Discord webhook per channel you want each alert type posted to.

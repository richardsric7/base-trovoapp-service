# Configuration Reference

This document lists every environment variable `app-backend` reads,
grouped by area. For each one: an example value (never a real secret —
these are safe placeholders), what it actually changes about the running
app, and concrete instructions for getting a real value.

Start from `.env.example` in this directory (`cp .env.example .env`) and
fill in what you need — most variables have sensible defaults for local
development, and only a subset is required to boot at all (see `main.go`'s
`requiredEnvironmentVariables` list, or [DEPLOYMENT.md](DEPLOYMENT.md)).

A note on secrets: several variables below are blockchain private
keys/mnemonics or fee-wallet addresses that control real (or testnet)
funds. These are normally issued and rotated by whoever administers Trovo's
on-chain treasury, not something you generate ad hoc — for local/dev work,
ask a team lead for a **testnet** set, or generate your own testnet-only
wallet and never reuse it for anything real.

---

## Server & runtime

**`PORT`**
- Example: `8080`
- What it does: TCP port the HTTP server listens on. If unset, defaults to `8080` (`0.0.0.0:8080`, or `localhost:8080` on Windows).
- How to get a real value: pick any free local port, or let your deployment platform inject it.

**`GIN_MODE`**
- Example: `release`
- What it does: When set to `release`, disables Gin's verbose debug logging/output. Leave unset for local development to get full request logs.
- How to get a real value: set to `release` in any deployed environment; omit locally.

**`ORGANISATION`**
- Example: `trovotech`
- What it does: A free-text organisation identifier used in a few places in the codebase (legacy naming).
- How to get a real value: any short identifier for your deployment; not validated.

**`MIGRATE_ONLY`**
- Example: `1`
- What it does: When `1`, the process only opens the primary DB, runs `AutoMigrate`, and exits — it does not start the HTTP server and skips the ~50 required-env-var check the serving process needs. Used by CI as a one-shot migration step (see [DEPLOYMENT.md](DEPLOYMENT.md)).
- How to get a real value: leave unset when running the server normally; set to `1` only when running a dedicated migration step.

**`DB_AUTOMIGRATE`**
- Example: `1`
- What it does: When not `0`, GORM `AutoMigrate` runs against every model at boot, creating/altering tables to match the Go structs. Set to `0` to skip this (faster boot; use once the schema is already up to date). See [DEPLOYMENT.md](DEPLOYMENT.md) for why this matters for deploy speed.
- How to get a real value: `1` for local dev (`.env.example` default); `0` on a serving instance once migrations have been run separately.
- **Multi-instance note:** if two or more instances (or a `MIGRATE_ONLY=1` job and a serving instance) boot with `DB_AUTOMIGRATE` enabled at the same time, they no longer race each other's `AutoMigrate` calls against the same schema — a shared database-backed lock (`internal/sharedconfig/distributed_lock.go`) makes the second one wait for the first to finish instead. This works identically on SQLite and Postgres. See [DEPLOYMENT.md](DEPLOYMENT.md) section 4 for the full explanation, including a related dependency fix that was needed for `AutoMigrate` to work correctly against Postgres at all on anything but the very first boot.

**`ENABLE_AUTH_MIDDLEWARE`**
- Example: `1`
- What it does: When explicitly set to `0`, disables signature/API-key verification entirely (every request is treated as authenticated) — **never do this outside local debugging**. Any other value (including unset) keeps auth enforced.
- How to get a real value: leave unset (auth enforced) everywhere except a throwaway local debug session.

**`LOG_IP_ADDRESS`**
- Example: `0`
- What it does: When `1`, logs the caller's IP address on certain requests, for debugging.
- How to get a real value: `0`/unset normally; `1` temporarily while investigating an issue.

**`LOG_TARGET_USER`** / **`LOG_TARGET_USER_PK`**
- Example: `debuguser` / `0xA1B2...` (a wallet address)
- What it does: When a request's username or wallet address matches, extra diagnostic logging is emitted for that request only — a way to trace one user's traffic without turning on global verbose logging.
- How to get a real value: set to a specific username/address you're actively debugging; leave unset otherwise.

**`DEFAULT_ASSET_IMAGE_URL`** / **`NATIVE_ASSET_IMAGE_URL`**
- Example: `https://cdn.trovo.app/assets/default.png`
- What it does: Fallback image URLs shown for assets that don't have their own icon configured.
- How to get a real value: point at any publicly reachable image URL your frontend can render.

**`NATIVE_ASSET_CODE`**
- Example: `ETH`
- What it does: The asset code for the chain's native currency (used for balance/pricing logic).
- How to get a real value: matches whatever native asset your `BASE_RPC_URL` network uses.

---

## Database (primary)

**`DB_TYPE`**
- Example: `postgres`
- What it does: Selects the database driver — `postgres` or `sqlite`. Defaults to `postgres` if unset.
- How to get a real value: `sqlite` for zero-setup local dev (see `SQLITE_DB_PATH`/`DB_CONNECTION_STRING` below); `postgres` for anything shared/production.

**`DB_CONNECTION_STRING`**
- Example (Postgres): `host=localhost user=trovo password=trovo dbname=trovo_wallet port=5432 sslmode=disable`
- Example (SQLite): `bantupay.sqlite`
- What it does: The connection string (Postgres) or file path (SQLite, when `DB_TYPE=sqlite`) for the primary application database. Required — the app refuses to start without it.
- How to get a real value: for Postgres, run a local instance (`docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=trovo -e POSTGRES_USER=trovo -e POSTGRES_DB=trovo_wallet postgres:16`) and build the DSN from those credentials; for SQLite, any writable file path works and the file is created automatically.

**`DB_MAX_OPEN_CONNECTIONS`** / **`DB_MAX_IDLE_CONNECTIONS`**
- Example: `50`
- What it does: Postgres connection-pool sizing. Both default to `50` if unset.
- How to get a real value: tune based on your Postgres instance's `max_connections` and expected concurrent load; the defaults are fine for local dev.

**`SQLITE_DB_PATH`**
- Example: `dbs/bantupay.sqlite`
- What it does: File path for a secondary, standalone SQLite connection helper (`OpenSqliteDB`) used in some local/test paths. Defaults to `dbs/bantupay.sqlite` if unset.
- How to get a real value: any writable local path; not used in a Postgres deployment.

**`CDB_CONNECTION_STRING`**
- Example: `postgres://trovo:trovo@localhost:5432/trovo_payment_history?sslmode=disable`
- What it does: Connection string for the separate "RoachDB" (CockroachDB-compatible Postgres) database that stores payment-history tracking tables (`TrackedWallet`, `TrackedAddress`). Required to boot the serving process (not required under `MIGRATE_ONLY=1`).
- How to get a real value: point at a second local Postgres database (or a real CockroachDB cluster in production) — it can be the same Postgres server as `DB_CONNECTION_STRING`, just a different database name.

---

## Redis / cache

**`ENABLE_CACHING`**
- Example: `1`
- What it does: When `1`, connects to Redis and caches selected HTTP responses (curated assets, rates, announcements, etc.) to reduce DB load. When unset/`0`, every request hits the DB directly.
- How to get a real value: `0` for simplest local dev; `1` once you have Redis running.
- **Multi-instance note:** if you run more than one instance of this service, `ENABLE_CACHING=1` (a working Redis) is also what powers live delivery of P2P order/offer/escrow/dispute updates over the websocket (`GET /v1/users/websocket/:identifier`) to a connection that happens to be on a *different* instance than the one that processed the update — see `internal/sharedconfig/realtime.go`. With caching disabled, everything still works correctly (push notifications and REST responses are unaffected), you just won't get that live socket nudge across instances.

**`REDIS_HOST`** / **`REDIS_PORT`** / **`REDIS_PASSWORD`**
- Example: `localhost` / `6379` / `` (empty)
- What it does: Connection details for the Redis instance used when `ENABLE_CACHING=1`.
- How to get a real value: run Redis locally with `docker run -d -p 6379:6379 redis:7`, or use a managed Redis (e.g. a DigitalOcean/AWS Redis instance) in production — its host/port/password come from that provider's dashboard.

**`CACHING_PARAMETER`**
- Example: `v2`
- What it does: An arbitrary string mixed into cache keys, letting you invalidate all cached entries at once by changing it (a cheap cache-busting knob).
- How to get a real value: any short string; change it whenever you need to force a full cache invalidation.

---

## Rate limiting

A Redis-backed fixed-window limiter (`INCR`+`EXPIRE`) protecting every
signature- or API-key-authenticated route worth protecting: the
payment-history, payment, and swap endpoints, the `GET /v1/users/:targetUser`
lookup, and all 33 white-label "service link" endpoints under
`/v1/servicelinks/...`/`/v1/trovo-api/...` (login/authorize/event handshakes,
user-info lookup, onboarding, KYC, minting, tokenization, stakeholder
documents, etc). It's Redis-backed rather than in-memory on purpose: an
in-memory counter is per-process, so with more than one instance behind a
load balancer it would silently multiply the effective limit by the
instance count. See `internal/middleware/rate_limit_middleware.go`.

**`RATE_LIMIT_ENABLED`**
- Example: `1`
- What it does: When unset/`1` (and Redis is enabled via `ENABLE_CACHING=1`), rate limiting is active. Set to `0` to disable it entirely, e.g. for local dev or load testing.
- How to get a real value: leave unset in production (defaults on); set `0` only when you deliberately want it off.
- **Requires Redis:** if `ENABLE_CACHING=0` or Redis is unreachable, the middleware no-ops (fails open) regardless of this flag — same graceful-degradation posture as the caching layer itself.

**`RATE_LIMIT_REQUESTS_PER_MINUTE`**
- Example: `30`
- What it does: A global default limit (requests per 60s window) applied to every rate-limited route that doesn't have its own override. Each route also has a hardcoded default in code (roughly: 60/min for read/lookup routes, 30/min for login/authorize/event handshakes, 10-20/min for mutating routes like onboarding/KYC/minting/tokenization) - this env var overrides all of them at once.
- How to get a real value: leave unset to use the per-route defaults baked into the code; set it to tune globally without a redeploy.

**`RATE_LIMIT_<KEY>_PER_MINUTE`** (per-route override)
- Example: `RATE_LIMIT_SWAP_PER_MINUTE=5`
- What it does: Overrides the limit for one specific route's key (uppercased, hyphens to underscores - e.g. the `"swap"` key becomes `RATE_LIMIT_SWAP_PER_MINUTE`, `"payment-history"` becomes `RATE_LIMIT_PAYMENT_HISTORY_PER_MINUTE`). Takes precedence over both the route's hardcoded default and the global `RATE_LIMIT_REQUESTS_PER_MINUTE`.
- How to get a real value: only set when one specific endpoint needs a different budget than the rest; look up the exact key string at the route's `middleware.RateLimitMiddleware(gc, "<key>", ...)` call site.

**Per-service-link override (no env var - set from tm-api).** Every
`ServiceLink` row (`internal/components/servicelinks/models/servicelink.go`)
carries a `RateLimitPerMinute` column, editable from tm-api's Service Links
admin page. `0` (the default) means no override - the route's own
default/env-based limit applies as normal. A positive value overrides
**every** rate-limited API-key route that service link calls with that
same per-minute budget, and takes precedence over both the route default
and any `RATE_LIMIT_*` env var, since a per-tenant setting was
deliberately configured for that partner and should always win. This is
how an enterprise partner with unusually high (or low) legitimate traffic
gets a budget different from every other partner's, without an env var
change or redeploy. Only applies to API-key-authenticated routes - a
route reached with `X-TW-SIGNER` (a wallet-app request, not a service
link) never does this lookup.

---

## Blockchain / network (Base chain)

**`BASE_RPC_URL`**
- Example: `https://sepolia.base.org`
- What it does: The JSON-RPC endpoint used for all on-chain reads/writes (balances, transaction submission, event queries). Required to boot.
- How to get a real value: for testnet development, use a public endpoint like `https://sepolia.base.org` (Base Sepolia testnet); for production, sign up with an RPC provider such as [Alchemy](https://www.alchemy.com/) or [Infura](https://www.infura.io/) and create a Base mainnet app to get a dedicated endpoint URL.

**`RPC_TIMEOUT`**
- Example: `30s` (the default)
- What it does: The longest any single request to `BASE_RPC_URL` may take. The whole process shares one RPC client; if the RPC stops answering, requests give up after this long instead of waiting forever, and after 3 failures in a row the client refuses RPC calls at once ("blockchain RPC unavailable", for 10 seconds at a time, until a request gets through again) so a stalled RPC does not pile up waiting requests and memory. It recovers on its own.
- How to get a real value: leave unset. Raise it only if your RPC provider is legitimately slow for some calls (e.g. large `eth_getLogs` ranges).

**`BASE_CHAIN_ID`**
- Example: `84532`
- What it does: The numeric chain ID for the network `BASE_RPC_URL` points at (Base Sepolia testnet = `84532`, Base mainnet = `8453`). Used to sign transactions for the correct network.
- How to get a real value: match it to whichever Base network `BASE_RPC_URL` points at.

**`BLOCKCHAIN_NETWORK_PASSPHRASE`** / **`BLOCKCHAIN_BASE_RESERVE`**
- Example: `` (leave empty)
- What it does: Vestigial — carried over from this codebase's pre-Base (Stellar) era. Both are handled gracefully when unset and are no longer required.
- How to get a real value: leave empty.

**`BLOCKCHAIN_DATA_CACHE_LIFETIME`**
- Example: `94608000`
- What it does: How long (in seconds) certain blockchain-derived data is cached. Defaults to `94608000` (≈3 years) if unset.
- How to get a real value: leave at the default unless you have a specific reason to invalidate on-chain data caches more often.

**`BLOCKCHAIN_SWAP_DESTINATION_MIN`**
- Example: `0.0001`
- What it does: The minimum destination-asset amount a swap must produce to be accepted (guards against dust/near-zero swaps).
- How to get a real value: a small positive decimal appropriate for your assets' smallest meaningful unit.

**`DOLLAR_ASSET`**
- Example: `USDB:0x0000000000000000000000000000000000000000`
- What it does: The asset code + contract address pair the app treats as "the dollar-pegged asset" for USD-denominated pricing. Required to boot.
- How to get a real value: `<ASSET_CODE>:<contract address>` for whichever stablecoin your deployment uses (e.g. a USDC/USDB contract on Base).

**`TROV_ASSET_CONTRACT_ADDRESS`** / **`USDB_B20_TOKEN_ADDRESS`**
- Example: `0x0000000000000000000000000000000000000000`
- What it does: Contract addresses for the platform's native TROV token and a B20-standard USDB token, used in pricing/order-book logic.
- How to get a real value: the deployed contract address for each token on your target network.

**`USE_ASSET_FOR_NATIVE_PRICE`**
- Example: `USDB:0x0000000000000000000000000000000000000000`
- What it does: An `assetCode:contractAddress` pair telling the pricing engine which asset to use as a proxy when pricing the chain's native asset.
- How to get a real value: same shape as `DOLLAR_ASSET`; usually the same value.

**`NAIRA_ASSET`** / **`ENABLE_NAIRA_ASSET_BY_DEFAULT`**
- Example: `NGNB:0x0000000000000000000000000000000000000000` / `0`
- What it does: When `ENABLE_NAIRA_ASSET_BY_DEFAULT=1`, new wallets are seeded with a Naira-pegged asset by default; `NAIRA_ASSET` (must be ≥60 chars, i.e. a full `code:address` string) identifies which one.
- How to get a real value: only needed if you're enabling this feature; use your deployed Naira-pegged token's `code:address`.

**`ENABLE_DOLLAR_ASSET_BY_DEFAULT`** / **`ENABLE_TROV_ASSET_BY_DEFAULT`**
- Example: `1`
- What it does: When `1`, new wallets are automatically seeded with the dollar asset / TROV asset respectively.
- How to get a real value: `0` or `1` depending on your product defaults; no external dependency.

---

## Wallet signers & mnemonics (secrets)

These are BIP-39 mnemonics for internal signer wallets the app uses to
perform specific categories of on-chain actions automatically.

**`MNEMONIC_BULK_PAYMENT`** / **`BULK_PAYMENT_SALT`**
- Example: mnemonic / `openssl rand -hex 16` output
- What it does: Signer mnemonic and salt for the bulk-payment feature's derived sub-wallets.
- How to get a real value: same pattern as above.

**`MNEMONIC_MARKET_MAKING`** / **`MARKET_MAKING_SALT`**
- Example: mnemonic / `openssl rand -hex 16` output
- What it does: Signer mnemonic and salt for market-making sub-wallets.
- How to get a real value: same pattern as above.

**`ENCODER_SALT`**
- Example: output of `openssl rand -hex 16`
- What it does: Salt used generically when deriving/encoding wallet addresses from usernames elsewhere in the codebase.
- How to get a real value: `openssl rand -hex 16`.

**`VERIFICATION_CODE_SALT`**
- Example: `WhateverYouWant` (ships as a placeholder in `.env.example` — replace it)
- What it does: Salt mixed into generated email/SMS verification codes.
- How to get a real value: `openssl rand -hex 16`.

## Wallets: Safe accounts, bundler and paymaster

Every user wallet is a Safe (v1.4.1 with the ERC-4337 Safe4337Module)
owned by the user's key, and every send is a UserOperation the backend
builds, the user's app signs, and the backend submits to our bundler.
The wallet pays its own gas - in the stablecoin the user chose on their
profile (through the paymaster, see [`paymaster/`](../paymaster/README.md)),
or in ETH. See `internal/aa` and `internal/components/users/services/wallet_operations.go`.
Stablecoin gas the paymaster collects is recorded in `fee_collections`
(type `GAS`); a wallet that owes the paymaster gas (see INTEGRATION.md,
"Gas debt") pays in ETH and settles the debt automatically. Nothing
needs configuring for either beyond `PAYMASTER_ADDRESS`.

**`BUNDLER_URL`**
- Example: `http://bundler.internal:4337/rpc`
- What it does: JSON-RPC endpoint of our self-hosted ERC-4337 bundler (EntryPoint v0.7). Every wallet operation is gas-estimated with and submitted to it. Without it, no sends can be prepared.
- How to get a real value: the URL of the bundler you run (e.g. Rundler, Alto or Skandha) for the same chain as `BASE_RPC_URL`. Configure it for EntryPoint `0x0000000071727De22E5E9d8BAf0edAc6f37da032`.

**`PAYMASTER_ADDRESS`**
- Example: `0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9`
- What it does: The deployed `TrovoTokenPaymaster`. Wallets paying gas in a stablecoin approve it to take the fee. Leave empty to disable stablecoin gas (everyone pays in ETH).
- How to get a real value: printed by the paymaster deploy script and saved in `paymaster/contracts/deployments/<chainId>.json` (see `paymaster/DEPLOYMENT.md`).

**`PAYMASTER_QUOTE_SERVICE_URL`** / **`PAYMASTER_QUOTE_SERVICE_API_KEY`**
- Example: `http://paymaster-quotes.internal:8090` / `3f9c0a...e1` (64 hex)
- What it does: Where the backend gets signed gas quotes in stablecoins, and the API key it sends (`X-API-Key`). If the quote service cannot price a user's chosen token, the operation falls back to ETH gas.
- How to get a real value: the quote service's address, and one of the keys in its `API_KEYS` Vault setting (see `paymaster/CONFIGURATION.md`).

**`WALLET_OPERATION_VALIDITY`** / **`SHARED_WALLET_OPERATION_VALIDITY`**
- Example: `10m` / `24h`
- What it does: How long the owners have to sign a prepared operation - a single-signer wallet (default 10 minutes) or a shared wallet waiting for approvers (default 24 hours, at most the quote service's `QUOTE_MAX_VALIDITY`). An operation not signed in time expires and must be started again.
- How to get a real value: keep the defaults unless your approval process needs longer.

**`WALLET_OPERATION_GAS_BUFFER_PERCENT`**
- Example: `20`
- What it does: Percentage added to the bundler's gas estimates (default 20). Unused gas is not charged (and a stablecoin pre-charge is refunded), so this only affects the maximum shown to the user.
- How to get a real value: keep the default.

**`ENTRYPOINT_ADDRESS`**, **`SAFE_MODULE_SETUP_ADDRESS`**, **`SAFE_4337_MODULE_ADDRESS`** (and the `SAFE_*` addresses under [Tokenization](#tokenization))
- Example: leave empty for the defaults
- What it does: Override the contracts wallets are built from, for local chains only. Defaults are the canonical Base / Base Sepolia deployments: EntryPoint v0.7 `0x0000000071727De22E5E9d8BAf0edAc6f37da032`, SafeProxyFactory v1.4.1, SafeL2 v1.4.1, SafeModuleSetup v0.3.0 `0x2dd68b007B46fBe91B9A7c3EDa5A7a1063cB5b47`, Safe4337Module v0.3.0 `0x75cf11467937ce3F2f357CE24ffc3DBF8fD5c226`, MultiSendCallOnly v1.4.1 `0x9641d764fc13c8B624c04430C7356C1C7C8102e2`. **Changing them in production changes every user's wallet address** - wallet-core computes addresses with the defaults.
- How to get a real value: leave empty on Base / Base Sepolia.

---

## Channel accounts (transaction fee payers)

> **Multi-instance note:** which channel account is free to hand out is
> coordinated across every running instance through a `channel_accounts`
> database table (address + status only — the actual signing keys parsed
> from `CHANNEL_ACCOUNTS` below stay in each instance's memory, never
> written to the DB). You don't need to configure anything extra for
> this — it works the same whether you're running one instance or ten,
> and works identically on SQLite (local dev) and Postgres (production).
> See `internal/sharedconfig/channel_accounts.go` if you're curious how it
> works under the hood.

**`CHANNEL_ACCOUNTS`**
- Example: `0xabc...priv1,0xdef...priv2,0x123...priv3`
- What it does: A comma-separated list of signer keys ("channel accounts") the app round-robins through to pay gas/sign routine transactions, so a single hot wallet isn't bottlenecked by sequence numbers.
- How to get a real value: generate N fresh testnet wallets (e.g. `cast wallet new`, or `ethers.Wallet.createRandom()` in Node), fund them lightly, and list their private keys here — never reuse mainnet keys with real funds for this in a non-production environment.

**`CHANNEL_ACCOUNT_FUNDER`**
- Example: a funded wallet's private key
- What it does: The wallet that tops up channel accounts when their balance drops below `CHANNEL_ACCOUNT_MIN_BALANCE`.
- How to get a real value: a wallet you control, funded with enough native asset to cover top-ups.

**`CHANNEL_ACCOUNT_FUNDING_AMOUNT`** / **`CHANNEL_ACCOUNT_MIN_BALANCE`** / **`CHANNEL_ACCOUNT_MIN_COUNT`**
- Example: `0.05` / `0.01` / `5`
- What it does: How much to fund a channel account with, the balance threshold that triggers a top-up, and the minimum number of channel accounts the pool should maintain (extra ones are auto-generated if `CHANNEL_ACCOUNTS` has fewer).
- How to get a real value: small values are fine for testnet; size for production based on expected transaction volume.

**`CHANNEL_ACCOUNT_RECEIPIENT`**
- Example: `ops@trovo.app`
- What it does: Email address that receives a CSV of any newly auto-generated channel account keys (so an operator can track/fund them). Requires Mailgun to be configured (see Email section).
- How to get a real value: an inbox your ops team monitors.

**`CHECK_CHANNEL_ACCOUNT_BALANCE`**
- Example: `1`
- What it does: When `1`, the startup routine checks and tops up channel-account balances; when `0`/unset, it skips the balance check (faster boot, useful for local dev).
- How to get a real value: `0` locally; `1` in a real deployment.

---

## Fee wallets & amounts

Each `*_FEE_WALLET` below is the wallet address that platform fees for
that specific flow are collected into. Set the **address**: the fees are
sent there by the user's own wallet, so app-backend never needs the key.
(A private key is still accepted for backward compatibility, and only its
address is used - prefer the address.) For local development these can
all point at the same test wallet.

| Variable | Fee for |
| --- | --- |
| `VAT_WALLET` | VAT collected on transactions |
| `SUBWALLET_CREATION_FEE_WALLET` | Creating a sub-wallet |
| `TOKENIZATION_APPLICATION_FEE_WALLET` | Submitting a tokenized-asset application |
| `TOKENIZATION_FEE_WALLET` | Tokenized-asset issuance |
| `CLOSED_GROUP_FEE_WALLET` | Closed-group (private circle) features |
| `ACCOUNT_RECOVERY_FEE_WALLET` | Account recovery |
| `PATRON_FEE_WALLET` | Patron/subscription payments |
| `SHARED_ACCESS_PAYMENT_FEE_WALLET` | Shared-access wallet payments |
| `PAYMENT_FEE_WALLET` | Ordinary payments |
| `SWAP_FEE_WALLET` | Swaps (only read when `SWAP_FEE_ENABLED=1`) |
| `MARKET_MAKING_FEE_WALLET` | Market-making trades (only read when `MARKET_MAKING_FEE_ENABLED=1`) |

- Example value (any of the above): `0x1111111111111111111111111111111111111`
- How to get a real value: an address you control (testnet address for dev); in production these are real treasury addresses managed by whoever administers Trovo's fee collection.

**`SWAP_FEE_ENABLED`** / **`MARKET_MAKING_FEE_ENABLED`**
- Example: `1`
- What it does: Toggles whether swap fees / market-making fees are charged at all.
- How to get a real value: `0` or `1` per your product configuration.

**`MARKET_MAKING_FEE_AMOUNT`**
- Example: `0.5`
- What it does: The fee amount (in whatever unit the market-making flow expects) charged per market-making trade.
- How to get a real value: set per your fee schedule.

**Account recovery fee** (not an environment variable)
- The fee charged when a user turns account recovery on is the `service_fees` row `ACCOUNT_RECOVERY_FEE`: `fee_fixed` is the amount in USD, `fee_asset_code` / `fee_contract_address` the stablecoin it is paid in (converted at the current rate). It is paid from the user's primary wallet in the same operation that turns recovery on, to `ACCOUNT_RECOVERY_FEE_WALLET`. An inactive row (or a fee of 0) charges nothing.

**`SHARED_ACCESS_FEE_ADDRESS`** / **`SHARED_ACCESS_FEE_AMOUNT`** / **`SHARED_ACCESS_FEE_ASSET_CODE`** / **`SHARED_ACCESS_FEE_ASSET_CONTRACT_ADDRESS`**
- Example: `0x1111...` / `0.5` / `USDB` / `0x0000...`
- What it does: Destination address, amount, and asset for the fee charged when setting up shared wallet access.
- How to get a real value: set per your fee schedule.

**`TOKENIZATION_APPLICATION_FEE_AMOUNT`** / **`TOKENIZATION_APPLICATION_FEE_ASSET`**
- Example: `50` / `USDB`
- What it does: Amount and asset charged when a tokenization application is submitted.
- How to get a real value: set per your fee schedule.

**`CRYPTO_WITHDRAWAL_SERVICE_FEE`**
- Example: `0.001`
- What it does: Fee added on top of a crypto withdrawal processed through the 1Liquidity/OneLiquidity withdrawal service.
- How to get a real value: set per your fee schedule.

---

## Wallet minimum balances & activation amounts

These control how much native asset is required/sent when creating or
activating various wallet types — all are plain decimal amounts you set
per your network's economics (testnet values can be tiny).

`WALLET_MINIMUM_BALANCE`, `STANDARD_WALLET_MINIMUM_BALANCE`,
`WALLET_SIGNER_ACTIVATION_AMOUNT`, `SUB_WALLET_ACTIVATION_AMOUNT`,
`BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT`,
`MM_SUB_WALLET_ACTIVATION_AMOUNT`, `ISSUING_SUB_WALLET_ACTIVATION_AMOUNT`,
`MIN_SENDABLE_AMOUNT`

- Example value (any of the above): `0.001`
- What they do: each gates a specific wallet-creation or payment flow's minimum/activation amount (see the variable name for which flow).
- How to get a real value: small positive decimals for testnet; sized appropriately for mainnet gas/dust economics in production.

**Sub-wallet seeds.** `ISSUING_SUB_WALLET_ACTIVATION_AMOUNT`,
`MM_SUB_WALLET_ACTIVATION_AMOUNT`, `BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT`
and the fallback `SUB_WALLET_ACTIVATION_AMOUNT` are rows of the
`activation_amounts` table (id = the name, `amount`, `inactive`), not
environment variables. The amount is **ETH** the primary wallet sends to
each new sub-wallet (and its linked distribution wallet) in the operation
that deploys them, so the sub-wallet can pay its first network fees.
Missing, inactive or `0` means no seed - the sub-wallet then pays gas in
the owner's gas-fee stablecoin once it holds some. Example: `0.0002`.

**Sub-wallet creation fee.** The `SUBWALLET_CREATION_FEE` row of the
`service_fees` table sets it: `fee_fixed` is the fee **in USD**, paid in
the stablecoin `fee_asset_code` / `fee_contract_address` (the dollar
asset, a USD-named stablecoin such as USDC at 1:1, or the naira asset at
the USD/cNGN rate - other assets need a DEX and are refused), and
`inactive = 1` turns it off. It is charged from the primary wallet in the
same operation, sent to `SUBWALLET_CREATION_FEE_WALLET`, and recorded in
`fee_collections` (type `SUBWALLET_CREATION`). The platform's tokenization
issuing profile (`TOKENIZATION_ISSUING_PROFILE`, default `atprofile`) is
exempt. Creation is refused when the primary
wallet cannot cover the fee, the seeds and (when gas is paid in ETH) the
network fee.

---

## Memo requirements

**`WALLETS_REQUIRE_16_BYTE_MEMO`** / **`WALLETS_REQUIRE_28_BYTE_MEMO`** / **`WALLETS_REQUIRE_VARIABLE_BYTE_MEMO`**
- Example: `0xExchangeWalletA,0xExchangeWalletB`
- What it does: Comma-separated lists of destination wallet addresses (typically third-party exchange deposit addresses) that require a memo of that specific byte length on payments sent to them — a legacy carryover from the pre-Base chain's memo requirements, kept for destinations that still need it.
- How to get a real value: leave empty unless you know a specific destination requires a memo of that length.

---

## Tokenization

**`ENABLE_ASSET_TOKENIZATION`**
- Example: `1`
- What it does: Feature flag gating the tokenized-asset endpoints.
- How to get a real value: `0` or `1` per your deployment.

**`TOKENIZATION_ISSUING_PROFILE`** / **`TOKENIZATION_ISSUING_PROFILE_WALLET`**
- Example: `atprofile` / `0x1234...abcd` (optional)
- What it does: The platform user whose sub-wallets are every tokenized asset's issuing and distribution wallets. Each pair is owned by this user's key and the asset's minting approvers (see [INTEGRATION.md](INTEGRATION.md#tokenized-assets-token-issuing-and-distribution-wallets-sale-offer)); the backend holds no key for them. `TOKENIZATION_ISSUING_PROFILE_WALLET`, if set, must be that user's primary wallet address - a configuration check; leave it empty otherwise. It is no longer a private key.
- How to get a real value: the username of the issuing profile you register.

**`OFFER_BOOK_ADDRESS`**
- Example: `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- What it does: The `TrovoOfferBook` contract tokenized assets are sold through. Required when `ENABLE_ASSET_TOKENIZATION=1` (startup fails otherwise).
- How to get a real value: printed by `market/contracts/scripts/deploy.js` and saved in `market/contracts/deployments/<chainId>.json` (see [market/DEPLOYMENT.md](../market/DEPLOYMENT.md)).

**`OFFER_AUTHORIZER_PRIVATE_KEY`**
- Example: `0x59c6995e...` (hex private key)
- What it does: Signs the fill authorization of each purchase the platform has checked (KYC, cap, sale window). Its address must be an authorizer on the offer book. A leaked key lets someone buy at the seller's price without those checks; it cannot move anyone's tokens. Required when tokenization is enabled.
- How to get a real value: a dedicated key from your secrets manager; give its address to the offer book (`AUTHORIZER_ADDRESSES` at deploy, or `setAuthorizer` by the owner Safe).

**`SAFE_PROXY_FACTORY_ADDRESS`** / **`SAFE_SINGLETON_ADDRESS`** / **`SAFE_MULTISEND_CALL_ONLY_ADDRESS`** (and `SAFE_FALLBACK_HANDLER_ADDRESS`)
- Example: leave empty for the defaults
- What it does: The Safe contracts wallets are created from (see "Wallets: Safe accounts, bundler and paymaster"); MultiSendCallOnly also batches the internal balance minting Safe's fiat deliveries. Defaults are Safe's canonical deployments on Base and Base Sepolia.
- How to get a real value: leave empty unless Safe's contracts live elsewhere on your network.

A tokenized asset's **token contract** is not configuration: it is registered per asset through `PUT /v1/trovo-manager/tokenization/contract/:tid` (see [INTEGRATION.md](INTEGRATION.md#tokenized-assets-token-contract-vs-issuing-safe)).

---

## Shutdown

**`SHUTDOWN_GRACE_PERIOD`**
- Example: `90s` (default `60s`)
- What it does: How long a stopping instance (SIGTERM/SIGINT) waits for in-flight requests, running background jobs and platform-key transactions before exiting. See DEPLOYMENT.md "Running two or more instances".
- How to get a real value: a little less than your platform's stop timeout (e.g. Kubernetes `terminationGracePeriodSeconds`).

---

## Account recovery

Opt-in account recovery (see [INTEGRATION.md](INTEGRATION.md#account-recovery-opt-in-guardian)). Without `RECOVERY_MODULE_ADDRESS` and a guardian, the recovery endpoints answer `error-account-recovery-not-configured` (503); nothing else is affected.

**`RECOVERY_MODULE_ADDRESS`**
- Example: `0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e`
- What it does: Candide's Social Recovery Module v0.2.0 that covered wallets enable. Its recovery period is fixed when it is deployed.
- How to get a real value: printed by `recovery/contracts/scripts/deploy.js` and saved in `recovery/contracts/deployments/<chainId>.json` (see [recovery/DEPLOYMENT.md](../recovery/DEPLOYMENT.md)).

**`RECOVERY_PERIOD_SECONDS`**
- Example: `604800` (7 days)
- What it does: Only for messages to users ("your wallets move to the new key after 7 days"); the module's own period is what counts. Set it to the value the module was deployed with.
- How to get a real value: `recoveryPeriodSeconds` in the module's deployment file (`recovery/contracts/deployments/<chainId>.json`).

**`ACCOUNT_RECOVERY_GUARDIAN_SAFE`** / **`ACCOUNT_RECOVERY_GUARDIAN_SIGNERS`**
- Example: `0x1234...abcd` / `0x59c6...;0x5de4...` (hex private keys, `;` or `,` separated)
- What it does: The platform's recovery guardian, made the only guardian of each covered wallet. With `ACCOUNT_RECOVERY_GUARDIAN_SAFE` (recommended) the guardian is that Safe and the backend executes its calls with enough of the `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` keys to meet its threshold (the first key pays gas); without it the guardian is the first key itself. The guardian can only start replacing a covered wallet's key (finalized after the recovery period, cancellable by the user) - it cannot move funds. Its keys pay the gas of starting and finalizing recoveries, so keep them funded with a little ETH.
- How to get a real value: create a Safe owned by dedicated keys from your secrets manager (e.g. 2-of-3 with two keys here and one held offline). Changing the guardian later does not move existing wallets to it: users would have to turn recovery off and on again.

---

## Faucets & referral rewards

**`GAS_FAUCET`** / **`GAS_FAUCET_MIN_BALANCE`**
- Example: private key / `0.5`
- What it does: A wallet that dispenses small amounts of native gas asset to new wallets, and the balance threshold below which it alerts (see `FAUCET_LOW_BALANCE_WEBHOOK`).
- How to get a real value: a funded testnet wallet for dev; a real operational wallet in production.

**`REWARD_FAUCET`** / **`REWARD_FAUCET_MIN_BALANCE`** / **`REWARD_FAUCET_GAS_MIN_BALANCE`** / **`REWARD_ASSET_CODE`** / **`REWARD_ASSET_CONTRACT_ADDRESS`**
- Example: private key / `10` / `0.5` / `BNR` / `0x0000...`
- What it does: A wallet (and its low-balance thresholds) that pays out referral/reward tokens, and which asset those rewards are denominated in.
- How to get a real value: a funded wallet for the reward asset you've deployed.

**`ENABLE_REFERRAL_REWARD`** / **`REFERRAL_REWARD_AMOUNT`**
- What it does: Present in `.env.example` but **not currently read anywhere in the Go source** — vestigial. Safe to leave unset.

---

## Internal balance authorization / compliance

**`INTERNAL_BALANCE_ISSUING_SIGNERS`**
- Example: `0xkey1,0xkey2`
- What it does: Owner keys of each country's internal balance minting Safe (`CountryConfig.internalTokenMinterSafe`, which owns the internal balance token `internalTokenIssuer`). For a fiat purchase of a tokenized asset whose payment is confirmed, enough of them sign one Safe transaction that mints the purchase amount to the Safe and buys from the asset's sale offer for the buyer; the first key pays its gas, so it must hold ETH. (`INTERNAL_BALANCE_AUTHORIZER_WALLET` is no longer used.)
- How to get a real value: the keys of that Safe's owners, from your secrets manager.

**`COMPLIANCE_ACCOUNT_ID`**
- Example: `acct_test_123`
- What it does: An account identifier passed to the 1Liquidity/OneLiquidity compliance/KYC verification API.
- How to get a real value: issued by 1Liquidity when your organisation is onboarded with them — see `ONELIQUIDITY_*` below for how to reach them.

---

## P2P escrow (Gnosis Safe)

**`P2P_ESCROW_WALLET_ADDRESS`**
- Example: `0x2222222222222222222222222222222222222`
- What it does: The Gnosis Safe multisig address that holds P2P escrow deposits — both the deposit destination and the address `execTransaction` calls are sent to for settlement.
- How to get a real value: deploy a Gnosis Safe (e.g. via [safe.global](https://safe.global/)) on your target network and use its address.

**`P2P_ESCROW_SIGNERS`**
- Example: `key1;key2;key3`
- What it does: A `;`-separated (legacy `,`-separated also accepted) list of at least 3 signer private keys/mnemonics for the escrow Safe. The first signer also broadcasts and pays gas for every settlement.
- How to get a real value: the private keys of the Safe's configured owners — for local/testnet dev, create a testnet Safe with test-only owner keys.

**`P2P_ESCROW_MULTISEND_ADDRESS`**
- Example: `` (leave empty to use the default)
- What it does: Overrides the canonical `MultiSendCallOnly` v1.3.0 contract address that batches same-token settlement legs — only needed if your network's Safe deployment uses a non-standard address.
- How to get a real value: leave empty unless Safe's docs tell you your network needs a different address.

---

## Crypto deposit minting & withdrawal service (1Liquidity)

**`ENABLE_CRYPTO_DEPOSIT_MINTING`** / **`CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY`**
- Example: `0` / `0x1111111111111111111111111111111111111`
- What it does: When `1`, enables a background routine that mints assets for confirmed crypto deposits; the public key identifies the wallet whose signature authorizes minting (must be a 42-char address when enabled).
- How to get a real value: `0` unless you're specifically testing this flow.

**`ENABLE_CRYPTO_WITHDRAWAL_SERVICE`**
- Example: `0`
- What it does: When `1`, enables background polling of 1Liquidity for supported withdrawal networks.
- How to get a real value: `0` unless integrating with 1Liquidity.

**`ONELIQUIDITY_BASE_URL`** / **`ONELIQUIDITY_TOKEN`** / **`ONELIQUIDITY_WITHDRAWAL_CURRENCY_LIST`**
- Example: `https://api.oneliquidity.io` / `sk_live_...` / `USDT,USDC`
- What it does: Base URL, bearer token, and comma-separated currency list for 1Liquidity's crypto on/off-ramp and compliance API (`internal/components/users/services/oneliquidity.go`).
- How to get a real value: sign up as a partner with [1Liquidity](https://www.oneliquidity.io/) — they issue the base URL and API token during integration onboarding.

---

## Stablerail (fiat onramp for cNGN)

**`CNGN_PRICE_API_URL`**
- Example: `https://api.example.com/cngn/price`
- What it does: URL the app polls for the current cNGN (Nigerian Naira stablecoin) price. Required to boot.
- How to get a real value: your cNGN price-feed provider's endpoint.

Note: Stablerail's own API key/config (bank onboarding, onramp) is stored
in the database (`StablerailConfig` table), not read from environment
variables — set it up via whatever admin tooling manages that table (see
`tm-api`), not `.env`.

---

## Email (Mailgun)

**`MAILGUN_PRIVATE_API_KEY`**
- Example: `<your-mailgun-private-api-key>` (Mailgun keys look like `key-` followed by 32 hex characters — not shown here in that shape, since that pattern trips GitHub's secret scanner even as an obvious placeholder)
- What it does: Authenticates all outbound email (verification codes, notifications) sent via Mailgun. Required to boot.
- How to get a real value: sign up at [mailgun.com](https://www.mailgun.com/), verify a sending domain, and copy the private API key from Settings → API Keys.

**`MAILGUN_DOMAIN`**
- Example: `mg.yourcompany.com`
- What it does: The Mailgun sending domain emails are sent from. Required to boot.
- How to get a real value: the domain you verified in your Mailgun account.

**`MAILGUN_VALIDATOR_API_KEY`**
- Example: `<your-mailgun-validator-api-key>` (same key shape/caveat as `MAILGUN_PRIVATE_API_KEY` above)
- What it does: Used to validate email addresses via Mailgun's validation API. Required only when `ENABLE_EMAIL_VALIDATION=1`.
- How to get a real value: same Mailgun dashboard as above (may require a separate validation add-on depending on your plan).

**`ENABLE_EMAIL_VALIDATION`**
- Example: `0`
- What it does: When `1`, validates email addresses via Mailgun before accepting them (requires `MAILGUN_VALIDATOR_API_KEY`). Required to boot (even if `0`).
- How to get a real value: `0` for local dev unless you're testing this specifically.

**`SHOW_MAIL_VALIDATION_RESULT`**
- Example: `0`
- What it does: When `1`, includes the raw mail-validation result in the API response (debugging aid).
- How to get a real value: `0` normally.

**`ENABLE_EMAIL_NOTIFICATIONS`**
- Example: `1`
- What it does: Master toggle for whether transactional email notifications are sent at all.
- How to get a real value: `1` once Mailgun is configured; `0` to silence emails during local testing.

**`MAIL_SENDER`** / **`DEFAULT_MAIL_SENDER`**
- Example: `Trovotech <no-reply@mg.yourcompany.com>`
- What it does: The "From" address used for verification emails (`MAIL_SENDER`) and account-deletion emails (`DEFAULT_MAIL_SENDER`) respectively. Both fall back to a sensible default derived from `MAILGUN_DOMAIN` if unset.
- How to get a real value: any address on your verified Mailgun domain.

**`EMAIL_VERIFICATION_SUBJECT`** / **`EMAIL_VERIFICATION_TEMPLATE`**
- Example: `Your TrovoApp Verification Code`
- What it does: Subject line (and an unused/legacy template name) for verification emails.
- How to get a real value: any text; has a built-in default if unset.

**`ACCOUNT_DELETION_REQUEST_TEMPLATE`** / **`ACCOUNT_DELETION_EMAIL_SUBJECT`** / **`ACCOUNT_DELETION_DAYS`**
- Example: `account-deletion-request-template` / `TrovoApp Account Deletion Request` / `30`
- What it does: Template name, subject, and grace-period (in days) for the account-deletion email flow. Each has a default if unset (see the log lines in `main.go`).
- How to get a real value: defaults are fine unless you need custom copy/timing.

**`SUPPORT_EMAIL`**
- Example: `support@trovo.app`
- What it does: Support contact address surfaced in some emails/responses.
- How to get a real value: your team's real support inbox.

---

## SMS

**`DEFAULT_SMS_PROVIDER`**
- Example: `termii`
- What it does: Chooses which SMS provider (`termii` or `infobip`) is used by default when a phone-number-specific provider mapping isn't found in the DB.
- How to get a real value: `termii` or `infobip`, whichever you've set up.

**`TERMII_SMS_API_KEY`** / **`TERMII_SMS_URL`** / **`TERMII_SMS_SENDER_ID`**
- Example: `TLxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` / `https://api.ng.termii.com/api/sms/send` / `TrovoWallet`
- What it does: Credentials and endpoint for sending SMS via [Termii](https://termii.com/).
- How to get a real value: sign up at [termii.com](https://termii.com/), and copy the API key from your dashboard's API settings; register a sender ID there too.

**`INFOBIP_SMS_API_KEY`** / **`INFOBIP_SMS_HOST`**
- Example: `abcdef0123456789abcdef0123456789abcdef01` / `xxxxxx.api.infobip.com`
- What it does: Credentials and base host for sending SMS via [Infobip](https://www.infobip.com/).
- How to get a real value: sign up at [infobip.com](https://www.infobip.com/), create an API key under Settings → API Keys; your account's base host is shown on the same page.

**`TROVOWALLET_SMS_FROM`**
- Example: `TrovoWallet`
- What it does: The "From" sender name shown on outgoing SMS. Defaults to `TrovoWallet` if unset.
- How to get a real value: any short alphanumeric sender name your SMS provider allows.

**`ENABLE_MOBILE_VERIFICATION`**
- Example: `0`
- What it does: Feature flag for SMS-based mobile number verification during registration.
- How to get a real value: `1` once SMS is configured and you want this flow active.

---

## Push notifications, storage & Google Cloud (Firebase)

**`GC`**
- Example: base64-encoded Firebase service-account JSON (a long string, e.g. `eyJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsIC4uLn0=`)
- What it does: A base64-encoded Google/Firebase service account JSON key, decoded at startup to authenticate both Firebase Cloud Messaging (push notifications) and Firebase/Google Cloud Storage. Required to boot.
- How to get a real value: in the [Firebase console](https://console.firebase.google.com/), go to Project Settings → Service Accounts → Generate new private key, download the JSON file, then base64-encode it: `base64 -i service-account.json | tr -d '\n'`.

**`GOOGLE_PROJECT_ID`**
- Example: `trovo-wallet-dev`
- What it does: The Google Cloud / Firebase project ID associated with the service account above. Required to boot.
- How to get a real value: shown on your Firebase project's Settings → General page ("Project ID").

**`STORAGE_BUCKET_NAME`**
- Example: `trovo-wallet-dev.appspot.com`
- What it does: The default Google Cloud Storage bucket used for general file uploads (e.g. profile pictures).
- How to get a real value: your Firebase project's default Storage bucket name (Firebase console → Storage).

**`STAKEHOLDER_DOCUMENTS_BUCKET_NAME`**
- Example: `trovo-stakeholder-docs-dev`
- What it does: A **separate, private** GCS bucket used only for stakeholder-portal document uploads/downloads. If unset, those specific endpoints return `503` rather than failing unpredictably.
- How to get a real value: create a dedicated private bucket in [Google Cloud Console](https://console.cloud.google.com/storage) and grant your `GC` service account access to it.

---

## Dynamic links / deep linking (Firebase Dynamic Links)

**`DYNAMIC_LINKS_DOMAIN_PREFIX`**
- Example: `trovoapp.page.link`
- What it does: The Firebase Dynamic Links domain used to build shortened referral/deep-link URLs. Required to boot.
- How to get a real value: set up in Firebase console → Dynamic Links, or your equivalent link-shortening domain.

**`DYNAMIC_LINKS_ANDROID_PACKAGE_NAME`** / **`DYNAMIC_LINKS_IOS_BUNDLE_ID`**
- Example: `app.trovo.wallet` / `app.trovo.wallet`
- What it does: The Android package name / iOS bundle ID that dynamic links should open in the mobile app; also served in the `.well-known` app-links files (see `internal/components/root`). Required to boot.
- How to get a real value: your actual published app's package name / bundle ID.

**`DYNAMIC_LINKS_FALLBACK_BASE_URL`**
- Example: `https://trovo.app`
- What it does: Where a dynamic link falls back to if the app isn't installed. Required to boot.
- How to get a real value: your marketing site or app-store landing page URL.

**`FBDL_SERVICE_API_KEY`**
- Example: `<your-firebase-web-api-key>` (Google/Firebase API keys share a recognizable `AIza...` shape that GitHub's secret scanner flags even as an obvious placeholder, so it isn't shown in that shape here)
- What it does: API key for the Firebase Dynamic Links REST API used to generate short links.
- How to get a real value: Firebase console → Project Settings → General → Web API Key (Dynamic Links uses your project's Web API key).

**`SHORT_LINKS_BASE_URL`**
- Example: `https://trovo.app`
- What it does: Base URL used when generating app short links outside of Firebase Dynamic Links. Defaults to `https://trovo.app` if unset.
- How to get a real value: your own domain, if different from the default.

**`WALLET_DOMAIN`**
- Example: `trovo.app`
- What it does: The domain suffix appended to wallet aliases (e.g. `username@trovo.app`) when talking to 1Liquidity and elsewhere. Required to boot.
- How to get a real value: your platform's canonical domain.

---

## JWT / internal admin auth

**`JWT_ACCESS_SECRET`**
- Example: output of `openssl rand -hex 32`
- What it does: Signing secret for JWTs validated by `JwtTokenAuthMiddleware`, used on the `/v1/trovo-manager/...` routes that `tm-api` calls into. Required to boot.
- How to get a real value: `openssl rand -hex 32` — generate once per environment and keep it identical wherever tokens need to be verified (i.e. shared with whatever issues these JWTs).

**`JWT_TOKEN_EXPIRY`** / **`JWT_REFRESH_TOKEN_EXPIRY`**
- Example: `3600` / `604800`
- What it does: Access-token and refresh-token lifetimes, in seconds. Required to boot.
- How to get a real value: pick expiries appropriate for your security posture (e.g. 1 hour / 7 days).

---

## IP geolocation

**`IPAPI_HOST`** / **`IPAPI_KEY`**
- Example: `api.ipapi.com` / `abcdef0123456789abcdef0123456789`
- What it does: Host and API key for an IP-geolocation lookup provider (e.g. used for announcement targeting/fraud signals). Required to boot.
- How to get a real value: sign up with an IP geolocation provider such as [ipapi.com](https://ipapi.com/) and copy the API key from their dashboard.

---

## Feature flags

**`ENABLE_PATRON`**
- Example: `0`
- What it does: When `1`, runs the background routine that activates pending patron/subscription memberships.
- How to get a real value: `1` if your deployment uses the patron feature.

**`ENABLE_WEBSOCKET_AUTH`**
- Example: `1`
- What it does: Whether the WebSocket stream endpoints (`/v1/stream/...`) enforce authentication.
- How to get a real value: `1` in any real deployment.

**`ENABLE_OLD_USER_MIGRATION`**
- What it does: Present in `.env.example` but **not currently read anywhere in the Go source** — vestigial. Safe to leave unset.

---

## Registration & abuse controls

**`REGISTRATION_THROTTLE_PER_IP`**
- Example: `5`
- What it does: Maximum registrations allowed per IP in the throttle window (`0`/unset disables throttling).
- How to get a real value: tune based on your abuse tolerance; `0` for local dev.

---

## Service links (white-label partner integrations)

**`SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY`** / **`SERVICE_LINK_LOGIN_REQUEST_VALIDITY`**
- Example: `300`
- What it does: How long (seconds) a pending service-link authorization/login request stays valid before expiring.
- How to get a real value: a few minutes' worth of seconds is typical (e.g. `300` = 5 minutes).

See [INTEGRATION.md](INTEGRATION.md#service-links) for how service-link API
keys themselves are provisioned (via the database, not an env var).

---

## Error/monitoring webhooks (Discord)

Each of these, when set, redirects a specific category of background-job
error alert to a Discord channel via an incoming webhook (some code paths
also fall back to a hardcoded Trovo-internal webhook if unset — for your
own deployment, always set these explicitly).

`CONNECTION_WARNING_WEBHOOK`, `EXPANSION_NETWORK_ERROR_WEBHOOK`,
`FAILED_PAYMENT_ERROR_WEBHOOK`, `FAUCET_LOW_BALANCE_WEBHOOK`,
`IMPORT_ERROR_WEBHOOK`, `REGISTRATION_ERROR_WEBHOOK`

- Example value (any of the above): `https://discord.com/api/webhooks/123456789012345678/AbCdEfGhIjKlMnOpQrStUvWxYz`
- What it does: alerts for that specific failure category (see the variable name) are posted to this Discord webhook.
- How to get a real value: in Discord, go to a channel's Settings → Integrations → Webhooks → New Webhook, and copy its URL.

---

## Test suite only

**`RICPK`** / **`RICSC`**
- Example: a testnet wallet's public key / private key
- What it does: Only read by `main_test.go` (the integration test suite that exercises a live running server) — not used by the application itself. They provide a funded test wallet's public key and secret key for the tests to transact with.
- How to get a real value: a funded testnet wallet, used only for running `main_test.go` against a live local server; never a production key.

---

## Present in `.env.example` but not currently read by the code

These appear as leftovers in `.env.example` from an earlier SMTP-based mail
setup (the app now uses Mailgun — see the Email section) or other retired
features. You can leave them unset without any effect:
`SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `MAIL_FROM`,
`FBDL_SERVICE_URLS`, `ENABLE_OLD_USER_MIGRATION`, `ENABLE_REFERRAL_REWARD`,
`REFERRAL_REWARD_AMOUNT`.

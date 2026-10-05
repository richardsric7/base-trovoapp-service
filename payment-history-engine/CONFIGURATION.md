# Configuration

Every parameter `payment-history-engine` reads. Each one says what it is,
why it is needed, whether you must set it, an example and how to get a real
value. Examples are never real secrets.

## How to set them

The engine reads **environment variables**:

- **A `.env` file** in the directory you run it from. Start from the
  template: `cp .env.example .env`, then edit it. It is loaded
  automatically at start.
- **Or the container / hosting platform's environment settings** (for
  example `docker run --env-file .env`). Use this for secrets in
  production.

If a required value is missing, the engine prints
`Required environment variable is missing <NAME>` and stops.

## The minimum to start

| Parameter | Why |
|---|---|
| [`DB_CONNECTION_STRING`](#db_connection_string) | where wallets are read and payment history is written (app-backend's database) |
| [`CDB_CONNECTION_STRING`](#cdb_connection_string) | the engine's own tracking database |
| [`BASE_RPC_URL`](#base_rpc_url) | the blockchain node it watches |
| [`NATIVE_ASSET_CODE`](#native_asset_code) | the label for ETH transfers |
| [`ENABLE_CACHING`](#enable_caching) | must be set (`0` to run without Redis) |

---

## Main database (shared with app-backend)

### `DB_CONNECTION_STRING`

- **What it is:** The address and login of app-backend's main database.
- **Why it's needed:** The engine reads `user_wallets` there to learn which
  wallets to watch, and writes every transfer it finds into
  `payment_history`, which the apps and Trovo Manager show. Without it the
  engine stops at start.
- **Required:** Yes.
- **Example:** `postgres://trovo:change-me@db.internal:5432/trovo?sslmode=require`
  (Postgres) or `./local-wallet.db` (SQLite, local only)
- **How to get it:** Use the same database as app-backend's
  [`DB_CONNECTION_STRING`](../app-backend/CONFIGURATION.md). Deploy
  app-backend first: it creates the `payment_history` table. For a local
  test database:
  ```bash
  docker run --name wallet-db -e POSTGRES_USER=wallet_user -e POSTGRES_PASSWORD=wallet_pass -e POSTGRES_DB=trovo_wallet -p 5432:5432 -d postgres:16
  ```
  then `postgres://wallet_user:wallet_pass@localhost:5432/trovo_wallet?sslmode=disable`.

### `DB_TYPE`

- **What it is:** Which kind of database `DB_CONNECTION_STRING` is:
  `postgres` or `sqlite`.
- **Why it's needed:** Production uses Postgres; `sqlite` lets you run
  locally against a file without a database server.
- **Required:** No, default `postgres`.
- **Example:** `postgres`
- **How to get it:** Leave unset in real deployments; `sqlite` only for
  local tests.

### `DB_MAX_OPEN_CONNECTIONS`

- **What it is:** The most connections the engine keeps open to the main
  database (Postgres only).
- **Why it's needed:** The database has a connection limit shared with
  app-backend, tm-api and payout-engine.
- **Required:** No, default `50`.
- **Example:** `20`
- **How to get it:** Ask your database administrator what share of the
  database's `max_connections` this service may use.

### `DB_MAX_IDLE_CONNECTIONS`

- **What it is:** How many unused connections are kept open, ready for
  reuse (Postgres only).
- **Why it's needed:** Reusing connections is faster than opening new ones.
- **Required:** No, default `50`.
- **Example:** `20`
- **How to get it:** Usually the same as `DB_MAX_OPEN_CONNECTIONS`.

## Tracking database ("RoachDB")

### `CDB_CONNECTION_STRING`

- **What it is:** The address of the engine's own tracking database (a
  CockroachDB or Postgres database, called "RoachDB" in the code).
- **Why it's needed:** It holds the list of tracked wallets and addresses
  (`tracked_wallets`, `tracked_addresses`, which app-backend also adds to
  when users create wallets) and how far the engine has scanned
  (`monitored_cursors`, `monitored_account_cursors`). Without it the engine
  stops at start.
- **Required:** Yes.
- **Example:** `postgresql://root@roach.internal:26257/payment_history?sslmode=require`
  or `./local-roach.db` (with `ROACH_DB_TYPE=sqlite`)
- **How to get it:** Use the same value as app-backend's
  `CDB_CONNECTION_STRING`. For a local CockroachDB:
  ```bash
  cockroach start-single-node --insecure --listen-addr=localhost:26257 --background
  cockroach sql --insecure -e "CREATE DATABASE payment_history;"
  ```
  then `postgresql://root@localhost:26257/payment_history?sslmode=disable`.
  A managed CockroachDB (CockroachDB Cloud) shows its connection string on
  the cluster's **Connect** page.

### `ROACH_DB_TYPE`

- **What it is:** Set to `sqlite` to use a local file as the tracking
  database instead of CockroachDB.
- **Why it's needed:** Lets you run locally without CockroachDB.
- **Required:** No. Unset means CockroachDB/Postgres.
- **Example:** `sqlite`
- **How to get it:** Only for local tests; leave unset in real deployments.

### `PURGE_TABLES`

- **What it is:** A destructive reset switch.
- **Why it's needed:** To wipe a broken local tracking database. With `1`,
  the engine deletes `tracked_wallets`, `tracked_addresses` and
  `monitored_cursors` from the tracking database and exits without doing
  anything else. There is no confirmation.
- **Required:** No. Leave unset.
- **Example:** `1`
- **How to get it:** Never set it in production.

## Blockchain

### `BASE_RPC_URL`

- **What it is:** The web address of a Base blockchain node.
- **Why it's needed:** The engine reads every new block and the token
  transfer logs from it. The `/ready` check also uses it.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key` (or
  `https://sepolia.base.org` for the test network)
- **How to get it:** For testing, `https://sepolia.base.org` needs no
  sign-up. For production, create an app for Base Mainnet with a provider
  such as Alchemy, Infura or QuickNode and copy its HTTPS URL. The engine
  makes many log queries, so a paid plan is advisable.

### `BASE_CHAIN_ID`

- **What it is:** The network number.
- **Why it's needed:** Used to work out the sender of each transaction.
  It must match the network of `BASE_RPC_URL`.
- **Required:** No, default `84532` (Base Sepolia). **Set `8453` on
  mainnet.**
- **Example:** `8453`
- **How to get it:** `8453` for Base Mainnet, `84532` for Base Sepolia.

### `ENTRYPOINT_ADDRESS`

- **What it is:** The address of the ERC-4337 EntryPoint that Trovo
  wallets' operations go through.
- **Why it's needed:** Trovo wallets send ETH inside these operations, which
  a simple block scan does not see; the engine decodes them to record those
  transfers.
- **Required:** No, default `0x0000000071727De22E5E9d8BAf0edAc6f37da032`
  (the official v0.7 EntryPoint on Base and Base Sepolia).
- **Example:** `0x0000000071727De22E5E9d8BAf0edAc6f37da032`
- **How to get it:** Keep the default; it must equal app-backend's
  `ENTRYPOINT_ADDRESS`. On a local test chain, use the address the paymaster
  local stack prints.

### `NATIVE_ASSET_CODE`

- **What it is:** The asset code recorded for ETH transfers.
- **Why it's needed:** Every payment history row has an asset code; ETH has
  no token contract to read one from, so this label is used.
- **Required:** Yes.
- **Example:** `ETH`
- **How to get it:** Use the same code app-backend uses for the native
  asset, so the apps show it consistently (`ETH` on Base).

### `RPC_TIMEOUT`

- **What it is:** The longest one request to the node may take.
- **Why it's needed:** A stuck node must not freeze the engine. After 3
  failures in a row, the engine pauses node calls for 10 seconds at a time
  until one succeeds. A block that could not be read is retried, never
  skipped.
- **Required:** No, default `30s`.
- **Example:** `30s`
- **How to get it:** Keep the default; raise it if your provider is slow
  for large log queries.

### `TRACK_ADDRESS_CONCURRENCY`

- **What it is:** How many tracked addresses have their history filled in
  at the same time.
- **Why it's needed:** Each fill-in uses node requests and memory; this
  caps both.
- **Required:** No, default `8`.
- **Example:** `8`
- **How to get it:** Keep the default; raise it only if your node plan
  allows more requests at once.

### `LAST_CURSOR`

- **What it is:** A block number to resume the live scan from.
- **Why it's needed:** Normally the engine resumes from where it stopped
  (saved in the tracking database). This forces a different start, for
  example after a manual fix. The more recent of this and the saved cursor
  wins.
- **Required:** No.
- **Example:** `21000000`
- **How to get it:** The block number from
  [basescan.org](https://basescan.org); leave unset normally.

### `BLOCKCHAIN_NETWORK_PASSPHRASE`, `BLOCKCHAIN_BASE_RESERVE`, `BLOCKCHAIN_SWAP_DESTINATION_MIN`

- **What they are:** Leftovers from the platform's previous blockchain
  (Stellar). They have no meaning on Base.
- **Why they're needed:** They aren't. The code reads them but works
  without them.
- **Required:** No.
- **Example:** (leave unset)
- **How to get them:** Don't set them.

## Caching (Redis)

### `ENABLE_CACHING`

- **What it is:** Turns Redis caching on (`1`) or off (`0`).
- **Why it's needed:** It must be set to something: the engine refuses to
  start if it is empty. With `1`, the Redis settings below become required
  and the engine stops if Redis does not answer.
- **Required:** Yes (`0` or `1`).
- **Example:** `0`
- **How to get it:** `0` unless you run Redis for this service.

### `REDIS_HOST`

- **What it is:** The Redis server's host name.
- **Why it's needed:** Where the cache lives.
- **Required:** Only with `ENABLE_CACHING=1`.
- **Example:** `redis.internal`
- **How to get it:** Your Redis server (the same one app-backend uses is
  fine). Locally: `docker run --name redis -p 6379:6379 -d redis:7`, then
  `localhost`.

### `REDIS_PORT`

- **What it is:** The Redis server's port.
- **Why it's needed:** To connect to Redis.
- **Required:** Only with `ENABLE_CACHING=1`.
- **Example:** `6379`
- **How to get it:** `6379` unless your Redis says otherwise.

### `REDIS_PASSWORD`

- **What it is:** The Redis password.
- **Why it's needed:** A protected Redis refuses connections without it.
- **Required:** No (only if your Redis has a password).
- **Example:** `change-me`
- **How to get it:** From whoever runs your Redis.

### `CACHING_PARAMETER`

- **What it is:** A short text added to every cache key.
- **Why it's needed:** Changing it makes all old cache entries unused at
  once (for example after the cached data's shape changes).
- **Required:** Only with `ENABLE_CACHING=1`.
- **Example:** `v1`
- **How to get it:** Any short text; change it to clear the cache.

## Health server

### `HEALTH_PORT`

- **What it is:** The port the health endpoints (`/health`, `/ready`) and
  the API reference (`/swagger/index.html`) are served on.
- **Why it's needed:** Your hosting platform checks `/health` and `/ready`
  to know the engine is alive and working.
- **Required:** No, default `8080`.
- **Example:** `8080`
- **How to get it:** Keep the default unless the port is taken.

### `APP_VERSION`

- **What it is:** The version reported by `/health` and `/ready`.
- **Why it's needed:** Tells you which build is running. Normally stamped
  into the image when it is built (the Dockerfile's `APP_VERSION` build
  argument); this variable is only used if that stamp is empty.
- **Required:** No.
- **Example:** `a1b2c3d`
- **How to get it:** `git rev-parse --short HEAD`.

## Alerts (Discord)

The engine posts some errors to Discord. Each has a built-in fallback
webhook in the code; set your own so the alerts reach your team.

### `CONNECTION_WARNING_WEBHOOK`

- **What it is:** A Discord webhook for "database connection pool is full"
  warnings.
- **Why it's needed:** A full pool means the engine is falling behind; your
  team should know.
- **Required:** No (a built-in webhook is used otherwise).
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** In Discord: the channel's **Edit Channel →
  Integrations → Webhooks → New Webhook → Copy Webhook URL**.

### `EXPANSION_NETWORK_ERROR_WEBHOOK`

- **What it is:** A Discord webhook for errors reaching the blockchain node.
- **Why it's needed:** So your team notices when the node is down.
- **Required:** No. It must be longer than 50 characters to replace the
  built-in one.
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** As above.

### `FAILED_PAYMENT_ERROR_WEBHOOK`

- **What it is:** A Discord webhook for errors recording a payment.
- **Why it's needed:** A failed record means a user's history is missing a
  transfer.
- **Required:** No. It must be longer than 50 characters to replace the
  built-in one.
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** As above.

## Optional and unused

These are read by code paths the engine does not currently use. Leave them
unset unless that changes.

### `ENABLE_EMAIL_VALIDATION`

- **What it is:** Turns on Mailgun email validation (`1`).
- **Why it's needed:** Not used by the engine's work; with `1`,
  `MAILGUN_VALIDATOR_API_KEY` becomes required at start.
- **Required:** No.
- **Example:** `0`
- **How to get it:** Leave unset.

### `MAILGUN_VALIDATOR_API_KEY`

- **What it is:** Mailgun's email validation API key.
- **Why it's needed:** Only with `ENABLE_EMAIL_VALIDATION=1`.
- **Required:** Only with `ENABLE_EMAIL_VALIDATION=1`.
- **Example:** `key-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Mailgun dashboard → **API Keys**.

### `MIN_SENDABLE_AMOUNT`

- **What it is:** A minimum amount quoted in one error message.
- **Why it's needed:** Not reachable through this engine.
- **Required:** No.
- **Example:** `1`
- **How to get it:** Leave unset.

### `IPAPI_HOST`, `IPAPI_KEY`

- **What they are:** An IP-geolocation service's address and key.
- **Why they're needed:** Read by a helper the engine never calls.
- **Required:** No.
- **Example:** `ip-api.com`, `abc123`
- **How to get them:** Leave unset.

### `ENABLE_AUTH_MIDDLEWARE`

- **What it is:** `0` turns off request-signature checks in an
  authentication middleware.
- **Why it's needed:** The middleware is not attached to any route (the
  engine only serves health endpoints), so this has no effect.
- **Required:** No.
- **Example:** (leave unset)
- **How to get it:** Leave unset.

## Used only by tests

`AA_LOCAL_STACK` (the local blockchain deployment file from the paymaster
project) and `AA_SAFE_ETH_TX_FILE` point the end-to-end tests at a local
chain; see [DEPLOYMENT.md](DEPLOYMENT.md#run-the-tests). The engine itself
never reads them.

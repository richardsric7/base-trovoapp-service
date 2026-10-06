# Deployment

A step-by-step guide to getting `app-backend` running, written for someone
doing it for the first time. Commands run from the `app-backend/` folder
unless a step says otherwise.

## 1. What you are deploying

`app-backend` is the **Trovo Wallet API**: the server behind the Trovo
mobile app (`app-mobile`) and web wallet (`app-web`). It holds user
accounts, creates each user's wallet (a Safe smart account on the Base
blockchain), sends payments and swaps, runs P2P trading, tokenized assets,
Public Markets, KYC and account recovery, and gives partner businesses an
API ("service links"). Trovo Manager (`tm-api`) shares its database.

**It needs these to be running first:**

| What | Why | Where |
|---|---|---|
| PostgreSQL (main database) | accounts, wallets, history and everything else | your Postgres; step 3 starts one locally |
| A tracking database (Postgres or CockroachDB) | the list of wallets payment-history-engine watches; app-backend adds new wallets to it | the same server can hold both databases |
| Redis | caching, rate limits, websocket fan-out between copies | step 3 starts one locally |
| A Base RPC node | reads the blockchain and sends transactions | a provider such as Alchemy, or `https://sepolia.base.org` for testing |
| A bundler, the paymaster contract and its quote service | every wallet transaction goes through them | [paymaster/DEPLOYMENT.md](../paymaster/DEPLOYMENT.md) |
| Mailgun and Firebase | emails, push notifications and file storage | accounts with each provider ([CONFIGURATION.md](CONFIGURATION.md)) |
| Optional contracts | the offer book (tokenization, swaps) and the recovery module | [market/DEPLOYMENT.md](../market/DEPLOYMENT.md), [recovery/DEPLOYMENT.md](../recovery/DEPLOYMENT.md) |

After app-backend, deploy payment-history-engine, payout-engine, tm-api and
the apps (see the root [README.md](../README.md) for the full order).

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Go 1.26 or newer | <https://go.dev/doc/install> | `go version` |
| Docker | <https://docs.docker.com/get-docker/> | `docker --version` |
| `psql` (PostgreSQL client) | <https://www.postgresql.org/download/> | `psql --version` |
| `openssl` (to make secrets) | usually preinstalled | `openssl version` |
| Foundry's `cast` (to make keys and phrases) | <https://book.getfoundry.sh/getting-started/installation> | `cast --version` |
| `make` (optional) | usually preinstalled | `make --version` |

## 3. Start a database and Redis (local only)

On a server, use your managed Postgres and Redis instead and skip this
step. Locally, Docker gives you both in two commands:

```bash
docker run -d --name trovo-postgres -p 5432:5432 \
  -e POSTGRES_USER=trovo -e POSTGRES_PASSWORD=change-me -e POSTGRES_DB=trovo postgres:16
docker run -d --name trovo-redis -p 6379:6379 redis:7
```

Then create the tracking database next to the main one:

```bash
psql "postgresql://trovo:change-me@localhost:5432/trovo" -c "CREATE DATABASE trovo_tracking;"
```

**You should see:** `CREATE DATABASE`.

## 4. Get the code and build it

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/app-backend
go mod download
go build -o main .
```

**You should see:** no output from `go build`, and a `main` file in the
folder.

## 5. Set the parameters

```bash
cp .env.example .env
```

`.env.example` lists the **required** settings first (app-backend refuses
to start without them), then the recommended and optional ones. Each is
explained, with an example and where to get it, in
[CONFIGURATION.md](CONFIGURATION.md). The ones you must create yourself:

| Parameter | How to make it |
|---|---|
| [`VERIFICATION_CODE_SALT`](CONFIGURATION.md#verification_code_salt), [`ENCODER_SALT`](CONFIGURATION.md#encoder_salt), [`MARKET_MAKING_SALT`](CONFIGURATION.md#market_making_salt), [`BULK_PAYMENT_SALT`](CONFIGURATION.md#bulk_payment_salt) | `openssl rand -hex 16` (a different one each) |
| [`MNEMONIC_MARKET_MAKING`](CONFIGURATION.md#mnemonic_market_making), [`MNEMONIC_BULK_PAYMENT`](CONFIGURATION.md#mnemonic_bulk_payment) | `cast wallet new-mnemonic` (a different one each) |
| [`CHANNEL_ACCOUNTS`](CONFIGURATION.md#channel_accounts), [`CHANNEL_ACCOUNT_FUNDER`](CONFIGURATION.md#channel_account_funder) | `cast wallet new` twice; use the private keys. They hold no money. |
| [`JWT_ACCESS_SECRET`](CONFIGURATION.md#jwt_access_secret) | `openssl rand -hex 32` |
| the ten [fee wallets](CONFIGURATION.md#fee-wallets) | addresses of wallets your finance team controls (for testing: `cast wallet new` and use the addresses) |

**Write the salts and phrases into your secrets manager now.** Changing
`ENCODER_SALT`, the two phrases or their salts after users sign up breaks
those users' security answers and wallets.

The rest come from your providers: the database and Redis addresses (step
3), [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url), Mailgun
([`MAILGUN_PRIVATE_API_KEY`](CONFIGURATION.md#mailgun_private_api_key)),
Firebase ([`GC`](CONFIGURATION.md#gc),
[`GOOGLE_PROJECT_ID`](CONFIGURATION.md#google_project_id)), and the paymaster stack
([`BUNDLER_URL`](CONFIGURATION.md#bundler_url),
[`PAYMASTER_ADDRESS`](CONFIGURATION.md#paymaster_address),
[`PAYMASTER_QUOTE_SERVICE_URL`](CONFIGURATION.md#paymaster_quote_service_url),
[`PAYMASTER_QUOTE_SERVICE_API_KEY`](CONFIGURATION.md#paymaster_quote_service_api_key)).

For a first local run without every provider: `GC` only has to be base64
(any value works until a push notification is sent), the Mailgun key is
used only when an email is sent, and
[`IPAPI_HOST`](CONFIGURATION.md#ipapi_host) /
[`IPAPI_KEY`](CONFIGURATION.md#ipapi_key) only have to be non-empty
(nothing calls the location lookup today).

Never commit `.env`. On a server, set the same values in the hosting
platform's environment settings or with `--env-file`.

## 6. Create the tables (once per release)

app-backend creates and updates its tables itself. Run the update as its
own step, from a machine close to the database:

```bash
MIGRATE_ONLY=1 ./main
```

**You should see**, after a while, `MIGRATE_ONLY=1 set — migrations
complete, exiting without starting the server`. It only needs the database
settings. It checks about 110 tables one by one, so over a slow or distant
connection it can take many minutes; run it on the server or in the same
region as the database.

Locally you can skip this: with `DB_AUTOMIGRATE=1` (the template's value)
the server updates the tables every time it starts.

## 7. Run it

**A. Directly:**

```bash
./main          # or: go run main.go / make run
```

**B. With Docker (recommended for servers):**

```bash
docker build -t trovo-wallet-api .
docker run --rm --env-file .env -e MIGRATE_ONLY=1 trovo-wallet-api        # step 6, in Docker
docker run -d --name trovo-wallet-api --restart unless-stopped \
  -p 8080:8080 --env-file .env -e DB_AUTOMIGRATE=0 trovo-wallet-api
```

Inside Docker, `localhost` means the container itself: use the database's
and Redis's real host names (or `host.docker.internal` with Docker
Desktop).

**At start app-backend:** loads `.env` if there is one, opens both
databases, checks the required settings (printing `Required environment
variable is missing <NAME>` for each one missing, then stopping), updates
the tables unless `DB_AUTOMIGRATE=0`, connects to Redis, starts its
background jobs, and listens on `PORT` (default `8080`). The last line is
`##service started##`.

## 8. Check it works

```bash
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/swagger/index.html
```

**You should see:** `200`. The API reference is at
<http://localhost:8080/swagger/index.html>. app-backend has no separate
`/health` route: point your hosting platform's health check at
`/swagger/index.html` or at the port itself (a TCP check).

Then:

1. Start the web wallet ([app-web/DEPLOYMENT.md](../app-web/DEPLOYMENT.md))
   or the mobile app with its API address set to this server, and sign up.
2. In `psql`, the new user is in `users` and their wallet in
   `user_wallets`; the tracking database's `tracked_wallets` has it too.
3. Send a small test payment between two test users: it needs the bundler
   and paymaster from step 1.

## 9. Running in production, updating and rolling back

- **Release order:** run the `MIGRATE_ONLY=1` step, then replace the
  serving containers (which run with `DB_AUTOMIGRATE=0`). The old version
  keeps serving while the tables update; a failed update exits with an
  error and leaves it serving.
- **More than one copy** is safe behind a load balancer. Turn on Redis
  (`ENABLE_CACHING=1`) so rate limits, caches and websocket messages are
  shared. Everything else is handled for you:
  - copies starting together take turns updating the tables (a lock row in
    `distributed_locks`);
  - each background job (sale start, order expiry, offer book index,
    callbacks to partners, crypto deposit minting, and others) runs on one
    copy at a time;
  - each platform signing key and platform Safe is used by one copy at a
    time, so two copies never send transactions with the same nonce
    (signatures made now and sent later reserve a nonce in
    `nonce_reservations`);
  - callbacks to partners are saved in `callback_deliveries` and retried
    (up to 10 times), so a stopped copy loses none.
- **Stopping** (`docker stop`, or an autoscaler removing a copy) is
  graceful: app-backend stops taking new work and waits up to
  [`SHUTDOWN_GRACE_PERIOD`](CONFIGURATION.md#shutdown_grace_period)
  (default 60 seconds) for running requests and transactions. Give your
  platform's stop timeout at least that long (`docker stop -t 90`).
- **Updating:**
  ```bash
  git pull
  docker build -t trovo-wallet-api .
  docker run --rm --env-file .env -e MIGRATE_ONLY=1 trovo-wallet-api
  docker rm -f trovo-wallet-api
  docker run -d --name trovo-wallet-api --restart unless-stopped \
    -p 8080:8080 --env-file .env -e DB_AUTOMIGRATE=0 trovo-wallet-api
  ```
  When a release changes tables tm-api or payment-history-engine also use,
  update app-backend first.
- **Rolling back:** `git checkout <previous-tag>`, rebuild and replace the
  container the same way. Table updates add tables and columns and never
  remove them, so an older version normally still runs against the newer
  tables; check the release notes for renamed columns first.
- There is no CI/CD pipeline in this repository: run `make ci` (tidy check,
  build, vet, lint, tests) before merging, and build and deploy the image
  by hand.

## 10. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `Required environment variable is missing X` and the program stops | `X` is empty | set it ([CONFIGURATION.md](CONFIGURATION.md) explains each) |
| `[main]Error opening DB` | wrong `DB_CONNECTION_STRING`, or the database is unreachable | test with `psql` from the same machine; inside Docker, don't use `localhost` |
| a crash at start mentioning `CHANNEL_ACCOUNT_FUNDER` or `CHANNEL_ACCOUNT_MIN_COUNT` | not a private key, or not a number | a key from `cast wallet new`; `0` |
| `[main]Error connecting to redis` and the program stops | `ENABLE_CACHING=1` but `REDIS_HOST`/`REDIS_PORT`/`REDIS_PASSWORD` are wrong or Redis is down | check Redis (`redis-cli -h <host> ping` answers `PONG`) |
| the table update takes 20 minutes or more | it runs far from the database | run `MIGRATE_ONLY=1` in the database's region; serve with `DB_AUTOMIGRATE=0` |
| `insufficient arguments` during the table update | an old Postgres driver (fixed in `v1.5.11`) | update to this repository's `go.mod` |
| every send fails | `BUNDLER_URL` unset, or the bundler, paymaster or quote service is down | check them ([paymaster/DEPLOYMENT.md](../paymaster/DEPLOYMENT.md)) |
| users can't log in after a redeploy | `JWT_ACCESS_SECRET` changed (logins are re-required), or the copies have different values | use one value on every copy |
| security answers or market-making wallets stop working | `ENCODER_SALT` or a phrase or salt changed | restore the original value from your secrets manager |
| new wallets have no payment history | payment-history-engine is not running, or its `CDB_CONNECTION_STRING` points elsewhere | use the same tracking database in both |
| `429 Too Many Requests` | a rate limit | raise it ([CONFIGURATION.md](CONFIGURATION.md#4-rate-limiting)) or the partner's own limit in Trovo Manager |

## Running the tests

```bash
make test       # go test ./internal/... -race
make ci         # tidy check, build, vet, lint, tests
```

Tests against a local blockchain (the `internal/aa`, `internal/gnosissafe`
and `internal/network` end-to-end tests) need the paymaster local stack
running and `AA_LOCAL_CHAIN=1` plus the `AA_LOCAL_*` settings
([CONFIGURATION.md](CONFIGURATION.md#25-used-only-by-tests),
[paymaster/DEPLOYMENT.md](../paymaster/DEPLOYMENT.md#local-development-stack)).
Run them one package at a time (`go test -p 1 ./internal/...`). To run the
Public Markets tests against Postgres instead of SQLite, set
`PUBLIC_MARKETS_TEST_POSTGRES`.

## Swagger (API reference)

If you change a route or its annotations, regenerate the reference and
commit `docs/`:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init --parseDependency=false
```

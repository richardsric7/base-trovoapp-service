# Deployment

A step-by-step guide to getting `payment-history-engine` running, written
for someone doing it for the first time.

## 1. What you are deploying

A background program that watches the Base blockchain for transfers to and
from Trovo users' wallets and writes each one into the payment history that
the Trovo apps and Trovo Manager show. It has no user-facing API; it only
serves two health-check pages (`/health`, `/ready`) for your hosting
platform.

**It needs these to be running first:**

| What | Why | Where |
|---|---|---|
| app-backend and its Postgres database | the engine reads `user_wallets` and writes `payment_history` there; app-backend creates those tables | [app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md) |
| A CockroachDB (or Postgres) database, "RoachDB" | the engine's own list of tracked wallets and scan positions; app-backend writes to it too | CockroachDB Cloud, or self-hosted |
| A Base node URL | to read blocks and transfers | a provider such as Alchemy or QuickNode |
| Redis (optional) | caching | the same Redis app-backend uses |

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Go 1.26 or newer (to build from source) | <https://go.dev/doc/install> | `go version` |
| Docker (recommended for servers) | <https://docs.docker.com/get-docker/> | `docker --version` |
| `psql` (to test database connections) | part of PostgreSQL: <https://www.postgresql.org/download/> | `psql --version` |

For local testing without database servers, the engine can use SQLite files
for both databases (see step 4).

## 3. Get the code and build it

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/payment-history-engine
go mod download
go build -o payment-history-engine .
go vet ./...
```

**You should see:** no output from `go build` and `go vet` (any output is an
error).

## 4. Set the parameters

```bash
cp .env.example .env
```

Open `.env` and fill in at least these (each is explained, with an example
and where to get it, in [CONFIGURATION.md](CONFIGURATION.md)):

| Parameter | Example |
|---|---|
| [`DB_CONNECTION_STRING`](CONFIGURATION.md#db_connection_string) | app-backend's database, e.g. `postgres://trovo:change-me@db.internal:5432/trovo?sslmode=require` |
| [`CDB_CONNECTION_STRING`](CONFIGURATION.md#cdb_connection_string) | `postgresql://root@roach.internal:26257/payment_history?sslmode=require` |
| [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url) | `https://base-mainnet.g.alchemy.com/v2/your-api-key` |
| [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id) | `8453` (mainnet; the default is the test network) |
| [`NATIVE_ASSET_CODE`](CONFIGURATION.md#native_asset_code) | `ETH` |
| [`ENABLE_CACHING`](CONFIGURATION.md#enable_caching) | `0` |

Recommended: your own Discord webhooks for the
[alerts](CONFIGURATION.md#alerts-discord).

**Local test without database servers:** set `DB_TYPE=sqlite`,
`DB_CONNECTION_STRING=./local-wallet.db`, `ROACH_DB_TYPE=sqlite` and
`CDB_CONNECTION_STRING=./local-roach.db`. (The `payment_history` and
`user_wallets` tables still have to exist, so point
`DB_CONNECTION_STRING` at a database app-backend has started against.)

Never commit `.env`: it holds secrets.

## 5. Run it

**A. Directly:**

```bash
./payment-history-engine
```

**B. With Docker (recommended for servers):**

```bash
docker build -t payment-history-engine --build-arg APP_VERSION=$(git rev-parse --short HEAD) .
docker run -d --name payment-history-engine --restart unless-stopped -p 8080:8080 --env-file .env payment-history-engine
```

`APP_VERSION` is a **build argument** stamped into the image and shown by
`/health`:

- **What it is:** the version of the code you built.
- **Why it's needed:** to tell which build is running.
- **Required:** No (default `unknown`).
- **Example:** `a1b2c3d`
- **How to get it:** `git rev-parse --short HEAD`.

With SQLite inside a container, the files disappear with the container;
mount a folder (`-v $(pwd)/data:/data`) and put the files in `/data/`.

**If it stops at start**, the last lines say why, for example
`Required environment variable is missing BASE_RPC_URL` or
`[main]Error opening DB ...`. See [Troubleshooting](#8-troubleshooting).

## 6. Check it works

```bash
curl -s localhost:8080/health    # answers 200 with the version: the process is up
curl -s localhost:8080/ready     # each dependency (databases, node, Redis) with "status": "up", and the scan's progress
```

Point your hosting platform's liveness check at `/health` and its
readiness check at `/ready`. When no wallets are tracked yet, `/ready`
reports `stream.status: "idle"`; that is normal on a new deployment.

Then make a small test transfer to a Trovo test wallet and, within a minute
or so, see it in the app's history (or in `payment_history`).

The API reference for these two endpoints is at
<http://localhost:8080/swagger/index.html>.

## 7. Running in production, updating and rolling back

- **One copy.** The engine keeps its scan position in the tracking database
  and is designed to run as a single instance.
- **Stopping is safe.** It resumes from the last saved block; a block it
  could not read is retried, never skipped.
- **Updating:**
  ```bash
  git pull
  docker build -t payment-history-engine --build-arg APP_VERSION=$(git rev-parse --short HEAD) .
  docker rm -f payment-history-engine
  docker run -d --name payment-history-engine --restart unless-stopped -p 8080:8080 --env-file .env payment-history-engine
  ```
  Deploy app-backend first if the release changed the `payment_history`
  table (app-backend owns it).
- **Rolling back:** `git checkout <previous-tag>`, rebuild and replace the
  container the same way.
- There is no CI/CD pipeline in this repository; these steps are run by
  hand (see the root [ARCHITECTURE.md](../ARCHITECTURE.md)).

## 8. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `Required environment variable is missing ENABLE_CACHING` | `ENABLE_CACHING` is empty | set `ENABLE_CACHING=0` (or `1` with Redis) |
| `Required environment variable is missing BASE_RPC_URL` / `NATIVE_ASSET_CODE` | not set | set it in `.env` |
| `[main]Error opening DB ...` | wrong `DB_CONNECTION_STRING`, or the database is unreachable | test it: `psql "<the value>"` |
| `[main]Error opening RoachDB ...` | wrong `CDB_CONNECTION_STRING` | test it: `psql "<the value>"` (CockroachDB speaks the Postgres protocol) |
| `redis host of port not configured properly` | `ENABLE_CACHING=1` without `REDIS_HOST` / `REDIS_PORT` | set them, or `ENABLE_CACHING=0` |
| `env CACHING_PARAMETER not configured` | `ENABLE_CACHING=1` without it | set `CACHING_PARAMETER=v1` |
| `/ready` fails on the node check | the node URL is wrong or down | `curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber"}' <BASE_RPC_URL>` |
| Transfers missing from history | wrong `BASE_CHAIN_ID` for the network, or the wallet is not tracked | set `8453` on mainnet; check the wallet is in `tracked_wallets` |
| Errors about missing `payment_history` | app-backend has not created its tables | start app-backend against the same database first |

## Run the tests

```bash
go test ./...
```

The end-to-end test that records ETH sent by Safe wallets needs the
paymaster project's local blockchain
([paymaster/DEPLOYMENT.md](../paymaster/DEPLOYMENT.md#local-development-stack)):

- **`AA_LOCAL_STACK`:** the deployment file written by the local stack
  (`paymaster/contracts/deployments/local-stack.json`). Without it the test
  is skipped.
- **`AA_SAFE_ETH_TX_FILE`:** a file app-backend's `TestSwapsOnLocalChain`
  writes (`<tx hash> <wallet>`) after leaving a wallet's ETH batch send on
  the local chain; run that test first with the same variable. Without it
  the test is skipped.

## Notes

- **Database driver fix.** An old Postgres driver version broke
  `AutoMigrate` on every restart after the first against an existing
  database. It is fixed (driver `v1.5.11` or newer); nothing to do.
- **The Dockerfile** builds a static binary in `golang:alpine` (with the
  `APP_VERSION` build argument stamped into it) and copies only that binary
  into a small `alpine` image.

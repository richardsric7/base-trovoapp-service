# Deployment

A step-by-step guide to getting `tm-api` running, written for someone doing
it for the first time.

## 1. What you are deploying

`tm-api` is the backend of **Trovo Manager**, the internal admin dashboard
(`tm-web` is its website). Staff use it to manage users, assets, payouts,
Public Markets, P2P, service links and the vault manager. It works directly
in app-backend's database and calls app-backend's API for actions such as
admin logins and payments.

**It needs these to be running first:**

| What | Why | Where |
|---|---|---|
| app-backend and its Postgres database | tm-api reads and writes that database, and calls app-backend's API | [app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md) |
| A Trovo Wallet account for each first admin | admins log in by approving a request in the Trovo Wallet app | the mobile app |
| Redis (recommended) | instant payout-engine wake-ups, rate limiting, and correct logins with more than one copy | the same Redis as payout-engine |
| HashiCorp Vault (for the vault manager) | stores platform signing keys | your Vault |

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Go 1.26 or newer | <https://go.dev/doc/install> | `go version` |
| Docker (recommended for servers) | <https://docs.docker.com/get-docker/> | `docker --version` |
| `psql` | part of PostgreSQL: <https://www.postgresql.org/download/> | `psql --version` |
| `make` (optional, for the check commands) | usually preinstalled | `make --version` |
| `golangci-lint` (optional, for `make lint`) | <https://golangci-lint.run/usage/install/> | `golangci-lint --version` |

## 3. Create Trovo Manager's service link (once)

tm-api identifies itself to app-backend with a **service link**: a row in
app-backend's `service_links` table with an API key. Trovo Manager's
**Service Links** page manages these, but on a brand-new system you cannot
log in yet, so create the first one directly in the database.

1. Choose the Trovo Wallet user that owns it (any existing user, for
   example a company account) and find their wallet address:
   ```bash
   psql "<app-backend database>" -c "SELECT username, address FROM users WHERE username = 'trovomanager';"
   ```
2. Generate an API key (48 characters):
   ```bash
   openssl rand -hex 24
   ```
3. Insert the service link with the permissions Trovo Manager uses (login,
   authorization, payments, push notifications, user info, token info and
   tokenized-asset authorization):
   ```sql
   INSERT INTO service_links (
     id, created_at, updated_at, owner_username, address, api_key, short_name, long_name,
     login_permission, payment_permission, token_info_permission, authorization_permission,
     allow_user_info, push_notification_permission, tokenized_asset_authorization_permission, verified
   ) VALUES (
     gen_random_uuid()::text, now(), now(), 'trovomanager', '<wallet address from step 1>',
     '<api key from step 2>', 'trovomanager', 'Trovo Manager',
     1, 1, 1, 1, 1, 1, 1, 1
   );
   ```
4. Use `owner_username` as [`SERVICE_LINK_USERNAME`](CONFIGURATION.md#service_link_username)
   and the API key as [`SERVICE_LINK_API_KEY`](CONFIGURATION.md#service_link_api_key).

## 4. Get the code and build it

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/tm-api
go mod download
go build -o main .
```

**You should see:** no output from `go build`. Optionally run all checks
with `make ci` (tidy check, build, vet, lint, tests).

## 5. Set the parameters

```bash
cp .env-sample .env
```

Edit `.env` (each setting is explained, with an example and where to get
it, in [CONFIGURATION.md](CONFIGURATION.md)). At least:

| Parameter | Example |
|---|---|
| [`ADMIN_CONNECTION_STRING`](CONFIGURATION.md#admin_connection_string) | app-backend's database: `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require` |
| [`TROVO_WALLET_BASE_URL`](CONFIGURATION.md#trovo_wallet_base_url) | `https://api.trovo.example.com` |
| [`SERVICE_LINK_USERNAME`](CONFIGURATION.md#service_link_username) / [`SERVICE_LINK_API_KEY`](CONFIGURATION.md#service_link_api_key) | from step 3 |
| [`LOGIN_CALLBACK_URL`](CONFIGURATION.md#login_callback_url) | `https://manager-api.trovo.example.com/api/v1/callbacks/login` |
| [`JWT_SECRET`](CONFIGURATION.md#jwt_secret) | `openssl rand -hex 32` |
| [`DEFAULT_SUPER_ADMINS`](CONFIGURATION.md#default_super_admins) | `ada,chidi` |
| [`TROVO_MANAGER_BASE_URL`](CONFIGURATION.md#trovo_manager_base_url), [`ADMIN_DASHBOARD_URL`](CONFIGURATION.md#admin_dashboard_url) | `https://manager.trovo.example.com` |
| [`NATIVE_ASSET_CODE`](CONFIGURATION.md#native_asset_code) | `ETH` |
| [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url), [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id) | `https://base-mainnet.g.alchemy.com/v2/your-api-key`, `8453` |
| [`GIN_MODE`](CONFIGURATION.md#gin_mode) | `release` (in production) |

Recommended: the [Redis settings](CONFIGURATION.md#7-redis-cache-wake-ups-and-rate-limiting),
[Vault](CONFIGURATION.md#9-vault-vault-manager-and-personal-settings),
[Mailgun](CONFIGURATION.md#11-email-mailgun),
[`CORS_ALLOWED_ORIGINS`](CONFIGURATION.md#cors_allowed_origins) and your
own [Discord webhooks](CONFIGURATION.md#13-discord-alerts).

Never commit `.env`. In production (`GIN_MODE=release`) the `.env` file is
ignored: set the values in your hosting platform or with `--env-file`.

## 6. Prepare the database (once)

tm-api creates its own tables when it starts. One table, careers, uses a
Postgres type that needs a one-time script:

```bash
psql "<ADMIN_CONNECTION_STRING>" -f migrations/create_career_roles.sql
```

## 7. Run it

**A. Directly:**

```bash
./main          # or: go run main.go / make run
```

**B. With Docker (recommended for servers):**

```bash
docker build --build-arg APP_VERSION=$(git rev-parse --short HEAD) -t tm-api .
docker run -d --name tm-api --restart unless-stopped -p 8082:8082 --env-file .env tm-api
```

`APP_VERSION` is a **build argument**:

- **What it is:** the version of the code, stamped into the program.
- **Why it's needed:** `/health` reports it, so you know which build runs.
- **Required:** No (default `unknown`).
- **Example:** `a1b2c3d`
- **How to get it:** `git rev-parse --short HEAD`.

The image build also regenerates the API reference (Swagger) from the code.

Inside Docker, `localhost` in a connection string means the container
itself: use the database's real host name (or `host.docker.internal` with
Docker Desktop).

**At start tm-api:** loads `.env` (not in `release` mode), connects to the
database and creates or updates its tables, makes the
`DEFAULT_SUPER_ADMINS` super admins, sets up roles and permissions, and
listens on `PORT` (default `8082`).

## 8. Check it works

```bash
curl http://localhost:8082/health     # 200 while the process runs (liveness), with the version
curl http://localhost:8082/ready      # checks the database and other dependencies (readiness)
```

Point your hosting platform's liveness check at `/health` and readiness
check at `/ready`. The API reference is at
<http://localhost:8082/swagger/index.html>.

Then deploy tm-web ([tm-web/DEPLOYMENT.md](../tm-web/DEPLOYMENT.md)),
open it and log in as one of the `DEFAULT_SUPER_ADMINS`: enter your Trovo
Wallet username, approve the request in the Trovo Wallet app, and the
dashboard opens.

## 9. Running in production, updating and rolling back

- **More than one copy** is safe behind a load balancer. Turn Redis on
  (`ENABLE_CACHING=1` and the `REDIS_*` settings): it carries login
  notifications between copies, shares the rate limit, and makes vault
  manager key changes run one at a time. Table updates at start are
  protected by a database lock, so copies starting together do not
  conflict.
- **Stopping** (`docker stop`) is graceful.
- **Updating:**
  ```bash
  git pull
  docker build --build-arg APP_VERSION=$(git rev-parse --short HEAD) -t tm-api .
  docker rm -f tm-api
  docker run -d --name tm-api --restart unless-stopped -p 8082:8082 --env-file .env tm-api
  ```
  Deploy the matching app-backend first when a release changes shared
  tables.
- **Rolling back:** `git checkout <previous-tag>`, rebuild and replace the
  container the same way.
- There is no CI/CD pipeline in this repository: run `make ci` before
  merging, and build and deploy the image by hand.

## 10. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `empty connection string. Check env variable: ADMIN_CONNECTION_STRING` | not set (or `.env` ignored because `GIN_MODE=release`) | set it in the environment |
| database connection errors at start | wrong connection string, or unreachable from the container | test with `psql "<value>"` from the same machine |
| login page says the user is not an admin | the username is not in `DEFAULT_SUPER_ADMINS` (or not added as an admin), or is not a Trovo Wallet user | add it to `DEFAULT_SUPER_ADMINS` and restart, or have an admin invite them |
| login request never reaches the phone | wrong `TROVO_WALLET_BASE_URL`, `SERVICE_LINK_USERNAME` or `SERVICE_LINK_API_KEY`, or the service link lacks `login_permission` | check step 3 and the values |
| approving on the phone takes several seconds to log in | `LOGIN_CALLBACK_URL` is wrong, so the page falls back to polling | set it to tm-api's public address + `/api/v1/callbacks/login` |
| browser shows CORS errors | `CORS_ALLOWED_ORIGINS` does not list tm-web's address | add it, or leave the setting unset |
| vault manager pages answer 503 | `VAULT_ADDR` / `VAULT_TOKEN` missing or Vault unreachable | set them and check Vault |
| payments from Trovo Manager fail with "native asset code is empty" | `NATIVE_ASSET_CODE` not set | set `ETH` |
| invite emails link to the wrong site | `TROVO_MANAGER_BASE_URL` / `ADMIN_DASHBOARD_URL` not set | set them to tm-web's address |

## Swagger (API reference)

The image build regenerates `docs/` from the code. If you change a route or
its annotations, regenerate locally too so the committed copy stays in step:

```bash
make swagger          # swag init + scripts/fix-swagger-yaml.sh
make swagger-clean    # the same, after deleting docs/
head -5 docs/swagger.yaml   # should start with "swagger:"
```

`scripts/fix-swagger-yaml.sh` moves the `swagger:` and `info:` blocks to the
top of `docs/swagger.yaml` when `swag` writes them in the wrong order, which
some tools reject. Always regenerate with `make swagger`, not bare
`swag init`.

## Notes

- **Database driver fix.** An old Postgres driver version broke the table
  updates on every restart after the first. It is fixed (driver `v1.5.11`
  or newer); nothing to do.
- **The Dockerfile** builds in `golang:alpine` (installing `swag`,
  regenerating the API reference, stamping `APP_VERSION`) and copies only
  the program into a small `alpine` image.

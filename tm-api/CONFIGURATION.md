# Configuration

Every parameter tm-api reads. Each one says what it is, why it is needed,
whether you must set it, an example and how to get a real value. Examples
are never real secrets.

## How to set them

tm-api reads **environment variables**:

- **Locally**, from a `.env` file: `cp .env-sample .env`, then edit it. It
  is loaded automatically **unless** `GIN_MODE=release`.
- **In production** (`GIN_MODE=release`), from the container or hosting
  platform's environment settings; the `.env` file is ignored.

`config/config.go` has an unused loader for `ORCHESTRA_*` variables (for
example `ORCHESTRA_PORT`). Nothing calls it: do not set `ORCHESTRA_*`
variables.

## The minimum to start

Only [`ADMIN_CONNECTION_STRING`](#admin_connection_string) stops tm-api
from starting when missing. For Trovo Manager to be **usable** you also
need:

| Parameter | Without it |
|---|---|
| [`TROVO_WALLET_BASE_URL`](#trovo_wallet_base_url), [`SERVICE_LINK_USERNAME`](#service_link_username), [`SERVICE_LINK_API_KEY`](#service_link_api_key) | admins cannot log in (logins are approved in the wallet app through app-backend) |
| [`LOGIN_CALLBACK_URL`](#login_callback_url) | logins are only confirmed by slower polling |
| [`JWT_SECRET`](#jwt_secret) | organization and stakeholder portal logins fail |
| [`TROVO_MANAGER_BASE_URL`](#trovo_manager_base_url), [`ADMIN_DASHBOARD_URL`](#admin_dashboard_url) | links in invite emails point to the wrong site |
| [`DEFAULT_SUPER_ADMINS`](#default_super_admins) | nobody has admin access on a new database |

---

## 1. Server

### `PORT`

- **What it is:** The port tm-api's web server listens on.
- **Why it's needed:** tm-web and app-backend's callbacks reach tm-api here.
- **Required:** No, default `8082`.
- **Example:** `8082`
- **How to get it:** Any free port; your hosting platform may set it for
  you.

### `GIN_MODE`

- **What it is:** The web framework's mode.
- **Why it's needed:** `release` turns off verbose debug logging and makes
  tm-api read settings only from the environment (no `.env` file).
- **Required:** No. Unset means development mode.
- **Example:** `release`
- **How to get it:** `release` in every deployed environment; unset on your
  own machine.

### `APP_ENV`

- **What it is:** The name of the environment.
- **Why it's needed:** It is attached to crash reports so you can tell
  production errors from test ones.
- **Required:** No.
- **Example:** `production` (or `development`, `staging`)
- **How to get it:** The name of where you are deploying.

### `APP_VERSION`

- **What it is:** The version of the code that is running.
- **Why it's needed:** Shown by `/health` and attached to crash reports.
- **Required:** No.
- **Example:** `a1b2c3d`
- **How to get it:** `git rev-parse --short HEAD`, passed in when the image
  is built or started.

### `DEBUG_ADDR`

- **What it is:** An address for a separate profiling server (Go's
  `pprof`).
- **Why it's needed:** Only for diagnosing performance problems. Never
  expose it publicly.
- **Required:** No (off when unset).
- **Example:** `localhost:6060`
- **How to get it:** Set only while investigating; `localhost` keeps it
  private.

## 2. Database

tm-api uses **one** database connection: app-backend's main database. Trovo
Manager's own tables (admins, roles, access logs) are created there next to
app-backend's.

### `ADMIN_CONNECTION_STRING`

- **What it is:** The address and login of the database, in Postgres
  key=value form.
- **Why it's needed:** Everything Trovo Manager shows and changes is in this
  database: app-backend's users, wallets, assets and payouts, and Trovo
  Manager's own admin tables, which tm-api creates there at start. Without
  it tm-api does not start.
- **Required:** Yes.
- **Example:** `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`
- **How to get it:** Point it at **app-backend's database** (the same
  database as app-backend's `DB_CONNECTION_STRING`, written in this
  key=value form). For a quick local check without Postgres, use
  `DB_TYPE=sqlite` and a file path; most screens will then be empty.

### `DB_TYPE`

- **What it is:** `postgres` or `sqlite`.
- **Why it's needed:** `sqlite` makes `ADMIN_CONNECTION_STRING` a file path,
  for a quick local check without a database server.
- **Required:** No, default `postgres`.
- **Example:** `postgres`
- **How to get it:** Leave unset in real deployments.

### `DB_MAX_OPEN_CONNECTIONS`

- **What it is:** The most database connections tm-api keeps open.
- **Why it's needed:** The database's connection limit is shared with
  app-backend and the engines.
- **Required:** No, default `50`.
- **Example:** `20`
- **How to get it:** Ask your database administrator what share of
  `max_connections` tm-api may use (per running copy).

### `DB_MAX_IDLE_CONNECTIONS`

- **What it is:** How many unused connections are kept open for reuse.
- **Why it's needed:** Reusing connections is faster than opening new ones.
- **Required:** No, default `50`.
- **Example:** `20`
- **How to get it:** Usually the same as `DB_MAX_OPEN_CONNECTIONS`.

## 3. Login tokens (JWT)

### `JWT_SECRET`

- **What it is:** The secret used to sign and check login tokens for
  **organization members** and the **stakeholder portal**.
- **Why it's needed:** Without it those users' tokens cannot be issued or
  checked. (Trovo admins' tokens are checked against app-backend instead;
  see [INTEGRATION.md](INTEGRATION.md).)
- **Required:** Yes, for the organization and stakeholder portals.
- **Example:** `5f2b8c9e0d1a4f7b3c6e8a2d9f0b1c4e7a3d6f9b2c5e8a1d4f7b0c3e6a9d2f5b`
- **How to get it:** `openssl rand -hex 32`. Use the **same value** on every
  running copy of tm-api, and keep it secret: anyone with it can make valid
  tokens.

### `JWT_ACCESS_SECRET`, `JWT_TOKEN_EXPIRY`, `JWT_REFRESH_TOKEN_EXPIRY`

- **What they are:** The secret and lifetimes (in minutes; defaults `10` and
  `15`) for a token function (`CreateToken`) that nothing calls.
- **Why they're needed:** They aren't, today.
- **Required:** No.
- **Example:** (leave unset)
- **How to get them:** Leave unset.

## 4. Connection to app-backend

tm-api calls app-backend's API for admin logins, payments, authorizations,
push notifications and tokenization (see [INTEGRATION.md](INTEGRATION.md)).

### `TROVO_WALLET_BASE_URL`

- **What it is:** The web address of app-backend.
- **Why it's needed:** Every call tm-api makes to app-backend goes here.
- **Required:** Yes, for logins and most admin actions.
- **Example:** `https://api.trovo.example.com`
- **How to get it:** app-backend's public (or internal) address for this
  environment.

### `SERVICE_LINK_USERNAME`

- **What it is:** The Trovo Wallet username that owns tm-api's service link
  (the `owner_username` of a row in app-backend's `service_links` table).
- **Why it's needed:** app-backend only accepts calls from registered
  service links. tm-api puts this name in the address of every login,
  authorization and user-info request, so app-backend knows which service
  is asking and sends admins' login requests to their wallet app.
- **Required:** Yes, for logins and most admin actions.
- **Example:** `trovomanager`
- **How to get it:** On a new system, create the service link once in the
  database as shown in [DEPLOYMENT.md step 3](DEPLOYMENT.md#3-create-trovo-managers-service-link-once);
  use its `owner_username` here. Afterwards admins can manage service
  links in Trovo Manager's **Service Links** page.

### `SERVICE_LINK_API_KEY`

- **What it is:** The API key paired with `SERVICE_LINK_USERNAME`.
- **Why it's needed:** It proves the calls really come from tm-api; it is
  sent in the `X-TW-SERVICE-LINK-API-KEY` header.
- **Required:** Yes, with `SERVICE_LINK_USERNAME`.
- **Example:** `sl_live_xxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** The `api_key` of that `service_links` row. Treat it as
  a secret.

### `LOGIN_CALLBACK_URL`

- **What it is:** The address app-backend calls when an admin approves a
  login on their phone: tm-api's public address followed by
  `/api/v1/callbacks/login`. tm-api adds `/trovo` to the end, so
  app-backend calls `.../api/v1/callbacks/login/trovo`.
- **Why it's needed:** That call is what tells the waiting Trovo Manager
  login page straight away that the admin approved. If it is wrong, logins
  still work through the page's slower polling fallback.
- **Required:** Yes, for instant login confirmation.
- **Example:** `https://manager-api.trovo.example.com/api/v1/callbacks/login`
- **How to get it:** tm-api's public address + `/api/v1/callbacks/login`.
  For local development behind a tunnel (for example ngrok), the tunnel's
  address + `/api/v1/callbacks/login`. (The stakeholder portal's wallet
  authorizations also send this value as their callback address; those
  requests are confirmed by polling as well.)

### `TROVO_MANAGER_BASE_URL`

- **What it is:** tm-web's public address (the Trovo Manager website, not
  tm-api), despite the name.
- **Why it's needed:** Organization invite emails link to
  `{TROVO_MANAGER_BASE_URL}/organizations/invite/validate-email/{invite}`,
  a tm-web page.
- **Required:** No, default `https://dashboard.dev.admin.trovo.app`. **Set
  your own**, or invites point at the development site.
- **Example:** `https://manager.trovo.example.com`
- **How to get it:** Where you deploy tm-web, without a trailing `/`.

### `TROVO_MANAGER_FRONTEND_URL`

- **What it is:** An old setting no code reads any more (`.env-sample` used
  to list it).
- **Why it's needed:** It isn't; tm-web's address is
  [`TROVO_MANAGER_BASE_URL`](#trovo_manager_base_url) and
  [`ADMIN_DASHBOARD_URL`](#admin_dashboard_url).
- **Required:** No.
- **Example:** (leave unset)
- **How to get it:** Leave unset.

### `NATIVE_ASSET_CODE`

- **What it is:** The code of the blockchain's own currency.
- **Why it's needed:** When tm-api asks app-backend to make a payment, only
  this asset may have no token contract address; every other asset needs
  one. Payments fail with "native asset code is empty" without it.
- **Required:** Yes, for payments made from Trovo Manager.
- **Example:** `ETH`
- **How to get it:** The same value as app-backend's `NATIVE_ASSET_CODE`
  (`ETH` on Base).

## 5. First admins

### `DEFAULT_SUPER_ADMINS`

- **What it is:** Wallet usernames to make super admins every time tm-api
  starts, separated by commas.
- **Why it's needed:** On a new database nobody can log in to Trovo
  Manager; this gives the first people access. Each must already be a
  Trovo Wallet user (in app-backend's `users`).
- **Required:** No, default `obi,toluwase`. **Set your own.**
- **Example:** `ada,chidi`
- **How to get it:** The Trovo Wallet usernames of your first
  administrators (each signs up in the app first).

## 6. Browser access (CORS)

### `CORS_ALLOWED_ORIGINS`

- **What it is:** The web addresses allowed to call tm-api from a browser,
  separated by commas.
- **Why it's needed:** Restricting it stops other websites calling tm-api
  with an admin's browser session.
- **Required:** No. **Unset allows any origin.** Once set, every address
  that needs access (dashboard, organization portal, preview deployments)
  must be listed, or it is blocked.
- **Example:** `https://manager.trovo.example.com,https://portal.trovo.example.com`
- **How to get it:** The addresses where tm-web is deployed.

## 7. Redis (cache, wake-ups and rate limiting)

### `ENABLE_CACHING`

- **What it is:** `1` turns on Redis.
- **Why it's needed:** With Redis, tm-api caches responses, wakes
  payout-engine instantly after admin actions, rate-limits requests across
  copies, and delivers admin login notifications correctly when you run
  more than one copy. Without it, all of these fall back to slower or
  per-copy behaviour; nothing breaks.
- **Required:** No, default off.
- **Example:** `1`
- **How to get it:** `1` in production with a Redis server; `0` locally.

### `REDIS_HOST`

- **What it is:** The Redis server's host name.
- **Why it's needed:** Where Redis is, when `ENABLE_CACHING=1`.
- **Required:** With `ENABLE_CACHING=1`.
- **Example:** `redis.internal`
- **How to get it:** Your Redis server; payout-engine must use the same one.
  Locally: `docker run --name redis -p 6379:6379 -d redis:7` and
  `localhost`.

### `REDIS_PORT`

- **What it is:** The Redis server's port.
- **Why it's needed:** To connect.
- **Required:** With `ENABLE_CACHING=1`.
- **Example:** `6379`
- **How to get it:** `6379` unless your Redis says otherwise.

### `REDIS_PASSWORD`

- **What it is:** The Redis password.
- **Why it's needed:** A protected Redis refuses connections without it.
- **Required:** Only if your Redis has one.
- **Example:** `change-me`
- **How to get it:** From whoever runs Redis.

### `REDIS_SKIP_INSECURE_VERIFY`

- **What it is:** `1` skips checking the Redis server's TLS certificate.
- **Why it's needed:** Only for a self-signed certificate you trust.
- **Required:** No, default off.
- **Example:** `0`
- **How to get it:** Leave unset in production.

### `CACHING_PARAMETER`

- **What it is:** Text added to cache keys.
- **Why it's needed:** Changing it makes all old cache entries unused.
- **Required:** No.
- **Example:** `v1`
- **How to get it:** Leave unset, or change it to clear the cache.

## 8. Rate limiting

A limit on how many requests per minute each caller may make, shared across
all copies through Redis. Without Redis the limit is off.

### `RATE_LIMIT_ENABLED`

- **What it is:** `0` turns the limit off.
- **Why it's needed:** The limit protects tm-api from abuse; turn it off
  only for load tests.
- **Required:** No, default on (when Redis is on).
- **Example:** `1`
- **How to get it:** Leave unset in production.

### `RATE_LIMIT_REQUESTS_PER_MINUTE`

- **What it is:** The requests allowed per caller per minute.
- **Why it's needed:** Sets how strict the limit is.
- **Required:** No, default `300`.
- **Example:** `300`
- **How to get it:** Keep the default; raise it if admins hit the limit.

### `RATE_LIMIT_GLOBAL_PER_MINUTE`

- **What it is:** The same limit for the `global` key only (the general
  form is `RATE_LIMIT_<KEY>_PER_MINUTE`; tm-api has only the `global` key).
  It wins over `RATE_LIMIT_REQUESTS_PER_MINUTE`.
- **Why it's needed:** Rarely; kept for consistency with app-backend.
- **Required:** No.
- **Example:** `300`
- **How to get it:** Leave unset.

## 9. Vault (vault manager and personal settings)

The vault manager ("Vault Signer") keeps platform signing keys in
HashiCorp Vault. Without Vault, its pages answer "unavailable" (503); the
rest of Trovo Manager works.

### `VAULT_ADDR`

- **What it is:** The web address of your Vault server.
- **Why it's needed:** Where signing keys and personal settings are kept.
- **Required:** For the vault manager.
- **Example:** `https://vault.internal.trovo.io:8200`
- **How to get it:** From whoever runs Vault (or install Vault:
  <https://developer.hashicorp.com/vault/install>).

### `VAULT_TOKEN`

- **What it is:** tm-api's Vault access token.
- **Why it's needed:** Vault only lets a token with the right policy read and
  write secrets.
- **Required:** For the vault manager.
- **Example:** `hvs.CAESIJexampleexampleexample`
- **How to get it:** Ask your Vault administrator for a token (or AppRole)
  for tm-api, limited to the paths it manages. Never a root token.

### `PERSONAL_ENVS`

- **What it is:** The kinds of personal secret users may manage themselves,
  as `PREFIX:Label` pairs separated by commas.
- **Why it's needed:** It is the list of personal secret types the
  self-service page offers.
- **Required:** No.
- **Example:** `AUTO_APPROVE:Auto Approval,MARKET_MAKER:Market Maker`
- **How to get it:** Add a `PREFIX:Label` pair for each new kind your team
  wants to allow.

### `PERSONAL_ENV_VAULT_MOUNT`

- **What it is:** The Vault key-value store personal secrets go in.
- **Why it's needed:** Part of where they are stored.
- **Required:** No.
- **Example:** `secret`
- **How to get it:** Your Vault's KV mount (`vault secrets list`).

### `PERSONAL_ENV_VAULT_PATH_PREFIX`

- **What it is:** The folder (path prefix) under that mount.
- **Why it's needed:** Keeps personal secrets apart from others.
- **Required:** No.
- **Example:** `personal-envs`
- **How to get it:** Keep the example unless your Vault layout differs.

## 10. System health screens (observability)

All optional. Each unset value turns its screen off ("not configured").
These are addresses **inside your server network**; browsers never call
them.

### `OBSERVABILITY_SERVICES`

- **What it is:** The backends the health screen checks, as `name=url`
  pairs separated by commas.
- **Why it's needed:** So the screen shows whether each service is up.
- **Required:** No.
- **Example:** `admin-api=http://admin-api:8080,wallet-api=http://wallet-api:8080`
- **How to get it:** Each service's name and internal address (for example
  its Docker service name and port).

### `PROMETHEUS_URL`

- **What it is:** The internal address of Prometheus (metrics).
- **Why it's needed:** For the host panel and trend charts.
- **Required:** No.
- **Example:** `http://prometheus:9090`
- **How to get it:** Where your Prometheus runs.

### `LOKI_URL`

- **What it is:** The internal address of Loki (logs).
- **Why it's needed:** For following a request across services.
- **Required:** No.
- **Example:** `http://loki:3100`
- **How to get it:** Where your Loki runs.

### `GLITCHTIP_URL`

- **What it is:** The internal address of GlitchTip (crash reports).
- **Why it's needed:** For the issues screen; needs `GLITCHTIP_API_TOKEN`
  and `GLITCHTIP_ORG` too.
- **Required:** No.
- **Example:** `http://glitchtip-web:8000`
- **How to get it:** Where your GlitchTip runs.

### `GLITCHTIP_API_TOKEN`

- **What it is:** A GlitchTip API token.
- **Why it's needed:** To read issues.
- **Required:** With `GLITCHTIP_URL`.
- **Example:** `gt_xxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** GlitchTip → your profile → **Auth Tokens** → create one
  with read access.

### `GLITCHTIP_ORG`

- **What it is:** Your GlitchTip organization's short name (slug).
- **Why it's needed:** Issues are read per organization.
- **Required:** With `GLITCHTIP_URL`.
- **Example:** `trovo`
- **How to get it:** The part after `/organizations/` in GlitchTip's address
  bar.

### `GLITCHTIP_PUBLIC_URL`

- **What it is:** GlitchTip's address as people's browsers reach it.
- **Why it's needed:** For "open in GlitchTip" links (the internal address
  above would not open in a browser).
- **Required:** No.
- **Example:** `https://glitchtip.trovo.example.com`
- **How to get it:** The address your team opens GlitchTip at.

### `GRAFANA_URL`

- **What it is:** Grafana's public address.
- **Why it's needed:** For "open in Grafana" links.
- **Required:** No.
- **Example:** `https://grafana.trovo.example.com`
- **How to get it:** The address your team opens Grafana at.

### `SENTRY_DSN`

- **What it is:** Where tm-api sends its own crash reports (GlitchTip or
  Sentry).
- **Why it's needed:** So crashes are reported, not only logged.
- **Required:** No.
- **Example:** `https://<public-key>@glitchtip.trovo.example.com/3`
- **How to get it:** Create a project in GlitchTip (or Sentry) and copy its
  DSN from the project settings.

## 11. Email (Mailgun)

### `MAILGUN_DOMAIN`

- **What it is:** Your Mailgun sending domain.
- **Why it's needed:** All Trovo Manager emails (invites, verification
  codes, notifications) are sent through it.
- **Required:** For email.
- **Example:** `mg.trovo.example.com`
- **How to get it:** Mailgun dashboard → **Sending → Domains** (add and
  verify your domain first).

### `MAILGUN_PRIVATE_API_KEY`

- **What it is:** Mailgun's private API key.
- **Why it's needed:** To send most emails.
- **Required:** For email.
- **Example:** `key-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Mailgun dashboard → **API Keys**.

### `MAILGUN_API_KEY`

- **What it is:** A second name for the Mailgun key, read by one helper
  (`internal/utils/email.go`).
- **Why it's needed:** That helper's emails fail without it.
- **Required:** For email.
- **Example:** `key-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** The same value as `MAILGUN_PRIVATE_API_KEY`.

### `MAIL_SENDER`

- **What it is:** The "From" address of emails.
- **Why it's needed:** Recipients see it.
- **Required:** No, default `TrovoP2P Team <no-reply@MAILGUN_DOMAIN>`.
- **Example:** `Trovo <no-reply@mg.trovo.example.com>`
- **How to get it:** An address on your Mailgun domain.

### `TROVO_ADMIN_MANAGER`

- **What it is:** The sender or reply-to address on some admin emails.
- **Why it's needed:** So replies reach your team.
- **Required:** No.
- **Example:** `admin-team@trovo.example.com`
- **How to get it:** Your team's shared inbox.

### `ADMIN_DASHBOARD_URL`

- **What it is:** The dashboard address placed in admin invite emails.
- **Why it's needed:** So invited admins can open Trovo Manager.
- **Required:** No, but invite links break without it.
- **Example:** `https://manager.trovo.example.com`
- **How to get it:** Where tm-web is deployed.

### `TROVO_CONTACT_EMAIL`

- **What it is:** The support address shown in some emails.
- **Why it's needed:** So recipients know where to ask for help.
- **Required:** No.
- **Example:** `support@trovo.example.com`
- **How to get it:** Your support inbox.

### `EMAIL_VERIFICATION_SUBJECT`

- **What it is:** The subject line of verification-code emails.
- **Why it's needed:** To word it your way.
- **Required:** No, default `Your Trovo Wallet Email Verification Code`.
- **Example:** `Your Trovo verification code`
- **How to get it:** Your choice.

### `ENABLE_EMAIL_VALIDATION`

- **What it is:** `1` checks email addresses with Mailgun's validation
  service when users register.
- **Why it's needed:** Catches mistyped or fake addresses (a paid Mailgun
  feature).
- **Required:** No, default off.
- **Example:** `0`
- **How to get it:** `1` only if you pay for Mailgun validation.

### `MAILGUN_VALIDATOR_API_KEY`

- **What it is:** The key for Mailgun's validation service.
- **Why it's needed:** Needed when `ENABLE_EMAIL_VALIDATION=1`.
- **Required:** With `ENABLE_EMAIL_VALIDATION=1`.
- **Example:** `pubkey-xxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Mailgun dashboard → **API Keys** (validation key).

### `SHOW_MAIL_VALIDATION_RESULT`

- **What it is:** `1` returns the validation result to the caller instead
  of only logging it.
- **Why it's needed:** Only for debugging validation.
- **Required:** No.
- **Example:** `0`
- **How to get it:** Leave unset.

## 12. Blockchain

### `BASE_RPC_URL`

- **What it is:** The web address of a Base blockchain node.
- **Why it's needed:** tm-api reads the chain for the vault manager's
  Safe operations and to check token contracts when they are registered
  (Public Markets). Without it, contracts cannot be registered.
- **Required:** For those features.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
- **How to get it:** A node provider (Alchemy, Infura, QuickNode): create an
  app for Base and copy its HTTPS URL. The same provider as app-backend is
  fine.

### `BASE_CHAIN_ID`

- **What it is:** The network number.
- **Why it's needed:** Transactions are signed for one network.
- **Required:** Yes, with `BASE_RPC_URL`.
- **Example:** `8453`
- **How to get it:** `8453` for Base Mainnet, `84532` for Base Sepolia.

### `RPC_TIMEOUT`

- **What it is:** The longest one request to the node may take.
- **Why it's needed:** A stuck node must not hold requests open. After 3
  failures in a row, node calls are refused for 10 seconds at a time until
  one succeeds.
- **Required:** No, default `30s`.
- **Example:** `30s`
- **How to get it:** Keep the default.

### `BLOCKCHAIN_NETWORK_PASSPHRASE`

- **What it is:** A leftover from the platform's previous blockchain
  (Stellar).
- **Why it's needed:** It isn't, on Base.
- **Required:** No.
- **Example:** (leave unset)
- **How to get it:** Leave unset.

## 13. Discord alerts

Each falls back to a built-in Discord webhook when unset (or shorter than
50 characters). Set your own so alerts reach your team.

### `500_ERROR_WEBHOOK`

- **What it is:** A Discord webhook for unexpected server errors.
- **Why it's needed:** So your team sees crashes.
- **Required:** No.
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** In Discord: channel **Edit Channel → Integrations →
  Webhooks → New Webhook → Copy Webhook URL**.

### `REGISTRATION_ERROR_WEBHOOK`

- **What it is:** A Discord webhook for user registration failures.
- **Why it's needed:** So failed sign-ups are noticed.
- **Required:** No.
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** As above.

### `EXPANSION_NETWORK_ERROR_WEBHOOK`

- **What it is:** A Discord webhook for vault manager and blockchain errors.
- **Why it's needed:** So signing problems are noticed.
- **Required:** No.
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** As above.

### `CONNECTION_WARNING_WEBHOOK`

- **What it is:** A Discord webhook for "database connection pool full"
  warnings.
- **Why it's needed:** A full pool slows everything down.
- **Required:** No.
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** As above.

## 14. Location lookup

### `IPAPI_HOST`

- **What it is:** The address of an IP-to-location service.
- **Why it's needed:** To add a country to user records from their IP
  address. Skipped when unset.
- **Required:** No.
- **Example:** `https://api.ipapi.com/api`
- **How to get it:** Your provider's API address (for example ipapi.com).

### `IPAPI_KEY`

- **What it is:** That service's API key.
- **Why it's needed:** The service needs it. Skipped when unset.
- **Required:** With `IPAPI_HOST`.
- **Example:** `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Sign up at the provider (for example
  <https://ipapi.com>) and copy your key.

## 15. Business settings

### `PROCEED_PAYOUT_APPROVALS_REQUIRED`

- **What it is:** How many different admins must approve a proceeds
  payout's holder list before it can be funded and paid.
- **Why it's needed:** Several people checking a payout guards against
  mistakes and fraud. The admin who prepared it and the one who set its fee
  do not count.
- **Required:** No, default `2` (at least `1`).
- **Example:** `2`
- **How to get it:** A business decision; keep it at 2 or more.

### `TAKER_FEE`

- **What it is:** The taker fee percentage shown in user details.
- **Why it's needed:** Displayed to admins.
- **Required:** No, default `0`.
- **Example:** `0.1`
- **How to get it:** Your current taker fee rate.

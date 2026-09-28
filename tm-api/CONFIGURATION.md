# Configuration

Every environment variable tm-api reads, grouped by purpose. Values below
are **realistic examples, never real secrets** — copy `.env-sample` as a
starting point (`cp .env-sample .env`) and fill in real values as described
here.

A variable with no default in the code fails at the point it's used (usually
loudly — a DB connection error, a 401, or a Discord webhook that silently no-ops).
Only `ADMIN_CONNECTION_STRING` is required for the process to start at all.

> **Note on `config/config.go`:** this repo also has an `envconfig`-based
> config loader (`config.Load()`, prefix `ORCHESTRA_`, so e.g. `ORCHESTRA_PORT`).
> It is **not called from anywhere** in the codebase — `main.go` reads
> environment variables directly via `os.Getenv` throughout. Don't set
> `ORCHESTRA_*` variables expecting them to do anything; they're dead code
> left over from an earlier version of the service.

## 1. Server / runtime

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `PORT` | `8082` | Port the HTTP server listens on. Defaults to `8082` if unset (see `internal/server/models/server.go`). | Pick any free port. |
| `GIN_MODE` | `release` | `release` puts Gin in production mode (less verbose logging) and skips loading `.env` (real env vars are expected instead). Anything else (or unset) is dev mode. | Set to `release` in any deployed environment; leave unset locally. |
| `APP_ENV` | `development` | Environment name attached to crash reports (`development`\|`staging`\|`production`). | Pick the value matching where you're running. |
| `APP_VERSION` | `a1b2c3d` | Build identifier reported by `/health` and attached to crash reports. CI passes the git SHA automatically via the Docker build arg. | Leave unset locally; CI sets it. |
| `DEBUG_ADDR` | `localhost:6060` | If set, starts a `pprof`/debug HTTP server on this address for profiling. Leave unset in normal operation. | Only set when actively profiling. |

## 2. Database

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `ADMIN_CONNECTION_STRING` | `host=localhost user=postgres password=postgres dbname=admin_dashboard port=5432 sslmode=disable` | **Required.** The Postgres DSN (libpq key=value format) tm-api connects to. GORM `AutoMigrate` runs against it on startup for tm-api's own tables. As of this snapshot of the code, this is **also** used for the shared TrovoWalletDB/P2P handles (see README's two-database section) — so in practice it needs to point at a database that already has app-backend's schema (`users`, curated assets, etc) if you want wallet/P2P features to work, not just an empty database. | For **AdminDB**-only local work: any fresh Postgres database, or `DB_TYPE=sqlite` with a local file path (see below). For a **fully working** local setup: ask whoever owns `app-backend`'s local dev environment for the same connection string / a seeded dump they use for `TrovoWalletDB`. In a real deployment, this is provisioned by whoever manages the shared Postgres instance — it should point at the **same values app-backend's own DB config uses**. |
| `DB_TYPE` | `postgres` | `postgres` (default) or `sqlite`. In `sqlite` mode, `ADMIN_CONNECTION_STRING` is a file path (or `:memory:`) instead of a Postgres DSN — useful for a quick local smoke test without a real Postgres instance, though most wallet/P2P-dependent features won't have real data to show. | Leave as `postgres` for anything beyond a quick local check. |
| `DB_MAX_OPEN_CONNECTIONS` | `50` | Max open DB connections in the pool. Defaults to `50`. | Tune based on your Postgres instance's `max_connections` and how many tm-api instances run concurrently. |
| `DB_MAX_IDLE_CONNECTIONS` | `50` | Max idle DB connections kept in the pool. Defaults to `50`. | Same as above. |

## 3. Authentication (JWT)

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `JWT_SECRET` | *(generate — see below)* | HMAC signing/verification secret for **organization-member** tokens (see `internal/middleware/organization_auth_middleware.go` and `internal/middleware/stakeholder_auth_middleware.go` — both sign and verify with this same variable). Note: Trovo-admin tokens (the ones tm-web's admin users carry) are **not** verified locally with this secret — see [INTEGRATION.md](./INTEGRATION.md) for how those are verified against app-backend instead. | Generate a strong random value: `openssl rand -hex 32`. Use the **same value** across every tm-api instance that needs to accept the same organization-member tokens. |
| `JWT_ACCESS_SECRET` | *(generate — see below)* | Signing secret used by `internal/middleware/createToken.go`'s `CreateToken`. This function is not currently called from anywhere in the codebase — treat it as unused/dead unless you're the one wiring it back up. | `openssl rand -hex 32` if you need it for something new; otherwise you can leave this unset. |
| `JWT_TOKEN_EXPIRY` | `10` | Access-token lifetime in **minutes**, used only by the same unused `CreateToken`. Defaults to `10`. | Not needed unless reviving `CreateToken`. |
| `JWT_REFRESH_TOKEN_EXPIRY` | `15` | Refresh-token lifetime in **minutes**, same unused code path. Defaults to `15`. | Not needed unless reviving `CreateToken`. |

## 4. app-backend integration (Trovo SDK / ServiceLink)

These configure `internal/trovosdk`, the HTTP client tm-api uses to call
app-backend's own API — see [INTEGRATION.md](./INTEGRATION.md) for what it's
used for (login approval, JWT verification, payments, authorizations, push
notifications, tokenized-asset data).

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `TROVO_WALLET_BASE_URL` | `https://walletdev.trovotechnologies.com` | Base URL of the app-backend instance tm-api calls (`/v1/servicelinks/...`). Also used directly for a couple of tokenization endpoints. | The URL of whichever app-backend environment you're pointing at (dev/staging/prod) — ask whoever manages app-backend's deployments, or check its own deployment docs. |
| `SERVICE_LINK_USERNAME` | `trovo-manager` | The service-link "username" tm-api authenticates as when calling app-backend. Must correspond to an existing `service_links` row (see README's shared-table section) with the permissions tm-api's calls need. | Provisioned as a row in the shared `service_links` table — ask an existing admin, or create one via the admin dashboard's Service Links screen once you have any working admin access. |
| `SERVICE_LINK_API_KEY` | `replace-with-service-link-api-key` | The API key/secret paired with `SERVICE_LINK_USERNAME`. | The `api_key` column of that same `service_links` row. Treat it as a secret. |
| `LOGIN_CALLBACK_URL` | `https://dashboarddev.trovotechnologies.com/api/v1` | Base URL app-backend calls back to (`{LOGIN_CALLBACK_URL}/callbacks/login/{serviceName}` and `.../callbacks/auth/{serviceName}`) after a user approves a login/authorization request on their phone. Must be a URL that reaches **this** tm-api instance from app-backend. | Your tm-api instance's own externally-reachable base URL + `/api/v1`. For local dev behind a tunnel, use the tunnel's URL. |
| `TROVO_MANAGER_BASE_URL` | `https://dashboarddev.trovotechnologies.com/api/v1` | This tm-api instance's own base URL, used when building links (e.g. in emails). | Same as `LOGIN_CALLBACK_URL`'s base, typically. |
| `TROVO_MANAGER_FRONTEND_URL` | `https://admin-app-dashboard-w9dxf.ondigitalocean.app` | tm-web's own base URL, used when building links back to the dashboard (e.g. invite emails). | The deployed URL of the `tm-web` instance this tm-api serves. |
| `NATIVE_ASSET_CODE` | `XLM` | The chain's native asset code, used by `internal/trovosdk` to decide whether an asset needs a contract address. | Whatever the underlying chain's native asset symbol is (ask the blockchain/wallet-core team). |

## 5. Admin bootstrap

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `DEFAULT_SUPER_ADMINS` | `obi,toluwase` | Comma-separated **Trovo Wallet usernames** to seed as `SUPER_ADMIN` tm-api admins on startup (`seedSuperAdmin` in `main.go`) — each must already exist as a user in the shared wallet DB. Defaults to `obi,toluwase` if unset. | The wallet usernames of whoever should have super-admin access to a fresh environment. |

## 6. CORS

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `CORS_ALLOWED_ORIGINS` | `https://admin-app-dashboard-w9dxf.ondigitalocean.app` | Comma-separated allow-list of origins that get a credentialed CORS response (`internal/middleware/cors.go`). **Leave unset to allow any origin** (the current default) — if you set this, every legitimate frontend origin (dashboard, organization portal, preview deployments) must be listed or it gets CORS-blocked. | List every `tm-web` origin (production, staging, PR previews) that needs to call this tm-api instance. |

## 7. Caching (Redis)

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `ENABLE_CACHING` | `0` | `1` enables the Redis-backed HTTP response cache; anything else disables it (in-memory no-op). Also gates whether `/health`'s cache dependency check runs. | `0` for local dev unless you're specifically testing caching behavior. |
| `REDIS_HOST` | `localhost` | Redis host, used only when `ENABLE_CACHING=1`. | Your Redis instance's hostname. |
| `REDIS_PORT` | `6379` | Redis port. | Your Redis instance's port. |
| `REDIS_PASSWORD` | *(empty)* | Redis auth password, if any. | Your Redis instance's password, if it requires one. |
| `REDIS_SKIP_INSECURE_VERIFY` | `0` | `1` skips TLS certificate verification for the Redis connection. Only use for a self-signed cert you trust; never in production against an untrusted network. | Leave `0`/unset unless you know why you need `1`. |
| `CACHING_PARAMETER` | *(empty)* | Overrides an internal cache key namespace parameter (`internal/cache/main.go`). Rarely needs setting. | Leave unset. |

> **Multi-instance note:** if you run more than one instance of this
> service, `ENABLE_CACHING=1` (a working Redis) is also what makes the
> admin QR-login flow's live notification (`GET
> /login/stream/{loginID}`) work correctly regardless of which instance
> receives the wallet's approval callback (`POST
> /callbacks/login/{serviceName}`) — see `internal/models/streams.go`.
> Without Redis, the login stream still works whenever the callback lands
> on the *same* instance that opened the SSE connection, and the client's
> existing polling fallback (`GET /users/verify/...`) covers the rest, so
> nothing breaks — you just lose the guaranteed-fast push across
> instances.

## 8. Rate limiting

A Redis-backed fixed-window limiter (`INCR`+`EXPIRE`, `internal/middleware/rate_limit_middleware.go`) applied globally via `router.Use(...)` in `main.go`. It's Redis-backed rather than in-memory so the limit is correct across more than one instance — an in-memory counter would be per-process, silently multiplying the effective limit by the instance count. This admin panel is JWT-protected with a much smaller abuse surface than the wallet API, so this is one generous global default rather than per-route budgets.

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `RATE_LIMIT_ENABLED` | `1` | When unset/`1` (and `ENABLE_CACHING=1`), the global rate limit is active. `0` disables it entirely. | Leave unset in production (defaults on); `0` only for local dev/load testing. |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | `300` | Requests-per-minute limit for the global default (key `"global"`). | Leave unset to use the built-in default (300/min); tune if that's too tight or too loose for your traffic. |
| `RATE_LIMIT_GLOBAL_PER_MINUTE` | `300` | Same effect as `RATE_LIMIT_REQUESTS_PER_MINUTE` above, but scoped to just the `"global"` key (`RATE_LIMIT_<KEY>_PER_MINUTE` is the general per-route override pattern; tm-api only registers the one `"global"` key today). Takes precedence over `RATE_LIMIT_REQUESTS_PER_MINUTE` if both are set. | Only set if you need this key's limit to differ from a global env-var change affecting every rate-limited route. |

> **Requires Redis:** if `ENABLE_CACHING=0` or Redis is unreachable, the
> middleware no-ops (fails open) regardless of `RATE_LIMIT_ENABLED` — same
> graceful-degradation posture as the caching layer itself.

## 9. Vault signer module

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `VAULT_ADDR` | `https://vault.example.com` | HashiCorp Vault server address. If unset/unreachable, every `/me/vault-signer/*` and `/admin/vault-signer/*` route returns `503` instead of failing at startup. | Your organization's Vault cluster address. |
| `VAULT_TOKEN` | `replace-with-vault-token` | Vault auth token tm-api uses. | Issued by whoever administers your Vault cluster — typically a scoped token/AppRole for this service, not a root token. |
| `PERSONAL_ENVS` | `AUTO_APPROVE:Auto Approval,MARKET_MAKER:Market Maker` | Comma-separated `prefix:Display Name` pairs defining which "personal env" secret prefixes self-service users are allowed to create/manage. | Defined by whoever owns the vault-signer feature's product requirements — add a new `prefix:Label` pair here to allow a new personal-secret category. |
| `PERSONAL_ENV_VAULT_MOUNT` | `secret` | Vault KV mount path for personal-env secrets. | Match whatever mount your Vault cluster uses for this. |
| `PERSONAL_ENV_VAULT_PATH_PREFIX` | `personal-envs` | Path prefix under that mount for personal-env secrets. | Match your Vault path layout, or leave as-is if unsure. |

## 10. Observability (Trovo Manager health screens)

All optional — each unset value disables its screen, which reports "not
configured" instead of showing an empty result. These are **internal service
names on the Docker network**; the browser never contacts them directly (see
`internal/components/observability`).

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `OBSERVABILITY_SERVICES` | `admin-api=http://admin-api-dev:8080,wallet-api=http://wallet-api-dev:8080,p2p-api=http://p2p-dev:8080` | `name=url` pairs of backends the health screen polls. Configured rather than auto-derived because container names don't follow one convention across environments. | List every backend service + its internal Docker/network URL for this environment. |
| `PROMETHEUS_URL` | `http://prometheus:9090` | Metrics source for the health screen's host panel and trend charts. | Your Prometheus instance's internal URL. |
| `LOKI_URL` | `http://loki:3100` | Log aggregation source, for tracing a request across services. | Your Loki instance's internal URL. |
| `GLITCHTIP_URL` | `http://glitchtip-web:8000` | Internal crash-reporter URL the backend calls, for the issues screen. Needed together with the two below. | Your GlitchTip instance's internal URL. |
| `GLITCHTIP_API_TOKEN` | *(empty)* | GlitchTip API token for reading issues. | Generate one in your GlitchTip project's settings. |
| `GLITCHTIP_ORG` | *(empty)* | GlitchTip organization slug. | Your GlitchTip org's slug. |
| `GLITCHTIP_PUBLIC_URL` | *(empty)* | The **browser-reachable** GlitchTip URL (distinct from `GLITCHTIP_URL`, which is internal-only) — used for "open in GlitchTip" links. | The public URL your team actually opens in a browser. |
| `GRAFANA_URL` | *(empty)* | Public Grafana URL, shown as an "open in Grafana" link. | Your team's Grafana URL. |
| `SENTRY_DSN` | *(empty)* | GlitchTip/Sentry-compatible DSN for crash reporting (`internal/observe`). Leave empty to disable — panics are still recovered and logged locally, just not reported externally. | Create a project in GlitchTip (or Sentry) and copy its DSN. |

## 11. Email (Mailgun)

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `MAILGUN_DOMAIN` | `mg.yourcompany.com` | Mailgun sending domain, used by most of `internal/mail` and `internal/utils/email.go`. | From your [Mailgun dashboard](https://app.mailgun.com/) → Sending → Domains. |
| `MAILGUN_PRIVATE_API_KEY` | `key-xxxxxxxx` | Mailgun private API key, used by `internal/mail/main.go`'s email-sending functions (invite/notification emails). | Mailgun dashboard → Settings → API Keys. |
| `MAILGUN_API_KEY` | `key-xxxxxxxx` | **A second, separately-named Mailgun key** read only by `internal/utils/email.go` — this looks like a naming inconsistency with `MAILGUN_PRIVATE_API_KEY` above rather than an intentionally distinct credential; set it to the same value unless you know otherwise. | Same Mailgun API key as above. |
| `MAIL_SENDER` | `noreply@yourcompany.com` | "From" address for most outgoing email. | A verified sender address on your Mailgun domain. |
| `TROVO_ADMIN_MANAGER` | `admin-team@yourcompany.com` | Secondary "admin" sender/reply-to address used on some templates. | Your team's shared admin inbox. |
| `ADMIN_DASHBOARD_URL` | `https://dashboarddev.trovotechnologies.com` | URL embedded in outgoing admin-invite emails as the `dashboard_url` template variable. | Your `tm-web` deployment's URL. |
| `TROVO_CONTACT_EMAIL` | `support@yourcompany.com` | Contact email embedded in some email templates. | Your team's support inbox. |
| `EMAIL_VERIFICATION_SUBJECT` | `Verify your email` | Subject line for the email-verification message. | Any subject line you want. |
| `ENABLE_EMAIL_VALIDATION` | `0` | `1` enables Mailgun's email-address validation API when registering a user (`internal/models/user_registration_info.go`). | `0` unless you specifically want registration-time email validation and are paying for Mailgun's validation API. |
| `MAILGUN_VALIDATOR_API_KEY` | `pubkey-xxxxxxxx` | API key for Mailgun's email validation API, required if `ENABLE_EMAIL_VALIDATION=1`. | Mailgun dashboard → Settings → API Keys (the public validation key). |
| `SHOW_MAIL_VALIDATION_RESULT` | `0` | `1` surfaces the validation result to the caller instead of only logging it. | `0` unless debugging validation behavior. |

## 12. Blockchain / network

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `BASE_RPC_URL` | `https://mainnet.base.org` | RPC endpoint used by `internal/network` for on-chain reads. | An RPC provider URL for the chain you're targeting (e.g. a Base RPC provider). |
| `BASE_CHAIN_ID` | `8453` | Numeric chain ID matching `BASE_RPC_URL`. | The chain ID for whichever network `BASE_RPC_URL` points at. |
| `BLOCKCHAIN_NETWORK_PASSPHRASE` | *(network-specific)* | Network passphrase used when constructing/validating blockchain transactions. | Ask the wallet-core/blockchain team for the correct value per network. |

## 13. Legacy Discord error-alert webhooks

Each of these is optional; if unset (or shorter than 50 characters), the
corresponding alert falls back to a shared, hardcoded Discord webhook already
in the source — a real value routes alerts to your own channel instead.

| Variable | Effect |
|---|---|
| `500_ERROR_WEBHOOK` | General unhandled-error alerts (auth, users, general, swagger components). |
| `REGISTRATION_ERROR_WEBHOOK` | User-registration failure alerts. |
| `EXPANSION_NETWORK_ERROR_WEBHOOK` | Vault-signer / network-expansion error alerts. |
| `CONNECTION_WARNING_WEBHOOK` | DB connection-pool fill-up warnings. |

**How to get a real value:** create a webhook on a Discord channel you own
(Channel Settings → Integrations → Webhooks → New Webhook) and copy its URL.

## 14. Geo / IP lookup

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `IPAPI_KEY` | `replace-with-ipapi-key` | API key for the IP-geolocation lookup used to enrich user records (`internal/models/geo.go`). If unset, geo lookup is skipped. | Sign up at [ipapi.com](https://ipapi.com/) (or whichever provider `IPAPI_HOST` points at) and copy your key. |
| `IPAPI_HOST` | `https://api.ipapi.com/api` | Base URL of the IP-geolocation API. If unset, geo lookup is skipped. | The base URL for your ipapi.com-compatible provider. |

## 15. Misc

| Variable | Example | Effect | How to get a real value |
|---|---|---|---|
| `TAKER_FEE` | `0.1` | Overrides the taker-fee percentage shown on user-info responses (`internal/components/users/services/get_user.go`). Leave unset to use the built-in default. | A business decision — set to whatever the current taker-fee rate is. |

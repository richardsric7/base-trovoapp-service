# Trovo Wallet API (`app-backend`)

## What this project does

`app-backend` is the main server of the Trovo Wallet platform. When someone
uses the Trovo mobile app or web wallet, every action (signing up, seeing a
balance, sending money, trading) is a request to this server. It:

- **keeps user accounts** and creates each user's wallet, a Safe smart
  account on the **Base** blockchain, plus sub-wallets and shared wallets;
- **sends payments and swaps** for users, paying gas in ETH or a
  stablecoin through the paymaster;
- runs the **P2P marketplace** (people buy and sell crypto for local
  currency, with the platform holding funds in escrow);
- issues and sells **tokenized assets** (real-world assets as tokens) and
  runs **Public Markets** (tokenized Nigerian stocks and bonds);
- handles **KYC** (identity checks), **account recovery**, bank deposits and
  withdrawals, and emails, SMS and push notifications;
- gives **partner businesses** an API ("service links") to offer these
  features in their own products.

It is a REST API written in Go. If you are new, read this file, then
[DEPLOYMENT.md](DEPLOYMENT.md) to run it, [CONFIGURATION.md](CONFIGURATION.md)
for every setting and [INTEGRATION.md](INTEGRATION.md) for everything it
connects to.

## Where this fits in the monorepo

| Project | Relationship to `app-backend` |
| --- | --- |
| **`app-mobile`**, **`app-web`** | The user apps. They call this API, signing each request with the user's key. |
| **`tm-api`** (and its website `tm-web`) | Trovo Manager, the admin dashboard. Shares this project's database for settings and catalog data, and calls this API for actions with business rules (admin logins, payments, P2P disputes). |
| **`payment-history-engine`** | Watches the blockchain and writes users' payment history. app-backend tells it which wallets to watch through a shared "tracking" database. |
| **`payout-engine`** | Pays dividends and interest on tokenized assets, working in this project's database. |
| **`paymaster`** | The paymaster contract and quote service every wallet transaction uses. |
| **`market`**, **`recovery`** | The offer book contract (sales, swaps) and the account recovery module. |
| **`wallet-core`** | Key and Safe address code shared with the apps. |

See [INTEGRATION.md](INTEGRATION.md) for how each connection works.

## Tech stack

- **Language**: Go 1.26 (see `go.mod`)
- **Web framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: PostgreSQL through [GORM](https://gorm.io/); app-backend
  creates and updates its own tables (see [DEPLOYMENT.md](DEPLOYMENT.md)).
  SQLite works for local experiments (`DB_TYPE=sqlite`).
- **Blockchain**: Base, through `go-ethereum`; wallets are Safe smart
  accounts sent through an ERC-4337 bundler.
- **Also**: Redis (cache, rate limits, websocket fan-out), Firebase (push
  notifications, file storage), Mailgun (email), a second Postgres or
  CockroachDB database ("RoachDB") for payment-history tracking.
- **API reference**: [swaggo/swag](https://github.com/swaggo/swag)-generated
  Swagger, served at `/swagger/index.html` once the app is running.

## Directory structure

```
app-backend/
├── main.go              # entrypoint: env checks, DB open + migrate, router wiring, background jobs
├── docs/                 # swag-generated OpenAPI output (docs.go, swagger.json, swagger.yaml) — do not hand-edit
├── Dockerfile            # multi-stage build → small alpine runtime image
├── Makefile               # CI-parity targets: build, vet, lint, test
└── internal/
    ├── components/        # one folder per business domain — see below
    ├── db/                 # DB connection setup + GORM AutoMigrate model list
    ├── middleware/          # authentication (signature, API key, JWT), CORS
    ├── network/              # Base/EVM chain client & helpers
    ├── basetxn/               # transaction building/signing for the Base chain
    ├── evmkeypair/             # EVM key parsing/generation
    ├── cache/                   # Redis-backed response cache
    ├── mail/, sms/, pns/         # email (Mailgun), SMS (Termii/Infobip), push notifications (Firebase)
    ├── dynamiclinks/               # Firebase Dynamic Links (referral/deep links)
    ├── sharedconfig/                # the GlobalConfig struct threaded through every handler
    └── errors/, validators/          # shared error types and input validators
```

Each folder under `internal/components/` is a self-contained business
domain, and each follows the same internal shape:
`controllers/` (HTTP handlers + route registration), `models/` (GORM
models), `services/` (business logic), and sometimes `db/` (queries) or
`blockchain/`/`safesigner/` (on-chain interactions specific to that domain).

| Component | Owns |
| --- | --- |
| `root` | A signed `GET /` check, and the Apple/Android app-link files under `/.well-known/` |
| `users` | Accounts, wallets, KYC, shared access, tokenized assets, subwallets, account recovery — the largest component |
| `payments` | Sending assets between wallets |
| `swaps` | Currency/asset swaps |
| `assets` | Curated asset catalog (read side — the catalog itself is largely maintained by `tm-api`) |
| `rates` | Exchange rates |
| `p2p` | The P2P marketplace: offers, orders, escrow, disputes, refunds, merchant status |
| `servicelinks` | White-label partner API-key integrations, plus the `tm-api` server-to-server endpoints |
| `callbacks` | Webhook receivers from external providers (1Liquidity, Doja, Flutterwave) |
| `announcements` | In-app announcement banners and minimum app-version checks |

### Loading users

A user can be loaded two ways. Pick the lighter one whenever you can:

| Loader | Loads | Use when |
| --- | --- | --- |
| `Username(x).GetFullUser`, `usersDB.GetUser`, `usersDB.GetUserFromPrimarySigner` | the `users` row **plus** wallets (with permissions), wallets shared with the user, patron membership, closed groups and fiat payment methods (about seven queries) | the code reads `UserWallets`, `WalletsSharedWithUser` or another association |
| `Username(x).GetSlimUser` (and `GetSimpleUser`), `usersDB.GetSlimUser`, `usersDB.GetSlimUserFromPrimarySigner` / `userModels.GetSlimUserBySigner` | the `users` row only, with every column, so it can be saved back | only the user's own fields are read (username, push token, KYC, suspension, referrer, fee exemption, ...) |

Both are cached in Redis for 2000 seconds under the user's username, email,
primary signer and ID. Full users live under `userObj …` and slim users
under `userLite …`, so a slim object is never served to code that expects
wallets. `InvalidateUserCache` / `InvalidateUserWalletCache` clear both. On
a slim user they read the wallet ids from the database to clear the
wallet-keyed entries too, so you can invalidate through whichever kind you
hold. Tests: `internal/components/users/db/user_slim_test.go`, which needs a
Redis at `TEST_REDIS_ADDR` (default `127.0.0.1:6379`) and is skipped
without one.

## Running it locally

The short version (the full guide, with a local Postgres and Redis in
Docker, is [DEPLOYMENT.md](DEPLOYMENT.md)):

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/app-backend
cp .env.example .env        # fill it in; CONFIGURATION.md explains each value
go mod download
go test ./internal/... -vet=off
go run main.go              # listens on PORT, default 8080
```

The server logs each part as it starts (`##users services initialized##`,
...) and finally `##service started##`.

### Swagger UI

Once the server is running, the full per-endpoint API reference is at:

```
http://localhost:8080/swagger/index.html
```

This is generated from `@Summary`/`@Router`/etc. comments above each
handler function (see the `docs/` folder) — it is always the most
up-to-date and precise source for what a given endpoint accepts and
returns. See [INTEGRATION.md](INTEGRATION.md) for how to authenticate
requests you try from the Swagger UI or elsewhere.

## Further reading

- **[DEPLOYMENT.md](DEPLOYMENT.md)** — how to build, containerize, and
  deploy this service, and what the migration/`DB_AUTOMIGRATE` behavior
  means for you.
- **[CONFIGURATION.md](CONFIGURATION.md)** — every environment variable
  this service reads, what it does, and how to get a real value for it.
- **[INTEGRATION.md](INTEGRATION.md)** — how `app-web`, `app-mobile`,
  `tm-api`, and white-label partners each integrate with this API.
- **[PUBLIC_MARKETS.md](PUBLIC_MARKETS.md)** — the Public Markets engine:
  order paths, net batches, partners (mock/REST/manual), reconciliation,
  exchanges and dividends.

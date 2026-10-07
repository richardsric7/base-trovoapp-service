# Integration

Everything `app-backend` connects to, and how. app-backend is the hub of
the platform: the apps, Trovo Manager and partner businesses call it; it
calls the blockchain and a set of outside services; and three other
projects share its databases.

For the exact request and response of any endpoint, the reference is the
Swagger UI (<http://localhost:8080/swagger/index.html> on a running
server), not this document.

| Connects to | Direction | Through | Needed? |
|---|---|---|---|
| [1. The apps (app-mobile, app-web)](#1-the-apps-app-mobile-and-app-web) | apps → app-backend | HTTPS (signed requests) and a websocket | yes |
| [2. Trovo Manager (tm-api)](#2-trovo-manager-tm-api) | both ways | the same database, plus HTTPS | yes, for admin |
| [3. Partner businesses (service links)](#3-partner-businesses-service-links) | partners → app-backend, app-backend → partners | HTTPS with an API key; callbacks | optional |
| [4. Main database](#4-main-database-postgres) | app-backend reads and writes | Postgres | yes |
| [5. payment-history-engine](#5-payment-history-engine-and-the-tracking-database) | app-backend adds wallets; the engine writes history | the tracking database and the main database | yes, for history |
| [6. payout-engine](#6-payout-engine) | they share payout tables | the main database | optional |
| [7. Redis](#7-redis) | app-backend reads and writes | Redis | recommended |
| [8. Base blockchain](#8-base-blockchain) | app-backend reads and sends | JSON-RPC over HTTPS | yes |
| [9. Bundler, paymaster and quote service](#9-bundler-paymaster-and-quote-service) | app-backend → them | HTTPS | yes, to send |
| [10. Platform contracts](#10-platform-contracts-offer-book-recovery-module-safes) | app-backend reads and sends | the blockchain | per feature |
| [11. Email (Mailgun)](#11-email-mailgun) | app-backend → Mailgun | HTTPS | yes |
| [12. SMS (Termii, Infobip)](#12-sms-termii-infobip) | app-backend → provider | HTTPS | optional |
| [13. Firebase](#13-firebase) | app-backend → Firebase | Google APIs | yes |
| [14. KYC (Sumsub, Doja)](#14-kyc-sumsub-doja) | both ways | HTTPS and webhooks | for KYC |
| [15. Bank deposits and withdrawals (Stablerail)](#15-bank-deposits-and-withdrawals-stablerail) | app-backend → Stablerail | HTTPS | optional |
| [16. Crypto deposits and withdrawals (1Liquidity)](#16-crypto-deposits-and-withdrawals-1liquidity) | both ways | HTTPS and callbacks | optional |
| [17. Public Markets partners and price feed](#17-public-markets-partners-and-price-feed) | both ways | HTTPS and callbacks | for Public Markets |
| [18. Naira price API](#18-naira-price-api) | app-backend → API | HTTPS | yes |
| [19. Location lookup (ipapi)](#19-location-lookup-ipapi) | not called today | — | settings required at start |
| [20. Discord](#20-discord-alerts) | app-backend → Discord | webhooks | optional |

---

## 1. The apps (app-mobile and app-web)

- **What it is and why:** the Trovo mobile app and web wallet. Everything a
  user does goes through app-backend.
- **Direction:** apps → app-backend; app-backend → apps over the websocket
  (live payments, P2P updates, login requests).
- **How they connect:** HTTPS REST. Each request is **signed** with the
  user's key (headers `X-TW-PUBLIC-KEY`, `X-TW-SIGNER`, `X-TW-TIMESTAMP`,
  `X-TW-SIGNATURE`; see [Authentication](#authentication)); no password or
  token is sent. The websocket is `GET /v1/stream/ws/<username>`.
- **Settings on this side:** none beyond `PORT`; put app-backend behind
  HTTPS.
- **Settings on the other side:** app-web's
  [`VITE_API_URL`](../app-web/CONFIGURATION.md) and `VITE_SOCKET_URL`;
  app-mobile's API address ([app-mobile/CONFIGURATION.md](../app-mobile/CONFIGURATION.md)).
  Example: `https://api.trovo.example.com`.
- **How to check it works:** sign up in the app; the user appears in the
  `users` table.
- **When it is down:** the apps cannot load balances or send anything.

## 2. Trovo Manager (tm-api)

- **What it is and why:** the admin dashboard's backend. Staff manage
  users, assets, fees, payouts, Public Markets and partners there.
- **Direction:** both ways, in three forms:
  1. **The same database.** tm-api reads and writes app-backend's tables
     directly for settings and catalog data (curated assets, fees, service
     links, Public Markets settings).
  2. **tm-api → app-backend, with an API key**, for actions with business
     rules: admin logins and authorizations (the admin approves in their
     Trovo app), payments, push notifications, P2P dispute resolution.
  3. **tm-api → app-backend, with the admin's login token**, for
     `/v1/trovo-manager/...` reads (wallet balances, banks, tokenization).
     app-backend issued that token at login and checks it with
     `JWT_ACCESS_SECRET`.
- **Settings on this side:** the [service link](#3-partner-businesses-service-links)
  Trovo Manager uses (a row in `service_links`), and
  [`JWT_ACCESS_SECRET`](CONFIGURATION.md#jwt_access_secret).
- **Settings on the other side:** tm-api's `ADMIN_CONNECTION_STRING` (this
  database), `TROVO_WALLET_BASE_URL` (this server's address),
  `SERVICE_LINK_USERNAME` and `SERVICE_LINK_API_KEY`
  ([tm-api/CONFIGURATION.md](../tm-api/CONFIGURATION.md)).
- **How to check it works:** log in to Trovo Manager: the approval request
  arrives in the Trovo app.
- **When it is down:** the apps keep working; only admin work stops.

More: [How tm-api relates to app-backend](#how-tm-api-relates-to-app-backend).

## 3. Partner businesses (service links)

- **What it is and why:** businesses that offer Trovo wallets, payments,
  tokenized assets or Public Markets inside their own products.
- **Direction:** partner → app-backend for requests; app-backend → partner
  for callbacks (login and authorization results, payment notifications).
- **How they connect:** HTTPS with the header `X-TW-SERVICE-LINK-API-KEY`.
  Routes are under `/v1/servicelinks/...` and `/v1/trovo-api/...`. Callbacks
  are saved in `callback_deliveries` and retried up to 10 times.
- **Settings on this side:** a row in `service_links` per partner, made in
  Trovo Manager's **Service Links** page: the API key, permission flags,
  callback addresses and an optional rate limit
  (`rateLimitPerMinute`). Validity of login and authorization requests:
  [`SERVICE_LINK_LOGIN_REQUEST_VALIDITY`](CONFIGURATION.md#service_link_login_request_validity),
  [`SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY`](CONFIGURATION.md#service_link_authorization_request_validity).
- **Settings on the other side:** the partner stores the API key and
  app-backend's address, and answers on its callback addresses.
- **How to check it works:**
  `curl -H "X-TW-SERVICE-LINK-API-KEY: <key>" https://api.trovo.example.com/v1/trovo-api/users/balance/<address>`
  answers instead of `401`.
- **When it is down (the partner):** callbacks wait and are retried.

More: [Service links](#service-links-white-label-partner-integration).

## 4. Main database (Postgres)

- **What it is and why:** where all accounts, wallets, history, offers,
  orders and settings live.
- **Direction:** app-backend reads and writes, and creates and updates the
  tables (other projects that share it never change its tables).
- **How they connect:** a direct Postgres connection.
- **Settings on this side:** [`DB_CONNECTION_STRING`](CONFIGURATION.md#db_connection_string)
  (example `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`),
  [`DB_TYPE`](CONFIGURATION.md#db_type),
  [`DB_AUTOMIGRATE`](CONFIGURATION.md#db_automigrate).
- **Settings on the other side:** the same database is
  tm-api's `ADMIN_CONNECTION_STRING`, payment-history-engine's and
  payout-engine's `DB_CONNECTION_STRING`.
- **How to check it works:** app-backend starts without `Error opening DB`.
- **When it is down:** app-backend cannot serve requests.

## 5. payment-history-engine and the tracking database

- **What it is and why:** payment-history-engine watches the blockchain for
  transfers to and from Trovo wallets and records them, so users see their
  history. It needs to know which wallets to watch.
- **Direction:** app-backend writes each new wallet (on sign-up, new
  sub-wallet and approvals) into the **tracking database**'s
  `tracked_wallets` and `tracked_addresses`; the engine writes
  `payment_history` in the main database, which app-backend reads. No API
  calls either way.
- **How they connect:** two shared databases.
- **Settings on this side:** [`CDB_CONNECTION_STRING`](CONFIGURATION.md#cdb_connection_string)
  (example `postgresql://root@roach.internal:26257/payment_history?sslmode=require`).
- **Settings on the other side:** the engine's `CDB_CONNECTION_STRING` (the
  same tracking database) and `DB_CONNECTION_STRING` (this main database)
  ([payment-history-engine/CONFIGURATION.md](../payment-history-engine/CONFIGURATION.md)).
- **How to check it works:** after a sign-up, `tracked_wallets` has the new
  wallet; after a transfer, `payment_history` has a row.
- **When it is down (the engine):** history stops updating, and catches up
  when the engine is back.

More: [Payment history: source vs. destination](#payment-history-source-vs-destination).

## 6. payout-engine

- **What it is and why:** pays dividends and interest to holders of
  tokenized assets.
- **Direction:** the engine reads and writes payout tables in the main
  database, which app-backend creates; app-backend shows users their
  payouts.
- **How they connect:** the shared main database.
- **Settings:** none on this side; the engine's `DB_CONNECTION_STRING`
  points at this database ([payout-engine/CONFIGURATION.md](../payout-engine/CONFIGURATION.md)).
- **When it is down:** payouts wait; nothing else is affected.

More: [Proceeds payouts](#proceeds-payouts-dividends-and-interest).

## 7. Redis

- **What it is and why:** caching (balances, prices), rate limits and
  passing websocket messages between copies of app-backend.
- **Direction:** app-backend reads and writes.
- **How they connect:** the Redis protocol.
- **Settings on this side:** [`ENABLE_CACHING`](CONFIGURATION.md#enable_caching),
  [`REDIS_HOST`](CONFIGURATION.md#redis_host) (example `redis.internal`),
  [`REDIS_PORT`](CONFIGURATION.md#redis_port) (example `6379`),
  [`REDIS_PASSWORD`](CONFIGURATION.md#redis_password).
- **Settings on the other side:** none; payout-engine and tm-api may use the
  same Redis.
- **How to check it works:** the start-up log says `Testing redis
  connection...` with no error after it.
- **When it is down:** with caching on, app-backend does not start; once
  running, cached reads fall back to the database and rate limits are off.

## 8. Base blockchain

- **What it is and why:** the network where wallets, tokens and payments
  live.
- **Direction:** app-backend reads (balances, contract state, events) and
  sends (platform transactions such as P2P settlements, minting, recovery).
- **How they connect:** JSON-RPC over HTTPS, with a time limit per request.
- **Settings on this side:** [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url)
  (example `https://base-mainnet.g.alchemy.com/v2/your-api-key`),
  [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id) (`8453` mainnet,
  `84532` Base Sepolia), [`RPC_TIMEOUT`](CONFIGURATION.md#rpc_timeout).
- **Settings on the other side:** an account with an RPC provider (Alchemy,
  QuickNode, Coinbase Developer Platform) or your own node.
- **How to check it works:** a wallet's balance loads in the app.
- **When it is down:** balances and sends fail; nothing is lost.

## 9. Bundler, paymaster and quote service

- **What it is and why:** each user wallet is a Safe smart account. Its
  transactions are "user operations" sent through a **bundler**; gas can be
  paid in a stablecoin through the **paymaster** contract, at a price
  signed by the **quote service**.
- **Direction:** app-backend → quote service (gets a signed quote),
  app-backend → bundler (estimates and submits operations).
- **How they connect:** HTTPS. The quote service needs an API key.
- **Settings on this side:** [`BUNDLER_URL`](CONFIGURATION.md#bundler_url),
  [`PAYMASTER_ADDRESS`](CONFIGURATION.md#paymaster_address),
  [`PAYMASTER_QUOTE_SERVICE_URL`](CONFIGURATION.md#paymaster_quote_service_url)
  (example `http://quote-service.internal:8090`),
  [`PAYMASTER_QUOTE_SERVICE_API_KEY`](CONFIGURATION.md#paymaster_quote_service_api_key).
- **Settings on the other side:** the quote service's API keys and the
  paymaster deployment ([paymaster/CONFIGURATION.md](../paymaster/CONFIGURATION.md)).
- **How to check it works:** a test payment between two users succeeds.
- **When it is down:** no wallet can send; paying gas in ETH still needs
  the bundler.

More: [Wallets and sending](#wallets-and-sending-safe-useroperations).

## 10. Platform contracts (offer book, recovery module, Safes)

- **What it is and why:** contracts app-backend uses for specific features:
  the **offer book** (tokenized asset sales, swaps, market-making), the
  **social recovery module** (account recovery) and Safe's standard
  contracts.
- **Direction:** app-backend reads them and sends transactions to them.
- **Settings on this side:** [`OFFER_BOOK_ADDRESS`](CONFIGURATION.md#offer_book_address),
  [`OFFER_BOOK_START_BLOCK`](CONFIGURATION.md#offer_book_start_block),
  [`RECOVERY_MODULE_ADDRESS`](CONFIGURATION.md#recovery_module_address),
  and the [Safe contract addresses](CONFIGURATION.md#safe-and-entrypoint-contract-addresses)
  (defaults are Safe's official Base deployments).
- **Settings on the other side:** the deployments in
  [market/DEPLOYMENT.md](../market/DEPLOYMENT.md) and
  [recovery/DEPLOYMENT.md](../recovery/DEPLOYMENT.md).
- **When one is missing:** that feature is off (or answers with an error).

More: [Tokenized assets](#tokenized-assets-token-issuing-and-distribution-wallets-sale-offer),
[Account recovery](#account-recovery-opt-in-guardian).

## 11. Email (Mailgun)

- **What it is and why:** verification codes, notifications and receipts
  by email.
- **Direction:** app-backend → Mailgun.
- **Settings on this side:** [`MAILGUN_PRIVATE_API_KEY`](CONFIGURATION.md#mailgun_private_api_key),
  [`MAILGUN_DOMAIN`](CONFIGURATION.md#mailgun_domain) (example
  `mg.trovo.example.com`), [`MAIL_SENDER`](CONFIGURATION.md#mail_sender)
  and the other [email settings](CONFIGURATION.md#17-email-mailgun).
- **Settings on the other side:** a verified sending domain and the email
  templates in your Mailgun account.
- **How to check it works:** sign up with a real address and receive the
  code.
- **When it is down:** emails fail; sign-ups that need a code wait.

## 12. SMS (Termii, Infobip)

- **What it is and why:** verification codes by text message.
- **Direction:** app-backend → provider.
- **Settings on this side:** [`DEFAULT_SMS_PROVIDER`](CONFIGURATION.md#default_sms_provider)
  and the provider's [keys](CONFIGURATION.md#18-sms).
- **Settings on the other side:** an account and approved sender id.
- **When it is down:** phone verification fails.

## 13. Firebase

- **What it is and why:** push notifications to the apps, file storage
  (logos, documents) and app links.
- **Direction:** app-backend → Firebase.
- **Settings on this side:** [`GC`](CONFIGURATION.md#gc) (the service
  account key), [`GOOGLE_PROJECT_ID`](CONFIGURATION.md#google_project_id),
  [`STORAGE_BUCKET_NAME`](CONFIGURATION.md#storage_bucket_name), and the
  [link settings](CONFIGURATION.md#19-firebase-storage-and-links).
- **Settings on the other side:** the same Firebase project in app-mobile
  and app-web.
- **When it is down:** notifications and uploads fail; payments still work.

## 14. KYC (Sumsub, Doja)

- **What it is and why:** identity checks before users reach higher limits
  or features.
- **Direction:** app-backend → provider (create applicants, read results);
  provider → app-backend (webhooks: `POST /v1/webhook/sumsub/kyc/individual`,
  `POST /v1/callbacks/doja/webhook`).
- **Settings on this side:** **not** environment variables: each provider's
  token and secret are a row in the `kyc_configs` table
  (`service_provider` = `sumsub` or `doja`, `token`, `secret_key`).
- **Settings on the other side:** the webhook addresses above, set in the
  provider's dashboard.
- **When it is down:** new KYC checks wait; existing levels stay.

## 15. Bank deposits and withdrawals (Stablerail)

- **What it is and why:** Naira bank deposits (on-ramp to cNGN) and
  withdrawals to bank accounts.
- **Direction:** app-backend → Stablerail; app-backend follows withdrawal
  status itself.
- **Settings on this side:** the `stablerail_configs` table (`api_key`,
  `base_url`, `enable_stablerail`) and
  [`STABLERAIL_MIN_WITHDRAWAL`](CONFIGURATION.md#stablerail_min_withdrawal).
- **When it is down:** bank deposits and withdrawals are unavailable.

More: [Bank deposits and withdrawals](#bank-deposits-and-withdrawals-stablerail).

## 16. Crypto deposits and withdrawals (1Liquidity)

- **What it is and why:** deposits from and withdrawals to outside crypto
  networks.
- **Direction:** both ways; 1Liquidity calls `POST /v1/callbacks/1l`.
- **Settings on this side:** [`ONELIQUIDITY_BASE_URL`](CONFIGURATION.md#oneliquidity_base_url),
  [`ONELIQUIDITY_TOKEN`](CONFIGURATION.md#oneliquidity_token) and the
  [other settings](CONFIGURATION.md#15-crypto-deposits-and-withdrawals-1liquidity).
- **When it is down:** these deposits and withdrawals wait.

## 17. Public Markets partners and price feed

- **What it is and why:** tokenized NGX stocks and FMDQ bonds are backed by
  units a **Custodian** holds; a **Dealing Member** buys and sells on the
  exchange; a **price feed** supplies prices.
- **Direction:** app-backend → partners (instructions); partners →
  app-backend (`/v1/custodian-partners/callbacks/...`,
  `/v1/dealing-member-partners/callbacks/execution`); app-backend → price
  feed.
- **Settings on this side:** partner details in Trovo Manager; the
  [price feed settings](CONFIGURATION.md#13-public-markets-tokenized-ngx-stocks-and-fmdq-bonds).
  Until partners' specifications are final they run as mocks.
- **When it is down:** orders wait in the instruction queue.

More: [Public Markets](#public-markets-tokenized-equities-and-bonds),
[PUBLIC_MARKETS.md](PUBLIC_MARKETS.md).

## 18. Naira price API

- **What it is and why:** converts between Naira and US dollars.
- **Direction:** app-backend → API (`{url}/api/convert/ngn-to-usd/{amount}`
  and `{url}/api/convert/usd-to-ngn/{amount}`).
- **Settings on this side:** [`CNGN_PRICE_API_URL`](CONFIGURATION.md#cngn_price_api_url).
- **When it is down:** Naira prices are unavailable.

## 19. Location lookup (ipapi)

- **What it is and why:** an IP-to-location service. The lookup code
  exists but nothing calls it today; its two settings are still checked
  at start.
- **Direction:** none today.
- **Settings on this side:** [`IPAPI_HOST`](CONFIGURATION.md#ipapi_host),
  [`IPAPI_KEY`](CONFIGURATION.md#ipapi_key) (any non-empty values).
- **When it is down:** nothing is affected.

## 20. Discord alerts

- **What it is and why:** chat alerts for errors (registration, payments,
  connection pool, low faucet balance).
- **Direction:** app-backend → Discord webhooks.
- **Settings on this side:** the [Discord settings](CONFIGURATION.md#24-discord-alerts);
  each falls back to a built-in webhook.

---

The sections below go into more depth on how the apps, tm-api and partners
authenticate and on each feature's flow.

## Authentication

`app-backend` has **three** distinct authentication schemes, chosen per
route based on who's expected to call it. There is no cookie-based session
and no classic OAuth bearer flow for end users — every scheme is stateless
and checked on every request.

### 1. Request-signing (`app-web` and `app-mobile`)

This is how the two end-user client apps authenticate almost every
request. There is no login token to attach — instead, **every request is
individually signed** with the calling wallet's private key. The client
sends four custom headers:

| Header | Meaning |
| --- | --- |
| `X-TW-PUBLIC-KEY` | The wallet address making the call |
| `X-TW-SIGNER` | The address whose private key actually signed this request (usually the same as `X-TW-PUBLIC-KEY`, but can differ for shared-access wallets where an approver signs on behalf of the wallet) |
| `X-TW-TIMESTAMP` | A current timestamp, included in what gets signed so a captured request can't be replayed later |
| `X-TW-SIGNATURE` | A base64-encoded EIP-191 `personal_sign` signature |

The signature covers `<full request path + query string> + <X-TW-SIGNER> +
<X-TW-TIMESTAMP>` (see `internal/middleware/authentication_middleware.go`'s
`authenticationChecks` and `internal/middleware/security_checks.go`'s
`SignHttp`/`VerifyHttpSignature` for the exact construction). The server
recomputes the expected signature from the request it received and rejects
the call with `401` if it doesn't match, or if `X-TW-SIGNER` fails basic
address-format validation.

Practically: whichever wallet SDK/library your client uses needs to be able
to sign an arbitrary string with the user's private key (never send the
private key itself to the server). The middleware enforcing this is
`AuthenticationMiddlewareUsingTimestamp()`, applied per-route in each
component's `controllers/main.go`.

In the Swagger UI and generated docs, routes that require this are marked
`@Security SignatureAuth`.

### 2. API key (white-label "service links", and a couple of server-to-server admin calls)

Routes under `/v1/servicelinks/...` and `/v1/trovo-api/...`, plus a small
number of admin-only endpoints elsewhere (e.g. P2P dispute resolution),
authenticate with a single static header instead:

```
X-TW-SERVICE-LINK-API-KEY: <the partner's API key>
```

The server looks this key up against the `service_links` table and rejects
the call with `401` if it's missing, unknown, or the service link is
suspended. See [Service links](#service-links-white-label-partner-integration)
below for how a partner gets one of these keys. The middleware is
`AuthenticationMiddlewareUsingAPIKey()`; these routes are marked
`@Security ServiceLinkApiKey` in the generated docs.

### 3. JWT bearer token (`tm-api`'s own admin routes)

A small set of routes under `/v1/trovo-manager/...` — internal admin
read/write endpoints `tm-api`'s staff UI calls — use a conventional JWT
bearer token instead:

```
Authorization: Bearer <jwt>
```

checked by `JwtTokenAuthMiddleware()`, signed with the `JWT_ACCESS_SECRET`
described in [CONFIGURATION.md](CONFIGURATION.md). These are marked
`@Security BearerAuth` in the generated docs. End-user clients and
white-label partners never use this scheme.

---

## How `tm-api` relates to `app-backend`

`tm-api` is Trovo's internal admin/catalog backend, and it integrates with
`app-backend` in **two different ways** at once:

1. **Direct database access, for admin/catalog data.** `tm-api` reads and
   writes tables like the curated-asset catalog, fee configuration, and
   service-link records directly against the **same physical Postgres
   database** `app-backend` uses. There's no API call involved for this —
   `app-backend`'s `assets`/`rates`/`servicelinks` components largely just
   *read* what `tm-api` wrote.
2. **HTTP calls into `app-backend`'s own API, for P2P transactional
   writes.** For creating orders, resolving disputes, and other P2P
   state-changing actions, `tm-api` calls this service's HTTP API (e.g.
   `POST /v1/p2p/disputes/{disputeID}/admin-resolve`, authenticated with
   `ServiceLinkApiKey`) rather than writing P2P tables directly.

**Why the split exists:** P2P orders and disputes involve real
in-flight, stateful business logic — escrow balances, order-state
transitions, notifications, on-chain settlement — the kind of thing that
needs to go through `app-backend`'s own validation and side effects (see
`internal/components/p2p/controllers/admin_handlers.go`'s own comment: *"app-backend
has no admin/staff user model of its own... this is a server-to-server
endpoint an admin-authorized tm-api call can reach"*). Writing directly to
the P2P tables from `tm-api` would risk leaving escrow state, balances, or
notifications out of sync with what `app-backend` believes is true.
Catalog/config data, by contrast, is comparatively inert — a curated
asset's metadata or a fee percentage doesn't have in-flight state that
another process could race with — so a direct, low-overhead DB write is
fine for that category.

**A consequence for you if you're adding a new admin feature:** if it's
read-only config/catalog data, writing directly to the DB from `tm-api`
(following existing patterns there) is consistent with how this already
works. If it's anything that changes P2P (or similar stateful,
in-flight) records, add a new `app-backend` endpoint instead of writing
those tables directly, and gate it behind `AuthenticationMiddlewareUsingAPIKey`
the same way `postAdminResolveDisputeHandler` does — see
`internal/components/p2p/controllers/admin_handlers.go` as the template.

---

## Service links (white-label partner integration)

A "service link" is how an external partner integrates Trovo Wallet
functionality into their own product under their own brand — a
white-label / embedded-wallet integration. The routes live in
`internal/components/servicelinks/controllers/` and are all under
`/v1/servicelinks/...` and `/v1/trovo-api/...`.

### Getting an API key

A service link is a row in the `service_links` table (owned by a specific
Trovo user account, `OwnerUsername`) carrying a generated `ApiKey`. There
is currently no self-serve signup endpoint for this — a service link
record is provisioned by Trovo (via `tm-api`'s admin tooling, direct DB
access as described above). Once provisioned, the partner receives their
`ApiKey` out of band and sends it as `X-TW-SERVICE-LINK-API-KEY` on every
request (see [Authentication](#2-api-key-white-label-service-links-and-a-couple-of-server-to-server-admin-calls) above).

### What a service link can do

Broadly, three categories of endpoint:

- **User-facing login/authorization handshakes** — `POST
  /v1/servicelinks/login/request/{targetUser}`,
  `POST /v1/servicelinks/authorize/request/{targetUser}`, and their
  `/verify/...` counterparts. These let the partner ask a Trovo user (by
  username) to approve the partner accessing their account, which the user
  confirms from inside the Trovo app itself (`POST
  /v1/users/servicelinks/login/approval/{targetUser}` and
  `.../authorize/approval/{targetUser}`, both signed by the user, not the
  partner). Approved sessions/authorizations are exchanged for short-lived
  tokens (`POST /v1/servicelinks/token/refresh`, `POST
  /v1/servicelinks/token/verify`) the partner then presents on subsequent
  calls.
- **Server-to-server "trovo-api" endpoints** — `/v1/trovo-api/users/...`,
  `/v1/trovo-api/assets/...`, `/v1/trovo-api/tokens/mint`, etc. These let
  an already-authorized partner onboard users, check balances, request
  payments, mint tokens, and manage their own asset listings, all
  authenticated purely by the API key (no per-request user approval).
- **Push/event delivery** — `POST
  /v1/servicelinks/{ownerUsername}/{targetUser}/push` and `POST
  /v1/servicelinks/events/request` let a partner push a notification to a
  user or register for event callbacks.

### The permission-flag model

Not every service link can do everything above. The `ServiceLink` row
(`internal/components/servicelinks/models/servicelink.go`) carries a set of
boolean-ish integer flags — `LoginPermission`, `PaymentPermission`,
`TokenInfoPermission`, `AuthorizationPermission`, `EventPermission`,
`AllowUserInfo`, `PushNotificationPermission`, `IncludePhoneNumbers`,
`IncludeUserBalances`, `TokenizedAssetAuthorizationPermission`,
`CreateUsersPermission`, `AllowReferralForRegisteredUsers`, `Verified`,
`Inactive`, `Suspended` — and each relevant handler checks the specific
flag(s) it needs before proceeding, returning a permission error if the
service link isn't authorized for that action. When you're orienting
yourself in this component for the first time, you don't need to memorize
every flag — just know that **each capability above is individually
gated**, so a partner integration is scoped to exactly what it was
provisioned for, and a new capability you add should follow the same
pattern (add a flag, check it in your new handler) rather than reusing an
existing, differently-scoped flag.

---

## Payment history: source vs. destination

`GET /v1/users/payments/:targetAddressForHistory` and
`GET /v1/trovo-api/users/payment-history/:walletAddress` (the service-link
equivalent) both return each row split into a **source** side (what left
`fromAddress`) and a **destination** side (what arrived at `toAddress`),
replacing the old flat `assetCode`/`contractAddress`/`amount` fields:

```json
{
  "sourceNetwork": "base",
  "sourceContractAddress": "0x...",
  "sourceAssetCode": "USDC",
  "sourceAmount": "100.0000000",
  "destinationNetwork": "base",
  "destinationContractAddress": "0x...",
  "destinationAssetCode": "USDC",
  "destinationAmount": "100.0000000"
}
```

For a plain payment the two sides are identical — that's not a
placeholder, it's the correct representation of "the same asset left and
arrived." A row where `sourceAssetCode`/`sourceNetwork` differ from
`destinationAssetCode`/`destinationNetwork` is, by definition, a swap.
`network` is currently always `"base"` on both sides; it's a real column
(not inferred/hardcoded client-side) so a future second network or
Base-native bridge doesn't require another breaking response-shape
change — just a different value here.

The read endpoints' query filters follow the same split:
`destinationAssetCode`/`destinationContractAddress`/`destinationAmount`
and `sourceAssetCode`/`sourceContractAddress`/`sourceAmount` (a
`min%max` range for the amount filters), replacing the old flat
`assetCode`/`contractAddress`/`amount` filter names.

Both of these endpoints, along with every other API-key-authenticated
service-link route (`/v1/servicelinks/...`, `/v1/trovo-api/...`),
`POST /v1/users/payment`, `POST /v1/shared-access/payment`,
`POST /v1/users/swap`, `POST /v1/shared-access/swap`, and
`GET /v1/users/:targetUser`, are rate-limited (see
[Rate limiting](CONFIGURATION.md#4-rate-limiting) in `CONFIGURATION.md`) —
a request over the limit gets `429 Too Many Requests` with a
`Retry-After` header. A service-link partner with unusually high (or low)
traffic needs can get a per-partner limit override instead of the
route's shared default — see CONFIGURATION.md's "Per-service-link
override", set from tm-api's Service Links admin page, not an env var.

---

## Live P2P updates over the websocket

`GET /v1/users/websocket/:identifier` is a websocket endpoint every
authenticated user can connect to for live, in-app updates. Once
connected, it pushes JSON messages shaped like:

```json
{ "stream": <payload>, "streamType": "<type>" }
```

The `streamType` you care about for live P2P updates is **`p2pEvent`** —
you'll get one of these whenever something happens on an order/offer/
escrow/dispute you're a party to (accepted, payment confirmed, dispute
opened, etc.), with the same title/body/data your push notification for
that event carries:

```json
{
  "streamType": "p2pEvent",
  "stream": {
    "title": "Order accepted",
    "body": "Your order was accepted by the merchant",
    "data": { "route": "orderDetail", "orderId": "..." }
  }
}
```

**How to connect:** open the websocket connection the same way you'd open
any other (`wss://.../v1/users/websocket/<your-wallet-address>`), send the
initial subscription handshake the endpoint's Swagger entry describes,
and then just listen — `p2pEvent` messages arrive automatically for as
long as the connection stays open. There is nothing else to subscribe to
per-order; every event for orders/offers you're involved in reaches this
one connection.

**Important: this is a live nudge, not a source of truth.** If your app
has no open connection at the moment an event happens (app closed,
connection dropped, etc.), that specific message is simply never
delivered — there's no catch-up/replay. Treat a `p2pEvent` message as a
hint to re-fetch the relevant order/offer from its REST endpoint, and
always fall back to the existing push notification (which you already
handle) plus normal polling/pull-to-refresh for anything you can't afford
to miss. This is the same tradeoff the admin panel's (`tm-api`) login
notification stream makes, and it's deliberate: the REST API and the
database are always the source of truth, this socket is purely a "check
now" signal.

If you're the first client integrating this (`app-web`/`app-mobile`
haven't wired it up yet as of this writing), the backend side is fully
built and tested — you just need to open the connection and handle
`streamType: "p2pEvent"` messages as described above. The same
connection also carries live payments (see "Live payments over the
websocket").

---

## Wallets and sending: Safe UserOperations

Every user wallet is a **Safe** (v1.4.1 with the ERC-4337 Safe4337Module)
owned by the user's key. Its address is fixed before it is deployed, so
it receives funds from registration on; its first send deploys it
("activates" it) in the same operation.

**Registration.** The app computes the user's address with wallet-core's
`primarySafeAddress(signer, "0")` and sends it as `X-TW-PUBLIC-KEY` (the
signer as `X-TW-SIGNER`). `POST /v1/users` refuses an address that is not
that Safe. Service-link onboarding may omit `publicKey`; it is derived.

**Sending (two calls, unchanged shape).**

1. The first call (e.g. `POST /v1/users/payment` without
   `transactionSignature`) validates the request, builds the operation
   and returns it as `transaction`: **base64 of the 32-byte SafeOp hash**,
   plus `messages` (including the network fee, and a note when the send
   also activates the wallet).
2. The app signs it exactly as before - `signBase64Txn` base64-decodes
   `transaction` and `personal_sign`s the bytes - and calls again with the
   same `transaction` and the base64 `transactionSignature`. The backend
   submits the stored operation (it is never rebuilt) to the bundler and
   returns its `userOpHash` as `transactionId`.

Operations expire (`WALLET_OPERATION_VALIDITY`, 10 minutes by default);
an expired or already-submitted `transaction` is refused. Shared wallets
with approvers get an approval request instead, and their operations
stay signable for `SHARED_WALLET_OPERATION_VALIDITY`.

**Gas.** The wallet pays its own gas: in the stablecoin its owner picked
(`User.gasFeeAsset`, one of the curated assets with `gasFeeEligible`)
through the paymaster when the wallet holds some, otherwise in ETH. A
payment that would leave too little to pay the fee in the same asset is
refused with `error-insufficient-for-network-fee`.

**Gas debt.** A new wallet's first operation cannot pre-pay its
stablecoin fee (the wallet approves the paymaster inside that operation),
so the paymaster takes it afterwards; if the operation spent the tokens
first, the fee is recorded on-chain as the wallet's **debt** and the
paymaster will not pay its gas again until it is settled. The backend
handles this: the next operation pays its gas in ETH and, when the wallet
holds enough of the token, also settles the debt (approve + `settleDebt`)
in the same operation. `messages` then says so ("This also pays … of
network fees this wallet still owed …"); when the wallet cannot settle
it yet, `messages` explains that fees are paid in ETH until it holds
enough.

**Fee records.** Once an operation is mined, the stablecoin gas the
paymaster collected (and any gas debt it settled) is recorded in
`fee_collections` with `fee_type` `GAS`, the token, amount, transaction
hash and the paymaster as `destination_wallet`. A charge the paymaster
could not collect is not recorded as a fee; it is alerted to Discord and
recorded once it is settled. ETH-paid gas goes to the bundler and is not
a platform fee.

**Sub-wallets** (`POST /v1/users/subwallet`, and the service-link
`POST /v1/trovo-api/users/subwallet`). Each sub-wallet is its own Safe,
deployed by the **primary** wallet, so only the primary wallet signs:

1. First call with `walletType` (0 normal, 1 issuing, 2 market making, 3
   bulk payment), `walletTag`, `walletDescription`. The backend picks the
   new Safe's address (the user's signer as owner, a random salt; types 2
   and 3 also get the platform's co-signer as a second owner, threshold
   1). An issuing wallet - or a request with any `linkedWalletAddress` -
   also gets a linked distribution Safe. The response carries `publicKey`
   (the new wallet), `linkedWalletAddress`, `transaction` (the primary
   wallet's operation that deploys them, seeds them with ETH and pays the
   creation fee - set in USD, charged in its stablecoin), `messages` and
   `feeAmount`/`feeCode` (the fee in that stablecoin). A `publicKey`
   sent by the client is ignored.
2. Second call with `transaction` and `primarySignature`. The wallets are
   recorded and the operation submitted; they are marked activated when it
   is mined. `subWalletSignature` / `linkedWalletSignature` /
   `channelAccountSignature` are no longer used (`subWalletMustSign` and
   `linkedWalletMustSign` are always 0).

If the primary wallet is not activated yet, the same operation activates
it too.

**Shared wallets.** Shared access is enforced by the wallet's Safe:
its owners are the wallet owner's key plus every APPROVER's key, and its
threshold is `numberOfApprovalsNeeded` (1 - the owner alone - when there
are no approvers). VIEW-ONLY and INITIATOR permissions stay in the
database. An issuing wallet's linked distribution wallet is deployed with
the issuing wallet's Safe enabled as a module, so the same operation that
changes the issuing wallet's owners changes the distribution wallet's too
(one set of signatures mirrors both).

- *Enable / modify / disable* keep their request shapes. While the owner
  alone controls the wallet, the owner signs the returned `transaction`
  (a wallet operation when approvers change, or a statement of the change
  when only view/initiator permissions change). Once the wallet has
  approvers, the initiator previews (`commit: 0`), then submits the change
  for approval (`commit: 1` with the same `transaction`).
- *Approvals* (`POST /v1/shared-access/approval/:ID`; `DELETE` rejects):
  the first call returns the
  request's `transaction` (base64 of the operation hash or statement);
  each approver signs it with `signBase64Txn`. A signature is checked
  against the approver's key, which must be one of the Safe's owners, and
  the Safe's current threshold decides how many approvals it takes. The
  last approval submits the operation with the collected signatures.
  Payments and crypto withdrawals from shared wallets work the same way.
- Operations waiting for approvers use their own EntryPoint nonce key and
  stay signable for `SHARED_WALLET_OPERATION_VALIDITY` (24h by default), so
  several can be pending on one wallet. Rejecting a request retires its
  operation.

**Gas-fee asset.** `GET /v1/users/settings/gas-fee-assets` lists the
curated assets users can pay network fees in (`gasFeeEligible`, set from
tm-web) and the user's current choice; `PUT /v1/users/settings/gas-fee-asset`
with `{"assetCode": "USDC"}` sets it (`""` clears it - fees in ETH). The
user profile carries it as `gasFeeAsset`.

**Other sends built the same way** (two calls, base64 SafeOp hash to sign):

- *Crypto withdrawal* (`POST /v1/crypto/withdrawals` and the
  shared-access variant): crypto deposits are minted as Trovo tokens, so
  a withdrawal **burns** `amountSubmitted` from the wallet (`burn(uint256)`
  of an OpenZeppelin `ERC20Burnable` token contract) and records the
  withdrawal request for payout. The first response carries `messages`.
- *Patron subscription*: paid from the primary wallet to the patron fee
  wallet (price + VAT). Payable in the dollar (`DOLLAR_ASSET`) or naira
  (`NAIRA_ASSET`, converted with the USD/cNGN rate) stablecoin; other
  payment assets need a DEX and are refused with
  `error-payment-asset-not-supported`.
- *P2P escrow deposit from the depositor's own wallet* goes through the
  payment flow; the deposit is recorded as pending under the operation's
  `userOpHash` and applied to the order only once it is mined, under the
  mined transaction's hash (the reconciliation sweep finishes deposits
  that take longer than the request).

Swaps fill market offers on the offer book - see "Swaps and market-making
offers" below. (Closed-group creation had no endpoint; its unused
transaction builder was removed.)

**What app-backend calls.** The bundler (`BUNDLER_URL`:
`eth_estimateUserOperationGas`, `eth_sendUserOperation`,
`eth_getUserOperationReceipt`) and the paymaster quote service
(`PAYMASTER_QUOTE_SERVICE_URL`, `POST /v1/quote` with `X-API-Key`; see
[`paymaster/INTEGRATION.md`](../paymaster/INTEGRATION.md)). A background
loop follows submitted operations to inclusion (`wallet_operations`
table) and marks wallets activated.

## Swaps and market-making offers

Swaps and market-making offers trade on `TrovoOfferBook`
(`OFFER_BOOK_ADDRESS`, [market/](../market/README.md)), the same contract
tokenized assets are sold through. Only tokens trade there: a swap or
offer with the native asset (ETH) answers `error-asset-not-swappable` /
`error-asset-not-tradable`.

**Market-making offers** (`POST /v1/users/trades`) place an offer with the
wallet's own funds, as one wallet operation the user signs:

- `SELL` *quantity* of the asset at *pricePerUnit* (in the currency):
  escrows the asset, priced in the currency.
- `BUY` *quantity* of the asset at *pricePerUnit*: escrows
  quantity x price of the currency, priced in the asset.

Like payments it takes two calls: without `transactionSignature` the
response carries `transaction` (sign it as for a payment) and messages;
send it back with `transactionSignature` to submit (shared wallets with
approvers send `commit: 1` and get an approval request). The offer keeps
the escrowed funds until it is filled or cancelled; proceeds go straight
to the wallet. `GET /v1/users/trades` lists the wallet's offers
(`{"offers": [...]}`), each with:

- `open`: whether it is still on the book;
- `remaining`: what it still sells, in the token it sells (whole tokens);
- `sold` and `received`: what has traded so far, in the token it sells and
  the token it buys; they stay correct after a cancellation, which sets
  `remaining` to 0;
- `fills`: how many trades it has had;
- `canceled` (1 once cancelled) and `blockchainOfferId` (set once its
  operation is mined).

`DELETE /v1/users/trades/:id` cancels one in the same two steps, returning
what is left.

**Swaps** (`POST /v1/users/swap`, `POST /v1/shared-access/swap`) buy the
destination token with the source token from the offers that sell it,
cheapest first (up to 8 offers per swap), in one wallet operation: the
swap fee and VAT transfers, allowing the book to take the payment, then
the fills. Every fill is authorized by the platform
(`OFFER_AUTHORIZER_PRIVATE_KEY`) for exactly that wallet and amount, and
the price is fixed when the swap is built: if an offer changes or runs out
before it is submitted, the operation reverts and nothing moves. The first
call returns `transaction`, `swappedEstimate` (what it receives),
`feeAmount`/`vatAmount` and messages (including any source amount too
small to buy more, which stays in the wallet); the second, with
`transactionSignature`, submits it. A swap larger than the book answers
`error-low-liquidity` with what is on offer.

**The order book, prices and charts** come from app-backend's index of the
offer book (`offer_book_*` tables, kept up to date by one instance at a
time): the order book of an asset against a currency has asks (offers
selling the asset; amount in the asset) and bids (offers buying it with
the currency; amount in the currency), both priced in the currency per
asset, as Stellar's order book had. Trade candles aggregate the fills.
The websockets `GET /v1/stream/orderbook` (`streamType: "orderBook"`) and
`GET /v1/stream/tradechart` (`"tradeChart"`) push them as they change.

**Public market data** (no sign-in; 120 requests a minute per caller, rate
limit key `market`) is what the apps' market screens show. Tokens are
named by their contract address; prices are in the counter token per one
base token, and amounts are decimal strings in whole tokens:

| Endpoint | Returns |
|---|---|
| `GET /v1/market/pairs` | `{"pairs": [...]}`: each pair's `base` and `counter` (`code`, `name`, `contractAddress`, `imageUrl`), `lastPrice`, `changePercent24h`, `high24h`, `low24h`, `baseVolume24h`, `counterVolume24h`, `trades24h`, `bestAsk`, `bestBid`, `lastTradeAt` (unix seconds). Bases are the curated tokens open on the book; counters are `NAIRA_ASSET` and `DOLLAR_ASSET`. |
| `GET /v1/market/orderbook?base=&counter=&limit=20` | `{"asks": [{price, amount}], "bids": [...]}`, amounts in the base token, best first (up to 100 levels) |
| `GET /v1/market/trades?base=&counter=&limit=50` | `{"trades": [{timestamp, price, baseAmount, counterAmount}]}`, newest first (up to 200) |
| `GET /v1/market/candles?base=&counter=&resolution=1h&limit=100` | `{"resolution", "candles": [{timestamp, open, high, low, close, baseVolume, counterVolume, trades}]}`, oldest first (up to 500); `resolution` is `15m`, `1h`, `4h`, `1d` or `1w` |

They answer 503 `error-market-not-configured` when `OFFER_BOOK_ADDRESS` is
not set, and 400 for an unknown pair or resolution.

## Live payments over the websocket

Besides `p2pEvent`, `GET /v1/users/websocket/:identifier` pushes the
user's wallets' new payments as they are recorded (by
payment-history-engine): `streamType` `"payment"`, or `"swap"` for swaps,
with a payment history record as `stream`. Send `cursor` (RFC 3339 time
or unix seconds) in the handshake to also get those since then. Like
`p2pEvent`, it is a live nudge; the payment history endpoints stay the
source of truth.

## Tokenized assets: token, issuing and distribution wallets, sale offer

On Base a tokenized asset has these addresses, and nothing in the API
treats one as another:

| Field | What it is | Used for |
|---|---|---|
| `contractAddress` | The asset's **token contract** (`TokenizedAsset`, see [`market/`](../market/README.md)) - the asset itself | Every balance, transfer, purchase, wallet authorization, curated-asset listing, subscription record, payout and deep link (`?action=tokenizedAsset&assetCode=...&contractAddress=...`) |
| `issuingWalletAddress` | The **issuing wallet**: a Safe sub-wallet of the issuing profile (`TOKENIZATION_ISSUING_PROFILE`), owned by the profile's key and the asset's minting approvers, with the minting threshold (approvers - 2). It owns the token contract, is its only minter, and is the seller of the sale offer | Minting, the sale offer, approving wallet authorizations |
| `marketMakingWallet` / `walletToHoldAssetsNotForSale` | The **distribution wallet**: the issuing wallet's linked Safe, with the same owners and the issuing Safe as a module | Supply not offered for sale; tokens returned by early exits |
| `offerBookOfferId` | The asset's sale offer on `TrovoOfferBook` (`OFFER_BOOK_ADDRESS`) | Every purchase; set once the mint is mined, which is when the asset counts as minted |

No platform key owns or signs for either Safe: the minting approvers do.

Lifecycle:

1. **Issuing and distribution wallets** - assigned when the asset's
   details are submitted: their addresses are computed from their owners
   (the issuing profile's key and the minting approvers) and recorded as
   sub-wallets of the issuing profile, with the approvers and initiators as
   their shared access. Nothing is deployed yet. While the token contract
   is not registered, changing the minting approvers moves the asset to new
   wallets; afterwards they are fixed (`409 error-minting-approvers-fixed`).
2. **Token contract** - operations deploy the asset's token with the
   issuing wallet as owner (Ownable) or `MINTER_ROLE` holder, the asset code
   as `symbol()` and zero supply, exposing `mint(address,uint256)`, then
   register it with `PUT /v1/trovo-manager/tokenization/contract/:tid
   {"contractAddress": "0x..."}` (minting initiators/approvers, until
   minted). The backend checks on-chain that the contract exists, its
   symbol is the asset code, `decimals()` works, `totalSupply()` is 0 and
   the issuing wallet can mint (`400 error-invalid-token-contract`
   otherwise). The offer book's owner must also list the token
   (`setTradable`) before minting; see [market/INTEGRATION.md](../market/INTEGRATION.md).
3. **Mint** - `POST /v1/trovo-manager/tokenization/mint/:tid` (status 3)
   builds the issuing wallet's first operation and creates a `TOKENIZE
   ASSET` approval request for it. The operation, which the minting
   approvers sign like any shared-wallet operation (base64 SafeOp hash):
   - deploys the issuing wallet (and, from it, the distribution wallet);
   - mints `maxNumberOfTokenAvailableForSale` to the issuing wallet, the
     rest less `feeInAsset` to the distribution wallet, and `feeInAsset`
     to `TOKENIZATION_FEE_WALLET`;
   - offers the amount for sale on the offer book at `pricePerToken`, in
     each tokenization payment stablecoin and in the country's internal
     balance token (all taken as 1:1 with the quote currency), with
     proceeds paid to `fundsHoldingWalletAddress` (required before
     minting).

   The issuing wallet pays the operation's gas, so it must hold ETH (or a
   stablecoin the issuing profile pays fees in) first:
   `error-insufficient-network-fee` names the address to fund. The request
   is refused once the token has supply (`error-already-minted`). When the
   threshold of approvers has signed, the operation is submitted; once it
   is mined, `offerBookOfferId` is recorded and the asset can go on sale
   (status 5 from its sales start date).
4. **Purchase with a stablecoin** - `SubscribeToTokenizedAsset` (the
   subscription endpoints, unchanged requests): the first call prices it -
   the most of the asset `amount` of `paymentAssetCode` (default CNGN)
   buys at the offer's price - after the platform checks (KYC, sale status,
   purchase cap), signs a fill authorization for exactly that purchase,
   and returns the buyer wallet's operation as `transaction` (approve the
   offer book for the payment, then fill). `messages` says what is paid
   and received; `swappedEstimate` is the amount of the asset. The second
   call with `transactionSignature` submits it; payment and delivery
   happen in one transaction or not at all. Shared wallets with approvers
   get an `ASSET SUBSCRIPTION` approval request. Sold out or too large:
   `error-no-liquidity` / `error-low-liquidity` (with the maximum).
5. **Purchase with fiat** - the first call creates the invoice and returns
   a statement as `transaction` for the buyer to sign (consent; their
   wallet sends nothing); the second stores the signature. When
   Flutterwave confirms payment, the internal balance token's minting Safe
   (`CountryConfig.internalTokenMinterSafe`, owned by
   `INTERNAL_BALANCE_ISSUING_SIGNERS`) mints the amount to itself and fills
   the offer with the buyer's wallet as recipient. If that fails (e.g. sold
   out after payment) the invoice stays `PENDING` and Discord is alerted
   for a manual delivery or refund.
6. **Early exit** - one operation of the holder's wallet returns the tokens
   to the distribution wallet; the payout is settled off-chain from the
   recorded early exit. Shared wallets get an approval request.

Purchases, early exits and payouts refuse assets whose token contract is
not registered (`error-token-contract-not-registered`) or not on sale yet
(`error-not-on-sale`).

---

## Account recovery: opt-in guardian

Trovo is non-custodial: users keep their secret key, and the platform holds
no key that controls their wallets. Account recovery is an **opt-in**
exception with a narrow power. It uses Candide's Social Recovery Module
(`RECOVERY_MODULE_ADDRESS`, deployed from [recovery/](../recovery/README.md)):
a user who turns it on makes the platform's recovery guardian
(`ACCOUNT_RECOVERY_GUARDIAN_SAFE` / `_SIGNERS`) the only guardian of their
wallets. The guardian can only *start* replacing a wallet's key; the
replacement takes effect after the module's recovery period, during which
the user is notified and can cancel it with their current key. It can never
transfer funds or execute anything else on the wallet.

**Covered wallets**: the user's primary wallet and their own sub-wallets
without approvers, once activated (deployed). Wallets with co-signers are
recovered by their co-signers instead (below), and adding approvers to a
covered wallet removes the guardian in the same operation.

### Turning it on and off (two steps, one signature per wallet)

`POST /v1/users/account/recovery` (on) and `DELETE /v1/users/account/recovery`
(off), signed with the user's key:

1. Without signatures the backend answers **202** with `transactions` (one
   operation per wallet, base64 hashes to sign), `wallets` (their aliases,
   same order) and `messages` to show. Turning it on enables the module and
   adds the guardian on each wallet not yet covered; the primary wallet's
   operation also pays the fee (`ACCOUNT_RECOVERY_FEE`, first time only).
   The primary wallet must be activated (`error primary account not yet
   activated`).
2. The app signs every transaction and sends the same body back with
   `transactionSignatures` (same order). **200** with `transactionId` (the
   operations' hashes, comma separated) when submitted. Turning it on again
   later covers wallets created since, without a fee.

`DELETE` answers 200 directly when no wallet still has the guardian.

### Recovering (the device with the new key)

1. Email OTP (`/v1/account/recovery/request-email-otp/:username`,
   `/v1/account/recovery/verify-email-otp/...`) and security answers, as before.
2. `POST /v1/users/account/recover` signed with the **new** key, with
   `newSignerAddress`, `emailOtp`, `securityAnswers` and `commit`:
   `commit: 0` checks and returns `messages` and the covered `wallets`
   (202); `commit: 1` has the guardian start replacing the old key with the
   new one on every covered wallet and answers **200** with `executeAfter`
   (when it takes effect) and `transactionId` (the guardian's transactions).
   The user is notified by push and email straight away.
3. `GET /v1/account/recovery/status/:username`, signed with the new key,
   returns `{status, executeAfter, completedAt}` for the recovery onto that
   key: `PENDING`, then `COMPLETED` (import the account with the new key),
   or `CANCELED` / `FAILED`.

After the period a background job (every 30 s, `process-account-recoveries`
lock) finalizes the recovery on each wallet, then switches the user's
signer to the new key. Wallets the user shares with approvers - their own
and those they approve on - each get a `REPLACE SIGNER` approval request
(an operation swapping the old key for the new one) that the other
co-signers approve like any other request. When the remaining co-signers
cannot reach the wallet's threshold, the request says so and Discord is
alerted.

### Cancelling (the device with the current key)

`POST /v1/users/account/recovery/cancel`, signed with the current key, in
the same two steps as turning recovery on: **202** with one `transactions`
entry per wallet being recovered, then **200** once the signed operations
are submitted. **404** `error-no-recovery-pending` when there is none. The
apps show it on the account recovery page, which recovery push
notifications (`route: accountRecovery`) open.

### Monitoring

The same job watches the module's `RecoveryExecuted` events: a recovery the
backend did not start (e.g. a misused guardian key) raises a Discord alert
and an urgent notification asking the owner to cancel it. A recovery the
guardian's transaction never confirmed is marked `FAILED` after 15 minutes
with an alert.

### Never-activated accounts

`POST /v1/users/inactive-account/recover` (users who never turned recovery
on) only works while none of the account's wallets is deployed or holds
any balance: it rebuilds the primary wallet for the new key, which changes
its address.

---

## Bank deposits and withdrawals (Stablerail)

Naira moves in and out of the platform as cNGN through Stablerail. All
routes are request-signed (like payments) and answer `503
error-fiat-unavailable` while Stablerail is not configured and enabled (see
`CONFIGURATION.md`).

- **Onboarding.** There is no route for it. When KYC level 1 completes with
  a BVN, the KYC callback starts the user's Stablerail onboarding; if
  Stablerail cannot be reached the callback saves a retry
  (`stablerail_onboard_user_retries`) that the `stablerail-onboarding-onramp`
  loop retries every 5 minutes (24 times, then a Discord alert; BVNs are
  masked in alerts). A repeated callback does not onboard twice.
  `GET /v1/users/stablerail/profile` tells the apps whether Stablerail is
  enabled, whether the user is onboarded (or the last onboarding's status)
  and the smallest withdrawal.
- **Deposit.** `POST /v1/users/stablerail/onrampcngn/:amount` (Naira)
  returns a virtual account to pay into; the cNGN is sent to the calling
  wallet once Stablerail sees the payment.
- **Withdraw.** `POST /v1/users/stablerail/withdraw` with `amount`,
  `accountNumber` (10 digits) and `bankCode` (from
  `GET /v1/users/stablerail/banks`) takes two calls, like a payment: the
  first returns `transaction` to sign and `messages` to show; the second
  (same body plus `transaction` and `transactionSignature`, or `commit` for a
  shared wallet, whose approvers then sign it as any other operation)
  submits a wallet operation (kind `BANK WITHDRAWAL`) sending the cNGN to the
  user's own Stablerail wallet (from `/getuserdetails`, cached). Once it is
  mined Stablerail is asked to pay the bank account (`/cngnofframp`); the
  `stablerail-offramp-status` loop retries that (10 times) and follows the
  payout to its end, and the user gets a push notification when it is paid
  or fails.
- **History.** `GET /v1/users/stablerail/withdrawals` lists the user's
  withdrawals with their status: `DEPOSITING` (transfer submitted),
  `DEPOSITED`, `REQUESTING`, `DEPOSIT_FAILED` (nothing left the wallet),
  `REQUEST_FAILED` (the cNGN is in the user's Stablerail wallet; support is
  alerted), then Stablerail's own statuses (`pending` ... `completed`,
  `failed`). The transfer itself also appears in normal payment history as a
  cNGN send to the Stablerail wallet.

Bank accounts are entered per withdrawal (Stablerail supports a fixed list
of banks); they are not the P2P payment methods, which accept any bank or
channel name.

## Proceeds payouts (dividends and interest)

A tokenized asset's proceeds are paid to its token holders by
`payout-engine`, a separate worker on this database (see
[payout-engine/README.md](../payout-engine/README.md)), driven from TM.
app-backend owns the tables (`proceed_payouts`,
`tokenized_asset_payout_schedules`, `proceed_payout_batches`,
`proceed_payout_approvals`, `payout_token_indexes`,
`payout_token_balances`, `payout_engine_states`) and shows holders their
payouts:

- `GET /v1/tokenization/payouts[?tokenizedAssetId=]` (request-signed, rate
  limited) lists the payouts to the user's wallets, newest first. It
  includes paid ones, with their transaction, and scheduled ones once a
  payout's holder schedule is locked (status `PENDING` / `QUEUED`), plus
  ones the engine could not pay yet (`FAILED`). Payouts that were
  cancelled or are still being prepared, and holders an admin excluded, are
  not listed. Each entry has the asset (code, name, token contract), the
  payout token and amount, the tokens held at the snapshot, the receiving
  wallet and its alias. app-mobile's dividend and interest screens use it.

The processing fee and its VAT on each payout are paid to the fee and VAT
wallets and recorded in `fee_collections` as `PROCEED_PAYOUT_FEE` /
`PROCEED_PAYOUT_FEE_VAT`.

## Public Markets (tokenized equities and bonds)

How the engine works is in [PUBLIC_MARKETS.md](PUBLIC_MARKETS.md). There
are three surfaces.

### The apps (`/v1/public-markets`)

The listing and prices are public (rate limited). Everything else uses
the request signature, like every other app route.

| Route | What it returns |
|---|---|
| `GET /v1/public-markets?market=&type=&search=` | Assets with price, day change, market session and whether they accept orders |
| `GET /v1/public-markets/assets/:assetCode` | One asset: statistics, custody chain, corporate actions |
| `GET /v1/public-markets/assets/:assetCode/prices?range=1D\|1W\|1M\|3M\|1Y\|All` | Chart points |
| `GET /v1/public-markets/assets/:assetCode/quote?side=buy&amount=` / `?side=sell&quantity=` | Fee, quantity or proceeds, and path (fast / netted / slow) |
| `POST /v1/public-markets/assets/:assetCode/buy` / `.../sell` | Two calls. With `walletAddress` and `amount` (or `quantity`), returns the quote and the wallet operation to sign. Called again with `transaction` and `transactionSignature`, it submits the operation and returns the order |
| `GET /v1/public-markets/portfolio` | Holdings, average cost, returns, income, open orders, activity |
| `GET /v1/public-markets/orders`, `GET /v1/public-markets/orders/:orderID` | Orders, and one order with its timeline |
| `GET /v1/public-markets/dividends?assetCode=` | Dividends and coupons: units at the record date, gross, withholding tax, net |

The order's status changes arrive as push notifications, and over the
user websocket as `publicMarketsEvent`. Trading needs KYC. Shared wallets
(wallets with approvers) cannot trade Public Markets.

### Exchanges (`/v1/trovo-api/public-markets`)

The exchange must be a verified service link that Trovo Manager has
onboarded for Public Markets. Every request carries:

- the API key, as `X-TW-SERVICE-LINK-API-KEY` or
  `Authorization: HMAC {apiKey}:{signature}`;
- `X-Trovotech-Timestamp` (Unix seconds, within 5 minutes);
- the signature, as `hex(HMAC-SHA256(signingSecret, timestamp + "." + body))`,
  in the `Authorization` header or `X-Trovotech-Signature`.

After a secret rotation, the previous secret is accepted until its expiry.
Every `POST` requires an `Idempotency-Key`:

- a repeat returns the original response with `Idempotent-Replayed: true`;
- the same key with a different body gets `422`;
- a key whose request is still running gets `409`;
- a `5xx` is not recorded, so it can be retried.

| Route | Purpose |
|---|---|
| `GET /assets`, `GET /assets/:assetCode` | Assets and reference prices |
| `POST /wallets` | Provision a customer wallet: `externalUserRef`, `legalName`, `taxIdentifier`, `residencyCountry`, `nationality`, `ndpaConsent`. Answers `201`, or `200` for a known `externalUserRef`. A missing field gets `400` with `field` |
| `GET /wallets/:walletID` | A wallet's holdings |
| `POST /creation` | `walletId`, `assetCode`, `amount` (NGN), `externalOrderRef`. Paid from the prefunded balance. Answers `202 pending` |
| `POST /redemption` | `walletId`, `assetCode`, `quantity` or `amount`, `externalOrderRef`. Proceeds go to the balance. Answers `202 pending` |
| `GET /orders/:orderID` | `pending`, `settled`, `rejected` or `failed` |
| `POST /confirmations` | Confirm a webhook (`eventId`, or `event` + `walletId`). Required for every `dividend.paid` |
| `GET /account` | Balance, deposit address (the treasury), the registered funding wallet and recent movements |
| `POST /deposits` | `txHash` of a CNGN transfer from the funding wallet to the treasury, which is verified on-chain and credited |

Webhooks go to the exchange's URL:

- events: `creation.settled`, `redemption.settled`, `order.rejected`,
  `dividend.paid`;
- headers: `X-Trovotech-Event`, `X-Trovotech-Delivery` (the event id),
  `X-Trovotech-Timestamp` and `X-Trovotech-Signature`, signed like the
  requests.

Failed deliveries are retried with backoff and then dead-lettered.
Trovo Manager can replay them.

### Custodian and Dealing Member callbacks

| Route | Sent by |
|---|---|
| `POST /v1/custodian-partners/callbacks/settlement` | Custodian: `instructionId`, `status` (`settlement_final` / `failed`), `settledQuantity`, `settlementDate`, `custodianReference` |
| `POST /v1/custodian-partners/callbacks/positions` | Custodian: `asOf`, `positions[]` (`assetCode`, `unitsHeld`) |
| `POST /v1/custodian-partners/callbacks/corporate-action` | Custodian: `assetCode`, `eventType`, `recordDate`, `payDate`, `amountPerUnit`, `currency` |
| `POST /v1/dealing-member-partners/callbacks/execution` | Dealing Member: `orderId`, `status` (`FILLED` / `PARTIAL` / `REJECTED`), `executedQuantity`, `executedPrice`, `executionTime` |

The partner names itself in `X-Partner-Code` and signs the body the same
way, with the secret from its credentials reference
(`X-Trovotech-Timestamp` + `X-Partner-Signature`). A redelivery answers
`{"status":"duplicate"}`. These payloads follow the integration
document's proposal and will be adjusted when the partners' specifications
are final. Until then, partners run as mocks.

## Swagger UI: the per-endpoint reference

Once the server is running, every documented endpoint — request
parameters, body shape, response codes, and which of the three auth
schemes above it requires — is browsable at:

```
http://localhost:8080/swagger/index.html
```

(replace the host/port with wherever the service is actually running). The
underlying spec is generated by [`swaggo/swag`](https://github.com/swaggo/swag)
from `@Summary`/`@Router`/etc. comments directly above each handler
function, regenerated with `swag init --parseDependency=false` (see
`docs/docs.go`'s header and [DEPLOYMENT.md](DEPLOYMENT.md)). If you're
integrating against a specific endpoint, always check its Swagger entry
first — it reflects the current handler code, not a hand-maintained
description that can drift.

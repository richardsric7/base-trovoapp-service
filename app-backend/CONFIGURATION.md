# Configuration

Every parameter `app-backend` reads. Each one says what it is, why it is
needed, whether you must set it, an example and how to get a real value.
Examples are never real secrets.

## How to set them

app-backend reads **environment variables**:

- **Locally**, from a `.env` file in `app-backend/`: `cp .env.example .env`,
  then edit it. It is loaded automatically at start.
- **In production**, from your container or hosting platform's environment
  settings (for example `docker run --env-file .env`, or the platform's
  "Environment variables" page). Inject secrets from your secrets manager
  (Vault) rather than baking them into an image.

### What happens when something is missing

At start, app-backend checks a list of required settings. If any is empty
it prints `Required environment variable is missing <NAME>` for each one
and stops (exit code 1). The settings marked **Required: Yes** below are on
that list, or stop the program in another way. `MIGRATE_ONLY=1` skips the
check (it only updates the database).

### A warning about secrets

Many settings below are private keys, mnemonics, salts or wallet addresses
that control real money. Use separate test values on test networks, keep
production values only in your secrets manager, and never commit them.
Some (marked **never change**) cannot be changed after launch without
breaking existing users.

### The required settings at a glance

| Area | Parameters |
|---|---|
| Database | [`DB_CONNECTION_STRING`](#db_connection_string), [`CDB_CONNECTION_STRING`](#cdb_connection_string) |
| Redis | [`ENABLE_CACHING`](#enable_caching), [`REDIS_HOST`](#redis_host), [`REDIS_PORT`](#redis_port) |
| Blockchain | [`BASE_RPC_URL`](#base_rpc_url), [`NATIVE_ASSET_CODE`](#native_asset_code), [`DOLLAR_ASSET`](#dollar_asset), [`BLOCKCHAIN_DATA_CACHE_LIFETIME`](#blockchain_data_cache_lifetime), [`CNGN_PRICE_API_URL`](#cngn_price_api_url) |
| Images | [`DEFAULT_ASSET_IMAGE_URL`](#default_asset_image_url), [`NATIVE_ASSET_IMAGE_URL`](#native_asset_image_url) |
| Salts and derived keys | [`VERIFICATION_CODE_SALT`](#verification_code_salt), [`ENCODER_SALT`](#encoder_salt), [`MNEMONIC_MARKET_MAKING`](#mnemonic_market_making), [`MARKET_MAKING_SALT`](#market_making_salt), [`MNEMONIC_BULK_PAYMENT`](#mnemonic_bulk_payment), [`BULK_PAYMENT_SALT`](#bulk_payment_salt) |
| Channel accounts | [`CHANNEL_ACCOUNTS`](#channel_accounts), [`CHANNEL_ACCOUNT_FUNDER`](#channel_account_funder), [`CHANNEL_ACCOUNT_MIN_COUNT`](#channel_account_min_count), [`CHECK_CHANNEL_ACCOUNT_BALANCE`](#check_channel_account_balance) |
| Fee wallets | [`VAT_WALLET`](#vat_wallet), [`PAYMENT_FEE_WALLET`](#payment_fee_wallet), [`SWAP_FEE_WALLET`](#swap_fee_wallet), [`SUBWALLET_CREATION_FEE_WALLET`](#subwallet_creation_fee_wallet), [`TOKENIZATION_APPLICATION_FEE_WALLET`](#tokenization_application_fee_wallet), [`TOKENIZATION_FEE_WALLET`](#tokenization_fee_wallet), [`CLOSED_GROUP_FEE_WALLET`](#closed_group_fee_wallet), [`ACCOUNT_RECOVERY_FEE_WALLET`](#account_recovery_fee_wallet), [`PATRON_FEE_WALLET`](#patron_fee_wallet), [`SHARED_ACCESS_PAYMENT_FEE_WALLET`](#shared_access_payment_fee_wallet), [`SWAP_FEE_ENABLED`](#swap_fee_enabled), [`MARKET_MAKING_FEE_ENABLED`](#market_making_fee_enabled) |
| Email | [`MAILGUN_PRIVATE_API_KEY`](#mailgun_private_api_key), [`MAILGUN_DOMAIN`](#mailgun_domain), [`ENABLE_EMAIL_VALIDATION`](#enable_email_validation) |
| Firebase and links | [`GC`](#gc), [`GOOGLE_PROJECT_ID`](#google_project_id), [`DYNAMIC_LINKS_DOMAIN_PREFIX`](#dynamic_links_domain_prefix), [`DYNAMIC_LINKS_ANDROID_PACKAGE_NAME`](#dynamic_links_android_package_name), [`DYNAMIC_LINKS_IOS_BUNDLE_ID`](#dynamic_links_ios_bundle_id), [`DYNAMIC_LINKS_FALLBACK_BASE_URL`](#dynamic_links_fallback_base_url), [`WALLET_DOMAIN`](#wallet_domain) |
| Login tokens | [`JWT_ACCESS_SECRET`](#jwt_access_secret), [`JWT_TOKEN_EXPIRY`](#jwt_token_expiry), [`JWT_REFRESH_TOKEN_EXPIRY`](#jwt_refresh_token_expiry) |
| Location lookup | [`IPAPI_HOST`](#ipapi_host), [`IPAPI_KEY`](#ipapi_key) |

Strongly recommended for a working wallet: [`BUNDLER_URL`](#bundler_url)
(without it no wallet can send anything) and [`BASE_CHAIN_ID`](#base_chain_id).

---

## 1. Server

### `PORT`

- **What it is:** The port the web server listens on.
- **Why it's needed:** The apps, tm-api and partners reach app-backend here.
- **Required:** No, default `8080`.
- **Example:** `8080`
- **How to get it:** Any free port; your hosting platform may set it.

### `GIN_MODE`

- **What it is:** The web framework's mode.
- **Why it's needed:** `release` turns off verbose debug logging.
- **Required:** No (development mode when unset).
- **Example:** `release`
- **How to get it:** `release` in every deployed environment.

### `ORGANISATION`

- **What it is:** A short organization name returned by the root endpoint.
- **Why it's needed:** Only informational.
- **Required:** No.
- **Example:** `trovotech`
- **How to get it:** Any short name.

### `MIGRATE_ONLY`

- **What it is:** `1` makes the program update the database tables and exit,
  without serving requests.
- **Why it's needed:** To update the database as a separate step before
  starting new servers (see [DEPLOYMENT.md](DEPLOYMENT.md)). It skips the
  required-settings check, since it only needs the database.
- **Required:** No.
- **Example:** `1`
- **How to get it:** Set `1` only for the one-off migration step.

### `DB_AUTOMIGRATE`

- **What it is:** Whether to create and update database tables at start.
- **Why it's needed:** New releases add tables and columns. With `0` the
  server starts faster but relies on a separate `MIGRATE_ONLY=1` step.
  Several copies starting together take turns (a database lock), so they do
  not conflict.
- **Required:** No (on unless set to `0`).
- **Example:** `1`
- **How to get it:** `1` for a single server or local development; `0` on
  servers when you run a separate migration step.

### `SHUTDOWN_GRACE_PERIOD`

- **What it is:** How long a stopping server waits for running requests,
  background jobs and platform transactions to finish.
- **Why it's needed:** Stopping mid-transaction could leave work half done.
- **Required:** No, default `60s`.
- **Example:** `90s`
- **How to get it:** A little less than your platform's stop timeout (for
  example Kubernetes `terminationGracePeriodSeconds`).

### `ENABLE_AUTH_MIDDLEWARE`

- **What it is:** `0` turns off all request signature and API-key checks.
- **Why it's needed:** Only for a throwaway local debugging session. **With
  `0`, anyone can act as any user.**
- **Required:** No (checks are on unless set to `0`).
- **Example:** (leave unset)
- **How to get it:** Never set it outside your own machine.

### `LOG_IP_ADDRESS`

- **What it is:** `1` logs callers' IP addresses on some requests.
- **Why it's needed:** Only while investigating a problem.
- **Required:** No.
- **Example:** `0`
- **How to get it:** Leave unset normally.

### `LOG_TARGET_USER`

- **What it is:** A username whose requests get extra logging.
- **Why it's needed:** To trace one user's problem without verbose logging
  for everyone.
- **Required:** No.
- **Example:** `ada`
- **How to get it:** The username you are investigating; unset otherwise.

### `LOG_TARGET_USER_PK`

- **What it is:** A wallet address whose requests get extra logging.
- **Why it's needed:** As `LOG_TARGET_USER`, by address.
- **Required:** No.
- **Example:** `0x1111111111111111111111111111111111111111`
- **How to get it:** The address you are investigating; unset otherwise.

## 2. Database

### `DB_CONNECTION_STRING`

- **What it is:** The address and login of the main Postgres database.
- **Why it's needed:** Everything (users, wallets, payments, assets,
  payouts, Public Markets) is stored there. tm-api, payout-engine and
  payment-history-engine use the same database. Without it app-backend does
  not start.
- **Required:** Yes.
- **Example:** `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`
  (or a file path such as `trovo.sqlite` with `DB_TYPE=sqlite`)
- **How to get it:** Create a Postgres 14+ database and user (a managed
  service such as DigitalOcean, AWS RDS or Google Cloud SQL shows the
  connection details), and write them in the form above. Locally:
  `docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=trovo -e POSTGRES_USER=trovo -e POSTGRES_DB=trovo postgres:16`
  and `host=localhost user=trovo password=trovo dbname=trovo port=5432 sslmode=disable`.

### `DB_TYPE`

- **What it is:** `postgres` or `sqlite`.
- **Why it's needed:** `sqlite` lets you run locally against a file with no
  database server.
- **Required:** No, default `postgres`.
- **Example:** `postgres`
- **How to get it:** Leave unset in real deployments.

### `DB_MAX_OPEN_CONNECTIONS`

- **What it is:** The most database connections one server keeps open.
- **Why it's needed:** The database's connection limit is shared by every
  copy of every service.
- **Required:** No, default `50`.
- **Example:** `50`
- **How to get it:** The database's `max_connections`, divided between all
  services and copies.

### `DB_MAX_IDLE_CONNECTIONS`

- **What it is:** How many unused connections are kept open for reuse.
- **Why it's needed:** Reusing connections is faster.
- **Required:** No, default `50`.
- **Example:** `50`
- **How to get it:** Usually the same as `DB_MAX_OPEN_CONNECTIONS`.

### `SQLITE_DB_PATH`

- **What it is:** A file path for a separate SQLite helper used in some
  local and test paths.
- **Why it's needed:** Local development only.
- **Required:** No, default `dbs/bantupay.sqlite`.
- **Example:** `dbs/local.sqlite`
- **How to get it:** Leave unset.

### `CDB_CONNECTION_STRING`

- **What it is:** The address of the tracking database ("RoachDB", a
  CockroachDB or Postgres database).
- **Why it's needed:** When a user creates a wallet, app-backend records it
  there so payment-history-engine starts watching it for transfers.
- **Required:** Yes.
- **Example:** `postgresql://root@roach.internal:26257/payment_history?sslmode=require`
- **How to get it:** A second database (CockroachDB Cloud shows the string
  on the cluster's **Connect** page; a second Postgres database also works).
  payment-history-engine must use the same value.

## 3. Redis

### `ENABLE_CACHING`

- **What it is:** `1` connects to Redis; `0` runs without it.
- **Why it's needed:** Redis caches frequent responses, holds the rate
  limits, and carries live updates (P2P, login) between copies of
  app-backend. Without it everything still works on a single server, but
  slower, with no rate limiting, and live updates only reach users on the
  same copy.
- **Required:** Yes (`0` or `1`).
- **Example:** `1`
- **How to get it:** `1` in production; `0` for simple local development.

### `REDIS_HOST`

- **What it is:** The Redis server's host name.
- **Why it's needed:** Where Redis is. It must be set even with
  `ENABLE_CACHING=0` (the start-up check requires it).
- **Required:** Yes.
- **Example:** `redis.internal`
- **How to get it:** Your Redis server (a managed Redis shows its host; or
  `docker run -d -p 6379:6379 redis:7` and `localhost`).

### `REDIS_PORT`

- **What it is:** The Redis server's port.
- **Why it's needed:** To connect.
- **Required:** Yes.
- **Example:** `6379`
- **How to get it:** `6379` unless your provider says otherwise.

### `REDIS_PASSWORD`

- **What it is:** The Redis password.
- **Why it's needed:** A protected Redis refuses connections without it.
- **Required:** Only if your Redis has one.
- **Example:** `change-me`
- **How to get it:** Your Redis provider's dashboard.

### `CACHING_PARAMETER`

- **What it is:** Text added to every cache key.
- **Why it's needed:** Changing it makes all old cache entries unused at
  once.
- **Required:** No.
- **Example:** `v2`
- **How to get it:** Any short text; change it to clear the cache.

## 4. Rate limiting

Each protected route allows a number of requests per caller per minute,
counted in Redis so the limit is the same across all copies. Protected
routes include payment history, payments, swaps, `GET /v1/users/:targetUser`
and every partner (service-link) route under `/v1/servicelinks/...` and
`/v1/trovo-api/...`. Without Redis the limits are off.

A partner can also be given its own limit, which wins over everything
below: the service link's `rateLimitPerMinute`, set in Trovo Manager's
**Service Links** page (0 = no override).

### `RATE_LIMIT_ENABLED`

- **What it is:** `0` turns all rate limits off.
- **Why it's needed:** Limits protect the platform from abuse and runaway
  partners; turn them off only for load tests.
- **Required:** No (on when Redis is on).
- **Example:** `1`
- **How to get it:** Leave unset in production.

### `RATE_LIMIT_REQUESTS_PER_MINUTE`

- **What it is:** One limit for every protected route, replacing each
  route's own default.
- **Why it's needed:** To loosen or tighten all limits at once. Each route
  otherwise has its own default in the code: about 60 per minute for reads,
  30 for login and authorization handshakes, 10 to 20 for actions such as
  onboarding, KYC and minting.
- **Required:** No (each route's default).
- **Example:** `30`
- **How to get it:** Leave unset unless the defaults do not fit.

### `RATE_LIMIT_<KEY>_PER_MINUTE`

- **What it is:** The limit for one route, by its key in capitals with
  `-` replaced by `_`: the `swap` key is `RATE_LIMIT_SWAP_PER_MINUTE`,
  `payment-history` is `RATE_LIMIT_PAYMENT_HISTORY_PER_MINUTE`.
- **Why it's needed:** When one route needs a different limit. It wins
  over `RATE_LIMIT_REQUESTS_PER_MINUTE`.
- **Required:** No.
- **Example:** `RATE_LIMIT_SWAP_PER_MINUTE=5`
- **How to get it:** Find the route's key in its
  `middleware.RateLimitMiddleware(gc, "<key>", ...)` call.

## 5. Blockchain

### `BASE_RPC_URL`

- **What it is:** The web address of a Base blockchain node.
- **Why it's needed:** Every balance, transfer, wallet operation and
  contract call goes through it.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
- **How to get it:** Create an app for Base Mainnet (or Base Sepolia for
  testing) with a node provider such as Alchemy, Infura or QuickNode and
  copy its HTTPS URL. `https://sepolia.base.org` works for testing.

### `BASE_CHAIN_ID`

- **What it is:** The network number.
- **Why it's needed:** Transactions are signed for one network.
- **Required:** No, but set it: the default is `84532` (Base Sepolia).
- **Example:** `8453`
- **How to get it:** `8453` for Base Mainnet, `84532` for Base Sepolia.

### `RPC_TIMEOUT`

- **What it is:** The longest one request to the node may take.
- **Why it's needed:** A stuck node must not freeze the server. After 3
  failures in a row, node calls are refused for 10 seconds at a time until
  one succeeds.
- **Required:** No, default `30s`.
- **Example:** `30s`
- **How to get it:** Keep the default.

### `NATIVE_ASSET_CODE`

- **What it is:** The code of the network's own currency.
- **Why it's needed:** Balances and payments treat this asset specially (it
  has no token contract).
- **Required:** Yes.
- **Example:** `ETH`
- **How to get it:** `ETH` on Base.

### `DOLLAR_ASSET`

- **What it is:** The US-dollar stablecoin used for dollar prices, written
  `CODE:contractAddress`.
- **Why it's needed:** Prices and balances are shown in dollars through it,
  and fees in USD are paid in it.
- **Required:** Yes.
- **Example:** `USDC:0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913` (USDC on Base Mainnet)
- **How to get it:** The stablecoin's code and its contract address on your
  network (Circle's documentation lists USDC's addresses).

### `USE_ASSET_FOR_NATIVE_PRICE`

- **What it is:** The asset whose market price stands in for ETH's, written
  `CODE:contractAddress`.
- **Why it's needed:** To price ETH in the order book summary and charts.
- **Required:** No.
- **Example:** `USDC:0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913`
- **How to get it:** Usually the same as `DOLLAR_ASSET`.

### `NAIRA_ASSET`

- **What it is:** The Naira stablecoin, written `CODE:contractAddress`.
- **Why it's needed:** Naira prices and conversions use it.
- **Required:** Only with `ENABLE_NAIRA_ASSET_BY_DEFAULT=1`.
- **Example:** `CNGN:0xC0FFEE0000000000000000000000000000000001`
- **How to get it:** cNGN's contract address from its issuer's
  documentation.

### `ENABLE_NAIRA_ASSET_BY_DEFAULT`

- **What it is:** `1` makes the start-up check require `NAIRA_ASSET`.
- **Why it's needed:** A safety check when Naira is part of your product.
- **Required:** No.
- **Example:** `1`
- **How to get it:** `1` if you use cNGN.

### `USDB_B20_TOKEN_ADDRESS`

- **What it is:** A dollar token's contract address used by the order book
  summary when a request does not name the buying token.
- **Why it's needed:** A default for that one query.
- **Required:** No.
- **Example:** `0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913`
- **How to get it:** The same contract as `DOLLAR_ASSET`.

### `TROV_ASSET_CONTRACT_ADDRESS`

- **What it is:** The contract of the TROV token.
- **Why it's needed:** Some tokenization steps check an applicant's TROV
  balance.
- **Required:** No.
- **Example:** `0xC0FFEE0000000000000000000000000000000003`
- **How to get it:** The TROV token's address on your network, if you use
  it.

### `CNGN_PRICE_API_URL`

- **What it is:** The address of a price service that converts between
  Naira and US dollars.
- **Why it's needed:** Naira amounts are priced in dollars and back.
  app-backend calls `{url}/api/convert/ngn-to-usd/{amount}` and
  `{url}/api/convert/usd-to-ngn/{amount}`.
- **Required:** Yes.
- **Example:** `https://prices.trovo.example.com`
- **How to get it:** Your Naira price provider's base address (it must
  answer those two paths).

### `BLOCKCHAIN_DATA_CACHE_LIFETIME`

- **What it is:** How long, in seconds, some blockchain data is cached.
- **Why it's needed:** Avoids reading data that rarely changes again and
  again.
- **Required:** Yes (the start-up check requires it).
- **Example:** `94608000` (about 3 years)
- **How to get it:** `94608000` unless you have a reason to change it.

### `BLOCKCHAIN_SWAP_DESTINATION_MIN`

- **What it is:** The smallest amount a swap may produce.
- **Why it's needed:** Stops dust-sized swaps.
- **Required:** No (`0`).
- **Example:** `0.0001`
- **How to get it:** A small amount meaningful for your assets.

### `BLOCKCHAIN_NETWORK_PASSPHRASE`, `BLOCKCHAIN_BASE_RESERVE`

- **What they are:** Leftovers from the platform's previous blockchain
  (Stellar).
- **Why they're needed:** They aren't, on Base.
- **Required:** No.
- **Example:** (leave unset)
- **How to get them:** Leave unset.

### `DEFAULT_ASSET_IMAGE_URL`

- **What it is:** The picture shown for an asset that has none.
- **Why it's needed:** So the apps never show a broken image.
- **Required:** Yes.
- **Example:** `https://cdn.trovo.example.com/assets/default.png`
- **How to get it:** Upload an image anywhere public (your website, cloud
  storage) and use its address.

### `NATIVE_ASSET_IMAGE_URL`

- **What it is:** The picture shown for ETH.
- **Why it's needed:** ETH has no token entry to take an image from.
- **Required:** Yes.
- **Example:** `https://cdn.trovo.example.com/assets/eth.png`
- **How to get it:** As above.

## 6. Wallets: Safe accounts, bundler and paymaster

Every user wallet is a Safe (v1.4.1 with the ERC-4337 Safe4337Module) owned
by the user's key. Every send is an operation app-backend prepares, the
user's app signs and app-backend submits to the **bundler**. The wallet pays
its own network fee, in ETH or (through the paymaster) in the stablecoin
the user chose.

### `BUNDLER_URL`

- **What it is:** The address of your ERC-4337 bundler (EntryPoint v0.7).
- **Why it's needed:** Every wallet operation is priced by and submitted to
  it. Without it no wallet can send anything.
- **Required:** Yes, for a working wallet (not on the start-up list).
- **Example:** `http://bundler.internal:4337/rpc`
- **How to get it:** Run a bundler (for example Rundler, Alto or Skandha)
  for the same network as `BASE_RPC_URL`, configured for EntryPoint
  `0x0000000071727De22E5E9d8BAf0edAc6f37da032`, and use its address.

### `PAYMASTER_ADDRESS`

- **What it is:** The address of the deployed `TrovoTokenPaymaster`.
- **Why it's needed:** Lets users pay network fees in stablecoins. Empty
  means everyone pays in ETH.
- **Required:** No.
- **Example:** `0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9`
- **How to get it:** From `paymaster/contracts/deployments/<chainId>.json`
  ([paymaster/DEPLOYMENT.md](../paymaster/DEPLOYMENT.md)).

### `PAYMASTER_QUOTE_SERVICE_URL`

- **What it is:** The address of the paymaster's quote service.
- **Why it's needed:** It prices each stablecoin-paid operation. If it
  cannot, the operation falls back to ETH.
- **Required:** With `PAYMASTER_ADDRESS`.
- **Example:** `http://quote-service.internal:8090`
- **How to get it:** Where you run the quote service.

### `PAYMASTER_QUOTE_SERVICE_API_KEY`

- **What it is:** The key app-backend sends to the quote service
  (`X-API-Key` header).
- **Why it's needed:** The quote service only answers known callers.
- **Required:** With `PAYMASTER_ADDRESS`.
- **Example:** `3f9c1d5e8a7b4c2d9e0f1a2b3c4d5e6f`
- **How to get it:** One of the quote service's `API_KEYS`
  ([paymaster/CONFIGURATION.md](../paymaster/CONFIGURATION.md#api_keys)).

### `WALLET_OPERATION_VALIDITY`

- **What it is:** How long a single-signer wallet has to sign a prepared
  operation.
- **Why it's needed:** Prices and quotes go stale; an unsigned operation
  expires and must be started again.
- **Required:** No, default `10m`.
- **Example:** `10m`
- **How to get it:** Keep the default.

### `SHARED_WALLET_OPERATION_VALIDITY`

- **What it is:** How long a shared wallet's approvers have to sign.
- **Why it's needed:** Approvers need time; it must be at most the quote
  service's `QUOTE_MAX_VALIDITY`.
- **Required:** No, default `24h`.
- **Example:** `24h`
- **How to get it:** Keep the default unless approvals take longer.

### `WALLET_OPERATION_GAS_BUFFER_PERCENT`

- **What it is:** Extra percentage added to the bundler's fee estimates.
- **Why it's needed:** Estimates can be slightly low. Unused gas is not
  charged, so this only raises the maximum shown to the user.
- **Required:** No, default `20`.
- **Example:** `20`
- **How to get it:** Keep the default.

### Safe and EntryPoint contract addresses

`ENTRYPOINT_ADDRESS`, `SAFE_PROXY_FACTORY_ADDRESS`, `SAFE_SINGLETON_ADDRESS`,
`SAFE_MODULE_SETUP_ADDRESS`, `SAFE_4337_MODULE_ADDRESS`,
`SAFE_MULTISEND_CALL_ONLY_ADDRESS`, `SAFE_FALLBACK_HANDLER_ADDRESS`:

- **What they are:** The addresses of the standard contracts wallets are
  built from.
- **Why they're needed:** Only a private test chain has different
  addresses. The defaults are the official deployments on Base and Base
  Sepolia: EntryPoint v0.7 `0x0000000071727De22E5E9d8BAf0edAc6f37da032`,
  SafeModuleSetup v0.3.0 `0x2dd68b007B46fBe91B9A7c3EDa5A7a1063cB5b47`,
  Safe4337Module v0.3.0 `0x75cf11467937ce3F2f357CE24ffc3DBF8fD5c226`,
  MultiSendCallOnly v1.4.1 `0x9641d764fc13c8B624c04430C7356C1C7C8102e2`,
  CompatibilityFallbackHandler v1.4.1 `0xfd0732Dc9E303f09fCEf3a7388Ad10A83459Ec99`,
  and Safe's v1.4.1 SafeProxyFactory and SafeL2.
- **Required:** No. **Never change them in production**: every user's
  wallet address is computed from them, and wallet-core in the apps uses
  the defaults.
- **Example:** (leave unset)
- **How to get them:** Leave unset; on a local chain, use the addresses the
  paymaster local stack prints.

## 7. Salts and derived keys (never change)

These values are mixed into hashes and key derivations. **Changing any of
them after launch breaks existing users**: their security answers stop
matching, or their market-making and bulk-payment wallets get different
keys.

### `VERIFICATION_CODE_SALT`

- **What it is:** A secret mixed into email and SMS verification codes.
- **Why it's needed:** Makes codes impossible to predict.
- **Required:** Yes.
- **Example:** `9f8e7d6c5b4a39281706f5e4d3c2b1a0`
- **How to get it:** `openssl rand -hex 16`.

### `ENCODER_SALT`

- **What it is:** A secret mixed into the hash of users' security-question
  answers.
- **Why it's needed:** Answers are stored only as hashes. **Never change
  it**: existing answers would stop matching and users could not recover
  their accounts.
- **Required:** Yes.
- **Example:** `1a2b3c4d5e6f708192a3b4c5d6e7f809`
- **How to get it:** `openssl rand -hex 16`, once.

### `MNEMONIC_MARKET_MAKING`

- **What it is:** A secret recovery phrase (12 or 24 words) used to derive
  each user's market-making wallet signing key.
- **Why it's needed:** Market-making sub-wallets are signed by a key
  derived from this, the salt below, the username and the wallet. **Never
  change it**: those wallets would get different keys.
- **Required:** Yes.
- **Example:** `test test test test test test test test test test test junk`
  (a well-known test phrase; never use it for real)
- **How to get it:** Generate a new phrase once (any BIP-39 tool, for
  example `cast wallet new-mnemonic`) and keep it in your secrets manager.

### `MARKET_MAKING_SALT`

- **What it is:** A secret mixed into the market-making key derivation.
- **Why it's needed:** As above. **Never change it.**
- **Required:** Yes.
- **Example:** `5e4d3c2b1a0f9e8d7c6b5a4938271605`
- **How to get it:** `openssl rand -hex 16`, once.

### `MNEMONIC_BULK_PAYMENT`

- **What it is:** A secret recovery phrase used to derive each user's
  bulk-payment wallet signing key.
- **Why it's needed:** As for market making. **Never change it.**
- **Required:** Yes.
- **Example:** `test test test test test test test test test test test junk`
  (test phrase)
- **How to get it:** Generate a separate phrase once and keep it in your
  secrets manager.

### `BULK_PAYMENT_SALT`

- **What it is:** A secret mixed into the bulk-payment key derivation.
- **Why it's needed:** As above. **Never change it.**
- **Required:** Yes.
- **Example:** `0f1e2d3c4b5a69788796a5b4c3d2e1f0`
- **How to get it:** `openssl rand -hex 16`, once.

## 8. Channel accounts (legacy)

Channel accounts are a leftover from the platform's previous blockchain,
where the platform paid fees for users. On Base, wallets pay their own
fees and no new transaction takes a channel account, but the settings are
still checked at start and the start-up routine can still fund and create
them. To avoid spending ETH on them, list one key, set
`CHANNEL_ACCOUNT_MIN_COUNT` to `0` (or no more than the keys listed) and
`CHECK_CHANNEL_ACCOUNT_BALANCE` to `0`.

### `CHANNEL_ACCOUNTS`

- **What it is:** Private keys of channel accounts, separated by commas.
- **Why it's needed:** Required at start; not used for new transactions on
  Base.
- **Required:** Yes.
- **Example:** `0x8b3a350cf5c34c9194ca85829a2df0ec3153be0318b5e2d3348e872092edffba`
  (a test key)
- **How to get it:** Generate one dedicated key
  (`node -e "console.log(require('ethers').Wallet.createRandom().privateKey)"`)
  that holds no funds.

### `CHANNEL_ACCOUNT_FUNDER`

- **What it is:** The private key of the wallet that funds channel accounts.
- **Why it's needed:** It must be a valid private key, or app-backend
  crashes at start. With the settings above it sends nothing.
- **Required:** Yes (in practice).
- **Example:** `0x47e179ec197488593b187f80a00eb0da91f1b9d0b13f8733639f19c30a34926a` (a test key)
- **How to get it:** Generate a dedicated key; fund it only if you want
  channel accounts funded.

### `CHANNEL_ACCOUNT_MIN_COUNT`

- **What it is:** The number of channel accounts to keep. If fewer are
  listed, new ones are generated, funded, and their keys emailed to
  `CHANNEL_ACCOUNT_RECEIPIENT`.
- **Why it's needed:** Legacy. It must be a number, or app-backend crashes
  at start.
- **Required:** Yes (in practice).
- **Example:** `0`
- **How to get it:** `0` on Base.

### `CHECK_CHANNEL_ACCOUNT_BALANCE`

- **What it is:** `1` checks channel accounts' balances at start and tops
  them up; `0` skips that.
- **Why it's needed:** Legacy; `0` avoids spending ETH.
- **Required:** Yes (`0` or `1`).
- **Example:** `0`
- **How to get it:** `0` on Base.

### `CHANNEL_ACCOUNT_FUNDING_AMOUNT`

- **What it is:** How much ETH each channel account is funded with.
- **Why it's needed:** Only when funding happens (see above).
- **Required:** No.
- **Example:** `0.001`
- **How to get it:** Leave unset on Base.

### `CHANNEL_ACCOUNT_MIN_BALANCE`

- **What it is:** The balance below which a channel account is topped up.
- **Why it's needed:** Only with `CHECK_CHANNEL_ACCOUNT_BALANCE=1`.
- **Required:** With `CHECK_CHANNEL_ACCOUNT_BALANCE=1`.
- **Example:** `0.0005`
- **How to get it:** Leave unset on Base.

### `CHANNEL_ACCOUNT_RECEIPIENT`

- **What it is:** The email address that receives generated channel
  accounts' keys.
- **Why it's needed:** Only when accounts are generated. Avoid it: keys
  should not travel by email.
- **Required:** No.
- **Example:** `ops@trovo.example.com`
- **How to get it:** Leave unset on Base (with `CHANNEL_ACCOUNT_MIN_COUNT=0`
  nothing is generated).

## 9. Fees

### Fee wallets

Each fee wallet is the **address** a kind of platform fee is paid to. The
user's own wallet sends the fee, so app-backend never needs the key (a
private key is still accepted for backward compatibility; only its address
is used). Use addresses your finance team controls; on a test network they
can all be one test wallet. None of them may hold an asset's tokens for
payouts (see payout-engine).

#### `VAT_WALLET`

- **What it is:** Where VAT charged on fees is paid.
- **Why it's needed:** Fees in countries with VAT add VAT on top; it must go
  somewhere separate for tax reporting. payout-engine uses the same value.
- **Required:** Yes.
- **Example:** `0x2546BcD3c84621e976D8185a91A922aE77ECEc30`
- **How to get it:** Your finance team's VAT wallet address.

#### `PAYMENT_FEE_WALLET`

- **What it is:** Where fees on ordinary payments are paid.
- **Why it's needed:** Payment fees go here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000010`
- **How to get it:** A treasury wallet address.

#### `SWAP_FEE_WALLET`

- **What it is:** Where swap fees are paid.
- **Why it's needed:** Swap fees go here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000011`
- **How to get it:** A treasury wallet address.

#### `SHARED_ACCESS_PAYMENT_FEE_WALLET`

- **What it is:** Where fees on payments from shared wallets are paid.
- **Why it's needed:** Those fees go here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000012`
- **How to get it:** A treasury wallet address.

#### `SUBWALLET_CREATION_FEE_WALLET`

- **What it is:** Where the fee for creating a sub-wallet is paid.
- **Why it's needed:** The fee itself is the `SUBWALLET_CREATION_FEE` row of
  the `service_fees` table (amount in USD, paid in the stablecoin named
  there; `inactive = 1` turns it off); this is where it goes.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000013`
- **How to get it:** A treasury wallet address.

#### `TOKENIZATION_APPLICATION_FEE_WALLET`

- **What it is:** Where the fee for applying to tokenize an asset is paid.
- **Why it's needed:** Application fees go here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000014`
- **How to get it:** A treasury wallet address.

#### `TOKENIZATION_FEE_WALLET`

- **What it is:** Where the tokenization fee, paid in the asset's own
  tokens at minting, is sent.
- **Why it's needed:** The platform's share of a new asset goes here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000015`
- **How to get it:** A treasury wallet address.

#### `CLOSED_GROUP_FEE_WALLET`

- **What it is:** Where closed-group fees are paid.
- **Why it's needed:** Those fees go here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000016`
- **How to get it:** A treasury wallet address.

#### `ACCOUNT_RECOVERY_FEE_WALLET`

- **What it is:** Where the fee for turning on account recovery is paid.
- **Why it's needed:** The fee itself is the `ACCOUNT_RECOVERY_FEE` row of
  `service_fees` (amount in USD; an inactive row or 0 charges nothing).
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000017`
- **How to get it:** A treasury wallet address.

#### `PATRON_FEE_WALLET`

- **What it is:** Where patron (subscription) payments are paid.
- **Why it's needed:** Patron payments go here.
- **Required:** Yes.
- **Example:** `0xA11CE00000000000000000000000000000000018`
- **How to get it:** A treasury wallet address.

### `SWAP_FEE_ENABLED`

- **What it is:** A switch that must be set for app-backend to start.
- **Why it's needed:** Only the start-up check reads it; whether a swap fee
  is charged is decided by the `service_fees` table.
- **Required:** Yes (any value).
- **Example:** `1`
- **How to get it:** `1`.

### `MARKET_MAKING_FEE_ENABLED`

- **What it is:** A switch that must be set for app-backend to start.
- **Why it's needed:** Only the start-up check reads it.
- **Required:** Yes (any value).
- **Example:** `0`
- **How to get it:** `0`.

### `TOKENIZATION_APPLICATION_FEE_AMOUNT`

- **What it is:** The amount charged when a tokenization application is
  submitted.
- **Why it's needed:** Sets that fee.
- **Required:** No.
- **Example:** `50`
- **How to get it:** Your fee schedule.

### `TOKENIZATION_APPLICATION_FEE_ASSET`

- **What it is:** The asset that fee is paid in.
- **Why it's needed:** Sets which asset is charged.
- **Required:** With `TOKENIZATION_APPLICATION_FEE_AMOUNT`.
- **Example:** `USDC`
- **How to get it:** The code of a stablecoin users hold.

### `CRYPTO_WITHDRAWAL_SERVICE_FEE`

- **What it is:** A fee added to crypto withdrawals through 1Liquidity.
- **Why it's needed:** Covers that service's cost.
- **Required:** No.
- **Example:** `0.001`
- **How to get it:** Your fee schedule.

### `MIN_SENDABLE_AMOUNT`

- **What it is:** The smallest amount a user may send.
- **Why it's needed:** Stops dust payments.
- **Required:** No (no minimum when unset or `0`).
- **Example:** `0.0001`
- **How to get it:** A small amount meaningful for your assets.

### `STANDARD_WALLET_MINIMUM_BALANCE`

- **What it is:** The ETH balance at which a wallet counts as funded for
  TROV activation rewards.
- **Why it's needed:** Decides whether a user gets that activation share.
- **Required:** No.
- **Example:** `0.0001`
- **How to get it:** A small ETH amount.

Other fees (the payout fee, Public Markets fees, gas) are set in Trovo
Manager or the database, not here. Accounts that pay no platform fees are
listed in the `fee_exempt_users` table, managed on Trovo Manager's **Fee
Exemptions** page.

## 10. Tokenization (tokenized real-world assets)

### `ENABLE_ASSET_TOKENIZATION`

- **What it is:** `1` turns on the tokenized-asset endpoints.
- **Why it's needed:** Without it, assets cannot be tokenized, bought or
  sold. With it, `OFFER_BOOK_ADDRESS` is required.
- **Required:** No.
- **Example:** `1`
- **How to get it:** `1` if you offer tokenized assets.

### `TOKENIZATION_ISSUING_PROFILE`

- **What it is:** The username of the Trovo account that issues tokenized
  assets.
- **Why it's needed:** Each asset's issuing and distribution wallets are
  sub-wallets of this account (owned by its key and the asset's minting
  approvers; app-backend holds no key). It pays no platform fees.
  payout-engine uses the same value.
- **Required:** No, default `atprofile`.
- **Example:** `atprofile`
- **How to get it:** Register a Trovo account for issuing and use its
  username.

### `TOKENIZATION_ISSUING_PROFILE_WALLET`

- **What it is:** That account's primary wallet address, as a check.
- **Why it's needed:** If set, app-backend checks it matches.
- **Required:** No.
- **Example:** `0xA11CE00000000000000000000000000000000020`
- **How to get it:** The issuing account's wallet address, or leave unset.

### `OFFER_BOOK_ADDRESS`

- **What it is:** The address of the `TrovoOfferBook` contract.
- **Why it's needed:** Tokenized assets are sold, and swaps and
  market-making offers trade, on it. app-backend also keeps a copy of its
  offers in the database. Without it, swaps and market offers answer
  `error-swaps-not-configured` / `error-market-not-configured`.
- **Required:** With `ENABLE_ASSET_TOKENIZATION=1`.
- **Example:** `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- **How to get it:** Printed by the market deploy script
  ([market/DEPLOYMENT.md](../market/DEPLOYMENT.md)).

### `OFFER_AUTHORIZER_PRIVATE_KEY`

- **What it is:** The key that signs Trovo's permission for each purchase
  and swap.
- **Why it's needed:** The offer book only accepts a purchase with a
  permission from an authorizer, given after app-backend's checks (KYC,
  limits, sale window). A leaked key lets someone skip those checks (at the
  seller's price; it cannot move anyone's tokens), so rotate it if leaked.
- **Required:** With `ENABLE_ASSET_TOKENIZATION=1`.
- **Example:** `0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a` (a test key)
- **How to get it:** Generate a dedicated key; give its address to the offer
  book as an authorizer ([market/CONFIGURATION.md](../market/CONFIGURATION.md#authorizer_addresses)).

### `OFFER_BOOK_START_BLOCK`

- **What it is:** The block app-backend starts reading the offer book from,
  the first time.
- **Why it's needed:** Starting at the deployment block skips years of empty
  history.
- **Required:** No (block 0: works, but slow on mainnet).
- **Example:** `23456789`
- **How to get it:** The block of the offer book's deployment transaction
  on basescan.

### `OFFER_BOOK_CONFIRMATIONS`

- **What it is:** How many blocks behind the newest the offer copy stays.
- **Why it's needed:** Avoids recording a purchase that a short chain
  reorganization undoes.
- **Required:** No, default `2`.
- **Example:** `2`
- **How to get it:** Keep the default.

### `INTERNAL_BALANCE_ISSUING_SIGNERS`

- **What it is:** Owner keys of each country's internal balance minting
  Safe, separated by `;`.
- **Why it's needed:** When a buyer pays in fiat, those keys sign one
  transaction that creates the internal balance token and buys the asset for
  the buyer. The first key pays the network fee, so it needs ETH.
- **Required:** For fiat purchases of tokenized assets.
- **Example:** `0x59c6...;0x5de4...` (test keys)
- **How to get it:** The minting Safe owners' keys, from your secrets
  manager.

A tokenized asset's **token contract** is not configuration: it is
registered per asset in Trovo Manager (see [INTEGRATION.md](INTEGRATION.md#tokenized-assets-token-issuing-and-distribution-wallets-sale-offer)).

## 11. Account recovery

Optional. Without the module and a guardian, the recovery endpoints answer
`error-account-recovery-not-configured` (503); nothing else is affected.

### `RECOVERY_MODULE_ADDRESS`

- **What it is:** The deployed Social Recovery Module.
- **Why it's needed:** Wallets that turn on recovery enable it.
- **Required:** For account recovery.
- **Example:** `0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e`
- **How to get it:** `recoveryModule` in
  `recovery/contracts/deployments/<chainId>.json` ([recovery/DEPLOYMENT.md](../recovery/DEPLOYMENT.md)).

### `RECOVERY_PERIOD_SECONDS`

- **What it is:** The module's waiting period, in seconds.
- **Why it's needed:** Only for messages to users ("your wallets move to the
  new key in 7 days"); the module's own period is what counts.
- **Required:** For account recovery.
- **Example:** `604800`
- **How to get it:** `recoveryPeriodSeconds` in the same file.

### `ACCOUNT_RECOVERY_GUARDIAN_SAFE`

- **What it is:** The recovery guardian Safe.
- **Why it's needed:** It is the only guardian of each covered wallet; it
  can start a key replacement (cancellable by the user), never move funds.
- **Required:** Recommended for account recovery (without it, the first
  signer key is the guardian).
- **Example:** `0xA11CE00000000000000000000000000000000002`
- **How to get it:** The Safe from [recovery/DEPLOYMENT.md step 5](../recovery/DEPLOYMENT.md#5-create-the-guardian).
  Changing it later does not move existing wallets.

### `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS`

- **What it is:** The guardian Safe's signer keys, separated by `;`.
- **Why it's needed:** app-backend signs the guardian's calls with enough of
  them; the first pays the network fee, so keep a little ETH on it.
- **Required:** For account recovery.
- **Example:** `0x59c6...;0x5de4...` (test keys)
- **How to get it:** From your secrets manager.

## 12. P2P escrow

### `P2P_ESCROW_WALLET_ADDRESS`

- **What it is:** The Safe that holds P2P trades' escrowed funds.
- **Why it's needed:** Sellers deposit here, and settlements are paid from
  here.
- **Required:** For P2P trading.
- **Example:** `0xA11CE00000000000000000000000000000000030`
- **How to get it:** Create a Safe at <https://app.safe.global> owned by the
  escrow signer keys, threshold 3.

### `P2P_ESCROW_SIGNERS`

- **What it is:** At least 3 signer keys of the escrow Safe, separated by
  `;`.
- **Why it's needed:** Each settlement needs their signatures; the first
  also sends it and pays its fee (keep ETH on it).
- **Required:** For P2P trading.
- **Example:** `0x59c6...;0x5de4...;0x7c85...` (test keys)
- **How to get it:** From your secrets manager (Trovo Manager's vault
  manager can hold them).

### `P2P_ESCROW_MULTISEND_ADDRESS`

- **What it is:** The MultiSendCallOnly contract used to batch settlement
  payments.
- **Why it's needed:** Only a network with a non-standard Safe deployment
  needs another address.
- **Required:** No (the standard address).
- **Example:** (leave unset)
- **How to get it:** Leave unset.

## 13. Public Markets (tokenized NGX stocks and FMDQ bonds)

See [PUBLIC_MARKETS.md](PUBLIC_MARKETS.md). Thresholds, fees, approvers,
market hours, withholding tax and partner details are **not** environment
variables: they are edited in Trovo Manager. Without `BASE_RPC_URL` the
engine only does database work.

### `PUBLIC_MARKETS_SIGNERS`

- **What it is:** The keys that own the Public Markets Safes (each asset's
  issuing Safe, the treasury and every exchange customer's wallet),
  separated by `;`.
- **Why it's needed:** Minting, payments and wallet creation need their
  signatures; the first key sends and pays the fee (keep ETH on it).
- **Required:** For Public Markets.
- **Example:** `0x59c6...;0x5de4...;0x7c85...` (test keys)
- **How to get it:** A managed secret (Vault), one key per signing
  officer.

### `PUBLIC_MARKETS_TREASURY_SAFE`

- **What it is:** The Safe that receives buyers' cNGN and exchange deposits,
  and pays sellers, dividends, fees and withholding tax.
- **Why it's needed:** All Public Markets money moves through it.
- **Required:** For Public Markets.
- **Example:** `0xA11CE00000000000000000000000000000000040`
- **How to get it:** A Safe owned by the Public Markets signers.

### `PUBLIC_MARKETS_PARTNER_WALLET_THRESHOLD`

- **What it is:** How many signers each exchange customer's wallet needs.
- **Why it's needed:** Sets the signing policy for those wallets (only for
  wallets created after a change).
- **Required:** No (the smaller of 3 and the number of signers).
- **Example:** `3`
- **How to get it:** Your signing policy.

### `PUBLIC_MARKETS_MOCK_EXECUTION_SECONDS`

- **What it is:** How long the mock broker takes to fill an order.
- **Why it's needed:** Only for partners in `MOCK` mode (until the real CSCS
  and broker integrations exist).
- **Required:** No, default `20`.
- **Example:** `20`
- **How to get it:** Keep the default; raise it to demo pending states.

### `PUBLIC_MARKETS_MOCK_SETTLEMENT_SECONDS`

- **What it is:** How long the mock custodian takes to confirm settlement.
- **Why it's needed:** As above (the real cycle is T+2).
- **Required:** No, default `60`.
- **Example:** `60`
- **How to get it:** Keep the default.

### `PUBLIC_MARKETS_PRICE_FEED`

- **What it is:** `rest` to take prices from a price vendor; empty for the
  built-in mock feed (small moves within NGX's ±10% daily band during
  market hours).
- **Why it's needed:** Prices value holdings and orders.
- **Required:** No.
- **Example:** `rest`
- **How to get it:** `rest` once a vendor is contracted.

### `PUBLIC_MARKETS_PRICE_FEED_URL`

- **What it is:** The price vendor's address (`GET /v1/quotes/{isin}`).
- **Why it's needed:** With `PUBLIC_MARKETS_PRICE_FEED=rest`.
- **Required:** With `rest`.
- **Example:** `https://prices.vendor.example.com`
- **How to get it:** From the vendor.

### `PUBLIC_MARKETS_PRICE_FEED_AUTH`

- **What it is:** How requests to the vendor are authenticated: `hmac` or
  `mtls`.
- **Why it's needed:** With `rest`.
- **Required:** With `rest`.
- **Example:** `hmac`
- **How to get it:** What the vendor supports.

### `PUBLIC_MARKETS_PRICE_FEED_CREDENTIALS`

- **What it is:** A credentials reference, like the partners' (below).
- **Why it's needed:** With `rest`.
- **Required:** With `rest`.
- **Example:** `env:PRICE_VENDOR`
- **How to get it:** See partner credentials below.

### Partner credentials (`credentialsRef`)

A Custodian, Dealing Member or the price feed is configured in Trovo
Manager with a **credentials reference** such as `env:CUSTODIAN_A` or
`vault://secret/public-markets/custodian-a#CUSTODIAN_A`. The part after
`env:` (or after `#`) is the **name** of environment variables you then
set on app-backend:

#### `<NAME>` (for example `CUSTODIAN_A`)

- **What it is:** The partner's key id and shared secret, as
  `keyId:secret`.
- **Why it's needed:** Requests to the partner are signed with it, and the
  partner's callbacks are checked with it (HMAC). Without it the partner
  stays on its mock.
- **Required:** For each real (not mock) partner that uses HMAC.
- **Example:** `ck_live_01:xxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Issued by the partner during onboarding. With a
  `vault://` reference, your deployment injects this field from Vault.

#### `<NAME>_CERT`, `<NAME>_KEY`

- **What they are:** Paths to the client certificate and private key (PEM
  files) for partners that use mutual TLS.
- **Why they're needed:** Such partners only accept connections that
  present this certificate.
- **Required:** Only for partners that use mutual TLS.
- **Example:** `/run/secrets/custodian-a.crt`, `/run/secrets/custodian-a.key`
- **How to get them:** The partner signs or issues the certificate during
  onboarding; mount both files into the container read-only.

## 14. Bank (Naira) deposits and withdrawals

### `STABLERAIL_MIN_WITHDRAWAL`

- **What it is:** The smallest bank withdrawal (in Naira) a user may ask
  for.
- **Why it's needed:** Stablerail has a minimum payout.
- **Required:** No, default `1000`.
- **Example:** `1000`
- **How to get it:** Your Stablerail account's minimum, or higher.

Stablerail's own API key and switch are **not** environment variables: they
are in the `stablerail_configs` table (`api_key`, `base_url`,
`enable_stablerail`), read on every request. There is no admin page for it
yet; set the row directly.

## 15. Crypto deposits and withdrawals (1Liquidity)

### `ENABLE_CRYPTO_DEPOSIT_MINTING`

- **What it is:** `1` turns on minting for confirmed crypto deposits.
- **Why it's needed:** For the 1Liquidity deposit flow.
- **Required:** No.
- **Example:** `0`
- **How to get it:** `0` unless you use 1Liquidity.

### `CRYPTO_DEPOSIT_MINTING_INITIATOR_PUBLIC_KEY`

- **What it is:** The wallet address that initiates that minting.
- **Why it's needed:** Required with `ENABLE_CRYPTO_DEPOSIT_MINTING=1`
  (a 42-character address).
- **Required:** With `ENABLE_CRYPTO_DEPOSIT_MINTING=1`.
- **Example:** `0xA11CE00000000000000000000000000000000050`
- **How to get it:** Your minting initiator's address.

### `ENABLE_CRYPTO_WITHDRAWAL_SERVICE`

- **What it is:** `1` turns on crypto withdrawals through 1Liquidity.
- **Why it's needed:** For that flow.
- **Required:** No.
- **Example:** `0`
- **How to get it:** `0` unless you use 1Liquidity.

### `ONELIQUIDITY_BASE_URL`

- **What it is:** 1Liquidity's API address.
- **Why it's needed:** For their deposit, withdrawal and compliance API.
- **Required:** With the 1Liquidity features.
- **Example:** `https://api.oneliquidity.example.com`
- **How to get it:** Issued by 1Liquidity at onboarding.

### `ONELIQUIDITY_TOKEN`

- **What it is:** 1Liquidity's API token.
- **Why it's needed:** Authenticates those calls.
- **Required:** With the 1Liquidity features.
- **Example:** `ol_xxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Issued by 1Liquidity.

### `ONELIQUIDITY_WITHDRAWAL_CURRENCY_LIST`

- **What it is:** Currencies withdrawable through 1Liquidity, separated by
  commas.
- **Why it's needed:** Limits the withdrawal options.
- **Required:** With `ENABLE_CRYPTO_WITHDRAWAL_SERVICE=1`.
- **Example:** `USDT,USDC`
- **How to get it:** The currencies your 1Liquidity account supports.

### `COMPLIANCE_ACCOUNT_ID`

- **What it is:** Your account id in 1Liquidity's compliance API.
- **Why it's needed:** Sent with compliance checks.
- **Required:** With the 1Liquidity features.
- **Example:** `acct_test_123`
- **How to get it:** Issued by 1Liquidity.

## 16. Faucets and rewards

### `GAS_FAUCET`

- **What it is:** The private key of a wallet that sends small amounts of
  ETH when the server makes a gas-faucet payment.
- **Why it's needed:** For flows that seed wallets with gas.
- **Required:** No.
- **Example:** `0x8b3a...` (a test key)
- **How to get it:** A dedicated, lightly funded wallet.

### `GAS_FAUCET_MIN_BALANCE`

- **What it is:** The balance below which the gas faucet alerts.
- **Why it's needed:** So someone refills it (`FAUCET_LOW_BALANCE_WEBHOOK`).
- **Required:** No.
- **Example:** `0.05`
- **How to get it:** A few days' worth of faucet payments.

### `REWARD_FAUCET`

- **What it is:** The private key of the wallet that pays reward tokens.
- **Why it's needed:** For reward payments.
- **Required:** No.
- **Example:** `0x47e1...` (a test key)
- **How to get it:** A dedicated wallet holding the reward token.

### `REWARD_FAUCET_MIN_BALANCE`

- **What it is:** The reward-token balance below which it alerts.
- **Why it's needed:** So someone refills it.
- **Required:** No.
- **Example:** `10`
- **How to get it:** A few days' worth of rewards.

### `REWARD_FAUCET_GAS_MIN_BALANCE`

- **What it is:** The ETH balance below which the reward faucet alerts.
- **Why it's needed:** It pays network fees in ETH.
- **Required:** No.
- **Example:** `0.05`
- **How to get it:** A few days' worth of fees.

### `REWARD_ASSET_CODE`

- **What it is:** The reward token's code.
- **Why it's needed:** Which token rewards are paid in.
- **Required:** With `REWARD_FAUCET`.
- **Example:** `TROV`
- **How to get it:** Your reward token's code.

### `REWARD_ASSET_CONTRACT_ADDRESS`

- **What it is:** The reward token's contract address.
- **Why it's needed:** As above.
- **Required:** With `REWARD_FAUCET`.
- **Example:** `0xC0FFEE0000000000000000000000000000000003`
- **How to get it:** The token's address on your network.

## 17. Email (Mailgun)

### `MAILGUN_PRIVATE_API_KEY`

- **What it is:** Mailgun's private API key.
- **Why it's needed:** Every email (verification codes, notifications,
  account deletion) is sent through Mailgun.
- **Required:** Yes.
- **Example:** `<your-mailgun-private-api-key>` (Mailgun keys start with
  `key-`; not shown here, as that pattern trips GitHub's secret scanner)
- **How to get it:** Sign up at <https://www.mailgun.com>, add and verify a
  sending domain, then **API Keys**.

### `MAILGUN_DOMAIN`

- **What it is:** Your Mailgun sending domain.
- **Why it's needed:** Emails are sent from it.
- **Required:** Yes.
- **Example:** `mg.trovo.example.com`
- **How to get it:** The domain verified in Mailgun.

### `ENABLE_EMAIL_VALIDATION`

- **What it is:** `1` checks email addresses with Mailgun's validation
  service before accepting them.
- **Why it's needed:** Catches mistyped and fake addresses (paid feature).
  It must be set (`0` or `1`).
- **Required:** Yes.
- **Example:** `0`
- **How to get it:** `1` only if you pay for Mailgun validation.

### `MAILGUN_VALIDATOR_API_KEY`

- **What it is:** The key for Mailgun's validation service.
- **Why it's needed:** With `ENABLE_EMAIL_VALIDATION=1`.
- **Required:** With `ENABLE_EMAIL_VALIDATION=1`.
- **Example:** `<your-mailgun-validator-api-key>`
- **How to get it:** Mailgun dashboard → **API Keys**.

### `SHOW_MAIL_VALIDATION_RESULT`

- **What it is:** `1` includes the validation result in responses.
- **Why it's needed:** Only for debugging.
- **Required:** No.
- **Example:** `0`
- **How to get it:** Leave unset.

### `ENABLE_EMAIL_NOTIFICATIONS`

- **What it is:** `1` sends notification emails (for example account
  deletion confirmations).
- **Why it's needed:** Turns those emails on.
- **Required:** No (off when unset).
- **Example:** `1`
- **How to get it:** `1` once Mailgun works.

### `MAIL_SENDER`

- **What it is:** The "From" address of verification emails.
- **Why it's needed:** Recipients see it.
- **Required:** No (a default on `MAILGUN_DOMAIN`).
- **Example:** `Trovo <no-reply@mg.trovo.example.com>`
- **How to get it:** An address on your Mailgun domain.

### `DEFAULT_MAIL_SENDER`

- **What it is:** The "From" address of account deletion and recovery
  emails.
- **Why it's needed:** As above.
- **Required:** No (a default on `MAILGUN_DOMAIN`).
- **Example:** `Trovo <no-reply@mg.trovo.example.com>`
- **How to get it:** As above.

### `EMAIL_VERIFICATION_SUBJECT`

- **What it is:** The subject of verification-code emails.
- **Why it's needed:** To word it your way.
- **Required:** No (a built-in subject).
- **Example:** `Your Trovo verification code`
- **How to get it:** Your choice.

### `ACCOUNT_DELETION_REQUEST_TEMPLATE`

- **What it is:** The Mailgun template for account deletion emails.
- **Why it's needed:** The email is sent from this template.
- **Required:** No, default `account-deletion-request-template`.
- **Example:** `account-deletion-request-template`
- **How to get it:** Create a template with that name in Mailgun
  (**Sending → Templates**).

### `ACCOUNT_DELETION_EMAIL_SUBJECT`

- **What it is:** That email's subject.
- **Why it's needed:** To word it your way.
- **Required:** No, default `TrovoApp Account Deletion Request`.
- **Example:** `Your Trovo account deletion request`
- **How to get it:** Your choice.

### `ACCOUNT_DELETION_DAYS`

- **What it is:** Days between a deletion request and the deletion.
- **Why it's needed:** Gives users time to change their mind.
- **Required:** No, default `30`.
- **Example:** `30`
- **How to get it:** Your policy.

### `SUPPORT_EMAIL`

- **What it is:** Your support address, shown in some emails.
- **Why it's needed:** So users know where to get help.
- **Required:** No.
- **Example:** `support@trovo.example.com`
- **How to get it:** Your support inbox.

## 18. SMS

### `DEFAULT_SMS_PROVIDER`

- **What it is:** `termii` or `infobip`: the provider used when a phone
  number has no specific provider set in the database.
- **Why it's needed:** Verification codes by SMS go through it.
- **Required:** For SMS.
- **Example:** `termii`
- **How to get it:** Whichever provider you signed up with.

### `TERMII_SMS_API_KEY`

- **What it is:** Termii's API key.
- **Why it's needed:** To send SMS through Termii.
- **Required:** With Termii.
- **Example:** `<your-termii-api-key>`
- **How to get it:** <https://termii.com> dashboard → API settings.

### `TERMII_SMS_URL`

- **What it is:** Termii's send address.
- **Why it's needed:** Where SMS requests go.
- **Required:** With Termii.
- **Example:** `https://api.ng.termii.com/api/sms/send`
- **How to get it:** Termii's documentation.

### `TERMII_SMS_SENDER_ID`

- **What it is:** The sender name shown on Termii SMS.
- **Why it's needed:** Termii requires a registered sender ID.
- **Required:** With Termii.
- **Example:** `Trovo`
- **How to get it:** Register one in the Termii dashboard.

### `INFOBIP_SMS_API_KEY`

- **What it is:** Infobip's API key.
- **Why it's needed:** To send SMS through Infobip.
- **Required:** With Infobip.
- **Example:** `<your-infobip-api-key>`
- **How to get it:** <https://www.infobip.com> → **Settings → API Keys**.

### `INFOBIP_SMS_HOST`

- **What it is:** Your Infobip account's API host.
- **Why it's needed:** Where SMS requests go.
- **Required:** With Infobip.
- **Example:** `xxxxxx.api.infobip.com`
- **How to get it:** Shown on the same Infobip page.

### `TROVOWALLET_SMS_FROM`

- **What it is:** The sender name on outgoing SMS.
- **Why it's needed:** Recipients see it.
- **Required:** No, default `TrovoWallet`.
- **Example:** `Trovo`
- **How to get it:** A short name your provider allows.

### `ENABLE_MOBILE_VERIFICATION`

- **What it is:** `1` requires a verified phone number for some partner
  (service-link) flows.
- **Why it's needed:** Turns that requirement on.
- **Required:** No.
- **Example:** `0`
- **How to get it:** `1` once SMS works and you want it.

## 19. Firebase, storage and links

### `GC`

- **What it is:** The Firebase service account key (a JSON file), encoded as
  base64.
- **Why it's needed:** Push notifications and file storage use it.
  payout-engine uses the same value.
- **Required:** Yes.
- **Example:** `ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAidHJvdm8tZXhhbXBsZSIKfQ==`
- **How to get it:** [Firebase console](https://console.firebase.google.com)
  → Project settings → **Service accounts → Generate new private key**,
  then `base64 -w0 key.json` (macOS: `base64 -i key.json | tr -d '\n'`).

### `GOOGLE_PROJECT_ID`

- **What it is:** The Firebase / Google Cloud project id.
- **Why it's needed:** Identifies the project for storage and messaging.
- **Required:** Yes.
- **Example:** `trovo-wallet-prod`
- **How to get it:** Firebase console → Project settings → **Project ID**.

### `STORAGE_BUCKET_NAME`

- **What it is:** The cloud storage bucket for uploads (for example profile
  pictures).
- **Why it's needed:** Uploaded files are stored there.
- **Required:** For uploads.
- **Example:** `trovo-wallet-prod.appspot.com`
- **How to get it:** Firebase console → **Storage** (the bucket name).

### `STAKEHOLDER_DOCUMENTS_BUCKET_NAME`

- **What it is:** A separate, private bucket for stakeholder-portal
  documents.
- **Why it's needed:** Those documents must not share a bucket with public
  files. Without it, those endpoints answer 503.
- **Required:** For the stakeholder portal's documents.
- **Example:** `trovo-stakeholder-docs-prod`
- **How to get it:** Create a private bucket in [Google Cloud
  Storage](https://console.cloud.google.com/storage) and give the `GC`
  service account access.

### `DYNAMIC_LINKS_DOMAIN_PREFIX`

- **What it is:** The link domain for app links.
- **Why it's needed:** Required by the start-up check.
- **Required:** Yes.
- **Example:** `trovo.page.link`
- **How to get it:** Your app-link domain.

### `DYNAMIC_LINKS_ANDROID_PACKAGE_NAME`

- **What it is:** The Android app's package name.
- **Why it's needed:** Links open the app, and it is served in the Android
  app-links file (`/.well-known/assetlinks.json`).
- **Required:** Yes.
- **Example:** `app.trovo.wallet`
- **How to get it:** The app's `applicationId` (app-mobile's Android
  settings).

### `DYNAMIC_LINKS_IOS_BUNDLE_ID`

- **What it is:** The iOS app's bundle id.
- **Why it's needed:** Served in the iOS app-links file
  (`/.well-known/apple-app-site-association`).
- **Required:** Yes.
- **Example:** `app.trovo.wallet`
- **How to get it:** The app's bundle identifier in Xcode / App Store
  Connect.

### `DYNAMIC_LINKS_FALLBACK_BASE_URL`

- **What it is:** Where a link goes when the app is not installed.
- **Why it's needed:** So links still lead somewhere useful.
- **Required:** Yes.
- **Example:** `https://trovo.example.com`
- **How to get it:** Your website or app-store page.

### `FBDL_SERVICE_API_KEY`

- **What it is:** The Firebase web API key used to shorten links.
- **Why it's needed:** Creates short referral and share links.
- **Required:** For short links.
- **Example:** `<your-firebase-web-api-key>`
- **How to get it:** Firebase console → Project settings → **Web API Key**.

### `SHORT_LINKS_BASE_URL`

- **What it is:** The base address of Trovo's own short links.
- **Why it's needed:** Used when links are made without Firebase.
- **Required:** No, default `https://trovo.app`.
- **Example:** `https://trovo.example.com`
- **How to get it:** Your own domain.

### `WALLET_DOMAIN`

- **What it is:** The domain added to usernames to make wallet aliases
  (`ada@trovo.app`).
- **Why it's needed:** Aliases are used with partners such as 1Liquidity.
- **Required:** Yes.
- **Example:** `trovo.app`
- **How to get it:** Your platform's domain.

## 20. Login tokens (JWT)

### `JWT_ACCESS_SECRET`

- **What it is:** The secret that signs and checks login tokens issued by
  app-backend (for example to Trovo Manager's admins).
- **Why it's needed:** Anyone with it can make valid tokens, so keep it
  secret and the same on every copy.
- **Required:** Yes.
- **Example:** `5f2b8c9e0d1a4f7b3c6e8a2d9f0b1c4e7a3d6f9b2c5e8a1d4f7b0c3e6a9d2f5b`
- **How to get it:** `openssl rand -hex 32`.

### `JWT_TOKEN_EXPIRY`

- **What it is:** How long a login token lasts, **in minutes**.
- **Why it's needed:** Short tokens limit the damage of a stolen one.
- **Required:** Yes.
- **Example:** `15`
- **How to get it:** Your security policy (10–60 minutes is typical).

### `JWT_REFRESH_TOKEN_EXPIRY`

- **What it is:** How long a refresh token lasts, **in minutes**.
- **Why it's needed:** It lets a user get a new token without logging in
  again until it expires.
- **Required:** Yes.
- **Example:** `1440` (one day)
- **How to get it:** Your security policy.

## 21. Partner (service link) settings

### `SERVICE_LINK_LOGIN_REQUEST_VALIDITY`

- **What it is:** How long a partner's login request stays valid, **in
  minutes**.
- **Why it's needed:** Unapproved requests expire.
- **Required:** No, default `3`.
- **Example:** `3`
- **How to get it:** Keep the default.

### `SERVICE_LINK_AUTHORIZATION_REQUEST_VALIDITY`

- **What it is:** How long a partner's authorization request stays valid,
  **in minutes**, unless the request asks for its own.
- **Why it's needed:** Unapproved requests expire.
- **Required:** No, default `3`.
- **Example:** `5`
- **How to get it:** Keep the default.

Partners' API keys are rows of the `service_links` table, not environment
variables (see [INTEGRATION.md](INTEGRATION.md#3-partner-businesses-service-links)).

## 22. Location lookup

No code calls the location lookup today (`GetGeoIP` in
`internal/components/users/models/geo.go` is never used), but both
settings are on the start-up check, so they must not be empty.

### `IPAPI_HOST`

- **What it is:** The address of an IP-to-location service.
- **Why it's needed:** Only to pass the start-up check; nothing calls it.
- **Required:** Yes (the start-up check requires it).
- **Example:** `https://api.ipapi.com/api`
- **How to get it:** Your provider's API address (for example ipapi.com),
  or any non-empty value.

### `IPAPI_KEY`

- **What it is:** That service's API key.
- **Why it's needed:** As above.
- **Required:** Yes.
- **Example:** `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx` (placeholder)
- **How to get it:** Sign up at the provider (for example
  <https://ipapi.com>) and copy your key, or any non-empty value.

## 23. Features

### `ENABLE_PATRON`

- **What it is:** `1` turns on patron (subscription) memberships and their
  background activation.
- **Why it's needed:** Turns the feature on.
- **Required:** No.
- **Example:** `0`
- **How to get it:** `1` if you offer patron memberships.

### `ENABLE_WEBSOCKET_AUTH`

- **What it is:** `1` requires authentication on live-update (websocket)
  connections.
- **Why it's needed:** Without it, anyone could listen to a user's live
  updates.
- **Required:** No, but set `1`.
- **Example:** `1`
- **How to get it:** `1` in every real deployment.

### `REGISTRATION_THROTTLE_PER_IP`

- **What it is:** The number of **seconds** that must pass between two
  registrations from the same IP address.
- **Why it's needed:** Slows down mass fake sign-ups.
- **Required:** No (off when unset or `0`).
- **Example:** `3600`
- **How to get it:** Your abuse tolerance; `0` for local development.

## 24. Discord alerts

Each alert type falls back to a built-in Discord webhook when unset. Set
your own so alerts reach your team.

| Parameter | Alerts about |
|---|---|
| `CONNECTION_WARNING_WEBHOOK` | the database connection pool filling up |
| `EXPANSION_NETWORK_ERROR_WEBHOOK` | blockchain node and network errors |
| `FAILED_PAYMENT_ERROR_WEBHOOK` | failed payments, recoveries and patron payments |
| `FAUCET_LOW_BALANCE_WEBHOOK` | faucet wallets running low |
| `REGISTRATION_ERROR_WEBHOOK` | failed or throttled registrations |

For each one:

- **What it is:** A Discord webhook address for that kind of alert.
- **Why it's needed:** So your team sees those problems.
- **Required:** No (the built-in webhook otherwise).
- **Example:** `https://discord.com/api/webhooks/123456789012345678/abcdefghijklmnopqrstuvwxyz`
- **How to get it:** In Discord: the channel's **Edit Channel →
  Integrations → Webhooks → New Webhook → Copy Webhook URL**.

## 25. Used only by tests

| Parameter | What it is |
|---|---|
| `AA_LOCAL_STACK` | the paymaster local stack's deployment file, for end-to-end tests on a local chain |
| `AA_LOCAL_MARKET` | the market contracts' local deployment file |
| `AA_LOCAL_RECOVERY` | the recovery module's local deployment file |
| `AA_RPC_URL`, `AA_QUOTE_SERVICE_URL`, `AA_QUOTE_SERVICE_API_KEY` | the local chain and quote service for `internal/aa` tests |
| `AA_SAFE_ETH_TX_FILE` | a file the swaps test writes for payment-history-engine's test |
| `PUBLIC_MARKETS_TEST_POSTGRES` | a Postgres database to run the Public Markets tests against (each test in its own schema) instead of SQLite |

See [DEPLOYMENT.md](DEPLOYMENT.md) and the contract projects'
DEPLOYMENT.md files for how to start the local chain.

## 26. Listed in `.env.example` but not read

These appear in old templates but no code reads them; leave them unset:
`SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `MAIL_FROM`,
`FBDL_SERVICE_URLS`, `ENABLE_OLD_USER_MIGRATION`, `ENABLE_REFERRAL_REWARD`,
`REFERRAL_REWARD_AMOUNT`, `ENABLE_DOLLAR_ASSET_BY_DEFAULT`,
`ENABLE_TROV_ASSET_BY_DEFAULT`, `EMAIL_VERIFICATION_TEMPLATE`,
`IMPORT_ERROR_WEBHOOK`, `MARKET_MAKING_FEE_AMOUNT`,
`MARKET_MAKING_FEE_WALLET`, `SHARED_ACCESS_FEE_ADDRESS`,
`SHARED_ACCESS_FEE_AMOUNT`, `SHARED_ACCESS_FEE_ASSET_CODE`,
`SHARED_ACCESS_FEE_ASSET_CONTRACT_ADDRESS`, `WALLET_MINIMUM_BALANCE`,
`WALLET_SIGNER_ACTIVATION_AMOUNT`, `INTERNAL_BALANCE_AUTHORIZER_WALLET`.
Sub-wallet seed amounts (`ISSUING_SUB_WALLET_ACTIVATION_AMOUNT`,
`MM_SUB_WALLET_ACTIVATION_AMOUNT`, `BULKPAYMENT_SUB_WALLET_ACTIVATION_AMOUNT`,
`SUB_WALLET_ACTIVATION_AMOUNT`) are rows of the `activation_amounts` table
(id = the name, `amount` in ETH, `inactive`), not environment variables.

## Memo requirements (exchange deposit addresses)

### `WALLETS_REQUIRE_16_BYTE_MEMO`, `WALLETS_REQUIRE_28_BYTE_MEMO`, `WALLETS_REQUIRE_VARIABLE_BYTE_MEMO`

- **What they are:** Lists of destination addresses (usually exchange
  deposit addresses), separated by commas, that need a memo of exactly 16
  characters, exactly 28 characters, or at least 9 characters.
- **Why they're needed:** Some exchanges credit deposits by memo; a payment
  to them without the right memo can be lost. app-backend refuses such a
  payment with a warning.
- **Required:** No.
- **Example:** `0xA11CE00000000000000000000000000000000060,0xA11CE00000000000000000000000000000000061`
- **How to get them:** The exchange's deposit instructions.

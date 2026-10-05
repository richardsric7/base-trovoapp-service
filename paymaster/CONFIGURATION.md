# paymaster — configuration

Every parameter the paymaster project reads. Each one says what it is, why
it is needed, whether you must set it, an example and how to get a real
value. All examples are illustrative: never reuse them.

The project has two parts, configured the same way:

| Part | What it is | Its settings |
|---|---|---|
| `contracts/scripts/deploy.js` | the script that deploys the paymaster contract (once per network) | [Deploy script secret](#2-deploy-script-secret) |
| `quote-service` | the always-on service that prices gas in stablecoins | [Quote service secret](#3-quote-service-secret) |

## How to set them

Both parts read their settings from **HashiCorp Vault** (a secure store for
secrets), so private keys and API keys never sit in files, images or shell
history:

1. A few **bootstrap environment variables** say where Vault is and which
   secret to read ([section 1](#1-bootstrap-environment-variables)).
2. Every other setting is a key of one **Vault secret** per part
   (sections [2](#2-deploy-script-secret) and [3](#3-quote-service-secret)).
   [Section 5](#5-writing-the-secrets-and-the-vault-policies) shows how to
   write them.

| Secret (default path) | Read by |
|---|---|
| `secret/trovo/paymaster/deploy` | the deploy script, once per deployment |
| `secret/trovo/paymaster/quote-service` | the quote service, at start and on reload |

For local development only, `PAYMASTER_CONFIG_SOURCE=env` reads the same
keys from ordinary environment variables (see [section 4](#4-local-development)).

### The minimum

| Part | Parameters |
|---|---|
| deploy script | [`VAULT_ADDR`](#vault_addr), [`VAULT_TOKEN`](#vault_token); secret: [`RPC_URL`](#rpc_url), [`CHAIN_ID`](#chain_id), [`DEPLOYER_PRIVATE_KEY`](#deployer_private_key), [`OWNER_ADDRESS`](#owner_address), [`QUOTE_SIGNER_ADDRESSES`](#quote_signer_addresses), [`PAUSER_ADDRESS`](#pauser_address), [`GAS_TOKENS` (deploy)](#gas_tokens-deploy) |
| quote service | [`VAULT_ADDR`](#vault_addr), [`VAULT_TOKEN`](#vault_token); secret: [`RPC_URL`](#rpc_url-quote-service), [`CHAIN_ID`](#chain_id-quote-service), [`PAYMASTER_ADDRESS`](#paymaster_address), [`QUOTE_SIGNER_PRIVATE_KEY`](#quote_signer_private_key), [`API_KEYS`](#api_keys), [`GAS_TOKENS`](#gas_tokens-quote-service), [`RATE_PAIRS`](#rate_pairs) |

---

## 1. Bootstrap environment variables

These are the only values read from the process environment.

### `VAULT_ADDR`

- **What it is:** The web address of your Vault server.
- **Why it's needed:** Both parts read all their other settings from Vault.
- **Required:** Yes, unless `PAYMASTER_CONFIG_SOURCE=env`.
- **Example:** `https://vault.internal.trovo.io:8200`
- **How to get it:** The same Vault tm-api's vault manager uses (tm-api's
  `VAULT_ADDR`); ask whoever runs Vault.

### `VAULT_TOKEN`

- **What it is:** A Vault access token.
- **Why it's needed:** Vault only returns a secret to a token whose policy
  allows reading it.
- **Required:** Yes, unless `PAYMASTER_CONFIG_SOURCE=env`.
- **Example:** `hvs.CAESIJexampleexampleexample`
- **How to get it:** Ask your Vault administrator for a token with the
  read-only policy in [section 5](#5-writing-the-secrets-and-the-vault-policies):
  `vault token create -policy=paymaster-quote-service -period=768h` for the
  service, or a short-lived one (`-ttl=1h`) with the deploy policy for a
  deployment. In production, have your platform inject it (Kubernetes
  auth, an AppRole login in a start-up step) instead of putting it in an
  image.

### `VAULT_KV_MOUNT`

- **What it is:** The name of Vault's key-value (KV version 2) store.
- **Why it's needed:** A secret's full location is mount + path.
- **Required:** No, default `secret`.
- **Example:** `secret`
- **How to get it:** Keep the default unless your Vault uses another mount
  (`vault secrets list`).

### `VAULT_NAMESPACE`

- **What it is:** The Vault namespace the secrets live in.
- **Why it's needed:** Only Vault Enterprise and HCP Vault have
  namespaces; without the right one a secret is not found.
- **Required:** No.
- **Example:** `trovo`
- **How to get it:** Your Vault administrator; leave unset on open-source
  Vault.

### `PAYMASTER_DEPLOY_SECRET_PATH`

- **What it is:** Where in Vault the deploy script's secret is stored
  (deploy script only).
- **Why it's needed:** It tells the script which secret to read; a path
  per network keeps test and production settings apart.
- **Required:** No, default `trovo/paymaster/deploy`.
- **Example:** `trovo/paymaster/base-sepolia/deploy`
- **How to get it:** Choose it when you write the secret.

### `QUOTE_SERVICE_SECRET_PATH`

- **What it is:** Where in Vault the quote service's secret is stored
  (quote service only).
- **Why it's needed:** It tells the service which secret to read.
- **Required:** No, default `trovo/paymaster/quote-service`.
- **Example:** `trovo/paymaster/base-mainnet/quote-service`
- **How to get it:** Choose it when you write the secret.

### `PAYMASTER_CONFIG_SOURCE`

- **What it is:** Where the settings come from: `vault` or `env`.
- **Why it's needed:** `env` lets you run everything on your own machine
  without Vault.
- **Required:** No, default `vault`.
- **Example:** `env`
- **How to get it:** `env` only for local development; never in production.

---

## 2. Deploy script secret

Keys of `secret/trovo/paymaster/deploy`, read by `npm run deploy` (see
[DEPLOYMENT.md](DEPLOYMENT.md)).

### `RPC_URL`

- **What it is:** The web address of a Base blockchain node.
- **Why it's needed:** The script sends the deployment, stake, deposit and
  setup transactions through it.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
- **How to get it:** From a node provider's dashboard (Alchemy, QuickNode,
  Coinbase Developer Platform): create an app for Base Mainnet or Base
  Sepolia and copy its HTTPS URL.

### `CHAIN_ID`

- **What it is:** The number of the network you mean to deploy to.
- **Why it's needed:** The script refuses to run if `RPC_URL` is another
  network.
- **Required:** Yes.
- **Example:** `8453` (Base Mainnet), `84532` (Base Sepolia)
- **How to get it:** `8453` for production, `84532` for testing.

### `DEPLOYER_PRIVATE_KEY`

- **What it is:** The private key of the wallet that deploys the contract.
- **Why it's needed:** It pays for the deployment, the stake and the first
  deposit. It owns the paymaster only until the end of the script, when
  ownership moves to `OWNER_ADDRESS`.
- **Required:** Yes.
- **Example:** `0x4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318`
  (illustrative)
- **How to get it:** Create a fresh wallet (`cast wallet new` with Foundry,
  or the `node -e` command in [DEPLOYMENT.md](DEPLOYMENT.md#4-create-the-keys-and-addresses)),
  fund it with `STAKE_AMOUNT_ETH` + `INITIAL_DEPOSIT_ETH` + about 0.01 ETH
  for fees, and empty it and delete the key from Vault afterwards.

### `ENTRYPOINT_ADDRESS`

- **What it is:** The address of the ERC-4337 EntryPoint v0.7: the shared
  contract that runs smart-wallet operations and pays for their gas.
- **Why it's needed:** The paymaster is tied to one EntryPoint and keeps its
  ETH deposit there.
- **Required:** No, default `0x0000000071727De22E5E9d8BAf0edAc6f37da032`
  (the official v0.7 deployment, the same on Base and Base Sepolia). The
  script checks it really is a v0.7 EntryPoint.
- **Example:** `0x0000000071727De22E5E9d8BAf0edAc6f37da032`
- **How to get it:** Keep the default. Only a private test chain needs its
  own (the local stack prints it).

### `OWNER_ADDRESS`

- **What it is:** The final owner of the paymaster: the treasury Safe.
- **Why it's needed:** The owner enables gas tokens and their rate limits,
  rotates quote signers, unpauses, and withdraws collected tokens and ETH.
- **Required:** Yes. It must be a contract (a Safe) unless
  `ALLOW_EOA_OWNER=true`.
- **Example:** `0xA11CE00000000000000000000000000000000003`
- **How to get it:** The treasury Safe's address; create it at
  <https://app.safe.global> (owners: the finance and platform signers,
  threshold at least 2).

### `ALLOW_EOA_OWNER`

- **What it is:** Allows a plain wallet as owner.
- **Why it's needed:** Convenient on a test network; on mainnet one leaked
  key would control the paymaster's funds.
- **Required:** No, default `false`.
- **Example:** `true`
- **How to get it:** Only on a test network.

### `QUOTE_SIGNER_ADDRESSES`

- **What it is:** The addresses whose signed gas quotes the paymaster
  accepts, separated by commas.
- **Why it's needed:** The paymaster only pays for an operation that
  carries a quote signed by one of these. During a rotation, list the
  current and the next signer.
- **Required:** Yes.
- **Example:** `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
- **How to get it:** The address of the key you give the quote service as
  [`QUOTE_SIGNER_PRIVATE_KEY`](#quote_signer_private_key)
  ([DEPLOYMENT.md step 4](DEPLOYMENT.md#4-create-the-keys-and-addresses)).

### `PAUSER_ADDRESS`

- **What it is:** The address allowed to pause the paymaster in an
  emergency.
- **Why it's needed:** Pausing must be fast (one signature); only the owner
  can unpause, so a pauser cannot do lasting harm.
- **Required:** Yes.
- **Example:** `0x90F79bf6EB2c4f870365E785982E1f101E93b906`
- **How to get it:** An on-call operations key, or a Safe needing one
  signature from several on-call people.

### `POST_OP_GAS_OVERHEAD`

- **What it is:** Extra gas added to each operation's charge.
- **Why it's needed:** It covers the paymaster's own work after the
  operation (charging the token), which the EntryPoint leaves out of the
  cost it reports. The owner can change it later
  (`setPostOpGasOverhead`).
- **Required:** No, default `45000`.
- **Example:** `45000`
- **How to get it:** Keep the default.

### `STAKE_AMOUNT_ETH`

- **What it is:** ETH locked at the EntryPoint as the paymaster's stake.
- **Why it's needed:** Bundlers only accept paymasters like this one (that
  touch token balances during checks) if they are staked. The stake is not
  spent.
- **Required:** No, default `0.1`.
- **Example:** `0.5`
- **How to get it:** Your bundler's documented minimum stake, or more.

### `UNSTAKE_DELAY_SEC`

- **What it is:** How long, in seconds, the stake stays locked after
  unlocking before it can be withdrawn.
- **Why it's needed:** Bundlers require at least one day.
- **Required:** No, default `86400`.
- **Example:** `86400`
- **How to get it:** Keep the default, or your bundler's minimum.

### `INITIAL_DEPOSIT_ETH`

- **What it is:** ETH deposited at the EntryPoint to pay users' gas.
- **Why it's needed:** This is the money that is actually spent: the
  paymaster pays gas in ETH from it and is repaid in stablecoins. Top it
  up from converted stablecoins later.
- **Required:** No, default `0.05`.
- **Example:** `1`
- **How to get it:** Enough for a few days of expected traffic.

### `GAS_TOKENS` (deploy)

- **What it is:** A JSON list of the stablecoins to enable for gas, with
  the range of exchange rates the paymaster will accept for each.
- **Why it's needed:** Only enabled tokens can pay gas. The rate range is a
  safety net: even a stolen quote-signer key cannot sign a rate outside it.
- **Required:** Yes.
- **Example:**
  ```json
  [
    {"symbol": "USDC", "address": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "minRate": "500000000", "maxRate": "20000000000"},
    {"symbol": "cNGN", "address": "<cNGN token address on Base>", "minRate": "500000000000", "maxRate": "40000000000000"}
  ]
  ```
- **How to get it:**
  - `address`: the token contract. USDC above is Base Mainnet's native
    USDC; take cNGN's and USDT's from the issuer's official documentation.
  - `minRate` / `maxRate`: in **token base units per 1 ETH**, spread
    included. USDC has 6 decimals, so `500000000` = 500 USDC per ETH and
    `20000000000` = 20,000 USDC per ETH. Keep them wide enough for normal
    price swings (about ¼× to 4× today's price) and review them when the
    market moves a lot. If quotes start failing with `rate-out-of-bounds`,
    the owner Safe updates them with `setToken`.

---

## 3. Quote service secret

Keys of `secret/trovo/paymaster/quote-service`.

### Chain and signing

#### `RPC_URL` (quote service)

- **What it is:** The web address of a Base node.
- **Why it's needed:** The service reads the paymaster's settings, its ETH
  deposit, Chainlink price feeds and DEX pools through it. It only reads;
  it never sends transactions.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
- **How to get it:** As for the deploy script's [`RPC_URL`](#rpc_url);
  a separate key per service makes usage easier to track.

#### `CHAIN_ID` (quote service)

- **What it is:** The network number.
- **Why it's needed:** It is part of every signed quote, and the service
  refuses to start if the node is another network.
- **Required:** Yes.
- **Example:** `8453`
- **How to get it:** `8453` (Base Mainnet) or `84532` (Base Sepolia).

#### `ENTRYPOINT_ADDRESS` (quote service)

- **What it is:** The EntryPoint v0.7 address.
- **Why it's needed:** It must be the one the paymaster was deployed with;
  checked at start.
- **Required:** No, default `0x0000000071727De22E5E9d8BAf0edAc6f37da032`.
- **Example:** `0x0000000071727De22E5E9d8BAf0edAc6f37da032`
- **How to get it:** Keep the default.

#### `PAYMASTER_ADDRESS`

- **What it is:** The address of the deployed `TrovoTokenPaymaster`.
- **Why it's needed:** Quotes are for this paymaster, and the service reads
  its tokens, signers and deposit.
- **Required:** Yes.
- **Example:** `0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9`
- **How to get it:** Printed by the deploy script and saved in
  `contracts/deployments/<chainId>.json`.

#### `QUOTE_SIGNER_PRIVATE_KEY`

- **What it is:** The private key that signs gas quotes.
- **Why it's needed:** The paymaster only pays for operations with a quote
  signed by one of its quote signers. The service refuses to start if this
  key's address is not one of them. It needs no ETH.
- **Required:** Yes.
- **Example:** `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
  (a well-known test key; never use it for real)
- **How to get it:** Generate it ([DEPLOYMENT.md step 4](DEPLOYMENT.md#4-create-the-keys-and-addresses))
  or with `openssl rand -hex 32` (prefix `0x`). Keep it only in Vault, and
  put its address in the deploy script's
  [`QUOTE_SIGNER_ADDRESSES`](#quote_signer_addresses).

#### `API_KEYS`

- **What it is:** The keys callers must send in the `X-API-Key` header,
  separated by commas (at least 16 characters each).
- **Why it's needed:** Only Trovo's own backend should get quotes. Several
  keys allow rotating one without downtime.
- **Required:** Yes.
- **Example:** `3f9c1d5e8a7b4c2d9e0f1a2b3c4d5e6f,8a02b7c4d1e9f3a6`
- **How to get it:** `openssl rand -hex 32`. Give app-backend its key as
  `PAYMASTER_QUOTE_SERVICE_API_KEY`
  ([app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)).

### Prices and spread

#### `GAS_TOKENS` (quote service)

- **What it is:** A JSON list of the gas tokens the service prices, how to
  price each one and the spread to add.
- **Why it's needed:** It says which stablecoins users can pay gas with and
  how to turn the price of ETH into each.
- **Required:** Yes.
- **Example:**
  ```json
  [
    {"symbol": "USDC", "address": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "decimals": 6, "route": ["ETH/USD", "USD/USDC"], "spreadBps": 100},
    {"symbol": "USDT", "address": "<USDT on Base>", "decimals": 6, "route": ["ETH/USD", "USD/USDT"], "spreadBps": 100},
    {"symbol": "cNGN", "address": "<cNGN on Base>", "decimals": 6, "route": ["ETH/USD", "USD/NGN", "NGN/CNGN"], "spreadBps": 150}
  ]
  ```
- **How to get it:** For each token:

  | Field | What to put | How to get it |
  |---|---|---|
  | `symbol` | the name used in requests (`"token": "cNGN"`) and logs | your choice; app-backend uses the same names |
  | `address` | the token contract; it must be enabled on the paymaster | the issuer's documentation |
  | `decimals` | the token's decimals | `cast call <token> "decimals()(uint8)" --rpc-url <RPC_URL>`, or basescan's **Read Contract** tab |
  | `route` | how to turn 1 ETH into this token: a chain of pairs from [`RATE_PAIRS`](#rate_pairs), starting at `ETH`; a step may be the reverse of a configured pair (`USD/USDC` uses `USDC/USD`) | from the pairs you configure |
  | `spreadBps` | this token's spread in basis points (100 = 1%); overrides `DEFAULT_SPREAD_BPS` | a business decision (see below) |

#### `RATE_PAIRS`

- **What it is:** A JSON object of exchange-rate pairs (`"BASE/QUOTE"`,
  meaning 1 BASE = x QUOTE) and the sources each one is read from.
- **Why it's needed:** It is where prices come from. Each pair's rate is the
  median of its sources, with outliers dropped, so one bad source cannot
  move the price.
- **Required:** Yes.
- **Example:**
  ```json
  {
    "ETH/USD": {
      "minSources": 2,
      "maxDeviationBps": 200,
      "sources": [
        {"name": "chainlink", "type": "chainlink", "feed": "<Chainlink ETH/USD feed on Base>", "maxAge": "30m"},
        {"name": "coingecko", "type": "http-json", "url": "https://pro-api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd", "path": "ethereum.usd", "headers": {"x-cg-pro-api-key": "vault:COINGECKO_API_KEY"}}
      ]
    },
    "USD/NGN": {
      "minSources": 1,
      "maxDeviationBps": 300,
      "sources": [
        {"name": "cngn-usdc-pool", "type": "uniswap-v3-twap", "pool": "<cNGN/USDC pool on Base>", "baseToken": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "window": "30m"},
        {"name": "exchange", "type": "http-json", "url": "<exchange ticker URL for USDT/NGN>", "path": "<path to the last price>", "headers": {"Authorization": "vault:EXCHANGE_API_KEY"}},
        {"name": "fx", "type": "http-json", "url": "https://open.er-api.com/v6/latest/USD", "path": "rates.NGN"}
      ]
    },
    "NGN/CNGN": {"sources": [{"name": "peg", "type": "fixed", "value": "1"}]},
    "USDC/USD": {"sources": [{"name": "peg", "type": "fixed", "value": "1"}]},
    "USDT/USD": {"sources": [{"name": "peg", "type": "fixed", "value": "1"}]}
  }
  ```
- **How to get it:** Choose the sources per pair from the types below.

  **Pair options:**

  | Field | Default | Meaning |
  |---|---|---|
  | `sources` | — | one or more sources (below) |
  | `minSources` | `1` | how many sources must answer (within the deviation limit) for the pair to be usable; use 2 or more for market prices once you have two good sources |
  | `maxDeviationBps` | `500` | sources further than this from the median are ignored as outliers (and flagged in `GET /v1/rates`); a negative value turns the check off |

  **Source types.** Every source has a `name` (unique within its pair) and
  a `type`, and may add `"invert": true` (use 1/x, for a feed quoted the
  other way round) and `"multiplier": "<decimal>"` (scale the answer):

  | `type` | Fields | Use it for | Where to find the values |
  |---|---|---|---|
  | `fixed` | `value` | a peg (`NGN/CNGN` = 1, `USDC/USD` = 1) or a temporary manual price | — |
  | `http-json` | `url`, `path` (dot path; numbers index lists, e.g. `data.0.price`), `headers` | any JSON API: exchange tickers, FX feeds, CoinGecko, the stablecoin issuer's rate | the API's documentation; a header value `vault:KEY` is read from key `KEY` of this same secret, so API keys stay in Vault |
  | `chainlink` | `feed`, `maxAge` (default `1h`) | a Chainlink price feed | feed addresses at <https://data.chain.link> (choose Base); set `maxAge` a little above the feed's heartbeat |
  | `uniswap-v3-twap` | `pool`, `baseToken`, `window` (default `30m`, 1m–168h) | the time-averaged price in any Uniswap v3-style pool (Uniswap v3, Aerodrome Slipstream, PancakeSwap v3); with `baseToken` = USDC on a cNGN/USDC pool this is the market USD/NGN rate at which gas actually converts | the pool address from the DEX's site or basescan; pick a pool with real liquidity; a longer window is harder to manipulate |

  API keys referenced as `vault:KEY` (for example `COINGECKO_API_KEY`,
  `EXCHANGE_API_KEY`) are extra keys of this secret; get them from each
  provider's dashboard. A new source type is one Go file in
  `quote-service/internal/rates` that calls `Register("my-type", factory)`.

#### `DEFAULT_SPREAD_BPS`

- **What it is:** The spread for tokens without their own `spreadBps`, in
  basis points (100 = 1%).
- **Why it's needed:** The paymaster pays gas in ETH but is paid in
  stablecoins, so each quote is a small currency exchange. The spread
  covers price movement until conversion, conversion costs and margin.
  The user pays the market price of gas plus this.
- **Required:** No, default `0`. Maximum `5000`.
- **Example:** `100`
- **How to get it:** A business decision; 100–150 is typical.

#### `RATE_REFRESH_INTERVAL`

- **What it is:** How often every rate source is read again.
- **Why it's needed:** Keeps prices current without overloading paid APIs.
- **Required:** No, default `30s`.
- **Example:** `30s`
- **How to get it:** Keep the default, or lengthen it if an API's rate
  limit is low.

#### `RATE_MAX_AGE`

- **What it is:** The oldest a pair's last good price may be before it is
  treated as unusable.
- **Why it's needed:** Quoting on a stale price could charge users too
  little or too much; with no fresh price, quotes for that token are
  refused (`rate-unavailable`).
- **Required:** No, default `5m`.
- **Example:** `5m`
- **How to get it:** Keep the default.

#### `RATE_SOURCE_TIMEOUT`

- **What it is:** How long to wait for one source to answer.
- **Why it's needed:** A slow API should not hold up the others.
- **Required:** No, default `10s`.
- **Example:** `5s`
- **How to get it:** Keep the default.

### Quotes

#### `QUOTE_VALIDITY`

- **What it is:** How long a quote stays valid when the caller does not
  ask for a particular time.
- **Why it's needed:** Long enough for a user to sign and send, short
  enough that the price does not drift far.
- **Required:** No, default `10m`.
- **Example:** `10m`
- **How to get it:** Keep the default.

#### `QUOTE_MAX_VALIDITY`

- **What it is:** The longest validity a caller may ask for.
- **Why it's needed:** Shared wallets' operations wait for approvers, so
  they need longer quotes. The rate is fixed for the whole time, so the
  spread must cover price movement over it.
- **Required:** No, default `24h`, at most `168h`.
- **Example:** `24h`
- **How to get it:** Keep the default.

#### `QUOTE_CLOCK_SKEW`

- **What it is:** How far before "now" a quote starts being valid.
- **Why it's needed:** Server and blockchain clocks differ slightly; this
  stops a fresh quote being rejected as "not valid yet".
- **Required:** No, default `60s`.
- **Example:** `60s`
- **How to get it:** Keep the default.

#### `PAYMASTER_VERIFICATION_GAS_LIMIT`

- **What it is:** The gas limit for the paymaster's checks, written into
  every quote (and signed).
- **Why it's needed:** The bundler and EntryPoint need it; too low and
  operations fail.
- **Required:** No, default `150000`.
- **Example:** `150000`
- **How to get it:** Keep the default.

#### `PAYMASTER_POST_OP_GAS_LIMIT`

- **What it is:** The gas limit for the paymaster's step after the
  operation (charging the token).
- **Why it's needed:** It must cover a token transfer and be at least the
  contract's `postOpGasOverhead`.
- **Required:** No, default `80000`.
- **Example:** `80000`
- **How to get it:** Keep the default.

### Deposit monitoring and operations

#### `DEPOSIT_CHECK_INTERVAL`

- **What it is:** How often the paymaster's ETH deposit is checked.
- **Why it's needed:** To alert before the deposit runs out.
- **Required:** No, default `1m`.
- **Example:** `1m`
- **How to get it:** Keep the default.

#### `DEPOSIT_LOW_WATERMARK_ETH`

- **What it is:** The deposit level, in ETH, below which an alert is sent.
- **Why it's needed:** When the deposit is empty, users can no longer pay
  gas in stablecoins. Quotes are only refused (`paymaster-deposit-low`)
  when the deposit cannot cover the specific operation.
- **Required:** No, default `0.05`.
- **Example:** `0.5`
- **How to get it:** Enough for a few hours of traffic, so treasury has
  time to top it up.

#### `ALERT_WEBHOOK_URL`

- **What it is:** A chat webhook (Discord or Slack) for low-deposit alerts.
- **Why it's needed:** So someone notices the deposit running low. The
  message carries both `content` (Discord) and `text` (Slack).
- **Required:** No. When empty, alerts only go to the log.
- **Example:** `https://discord.com/api/webhooks/123456789/abcdef`
- **How to get it:** Discord: channel **Settings → Integrations → Webhooks
  → New Webhook → Copy Webhook URL**. Slack: create an app with **Incoming
  Webhooks** and copy its URL.

#### `ALERT_COOLDOWN`

- **What it is:** The shortest time between two low-deposit alerts.
- **Why it's needed:** Stops the channel being flooded while the deposit
  stays low.
- **Required:** No, default `1h`.
- **Example:** `1h`
- **How to get it:** Keep the default.

#### `HTTP_ADDR`

- **What it is:** The address and port the service listens on.
- **Why it's needed:** app-backend calls the service here. Changing it
  needs a restart.
- **Required:** No, default `:8090` (port 8090 on every network interface).
- **Example:** `:8090`
- **How to get it:** Keep the default unless the port is taken.

#### `CONFIG_RELOAD_INTERVAL`

- **What it is:** How often to re-read the secret and apply changes
  (tokens, spreads, sources, keys) without a restart.
- **Why it's needed:** Lets you change spreads or rotate keys with no
  downtime. Sending the process `SIGHUP` also reloads. An invalid new
  configuration is logged and the running one is kept.
- **Required:** No, default `0s` (off).
- **Example:** `5m`
- **How to get it:** `5m` in production is a good choice.

---

## 4. Local development

### `PAYMASTER_CONFIG_SOURCE=env`

Runs the quote service without Vault, with every key as an environment
variable:

```bash
cd paymaster/quote-service
export PAYMASTER_CONFIG_SOURCE=env RPC_URL=http://127.0.0.1:8545 CHAIN_ID=31337 \
  PAYMASTER_ADDRESS=0x... QUOTE_SIGNER_PRIVATE_KEY=0x... API_KEYS=local-dev-key-0123456789 \
  GAS_TOKENS='[{"symbol":"USDC","address":"0x...","decimals":6,"route":["ETH/USD"]}]' \
  RATE_PAIRS='{"ETH/USD":{"sources":[{"type":"fixed","value":"3000"}]}}'
go run .
```

### Local stack (`contracts/scripts/local-stack.js`)

These are read only by the local test stack ([DEPLOYMENT.md](DEPLOYMENT.md#local-development-stack)).

#### `LOCAL_RPC_URL`

- **What it is:** The address of the local test blockchain.
- **Why it's needed:** Hardhat's `localhost` network uses it.
- **Required:** No, default `http://127.0.0.1:8545`.
- **Example:** `http://127.0.0.1:8545`
- **How to get it:** Keep the default (where `npx hardhat node` listens).

#### `DEV_BUNDLER_PORT`

- **What it is:** The port of the small development bundler the local
  stack starts.
- **Why it's needed:** app-backend sends wallet operations to it in local
  tests (its `BUNDLER_URL`).
- **Required:** No, default `4337`.
- **Example:** `4337`
- **How to get it:** Keep the default unless the port is taken.

#### `QUOTE_SIGNER_ADDRESS`

- **What it is:** The quote signer the local paymaster accepts.
- **Why it's needed:** The local quote service signs with the matching key.
- **Required:** No, default `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
  (Hardhat's test account #1, whose key is
  `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`).
- **Example:** `0x70997970C51812dc3A010C7d01b50e0d17dc79C8`
- **How to get it:** Keep the default.

---

## 5. Writing the secrets and the Vault policies

Write each secret from a JSON file:

```bash
vault kv put -mount=secret trovo/paymaster/quote-service @quote-service.json
vault kv put -mount=secret trovo/paymaster/deploy @deploy.json
```

`quote-service.json` looks like:

```json
{
  "RPC_URL": "https://base-mainnet.g.alchemy.com/v2/your-api-key",
  "CHAIN_ID": "8453",
  "PAYMASTER_ADDRESS": "0x...",
  "QUOTE_SIGNER_PRIVATE_KEY": "0x...",
  "API_KEYS": "...",
  "GAS_TOKENS": "[{\"symbol\":\"USDC\",\"address\":\"0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913\",\"decimals\":6,\"route\":[\"ETH/USD\",\"USD/USDC\"],\"spreadBps\":100}]",
  "RATE_PAIRS": "{\"ETH/USD\":{\"sources\":[{\"name\":\"chainlink\",\"type\":\"chainlink\",\"feed\":\"0x...\"}]},\"USDC/USD\":{\"sources\":[{\"name\":\"peg\",\"type\":\"fixed\",\"value\":\"1\"}]}}",
  "DEFAULT_SPREAD_BPS": "100",
  "CONFIG_RELOAD_INTERVAL": "5m"
}
```

**Every value must be a string.** `GAS_TOKENS` and `RATE_PAIRS` are JSON
*inside a string* (with `\"` escapes, as above). If you write them as
nested JSON objects instead, Vault stores them as objects and the service
cannot read them. A quick way to produce the escaped string from a
readable file: `jq -c . gas-tokens.json | jq -R .`

Read-only policy for the quote service (save as
`paymaster-quote-service.hcl`, load with
`vault policy write paymaster-quote-service paymaster-quote-service.hcl`):

```hcl
path "secret/data/trovo/paymaster/quote-service" {
  capabilities = ["read"]
}
```

The deploy secret holds a funded private key: give it its own policy
(`path "secret/data/trovo/paymaster/deploy" { capabilities = ["read"] }`),
only to whoever runs deployments, and delete or rotate
`DEPLOYER_PRIVATE_KEY` after the deployment.

## What other projects need

| Project | Parameter | Value |
|---|---|---|
| app-backend | `PAYMASTER_ADDRESS` | the deployed paymaster |
| app-backend | `PAYMASTER_QUOTE_SERVICE_URL` | where the quote service listens, e.g. `http://quote-service.internal:8090` |
| app-backend | `PAYMASTER_QUOTE_SERVICE_API_KEY` | one of [`API_KEYS`](#api_keys) |
| app-backend | `BUNDLER_URL` | your bundler (see [INTEGRATION.md](INTEGRATION.md#4-the-bundler)) |

See [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md).

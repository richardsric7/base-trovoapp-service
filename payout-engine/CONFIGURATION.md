# Configuration

Every parameter `payout-engine` reads. Each one says what it is, why it is
needed, whether you must set it, an example value and how to get a real
value. Examples are never real secrets.

## How to set them

The engine reads **environment variables**. Two ways to provide them:

- **A `.env` file** next to the binary (or in the directory you run it
  from). Start from the template: `cp .env.example .env`, then edit it. The
  engine loads it automatically at start.
- **The container or hosting platform's environment settings** (for
  example `docker run --env-file .env`, or the "Environment variables" page
  of your hosting platform). In production, inject secrets this way rather
  than baking a `.env` file into an image.

If a **required** value is missing or invalid, the engine stops at start
and prints which one and why.

## The minimum to start

| Parameter | Why |
|---|---|
| [`DB_CONNECTION_STRING`](#db_connection_string) | the database the payouts live in |
| [`BASE_RPC_URL`](#base_rpc_url) | to read and write the blockchain |
| [`PROCEED_PAYOUT_SAFE_ADDRESS`](#proceed_payout_safe_address) | the wallet payouts are paid from |
| [`PROCEED_PAYOUT_SIGNERS`](#proceed_payout_signers) or [`PROCEED_PAYOUT_SIGNERS_VAULT_PATH`](#proceed_payout_signers_vault_path) | the keys that sign payments |

Everything else is optional or has a sensible default.

---

## Database

### `DB_CONNECTION_STRING`

- **What it is:** The address and login of app-backend's main Postgres
  database.
- **Why it's needed:** Payouts, holder schedules, batches and approvals are
  all stored there, and the engine also reads assets, users, wallets, fees
  and country settings from it. Without it the engine cannot start.
- **Required:** Yes.
- **Example:** `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`
- **How to get it:** Use exactly the same value as app-backend's
  [`DB_CONNECTION_STRING`](../app-backend/CONFIGURATION.md). Deploy
  app-backend first: it creates the payout tables when it starts.

### `DB_TYPE`

- **What it is:** Which kind of database `DB_CONNECTION_STRING` points to.
- **Why it's needed:** Production uses Postgres. `sqlite` lets you try the
  engine against a local file without a database server (tests and
  experiments only).
- **Required:** No, default `postgres`.
- **Example:** `postgres`
- **How to get it:** Leave it unset in every real deployment.

### `DB_MAX_OPEN_CONNECTIONS`

- **What it is:** The most database connections the engine keeps open at
  once.
- **Why it's needed:** It stops the engine from using up the database's
  connection limit, which it shares with app-backend and tm-api.
- **Required:** No, default `10`.
- **Example:** `10`
- **How to get it:** Keep `10` unless your database administrator asks for
  fewer.

## Blockchain

### `BASE_RPC_URL`

- **What it is:** The web address of a Base blockchain node (a "JSON-RPC
  endpoint").
- **Why it's needed:** The engine reads token transfers to work out who
  holds a token, checks balances, estimates fees and sends the payment
  transactions through it. It must allow log queries (`eth_getLogs`) over
  [`PROCEED_PAYOUT_LOG_CHUNK`](#proceed_payout_log_chunk) blocks at a time.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
- **How to get it:** Use the same provider as app-backend's `BASE_RPC_URL`,
  or create a key with a node provider: sign up at a provider such as
  Alchemy, Infura or QuickNode, create an app for **Base Mainnet** (or
  **Base Sepolia** for testing) and copy its HTTPS URL.

### `BASE_CHAIN_ID`

- **What it is:** The number that identifies the blockchain network.
- **Why it's needed:** Every transaction is signed for one network, so a
  transaction cannot be replayed on another. Setting it also guards against
  pointing the engine at the wrong network by mistake.
- **Required:** No. When unset, the engine asks the node.
- **Example:** `8453` (Base Mainnet), `84532` (Base Sepolia)
- **How to get it:** `8453` for production, `84532` for the test network.

### `PROCEED_PAYOUT_SAFE_ADDRESS`

- **What it is:** The address of the Safe (a multi-signature wallet) that
  pays proceeds to token holders.
- **Why it's needed:** Each payout is funded by sending the payout token
  (for example cNGN) to this Safe, and the holders, the fee and the VAT are
  paid from it. Nothing can be paid without it.
- **Required:** Yes.
- **Example:** `0x0900c8af28C7CA2CAd2Fc7941E6d9d5A1b2Ae72e`
- **How to get it:**
  1. Generate the signer keys first (see
     [`PROCEED_PAYOUT_SIGNERS`](#proceed_payout_signers)).
  2. Open [Safe{Wallet}](https://app.safe.global), connect a wallet on
     Base, choose **Create account**, add the signer addresses as owners
     and set the threshold to **3**.
  3. Copy the new Safe's address.
  Use this Safe for payouts only.

### `SAFE_MULTISEND_CALL_ONLY_ADDRESS`

- **What it is:** The address of Safe's standard "MultiSendCallOnly"
  helper contract.
- **Why it's needed:** The engine pays many holders in one transaction by
  bundling the transfers through this contract.
- **Required:** No, default `0x9641d764fc13c8B624c04430C7356C1C7C8102e2`
  (the official Safe v1.4.1 deployment, the same on every network).
- **Example:** `0x9641d764fc13c8B624c04430C7356C1C7C8102e2`
- **How to get it:** Leave it unset. Only a private test chain with its own
  Safe deployment needs another value (the address from that deployment).

### `OFFER_BOOK_ADDRESS`

- **What it is:** The address of Trovo's offer book contract
  (`TrovoOfferBook`), where tokens being sold are held.
- **Why it's needed:** Tokens in an open sell offer sit in the offer book,
  not in the seller's wallet. With this set, the engine pays those tokens'
  share to the sellers instead of to the offer book.
- **Required:** No. Without it, tokens in open offers earn nothing.
- **Example:** `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- **How to get it:** The same value as app-backend's `OFFER_BOOK_ADDRESS`;
  it is printed when the market contracts are deployed (see
  [market/DEPLOYMENT.md](../market/DEPLOYMENT.md)).

### `PROCEED_PAYOUT_SWEEP_ADDRESS`

- **What it is:** The wallet that leftover funds are moved to when an
  admin presses **Sweep** in Trovo Manager.
- **Why it's needed:** After payouts finish, the payout Safe can hold
  leftovers (shares of excluded holders, rounding). A sweep moves a token's
  whole balance here. It is refused while a payout of that token is being
  funded or paid.
- **Required:** No. Without it, sweeps are refused.
- **Example:** `0x8ba1f109551bD432803012645Ac136ddd64DBA72`
- **How to get it:** The address of a treasury wallet or Safe your finance
  team controls.

### `PROCEED_PAYOUT_EXCLUDED_ADDRESSES`

- **What it is:** A list of wallet addresses that are never paid.
- **Why it's needed:** Some holders, such as company wallets, should not
  receive proceeds. These are already excluded automatically: the payout
  Safe, the asset's market-making wallet, the issuing profile's wallets
  and the offer book. Use this for any others.
- **Required:** No.
- **Example:** `0x1111111111111111111111111111111111111111,0x2222222222222222222222222222222222222222`
- **How to get it:** Ask the business which wallets must not be paid.
  Separate addresses with commas, semicolons or spaces.

### `TOKENIZATION_ISSUING_PROFILE`

- **What it is:** The username of the Trovo account that issues tokenized
  assets.
- **Why it's needed:** That account's wallets hold unsold tokens; they are
  excluded so the company does not pay itself.
- **Required:** No, default `atprofile`.
- **Example:** `atprofile`
- **How to get it:** The same value as app-backend's
  `TOKENIZATION_ISSUING_PROFILE`.

### `VAT_WALLET`

- **What it is:** The wallet the VAT on payout fees is paid to.
- **Why it's needed:** When a payout charges a fee and the asset's country
  charges VAT (`country_configs.vat_percent`), the VAT part is paid here.
- **Required:** Only when payouts have a fee in a country with VAT. When
  unset, the wallet in the `VAT` row of app-backend's `service_fees` table
  is used.
- **Example:** `0x2546BcD3c84621e976D8185a91A922aE77ECEc30`
- **How to get it:** The same value as app-backend's `VAT_WALLET`, given by
  your finance team. It must not hold the asset's token.

The payout **fee** wallet and default fee are not environment variables:
set them in Trovo Manager under **Proceeds Payouts → Fee & engine**.

## Signing keys

The payout Safe needs 3 signatures for every payment. The keys are kept as
one secret: at least 3 private keys separated by `;`. The first three sign,
and the first one also sends the transaction and pays its network fee in
ETH, so keep that wallet topped up (the funding check shows how much it
needs).

Give the keys **one** of two ways: directly in
`PROCEED_PAYOUT_SIGNERS`, or by telling the engine where to read them in
Vault with `PROCEED_PAYOUT_SIGNERS_VAULT_PATH`. Vault is better: a key
rotated in Trovo Manager's vault manager is picked up without a restart.

### `PROCEED_PAYOUT_SIGNERS`

- **What it is:** The payout Safe's signing keys: at least 3 private keys
  separated by `;`.
- **Why it's needed:** Without them the engine cannot approve or send any
  payment from the Safe.
- **Required:** Yes, unless `PROCEED_PAYOUT_SIGNERS_VAULT_PATH` is set.
- **Example:** `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d;0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a;0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6`
  (well-known test keys, never use them for real funds)
- **How to get it:**
  1. In Trovo Manager open the vault manager (Vault Signer) and create a
     managed secret with at least 4 keys; set its wallet address to the
     payout Safe.
  2. Have your deployment inject that secret into this variable.
  Never commit it to git or put it in a shared document.

### `PROCEED_PAYOUT_SIGNERS_VAULT_PATH`

- **What it is:** Where in HashiCorp Vault the signer secret is stored,
  written `mount/path#field`.
- **Why it's needed:** With it, the engine reads the keys from Vault each
  time it signs, so a rotated key is used straight away and the keys never
  sit in the environment.
- **Required:** No (use it instead of `PROCEED_PAYOUT_SIGNERS`). When set,
  `VAULT_ADDR` and `VAULT_TOKEN` are required.
- **Example:** `secret/trovo/payout-engine#PROCEED_PAYOUT_SIGNERS`
  (mount `secret`, path `trovo/payout-engine`, field
  `PROCEED_PAYOUT_SIGNERS`; the field part is optional and defaults to
  `PROCEED_PAYOUT_SIGNERS`)
- **How to get it:** The path the vault manager writes the managed secret
  to; ask whoever runs Vault, or read it in Vault's UI under the KV engine.

### `VAULT_ADDR`

- **What it is:** The web address of your Vault server.
- **Why it's needed:** To read the signer secret from Vault.
- **Required:** Only with `PROCEED_PAYOUT_SIGNERS_VAULT_PATH`.
- **Example:** `https://vault.internal.trovo.io:8200`
- **How to get it:** The same value tm-api's vault manager uses
  (`VAULT_ADDR` in [tm-api/CONFIGURATION.md](../tm-api/CONFIGURATION.md)).

### `VAULT_TOKEN`

- **What it is:** A Vault access token.
- **Why it's needed:** Vault only returns the secret to a token whose
  policy allows reading that path.
- **Required:** Only with `PROCEED_PAYOUT_SIGNERS_VAULT_PATH`.
- **Example:** `hvs.CAESIJexampleexampleexample`
- **How to get it:** Ask your Vault administrator for a token with a
  read-only policy on the signer path, for example:
  `vault token create -policy=payout-engine-read -period=768h`.

### `VAULT_NAMESPACE`

- **What it is:** The Vault namespace the secret lives in.
- **Why it's needed:** Only Vault Enterprise and HCP Vault use namespaces;
  without the right one the path is not found.
- **Required:** No.
- **Example:** `admin`
- **How to get it:** Your Vault administrator; leave unset on open-source
  Vault.

## Redis (optional)

Redis lets Trovo Manager wake the engine straight away when an admin acts,
and lets the engine report progress and its heartbeat. Without Redis the
engine still works; it notices admin actions on its next database check
(every [`PROCEED_PAYOUT_POLL_SECONDS`](#proceed_payout_poll_seconds)).

### `REDIS_HOST`

- **What it is:** The host name of the Redis server.
- **Why it's needed:** Turns on instant wake-ups, progress events and the
  heartbeat. tm-api publishes the wake-ups, so it must be the **same
  Redis** tm-api uses.
- **Required:** No.
- **Example:** `redis.internal`
- **How to get it:** The same value as tm-api's `REDIS_HOST`.

### `REDIS_PORT`

- **What it is:** The Redis server's port.
- **Why it's needed:** To connect to Redis.
- **Required:** No, default `6379`.
- **Example:** `6379`
- **How to get it:** The same value as tm-api's `REDIS_PORT`.

### `REDIS_PASSWORD`

- **What it is:** The Redis password.
- **Why it's needed:** A protected Redis refuses connections without it.
- **Required:** Only if your Redis has a password.
- **Example:** `change-me`
- **How to get it:** The same value as tm-api's `REDIS_PASSWORD`.

### `REDIS_TLS`

- **What it is:** Whether to connect to Redis over an encrypted (TLS)
  connection.
- **Why it's needed:** Managed Redis services usually require TLS.
- **Required:** No, default off. Set `1` to turn it on.
- **Example:** `1`
- **How to get it:** `1` if your Redis address starts with `rediss://` or
  the provider says TLS is required; otherwise leave unset.

## Push notifications

### `GC`

- **What it is:** The Firebase service account key (a JSON file), encoded
  as base64.
- **Why it's needed:** The engine sends each paid Trovo user one push
  notification per payout. Without it, payments still happen but nobody
  is notified.
- **Required:** No.
- **Example:** `ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAidHJvdm8tZXhhbXBsZSIKfQ==`
- **How to get it:** Use the same value as app-backend's `GC`. To make one:
  Firebase console → Project settings → Service accounts → **Generate new
  private key**, then `base64 -w0 key.json`.

## Tuning

These have safe defaults. Change them only for a reason given below.

### `PROCEED_PAYOUT_CONFIRMATIONS`

- **What it is:** How many blocks behind the newest block the holder
  snapshot stays.
- **Why it's needed:** Very recent blocks can occasionally be replaced
  (a "reorg"); staying a few blocks behind avoids paying holders based on
  a transfer that disappears.
- **Required:** No, default `3`.
- **Example:** `3`
- **How to get it:** Keep the default on Base.

### `PROCEED_PAYOUT_LOG_CHUNK`

- **What it is:** How many blocks the engine reads token transfers for in
  one request.
- **Why it's needed:** Node providers limit how big a log query may be. If
  yours rejects queries ("block range too large"), lower this.
- **Required:** No, default `5000`.
- **Example:** `2000`
- **How to get it:** Your RPC provider's documented `eth_getLogs` block
  range limit.

### `PROCEED_PAYOUT_CATCH_UP_WITHIN`

- **What it is:** How close (in blocks) to the safe head the scan must be
  before the schedule is locked.
- **Why it's needed:** New blocks keep arriving while the engine scans;
  this decides when it is "caught up enough" to lock the list of holders.
- **Required:** No, default `2`.
- **Example:** `2`
- **How to get it:** Keep the default.

### `PROCEED_PAYOUT_BATCH_SIZE`

- **What it is:** How many holder payments go into one transaction.
- **Why it's needed:** Bigger batches mean fewer transactions and lower
  total fees, but each must fit in a block. Lower it if batches fail for
  running out of gas.
- **Required:** No, default `150` (minimum `1`).
- **Example:** `150`
- **How to get it:** Keep the default unless batches fail.

### `PROCEED_PAYOUT_MAX_BATCH_GAS`

- **What it is:** The most gas one batch may be estimated to use.
- **Why it's needed:** A batch estimated above this is split in two, so
  every transaction fits comfortably in a block.
- **Required:** No, default `12000000`.
- **Example:** `12000000`
- **How to get it:** Keep the default.

### `PROCEED_PAYOUT_MIN_UNITS`

- **What it is:** The smallest share worth paying, in the token's smallest
  unit.
- **Why it's needed:** Tiny shares can cost more in fees than they are
  worth; shares below this are skipped.
- **Required:** No, default `1` (pay everything above zero).
- **Example:** `1000000000000000` (0.001 of an 18-decimal token)
- **How to get it:** Decide the smallest amount worth sending and multiply
  by 10 to the power of the token's decimals.

### `PROCEED_PAYOUT_GAS_PRICE_MULTIPLIER`

- **What it is:** A multiplier on the network fee the node suggests.
- **Why it's needed:** Paying a little above the suggestion gets batches
  mined promptly when the network is busy.
- **Required:** No, default `1.25`. Values below `1` are ignored.
- **Example:** `1.5`
- **How to get it:** Raise it only if batches often wait a long time to be
  mined.

### `PROCEED_PAYOUT_POLL_SECONDS`

- **What it is:** How often (in seconds) the engine checks the database
  for work.
- **Why it's needed:** Without Redis this is how quickly it notices admin
  actions; with Redis it is a safety net.
- **Required:** No, default `10`.
- **Example:** `10`
- **How to get it:** Keep the default.

### `PROCEED_PAYOUT_RECEIPT_TIMEOUT_SECONDS`

- **What it is:** How long (in seconds) the engine waits for a sent batch
  to be mined.
- **Why it's needed:** After this, it stops waiting and checks the batch
  again on its next round: settled, dropped or sent again. Nothing is paid
  twice.
- **Required:** No, default `180`.
- **Example:** `180`
- **How to get it:** Keep the default.

### `PROCEED_PAYOUT_DEFAULT_START_BLOCK`

- **What it is:** The earliest block a token's first scan may start from.
- **Why it's needed:** A token's first scan starts an hour before the
  asset was created, or at this block if later. Setting it to around the
  time Trovo went live on the network skips empty history and makes the
  first scan faster.
- **Required:** No, default `0`.
- **Example:** `21000000`
- **How to get it:** On [basescan.org](https://basescan.org) look up the
  block number of the first token deployment, or leave `0`.

## Identity

### `INSTANCE_ID`

- **What it is:** A name for this running copy of the engine.
- **Why it's needed:** Only one copy works at a time; Trovo Manager shows
  which one holds the heartbeat, so a clear name helps when you run a
  standby.
- **Required:** No, default the machine's host name.
- **Example:** `payout-engine-1`
- **How to get it:** Any short name, different for each copy.

### `APP_VERSION`

- **What it is:** The version of the engine that is running.
- **Why it's needed:** Shown in Trovo Manager's engine status so you can
  tell which build is live.
- **Required:** No, default `dev`.
- **Example:** `2026.10.05-1161fa1`
- **How to get it:** Set it in your deployment to the release tag or git
  commit (`git rev-parse --short HEAD`).

## Used only by tests

`AA_LOCAL_STACK` and `AA_LOCAL_MARKET` point the end-to-end test at a local
blockchain; see [DEPLOYMENT.md](DEPLOYMENT.md#run-the-tests). The engine
itself never reads them.

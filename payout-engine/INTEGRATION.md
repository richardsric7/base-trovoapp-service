# Integration

Everything `payout-engine` connects to, and how. The engine has no API of
its own: nothing calls it directly. It works through a shared database,
Redis messages, the blockchain and two outside services.

| Connects to | Direction | Through | Needed? |
|---|---|---|---|
| [app-backend's database](#1-app-backends-database) | engine reads and writes | Postgres | yes |
| [tm-api (Trovo Manager)](#2-tm-api-trovo-manager) | tm-api changes payouts, engine does the work | the same database, plus Redis wake-ups | yes |
| [Redis](#3-redis) | tm-api → engine, engine → anyone listening | Redis publish/subscribe | optional |
| [Base blockchain](#4-base-blockchain) | engine reads and sends transactions | JSON-RPC over HTTPS | yes |
| [Vault and the vault manager](#5-vault-and-the-vault-manager) | engine reads the signing keys | Vault HTTP API | yes, or keys in the environment |
| [Firebase Cloud Messaging](#6-firebase-cloud-messaging) | engine sends push notifications | Firebase Admin SDK | optional |
| [market contracts](#7-market-contracts-offer-book) | engine reads open sell offers | the database and the chain | optional |

---

## 1. app-backend's database

- **What it is and why:** app-backend's main Postgres database. All payout
  data lives there (payouts, holder schedules, batches, approvals), and the
  engine reads the assets, users, wallets, fees and country settings it
  needs to compute each payout.
- **Direction:** the engine reads and writes. app-backend owns the tables
  (creates and changes them); tm-api also writes to them on admin actions.
- **How they connect:** a direct Postgres connection. Every status change,
  by the engine or by tm-api, is a conditional update (`... WHERE status =
  <expected>`), so if both act at the same moment one of them simply loses
  instead of overwriting the other.
- **Settings on this side:** [`DB_CONNECTION_STRING`](CONFIGURATION.md#db_connection_string)
  (example `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`).
- **Settings on the other side:** it must equal app-backend's
  `DB_CONNECTION_STRING` (see [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md))
  and tm-api's connection to the same database. Deploy app-backend first so
  the tables exist.
- **How to check it works:** the engine starts without a `database:` error,
  and Trovo Manager shows its heartbeat.
- **When it is down:** the engine cannot do anything; it logs the error and
  carries on from where it was when the database is back.

Which tables each side uses:

| Table | The engine | tm-api | app-backend |
|---|---|---|---|
| `proceed_payouts` | moves a payout through `PREPARING`, `LOCKED`, `PAYING`, `COMPLETED*` (and back to `APPROVED` when funding fails); fills in the snapshot, amounts, fee, VAT and counts | registers it (`REGISTERED`) when a trustee authorizes a distribution; moves it on admin actions (`PREPARE_REQUESTED`, `APPROVED`, `FUNDING_CHECK_REQUESTED`, `PAUSED`, `CANCELLED`, …); sets the fee | owns the table |
| `tokenized_asset_payout_schedules` | writes the holder list when locking; marks lines `QUEUED` / `PAID` / `FAILED` | admin exclusions, inclusions and mark-paid | `GET /v1/tokenization/payouts` lists a user's lines |
| `proceed_payout_approvals` | counts approvals for the current checksum at the funding check | records approvals (distinct admins) | |
| `proceed_payout_batches` | records each payment transaction (signed) before sending it, then settles it | shows them in Trovo Manager | |
| `payout_token_indexes`, `payout_token_balances` | the per-token record of transfers and balances | | |
| `payout_engine_states` (row 1) | heartbeat, current activity, last error; performs sweeps | kill switch (`halted`), sweep requests | |
| `fee_collections` | records the payout fee and its VAT when paid (`PROCEED_PAYOUT_FEE`, `PROCEED_PAYOUT_FEE_VAT`) | Trovo Manager's fee and VAT report | owns the table |
| `service_fees` | reads `PROCEED_PAYOUT_FEE` (fee wallet) and `VAT` | writes `PROCEED_PAYOUT_FEE` (fee wallet, default fee) | |
| `country_configs` | reads the asset country's `vat_percent` | | |
| `tokenized_assets`, `users`, `user_wallets`, `offer_book_offers` | reads the token, market-making wallet, issuing-profile wallets, usernames and push tokens, and open offers | | |

## 2. tm-api (Trovo Manager)

- **What it is and why:** tm-api is the backend of Trovo Manager, where
  admins run payouts. Every admin step (register, prepare, approve, confirm
  funding, pause, cancel, sweep, stop) is a tm-api endpoint under
  `/api/v1/proceed-payouts/…` that changes the database; the engine then
  does the work that step needs.
- **Direction:** tm-api → engine, through the database (the real
  instruction) and a Redis message (a wake-up so the engine does not wait
  for its next check).
- **How they connect:** no direct calls. tm-api writes the new status, then
  publishes `{"action": "prepare", "payoutId": 1}` on Redis channel
  `payout-engine:commands`. The engine re-reads the database when woken,
  or every [`PROCEED_PAYOUT_POLL_SECONDS`](CONFIGURATION.md#proceed_payout_poll_seconds)
  anyway.
- **Settings on this side:** [`DB_CONNECTION_STRING`](CONFIGURATION.md#db_connection_string),
  and the [Redis settings](CONFIGURATION.md#redis-optional).
- **Settings on the other side:** tm-api must use the same database and
  the same Redis (`REDIS_HOST`, `REDIS_PORT`, `REDIS_PASSWORD`, `REDIS_TLS`
  in [tm-api/CONFIGURATION.md](../tm-api/CONFIGURATION.md)). The number of
  approvals each payout needs is tm-api's `PROCEED_PAYOUT_APPROVALS_REQUIRED`
  (default 2).
- **How to check it works:** press **Prepare** on a payout in Trovo
  Manager; its status changes to **Preparing** within a second with
  Redis, or within 10 seconds without.
- **When it is down:** admins cannot change payouts, but the engine
  finishes the work already requested.

The stakeholder portal's distribution page and payout CSV also read the
payout through tm-api. A distribution follows its payout: `processing`
from preparation, `completed` when the payout completes, `failed` when it
is cancelled. See [tm-api/INTEGRATION.md](../tm-api/INTEGRATION.md).

## 3. Redis

- **What it is and why:** an in-memory message service. It carries
  wake-ups from tm-api and the engine's progress events and heartbeat.
- **Direction:**

  | Key / channel | From → to | Content |
  |---|---|---|
  | `payout-engine:commands` | tm-api → engine | `{"action": "prepare", "payoutId": 1}`; only wakes the engine, which re-reads the database |
  | `payout-engine:events` | engine → any listener | `{"type": "locked", "payoutId": 1, "message": "…", "at": "…"}`; type is one of `preparing`, `prepare-refused`, `locked`, `funded`, `funding-failed`, `batch-mined`, `batch-reverted`, `completed`, `sweep`, `error` |
  | `payout-engine:heartbeat` | engine → any reader | which instance is active, expiring on its own |

- **How they connect:** a Redis connection, optionally over TLS.
- **Settings on this side:** [`REDIS_HOST`](CONFIGURATION.md#redis_host)
  (example `redis.internal`), [`REDIS_PORT`](CONFIGURATION.md#redis_port)
  (`6379`), [`REDIS_PASSWORD`](CONFIGURATION.md#redis_password),
  [`REDIS_TLS`](CONFIGURATION.md#redis_tls).
- **Settings on the other side:** tm-api's Redis settings must point to
  the same server.
- **How to check it works:** `redis-cli -h <REDIS_HOST> SUBSCRIBE payout-engine:events`,
  then act on a payout in Trovo Manager; an event appears.
- **When it is down:** nothing breaks; the engine notices admin actions on
  its next database check instead.

## 4. Base blockchain

- **What it is and why:** the network the tokens and payments live on.
- **Direction:** the engine reads and sends.
  - **Reads:** token `Transfer` events (`eth_getLogs`) to work out holders,
    balances, token decimals, the Safe's nonce, owners and threshold, and
    network fees.
  - **Sends:** Safe `execTransaction` calls, sent by the first signer. Each
    one pays a batch of holders through Safe's `MultiSendCallOnly`
    contract. The engine only ever moves the payout token out of the
    payout Safe (and, for a sweep, to the configured sweep address).
- **How they connect:** JSON-RPC over HTTPS to a node provider.
- **Settings on this side:** [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url)
  (example `https://base-mainnet.g.alchemy.com/v2/your-api-key`),
  [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id) (`8453`),
  [`PROCEED_PAYOUT_SAFE_ADDRESS`](CONFIGURATION.md#proceed_payout_safe_address),
  [`SAFE_MULTISEND_CALL_ONLY_ADDRESS`](CONFIGURATION.md#safe_multisend_call_only_address).
- **Settings on the other side:** your node provider account (see
  [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url) for how to get one).
- **How to check it works:**
  `curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId"}' <BASE_RPC_URL>`
  returns `"result":"0x2105"` (8453) on Base Mainnet.
- **When it is down:** preparing and paying pause and resume when the node
  answers again. A batch sent just before is checked afterwards; nothing
  is paid twice.

## 5. Vault and the vault manager

- **What it is and why:** HashiCorp Vault stores the payout Safe's signing
  keys. Trovo Manager's vault manager (Vault Signer, in tm-api) creates and
  rotates them, and swaps the Safe's owner on the blockchain when it
  rotates one.
- **Direction:** the engine reads the keys from Vault each time it signs.
- **How they connect:** Vault's HTTP API with a token (KV version 2).
- **Settings on this side:** [`PROCEED_PAYOUT_SIGNERS_VAULT_PATH`](CONFIGURATION.md#proceed_payout_signers_vault_path)
  (example `secret/trovo/payout-engine#PROCEED_PAYOUT_SIGNERS`),
  [`VAULT_ADDR`](CONFIGURATION.md#vault_addr),
  [`VAULT_TOKEN`](CONFIGURATION.md#vault_token),
  [`VAULT_NAMESPACE`](CONFIGURATION.md#vault_namespace). Or skip Vault and
  inject the keys as [`PROCEED_PAYOUT_SIGNERS`](CONFIGURATION.md#proceed_payout_signers).
- **Settings on the other side:** the managed secret in the vault manager,
  with its wallet address set to the payout Safe; a Vault policy that lets
  the engine's token read that path.
- **How to check it works:** the start line says
  `signers from vault secret/...`, and the funding check in Trovo Manager
  confirms the signers own the Safe.
- **When it is down:** the engine cannot sign; payments wait and continue
  when Vault is back.

## 6. Firebase Cloud Messaging

- **What it is and why:** Google's push notification service, used by the
  Trovo mobile app.
- **Direction:** engine → Firebase → the user's phone. Each paid Trovo user
  gets one notification per payout, sent to the push token stored in
  `users`.
- **How they connect:** the Firebase Admin SDK with a service account key.
- **Settings on this side:** [`GC`](CONFIGURATION.md#gc).
- **Settings on the other side:** the same Firebase project as app-backend
  and the mobile app (app-backend's `GC`).
- **How to check it works:** after a test payout completes, the test user's
  phone shows the notification.
- **When it is down or unset:** payments still go out; the notifications
  are skipped.

## 7. market contracts (offer book)

- **What it is and why:** Trovo's offer book holds tokens that sellers have
  put up for sale. Those tokens still belong to the sellers, so their
  proceeds should go to the sellers.
- **Direction:** the engine reads open offers (`offer_book_offers` in the
  database, kept up to date by app-backend) and the offer book's address.
- **Settings on this side:** [`OFFER_BOOK_ADDRESS`](CONFIGURATION.md#offer_book_address).
- **Settings on the other side:** it must equal app-backend's
  `OFFER_BOOK_ADDRESS`, printed when the market contracts are deployed
  ([market/DEPLOYMENT.md](../market/DEPLOYMENT.md)).
- **How to check it works:** in a test payout, a holder with an open sell
  offer is paid for the tokens in the offer.
- **When it is unset:** tokens in open offers earn nothing in payouts.

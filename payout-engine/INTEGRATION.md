# Integration

payout-engine exposes no API. It shares app-backend's database with
app-backend (which owns the schema) and tm-api (which changes payouts on
admin actions), and listens on Redis for wake-ups.

## app-backend's database

| Table | The engine | tm-api | app-backend |
|---|---|---|---|
| `proceed_payouts` | moves a payout through `PREPARING`, `LOCKED`, `PAYING`, `COMPLETED*` (and back to `APPROVED` when funding fails); fills the snapshot, amounts, fee, VAT and counts | registers it (`REGISTERED`) when a trustee authorizes a distribution; moves it on admin actions (`PREPARE_REQUESTED`, `APPROVED`, `FUNDING_CHECK_REQUESTED`, `PAUSED`, `CANCELLED`, …); sets the fee | owns the model and migration |
| `tokenized_asset_payout_schedules` | writes the schedule when locking; marks lines `QUEUED` / `PAID` / `FAILED` | admin exclusions, inclusions and mark-paid | `GET /v1/tokenization/payouts` lists a user's lines |
| `proceed_payout_approvals` | checks the count for the current checksum at the funding check | records approvals (distinct admins) | |
| `proceed_payout_batches` | records each Safe transaction (with its signed transaction) before sending, then settles it | reads them for TM | |
| `payout_token_indexes`, `payout_token_balances` | the per-token Transfer replay | | |
| `payout_engine_states` (row 1) | heartbeat, activity, last error; performs sweeps | kill switch (`halted`), sweep requests | |
| `fee_collections` | records the payout fee and its VAT when paid (`PROCEED_PAYOUT_FEE`, `PROCEED_PAYOUT_FEE_VAT`) | TM's payout fee and VAT report | owns the table (TM's fee report) |
| `service_fees` | reads `PROCEED_PAYOUT_FEE` (fee wallet) and `VAT` | writes `PROCEED_PAYOUT_FEE` (fee wallet, default fee) | |
| `country_configs` | reads `vat_percent` of the asset's country | | |
| `tokenized_assets`, `users`, `user_wallets`, `offer_book_offers` | reads the token, market-making wallet, issuing-profile wallets, usernames (for notifications) and open offers | | |

Every status change on either side is a conditional update (`WHERE status
= <expected>`), so a change made at the same moment by the other side
loses cleanly instead of overwriting.

## Redis

| Key / channel | Direction | Content |
|---|---|---|
| `payout-engine:commands` | tm-api → engine | `{"action": "prepare", "payoutId": 1}` after each change. It only wakes the engine, which then re-reads the database. |
| `payout-engine:events` | engine → | `{"type": "locked", "payoutId": 1, "message": "…", "at": "…"}`, where the type is one of preparing, prepare-refused, locked, funded, funding-failed, batch-mined, batch-reverted, completed, sweep or error |
| `payout-engine:heartbeat` | engine → | the active instance, with a TTL |

Redis is optional. Without it the engine polls every
`PROCEED_PAYOUT_POLL_SECONDS`.

## tm-api

`/api/v1/proceed-payouts/…` (Trovo admins only, audited) changes the
tables above. The stakeholder portal's distribution detail and payout CSV
read the payout through tm-api's `proceedpayouts.DistributionClient`. A
distribution follows its payout: `processing` from preparation,
`completed` when the payout completes, `failed` when it is cancelled. See
tm-api's [INTEGRATION.md](../tm-api/INTEGRATION.md).

## Base

- Reads: `eth_getLogs` (Transfer events of the asset token), balances,
  `decimals`, the Safe's nonce, owners and threshold, gas prices.
- Writes: Safe `execTransaction` calls from the first signer. Each call
  delegate-calls `MultiSendCallOnly` with ERC-20 `transfer`s from the
  payout Safe. Nothing else is sent: the engine moves the payout token out
  of the payout Safe only (and the sweep, to the configured sweep address).

## Push notifications

Firebase (`GC`): one notification per paid Trovo user per payout, to the
user's push token from `users`.

## Signers

`PROCEED_PAYOUT_SIGNERS` is a managed secret of TM's vault manager (Vault
Signer). Rotating an entry there swaps that owner on the payout Safe. Read
from Vault (`PROCEED_PAYOUT_SIGNERS_VAULT_PATH`), the engine uses the new
keys on its next transaction.

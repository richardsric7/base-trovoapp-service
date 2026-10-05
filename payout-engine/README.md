# Payout Engine (`payout-engine`)

## What this project does

When a company whose asset is tokenized on Trovo shares out money (a
dividend, rent, interest), every person holding that asset's tokens should
get their share. `payout-engine` is the program that pays them. It works
out who held the tokens and how many, waits for Trovo admins to check and
approve the list, makes sure the money has arrived, and then sends each
holder their share in a stablecoin such as cNGN from one dedicated wallet
on the Base blockchain. Paid users get a push notification.

It runs in the background with **no website and no API**. Admins control it
from Trovo Manager (`tm-web` → `tm-api`); it shares app-backend's database,
and Redis (optional) lets Trovo Manager wake it instantly. Without it,
approved payouts are never paid.

- How to deploy it: [DEPLOYMENT.md](DEPLOYMENT.md)
- Every setting: [CONFIGURATION.md](CONFIGURATION.md)
- What it connects to: [INTEGRATION.md](INTEGRATION.md)

## How a payout works

A payout exists once a trustee authorizes a stakeholder distribution in the
stakeholder portal: `tm-api` registers it in `proceed_payouts` as
`REGISTERED`. Admins then move it through these stages, and the engine does
the work at each one:

| Stage (status) | Who | What happens |
|---|---|---|
| `REGISTERED` | admin: **Prepare** | |
| `PREPARE_REQUESTED` → `PREPARING` | engine | Replays the asset token's `Transfer` logs into a per-token balance index (`payout_token_indexes` / `payout_token_balances`), resuming where the last payout of that token stopped, and catches up to the chain head, so transfers made while preparing are included. It excludes platform wallets (issuing-profile Safes, the market-making wallet, the offer book, the payout Safe, `PROCEED_PAYOUT_EXCLUDED_ADDRESSES`) and credits offer-book escrow back to the sellers. Then it computes the fee and VAT and each holder's share, and locks the schedule (`tokenized_asset_payout_schedules`) with a checksum. |
| `LOCKED` | admins: **Approve** | `PROCEED_PAYOUT_APPROVALS_REQUIRED` distinct admins (default 2) approve this checksum. The admin who prepared the schedule, or set the fee, cannot approve. |
| `APPROVED` | admin: **Confirm funding** | after sending the payout to the payout Safe |
| `FUNDING_CHECK_REQUESTED` | engine | Checks that the schedule still matches its checksum, the approvals are present, the signers own the Safe, the Safe holds what is still owed (on top of other payouts in progress), and the executing signer has gas. On success it moves to `PAYING`; otherwise back to `APPROVED`, with the shortfall in the note. |
| `PAYING` | engine | Pays in `MultiSendCallOnly` batches (`PROCEED_PAYOUT_BATCH_SIZE`). Each batch is recorded with its signed transaction **before** it is broadcast. A failing transfer is isolated by bisection and marked `FAILED`; the rest are paid. Each paid Trovo user gets one push notification. |
| `COMPLETED` / `COMPLETED_WITH_FAILURES` | engine | Admins can **Retry failed** transfers, **Exclude** a holder, or **Mark paid** one paid another way. |

At any time admins can **Pause**, **Resume** (the funding check runs again),
**Cancel** before or between batches, **Stop** all payouts (the kill switch),
or **Sweep** a token from the payout Safe.

### Fee and VAT

Each payout has its own processing fee:

- **FIXED**: an amount of the payout token.
- **PERCENT**: a share of the payout, capped when the cap is above 0.

New payouts default to the fee configured in TM: the `PROCEED_PAYOUT_FEE`
row of `service_fees`, which also holds the fee wallet. With nothing
configured the fee is FIXED 0, cap 0.

VAT is charged on the fee at the asset country's `country_configs.vat_percent`.
The fee and VAT come out of the payout, before the holders' share, and are
paid with the holders to the fee wallet and the VAT wallet (`VAT_WALLET`, or
the `VAT` service fee). They are part of the approved schedule. Once paid,
they are recorded in `fee_collections` as `PROCEED_PAYOUT_FEE` and
`PROCEED_PAYOUT_FEE_VAT`, which feed TM's payout fee and VAT report.

### Safety

- **No double payment after a restart.** A batch is in the database, with
  its signed transaction, before it is sent. On restart the engine settles
  it from its receipt if it was mined. If the Safe nonce was used by
  another transaction, it marks the batch dropped and pays its items again.
  Otherwise it broadcasts the same transaction again.
- **No preparing a schedule twice once paying has started.** A schedule
  that has started paying is never prepared again; tm-api and the engine
  both refuse.
- **One active instance.** Only one engine instance works at a time (a
  heartbeat claim on `payout_engine_states`); others wait as standbys.
- **Rotation without a restart.** Signers come from the managed secret
  `PROCEED_PAYOUT_SIGNERS` (see [CONFIGURATION.md](CONFIGURATION.md)).
  When read from Vault they are re-read on every use, so the vault manager
  can rotate them without a restart.

## Where it fits

```
trustee authorizes distribution ─▶ tm-api ─▶ proceed_payouts (REGISTERED)
admin actions in tm-web ─▶ tm-api ─▶ app-backend DB + Redis wake-up ─▶ payout-engine
payout-engine ─▶ Base (payout Safe → MultiSendCallOnly → holders, fee, VAT)
app-backend GET /v1/tokenization/payouts ◀─ app-mobile (dividend screens)
```

See [INTEGRATION.md](INTEGRATION.md) for the exact tables and channels.

## Tech stack

- Go (module `trovo-payout-engine`), GORM on app-backend's Postgres (SQLite
  for tests), go-ethereum, go-redis v9, HashiCorp Vault API, Firebase
  Messaging.
- Safe v1.4.1 transactions (EIP-712 signatures from the first three
  configured signers; the first executes and pays the gas).

## Directory structure

```
payout-engine/
├── main.go                 # config, DB, RPC, signers, push, Redis; runs the engine
└── internal/
    ├── config/             # environment, signer sources (env / Vault)
    ├── store/              # app-backend tables the engine uses (mirrors)
    ├── chain/              # RPC helpers: Transfer logs, balances, block by time
    ├── safe/               # Safe transaction hashing, signing, execution
    ├── engine/             # prepare, fund check, pay, recover, notify, sweep
    ├── bus/                # Redis: wake-ups, events, heartbeat
    └── push/               # Firebase push notifications
```

## Running it locally

```bash
cd payout-engine
cp .env.example .env        # fill it in; see CONFIGURATION.md
go test ./...               # unit tests; the on-chain test needs the local stack (below)
go run .
```

The end-to-end test (`internal/engine/engine_e2e_test.go`) runs the whole
lifecycle on a local Hardhat chain. It covers fee and VAT, a failing holder
isolated by bisection, a dropped batch, notifications and the sweep. Run it
against the local stack from `paymaster/contracts` and `market/contracts`:

```bash
AA_LOCAL_STACK=../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_MARKET=../market/contracts/deployments/local-market.json go test ./...
```

Without those variables the test is skipped.

## Further reading

- [DEPLOYMENT.md](DEPLOYMENT.md): building, Docker, running in production.
- [CONFIGURATION.md](CONFIGURATION.md): every environment variable.
- [INTEGRATION.md](INTEGRATION.md): tables, Redis channels, tm-api and
  app-backend.

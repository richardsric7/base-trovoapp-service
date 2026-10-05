# Public Markets

Public Markets lets Trovo App users and partner exchanges buy and sell
tokenized NGX equities and FMDQ bonds. Each token is backed one-for-one by
real units that a Custodian holds in an omnibus (pool) account at CSCS. The
code lives in `internal/components/publicmarkets/`:

| Package | What it holds |
|---|---|
| `models` | The tables (all prefixed `public_market_`, plus `approved_dealing_members`, `custodian_positions`, `beneficial_ownership_*`, `price_oracle_snapshots`, `wallet_provisioning_requests`) and their states |
| `partners` | The Custodian, Dealing Member and price-feed interfaces, their REST clients and the HMAC signing |
| `services` | The engine: quotes, orders, net batches, instructions, partner events, the ownership ledger, reconciliation, prices, exchange accounts and webhooks, dividends, and the background jobs |
| `controllers` | The routes: the apps' `/v1/public-markets`, the exchanges' `/v1/trovo-api/public-markets` and the partner callbacks |

Trovo Manager (tm-api + tm-web) operates it, and app-mobile and app-web
show it to users. The engine runs inside app-backend. Each of its loops
runs on one instance at a time (a database lock), like the other
background jobs.

## How the architecture documents map to Base

The architecture and integration documents were written for Stellar. On
Base:

- **Asset token**: one `TokenizedAsset` ERC-20 per asset, owned by the
  asset's **issuing Safe**. Minting and burning are Safe transactions
  signed by the Public Markets signers (`PUBLIC_MARKETS_SIGNERS`).
- **Treasury**: one Safe (`PUBLIC_MARKETS_TREASURY_SAFE`) receives buyers'
  funding stablecoin (CNGN) and pays sellers, dividends, fees and
  withholding tax.
- **Exchange customer wallets**: each one is a Safe owned by the Public
  Markets signers (threshold `min(3, signers)`). It is predicted when the
  wallet is provisioned and deployed the first time it has to send.
  Exchanges never hold keys.
- **Trovo App wallets** are the users' own Safes. A buy is a CNGN transfer
  to the treasury and a sell is a token transfer to the issuing Safe. Both
  are ordinary signed wallet operations (two calls, like every other send).
- **Beneficial ownership ledger**: the engine indexes every token
  Transfer into `beneficial_ownership_movements` (one row per wallet per
  transfer) and keeps running balances in `beneficial_ownership_ledgers`.
  Dividend snapshots are taken from these movements at the record-date
  block.

## Two paths for an order

**Fast.** The market is open, the price is fresh and the Custodian already
holds enough unissued units ("inventory" = latest position − supply −
units reserved by open orders). The order mints from inventory, or for a
sale burns and pays from the treasury, within seconds.

**Slow.** Otherwise the order waits for the session's **net batch**. Every
`BatchIntervalMinutes`, the engine nets each asset's waiting orders:

- **BUY** when creations exceed inventory plus redemptions. It buys the
  shortfall in whole units, plus the asset's inventory target.
- **SELL** when app sellers must be paid more than the treasury's free
  cash. Free cash is the balance minus what it owes: exchange balances,
  unswept fees, pending payouts and approved dividends.
- **INTERNAL** when everything nets. Nothing goes to the market.

A BUY batch whose value, added to the asset's other BUY batches that day,
exceeds `NetCreationThreshold` waits for `ApprovalsRequired` approvals.
Approvers must be listed in `NetCreationApprovers`.

A released batch sends an order to the Dealing Member and an instruction to
the Custodian. When the trade executes, orders move to `executed`. When the
Custodian confirms `settlement_final`, the position grows (or shrinks) and
tokens are minted, burned or paid. **No token is minted before the
Custodian confirms settlement.**

Order states: `awaiting-payment` → `queued` → (`filled-from-inventory` |
`pending-execution` → `executed`) → `settlement_final` → `submitted` →
`chain_final` → `complete`. The other end states are `rejected` (refunded),
`failed` and `cancelled`. Exchanges see only `pending`, `settled`,
`rejected` and `failed`.

## Fail-closed reconciliation

Every day at `ReconciliationHour` (WAT), and whenever Trovo Manager asks,
each asset is checked:

- supply must be ≤ the Custodian's units;
- the ledger must equal supply (once the indexer has caught up);
- a position must exist.

Any drift **halts** the asset, and new orders are refused. An asset
resumes only when a fresh run matches. The difference between position and
supply is working inventory.

## Partners: mock, REST or manual

The real Custodian and Dealing Member APIs (and CSCS behind the Custodian)
are not agreed yet, so each partner has a **mode** (set in Trovo Manager):

- **MOCK** (default): in-process mocks. The broker fills at the reference
  price ±5 bps after `PUBLIC_MARKETS_MOCK_EXECUTION_SECONDS`. The Custodian
  settles after `PUBLIC_MARKETS_MOCK_SETTLEMENT_SECONDS`, keeps its own
  book and sends a position feed at 18:00 WAT. Their callbacks go through
  exactly the same processing as real webhooks.
- **REST**: the clients in `partners/partners.go` call the endpoints the
  integration document proposes:
  - `POST /v1/instructions/creation|redemption`
  - `GET /v1/positions/{code}`
  - `POST /v1/orders`

  Calls carry an Idempotency-Key and an HMAC (or mTLS) signature. 5xx and
  429 responses are retried; other 4xx responses are permanent. When the
  final specifications arrive, only these clients and the callback payload
  parsing in `services/callbacks.go` should need to change.
- **MANUAL**: nothing is sent. Operations act on the partner's portal and
  record the confirmation in Trovo Manager, which stores it as a partner
  event for the engine to process.

Instructions are queued durably and retried with backoff (30 s · 2ⁿ, up to
`InstructionMaxAttempts`) before being escalated. A settlement not
confirmed within `SettlementSLAHours` raises an alert.

Every partner webhook is stored once, unique on (partner, kind,
reference), so redeliveries are acknowledged and ignored.

## Exchanges

An exchange is a **service link** (its API key) that Trovo Manager has
onboarded for Public Markets. Onboarding gives it:

- a signing secret;
- a registered funding wallet;
- a tier (which sets the service link's rate limit);
- a webhook URL.

Exchange orders are paid from a **prefunded balance**. The integration
document does not define how exchanges pay, so this is a stated
assumption. The exchange sends CNGN from its funding wallet to the
treasury and submits the transaction to `POST /deposits`. The engine
verifies the transfer on-chain and credits the balance. Creations debit
the balance, redemptions credit it, and withdrawals are paid back to the
funding wallet on a Trovo Manager request.

## Dividends and coupons

Dividends do **not** use the proceeds payout engine. Each corporate action
goes through these steps:

1. **Declared.** The Custodian's corporate-action webhook declares it, or
   Operations declare it in Trovo Manager. BONUS, RIGHTS and SPLIT are
   marked `NEEDS_MANUAL`.
2. **Snapshotted.** After the record date ends (Lagos time), the engine
   takes a snapshot of holders at that block. Each owner gets an
   entitlement with withholding tax:
   - exchange customers: the resident, non-resident or missing-tax-ID
     rate, from their residency and tax ID;
   - app users: from their country.

   Platform wallets are `RETAINED`. The snapshot has a SHA-256 checksum.
3. **Approved.** It needs approvals on **that checksum** from
   `DividendApprovers`. An approval on an older snapshot does not count.
4. **Paid.** On or after the pay date, once the treasury can cover it:
   - exchange balances are credited, with a `dividend.paid` webhook that
     the exchange must confirm;
   - app holders are paid from the treasury in batches of 100 transfers;
   - withholding tax goes to `WHTWallet`;
   - holders get a push notification.

## Tests

```bash
go test ./internal/components/publicmarkets/... -count=1
# the same tests against Postgres (each test gets its own schema):
PUBLIC_MARKETS_TEST_POSTGRES="host=localhost user=... dbname=... sslmode=disable" \
  go test ./internal/components/publicmarkets/... -count=1 -p 1
```

The engine tests run the real engine, mocks and database against an
in-memory chain that applies the Safes' transfer, mint and burn calls.
They cover:

- fast and slow creation through the mock partners;
- batch approvals and a rejected batch's refund;
- netted redemption;
- reconciliation halt and resume;
- idempotent partner events;
- provisioning field errors and signatures;
- a dividend from snapshot to payment.

The controller tests cover exchange authentication, Idempotency-Key
replay and the partner callbacks.

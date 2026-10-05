# Deployment

## Prerequisites

- [Go](https://go.dev/doc/install), the version in `go.mod` (1.26)
- app-backend deployed against the same database. Its AutoMigrate creates
  the payout tables: `proceed_payouts`, `tokenized_asset_payout_schedules`,
  `proceed_payout_batches`, `proceed_payout_approvals`,
  `payout_token_indexes`, `payout_token_balances`, `payout_engine_states`.
- tm-api deployed (admins drive payouts from TM), with the same Redis if
  you use one.
- A Base RPC endpoint.
- The payout Safe and its signers (see "Before the first payout" below).

## Build and test

```bash
cd payout-engine
go build ./...
go vet ./...
go test ./...      # unit tests; the on-chain e2e test is skipped without the local stack
```

To run the on-chain e2e test, start the local stack first (a Hardhat node
with Safe, the EntryPoint, the paymaster and the market contracts):

```bash
cd ../paymaster/contracts && npx hardhat node &                      # leave running
npx hardhat run scripts/local-stack.js --network localhost
cd ../../market/contracts && npx hardhat run scripts/local-market.js --network localhost
cd ../../payout-engine
AA_LOCAL_STACK=../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_MARKET=../market/contracts/deployments/local-market.json go test ./...
```

## Run

```bash
cp .env.example .env    # fill it in (CONFIGURATION.md)
go run .
```

At startup it logs the chain, the payout Safe and where the signers come
from. An invalid or missing required setting stops it with the reason.

## Docker

```bash
docker build -t payout-engine .
docker run --env-file .env payout-engine
```

The image has no exposed ports, because the engine has no HTTP interface.
The Dockerfile follows the other Go services here. It was **not** built in
the environment where it was written (no Docker daemon there), so build it
once in CI before relying on it.

## Running in production

- **Run one replica, or more as standbys.** Only one instance works at a
  time: it claims `payout_engine_states` (row 1) with a heartbeat, and a
  standby takes over when the heartbeat is more than 90 seconds old. Extra
  replicas only add failover. Give each an `INSTANCE_ID`, or the hostname
  is used.
- **Stopping is safe at any time** (SIGTERM). A batch already sent is
  settled, dropped or sent again by whichever instance runs next.
- **Keep the executing signer funded.** The first signer pays every batch's
  gas. The funding check refuses to start a payout without enough ETH for
  it and shows the amount in the payout's note in TM.
- **Monitor it.** TM shows the engine's heartbeat, current activity and
  last error (**Proceeds Payouts → Fee & engine**). Events are also
  published on Redis `payout-engine:events`.

## Before the first payout

1. **Deploy the payout Safe** (v1.4.1) on Base. Its owners are the
   `PROCEED_PAYOUT_SIGNERS` addresses, with a threshold of at most 3. Set
   `PROCEED_PAYOUT_SAFE_ADDRESS`.
2. **Register the signer secret** in TM's vault manager as a managed secret
   with the payout Safe as its wallet. Rotating a signer then swaps the
   Safe owner too. Either inject the secret as `PROCEED_PAYOUT_SIGNERS`, or
   point `PROCEED_PAYOUT_SIGNERS_VAULT_PATH` at it.
3. **Set the payout fee wallet** in TM (Proceeds Payouts → Fee & engine),
   and `VAT_WALLET` if payout fees are charged. The fee and VAT wallets
   must not hold the asset's token.
4. Make sure each payout currency (e.g. CNGN) is in app-backend's
   `tokenization_currencies` with its contract address. A payout is
   registered against the first of the distribution's currency, the
   asset's proceed payout currency and its country's internal balance token
   that has a contract.

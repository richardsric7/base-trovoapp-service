# market - deployment

## Build and test

```bash
cd market/contracts
npm ci
npx hardhat compile
npx hardhat test          # TrovoOfferBook and TokenizedAsset
```

The Solidity compiler (0.8.28) comes from the pinned `solc` npm package, so
builds do not download a compiler.

## Deploy

1. Write the deploy secret (see [CONFIGURATION.md](CONFIGURATION.md)).
2. Dry run, then deploy:

   ```bash
   npm run build
   VAULT_ADDR=… VAULT_TOKEN=… node scripts/deploy.js --dry-run
   VAULT_ADDR=… VAULT_TOKEN=… node scripts/deploy.js
   ```

   The script deploys `TrovoOfferBook` with the authorizers, lists
   `TRADABLE_TOKENS`, hands ownership to `OWNER_ADDRESS` and writes
   `deployments/<chainId>.json`.
3. Set `OFFER_BOOK_ADDRESS` and `OFFER_AUTHORIZER_PRIVATE_KEY` in
   app-backend.

`TokenizedAsset` is deployed per asset by operations, with the asset's
issuing wallet as owner, and registered in app-backend (see
[INTEGRATION.md](INTEGRATION.md)).

## Local deployment for end-to-end tests

On the local node of the paymaster's local stack:

```bash
cd paymaster/contracts
npx hardhat node                                           # terminal 1
npx hardhat run scripts/local-stack.js --network localhost  # terminal 2 (keeps the dev bundler running)
cd ../../market/contracts
npx hardhat run scripts/local-market.js --network localhost
```

`local-market.js` deploys the book (owner hardhat #0, authorizer hardhat #2),
a mock cNGN and a mock internal balance token, lists them, and writes
`deployments/local-market.json`. Then, from `app-backend/`:

```bash
AA_LOCAL_STACK=$PWD/../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_MARKET=$PWD/../market/contracts/deployments/local-market.json \
go test ./internal/components/users/services/ -run TestTokenizationOnLocalChain -v
```

The test runs app-backend's own mint plan, pricing and authorization code:
an issuing Safe owned by a profile key and four approvers mints (signed by
two approvers: deploying itself and the distribution Safe, minting, opening
the offer), a new wallet buys with cNGN in its activation operation, a
replayed authorization is refused, and a minting Safe makes a fiat delivery
paid in the internal balance token. The dev bundler's fixed gas estimates
are too low for a mint, so the test raises the gas buffer; a real bundler
estimates by simulation.

## Verification done so far

- 18 Hardhat tests: escrow, pricing and rounding, fills paid in a
  stablecoin and in internal balance for another recipient, every
  authorization check (wrong key, someone else's authorization, any changed
  field, reused nonce, expired), limits, delisting and pause, repricing,
  cancel (also while paused or delisted), fee-on-transfer tokens,
  owner-only administration.
- The Go EIP-712 digest and signatures (app-backend `internal/aa`) are
  accepted by the contract in the end-to-end test.
- Not yet: Base Sepolia, and an external audit.

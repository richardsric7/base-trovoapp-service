# recovery - deployment

## Build and test

```bash
cd recovery/contracts
npm ci
npx hardhat compile
npx hardhat test          # the module against real Safe v1.4.1 wallets
```

`.npmrc` sets `legacy-peer-deps` because `@safe-global/safe-contracts`
1.4.1 (used only by the tests, for real Safe wallets) declares an old
`ethers` peer; the toolbox's own peers are pinned explicitly in
`package.json`.

## Deploy

1. Write the deploy secret (see [CONFIGURATION.md](CONFIGURATION.md)).
2. Dry run, then deploy:

   ```bash
   npm run build
   VAULT_ADDR=… VAULT_TOKEN=… node scripts/deploy.js --dry-run
   VAULT_ADDR=… VAULT_TOKEN=… node scripts/deploy.js
   ```

   The script deploys `SocialRecoveryModule` with `RECOVERY_PERIOD_SECONDS`,
   checks it reports version `0.2.0` and writes
   `deployments/<chainId>.json` (`recoveryModule`, `recoveryPeriodSeconds`).
   The module has no owner and no other settings.
3. Create the guardian Safe (e.g. 2-of-3 recovery keys, one held offline)
   and fund its signer keys with a little ETH for gas.
4. Set in app-backend: `RECOVERY_MODULE_ADDRESS`, `RECOVERY_PERIOD_SECONDS`,
   `ACCOUNT_RECOVERY_GUARDIAN_SAFE`, `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS`
   (see [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md#account-recovery)),
   and the `ACCOUNT_RECOVERY_FEE` service fee.

Changing the recovery period means deploying a new module; wallets keep
the module they enabled until their users turn recovery off and on again.

## Local deployment for end-to-end tests

On the local node of the paymaster's local stack:

```bash
cd paymaster/contracts
npx hardhat node                                           # terminal 1
npx hardhat run scripts/local-stack.js --network localhost  # terminal 2 (keeps the dev bundler running)
cd ../../recovery/contracts
npx hardhat run scripts/local-recovery.js --network localhost   # 60 second period
```

Then, from `app-backend/`:

```bash
AA_LOCAL_STACK=$PWD/../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_RECOVERY=$PWD/../recovery/contracts/deployments/local-recovery.json \
go test ./internal/components/users/services/ -run TestAccountRecoveryOnLocalChain -v
```

The test runs app-backend's own code against the module, Safe 4337
wallets and the dev bundler: a user turns recovery on for a primary wallet
and a sub-wallet (each in its first operation), shares the sub-wallet with
a co-signer (which removes the guardian from it), then the guardian Safe
starts a recovery that the user cancels, starts another one that is
finalized after the period, the account moves to the new key (which then
pays from the wallet while the old key cannot), the co-signer approves the
shared wallet's REPLACE SIGNER request, and a recovery started directly
with the guardian's key raises an alert and is cancelled by the owner.

## Verification done so far

- 6 Hardhat tests: replacing the lost key after the period with the funds
  intact (new key works, old key does not), cancel during the period and
  nonce invalidation, the guardian cannot move funds, become an owner or
  change guardians, a two-owner wallet keeps its co-signer, revoking the
  guardian ends coverage, and refused configurations.
- The app-backend end-to-end test above.
- Not yet: Base Sepolia. The module's audits are Candide's (Ackee,
  Nethermind, Certora); Trovo's integration has not been audited.

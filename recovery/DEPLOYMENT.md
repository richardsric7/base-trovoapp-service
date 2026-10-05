# recovery - deployment

A step-by-step guide to deploying the account recovery module, written for
someone doing it for the first time.

## 1. What you are deploying

One smart contract on the Base blockchain: **Candide's Social Recovery
Module v0.2.0**. Users who turn on account recovery in the Trovo app add
it to their wallets. If they lose their secret key, Trovo's **recovery
guardian** can, after a waiting period the user can cancel, replace the
lost key with a new one. The guardian can never move funds.

You deploy the module **once per network**, then create the **guardian**
(a Safe wallet controlled by dedicated recovery keys) and give both to
app-backend. Deploy to the test network (Base Sepolia) first.

Nothing else needs to be running first.

## 2. Before you start

| Tool / account | Install or get it | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Node.js 20 or 22 (LTS) | <https://nodejs.org> | `node --version` |
| Vault CLI | <https://developer.hashicorp.com/vault/install> | `vault --version` |
| A Base node URL | a provider such as Alchemy or QuickNode | see [`RPC_URL`](CONFIGURATION.md#rpc_url) |
| A deployer wallet with a little ETH | see [`DEPLOYER_PRIVATE_KEY`](CONFIGURATION.md#deployer_private_key) | its balance on [basescan.org](https://basescan.org) |
| Safe{Wallet} | <https://app.safe.global> | for the guardian (step 5) |

## 3. Get the code, build and test

1. Get the code and install the dependencies:
   ```bash
   git clone https://github.com/richardsric7/base-trovoapp-service.git
   cd base-trovoapp-service/recovery/contracts
   npm ci
   ```
   (`.npmrc` turns on `legacy-peer-deps` because the Safe contracts used by
   the tests declare an old peer dependency; that is expected.)
2. Compile: `npm run build`. **You should see:** `Compiled ... Solidity
   files successfully`.
3. Test: `npm test`. **You should see:** 6 passing, `0 failing`. The tests
   run the module against real Safe v1.4.1 wallets.

## 4. Deploy the module (once per network)

1. **Choose the recovery period** (see
   [`RECOVERY_PERIOD_SECONDS`](CONFIGURATION.md#recovery_period_seconds)):
   7 days (`604800`) is a sensible default. It cannot be changed later
   without deploying a new module.
2. **Store the settings in Vault:**
   ```bash
   export VAULT_ADDR=https://vault.internal.trovo.io:8200
   vault login
   vault kv put secret/trovo/recovery/deploy \
     RPC_URL=https://sepolia.base.org CHAIN_ID=84532 \
     DEPLOYER_PRIVATE_KEY=0x... RECOVERY_PERIOD_SECONDS=604800
   ```
3. **Dry run** (checks everything, deploys nothing):
   ```bash
   export VAULT_TOKEN=hvs....
   npm run deploy:dry-run
   ```
   **You should see:**
   ```
   config: vault secret/trovo/recovery/deploy
   chain 84532, deployer 0x..., recovery period 604800s (7 days)
   dry run: nothing deployed
   ```
4. **Deploy:** `npm run deploy`. **You should see:**
   `deployed SocialRecoveryModule v0.2.0 at 0x...` and
   `next: set RECOVERY_MODULE_ADDRESS=0x... in app-backend`.
5. **Keep the record:** commit `contracts/deployments/<chainId>.json`
   (it holds `recoveryModule` and `recoveryPeriodSeconds`).

The module has no owner and no other settings.

## 5. Create the guardian

The guardian should be a Safe owned by dedicated recovery keys, so no
single leaked key can start recoveries.

1. Create three keys (repeat the command three times; store each key in
   your secrets manager):
   ```bash
   node -e "const {Wallet}=require('ethers');const w=Wallet.createRandom();console.log('address',w.address);console.log('key',w.privateKey)"
   ```
   Keep one of them **offline** (for example printed and locked away).
2. At <https://app.safe.global> choose **Create account** on the same
   network, add the three addresses as owners and set the threshold to 2.
   Copy the Safe's address.
3. Send each of the two online keys a little ETH (about 0.005 ETH): they
   pay the network fees of starting and finalizing recoveries; the first
   one pays most.

## 6. Set the parameters in app-backend

| Parameter | Value | Example |
|---|---|---|
| `RECOVERY_MODULE_ADDRESS` | the module from step 4 | `0xB7f8BC63BbcaD18155201308C8f3540b07f84F5e` |
| `RECOVERY_PERIOD_SECONDS` | the same period | `604800` |
| `ACCOUNT_RECOVERY_GUARDIAN_SAFE` | the Safe from step 5 | `0xA11CE00000000000000000000000000000000002` |
| `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` | the two online keys, `;` separated | `0x59c6...;0x5de4...` |
| `ACCOUNT_RECOVERY_FEE_WALLET` | where the opt-in fee is paid | a treasury wallet |

Each is described in
[app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md). Set the
fee itself in the `ACCOUNT_RECOVERY_FEE` row of app-backend's
`service_fees` table (an inactive row or a fee of 0 charges nothing), then
restart app-backend.

## 7. Check it works

1. On [basescan.org](https://basescan.org) (or sepolia.basescan.org), the
   module's address has contract code and the deployment transaction from
   step 4 succeeded. The period is not readable on-chain (the contract
   keeps it private); `recoveryPeriodSeconds` in the deployment file is
   the record. If you verify the contract's source on basescan, its
   **Read Contract** tab shows `VERSION` = `0.2.0`.
2. On the test network: turn recovery on for a test account in the app,
   recover the account on another device, cancel the recovery from the
   first device, then recover again and wait for the period: the account
   moves to the new key.
3. Without the settings above, app-backend's recovery endpoints answer
   `error-account-recovery-not-configured`; with them, they work.

## 8. Updating and rolling back

The module cannot be changed. To change the recovery period, deploy a new
module and point `RECOVERY_MODULE_ADDRESS` at it. Wallets keep the module
they enabled until their users turn recovery off and on again. The same
applies to changing the guardian.

## 9. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `VAULT_ADDR is not set` / `VAULT_TOKEN is not set` | bootstrap variables missing | `export` them in the same terminal |
| `RECOVERY_PERIOD_SECONDS is required` (or another key) | missing from the secret | `vault kv get secret/trovo/recovery/deploy` and add it |
| `RECOVERY_PERIOD_SECONDS is under a day` | a test period on a real network | use at least `86400`, or `ALLOW_SHORT_PERIOD=true` on a test network |
| `RPC_URL is chain ..., CHAIN_ID says ...` | node and `CHAIN_ID` disagree | fix the wrong one |
| `... not found - run "npm run build" first` | not compiled | `npm run build` |
| `insufficient funds` | the deployer has no ETH | fund it and run again |
| app-backend: `error-account-recovery-not-configured` | module or guardian not set | step 6 |

## Local blockchain for end-to-end tests

app-backend's end-to-end test uses the module on a local blockchain started
from the paymaster project:

```bash
# terminal 1: a local blockchain (leave it running)
cd paymaster/contracts && npx hardhat node

# terminal 2: Safe, EntryPoint, paymaster and a dev bundler, then the module
cd paymaster/contracts && npx hardhat run scripts/local-stack.js --network localhost
cd ../../recovery/contracts && npx hardhat run scripts/local-recovery.js --network localhost
```

`local-recovery.js` deploys the module with a 60-second period (see
[`RECOVERY_PERIOD_SECONDS` (local script)](CONFIGURATION.md#recovery_period_seconds-local-script))
and writes `deployments/local-recovery.json`. Then, from `app-backend/`:

```bash
AA_LOCAL_STACK=$PWD/../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_RECOVERY=$PWD/../recovery/contracts/deployments/local-recovery.json \
go test ./internal/components/users/services/ -run TestAccountRecoveryOnLocalChain -v
```

- **`AA_LOCAL_STACK`:** the file listing the local Safe, EntryPoint and
  paymaster (written by `local-stack.js`).
- **`AA_LOCAL_RECOVERY`:** the file listing the local module (written by
  `local-recovery.js`).

The test turns recovery on for a primary wallet and a sub-wallet, shares
the sub-wallet with a co-signer (which removes the guardian from it), has
the guardian start a recovery the user cancels, then another that
completes after the period, checks the new key works and the old one does
not, and checks that a recovery started outside Trovo raises an alert.

## Verification done so far

- 6 Hardhat tests: replacing a lost key after the period with funds intact,
  cancelling during the period, the guardian cannot move funds, become an
  owner or change guardians, a two-owner wallet keeps its co-signer,
  removing the guardian ends coverage, and refused configurations.
- The app-backend end-to-end test above.
- Not yet: a Base Sepolia deployment. The module's audits are Candide's
  (Ackee, Nethermind, Certora); Trovo's integration has not been audited.

# recovery - configuration

Every parameter the recovery project reads. Each one says what it is, why
it is needed, whether you must set it, an example and how to get a real
value. Examples are never real keys.

The recovery module has a single setting of its own, its **recovery
period**, fixed when it is deployed. What you configure is the **deploy
script** (`contracts/scripts/deploy.js`).

## How to set them

The deploy script reads its settings from **HashiCorp Vault** (a secure
store for secrets), so the deployer's private key never sits in files or
shell history:

1. A few **bootstrap environment variables** tell the script where Vault is
   and how to log in ([below](#bootstrap-environment-variables)).
2. Everything else is one **Vault secret** with several keys
   ([below](#keys-of-the-deploy-secret)).

For local experiments only, `RECOVERY_CONFIG_SOURCE=env` reads those keys
from ordinary environment variables instead.

### The minimum to deploy

| Where | Parameter |
|---|---|
| environment | [`VAULT_ADDR`](#vault_addr), [`VAULT_TOKEN`](#vault_token) |
| Vault secret | [`RPC_URL`](#rpc_url), [`CHAIN_ID`](#chain_id), [`DEPLOYER_PRIVATE_KEY`](#deployer_private_key), [`RECOVERY_PERIOD_SECONDS`](#recovery_period_seconds) |

---

## Bootstrap environment variables

### `VAULT_ADDR`

- **What it is:** The web address of your Vault server.
- **Why it's needed:** The script reads its settings from Vault.
- **Required:** Yes, unless `RECOVERY_CONFIG_SOURCE=env`.
- **Example:** `https://vault.internal.trovo.io:8200`
- **How to get it:** The same Vault the other services use (tm-api's
  `VAULT_ADDR`); ask whoever runs Vault.

### `VAULT_TOKEN`

- **What it is:** A Vault access token.
- **Why it's needed:** Vault only returns the secret to a token allowed to
  read it.
- **Required:** Yes, unless `RECOVERY_CONFIG_SOURCE=env`.
- **Example:** `hvs.CAESIJexampleexampleexample`
- **How to get it:** Ask your Vault administrator for a short-lived token
  with the read-only policy in [Writing the secret](#writing-the-secret),
  for example `vault token create -policy=recovery-deploy -ttl=1h`.

### `VAULT_KV_MOUNT`

- **What it is:** The name of Vault's key-value (KV version 2) store.
- **Why it's needed:** The secret's full location is mount + path.
- **Required:** No, default `secret`.
- **Example:** `secret`
- **How to get it:** Keep the default unless your Vault uses another mount
  (`vault secrets list` shows them).

### `VAULT_NAMESPACE`

- **What it is:** The Vault namespace the secret lives in.
- **Why it's needed:** Only Vault Enterprise and HCP Vault have
  namespaces; without the right one the secret is not found.
- **Required:** No.
- **Example:** `trovo`
- **How to get it:** Your Vault administrator; leave unset on open-source
  Vault.

### `RECOVERY_DEPLOY_SECRET_PATH`

- **What it is:** Where in Vault the deploy secret is stored.
- **Why it's needed:** It tells the script which secret to read.
- **Required:** No, default `trovo/recovery/deploy`.
- **Example:** `trovo/recovery/deploy-sepolia`
- **How to get it:** Keep the default; use a different path per network if
  you deploy to both a test network and mainnet.

### `RECOVERY_CONFIG_SOURCE`

- **What it is:** Where the script reads the keys below from: `vault` or
  `env`.
- **Why it's needed:** `env` lets you try a deployment on a local or test
  network without Vault.
- **Required:** No, default `vault`.
- **Example:** `env`
- **How to get it:** `env` only on your own machine for a local or test
  network. Never for mainnet.

## Keys of the deploy secret

### `RPC_URL`

- **What it is:** The web address of a Base blockchain node.
- **Why it's needed:** The script sends the deployment transaction through
  it.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
  (mainnet), `https://sepolia.base.org` (test network)
- **How to get it:** Sign up with a node provider (Alchemy, Infura,
  QuickNode), create an app for Base Mainnet or Base Sepolia and copy its
  HTTPS URL.

### `CHAIN_ID`

- **What it is:** The number of the network you mean to deploy to.
- **Why it's needed:** The script checks the node really is that network.
- **Required:** Yes.
- **Example:** `8453` (Base Mainnet), `84532` (Base Sepolia)
- **How to get it:** `8453` for production, `84532` for testing.

### `DEPLOYER_PRIVATE_KEY`

- **What it is:** The private key of the wallet that sends the deployment
  transaction.
- **Why it's needed:** Someone has to pay the network fee. This wallet has
  no role in the module afterwards (the module has no owner).
- **Required:** Yes.
- **Example:** `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
  (a well-known test key; never use it for real funds)
- **How to get it:** Create a fresh wallet, from `recovery/contracts` after
  `npm ci`:
  ```bash
  node -e "const {Wallet}=require('ethers');const w=Wallet.createRandom();console.log('address',w.address);console.log('key',w.privateKey)"
  ```
  Send it about 0.005 ETH on Base (a faucet for Base Sepolia is linked
  from <https://docs.base.org/tools/network-faucets>).

### `RECOVERY_PERIOD_SECONDS`

- **What it is:** How long, in seconds, a started recovery must wait before
  it takes effect.
- **Why it's needed:** This is the user's window to cancel a recovery they
  did not ask for. Longer is safer; shorter means a user who really lost
  their key waits less. It cannot be changed after deployment.
- **Required:** Yes. At least one day (`86400`) unless
  `ALLOW_SHORT_PERIOD=true`.
- **Example:** `604800` (7 days)
- **How to get it:** A business decision; 7 days is a sensible default.
  Days × 86400 = seconds.

### `ALLOW_SHORT_PERIOD`

- **What it is:** Allows a recovery period shorter than one day.
- **Why it's needed:** Testing on a test network is easier with a period of
  minutes. On mainnet a short period would leave users almost no time to
  cancel a recovery they did not start.
- **Required:** No, default off.
- **Example:** `true`
- **How to get it:** Only on a test network.

### Writing the secret

```bash
vault kv put secret/trovo/recovery/deploy \
  RPC_URL=https://sepolia.base.org CHAIN_ID=84532 \
  DEPLOYER_PRIVATE_KEY=0x... RECOVERY_PERIOD_SECONDS=604800
```

A read-only policy for the person or machine that deploys (save as
`recovery-deploy.hcl`, load with `vault policy write recovery-deploy recovery-deploy.hcl`):

```hcl
path "secret/data/trovo/recovery/deploy" {
  capabilities = ["read"]
}
```

## Local development only

Read only by the local test setup
([DEPLOYMENT.md](DEPLOYMENT.md#local-blockchain-for-end-to-end-tests)).

### `LOCAL_RPC_URL`

- **What it is:** The address of a local test blockchain.
- **Why it's needed:** Hardhat's `localhost` network uses it.
- **Required:** No, default `http://127.0.0.1:8545` (where
  `npx hardhat node` listens).
- **Example:** `http://127.0.0.1:8545`
- **How to get it:** Keep the default.

### `RECOVERY_PERIOD_SECONDS` (local script)

- **What it is:** The recovery period the local deployment script
  (`scripts/local-recovery.js`) uses.
- **Why it's needed:** The end-to-end test waits for a recovery to finish,
  so a short period keeps it fast.
- **Required:** No, default `60` seconds.
- **Example:** `60`
- **How to get it:** Keep the default.

## What app-backend needs afterwards

| Parameter | Value |
|---|---|
| `RECOVERY_MODULE_ADDRESS` | the deployed module (printed by the script; `recoveryModule` in `contracts/deployments/<chainId>.json`) |
| `RECOVERY_PERIOD_SECONDS` | the same period (`recoveryPeriodSeconds` in that file), used in messages to users |
| `ACCOUNT_RECOVERY_GUARDIAN_SAFE` | the guardian Safe ([DEPLOYMENT.md step 5](DEPLOYMENT.md#5-create-the-guardian)) |
| `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` | the guardian Safe's signer keys, `;` separated |
| `ACCOUNT_RECOVERY_FEE_WALLET` | where the fee for turning recovery on is paid |

Each is described in [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md).
The fee itself is the `ACCOUNT_RECOVERY_FEE` row of app-backend's
`service_fees` table.

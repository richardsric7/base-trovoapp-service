# recovery - configuration

Only the deploy script is configured; the module's one setting, its
recovery period, is fixed at deployment.

## Deploy script (`contracts/scripts/deploy.js`)

Settings come from a Vault KV v2 secret, like the other deploy scripts.
Bootstrap environment:

| Variable | Required | Example | What it does |
|---|---|---|---|
| `VAULT_ADDR` | yes (unless `RECOVERY_CONFIG_SOURCE=env`) | `https://vault.internal.trovo.io:8200` | Vault server URL. |
| `VAULT_TOKEN` | yes (unless `RECOVERY_CONFIG_SOURCE=env`) | `hvs.CAESIJ…` | Token allowed to read the secret. |
| `VAULT_KV_MOUNT` | no (default `secret`) | `secret` | KV v2 mount. |
| `VAULT_NAMESPACE` | no | `trovo` | Vault Enterprise namespace. |
| `RECOVERY_DEPLOY_SECRET_PATH` | no (default `trovo/recovery/deploy`) | `trovo/recovery/deploy` | Path of the secret below. |
| `RECOVERY_CONFIG_SOURCE` | no (default `vault`) | `env` | `env` reads the keys below from the process environment. **Local development only.** |

Keys of the secret:

| Key | Required | Example | What it does |
|---|---|---|---|
| `RPC_URL` | yes | `https://mainnet.base.org` | Base RPC endpoint. |
| `CHAIN_ID` | yes | `8453` | Checked against the RPC. |
| `DEPLOYER_PRIVATE_KEY` | yes | `0x…` | Pays for the deployment. Has no role in the module afterwards. |
| `RECOVERY_PERIOD_SECONDS` | yes | `604800` | How long a started recovery waits before it can be finalized - the user's window to cancel. At least a day unless `ALLOW_SHORT_PERIOD=true`. |
| `ALLOW_SHORT_PERIOD` | no | `true` | Allow a period under a day - testnets only. |

Writing the secret:

```bash
vault kv put secret/trovo/recovery/deploy \
  RPC_URL=https://sepolia.base.org CHAIN_ID=84532 \
  DEPLOYER_PRIVATE_KEY=0x… RECOVERY_PERIOD_SECONDS=604800
```

A read-only policy for the deploy operator:

```hcl
path "secret/data/trovo/recovery/deploy" {
  capabilities = ["read"]
}
```

## What app-backend needs afterwards

- `RECOVERY_MODULE_ADDRESS` - the deployed module.
- `RECOVERY_PERIOD_SECONDS` - the same period, for messages to users.
- `ACCOUNT_RECOVERY_GUARDIAN_SAFE` and `ACCOUNT_RECOVERY_GUARDIAN_SIGNERS` -
  the guardian Safe and its signer keys.

See [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md#account-recovery).

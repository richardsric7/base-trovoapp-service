# market - configuration

Only the deploy script is configured; the contracts take their settings at
deployment and from their owner afterwards.

## Deploy script (`contracts/scripts/deploy.js`)

Settings come from a Vault KV v2 secret, like the paymaster's deploy
script. Bootstrap environment:

| Variable | Required | Example | What it does |
|---|---|---|---|
| `VAULT_ADDR` | yes (unless `MARKET_CONFIG_SOURCE=env`) | `https://vault.internal.trovo.io:8200` | Vault server URL. |
| `VAULT_TOKEN` | yes (unless `MARKET_CONFIG_SOURCE=env`) | `hvs.CAESIJ…` | Token allowed to read the secret. |
| `VAULT_KV_MOUNT` | no (default `secret`) | `secret` | KV v2 mount. |
| `VAULT_NAMESPACE` | no | `trovo` | Vault Enterprise namespace. |
| `MARKET_DEPLOY_SECRET_PATH` | no (default `trovo/market/deploy`) | `trovo/market/deploy` | Path of the secret below. |
| `MARKET_CONFIG_SOURCE` | no (default `vault`) | `env` | `env` reads the keys below from the process environment. **Local development only.** |

Keys of the secret:

| Key | Required | Example | What it does |
|---|---|---|---|
| `RPC_URL` | yes | `https://mainnet.base.org` | Base RPC endpoint. |
| `CHAIN_ID` | yes | `8453` | Checked against the RPC. |
| `DEPLOYER_PRIVATE_KEY` | yes | `0x…` | Pays for deployment; owner only until the end of the script. |
| `OWNER_ADDRESS` | yes | `0x…` | The platform admin Safe that ends up owning the book (lists tokens, sets authorizers, pauses). Must be a contract unless `ALLOW_EOA_OWNER=true`. |
| `ALLOW_EOA_OWNER` | no | `true` | Allow a plain address as owner - testnets only. |
| `AUTHORIZER_ADDRESSES` | yes | `0xA…,0xB…` | Addresses whose fill authorizations the book accepts: the address of app-backend's `OFFER_AUTHORIZER_PRIVATE_KEY`. Several may be active (rotation). |
| `TRADABLE_TOKENS` | no | `0xcNGN…,0xNGNI…` | Tokens to list at deployment: the tokenization payment stablecoins and the internal balance tokens. Each tokenized asset's token is listed when it is created (see INTEGRATION.md). |

Writing the secret:

```bash
vault kv put secret/trovo/market/deploy \
  RPC_URL=https://sepolia.base.org CHAIN_ID=84532 \
  DEPLOYER_PRIVATE_KEY=0x… OWNER_ADDRESS=0x… \
  AUTHORIZER_ADDRESSES=0x… TRADABLE_TOKENS=0x…,0x…
```

A read-only policy for the deploy operator:

```hcl
path "secret/data/trovo/market/deploy" {
  capabilities = ["read"]
}
```

## What app-backend needs afterwards

- `OFFER_BOOK_ADDRESS` - the deployed book (printed by the script and saved
  in `contracts/deployments/<chainId>.json`).
- `OFFER_AUTHORIZER_PRIVATE_KEY` - the key whose address is in
  `AUTHORIZER_ADDRESSES`.

See [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md#tokenization).

# market - configuration

Every parameter the market project reads. Each one says what it is, why it
is needed, whether you must set it, an example and how to get a real value.
Examples are never real keys.

The contracts themselves have no configuration file: they receive their
settings when they are deployed, and from their owner afterwards. What you
configure is the **deploy script** (`contracts/scripts/deploy.js`), which
deploys the offer book (`TrovoOfferBook`).

## How to set them

The deploy script reads its settings from **HashiCorp Vault** (a secure
store for secrets), so private keys never sit in files or shell history:

1. A few **bootstrap environment variables** tell the script where Vault is
   and how to log in ([below](#bootstrap-environment-variables)).
2. Everything else is stored as one **Vault secret** with several keys
   ([below](#keys-of-the-deploy-secret)).

For local experiments only, `MARKET_CONFIG_SOURCE=env` makes the script read
those keys from ordinary environment variables instead of Vault.

### The minimum to deploy

| Where | Parameter |
|---|---|
| environment | [`VAULT_ADDR`](#vault_addr), [`VAULT_TOKEN`](#vault_token) |
| Vault secret | [`RPC_URL`](#rpc_url), [`CHAIN_ID`](#chain_id), [`DEPLOYER_PRIVATE_KEY`](#deployer_private_key), [`OWNER_ADDRESS`](#owner_address), [`AUTHORIZER_ADDRESSES`](#authorizer_addresses) |

---

## Bootstrap environment variables

### `VAULT_ADDR`

- **What it is:** The web address of your Vault server.
- **Why it's needed:** The script reads its settings and keys from Vault.
- **Required:** Yes, unless `MARKET_CONFIG_SOURCE=env`.
- **Example:** `https://vault.internal.trovo.io:8200`
- **How to get it:** The same Vault the other services use (tm-api's
  `VAULT_ADDR`); ask whoever runs Vault.

### `VAULT_TOKEN`

- **What it is:** A Vault access token.
- **Why it's needed:** Vault only returns the secret to a token allowed to
  read it.
- **Required:** Yes, unless `MARKET_CONFIG_SOURCE=env`.
- **Example:** `hvs.CAESIJexampleexampleexample`
- **How to get it:** Ask your Vault administrator for a short-lived token
  with the read-only policy shown in
  [Writing the secret](#writing-the-secret), for example
  `vault token create -policy=market-deploy -ttl=1h`.

### `VAULT_KV_MOUNT`

- **What it is:** The name of Vault's key-value (KV version 2) store.
- **Why it's needed:** The secret's full location is mount + path.
- **Required:** No, default `secret`.
- **Example:** `secret`
- **How to get it:** Keep the default unless your Vault uses another mount
  name (`vault secrets list` shows them).

### `VAULT_NAMESPACE`

- **What it is:** The Vault namespace the secret lives in.
- **Why it's needed:** Only Vault Enterprise and HCP Vault have namespaces;
  without the right one, the secret is not found.
- **Required:** No.
- **Example:** `trovo`
- **How to get it:** Your Vault administrator; leave unset on open-source
  Vault.

### `MARKET_DEPLOY_SECRET_PATH`

- **What it is:** Where in Vault the deploy secret is stored.
- **Why it's needed:** It tells the script which secret to read.
- **Required:** No, default `trovo/market/deploy`.
- **Example:** `trovo/market/deploy-sepolia`
- **How to get it:** Keep the default; use a different path per network
  if you deploy to both a test network and mainnet.

### `MARKET_CONFIG_SOURCE`

- **What it is:** Where the script reads the keys below from: `vault` or
  `env`.
- **Why it's needed:** `env` lets you try a deployment on a local or test
  network without Vault.
- **Required:** No, default `vault`.
- **Example:** `env`
- **How to get it:** Use `env` only on your own machine for a local or
  test network. Never for mainnet.

## Keys of the deploy secret

### `RPC_URL`

- **What it is:** The web address of a Base blockchain node.
- **Why it's needed:** The script sends the deployment transactions
  through it.
- **Required:** Yes.
- **Example:** `https://base-mainnet.g.alchemy.com/v2/your-api-key`
  (mainnet), `https://sepolia.base.org` (test network)
- **How to get it:** Sign up with a node provider (Alchemy, Infura,
  QuickNode), create an app for Base Mainnet or Base Sepolia and copy its
  HTTPS URL. For a test network the public `https://sepolia.base.org`
  works.

### `CHAIN_ID`

- **What it is:** The number of the network you mean to deploy to.
- **Why it's needed:** The script checks the node really is that network,
  so you cannot deploy to the wrong one by mistake.
- **Required:** Yes.
- **Example:** `8453` (Base Mainnet), `84532` (Base Sepolia)
- **How to get it:** `8453` for production, `84532` for testing.

### `DEPLOYER_PRIVATE_KEY`

- **What it is:** The private key of the wallet that sends the deployment
  transactions.
- **Why it's needed:** Someone has to pay the network fees and send the
  transactions. This wallet owns the book only for the few seconds of the
  script; ownership then passes to `OWNER_ADDRESS`.
- **Required:** Yes.
- **Example:** `0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`
  (a well-known test key; never use it for real funds)
- **How to get it:** Create a fresh wallet just for deploying, from
  `market/contracts` after `npm ci`:
  ```bash
  node -e "const {Wallet}=require('ethers');const w=Wallet.createRandom();console.log('address',w.address);console.log('key',w.privateKey)"
  ```
  Send it about 0.005 ETH on Base (on Base Sepolia, use a faucet such as
  the one linked from <https://docs.base.org/tools/network-faucets>).

### `OWNER_ADDRESS`

- **What it is:** The address that owns the offer book after deployment:
  the platform admin Safe.
- **Why it's needed:** The owner lists which tokens may be traded, sets
  which keys may authorize purchases, and can pause trading in an
  emergency. It can never touch anyone's escrowed tokens.
- **Required:** Yes. It must be a contract (a Safe) unless
  `ALLOW_EOA_OWNER=true`.
- **Example:** `0xA11CE00000000000000000000000000000000001`
- **How to get it:** The address of Trovo's platform admin Safe on that
  network. If there isn't one yet, create it at <https://app.safe.global>
  (Create account; owners: the platform administrators; a threshold of at
  least 2).

### `ALLOW_EOA_OWNER`

- **What it is:** Allows a plain wallet (not a Safe) as the owner.
- **Why it's needed:** Convenient on a test network; dangerous on mainnet,
  where one leaked key would control the book.
- **Required:** No, default `false`.
- **Example:** `true`
- **How to get it:** Leave unset for mainnet. Set `true` only on a test
  network.

### `AUTHORIZER_ADDRESSES`

- **What it is:** The addresses whose signatures the offer book accepts as
  permission for a purchase, separated by commas.
- **Why it's needed:** Every purchase needs a signed permission from Trovo,
  given only after app-backend's checks (KYC, sale window, purchase cap).
  The book only accepts permissions signed by these addresses.
- **Required:** Yes (at least one).
- **Example:** `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`
- **How to get it:**
  1. Generate a key with the same command as for `DEPLOYER_PRIVATE_KEY`.
  2. Put its **address** here.
  3. Give its **private key** to app-backend as `OFFER_AUTHORIZER_PRIVATE_KEY`
     (see [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)).
  To rotate later, list both the old and new addresses for a while.

### `TRADABLE_TOKENS`

- **What it is:** Token addresses to allow trading in right away,
  separated by commas.
- **Why it's needed:** The book only trades tokens its owner has listed.
  The payment tokens buyers pay with (stablecoins such as cNGN, and the
  internal balance tokens) must be listed, or nothing can be bought with
  them.
- **Required:** No. Tokens can be listed later by the owner Safe.
- **Example:** `0xC0FFEE0000000000000000000000000000000001,0xC0FFEE0000000000000000000000000000000002`
- **How to get it:** The contract addresses of the payment stablecoins
  (for example cNGN's address from its issuer's documentation or
  [basescan.org](https://basescan.org)) and of the internal balance tokens
  (app-backend's configuration). Each tokenized asset's own token is listed
  later, when that asset is created (see [DEPLOYMENT.md](DEPLOYMENT.md#5-for-each-new-tokenized-asset)).

### Writing the secret

```bash
vault kv put secret/trovo/market/deploy \
  RPC_URL=https://sepolia.base.org CHAIN_ID=84532 \
  DEPLOYER_PRIVATE_KEY=0x... OWNER_ADDRESS=0x... \
  AUTHORIZER_ADDRESSES=0x... TRADABLE_TOKENS=0x...,0x...
```

A read-only policy for the person or machine that deploys:

```hcl
path "secret/data/trovo/market/deploy" {
  capabilities = ["read"]
}
```

Save it as `market-deploy.hcl` and load it with
`vault policy write market-deploy market-deploy.hcl`.

## Local development only

These are read only by the local test setup
([DEPLOYMENT.md](DEPLOYMENT.md#local-blockchain-for-end-to-end-tests)).

### `LOCAL_RPC_URL`

- **What it is:** The address of a local test blockchain.
- **Why it's needed:** Hardhat's `localhost` network uses it to reach the
  local node.
- **Required:** No, default `http://127.0.0.1:8545` (where
  `npx hardhat node` listens).
- **Example:** `http://127.0.0.1:8545`
- **How to get it:** Keep the default unless your local node runs
  elsewhere.

### `OFFER_AUTHORIZER_ADDRESS`

- **What it is:** The authorizer address the local deployment script
  (`scripts/local-market.js`) gives the local offer book.
- **Why it's needed:** The end-to-end tests sign purchase permissions with
  the matching key.
- **Required:** No, default `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`
  (Hardhat's built-in test account #2).
- **Example:** `0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC`
- **How to get it:** Keep the default.

## What other projects need afterwards

| Project | Parameter | Value |
|---|---|---|
| app-backend | `OFFER_BOOK_ADDRESS` | the deployed book, printed by the script and saved in `contracts/deployments/<chainId>.json` |
| app-backend | `OFFER_AUTHORIZER_PRIVATE_KEY` | the private key whose address is in `AUTHORIZER_ADDRESSES` |
| payout-engine | `OFFER_BOOK_ADDRESS` | the same book |

See [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md) and
[payout-engine/CONFIGURATION.md](../payout-engine/CONFIGURATION.md#offer_book_address).

# market - deployment

A step-by-step guide to deploying the market contracts, written for someone
doing it for the first time.

## 1. What you are deploying

Smart contracts on the Base blockchain that tokenized assets are issued and
sold with:

- **`TrovoOfferBook`** (deployed **once per network**): where tokens are
  offered for sale at a fixed price and bought. Every purchase needs a
  signed permission from Trovo's backend.
- **`TokenizedAsset`** (deployed **once per tokenized asset**): the token
  of one asset, which only that asset's issuing wallet can create (mint).

A deployed contract cannot be changed. Deploy to the test network (Base
Sepolia) first and check everything before mainnet. **The contracts have
not been externally audited yet: get them audited before mainnet.**

Nothing else needs to be running first, but you need a Base node URL and,
for mainnet, the platform admin Safe (see step 2).

## 2. Before you start

| Tool / account | Install or get it | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Node.js 20 or 22 (LTS) | <https://nodejs.org> | `node --version` shows `v20` or `v22` |
| Vault CLI (to store the settings) | <https://developer.hashicorp.com/vault/install> | `vault --version` |
| A Base node URL | a provider such as Alchemy or QuickNode | see [`RPC_URL`](CONFIGURATION.md#rpc_url) |
| A deployer wallet with a little ETH | see [`DEPLOYER_PRIVATE_KEY`](CONFIGURATION.md#deployer_private_key) | its balance on [basescan.org](https://basescan.org) |
| The platform admin Safe | <https://app.safe.global> | see [`OWNER_ADDRESS`](CONFIGURATION.md#owner_address) |

## 3. Get the code, build and test

1. Get the code and install the dependencies:
   ```bash
   git clone https://github.com/richardsric7/base-trovoapp-service.git
   cd base-trovoapp-service/market/contracts
   npm ci
   ```
   **You should see:** `added ... packages` and no `ERR!` lines.
2. Compile:
   ```bash
   npm run build
   ```
   **You should see:** `Compiled ... Solidity files successfully`. The
   compiler (Solidity 0.8.28) comes from an npm package, so nothing is
   downloaded.
3. Run the tests:
   ```bash
   npm test
   ```
   **You should see:** every test passing (18 or more), `0 failing`.

## 4. Deploy the offer book (once per network)

1. **Create the authorizer key** that app-backend will sign purchase
   permissions with (see [`AUTHORIZER_ADDRESSES`](CONFIGURATION.md#authorizer_addresses)):
   ```bash
   node -e "const {Wallet}=require('ethers');const w=Wallet.createRandom();console.log('address',w.address);console.log('key',w.privateKey)"
   ```
   Keep the key safe; it goes to app-backend in step 6.
2. **Store the settings in Vault** (each key is explained in
   [CONFIGURATION.md](CONFIGURATION.md#keys-of-the-deploy-secret)):
   ```bash
   export VAULT_ADDR=https://vault.internal.trovo.io:8200
   vault login
   vault kv put secret/trovo/market/deploy \
     RPC_URL=https://sepolia.base.org CHAIN_ID=84532 \
     DEPLOYER_PRIVATE_KEY=0x... OWNER_ADDRESS=0x... \
     AUTHORIZER_ADDRESSES=0x... TRADABLE_TOKENS=0x...,0x...
   ```
3. **Do a dry run.** It checks everything and deploys nothing:
   ```bash
   export VAULT_TOKEN=hvs....     # a token that can read the secret
   npm run deploy:dry-run
   ```
   **You should see:**
   ```
   config: vault secret/trovo/market/deploy
   chain 84532, deployer 0x... (0.01 ETH)
   owner 0x..., authorizers 0x...
   tradable tokens: 0x..., 0x...
   dry run: nothing deployed
   ```
   Check the chain, owner and authorizers are what you expect, and that
   the deployer has ETH.
4. **Deploy:**
   ```bash
   npm run deploy
   ```
   **You should see:** `deployed TrovoOfferBook at 0x...`, a transaction
   for each listed token and for the ownership transfer,
   `wrote .../deployments/84532.json` and
   `next: set OFFER_BOOK_ADDRESS=0x... in app-backend`.
5. **Keep the record.** Commit `contracts/deployments/<chainId>.json`; it
   is the record of what was deployed where.

## 5. For each new tokenized asset

Each asset gets its own token contract, owned by the asset's issuing
wallet. The issuing wallet's address appears in Trovo Manager once the
asset's details are submitted.

1. **Deploy the token.** There is no script for this; the simplest way is
   Remix in the browser:
   1. Open <https://remix.ethereum.org>, create `TokenizedAsset.sol` and
      paste the contents of [`contracts/src/TokenizedAsset.sol`](contracts/src/TokenizedAsset.sol).
   2. In **Solidity compiler**, choose version `0.8.28`, open **Advanced
      configurations**, set EVM version `cancun`, turn on optimization
      (200 runs), and compile.
   3. In **Deploy & run transactions**, choose **Injected Provider**
      (your browser wallet, switched to Base or Base Sepolia).
   4. Fill in the constructor:

      | Field | What to enter | Example |
      |---|---|---|
      | `name_` | the asset's name | `Lekki Gardens Phase 2` |
      | `symbol_` | the asset code, **exactly** as in Trovo | `LKG2` |
      | `decimals_` | the token's decimals agreed for the asset | `6` |
      | `initialOwner` | the asset's **issuing wallet** address from Trovo Manager | `0x...` |

   5. Click **Deploy**, confirm in your wallet and copy the new contract's
      address.
2. **Register it** in Trovo Manager on the asset's tokenization page
   (contract address field). app-backend checks the contract: it exists,
   its symbol is the asset code, its supply is 0 and the issuing wallet
   can mint.
3. **List it on the offer book.** In <https://app.safe.global>, open the
   platform admin Safe, choose **New transaction → Transaction Builder**,
   enter the offer book address, paste the offer book's ABI (from
   `contracts/artifacts/src/TrovoOfferBook.sol/TrovoOfferBook.json`, the
   `abi` part), choose `setTradable`, enter the token address and `true`,
   then create, sign and execute it with the Safe's owners. app-backend
   refuses to mint an unlisted token (`error-token-not-listed`).
4. Minting is then done from Trovo Manager by the minting approvers.

## 6. Set the parameters in other projects

| Project | Parameter | Value |
|---|---|---|
| app-backend | `OFFER_BOOK_ADDRESS` | the address from step 4 |
| app-backend | `OFFER_AUTHORIZER_PRIVATE_KEY` | the key from step 4.1 |
| payout-engine | `OFFER_BOOK_ADDRESS` | the address from step 4 |

Restart those services after changing them.

## 7. Check it works

1. Open the offer book on [basescan.org](https://basescan.org) (or
   [sepolia.basescan.org](https://sepolia.basescan.org)): under **Read
   Contract**, `owner` is the admin Safe, `authorizers(<your authorizer
   address>)` is `true`, and `tradable(<a payment token>)` is `true`.
2. On the test network, run a test asset end to end: register its token,
   list it, mint it from Trovo Manager, and buy some from a test account
   in the app.

## 8. Updating and rolling back

The contracts cannot be upgraded. To change the offer book you deploy a new
one, list the tokens on it, point app-backend's and payout-engine's
`OFFER_BOOK_ADDRESS` at it, and leave the old one for sellers to cancel
their offers (cancelling always works, even when paused).

Everyday changes do not need a redeploy; the owner Safe makes them:

- **Rotate an authorizer:** `setAuthorizer(new, true)`, switch app-backend's
  key, then `setAuthorizer(old, false)`.
- **Stop trading in an emergency:** `pause()`; `unpause()` to resume.
  Sellers can still cancel while paused.
- **Stop trading one token:** `setTradable(token, false)`.

## 9. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `VAULT_ADDR is not set` / `VAULT_TOKEN is not set` | the bootstrap variables are missing | `export` them in the same terminal |
| `RPC_URL is required` (or another key) | the key is missing from the Vault secret | `vault kv get secret/trovo/market/deploy` and add it |
| `RPC_URL is chain 8453, CHAIN_ID says 84532` | the node and `CHAIN_ID` are different networks | fix whichever is wrong |
| `OWNER_ADDRESS ... has no code` | the owner is a plain wallet, not a Safe | use the admin Safe's address (or `ALLOW_EOA_OWNER=true` on a test network) |
| `no token contract at 0x...` | a `TRADABLE_TOKENS` address is wrong or on another network | check it on basescan for that network |
| `... not found - run "npm run build" first` | not compiled | `npm run build` |
| `insufficient funds` | the deployer has no ETH | send it ETH and run again |
| app-backend: `error-token-not-listed` | the asset's token is not listed | step 5.3 |
| app-backend: `error-invalid-token-contract` | wrong symbol, owner or supply | redeploy the token with the asset code as symbol and the issuing wallet as owner |

## Local blockchain for end-to-end tests

app-backend's and payout-engine's end-to-end tests use the market contracts
on a local blockchain started from the paymaster project:

```bash
# terminal 1: start a local blockchain (leave it running)
cd paymaster/contracts && npx hardhat node

# terminal 2: deploy Safe, the EntryPoint and the paymaster, then the market
cd paymaster/contracts && npx hardhat run scripts/local-stack.js --network localhost
cd ../../market/contracts && npx hardhat run scripts/local-market.js --network localhost
```

`local-market.js` deploys the book (owner: Hardhat account #0; authorizer:
[`OFFER_AUTHORIZER_ADDRESS`](CONFIGURATION.md#offer_authorizer_address),
Hardhat account #2 by default), a test cNGN and a test internal balance
token, lists them and writes `deployments/local-market.json`. Then, from
`app-backend/`:

```bash
AA_LOCAL_STACK=$PWD/../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_MARKET=$PWD/../market/contracts/deployments/local-market.json \
go test ./internal/components/users/services/ -run TestTokenizationOnLocalChain -v
```

- **`AA_LOCAL_STACK`:** the file listing the local Safe, EntryPoint and
  paymaster addresses (written by `local-stack.js`).
- **`AA_LOCAL_MARKET`:** the file listing the local offer book and tokens
  (written by `local-market.js`).

The test runs app-backend's own mint, pricing and authorization code: an
issuing Safe mints and opens the offer, a new wallet buys with cNGN, a
reused authorization is refused, and a minting Safe makes a fiat purchase
paid in the internal balance token.

## Verification done so far

- 18 Hardhat tests: escrow, pricing and rounding, purchases paid in a
  stablecoin and in the internal balance token, every authorization check
  (wrong key, someone else's authorization, any changed field, reused
  nonce, expired), limits, delisting and pause, repricing, cancelling,
  fee-on-transfer tokens, owner-only administration.
- app-backend's signatures are accepted by the contract in the end-to-end
  test.
- Not yet: a Base Sepolia deployment, and an external audit.

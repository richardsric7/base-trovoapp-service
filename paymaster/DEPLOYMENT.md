# paymaster — deployment

A step-by-step guide to deploying the paymaster, written for someone doing
it for the first time. Nothing here has been deployed to a public network
yet; this is the procedure.

## 1. What you are deploying

Two parts that together let users pay their wallet's network fees ("gas")
in a stablecoin (cNGN, USDC, USDT) instead of ETH:

- **`TrovoTokenPaymaster`**, a smart contract on Base, deployed **once per
  network**. It pays the network in ETH from its own deposit and charges
  the user's wallet the equivalent in the chosen stablecoin.
- **The quote service**, a small web service you run **all the time**. It
  works out the current price of ETH in each stablecoin, adds Trovo's
  spread and signs a price ("quote") for each operation. app-backend asks
  it for a quote whenever a user pays gas in a stablecoin.

**What must exist first:**

| What | Why |
|---|---|
| A Base node URL | to deploy the contract and read prices |
| HashiCorp Vault | holds the settings and keys of both parts |
| The treasury Safe | owns the paymaster and its collected funds |
| An ERC-4337 bundler (self-hosted) | submits users' operations; see [INTEGRATION.md](INTEGRATION.md#4-the-bundler) |
| app-backend | the only caller of the quote service; configured in step 10 |

**Get the contract externally audited before mainnet**
(`contracts/src/TrovoTokenPaymaster.sol`). Deploy to Base Sepolia first.

## 2. Before you start

| Tool | Why | Install | Check |
|---|---|---|---|
| Git | get the code | <https://git-scm.com/downloads> | `git --version` |
| Node.js 20 or newer, with npm | build, test and deploy the contract | <https://nodejs.org> | `node --version` |
| Go 1.26 or newer | build the quote service | <https://go.dev/dl> | `go version` |
| Vault CLI | write the secrets | <https://developer.hashicorp.com/vault/install> | `vault --version` |
| Docker (recommended) | run the quote service | <https://docs.docker.com/get-docker> | `docker --version` |
| Foundry `cast` (optional) | create keys, read contracts | <https://book.getfoundry.sh/getting-started/installation> | `cast --version` |

The contract compiles with the `solc` npm package, so no compiler is
downloaded during the build.

## 3. Get the code and run the tests

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/paymaster/contracts
npm ci
npm test          # 24 tests against the official EntryPoint v0.7
npm run build

cd ../quote-service
go test ./...
```

**You should see:** `24 passing` for the contract, and `ok` for each Go
package.

## 4. Create the keys and addresses

You need four things. Write each down where indicated.

| What | How | Goes to |
|---|---|---|
| **Treasury Safe** (the owner) | <https://app.safe.global> → Create account, owners = finance and platform signers, threshold ≥ 2 | deploy secret [`OWNER_ADDRESS`](CONFIGURATION.md#owner_address) |
| **Pauser** | an on-call operations wallet, or a Safe needing one signature | deploy secret [`PAUSER_ADDRESS`](CONFIGURATION.md#pauser_address) |
| **Quote signer key** | the command below | its **address** → deploy secret [`QUOTE_SIGNER_ADDRESSES`](CONFIGURATION.md#quote_signer_addresses); its **key** → quote service secret [`QUOTE_SIGNER_PRIVATE_KEY`](CONFIGURATION.md#quote_signer_private_key) |
| **Deployer key** | the command below; fund it with the stake + deposit + ~0.01 ETH | deploy secret [`DEPLOYER_PRIVATE_KEY`](CONFIGURATION.md#deployer_private_key) |

To create a key, from `paymaster/contracts`:

```bash
node -e "const {Wallet}=require('ethers');const w=Wallet.createRandom();console.log('address',w.address);console.log('key',w.privateKey)"
```

(or `cast wallet new`). Store each key straight into Vault; never in a file
in the repository.

## 5. Write the deploy secret

Put the settings in a file `deploy.json` (each key is explained in
[CONFIGURATION.md section 2](CONFIGURATION.md#2-deploy-script-secret);
every value is a string):

```json
{
  "RPC_URL": "https://sepolia.base.org",
  "CHAIN_ID": "84532",
  "DEPLOYER_PRIVATE_KEY": "0x...",
  "OWNER_ADDRESS": "0x...",
  "QUOTE_SIGNER_ADDRESSES": "0x...",
  "PAUSER_ADDRESS": "0x...",
  "STAKE_AMOUNT_ETH": "0.1",
  "INITIAL_DEPOSIT_ETH": "0.05",
  "GAS_TOKENS": "[{\"symbol\":\"USDC\",\"address\":\"0x...\",\"minRate\":\"500000000\",\"maxRate\":\"20000000000\"}]"
}
```

Then:

```bash
export VAULT_ADDR=https://vault.internal.trovo.io:8200
vault login
vault kv put -mount=secret trovo/paymaster/deploy @deploy.json
rm deploy.json
```

## 6. Deploy the contract

```bash
cd paymaster/contracts
export VAULT_TOKEN=hvs....        # a token that can read trovo/paymaster/deploy
npm run deploy:dry-run            # checks everything, deploys nothing
```

**You should see** the plan, ending with `dry run: nothing deployed`:

```
config: vault secret/trovo/paymaster/deploy
chain 84532, deployer 0x... (0.2 ETH), EntryPoint 0x0000000071727De22E5E9d8BAf0edAc6f37da032
owner 0x..., pauser 0x..., quote signers 0x...
stake 0.1 ETH (unstake delay 86400s), deposit 0.05 ETH, postOp overhead 45000 gas
gas token USDC 0x...: rate bounds [500000000, 20000000000] base units per ETH
dry run: nothing deployed
```

Check every line, then deploy:

```bash
npm run deploy
```

The script:

1. checks the network, that the EntryPoint is v0.7, that each token and the
   owner are contracts, and that the deployer has enough ETH;
2. deploys `TrovoTokenPaymaster` with the deployer as temporary owner;
3. stakes at the EntryPoint and funds the deposit;
4. enables each gas token with its rate bounds;
5. hands ownership to the treasury Safe and checks it took effect;
6. writes `deployments/<chainId>.json`.

**You should see** `deployed TrovoTokenPaymaster at 0x...`, one line per
setup transaction, `wrote .../deployments/84532.json` and
`next: set PAYMASTER_ADDRESS=0x... in the quote service's Vault secret`.

Then:

- commit `contracts/deployments/<chainId>.json` (the record of the
  deployment);
- move any ETH left on the deployer to the treasury, and delete
  `DEPLOYER_PRIVATE_KEY` from Vault.

## 7. Write the quote service secret

Create `quote-service.json` (each key is explained in
[CONFIGURATION.md section 3](CONFIGURATION.md#3-quote-service-secret); see
the full example there):

```json
{
  "RPC_URL": "https://sepolia.base.org",
  "CHAIN_ID": "84532",
  "PAYMASTER_ADDRESS": "0x... (from step 6)",
  "QUOTE_SIGNER_PRIVATE_KEY": "0x... (from step 4)",
  "API_KEYS": "<openssl rand -hex 32>",
  "GAS_TOKENS": "[{\"symbol\":\"USDC\",\"address\":\"0x...\",\"decimals\":6,\"route\":[\"ETH/USD\",\"USD/USDC\"],\"spreadBps\":100}]",
  "RATE_PAIRS": "{\"ETH/USD\":{\"sources\":[{\"name\":\"chainlink\",\"type\":\"chainlink\",\"feed\":\"0x...\"}]},\"USDC/USD\":{\"sources\":[{\"name\":\"peg\",\"type\":\"fixed\",\"value\":\"1\"}]}}",
  "CONFIG_RELOAD_INTERVAL": "5m",
  "DEPOSIT_LOW_WATERMARK_ETH": "0.05",
  "ALERT_WEBHOOK_URL": "https://discord.com/api/webhooks/..."
}
```

```bash
vault kv put -mount=secret trovo/paymaster/quote-service @quote-service.json
rm quote-service.json
```

Create the service's read-only policy and token
([CONFIGURATION.md section 5](CONFIGURATION.md#5-writing-the-secrets-and-the-vault-policies)).

## 8. Run the quote service

**With Docker (recommended):**

```bash
docker build -t trovo-paymaster-quote-service paymaster/quote-service
docker run -d --name quote-service --restart unless-stopped -p 8090:8090 \
  -e VAULT_ADDR -e VAULT_TOKEN -e QUOTE_SERVICE_SECRET_PATH=trovo/paymaster/quote-service \
  trovo-paymaster-quote-service
```

**Or directly:**

```bash
cd paymaster/quote-service
go build -o quote-service .
VAULT_ADDR=... VAULT_TOKEN=... ./quote-service
```

**You should see** (`docker logs quote-service`):

```
quote service: config from vault secret/trovo/paymaster/quote-service, chain 84532, paymaster 0x..., signer 0x...
[startup] USDC: 1 ETH = 3012.45 USDC at market, 3042.57 quoted (spread 100 bps)
listening on :8090
```

At start the service refuses to run if the node is another network than
`CHAIN_ID`, if `PAYMASTER_ADDRESS` is not the paymaster, if the paymaster's
EntryPoint differs, or if the signer key is not one of the paymaster's
quote signers. The error says which.

Run **at least two copies** behind a load balancer in production; they
share nothing, so any number can run.

## 9. Check it works

```bash
curl http://localhost:8090/healthz                     # 200: the process is up
curl http://localhost:8090/readyz                      # 200: every token has a fresh price and the deposit was read
curl -H "X-API-Key: <one of API_KEYS>" http://localhost:8090/v1/status   # signer, paymaster, deposit
curl -H "X-API-Key: <one of API_KEYS>" http://localhost:8090/v1/tokens   # price of 1 ETH in each gas token
```

Use `/readyz` as the load balancer's health check. The API reference is at
`/swagger/index.html`.

## 10. Connect app-backend and the bundler

Set in app-backend (see [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)):

| Parameter | Example |
|---|---|
| `PAYMASTER_ADDRESS` | the address from step 6 |
| `PAYMASTER_QUOTE_SERVICE_URL` | `http://quote-service.internal:8090` |
| `PAYMASTER_QUOTE_SERVICE_API_KEY` | one of the `API_KEYS` |
| `BUNDLER_URL` | your bundler |

Configure the bundler as described in
[INTEGRATION.md](INTEGRATION.md#4-the-bundler), then on Base Sepolia send a
payment from a test wallet with its gas asset set to the stablecoin, and
check the paymaster's `GasPaidInToken` event on basescan.

## 11. Day-to-day operations

**Keep the deposit funded.** When the deposit falls below
`DEPOSIT_LOW_WATERMARK_ETH`, the service logs it and posts to
`ALERT_WEBHOOK_URL`. Then treasury:

1. from the treasury Safe, calls `withdrawToken(<token>, <treasury>, <amount>)`
   to take out collected stablecoins;
2. converts them to ETH (exchange or DEX);
3. calls `deposit()` on the paymaster with the ETH (anyone may).

`GET /v1/status` shows the current deposit.

**Rotate the quote signer:**

1. Treasury Safe: `setQuoteSigner(<new address>, true)`.
2. Put the new key in the quote service secret; it is applied on the next
   reload (`CONFIG_RELOAD_INTERVAL`, or `docker kill -s HUP quote-service`).
3. After `QUOTE_MAX_VALIDITY` (24 hours by default), when old quotes have
   expired: `setQuoteSigner(<old address>, false)`.

**Other owner actions** (treasury Safe transactions, made in Safe{Wallet}
with **New transaction → Transaction Builder** and the ABI from
`contracts/artifacts/src/TrovoTokenPaymaster.sol/TrovoTokenPaymaster.json`):
`setToken` (add a token or change its rate bounds), `setPauser`,
`setPostOpGasOverhead`, `withdrawTo`, `unpause`.

**Emergency:** the pauser calls `pause()`. Users then pay gas in ETH, or
wait. Only the owner can `unpause()`.

## 12. Updating and rolling back

- **Quote service:** rebuild and replace the container (`git pull`,
  `docker build ...`, `docker rm -f quote-service`, `docker run ...`).
  Rolling back is the same with the previous release. Setting changes
  need no redeploy: edit the secret and let it reload.
- **Contract:** it cannot be upgraded. A new version is a new deployment
  (steps 5–6), then the quote service's and app-backend's
  `PAYMASTER_ADDRESS` change. Withdraw the old paymaster's tokens and
  deposit, and unlock and withdraw its stake after the unstake delay.

## 13. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `VAULT_ADDR is not set` / `VAULT_TOKEN is not set` | bootstrap variables missing | set them in the same shell or container |
| deploy: `RPC_URL is chain ..., CHAIN_ID says ...` | node and `CHAIN_ID` disagree | fix the wrong one |
| deploy: owner `has no code` | the owner is not a Safe | use the treasury Safe (or `ALLOW_EOA_OWNER=true` on a test network) |
| deploy: `insufficient funds` | the deployer lacks ETH | fund it with stake + deposit + ~0.01 ETH |
| quote service: signer is not a quote signer | `QUOTE_SIGNER_PRIVATE_KEY` does not match `QUOTE_SIGNER_ADDRESSES` | use the matching key, or add its address with `setQuoteSigner` |
| quote service: cannot parse `GAS_TOKENS` / `RATE_PAIRS` | written as JSON objects instead of strings | write them as JSON strings ([CONFIGURATION.md section 5](CONFIGURATION.md#5-writing-the-secrets-and-the-vault-policies)) |
| `/readyz` returns 503 | a token has no fresh price, or the deposit has not been read | `GET /v1/rates` shows which source fails |
| quotes refused with `rate-out-of-bounds` | the market moved outside the paymaster's rate bounds | treasury Safe: `setToken` with wider bounds |
| quotes refused with `paymaster-deposit-low` | the deposit cannot cover the operation | top up the deposit (step 11) |
| app-backend gets `401 unauthorized` | API key mismatch | app-backend's `PAYMASTER_QUOTE_SERVICE_API_KEY` must be one of `API_KEYS` |
| bundler rejects with `AA34` | the operation changed after it was quoted | a bug in the caller: nothing may change after the quote |

## Local development stack

`scripts/local-stack.js` deploys everything a Trovo wallet needs to a local
blockchain (EntryPoint v0.7, Safe v1.4.1 with the Safe4337Module, a test
USDC and the paymaster) and runs a small **development** bundler, so
app-backend and the quote service can run end to end on your machine:

```bash
cd paymaster/contracts
npx hardhat node                                              # terminal 1 (leave running)
npx hardhat run scripts/local-stack.js --network localhost     # terminal 2 (leave running: it is the bundler)
```

It prints, and writes to `deployments/local-stack.json`, the settings for
app-backend (`ENTRYPOINT_ADDRESS`, `SAFE_*_ADDRESS`, `BUNDLER_URL`,
`PAYMASTER_ADDRESS`, `BASE_CHAIN_ID`) and the quote service
(`PAYMASTER_ADDRESS`, `ENTRYPOINT_ADDRESS`, the test USDC for
`GAS_TOKENS`). The local quote signer is Hardhat account #1 (key
`0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d`).
The development bundler sends each operation alone with fixed gas
estimates: never use it outside development. Its settings
(`LOCAL_RPC_URL`, `DEV_BUNDLER_PORT`, `QUOTE_SIGNER_ADDRESS`) are in
[CONFIGURATION.md](CONFIGURATION.md#local-stack-contractsscriptslocal-stackjs).

app-backend's `internal/aa` package has an end-to-end test against this
stack (activation paid in USDC, a second USDC-paid send, an ETH-paid send,
a sub-wallet, a linked distribution Safe, a 2-of-3 approved payment and the
gas-debt cycle). Run it from `app-backend/` with the quote service running
against the local stack:

```bash
AA_LOCAL_STACK=$PWD/../paymaster/contracts/deployments/local-stack.json \
AA_RPC_URL=http://127.0.0.1:8545 \
AA_QUOTE_SERVICE_URL=http://127.0.0.1:8090 AA_QUOTE_SERVICE_API_KEY=<key> \
go test ./internal/aa/ -run TestLocalStackEndToEnd -v
```

- **`AA_LOCAL_STACK`:** the file written by `local-stack.js`.
- **`AA_RPC_URL`:** the local blockchain.
- **`AA_QUOTE_SERVICE_URL`:** the locally running quote service.
- **`AA_QUOTE_SERVICE_API_KEY`:** one of its `API_KEYS`.

## Verification done so far

- Contract: 24 Hardhat tests against the official EntryPoint v0.7
  (charging and refunds, the activation charge, debt, reverted operations,
  every rejection path, pause, signer rotation, owner-only administration).
- Deploy script: run end to end against a local chain with a test Vault.
- Quote service: unit tests, including matching the contract's quote hash
  and signatures; an end-to-end run where its quotes paid for a transfer
  (cNGN) and a new wallet's activation (USDC).
- app-backend's operation builder: end to end against the local stack with
  the real quote service.
- Not yet: Base Sepolia with a real bundler, and an external audit.

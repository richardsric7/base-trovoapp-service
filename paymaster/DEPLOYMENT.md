# paymaster — deployment

Nothing here has been deployed yet. This is the procedure; every value it
mentions is described in [CONFIGURATION.md](CONFIGURATION.md).

## Prerequisites

| Tool | Why | Install |
|---|---|---|
| Node.js ≥ 20 + npm | Hardhat build, tests, deploy script | https://nodejs.org |
| Go ≥ 1.26 | Quote service | https://go.dev/dl |
| Vault CLI | Writing the secrets | https://developer.hashicorp.com/vault/install |
| Foundry `cast` (optional) | Generating keys, reading contracts | https://book.getfoundry.sh/getting-started/installation |
| Docker (optional) | Container image for the quote service | https://docs.docker.com/get-docker |

The Hardhat config compiles with the `solc` npm package (solc-js) instead
of downloading a compiler binary, so builds work behind restrictive proxies.

## Order of operations

1. **Audit the contract** (`contracts/src/TrovoTokenPaymaster.sol`) before
   mainnet.
2. Create the treasury **owner Safe** and choose a **pauser** address.
3. Generate the **quote signer key** (`cast wallet new`) — you need its
   address for the deployment and its key for the quote service.
4. Deploy the contract (below).
5. Configure and start the quote service (below).
6. Configure app-backend to use it (see [INTEGRATION.md](INTEGRATION.md)).
7. Run the bundler with the paymaster allowed (see INTEGRATION.md).

## Deploying the contract

```bash
cd paymaster/contracts
npm ci
npm test                  # 20 tests against the official EntryPoint v0.7
npm run build

export VAULT_ADDR=https://vault.internal.trovo.io:8200
export VAULT_TOKEN=…      # a token allowed to read trovo/paymaster/deploy
export PAYMASTER_DEPLOY_SECRET_PATH=trovo/paymaster/deploy   # default

npm run deploy:dry-run    # prints the plan and checks everything; deploys nothing
npm run deploy
```

The script:

1. reads the secret and checks the chain id, that the EntryPoint is a v0.7
   EntryPoint, that each token has code, that the owner is a contract (a
   Safe), and that the deployer holds enough ETH;
2. deploys `TrovoTokenPaymaster` with the deployer as temporary owner;
3. stakes (`addStake`) and funds the deposit (`deposit`);
4. enables every token in `GAS_TOKENS` with its rate bounds (`setToken`);
5. transfers ownership to `OWNER_ADDRESS` and verifies it;
6. writes `deployments/<chainId>.json` (addresses, configuration and
   transaction hashes) — commit it.

Then put `PAYMASTER_ADDRESS` into the quote service's secret and remove the
deployer key from Vault.

Later changes (new token, new bounds, signer rotation, top-ups,
withdrawals) are owner-Safe transactions: `setToken`, `setQuoteSigner`,
`setPauser`, `setPostOpGasOverhead`, `deposit` (anyone can deposit),
`withdrawToken`, `withdrawTo`, `unpause`.

### Rotating the quote signer

1. Owner Safe: `setQuoteSigner(<new>, true)`.
2. Put the new key in the quote service secret and reload (`SIGHUP` or
   `CONFIG_RELOAD_INTERVAL`).
3. After `QUOTE_MAX_VALIDITY` has passed (outstanding quotes expired):
   `setQuoteSigner(<old>, false)`.

## Running the quote service

```bash
cd paymaster/quote-service
go test ./...
go build -o quote-service .

export VAULT_ADDR=… VAULT_TOKEN=… QUOTE_SERVICE_SECRET_PATH=trovo/paymaster/quote-service
./quote-service
```

Or with Docker:

```bash
docker build -t trovo-paymaster-quote-service paymaster/quote-service
docker run -p 8090:8090 -e VAULT_ADDR -e VAULT_TOKEN -e QUOTE_SERVICE_SECRET_PATH trovo-paymaster-quote-service
```

At startup the service refuses to run if the RPC chain id differs from
`CHAIN_ID`, if `PAYMASTER_ADDRESS` does not answer like the paymaster, if
its EntryPoint differs from `ENTRYPOINT_ADDRESS`, or if the signer key is
not one of the paymaster's quote signers. It then fetches every rate and
logs each token's market and quoted price of 1 ETH.

Health checks:

- `GET /healthz` — liveness (always 200 while the process serves).
- `GET /readyz` — 200 only when every gas token has a fresh rate and the
  deposit has been read; use it as the readiness probe.

Swagger UI: `/swagger/index.html`. Regenerate the docs after changing an
endpoint: `swag init -g main.go -o docs` (from `quote-service/`).

Run at least two replicas behind a load balancer; they are stateless
(every replica refreshes its own rates).

## Keeping the deposit funded

The service logs and posts to `ALERT_WEBHOOK_URL` when the deposit drops
below `DEPOSIT_LOW_WATERMARK_ETH`. Treasury then:

1. Owner Safe: `withdrawToken(<token>, <treasury>, <amount>)` for the
   collected stablecoins.
2. Converts them to ETH (exchange or DEX).
3. Calls `deposit()` on the paymaster with the ETH (any account may).

`GET /v1/status` shows the current deposit.

## Verification done so far

- Contract: 20 Hardhat tests against the official EntryPoint v0.7
  (charging and refunds, activation charge, debt, reverted operations,
  every rejection path, pause, signer rotation, owner-only administration).
- Deploy script: run end to end against a local chain with a mock Vault.
- Quote service: unit tests, including matching the contract's `getHash`
  and ethers' signatures from a generated fixture; and an end-to-end run
  against a local chain in which the service's quotes paid for a deployed
  wallet's transfer (cNGN) and a new wallet's activation (USDC).
- Not yet: a deployment to Base Sepolia with a real bundler, and an
  external audit.

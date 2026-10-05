# Deployment

A step-by-step guide to getting `payout-engine` running, written for
someone doing it for the first time.

## 1. What you are deploying

`payout-engine` is a background program (no website, no API) that pays a
tokenized asset's proceeds (dividends, rent, interest) to the people who
hold the asset's tokens. Trovo admins start and approve each payout in
Trovo Manager; the engine does the work: it lists the holders, waits for
approval, checks the money is there and sends the payments on the Base
blockchain from one dedicated Safe wallet.

**It needs these to be running first:**

| What | Why | Where to set it up |
|---|---|---|
| app-backend, with its Postgres database | The engine works in app-backend's database; app-backend creates the payout tables when it starts. | [app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md) |
| tm-api and tm-web (Trovo Manager) | Admins create, approve and fund payouts there. | [tm-api/DEPLOYMENT.md](../tm-api/DEPLOYMENT.md), [tm-web/DEPLOYMENT.md](../tm-web/DEPLOYMENT.md) |
| A Base blockchain node (RPC) | To read token transfers and send payments. | a provider such as Alchemy or QuickNode |
| Redis (optional) | Instant wake-ups from Trovo Manager. Without it the engine checks every 10 seconds. | the same Redis tm-api uses |

## 2. Before you start

Install these tools and check each one works:

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Go 1.26 or newer (to build from source) | <https://go.dev/doc/install> | `go version` shows `go1.26` or newer |
| Docker (to build and run the container) | <https://docs.docker.com/get-docker/> | `docker --version` |

And have these ready (each is explained in
[CONFIGURATION.md](CONFIGURATION.md)):

- app-backend's database connection string.
- A Base RPC URL.
- Access to Trovo Manager as an admin, and to the vault manager (Vault
  Signer) there.
- A little ETH on Base to pay network fees (see step 4).

## 3. Get the code and build it

1. Get the code:
   ```bash
   git clone https://github.com/richardsric7/base-trovoapp-service.git
   cd base-trovoapp-service/payout-engine
   ```
2. Build it:
   ```bash
   go build ./...
   ```
   **You should see:** no output. Any output is an error.
3. Run the checks and tests:
   ```bash
   go vet ./...
   go test ./...
   ```
   **You should see:** `ok` for each package. The on-chain end-to-end test
   says it is skipped; that is expected (see
   [Run the tests](#run-the-tests) to run it).

## 4. Create the payout Safe and its keys (once)

The engine pays from a Safe that needs 3 signatures for every payment.

1. **Create the signing keys.** In Trovo Manager open the vault manager
   (Vault Signer) and create a managed secret with at least 4 keys. Note
   the addresses of the first three; they will be the Safe's owners.
2. **Create the Safe.** Open <https://app.safe.global>, connect any wallet
   on Base (Base Sepolia for testing), choose **Create account**, add the
   signer addresses as owners and set the threshold to **3**. Copy the new
   Safe's address: this is
   [`PROCEED_PAYOUT_SAFE_ADDRESS`](CONFIGURATION.md#proceed_payout_safe_address).
3. **Link them.** Back in the vault manager, set the managed secret's
   wallet address to the Safe. Rotating a key there now also swaps the
   Safe's owner.
4. **Fund the first signer with ETH.** The first signer sends every payment
   transaction and pays its network fee. Send it about 0.01 ETH on Base to
   start. The funding check in Trovo Manager tells you if more is needed.
5. **Set the fee wallet.** In Trovo Manager open **Proceeds Payouts → Fee &
   engine** and set the payout fee wallet and default fee.
6. **Check the payout currency.** Each payout currency (for example CNGN)
   must be in app-backend's `tokenization_currencies` table with its token
   contract address.

## 5. Set the parameters

1. Copy the template:
   ```bash
   cp .env.example .env
   ```
2. Open `.env` in a text editor and fill in at least these (each is
   explained, with an example and where to get it, in
   [CONFIGURATION.md](CONFIGURATION.md)):

   | Parameter | Example |
   |---|---|
   | [`DB_CONNECTION_STRING`](CONFIGURATION.md#db_connection_string) | `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require` |
   | [`BASE_RPC_URL`](CONFIGURATION.md#base_rpc_url) | `https://base-mainnet.g.alchemy.com/v2/your-api-key` |
   | [`BASE_CHAIN_ID`](CONFIGURATION.md#base_chain_id) | `8453` |
   | [`PROCEED_PAYOUT_SAFE_ADDRESS`](CONFIGURATION.md#proceed_payout_safe_address) | the Safe from step 4 |
   | [`PROCEED_PAYOUT_SIGNERS_VAULT_PATH`](CONFIGURATION.md#proceed_payout_signers_vault_path) + [`VAULT_ADDR`](CONFIGURATION.md#vault_addr) + [`VAULT_TOKEN`](CONFIGURATION.md#vault_token), **or** [`PROCEED_PAYOUT_SIGNERS`](CONFIGURATION.md#proceed_payout_signers) | `secret/trovo/payout-engine#PROCEED_PAYOUT_SIGNERS` |

   Recommended as well: [`OFFER_BOOK_ADDRESS`](CONFIGURATION.md#offer_book_address),
   [`VAT_WALLET`](CONFIGURATION.md#vat_wallet), the
   [Redis settings](CONFIGURATION.md#redis-optional) and
   [`GC`](CONFIGURATION.md#gc) (push notifications).

Never commit `.env` to git: it holds secrets.

## 6. Run it

Pick one.

**A. Directly with Go** (good for a first try):

```bash
go run .
```

**B. With Docker** (recommended for servers):

```bash
docker build -t payout-engine .
docker run -d --name payout-engine --restart unless-stopped --env-file .env payout-engine
```

The container opens no ports: the engine has no web interface.

**You should see** a start line like this (with Docker: `docker logs payout-engine`):

```
payout-engine 2026.10.05 (payout-engine-1) on chain 8453: payout Safe 0x0900…e72e, signers from vault secret/trovo/payout-engine#PROCEED_PAYOUT_SIGNERS
```

If it stops instead, the last line says why, for example
`configuration: BASE_RPC_URL is not set`. See
[Troubleshooting](#9-troubleshooting).

## 7. Check it works

1. In Trovo Manager open **Proceeds Payouts → Fee & engine**. The engine
   status shows this instance, a recent heartbeat and no error.
2. Optionally, with Redis, watch its events:
   `redis-cli SUBSCRIBE payout-engine:events`.
3. Run a small test payout on Base Sepolia end to end before the first real
   one: register it, **Prepare**, approve it with two other admins, send
   the payout token to the payout Safe, **Confirm funding**, and watch it
   reach **Completed**.

## 8. Running in production, updating and rolling back

- **How many copies:** run one, or two for failover. Only one works at a
  time: it holds a heartbeat in the database (`payout_engine_states`), and
  a standby takes over when the heartbeat is more than 90 seconds old.
  Give each copy its own [`INSTANCE_ID`](CONFIGURATION.md#instance_id).
- **Stopping is safe at any time** (Ctrl+C, `docker stop`). A payment
  batch already sent is checked on the next start: settled, dropped or
  sent again. Nothing is paid twice.
- **Updating:** pull the new code, rebuild the image, then replace the
  container:
  ```bash
  git pull
  docker build -t payout-engine .
  docker rm -f payout-engine
  docker run -d --name payout-engine --restart unless-stopped --env-file .env payout-engine
  ```
  Deploy the matching app-backend first if the release changed the payout
  tables.
- **Rolling back:** check out the previous release
  (`git checkout <previous-tag>`), rebuild and replace the container the
  same way.
- **Keep the first signer funded** with ETH, and watch the engine status
  in Trovo Manager.

## 9. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `configuration: BASE_RPC_URL is not set` (or another name) | a required parameter is empty | set it in `.env` or the container environment |
| `configuration: PROCEED_PAYOUT_SAFE_ADDRESS is not an address` | a typo in an address | addresses are `0x` plus 40 hexadecimal characters |
| `signers: PROCEED_PAYOUT_SIGNERS must contain at least 3 signers` | fewer than 3 keys, or wrong separator | separate keys with `;` |
| `signers: PROCEED_PAYOUT_SIGNERS_VAULT_PATH is set but VAULT_ADDR / VAULT_TOKEN are not` | Vault path without Vault login | set `VAULT_ADDR` and `VAULT_TOKEN` |
| `database: ...` | wrong connection string, or the database is not reachable | test it with `psql "<the value>"` from the same machine |
| `rpc: ...` or `reading the chain id: ...` | the RPC URL is wrong or the provider is down | open the URL's provider dashboard; try `curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","id":1,"method":"eth_chainId"}' <BASE_RPC_URL>` |
| `GC is not set: holders will not get push notifications` | no Firebase key | set [`GC`](CONFIGURATION.md#gc), or ignore if you do not want notifications |
| Preparing fails with "block range too large" | the RPC limits log queries | lower [`PROCEED_PAYOUT_LOG_CHUNK`](CONFIGURATION.md#proceed_payout_log_chunk) |
| A payout goes back to **Approved** with a shortfall note | the Safe is short of the token, or the first signer is short of ETH | send what the note says, then **Confirm funding** again |
| Admin actions take ~10 seconds to start | Redis is not configured | set the [Redis settings](CONFIGURATION.md#redis-optional) (same Redis as tm-api) |

## Run the tests

`go test ./...` runs the unit tests. The end-to-end test needs a local
blockchain with Safe, the EntryPoint, the paymaster and the market
contracts:

```bash
# terminal 1: a local blockchain (leave it running)
cd ../paymaster/contracts && npx hardhat node

# terminal 2: deploy the contracts to it, then run the tests
cd ../paymaster/contracts && npx hardhat run scripts/local-stack.js --network localhost
cd ../../market/contracts && npx hardhat run scripts/local-market.js --network localhost
cd ../../payout-engine
AA_LOCAL_STACK=../paymaster/contracts/deployments/local-stack.json \
AA_LOCAL_MARKET=../market/contracts/deployments/local-market.json go test ./...
```

- **`AA_LOCAL_STACK`:** the file the local-stack script writes, listing the
  local Safe, EntryPoint and paymaster addresses. Without it the
  end-to-end test is skipped.
- **`AA_LOCAL_MARKET`:** the file the local-market script writes, listing
  the local offer book and token addresses.

The Dockerfile was not built in the environment where it was written (no
Docker there); build it once before relying on it.

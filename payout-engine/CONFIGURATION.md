# Configuration

Every environment variable the engine reads, found with:

```bash
grep -rn "os.Getenv\|env(\"\|envUint(\"\|envAddress(\"" --include="*.go" . | grep -v _test.go
```

`main.go` loads a `.env` file if there is one (`godotenv`), otherwise the
process environment. Copy `.env.example` to `.env` to start. A missing
**required** value stops the engine at startup with the reason.

## Database

### `DB_CONNECTION_STRING`
- **Required**: yes
- **Example**: `host=db.internal user=trovo password=change-me dbname=trovo port=5432 sslmode=require`
- **What it is**: app-backend's primary database. The engine reads and
  writes the payout tables, and reads assets, users, wallets, service fees,
  country configs and offer-book offers.
- **How to get a real value**: use the same value as app-backend's
  `DB_CONNECTION_STRING` (see [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md)).
  app-backend creates the tables (its AutoMigrate), so deploy app-backend
  first.

### `DB_TYPE`
- **Required**: no (defaults to `postgres`)
- **Example**: `postgres`
- **What it is**: `sqlite` makes `DB_CONNECTION_STRING` a SQLite file path
  (tests and local experiments only).

### `DB_MAX_OPEN_CONNECTIONS`
- **Required**: no (defaults to `10`)
- **Example**: `10`
- **What it is**: the size of the engine's Postgres connection pool.

## Chain

### `BASE_RPC_URL`
- **Required**: yes
- **Example**: `https://base-mainnet.g.alchemy.com/v2/your-key`
- **What it is**: the Base JSON-RPC endpoint used for logs, balances, gas
  and sending batches. It must serve `eth_getLogs` over ranges of
  `PROCEED_PAYOUT_LOG_CHUNK` blocks.
- **How to get a real value**: the same provider app-backend uses (its
  `BASE_RPC_URL`), or a dedicated key from your RPC provider's dashboard.

### `BASE_CHAIN_ID`
- **Required**: no (read from the RPC when unset)
- **Example**: `8453` (Base), `84532` (Base Sepolia)

### `PROCEED_PAYOUT_SAFE_ADDRESS`
- **Required**: yes
- **Example**: `0x0900c8af28C7CA2CAd2Fc7941E6d9d5A1b2Ae72e`
- **What it is**: the Safe (v1.4.1) dedicated to paying proceeds. Payouts
  are funded by sending the payout token here, and holders, the fee and the
  VAT are paid from here.
- **How to get a real value**: deploy a Safe whose owners are the
  `PROCEED_PAYOUT_SIGNERS` addresses, with a threshold of at most 3 (e.g.
  with the Safe{Wallet} app on Base). Register `PROCEED_PAYOUT_SIGNERS` in
  TM's vault manager as a managed secret with this Safe as its wallet
  address, so rotating a signer also swaps the Safe owner. Use it for
  nothing else.

### `SAFE_MULTISEND_CALL_ONLY_ADDRESS`
- **Required**: no (defaults to `0x9641d764fc13c8B624c04430C7356C1C7C8102e2`, Safe v1.4.1's canonical MultiSendCallOnly)
- **What it is**: the contract each batch is delegate-called through.

### `OFFER_BOOK_ADDRESS`
- **Required**: no
- **Example**: `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- **What it is**: the `TrovoOfferBook`. When set, tokens escrowed in open
  sell offers are credited to their sellers in the schedule rather than to
  the offer book.
- **How to get a real value**: app-backend's `OFFER_BOOK_ADDRESS` (see
  [market/DEPLOYMENT.md](../market/DEPLOYMENT.md)).

### `PROCEED_PAYOUT_SWEEP_ADDRESS`
- **Required**: no (sweeps are refused without it)
- **Example**: `0x8ba1f109551bD432803012645Ac136ddd64DBA72`
- **What it is**: where an admin's **Sweep** in TM moves the payout Safe's
  whole balance of a token (e.g. excluded holders' shares and rounding dust
  once payouts are done). It is refused while a payout of that token is
  being funded or paid.
- **How to get a real value**: a treasury wallet or Safe the business
  controls.

### `PROCEED_PAYOUT_EXCLUDED_ADDRESSES`
- **Required**: no
- **Example**: `0xabc…,0xdef…`
- **What it is**: addresses never paid, beyond those excluded
  automatically: the payout Safe, the asset's market-making wallet, the
  issuing profile's wallets and the offer book. Separate them with commas,
  semicolons or spaces.

### `TOKENIZATION_ISSUING_PROFILE`
- **Required**: no (defaults to `atprofile`)
- **What it is**: the username whose wallets (issuing and distribution
  Safes) are excluded from payouts. It is the same as app-backend's
  `TOKENIZATION_ISSUING_PROFILE`.

## Signers

The payout Safe is signed by a managed secret: private keys separated by
`;`, at least 3. The first three sign every
transaction, and the first one sends it and pays the gas, so keep it funded
with ETH (the funding check says how much is needed).

### `PROCEED_PAYOUT_SIGNERS`
- **Required**: yes, unless `PROCEED_PAYOUT_SIGNERS_VAULT_PATH` is set
- **Example**: `0x59c6…690d;0x5de4…365a;0x7c85…07a6`
- **What it is**: the managed secret itself, injected into the
  environment at deploy time.
- **How to get a real value**: TM's vault manager (Vault Signer) holds it;
  inject the secret from Vault into the container. Never commit it.

### `PROCEED_PAYOUT_SIGNERS_VAULT_PATH`, `VAULT_ADDR`, `VAULT_TOKEN`, `VAULT_NAMESPACE`
- **Required**: no
- **Example**: `secret/trovo/payout-engine#PROCEED_PAYOUT_SIGNERS`, `https://vault.internal:8200`, `hvs.…`, `admin`
- **What it is**: read the managed secret from Vault (KV v2,
  `mount/path#field`, the field defaulting to `PROCEED_PAYOUT_SIGNERS`) on
  every use instead of the environment. A rotation then takes effect
  without a restart.
- **How to get a real value**: the Vault path the vault manager writes, and
  a Vault token whose policy can read it.

## Tuning

| Variable | Default | Effect |
|---|---|---|
| `PROCEED_PAYOUT_CONFIRMATIONS` | `3` | blocks behind the head the holder snapshot stays (reorg safety) |
| `PROCEED_PAYOUT_LOG_CHUNK` | `5000` | blocks per `eth_getLogs` call |
| `PROCEED_PAYOUT_CATCH_UP_WITHIN` | `2` | the snapshot locks once the scan is this close to the safe head |
| `PROCEED_PAYOUT_BATCH_SIZE` | `150` | transfers per Safe transaction |
| `PROCEED_PAYOUT_MAX_BATCH_GAS` | `12000000` | a batch is split if its estimate is above this |
| `PROCEED_PAYOUT_MIN_UNITS` | `1` | shares below this many base units are skipped |
| `PROCEED_PAYOUT_GAS_PRICE_MULTIPLIER` | `1` | multiplies the suggested fee caps (≥ 1) |
| `PROCEED_PAYOUT_POLL_SECONDS` | `10` | how often it checks the database without a Redis wake-up |
| `PROCEED_PAYOUT_RECEIPT_TIMEOUT_SECONDS` | `180` | how long it waits for a batch to be mined; after that the batch is recovered on the next round (settled, dropped or sent again) |
| `PROCEED_PAYOUT_DEFAULT_START_BLOCK` | `0` | the earliest block a token's first scan starts from (it starts at the asset's creation, an hour early, or this block if later) |

## Fees and VAT

The payout fee wallet and default fee are set in TM (**Proceeds Payouts →
Fee & engine**) and stored in app-backend's `service_fees` row
`PROCEED_PAYOUT_FEE`. The VAT rate is the asset country's
`country_configs.vat_percent`.

### `VAT_WALLET`
- **Required**: only when a payout has a fee and its country charges VAT
- **Example**: `0x2546BcD3c84621e976D8185a91A922aE77ECEc30`
- **What it is**: where the VAT on payout fees is paid. It is the same
  variable, and should be the same wallet, as app-backend's `VAT_WALLET`.
  When unset, the `VAT` row of `service_fees` is used.

## Redis (optional)

Without Redis the engine still works; it notices admin actions on its next
poll.

| Variable | Example | Effect |
|---|---|---|
| `REDIS_HOST` | `redis.internal` | enables the wake-ups, events and heartbeat |
| `REDIS_PORT` | `6379` | |
| `REDIS_PASSWORD` | `change-me` | |
| `REDIS_TLS` | `1` | use TLS |

Use the same Redis as tm-api, which publishes the wake-ups.

## Notifications

### `GC`
- **Required**: no (paid holders are not notified without it)
- **Example**: `ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIs…` (base64)
- **What it is**: the Firebase service account JSON, base64-encoded. The
  engine sends one push notification to each paid Trovo user per payout.
- **How to get a real value**: the same value as app-backend's `GC`.

## Identity

| Variable | Default | Effect |
|---|---|---|
| `INSTANCE_ID` | the hostname | names this instance in the heartbeat and in TM |
| `APP_VERSION` | `dev` | shown in TM's engine state |

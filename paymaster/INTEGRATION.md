# paymaster — integration

How the paymaster connects to the rest of Trovo.

```
app-mobile / app-web ──(sign userOpHash)──┐
                                          ▼
app-backend ──POST /v1/quote──► quote-service ──reads──► Base RPC (paymaster, EntryPoint, Chainlink, DEX pools)
     │                               │
     │                               └──GET──► exchange / FX APIs (http-json sources)
     │
     └──eth_sendUserOperation──► bundler (self-hosted) ──handleOps──► EntryPoint v0.7 ──► TrovoTokenPaymaster
```

| Connects to | Direction | What for |
|---|---|---|
| [app-backend](#1-app-backend) | app-backend → quote service | gets a signed gas quote for each stablecoin-paid operation |
| [Base blockchain](#2-base-blockchain) | quote service reads; EntryPoint calls the paymaster | prices, deposit, paying gas |
| [Price sources](#3-price-sources) | quote service → outside APIs | exchange rates |
| [The bundler](#4-the-bundler) | app-backend → bundler → EntryPoint → paymaster | submits users' operations |
| [Vault](#5-vault) | both parts read | settings and keys |
| [Alerts](#6-alerts) | quote service → chat | low-deposit warnings |
| [The apps](#7-app-mobile-and-app-web) | apps → app-backend | users choose their gas token and sign |

---

## 1. app-backend

- **What it is and why:** Trovo's main backend builds every wallet
  operation. When a user pays gas in a stablecoin, it needs a quote from the
  quote service. It is the only intended caller.
- **Direction:** app-backend → quote service (HTTP).
- **How they connect:** HTTPS (or HTTP inside a private network) to the
  quote service's address, with an `X-API-Key` header.
- **Settings on this side:** [`API_KEYS`](CONFIGURATION.md#api_keys),
  [`HTTP_ADDR`](CONFIGURATION.md#http_addr).
- **Settings on the other side (app-backend):**

  | Parameter | Example | How to get it |
  |---|---|---|
  | `PAYMASTER_QUOTE_SERVICE_URL` | `http://quote-service.internal:8090` | where you run the quote service |
  | `PAYMASTER_QUOTE_SERVICE_API_KEY` | `3f9c1d5e8a7b4c2d9e0f1a2b3c4d5e6f` | one of the quote service's `API_KEYS` |
  | `PAYMASTER_ADDRESS` | `0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9` | `contracts/deployments/<chainId>.json` |

  See [app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md).
- **How to check it works:** a test payment with the gas asset set to a
  stablecoin succeeds, and the quote service's log shows the request.
- **When it is down:** app-backend cannot quote in stablecoins; users with
  ETH pay gas in ETH, others are asked to retry later.

### What app-backend does

When a user sends from a Safe wallet and their gas asset (chosen on their
profile) is an enabled stablecoin the wallet holds, app-backend:

1. builds the UserOperation (sender = the Safe, callData = the Safe4337
   module call, gas limits and fees from the bundler's estimate);
2. for the wallet's **activation** operation (the Safe is not deployed yet)
   includes `approve(paymaster, max)` for the gas token in the batch, and
   passes the Safe factory call as `factory`/`factoryData`;
3. calls `POST /v1/quote` with `X-API-Key: <one of API_KEYS>`:

   ```json
   {
     "token": "cNGN",
     "validitySeconds": 600,
     "userOp": {
       "sender": "0x…", "nonce": "0x5",
       "callData": "0x…",
       "callGasLimit": "0x30d40", "verificationGasLimit": "0x61a80", "preVerificationGas": "0xea60",
       "maxFeePerGas": "0x3b9aca00", "maxPriorityFeePerGas": "0xf4240"
     }
   }
   ```

   Numbers may be `0x` hex or decimal strings. `factory` + `factoryData`
   (or a packed `initCode`) only for the activation operation.
   Shared wallets (operations that wait for approvers) ask for a longer
   `validitySeconds`, up to `QUOTE_MAX_VALIDITY`.

4. copies `paymaster`, `paymasterVerificationGasLimit`,
   `paymasterPostOpGasLimit` and `paymasterData` from the response into the
   UserOperation (v0.7 RPC form), or `paymasterAndData` into the packed
   form — and changes nothing else afterwards: the quote covers every
   field, so any change makes the paymaster reject it (`AA34`);
5. has the user's app sign the userOpHash, submits to the bundler, and
   records the charge (`maxTokenCost` is the most the user can be charged;
   the actual charge is in the paymaster's `GasPaidInToken` event) as a
   `GAS` fee collection.

### The quote response

| Field | Meaning |
|---|---|
| `paymaster`, `paymasterVerificationGasLimit`, `paymasterPostOpGasLimit`, `paymasterData` | Put into the UserOperation. |
| `paymasterAndData` | Same, packed (for the EntryPoint's `PackedUserOperation`). |
| `exchangeRate` | Token base units per 1 ETH (spread included), hex. |
| `marketRate`, `quotedRate`, `spreadBps` | For display/records: tokens per ETH at market and as charged. |
| `maxTokenCost`, `maxTokenCostFormatted` | The most this operation can cost the user in the token — show it before they sign. |
| `validAfter`, `validUntil` | Quote window (unix seconds). |

### Errors

Errors come back as `{"code": "...", "message": "..."}`:

| HTTP | `code` | What app-backend should do |
|---|---|---|
| 400 | `invalid-request`, `invalid-user-operation`, `unknown-token`, `invalid-validity` | A bug in the request — log it. |
| 401 | `unauthorized` | Wrong API key configuration. |
| 409 | `outstanding-debt` | The wallet owes the paymaster from a failed activation charge. Pay this operation's gas in ETH, including a `settleDebt` call, or ask ops to forgive it. |
| 503 | `rate-unavailable`, `rate-out-of-bounds`, `token-disabled`, `chain-unavailable`, `paymaster-deposit-low` | Temporarily unable to quote in this token: fall back to ETH gas if the wallet has ETH, otherwise ask the user to retry later. `rate-out-of-bounds` needs the owner Safe to update `setToken`. |

### Other endpoints

With the same API key: `GET /v1/tokens` (the price of 1 ETH in each gas
token, for fee estimates before building an operation), `GET /v1/rates`
(every pair and each source's answer, for operations staff), `GET /v1/status`
(signer, paymaster, deposit). `GET /healthz` and `GET /readyz` need no key.
The full reference is the service's Swagger page, `/swagger/index.html`.

## 2. Base blockchain

- **What it is and why:** where the paymaster contract, the EntryPoint and
  the price feeds live.
- **Direction:** the quote service only reads (it never sends transactions
  and its key holds no funds). The EntryPoint calls the paymaster to check
  and charge each operation.
- **What the quote service reads:** the paymaster's `tokens()`,
  `quoteSigners()`, `debtTokenCount()` and `entryPoint()`; the EntryPoint's
  `balanceOf(paymaster)` (the deposit); Chainlink `latestRoundData()`;
  DEX pools' `observe()`, `token0()`, `token1()`; tokens' `decimals()`.
- **On-chain behaviour:**

  - `TrovoTokenPaymaster` is called only by the EntryPoint (validation and
    postOp) and by its owner / pauser for administration. Anyone can call
    `deposit()` and `settleDebt()`.
  - Wallets approve the paymaster for their gas token (in the activation
    operation, or when the user switches gas token — the next operation's
    batch includes the approval).
  - Collected tokens stay in the paymaster until the owner withdraws them.
- **Settings:** [`RPC_URL`](CONFIGURATION.md#rpc_url-quote-service),
  [`CHAIN_ID`](CONFIGURATION.md#chain_id-quote-service),
  [`PAYMASTER_ADDRESS`](CONFIGURATION.md#paymaster_address),
  [`ENTRYPOINT_ADDRESS`](CONFIGURATION.md#entrypoint_address-quote-service).
- **When the node is down:** the quote service cannot refresh prices; after
  [`RATE_MAX_AGE`](CONFIGURATION.md#rate_max_age) quotes are refused with
  `chain-unavailable` or `rate-unavailable`.

## 3. Price sources

- **What it is and why:** outside services the quote service reads
  exchange rates from: Chainlink feeds and DEX pools on Base, and JSON APIs
  (exchanges, FX feeds, CoinGecko).
- **Direction:** quote service → each source, every
  [`RATE_REFRESH_INTERVAL`](CONFIGURATION.md#rate_refresh_interval).
- **Settings:** [`RATE_PAIRS`](CONFIGURATION.md#rate_pairs); API keys are
  extra keys of the quote service's Vault secret, referenced as
  `vault:KEY`.
- **How to check it works:** `GET /v1/rates` shows each source's last
  answer and flags outliers.
- **When a source is down:** the pair uses its other sources; below
  `minSources`, quotes for tokens that need it are refused
  (`rate-unavailable`).

## 4. The bundler

- **What it is and why:** an ERC-4337 bundler takes users' operations and
  submits them to the EntryPoint, which calls the paymaster. Trovo runs its
  own.
- **Direction:** app-backend → bundler → EntryPoint → paymaster.
- **Settings on the other side:** app-backend's `BUNDLER_URL`.
- **Bundler configuration:**

  The paymaster is staked, reads and writes token balances during
  validation, and is called by our own (self-hosted) bundler. Configure the
  bundler with EntryPoint v0.7 at `0x0000000071727De22E5E9d8BAf0edAc6f37da032`.
  A wallet's activation operation also deploys the Safe through the Safe
  proxy factory; if the bundler enforces ERC-7562 storage rules strictly for
  unstaked factories, allow-list the Safe proxy factory (or the paymaster)
  in its configuration. Test the full flow on Base Sepolia before mainnet.
- **How to check it works:** a stablecoin-paid operation on Base Sepolia is
  mined and the paymaster emits `GasPaidInToken`.

## 5. Vault

- **What it is and why:** stores both parts' settings and keys.
- **Direction:** the deploy script reads once; the quote service reads at
  start and on each reload.
- **Settings:** [`VAULT_ADDR`](CONFIGURATION.md#vault_addr),
  [`VAULT_TOKEN`](CONFIGURATION.md#vault_token) and the secret paths
  ([CONFIGURATION.md section 1](CONFIGURATION.md#1-bootstrap-environment-variables)).
- **When it is down:** a running quote service keeps its current settings;
  a new one cannot start.

## 6. Alerts

- **What it is and why:** a Discord or Slack webhook that tells treasury
  the paymaster's ETH deposit is running low.
- **Settings:** [`ALERT_WEBHOOK_URL`](CONFIGURATION.md#alert_webhook_url),
  [`DEPOSIT_LOW_WATERMARK_ETH`](CONFIGURATION.md#deposit_low_watermark_eth),
  [`ALERT_COOLDOWN`](CONFIGURATION.md#alert_cooldown).
- **When unset:** the warning only goes to the log.

## 7. app-mobile and app-web

The apps never call the quote service. The user picks their gas asset on
their profile, sees the most the operation can cost
(`maxTokenCostFormatted`) and signs the operation hash, which covers the
quote, like any other send.

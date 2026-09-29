# paymaster — integration

```
app-mobile / app-web ──(sign userOpHash)──┐
                                          ▼
app-backend ──POST /v1/quote──► quote-service ──reads──► Base RPC (paymaster, EntryPoint, Chainlink, DEX pools)
     │                               │
     │                               └──GET──► exchange / FX APIs (http-json sources)
     │
     └──eth_sendUserOperation──► bundler (self-hosted) ──handleOps──► EntryPoint v0.7 ──► TrovoTokenPaymaster
```

## Who calls the quote service

**app-backend** (the only intended caller). When a user sends from a
Safe wallet and their profile's gas asset is an enabled stablecoin that the
wallet holds, app-backend:

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

Response fields app-backend uses:

| Field | Meaning |
|---|---|
| `paymaster`, `paymasterVerificationGasLimit`, `paymasterPostOpGasLimit`, `paymasterData` | Put into the UserOperation. |
| `paymasterAndData` | Same, packed (for the EntryPoint's `PackedUserOperation`). |
| `exchangeRate` | Token base units per 1 ETH (spread included), hex. |
| `marketRate`, `quotedRate`, `spreadBps` | For display/records: tokens per ETH at market and as charged. |
| `maxTokenCost`, `maxTokenCostFormatted` | The most this operation can cost the user in the token — show it before they sign. |
| `validAfter`, `validUntil` | Quote window (unix seconds). |

Errors come back as `{"code": "...", "message": "..."}`:

| HTTP | `code` | What app-backend should do |
|---|---|---|
| 400 | `invalid-request`, `invalid-user-operation`, `unknown-token`, `invalid-validity` | A bug in the request — log it. |
| 401 | `unauthorized` | Wrong API key configuration. |
| 409 | `outstanding-debt` | The wallet owes the paymaster from a failed activation charge. Pay this operation's gas in ETH, including a `settleDebt` call, or ask ops to forgive it. |
| 503 | `rate-unavailable`, `rate-out-of-bounds`, `token-disabled`, `chain-unavailable`, `paymaster-deposit-low` | Temporarily unable to quote in this token: fall back to ETH gas if the wallet has ETH, otherwise ask the user to retry later. `rate-out-of-bounds` needs the owner Safe to update `setToken`. |

Other endpoints (same API key): `GET /v1/tokens` (price of 1 ETH in each
gas token, for fee estimates before building an operation), `GET /v1/rates`
(every pair and each source's answer, for operations), `GET /v1/status`
(signer, paymaster, deposit). `GET /healthz` and `GET /readyz` need no key.

app-backend configuration for this (added with the stablecoin gas
integration): the quote service base URL and API key, and the bundler URL.

## What the quote service calls

- **Base RPC** (`RPC_URL`): paymaster `tokens()`, `quoteSigners()`,
  `debtTokenCount()`, `entryPoint()`; EntryPoint `balanceOf(paymaster)`;
  Chainlink `latestRoundData()`; pool `observe()`/`token0()`/`token1()` and
  token `decimals()`. Read-only; the service never sends transactions and
  its key holds no funds.
- **Rate APIs** configured as `http-json` sources, with API keys from its
  Vault secret.
- **Vault**: reads its own secret at startup and on reload.
- **Alert webhook** (`ALERT_WEBHOOK_URL`) for low-deposit alerts.

## The bundler

The paymaster is staked, reads and writes token balances during
validation, and is called by our own (self-hosted) bundler. Configure the
bundler with EntryPoint v0.7 at `0x0000000071727De22E5E9d8BAf0edAc6f37da032`.
A wallet's activation operation also deploys the Safe through the Safe
proxy factory; if the bundler enforces ERC-7562 storage rules strictly for
unstaked factories, allow-list the Safe proxy factory (or the paymaster)
in its configuration. Test the full flow on Base Sepolia before mainnet.

## On-chain

- `TrovoTokenPaymaster` is called only by the EntryPoint (validation and
  postOp) and by its owner / pauser for administration. Anyone can call
  `deposit()` and `settleDebt()`.
- Wallets approve the paymaster for their gas token (in the activation
  operation, or when the user switches gas token — the next operation's
  batch includes the approval).
- Collected tokens stay in the paymaster until the owner withdraws them.

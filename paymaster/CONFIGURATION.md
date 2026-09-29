# paymaster — configuration

Both parts are configured from **HashiCorp Vault (KV v2)**. The process
environment only says where Vault is and which secret to read; every other
setting — including private keys, API keys and spreads — lives in the
secret.

| Secret (default path) | Read by |
|---|---|
| `secret/trovo/paymaster/deploy` | `contracts/scripts/deploy.js` (once per deployment) |
| `secret/trovo/paymaster/quote-service` | `quote-service` (at startup and on reload) |

All example values below are illustrative — never reuse them.

---

## 1. Bootstrap environment variables (both parts)

These are the only variables read from the process environment.

| Variable | Required | Example | What it does |
|---|---|---|---|
| `VAULT_ADDR` | yes (unless `PAYMASTER_CONFIG_SOURCE=env`) | `https://vault.internal.trovo.io:8200` | Vault server URL. Same variable tm-api's vault signer uses. |
| `VAULT_TOKEN` | yes (unless `PAYMASTER_CONFIG_SOURCE=env`) | `hvs.CAESIJ…` | Token allowed to read the secret (see the policy in §5). In production inject it from your orchestrator (Kubernetes auth, AppRole login in an init step) rather than baking it into an image. |
| `VAULT_KV_MOUNT` | no (default `secret`) | `secret` | Mount path of the KV v2 engine. |
| `VAULT_NAMESPACE` | no | `trovo` | Vault Enterprise / HCP namespace, if you use one. |
| `PAYMASTER_DEPLOY_SECRET_PATH` | no (default `trovo/paymaster/deploy`) | `trovo/paymaster/base-sepolia/deploy` | Deploy script only: path of its secret inside the mount. |
| `QUOTE_SERVICE_SECRET_PATH` | no (default `trovo/paymaster/quote-service`) | `trovo/paymaster/base-mainnet/quote-service` | Quote service only: path of its secret inside the mount. |
| `PAYMASTER_CONFIG_SOURCE` | no (default `vault`) | `env` | `env` reads every key below from the process environment instead of Vault. **Local development only** — never in production. |

---

## 2. Deploy script secret (`trovo/paymaster/deploy`)

Used by `node scripts/deploy.js` (see [DEPLOYMENT.md](DEPLOYMENT.md)).

| Key | Required | Example | What it does / how to get it |
|---|---|---|---|
| `RPC_URL` | yes | `https://base-mainnet.g.alchemy.com/v2/<key>` | JSON-RPC endpoint of the target chain. From your node provider's dashboard (Alchemy, QuickNode, Coinbase Developer Platform) or your own node. |
| `CHAIN_ID` | yes | `8453` | Expected chain id; the script refuses to run if `RPC_URL` is a different chain. Base mainnet `8453`, Base Sepolia `84532`. |
| `DEPLOYER_PRIVATE_KEY` | yes | `0x4c0883a6…` (64 hex) | Key that deploys the contract and pays for deployment, stake and deposit. It is only the owner until the end of the script, when ownership moves to `OWNER_ADDRESS`. Generate with `cast wallet new` (Foundry) and fund it with `STAKE_AMOUNT_ETH + INITIAL_DEPOSIT_ETH` plus ~0.01 ETH for gas. Empty it afterwards. |
| `ENTRYPOINT_ADDRESS` | no (default `0x0000000071727De22E5E9d8BAf0edAc6f37da032`) | `0x0000000071727De22E5E9d8BAf0edAc6f37da032` | ERC-4337 EntryPoint v0.7. The default is the canonical deployment (same on Base and Base Sepolia); the script checks it is really a v0.7 EntryPoint. |
| `OWNER_ADDRESS` | yes | `0x9A3c…` (treasury Safe) | Final owner: can change gas tokens and rate bounds, rotate quote signers, pause/unpause, withdraw collected tokens and ETH. Must be a Safe (the script refuses an address without code unless `ALLOW_EOA_OWNER=true`). |
| `ALLOW_EOA_OWNER` | no (default `false`) | `true` | Allows a plain key as owner. Testnets only. |
| `QUOTE_SIGNER_ADDRESSES` | yes | `0x7099…79C8,0x3C44…93BC` | Comma-separated addresses whose quote signatures the paymaster accepts: the address of the quote service's `QUOTE_SIGNER_PRIVATE_KEY` (and, during a rotation, the next one). |
| `PAUSER_ADDRESS` | yes | `0x90F7…b906` | Address that may pause the paymaster in an emergency (only the owner can unpause). An on-call operations key or a 1-of-N Safe. |
| `POST_OP_GAS_OVERHEAD` | no (default `45000`) | `45000` | Gas added to each operation's charge to cover the paymaster's own post-operation work (which the EntryPoint does not include in the cost it reports). Owner can change it later with `setPostOpGasOverhead`. |
| `STAKE_AMOUNT_ETH` | no (default `0.1`) | `0.5` | ETH staked at the EntryPoint. Bundlers require paymasters that touch token storage to be staked; check your bundler's minimum. The stake is not spent. |
| `UNSTAKE_DELAY_SEC` | no (default `86400`) | `86400` | Delay between `unlockStake` and `withdrawStake`. Bundlers require at least 1 day. |
| `INITIAL_DEPOSIT_ETH` | no (default `0.05`) | `1` | ETH deposited at the EntryPoint to pay gas. This is what gets spent; top it up from converted stablecoins. |
| `GAS_TOKENS` | yes | see below | JSON array of the tokens to enable and their rate bounds. |

`GAS_TOKENS` example:

```json
[
  {"symbol": "USDC", "address": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "minRate": "500000000", "maxRate": "20000000000"},
  {"symbol": "cNGN", "address": "<cNGN token address on Base>", "minRate": "500000000000", "maxRate": "40000000000000"}
]
```

- `address`: the token contract (USDC above is Base mainnet's native USDC;
  take cNGN's and USDT's from the issuer's official documentation).
- `minRate` / `maxRate`: the range of exchange rates the paymaster will
  accept, in **token base units per 1 ETH** (spread included). USDC has 6
  decimals, so `500000000` = 500 USDC per ETH and `20000000000` = 20,000
  USDC per ETH. They are a safety net against a compromised quote signer,
  not a price feed: keep them wide enough for normal volatility (e.g. ¼× to
  4× today's price) and review them when the market moves a lot. The quote
  service refuses to sign rates outside them and reports
  `rate-out-of-bounds` so you know to update them (`setToken` from the
  owner Safe).

---

## 3. Quote service secret (`trovo/paymaster/quote-service`)

### Chain and signing

| Key | Required | Example | What it does / how to get it |
|---|---|---|---|
| `RPC_URL` | yes | `https://base-mainnet.g.alchemy.com/v2/<key>` | JSON-RPC endpoint for reading the paymaster, the EntryPoint deposit, Chainlink feeds and DEX pools. |
| `CHAIN_ID` | yes | `8453` | Must match `RPC_URL`'s chain; part of every quote hash. |
| `ENTRYPOINT_ADDRESS` | no (default `0x0000000071727De22E5E9d8BAf0edAc6f37da032`) | `0x0000000071727De22E5E9d8BAf0edAc6f37da032` | Must be the EntryPoint the paymaster was deployed with (checked at startup). |
| `PAYMASTER_ADDRESS` | yes | `0xDc64…F6C9` | The deployed `TrovoTokenPaymaster` (printed by the deploy script and saved in `contracts/deployments/<chainId>.json`). |
| `QUOTE_SIGNER_PRIVATE_KEY` | yes | `0x59c6…690d` (64 hex) | Key that signs quotes. Its address must be in the paymaster's quote signers (startup fails otherwise). Generate with `cast wallet new` or `openssl rand -hex 32`. It needs no ETH. Keep it only in Vault. |
| `API_KEYS` | yes | `3f9c…e1,8a02…7d` | Comma-separated keys accepted in the `X-API-Key` header (≥ 16 characters each). Give app-backend its own key; several allow rotation. Generate with `openssl rand -hex 32`. |

### Prices and spread

| Key | Required | Example | What it does |
|---|---|---|---|
| `GAS_TOKENS` | yes | see below | JSON array of gas tokens, how to price them, and their spread. |
| `RATE_PAIRS` | yes | see below | JSON object of rate pairs and their pluggable sources. |
| `DEFAULT_SPREAD_BPS` | no (default `0`) | `100` | Spread for tokens without their own `spreadBps`, in basis points (100 = 1%). Users pay the market price of gas plus this. Max 5000. |
| `RATE_REFRESH_INTERVAL` | no (default `30s`) | `30s` | How often every source is re-read. |
| `RATE_MAX_AGE` | no (default `5m`) | `5m` | A pair whose last successful aggregation is older than this is unusable, and quotes for tokens that need it are refused (`rate-unavailable`). |
| `RATE_SOURCE_TIMEOUT` | no (default `10s`) | `5s` | Per-source request timeout. |

`GAS_TOKENS` example:

```json
[
  {"symbol": "USDC", "address": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "decimals": 6, "route": ["ETH/USD", "USD/USDC"], "spreadBps": 100},
  {"symbol": "USDT", "address": "<USDT on Base>", "decimals": 6, "route": ["ETH/USD", "USD/USDT"], "spreadBps": 100},
  {"symbol": "cNGN", "address": "<cNGN on Base>", "decimals": 6, "route": ["ETH/USD", "USD/NGN", "NGN/CNGN"], "spreadBps": 150}
]
```

| Field | Meaning |
|---|---|
| `symbol` | Name used in API requests (`"token": "cNGN"`) and logs. |
| `address` | Token contract; must be enabled on the paymaster (`setToken`). |
| `decimals` | Token decimals (read it with `cast call <token> "decimals()(uint8)"`). |
| `route` | How to convert 1 ETH into this token: a chain of pairs from `RATE_PAIRS`, starting at `ETH`. A leg may be the inverse of a configured pair (`USD/USDC` uses a configured `USDC/USD`). |
| `spreadBps` | This token's spread (overrides `DEFAULT_SPREAD_BPS`). The paymaster rate = market rate × (1 + spreadBps/10000). |

### Rate pairs and sources

`RATE_PAIRS` maps `"BASE/QUOTE"` (1 BASE = x QUOTE) to its sources:

```json
{
  "ETH/USD": {
    "minSources": 2,
    "maxDeviationBps": 200,
    "sources": [
      {"name": "chainlink", "type": "chainlink", "feed": "<Chainlink ETH/USD feed on Base>", "maxAge": "30m"},
      {"name": "coingecko", "type": "http-json", "url": "https://pro-api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd", "path": "ethereum.usd", "headers": {"x-cg-pro-api-key": "vault:COINGECKO_API_KEY"}}
    ]
  },
  "USD/NGN": {
    "minSources": 1,
    "maxDeviationBps": 300,
    "sources": [
      {"name": "cngn-usdc-pool", "type": "uniswap-v3-twap", "pool": "<cNGN/USDC pool on Base>", "baseToken": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "window": "30m"},
      {"name": "exchange", "type": "http-json", "url": "<exchange ticker URL for USDT/NGN>", "path": "<path to the last price>", "headers": {"Authorization": "vault:EXCHANGE_API_KEY"}},
      {"name": "fx", "type": "http-json", "url": "https://open.er-api.com/v6/latest/USD", "path": "rates.NGN"}
    ]
  },
  "NGN/CNGN": {"sources": [{"name": "peg", "type": "fixed", "value": "1"}]},
  "USDC/USD": {"sources": [{"name": "peg", "type": "fixed", "value": "1"}]},
  "USDT/USD": {"sources": [{"name": "peg", "type": "fixed", "value": "1"}]}
}
```

Pair options:

| Field | Default | Meaning |
|---|---|---|
| `sources` | — | One or more source definitions (below). |
| `minSources` | `1` | How many sources must answer (within the deviation limit) for the pair to be usable. Use ≥ 2 for anything market-driven once you have two good sources. |
| `maxDeviationBps` | `500` | Sources further than this from the median are ignored as outliers (and flagged in `GET /v1/rates`). Negative disables the check. |

The pair's rate is the median of the remaining sources.

Source types — every source has a `name` (unique within its pair) and a
`type`, and may add `"invert": true` (use 1/x, e.g. a feed quoted the other
way round) and `"multiplier": "<decimal>"` (scale the answer):

| `type` | Fields | Use it for |
|---|---|---|
| `fixed` | `value` | A peg (`NGN/CNGN` = 1, `USDC/USD` = 1) or a temporary manual override. |
| `http-json` | `url`, `path` (dot path, numeric segments index arrays, e.g. `data.0.price`), `headers` | Any JSON API: exchange tickers, FX feeds, the stablecoin issuer's own rate endpoint, CoinGecko. Header values `vault:KEY` are read from key `KEY` of this same secret, so API keys stay in Vault. |
| `chainlink` | `feed`, `maxAge` (default `1h`) | A Chainlink price feed. Find feed addresses at [data.chain.link](https://data.chain.link) (select Base). Set `maxAge` a little above the feed's heartbeat. |
| `uniswap-v3-twap` | `pool`, `baseToken`, `window` (default `30m`, 1m–168h) | The time-weighted average price of `baseToken` in the pool's other token, from any Uniswap v3-compatible pool (Uniswap v3, Aerodrome Slipstream, PancakeSwap v3). This is the **cNGN market-rate discovery**: with `baseToken` = USDC on a cNGN/USDC pool you get cNGN per USDC, i.e. the market USD/NGN rate at which gas actually converts. Pick a pool with real liquidity; a TWAP over a longer window is harder to manipulate. The pool must have enough observation cardinality for the window (the source reports an error otherwise). |

Adding a new source type is one Go file in `quote-service/internal/rates`
calling `Register("my-type", factory)`.

Secret values referenced by sources (e.g. `COINGECKO_API_KEY`,
`EXCHANGE_API_KEY` above) are extra keys in the same secret.

### Quotes

| Key | Required | Example | What it does |
|---|---|---|---|
| `QUOTE_VALIDITY` | no (default `10m`) | `10m` | Validity of a quote when the caller does not ask for one. Enough for a single-signer wallet to sign and submit. |
| `QUOTE_MAX_VALIDITY` | no (default `24h`, max `168h`) | `24h` | Longest validity a caller may request — for shared wallets, whose operations wait for approvers. The rate is fixed for the whole window, so the spread must cover price movement over it. |
| `QUOTE_CLOCK_SKEW` | no (default `60s`) | `60s` | Quotes become valid this long before "now", tolerating clock differences with the chain. |
| `PAYMASTER_VERIFICATION_GAS_LIMIT` | no (default `150000`) | `150000` | Gas limit for the paymaster's validation, set in every quote (part of the signed data). |
| `PAYMASTER_POST_OP_GAS_LIMIT` | no (default `80000`) | `80000` | Gas limit for the paymaster's post-operation step. Must cover a token transfer; must be ≥ the contract's `postOpGasOverhead` in practice. |

### Deposit monitoring and operations

| Key | Required | Example | What it does |
|---|---|---|---|
| `DEPOSIT_CHECK_INTERVAL` | no (default `1m`) | `1m` | How often the paymaster's EntryPoint deposit is read. |
| `DEPOSIT_LOW_WATERMARK_ETH` | no (default `0.05`) | `0.5` | Alert when the deposit falls below this. Size it for a few hours of traffic. Quotes are refused (`paymaster-deposit-low`) only when the deposit cannot cover the specific operation. |
| `ALERT_WEBHOOK_URL` | no | `https://discord.com/api/webhooks/…` | Incoming webhook (Discord or Slack format — the body carries both `content` and `text`) for low-deposit alerts. Empty = log only. |
| `ALERT_COOLDOWN` | no (default `1h`) | `1h` | Minimum time between two low-deposit alerts. |
| `HTTP_ADDR` | no (default `:8090`) | `:8090` | Listen address. Changing it needs a restart. |
| `CONFIG_RELOAD_INTERVAL` | no (default `0s` = off) | `5m` | Re-read the secret periodically and apply it (tokens, spreads, sources, keys) without a restart. Sending `SIGHUP` to the process also reloads. An invalid new configuration is logged and the running one is kept. |

---

## 4. Local development without Vault

```bash
cd paymaster/quote-service
export PAYMASTER_CONFIG_SOURCE=env RPC_URL=http://127.0.0.1:8545 CHAIN_ID=31337 \
  PAYMASTER_ADDRESS=0x… QUOTE_SIGNER_PRIVATE_KEY=0x… API_KEYS=local-dev-key-0123456789 \
  GAS_TOKENS='[{"symbol":"USDC","address":"0x…","decimals":6,"route":["ETH/USD"]}]' \
  RATE_PAIRS='{"ETH/USD":{"sources":[{"type":"fixed","value":"3000"}]}}'
go run .
```

## 5. Writing the secrets and the Vault policy

```bash
# quote service (values as JSON so the nested GAS_TOKENS / RATE_PAIRS stay intact)
vault kv put -mount=secret trovo/paymaster/quote-service @quote-service.json

# deploy script
vault kv put -mount=secret trovo/paymaster/deploy @deploy.json
```

Policy for the quote service's token (read-only, its own path only):

```hcl
path "secret/data/trovo/paymaster/quote-service" {
  capabilities = ["read"]
}
```

The deploy secret holds a funded private key: give it its own policy, only
to whoever runs deployments, and delete or rotate `DEPLOYER_PRIVATE_KEY`
after the deployment.

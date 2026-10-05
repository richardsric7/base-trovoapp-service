# paymaster

## What this project does

Every action on the blockchain (sending money, buying a token) costs a
small network fee, called "gas", normally paid in ETH. Most Trovo users
hold stablecoins such as cNGN or USDC, not ETH. This project lets them pay
that fee **in the stablecoin they already have**: Trovo pays the network
in ETH and charges the user the equivalent in their stablecoin, at a fair
current price plus a small spread.

It has two parts: a blockchain contract (the "paymaster") that pays the
fees, and a small always-on web service (the "quote service") that works
out the price and signs it. app-backend asks the quote service for a price
whenever a user pays gas in a stablecoin. Without it, users must hold ETH
to do anything.

- How to deploy: [DEPLOYMENT.md](DEPLOYMENT.md)
- Every setting: [CONFIGURATION.md](CONFIGURATION.md)
- How other projects use it: [INTEGRATION.md](INTEGRATION.md)

## The two parts

| Part | What it is |
|---|---|
| [`contracts/`](contracts) | `TrovoTokenPaymaster`, an ERC-4337 (EntryPoint v0.7) paymaster. It pays the EntryPoint in ETH from its own deposit and charges the wallet the equivalent in the chosen token, at an exchange rate signed by the quote service. Hardhat project with tests and a Vault-configured deploy script. |
| [`quote-service/`](quote-service) | Go HTTP service. Discovers exchange rates from pluggable sources (Chainlink, DEX pool TWAP, JSON APIs, fixed pegs), adds Trovo's spread, and signs quotes for individual UserOperations. Also watches the paymaster's ETH deposit and alerts when it runs low. Configured from Vault. |

## How a gas payment works

1. app-backend builds the wallet's UserOperation (a Safe transaction) and
   asks the quote service for a quote in the gas token the user picked on
   their profile (`POST /v1/quote`).
2. The quote service prices 1 ETH in that token (e.g. cNGN: ETH/USD ×
   USD/NGN × NGN/cNGN, from several sources, median, outliers dropped),
   adds the token's spread, and signs a quote bound to that exact operation.
3. app-backend puts the returned `paymasterAndData` in the operation; the
   user's app signs the operation hash (which covers the quote) with the
   wallet key, as for any other send.
4. The bundler submits it. The paymaster checks the quote signature, the
   token and the rate bounds, and:
   - for a deployed wallet that already approved the paymaster, pulls the
     maximum possible cost up front and refunds the unused part afterwards;
   - for the wallet's first (activation) operation — whose calldata
     approves the paymaster — charges the actual cost after the calldata
     ran.
5. Collected stablecoins accumulate in the paymaster; the treasury
   withdraws them, converts them to ETH and tops the deposit up. The quote
   service's deposit monitor alerts when the deposit falls below the
   watermark.

Wallets that hold no enabled stablecoin pay gas in ETH as usual, without
the paymaster.

### Why the spread matters

The paymaster pays Base in ETH but is paid in stablecoins, so every quote
is a small currency exchange. The spread (per token, in basis points)
covers price movement between quoting and conversion, conversion costs,
and margin. For cNGN the market rate is discovered from its on-chain
market (the cNGN/USDC pool TWAP) and/or exchange and FX APIs, so users get
the actual cNGN price of gas plus Trovo's spread — see
[CONFIGURATION.md](CONFIGURATION.md#rate_pairs).

## Safety properties

- A quote only pays for the exact operation it was signed for (every field
  but the wallet signature is hashed), within its validity window.
- The owner (the treasury Safe) sets per-token minimum/maximum rates; a
  compromised quote signer cannot sign rates outside them.
- Quote signers can be rotated (several may be active at once), a separate
  pauser can stop the paymaster immediately, and only the owner can unpause,
  change configuration or withdraw funds.
- A failed charge never reverts the user's action; it is recorded as debt,
  and the wallet cannot use the paymaster again until it is settled.
- The contract is not upgradeable. **It has not been externally audited
  yet — get it audited before mainnet.**

## Layout

```
paymaster/
├── contracts/
│   ├── src/TrovoTokenPaymaster.sol    # the paymaster
│   ├── src/test/                      # test-only contracts (mock token, EntryPoint imports)
│   ├── test/                          # Hardhat tests against the official EntryPoint v0.7
│   ├── scripts/deploy.js              # Vault-configured deployment
│   ├── scripts/fixtures.js            # regenerates the Go hash fixture
│   └── hardhat.config.js
└── quote-service/
    ├── main.go
    ├── internal/config/               # bootstrap env + Vault loading + validation
    ├── internal/rates/                # pluggable rate sources and aggregation
    ├── internal/quote/                # quote hash (mirrors getHash), pricing, signing
    ├── internal/chain/                # paymaster / EntryPoint reads
    ├── internal/monitor/              # deposit monitor + alerts
    ├── internal/api/                  # HTTP API (Swagger at /swagger/index.html)
    └── Dockerfile
```

## Running the tests

```bash
cd contracts && npm ci && npm test          # 20 Hardhat tests
cd quote-service && go test ./...           # includes a check against the contract's own getHash
```

If you change the quote hash in the contract, regenerate the fixture
(`cd contracts && npx hardhat run scripts/fixtures.js`) and re-run the Go
tests.

## More

- [DEPLOYMENT.md](DEPLOYMENT.md) — deploying the contract and running the service
- [CONFIGURATION.md](CONFIGURATION.md) — every setting, with examples
- [INTEGRATION.md](INTEGRATION.md) — how app-backend and the bundler use it

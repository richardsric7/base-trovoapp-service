# Trovo Wallet platform

This repository holds every part of the **Trovo Wallet** platform: the
mobile and web wallets people use, the servers behind them, the admin
dashboard staff use, and the smart contracts on the **Base** blockchain.

Trovo Wallet lets people hold and send crypto and stablecoins, trade with
each other (P2P), buy tokenized real-world assets and tokenized Nigerian
stocks and bonds (Public Markets), and receive dividends and interest on
them. Each user's wallet is a Safe smart account on Base.

New here? Read this page, then [ARCHITECTURE.md](ARCHITECTURE.md) for how
the parts talk to each other, then the README of the project you work on.
Every project has the same four documents:

| Document | What it tells you |
|---|---|
| `README.md` | what the project does, and a quick start |
| `DEPLOYMENT.md` | step by step, how to get it running |
| `CONFIGURATION.md` | every setting: what it is, why it's needed, an example and how to get it |
| `INTEGRATION.md` | everything it connects to, and how |

How to write them: [CONTRIBUTING.md](CONTRIBUTING.md).

## The projects

| Project | What it does | Built with |
|---|---|---|
| [`app-backend`](app-backend/README.md) | The main server (the Trovo Wallet API). Every action in the apps goes through it: accounts, wallets, payments, swaps, P2P trading, tokenized assets, Public Markets, KYC, account recovery, and the API partner businesses use. | Go |
| [`app-mobile`](app-mobile/README.md) | The Trovo mobile app for Android and iOS. | Flutter (Dart) |
| [`app-web`](app-web/README.md) | The Trovo web wallet: the same features in a browser. | React (Vite, TypeScript) |
| [`tm-api`](tm-api/README.md) | The server behind Trovo Manager, the admin dashboard: manages users, assets, fees, payouts, Public Markets, partners and platform keys. | Go |
| [`tm-web`](tm-web/README.md) | Trovo Manager's website, used by staff. | Next.js (TypeScript) |
| [`payment-history-engine`](payment-history-engine/README.md) | Watches the blockchain for transfers to and from Trovo wallets and records them, so users see their payment history. | Go |
| [`payout-engine`](payout-engine/README.md) | Pays dividends and interest on tokenized assets to their holders, after staff approve each payout. | Go |
| [`paymaster`](paymaster/README.md) | Lets users pay blockchain fees ("gas") in a stablecoin instead of ETH: a contract, plus a quote service that prices each fee. | Solidity, Go |
| [`market`](market/README.md) | The offer book contract where tokenized assets are sold and swapped, and the token contract each tokenized asset uses. | Solidity |
| [`recovery`](recovery/README.md) | The account recovery contract: if a user loses their phone, the platform can give their wallet a new key after a waiting period, but can never move their money. | Solidity |
| [`wallet-core`](wallet-core/README.md) | Shared code for keys, signing and wallet addresses, built into app-web (WebAssembly) and app-mobile (a native library). Not deployed on its own. | Rust |

## Deploy order

Each project needs the ones above it. Follow each one's `DEPLOYMENT.md`.

1. **Databases and Redis**: PostgreSQL (main database, plus a tracking
   database for payment history) and Redis.
2. **Contracts** on Base:
   [`paymaster`](paymaster/DEPLOYMENT.md) (and its quote service and a
   bundler), then, if you use those features,
   [`market`](market/DEPLOYMENT.md) and [`recovery`](recovery/DEPLOYMENT.md).
3. [`app-backend`](app-backend/DEPLOYMENT.md): it creates the main
   database's tables.
4. [`payment-history-engine`](payment-history-engine/DEPLOYMENT.md) and
   [`payout-engine`](payout-engine/DEPLOYMENT.md): they work in tables
   app-backend created.
5. [`tm-api`](tm-api/DEPLOYMENT.md), then [`tm-web`](tm-web/DEPLOYMENT.md).
6. [`app-web`](app-web/DEPLOYMENT.md) and [`app-mobile`](app-mobile/DEPLOYMENT.md),
   pointed at app-backend's address.

`wallet-core` is built into the apps; rebuild it only when you change it
([wallet-core/DEPLOYMENT.md](wallet-core/DEPLOYMENT.md)).

## Getting the code

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service
```

There is no CI/CD pipeline in this repository yet. Each project's
`DEPLOYMENT.md` lists its build and test commands; run them before merging.

# Trovo Wallet Web App (`app-web`)

## What this project does

`app-web` is the **Trovo web wallet**, the website version of the Trovo
app. In a browser, people can:

- create a wallet (or import one with its 12-word recovery phrase) and see
  their balances;
- send and receive money, and swap one asset for another;
- trade with other people on the **P2P** marketplace;
- buy **tokenized assets** and **Public Markets** (tokenized Nigerian
  stocks and bonds), and see dividends;
- deposit and withdraw Naira through a bank, and turn on account recovery.

It is only a website: it has no server of its own. All data comes from
[`app-backend`](../app-backend/README.md). Keys are made and used only in
the browser, with [`wallet-core`](../wallet-core/README.md) compiled to
WebAssembly; app-backend never sees them.

Next: [DEPLOYMENT.md](DEPLOYMENT.md) to run it,
[CONFIGURATION.md](CONFIGURATION.md) for its settings and
[INTEGRATION.md](INTEGRATION.md) for what it connects to.

## Tech stack

| Concern | Library / tool |
|---|---|
| Framework | React 19 |
| Language | TypeScript 7 |
| Build tool / dev server | Vite 8 (`@vitejs/plugin-react`) |
| Routing | React Router 7 (`react-router-dom`) |
| State and data | Redux Toolkit and RTK Query, with an Axios-based `baseQuery` |
| Styling | Tailwind CSS 4 (through `@tailwindcss/postcss`) and `tw-elements-react` |
| Keys and signing | `wallet-core`, compiled to WebAssembly, copied into `src/walletCore/` |
| PDF receipts | `@react-pdf/renderer` |
| Card payments | `flutterwave-react-v3` |
| Production server | nginx (Docker image) |

## Directory structure (`src/`)

```
src/
├── pages/            Route-level screens, one folder per flow:
│                        onboarding/, createAccount/, importWallet/,
│                        accountRecovery/, dashboard/, pay/, pdfPages/
├── components/        Shared/reusable UI components (buttons, modals,
│                        tables, wallet cards, dropdowns, etc.), plus a
│                        components/p2p/ subfolder for P2P-marketplace-
│                        specific components
├── routingSetup/      React Router setup: appRouter.tsx (route table) and
│                        routeGuard.tsx (ProtectedRoutes — redirects to
│                        /login or /welcome when the user isn't logged in)
├── store/             Redux state + data layer
│   ├── api/             RTK Query "API slices" per backend domain:
│   │                      authApi.ts, walletApis.ts, p2pApis.ts,
│   │                      tokenizationApis.ts, sharedAccessApis.ts,
│   │                      bankApis.ts, publicMarketsApis.ts,
│   │                      cacheApi.ts, and baseapi/ (the shared
│   │                      axios-backed baseQuery + auth/signing headers)
│   ├── reduxStore/       Store configuration (configureStore, preloaded
│   │                      state)
│   ├── config.ts         BASE_URL / SOCKET_URL (reads VITE_API_URL /
│   │                      VITE_SOCKET_URL)
│   ├── constants.ts       localStorage keys and misc constants
│   └── *Slice.ts          Redux slices (auth, sidebar, cache, app state)
├── walletCore/        Vendored wasm-bindgen output from the wallet-core
│                        Rust crate (wallet_core.js, wallet_core_bg.wasm,
│                        and .d.ts types) — do not hand-edit, see
│                        INTEGRATION.md
├── utils/              Helpers, incl. trovoSDK.ts (thin wrapper around
│                        walletCore for key generation / signing) and
│                        storage.ts (localStorage helpers)
├── hooks/              Custom React hooks
├── types/              Shared TypeScript types
├── assets/, fonts/     Static images/fonts imported by components
```

Other notable root files:
- `vite.config.mts` — Vite config (dev server on port 3000, React plugin,
  tsconfig path resolution). No Node globals polyfill is needed: the app's
  code and dependencies run in the browser without one.
- `postcss.config.js` — Tailwind 4 through PostCSS (there is no
  `tailwind.config.js`; Tailwind 4 reads its settings from the CSS).
- `Dockerfile` / `nginx.conf.template` — production container build; see
  [`DEPLOYMENT.md`](./DEPLOYMENT.md).

## Running it locally

The short version (the full guide is [DEPLOYMENT.md](DEPLOYMENT.md)), with
Node.js 20.19 or newer and app-backend running at `http://localhost:8080`:

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/app-web
npm ci
cp .env.example .env.local   # points the website at app-backend
npm start                    # http://localhost:3000
```

To check and build the production files: `npm run build` (output in
`dist/`), then `npm run preview` to look at them.

## Further documentation

- [`DEPLOYMENT.md`](./DEPLOYMENT.md) — how to build and deploy this app,
  including how the Docker/nginx setup injects runtime config.
- [`CONFIGURATION.md`](./CONFIGURATION.md) — every environment variable
  this app reads, what it does, and how to get a real value.
- [`INTEGRATION.md`](./INTEGRATION.md) — how this app talks to
  `app-backend` (signing), uses the `wallet-core` WebAssembly package,
  and its other connections.

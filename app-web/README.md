# Trovo Wallet Web App (`app-web`)

`app-web` is the end-user web client for **Trovo Wallet** — the customer-facing
app people use to manage their wallet, send and receive payments, swap
assets, trade on the P2P marketplace, hold tokenized assets, and buy and
sell tokenized Nigerian stocks and bonds (Public Markets).

It is a **pure frontend**: it has no server of its own and exposes no API.
Every piece of data it shows comes from HTTP (and WebSocket) calls to the
[`app-backend`](../app-backend) Go API. Some cryptographic operations
(key generation, mnemonic handling, request signing) are done client-side
using a Rust crate from [`wallet-core`](../wallet-core), compiled to
WebAssembly and vendored into this project — see
[`INTEGRATION.md`](./INTEGRATION.md) for details.

This app is unrelated to the internal admin tooling in the monorepo
(`tm-api` / `tm-web`) — see [`INTEGRATION.md`](./INTEGRATION.md#c-no-relationship-to-tm-api--tm-web)
for that distinction.

## Tech stack

| Concern | Library / tool |
|---|---|
| Framework | React 18 (`react` / `react-dom` ^18.2) |
| Language | TypeScript ^5.9 |
| Build tool / dev server | Vite ^7.3 (`@vitejs/plugin-react`) |
| Routing | React Router (`react-router-dom` ^6.21) |
| State management | Redux Toolkit (`@reduxjs/toolkit`, `react-redux`) |
| Data fetching / server cache | RTK Query (`createApi`, via `@reduxjs/toolkit/query/react`), with a **custom Axios-based `baseQuery`** (not the default `fetchBaseQuery`) |
| HTTP client | Axios |
| Styling | Tailwind CSS (`tailwindcss`, `postcss`, `autoprefixer`) + `tw-elements-react` |
| Client-side crypto | `wallet_core` — a Rust crate from the sibling `wallet-core` project, compiled to WASM and vendored under `src/walletCore/` |
| PDF generation | `@react-pdf/renderer` (e.g. exporting receipts/statements) |
| Payments (fiat on-ramp) | `flutterwave-react-v3` (Flutterwave) |
| Testing | `@testing-library/react`, `@testing-library/jest-dom` |

The project was originally bootstrapped with Create React App and has since
been migrated to Vite; some CRA-era artifacts (e.g. `public/manifest.json`,
`browserslist` in `package.json`) are still present.

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
- `tailwind.config.js` / `postcss.config.js` — Tailwind/PostCSS setup.
- `Dockerfile` / `nginx.conf.template` — production container build; see
  [`DEPLOYMENT.md`](./DEPLOYMENT.md).

## Running it locally

Run these from inside the `app-web/` directory.

```bash
# 1. Use Node 20 (matches CI; see DEPLOYMENT.md for details)
nvm use 20   # or install Node 20.x another way

# 2. Install dependencies (exact versions, from package-lock.json)
npm ci

# 3. Point the app at a backend (optional — see below for the default)
cp .env.example .env.local
# edit .env.local if you want to hit a different app-backend instance

# 4. Start the dev server
npm start

# 5. Open the app
# Vite opens it automatically at http://localhost:3000
```

By default (no `.env.local`), the app talks to same-origin `/api` and
`/ws/v1` — that only works behind the nginx container described in
`DEPLOYMENT.md`. For local development against a real backend, set
`VITE_API_URL` / `VITE_SOCKET_URL` in `.env.local` — see
[`CONFIGURATION.md`](./CONFIGURATION.md) for the full list of environment
variables and how to get real values for them.

To type-check and build a production bundle locally:

```bash
npm run build   # runs `tsc && vite build`, output goes to dist/
npm run preview # serve the built dist/ locally, for a quick smoke test
```

## Further documentation

- [`DEPLOYMENT.md`](./DEPLOYMENT.md) — how to build and deploy this app,
  including how the Docker/nginx setup injects runtime config.
- [`CONFIGURATION.md`](./CONFIGURATION.md) — every environment variable
  this app reads, what it does, and how to get a real value.
- [`INTEGRATION.md`](./INTEGRATION.md) — how this app talks to
  `app-backend` (auth/signing) and to the `wallet-core` WASM package, and
  its (non-)relationship to `tm-api` / `tm-web`.

# Integration

Everything `app-web` connects to, and how. app-web is a website with no
server code of its own: it talks to app-backend, and uses wallet-core (built
into it) to make keys and sign in the browser.

| Connects to | Direction | Through | Needed? |
|---|---|---|---|
| [1. app-backend](#1-app-backend) | website → app-backend | HTTPS (REST) | yes |
| [2. wallet-core](#2-wallet-core) | built into the website | WebAssembly | yes |
| [3. Flutterwave](#3-flutterwave) | website → Flutterwave pop-up | Flutterwave's script | optional |
| [4. Block explorer](#4-block-explorer) | links only | basescan.org | — |
| [5. Not connected: tm-api and tm-web](#5-not-connected-tm-api-and-tm-web) | none | — | — |

---

## 1. app-backend

- **What it is and why:** the server that holds accounts and wallets and
  sends transactions. Every screen's data comes from it.
- **Direction:** website → app-backend. The website opens no websocket
  today; screens that change (an order, a recovery) ask again every few
  seconds.
- **How they connect:** HTTPS REST through Axios and RTK Query (details
  below). Requests that prove who the user is are **signed** in the
  browser with the user's key (headers `X-TW-SIGNER`, `X-TW-PUBLIC-KEY`,
  `X-TW-TIMESTAMP`, `X-TW-SIGNATURE`); the key never leaves the browser.
  Transactions come back from app-backend as a hash, which the website
  signs and sends back.
- **Settings on this side:** [`VITE_API_URL`](CONFIGURATION.md#vite_api_url)
  for the dev server (example `http://localhost:8080`), or
  [`BACKEND_URL`](CONFIGURATION.md#backend_url) on the Docker container
  (example `https://api.trovo.example.com`).
- **Settings on the other side:** none. app-backend accepts calls from any
  website address (CORS `*`).
- **How to check it works:** sign up; with Docker,
  `curl http://<website>/api/swagger/index.html` answers `200`.
- **When it is down:** pages show errors; nothing is stored in the website.

## 2. wallet-core

- **What it is and why:** Trovo's Rust code for keys, recovery phrases,
  signing and Safe wallet addresses, compiled to WebAssembly and copied
  into `src/walletCore/`. The same code runs in the mobile app.
- **Direction:** the website calls it, in the browser.
- **How they connect:** `src/utils/trovoSDK.ts` imports
  `src/walletCore/wallet_core.js`, which loads `wallet_core_bg.wasm`.
- **Settings:** none. To update it, rebuild wallet-core for the web and copy
  the output over `src/walletCore/` ([wallet-core/DEPLOYMENT.md](../wallet-core/DEPLOYMENT.md)).
- **How to check it works:** creating an account shows a 12-word phrase and
  an address.
- **When it is missing or broken:** the website cannot start.

## 3. Flutterwave

- **What it is and why:** a card payment pop-up for buying crypto with
  money (`src/pages/dashboard/home.tsx`, `flutterwave-react-v3`).
- **Direction:** website → Flutterwave.
- **Settings on this side:** the public key, written in the code (a
  placeholder today; see [CONFIGURATION.md](CONFIGURATION.md#flutterwave-public-key)).
- **Settings on the other side:** a Flutterwave merchant account.
- **When it is down or not set up:** card purchase does not work; nothing
  else is affected.

## 4. Block explorer

- **What it is and why:** links from a transaction to
  `https://basescan.org/tx/...` (or Sepolia's) so users can look it up.
- **Direction:** links only; the website does not call it.

## 5. Not connected: tm-api and tm-web

app-web has no code, address or call related to Trovo Manager (`tm-api`,
`tm-web`). Admin actions reach users through app-backend's database and
API.

---

The sections below go into more depth on how app-web talks to app-backend
and uses wallet-core.

## Talking to app-backend in detail

### Data-fetching layer

Requests are made through **RTK Query** (`@reduxjs/toolkit/query/react`),
via a shared API slice at `src/store/api/baseapi/index.ts`:

```ts
export const baseApi = createApi({
  baseQuery: axiosBaseQuery({ baseUrl: BASE_URL }),
  tagTypes: Object.values(tagTypes),
  endpoints: () => ({})
})
```

Feature-specific endpoints are added with `baseApi.injectEndpoints(...)`
in separate files under `src/store/api/`: `authApi.ts`, `walletApis.ts`,
`p2pApis.ts`, `tokenizationApis.ts`, `sharedAccessApis.ts`,
`cacheApi.ts`. Rather than RTK Query's default `fetchBaseQuery`, this
project uses a **custom `axiosBaseQuery`** (`src/store/api/baseapi/axiosBaseQuery.ts`)
that wraps Axios, so all requests ultimately go through Axios.

### Authentication / request signing — two schemes, used together

Reading `axiosBaseQuery.ts` directly (rather than guessing), this app
actually applies **two separate auth mechanisms**, which can both be
present on the same request:

**1. Bearer token (session auth), applied globally via an Axios
interceptor:**

```ts
axios.interceptors.request.use(async (config: any) => {
  const access_token = getStorage(TOKEN)?.accessToken ?? '';
  config.headers = {
    ...config.headers,
    Authorization: `Bearer ${access_token}`
  };
  return config;
});
```

The access token is read from `localStorage` (via `getStorage`, keyed by
the `TOKEN` constant in `src/store/constants.ts`) and sent as a standard
`Authorization: Bearer <token>` header on **every** Axios request, once
the user has logged in and a token has been stored. (There's a
commented-out 401-refresh interceptor in `src/store/reduxStore/index.ts`
that isn't currently active — token refresh-on-401 is not implemented at
the moment.)

**2. Per-request signature auth, for endpoints that pass `creds`:**

Several endpoints (registration, account recovery, security questions,
etc. — see `authApi.ts`) pass a `creds` object
(`{ signer, address, secretKey }`) alongside the request payload. When
present, `axiosBaseQuery` calls `getRequestHeaders(url, signer, address,
secretKey)`, which builds a message to sign as:

```ts
const toSign = uri + signer + serverTs;   // serverTs = current unix time, seconds
const signature = signHTTP(toSign, secretKey);
```

`signHTTP` (in `src/utils/trovoSDK.ts`) signs that string using EIP-191
`personal_sign`, implemented by the `wallet-core` WASM module (see part
"Using wallet-core in detail" below) — the private key never leaves the browser; only the
signature is sent. The resulting headers are:

| Header | Value |
|---|---|
| `X-TW-SIGNATURE` | Base64 EIP-191 signature of `uri + signer + timestamp` |
| `X-TW-PUBLIC-KEY` | The account's public address |
| `X-TW-SIGNER` | The signer identifier passed in `creds` |
| `X-TW-DEVICE-ID` | `navigator.userAgent` (a crude device identifier — not a stable device ID) |
| `X-TW-TIMESTAMP` | Unix timestamp (seconds) used in the signed message |

This is used for flows that need to prove control of a wallet's private
key directly (e.g. account registration, account-recovery/OTP flows)
rather than relying on an already-issued bearer token — presumably
because the user isn't authenticated with a session token yet at that
point in the flow. Note: `axiosBaseQuery.ts` currently prints the
string-to-sign, the address, the path, the signature and the timestamp of
every signed request to the browser console (`console.log`); remove or
hide these before a hardened production build. The private key is not
printed.

### Account recovery

Recovering an account on the web starts the recovery (`POST
/v1/users/account/recover`, `commit: 1`) and then polls `GET
/v1/account/recovery/status/:username` with the new key: the recovery
takes effect after its waiting period (the account owner can cancel it
meanwhile), and only then does the page send the user to import the
account with the new key. See
[app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md#account-recovery-opt-in-guardian).

### Bank (NGN) deposits and withdrawals

The wallet panel's **Deposit/Withdraw** button opens `/dashboard/bank`
(`src/pages/dashboard/bank.tsx`, `src/store/api/bankApis.ts`) for the active
wallet: deposit Naira into a virtual account for cNGN, withdraw cNGN to a
Nigerian bank account (build, confirm, `signBase64Txn`, submit - like a
payment) and see past withdrawals. Users are onboarded with Stablerail when
KYC level 1 completes in the mobile app; there is no BVN step on the web. See
`app-backend/INTEGRATION.md` ("Bank deposits and withdrawals").

### Public Markets (tokenized NGX stocks and FMDQ bonds)

The **Public Markets** sidebar item and the home page strip open
`/dashboard/public-markets` (`src/pages/dashboard/publicMarkets/`,
`src/store/api/publicMarketsApis.ts`, types in `src/types/publicMarkets.ts`):

| Route | Page | Calls |
| --- | --- | --- |
| `public-markets` | catalogue with search and All / Equities / Bonds / Top gainers | `GET /v1/public-markets` (unsigned) |
| `public-markets/asset/:code` | price, chart (inline SVG), position, key figures, custody chain, session, corporate actions | `GET .../assets/:code`, `.../prices?range=` (unsigned); portfolio (signed) |
| `public-markets/asset/:code/trade` | buy (amount of cNGN) or sell (quantity) from one of the user's own wallets, live quote | `GET .../quote`, `POST .../buy` / `.../sell` |
| `public-markets/order/:orderId` | order progress, refreshed every 5 s until final | `GET /v1/public-markets/orders/:id` |
| `public-markets/portfolio` | "My Stocks": value, returns, income, holdings, orders, activity | `GET .../portfolio`, `GET .../orders` |
| `public-markets/dividends` | dividends and coupons, net of withholding tax | `GET .../dividends` |

A trade is two calls, as with bank withdrawals: the unsigned `buy`/`sell`
returns `{quote, transaction, messages}`; after the user confirms, the
transaction is signed with `signBase64Txn(secretKey, transaction, '')` and
the same call is repeated with `transaction` and `transactionSignature`,
returning `{quote, order}`. Shared wallets are not offered for trading.
Dividends are paid by Public Markets' own distribution engine, not the
tokenized-asset proceeds payout. See `app-backend/PUBLIC_MARKETS.md`.

## Using wallet-core in detail

Yes — confirmed by code, not inferred. `app-web` vendors the compiled
output of the `wallet-core` Rust crate under `src/walletCore/`
(`wallet_core.js`, `wallet_core_bg.wasm`, and `.d.ts` type declarations —
this is `wasm-bindgen`-generated glue code plus the compiled WASM binary).

It's used for **client-side key generation, mnemonic handling, and
message signing** — never for making network calls itself. The
integration point is `src/utils/trovoSDK.ts`, which imports directly from
`../walletCore/wallet_core.js`:

```ts
import init, {
  generateKeypair,
  keypairFromPrivateKey,
  addressFromPrivateKey,
  signPersonal,
  signPersonalBytes,
  generateMnemonic as wcGenerateMnemonic,
  keypairFromMnemonic,
  primarySafeAddress,
} from '../walletCore/wallet_core.js';

await init();
```

That top-level `await init()` (which asynchronously fetches/instantiates
the `.wasm` binary) blocks evaluation of `trovoSDK.ts` — and transitively
every module that imports it, up to the app's entry point — until the
WASM module is ready, so every exported helper in `trovoSDK.ts` can be
called synchronously afterwards. Per the comments in that file, this is
the same Rust core the `app-mobile` Flutter app calls through Dart FFI,
so wallet key derivation and signing logic is implemented once (in Rust)
and reused across platforms instead of being reimplemented per-platform.

Concretely, `trovoSDK.ts` wraps the WASM exports into these app-level
functions, all backed by `wallet-core`:

| `trovoSDK.ts` function | Backed by (wallet-core) | Purpose |
|---|---|---|
| `createAccount()` | `generateKeypair()` | Generate a brand-new keypair for a new wallet |
| `importAccount(secretKey)` | `addressFromPrivateKey()` | Derive an address from an existing private key |
| `parseSecretKey(secretKey)` | `keypairFromPrivateKey()` | Parse/validate a private key into a full keypair |
| `generateMnemonic()` | `generateMnemonic()` | Generate a fresh 12-word BIP39 mnemonic |
| `getCredsFromPassPhrase(passphrase)` | `keypairFromMnemonic()` | Derive a keypair from a BIP39 mnemonic (BIP44 path `m/44'/60'/0'/0/0`) |
| `signHTTP(toSign, secretKey)` | `signPersonal()` | EIP-191 `personal_sign` over a string — used for the API request-signing headers described above |
| `signBase64Txn(secretKey, transactionDigest, _networkPassphrase)` | `signPersonalBytes()` | EIP-191 `personal_sign` over raw bytes decoded from a base64 digest — used to sign transaction digests computed upstream by `app-backend`; the `_networkPassphrase` parameter is unused/vestigial (kept only so call sites don't need to change their argument count) |
| `primaryWalletAddress(signer)` | `primarySafeAddress(signer, "0")` | The user's primary wallet address: a Safe owned by the key, at a fixed address before it is deployed. Registration sends it as `X-TW-PUBLIC-KEY` (the key's address is `X-TW-SIGNER`); the backend refuses any other |

In short: all private-key material (creation, import, mnemonic
derivation, and signing) is handled by the `wallet-core` WASM module
in-browser; private keys are never sent to `app-backend` — only derived
public addresses and signatures are.

### Wallets are Safes: signer vs address

Every wallet is a Safe smart account. The user's **key** (from the
mnemonic) is its owner and signs everything; the **wallet address** is the
Safe. So `X-TW-SIGNER` is always `appUser.primarySigner` and
`X-TW-PUBLIC-KEY` is the wallet address (`appUser.address` for the primary
wallet) - they are no longer the same value.

- **Sends** keep their two-call shape: `transaction` in the first
  response is the base64 hash of the wallet operation; `signBase64Txn`
  signs it unchanged.
- **Sub-wallets** are Safes owned by the same key: the "Add Subwallet"
  form only takes a tag and description, the backend picks the address,
  and only `primarySignature` is sent (there is no sub-wallet key to
  generate, import or back up).
- A wallet's `activated` is `false` until its first send deploys it; the
  wallet card shows this. It can receive funds before that.
- **Settings -> Network fees** (`/dashboard/settings`) lets the user pick
  the stablecoin their wallets pay network fees in
  (`GET /v1/users/settings/gas-fee-assets`,
  `PUT /v1/users/settings/gas-fee-asset`; `src/store/api/settingsApis.ts`).

# Integration with the rest of the monorepo

`app-web` is a pure frontend with two integration points into the rest of
the monorepo: it calls `app-backend`'s REST/WebSocket API over HTTP, and
it uses a WASM package built from the `wallet-core` Rust crate for
client-side key generation and request signing. It has no server-side
code and no direct relationship with the admin-side projects (`tm-api` /
`tm-web`).

## (a) Talking to `app-backend`

### Base URL

The API base URL is `src/store/config.ts`'s `BASE_URL`, read from
`VITE_API_URL` (build-time) with a same-origin `/api` fallback used in
production behind the nginx proxy. See
[`CONFIGURATION.md`](./CONFIGURATION.md) and
[`DEPLOYMENT.md`](./DEPLOYMENT.md) for how that value is set per
environment.

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
(b) below) — the private key never leaves the browser; only the
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
point in the flow. Note: `axiosBaseQuery.ts` currently has several
`console.log` calls that print the string-to-sign, the public key, the
URI, the computed signature, and the timestamp for every signed request —
worth flagging as something you'd want removed or gated behind a debug
flag before a hardened production build, since it logs signing material
to the browser console.

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

## (b) Using `wallet-core` (WASM)

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
| `signHTTP(toSign, secretKey)` | `signPersonal()` | EIP-191 `personal_sign` over a string — used for the API request-signing headers described in (a) |
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

## (c) No relationship to `tm-api` / `tm-web`

Confirmed: there are no imports, URLs, or references to `tm-api` or
`tm-web` anywhere in `app-web/src`. `tm-api` and `tm-web` form a separate
admin-side stack (internal dashboard + its own backend) elsewhere in the
monorepo, and `app-web` — the end-user wallet app — talks only to
`app-backend`. The two stacks appear to be entirely independent frontends
against (presumably) independent or at most incidentally-related
backends; nothing in this project's code assumes or depends on the
admin side existing.

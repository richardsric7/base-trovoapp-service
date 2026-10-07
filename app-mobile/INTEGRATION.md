# Integration

Everything `app-mobile` connects to, and how. The app talks to app-backend
for everything, makes keys and signs on the phone with wallet-core, and
uses Firebase for push notifications, crash reports and remote values.

| Connects to | Direction | Through | Needed? |
|---|---|---|---|
| [1. app-backend](#1-app-backend) | app → app-backend | HTTPS (REST) | yes |
| [2. wallet-core](#2-wallet-core) | built into the app | native library (Dart FFI) | yes |
| [3. Firebase](#3-firebase) | both ways | Firebase SDKs | yes |
| [4. Market data (DEX Trade)](#4-market-data-dex-trade) | app → app-backend | HTTPS | optional |
| [5. Not connected: tm-api and tm-web](#5-not-connected-tm-api-and-tm-web) | none | — | — |

---

## 1. app-backend

- **What it is and why:** the server that holds accounts and wallets and
  sends transactions. Every screen's data comes from it.
- **Direction:** app → app-backend. The app opens no websocket; screens
  that change (orders, recoveries, logins waiting for approval) ask again
  every few seconds, and push notifications tell the app something
  happened.
- **How they connect:** HTTPS REST through `lib/network/requests.dart`
  (and `P2PApi` on top of it). Requests are **signed** on the phone with
  the user's key (headers `X-TW-SIGNER`, `X-TW-PUBLIC-KEY`,
  `X-TW-TIMESTAMP`, `X-TW-SIGNATURE`, plus `X-TW-DEVICE-ID` and
  `X-TW-APP-VERSION`); the key never leaves the phone. Transactions come
  back from app-backend as a hash, which the app signs and sends back.
- **Settings on this side:** the Testnet and Mainnet addresses in
  `getTrovoAppBaseURL()` and the in-app network switch
  ([CONFIGURATION.md](CONFIGURATION.md#1-network-and-servers)).
- **Settings on the other side:** app-backend must be on HTTPS. Its
  dynamic-link settings (`DYNAMIC_LINKS_ANDROID_PACKAGE_NAME`,
  `DYNAMIC_LINKS_IOS_BUNDLE_ID`) must be this app's id, `com.trovo.wallet`.
- **How to check it works:** create a wallet; the account appears in
  app-backend's `users` table.
- **When it is down:** screens show network errors; keys stay safe on the
  phone.

## 2. wallet-core

- **What it is and why:** Trovo's Rust code for keys, recovery phrases,
  signing and Safe wallet addresses. The same code runs in app-web, so
  both apps derive the same wallet from the same phrase.
- **Direction:** the app calls it on the phone.
- **How they connect:** Dart FFI: `lib/functions/wallet_core_ffi_io.dart`
  loads `libwallet_core.so` (Android) or the linked `libwallet_core.a`
  (iOS); `lib/functions/trovo-sdk.dart` (`TrovoWalletSDK`) wraps it.
- **Settings:** the compiled library must be in the app; it is **not in
  the repository** ([DEPLOYMENT.md](DEPLOYMENT.md#4-build-wallet-core-for-the-phone)).
- **How to check it works:** creating a wallet shows a 12-word phrase.
- **When it is missing:** creating, importing and signing fail with
  "library not found".

## 3. Firebase

- **What it is and why:** push notifications (`firebase_messaging`),
  crash reports (`firebase_crashlytics`), analytics (`firebase_analytics`)
  and remote values such as `share_wallet_referral_label`
  (`firebase_remote_config`).
- **Direction:** the app registers for notifications and sends crash
  reports and analytics; app-backend sends notifications through Firebase
  to the phone.
- **How they connect:** the Firebase SDKs, with the project in
  `lib/firebase_options.dart` (one for Testnet, one for Mainnet). The app
  gives app-backend its notification token.
- **Settings on this side:** [Firebase project settings](CONFIGURATION.md#firebase-project-settings-libfirebase_optionsdart).
- **Settings on the other side:** app-backend's `GC` and
  `GOOGLE_PROJECT_ID` must be the **same** Firebase project, or
  notifications never arrive.
- **How to check it works:** receive a payment with the app in the
  background: a notification arrives.
- **When it is down:** no notifications; everything else works.

## 4. Market data (DEX Trade)

- **What it is and why:** the DEX Trade screens (side menu → **DEX
  Trade**, `lib/screens/market_trade/`) show the tokens open for trading,
  each pair's price chart, order book and recent trades, and let the user
  place and cancel buy and sell orders.
- **Direction:** app → app-backend only. There is no outside price feed:
  every figure comes from trades on Trovo's own offer book.
- **How they connect:** `lib/network/market_requests.dart` (`MarketApi`):
  `GET /v1/market/pairs`, `/v1/market/orderbook`, `/v1/market/trades` and
  `/v1/market/candles` (public), and the signed `GET`, `POST` and `DELETE
  /v1/users/trades` for the wallet's own orders. Placing or cancelling an
  order takes two calls: the first returns the operation, the app signs it
  with wallet-core (or, on a shared wallet, sends `commit: 1` for
  approval) and sends it back. Starred pairs are kept on the phone only.
- **Settings on the other side:** app-backend's `OFFER_BOOK_ADDRESS`,
  `NAIRA_ASSET` and `DOLLAR_ASSET`, and tokens marked tradable on the book
  ([app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md)).
- **How to check it works:** open DEX Trade: the pairs list shows prices;
  open a pair: the Order Book tab lists offers, and a small Buy order shows
  under Orders a few seconds after you confirm it.
- **When it is down:** the screens say the market could not be loaded;
  the rest of the app works. Without `OFFER_BOOK_ADDRESS` app-backend
  answers 503 and the list stays empty.

## 5. Not connected: tm-api and tm-web

The app has no address, call or code for Trovo Manager. Admin actions reach
users through app-backend (for example push notifications, or a login
request an admin approves in this app, which goes through app-backend).

---

The sections below go into more depth on how the app talks to app-backend
and uses wallet-core.

## Talking to app-backend in detail

All backend communication goes through `lib/network/requests.dart`
(`makeGetRequest`, `makePostRequest`, `makePutRequest`, etc.) and the
typed wrapper `lib/network/p2p_requests.dart` (`P2PApi`) built on top of
it. There is no second networking layer or generated API client.

### Base URL

`getTrovoAppBaseURL()` (in `lib/network/requests.dart`) picks one of two
hardcoded hosts based on the user's persisted `walletMode` setting:

```dart
Future<String> getTrovoAppBaseURL() async {
  String trovoBaseURL;
  if (await StoreData().storeGetData('walletMode') == "Testnet") {
    trovoBaseURL = 'https://api.dev.trovo.app';
  } else {
    trovoBaseURL = 'https://api.trovotechnologies.com';
  }
  return trovoBaseURL;
}
```

See [CONFIGURATION.md](./CONFIGURATION.md) for the full picture (there is
no env-based override — this is a compiled-in constant, and it's the
same `app-backend` REST API that `app-web` talks to).

### Authentication / request signing scheme

This app does **not** use a bearer token or session cookie. Every signed
request is authenticated by having the client sign a canonical string
with the wallet's private key, done in `getRequestHeader()`
(`lib/network/requests.dart`):

```dart
getRequestHeader({uri, signer, address, secretKey}) async {
  var deviceID = await getDeviceDetails();
  var appVersion = await getAppVersion();
  var serverTs = (DateTime.now().toUtc().millisecondsSinceEpoch / 1000)
      .round()
      .toString();
  var toSign = uri + signer + serverTs;
  var signHTTP = TrovoWalletSDK().signHTTP(toSign: toSign, secretKey: secretKey);

  return {
    "X-TW-SIGNATURE": signHTTP,
    "X-TW-PUBLIC-KEY": address,
    "X-TW-SIGNER": signer,
    "X-TW-DEVICE-ID": deviceID,
    "X-TW-APP-VERSION": appVersion,
    "X-TW-TIMESTAMP": serverTs,
  };
}
```

In plain terms: the client concatenates the request path, the signer
identity, and a Unix timestamp, signs that string with the wallet's
private key (`TrovoWalletSDK().signHTTP`, see below), and sends the
signature, the wallet's public address, the signer, a device ID, the app
version, and the timestamp as headers. app-backend recomputes the string,
checks the signature against `X-TW-SIGNER` and refuses the request (`401`)
if it does not match (see app-backend's
[Authentication](../app-backend/INTEGRATION.md#authentication)).

The actual signing (`TrovoWalletSDK().signHTTP` in
`lib/functions/trovo-sdk.dart`) delegates to `WalletCoreFFI` — see below.

### Account recovery

Turning account recovery on or off, and cancelling a recovery in progress,
change each covered wallet in its own operation: the backend's first answer
(202) carries `transactions` (one per wallet) and `wallets`, and the app
signs every one (`lib/screens/account_recovery/recovery_operations.dart`)
and sends them back as `transactionSignatures`. A device recovering an
account polls `GET /v1/account/recovery/status/:username` with the new key
until the recovery completes after its waiting period. See
[recovery/INTEGRATION.md](../recovery/INTEGRATION.md) and
[app-backend/INTEGRATION.md](../app-backend/INTEGRATION.md#account-recovery-opt-in-guardian).

### Bank (NGN) deposits and withdrawals

cNGN moves to and from Nigerian bank accounts through Stablerail
(`lib/network/fiat_requests.dart`, `lib/screens/fiat/bank_transfer_view.dart`).
It is reached from an asset's existing **Deposit/Withdraw** sheet
(`wrapped_asset.dart`): for cNGN a "Bank account (NGN)" row (Deposit,
Withdraw, History) appears when `GET /v1/users/stablerail/profile` says
Stablerail is enabled. There is no BVN screen: the backend onboards the user
when KYC level 1 completes, and until then the bank screen points to the KYC
screen. A withdrawal is built, shown, confirmed (biometrics where available),
signed with `signBase64Txn` and submitted, like a payment. See
`app-backend/INTEGRATION.md` ("Bank deposits and withdrawals").

### Dividends and interest

An asset's **Dividend and Interest** screen (`dividend_and_yield.dart`), its
history (`dividend_history.dart`; `yield_history.dart` shows the same list
as "Interest History") and a payout's detail
(`dividend_payment_detail.dart`) read
`GET /v1/tokenization/payouts` (`lib/network/payout_requests.dart`,
`lib/models/proceed_payout.dart`). They are matched to the asset by its
token contract, and show what was paid, with its transaction on Basescan,
and what is scheduled but not yet paid. The payouts themselves are made by
`payout-engine`; see `app-backend/INTEGRATION.md` ("Proceeds payouts").

### Public Markets (tokenized NGX equities and FMDQ bonds)

`lib/network/public_markets_requests.dart` calls app-backend's
`/v1/public-markets`. The asset list, an asset's page and its price chart
are unsigned requests. The quote, buy, sell, portfolio, orders and
dividends are signed requests from the primary wallet.

The screens live in `lib/screens/public_markets/`:

- the home screen's Public Markets block;
- the list (search, and filters for equities, bonds and top gainers);
- an asset's page (chart, the user's position, key statistics, how the
  tokens are owned and backed, the trading session, corporate actions);
- the trade screen;
- the order's status page, which refreshes until the order is final;
- My Stocks;
- dividend history and details.

Buying and selling follow the usual wallet flow:

1. Without a signature, the buy or sell call returns the quote and the
   operation to sign.
2. The review sheet shows what happens: whether the order fills now or
   at the next session, the fee and the amounts. The user confirms with
   biometrics.
3. The signed operation (`TrovoWalletSDK().signBase64Txn`) is sent back
   and the order opens.

Only the user's own wallets can trade. Shared wallets with approvers
cannot. Pushes carry `route: publicMarketsOrder` (with `orderId`) or
`publicMarketsDividend` (with `assetCode`), and tapping one opens the
order or the dividends.

## Using wallet-core in detail

This app is **not** fully code-isolated — it's meant to share its
crypto/signing logic with `app-web` through the monorepo's
[`wallet-core`](../wallet-core) Rust crate, though as of this writing that
sharing is only partially wired up on the mobile side.

- `wallet-core/src/core.rs` implements key generation, address
  derivation, and EIP-191-style signing once, in Rust.
- `app-web` consumes it compiled to **WASM** (`app-web/src/utils/trovoSDK.ts`
  imports a vendored WASM build under `app-web/src/walletCore/`).
- `app-mobile` is meant to consume the **same** logic compiled as a
  native library and called through Dart FFI:
  - `lib/functions/wallet_core_ffi.dart` contains hand-written
    `dart:ffi` bindings over `wallet-core/src/ffi.rs`'s C ABI.
  - `lib/functions/trovo-sdk.dart` (`TrovoWalletSDK`) is the façade the
    rest of the app calls (`createAccount`, `signHTTP`, `importAccount`,
    `signBase64Txn`), and it calls into `WalletCoreFFI` — mirroring
    `app-web`'s `trovoSDK.ts` API and derivation scheme (derivation index
    0, matching `app-web` — see the comments in `trovo-sdk.dart` about a
    prior mismatch between the two platforms that was since fixed).

**Caveat — this integration is not fully wired up yet.** The native
library `wallet_core_ffi.dart` expects to load
(`android/app/src/main/jniLibs/<abi>/libwallet_core.so` on Android,
statically linked into the Xcode project on iOS) is not present anywhere
in this repo, and nothing in this repo's build currently compiles or
places it — this matches what `wallet-core/README.md` itself says: the
Dart-side bindings are code-complete and in active use by
`trovo-sdk.dart`, but the FFI path is "not yet wired into this repo's
mobile build." Do not assume a fresh `flutter build`/`flutter run`
produces a working native signing path end-to-end without that build
step being added first (see `wallet-core/README.md`, "Building for
app-mobile", for what that step needs to do).

Other than `wallet-core`, `app-mobile` does not share Dart code,
generated models, or any other package with the rest of the monorepo —
its models (`lib/models/`) are hand-written and specific to this app, and
there is no shared/common package referenced from `pubspec.yaml` that
points back into this monorepo.

### Wallets are Safes: signer vs address

Every wallet is a Safe smart account owned by the user's key. The key (from
the mnemonic) signs everything and is sent as `X-TW-SIGNER`; the wallet
address (the Safe) is `X-TW-PUBLIC-KEY`. They are no longer the same value.

- **Registration/import** compute the primary wallet address with
  `TrovoWalletSDK().primaryWalletAddress(signer)` (wallet-core's
  `wc_primary_safe_address(signer, "0")` over FFI) and keep it in
  `tempAddress`, with the key's address in `tempSigner`.
- **Sends** are unchanged for the app: `transaction` is the base64 hash of
  the wallet operation and `signBase64Txn` signs it.
- **Sub-wallets** (wallets page and the add-wallet dialog) take only a tag,
  description and type. The backend picks the Safe (and an issuing
  wallet's linked distribution Safe); only `primarySignature` is sent, and
  no new key is generated, imported, stored or backed up.
- A wallet's `activated` is `false` until its first send deploys it; the
  wallet details screen says so.
- **Settings -> Network fees** (`lib/widgets/network_fee_setting.dart`)
  picks the stablecoin wallets pay network fees in
  (`GET /v1/users/settings/gas-fee-assets`,
  `PUT /v1/users/settings/gas-fee-asset`).

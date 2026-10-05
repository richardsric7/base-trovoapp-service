# Integration with the rest of the monorepo

## Talking to `app-backend`

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
version, and the timestamp as headers. `app-backend` is presumably able
to verify the signature against the claimed public key/signer and reject
stale timestamps or bad signatures — that verification logic lives in
`app-backend`, not here.

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

## Shared code with the rest of the monorepo: `wallet-core`

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

## Relationship to `tm-api` / `tm-web`

None. `tm-api`/`tm-web` are a separate internal admin console with their
own backend; a search of both (`grep -rln "app-mobile" tm-api tm-web`)
turns up only a mention in `tm-web/README.md`'s own integration notes,
where it explicitly states `tm-web` has no relationship with
`app-backend`, `app-web`, or `app-mobile`. Nothing in this app references
`tm-api`/`tm-web`, and nothing in `tm-api`/`tm-web` builds against or
calls this app.

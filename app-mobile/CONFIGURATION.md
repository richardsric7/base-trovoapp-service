# Configuration

Every parameter `app-mobile` uses. Each one says what it is, why it is
needed, whether you must set it, an example and how to get a real value.

## How to set them

The Trovo mobile app has **no environment variables** and no `.env` file:
nothing in `lib/` reads `--dart-define`, `String.fromEnvironment` or a
dotenv package. Its settings are of four kinds:

| Kind | Where | Changing it needs |
|---|---|---|
| **Written in the code** (server addresses, Firebase projects) | `lib/network/requests.dart`, `lib/firebase_options.dart` | an edit and a new build |
| **A switch in the app** (Testnet or Mainnet) | Settings screen, saved on the phone | nothing: the user flips it |
| **Build files** (app id, version, signing key) | `pubspec.yaml`, `android/`, `ios/` | a new build |
| **Remote values** | Firebase Remote Config console | nothing: the app fetches them |

The app ships with **two sets** of addresses and Firebase projects, one for
**Testnet** (testing) and one for **Mainnet** (real money); the in-app
switch picks between them. A new install starts on **Testnet**.

---

## 1. Network and servers

### Testnet / Mainnet switch (`walletMode`)

- **What it is:** Which set of servers and Firebase project the app uses:
  `Testnet`, or anything else for Mainnet. Saved on the phone under the key
  `walletMode` (`lib/storage/state.dart`, `changeWalletMode`).
- **Why it's needed:** The same app build can be used for testing and for
  real money.
- **Required:** No, default `Testnet`.
- **Example:** `Testnet`
- **How to get it:** The user changes it in **Settings**; the app then
  restarts from the splash screen.

### app-backend address (Testnet)

- **What it is:** The address of the test app-backend, written in
  `getTrovoAppBaseURL()` in `lib/network/requests.dart`. Every API call
  starts with it.
- **Why it's needed:** All of the app's data comes from app-backend.
- **Required:** Yes (already set).
- **Example:** `https://api.dev.trovo.app` (the current value)
- **How to get it:** Your test app-backend's public address
  ([app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md)), HTTPS, no
  trailing `/`. Android and iOS refuse plain `http://` addresses by
  default (the app allows no exception), so to use an app-backend on your
  own computer, give it an HTTPS address with a tunnel (for example
  `ngrok http 8080` prints `https://<random>.ngrok-free.app`), put that
  here and rebuild.

### app-backend address (Mainnet)

- **What it is:** The production app-backend's address, in the same
  function.
- **Why it's needed:** Real users' requests go here.
- **Required:** Yes (already set).
- **Example:** `https://api.trovotechnologies.com` (the current value)
- **How to get it:** Your production app-backend's public address.

The app opens no websocket and calls no blockchain node itself: keys are
made and used on the phone with wallet-core, and app-backend does
everything on the chain. The market screens' prices, charts and order
book come from app-backend too (`/v1/market/*`); there is no outside price
feed.

## 2. Firebase (push notifications, analytics, crash reports, remote values)

### Firebase project settings (`lib/firebase_options.dart`)

- **What it is:** The identifiers of the Firebase project the app connects
  to, one set per platform and network (`androidTestNet`, `androidMainnet`,
  `iosTestNet`, `iosMainnet`): `apiKey`, `appId`, `messagingSenderId`,
  `projectId`, `storageBucket` (and `iosBundleId` on iOS).
- **Why it's needed:** Push notifications (`firebase_messaging`), analytics
  (`firebase_analytics`), crash reports (`firebase_crashlytics`) and remote
  values (`firebase_remote_config`) all go to that project. app-backend
  sends push notifications through the **same** project (its `GC` and
  `GOOGLE_PROJECT_ID` settings).
- **Required:** Yes (already set).
- **Example:**
  ```dart
  static const FirebaseOptions androidTestNet = FirebaseOptions(
    apiKey: '<firebase-api-key>',
    appId: '1:123456789012:android:0123456789abcdef',
    messagingSenderId: '123456789012',
    projectId: 'trovo-wallet-test',
    storageBucket: 'trovo-wallet-test.appspot.com',
  );
  ```
- **How to get it:** Don't type these by hand. Install the FlutterFire tool
  (`dart pub global activate flutterfire_cli`), log in to Firebase
  (`firebase login`), then run `flutterfire configure` in `app-mobile/` and
  pick your project; it writes the values. Because this app keeps two
  projects, copy the generated values into the matching `...TestNet` or
  `...Mainnet` block. These values identify the app; they are not secrets
  (Firebase protects data with its own rules), which is why they are in the
  code.

### Remote value `share_wallet_referral_label`

- **What it is:** The text of the referral share label, fetched from
  Firebase Remote Config (`lib/bottom_bar/bottom_pages/settings.dart`).
- **Why it's needed:** Marketing can change it without a new app release.
- **Required:** No (empty when not set).
- **Example:** `Join me on Trovo and get a bonus`
- **How to get it:** Firebase console → your project → **Remote Config** →
  add a parameter with that name → **Publish**. The app fetches values at
  start, at most once a minute.

## 3. Build files

### App id / bundle id

- **What it is:** The app's unique id in the stores: `com.trovo.wallet`
  (`android/app/build.gradle.kts`: `applicationId` and `namespace`; iOS:
  `PRODUCT_BUNDLE_IDENTIFIER` in `ios/Runner.xcodeproj`).
- **Why it's needed:** The stores, push notifications and Firebase
  identify the app by it.
- **Required:** Yes (already set).
- **Example:** `com.trovo.wallet`
- **How to get it:** Keep it. Changing it makes a different app in the
  stores and needs matching Firebase apps.

### App version (`pubspec.yaml`)

- **What it is:** `version: 0.0.121+165`: the version people see (`0.0.121`)
  and the build number (`165`).
- **Why it's needed:** The stores refuse an upload whose build number is
  not higher than the last one.
- **Required:** Yes.
- **Example:** `version: 0.0.122+166`
- **How to get it:** Raise the build number (after `+`) for every upload,
  and the version when you release new features.

### Android signing key (`android/key.properties`)

- **What it is:** A file pointing to the key that signs Android release
  builds, with four values: `storeFile`, `storePassword`, `keyAlias`,
  `keyPassword` (read by `android/app/build.gradle.kts`).
- **Why it's needed:** Google Play only accepts signed builds, and every
  update must be signed with the same key.
- **Required:** For release builds (debug builds use a built-in key).
- **Example:**
  ```properties
  storeFile=/home/ada/keys/trovo-upload.jks
  storePassword=<your-keystore-password>
  keyAlias=upload
  keyPassword=<your-key-password>
  ```
- **How to get it:** Create the key once (see
  [DEPLOYMENT.md](DEPLOYMENT.md#7-build-an-android-release)), keep it and
  its passwords in your password manager, and write this file on the build
  machine. `key.properties` is in `.gitignore`: never commit it or the key.

### iOS signing team

- **What it is:** The Apple Developer team that signs iOS builds, chosen in
  Xcode (**Runner → Signing & Capabilities → Team**).
- **Why it's needed:** Apple only installs signed apps.
- **Required:** For running on an iPhone and for releases.
- **Example:** `Trovo Technologies Ltd (AB12CD34EF)`
- **How to get it:** An Apple Developer Program membership
  (<https://developer.apple.com/programs/>); your account must be on the
  team.

### wallet-core native library

- **What it is:** The compiled wallet-core library the app calls to make
  keys, recovery phrases and signatures:
  `android/app/src/main/jniLibs/<abi>/libwallet_core.so` on Android, and
  `libwallet_core.a` linked into the Xcode project on iOS
  (`lib/functions/wallet_core_ffi_io.dart`).
- **Why it's needed:** Without it, creating or importing a wallet and
  signing fail with "library not found".
- **Required:** Yes, for any build that creates wallets or signs. **It is
  not in the repository**: you build it from `wallet-core/`.
- **Example:** `android/app/src/main/jniLibs/arm64-v8a/libwallet_core.so`
- **How to get it:** [DEPLOYMENT.md](DEPLOYMENT.md#4-build-wallet-core-for-the-phone)
  step 4.

## Not used

`lib/config/app_settings.config.dart` defines an `appSettings` map that no
code reads. A comment in `lib/screens/kyc_screen.dart` mentions ".env
variables"; there are none.

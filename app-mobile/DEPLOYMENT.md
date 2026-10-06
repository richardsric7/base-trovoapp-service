# Deployment

A step-by-step guide to building the Trovo mobile app and releasing it,
written for someone who has never built a Flutter app. Commands run from
the `app-mobile/` folder unless a step says otherwise.

## 1. What you are deploying

`app-mobile` is the **Trovo app for Android and iOS**. Users create a
wallet (keys stay on the phone), send and receive money, trade P2P, buy
tokenized assets and Public Markets stocks and bonds, and approve logins
and transactions. All data comes from app-backend.

"Deploying" it means building a signed app file and uploading it to
**Google Play** (Android) and the **App Store** (iOS).

**It needs these first:**

| What | Why | Where |
|---|---|---|
| app-backend, at an HTTPS address | every screen's data | [app-backend/DEPLOYMENT.md](../app-backend/DEPLOYMENT.md) |
| wallet-core's native library | keys and signatures on the phone | step 4 below |
| Firebase projects (test and production) | push notifications, crash reports | [CONFIGURATION.md](CONFIGURATION.md#2-firebase-push-notifications-analytics-crash-reports-remote-values) |

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Flutter 3.47 or newer, stable channel | <https://docs.flutter.dev/get-started/install> | `flutter --version` |
| Android Studio (Android SDK and emulator) | <https://developer.android.com/studio> | `flutter doctor` shows ✓ Android toolchain |
| JDK 17 | <https://adoptium.net/> (Temurin 17) | `java -version` |
| Rust (to build wallet-core) | <https://rustup.rs> | `cargo --version` |
| Android NDK (to build wallet-core for Android) | Android Studio → **SDK Manager → SDK Tools → NDK (Side by side)** | `ls $ANDROID_HOME/ndk` |
| For iOS only: a Mac with Xcode and CocoaPods | Mac App Store; `sudo gem install cocoapods` | `xcodebuild -version`, `pod --version` |

Run `flutter doctor` and fix everything it marks ✗ for the platforms you
build (accept Android licenses with `flutter doctor --android-licenses`).
Without a Mac you can still do all the Android steps.

## 3. Get the code and packages

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/app-mobile
flutter pub get
```

**You should see:** `Got dependencies!`.

## 4. Build wallet-core for the phone

The app makes keys and signs with wallet-core, Trovo's Rust library. The
compiled library is **not in the repository**; without it, creating or
importing a wallet fails with "library not found". Build it once, and again
whenever `wallet-core/src/` changes.

**Android** (from the `wallet-core/` folder):

```bash
cd ../wallet-core
rustup target add aarch64-linux-android armv7-linux-androideabi x86_64-linux-android
cargo install cargo-ndk
export ANDROID_NDK_HOME=$ANDROID_HOME/ndk/<version>     # the NDK you installed
cargo ndk -t arm64-v8a -t armeabi-v7a -t x86_64 \
  -o ../app-mobile/android/app/src/main/jniLibs build --release
cd ../app-mobile
```

**You should see:** `libwallet_core.so` in
`android/app/src/main/jniLibs/arm64-v8a/`, `armeabi-v7a/` and `x86_64/`.

**iOS** (on a Mac): build the static library for each target and link it
into Xcode:

```bash
cd ../wallet-core
rustup target add aarch64-apple-ios aarch64-apple-ios-sim
cargo build --release --target aarch64-apple-ios        # iPhones
cargo build --release --target aarch64-apple-ios-sim    # simulators on Apple-silicon Macs
cd ../app-mobile
```

Then open `ios/Runner.xcworkspace` in Xcode, add
`../wallet-core/target/aarch64-apple-ios/release/libwallet_core.a` under
**Runner → Build Phases → Link Binary With Libraries**, and in **Build
Settings → Other Linker Flags** add `-force_load` followed by that file's
path, so its functions are kept (the app finds them with
`DynamicLibrary.process()`). To build for the simulator as well, combine
the two into an XCFramework (`xcodebuild -create-xcframework -library ...
-library ... -output WalletCore.xcframework`) and link that instead.

Details: [wallet-core/DEPLOYMENT.md](../wallet-core/DEPLOYMENT.md). These
steps are not automated in this repository yet; if you automate them, keep
the output paths the same.

## 5. Choose the servers

The app has two app-backend addresses written in
`lib/network/requests.dart` (Testnet and Mainnet) and two Firebase
projects in `lib/firebase_options.dart`; the Settings screen switches
between them, and a new install starts on Testnet. Check they point at
your servers before building ([CONFIGURATION.md](CONFIGURATION.md)).

## 6. Run it on a phone or emulator

1. Start an Android emulator (Android Studio → **Device Manager → Create
   device**), an iOS Simulator (`open -a Simulator`), or plug in a phone
   with developer mode / USB debugging on.
2. Check Flutter sees it:
   ```bash
   flutter devices
   ```
3. Run:
   ```bash
   flutter run
   ```
   **You should see** the app open on the splash screen. Press `r` in the
   terminal to reload after a code change.
4. Create a wallet: you get a 12-word recovery phrase, and the account
   appears in app-backend's `users` table.

Before a release, run the checks:

```bash
flutter analyze
flutter test
```

## 7. Build an Android release

1. **Create the signing key (once, ever).** Every update must be signed
   with the same key, so keep it safe:
   ```bash
   keytool -genkey -v -keystore ~/keys/trovo-upload.jks \
     -keyalg RSA -keysize 2048 -validity 10000 -alias upload
   ```
   It asks for a password and your name and organization. Put the `.jks`
   file and both passwords in your password manager. (With Google Play App
   Signing, this is your **upload** key; Google holds the app signing key.)
2. **Tell the build where it is:** create `android/key.properties`:
   ```properties
   storeFile=/home/ada/keys/trovo-upload.jks
   storePassword=<your-keystore-password>
   keyAlias=upload
   keyPassword=<your-key-password>
   ```
   It is in `.gitignore`; never commit it ([CONFIGURATION.md](CONFIGURATION.md#android-signing-key-androidkeyproperties)).
3. **Raise the version** in `pubspec.yaml` (`version: 0.0.122+166`: the
   number after `+` must be higher than the last upload).
4. **Build:**
   ```bash
   flutter build appbundle      # for Google Play
   flutter build apk            # an installable file for testers
   ```
   **You should see:** `✓ Built build/app/outputs/bundle/release/app-release.aab`
   (or `.../flutter-apk/app-release.apk`).
5. **Upload** the `.aab` in the [Google Play Console](https://play.google.com/console)
   → your app → **Testing → Internal testing → Create new release**. Test
   it, then promote it to **Production**.

## 8. Build an iOS release (Mac only)

1. Join the [Apple Developer Program](https://developer.apple.com/programs/)
   and make sure an App Store Connect app exists for `com.trovo.wallet`.
2. `cd ios && pod install && cd ..`
3. Open `ios/Runner.xcworkspace` in Xcode (the workspace, not the
   `.xcodeproj`). In **Runner → Signing & Capabilities**, choose your
   **Team** and keep **Automatically manage signing** on.
4. Raise the version in `pubspec.yaml` as for Android.
5. Build:
   ```bash
   flutter build ipa
   ```
   **You should see:** `Built IPA to build/ios/ipa`.
6. Upload with the **Transporter** app (Mac App Store) or Xcode
   (**Product → Archive → Distribute App**), test it in **TestFlight**,
   then submit it for review in App Store Connect.

## 9. Check it works

On a test phone with the release build:

1. Create a wallet and write down the phrase.
2. Receive a small test payment, then send one.
3. Put the app in the background, send it a payment from another account:
   a push notification arrives.
4. In Firebase **Crashlytics**, the app shows up (after its first start).

## 10. Updating and rolling back

- **Updating:** raise the version, build (steps 7 and 8) and upload. Users
  get it through the stores. Deploy app-backend first when the release
  needs new API routes, and keep app-backend working for the previous app
  version: users update slowly.
- **Rolling back:** the stores cannot go back to an older build. Build the
  previous code (`git checkout <previous-tag>`) with a **higher** build
  number and release it; on Google Play you can also halt a staged
  rollout.
- There is no CI/CD pipeline in this repository: run `flutter analyze` and
  `flutter test` before merging, and build and upload by hand.

## 11. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| "library not found" / `libwallet_core.so` errors when creating a wallet | wallet-core not built into the app | step 4 |
| iOS: `symbol not found` for wallet-core functions | the library is not linked, or not force-loaded | step 4, `-force_load` |
| every screen shows a network error | wrong app-backend address for the chosen network, or not HTTPS | [CONFIGURATION.md](CONFIGURATION.md#1-network-and-servers) |
| no push notifications | the Firebase project in the app differs from app-backend's, or notification permission was refused | use the same project in both; allow notifications |
| `Keystore file not found` / signing errors in a release build | `android/key.properties` missing or wrong | step 7.2 |
| Google Play: "version code already used" | build number not raised | raise the number after `+` in `pubspec.yaml` |
| strange errors after changing Flutter version | old build files | `flutter clean && flutter pub get` |
| iOS: `pod install` errors | CocoaPods out of date | `sudo gem install cocoapods`, then `pod repo update` |

## Web build (testing only)

The project also builds for the web (`flutter build web`), used only to
check that the app compiles and its screens load. On the web, wallet
functions are not available (`lib/functions/wallet_core_ffi_web.dart`
reports "not available"); the real web wallet is
[app-web](../app-web/README.md).

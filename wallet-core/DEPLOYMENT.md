# Deployment

A step-by-step guide to building `wallet-core` and putting it into the
apps, written for someone who has never used Rust. Commands run from the
`wallet-core/` folder unless a step says otherwise.

## 1. What you are deploying

wallet-core is Trovo's code for **keys, 12-word recovery phrases,
signatures and wallet addresses**, written once in Rust. It is not a
server: nothing runs on its own and there is nothing to health-check.
"Deploying" it means compiling it and copying the result into the apps:

| Output | Goes into | Committed? |
|---|---|---|
| a WebAssembly package (4 files) | `app-web/src/walletCore/` | **yes**: app-web builds without Rust |
| `libwallet_core.so` per Android processor type | `app-mobile/android/app/src/main/jniLibs/<abi>/` | no: build it before building the app |
| `libwallet_core.a` for iOS | linked into `app-mobile/ios` in Xcode | no |

Redo the relevant steps whenever `src/` changes, then rebuild and release
the apps.

## 2. Before you start

| Tool | Install | Check |
|---|---|---|
| Git | <https://git-scm.com/downloads> | `git --version` |
| Rust (rustup, cargo) | `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs \| sh` (or <https://rustup.rs>) | `cargo --version` |
| For the web: the WebAssembly target | `rustup target add wasm32-unknown-unknown` | `rustup target list --installed` |
| For the web: `wasm-bindgen` 0.2.128 | `cargo install wasm-bindgen-cli --version 0.2.128` | `wasm-bindgen --version` |
| For Android: Android NDK and `cargo-ndk` | Android Studio → SDK Manager → NDK; `cargo install cargo-ndk` | `cargo ndk --version` |
| For iOS: a Mac with Xcode | Mac App Store | `xcodebuild -version` |

The `wasm-bindgen` version must equal the one in `Cargo.lock`
([CONFIGURATION.md](CONFIGURATION.md#wasm-bindgen-tool-version)).

## 3. Get the code and run the tests

```bash
git clone https://github.com/richardsric7/base-trovoapp-service.git
cd base-trovoapp-service/wallet-core
cargo test
```

**You should see:** `test result: ok. 11 passed; 0 failed`. The tests check
signing and verifying, phrases giving the same key every time, the address
format, rejecting bad keys, and wallet (Safe) addresses against a file made
from real Safe contracts (`testdata/safe_fixture.json`; app-backend's tests
use an identical copy).

## 4. Build for app-web

```bash
cargo build --release --target wasm32-unknown-unknown
wasm-bindgen --target web --out-dir pkg-web \
  target/wasm32-unknown-unknown/release/wallet_core.wasm
cp pkg-web/wallet_core.js pkg-web/wallet_core.d.ts \
   pkg-web/wallet_core_bg.wasm pkg-web/wallet_core_bg.wasm.d.ts \
   ../app-web/src/walletCore/
```

**You should see:** `Finished release profile`, then four files in
`pkg-web/`. If nothing in `src/` changed, the copied files are identical to
the committed ones (`git status` in `app-web/` shows no change).

Then check and commit app-web:

```bash
cd ../app-web && npm run build && cd ../wallet-core
```

## 5. Build for app-mobile (Android)

```bash
rustup target add aarch64-linux-android armv7-linux-androideabi x86_64-linux-android
export ANDROID_NDK_HOME=$ANDROID_HOME/ndk/27.0.12077973     # your NDK folder
cargo ndk -t arm64-v8a -t armeabi-v7a -t x86_64 \
  -o ../app-mobile/android/app/src/main/jniLibs build --release
```

**You should see:** `libwallet_core.so` in
`../app-mobile/android/app/src/main/jniLibs/arm64-v8a/`, `armeabi-v7a/` and
`x86_64/`. The app loads it with `DynamicLibrary.open('libwallet_core.so')`.

## 6. Build for app-mobile (iOS, Mac only)

```bash
rustup target add aarch64-apple-ios aarch64-apple-ios-sim
cargo build --release --target aarch64-apple-ios
cargo build --release --target aarch64-apple-ios-sim
```

**You should see:** `target/aarch64-apple-ios/release/libwallet_core.a` (and
the `-sim` one). Link it into the app as described in
[app-mobile/DEPLOYMENT.md](../app-mobile/DEPLOYMENT.md#4-build-wallet-core-for-the-phone):
iOS apps include native code inside the app itself, which is why the app
finds the functions with `DynamicLibrary.process()`.

## 7. Check it works

- **Web:** in app-web, create an account: a 12-word phrase and an address
  appear.
- **Mobile:** run the app and create a wallet the same way.
- **Both agree:** import the same phrase in app-web and in app-mobile: the
  wallet address is the same, and app-backend accepts the sign-up (it
  calculates the address itself and refuses a mismatch).

## 8. Updating and rolling back

- **Updating:** change `src/`, run `cargo test`, then steps 4 to 6, then
  release the apps.
- **Rolling back:** `git checkout <previous-tag> -- src/`, rebuild, and
  release the apps again; for the web, you can also restore the previous
  `app-web/src/walletCore/` from git.
- **Never change** the key path or the Safe addresses without a migration
  plan: existing users' keys and wallet addresses depend on them
  ([CONFIGURATION.md](CONFIGURATION.md#3-values-written-in-the-code)).
- There is no CI/CD pipeline in this repository: these steps are manual.

## 9. Troubleshooting

| You see | Cause | Fix |
|---|---|---|
| `it looks like the Rust project used to create this wasm file was linked against version of wasm-bindgen that uses a different bindgen format` | `wasm-bindgen` tool version differs from `Cargo.lock` | install the matching version (section 2) |
| `can't find crate for core` / `target may not be installed` | the target is not installed | `rustup target add <target>` |
| `cargo ndk` cannot find the NDK | `ANDROID_NDK_HOME` not set | set it to your NDK folder |
| the phone app says "library not found" | the `.so` is missing for that phone's processor type | step 5 with all three ABIs |
| app-backend refuses a sign-up from the app | the app and app-backend calculate different wallet addresses | keep wallet-core's Safe addresses equal to app-backend's `SAFE_*` settings |
| app-web shows an old behaviour after a change | the WebAssembly files were not copied | step 4 |

## How the parts fit together

- `src/core.rs`: keys, phrases, signing (plain Rust).
- `src/safe.rs`: wallet (Safe) address calculation.
- `src/wasm.rs`: the JavaScript functions (only in WebAssembly builds).
- `src/ffi.rs`: the C functions Dart calls (only in phone and desktop
  builds).
- `Cargo.toml`'s `crate-type = ["cdylib", "staticlib", "rlib"]` makes the
  shared library (`.so`, `.wasm`), the static library (`.a`) and the Rust
  library `cargo test` uses, in one build.

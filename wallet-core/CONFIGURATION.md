# Configuration

Every parameter `wallet-core` uses. Each one says what it is, why it is
needed, whether you must set it, an example and how to get a real value.

## How to set them

wallet-core is a library, not a running program: it reads **no environment
variables**, has no settings file and no Cargo features (`Cargo.toml` has
no `[features]`; `src/` has no `env::var` or `cfg(feature)`). Given the
same inputs, it always gives the same outputs.

What you choose is **how you build it** (sections 1 and 2), and a few
values are **written in the code** and must agree with app-backend
(section 3).

---

## 1. Build choices

### Build target (`--target`)

- **What it is:** Which kind of computer or phone to compile for. It picks
  the binding code: WebAssembly builds include `src/wasm.rs` (for app-web),
  every other target includes `src/ffi.rs` (for app-mobile).
- **Why it's needed:** A browser, an Android phone and an iPhone each need
  their own compiled file.
- **Required:** Yes, for the web and the phones (none for `cargo test`,
  which builds for your own computer).
- **Example:** `wasm32-unknown-unknown`
- **How to get it:**

  | For | Target | Output |
  |---|---|---|
  | app-web | `wasm32-unknown-unknown` | `target/wasm32-unknown-unknown/release/wallet_core.wasm` |
  | Android phones (64-bit) | `aarch64-linux-android` | `libwallet_core.so` |
  | Android phones (32-bit) | `armv7-linux-androideabi` | `libwallet_core.so` |
  | Android emulator | `x86_64-linux-android` | `libwallet_core.so` |
  | iPhones | `aarch64-apple-ios` | `libwallet_core.a` |
  | iOS simulator (Apple-silicon Mac) | `aarch64-apple-ios-sim` | `libwallet_core.a` |

  Install each once with `rustup target add <target>`.

### `wasm-bindgen` tool version

- **What it is:** The version of the `wasm-bindgen` command that turns the
  `.wasm` file into a JavaScript module.
- **Why it's needed:** It must be **exactly** the version of the
  `wasm-bindgen` library in `Cargo.lock`, or it stops with a version
  mismatch error.
- **Required:** Yes, for the web build.
- **Example:** `0.2.128`
- **How to get it:** `grep -A1 'name = "wasm-bindgen"' Cargo.lock`, then
  `cargo install wasm-bindgen-cli --version <that version>`.

### `wasm-bindgen --target web` and `--out-dir`

- **What they are:** `--target web` makes a module a browser loads directly
  (`import init from './wallet_core.js'; await init()`), which is how
  app-web's `src/utils/trovoSDK.ts` uses it. `--out-dir` is the folder the
  four output files go to.
- **Why they're needed:** Other `--target` values make modules app-web
  cannot load the same way.
- **Required:** Yes, for the web build.
- **Example:** `--target web --out-dir pkg-web`
- **How to get it:** Use these values; `pkg-web/` is in `.gitignore`.

### Android ABIs (`cargo ndk -t`)

- **What they are:** Which Android processor types to build for:
  `arm64-v8a`, `armeabi-v7a`, `x86_64`.
- **Why they're needed:** An Android app needs one library per processor
  type it runs on; a missing one fails on those phones.
- **Required:** For Android; `arm64-v8a` at least (almost all phones),
  `x86_64` for emulators.
- **Example:** `-t arm64-v8a -t armeabi-v7a -t x86_64`
- **How to get it:** Use all three.

### `ANDROID_NDK_HOME`

- **What it is:** Where the Android NDK (the compilers for Android) is
  installed.
- **Why it's needed:** `cargo ndk` uses it to build the Android library.
- **Required:** For Android builds.
- **Example:** `/home/ada/Android/Sdk/ndk/27.0.12077973`
- **How to get it:** Install the NDK in Android Studio (**SDK Manager → SDK
  Tools → NDK (Side by side)**), then look in `<Android SDK>/ndk/`.
  app-mobile's Android build uses NDK `27.0.12077973`
  (`app-mobile/android/app/build.gradle.kts`).

## 2. Release profile

`Cargo.toml` builds releases small (`opt-level = "s"`, `lto = true`). Keep
it: the web file is downloaded by every visitor.

## 3. Values written in the code

### Safe contract addresses (`src/safe.rs`)

- **What they are:** The addresses of Safe's standard contracts that every
  Trovo wallet is built from: `SAFE_PROXY_FACTORY`
  (`0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67`), `SAFE_L2_SINGLETON`
  (`0x29fcB43b46531BcA003ddC8FCB67FFE91900C762`), `SAFE_MODULE_SETUP`
  (`0x2dd68b007B46fBe91B9A7c3EDa5A7a1063cB5b47`), `SAFE_4337_MODULE`
  (`0x75cf11467937ce3F2f357CE24ffc3DBF8fD5c226`), and the Safe proxy's
  creation code.
- **Why they're needed:** A wallet's address is calculated from them before
  the wallet exists. The apps calculate it with wallet-core, app-backend
  calculates it again, and **app-backend refuses a sign-up whose address
  does not match its own**.
- **Required:** Yes (already set to Safe's official deployments, the same
  on Base and Base Sepolia).
- **Example:** `0x4e1DCf7AD4e460CfD30791CCC4F9c8a4f820ec67`
- **How to get it:** Keep them. They must equal app-backend's
  `SAFE_PROXY_FACTORY_ADDRESS`, `SAFE_SINGLETON_ADDRESS`,
  `SAFE_MODULE_SETUP_ADDRESS` and `SAFE_4337_MODULE_ADDRESS`
  ([app-backend/CONFIGURATION.md](../app-backend/CONFIGURATION.md#safe-and-entrypoint-contract-addresses)),
  which default to the same values. Change both together, then rebuild
  wallet-core into both apps.

### Key derivation path

- **What it is:** Keys come from the 12-word phrase at path
  `m/44'/60'/0'/0/{index}` (the standard Ethereum path), index `0` for the
  user's key.
- **Why it's needed:** The same phrase must give the same key in the web
  app, the mobile app and other Ethereum wallets.
- **Required:** Yes (already set).
- **Example:** `m/44'/60'/0'/0/0`
- **How to get it:** Never change it: existing users would get different
  keys.

# Deployment

"Deployment" for `wallet-core` doesn't mean a running service — there's
no server, no URL, nothing to health-check. It means: compile this Rust
crate into an artifact another project can link against, and get that
artifact into the consuming project. This document covers both artifact
types the crate produces (see `Cargo.toml`'s `crate-type`):

1. A **WebAssembly package** for `app-web` (and, if wired up in future,
   `tm-web`).
2. A **native shared/static library** for `app-mobile`, called through
   Dart's `dart:ffi`.

If you've never touched Rust or wasm-pack before, read section 1 fully —
it explains the moving parts (`rustup`, targets, `wasm-bindgen`) that
section 2 also depends on.

## 1. Prerequisites

### 1.1 Install Rust (via rustup)

[rustup](https://rustup.rs/) is Rust's official toolchain installer — it
installs the compiler (`rustc`), package manager/build tool (`cargo`),
and lets you add extra compilation targets later.

```sh
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

Follow the interactive prompts (the default option is fine), then either
restart your shell or run the `source` command it prints. Verify it
worked:

```sh
rustc --version
cargo --version
```

### 1.2 Add the WebAssembly target

By default, `rustc` only compiles for your machine's own OS/architecture.
To produce WebAssembly, you add a "target" — a cross-compilation profile:

```sh
rustup target add wasm32-unknown-unknown
```

This is a one-time step per machine. `wasm32-unknown-unknown` means: WASM
32-bit, with no assumed host OS/environment (as opposed to
`wasm32-wasi`, which targets a WASI runtime) — the right choice for code
that will run inside a browser's own WASM engine.

### 1.3 Install `wasm-bindgen-cli`

Rust's compiler can emit a raw `.wasm` binary, but JavaScript can't call
into it usefully on its own — `wasm-bindgen` generates the glue: a
`.js`/`.d.ts` file that wraps the raw WASM exports in an ergonomic,
typed JS/TS API (this is what turns a Rust `Result<T, E>` into a thrown
JS exception, a Rust `String` into a JS `string`, etc).

The **CLI version must exactly match** the `wasm-bindgen` crate version
this project depends on (`Cargo.lock` currently pins `0.2.128`) — a
mismatch fails the build with a version-mismatch error. Check
`Cargo.lock` for the current pinned version before running this:

```sh
grep -A2 'name = "wasm-bindgen"' Cargo.lock   # confirm the version first
cargo install wasm-bindgen-cli --version 0.2.128
```

> **Note on `wasm-pack`**: `src/wasm.rs`'s doc comment mentions
> `wasm-pack build --target web` as a possible build command. In
> practice, the `pkg-web/` output actually checked into `app-web`
> (`app-web/src/walletCore/`) has **no `package.json`** — which
> `wasm-pack` always generates — meaning it was produced by the raw
> `cargo build` + `wasm-bindgen` CLI steps below, not `wasm-pack`. Either
> approach produces compatible output, but the commands below are the
> ones that reproduce what's actually in `pkg-web/` today. If you'd
> rather use `wasm-pack` (it wraps the same two steps into one command
> and also installs `wasm-bindgen-cli` for you), install it per
> [rustwasm.github.io/wasm-pack/installer](https://rustwasm.github.io/wasm-pack/installer/)
> and run `wasm-pack build --target web --out-dir pkg-web` instead of
> steps 2.2–2.3 below — but note it *will* additionally emit a
> `pkg-web/package.json`, which isn't there today.

## 2. Building the WASM package (for app-web)

From the `wallet-core/` directory:

```sh
# 2.1 Run the test suite first (native, doesn't need wasm) - see section 4.
cargo test

# 2.2 Compile to WebAssembly, release profile
cargo build --release --target wasm32-unknown-unknown

# 2.3 Generate the JS/TS bindings from the compiled .wasm
wasm-bindgen --target web --out-dir pkg-web \
  target/wasm32-unknown-unknown/release/wallet_core.wasm
```

`--target web` tells `wasm-bindgen` to generate an ES module meant to be
loaded directly by a browser (`import init, { ... } from './wallet_core.js'`
then `await init()`), which matches how `app-web`'s `trovoSDK.ts` uses it
(see INTEGRATION.md).

This produces, in `pkg-web/`:

- `wallet_core.js` — the JS glue/wrapper module.
- `wallet_core.d.ts` — its TypeScript type declarations.
- `wallet_core_bg.wasm` — the actual compiled WebAssembly binary.
- `wallet_core_bg.wasm.d.ts` — low-level type declarations for the raw
  wasm exports (used internally by `wallet_core.js`; consumers don't
  import this directly).

### 2.1 Getting the build into app-web

**What's actually wired up today**: this is a manual copy, not an
automated dependency. `app-web/package.json` has no `wallet-core` entry
and `app-web`'s own `npm run build` (`tsc && vite build`) does not invoke
Rust or wasm-bindgen at all — that's intentional, so `app-web` contributors
don't need a Rust toolchain just to build the frontend. Instead, the four
files above are committed directly into `app-web/src/walletCore/`, and
whoever changes `wallet-core/src/` re-runs the build above and re-copies:

```sh
cp pkg-web/wallet_core.js pkg-web/wallet_core.d.ts \
   pkg-web/wallet_core_bg.wasm pkg-web/wallet_core_bg.wasm.d.ts \
   ../app-web/src/walletCore/
```

`app-web/src/utils/trovoSDK.ts` then imports directly from
`'../walletCore/wallet_core.js'`.

**If you wanted to wire this into `tm-web` (or any other frontend)**:
`tm-web` has no reference to `wallet-core` today (confirmed — nothing in
`tm-web/package.json` or its source tree). The lowest-effort path is the
same vendoring pattern above (build, then `cp` the four files into
`tm-web/src/<somewhere>/`). For a real npm dependency instead of a
vendored copy, you'd need `wasm-pack build --target web --out-dir pkg-web`
(not the raw `wasm-bindgen` CLI) so a `package.json` is generated in
`pkg-web/`, then add to the consumer's `package.json`:

```json
"dependencies": {
  "wallet-core": "file:../wallet-core/pkg-web"
}
```

and `npm install` to symlink it in, importing it as `import { ... } from
"wallet-core"` rather than a relative path. **This local-path-dependency
approach is not what `app-web` does today** — it's described here only
as the mechanism you'd reach for if you wanted an automated dependency
instead of copy-pasting build output.

## 3. Building the native library (for app-mobile / Dart FFI)

`app-mobile`'s Dart code (`app-mobile/lib/functions/wallet_core_ffi.dart`)
already contains hand-written `dart:ffi` bindings expecting a native
`wallet-core` library to be present at runtime — but **no build step in
this repo produces or places that library today**. This section explains
how you'd build and place it if you're picking that work up.

### 3.1 Add the relevant target triple(s)

Android and iOS each need their own cross-compilation targets:

```sh
# Android (one per ABI you need to support)
rustup target add aarch64-linux-android      # 64-bit ARM (most real devices)
rustup target add armv7-linux-androideabi    # 32-bit ARM (older devices)
rustup target add x86_64-linux-android       # 64-bit x86 (emulators)

# iOS
rustup target add aarch64-apple-ios          # physical devices
rustup target add aarch64-apple-ios-sim      # Apple Silicon simulators
rustup target add x86_64-apple-ios           # Intel simulators
```

### 3.2 Build

**Android** needs the Android NDK's linkers configured (`cargo` doesn't
know where they are by default). The simplest way is
[`cargo-ndk`](https://github.com/bbqsrc/cargo-ndk), which wraps this for
you:

```sh
cargo install cargo-ndk
cargo ndk -t arm64-v8a -t armeabi-v7a -t x86_64 -o jniLibs build --release
```

This produces `libwallet_core.so` per ABI under `jniLibs/<abi>/`, which
is exactly where `wallet_core_ffi.dart`'s comment says Android expects
it: `app-mobile/android/app/src/main/jniLibs/<abi>/libwallet_core.so`
(that directory does not currently exist in `app-mobile/`).

**iOS** builds a static library per target and links it directly into
the Xcode project (Flutter iOS builds statically link native code rather
than loading a bundled dynamic library at runtime — this is why
`wallet_core_ffi.dart` uses `DynamicLibrary.process()` on iOS/macOS
instead of `DynamicLibrary.open(...)`):

```sh
cargo build --release --target aarch64-apple-ios
# -> target/aarch64-apple-ios/release/libwallet_core.a
```

Combine the per-architecture `.a` files (e.g. via `lipo` or an XCFramework)
and add the result to the Xcode project's linked libraries. Consult
Flutter's own [native/FFI plugin build docs](https://docs.flutter.dev/platform-integration/platform-channels)
for the current recommended packaging approach, since none of this is
set up in `app-mobile/` yet.

### 3.3 Current state — be honest about this

As of this writing: **no `.so`/`.a`/`.dylib` file for `wallet-core`
exists anywhere in `app-mobile/`, no Gradle/Podfile/Xcode config
references `libwallet_core`, and no CI step builds one** (confirmed by
searching `app-mobile/android` and `app-mobile/ios` for `jniLibs`,
`libwallet_core`, and `XCFramework` — no matches). The Dart FFI bindings
in `wallet_core_ffi.dart` are real, typed, and already called from
`trovo-sdk.dart` — but calling any of them today would fail at runtime
with a "library not found" error, because
`DynamicLibrary.open('libwallet_core.so')` (or `.process()` on iOS, which
would fail to find the expected symbols) has nothing to actually load.
Building and wiring in this native library is unfinished work, not
something broken.

## 4. Running tests

The `rlib` build target (see `Cargo.toml`) is what `cargo test` compiles
and runs — plain native code, no wasm or mobile toolchain needed:

```sh
cargo test
```

This exercises (see the `tests` module at the bottom of `src/core.rs`):

- Keypair generation + EIP-191 sign/verify round-tripping.
- Signature-based signer recovery.
- Private-key import producing the same address as the original keypair.
- Tolerance for an uppercase `0x`-prefixed hex private key (a
  compatibility case for `app-mobile`'s legacy Stellar-secret-key
  convention).
- BIP39/BIP44 mnemonic derivation determinism (same mnemonic + index →
  same address, different index → different address).
- An EIP-55 checksum encoding against the spec's own test vector.
- Rejection of malformed private keys.

Run this before generating a new build for either the web or mobile path
— it's fast, needs no extra setup, and catches logic regressions before
they reach a platform-specific binding layer.

## 5. CI / build pipeline

**There is no CI pipeline for `wallet-core` today.** Both workflow files
in `/home/user/src-monorepo/.github/workflows/` were checked directly:

- **`deploy.yml`** ("Build & Push") builds and pushes Docker images for
  `backend`, `web`, `trovotech-io`, and `trovo-app-website` on pushes to
  `dev`/`staging`/`main`. It never mentions `wallet-core`, Rust, `cargo`,
  or `wasm-pack`/`wasm-bindgen`.
- **`pr-checks.yml`** runs lint/build/test jobs for `backend` (Go),
  `web` (Vite/React), and `mobile` (Flutter) on every PR. Again, no job
  installs a Rust toolchain, runs `cargo test`, or builds this crate in
  any form.

Neither workflow's `paths-filter` patterns (`backend/**`, `web/**`,
`mobile/**`, etc.) match this repo's actual top-level directory names
(`wallet-core/`, `app-web/`, `app-mobile/`, `tm-web/`, `tm-api/`,
`app-backend/`) — so even a change *inside* `wallet-core/` wouldn't be
detected as relevant to any existing job as currently configured. In
short: there is no automated build, test, or artifact-publishing step
for this crate — `cargo test`, building `pkg-web/`, and copying it into
`app-web` are all manual, local steps today.

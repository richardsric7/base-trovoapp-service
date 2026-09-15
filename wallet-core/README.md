# wallet-core

Shared Base (EVM/secp256k1) wallet logic - key generation, address
derivation, EIP-191 `personal_sign`/verify, and BIP39/BIP44 mnemonic
derivation - used by both `app-web` (compiled to WebAssembly) and
`app-mobile` (compiled to a native library, called through Dart's
`dart:ffi`). One implementation, reviewed once, instead of the Stellar
keypair/signing logic each app previously carried independently
(`app-web/src/utils/trovoSDK.ts`, `app-mobile/lib/functions/trovo-sdk.dart`).

## Layout

- `src/core.rs` - the actual logic, plain Rust, no platform-specific types.
- `src/wasm.rs` - `wasm-bindgen` bindings over `core`, compiled in only for
  `--target wasm32-unknown-unknown`.
- `src/ffi.rs` - a C ABI (`extern "C"`) surface over `core`, returning
  heap-allocated JSON C strings (freed via `wc_free_string`), compiled in
  for every other target.

## Building for app-web (WASM)

```sh
rustup target add wasm32-unknown-unknown   # once
cargo install wasm-bindgen-cli --version <version matching the wasm-bindgen crate in Cargo.toml>

cargo build --release --target wasm32-unknown-unknown
wasm-bindgen --target web --out-dir pkg-web \
  target/wasm32-unknown-unknown/release/wallet_core.wasm
```

That produces `pkg-web/wallet_core.js` (+ `.d.ts`) and
`pkg-web/wallet_core_bg.wasm`. Copy those into `app-web/src/walletCore/`
(vendored/committed there, not a build step in `app-web`'s own
`npm run build`, so app-web doesn't need a Rust toolchain):

```sh
cp pkg-web/wallet_core.js pkg-web/wallet_core.d.ts \
   pkg-web/wallet_core_bg.wasm pkg-web/wallet_core_bg.wasm.d.ts \
   ../app-web/src/walletCore/
```

`app-web/src/utils/trovoSDK.ts` imports from `../walletCore/wallet_core.js`
directly - re-run the copy step above after changing anything under `src/`
that affects `wasm.rs`'s exported surface.

## Building for app-mobile (Dart FFI)

Build the native shared library for each target platform's ABI(s) (see
each target's own cross-compilation setup - e.g. `cargo-ndk` for Android,
`cargo build --target aarch64-apple-ios` for iOS) and place the resulting
library where `app-mobile`'s FFI loader
(`app-mobile/lib/functions/wallet_core_ffi.dart`) expects it per platform
(Android: `jniLibs/<abi>/libwallet_core.so`; iOS: linked into the app
binary as a static/XCFramework; see that file's loader for the exact
per-platform paths). Example for a local/native build:

```sh
cargo build --release
# -> target/release/libwallet_core.so (Linux) / .dylib (macOS) / wallet_core.dll (Windows)
```

The C surface is declared in `src/ffi.rs`; `app-mobile`'s Dart bindings
(`app-mobile/lib/functions/wallet_core_ffi.dart`) mirror those signatures
by hand (no cbindgen/flutter_rust_bridge codegen step wired up yet - see
that file's doc comment).

## Testing

```sh
cargo test
```

Exercises key generation, EIP-191 sign/verify round-tripping, BIP44
mnemonic derivation determinism, and an EIP-55 checksum test vector - see
`src/core.rs`'s `tests` module.

//! wallet_core: the shared Base (EVM/secp256k1) wallet logic behind
//! app-web (compiled to WASM) and app-mobile (compiled to a native
//! library, called through Dart FFI). This replaces the Stellar keypair/
//! signing logic each app previously carried independently
//! (app-web's src/utils/trovoSDK.ts, app-mobile's
//! lib/functions/trovo-sdk.dart) with one implementation, so key
//! generation, address derivation, and message signing are identical
//! (and reviewed once) across platforms instead of duplicated per-app.
//!
//! Layout: [`core`] holds the platform-agnostic logic as plain Rust
//! functions/errors; [`wasm`] and [`ffi`] are thin per-platform bindings
//! over it, compiled in only for their respective targets.

pub mod core;

#[cfg(target_arch = "wasm32")]
pub mod wasm;

#[cfg(not(target_arch = "wasm32"))]
pub mod ffi;

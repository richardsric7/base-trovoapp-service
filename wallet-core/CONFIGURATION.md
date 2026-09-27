# Configuration

**This crate has no runtime or feature-flag configuration surface.**
That's an honest finding, not an oversight in this document — `wallet-core`
is a small, fully static-behavior library: given the same inputs, its
functions always do the same thing, regardless of environment variables,
Cargo features, or build-time flags (beyond which *target* you compile
for, which is a build-command choice, not a configuration toggle — see
DEPLOYMENT.md).

## What was checked

- **`Cargo.toml`** has no `[features]` section at all. There is nothing
  to enable or disable with `cargo build --features <name>`.
- **`src/core.rs`, `src/wasm.rs`, `src/ffi.rs`, `src/lib.rs`** were
  grepped for both of Rust's common feature-flag patterns —
  `cfg(feature = ...)` and `cfg!(feature = ...)` — and for any use of
  `std::env::var`/`env::var`. Neither pattern appears anywhere in `src/`.

## The one thing that *does* vary: compile target

The only thing that changes this crate's compiled output is which
**target** you build for, controlled entirely by `cargo build`'s
`--target` flag (not a feature flag):

- Building for `wasm32-unknown-unknown` compiles in `src/wasm.rs` (the
  `wasm-bindgen` bindings) and excludes `src/ffi.rs`, per the
  `#[cfg(target_arch = "wasm32")]` / `#[cfg(not(target_arch = "wasm32"))]`
  guards in `src/lib.rs`.
- Building for any native target (the default, or an explicit target
  like `aarch64-linux-android`) does the reverse: `src/ffi.rs` is
  compiled in, `src/wasm.rs` is not.
- `crate-type = ["cdylib", "staticlib", "rlib"]` in `Cargo.toml` means
  every build produces all three output kinds together — which one a
  given consumer actually picks up (the `.wasm`, the `.so`/`.a`, or the
  `.rlib` used by `cargo test`) depends on their own build/link step, not
  on anything you configure here.

This is "target selection," covered fully in DEPLOYMENT.md, not
"configuration" in the feature-flag/env-var sense — there's no dial here
for a consumer to turn based on their own needs (e.g. no
"lite" vs "full" feature set, no debug-logging toggle, no alternate
crypto backend). If a consumer needs different behavior from this crate,
that's new code in `src/core.rs`, not a flag to flip.

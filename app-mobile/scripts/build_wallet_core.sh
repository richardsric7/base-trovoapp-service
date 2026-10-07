#!/usr/bin/env bash
# Builds wallet-core (../wallet-core, Rust) into the native library the app
# calls for keys and signatures, and puts it where the app loads it from.
#
#   scripts/build_wallet_core.sh android   # jniLibs for arm64-v8a, armeabi-v7a, x86_64
#   scripts/build_wallet_core.sh ios       # ios/WalletCore/WalletCore.xcframework (Mac only)
#   scripts/build_wallet_core.sh host      # this computer, for `flutter test`
#
# Run it again whenever wallet-core/src changes. The outputs are not
# committed (see .gitignore). Needs Rust (https://rustup.rs); Android also
# needs the Android NDK (ANDROID_NDK_HOME, or ANDROID_HOME with one NDK
# installed) and cargo-ndk; iOS needs Xcode.
set -euo pipefail

here="$(cd "$(dirname "$0")/.." && pwd)"          # app-mobile/
core="$here/../wallet-core"
[ -f "$core/Cargo.toml" ] || { echo "wallet-core not found at $core" >&2; exit 1; }

case "${1:-}" in
  android)
    if [ -z "${ANDROID_NDK_HOME:-}" ] && [ -n "${ANDROID_HOME:-}" ] && [ -d "$ANDROID_HOME/ndk" ]; then
      ANDROID_NDK_HOME="$ANDROID_HOME/ndk/$(ls "$ANDROID_HOME/ndk" | sort -V | tail -1)"
      export ANDROID_NDK_HOME
    fi
    [ -n "${ANDROID_NDK_HOME:-}" ] || { echo "Set ANDROID_NDK_HOME (Android Studio > SDK Manager > SDK Tools > NDK)" >&2; exit 1; }
    command -v cargo-ndk >/dev/null || cargo install cargo-ndk
    rustup target add aarch64-linux-android armv7-linux-androideabi x86_64-linux-android
    (cd "$core" && cargo ndk -t arm64-v8a -t armeabi-v7a -t x86_64 \
      -o "$here/android/app/src/main/jniLibs" build --release)
    ls -l "$here"/android/app/src/main/jniLibs/*/libwallet_core.so
    ;;
  ios)
    [ "$(uname)" = "Darwin" ] || { echo "The iOS library can only be built on a Mac" >&2; exit 1; }
    rustup target add aarch64-apple-ios aarch64-apple-ios-sim x86_64-apple-ios
    (cd "$core" && for t in aarch64-apple-ios aarch64-apple-ios-sim x86_64-apple-ios; do
      cargo build --release --target "$t"
    done)
    sim="$core/target/ios-sim-universal/libwallet_core.a"
    mkdir -p "$(dirname "$sim")"
    lipo -create "$core/target/aarch64-apple-ios-sim/release/libwallet_core.a" \
      "$core/target/x86_64-apple-ios/release/libwallet_core.a" -output "$sim"
    out="$here/ios/WalletCore/WalletCore.xcframework"
    rm -rf "$out"
    xcodebuild -create-xcframework \
      -library "$core/target/aarch64-apple-ios/release/libwallet_core.a" \
      -library "$sim" \
      -output "$out"
    echo "Built $out. Run 'cd ios && pod install' once so Xcode links it."
    ;;
  host)
    (cd "$core" && cargo build --release)
    echo "Built $core/target/release. Run the FFI test with:"
    echo "  LD_LIBRARY_PATH=$core/target/release flutter test test/wallet_core_ffi_test.dart"
    ;;
  *)
    echo "usage: $0 android|ios|host" >&2
    exit 2
    ;;
esac

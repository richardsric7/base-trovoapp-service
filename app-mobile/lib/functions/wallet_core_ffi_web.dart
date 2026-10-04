// Flutter Web stand-in for wallet_core_ffi_io.dart's WalletCoreFFI.
//
// dart:ffi (and the `ffi` package it depends on: Pointer, Allocator,
// NativeType, ...) cannot compile for the web target at all - those types
// only exist for native/AOT compilation, so any code that imports
// wallet_core_ffi_io.dart unconditionally fails `flutter build web` outright,
// not just at runtime.
//
// app-web already has a working web implementation of this same Rust
// wallet-core crate, compiled to WASM and called from
// app-web/src/utils/trovoSDK.ts. Wiring Flutter Web to that same WASM
// module via dart:js_interop - matching app-web's calling convention, and
// verifying signature/address output against native FFI for the same
// inputs - is real work with real correctness stakes for a crypto-signing
// path, out of scope for this dependency-bump pass. Every method here
// throws instead, so the rest of the app (navigation, the ~20 bumped
// packages, everything not on this screen) is still real signal from
// `flutter run -d chrome`, while wallet operations honestly report
// "not available" rather than silently doing nothing or crashing on a
// missing native symbol.
import 'dart:typed_data';

class WalletCoreException implements Exception {
  final String message;
  WalletCoreException(this.message);

  @override
  String toString() => 'WalletCoreException: $message';
}

Never _unsupported(String method) {
  throw UnsupportedError(
    'WalletCoreFFI.$method is not available on the Flutter Web target - '
    'wallet-core signing/key-derivation is native-only here today (see '
    'wallet_core_ffi_web.dart). Use a native build for wallet operations.',
  );
}

class WalletCoreFFI {
  static WalletCoreFFI? _instance;
  static WalletCoreFFI get instance => _instance ??= WalletCoreFFI._();

  WalletCoreFFI._();

  ({String privateKeyHex, String address}) generateKeypair() =>
      _unsupported('generateKeypair');

  ({String privateKeyHex, String address}) keypairFromPrivateKey(
    String privateKeyHex,
  ) => _unsupported('keypairFromPrivateKey');

  String addressFromPrivateKey(String privateKeyHex) =>
      _unsupported('addressFromPrivateKey');

  String signPersonal(String privateKeyHex, String message) =>
      _unsupported('signPersonal');

  String signPersonalBytes(String privateKeyHex, Uint8List message) =>
      _unsupported('signPersonalBytes');

  bool? verifyPersonal(String address, String message, String signatureB64) =>
      _unsupported('verifyPersonal');

  String recoverPersonalSigner(String message, String signatureB64) =>
      _unsupported('recoverPersonalSigner');

  String generateMnemonic() => _unsupported('generateMnemonic');

  ({String privateKeyHex, String address}) keypairFromMnemonic(
    String mnemonic,
    int index,
  ) => _unsupported('keypairFromMnemonic');

  String primarySafeAddress(String owner, String saltNonce) =>
      _unsupported('primarySafeAddress');
}

// Dart FFI bindings for the shared wallet-core Rust crate's C ABI
// (../../../wallet-core/src/ffi.rs) - the Base (EVM/secp256k1) key
// generation, address derivation, and EIP-191 signing logic also used by
// app-web (compiled to WASM there instead - see
// app-web/src/utils/trovoSDK.ts and /wallet-core's crate doc).
//
// These bindings are hand-written rather than generated (no
// cbindgen/flutter_rust_bridge codegen wired up yet), so they must be kept
// in sync by hand with src/ffi.rs's exported function signatures.
//
// Build step this depends on (not yet wired into this repo's mobile build
// - see /wallet-core/README.md "Building for app-mobile"): compile
// wallet-core for each target platform/ABI and place the resulting
// library where the loader below expects it - Android:
// android/app/src/main/jniLibs/<abi>/libwallet_core.so; iOS/macOS: link
// libwallet_core.a into the Xcode project (DynamicLibrary.process() then
// finds its symbols, since iOS/macOS Flutter builds statically link
// rather than loading a bundled .so/.dylib at runtime).
import 'dart:convert';
import 'dart:ffi';
import 'dart:io';
import 'dart:typed_data';

import 'package:ffi/ffi.dart';

typedef _WcFreeStringNative = Void Function(Pointer<Utf8>);
typedef _WcFreeStringDart = void Function(Pointer<Utf8>);

typedef _WcNoArgNative = Pointer<Utf8> Function();
typedef _WcNoArgDart = Pointer<Utf8> Function();

typedef _WcOneArgNative = Pointer<Utf8> Function(Pointer<Utf8>);
typedef _WcOneArgDart = Pointer<Utf8> Function(Pointer<Utf8>);

typedef _WcTwoArgNative = Pointer<Utf8> Function(Pointer<Utf8>, Pointer<Utf8>);
typedef _WcTwoArgDart = Pointer<Utf8> Function(Pointer<Utf8>, Pointer<Utf8>);

typedef _WcVerifyNative =
    Int32 Function(Pointer<Utf8>, Pointer<Utf8>, Pointer<Utf8>);
typedef _WcVerifyDart = int Function(Pointer<Utf8>, Pointer<Utf8>, Pointer<Utf8>);

typedef _WcMnemonicDeriveNative = Pointer<Utf8> Function(Pointer<Utf8>, Uint32);
typedef _WcMnemonicDeriveDart = Pointer<Utf8> Function(Pointer<Utf8>, int);

typedef _WcSignBytesNative =
    Pointer<Utf8> Function(Pointer<Utf8>, Pointer<Uint8>, Size);
typedef _WcSignBytesDart =
    Pointer<Utf8> Function(Pointer<Utf8>, Pointer<Uint8>, int);

DynamicLibrary _openWalletCoreLibrary() {
  if (Platform.isAndroid) {
    return DynamicLibrary.open('libwallet_core.so');
  }
  if (Platform.isIOS || Platform.isMacOS) {
    // Statically linked into the app binary on Apple platforms - its
    // symbols are already part of the running process.
    return DynamicLibrary.process();
  }
  if (Platform.isLinux) {
    return DynamicLibrary.open('libwallet_core.so');
  }
  if (Platform.isWindows) {
    return DynamicLibrary.open('wallet_core.dll');
  }
  throw UnsupportedError(
    'wallet-core has no known native library path for this platform',
  );
}

/// Thrown when a wallet-core FFI call returns `{"error": "..."}`.
class WalletCoreException implements Exception {
  final String message;
  WalletCoreException(this.message);

  @override
  String toString() => 'WalletCoreException: $message';
}

/// Low-level, typed wrapper over wallet-core's C ABI. Prefer
/// TrovoWalletSDK (trovo-sdk.dart) at call sites - this class exists so
/// the raw FFI plumbing (library loading, pointer lifetime, JSON
/// decoding) lives in exactly one place.
class WalletCoreFFI {
  static WalletCoreFFI? _instance;
  static WalletCoreFFI get instance => _instance ??= WalletCoreFFI._();

  late final _WcFreeStringDart _wcFreeString;
  late final _WcNoArgDart _wcGenerateKeypair;
  late final _WcOneArgDart _wcKeypairFromPrivateKey;
  late final _WcOneArgDart _wcAddressFromPrivateKey;
  late final _WcTwoArgDart _wcSignPersonal;
  late final _WcSignBytesDart _wcSignPersonalBytes;
  late final _WcVerifyDart _wcVerifyPersonal;
  late final _WcTwoArgDart _wcRecoverPersonalSigner;
  late final _WcNoArgDart _wcGenerateMnemonic;
  late final _WcMnemonicDeriveDart _wcKeypairFromMnemonic;

  WalletCoreFFI._() {
    final lib = _openWalletCoreLibrary();

    _wcFreeString = lib
        .lookup<NativeFunction<_WcFreeStringNative>>('wc_free_string')
        .asFunction();
    _wcGenerateKeypair = lib
        .lookup<NativeFunction<_WcNoArgNative>>('wc_generate_keypair')
        .asFunction();
    _wcKeypairFromPrivateKey = lib
        .lookup<NativeFunction<_WcOneArgNative>>('wc_keypair_from_private_key')
        .asFunction();
    _wcAddressFromPrivateKey = lib
        .lookup<NativeFunction<_WcOneArgNative>>('wc_address_from_private_key')
        .asFunction();
    _wcSignPersonal = lib
        .lookup<NativeFunction<_WcTwoArgNative>>('wc_sign_personal')
        .asFunction();
    _wcSignPersonalBytes = lib
        .lookup<NativeFunction<_WcSignBytesNative>>('wc_sign_personal_bytes')
        .asFunction();
    _wcVerifyPersonal = lib
        .lookup<NativeFunction<_WcVerifyNative>>('wc_verify_personal')
        .asFunction();
    _wcRecoverPersonalSigner = lib
        .lookup<NativeFunction<_WcTwoArgNative>>('wc_recover_personal_signer')
        .asFunction();
    _wcGenerateMnemonic = lib
        .lookup<NativeFunction<_WcNoArgNative>>('wc_generate_mnemonic')
        .asFunction();
    _wcKeypairFromMnemonic = lib
        .lookup<NativeFunction<_WcMnemonicDeriveNative>>(
          'wc_keypair_from_mnemonic',
        )
        .asFunction();
  }

  /// Calls a wallet-core function that returns a JSON C string, decodes
  /// it, frees the native buffer, and throws [WalletCoreException] if the
  /// JSON carries an `error` key.
  Map<String, dynamic> _callJson(Pointer<Utf8> Function() call) {
    final resultPtr = call();
    try {
      final decoded =
          jsonDecode(resultPtr.toDartString()) as Map<String, dynamic>;
      if (decoded.containsKey('error')) {
        throw WalletCoreException(decoded['error'] as String);
      }
      return decoded;
    } finally {
      _wcFreeString(resultPtr);
    }
  }

  ({String privateKeyHex, String address}) generateKeypair() {
    final json = _callJson(() => _wcGenerateKeypair());
    return (
      privateKeyHex: json['privateKeyHex'] as String,
      address: json['address'] as String,
    );
  }

  ({String privateKeyHex, String address}) keypairFromPrivateKey(
    String privateKeyHex,
  ) {
    final priv = privateKeyHex.toNativeUtf8();
    try {
      final json = _callJson(() => _wcKeypairFromPrivateKey(priv));
      return (
        privateKeyHex: json['privateKeyHex'] as String,
        address: json['address'] as String,
      );
    } finally {
      calloc.free(priv);
    }
  }

  String addressFromPrivateKey(String privateKeyHex) {
    final priv = privateKeyHex.toNativeUtf8();
    try {
      final json = _callJson(() => _wcAddressFromPrivateKey(priv));
      return json['address'] as String;
    } finally {
      calloc.free(priv);
    }
  }

  String signPersonal(String privateKeyHex, String message) {
    final priv = privateKeyHex.toNativeUtf8();
    final msg = message.toNativeUtf8();
    try {
      final json = _callJson(() => _wcSignPersonal(priv, msg));
      return json['signature'] as String;
    } finally {
      calloc.free(priv);
      calloc.free(msg);
    }
  }

  /// Signs raw bytes (e.g. a pre-computed transaction digest). Use this
  /// instead of [signPersonal] for binary data - a transaction digest is
  /// arbitrary bytes, not necessarily valid UTF-8, so round-tripping it
  /// through a Dart String (as [signPersonal] does) would corrupt it.
  String signPersonalBytes(String privateKeyHex, Uint8List message) {
    final priv = privateKeyHex.toNativeUtf8();
    final msgPtr = calloc<Uint8>(message.length == 0 ? 1 : message.length);
    final msgBytes = msgPtr.asTypedList(message.length);
    msgBytes.setAll(0, message);
    try {
      final json = _callJson(
        () => _wcSignPersonalBytes(priv, msgPtr, message.length),
      );
      return json['signature'] as String;
    } finally {
      calloc.free(priv);
      calloc.free(msgPtr);
    }
  }

  /// Returns true/false, or null if the arguments/signature themselves
  /// were malformed (distinct from a confirmed-invalid signature).
  bool? verifyPersonal(String address, String message, String signatureB64) {
    final addr = address.toNativeUtf8();
    final msg = message.toNativeUtf8();
    final sig = signatureB64.toNativeUtf8();
    try {
      final result = _wcVerifyPersonal(addr, msg, sig);
      if (result < 0) return null;
      return result == 1;
    } finally {
      calloc.free(addr);
      calloc.free(msg);
      calloc.free(sig);
    }
  }

  String recoverPersonalSigner(String message, String signatureB64) {
    final msg = message.toNativeUtf8();
    final sig = signatureB64.toNativeUtf8();
    try {
      final json = _callJson(() => _wcRecoverPersonalSigner(msg, sig));
      return json['address'] as String;
    } finally {
      calloc.free(msg);
      calloc.free(sig);
    }
  }

  String generateMnemonic() {
    final json = _callJson(() => _wcGenerateMnemonic());
    return json['mnemonic'] as String;
  }

  ({String privateKeyHex, String address}) keypairFromMnemonic(
    String mnemonic,
    int index,
  ) {
    final m = mnemonic.toNativeUtf8();
    try {
      final json = _callJson(() => _wcKeypairFromMnemonic(m, index));
      return (
        privateKeyHex: json['privateKeyHex'] as String,
        address: json['address'] as String,
      );
    } finally {
      calloc.free(m);
    }
  }
}

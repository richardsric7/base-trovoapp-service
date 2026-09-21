library trovo_app_sdk_flutter;

import 'dart:convert';

import 'wallet_core_ffi.dart';

// Base (EVM/secp256k1) wallet operations, backed by the shared wallet-core
// Rust crate (compiled to a native library, called through dart:ffi - see
// wallet_core_ffi.dart) instead of this file's own stellar_flutter_sdk
// calls. This is the same core the app-web React app calls through WASM,
// so key derivation and signing are implemented once and reviewed once,
// not duplicated per platform - see /wallet-core's crate doc.
class TrovoWalletSDK {
  Account createAccount() {
    // create a completely new and unique pair of keys.
    final kp = WalletCoreFFI.instance.generateKeypair();
    return Account(kp.address, kp.privateKeyHex);
  }

  String signHTTP({toSign, secretKey}) {
    try {
      return WalletCoreFFI.instance.signPersonal('$toSign', secretKey);
    } catch (e) {
      return '';
    }
  }

  importAccount(secretKey) {
    try {
      return WalletCoreFFI.instance.addressFromPrivateKey(secretKey);
    } catch (e) {
      return '';
    }
  }

  // signBase64Txn signs a transaction digest. On Base there is no XDR
  // envelope to parse and re-hash: transactionDigest is a base64-encoded
  // digest computed upstream (by the backend - see
  // internal/middleware/security_checks.go's SignBase64Txn), and this
  // function's job is only to sign it with EIP-191 personal_sign.
  // networkPassphrase is vestigial (kept only so call sites don't need to
  // change their argument count).
  signBase64Txn(
    String secretKey,
    String transactionDigest,
    String networkPassphrase,
  ) {
    try {
      final bytes = base64.decode(transactionDigest);
      return WalletCoreFFI.instance.signPersonalBytes(secretKey, bytes);
    } catch (e) {
      return '';
    }
  }

  // Derivation index 0 for both generation and retrieval, matching
  // app-web's trovoSDK.ts - the two platforms previously used different,
  // uncoordinated derivation schemes for "the same" mnemonic (this file
  // derived at index 1 via Stellar's SLIP-10-based Wallet.getKeyPair;
  // app-web hashed the raw BIP39 seed directly with no index concept at
  // all). wallet-core's BIP44 derivation (m/44'/60'/0'/0/{index})
  // harmonizes both onto one scheme, so a mnemonic generated on either
  // platform now derives the same address on both.
  static const int _mnemonicDerivationIndex = 0;

  Future<String> generateCredentialsFromPassPhrase() async {
    return WalletCoreFFI.instance.generateMnemonic();
  }

  Future<Account> retrieveCredentialsFromPassPhrase(String passPhrase) async {
    final kp = WalletCoreFFI.instance.keypairFromMnemonic(
      passPhrase,
      _mnemonicDerivationIndex,
    );
    return Account(kp.address, kp.privateKeyHex);
  }

  Account parseSecretKey(secretKey) {
    final kp = WalletCoreFFI.instance.keypairFromPrivateKey(secretKey);
    return Account(kp.address, kp.privateKeyHex);
  }
}

class Account {
  late String address;
  late String secretKey;

  Account(String p, String s) {
    this.address = p;
    this.secretKey = s;
  }

  @override
  String toString() {
    return "Address: " + this.address + " Secret-Key: " + this.secretKey;
  }
}

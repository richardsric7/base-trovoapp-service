library trovo_app_sdk_flutter;

import 'dart:convert';
import 'dart:typed_data';

import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart';

class TrovoWalletSDK {
  Account createAccount() {
    // create a completely new and unique pair of keys.
    KeyPair keyPair = KeyPair.random();

    Account account = Account(keyPair.accountId, keyPair.secretSeed);

    return account;
  }

  String signHTTP({toSign, secretKey}) {
    try {
      KeyPair keyPair = KeyPair.fromSecretSeed(secretKey);

      List<int> list = utf8.encode('$toSign');
      Uint8List bytes = Uint8List.fromList(list);

      // sign with the keypair
      var signedData = keyPair.sign(bytes);

      var signedBase64Str = base64.encode(signedData.toList());

      return signedBase64Str;
    } catch (e) {
      return '';
    }
  }

  importAccount(secretKey) {
    try {
      KeyPair keyPair = KeyPair.fromSecretSeed(secretKey);
      return keyPair.accountId;
    } catch (e) {
      return '';
    }
  }

  signBase64Txn(
    String secretKey,
    String transcationXDR,
    String networkPassphrase,
  ) {
    try {
      KeyPair keyPair = KeyPair.fromSecretSeed(secretKey);

      var txn = AbstractTransaction.fromEnvelopeXdrString(transcationXDR);

      Network network = new Network(networkPassphrase);

      var bytes = txn.hash(network);

      // sign with the keypair
      var signedData = keyPair.sign(bytes);

      var signedBase64Str = base64.encode(signedData.toList());

      return signedBase64Str;
    } catch (e) {
      return '';
    }
  }

  Future<String> generateCredentialsFromPassPhrase() async {
    String mnemonic = await Wallet.generate24WordsMnemonic();
    Wallet wallet = await Wallet.from(mnemonic);
    KeyPair keyPair = await wallet.getKeyPair(index: 1);
    Account(keyPair.accountId, keyPair.secretSeed);
    return mnemonic;
  }

  Future<Account> retrieveCredentialsFromPassPhrase(String passPhrase) async {
    Wallet wallet = await Wallet.from(passPhrase);
    KeyPair keyPair = await wallet.getKeyPair(index: 1);
    return Account(keyPair.accountId, keyPair.secretSeed);
  }

  Account parseSecretKey(secretKey) {
    KeyPair keyPair = KeyPair.fromSecretSeed(secretKey);
    return Account(keyPair.accountId, keyPair.secretSeed);
  }
}

class Account {
  late String publicKey;
  late String secretKey;

  Account(String p, String s) {
    this.publicKey = p;
    this.secretKey = s;
  }

  @override
  String toString() {
    return "Public-Key: " + this.publicKey + " Secret-Key: " + this.secretKey;
  }
}

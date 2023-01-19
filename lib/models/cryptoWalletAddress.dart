import 'package:charts_flutter/flutter.dart';

class CryptoWalletDepositAddresses {
  String? id;
  DateTime? createdAt;
  String? trovoWalletPublicKey;
  String? currency;
  double? depositAddress;
  String? network;
  String? qrCode;
  CryptoWalletDepositAddresses({
    this.id,
    this.createdAt,
    this.trovoWalletPublicKey,
    this.currency,
    this.depositAddress,
    this.network,
    this.qrCode,
  });

  CryptoWalletDepositAddresses deserializeJson(m) {
    return CryptoWalletDepositAddresses(
      id: m["id"],
      createdAt: m["createdAt"],
      trovoWalletPublicKey: m["trovoWalletPublicKey"],
      qrCode: m["qrCode"],
      currency: m["currency"],
      depositAddress: m["depositAddress"],
      network: m["network"],
    );
  }
}

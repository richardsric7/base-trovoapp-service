import 'package:trovo_app/models/crypto_wallet_address.dart';

class Asset {
  String? assetIssuer;
  String? assetCode;
  String? qrCode;
  String? imageUrl;
  double? amount;
  double? usdPrice;
  double? nativePrice;
  int userPreferredIndex =
      0; // index where user prefers it to appear on asset list
  Map? inTrade;
  List<CryptoWalletDepositAddress>? cryptoWalletDepositAddresses;
  int? assetClassId;

  Asset({
    this.assetCode,
    this.assetIssuer,
    this.qrCode,
    this.imageUrl,
    this.amount,
    this.usdPrice,
    this.inTrade,
    this.nativePrice,
    this.cryptoWalletDepositAddresses,
    this.assetClassId,
  });

  Asset deserializeJson(Map<String, dynamic> m) {
    return Asset(
      assetCode: m["assetCode"],
      assetIssuer: m["assetIssuer"],
      amount: double.parse(m["amount"]),
      qrCode: m["qrCode"],
      imageUrl: m["imageUrl"],
      usdPrice: double.parse(m["usdPrice"]),
      inTrade: m["inTrade"],
      nativePrice: double.parse(m["nativePrice"]),
      cryptoWalletDepositAddresses:
          deserializeDepositAddresses(m['cryptoWalletDepositAddresses']),
    );
  }

  List<CryptoWalletDepositAddress> deserializeDepositAddresses(m) {
    var addresses = <CryptoWalletDepositAddress>[];
    if (m != null) {
      for (var i = 0; i < m.length; i++) {
        addresses.add(CryptoWalletDepositAddress().deserializeJson(m[i]));
      }
    }
    return addresses;
  }
}

import 'package:trovo_wallet/models/cryptoWalletAddress.dart';

class Asset {
  String? assetIssuer;
  String? assetCode;
  String? qrCode;
  String? imageUrl;
  double? amount;
  double? usdPrice;
  double? nativePrice;
  Map? inTrade;
  // List<CryptoWalletDepositAddresses>? cryptoWalletDepositAddresses;
  Asset({
    this.assetCode,
    this.assetIssuer,
    this.qrCode,
    this.imageUrl,
    this.amount,
    this.usdPrice,
    this.inTrade,
    this.nativePrice,
    // this.cryptoWalletDepositAddresses,
  });

  Asset deserializeJson(Map<String, dynamic> m) {
    print('deserializing asset ====> $m');
    return Asset(
      assetCode: m["assetCode"],
      assetIssuer: m["assetIssuer"],
      amount: double.parse(m["amount"]),
      qrCode: m["qrCode"],
      imageUrl: m["imageUrl"],
      usdPrice: double.parse(m["usdPrice"]),
      inTrade: m["inTrade"],
      nativePrice: double.parse(m["nativePrice"]),
      // cryptoWalletDepositAddresses: CryptoWalletDepositAddresses()
      //     .deserializeJson(m['cryptoWalletDepositAddresses']),
    );
  }
}

class CryptoWalletDepositAddress {
  String? id;
  DateTime? createdAt;
  String? trovoWalletAddress;
  String? currency;
  String? depositAddress;
  String? network;
  String? qrCode;
  CryptoWalletDepositAddress({
    this.id,
    this.createdAt,
    this.trovoWalletAddress,
    this.currency,
    this.depositAddress,
    this.network,
    this.qrCode,
  });

  CryptoWalletDepositAddress deserializeJson(m) {
    return CryptoWalletDepositAddress(
      id: m["id"],
      createdAt: DateTime.tryParse(m["createdAt"]),
      trovoWalletAddress: m["trovoWalletAddress"],
      qrCode: m["qrCode"],
      currency: m["currency"],
      depositAddress: m["depositAddress"],
      network: m["network"],
    );
  }
}

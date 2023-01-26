class CryptoWalletDepositAddress {
  String? id;
  DateTime? createdAt;
  String? trovoWalletPublicKey;
  String? currency;
  String? depositAddress;
  String? network;
  String? qrCode;
  CryptoWalletDepositAddress({
    this.id,
    this.createdAt,
    this.trovoWalletPublicKey,
    this.currency,
    this.depositAddress,
    this.network,
    this.qrCode,
  });

  CryptoWalletDepositAddress deserializeJson(m) {
    return CryptoWalletDepositAddress(
      id: m["id"],
      createdAt: DateTime.tryParse(m["createdAt"]),
      trovoWalletPublicKey: m["trovoWalletPublicKey"],
      qrCode: m["qrCode"],
      currency: m["currency"],
      depositAddress: m["depositAddress"],
      network: m["network"],
    );
  }
}

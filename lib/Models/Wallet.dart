class Wallet {
  DateTime? createdAt;
  String? publicKey;
  String? secretKey;
  String? tag;
  String? description;
  String? alias;
  String? signer;
  String? userId;
  int? managedAccessEnabled;
  int? primaryWallet;

  Wallet({
    this.createdAt,
    this.publicKey,
    this.secretKey,
    this.tag,
    this.description,
    this.alias,
    this.signer,
    this.userId,
    this.managedAccessEnabled,
    this.primaryWallet,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "createdAt": createdAt!.toIso8601String(),
      "publicKey": publicKey,
      "secretKey": secretKey,
      "tag": tag,
      "description": description,
      "alias": alias,
      "signer": signer,
      "userId": userId,
      "managedAccessEnabled": managedAccessEnabled,
      "primaryWallet": primaryWallet
    };
  }

  Wallet deserializeJson(Map<String, dynamic> m) {
    return Wallet(
      createdAt: DateTime.parse(m["createdAt"]),
      publicKey: m["publicKey"],
      secretKey: m["secretKey"],
      tag: m["tag"],
      description: m["description"],
      alias: m["alias"],
      signer: m["signer"],
      userId: m["userId"],
      managedAccessEnabled: m["managedAccessEnabled"],
      primaryWallet: m["primaryWallet"],
    );
  }
}

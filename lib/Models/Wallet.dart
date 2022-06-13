class Wallet {
  DateTime createdAt;
  String publicKey;
  String tag;
  String description;
  String alias;
  String signer;
  String userId;
  int managedAccessEnabled;
  int primaryWallet;

  Wallet({
    required this.createdAt,
    required this.publicKey,
    required this.tag,
    required this.description,
    required this.alias,
    required this.signer,
    required this.userId,
    required this.managedAccessEnabled,
    required this.primaryWallet,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "createdAt": createdAt.toIso8601String(),
      "publicKey": publicKey,
      "tag": tag,
      "description": description,
      "alias": alias,
      "signer": signer,
      "userId": userId,
      "managedAccessEnabled": managedAccessEnabled,
      "primaryWallet": primaryWallet
    };
  }

  deserializeJson(Map<String, dynamic> m) {
    return Wallet(
      createdAt: DateTime.parse(m["createdAt"]),
      publicKey: m["publicKey"],
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

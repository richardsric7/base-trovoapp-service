class PatronMembership {
  String? username;
  String? patronPackageId;
  String? patronTierId;
  int? price;
  DateTime? validTill;
  String? logo;
  PatronMembership({
    this.username,
    this.patronPackageId,
    this.patronTierId,
    this.price,
    this.validTill,
    this.logo,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "validTill": validTill!.toIso8601String(),
      "patronPackageId": patronPackageId,
      "patronTierId": patronTierId,
      "price": price,
      "username": username,
    };
  }

  PatronMembership deserializeJson(Map<String, dynamic> m) {
    print('===> deserializing patronMembership ${m}');
    return PatronMembership(
      validTill: DateTime.parse(m["validTill"]),
      patronPackageId: m["patronPackageId"],
      patronTierId: m["patronTierId"],
      price: m["price"],
      username: m["username"],
    );
  }
}

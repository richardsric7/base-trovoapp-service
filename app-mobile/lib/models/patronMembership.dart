class PatronMembership {
  String? username;
  String? patronPackageId;
  String? patronTierId;
  int? price;
  DateTime? validTill;
  // String? logo;
  int? id;
  PatronMembership({
    this.username,
    this.patronPackageId,
    this.patronTierId,
    this.price,
    this.validTill,
    // this.logo,
    this.id,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "validTill": validTill!.toIso8601String(),
      "patronPackageId": patronPackageId,
      "patronTierId": patronTierId,
      "price": price,
      "username": username,
      "id": id,
    };
  }

  String getLogo() {
    return patronPackageId!.toLowerCase() == 'gold'
        ? "assets/images/gold.png"
        : patronPackageId!.toLowerCase() == 'diamond'
        ? "assets/images/diamond.png"
        : "assets/images/platinum.png";
  }

  PatronMembership deserializeJson(Map<String, dynamic> m) {
    return PatronMembership(
      validTill: DateTime.parse(m["validTill"]),
      patronPackageId: m["patronPackageId"],
      patronTierId: m["patronTierId"],
      price: m["price"],
      username: m["username"],
      id: m["id"],
    );
  }
}

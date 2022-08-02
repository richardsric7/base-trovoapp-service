class TransactionInfo {
  DateTime? transactionDate;
  String? transactionType;
  String? from;
  String? fromPublicKey;
  String? to;
  String? toPublicKey;
  String? memo;
  String? assetCode;
  String? assetIssuer;
  double? amount;
  String? transactionId;

  TransactionInfo({
    this.transactionDate,
    this.transactionType,
    this.from,
    this.fromPublicKey,
    this.to,
    this.toPublicKey,
    this.memo,
    this.assetCode,
    this.assetIssuer,
    this.amount,
    this.transactionId,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "transactionDate": transactionDate!.toIso8601String(),
      "transactionType": transactionType,
      "from": from,
      "fromPublicKey": fromPublicKey,
      "to": to,
      "toPublicKey": toPublicKey,
      "memo": memo,
      "assetCode": assetCode,
      "assetIssuer": assetIssuer,
      "amount": amount,
      "transactionId": transactionId
    };
  }

  TransactionInfo deserializeJson(Map<String, dynamic> m) {
    return TransactionInfo(
      transactionDate: DateTime.parse(m["transactionDate"]),
      transactionType: m["transactionType"],
      from: m["from"],
      fromPublicKey: m["fromPublicKey"],
      to: m["to"],
      toPublicKey: m["toPublicKey"],
      memo: m["memo"],
      assetCode: m["assetCode"],
      assetIssuer: m["assetIssuer"],
      amount: double.tryParse(m["amount"]),
      transactionId: m["transactionId"],
    );
  }
}

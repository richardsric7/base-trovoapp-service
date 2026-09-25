import 'package:trovo_app/bottom_bar/bottom_pages/payment_history.dart';

class TransactionInfo {
  DateTime? transactionDate;
  String? transactionType;
  TransactionDirection? transactionDirection;
  String? from;
  String? fromAddress;
  String? to;
  String? toAddress;
  String? memo;
  String? assetCode;
  String? contractAddress;
  double? amount;
  String? transactionId;

  TransactionInfo({
    this.transactionDate,
    this.transactionType,
    this.transactionDirection,
    this.from,
    this.fromAddress,
    this.to,
    this.toAddress,
    this.memo,
    this.assetCode,
    this.contractAddress,
    this.amount,
    this.transactionId,
  });

  toJSONEncodable() {
    return <String, dynamic>{
      "transactionDate": transactionDate!.toIso8601String(),
      "transactionType": transactionType,
      "from": from,
      "fromAddress": fromAddress,
      "to": to,
      "toAddress": toAddress,
      "memo": memo,
      "assetCode": assetCode,
      "contractAddress": contractAddress,
      "amount": amount,
      "transactionId": transactionId,
    };
  }

  TransactionInfo deserializeJson(Map<String, dynamic> m) {
    return TransactionInfo(
      transactionDate: DateTime.parse(m["transactionDate"]),
      transactionType: m["transactionType"],
      from: m["from"],
      fromAddress: m["fromAddress"],
      to: m["to"],
      toAddress: m["toAddress"],
      memo: m["memo"],
      assetCode: m["assetCode"],
      contractAddress: m["contractAddress"],
      amount: double.tryParse(m["amount"]),
      transactionId: m["transactionId"],
    );
  }
}

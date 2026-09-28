import 'package:trovo_app/bottom_bar/bottom_pages/payment_history.dart';

const String kNetworkBase = 'base';

class TransactionInfo {
  DateTime? transactionDate;
  String? transactionType;
  TransactionDirection? transactionDirection;
  String? from;
  String? fromAddress;
  String? to;
  String? toAddress;
  String? memo;

  // Source: funds leaving fromAddress. Destination: funds arriving at
  // toAddress. For a plain payment the two sides are identical; a row
  // where they differ is, by definition, a swap (see isSwap below).
  String? sourceNetwork;
  String? sourceContractAddress;
  String? sourceAssetCode;
  double? sourceAmount;

  String? destinationNetwork;
  String? destinationContractAddress;
  String? destinationAssetCode;
  double? destinationAmount;

  String? transactionId;

  // assetCode/contractAddress/amount alias the destination side, so every
  // screen written against this model's original (single-asset) shape
  // keeps working unchanged - they were always describing what arrived at
  // toAddress, which is exactly what "destination" means.
  String? get assetCode => destinationAssetCode;
  String? get contractAddress => destinationContractAddress;
  double? get amount => destinationAmount;

  bool get isSwap =>
      sourceAssetCode != destinationAssetCode ||
      (sourceNetwork ?? kNetworkBase) != (destinationNetwork ?? kNetworkBase);

  TransactionInfo({
    this.transactionDate,
    this.transactionType,
    this.transactionDirection,
    this.from,
    this.fromAddress,
    this.to,
    this.toAddress,
    this.memo,
    this.sourceNetwork,
    this.sourceContractAddress,
    this.sourceAssetCode,
    this.sourceAmount,
    this.destinationNetwork,
    this.destinationContractAddress,
    this.destinationAssetCode,
    this.destinationAmount,
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
      "sourceNetwork": sourceNetwork,
      "sourceContractAddress": sourceContractAddress,
      "sourceAssetCode": sourceAssetCode,
      "sourceAmount": sourceAmount,
      "destinationNetwork": destinationNetwork,
      "destinationContractAddress": destinationContractAddress,
      "destinationAssetCode": destinationAssetCode,
      "destinationAmount": destinationAmount,
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
      sourceNetwork: m["sourceNetwork"],
      sourceContractAddress: m["sourceContractAddress"],
      sourceAssetCode: m["sourceAssetCode"],
      sourceAmount: _parseAmount(m["sourceAmount"]),
      destinationNetwork: m["destinationNetwork"],
      destinationContractAddress: m["destinationContractAddress"],
      destinationAssetCode: m["destinationAssetCode"],
      destinationAmount: _parseAmount(m["destinationAmount"]),
      transactionId: m["transactionId"],
    );
  }
}

double? _parseAmount(dynamic value) {
  if (value == null) return null;
  return double.tryParse(value.toString());
}

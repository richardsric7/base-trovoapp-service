class DepositTransactionModel {
  late DateTime createdAt;
  late DateTime updatedAt;
  late String currency;
  late String trovoWalletPublicKey;
  late String transactionId;
  late String fromAddress;
  late String toAddress;
  late bool isCompleted;
  late bool isVerified;
  late bool isValid;
  late double amount;

  DepositTransactionModel({
    required this.createdAt,
    required this.updatedAt,
    required this.currency,
    required this.trovoWalletPublicKey,
    required this.transactionId,
    required this.fromAddress,
    required this.toAddress,
    required this.isCompleted,
    required this.isVerified,
    required this.isValid,
    required this.amount,
  });

  DepositTransactionModel.deserializeJson(m)
      : this(
          createdAt: DateTime.parse(m["createdAt"]),
          updatedAt: DateTime.parse(m["updatedAt"]),
          currency: m["currency"],
          trovoWalletPublicKey: m["trovoWalletPublicKey"],
          transactionId: m["txId"],
          fromAddress: m["fromAddress"],
          toAddress: m["toAddress"],
          isCompleted: m["isCompleted"],
          isVerified: m["isVerified"],
          isValid: m["isValid"],
          amount: double.parse(m["amount"]),
        );
}

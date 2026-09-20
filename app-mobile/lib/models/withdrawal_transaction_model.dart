class WithdrawalTransactionModel {
  late DateTime createdAt;
  late String id;
  late String currency;
  late String walletAddress;
  late String walletAlias;
  late String userId;
  late String transactionId;
  late String withdrawalAddress;
  late String withdrawalNetwork;
  late String withdrawalStatus;
  late String withdrawalMemo;
  late double withdrawalNetworkFee;
  late double withdrawalServiceFee;
  late double amountSubmitted;
  late double amountToWithdraw;
  WithdrawalTransactionModel({
    required this.createdAt,
    required this.id,
    required this.currency,
    required this.walletAddress,
    required this.walletAlias,
    required this.userId,
    required this.transactionId,
    required this.withdrawalAddress,
    required this.withdrawalNetwork,
    required this.withdrawalStatus,
    required this.withdrawalMemo,
    required this.withdrawalNetworkFee,
    required this.withdrawalServiceFee,
    required this.amountSubmitted,
    required this.amountToWithdraw,
  });

  WithdrawalTransactionModel.deserializeJson(Map<String, dynamic> m)
    : this(
        createdAt: DateTime.parse(m["createdAt"]),
        id: m["id"],
        currency: m["currency"],
        walletAddress: m["walletAddress"],
        walletAlias: m["walletAlias"],
        userId: m["userId"],
        transactionId: m["transactionId"],
        withdrawalAddress: m["withdrawalAddress"],
        withdrawalNetwork: m["withdrawalNetwork"],
        withdrawalStatus: m["withdrawalStatus"],
        withdrawalMemo: m["withdrawalMemo"],
        withdrawalNetworkFee: m["withdrawalNetworkFee"],
        withdrawalServiceFee: m["withdrawalServiceFee"],
        amountSubmitted: m["amountSubmitted"],
        amountToWithdraw: m["amountToWithdraw"],
      );
}

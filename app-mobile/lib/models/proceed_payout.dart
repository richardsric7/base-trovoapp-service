// ProceedPayoutReceipt is a proceeds payout (dividend or interest) to one
// of the user's wallets, from GET /v1/tokenization/payouts: paid ones carry
// their transaction; scheduled ones (PENDING / QUEUED) are in a payout whose
// holder schedule is locked and waiting to be paid.
class ProceedPayoutReceipt {
  final String id;
  final int proceedPayoutId;
  final String tokenizedAssetId;
  final String assetCode;
  final String assetName;
  final String tokenContractAddress;
  final String payoutAssetCode;
  final String amount; // in the payout token
  final String tokensHeld; // at the snapshot
  final String amountPerToken;
  final String beneficiaryAddress;
  final String walletAlias;
  final String status; // PENDING, QUEUED, PAID, FAILED
  final String txHash;
  final DateTime? paidAt;
  final DateTime? scheduledAt;

  ProceedPayoutReceipt({
    required this.id,
    required this.proceedPayoutId,
    required this.tokenizedAssetId,
    required this.assetCode,
    required this.assetName,
    required this.tokenContractAddress,
    required this.payoutAssetCode,
    required this.amount,
    required this.tokensHeld,
    required this.amountPerToken,
    required this.beneficiaryAddress,
    required this.walletAlias,
    required this.status,
    required this.txHash,
    this.paidAt,
    this.scheduledAt,
  });

  bool get paid => status == 'PAID';
  bool get scheduled => status == 'PENDING' || status == 'QUEUED';
  double get amountValue => double.tryParse(amount) ?? 0;

  // the date to show: when it was paid, or when it was scheduled
  DateTime? get date => paidAt ?? scheduledAt;

  static DateTime? _date(dynamic v) {
    if (v == null) return null;
    final d = DateTime.tryParse('$v');
    if (d == null || d.year < 2000) return null;
    return d.toLocal();
  }

  factory ProceedPayoutReceipt.fromMap(Map m) => ProceedPayoutReceipt(
    id: '${m['id'] ?? ''}',
    proceedPayoutId: (m['proceedPayoutId'] as num?)?.toInt() ?? 0,
    tokenizedAssetId: '${m['tokenizedAssetId'] ?? ''}',
    assetCode: '${m['assetCode'] ?? ''}',
    assetName: '${m['assetName'] ?? ''}',
    tokenContractAddress: '${m['tokenContractAddress'] ?? ''}',
    payoutAssetCode: '${m['payoutAssetCode'] ?? ''}',
    amount: '${m['amount'] ?? '0'}',
    tokensHeld: '${m['tokensHeld'] ?? '0'}',
    amountPerToken: '${m['amountPerToken'] ?? ''}',
    beneficiaryAddress: '${m['beneficiaryAddress'] ?? ''}',
    walletAlias: '${m['walletAlias'] ?? ''}',
    status: '${m['status'] ?? ''}',
    txHash: '${m['txHash'] ?? ''}',
    paidAt: _date(m['paidAt']),
    scheduledAt: _date(m['scheduledAt']),
  );

  Map<String, dynamic> toMap() => {
    'id': id,
    'proceedPayoutId': proceedPayoutId,
    'tokenizedAssetId': tokenizedAssetId,
    'assetCode': assetCode,
    'assetName': assetName,
    'tokenContractAddress': tokenContractAddress,
    'payoutAssetCode': payoutAssetCode,
    'amount': amount,
    'tokensHeld': tokensHeld,
    'amountPerToken': amountPerToken,
    'beneficiaryAddress': beneficiaryAddress,
    'walletAlias': walletAlias,
    'status': status,
    'txHash': txHash,
    'paidAt': paidAt?.toIso8601String(),
    'scheduledAt': scheduledAt?.toIso8601String(),
  };
}

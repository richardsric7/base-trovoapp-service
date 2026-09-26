class P2PRefund {
  String? id;
  String? orderId;
  String? sender;
  String? token;
  String? contractAddress;
  String? amount;
  String? reason;
  bool? claimed;
  DateTime? claimedAt;
  String? transactionHash;
  DateTime? createdAt;

  P2PRefund({
    this.id,
    this.orderId,
    this.sender,
    this.token,
    this.contractAddress,
    this.amount,
    this.reason,
    this.claimed,
    this.claimedAt,
    this.transactionHash,
    this.createdAt,
  });

  P2PRefund deserializeJson(Map<String, dynamic> m) {
    return P2PRefund(
      id: m['id'],
      orderId: m['orderId'],
      sender: m['sender'],
      token: m['token'],
      contractAddress: m['contractAddress'],
      amount: m['amount'],
      reason: m['reason'],
      claimed: m['claimed'] == true,
      claimedAt: m['claimedAt'] != null ? DateTime.tryParse(m['claimedAt']) : null,
      transactionHash: m['transactionHash'],
      createdAt: m['createdAt'] != null ? DateTime.tryParse(m['createdAt']) : null,
    );
  }

  static List<P2PRefund> deserializeList(List? m) {
    var list = <P2PRefund>[];
    if (m != null) {
      for (var i = 0; i < m.length; i++) {
        list.add(P2PRefund().deserializeJson(m[i]));
      }
    }
    return list;
  }
}

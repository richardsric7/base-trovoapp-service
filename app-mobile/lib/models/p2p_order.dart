import 'p2p_offer.dart';

class P2POrder {
  String? id;
  String? offerId;
  String? customerUserId;
  String? customerUsername;
  String? merchantUserId;
  String? merchantUsername;
  String? offerType;
  String? asset;
  String? assetContractAddress;
  P2PPaymentMethod? paymentMethodSnapshot;
  String? country;
  String? countryCode;
  String? currency;
  String? price;
  String? specifiedAssetAmount;
  String? paymentAmount;

  String? buyerTotalFees;
  String? buyerTotalVat;
  String? buyerTotalCharges;
  String? buyerNetAssetAmount;

  String? sellerTotalFees;
  String? sellerTotalVat;
  String? sellerTotalCharges;
  String? sellerEscrowAssetAmount;

  String? combinedPlatformFee;
  String? combinedRegulatoryFee;
  String? combinedVat;

  String? escrowDepositShortlink;
  String? escrowDepositQrCode;
  String? escrowDepositTransactionHash;
  String? assetReleaseTransactionHash;

  String? escrowDepositStatus;
  String? expectedEscrowAmount;
  String? depositedEscrowAmount;
  String? refundableAmount;

  String? assetDepositor;
  String? assetRecipient;
  String? fiatPayer;
  String? fiatRecipient;

  String? orderStatus;
  bool? isDisputed;
  DateTime? expiresAt;
  DateTime? completedAt;
  DateTime? paymentConfirmedAt;
  DateTime? assetReleasedAt;
  DateTime? createdAt;
  DateTime? updatedAt;

  P2POrder({
    this.id,
    this.offerId,
    this.customerUserId,
    this.customerUsername,
    this.merchantUserId,
    this.merchantUsername,
    this.offerType,
    this.asset,
    this.assetContractAddress,
    this.paymentMethodSnapshot,
    this.country,
    this.countryCode,
    this.currency,
    this.price,
    this.specifiedAssetAmount,
    this.paymentAmount,
    this.buyerTotalFees,
    this.buyerTotalVat,
    this.buyerTotalCharges,
    this.buyerNetAssetAmount,
    this.sellerTotalFees,
    this.sellerTotalVat,
    this.sellerTotalCharges,
    this.sellerEscrowAssetAmount,
    this.combinedPlatformFee,
    this.combinedRegulatoryFee,
    this.combinedVat,
    this.escrowDepositShortlink,
    this.escrowDepositQrCode,
    this.escrowDepositTransactionHash,
    this.assetReleaseTransactionHash,
    this.escrowDepositStatus,
    this.expectedEscrowAmount,
    this.depositedEscrowAmount,
    this.refundableAmount,
    this.assetDepositor,
    this.assetRecipient,
    this.fiatPayer,
    this.fiatRecipient,
    this.orderStatus,
    this.isDisputed,
    this.expiresAt,
    this.completedAt,
    this.paymentConfirmedAt,
    this.assetReleasedAt,
    this.createdAt,
    this.updatedAt,
  });

  static const statusAwaitingApproval = 'AWAITING_APPROVAL';
  static const statusAwaitingEscrowDeposit = 'AWAITING_ESCROW_DEPOSIT';
  static const statusAwaitingPayment = 'AWAITING_PAYMENT';
  static const statusAwaitingPaymentConfirmation =
      'AWAITING_PAYMENT_CONFIRMATION';
  static const statusCompleted = 'COMPLETED';
  static const statusRejected = 'REJECTED';
  static const statusCancelled = 'CANCELLED';
  static const statusExpired = 'EXPIRED';

  // The mobile client's cached profile (UserInfo) has no raw user-id field,
  // only `username` - so role checks compare by username, then mirror the
  // backend's fixed offerType-based role mapping (Plan Section 11) rather
  // than comparing against the opaque customerUserId/merchantUserId/
  // assetDepositor/etc. ids the API returns (those are meaningful to the
  // backend, not to a client that only knows its own username).
  bool isCustomer(String myUsername) => customerUsername == myUsername;
  bool isMerchant(String myUsername) => merchantUsername == myUsername;

  bool isAssetDepositor(String myUsername) {
    final iAmCustomer = isCustomer(myUsername);
    return offerType == 'SELL' ? !iAmCustomer : iAmCustomer;
  }

  bool isAssetRecipient(String myUsername) => !isAssetDepositor(myUsername);

  bool isFiatPayer(String myUsername) {
    final iAmCustomer = isCustomer(myUsername);
    return offerType == 'SELL' ? iAmCustomer : !iAmCustomer;
  }

  bool isFiatRecipient(String myUsername) => !isFiatPayer(myUsername);

  P2POrder deserializeJson(Map<String, dynamic> m) {
    DateTime? parseDate(dynamic v) => v != null ? DateTime.tryParse(v) : null;
    return P2POrder(
      id: m['id'],
      offerId: m['offerId'],
      customerUserId: m['customerUserId'],
      customerUsername: m['customerUsername'],
      merchantUserId: m['merchantUserId'],
      merchantUsername: m['merchantUsername'],
      offerType: m['offerType'],
      asset: m['asset'],
      assetContractAddress: m['assetContractAddress'],
      paymentMethodSnapshot: P2PPaymentMethod().deserializeJson(
        m['paymentMethodSnapshot'],
      ),
      country: m['country'],
      countryCode: m['countryCode'],
      currency: m['currency'],
      price: m['price'],
      specifiedAssetAmount: m['specifiedAssetAmount'],
      paymentAmount: m['paymentAmount'],
      buyerTotalFees: m['buyerTotalFees'],
      buyerTotalVat: m['buyerTotalVat'],
      buyerTotalCharges: m['buyerTotalCharges'],
      buyerNetAssetAmount: m['buyerNetAssetAmount'],
      sellerTotalFees: m['sellerTotalFees'],
      sellerTotalVat: m['sellerTotalVat'],
      sellerTotalCharges: m['sellerTotalCharges'],
      sellerEscrowAssetAmount: m['sellerEscrowAssetAmount'],
      combinedPlatformFee: m['combinedPlatformFee'],
      combinedRegulatoryFee: m['combinedRegulatoryFee'],
      combinedVat: m['combinedVat'],
      escrowDepositShortlink: m['escrowDepositShortlink'],
      escrowDepositQrCode: m['escrowDepositQrCode'],
      escrowDepositTransactionHash: m['escrowDepositTransactionHash'],
      assetReleaseTransactionHash: m['assetReleaseTransactionHash'],
      escrowDepositStatus: m['escrowDepositStatus'],
      expectedEscrowAmount: m['expectedEscrowAmount'],
      depositedEscrowAmount: m['depositedEscrowAmount'],
      refundableAmount: m['refundableAmount'],
      assetDepositor: m['assetDepositor'],
      assetRecipient: m['assetRecipient'],
      fiatPayer: m['fiatPayer'],
      fiatRecipient: m['fiatRecipient'],
      orderStatus: m['orderStatus'],
      isDisputed: m['isDisputed'] == true,
      expiresAt: parseDate(m['expiresAt']),
      completedAt: parseDate(m['completedAt']),
      paymentConfirmedAt: parseDate(m['paymentConfirmedAt']),
      assetReleasedAt: parseDate(m['assetReleasedAt']),
      createdAt: parseDate(m['createdAt']),
      updatedAt: parseDate(m['updatedAt']),
    );
  }

  static List<P2POrder> deserializeList(List? m) {
    var list = <P2POrder>[];
    if (m != null) {
      for (var i = 0; i < m.length; i++) {
        list.add(P2POrder().deserializeJson(m[i]));
      }
    }
    return list;
  }
}

class P2PDispute {
  String? id;
  String? orderId;
  String? openedBy;
  DateTime? openedAt;
  String? subject;
  String? description;
  List<String>? evidence;
  String? status;
  String? resolution;
  DateTime? resolvedAt;
  String? resolvedBy;

  P2PDispute({
    this.id,
    this.orderId,
    this.openedBy,
    this.openedAt,
    this.subject,
    this.description,
    this.evidence,
    this.status,
    this.resolution,
    this.resolvedAt,
    this.resolvedBy,
  });

  P2PDispute deserializeJson(Map<String, dynamic> m) {
    return P2PDispute(
      id: m['id'],
      orderId: m['orderId'],
      openedBy: m['openedBy'],
      openedAt: m['openedAt'] != null ? DateTime.tryParse(m['openedAt']) : null,
      subject: m['subject'],
      description: m['description'],
      evidence: (m['evidence'] as List?)?.map((e) => e.toString()).toList(),
      status: m['status'],
      resolution: m['resolution'],
      resolvedAt: m['resolvedAt'] != null ? DateTime.tryParse(m['resolvedAt']) : null,
      resolvedBy: m['resolvedBy'],
    );
  }
}

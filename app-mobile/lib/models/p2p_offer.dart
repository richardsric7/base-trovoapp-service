class P2PPaymentMethod {
  String? paymentChannel;
  String? provider;
  String? account;

  P2PPaymentMethod({this.paymentChannel, this.provider, this.account});

  P2PPaymentMethod deserializeJson(Map<String, dynamic>? m) {
    if (m == null) return this;
    paymentChannel = m['paymentChannel'];
    provider = m['provider'];
    account = m['account'];
    return this;
  }

  Map<String, dynamic> toJson() => {
    'paymentChannel': paymentChannel ?? '',
    'provider': provider ?? '',
    'account': account ?? '',
  };
}

class P2POffer {
  String? id;
  String? merchantUsername;
  String? merchantUserId;
  String? offerType; // BUY or SELL
  String? asset;
  String? contractAddress;
  P2PPaymentMethod? paymentMethod;
  String? country;
  String? countryCode;
  String? currency;
  String? priceType;
  String? price;
  String? priceMargin;
  String? minOrderAmount;
  String? maxOrderAmount;
  String? availableLiquidity;
  String? reservedLiquidity;
  String? remark;
  String? availabilityStatus; // ONLINE or OFFLINE
  String? status;
  int? version;
  DateTime? createdAt;
  DateTime? updatedAt;

  P2POffer({
    this.id,
    this.merchantUsername,
    this.merchantUserId,
    this.offerType,
    this.asset,
    this.contractAddress,
    this.paymentMethod,
    this.country,
    this.countryCode,
    this.currency,
    this.priceType,
    this.price,
    this.priceMargin,
    this.minOrderAmount,
    this.maxOrderAmount,
    this.availableLiquidity,
    this.reservedLiquidity,
    this.remark,
    this.availabilityStatus,
    this.status,
    this.version,
    this.createdAt,
    this.updatedAt,
  });

  bool get isOnline => availabilityStatus == 'ONLINE';
  bool get isBuy => offerType == 'BUY';

  P2POffer deserializeJson(Map<String, dynamic> m) {
    return P2POffer(
      id: m['id'],
      merchantUsername: m['merchantUsername'],
      merchantUserId: m['merchantUserId'],
      offerType: m['offerType'],
      asset: m['asset'],
      contractAddress: m['contractAddress'],
      paymentMethod: P2PPaymentMethod().deserializeJson(m['paymentMethod']),
      country: m['country'],
      countryCode: m['countryCode'],
      currency: m['currency'],
      priceType: m['priceType'],
      price: m['price'],
      priceMargin: m['priceMargin'],
      minOrderAmount: m['minOrderAmount'],
      maxOrderAmount: m['maxOrderAmount'],
      availableLiquidity: m['availableLiquidity'],
      reservedLiquidity: m['reservedLiquidity'],
      remark: m['remark'],
      availabilityStatus: m['availabilityStatus'],
      status: m['status'],
      version: m['version'],
      createdAt: m['createdAt'] != null ? DateTime.tryParse(m['createdAt']) : null,
      updatedAt: m['updatedAt'] != null ? DateTime.tryParse(m['updatedAt']) : null,
    );
  }

  static List<P2POffer> deserializeList(List? m) {
    var list = <P2POffer>[];
    if (m != null) {
      for (var i = 0; i < m.length; i++) {
        list.add(P2POffer().deserializeJson(m[i]));
      }
    }
    return list;
  }
}

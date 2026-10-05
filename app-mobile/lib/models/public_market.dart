// Public Markets: tokenized NGX equities and FMDQ bonds (app-backend
// /v1/public-markets). Amounts and quantities are decimal strings.

double pmNum(dynamic v) => double.tryParse('${v ?? ''}') ?? 0;

String _s(dynamic v) => v == null ? '' : '$v';

class PMAsset {
  final Map raw;
  PMAsset(this.raw);

  String get id => _s(raw['id']);
  String get assetCode => _s(raw['assetCode']);
  String get ticker => _s(raw['ticker']);
  String get name => _s(raw['name']);
  String get shortName => _s(raw['shortName']).isEmpty ? name : _s(raw['shortName']);
  String get market => _s(raw['market']);
  String get assetType => _s(raw['assetType']);
  String get sector => _s(raw['sector']);
  String get isin => _s(raw['isin']);
  String get description => _s(raw['description']);
  String get unitDescription => _s(raw['unitDescription']);
  String get price => _s(raw['price']);
  String get priceSource => _s(raw['priceSource']);
  bool get priceLive => raw['priceLive'] == true;
  String get dayChangePercent => _s(raw['dayChangePercent']);
  String get previousClose => _s(raw['previousClose']);
  String get dayHigh => _s(raw['dayHigh']);
  String get dayLow => _s(raw['dayLow']);
  String get dayVolume => _s(raw['dayVolume']);
  String get marketCap => _s(raw['marketCap']);
  String get peRatio => _s(raw['peRatio']);
  String get dividendYield => _s(raw['dividendYield']);
  String get coupon => _s(raw['coupon']);
  String get maturityDate => _s(raw['maturityDate']);
  String get status => _s(raw['status']); // open | halted | coming-soon
  String get contractAddress => _s(raw['contractAddress']);
  int get tokenDecimals => (raw['tokenDecimals'] as num?)?.toInt() ?? 0;
  String get minimumBuy => _s(raw['minimumBuy']);
  String get feePercent => _s(raw['feePercent']);
  String get fundingAsset => _s(raw['fundingAsset']).isEmpty ? 'CNGN' : _s(raw['fundingAsset']);
  String get tokensInCirculation => _s(raw['tokensInCirculation']);
  Map get logo => (raw['logo'] as Map?) ?? {};
  Map get session => (raw['session'] as Map?) ?? {};
  Map get custody => (raw['custody'] as Map?) ?? {};
  List<Map> get corporateActions => ((raw['corporateActions'] as List?) ?? []).map((e) => Map.from(e as Map)).toList();
  bool get isBond => assetType == 'BOND';
  bool get open => status == 'open';
  bool get marketOpen => session['open'] == true;
  String get initials => _s(logo['initials']).isEmpty ? (ticker.length > 2 ? ticker.substring(0, 2) : ticker) : _s(logo['initials']);
}

class PMQuote {
  final Map raw;
  PMQuote(this.raw);
  String get side => _s(raw['side']);
  String get amount => _s(raw['amount']);
  String get fee => _s(raw['fee']);
  String get feePercent => _s(raw['feePercent']);
  String get netAmount => _s(raw['netAmount']);
  String get quantity => _s(raw['quantity']);
  String get price => _s(raw['price']);
  String get priceSource => _s(raw['priceSource']);
  String get fundingAsset => _s(raw['fundingAsset']);
  bool get marketOpen => raw['marketOpen'] == true;
  String get path => _s(raw['path']); // FAST | SLOW | NETTED
  String get note => _s(raw['note']);
  String get nextSessionAt => _s(raw['nextSessionAt']);
  String get minimumBuy => _s(raw['minimumBuy']);
  String get custodianName => _s(raw['custodianName']);
  String get settlementNote => _s(raw['settlementNote']);
}

class PMOrder {
  final Map raw;
  PMOrder(this.raw);
  String get id => _s(raw['id']);
  String get type => _s(raw['type']); // CREATION | REDEMPTION
  bool get isBuy => type == 'CREATION';
  String get assetCode => _s(raw['assetCode']);
  String get state => _s(raw['state']);
  String get path => _s(raw['path']);
  String get amount => _s(raw['amount']);
  String get fee => _s(raw['fee']);
  String get netAmount => _s(raw['netAmount']);
  String get quantity => _s(raw['quantity']);
  String get referencePrice => _s(raw['referencePrice']);
  String get executedPrice => _s(raw['executedPrice']);
  String get fundingAssetCode => _s(raw['fundingAssetCode']);
  String get walletAddress => _s(raw['walletAddress']);
  String get walletAlias => _s(raw['walletAlias']);
  String get note => _s(raw['note']);
  String get createdAt => _s(raw['createdAt']);
  String get tokenTxHash => _s(raw['tokenTxHash']);
  String get payoutTxHash => _s(raw['payoutTxHash']);
  String get paymentTxHash => _s(raw['paymentTxHash']);
  List<Map> get events => ((raw['events'] as List?) ?? []).map((e) => Map.from(e as Map)).toList();
  bool get done => const ['complete', 'rejected', 'failed', 'cancelled'].contains(state);
  bool get ok => state == 'complete';
}

class PMHolding {
  final Map raw;
  PMHolding(this.raw);
  PMAsset get asset => PMAsset((raw['asset'] as Map?) ?? {});
  String get quantity => _s(raw['quantity']);
  String get marketValue => _s(raw['marketValue']);
  String get averageCost => _s(raw['averageCost']);
  String get costBasis => _s(raw['costBasis']);
  String get totalReturn => _s(raw['totalReturn']);
  String get returnPercent => _s(raw['returnPercent']);
  String get todayChange => _s(raw['todayChange']);
  String get incomeReceived => _s(raw['incomeReceived']);
  Map get wallets => (raw['wallets'] as Map?) ?? {};
}

class PMPortfolio {
  final Map raw;
  PMPortfolio(this.raw);
  String get value => _s(raw['value']);
  String get costBasis => _s(raw['costBasis']);
  String get totalReturn => _s(raw['totalReturn']);
  String get returnPercent => _s(raw['returnPercent']);
  String get todayChange => _s(raw['todayChange']);
  String get income => _s(raw['income']);
  List<PMHolding> get holdings => ((raw['holdings'] as List?) ?? []).map((e) => PMHolding(e as Map)).toList();
  List<PMOrder> get openOrders => ((raw['openOrders'] as List?) ?? []).map((e) => PMOrder(e as Map)).toList();
  List<Map> get activity => ((raw['activity'] as List?) ?? []).map((e) => Map.from(e as Map)).toList();
}

class PMDividend {
  final Map raw;
  PMDividend(this.raw);
  String get id => _s(raw['id']);
  String get assetCode => _s(raw['assetCode']);
  String get eventType => _s(raw['eventType']);
  String get description => _s(raw['description']);
  String get recordDate => _s(raw['recordDate']);
  String get payDate => _s(raw['payDate']);
  String get amountPerUnit => _s(raw['amountPerUnit']);
  String get custodianName => _s(raw['custodianName']);
  String get units => _s(raw['units']);
  String get grossAmount => _s(raw['grossAmount']);
  String get whtPercent => _s(raw['whtPercent']);
  String get whtAmount => _s(raw['whtAmount']);
  String get netAmount => _s(raw['netAmount']);
  String get status => _s(raw['status']);
  String get walletAddress => _s(raw['walletAddress']);
  String get txHash => _s(raw['txHash']);
  String get paidAt => _s(raw['paidAt']);
  bool get paid => status == 'PAID';
}

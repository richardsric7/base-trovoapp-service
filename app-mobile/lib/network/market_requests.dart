import 'dart:convert';

import 'package:candlesticks/candlesticks.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';

double _num(dynamic v) => double.tryParse('${v ?? ''}') ?? 0;

// MarketToken is one side of a pair: a token on the offer book, named by
// its contract address.
class MarketToken {
  final String code;
  final String name;
  final String contractAddress;
  final String imageUrl;

  MarketToken(Map m)
      : code = '${m['code'] ?? ''}',
        name = '${m['name'] ?? ''}',
        contractAddress = '${m['contractAddress'] ?? ''}',
        imageUrl = '${m['imageUrl'] ?? ''}';
}

// MarketPair is a pair's 24-hour summary from GET /v1/market/pairs. Prices
// are in the counter token per one base token.
class MarketPair {
  final MarketToken base;
  final MarketToken counter;
  final String lastPrice;
  final double changePercent24h;
  final String high24h;
  final String low24h;
  final String baseVolume24h;
  final String counterVolume24h;
  final int trades24h;
  final String bestAsk;
  final String bestBid;
  final DateTime? lastTradeAt;

  MarketPair(Map m)
      : base = MarketToken(Map.from(m['base'] ?? {})),
        counter = MarketToken(Map.from(m['counter'] ?? {})),
        lastPrice = '${m['lastPrice'] ?? '0'}',
        changePercent24h = _num(m['changePercent24h']),
        high24h = '${m['high24h'] ?? '0'}',
        low24h = '${m['low24h'] ?? '0'}',
        baseVolume24h = '${m['baseVolume24h'] ?? '0'}',
        counterVolume24h = '${m['counterVolume24h'] ?? '0'}',
        trades24h = int.tryParse('${m['trades24h'] ?? 0}') ?? 0,
        bestAsk = '${m['bestAsk'] ?? '0'}',
        bestBid = '${m['bestBid'] ?? '0'}',
        lastTradeAt = m['lastTradeAt'] == null ? null : DateTime.fromMillisecondsSinceEpoch(int.parse('${m['lastTradeAt']}') * 1000);

  String get symbol => '${base.code}/${counter.code}';
  bool get isUp => changePercent24h >= 0;

  // toMap keeps the pair in appState.viewData between screens.
  Map toMap() => {
        'base': {'code': base.code, 'name': base.name, 'contractAddress': base.contractAddress, 'imageUrl': base.imageUrl},
        'counter': {'code': counter.code, 'name': counter.name, 'contractAddress': counter.contractAddress, 'imageUrl': counter.imageUrl},
        'lastPrice': lastPrice,
        'changePercent24h': '$changePercent24h',
        'high24h': high24h,
        'low24h': low24h,
        'baseVolume24h': baseVolume24h,
        'counterVolume24h': counterVolume24h,
        'trades24h': trades24h,
        'bestAsk': bestAsk,
        'bestBid': bestBid,
        if (lastTradeAt != null) 'lastTradeAt': lastTradeAt!.millisecondsSinceEpoch ~/ 1000,
      };
}

// MarketLevel is one price level of the order book: amount is in the base
// token.
class MarketLevel {
  final double price;
  final double amount;
  MarketLevel(Map m)
      : price = _num(m['price']),
        amount = _num(m['amount']);
  double get total => price * amount;
}

class MarketOrderBook {
  final List<MarketLevel> asks; // cheapest first
  final List<MarketLevel> bids; // highest first
  MarketOrderBook(Map m)
      : asks = ((m['asks'] ?? []) as List).map((l) => MarketLevel(l as Map)).toList(),
        bids = ((m['bids'] ?? []) as List).map((l) => MarketLevel(l as Map)).toList();

  double? get spread => asks.isEmpty || bids.isEmpty ? null : asks.first.price - bids.first.price;
}

class MarketTrade {
  final DateTime at;
  final double price;
  final double baseAmount;
  final double counterAmount;
  MarketTrade(Map m)
      : at = DateTime.fromMillisecondsSinceEpoch(int.parse('${m['timestamp'] ?? 0}') * 1000),
        price = _num(m['price']),
        baseAmount = _num(m['baseAmount']),
        counterAmount = _num(m['counterAmount']);
}

// MarketOffer is one of the wallet's own market-making offers from
// GET /v1/users/trades.
class MarketOffer {
  final String id;
  final String offerType; // BUY | SELL
  final String assetCode;
  final String contractAddress;
  final String currencyCode;
  final String currencyIssuer;
  final double pricePerUnit;
  final double quantity;
  final double remaining; // what the offer still sells, in the token it sells
  final bool open;
  final bool canceled;
  final String transactionId;
  final String blockchainOfferId;
  final double sold; // sold so far, in the token it sells
  final double received; // paid so far, in the token it buys
  final int fills;

  MarketOffer(Map m)
      : id = '${m['id'] ?? ''}',
        offerType = '${m['offerType'] ?? ''}'.toUpperCase(),
        assetCode = '${m['assetCode'] ?? ''}',
        contractAddress = '${m['contractAddress'] ?? ''}',
        currencyCode = '${m['currencyCode'] ?? ''}',
        currencyIssuer = '${m['currencyIssuer'] ?? ''}',
        pricePerUnit = _num(m['pricePerUnit']),
        quantity = _num(m['quantity']),
        remaining = _num(m['remaining']),
        open = m['open'] == true,
        canceled = '${m['canceled']}' == '1',
        transactionId = '${m['transactionId'] ?? ''}',
        blockchainOfferId = '${m['blockchainOfferId'] ?? ''}',
        sold = _num(m['sold']),
        received = _num(m['received']),
        fills = int.tryParse('${m['fills'] ?? 0}') ?? 0;

  bool get isBuy => offerType == 'BUY';

  // what has traded so far: base tokens and counter tokens.
  double get filledBase => isBuy ? received : sold;
  double get filledCounter => isBuy ? sold : received;
  // the average price of the fills, else the asked price.
  double get averagePrice => filledBase > 0 ? filledCounter / filledBase : pricePerUnit;

  // status: on the book, waiting to be mined, filled or cancelled.
  String get status {
    if (canceled) return fills > 0 ? 'Cancelled (part filled)' : 'Cancelled';
    if (open && remaining == 0 && fills > 0) return 'Filled';
    if (open) return fills > 0 ? 'Part filled' : 'Open';
    if (blockchainOfferId.isNotEmpty) return fills > 0 ? 'Filled' : 'Closed';
    return 'Pending';
  }

  bool get cancellable => open && !canceled && remaining > 0;

  bool forPair(MarketPair p) =>
      contractAddress.toLowerCase() == p.base.contractAddress.toLowerCase() &&
      currencyIssuer.toLowerCase() == p.counter.contractAddress.toLowerCase();
}

// formatMarketNumber shows prices and amounts: two decimals from 1 up,
// four significant digits below.
String formatMarketNumber(double v) {
  if (v == 0) return '0';
  if (v.abs() >= 1) {
    final parts = v.toStringAsFixed(2).split('.');
    final whole = parts[0].replaceAllMapped(RegExp(r'\B(?=(\d{3})+(?!\d))'), (m) => ',');
    return parts[1] == '00' ? whole : '$whole.${parts[1]}';
  }
  return double.parse(v.toStringAsPrecision(4)).toString();
}

String formatMarketString(String v) => formatMarketNumber(_num(v));

// MarketFavorites keeps the pairs a user starred, on the phone.
class MarketFavorites {
  static const _key = 'marketFavoritePairs';

  static String id(MarketPair p) => '${p.base.contractAddress}/${p.counter.contractAddress}'.toLowerCase();

  static Future<Set<String>> load() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      return (prefs.getStringList(_key) ?? []).toSet();
    } catch (_) {
      return {};
    }
  }

  static Future<Set<String>> toggle(MarketPair p) async {
    final set = await load();
    final k = id(p);
    if (!set.remove(k)) set.add(k);
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setStringList(_key, set.toList());
    } catch (_) {}
    return set;
  }
}

// candle resolutions the backend serves, by the label the chart shows.
const marketResolutions = {'15 min': '15m', 'Hour': '1h', '4 hours': '4h', 'Day': '1d', 'Week': '1w'};

// MarketApi wraps /v1/market (public market data) and /v1/users/trades (the
// wallet's own offers: signed requests). Placing and cancelling an offer
// take two calls, like every wallet operation: the first returns the
// operation to sign, the second sends it signed (or committed, for a
// shared wallet).
class MarketApi {
  final DataProvider appState;
  MarketApi(this.appState);

  String get _signer => appState.primaryWallet.signer!;
  String get _secretKey => appState.secretKeys[0];
  String get _address => appState.primaryWallet.address!;

  static String message(Map r, String fallback) {
    final d = r['data'];
    if (d is Map) return '${d['message'] ?? d['error'] ?? fallback}';
    return fallback;
  }

  static String _pairQuery(MarketPair p, int limit) =>
      'base=${Uri.encodeComponent(p.base.contractAddress)}&counter=${Uri.encodeComponent(p.counter.contractAddress)}&limit=$limit';

  // pairs returns null when the market cannot be reached or is not set up.
  Future<List<MarketPair>?> pairs() async {
    final r = await makeUnSecuredGetRequest('/v1/market/pairs');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return ((r['data']['pairs'] ?? []) as List).map((p) => MarketPair(p as Map)).toList();
  }

  Future<MarketOrderBook?> orderBook(MarketPair p, {int limit = 20}) async {
    final r = await makeUnSecuredGetRequest('/v1/market/orderbook?${_pairQuery(p, limit)}');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return MarketOrderBook(r['data'] as Map);
  }

  Future<List<MarketTrade>?> trades(MarketPair p, {int limit = 50}) async {
    final r = await makeUnSecuredGetRequest('/v1/market/trades?${_pairQuery(p, limit)}');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return ((r['data']['trades'] ?? []) as List).map((t) => MarketTrade(t as Map)).toList();
  }

  // candles returns them newest first, as the chart widget expects.
  Future<List<Candle>?> candles(MarketPair p, String resolution, {int limit = 200}) async {
    final r = await makeUnSecuredGetRequest('/v1/market/candles?${_pairQuery(p, limit)}&resolution=$resolution');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    final list = ((r['data']['candles'] ?? []) as List).map((c) {
      final m = c as Map;
      return Candle(
        date: DateTime.fromMillisecondsSinceEpoch(int.parse('${m['timestamp'] ?? 0}') * 1000),
        open: _num(m['open']),
        high: _num(m['high']),
        low: _num(m['low']),
        close: _num(m['close']),
        volume: _num(m['baseVolume']),
      );
    }).toList();
    return list.reversed.toList();
  }

  Future<List<MarketOffer>?> myOffers() async {
    final r = await makeGetRequest(uri: '/v1/users/trades', signer: _signer, secretKey: _secretKey, address: _address);
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return ((r['data']['offers'] ?? []) as List).map((o) => MarketOffer(o as Map)).toList();
  }

  // placeOffer: without transaction it builds the operation
  // ({transaction, messages, ...}); with it, it submits.
  Future<Map> placeOffer(MarketPair p,
      {required bool buy,
      required String price,
      required String quantity,
      String transaction = '',
      String signature = '',
      bool commit = false}) {
    final body = {
      'offerType': buy ? 'BUY' : 'SELL',
      'assetCode': p.base.code,
      'contractAddress': p.base.contractAddress,
      'currencyCode': p.counter.code,
      'currencyIssuer': p.counter.contractAddress,
      'pricePerUnit': price,
      'quantity': quantity,
      if (transaction.isNotEmpty) 'transaction': transaction,
      if (signature.isNotEmpty) 'transactionSignature': signature,
      if (commit) 'commit': 1,
    };
    return makePostRequest(uri: '/v1/users/trades', body: jsonEncode(body), signer: _signer, secretKey: _secretKey, address: _address);
  }

  // cancelOffer works in the same two steps as placeOffer.
  Future<Map> cancelOffer(String id, {String transaction = '', String signature = '', bool commit = false}) {
    final body = {
      if (transaction.isNotEmpty) 'transaction': transaction,
      if (signature.isNotEmpty) 'transactionSignature': signature,
      if (commit) 'commit': 1,
    };
    return makeDeleteRequest(
        uri: '/v1/users/trades/${Uri.encodeComponent(id)}', body: jsonEncode(body), signer: _signer, secretKey: _secretKey, address: _address);
  }
}

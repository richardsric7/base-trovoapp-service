import 'dart:convert';

import 'package:trovo_app/models/public_market.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';

// PublicMarketsApi wraps /v1/public-markets: the listing and prices are
// public, everything else is a signed request. Buying and selling are two
// calls, like every wallet operation: trade() without a signature returns
// the quote and the operation to sign; submit() sends it signed and places
// the order.
class PublicMarketsApi {
  final DataProvider appState;
  PublicMarketsApi(this.appState);

  String get _signer => appState.primaryWallet.signer!;
  String get _secretKey => appState.secretKeys[0];
  String get _address => appState.primaryWallet.address!;

  Future<Map> _get(String uri) => makeGetRequest(uri: uri, signer: _signer, secretKey: _secretKey, address: _address);

  Future<Map> _post(String uri, Map body) =>
      makePostRequest(uri: uri, body: jsonEncode(body), signer: _signer, secretKey: _secretKey, address: _address);

  static String message(Map r, String fallback) {
    final d = r['data'];
    if (d is Map) return '${d['message'] ?? d['error'] ?? fallback}';
    return fallback;
  }

  // assets lists the assets (market: NGX | FMDQ, type: EQUITY | BOND).
  Future<List<PMAsset>?> assets({String market = '', String type = '', String search = ''}) async {
    final q = Uri(queryParameters: {
      if (market.isNotEmpty) 'market': market,
      if (type.isNotEmpty) 'type': type,
      if (search.isNotEmpty) 'search': search,
    }).query;
    final r = await makeUnSecuredGetRequest('/v1/public-markets${q.isEmpty ? '' : '?$q'}');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return ((r['data']['assets'] ?? []) as List).map((a) => PMAsset(a as Map)).toList();
  }

  Future<PMAsset?> asset(String code) async {
    final r = await makeUnSecuredGetRequest('/v1/public-markets/assets/${Uri.encodeComponent(code)}');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return PMAsset(r['data'] as Map);
  }

  // prices returns the chart points ([{at, price}]) for 1D, 1W, 1M, 3M, 1Y or All.
  Future<List<Map>> prices(String code, String range) async {
    final r = await makeUnSecuredGetRequest('/v1/public-markets/assets/${Uri.encodeComponent(code)}/prices?range=$range');
    if (r['statusCode'] != 200 || r['data'] is! Map) return [];
    return ((r['data']['points'] ?? []) as List).map((p) => Map.from(p as Map)).toList();
  }

  Future<PMQuote?> quote(String code, {required bool buy, required String value}) async {
    final r = await _get('/v1/public-markets/assets/${Uri.encodeComponent(code)}/quote?side=${buy ? 'buy' : 'sell'}&${buy ? 'amount' : 'quantity'}=$value');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return PMQuote(r['data'] as Map);
  }

  // trade builds the operation: {quote, transaction, messages}.
  Future<Map> trade(String code, {required bool buy, required String walletAddress, required String value}) =>
      _post('/v1/public-markets/assets/${Uri.encodeComponent(code)}/${buy ? 'buy' : 'sell'}',
          {'walletAddress': walletAddress, buy ? 'amount' : 'quantity': value});

  // submit sends the signed operation: {quote, order}.
  Future<Map> submit(String code,
          {required bool buy, required String walletAddress, required String value, required String transaction, required String signature}) =>
      _post('/v1/public-markets/assets/${Uri.encodeComponent(code)}/${buy ? 'buy' : 'sell'}', {
        'walletAddress': walletAddress,
        buy ? 'amount' : 'quantity': value,
        'transaction': transaction,
        'transactionSignature': signature,
      });

  Future<PMPortfolio?> portfolio() async {
    final r = await _get('/v1/public-markets/portfolio');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return PMPortfolio(r['data'] as Map);
  }

  Future<List<PMOrder>> orders() async {
    final r = await _get('/v1/public-markets/orders');
    if (r['statusCode'] != 200 || r['data'] is! Map) return [];
    return ((r['data']['orders'] ?? []) as List).map((o) => PMOrder(o as Map)).toList();
  }

  Future<PMOrder?> order(String id) async {
    final r = await _get('/v1/public-markets/orders/${Uri.encodeComponent(id)}');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return PMOrder(r['data'] as Map);
  }

  Future<List<PMDividend>?> dividends({String assetCode = ''}) async {
    final r = await _get('/v1/public-markets/dividends${assetCode.isEmpty ? '' : '?assetCode=${Uri.encodeComponent(assetCode)}'}');
    if (r['statusCode'] != 200 || r['data'] is! Map) return null;
    return ((r['data']['dividends'] ?? []) as List).map((d) => PMDividend(d as Map)).toList();
  }
}

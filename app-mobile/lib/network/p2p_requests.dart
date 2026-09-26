import 'dart:convert';

import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/models/p2p_order.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';

// P2PApi wraps the existing signed-request helpers (makeGetRequest/
// makePostRequest) for the /v1/p2p/... routes - it does not introduce a
// second networking layer, it just gives P2P screens typed call sites.
class P2PApi {
  final DataProvider appState;
  P2PApi(this.appState);

  String get _signer => appState.primaryWallet.signer!;
  String get _secretKey => appState.secretKeys[0];

  Future<Map> _get(String uri, {String address = ''}) {
    return makeGetRequest(
      uri: uri,
      signer: _signer,
      secretKey: _secretKey,
      address: address.isEmpty ? appState.primaryWallet.address! : address,
    );
  }

  Future<Map> _post(String uri, Map body, {String address = ''}) {
    return makePostRequest(
      uri: uri,
      body: jsonEncode(body),
      signer: _signer,
      secretKey: _secretKey,
      address: address.isEmpty ? appState.primaryWallet.address! : address,
    );
  }

  // ---- Offers / marketplace ----

  Future<Map<String, dynamic>> listMarketplaceOffers({
    String? offerType,
    String? asset,
    String? countryCode,
    String? currency,
    int page = 1,
    int pageSize = 20,
  }) async {
    var qp = <String, String>{'page': '$page', 'pageSize': '$pageSize'};
    if (offerType != null) qp['offerType'] = offerType;
    if (asset != null) qp['asset'] = asset;
    if (countryCode != null) qp['countryCode'] = countryCode;
    if (currency != null) qp['currency'] = currency;
    var query = qp.entries.map((e) => '${e.key}=${e.value}').join('&');
    var response = await _get('/v1/p2p/offers?$query');
    var data = response['data'];
    return {
      'offers': P2POffer.deserializeList(data['data']),
      'total': data['total'] ?? 0,
      'statusCode': response['statusCode'],
    };
  }

  Future<P2POffer?> getOffer(String offerId) async {
    var response = await _get('/v1/p2p/offers/$offerId');
    if (response['statusCode'] != 200) return null;
    return P2POffer().deserializeJson(response['data']);
  }

  Future<Map<String, dynamic>?> quoteOrderFees(
    String offerId,
    String amount,
  ) async {
    var response = await _get('/v1/p2p/offers/$offerId/quote?amount=$amount');
    if (response['statusCode'] != 200) return null;
    return response['data'];
  }

  Future<Map> createOffer(Map body) => _post('/v1/p2p/offers', body);

  Future<Map> activateOffer(String offerId) =>
      _post('/v1/p2p/offers/$offerId/activate', {});

  Future<Map> pauseOffer(String offerId) =>
      _post('/v1/p2p/offers/$offerId/pause', {});

  Future<List<P2POffer>> listMyOffers() async {
    var response = await _get('/v1/p2p/my-offers');
    return P2POffer.deserializeList(response['data']?['data']);
  }

  // ---- Orders ----

  Future<Map> createOrder({
    required String offerId,
    required String specifiedAssetAmount,
  }) => _post('/v1/p2p/orders', {
    'offerId': offerId,
    'specifiedAssetAmount': specifiedAssetAmount,
  });

  Future<Map<String, dynamic>> listMyOrders({
    String? role,
    String? status,
    int page = 1,
    int pageSize = 20,
  }) async {
    var qp = <String, String>{'page': '$page', 'pageSize': '$pageSize'};
    if (role != null) qp['role'] = role;
    if (status != null) qp['status'] = status;
    var query = qp.entries.map((e) => '${e.key}=${e.value}').join('&');
    var response = await _get('/v1/p2p/orders?$query');
    var data = response['data'];
    return {
      'orders': P2POrder.deserializeList(data['data']),
      'total': data['total'] ?? 0,
    };
  }

  Future<P2POrder?> getOrder(String orderId) async {
    var response = await _get('/v1/p2p/orders/$orderId');
    if (response['statusCode'] != 200) return null;
    return P2POrder().deserializeJson(response['data']);
  }

  Future<Map> acceptOrder(String orderId, {required String address}) =>
      _post('/v1/p2p/orders/$orderId/accept', {}, address: address);

  Future<Map> rejectOrder(String orderId, {required String address}) =>
      _post('/v1/p2p/orders/$orderId/reject', {}, address: address);

  Future<Map> cancelOrder(String orderId, {required String address}) =>
      _post('/v1/p2p/orders/$orderId/cancel', {}, address: address);

  // ---- Escrow deposit (two-phase build -> sign -> commit, same contract
  // as /v1/users/payment) ----

  Future<Map> escrowDepositBuild(String orderId, {required String address}) =>
      _post('/v1/p2p/orders/$orderId/escrow-deposit', {}, address: address);

  Future<Map> escrowDepositCommit(
    String orderId, {
    required String address,
    required String transaction,
    required String transactionSignature,
  }) => _post('/v1/p2p/orders/$orderId/escrow-deposit', {
    'transaction': transaction,
    'transactionSignature': transactionSignature,
    'commit': 1,
  }, address: address);

  // ---- Fiat payment stage ----

  Future<Map> markPaymentSent(String orderId, {required String address}) =>
      _post('/v1/p2p/orders/$orderId/payment-sent', {}, address: address);

  Future<Map> confirmPaymentReceived(
    String orderId, {
    required String address,
  }) => _post('/v1/p2p/orders/$orderId/payment-confirmed', {}, address: address);

  // ---- Disputes ----

  Future<P2PDispute?> getOpenDisputeForOrder(String orderId) async {
    var response = await _get('/v1/p2p/orders/$orderId/dispute');
    if (response['statusCode'] != 200) return null;
    return P2PDispute().deserializeJson(response['data']);
  }

  Future<Map> openDispute({
    required String orderId,
    required String subject,
    required String description,
    List<String> evidence = const [],
  }) => _post('/v1/p2p/orders/$orderId/disputes', {
    'subject': subject,
    'description': description,
    'evidence': evidence,
  });

  Future<Map> merchantConfirmsPayment(String disputeId) =>
      _post('/v1/p2p/disputes/$disputeId/merchant-confirms-payment', {});

  Future<Map> buyerConfirmsNotPaid(String disputeId) =>
      _post('/v1/p2p/disputes/$disputeId/buyer-confirms-not-paid', {});
}

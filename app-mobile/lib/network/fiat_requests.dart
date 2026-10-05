import 'dart:convert';

import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';

// FiatApi wraps the signed-request helpers for the /v1/users/stablerail/...
// routes: Naira deposits (onramp) and withdrawals to a Nigerian bank account
// (offramp). Users are onboarded with Stablerail by the backend when KYC
// level 1 (BVN) completes - there is no BVN step here.
class FiatApi {
  final DataProvider appState;
  // walletAddress is the wallet deposits go to and withdrawals leave from
  // (the primary wallet when not given).
  final String? walletAddress;
  FiatApi(this.appState, {this.walletAddress});

  String get _signer => appState.primaryWallet.signer!;
  String get _secretKey => appState.secretKeys[0];
  String get _address => walletAddress ?? appState.primaryWallet.address!;

  Future<Map> _get(String uri) => makeGetRequest(uri: uri, signer: _signer, secretKey: _secretKey, address: _address);

  Future<Map> _post(String uri, Map body) =>
      makePostRequest(uri: uri, body: jsonEncode(body), signer: _signer, secretKey: _secretKey, address: _address);

  // profile: {enabled, onboarded, onboardingStatus, minimumWithdrawal}
  Future<Map> profile() => _get('/v1/users/stablerail/profile');

  // banks: [{bank_code, bank_name}]
  Future<List<Map>> banks() async {
    final r = await _get('/v1/users/stablerail/banks');
    if (r['statusCode'] != 200 || r['data'] is! List) return [];
    final list = (r['data'] as List).map((b) => Map.from(b as Map)).toList();
    list.sort((a, b) => '${a['bank_name']}'.compareTo('${b['bank_name']}'));
    return list;
  }

  // deposit returns the virtual account to pay into.
  Future<Map> deposit(String amount) => _post('/v1/users/stablerail/onrampcngn/$amount', {});

  // withdrawBuild returns the transaction to sign (and messages to show);
  // withdrawSubmit sends it signed.
  Future<Map> withdrawBuild({required String amount, required String accountNumber, required String bankCode}) =>
      _post('/v1/users/stablerail/withdraw', {'amount': amount, 'accountNumber': accountNumber, 'bankCode': bankCode});

  Future<Map> withdrawSubmit({
    required String amount,
    required String accountNumber,
    required String bankCode,
    required String transaction,
    required String transactionSignature,
  }) =>
      _post('/v1/users/stablerail/withdraw', {
        'amount': amount,
        'accountNumber': accountNumber,
        'bankCode': bankCode,
        'transaction': transaction,
        'transactionSignature': transactionSignature,
      });

  Future<List<Map>> withdrawals() async {
    final r = await _get('/v1/users/stablerail/withdrawals');
    if (r['statusCode'] != 200) return [];
    return ((r['data']?['withdrawals'] ?? []) as List).map((w) => Map.from(w as Map)).toList();
  }
}

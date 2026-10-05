import 'package:trovo_app/models/proceed_payout.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';

// PayoutApi reads the proceeds payouts (dividends / interest) paid or
// scheduled to the user's wallets.
class PayoutApi {
  final DataProvider appState;
  PayoutApi(this.appState);

  // receipts for one asset token (by its contract, else its code), newest
  // first; null when the request failed
  Future<List<ProceedPayoutReceipt>?> receipts({
    String? assetCode,
    String? tokenContract,
  }) async {
    final r = await makeGetRequest(
      uri: '/v1/tokenization/payouts',
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0],
      address: appState.primaryWallet.address!,
    );
    if (r['statusCode'] != 200 ||
        r['data'] is! Map ||
        (r['data'] as Map)['payouts'] is! List) {
      return null;
    }
    final all = ((r['data'] as Map)['payouts'] as List)
        .map((m) => ProceedPayoutReceipt.fromMap(m as Map))
        .toList();
    final contract = (tokenContract ?? '').toLowerCase();
    final code = (assetCode ?? '').toUpperCase();
    return all.where((p) {
      if (contract.isNotEmpty && p.tokenContractAddress.isNotEmpty) {
        return p.tokenContractAddress.toLowerCase() == contract;
      }
      return code.isEmpty || p.assetCode.toUpperCase() == code;
    }).toList();
  }
}

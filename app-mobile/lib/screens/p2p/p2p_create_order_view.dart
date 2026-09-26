import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

// P2PCreateOrderView (Plan Section 95.3): amount entry with a live,
// real (not estimated) itemized fee breakdown fetched from the backend's
// quote endpoint before the customer commits (Plan Section 97.6).
class P2PCreateOrderView extends StatefulWidget {
  const P2PCreateOrderView({Key? key}) : super(key: key);

  @override
  State<P2PCreateOrderView> createState() => _P2PCreateOrderViewState();
}

class _P2PCreateOrderViewState extends State<P2PCreateOrderView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  bool loaded = false;
  P2POffer? offer;
  final amountController = TextEditingController();
  Timer? _debounce;
  Map<String, dynamic>? quote;
  String? quoteError;
  bool quoting = false;
  bool submitting = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    if (!loaded) {
      loaded = true;
      offer = appState.viewData?[P2PCreateOrderViewPageConfig.key]?['offer'];
    }
  }

  @override
  void dispose() {
    _debounce?.cancel();
    amountController.dispose();
    super.dispose();
  }

  void _onAmountChanged(String value) {
    _debounce?.cancel();
    setState(() {
      quote = null;
      quoteError = null;
    });
    if (value.trim().isEmpty || offer?.id == null) return;
    _debounce = Timer(const Duration(milliseconds: 500), () async {
      setState(() => quoting = true);
      final result = await api.quoteOrderFees(offer!.id!, value.trim());
      setState(() {
        quoting = false;
        if (result == null) {
          quoteError = 'This amount is outside the offer\'s allowed range.';
        } else {
          quote = result;
        }
      });
    });
  }

  Future<void> _submit() async {
    if (offer?.id == null || quote == null) return;
    setState(() => submitting = true);
    showLoader(context);
    try {
      final response = await api.createOrder(
        offerId: offer!.id!,
        specifiedAssetAmount: amountController.text.trim(),
      );
      hideLoader(context);
      setState(() => submitting = false);
      if (response['statusCode'] == 201) {
        appState.viewData ??= {};
        appState.viewData![P2POrderDetailViewPageConfig.key] = {
          'orderId': response['data']['id'],
        };
        appState.currentAction = PageAction(
          state: PageState.addAll,
          pages: [P2PMarketplaceViewPageConfig, P2POrderDetailViewPageConfig],
        );
      } else {
        popup(
          context,
          title: 'Could not create order',
          message: response['data']?['message'] ?? 'Please try again.',
        );
      }
    } catch (e) {
      hideLoader(context);
      setState(() => submitting = false);
      popup(context, title: 'Error', message: e.toString());
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);

    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('Create order', style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: offer == null
          ? const P2PEmptyState(icon: Icons.error_outline, message: 'No offer selected.')
          : Column(
              children: [
                Expanded(
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(P2PTheme.space4),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Amount (${offer!.asset})',
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        ),
                        const SizedBox(height: P2PTheme.space2),
                        TextField(
                          controller: amountController,
                          keyboardType: const TextInputType.numberWithOptions(decimal: true),
                          onChanged: _onAmountChanged,
                          decoration: InputDecoration(
                            hintText: 'Between ${offer!.minOrderAmount} and ${offer!.maxOrderAmount}',
                            filled: true,
                            fillColor: Colors.white,
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
                              borderSide: BorderSide.none,
                            ),
                          ),
                        ),
                        const SizedBox(height: P2PTheme.space4),
                        if (quoting) const Center(child: CircularProgressIndicator()),
                        if (quoteError != null)
                          Text(quoteError!, style: const TextStyle(color: P2PTheme.danger)),
                        if (quote != null) _breakdown(quote!),
                      ],
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(P2PTheme.space4),
                  child: SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: P2PTheme.brandDark,
                        padding: const EdgeInsets.symmetric(vertical: P2PTheme.space4),
                      ),
                      onPressed: (quote != null && !submitting) ? _submit : null,
                      child: const Text(
                        'Review and submit order',
                        style: TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                      ),
                    ),
                  ),
                ),
              ],
            ),
    );
  }

  Widget _breakdown(Map<String, dynamic> q) {
    final isBuy = offer!.isBuy;
    final asset = offer!.asset;
    return P2PListCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Order breakdown', style: TextStyle(fontWeight: FontWeight.w700)),
          const SizedBox(height: P2PTheme.space3),
          _row('You pay (fiat)', '${q['paymentAmount']} ${offer!.currency}'),
          const Divider(height: P2PTheme.space4),
          if (isBuy) ...[
            _row('You receive', '${q['buyerNetAssetAmount']} $asset'),
            _row('Platform + regulatory fee', '${q['buyerTotalFees']} $asset'),
            _row('VAT on fee', '${q['buyerTotalVat']} $asset'),
          ] else ...[
            _row('Escrow deposit required', '${q['sellerEscrowAssetAmount']} $asset'),
            _row('Platform + regulatory fee', '${q['sellerTotalFees']} $asset'),
            _row('VAT on fee', '${q['sellerTotalVat']} $asset'),
          ],
        ],
      ),
    );
  }

  Widget _row(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: P2PTheme.space1),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.black54)),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }
}

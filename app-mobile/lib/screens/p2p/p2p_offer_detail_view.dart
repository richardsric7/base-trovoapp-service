import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

// P2POfferDetailView (Plan Section 95.2): full offer terms + a Buy/Sell CTA
// into order creation.
class P2POfferDetailView extends StatefulWidget {
  const P2POfferDetailView({Key? key}) : super(key: key);

  @override
  State<P2POfferDetailView> createState() => _P2POfferDetailViewState();
}

class _P2POfferDetailViewState extends State<P2POfferDetailView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  String? offerId;
  P2POffer? offer;
  bool loading = true;
  bool loaded = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    if (!loaded) {
      loaded = true;
      offerId = appState.viewData?[P2POfferDetailViewPageConfig.key]?['offerId'];
      _load();
    }
  }

  Future<void> _load() async {
    if (offerId == null) return;
    setState(() => loading = true);
    final result = await api.getOffer(offerId!);
    setState(() {
      offer = result;
      loading = false;
    });
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
        title: Text('Offer details', style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: loading
          ? const Center(child: CircularProgressIndicator())
          : offer == null
              ? const P2PEmptyState(
                  icon: Icons.error_outline,
                  message: 'This offer could not be found.',
                )
              : _content(offer!),
    );
  }

  Widget _content(P2POffer o) {
    final isCustomerBuying = o.isBuy; // merchant is buying -> customer sells to them... UI copy adapts below
    return Column(
      children: [
        Expanded(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(P2PTheme.space4),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                P2PListCard(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          CircleAvatar(
                            backgroundColor: P2PTheme.brandDark,
                            child: Text(
                              (o.merchantUsername ?? '?').substring(0, 1).toUpperCase(),
                              style: const TextStyle(color: Colors.white),
                            ),
                          ),
                          const SizedBox(width: P2PTheme.space3),
                          Text(
                            o.merchantUsername ?? '',
                            style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16),
                          ),
                        ],
                      ),
                      const SizedBox(height: P2PTheme.space4),
                      _row('Price', '${o.price} ${o.currency}'),
                      _row('Available', '${o.availableLiquidity} ${o.asset}'),
                      _row('Limits', '${o.minOrderAmount} - ${o.maxOrderAmount} ${o.asset}'),
                      _row('Payment method', o.paymentMethod?.paymentChannel ?? '-'),
                      _row('Provider', o.paymentMethod?.provider ?? '-'),
                      if (o.remark != null && o.remark!.isNotEmpty) ...[
                        const SizedBox(height: P2PTheme.space2),
                        Text(o.remark!, style: const TextStyle(color: Colors.black54)),
                      ],
                    ],
                  ),
                ),
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
              onPressed: () {
                appState.viewData ??= {};
                appState.viewData![P2PCreateOrderViewPageConfig.key] = {
                  'offer': o,
                };
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: P2PCreateOrderViewPageConfig,
                );
              },
              child: Text(
                isCustomerBuying ? 'Sell to this offer' : 'Buy from this offer',
                style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _row(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: P2PTheme.space1),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: const TextStyle(color: Colors.black54)),
          Flexible(
            child: Text(
              value,
              textAlign: TextAlign.right,
              style: const TextStyle(fontWeight: FontWeight.w600),
            ),
          ),
        ],
      ),
    );
  }
}

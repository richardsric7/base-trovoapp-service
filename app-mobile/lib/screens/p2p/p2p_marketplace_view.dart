import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

// P2PMarketplaceView is the P2P entry point (Plan Section 95.1): a
// Buy/Sell toggle over a filterable offer list, each offer rendered as a
// P2PListCard. Replaces the "coming soon" buttons in home.dart and
// tokenized_asset_details.dart.
class P2PMarketplaceView extends StatefulWidget {
  const P2PMarketplaceView({Key? key}) : super(key: key);

  @override
  State<P2PMarketplaceView> createState() => _P2PMarketplaceViewState();
}

class _P2PMarketplaceViewState extends State<P2PMarketplaceView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;

  String offerType = 'BUY';
  bool loading = true;
  String? error;
  List<P2POffer> offers = [];
  final Map<String, Future<Map<String, dynamic>?>> _perfFutures = {};

  // Cached per merchantId so scrolling the same offer list doesn't refetch
  // performance for a merchant already fetched (Plan Section 8/26's trust
  // signal, shown on every offer card).
  Future<Map<String, dynamic>?> _perfFor(String? merchantId) {
    if (merchantId == null || merchantId.isEmpty) return Future.value(null);
    return _perfFutures.putIfAbsent(
      merchantId,
      () => api.getMerchantPerformance(merchantId),
    );
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
  }

  Future<void> _load() async {
    setState(() {
      loading = true;
      error = null;
    });
    try {
      var result = await api.listMarketplaceOffers(offerType: offerType);
      setState(() {
        offers = result['offers'];
        loading = false;
      });
    } catch (e) {
      setState(() {
        error = 'p2pcouldnotloadoffers'.tr();
        loading = false;
      });
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
        title: Text('p2pmarketplacetitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
        actions: [
          IconButton(
            icon: const Icon(Icons.receipt_long),
            tooltip: 'p2pmyorderstooltip'.tr(),
            onPressed: () {
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: P2PMyOrdersViewPageConfig,
              );
            },
          ),
          IconButton(
            icon: const Icon(Icons.storefront),
            tooltip: 'p2pmyofferstooltip'.tr(),
            onPressed: () {
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: P2PMyOffersViewPageConfig,
              );
            },
          ),
          IconButton(
            icon: const Icon(Icons.savings_outlined),
            tooltip: 'p2pmyrefundstooltip'.tr(),
            onPressed: () {
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: P2PMyRefundsViewPageConfig,
              );
            },
          ),
        ],
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(P2PTheme.space4),
            child: Row(
              children: [
                Expanded(child: _toggleButton('BUY', 'p2pbuy'.tr())),
                const SizedBox(width: P2PTheme.space2),
                Expanded(child: _toggleButton('SELL', 'p2psell'.tr())),
              ],
            ),
          ),
          Expanded(child: _body()),
        ],
      ),
    );
  }

  Widget _toggleButton(String value, String label) {
    final selected = offerType == value;
    return GestureDetector(
      onTap: () {
        setState(() => offerType = value);
        _load();
      },
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: P2PTheme.space3),
        decoration: BoxDecoration(
          color: selected ? P2PTheme.brandDark : Colors.white,
          borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
        ),
        child: Center(
          child: Text(
            label,
            style: TextStyle(
              color: selected ? Colors.white : Colors.black87,
              fontWeight: FontWeight.w600,
            ),
          ),
        ),
      ),
    );
  }

  Widget _body() {
    if (loading) {
      return const Center(child: CircularProgressIndicator());
    }
    if (error != null) {
      return P2PEmptyState(
        icon: Icons.wifi_off,
        message: error!,
        ctaLabel: 'retry'.tr(),
        onCta: _load,
      );
    }
    if (offers.isEmpty) {
      return P2PEmptyState(
        icon: Icons.storefront_outlined,
        message: 'p2pnoofferstext'.tr(),
        ctaLabel: 'refresh'.tr(),
        onCta: _load,
      );
    }
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4),
        itemCount: offers.length,
        itemBuilder: (context, i) => _offerCard(offers[i]),
      ),
    );
  }

  Widget _offerCard(P2POffer offer) {
    return P2PListCard(
      onTap: () {
        appState.viewData ??= {};
        appState.viewData![P2POfferDetailViewPageConfig.key] = {
          'offerId': offer.id,
        };
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: P2POfferDetailViewPageConfig,
        );
      },
      child: Row(
        children: [
          CircleAvatar(
            backgroundColor: P2PTheme.brandDark,
            child: Text(
              (offer.merchantUsername ?? '?').substring(0, 1).toUpperCase(),
              style: const TextStyle(color: Colors.white),
            ),
          ),
          const SizedBox(width: P2PTheme.space3),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        offer.merchantUsername ?? '',
                        style: const TextStyle(fontWeight: FontWeight.w600),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    Container(
                      width: 8,
                      height: 8,
                      decoration: BoxDecoration(
                        color: offer.isOnline ? P2PTheme.success : Colors.grey,
                        shape: BoxShape.circle,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: P2PTheme.space1),
                Text(
                  '${offer.price} ${offer.currency} / ${offer.asset}',
                  style: const TextStyle(color: Colors.black54, fontSize: 13),
                ),
                const SizedBox(height: P2PTheme.space1),
                Text(
                  'p2plimitrange'.tr(args: [
                    offer.minOrderAmount ?? '',
                    offer.maxOrderAmount ?? '',
                    offer.asset ?? '',
                  ]),
                  style: const TextStyle(color: Colors.black38, fontSize: 12),
                ),
                FutureBuilder<Map<String, dynamic>?>(
                  future: _perfFor(offer.merchantUserId),
                  builder: (context, snapshot) {
                    final perf = snapshot.data;
                    final completedTrades = perf?['completedTrades'] ?? 0;
                    if (perf == null || completedTrades == 0) {
                      return const SizedBox.shrink();
                    }
                    return Padding(
                      padding: const EdgeInsets.only(top: 2),
                      child: Text(
                        'p2pperfsummary'.tr(args: ['$completedTrades', '${perf['completionRate']}']),
                        style: const TextStyle(color: Colors.black38, fontSize: 11),
                      ),
                    );
                  },
                ),
              ],
            ),
          ),
          if (offer.paymentMethod?.paymentChannel != null &&
              offer.paymentMethod!.paymentChannel!.isNotEmpty)
            Container(
              padding: const EdgeInsets.symmetric(
                horizontal: P2PTheme.space2,
                vertical: 2,
              ),
              decoration: BoxDecoration(
                color: P2PTheme.neutralBg,
                borderRadius: BorderRadius.circular(P2PTheme.space2),
              ),
              child: Text(
                offer.paymentMethod!.paymentChannel!,
                style: const TextStyle(fontSize: 11),
              ),
            ),
        ],
      ),
    );
  }
}

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
import 'package:trovo_app/widgets/p2p_merchant_gate.dart';

// P2PMyOffersView (Plan Section 95.11): a merchant's own offers, with
// online/offline toggle.
class P2PMyOffersView extends StatefulWidget {
  const P2PMyOffersView({Key? key}) : super(key: key);

  @override
  State<P2PMyOffersView> createState() => _P2PMyOffersViewState();
}

class _P2PMyOffersViewState extends State<P2PMyOffersView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  bool loading = true;
  List<P2POffer> offers = [];
  bool? merchantOnline;
  bool togglingMerchant = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
    _loadMerchantStatus();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final result = await api.listMyOffers();
    setState(() {
      offers = result;
      loading = false;
    });
  }

  Future<void> _loadMerchantStatus() async {
    final status = await api.getMerchantStatus();
    if (!mounted) return;
    setState(() => merchantOnline = status?['merchantOnline'] == true);
  }

  Future<void> _toggleMerchantOnline() async {
    if (merchantOnline == null) return;
    setState(() => togglingMerchant = true);
    await api.setMerchantOnlineStatus(!merchantOnline!);
    if (!mounted) return;
    setState(() {
      merchantOnline = !merchantOnline!;
      togglingMerchant = false;
    });
  }

  Future<void> _toggle(P2POffer offer) async {
    if (offer.isOnline) {
      await api.pauseOffer(offer.id!);
    } else {
      await api.activateOffer(offer.id!);
    }
    _load();
  }

  Future<void> _close(P2POffer offer) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('p2pclosethisoffer'.tr()),
        content: Text('p2pcannotbeundone'.tr()),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context, false), child: Text('cancel'.tr())),
          TextButton(onPressed: () => Navigator.pop(context, true), child: Text('p2pcloseofferbutton'.tr())),
        ],
      ),
    );
    if (confirmed != true) return;
    await api.closeOffer(offer.id!);
    _load();
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
        title: Text('p2pmyofferstitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              appState.currentAction = PageAction(
                state: PageState.addPage,
                page: P2PCreateOfferViewPageConfig,
              );
            },
          ),
        ],
      ),
      body: P2PMerchantGate(child: Column(
        children: [
          _merchantToggleCard(),
          Expanded(child: loading
          ? const Center(child: CircularProgressIndicator())
          : offers.isEmpty
              ? P2PEmptyState(
                  icon: Icons.storefront_outlined,
                  message: 'p2pnooffersyet'.tr(),
                  ctaLabel: 'p2pcreatefirstoffer'.tr(),
                  onCta: () {
                    appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: P2PCreateOfferViewPageConfig,
                    );
                  },
                )
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4, vertical: P2PTheme.space2),
                    itemCount: offers.length,
                    itemBuilder: (context, i) {
                      final o = offers[i];
                      final isClosed = o.status == 'CLOSED';
                      return P2PListCard(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Expanded(
                                  child: Column(
                                    crossAxisAlignment: CrossAxisAlignment.start,
                                    children: [
                                      Text('${o.offerType} ${o.asset}', style: const TextStyle(fontWeight: FontWeight.w700)),
                                      const SizedBox(height: P2PTheme.space1),
                                      Text('${o.price} ${o.currency}', style: const TextStyle(color: Colors.black54)),
                                      const SizedBox(height: P2PTheme.space1),
                                      Text(o.status ?? '', style: const TextStyle(fontSize: 12, color: Colors.black38)),
                                    ],
                                  ),
                                ),
                                if (!isClosed)
                                  Switch(
                                    value: o.isOnline,
                                    activeColor: P2PTheme.success,
                                    onChanged: (_) => _toggle(o),
                                  ),
                              ],
                            ),
                            if (!isClosed) ...[
                              const SizedBox(height: P2PTheme.space2),
                              Row(
                                children: [
                                  TextButton(
                                    onPressed: () {
                                      appState.viewData ??= {};
                                      appState.viewData![P2PEditOfferViewPageConfig.key] = {
                                        'offer': o,
                                      };
                                      appState.currentAction = PageAction(
                                        state: PageState.addPage,
                                        page: P2PEditOfferViewPageConfig,
                                      );
                                    },
                                    child: Text('edit'.tr()),
                                  ),
                                  TextButton(
                                    onPressed: () => _close(o),
                                    style: TextButton.styleFrom(foregroundColor: P2PTheme.danger),
                                    child: Text('close'.tr()),
                                  ),
                                ],
                              ),
                            ],
                          ],
                        ),
                      );
                    },
                  ),
                ),
          ),
        ],
      )),
    );
  }

  // Offline gets the same warning-tinted treatment as a danger banner
  // (not just a neutral card) so a merchant can't miss it - going
  // offline silently hides every one of their offers from search, which
  // is easy to forget about otherwise.
  Widget _merchantToggleCard() {
    if (merchantOnline == null) return const SizedBox.shrink();
    final online = merchantOnline!;
    final textColor = online ? P2PTheme.brandDark : P2PTheme.danger;
    return Padding(
      padding: const EdgeInsets.fromLTRB(P2PTheme.space4, P2PTheme.space2, P2PTheme.space4, 0),
      child: Container(
        padding: const EdgeInsets.all(P2PTheme.space4),
        decoration: BoxDecoration(
          color: online ? Colors.white : P2PTheme.danger.withOpacity(0.08),
          borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
          border: online ? null : Border.all(color: P2PTheme.danger.withOpacity(0.3)),
          boxShadow: online ? P2PTheme.cardShadow : null,
        ),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    online ? 'p2pyouareonline'.tr() : '⚠ ${'p2pyouareoffline'.tr()}',
                    style: TextStyle(fontWeight: FontWeight.w700, color: textColor),
                  ),
                  const SizedBox(height: P2PTheme.space1),
                  Text(
                    online ? 'p2ponlinehint'.tr() : 'p2pofflinehint'.tr(),
                    style: TextStyle(color: online ? Colors.black54 : P2PTheme.danger, fontSize: 13),
                  ),
                ],
              ),
            ),
            Switch(
              value: online,
              activeColor: P2PTheme.success,
              onChanged: togglingMerchant ? null : (_) => _toggleMerchantOnline(),
            ),
          ],
        ),
      ),
    );
  }
}

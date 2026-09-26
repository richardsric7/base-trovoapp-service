import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

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

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final result = await api.listMyOffers();
    setState(() {
      offers = result;
      loading = false;
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

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);

    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('My offers', style: TextStyle(color: notifier.getblck)),
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
      body: loading
          ? const Center(child: CircularProgressIndicator())
          : offers.isEmpty
              ? P2PEmptyState(
                  icon: Icons.storefront_outlined,
                  message: 'You have no offers yet.',
                  ctaLabel: 'Create your first offer',
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
                      return P2PListCard(
                        child: Row(
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
                            Switch(
                              value: o.isOnline,
                              activeColor: P2PTheme.success,
                              onChanged: (_) => _toggle(o),
                            ),
                          ],
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}

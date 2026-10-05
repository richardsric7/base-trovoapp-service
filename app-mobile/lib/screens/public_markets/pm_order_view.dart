import 'dart:async';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/public_market.dart';
import 'package:trovo_app/network/public_markets_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/screens/public_markets/pm_widgets.dart';
import 'package:trovo_app/storage/state.dart';

// PMOrderView follows an order to completion: its state, amounts and
// timeline, refreshed every few seconds until it is final.
// viewData: {orderId, fresh (just placed)}.
class PMOrderView extends StatefulWidget {
  const PMOrderView({Key? key}) : super(key: key);

  @override
  State<PMOrderView> createState() => _PMOrderViewState();
}

class _PMOrderViewState extends State<PMOrderView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  String id = '';
  bool fresh = false;
  PMOrder? order;
  Timer? timer;
  bool started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (!started) {
      started = true;
      final args = pmArgs(appState, PMOrderViewPageConfig);
      id = '${args['orderId'] ?? ''}';
      fresh = args['fresh'] == true;
      _load();
      timer = Timer.periodic(const Duration(seconds: 5), (_) {
        if (order?.done != true) _load();
      });
    }
  }

  @override
  void dispose() {
    timer?.cancel();
    super.dispose();
  }

  Future<void> _load() async {
    final o = await PublicMarketsApi(appState).order(id);
    if (mounted && o != null) setState(() => order = o);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    final o = order;
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('pmorder'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: o == null
          ? const Center(child: CircularProgressIndicator())
          : ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
              Center(
                child: Column(children: [
                  Icon(o.ok ? Icons.check_circle : (o.done ? Icons.cancel : Icons.hourglass_top), size: 64, color: pmStateColor(o.state)),
                  const SizedBox(height: P2PTheme.space3),
                  Text(o.ok ? (o.isBuy ? 'pmbuydone' : 'pmselldone').tr() : (o.done ? pmStateLabel(o.state) : 'pmplacing'.tr()),
                      style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w700, color: P2PTheme.brandDark)),
                  const SizedBox(height: 6),
                  Text(
                    o.isBuy ? '+ ${pmQty(o.quantity)} ${o.assetCode}' : '+ ${pmMoney(o.netAmount, o.fundingAssetCode)}',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.w600, color: o.ok ? P2PTheme.success : Colors.black54),
                  ),
                  if (!o.done) ...[
                    const SizedBox(height: 6),
                    Text(o.path == 'SLOW' ? 'pmslownote'.tr() : 'pmfastnote'.tr(), textAlign: TextAlign.center, style: const TextStyle(color: Colors.black54)),
                  ],
                ]),
              ),
              const SizedBox(height: P2PTheme.space6),
              PMSection(
                title: 'pmorderdetails'.tr(),
                child: Column(children: [
                  PMRow('pmorderid'.tr(), o.id),
                  PMRow('pmstatus'.tr(), pmStateLabel(o.state), valueColor: pmStateColor(o.state)),
                  PMRow(o.isBuy ? 'pmyoupaid'.tr() : 'pmgrossproceeds'.tr(), pmMoney(o.amount, o.fundingAssetCode)),
                  PMRow('pmtrovofee'.tr(), pmMoney(o.fee, o.fundingAssetCode)),
                  PMRow(o.isBuy ? 'pminvested'.tr() : 'pmyoureceive'.tr(), pmMoney(o.netAmount, o.fundingAssetCode)),
                  PMRow('pmquantity'.tr(), '${pmQty(o.quantity)} ${o.assetCode}'),
                  PMRow('pmreferenceprice'.tr(), pmMoney(o.referencePrice)),
                  if (o.executedPrice.isNotEmpty) PMRow('pmexecutedprice'.tr(), pmMoney(o.executedPrice)),
                  PMRow('pmwallet'.tr(), o.walletAlias.isNotEmpty ? o.walletAlias : o.walletAddress),
                  if (o.note.isNotEmpty) PMRow('pmnote'.tr(), o.note),
                ]),
              ),
              PMSection(
                title: 'pmtimeline'.tr(),
                child: Column(children: [
                  for (final e in o.events)
                    Padding(
                      padding: const EdgeInsets.only(bottom: P2PTheme.space3),
                      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
                        Icon(Icons.circle, size: 10, color: pmStateColor('${e['state']}')),
                        const SizedBox(width: P2PTheme.space3),
                        Expanded(
                          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                            Text(pmStateLabel('${e['state']}'), style: const TextStyle(fontWeight: FontWeight.w600)),
                            if ('${e['note'] ?? ''}'.isNotEmpty) Text('${e['note']}', style: const TextStyle(color: Colors.black54, fontSize: 12)),
                            Text(pmDateTime('${e['at']}'), style: const TextStyle(color: Colors.black38, fontSize: 11)),
                          ]),
                        ),
                      ]),
                    ),
                ]),
              ),
            ]),
      bottomNavigationBar: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(P2PTheme.space4),
          child: Row(children: [
            Expanded(
              child: OutlinedButton(
                onPressed: () => appState.currentAction = PageAction(state: PageState.replace, page: PMPortfolioViewPageConfig),
                child: Text('pmviewmystocks'.tr()),
              ),
            ),
            const SizedBox(width: P2PTheme.space3),
            Expanded(
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark),
                onPressed: () => appState.currentAction = PageAction(state: PageState.pop),
                child: Text('pmdone'.tr(), style: const TextStyle(color: Colors.white)),
              ),
            ),
          ]),
        ),
      ),
    );
  }
}

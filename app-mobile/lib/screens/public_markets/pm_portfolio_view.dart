import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/public_market.dart';
import 'package:trovo_app/network/public_markets_requests.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/screens/public_markets/pm_widgets.dart';
import 'package:trovo_app/storage/state.dart';

// PMPortfolioView is "My Stocks": the value of the user's Public Markets
// holdings across their wallets, returns and income, each holding, open
// orders and recent activity.
class PMPortfolioView extends StatefulWidget {
  const PMPortfolioView({Key? key}) : super(key: key);

  @override
  State<PMPortfolioView> createState() => _PMPortfolioViewState();
}

class _PMPortfolioViewState extends State<PMPortfolioView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  PMPortfolio? portfolio;
  bool loading = true;
  bool started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (!started) {
      started = true;
      _load();
    }
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final p = await PublicMarketsApi(appState).portfolio();
    if (!mounted) return;
    setState(() {
      portfolio = p;
      loading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    final p = portfolio;
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('pmmystocks'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
        actions: [
          TextButton(onPressed: () => pmOpen(appState, PMDividendsViewPageConfig, {'assetCode': ''}), child: Text('pmdividends'.tr())),
        ],
      ),
      body: loading && p == null
          ? const Center(child: CircularProgressIndicator())
          : p == null
              ? P2PEmptyState(icon: Icons.wifi_off, message: 'pmloadfailed'.tr(), ctaLabel: 'retry'.tr(), onCta: _load)
              : RefreshIndicator(onRefresh: _load, child: _content(p)),
    );
  }

  Widget _content(PMPortfolio p) {
    final gain = pmNum(p.totalReturn) >= 0;
    return ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
      Container(
        padding: const EdgeInsets.all(P2PTheme.space6),
        decoration: BoxDecoration(color: P2PTheme.brandDark, borderRadius: BorderRadius.circular(P2PTheme.cardRadius)),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text('pmportfoliovalue'.tr(), style: const TextStyle(color: Colors.white70)),
          const SizedBox(height: 6),
          Text(pmMoney(p.value), style: const TextStyle(color: Colors.white, fontSize: 26, fontWeight: FontWeight.w700)),
          const SizedBox(height: 6),
          Text('${gain ? '+' : ''}${pmMoney(p.totalReturn)} (${pmNum(p.returnPercent).toStringAsFixed(2)}%) · ${'pmtoday'.tr()} ${pmMoney(p.todayChange)}',
              style: TextStyle(color: gain ? const Color(0xFF7BE0A8) : const Color(0xFFFFA48A))),
          const SizedBox(height: 6),
          Text('${'pmincome'.tr()} ${pmMoney(p.income)}', style: const TextStyle(color: Colors.white70, fontSize: 12)),
        ]),
      ),
      const SizedBox(height: P2PTheme.space4),
      if (p.holdings.isEmpty)
        P2PEmptyState(icon: Icons.show_chart, message: 'pmnoholdings'.tr(), ctaLabel: 'pmexplore'.tr(), onCta: () => pmOpen(appState, PublicMarketsViewPageConfig))
      else
        for (final h in p.holdings)
          P2PListCard(
            onTap: () => pmOpen(appState, PMAssetViewPageConfig, {'assetCode': h.asset.assetCode}),
            child: Row(children: [
              PMLogo(asset: h.asset),
              const SizedBox(width: P2PTheme.space3),
              Expanded(
                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Text(h.asset.ticker, style: const TextStyle(fontWeight: FontWeight.w700, color: P2PTheme.brandDark)),
                  Text('${pmQty(h.quantity)} · ${'pmavg'.tr()} ${pmMoney(h.averageCost)}', style: const TextStyle(color: Colors.black54, fontSize: 12)),
                ]),
              ),
              Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
                Text(pmMoney(h.marketValue), style: const TextStyle(fontWeight: FontWeight.w600, color: P2PTheme.brandDark)),
                Text('${pmNum(h.returnPercent) >= 0 ? '+' : ''}${pmNum(h.returnPercent).toStringAsFixed(2)}%',
                    style: TextStyle(fontSize: 12, color: pmNum(h.returnPercent) >= 0 ? P2PTheme.success : P2PTheme.danger)),
              ]),
            ]),
          ),
      if (p.openOrders.isNotEmpty)
        PMSection(
          title: 'pmopenorders'.tr(),
          child: Column(children: [
            for (final o in p.openOrders)
              ListTile(
                contentPadding: EdgeInsets.zero,
                onTap: () => pmOpen(appState, PMOrderViewPageConfig, {'orderId': o.id}),
                title: Text('${o.isBuy ? 'pmbuy'.tr() : 'pmsell'.tr()} ${o.assetCode}', style: const TextStyle(fontWeight: FontWeight.w600)),
                subtitle: Text(o.isBuy ? pmMoney(o.amount, o.fundingAssetCode) : '${pmQty(o.quantity)} ${o.assetCode}'),
                trailing: PMBadge(pmStateLabel(o.state), color: pmStateColor(o.state)),
              ),
          ]),
        ),
      if (p.activity.isNotEmpty)
        PMSection(
          title: 'pmactivity'.tr(),
          child: Column(children: [
            for (final a in p.activity.take(30))
              ListTile(
                contentPadding: EdgeInsets.zero,
                onTap: () {
                  final ref = '${a['reference'] ?? ''}';
                  if (a['kind'] == 'BUY' || a['kind'] == 'SELL') pmOpen(appState, PMOrderViewPageConfig, {'orderId': ref});
                },
                title: Text('${a['title']}', style: const TextStyle(fontWeight: FontWeight.w600)),
                subtitle: Text('${a['detail'] ?? ''}\n${pmDateTime('${a['at']}')}', style: const TextStyle(fontSize: 12)),
                isThreeLine: true,
                trailing: Text(
                  '${a['received'] ?? ''}'.isNotEmpty ? '+ ${a['received']}' : '− ${a['spent'] ?? ''}',
                  style: TextStyle(color: '${a['received'] ?? ''}'.isNotEmpty ? P2PTheme.success : P2PTheme.brandDark, fontWeight: FontWeight.w600),
                ),
              ),
          ]),
        ),
    ]);
  }
}

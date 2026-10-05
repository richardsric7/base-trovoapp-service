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

// PMAssetView is an asset's page: price and chart, the user's position,
// key statistics, how the tokens are owned and backed, the trading session
// and corporate actions, with Buy / Sell. viewData: {assetCode}.
class PMAssetView extends StatefulWidget {
  const PMAssetView({Key? key}) : super(key: key);

  @override
  State<PMAssetView> createState() => _PMAssetViewState();
}

class _PMAssetViewState extends State<PMAssetView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late PublicMarketsApi api;
  String code = '';
  PMAsset? asset;
  PMHolding? holding;
  List<Map> points = [];
  String range = '1D';
  bool loading = true;
  bool started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    appState = Provider.of<DataProvider>(context, listen: false);
    api = PublicMarketsApi(appState);
    if (!started) {
      started = true;
      code = '${pmArgs(appState, PMAssetViewPageConfig)['assetCode'] ?? ''}';
      _load();
    }
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final results = await Future.wait([api.asset(code), api.prices(code, range), api.portfolio()]);
    if (!mounted) return;
    final p = results[2] as PMPortfolio?;
    setState(() {
      asset = results[0] as PMAsset?;
      points = results[1] as List<Map>;
      holding = p?.holdings.where((h) => h.asset.assetCode == asset?.assetCode).firstOrNull;
      loading = false;
    });
  }

  Future<void> _range(String r) async {
    setState(() => range = r);
    final pts = await api.prices(code, r);
    if (mounted) setState(() => points = pts);
  }

  void _trade(bool buy) => pmOpen(appState, PMTradeViewPageConfig, {'assetCode': asset!.assetCode, 'side': buy ? 'buy' : 'sell'});

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    final a = asset;
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text(a?.ticker ?? code, style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: loading && a == null
          ? const Center(child: CircularProgressIndicator())
          : a == null
              ? P2PEmptyState(icon: Icons.error_outline, message: 'pmassetnotfound'.tr(), ctaLabel: 'retry'.tr(), onCta: _load)
              : RefreshIndicator(onRefresh: _load, child: _content(a)),
      bottomNavigationBar: a == null ? null : _actions(a),
    );
  }

  Widget _actions(PMAsset a) {
    final held = holding != null && pmNum(holding!.quantity) > 0;
    return SafeArea(
      child: Padding(
        padding: const EdgeInsets.all(P2PTheme.space4),
        child: a.open
            ? Row(children: [
                if (held) ...[
                  Expanded(
                    child: OutlinedButton(
                      style: OutlinedButton.styleFrom(padding: const EdgeInsets.symmetric(vertical: 14)),
                      onPressed: () => _trade(false),
                      child: Text('pmsell'.tr()),
                    ),
                  ),
                  const SizedBox(width: P2PTheme.space3),
                ],
                Expanded(
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark, padding: const EdgeInsets.symmetric(vertical: 14)),
                    onPressed: () => _trade(true),
                    child: Text(held ? 'pmbuymore'.tr() : 'pmbuy'.tr(), style: const TextStyle(color: Colors.white)),
                  ),
                ),
              ])
            : Text(a.status == 'halted' ? 'pmhaltednote'.tr() : 'pmcomingsoonnote'.tr(), textAlign: TextAlign.center, style: const TextStyle(color: Colors.black54)),
      ),
    );
  }

  Widget _content(PMAsset a) {
    final session = a.session;
    final custody = a.custody;
    final remaining = (session['minutesRemaining'] as num?)?.toInt();
    final total = (session['minutesTotal'] as num?)?.toInt();
    return ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
      Row(children: [
        PMLogo(asset: a, size: 52),
        const SizedBox(width: P2PTheme.space3),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(a.name, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16, color: P2PTheme.brandDark)),
            const SizedBox(height: 4),
            Row(children: [PMBadge(a.market), const SizedBox(width: 6), PMBadge(a.isBond ? 'pmbond'.tr() : 'pmequity'.tr(), color: Colors.black54)]),
          ]),
        ),
      ]),
      const SizedBox(height: P2PTheme.space4),
      Text(pmNum(a.price) > 0 ? pmMoney(a.price) : '—', style: const TextStyle(fontSize: 28, fontWeight: FontWeight.w700, color: P2PTheme.brandDark)),
      Row(children: [
        PMChange(percent: a.dayChangePercent, fontSize: 14),
        const SizedBox(width: 8),
        Text(a.priceLive ? 'pmpricelive'.tr(args: [a.market]) : 'pmpricelastclose'.tr(), style: const TextStyle(color: Colors.black45, fontSize: 12)),
      ]),
      const SizedBox(height: P2PTheme.space3),
      PMChart(points: points),
      Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
        for (final r in const ['1D', '1W', '1M', '3M', '1Y', 'All'])
          TextButton(
            onPressed: () => _range(r),
            style: TextButton.styleFrom(
              backgroundColor: range == r ? P2PTheme.brandDark : Colors.transparent,
              minimumSize: const Size(44, 32),
            ),
            child: Text(r, style: TextStyle(color: range == r ? Colors.white : P2PTheme.brandDark, fontSize: 12)),
          ),
      ]),
      const SizedBox(height: P2PTheme.space4),
      if (holding != null && pmNum(holding!.quantity) > 0)
        PMSection(
          title: 'pmyourposition'.tr(),
          trailing: TextButton(
            onPressed: () => pmOpen(appState, PMDividendsViewPageConfig, {'assetCode': a.assetCode}),
            child: Text('pmdividends'.tr()),
          ),
          child: Column(children: [
            PMRow('pmholding'.tr(), '${pmQty(holding!.quantity)} ${a.ticker}'),
            PMRow('pmmarketvalue'.tr(), pmMoney(holding!.marketValue)),
            PMRow('pmaveragecost'.tr(), pmMoney(holding!.averageCost)),
            PMRow('pmtotalreturn'.tr(), '${pmMoney(holding!.totalReturn)} (${pmNum(holding!.returnPercent).toStringAsFixed(2)}%)',
                valueColor: pmNum(holding!.totalReturn) >= 0 ? P2PTheme.success : P2PTheme.danger),
            PMRow(a.isBond ? 'pminterestreceived'.tr() : 'pmdividendsreceived'.tr(), pmMoney(holding!.incomeReceived)),
          ]),
        ),
      PMSection(
        title: 'pmkeystats'.tr(),
        child: Column(children: [
          PMRow('pmpreviousclose'.tr(), pmMoney(a.previousClose)),
          PMRow('pmdayrange'.tr(), pmNum(a.dayLow) > 0 ? '${pmQty(a.dayLow)} – ${pmQty(a.dayHigh)}' : '—'),
          if (a.dayVolume.isNotEmpty) PMRow('pmvolume'.tr(), a.dayVolume),
          if (a.marketCap.isNotEmpty) PMRow('pmmarketcap'.tr(), a.marketCap),
          if (a.peRatio.isNotEmpty) PMRow('pmperatio'.tr(), a.peRatio),
          if (a.dividendYield.isNotEmpty) PMRow('pmdividendyield'.tr(), a.dividendYield),
          if (a.coupon.isNotEmpty) PMRow('pmcoupon'.tr(), a.coupon),
          if (a.maturityDate.isNotEmpty) PMRow('pmmaturity'.tr(), pmDate(a.maturityDate)),
          PMRow('pmminimumbuy'.tr(), pmMoney(a.minimumBuy, a.fundingAsset)),
          PMRow('pmtrovofee'.tr(), '${a.feePercent}%'),
        ]),
      ),
      if (a.description.isNotEmpty) PMSection(title: 'pmabout'.tr(), child: Text(a.description, style: const TextStyle(height: 1.4))),
      PMSection(
        title: 'pmhowyouown'.tr(),
        child: Column(children: [
          _chain(Icons.account_balance_wallet_outlined, 'pmownyou'.tr(), 'pmownyoudetail'.tr(args: [a.ticker])),
          _chain(Icons.token_outlined, 'Trovotech', 'pmowntrovotech'.tr()),
          _chain(Icons.account_balance_outlined, '${custody['custodian'] ?? ''} · ${custody['nominee'] ?? ''}', 'pmowncustodian'.tr()),
          _chain(Icons.storefront_outlined, a.market == 'NGX' ? 'Nigerian Exchange (NGX)' : 'FMDQ Securities Exchange', 'pmownexchange'.tr()),
          const SizedBox(height: P2PTheme.space2),
          Text('pmbeneficialnote'.tr(), style: const TextStyle(color: Colors.black54, fontSize: 12)),
        ]),
      ),
      PMSection(
        title: 'pmcustodybacking'.tr(),
        child: Column(children: [
          if (a.tokensInCirculation.isNotEmpty && custody['unitsHeld'] != null)
            Text('pmbackedline'.tr(args: [pmQty(a.tokensInCirculation), a.ticker, pmQty(custody['unitsHeld']), '${custody['custodian'] ?? ''}']),
                style: const TextStyle(height: 1.4)),
          if (custody['lastReconciliation'] is Map)
            PMRow('pmlastreconciled'.tr(), pmDateTime('${custody['lastReconciliation']['at']}'),
                valueColor: custody['lastReconciliation']['result'] == 'MATCHED' ? P2PTheme.success : P2PTheme.danger),
          PMRow('pmdepository'.tr(), '${custody['depository'] ?? ''}'),
          PMRow('pmsettlement'.tr(), '${custody['settlement'] ?? ''}'),
          PMRow('pmdealingmember'.tr(), '${custody['dealingMember'] ?? ''}'),
        ]),
      ),
      PMSection(
        title: 'pmtradingsession'.tr(),
        child: session['open'] == true && remaining != null && total != null
            ? Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('pmsessionremaining'.tr(args: ['${remaining ~/ 60}h ${remaining % 60}m', a.market])),
                const SizedBox(height: 8),
                LinearProgressIndicator(value: total > 0 ? 1 - remaining / total : 0, color: P2PTheme.success, backgroundColor: P2PTheme.neutralBg),
              ])
            : Text('pmsessionclosed'.tr(args: [pmDateTime('${session['opensAt'] ?? ''}')]), style: const TextStyle(color: Colors.black54)),
      ),
      if (a.corporateActions.isNotEmpty)
        PMSection(
          title: 'pmcorporateactions'.tr(),
          child: Column(children: [
            for (final c in a.corporateActions)
              PMRow('${c['eventType']} · ${'pmrecorddate'.tr()} ${c['recordDate']}',
                  pmNum(c['amountPerUnit']) > 0 ? '${pmNaira(c['amountPerUnit'])} / ${'pmunit'.tr()} · ${c['status']}' : '${c['status']}'),
          ]),
        ),
      PMSection(
        title: 'pmassetinfo'.tr(),
        child: Column(children: [
          PMRow('ISIN', a.isin),
          PMRow('pmsector'.tr(), a.sector),
          PMRow('pmunitis'.tr(), a.unitDescription),
          PMRow('pmtokencode'.tr(), a.assetCode),
          PMRow('pmtokencontract'.tr(), a.contractAddress),
        ]),
      ),
      Text('pmriskdisclosure'.tr(), style: const TextStyle(color: Colors.black45, fontSize: 12, height: 1.4)),
      const SizedBox(height: P2PTheme.space8),
    ]);
  }

  Widget _chain(IconData icon, String title, String detail) => Padding(
        padding: const EdgeInsets.only(bottom: P2PTheme.space3),
        child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
          CircleAvatar(radius: 16, backgroundColor: P2PTheme.neutralBg, child: Icon(icon, size: 18, color: P2PTheme.brandDark)),
          const SizedBox(width: P2PTheme.space3),
          Expanded(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(title, style: const TextStyle(fontWeight: FontWeight.w600, color: P2PTheme.brandDark)),
              Text(detail, style: const TextStyle(color: Colors.black54, fontSize: 12)),
            ]),
          ),
        ]),
      );
}

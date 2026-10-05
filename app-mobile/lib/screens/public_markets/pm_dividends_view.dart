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

// PMDividendsView lists the user's Public Markets dividends and coupons
// (paid and coming), each with its breakdown: units held on the record
// date, gross, withholding tax and what was paid. viewData: {assetCode}
// ('' for all assets).
class PMDividendsView extends StatefulWidget {
  const PMDividendsView({Key? key}) : super(key: key);

  @override
  State<PMDividendsView> createState() => _PMDividendsViewState();
}

class _PMDividendsViewState extends State<PMDividendsView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  String assetCode = '';
  List<PMDividend>? dividends;
  bool loading = true;
  bool started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (!started) {
      started = true;
      assetCode = '${pmArgs(appState, PMDividendsViewPageConfig)['assetCode'] ?? ''}';
      _load();
    }
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final r = await PublicMarketsApi(appState).dividends(assetCode: assetCode);
    if (!mounted) return;
    setState(() {
      dividends = r;
      loading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    final list = dividends ?? [];
    final total = list.where((d) => d.paid).fold<double>(0, (s, d) => s + pmNum(d.netAmount));
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text(assetCode.isEmpty ? 'pmdividendhistory'.tr() : '${'pmdividendhistory'.tr()} · $assetCode', style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: loading && dividends == null
          ? const Center(child: CircularProgressIndicator())
          : dividends == null
              ? P2PEmptyState(icon: Icons.wifi_off, message: 'pmloadfailed'.tr(), ctaLabel: 'retry'.tr(), onCta: _load)
              : list.isEmpty
                  ? P2PEmptyState(icon: Icons.payments_outlined, message: 'pmnodividends'.tr())
                  : RefreshIndicator(
                      onRefresh: _load,
                      child: ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
                        PMSection(
                          title: 'pmtotalreceived'.tr(),
                          child: Text(pmMoney(total, 'CNGN'), style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w700, color: P2PTheme.success)),
                        ),
                        for (final d in list)
                          P2PListCard(
                            onTap: () => _detail(d),
                            child: Row(children: [
                              Expanded(
                                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                                  Text('${d.assetCode} · ${d.description.isNotEmpty ? d.description : d.eventType.toLowerCase()}',
                                      style: const TextStyle(fontWeight: FontWeight.w600, color: P2PTheme.brandDark)),
                                  const SizedBox(height: 4),
                                  Text('pmdivline'.tr(args: [pmNaira(d.amountPerUnit), pmQty(d.units), d.assetCode, pmNaira(d.whtAmount)]),
                                      style: const TextStyle(color: Colors.black54, fontSize: 12)),
                                ]),
                              ),
                              Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
                                Text('+ ${pmMoney(d.netAmount, 'CNGN')}', style: TextStyle(fontWeight: FontWeight.w700, color: d.paid ? P2PTheme.success : Colors.black45)),
                                Text(d.paid ? pmDate(d.paidAt) : '${'pmpayable'.tr()} ${d.payDate}', style: const TextStyle(fontSize: 11, color: Colors.black45)),
                              ]),
                            ]),
                          ),
                      ]),
                    ),
    );
  }

  void _detail(PMDividend d) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(P2PTheme.sheetRadius))),
      builder: (context) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(P2PTheme.space6),
          child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(d.paid ? 'pmdivreceived'.tr(args: [pmMoney(d.netAmount, 'CNGN')]) : 'pmdivcoming'.tr(args: [pmMoney(d.netAmount, 'CNGN')]),
                style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w700, color: P2PTheme.brandDark)),
            const SizedBox(height: 6),
            Text('pmdivexplain'.tr(args: [d.assetCode, d.description.isNotEmpty ? d.description : d.eventType.toLowerCase(), pmQty(d.units)]),
                style: const TextStyle(color: Colors.black54)),
            const SizedBox(height: P2PTheme.space4),
            PMRow('pmperunit'.tr(), pmNaira(d.amountPerUnit)),
            PMRow('pmrecorddateunits'.tr(args: [d.recordDate]), '${pmQty(d.units)} ${d.assetCode}'),
            PMRow('pmgross'.tr(), pmNaira(d.grossAmount)),
            PMRow('pmwht'.tr(args: [d.whtPercent]), '− ${pmNaira(d.whtAmount)}', valueColor: P2PTheme.danger),
            PMRow(d.paid ? 'pmpaidto'.tr() : 'pmtobepaid'.tr(), pmMoney(d.netAmount, 'CNGN'), valueColor: P2PTheme.success),
            PMRow('pmpaydate'.tr(), d.paid ? pmDate(d.paidAt) : d.payDate),
            if (d.txHash.isNotEmpty) PMRow('pmtransaction'.tr(), d.txHash),
            const SizedBox(height: P2PTheme.space3),
            Text('pmdivhow'.tr(args: [d.custodianName]), style: const TextStyle(color: Colors.black45, fontSize: 12, height: 1.4)),
          ]),
        ),
      ),
    );
  }
}

import 'dart:async';

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

// PublicMarketsView lists tokenized NGX equities and FMDQ bonds, with a
// search and the prototype's filters (all, equities, bonds, top gainers).
class PublicMarketsView extends StatefulWidget {
  const PublicMarketsView({Key? key}) : super(key: key);

  @override
  State<PublicMarketsView> createState() => _PublicMarketsViewState();
}

class _PublicMarketsViewState extends State<PublicMarketsView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  List<PMAsset>? assets;
  bool loading = true;
  bool started = false;
  String filter = 'all';
  final search = TextEditingController();
  Timer? debounce;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    appState = Provider.of<DataProvider>(context, listen: false);
    if (!started) {
      started = true;
      _load();
    }
  }

  @override
  void dispose() {
    debounce?.cancel();
    search.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final type = filter == 'eq' ? 'EQUITY' : (filter == 'bd' ? 'BOND' : '');
    final r = await PublicMarketsApi(appState).assets(type: type, search: search.text.trim());
    if (!mounted) return;
    setState(() {
      assets = r;
      loading = false;
    });
  }

  List<PMAsset> get _shown {
    final list = List<PMAsset>.from(assets ?? []);
    if (filter == 'gain') {
      list.sort((a, b) => pmNum(b.dayChangePercent).compareTo(pmNum(a.dayChangePercent)));
      return list.where((a) => pmNum(a.dayChangePercent) > 0).toList();
    }
    return list;
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('pmtitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
        actions: [
          IconButton(
            tooltip: 'pmmystocks'.tr(),
            icon: const Icon(Icons.pie_chart_outline),
            onPressed: () => pmOpen(appState, PMPortfolioViewPageConfig),
          ),
        ],
      ),
      body: Column(children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(P2PTheme.space4, P2PTheme.space3, P2PTheme.space4, 0),
          child: TextField(
            controller: search,
            onChanged: (_) {
              debounce?.cancel();
              debounce = Timer(const Duration(milliseconds: 400), _load);
            },
            decoration: InputDecoration(
              hintText: 'pmsearch'.tr(),
              prefixIcon: const Icon(Icons.search),
              filled: true,
              fillColor: Colors.white,
              border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
            ),
          ),
        ),
        SizedBox(
          height: 52,
          child: ListView(
            scrollDirection: Axis.horizontal,
            padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4, vertical: P2PTheme.space2),
            children: [
              for (final f in const [['all', 'pmfilterall'], ['eq', 'pmfilterequities'], ['bd', 'pmfilterbonds'], ['gain', 'pmfiltergainers']])
                Padding(
                  padding: const EdgeInsets.only(right: P2PTheme.space2),
                  child: ChoiceChip(
                    label: Text(f[1].tr()),
                    selected: filter == f[0],
                    selectedColor: P2PTheme.brandDark,
                    labelStyle: TextStyle(color: filter == f[0] ? Colors.white : P2PTheme.brandDark),
                    onSelected: (_) {
                      setState(() => filter = f[0]);
                      _load();
                    },
                  ),
                ),
            ],
          ),
        ),
        Expanded(child: _list()),
      ]),
    );
  }

  Widget _list() {
    if (loading && assets == null) return const Center(child: CircularProgressIndicator());
    if (assets == null) {
      return P2PEmptyState(icon: Icons.wifi_off, message: 'pmloadfailed'.tr(), ctaLabel: 'retry'.tr(), onCta: _load);
    }
    final shown = _shown;
    if (shown.isEmpty) return P2PEmptyState(icon: Icons.show_chart, message: 'pmnoassets'.tr());
    return RefreshIndicator(
      onRefresh: _load,
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4),
        itemCount: shown.length,
        itemBuilder: (context, i) => PMAssetTile(asset: shown[i], onTap: () => pmOpen(appState, PMAssetViewPageConfig, {'assetCode': shown[i].assetCode})),
      ),
    );
  }
}

// PMAssetTile is an asset in a list (also on the home screen).
class PMAssetTile extends StatelessWidget {
  final PMAsset asset;
  final VoidCallback onTap;
  const PMAssetTile({super.key, required this.asset, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return P2PListCard(
      onTap: onTap,
      child: Row(children: [
        PMLogo(asset: asset),
        const SizedBox(width: P2PTheme.space3),
        Expanded(
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(children: [
              Text(asset.ticker, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 15, color: P2PTheme.brandDark)),
              const SizedBox(width: 6),
              PMBadge(asset.market, color: asset.market == 'FMDQ' ? P2PTheme.warning : P2PTheme.primary),
              if (asset.status == 'halted') ...[const SizedBox(width: 6), PMBadge('pmhalted'.tr(), color: P2PTheme.danger)],
              if (asset.status == 'coming-soon') ...[const SizedBox(width: 6), PMBadge('pmcomingsoon'.tr(), color: Colors.black45)],
            ]),
            const SizedBox(height: 2),
            Text(asset.shortName, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(color: Colors.black54, fontSize: 13)),
          ]),
        ),
        Column(crossAxisAlignment: CrossAxisAlignment.end, children: [
          Text(pmNum(asset.price) > 0 ? pmMoney(asset.price) : '—', style: const TextStyle(fontWeight: FontWeight.w600, color: P2PTheme.brandDark)),
          const SizedBox(height: 2),
          PMChange(percent: asset.dayChangePercent),
        ]),
      ]),
    );
  }
}

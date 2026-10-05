import 'dart:async';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/public_market.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/public_markets_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/screens/public_markets/pm_widgets.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

// PMTradeView buys (an amount of the funding stablecoin, fee included) or
// sells (a quantity of tokens) from one of the user's own wallets. The
// quote says whether it fills now (Custodian inventory, or netted against
// the day's demand) or at the next session. Review builds the wallet
// operation, the user confirms it with biometrics, the signed operation is
// submitted and the order's status page opens. viewData: {assetCode, side}.
class PMTradeView extends StatefulWidget {
  const PMTradeView({Key? key}) : super(key: key);

  @override
  State<PMTradeView> createState() => _PMTradeViewState();
}

class _PMTradeViewState extends State<PMTradeView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late PublicMarketsApi api;
  bool started = false;
  bool buy = true;
  String code = '';
  PMAsset? asset;
  PMHolding? holding;
  Wallet? wallet;
  PMQuote? quote;
  String quoteError = '';
  bool quoting = false;
  bool busy = false;
  final input = TextEditingController();
  Timer? debounce;

  // the user's own wallets (shared wallets with approvers cannot trade here)
  List<Wallet> get wallets => (appState.userInfo?.wallets ?? []).where((w) => w.owner == null && (w.sharedAccessEnabled ?? 0) == 0).toList();

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    appState = Provider.of<DataProvider>(context, listen: false);
    api = PublicMarketsApi(appState);
    if (!started) {
      started = true;
      final args = pmArgs(appState, PMTradeViewPageConfig);
      code = '${args['assetCode'] ?? ''}';
      buy = args['side'] != 'sell';
      wallet = wallets.where((w) => w.primaryWallet == 1).firstOrNull ?? wallets.firstOrNull;
      _load();
    }
  }

  @override
  void dispose() {
    debounce?.cancel();
    input.dispose();
    super.dispose();
  }

  Future<void> _load() async {
    final results = await Future.wait([api.asset(code), api.portfolio()]);
    if (!mounted) return;
    setState(() {
      asset = results[0] as PMAsset?;
      holding = (results[1] as PMPortfolio?)?.holdings.where((h) => h.asset.assetCode == asset?.assetCode).firstOrNull;
    });
  }

  double get _fundingBalance {
    final code = asset?.fundingAsset ?? 'CNGN';
    final a = wallet?.claimedAssets?.where((x) => (x.assetCode ?? '').toUpperCase() == code).firstOrNull;
    return a?.amount ?? 0;
  }

  double get _heldInWallet => pmNum(holding?.wallets[wallet?.address] ?? holding?.wallets[wallet?.address?.toLowerCase()]);

  void _changed() {
    debounce?.cancel();
    debounce = Timer(const Duration(milliseconds: 450), _quote);
  }

  Future<void> _quote() async {
    final v = input.text.replaceAll(',', '').trim();
    if (pmNum(v) <= 0) {
      setState(() {
        quote = null;
        quoteError = '';
      });
      return;
    }
    setState(() => quoting = true);
    final q = await api.quote(code, buy: buy, value: v);
    if (!mounted) return;
    setState(() {
      quoting = false;
      quote = q;
      quoteError = q == null ? 'pmquotefailed'.tr() : '';
    });
  }

  void _preset(double share) {
    final max = buy ? _fundingBalance : _heldInWallet;
    if (max <= 0) return;
    final v = buy ? (max * share).floorToDouble() : max * share;
    input.text = buy ? v.toStringAsFixed(0) : v.toStringAsFixed(asset?.tokenDecimals ?? 4).replaceFirst(RegExp(r'\.?0+$'), '');
    _quote();
  }

  String? _problem() {
    final v = pmNum(input.text.replaceAll(',', ''));
    if (v <= 0) return buy ? 'pmenteramount'.tr() : 'pmenterquantity'.tr();
    if (wallet == null) return 'pmnowallet'.tr();
    if (buy && pmNum(asset?.minimumBuy) > 0 && v < pmNum(asset?.minimumBuy)) return 'pmbelowminimum'.tr(args: [pmMoney(asset!.minimumBuy, asset!.fundingAsset)]);
    if (buy && _fundingBalance > 0 && v > _fundingBalance) return 'pminsufficient'.tr(args: [asset!.fundingAsset]);
    if (!buy && v > _heldInWallet) return 'pmnotenoughtokens'.tr();
    return null;
  }

  Future<void> _review() async {
    final problem = _problem();
    if (problem != null) {
      popup(context, title: 'pmcheckorder'.tr(), message: problem);
      return;
    }
    final value = input.text.replaceAll(',', '').trim();
    setState(() => busy = true);
    showLoader(context);
    try {
      final built = await api.trade(code, buy: buy, walletAddress: wallet!.address!, value: value);
      hideLoader(context);
      final tx = built['data'] is Map ? built['data']['transaction'] : null;
      if (built['statusCode'] != 200 || tx == null || '$tx'.isEmpty) {
        setState(() => busy = false);
        popup(context, title: 'pmorderfailed'.tr(), message: PublicMarketsApi.message(built, 'p2ppleasetryagain'.tr()));
        return;
      }
      final q = PMQuote(Map.from(built['data']['quote'] ?? {}));
      final messages = ((built['data']['messages'] ?? []) as List).map((m) => '$m').toList();
      if (!await _confirm(q, messages)) {
        setState(() => busy = false);
        return;
      }
      showLoader(context);
      final signature = TrovoWalletSDK().signBase64Txn(appState.secretKeys[0], '$tx', '');
      final sent = await api.submit(code, buy: buy, walletAddress: wallet!.address!, value: value, transaction: '$tx', signature: signature);
      hideLoader(context);
      setState(() => busy = false);
      final order = sent['data'] is Map ? sent['data']['order'] : null;
      if (sent['statusCode'] == 200 && order is Map) {
        appState.viewData ??= {};
        appState.viewData![PMOrderViewPageConfig.key] = {'orderId': order['id'], 'fresh': true};
        appState.currentAction = PageAction(state: PageState.replace, page: PMOrderViewPageConfig);
      } else {
        popup(context, title: 'pmorderfailed'.tr(), message: PublicMarketsApi.message(sent, 'p2ppleasetryagain'.tr()));
      }
    } catch (e) {
      hideLoader(context);
      setState(() => busy = false);
      popup(context, title: 'error'.tr(), message: e.toString());
    }
  }

  // _confirm shows the order review and asks for biometrics.
  Future<bool> _confirm(PMQuote q, List<String> messages) async {
    final a = asset!;
    final ok = await showModalBottomSheet<bool>(
          context: context,
          isScrollControlled: true,
          shape: const RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(P2PTheme.sheetRadius))),
          builder: (context) => SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(P2PTheme.space6),
              child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(buy ? 'pmreviewbuy'.tr(args: [a.ticker]) : 'pmreviewsell'.tr(args: [a.ticker]),
                    style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w700, color: P2PTheme.brandDark)),
                const SizedBox(height: P2PTheme.space4),
                if (buy) ...[
                  PMRow('pmyoupay'.tr(), pmMoney(q.amount, q.fundingAsset)),
                  PMRow('pmtrovofee'.tr(), '${pmMoney(q.fee, q.fundingAsset)} (${q.feePercent}%)'),
                  PMRow('pmestimatedquantity'.tr(), '${pmQty(q.quantity)} ${a.ticker}'),
                ] else ...[
                  PMRow('pmyousell'.tr(), '${pmQty(q.quantity)} ${a.ticker}'),
                  PMRow('pmtrovofee'.tr(), '${pmMoney(q.fee, q.fundingAsset)} (${q.feePercent}%)'),
                  PMRow('pmyoureceive'.tr(), pmMoney(q.netAmount, q.fundingAsset)),
                ],
                PMRow('pmreferenceprice'.tr(), '${pmMoney(q.price)} · ${q.priceSource}'),
                PMRow('pmwallet'.tr(), wallet?.alias ?? wallet?.address ?? ''),
                const SizedBox(height: P2PTheme.space3),
                _pathNote(q),
                for (final m in messages) Padding(padding: const EdgeInsets.only(top: 6), child: Text(m, style: const TextStyle(fontSize: 12, color: Colors.black54))),
                const SizedBox(height: P2PTheme.space4),
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark, padding: const EdgeInsets.symmetric(vertical: 14)),
                    onPressed: () => Navigator.of(context).pop(true),
                    child: Text('pmauthorize'.tr(), style: const TextStyle(color: Colors.white)),
                  ),
                ),
                TextButton(onPressed: () => Navigator.of(context).pop(false), child: Center(child: Text('cancel'.tr()))),
              ]),
            ),
          ),
        ) ??
        false;
    if (!ok) return false;
    try {
      final authenticator = Authenticator();
      if (await authenticator.canCheckBiometrics()) return await authenticator.authenticateMe();
    } catch (_) {}
    return true;
  }

  Widget _pathNote(PMQuote q) {
    final fast = q.path == 'FAST' || q.path == 'NETTED';
    return Container(
      padding: const EdgeInsets.all(P2PTheme.space3),
      decoration: BoxDecoration(color: (fast ? P2PTheme.success : P2PTheme.warning).withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12)),
      child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Icon(fast ? Icons.bolt : Icons.schedule, size: 18, color: fast ? P2PTheme.success : P2PTheme.warning),
        const SizedBox(width: 8),
        Expanded(child: Text(q.note.isNotEmpty ? q.note : q.settlementNote, style: const TextStyle(fontSize: 12, height: 1.4))),
      ]),
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    final a = asset;
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text(a == null ? code : (buy ? 'pmbuyname' : 'pmsellname').tr(args: [a.shortName]), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: a == null
          ? const Center(child: CircularProgressIndicator())
          : ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
              SegmentedButton<bool>(
                segments: [ButtonSegment(value: true, label: Text('pmbuy'.tr())), ButtonSegment(value: false, label: Text('pmsell'.tr()))],
                selected: {buy},
                onSelectionChanged: (s) {
                  setState(() {
                    buy = s.first;
                    input.clear();
                    quote = null;
                  });
                },
              ),
              const SizedBox(height: P2PTheme.space4),
              PMSection(
                title: 'pmwallet'.tr(),
                child: DropdownButton<String>(
                  isExpanded: true,
                  value: wallet?.address,
                  underline: const SizedBox(),
                  items: [
                    for (final w in wallets)
                      DropdownMenuItem(value: w.address, child: Text(w.alias?.isNotEmpty == true ? w.alias! : (w.primaryWallet == 1 ? 'Main Wallet' : '${w.address}'))),
                  ],
                  onChanged: (v) => setState(() {
                    wallet = wallets.firstWhere((w) => w.address == v);
                    _changed();
                  }),
                ),
              ),
              PMSection(
                title: buy ? 'pmamount'.tr(args: [a.fundingAsset]) : 'pmquantityof'.tr(args: [a.ticker]),
                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  TextField(
                    controller: input,
                    keyboardType: const TextInputType.numberWithOptions(decimal: true),
                    inputFormatters: [FilteringTextInputFormatter.allow(RegExp(r'[0-9.,]'))],
                    onChanged: (_) => _changed(),
                    style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w700, color: P2PTheme.brandDark),
                    decoration: const InputDecoration(border: InputBorder.none, hintText: '0'),
                  ),
                  Text(
                    buy ? 'pmavailable'.tr(args: [pmMoney(_fundingBalance, a.fundingAsset)]) : 'pmavailable'.tr(args: ['${pmQty(_heldInWallet)} ${a.ticker}']),
                    style: const TextStyle(color: Colors.black54, fontSize: 12),
                  ),
                  const SizedBox(height: P2PTheme.space2),
                  Row(children: [
                    for (final p in const [0.25, 0.5, 0.75, 1.0])
                      Expanded(
                        child: Padding(
                          padding: const EdgeInsets.only(right: 6),
                          child: OutlinedButton(onPressed: () => _preset(p), child: Text(p == 1 ? 'pmmax'.tr() : '${(p * 100).toInt()}%')),
                        ),
                      ),
                  ]),
                ]),
              ),
              if (quoting) const LinearProgressIndicator(),
              if (quoteError.isNotEmpty) Text(quoteError, style: const TextStyle(color: P2PTheme.danger)),
              if (quote != null)
                PMSection(
                  title: buy ? 'pmyouget'.tr() : 'pmyoureceive'.tr(),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text(buy ? '${pmQty(quote!.quantity)} ${a.ticker}' : pmMoney(quote!.netAmount, quote!.fundingAsset),
                        style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w700, color: P2PTheme.brandDark)),
                    const SizedBox(height: 6),
                    PMRow('pmreferenceprice'.tr(), '${pmMoney(quote!.price)} ${'pmpertoken'.tr()} · ${quote!.priceSource}'),
                    PMRow('pmtrovofee'.tr(), '${pmMoney(quote!.fee, quote!.fundingAsset)} (${quote!.feePercent}%)'),
                    const SizedBox(height: 6),
                    _pathNote(quote!),
                  ]),
                ),
            ]),
      bottomNavigationBar: a == null
          ? null
          : SafeArea(
              child: Padding(
                padding: const EdgeInsets.all(P2PTheme.space4),
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark, padding: const EdgeInsets.symmetric(vertical: 14)),
                  onPressed: busy || !a.open ? null : _review,
                  child: Text('pmrevieworder'.tr(), style: const TextStyle(color: Colors.white)),
                ),
              ),
            ),
    );
  }
}

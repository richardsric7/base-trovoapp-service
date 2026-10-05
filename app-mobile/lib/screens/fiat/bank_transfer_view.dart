import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/network/fiat_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

// BankTransferView moves Naira between a Nigerian bank account and the
// wallet's cNGN through Stablerail: Deposit (pay into a virtual account,
// cNGN arrives in the wallet), Withdraw (cNGN leaves the wallet, Naira is
// paid to a bank account) and the withdrawals' history. It is opened from
// the cNGN Deposit/Withdraw sheet (WrappedAsset); viewData['bankTab'] picks
// the tab (0 deposit, 1 withdraw, 2 history) and viewData['walletAddress']
// the wallet. Users are onboarded by the backend when KYC level 1 (BVN)
// completes; until then the screen points them to KYC.
class BankTransferView extends StatefulWidget {
  const BankTransferView({Key? key}) : super(key: key);

  @override
  State<BankTransferView> createState() => _BankTransferViewState();
}

class _BankTransferViewState extends State<BankTransferView> with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late FiatApi api;
  late TabController tabs;
  bool loading = true;
  Map profile = {};
  List<Map> banks = [];
  List<Map> withdrawals = [];
  bool busy = false;

  final depositAmount = TextEditingController();
  final withdrawAmount = TextEditingController();
  final accountNumber = TextEditingController();
  String? bankCode;
  Map? virtualAccount;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    api = FiatApi(appState, walletAddress: appState.viewData?['walletAddress'] as String?);
    final tab = (appState.viewData?['bankTab'] as int?) ?? 0;
    tabs = TabController(length: 3, vsync: this, initialIndex: tab.clamp(0, 2));
    _load();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final p = await api.profile();
    final enabled = p['statusCode'] == 200 && p['data']?['enabled'] == true;
    final onboarded = enabled && p['data']?['onboarded'] == true;
    final results = onboarded ? await Future.wait([api.banks(), api.withdrawals()]) : [<Map>[], <Map>[]];
    if (!mounted) return;
    setState(() {
      profile = p['statusCode'] == 200 ? Map.from(p['data'] ?? {}) : {};
      banks = results[0];
      withdrawals = results[1];
      loading = false;
    });
  }

  String _message(Map response, String fallback) => '${response['data']?['message'] ?? fallback}';

  Future<void> _deposit() async {
    final amount = depositAmount.text.trim();
    if (double.tryParse(amount) == null || double.parse(amount) <= 0) {
      popup(context, title: 'bankinvalidamount'.tr(), message: 'bankenteramount'.tr());
      return;
    }
    setState(() => busy = true);
    showLoader(context);
    final r = await api.deposit(amount);
    hideLoader(context);
    setState(() => busy = false);
    if (r['statusCode'] == 200 && r['data']?['data']?['virtualAccount'] != null) {
      setState(() => virtualAccount = Map.from(r['data']['data']['virtualAccount']));
    } else {
      popup(context, title: 'bankdepositfailed'.tr(), message: _message(r, 'p2ppleasetryagain'.tr()));
    }
  }

  // _withdraw builds the transfer, shows what it does, confirms with the
  // device's biometrics (or a dialog) and submits it signed - the same
  // build -> sign -> submit flow as payments.
  Future<void> _withdraw() async {
    final amount = withdrawAmount.text.trim();
    final account = accountNumber.text.trim();
    if (bankCode == null || account.length != 10 || double.tryParse(amount) == null) {
      popup(context, title: 'bankmissingdetails'.tr(), message: 'bankmissingdetailsmessage'.tr());
      return;
    }
    setState(() => busy = true);
    showLoader(context);
    try {
      final built = await api.withdrawBuild(amount: amount, accountNumber: account, bankCode: bankCode!);
      hideLoader(context);
      final tx = built['data']?['transaction'];
      if (built['statusCode'] != 200 || tx == null || tx == '') {
        setState(() => busy = false);
        popup(context, title: 'bankwithdrawalfailed'.tr(), message: _message(built, 'p2ppleasetryagain'.tr()));
        return;
      }
      final messages = ((built['data']['messages'] ?? []) as List).map((m) => '$m').toList();
      if (!await _confirm(messages)) {
        setState(() => busy = false);
        return;
      }
      showLoader(context);
      final signature = TrovoWalletSDK().signBase64Txn(appState.secretKeys[0], tx, '');
      final sent = await api.withdrawSubmit(
        amount: amount,
        accountNumber: account,
        bankCode: bankCode!,
        transaction: tx,
        transactionSignature: signature,
      );
      hideLoader(context);
      setState(() => busy = false);
      if (sent['statusCode'] == 200) {
        withdrawAmount.clear();
        await _load();
        tabs.animateTo(2);
        popup(context, title: 'bankwithdrawalsent'.tr(), message: 'bankwithdrawalsentmessage'.tr(), bodyColor: P2PTheme.success);
      } else {
        popup(context, title: 'bankwithdrawalfailed'.tr(), message: _message(sent, 'p2ppleasetryagain'.tr()));
      }
    } catch (e) {
      hideLoader(context);
      setState(() => busy = false);
      popup(context, title: 'error'.tr(), message: e.toString());
    }
  }

  Future<bool> _confirm(List<String> messages) async {
    final ok = await showDialog<bool>(
          context: context,
          builder: (context) => AlertDialog(
            title: Text('bankconfirmwithdrawal'.tr()),
            content: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
              for (final m in messages) Padding(padding: const EdgeInsets.only(bottom: P2PTheme.space2), child: Text(m)),
            ]),
            actions: [
              TextButton(onPressed: () => Navigator.of(context).pop(false), child: Text('cancel'.tr())),
              TextButton(onPressed: () => Navigator.of(context).pop(true), child: Text('p2pconfirm'.tr())),
            ],
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

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('banktitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
        bottom: loading || profile['onboarded'] != true
            ? null
            : TabBar(
                controller: tabs,
                labelColor: P2PTheme.brandDark,
                unselectedLabelColor: Colors.black45,
                indicatorColor: P2PTheme.brandDark,
                tabs: [Tab(text: 'deposit'.tr()), Tab(text: 'withdraw'.tr()), Tab(text: 'bankhistory'.tr())],
              ),
      ),
      body: loading
          ? const Center(child: CircularProgressIndicator())
          : profile['enabled'] != true
              ? P2PEmptyState(icon: Icons.account_balance_outlined, message: 'bankunavailable'.tr())
              : profile['onboarded'] != true
                  ? P2PEmptyState(
                      icon: Icons.verified_user_outlined,
                      message: profile['onboardingStatus'] == 'failed'
                          ? 'bankonboardingfailed'.tr()
                          : (profile['onboardingStatus'] ?? '') != ''
                              ? 'bankonboardinginprogress'.tr()
                              : 'bankcompletekyc'.tr(),
                      ctaLabel: (profile['onboardingStatus'] ?? '') == '' ? 'bankgotokyc'.tr() : null,
                      onCta: () => appState.currentAction = PageAction(state: PageState.addPage, page: KycScreenViewPageConfig),
                    )
                  : TabBarView(controller: tabs, children: [_depositTab(), _withdrawTab(), _historyTab()]),
    );
  }

  Widget _depositTab() => ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
        Text('bankdepositintro'.tr()),
        const SizedBox(height: P2PTheme.space4),
        _field('bankamountngn'.tr(), depositAmount, number: true),
        const SizedBox(height: P2PTheme.space4),
        _button('bankgetaccount'.tr(), _deposit),
        if (virtualAccount != null) ...[
          const SizedBox(height: P2PTheme.space6),
          P2PListCard(
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text('bankpayinto'.tr(), style: const TextStyle(fontWeight: FontWeight.w700)),
              const SizedBox(height: P2PTheme.space2),
              _copyRow('bankaccountnumber'.tr(), '${virtualAccount!['accountNumber'] ?? ''}'),
              _copyRow('bankname'.tr(), '${virtualAccount!['bankName'] ?? ''}'),
              _copyRow('bankaccountname'.tr(), '${virtualAccount!['accountName'] ?? ''}'),
              _copyRow('bankamounttopay'.tr(), 'NGN ${virtualAccount!['amount'] ?? ''}'),
              const SizedBox(height: P2PTheme.space2),
              Text('bankpaywithin'.tr(), style: const TextStyle(color: Colors.black54, fontSize: 12)),
            ]),
          ),
        ],
      ]);

  Widget _withdrawTab() {
    // the accounts paid before, most recent first, to pick again
    final recent = <String, Map>{};
    for (final w in withdrawals) {
      final key = '${w['bankCode']}:${w['accountNumber']}';
      if ((w['accountNumber'] ?? '') != '' && !recent.containsKey(key)) recent[key] = w;
    }
    return ListView(padding: const EdgeInsets.all(P2PTheme.space4), children: [
      Text('bankwithdrawintro'.tr(args: ['${profile['minimumWithdrawal'] ?? ''}'])),
      if (recent.isNotEmpty) ...[
        const SizedBox(height: P2PTheme.space4),
        Text('bankrecentaccounts'.tr(), style: const TextStyle(fontWeight: FontWeight.w600)),
        const SizedBox(height: P2PTheme.space2),
        Wrap(spacing: P2PTheme.space2, runSpacing: P2PTheme.space2, children: [
          for (final w in recent.values.take(5))
            ActionChip(
              label: Text('${w['bankName'] ?? w['bankCode']} ${w['accountNumber']}'),
              onPressed: () => setState(() {
                bankCode = banks.any((b) => b['bank_code'] == w['bankCode']) ? w['bankCode'] : null;
                accountNumber.text = '${w['accountNumber']}';
              }),
            ),
        ]),
      ],
      const SizedBox(height: P2PTheme.space4),
      DropdownButtonFormField<String>(
        initialValue: bankCode,
        isExpanded: true,
        decoration: InputDecoration(labelText: 'bankname'.tr(), border: const OutlineInputBorder(), filled: true, fillColor: Colors.white),
        items: [for (final b in banks) DropdownMenuItem(value: '${b['bank_code']}', child: Text('${b['bank_name']}'))],
        onChanged: (v) => setState(() => bankCode = v),
      ),
      const SizedBox(height: P2PTheme.space3),
      _field('bankaccountnumber'.tr(), accountNumber, number: true, maxLength: 10),
      const SizedBox(height: P2PTheme.space3),
      _field('bankamountngn'.tr(), withdrawAmount, number: true),
      const SizedBox(height: P2PTheme.space4),
      _button('bankwithdraw'.tr(), _withdraw),
    ]);
  }

  Widget _historyTab() => withdrawals.isEmpty
      ? P2PEmptyState(icon: Icons.receipt_long_outlined, message: 'banknowithdrawals'.tr())
      : RefreshIndicator(
          onRefresh: _load,
          child: ListView.builder(
            padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4, vertical: P2PTheme.space2),
            itemCount: withdrawals.length,
            itemBuilder: (context, i) {
              final w = withdrawals[i];
              final created = DateTime.tryParse('${w['createdAt']}')?.toLocal();
              return P2PListCard(
                child: Row(children: [
                  Expanded(
                    child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                      Text('NGN ${w['amount']}', style: const TextStyle(fontWeight: FontWeight.w700)),
                      const SizedBox(height: P2PTheme.space1),
                      Text('${w['bankName'] ?? w['bankCode']} ${w['accountNumber']}', style: const TextStyle(color: Colors.black54, fontSize: 12)),
                      if (created != null)
                        Text(DateFormat.yMMMd().add_jm().format(created), style: const TextStyle(color: Colors.black45, fontSize: 11)),
                    ]),
                  ),
                  _statusPill('${w['status']}'),
                ]),
              );
            },
          ),
        );

  // bank withdrawal statuses: ours while the cNGN moves, then Stablerail's
  Widget _statusPill(String status) {
    final s = status.toLowerCase();
    final done = s == 'completed' || s == 'complete' || s == 'success' || s == 'successful';
    final failed = s.contains('fail') || s == 'cancelled' || s == 'canceled' || s == 'rejected' || s == 'reversed';
    final color = done ? P2PTheme.success : (failed ? P2PTheme.danger : P2PTheme.warning);
    final label = done ? 'bankstatuspaid'.tr() : (failed ? 'bankstatusfailed'.tr() : 'bankstatusprocessing'.tr());
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space2, vertical: P2PTheme.space1),
      decoration: BoxDecoration(color: color.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(999)),
      child: Text(label, style: TextStyle(color: color, fontSize: 12, fontWeight: FontWeight.w600)),
    );
  }

  Widget _field(String label, TextEditingController c, {bool number = false, int? maxLength}) => TextField(
        controller: c,
        maxLength: maxLength,
        keyboardType: number ? const TextInputType.numberWithOptions(decimal: true) : TextInputType.text,
        decoration: InputDecoration(labelText: label, border: const OutlineInputBorder(), filled: true, fillColor: Colors.white, counterText: ''),
      );

  Widget _button(String label, VoidCallback onPressed) => ElevatedButton(
        style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark, padding: const EdgeInsets.symmetric(vertical: P2PTheme.space3)),
        onPressed: busy ? null : onPressed,
        child: Text(label, style: const TextStyle(color: Colors.white)),
      );

  Widget _copyRow(String label, String value) => Padding(
        padding: const EdgeInsets.symmetric(vertical: P2PTheme.space1),
        child: Row(children: [
          Expanded(child: Text(label, style: const TextStyle(color: Colors.black54))),
          Text(value, style: const TextStyle(fontWeight: FontWeight.w600)),
          IconButton(
            icon: const Icon(Icons.copy, size: 16),
            onPressed: () => Clipboard.setData(ClipboardData(text: value.replaceFirst('NGN ', ''))),
          ),
        ]),
      );
}

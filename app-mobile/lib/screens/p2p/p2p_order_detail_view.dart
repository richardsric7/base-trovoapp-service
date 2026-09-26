import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/p2p_order.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import 'package:url_launcher/url_launcher.dart';

// P2POrderDetailView is the hub screen (Plan Sections 95.4-95.9,
// consolidated): the order timeline is the feature's emotional centerpiece
// (Plan Section 97.2), with the one relevant action for the order's current
// state surfaced directly below it, rather than splitting pending-approval/
// merchant-approval/awaiting-payment/etc. into separate screens.
class P2POrderDetailView extends StatefulWidget {
  const P2POrderDetailView({Key? key}) : super(key: key);

  @override
  State<P2POrderDetailView> createState() => _P2POrderDetailViewState();
}

class _P2POrderDetailViewState extends State<P2POrderDetailView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  bool loaded = false;
  String? orderId;
  P2POrder? order;
  P2PDispute? dispute;
  Map<String, dynamic>? customerPerf;
  bool loading = true;
  bool actionInFlight = false;

  static const _steps = [
    P2POrder.statusAwaitingApproval,
    P2POrder.statusAwaitingEscrowDeposit,
    P2POrder.statusAwaitingPayment,
    P2POrder.statusAwaitingPaymentConfirmation,
    P2POrder.statusCompleted,
  ];

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    if (!loaded) {
      loaded = true;
      orderId = appState.viewData?[P2POrderDetailViewPageConfig.key]?['orderId'];
      _load();
    }
  }

  Future<void> _load() async {
    if (orderId == null) return;
    setState(() => loading = true);
    final o = await api.getOrder(orderId!);
    P2PDispute? d;
    if (o?.isDisputed == true) {
      d = await api.getOpenDisputeForOrder(orderId!);
    }
    setState(() {
      order = o;
      dispute = d;
      loading = false;
    });
    if (o != null && o.isMerchant(myUsername)) {
      final perf = await api.getCustomerPerformance(o.id!);
      if (mounted) setState(() => customerPerf = perf);
    }
  }

  String get myUsername => appState.userInfo?.username ?? '';

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);

    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('p2pordertitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: loading
          ? const Center(child: CircularProgressIndicator())
          : order == null
              ? P2PEmptyState(icon: Icons.error_outline, message: 'p2pordernotfound'.tr())
              : RefreshIndicator(
                  onRefresh: _load,
                  child: SingleChildScrollView(
                    padding: const EdgeInsets.all(P2PTheme.space4),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _header(order!),
                        const SizedBox(height: P2PTheme.space4),
                        if (order!.isDisputed == true) _disputeBanner(),
                        _timeline(order!),
                        const SizedBox(height: P2PTheme.space4),
                        _actionPanel(order!),
                      ],
                    ),
                  ),
                ),
    );
  }

  Widget _header(P2POrder o) {
    return P2PListCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '${o.specifiedAssetAmount} ${o.asset}',
                style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 18),
              ),
              P2PStatusPill(status: o.orderStatus ?? '', disputed: o.isDisputed ?? false),
            ],
          ),
          const SizedBox(height: P2PTheme.space1),
          Text('${o.paymentAmount} ${o.currency}', style: const TextStyle(color: Colors.black54)),
          const SizedBox(height: P2PTheme.space1),
          Text(
            o.isCustomer(myUsername)
                ? 'p2pmerchantlabel'.tr(args: [o.merchantUsername ?? ''])
                : 'p2pcustomerlabel'.tr(args: [o.customerUsername ?? '']),
            style: const TextStyle(color: Colors.black38, fontSize: 12),
          ),
        ],
      ),
    );
  }

  Widget _disputeBanner() {
    return Container(
      margin: const EdgeInsets.only(bottom: P2PTheme.space4),
      padding: const EdgeInsets.all(P2PTheme.space3),
      decoration: BoxDecoration(
        color: P2PTheme.danger.withOpacity(0.08),
        borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
        border: Border.all(color: P2PTheme.danger.withOpacity(0.3)),
      ),
      child: Row(
        children: [
          const Icon(Icons.error_outline, color: P2PTheme.danger),
          const SizedBox(width: P2PTheme.space2),
          Expanded(
            child: Text(
              'p2pdisputeopenbanner'.tr(),
              style: const TextStyle(color: P2PTheme.danger, fontWeight: FontWeight.w600),
            ),
          ),
        ],
      ),
    );
  }

  // The order timeline is the single highest-leverage design element in
  // the feature (Plan Section 97.2) - a vertical stepper over the
  // lifecycle states, so the customer/merchant always sees what's
  // happening now and what happens next.
  Widget _timeline(P2POrder o) {
    final currentIndex = _steps.indexOf(o.orderStatus ?? '');
    final terminalNonSuccess = ['REJECTED', 'CANCELLED', 'EXPIRED'].contains(o.orderStatus);
    return P2PListCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (terminalNonSuccess)
            Text(
              P2PTheme.statusLabel(o.orderStatus!),
              style: const TextStyle(fontWeight: FontWeight.w700, color: Colors.black54),
            )
          else
            for (var i = 0; i < _steps.length; i++)
              _timelineStep(
                label: P2PTheme.statusLabel(_steps[i]),
                done: currentIndex >= 0 && i < currentIndex,
                current: i == currentIndex,
                isLast: i == _steps.length - 1,
              ),
        ],
      ),
    );
  }

  Widget _timelineStep({
    required String label,
    required bool done,
    required bool current,
    required bool isLast,
  }) {
    final color = done
        ? P2PTheme.success
        : current
            ? P2PTheme.brandDark
            : Colors.black26;
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Column(
          children: [
            Icon(
              done ? Icons.check_circle : (current ? Icons.radio_button_checked : Icons.radio_button_unchecked),
              color: color,
              size: 20,
            ),
            if (!isLast) Container(width: 2, height: 28, color: color.withOpacity(0.3)),
          ],
        ),
        const SizedBox(width: P2PTheme.space3),
        Padding(
          padding: const EdgeInsets.only(bottom: P2PTheme.space2),
          child: Text(
            label,
            style: TextStyle(
              color: current ? P2PTheme.brandDark : Colors.black87,
              fontWeight: current ? FontWeight.w700 : FontWeight.w400,
            ),
          ),
        ),
      ],
    );
  }

  Widget _actionPanel(P2POrder o) {
    final status = o.orderStatus;
    final children = <Widget>[];

    if (status == P2POrder.statusAwaitingApproval) {
      if (o.isMerchant(myUsername)) {
        if (customerPerf != null && (customerPerf!['completedTrades'] ?? 0) > 0) {
          children.add(Padding(
            padding: const EdgeInsets.only(bottom: P2PTheme.space2),
            child: Text(
              'p2pcustomerperfsummary'.tr(args: ['${customerPerf!['completedTrades']}', '${customerPerf!['completionRate']}']),
              style: const TextStyle(color: Colors.black38, fontSize: 12),
            ),
          ));
        }
        children.add(_buttonRow([
          _primaryButton('p2pacceptorder'.tr(), () => _simpleAction(() => api.acceptOrder(o.id!, address: appState.primaryWallet.address!))),
          _secondaryButton('p2preject'.tr(), () => _simpleAction(() => api.rejectOrder(o.id!, address: appState.primaryWallet.address!))),
        ]));
      } else {
        children.add(_infoText('p2pwaitingformerchantaccept'.tr()));
        children.add(_secondaryButton('p2pcancelorder'.tr(), () => _simpleAction(() => api.cancelOrder(o.id!, address: appState.primaryWallet.address!))));
      }
    } else if (status == P2POrder.statusAwaitingEscrowDeposit) {
      if (o.isAssetDepositor(myUsername)) {
        children.add(_infoText('p2pdepositintoescrow'.tr(args: ['${o.sellerEscrowAssetAmount}', '${o.asset}'])));
        children.add(_buttonRow([
          _primaryButton('p2pdepositfromwallet'.tr(), _depositFromOwnWallet),
          if (o.escrowDepositShortlink != null && o.escrowDepositShortlink!.isNotEmpty)
            _secondaryButton('p2psharedepositlink'.tr(), () {
              appState.viewData ??= {};
              appState.viewData![P2PEscrowShareViewPageConfig.key] = {'order': o};
              appState.currentAction = PageAction(state: PageState.addPage, page: P2PEscrowShareViewPageConfig);
            })
          else
            _secondaryButton(
              'p2pgeneratedepositlink'.tr(),
              () => _simpleAction(() => api.regenerateEscrowShortlink(o.id!, address: appState.primaryWallet.address!)),
            ),
        ]));
      } else {
        children.add(_infoText('p2pwaitingforescrow'.tr()));
      }
      if (o.isMerchant(myUsername)) {
        children.add(const SizedBox(height: P2PTheme.space2));
        children.add(_secondaryButton(
          'p2pcancelorder'.tr(),
          () => _simpleAction(() => api.merchantCancelOrder(o.id!, address: appState.primaryWallet.address!)),
        ));
      }
    } else if (status == P2POrder.statusAwaitingPayment) {
      if (o.isFiatPayer(myUsername)) {
        children.add(_paymentMethodCard(o));
        children.add(_primaryButton('p2phavesentpayment'.tr(), () => _simpleAction(() => api.markPaymentSent(o.id!, address: appState.primaryWallet.address!))));
      } else {
        children.add(_infoText('p2pescrowconfirmedwaiting'.tr()));
      }
      children.add(const SizedBox(height: P2PTheme.space2));
      children.add(_disputeButton(o));
    } else if (status == P2POrder.statusAwaitingPaymentConfirmation) {
      if (o.isFiatRecipient(myUsername)) {
        children.add(_infoText('p2pbuyermarkedsent'.tr()));
        children.add(_primaryButton('p2pconfirmpaymentreceived'.tr(), () => _simpleAction(() => api.confirmPaymentReceived(o.id!, address: appState.primaryWallet.address!))));
      } else {
        children.add(_infoText('p2pwaitingsellerconfirm'.tr()));
      }
      children.add(const SizedBox(height: P2PTheme.space2));
      children.add(_disputeButton(o));
    } else if (status == P2POrder.statusCompleted) {
      children.add(_infoText('p2porderiscomplete'.tr()));
      if (o.assetReleaseTransactionHash != null && o.assetReleaseTransactionHash!.isNotEmpty) {
        children.add(_txLink('p2pviewassetrelease'.tr(), o.assetReleaseTransactionHash!));
      }
      if (o.escrowDepositTransactionHash != null && o.escrowDepositTransactionHash!.isNotEmpty) {
        children.add(_txLink('p2pviewescrowdeposit'.tr(), o.escrowDepositTransactionHash!));
      }
    }

    if (o.isDisputed == true && dispute != null) {
      children.add(const SizedBox(height: P2PTheme.space3));
      children.add(_disputeResolutionActions(o));
    }

    if (children.isEmpty) return const SizedBox.shrink();
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: children);
  }

  Widget _disputeButton(P2POrder o) {
    if (o.isDisputed == true) return const SizedBox.shrink();
    return TextButton.icon(
      onPressed: () {
        appState.viewData ??= {};
        appState.viewData![P2PDisputeViewPageConfig.key] = {'orderId': o.id};
        appState.currentAction = PageAction(state: PageState.addPage, page: P2PDisputeViewPageConfig);
      },
      icon: const Icon(Icons.flag_outlined, color: P2PTheme.danger),
      label: Text('p2praiseadispute'.tr(), style: const TextStyle(color: P2PTheme.danger)),
    );
  }

  Widget _disputeResolutionActions(P2POrder o) {
    final iAmMerchant = o.isMerchant(myUsername);
    final iAmCustomer = o.isCustomer(myUsername);
    return P2PListCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('p2presolvethisdispute'.tr(), style: const TextStyle(fontWeight: FontWeight.w700)),
          const SizedBox(height: P2PTheme.space2),
          if (iAmMerchant)
            _primaryButton(
              'p2preceivedpayment'.tr(),
              () => _simpleAction(() => api.merchantConfirmsPayment(dispute!.id!)),
            ),
          if (iAmCustomer)
            _secondaryButton(
              'p2pnotpaidyet'.tr(),
              () => _simpleAction(() => api.buyerConfirmsNotPaid(dispute!.id!)),
            ),
        ],
      ),
    );
  }

  Widget _paymentMethodCard(P2POrder o) {
    final pm = o.paymentMethodSnapshot;
    return P2PListCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('p2ppaymerchantusing'.tr(), style: const TextStyle(fontWeight: FontWeight.w700)),
          const SizedBox(height: P2PTheme.space2),
          Text('${pm?.paymentChannel ?? ''} - ${pm?.provider ?? ''}'),
          const SizedBox(height: P2PTheme.space1),
          Text(pm?.account ?? '', style: const TextStyle(fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }

  Widget _txLink(String label, String txHash) {
    return TextButton(
      onPressed: () async {
        final uri = Uri.parse('${getExplorerBaseUrl(appState.walletMode)}$txHash');
        if (!await launchUrl(uri, mode: LaunchMode.externalApplication)) {
          popup(context, title: 'p2ptransactionhash'.tr(), message: txHash);
        }
      },
      child: Text(label),
    );
  }

  Widget _infoText(String text) => Padding(
        padding: const EdgeInsets.only(bottom: P2PTheme.space2),
        child: Text(text, style: const TextStyle(color: Colors.black54)),
      );

  Widget _buttonRow(List<Widget> buttons) => Row(
        children: [for (final b in buttons) Expanded(child: Padding(padding: const EdgeInsets.only(right: P2PTheme.space2), child: b))],
      );

  Widget _primaryButton(String label, VoidCallback onPressed) => ElevatedButton(
        style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark, padding: const EdgeInsets.symmetric(vertical: P2PTheme.space3)),
        onPressed: actionInFlight ? null : onPressed,
        child: Text(label, style: const TextStyle(color: Colors.white)),
      );

  Widget _secondaryButton(String label, VoidCallback onPressed) => OutlinedButton(
        style: OutlinedButton.styleFrom(padding: const EdgeInsets.symmetric(vertical: P2PTheme.space3)),
        onPressed: actionInFlight ? null : onPressed,
        child: Text(label),
      );

  Future<void> _simpleAction(Future<Map> Function() action) async {
    setState(() => actionInFlight = true);
    showLoader(context);
    try {
      final response = await action();
      hideLoader(context);
      setState(() => actionInFlight = false);
      if (response['statusCode'] != null && response['statusCode'] < 300) {
        _load();
      } else {
        popup(context, title: 'p2pcouldnotcompleteaction'.tr(), message: response['data']?['message'] ?? 'p2ppleasetryagain'.tr());
      }
    } catch (e) {
      hideLoader(context);
      setState(() => actionInFlight = false);
      popup(context, title: 'error'.tr(), message: e.toString());
    }
  }

  // Escrow deposit from the depositor's own wallet reuses the exact same
  // two-phase build -> sign -> commit contract as the generic Send flow
  // (Plan Section 28/29): the first call returns an unsigned transaction,
  // which is signed locally after a biometric confirmation gate, then
  // committed with a second call.
  Future<void> _depositFromOwnWallet() async {
    final authenticator = Authenticator();
    bool confirmed = false;
    try {
      confirmed = await authenticator.authenticateMe();
    } catch (_) {
      confirmed = false;
    }
    if (!confirmed) {
      confirmed = await _confirmWithoutBiometrics();
    }
    if (!confirmed) return;

    setState(() => actionInFlight = true);
    showLoader(context);
    try {
      final address = appState.primaryWallet.address!;
      final built = await api.escrowDepositBuild(order!.id!, address: address);
      if (built['statusCode'] != 202 || built['data']['transaction'] == null || built['data']['transaction'] == '') {
        hideLoader(context);
        setState(() => actionInFlight = false);
        popup(context, title: 'p2pcouldnotbuilddeposit'.tr(), message: built['data']?['message'] ?? 'p2ppleasetryagain'.tr());
        return;
      }
      final signature = TrovoWalletSDK().signBase64Txn(
        appState.secretKeys[0],
        built['data']['transaction'],
        built['data']['networkPassPhrase'] ?? '',
      );
      final committed = await api.escrowDepositCommit(
        order!.id!,
        address: address,
        transaction: built['data']['transaction'],
        transactionSignature: signature,
      );
      hideLoader(context);
      setState(() => actionInFlight = false);
      if (committed['statusCode'] == 202) {
        _load();
      } else {
        popup(context, title: 'p2pdepositfailed'.tr(), message: committed['data']?['message'] ?? 'p2ppleasetryagain'.tr());
      }
    } catch (e) {
      hideLoader(context);
      setState(() => actionInFlight = false);
      popup(context, title: 'error'.tr(), message: e.toString());
    }
  }

  Future<bool> _confirmWithoutBiometrics() async {
    return await showDialog<bool>(
          context: context,
          builder: (context) => AlertDialog(
            title: Text('p2pconfirmdeposit'.tr()),
            content: Text('p2pbiometricunavailableconfirm'.tr()),
            actions: [
              TextButton(onPressed: () => Navigator.of(context).pop(false), child: Text('cancel'.tr())),
              TextButton(onPressed: () => Navigator.of(context).pop(true), child: Text('p2pconfirm'.tr())),
            ],
          ),
        ) ??
        false;
  }
}

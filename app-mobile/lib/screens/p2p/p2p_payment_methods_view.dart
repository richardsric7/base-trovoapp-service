import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/p2p_merchant_gate.dart';

// P2PPaymentMethodsView manages a merchant's own saved fiat settlement
// channels: create, edit details (propagates to any live offer using it),
// and activate/deactivate. Never a delete - a payment method still
// referenced by a live offer cannot be deactivated either, the backend
// enforces that and this screen just surfaces the resulting error.
class P2PPaymentMethodsView extends StatefulWidget {
  const P2PPaymentMethodsView({Key? key}) : super(key: key);

  @override
  State<P2PPaymentMethodsView> createState() => _P2PPaymentMethodsViewState();
}

class _P2PPaymentMethodsViewState extends State<P2PPaymentMethodsView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;

  bool loading = true;
  List<MerchantPaymentMethod> methods = [];
  String error = '';

  bool showNew = false;
  final newChannel = TextEditingController();
  final newProvider = TextEditingController();
  final newAccount = TextEditingController();
  bool creating = false;

  String editingId = '';
  final editChannel = TextEditingController();
  final editProvider = TextEditingController();
  final editAccount = TextEditingController();

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final result = await api.listMyPaymentMethods();
    if (!mounted) return;
    setState(() {
      methods = result;
      loading = false;
    });
  }

  Future<void> _createNew() async {
    setState(() {
      creating = true;
      error = '';
    });
    final response = await api.createPaymentMethod({
      'paymentChannel': newChannel.text.trim(),
      'provider': newProvider.text.trim(),
      'account': newAccount.text.trim(),
    });
    if (!mounted) return;
    if (response['statusCode'] == 201) {
      setState(() {
        showNew = false;
        creating = false;
        newChannel.clear();
        newProvider.clear();
        newAccount.clear();
      });
      _load();
    } else {
      setState(() {
        creating = false;
        error = response['data']?['message'] ?? 'p2pcouldnotsavepaymentmethod'.tr();
      });
    }
  }

  void _startEdit(MerchantPaymentMethod pm) {
    setState(() {
      editingId = pm.id!;
      editChannel.text = pm.paymentChannel ?? '';
      editProvider.text = pm.provider ?? '';
      editAccount.text = pm.account ?? '';
    });
  }

  Future<void> _saveEdit() async {
    setState(() => error = '');
    final response = await api.updatePaymentMethod(editingId, {
      'paymentChannel': editChannel.text.trim(),
      'provider': editProvider.text.trim(),
      'account': editAccount.text.trim(),
    });
    if (!mounted) return;
    if (response['statusCode'] == 200) {
      setState(() => editingId = '');
      _load();
    } else {
      setState(() => error = response['data']?['message'] ?? 'p2pcouldnotsavepaymentmethod'.tr());
    }
  }

  Future<void> _toggleActive(MerchantPaymentMethod pm) async {
    setState(() => error = '');
    final response = await api.setPaymentMethodActive(pm.id!, !pm.isActive);
    if (!mounted) return;
    if (response['statusCode'] == 200) {
      _load();
    } else {
      setState(() => error = response['data']?['message'] ?? 'p2ppaymentmethodinuse'.tr());
    }
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
        title: Text('p2ppaymentmethodstitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: P2PMerchantGate(child: RefreshIndicator(
        onRefresh: _load,
        child: ListView(
          padding: const EdgeInsets.all(P2PTheme.space4),
          children: [
            if (error.isNotEmpty)
              Padding(
                padding: const EdgeInsets.only(bottom: P2PTheme.space3),
                child: Text(error, style: const TextStyle(color: P2PTheme.danger)),
              ),
            if (!showNew)
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark),
                  onPressed: () => setState(() => showNew = true),
                  child: Text('p2paddnewpaymentmethod'.tr(), style: const TextStyle(color: Colors.white)),
                ),
              )
            else
              _newForm(),
            const SizedBox(height: P2PTheme.space4),
            if (loading)
              const Center(child: CircularProgressIndicator())
            else if (methods.isEmpty)
              P2PEmptyState(icon: Icons.account_balance_outlined, message: 'p2pnopaymentmethodsyet'.tr())
            else
              ...methods.map(_methodCard),
          ],
        ),
      )),
    );
  }

  Widget _newForm() {
    return Container(
      padding: const EdgeInsets.all(P2PTheme.space3),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(P2PTheme.space2)),
      child: Column(
        children: [
          _input('p2ppaymentchannel'.tr(), newChannel),
          _input('p2pproviderlabel'.tr(), newProvider),
          _input('p2paccountdetails'.tr(), newAccount),
          Row(
            children: [
              Expanded(
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark),
                  onPressed: creating ? null : _createNew,
                  child: Text('p2psavepaymentmethod'.tr(), style: const TextStyle(color: Colors.white)),
                ),
              ),
              const SizedBox(width: P2PTheme.space2),
              TextButton(onPressed: () => setState(() => showNew = false), child: Text('cancel'.tr())),
            ],
          ),
        ],
      ),
    );
  }

  Widget _methodCard(MerchantPaymentMethod pm) {
    return P2PListCard(
      child: editingId == pm.id
          ? Column(
              children: [
                _input('p2ppaymentchannel'.tr(), editChannel),
                _input('p2pproviderlabel'.tr(), editProvider),
                _input('p2paccountdetails'.tr(), editAccount),
                Row(
                  children: [
                    Expanded(
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark),
                        onPressed: _saveEdit,
                        child: Text('p2psavechanges'.tr(), style: const TextStyle(color: Colors.white)),
                      ),
                    ),
                    const SizedBox(width: P2PTheme.space2),
                    TextButton(onPressed: () => setState(() => editingId = ''), child: Text('cancel'.tr())),
                  ],
                ),
              ],
            )
          : Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('${pm.paymentChannel} - ${pm.provider}', style: const TextStyle(fontWeight: FontWeight.w700)),
                const SizedBox(height: P2PTheme.space1),
                Text(pm.account ?? '', style: const TextStyle(color: Colors.black54)),
                const SizedBox(height: P2PTheme.space1),
                Text(
                  pm.isActive ? 'p2pactive'.tr() : 'p2pdisabled'.tr(),
                  style: TextStyle(fontSize: 12, color: pm.isActive ? P2PTheme.success : P2PTheme.danger),
                ),
                const SizedBox(height: P2PTheme.space2),
                Row(
                  children: [
                    TextButton(onPressed: () => _startEdit(pm), child: Text('edit'.tr())),
                    TextButton(
                      onPressed: () => _toggleActive(pm),
                      child: Text(pm.isActive ? 'p2pdisable'.tr() : 'p2penable'.tr()),
                    ),
                  ],
                ),
              ],
            ),
    );
  }

  Widget _input(String label, TextEditingController controller) {
    return Padding(
      padding: const EdgeInsets.only(bottom: P2PTheme.space2),
      child: TextFormField(
        controller: controller,
        decoration: InputDecoration(
          labelText: label,
          filled: true,
          fillColor: const Color(0xFFF9FAFB),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(P2PTheme.space2),
            borderSide: BorderSide.none,
          ),
        ),
      ),
    );
  }
}

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';

const String _addNew = '__add_new__';

// P2PPaymentMethodSelect is the SELL-offer payment-method dropdown: "Add
// new payment method" is always the first item, and picking it opens an
// inline quick-create form - on save, the new payment method is selected
// automatically and the merchant continues filling out the rest of the
// offer without leaving the screen.
class P2PPaymentMethodSelect extends StatefulWidget {
  final P2PApi api;
  final String? value;
  final ValueChanged<String?> onChanged;
  const P2PPaymentMethodSelect({super.key, required this.api, required this.value, required this.onChanged});

  @override
  State<P2PPaymentMethodSelect> createState() => _P2PPaymentMethodSelectState();
}

class _P2PPaymentMethodSelectState extends State<P2PPaymentMethodSelect> {
  List<MerchantPaymentMethod> methods = [];
  bool loading = true;
  bool addingNew = false;
  bool saving = false;
  String error = '';

  final paymentChannel = TextEditingController();
  final provider = TextEditingController();
  final account = TextEditingController();

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final result = await widget.api.listMyPaymentMethods();
    if (!mounted) return;
    setState(() {
      methods = result;
      loading = false;
    });
  }

  Future<void> _saveNew() async {
    setState(() {
      saving = true;
      error = '';
    });
    try {
      final response = await widget.api.createPaymentMethod({
        'paymentChannel': paymentChannel.text.trim(),
        'provider': provider.text.trim(),
        'account': account.text.trim(),
      });
      if (response['statusCode'] == 201) {
        final pm = MerchantPaymentMethod().deserializeJson(response['data']);
        if (!mounted) return;
        setState(() {
          methods = [...methods, pm];
          addingNew = false;
          saving = false;
          paymentChannel.clear();
          provider.clear();
          account.clear();
        });
        widget.onChanged(pm.id);
      } else {
        setState(() {
          saving = false;
          error = response['data']?['message'] ?? 'p2pcouldnotsavepaymentmethod'.tr();
        });
      }
    } catch (e) {
      if (!mounted) return;
      setState(() {
        saving = false;
        error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final activeMethods = methods.where((m) => m.isActive).toList();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('p2ppaymentmethod'.tr(), style: const TextStyle(fontWeight: FontWeight.w600)),
        const SizedBox(height: P2PTheme.space1),
        if (loading)
          const Padding(padding: EdgeInsets.symmetric(vertical: 8), child: LinearProgressIndicator())
        else
          Container(
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(P2PTheme.space2),
            ),
            child: DropdownButtonFormField<String>(
              value: addingNew ? _addNew : widget.value,
              isExpanded: true,
              decoration: InputDecoration(
                hintText: 'p2pselectpaymentmethod'.tr(),
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(P2PTheme.space2),
                  borderSide: BorderSide.none,
                ),
              ),
              items: [
                DropdownMenuItem(value: _addNew, child: Text('p2paddnewpaymentmethod'.tr())),
                ...activeMethods.map((m) => DropdownMenuItem(
                      value: m.id,
                      child: Text(
                        '${m.paymentChannel} - ${m.provider} (${m.account})',
                        overflow: TextOverflow.ellipsis,
                      ),
                    )),
              ],
              onChanged: (v) {
                if (v == _addNew) {
                  setState(() => addingNew = true);
                  return;
                }
                setState(() => addingNew = false);
                widget.onChanged(v);
              },
            ),
          ),
        if (addingNew) ...[
          const SizedBox(height: P2PTheme.space3),
          Container(
            padding: const EdgeInsets.all(P2PTheme.space3),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(P2PTheme.space2),
            ),
            child: Column(
              children: [
                _field('p2ppaymentchannel'.tr(), paymentChannel),
                _field('p2pproviderlabel'.tr(), provider),
                _field('p2paccountdetails'.tr(), account),
                if (error.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(bottom: P2PTheme.space2),
                    child: Text(error, style: const TextStyle(color: P2PTheme.danger)),
                  ),
                Row(
                  children: [
                    Expanded(
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark),
                        onPressed: saving ? null : _saveNew,
                        child: Text('p2psavepaymentmethod'.tr(), style: const TextStyle(color: Colors.white)),
                      ),
                    ),
                    const SizedBox(width: P2PTheme.space2),
                    TextButton(
                      onPressed: () => setState(() => addingNew = false),
                      child: Text('cancel'.tr()),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ],
    );
  }

  Widget _field(String label, TextEditingController controller) {
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

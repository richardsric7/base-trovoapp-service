import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

// P2PCreateOfferView (Plan Section 95.12): merchant offer creation, a
// single-page form (asset/type/price/limits/payment method/country) -
// following this app's standard Form/GlobalKey pattern (Plan Section 91).
class P2PCreateOfferView extends StatefulWidget {
  const P2PCreateOfferView({Key? key}) : super(key: key);

  @override
  State<P2PCreateOfferView> createState() => _P2PCreateOfferViewState();
}

class _P2PCreateOfferViewState extends State<P2PCreateOfferView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  final formKey = GlobalKey<FormState>();

  String offerType = 'SELL';
  final asset = TextEditingController();
  final country = TextEditingController();
  final countryCode = TextEditingController();
  final currency = TextEditingController();
  final price = TextEditingController();
  final minOrderAmount = TextEditingController();
  final maxOrderAmount = TextEditingController();
  final availableLiquidity = TextEditingController();
  final paymentChannel = TextEditingController();
  final provider = TextEditingController();
  final account = TextEditingController();
  final remark = TextEditingController();
  bool submitting = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
  }

  String? _required(String? v) => (v == null || v.trim().isEmpty) ? 'Required' : null;

  Future<void> _submit() async {
    if (!formKey.currentState!.validate()) return;
    setState(() => submitting = true);
    showLoader(context);
    try {
      final response = await api.createOffer({
        'offerType': offerType,
        'asset': asset.text.trim().toUpperCase(),
        'paymentMethod': {
          'paymentChannel': paymentChannel.text.trim(),
          'provider': provider.text.trim(),
          'account': account.text.trim(),
        },
        'country': country.text.trim(),
        'countryCode': countryCode.text.trim().toUpperCase(),
        'currency': currency.text.trim().toUpperCase(),
        'priceType': 'FIXED',
        'price': price.text.trim(),
        'minOrderAmount': minOrderAmount.text.trim(),
        'maxOrderAmount': maxOrderAmount.text.trim(),
        'availableLiquidity': availableLiquidity.text.trim(),
        'remark': remark.text.trim(),
      });
      hideLoader(context);
      setState(() => submitting = false);
      if (response['statusCode'] == 201) {
        appState.currentAction = PageAction(
          state: PageState.replace,
          page: P2PMyOffersViewPageConfig,
        );
      } else {
        popup(
          context,
          title: 'Could not create offer',
          message: response['data']?['message'] ?? 'Please check your inputs and try again.',
        );
      }
    } catch (e) {
      hideLoader(context);
      setState(() => submitting = false);
      popup(context, title: 'Error', message: e.toString());
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
        title: Text('Create offer', style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: Form(
        key: formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(P2PTheme.space4),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: RadioListTile<String>(
                      value: 'SELL',
                      groupValue: offerType,
                      title: const Text('Sell'),
                      onChanged: (v) => setState(() => offerType = v!),
                    ),
                  ),
                  Expanded(
                    child: RadioListTile<String>(
                      value: 'BUY',
                      groupValue: offerType,
                      title: const Text('Buy'),
                      onChanged: (v) => setState(() => offerType = v!),
                    ),
                  ),
                ],
              ),
              _field('Asset code (e.g. USDC)', asset),
              _field('Price per unit', price, numeric: true),
              _field('Currency (e.g. NGN)', currency),
              _field('Min order amount', minOrderAmount, numeric: true),
              _field('Max order amount', maxOrderAmount, numeric: true),
              _field('Available liquidity', availableLiquidity, numeric: true),
              _field('Country', country),
              _field('Country code (ISO-3)', countryCode),
              const SizedBox(height: P2PTheme.space2),
              const Text('Payment method', style: TextStyle(fontWeight: FontWeight.w600)),
              _field('Payment channel (e.g. BANK_TRANSFER)', paymentChannel),
              _field('Provider (e.g. GTBANK)', provider),
              _field('Account details', account),
              _field('Remark (optional)', remark, required: false),
              const SizedBox(height: P2PTheme.space4),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(
                    backgroundColor: P2PTheme.brandDark,
                    padding: const EdgeInsets.symmetric(vertical: P2PTheme.space4),
                  ),
                  onPressed: submitting ? null : _submit,
                  child: const Text(
                    'Create offer',
                    style: TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                  ),
                ),
              ),
              const SizedBox(height: P2PTheme.space6),
            ],
          ),
        ),
      ),
    );
  }

  Widget _field(
    String label,
    TextEditingController controller, {
    bool numeric = false,
    bool required = true,
  }) {
    return Padding(
      padding: const EdgeInsets.only(bottom: P2PTheme.space3),
      child: TextFormField(
        controller: controller,
        keyboardType: numeric ? const TextInputType.numberWithOptions(decimal: true) : TextInputType.text,
        validator: required ? _required : null,
        decoration: InputDecoration(
          labelText: label,
          filled: true,
          fillColor: Colors.white,
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(P2PTheme.space2),
            borderSide: BorderSide.none,
          ),
        ),
      ),
    );
  }
}

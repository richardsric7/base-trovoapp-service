import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_offer.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/p2p_merchant_gate.dart';
import 'package:trovo_app/widgets/p2p_payment_method_select.dart';
import 'package:trovo_app/widgets/p2p_wallet_payout_select.dart';
import 'package:trovo_app/widgets/popups.dart';

// P2PEditOfferView lets a merchant edit an existing offer's terms (Plan
// Section 13's implied listing-management control). offerType/asset are
// fixed at creation (they drive curated-asset validation and the role
// mapping in Plan Section 11), so they're shown read-only here.
class P2PEditOfferView extends StatefulWidget {
  const P2PEditOfferView({Key? key}) : super(key: key);

  @override
  State<P2PEditOfferView> createState() => _P2PEditOfferViewState();
}

class _P2PEditOfferViewState extends State<P2PEditOfferView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  final formKey = GlobalKey<FormState>();

  P2POffer? offer;
  bool loaded = false;
  bool submitting = false;

  final price = TextEditingController();
  final minOrderAmount = TextEditingController();
  final maxOrderAmount = TextEditingController();
  final addLiquidity = TextEditingController();
  final remark = TextEditingController();
  String? paymentMethodId;
  String? merchantPayoutAddress;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    if (!loaded) {
      loaded = true;
      offer = appState.viewData?[P2PEditOfferViewPageConfig.key]?['offer'];
      if (offer != null) {
        price.text = offer!.price ?? '';
        minOrderAmount.text = offer!.minOrderAmount ?? '';
        maxOrderAmount.text = offer!.maxOrderAmount ?? '';
        paymentMethodId = offer!.paymentMethodId;
        merchantPayoutAddress = offer!.merchantPayoutAddress;
        remark.text = offer!.remark ?? '';
      }
    }
  }

  String? _required(String? v) => (v == null || v.trim().isEmpty) ? 'p2prequired'.tr() : null;

  Future<void> _submit() async {
    if (offer == null || !formKey.currentState!.validate()) return;
    if (offer!.offerType == 'SELL' && (paymentMethodId == null || paymentMethodId!.isEmpty)) {
      popup(context, title: 'p2pcouldnotsaveoffer'.tr(), message: 'p2pselectpaymentmethod'.tr());
      return;
    }
    if (offer!.offerType == 'BUY' && (merchantPayoutAddress == null || merchantPayoutAddress!.isEmpty)) {
      popup(context, title: 'p2pcouldnotsaveoffer'.tr(), message: 'p2pselectwallet'.tr());
      return;
    }
    setState(() => submitting = true);
    showLoader(context);
    try {
      final response = await api.updateOffer(offer!.id!, {
        'offerType': offer!.offerType,
        'asset': offer!.asset,
        if (offer!.offerType == 'SELL') 'paymentMethodId': paymentMethodId,
        if (offer!.offerType == 'BUY') 'merchantPayoutAddress': merchantPayoutAddress,
        'country': offer!.country,
        'countryCode': offer!.countryCode,
        'currency': offer!.currency,
        'priceType': offer!.priceType,
        'price': price.text.trim(),
        'minOrderAmount': minOrderAmount.text.trim(),
        'maxOrderAmount': maxOrderAmount.text.trim(),
        'availableLiquidity': addLiquidity.text.trim(),
        'remark': remark.text.trim(),
      });
      hideLoader(context);
      setState(() => submitting = false);
      if (response['statusCode'] == 200) {
        appState.currentAction = PageAction(
          state: PageState.replace,
          page: P2PMyOffersViewPageConfig,
        );
      } else {
        popup(
          context,
          title: 'p2pcouldnotsaveoffer'.tr(),
          message: response['data']?['message'] ?? 'p2pcheckinputsretry'.tr(),
        );
      }
    } catch (e) {
      hideLoader(context);
      setState(() => submitting = false);
      popup(context, title: 'error'.tr(), message: e.toString());
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
        title: Text('p2peditoffertitle'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: offer == null
          ? const Center(child: CircularProgressIndicator())
          : P2PMerchantGate(child: Form(
              key: formKey,
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(P2PTheme.space4),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '${offer!.offerType} ${offer!.asset} (${offer!.currency})',
                      style: const TextStyle(color: Colors.black54),
                    ),
                    const SizedBox(height: P2PTheme.space3),
                    _field('p2ppriceperunit'.tr(), price, numeric: true),
                    _field('p2pminorderamount'.tr(), minOrderAmount, numeric: true),
                    _field('p2pmaxorderamount'.tr(), maxOrderAmount, numeric: true),
                    _field(
                      'p2paddliquiditylabel'.tr(args: [offer!.availableLiquidity ?? '', offer!.asset ?? '']),
                      addLiquidity,
                      numeric: true,
                      required: false,
                    ),
                    const SizedBox(height: P2PTheme.space2),
                    offer!.offerType == 'SELL'
                        ? P2PPaymentMethodSelect(
                            api: api,
                            value: paymentMethodId,
                            onChanged: (v) => setState(() => paymentMethodId = v),
                          )
                        : P2PWalletPayoutSelect(
                            value: merchantPayoutAddress,
                            onChanged: (v) => setState(() => merchantPayoutAddress = v),
                          ),
                    const SizedBox(height: P2PTheme.space3),
                    _field('p2premark'.tr(), remark, required: false),
                    const SizedBox(height: P2PTheme.space4),
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                          backgroundColor: P2PTheme.brandDark,
                          padding: const EdgeInsets.symmetric(vertical: P2PTheme.space4),
                        ),
                        onPressed: submitting ? null : _submit,
                        child: Text(
                          'p2psavechanges'.tr(),
                          style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                        ),
                      ),
                    ),
                    const SizedBox(height: P2PTheme.space6),
                  ],
                ),
              ),
            )),
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

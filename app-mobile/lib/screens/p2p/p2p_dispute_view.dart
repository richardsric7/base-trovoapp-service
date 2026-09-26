import 'package:easy_localization/easy_localization.dart';
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

// P2PDisputeView (Plan Section 95.9): file a dispute (subject + description).
// No live chat - the legacy schema's order_chat_messages table was
// explicitly decided to be irrelevant prior art, not something to add here.
class P2PDisputeView extends StatefulWidget {
  const P2PDisputeView({Key? key}) : super(key: key);

  @override
  State<P2PDisputeView> createState() => _P2PDisputeViewState();
}

class _P2PDisputeViewState extends State<P2PDisputeView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  final formKey = GlobalKey<FormState>();
  final description = TextEditingController();
  String subject = 'PAYMENT_NOT_RECEIVED';
  bool loaded = false;
  String? orderId;
  bool submitting = false;

  static Map<String, String> get subjects => {
        'NO_PAYMENT': 'p2pdisputesubjectnopayment'.tr(),
        'UNDER_PAYMENT': 'p2pdisputesubjectunderpayment'.tr(),
        'OVER_PAYMENT': 'p2pdisputesubjectoverpayment'.tr(),
        'PAYMENT_NOT_RECEIVED': 'p2pdisputesubjectpaymentnotreceived'.tr(),
        'PAYMENT_MARKED_SENT_IN_ERROR': 'p2pdisputesubjectpaymentmarkedsentinerror'.tr(),
        'WRONG_PAYMENT_AMOUNT': 'p2pdisputesubjectwrongpaymentamount'.tr(),
        'OTHER': 'p2pdisputesubjectother'.tr(),
      };

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    if (!loaded) {
      loaded = true;
      orderId = appState.viewData?[P2PDisputeViewPageConfig.key]?['orderId'];
    }
  }

  Future<void> _submit() async {
    if (orderId == null) return;
    if (!formKey.currentState!.validate()) return;
    setState(() => submitting = true);
    showLoader(context);
    try {
      final response = await api.openDispute(
        orderId: orderId!,
        subject: subject,
        description: description.text.trim(),
      );
      hideLoader(context);
      setState(() => submitting = false);
      if (response['statusCode'] == 201) {
        appState.viewData ??= {};
        appState.viewData![P2POrderDetailViewPageConfig.key] = {'orderId': orderId};
        appState.currentAction = PageAction(state: PageState.replace, page: P2POrderDetailViewPageConfig);
      } else {
        popup(context, title: 'p2pcouldnotopendispute'.tr(), message: response['data']?['message'] ?? 'p2ppleasetryagain'.tr());
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
        title: Text('p2praiseadispute'.tr(), style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: Form(
        key: formKey,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(P2PTheme.space4),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text('p2pwhatwentwrong'.tr(), style: const TextStyle(fontWeight: FontWeight.w600)),
              const SizedBox(height: P2PTheme.space2),
              Container(
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(P2PTheme.space2),
                ),
                padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space3),
                child: DropdownButtonFormField<String>(
                  value: subject,
                  decoration: const InputDecoration(border: InputBorder.none),
                  items: subjects.entries
                      .map((e) => DropdownMenuItem(value: e.key, child: Text(e.value)))
                      .toList(),
                  onChanged: (v) => setState(() => subject = v!),
                ),
              ),
              const SizedBox(height: P2PTheme.space4),
              Text('p2pdescribewhathappened'.tr(), style: const TextStyle(fontWeight: FontWeight.w600)),
              const SizedBox(height: P2PTheme.space2),
              TextFormField(
                controller: description,
                maxLines: 5,
                validator: (v) => (v == null || v.trim().length < 10) ? 'p2pneedmoredetails'.tr() : null,
                decoration: InputDecoration(
                  filled: true,
                  fillColor: Colors.white,
                  hintText: 'p2pdisputehint'.tr(),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(P2PTheme.space2),
                    borderSide: BorderSide.none,
                  ),
                ),
              ),
              const SizedBox(height: P2PTheme.space6),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  style: ElevatedButton.styleFrom(
                    backgroundColor: P2PTheme.danger,
                    padding: const EdgeInsets.symmetric(vertical: P2PTheme.space4),
                  ),
                  onPressed: submitting ? null : _submit,
                  child: Text('p2psubmitdispute'.tr(), style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

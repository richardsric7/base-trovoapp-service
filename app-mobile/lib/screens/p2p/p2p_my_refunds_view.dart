import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_refund.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

// P2PMyRefundsView: a refund is recorded (overpayment, an order that
// expired before payment, or a dispute resolved in the buyer's favor) but
// funds only actually move back to the depositor when they claim it here -
// this screen is that missing claim step (Plan Section 45).
class P2PMyRefundsView extends StatefulWidget {
  const P2PMyRefundsView({Key? key}) : super(key: key);

  @override
  State<P2PMyRefundsView> createState() => _P2PMyRefundsViewState();
}

class _P2PMyRefundsViewState extends State<P2PMyRefundsView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  bool loading = true;
  List<P2PRefund> refunds = [];
  String? claimingId;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final result = await api.listMyRefunds();
    setState(() {
      refunds = result;
      loading = false;
    });
  }

  Future<void> _claim(P2PRefund refund) async {
    setState(() => claimingId = refund.id);
    showLoader(context);
    try {
      final response = await api.claimRefund(refund.id!, address: appState.primaryWallet.address!);
      hideLoader(context);
      setState(() => claimingId = null);
      if (response['statusCode'] == 200) {
        _load();
      } else {
        popup(context, title: 'Could not claim refund', message: response['data']?['message'] ?? 'Please try again.');
      }
    } catch (e) {
      hideLoader(context);
      setState(() => claimingId = null);
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
        title: Text('My refunds', style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: loading
          ? const Center(child: CircularProgressIndicator())
          : refunds.isEmpty
              ? const P2PEmptyState(icon: Icons.savings_outlined, message: 'You have no refunds.')
              : RefreshIndicator(
                  onRefresh: _load,
                  child: ListView.builder(
                    padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4, vertical: P2PTheme.space2),
                    itemCount: refunds.length,
                    itemBuilder: (context, i) {
                      final r = refunds[i];
                      return P2PListCard(
                        child: Row(
                          children: [
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text('${r.amount} ${r.token}', style: const TextStyle(fontWeight: FontWeight.w700)),
                                  const SizedBox(height: P2PTheme.space1),
                                  Text(r.reason ?? '', style: const TextStyle(color: Colors.black54, fontSize: 12)),
                                ],
                              ),
                            ),
                            if (r.claimed == true)
                              const P2PStatusPill(status: 'COMPLETED')
                            else
                              ElevatedButton(
                                style: ElevatedButton.styleFrom(backgroundColor: P2PTheme.brandDark),
                                onPressed: claimingId == r.id ? null : () => _claim(r),
                                child: const Text('Claim', style: TextStyle(color: Colors.white)),
                              ),
                          ],
                        ),
                      );
                    },
                  ),
                ),
    );
  }
}

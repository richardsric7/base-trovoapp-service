import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_order.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

// P2PMyOrdersView (Plan Section 95.10): My Orders/My Trades - a filterable,
// paginated list, the direct mobile analog of payment_history.dart's
// filter+list shape (Plan Section 91).
class P2PMyOrdersView extends StatefulWidget {
  const P2PMyOrdersView({Key? key}) : super(key: key);

  @override
  State<P2PMyOrdersView> createState() => _P2PMyOrdersViewState();
}

class _P2PMyOrdersViewState extends State<P2PMyOrdersView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late P2PApi api;
  bool loading = true;
  String role = '';
  List<P2POrder> orders = [];

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
  }

  Future<void> _load() async {
    setState(() => loading = true);
    final result = await api.listMyOrders(role: role.isEmpty ? null : role);
    setState(() {
      orders = result['orders'];
      loading = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    final myUsername = appState.userInfo?.username ?? '';

    return Scaffold(
      backgroundColor: P2PTheme.neutralBg,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        title: Text('My orders', style: TextStyle(color: notifier.getblck)),
        iconTheme: IconThemeData(color: notifier.getblck),
      ),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(P2PTheme.space4),
            child: Row(
              children: [
                _tab('', 'All'),
                const SizedBox(width: P2PTheme.space2),
                _tab('customer', 'Buying/Selling as customer'),
                const SizedBox(width: P2PTheme.space2),
                _tab('merchant', 'As merchant'),
              ],
            ),
          ),
          Expanded(
            child: loading
                ? const Center(child: CircularProgressIndicator())
                : orders.isEmpty
                    ? P2PEmptyState(
                        icon: Icons.receipt_long_outlined,
                        message: 'You have no orders yet.',
                        ctaLabel: 'Browse the marketplace',
                        onCta: () {
                          appState.currentAction = PageAction(
                            state: PageState.pop,
                          );
                        },
                      )
                    : RefreshIndicator(
                        onRefresh: _load,
                        child: ListView.builder(
                          padding: const EdgeInsets.symmetric(horizontal: P2PTheme.space4),
                          itemCount: orders.length,
                          itemBuilder: (context, i) => _orderCard(orders[i], myUsername),
                        ),
                      ),
          ),
        ],
      ),
    );
  }

  Widget _tab(String value, String label) {
    final selected = role == value;
    return Expanded(
      child: GestureDetector(
        onTap: () {
          setState(() => role = value);
          _load();
        },
        child: Container(
          padding: const EdgeInsets.symmetric(vertical: P2PTheme.space2),
          decoration: BoxDecoration(
            color: selected ? P2PTheme.brandDark : Colors.white,
            borderRadius: BorderRadius.circular(P2PTheme.space2),
          ),
          child: Center(
            child: Text(
              label,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: 11,
                color: selected ? Colors.white : Colors.black87,
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _orderCard(P2POrder o, String myUsername) {
    final counterparty = o.isCustomer(myUsername) ? o.merchantUsername : o.customerUsername;
    return P2PListCard(
      onTap: () {
        appState.viewData ??= {};
        appState.viewData![P2POrderDetailViewPageConfig.key] = {'orderId': o.id};
        appState.currentAction = PageAction(
          state: PageState.addPage,
          page: P2POrderDetailViewPageConfig,
        );
      },
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('${o.specifiedAssetAmount} ${o.asset}', style: const TextStyle(fontWeight: FontWeight.w700)),
                const SizedBox(height: P2PTheme.space1),
                Text('with $counterparty', style: const TextStyle(color: Colors.black54, fontSize: 12)),
              ],
            ),
          ),
          P2PStatusPill(status: o.orderStatus ?? '', disputed: o.isDisputed ?? false),
        ],
      ),
    );
  }
}

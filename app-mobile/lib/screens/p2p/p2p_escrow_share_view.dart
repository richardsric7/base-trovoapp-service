import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:share_plus/share_plus.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/p2p_order.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';

// P2PEscrowShareView (Plan Section 95.5/97.7): the feature's "money moment"
// - the screen most likely seen by a third party who isn't even a Trovo
// user yet (Plan Section 33), so it gets the most visual care in the whole
// feature. Modeled directly on request_specific_payment_details.dart's
// shape: rounded info card, server-generated QR, copy + share actions.
class P2PEscrowShareView extends StatefulWidget {
  const P2PEscrowShareView({Key? key}) : super(key: key);

  @override
  State<P2PEscrowShareView> createState() => _P2PEscrowShareViewState();
}

class _P2PEscrowShareViewState extends State<P2PEscrowShareView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  bool loaded = false;
  P2POrder? order;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (!loaded) {
      loaded = true;
      order = appState.viewData?[P2PEscrowShareViewPageConfig.key]?['order'];
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);

    return Scaffold(
      backgroundColor: P2PTheme.brandDark,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.white),
        title: Text('p2pshareescrowdeposit'.tr(), style: const TextStyle(color: Colors.white)),
      ),
      body: order == null
          ? P2PEmptyState(icon: Icons.error_outline, message: 'p2pnoorderselected'.tr())
          : SingleChildScrollView(
              padding: const EdgeInsets.all(P2PTheme.space6),
              child: Column(
                children: [
                  Text(
                    '${order!.sellerEscrowAssetAmount} ${order!.asset}',
                    style: const TextStyle(color: Colors.white, fontSize: 28, fontWeight: FontWeight.w700),
                  ),
                  const SizedBox(height: P2PTheme.space1),
                  Text('p2panyonewithlink'.tr(), style: const TextStyle(color: Colors.white70)),
                  const SizedBox(height: P2PTheme.space8),
                  if (order!.escrowDepositQrCode != null && order!.escrowDepositQrCode!.isNotEmpty)
                    Container(
                      padding: const EdgeInsets.all(P2PTheme.space4),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(P2PTheme.sheetRadius),
                      ),
                      child: Image.network(order!.escrowDepositQrCode!, width: 220, height: 220),
                    ),
                  const SizedBox(height: P2PTheme.space6),
                  Container(
                    padding: const EdgeInsets.all(P2PTheme.space4),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(P2PTheme.cardRadius),
                    ),
                    child: Row(
                      children: [
                        Expanded(
                          child: Text(
                            order!.escrowDepositShortlink ?? '',
                            overflow: TextOverflow.ellipsis,
                            style: const TextStyle(fontSize: 13),
                          ),
                        ),
                        IconButton(
                          icon: const Icon(Icons.copy, size: 18),
                          onPressed: () {
                            Clipboard.setData(ClipboardData(text: order!.escrowDepositShortlink ?? ''));
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(content: Text('p2plinkcopied'.tr())),
                            );
                          },
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: P2PTheme.space4),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton.icon(
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(vertical: P2PTheme.space4),
                      ),
                      icon: const Icon(Icons.share, color: P2PTheme.brandDark),
                      label: Text('p2psharebutton'.tr(), style: const TextStyle(color: P2PTheme.brandDark, fontWeight: FontWeight.w600)),
                      onPressed: () {
                        Share.share(order!.escrowDepositShortlink ?? '');
                      },
                    ),
                  ),
                  const SizedBox(height: P2PTheme.space2),
                  TextButton(
                    onPressed: () => Navigator.of(context).pop(),
                    child: Text('p2pbacktoorder'.tr(), style: const TextStyle(color: Colors.white70)),
                  ),
                ],
              ),
            ),
    );
  }
}

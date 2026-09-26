import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/network/p2p_requests.dart';
import 'package:trovo_app/storage/state.dart';

// P2PMerchantGate wraps a merchant-only screen (create/edit offer, my
// offers): a caller who isn't a merchant yet sees a notice explaining
// that only merchants can do this, with a button to request merchant
// status right there, instead of the screen's own content. The backend
// enforces the same rule itself (CreateOffer rejects a non-merchant) -
// this is the UX half of that rule, not a substitute for it.
class P2PMerchantGate extends StatefulWidget {
  final Widget child;
  final VoidCallback? onBackToMarketplace;
  const P2PMerchantGate({super.key, required this.child, this.onBackToMarketplace});

  @override
  State<P2PMerchantGate> createState() => _P2PMerchantGateState();
}

class _P2PMerchantGateState extends State<P2PMerchantGate> {
  late DataProvider appState;
  late P2PApi api;
  bool loading = true;
  bool requesting = false;
  bool isMerchant = false;
  String error = '';

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    api = P2PApi(appState);
    _load();
  }

  Future<void> _load() async {
    final status = await api.getMerchantStatus();
    if (!mounted) return;
    setState(() {
      isMerchant = status?['isMerchant'] == true;
      loading = false;
    });
  }

  Future<void> _requestStatus() async {
    setState(() {
      requesting = true;
      error = '';
    });
    try {
      final response = await api.requestMerchantStatus();
      if (!mounted) return;
      if (response['statusCode'] == 200) {
        setState(() {
          isMerchant = true;
          requesting = false;
        });
      } else {
        setState(() {
          requesting = false;
          error = response['data']?['message'] ?? 'p2pcouldnotrequestmerchant'.tr();
        });
      }
    } catch (e) {
      if (!mounted) return;
      setState(() {
        requesting = false;
        error = e.toString();
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    appState = Provider.of<DataProvider>(context, listen: true);

    if (loading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (isMerchant) {
      return widget.child;
    }

    return Center(
      child: Padding(
        padding: const EdgeInsets.all(P2PTheme.space8),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(Icons.storefront_outlined, size: 48, color: P2PTheme.warning.withOpacity(0.6)),
            const SizedBox(height: P2PTheme.space4),
            Text(
              'p2ponlymerchants'.tr(),
              textAlign: TextAlign.center,
              style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 17, color: P2PTheme.brandDark),
            ),
            const SizedBox(height: P2PTheme.space2),
            Text(
              'p2ponlymerchantsbody'.tr(),
              textAlign: TextAlign.center,
              style: const TextStyle(color: Colors.black54, fontSize: 15),
            ),
            if (error.isNotEmpty) ...[
              const SizedBox(height: P2PTheme.space3),
              Text(error, textAlign: TextAlign.center, style: const TextStyle(color: P2PTheme.danger)),
            ],
            const SizedBox(height: P2PTheme.space6),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                style: ElevatedButton.styleFrom(
                  backgroundColor: P2PTheme.brandDark,
                  padding: const EdgeInsets.symmetric(vertical: P2PTheme.space4),
                ),
                onPressed: requesting ? null : _requestStatus,
                child: Text(
                  requesting ? 'p2prequesting'.tr() : 'p2pbecomeamerchant'.tr(),
                  style: const TextStyle(color: Colors.white, fontWeight: FontWeight.w600),
                ),
              ),
            ),
            if (widget.onBackToMarketplace != null) ...[
              const SizedBox(height: P2PTheme.space4),
              TextButton(onPressed: widget.onBackToMarketplace, child: Text('p2pbacktomarketplace'.tr())),
            ],
          ],
        ),
      ),
    );
  }
}

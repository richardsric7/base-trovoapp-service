import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/models/public_market.dart';
import 'package:trovo_app/network/public_markets_requests.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/screens/public_markets/pm_widgets.dart';
import 'package:trovo_app/screens/public_markets/public_markets_view.dart';
import 'package:trovo_app/storage/state.dart';

// PMHomeSection is the home screen's "Public Markets" block: the first
// few listed assets and a link to the full list. It hides itself when
// nothing is listed (or the request fails).
class PMHomeSection extends StatefulWidget {
  const PMHomeSection({Key? key}) : super(key: key);

  @override
  State<PMHomeSection> createState() => _PMHomeSectionState();
}

class _PMHomeSectionState extends State<PMHomeSection> {
  List<PMAsset> assets = [];
  bool started = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (!started) {
      started = true;
      PublicMarketsApi(Provider.of<DataProvider>(context, listen: false)).assets().then((r) {
        if (mounted && r != null) setState(() => assets = r.where((a) => a.status != 'coming-soon').take(3).toList());
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (assets.isEmpty) return const SizedBox.shrink();
    final appState = Provider.of<DataProvider>(context, listen: false);
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Row(children: [
        Text('pmtitle'.tr(), style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w600, color: P2PTheme.brandDark)),
        const SizedBox(width: 6),
        Text('pmhomesubtitle'.tr(), style: const TextStyle(fontSize: 12, color: Colors.black45)),
        const Spacer(),
        TextButton(
          onPressed: () => pmOpen(appState, PublicMarketsViewPageConfig),
          style: TextButton.styleFrom(padding: EdgeInsets.zero, minimumSize: Size.zero, tapTargetSize: MaterialTapTargetSize.shrinkWrap),
          child: Text('pmviewall'.tr()),
        ),
      ]),
      const SizedBox(height: P2PTheme.space2),
      for (final a in assets) PMAssetTile(asset: a, onTap: () => pmOpen(appState, PMAssetViewPageConfig, {'assetCode': a.assetCode})),
    ]);
  }
}

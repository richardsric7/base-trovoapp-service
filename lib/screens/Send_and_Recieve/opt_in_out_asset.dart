import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/models/curated_asset.dart';
import 'package:trovo_wallet/models/user.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:collection/src/list_extensions.dart';

class OptInOutAsset extends StatefulWidget {
  const OptInOutAsset({Key? key}) : super(key: key);

  @override
  State<OptInOutAsset> createState() => _OptInOutAssetState();
}

class _OptInOutAssetState extends State<OptInOutAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet? wallet;
  late bool hasAvailableBalance;
  int? selectedWalletIndex = null;
  List<CuratedAsset>? curatedAssets;

  List<DropdownMenuItem<int>> walletDropdownItems(bool isSelected) {
    return appState.userInfo!
        .transactionableWallets()
        .mapIndexed<DropdownMenuItem<int>>((index, wallet) {
      return DropdownMenuItem(
        child: Row(
          children: [
            Container(
              constraints: isSelected
                  ? BoxConstraints(maxWidth: width / 4)
                  : BoxConstraints(maxWidth: width / 2.5),
              child: Text(
                wallet.alias!,
                overflow:
                    isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
              ),
            ),
            if (wallet.isSharedWallet) ...[
              SizedBox(
                width: 2,
              ),
              Icon(
                Icons.people_outline,
                size: 17,
                color: notifier.getbluecolor,
              )
            ],
            if (!isSelected && index == selectedWalletIndex) ...[
              SizedBox(
                width: 2,
              ),
              Icon(
                Icons.check,
                size: 18,
                color: notifier.getbluecolor,
              )
            ],
          ],
        ),
        value: index,
        // ),
      );
    }).toList();
  }

  @override
  void initState() {
    super.initState();

    appState = Provider.of<DataProvider>(context, listen: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    if (selectedWalletIndex != null) {
      wallet =
          appState.userInfo!.transactionableWallets()[selectedWalletIndex!];
    }
    curatedAssets = appState.userInfo!.curatedSwapList!
        .where((asset) => asset.assetCode != '' && asset.assetIssuer != '')
        .toList();

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(height / 15),
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            title: Text(
              'Add/Remove Asset',
              style: TextStyle(
                  fontSize: 22,
                  fontWeight: FontWeight.bold,
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold),
            ),
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              SizedBox(
                height: height / 30,
              ),
              Row(
                children: [
                  SizedBox(
                    width: width / 15,
                  ),
                  Text(
                    'Select wallet',
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500),
                  ),
                  SizedBox(
                    width: width / 15,
                  ),
                  Expanded(
                      child: dropdown(
                    (newValue) {
                      setState(() {
                        selectedWalletIndex = int.parse(newValue.toString());
                      });
                    },
                    walletDropdownItems(false),
                    selectedWalletIndex.toString().isEmpty
                        ? null
                        : selectedWalletIndex,
                    null,
                    context,
                    (context) {
                      return walletDropdownItems(true);
                    },
                  )),
                  SizedBox(
                    width: width / 15,
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              if (selectedWalletIndex != null) ...[
                showCuratedAssets(),
              ] else ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(15.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20, vertical: 10.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              Row(
                                children: [
                                  Container(
                                    width: width / 1.28,
                                    child: Text(
                                      'Select the wallet on which you want to add or remove an asset.',
                                      textAlign: TextAlign.justify,
                                      style: TextStyle(
                                          fontSize: 15,
                                          color: notifier.getbluewhitecolor,
                                          fontFamily: fontbody),
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ]
            ],
          ),
        ),
      ),
    );
  }

  Widget showCuratedAssets() {
    return SingleChildScrollView(
      child: Column(
        children: [
          if (curatedAssets!.length > 0) ...[
            for (var i = 0; i < curatedAssets!.length; i++) ...[
              curatedAssetTiles(curatedAssets![i], wallet!.claimedAssets!),
            ],
            SizedBox(
              height: height / 22,
            ),
          ] else ...[
            Container(
              height: height / 3,
              child: Padding(
                  padding: const EdgeInsets.fromLTRB(10, 28.0, 10, 0),
                  child: Center(
                    child: Text(
                      LanguageEn.noassets,
                      style: TextStyle(
                        fontSize: 13,
                        fontWeight: FontWeight.bold,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                  )),
            ),
          ],
          SizedBox(
            height: height / 22,
          ),
        ],
      ),
    );
  }

  Widget curatedAssetTiles(CuratedAsset asset, List<Asset> claimedAssets) {
    bool removable = claimedAssets.any(
      (claimed) =>
          claimed.assetCode == asset.assetCode &&
          claimed.assetIssuer == asset.assetIssuer,
    );
    print('asset => ${asset.assetCode}');
    print('${wallet!.alias} removable =======> $removable');
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
            title: Row(
              children: [
                Image.network(
                  asset.imageUrl!,
                  height: 35,
                  width: 35,
                  errorBuilder: (context, error, stackTrace) {
                    return Image.asset(
                      'assets/images/trovo.png',
                      height: 35,
                      width: 35,
                    );
                  },
                ),
                SizedBox(width: 20),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      asset.assetCode!,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        asset.assetName!,
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        asset.assetClass!['assetClass'],
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                  ],
                ),
                // Column(
                //   crossAxisAlignment: CrossAxisAlignment.end,
                //   mainAxisAlignment: MainAxisAlignment.center,
                //   children: [
                //     Text(
                //       asset.assetClass!['assetClass'],
                //       style: TextStyle(
                //         fontSize: 12,
                //         fontFamily: fontsemibold,
                //         color: notifier.getblck,
                //       ),
                //     ),
                //     Padding(
                //       padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                //       child: Text(
                //         asset.organization!,
                //         style: TextStyle(
                //           fontSize: 9,
                //           fontFamily: fontbody,
                //           color: notifier.getblck,
                //         ),
                //       ),
                //     ),
                //   ],
                // )
              ],
            ),
            trailing: ElevatedButton(
              onPressed: () async {
                appState.viewData = {
                  'assetCode': asset.assetCode,
                  'assetIssuer': asset.assetIssuer,
                  'walletPublicKey': wallet!.publicKey,
                };
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: removable
                      ? OptOutAssetViewPageConfig
                      : OptInAssetViewPageConfig,
                );
              },
              style: ButtonStyle(
                backgroundColor: MaterialStateProperty.all<Color>(
                    removable ? Colors.red[400]! : notifier.getbluewhitecolor),
              ),
              child: Text(
                removable ? 'Remove' : 'Add',
                style: TextStyle(
                  fontFamily: fontsemibold,
                  fontSize: 9,
                ),
              ),
            )),
      ),
    );
  }

  @override
  void dispose() {
    super.dispose();
    print('disposing...');
    appState.viewData![OptInOutAssetViewPageConfig.key] = null;
  }
}

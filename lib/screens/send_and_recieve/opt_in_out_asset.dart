import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
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
  late RefreshController _refreshController;
  int? selectedWalletIndex = null;
  Map<String, Map<String, dynamic>> assets = {};

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
    _refreshController = RefreshController(initialRefresh: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    if (selectedWalletIndex != null) {
      wallet =
          appState.userInfo!.transactionableWallets()[selectedWalletIndex!];

      // get all assets on the curated swap list minus XBN
      var filteredList = appState.userInfo!.curatedSwapList!
          .where((asset) => asset.assetCode != '' && asset.assetIssuer != '')
          .toList();

      // add all the assets to a map object
      filteredList.forEach((asset) {
        assets['${asset.assetCode!}|${asset.assetIssuer}'] = {
          'isRemovable': false,
          'imageUrl': asset.imageUrl,
          'assetName': asset.assetName,
          'assetClass': asset.assetClass!['assetClass'],
          'assetIssuer': asset.assetIssuer,
          'assetCode': asset.assetCode,
        };
      });

      // get all the claimed assets on the selected wallet
      var filteredList2 = wallet!.claimedAssets!
          .where((asset) => asset.assetCode != '' && asset.assetIssuer != '')
          .toList();

      // add all the assets to the map object.
      // since we are using map there wont be any duplicate assets
      // rather new records will overwrite existing ones
      // however we must be careful not to overwrite records that
      // we wont want overwritten.
      filteredList2.forEach((asset) {
        // so we check if the record doesnt already exist on the map and if that
        // is true (which means we are dealing with an uncurated asset)
        // enter just the records that we want entered and leave the other ones
        // null
        if (assets['${asset.assetCode!}|${asset.assetIssuer}'] == null) {
          assets['${asset.assetCode!}|${asset.assetIssuer}'] = {
            'isRemovable': true,
            'imageUrl': asset.imageUrl,
            'assetName': truncatePublicKey(asset.assetIssuer),
            'assetIssuer': asset.assetIssuer,
            'assetCode': asset.assetCode,
          };
        } else {
          // if the record already exists which means its a curated asset that
          // we have already added to our claimed assets then we set the removable
          // to true
          var entry = assets['${asset.assetCode!}|${asset.assetIssuer}'];
          entry!['isRemovable'] = true;
        }
      });
    }

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "addremoveasset".tr(),
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        body: SmartRefresher(
          enablePullDown: true,
          controller: _refreshController,
          onRefresh: refreshData,
          child: SingleChildScrollView(
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
                      "selectwallet".tr(),
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
                                        "selectwallet2".tr(),
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
      ),
    );
  }

  Widget showCuratedAssets() {
    return SingleChildScrollView(
      child: Column(
        children: [
          for (var entry in assets.entries) ...[
            curatedAssetTiles(entry),
          ],
          SizedBox(
            height: height / 50,
          ),
          Button(
            "other".tr(),
            notifier.getbluecolor,
            wihitecolor,
            onTap: () {
              addCustomAssetPopup(context, (assetCode, assetIssuer) {
                appState.viewData = {
                  'assetCode': assetCode,
                  'assetIssuer': assetIssuer,
                  'walletPublicKey': wallet!.publicKey,
                };
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: OptInAssetViewPageConfig,
                );
              });
            },
          ),
          SizedBox(
            height: height / 22,
          ),
        ],
      ),
    );
  }

  Widget curatedAssetTiles(MapEntry entry) {
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
                  entry.value['imageUrl']!,
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
                      entry.value['assetCode']!,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getblck,
                      ),
                    ),
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        entry.value['assetName'],
                        style: TextStyle(
                          fontSize: 9,
                          fontFamily: fontbody,
                          color: notifier.getblck,
                        ),
                      ),
                    ),
                    if (entry.value['assetClass'] != null) ...[
                      Padding(
                        padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                        child: Text(
                          entry.value['assetClass'],
                          style: TextStyle(
                            fontSize: 9,
                            fontFamily: fontbody,
                            color: notifier.getblck,
                          ),
                        ),
                      ),
                    ]
                  ],
                ),
              ],
            ),
            trailing: ElevatedButton(
              onPressed: () async {
                appState.viewData = {
                  'assetCode': entry.value['assetCode']!,
                  'assetIssuer': entry.value['assetIssuer']!,
                  'walletPublicKey': wallet!.publicKey,
                };

                print('viewData: ${appState.viewData}');

                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: entry.value['isRemovable']
                      ? OptOutAssetViewPageConfig
                      : OptInAssetViewPageConfig,
                );
              },
              style: ButtonStyle(
                backgroundColor: MaterialStateProperty.all<Color>(
                  entry.value['isRemovable']
                      ? Colors.red[400]!
                      : notifier.isDark
                          ? notifier.getstructuredbluecolor50.backColor
                          : notifier.getbluewhitecolor,
                ),
              ),
              child: Text(
                entry.value['isRemovable'] ? "remove".tr() : "add".tr(),
                style: TextStyle(
                  fontFamily: fontsemibold,
                  fontSize: 9,
                ),
              ),
            )),
      ),
    );
  }

  void refreshData() async {
    try {
      await appState.refreshData();
      _refreshController.refreshCompleted();
      appState.updateListeners();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  @override
  void dispose() {
    super.dispose();
    print('disposing...');
    appState.viewData![OptInOutAssetViewPageConfig.key] = null;
  }
}

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/utils.dart' hide Trans;
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/asset.dart';
import 'package:trovo_app/models/user.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/top_drop_downs.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ReceiveAsset extends StatefulWidget {
  const ReceiveAsset({Key? key}) : super(key: key);

  @override
  State<ReceiveAsset> createState() => _ReceiveAssetState();
}

class _ReceiveAssetState extends State<ReceiveAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  late Wallet wallet;
  late Asset? asset;
  String selectedWallet = '';
  String selectedAsset = '';

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    userInfo = appState.userInfo!;
    wallet = userInfo.getWallet(
      appState.viewData!['walletPublicKey'],
    );
    asset = wallet.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );

    // free the memory..... lol
    appState.viewData = {};
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

    if (selectedWallet.isEmpty) {
      selectedWallet = wallet.publicKey!;
    }

    if (asset == null) {
      asset = wallet.claimedAssets!.firstWhere(
        (asset) => asset.assetCode == '' && asset.assetIssuer == '',
      );
      selectedAsset = '';
    }

    if (selectedAsset.isEmpty) {
      selectedAsset =
          "${getAssetCode(asset!.assetCode)}|${getAssetIssuer(asset!.assetIssuer)}";
    }

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
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            actions: [
              TopDropdowns(
                onWalletChanged: (newValue) {
                  selectedWallet = newValue;
                  this.wallet = userInfo.getWallet(selectedWallet);

                  this.asset = wallet.claimedAssets!.firstWhereOrNull((x) =>
                      "${getAssetCode(x.assetCode)}|${getAssetIssuer(x.assetIssuer)}" ==
                      selectedAsset);

                  setState(() {});
                },
                onAssetChanged: (newValue) {
                  setState(() {
                    selectedAsset = newValue;
                    newValue =
                        newValue.toString().contains('XBN') ? '|' : newValue;
                    for (var asset in wallet.claimedAssets!) {
                      var splitNewValue = newValue.toString().split('|');
                      if (asset.assetCode == splitNewValue[0] &&
                          asset.assetIssuer == splitNewValue[1]) {
                        this.asset = asset;
                      }
                    }
                  });
                },
                claimedAssets: wallet.claimedAssets!,
                selectedAsset: selectedAsset,
                selectedWallet: selectedWallet,
              ),
              if (appState.walletMode == "Testnet") ...[
                Visibility(
                  visible: true,
                  child: Container(
                    color: Color(0xFFAA453E),
                    width: 18,
                    child: RotatedBox(
                      quarterTurns: 1,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(horizontal: 2),
                            child: Text(
                              "testnet".tr(),
                              style: TextStyle(
                                fontFamily: fontsemibold,
                                color: wihitecolor,
                                fontSize: 11,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ]
            ],
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  SizedBox(
                    width: 20,
                  ),
                  Text(
                    "${"receive".tr()} " + getAssetCode(asset!.assetCode),
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: height / 30,
              ),
              showReceivingWallet(),
              SizedBox(
                height: height / 50,
              ),
              showPublicKey(),
              SizedBox(
                height: height / 50,
              ),
              showQrCode(),
              SizedBox(
                height: height / 20,
              ),
              Button(
                "requestspecificamount".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.viewData = {
                    'walletPublicKey': wallet.publicKey,
                    'assetCode': asset!.assetCode,
                    'assetIssuer': asset!.assetIssuer,
                  };

                  appState.currentAction = PageAction(
                      state: PageState.addPage,
                      page: RequestSpecificPaymentViewPageConfig);
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                "dashboard".tr(),
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: BottomHomePageConfig);
                },
              ),
              SizedBox(height: height / 20),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Padding showReceivingWallet() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20, vertical: 10.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    "receivingwallet".tr(),
                    style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(height: height / 90),
                  Row(
                    children: [
                      Container(
                        width: 250,
                        child: Text(
                          wallet.alias!,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                              fontSize: 20,
                              fontWeight: FontWeight.w600,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold),
                        ),
                      ),
                      IconButton(
                        onPressed: () {
                          Clipboard.setData(
                            ClipboardData(
                              text: wallet.alias!,
                            ),
                          );
                          showSnackBar("walletalias".tr(), context);
                        },
                        icon: Icon(Icons.copy,
                            size: 20, color: notifier.getbluewhitecolor),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Padding showPublicKey() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20, vertical: 10.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    "receivefromnontrovowallet".tr(),
                    style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(height: height / 90),
                  Row(
                    children: [
                      Container(
                        width: 250,
                        child: Text(
                          wallet.publicKey!,
                          style: TextStyle(
                              fontSize: 13,
                              fontWeight: FontWeight.w600,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold),
                        ),
                      ),
                      IconButton(
                        onPressed: () {
                          Clipboard.setData(
                            ClipboardData(
                              text: wallet.publicKey!,
                            ),
                          );
                          showSnackBar("publickey".tr(), context);
                        },
                        icon: Icon(Icons.copy,
                            size: 20, color: notifier.getbluewhitecolor),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget showQrCode() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
          ),
          child: Image.network(asset!.qrCode!)),
    );
  }
}

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
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
import 'package:trovo_app/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class DividendPaymentDetailsView extends StatefulWidget {
  const DividendPaymentDetailsView({Key? key}) : super(key: key);

  @override
  State<DividendPaymentDetailsView> createState() =>
      _DividendPaymentDetailsView();
}

class _DividendPaymentDetailsView extends State<DividendPaymentDetailsView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  Wallet? activeWallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int touchedIndex = -1;
  String password = '';
  String? publicKey;
  double? amount;
  String? assetCode;
  String memo = '';
  late Asset? asset;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    activeWallet = appState.activeWallet;
    assetCode = appState.viewData!['assetCode'];

    asset = activeWallet!.claimedAssets!.firstWhere(
      (asset) =>
          asset.assetCode == appState.viewData!['assetCode'] &&
          asset.assetIssuer == appState.viewData!['assetIssuer'],
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 30),
              Text(
                'Dividend Payment ${"details".tr()}',
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                  fontSize: 22,
                ),
              ),
              SizedBox(height: height / 30),
              Text(
                '+3,000.00 CNGN',
                textAlign: TextAlign.center,
                style: TextStyle(
                  color: notifier.getgreencolor,
                  fontFamily: fontsemibold,
                  fontSize: 20,
                ),
              ),
              SizedBox(height: height / 50),
              Stack(
                alignment: AlignmentDirectional.center,
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                    child: Container(
                      decoration: BoxDecoration(
                        borderRadius: const BorderRadius.all(
                          Radius.circular(15.0),
                        ),
                        color: notifier.isDark
                            ? darktilewhitecolor
                            : notifier.getaddsubwalletgrey,
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(height: height / 90),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Text(
                              "Receiving Wallet",
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 16,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          SizedBox(
                            width: width / 1.3,
                            child: Column(
                              children: [
                                Row(
                                  children: [
                                    Expanded(
                                      flex: 3,
                                      child: Padding(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0,
                                        ),
                                        child: Text(
                                          'Kennis',
                                          style: TextStyle(
                                            fontWeight: FontWeight.w500,
                                            color: notifier.getbluewhitecolor,
                                            fontSize: 18.sp,
                                            fontFamily: fontbody,
                                          ),
                                        ),
                                      ),
                                    ),
                                    Expanded(
                                      flex: 1,
                                      child: IconButton(
                                        padding: EdgeInsets.zero,
                                        onPressed: () => {
                                          Clipboard.setData(
                                            ClipboardData(
                                              text: 'viewData.toPublicKey',
                                            ),
                                          ),
                                          showSnackBar(
                                            "tousername".tr(),
                                            context,
                                          ),
                                        },
                                        icon: Icon(Icons.copy, size: 20),
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                  ],
                                ),
                                Row(
                                  children: [
                                    Expanded(
                                      flex: 3,
                                      child: Padding(
                                        padding: const EdgeInsets.symmetric(
                                          horizontal: 20.0,
                                        ),
                                        child: Text(
                                          truncate(
                                                'viewData.toPublicKey!',
                                                length: 5,
                                              ) +
                                              'viewData.toPublicKey!'.substring(
                                                'viewData.toPublicKey!'.length -
                                                    5,
                                              ),
                                          style: TextStyle(
                                            fontWeight: FontWeight.w500,
                                            color: notifier.getbluewhitecolor,
                                            fontSize: 13,
                                            fontFamily: fontbody,
                                          ),
                                        ),
                                      ),
                                    ),
                                    Expanded(
                                      flex: 1,
                                      child: IconButton(
                                        padding: EdgeInsets.zero,
                                        onPressed: () => {
                                          Clipboard.setData(
                                            ClipboardData(
                                              text: 'viewData.toPublicKey!',
                                            ),
                                          ),
                                          showSnackBar(
                                            "topublickey2".tr(),
                                            context,
                                          ),
                                        },
                                        icon: Icon(Icons.copy, size: 20),
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                            ),
                          ),
                          Divider(height: 5),
                          SizedBox(height: height / 90),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Text(
                              "Asset Token",
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 16,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          SizedBox(height: 5),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Row(
                              spacing: 10,
                              children: [
                                if (asset?.imageUrl != null) ...[
                                  ClipRRect(
                                    borderRadius: BorderRadius.circular(100.0),
                                    child: Image.network(
                                      asset!.imageUrl!,
                                      width: 40,
                                      height: 40,
                                      fit: BoxFit.fill,
                                    ),
                                  ),
                                ] else ...[
                                  Image.asset(
                                    'assets/images/trovo.png',
                                    height: 35,
                                    width: 35,
                                  ),
                                ],
                                Text(
                                  asset!.assetCode!.toUpperCase(),
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontFamily: fontsemibold,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Divider(height: 5),
                          SizedBox(height: height / 90),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Text(
                              "blockchainproof".tr(),
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 16,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          SizedBox(height: 5),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Row(
                              children: [
                                Expanded(
                                  flex: 5,
                                  child: GestureDetector(
                                    onTap: () => appState.goToWebView(
                                      getExplorerBaseUrl(appState.walletMode) +
                                          'viewData.transactionId!',
                                    ),
                                    child: Text(
                                      'viewData.transactionId!',
                                      style: TextStyle(
                                        decoration: TextDecoration.underline,
                                        color: notifier.getbluewhitecolor,
                                        fontSize: 12,
                                        fontWeight: FontWeight.w500,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ),
                                ),
                                Expanded(
                                  flex: 1,
                                  child: IconButton(
                                    onPressed: () => {
                                      Clipboard.setData(
                                        ClipboardData(
                                          text: 'viewData.transactionId!',
                                        ),
                                      ),
                                      showSnackBar(
                                        "transactionid".tr(),
                                        context,
                                      ),
                                    },
                                    icon: Icon(Icons.copy, size: 20),
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          Divider(height: 5),
                          SizedBox(height: height / 50),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Text(
                              "date".tr(),
                              style: TextStyle(
                                fontWeight: FontWeight.w500,
                                color: notifier.getbluewhitecolor,
                                fontSize: 16,
                                fontFamily: fontsemibold,
                              ),
                            ),
                          ),
                          SizedBox(height: height / 50),
                          Padding(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 20.0,
                            ),
                            child: Text(
                              '22 January, 2024  14:35 PM',
                              style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontSize: 13,
                                fontWeight: FontWeight.w500,
                                fontFamily: fontbody,
                              ),
                            ),
                          ),
                          SizedBox(height: height / 50),
                        ],
                      ),
                    ),
                  ),
                  Image.asset(
                    'assets/images/trovo_white.png',
                    height: height / 6.5,
                    color: notifier.isDark
                        ? notifier.getdarkgrey
                        : notifier.getsplashgrey,
                  ),
                ],
              ),
              SizedBox(height: height / 20),
              Button(
                "generatereceipt".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: ShareReceiptViewPageConfig,
                  );

                  // appState.viewData![ShareReceiptViewPageConfig.key] = viewData;
                },
              ),
              SizedBox(height: height / 10),
            ],
          ),
        ),
      ),
    );
  }

  Widget showUserInfo() {
    return Row(
      children: [
        SizedBox(width: width / 20),
        Container(
          width: width / 1.7,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [],
          ),
        ),
      ],
    );
  }
}

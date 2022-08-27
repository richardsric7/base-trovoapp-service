import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TransactionSuccess extends StatefulWidget {
  const TransactionSuccess({Key? key}) : super(key: key);

  @override
  State<TransactionSuccess> createState() => _TransactionSuccess();
}

class _TransactionSuccess extends State<TransactionSuccess>
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
  var viewData;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    activeWallet = appState.activeWallet;
    viewData = appState.viewData![TransactionSuccessViewPageConfig.key];

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 10),
              Center(
                child: Image.asset("assets/images/success.gif",
                    height: height / 10),
              ),
              SizedBox(height: height / 50),
              Text(
                LanguageEn.yourtransactionwassuccessful,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 22.sp),
              ),
              SizedBox(height: height / 30),
              Text(
                '- ${viewData['amount']} ${viewData['assetCode'].toString().isEmpty ? 'XBN' : viewData['assetCode']}',
                style: TextStyle(
                    color: Colors.red,
                    fontFamily: fontsemibold,
                    fontSize: 20.sp),
              ),
              SizedBox(
                height: height / 50,
              ),
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                child: Container(
                  decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                    color: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20.0, 15, 0, 0),
                        child: Text(
                          LanguageEn.sentto,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      showUserInfo(),
                      SizedBox(
                        height: height / 50,
                      ),
                      Divider(
                        height: 5,
                      ),
                      SizedBox(
                        height: height / 90,
                      ),
                      if (viewData['memo'].toString().isNotEmpty) ...[
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 20.0, vertical: 10),
                          child: Text(
                            LanguageEn.formemo,
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 16.sp,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        SizedBox(
                          height: 5,
                        ),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Text(
                            viewData['memo'],
                            style: TextStyle(
                              fontWeight: FontWeight.w500,
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontbody,
                            ),
                          ),
                        ),
                        SizedBox(
                          height: height / 50,
                        ),
                        Divider(
                          height: 5,
                        ),
                      ],
                      SizedBox(
                        height: height / 90,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20.0, vertical: 10),
                        child: Text(
                          LanguageEn.blockchainproof,
                          style: TextStyle(
                            fontWeight: FontWeight.w500,
                            color: notifier.getbluewhitecolor,
                            fontSize: 16.sp,
                            fontFamily: fontsemibold,
                          ),
                        ),
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: Row(
                          children: [
                            Expanded(
                              flex: 5,
                              child: GestureDetector(
                                onTap: () => appState.goToWebView(
                                    bantuBlockchainExplorerBaseUrl +
                                        viewData['transactionId']),
                                child: Text(
                                  viewData['transactionId'],
                                  style: TextStyle(
                                    decoration: TextDecoration.underline,
                                    color: notifier.getbluewhitecolor,
                                    fontSize: 12.sp,
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
                                      text: viewData['transactionId'],
                                    ),
                                  ),
                                  showSnackBar('Transaction ID', context),
                                },
                                icon: Icon(Icons.copy),
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ],
                        ),
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                LanguageEn.dashboard,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.replaceAll,
                    page: BottomHomePageConfig,
                  );
                },
              ),
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget showUserInfo() {
    if (viewData['destination'].toString().length == 56) {
      // destination user is not known so we display only
      // destination public key
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
            children: [
              Container(
                width: width / 1.5,
                child: Text(
                  viewData['destination'].toString(),
                  style: TextStyle(
                    fontWeight: FontWeight.w500,
                    color: notifier.getbluewhitecolor,
                    fontSize: 15.sp,
                    fontFamily: fontbody,
                  ),
                ),
              ),
            ],
          ),
        ),
      );
    }

    return Row(
      children: [
        Padding(
          padding: EdgeInsets.fromLTRB(width / 18, 0, 0, 0),
          child: viewData['destinationThumbnail'].toString().isEmpty
              ? CircleAvatar(
                  radius: 30,
                  backgroundColor: notifier.getaddsubwalletgrey,
                  foregroundImage: AssetImage("assets/images/default-user.png"),
                )
              : CircleAvatar(
                  radius: 30,
                  backgroundColor: notifier.getaddsubwalletgrey,
                  foregroundImage: NetworkImage(
                    viewData['destinationThumbnail'].toString(),
                  ),
                ),
        ),
        SizedBox(
          width: width / 70,
        ),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              viewData['destination'].toString(),
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: notifier.getbluewhitecolor,
                fontSize: 19.sp,
                fontFamily: fontbody,
              ),
            ),
            SizedBox(
              height: 5,
            ),
            Text(
              '${viewData['destinationFirstName']} ${viewData['destinationLastName']}',
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontSize: 12.sp,
                fontWeight: FontWeight.w500,
                fontFamily: fontbody,
              ),
            ),
          ],
        )
      ],
    );
  }
}

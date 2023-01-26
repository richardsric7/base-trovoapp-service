import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class RequestSpecificPaymentDetails extends StatefulWidget {
  const RequestSpecificPaymentDetails({Key? key}) : super(key: key);

  @override
  State<RequestSpecificPaymentDetails> createState() =>
      RequestSpecificPaymentDetailsState();
}

class RequestSpecificPaymentDetailsState
    extends State<RequestSpecificPaymentDetails> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  GlobalKey shareArea = GlobalKey();
  var viewData;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    viewData = appState.viewData!;

    // free the memory..... lol
    appState.viewData = {};
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
                    "Receive ${viewData['amount']} ${getAssetCode(viewData['assetCode'])}",
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
              Column(
                children: [
                  showReceivingWallet(),
                  SizedBox(
                    height: height / 50,
                  ),
                  if (viewData['memo'].toString().isNotEmpty) ...[
                    showMemo(),
                  ],
                  showQrCode(),
                  SizedBox(
                    height: height / 20,
                  ),
                ],
              ),
              Button(
                LanguageEn.share,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  share(
                    'Scan Qrcode or tap link to pay ${viewData['amount']} ${getAssetCode(viewData['assetCode'])} to [${viewData['walletAlias']}] => ${viewData['dynamicLink']}',
                    shareArea,
                  );
                },
              ),
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                LanguageEn.dashboard,
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
                    LanguageEn.receivingwallet,
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
                          viewData['walletAlias'],
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
                              text: viewData['walletAlias'],
                            ),
                          );
                          showSnackBar('Wallet alias', context);
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

  Padding showMemo() {
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
                    LanguageEn.formemo,
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
                          viewData['memo'],
                          style: TextStyle(
                              fontSize: 20,
                              fontWeight: FontWeight.w600,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold),
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
    );
  }

  Widget showQrCode() {
    return RepaintBoundary(
      key: shareArea,
      child: Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Container(
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Image.memory(
                base64.decode(viewData['qrCode'].split(',').last))),
      ),
    );
  }
}

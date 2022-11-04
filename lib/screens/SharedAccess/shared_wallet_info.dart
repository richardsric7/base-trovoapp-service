import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';

class SharedWalletInfo extends StatefulWidget {
  const SharedWalletInfo({Key? key}) : super(key: key);

  @override
  State<SharedWalletInfo> createState() => _SharedWalletInfoState();
}

class _SharedWalletInfoState extends State<SharedWalletInfo> {
  late ColorNotifier notifier;
  late DataProvider appState;
  var viewData;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    viewData = appState.viewData![SharedWalletInfoViewPageConfig.key];
    print(viewData);

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          'Shared Access',
          notifier.getbluewhitecolor,
          height: height / 15,
        ),
        body: Column(
          children: [
            SizedBox(
              height: height / 20,
            ),
            Text(
              viewData['walletAlias'],
              style: TextStyle(
                fontSize: 20,
                fontFamily: fontsemibold,
                color: notifier.getbluecolor,
              ),
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
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 20.0, vertical: 15.0),
                      child: Column(
                        children: [
                          Text(
                            'Description',
                            style: TextStyle(
                              fontSize: 17,
                              fontFamily: fontsemibold,
                              color: notifier.getbluecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 90,
                          ),
                          Text(
                            viewData['walletDescription'],
                            style: TextStyle(
                              fontSize: 16,
                              fontFamily: fontbody,
                              color: notifier.getbluecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Text(
                            'Owner',
                            style: TextStyle(
                              fontSize: 17,
                              fontFamily: fontsemibold,
                              color: notifier.getbluecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 90,
                          ),
                          Text(
                            viewData['owner'],
                            style: TextStyle(
                              fontSize: 16,
                              fontFamily: fontbody,
                              color: notifier.getbluecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 50,
                          ),
                          Text(
                            'Permissions',
                            style: TextStyle(
                              fontSize: 17,
                              fontFamily: fontsemibold,
                              color: notifier.getbluecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 90,
                          ),
                          Container(
                            width: width / 1.3,
                            child: Wrap(
                                alignment: WrapAlignment.center,
                                children: [
                                  Text(
                                    'You have ',
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontFamily: fontbody,
                                      color: notifier.getbluecolor,
                                    ),
                                  ),
                                  Text(
                                    viewData['permissions'][0]
                                        .toString()
                                        .toLowerCase(),
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluecolor,
                                    ),
                                  ),
                                  if (viewData['permissions'].length > 1) ...[
                                    SizedBox(
                                      width: width / 90,
                                    ),
                                    Text(
                                      'and',
                                      style: TextStyle(
                                        fontSize: 16,
                                        fontFamily: fontbody,
                                        color: notifier.getbluecolor,
                                      ),
                                    ),
                                    SizedBox(
                                      width: width / 90,
                                    ),
                                    Text(
                                      viewData['permissions'][1]
                                          .toString()
                                          .toLowerCase(),
                                      style: TextStyle(
                                        fontSize: 16,
                                        fontFamily: fontsemibold,
                                        color: notifier.getbluecolor,
                                      ),
                                    ),
                                    SizedBox(
                                      width: width / 90,
                                    ),
                                  ],
                                  Text(
                                    'access on this wallet',
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontFamily: fontbody,
                                      color: notifier.getbluecolor,
                                    ),
                                  ),
                                ]),
                          ),
                          SizedBox(
                            height: height / 90,
                          ),
                          SizedBox(height: 2),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            SizedBox(
              height: height / 20,
            ),
            Button(
              'View wallet',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                // appState.currentAction = PageAction(
                //     state: PageState.addPage,
                //     page: RecoverAccountViewPageConfig);
              },
            ),
            SizedBox(
              height: height / 50,
            ),
            ButtonOutlined(
              'View transaction history',
              notifier.getbluecolor80,
              wihitecolor,
              onTap: () {
                // appState.currentAction = PageAction(
                //     state: PageState.addPage,
                //     page: RecoverAccountViewPageConfig);
              },
            ),
            SizedBox(height: height / 50),
            ButtonOutlined(
              'Initiate access update',
              notifier.getwihitecolor,
              notifier.getbluewhitecolor,
              onTap: () {
                // appState.currentAction = PageAction(
                //     state: PageState.addPage, page: CreatePasswordPageConfig);
              },
            ),
          ],
        ),
      ),
    );
  }
}

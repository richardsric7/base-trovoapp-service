import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TokenizationWelcome extends StatefulWidget {
  const TokenizationWelcome({Key? key}) : super(key: key);

  @override
  State<TokenizationWelcome> createState() => _TokenizationWelcomeState();
}

class _TokenizationWelcomeState extends State<TokenizationWelcome>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late TabController tabController;

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
    tabController = TabController(length: 2, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBarWithoutLeading(
              context,
              notifier.getwihitecolor,
              txt: 'Asset Tokenization',
              titlecolor: notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Card(
                shadowColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: notifier.isDark
                    ? notifier.getbluecolor90
                    : notifier.getaddsubwalletgrey,
                child: Center(
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        LanguageEn.welcometoassettokenization,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        LanguageEn.welcometoassettokenization2,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 13,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Text(
                        LanguageEn.welcometoassettokenization3,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 13,
                          height: 1.4,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(
                        height: height / 70,
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceAround,
                        children: [
                          SmallButton(
                            'Tokenize asset',
                            notifier.getbluewhitecolor!,
                            notifier.getwihitecolor,
                            onTap: () {
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: TokenizeAssetViewPageConfig,
                              );
                            },
                          ),
                          SmallButtonOutlined(
                            'See my assets',
                            notifier.isDark
                                ? notifier.getbluecolor90
                                : notifier.getaddsubwalletgrey,
                            notifier.getbluewhitecolor,
                            onTap: () {},
                          ),
                        ],
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                    ],
                  ),
                ),
              ),
            ),
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
              child: TabBar(
                controller: tabController,
                labelColor: notifier.getbluewhitecolor,
                indicatorColor: notifier.getbluewhitecolor,
                labelStyle: TextStyle(
                  fontSize: 15.sp,
                  fontFamily: fontsemibold,
                ),
                tabs: [
                  Tab(
                    height: 20,
                    text: 'Primary Offering',
                  ),
                  Tab(
                    height: 20,
                    text: 'All Assets',
                  ),
                ],
              ),
            ),
            Container(
              height: height / 2.25,
              child: TabBarView(controller: tabController, children: [
                SingleChildScrollView(
                  child: Column(
                    children: [
                      for (var i = 10; i >= 0; i--) ...[
                        GestureDetector(
                          onTap: () {
                            appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: TokenizedAssetDetailViewPageConfig,
                            );
                          },
                          child:
                              assetTile('', 'ASSET $i', 'Property', i % 2 == 0),
                        ),
                      ],
                      SizedBox(height: height / 20),
                    ],
                  ),
                ),
                SingleChildScrollView(
                  child: Column(
                    children: [
                      for (var i = 0; i <= 10; i++) ...[
                        assetTile('', 'ASSET $i', 'Property', i % 2 == 0),
                      ],
                      SizedBox(height: height / 20),
                    ],
                  ),
                ),
              ]),
            ),
          ],
        ),
      ),
    );
  }

  Widget assetTile(String imageUrl, String name, String type, isSubscribed) {
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
                imageUrl,
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
                    name,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getblck,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      type,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
          trailing: TextButton(
            onPressed: () async {},
            child: Container(
              width: width / 4,
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    isSubscribed ? 'Subscribed' : 'Subscribe',
                    style: TextStyle(
                        fontFamily: fontsemibold,
                        fontSize: 12,
                        color: notifier.getblck),
                  ),
                  Icon(
                      isSubscribed
                          ? Icons.check_circle
                          : Icons.add_circle_rounded,
                      size: 20,
                      color: notifier.getblck),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

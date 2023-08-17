import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TokenizedAssetsList extends StatefulWidget {
  const TokenizedAssetsList({Key? key}) : super(key: key);

  @override
  State<TokenizedAssetsList> createState() => _TokenizedAssetsListState();
}

class _TokenizedAssetsListState extends State<TokenizedAssetsList>
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
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              'Tokenized Assets',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
              child: TabBar(
                controller: tabController,
                labelColor: notifier.getbluewhitecolor,
                indicatorColor: notifier.getbluewhitecolor,
                labelStyle: TextStyle(
                  fontSize: 13.sp,
                  fontFamily: fontsemibold,
                ),
                tabs: [
                  Tab(
                    height: 20,
                    text: 'Assets Tokens',
                  ),
                  Tab(
                    height: 20,
                    text: 'Tokenized Assets',
                  ),
                ],
              ),
            ),
            Container(
              height: height / 1.16,
              child: TabBarView(controller: tabController, children: [
                SingleChildScrollView(
                  child: Column(
                    children: [
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: MyAssetTokenDetailsViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'ANMF',
                          'Agriculture',
                          '12.4304324',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: MyAssetTokenDetailsViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'BCNH',
                          'Property',
                          '9.1204334',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: MyAssetTokenDetailsViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'CVTL',
                          'Health',
                          '2.3292323',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: MyAssetTokenDetailsViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'DRNFL',
                          'Beverage',
                          '5.0302344',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: MyAssetTokenDetailsViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'ZAAD',
                          'Property',
                          '3.0023231',
                        ),
                      ),
                      SizedBox(height: height / 20),
                    ],
                  ),
                ),
                SingleChildScrollView(
                  child: Column(
                    children: [
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'Titan Properties',
                          'Property',
                          'Pending approval',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'ATLANTIS 1',
                          'Property',
                          'Approved',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'Kings Home',
                          'Property',
                          'Pending approval',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: AssetDashboardViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'Infinity Productions',
                          'Property',
                          'Pending  liquidation',
                        ),
                      ),
                      GestureDetector(
                        onTap: () {
                          appState.currentAction = PageAction(
                            state: PageState.addPage,
                            page: TokenizeAssetViewPageConfig,
                          );
                        },
                        child: assetTile(
                          '',
                          'Heart Realty',
                          'Property',
                          'Rejected',
                        ),
                      ),
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

  Widget assetTile(
    String imageUrl,
    String name,
    String type,
    String status,
  ) {
    return Card(
      elevation: notifier.isDark ? 0 : 3,
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
                height: 30,
                width: 30,
                errorBuilder: (context, error, stackTrace) {
                  return Image.asset(
                    'assets/images/trovo.png',
                    height: 35,
                    width: 35,
                  );
                },
              ),
              SizedBox(width: 10),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    truncate(name, length: 14),
                    style: TextStyle(
                      fontSize: 13,
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
          trailing: Container(
            width: width / 3.1,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                Text(
                  status,
                  style: TextStyle(
                      fontFamily: fontsemibold,
                      fontSize: 11,
                      color: notifier.getblck),
                )
              ],
            ),
          ),
        ),
      ),
    );
  }
}

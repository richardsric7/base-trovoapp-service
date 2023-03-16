import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/screens/asset-tokenization/total_sales.dart';
import 'package:trovo_wallet/widgets/bar_chart.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/price_points.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetDashboard extends StatefulWidget {
  const AssetDashboard({Key? key}) : super(key: key);

  @override
  State<AssetDashboard> createState() => _AssetDashboardState();
}

class _AssetDashboardState extends State<AssetDashboard>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;

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
              'ATLANTIS 1',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Atlantis Estate 1 tokens are fractional tokens that represent part ownership (via investment) of our real estate development project at Atlantis Estate, Lekki, Lagos, Nigeria. ',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
            SizedBox(height: height / 50),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 2.3,
                  height: height / 5.5,
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(15.0),
                    ),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: TextButton(
                      onPressed: () {
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: AssetSubscribersViewPageConfig,
                        );
                      },
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            'Total Subscriptions',
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
                            '100',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 21,
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              pill(
                                '+23.4%',
                                backColor: Color(0xFF4F9A94),
                                foreColor: wihitecolor,
                              ),
                              Icon(
                                Icons.arrow_forward,
                                color: notifier.getbluewhitecolor,
                                size: 18,
                              ),
                            ],
                          )
                        ],
                      ),
                    ),
                  ),
                ),
                SizedBox(
                  width: width / 50,
                ),
                Container(
                  width: width / 2.3,
                  height: height / 5.5,
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(15.0),
                    ),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: TextButton(
                      onPressed: () {
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: TotalSalesViewPageConfig,
                        );
                      },
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            'Total Sale',
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
                            '467 000 TROV',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            '\$4 390.23',
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
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              pill(
                                '+23.4%',
                                backColor: Color(0xFF4F9A94),
                                foreColor: wihitecolor,
                              ),
                              Icon(
                                Icons.arrow_forward,
                                color: notifier.getbluewhitecolor,
                                size: 18,
                              ),
                            ],
                          )
                        ],
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 2.3,
                  height: height / 5.5,
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(15.0),
                    ),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: TextButton(
                      onPressed: () {},
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            'Asset Value',
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
                            'N2 000 000',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            '\$4 390.23',
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
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Container(),
                              Icon(
                                Icons.edit,
                                color: notifier.getbluewhitecolor,
                                size: 18,
                              ),
                            ],
                          )
                        ],
                      ),
                    ),
                  ),
                ),
                SizedBox(
                  width: width / 50,
                ),
                Container(
                  width: width / 2.3,
                  height: height / 5.5,
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(15.0),
                    ),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: TextButton(
                      onPressed: () {},
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            'Price Per Asset',
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
                            '100 TROV',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 70,
                          ),
                          Text(
                            '\$2.20',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 13,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            height: height / 30,
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Container(),
                              Container(),
                            ],
                          )
                        ],
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(height: height / 50),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 9,
                  height: height / 40,
                  child: pill('',
                      backColor: Colors.blueAccent,
                      foreColor: Colors.blueAccent,
                      hideDirectionUp: true),
                ),
                Text(
                  'Sales',
                  textAlign: TextAlign.center,
                  softWrap: true,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                    fontSize: 15,
                  ),
                ),
              ],
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Container(
                  width: width / 9,
                  height: height / 40,
                  child: pill('',
                      backColor: Colors.green,
                      foreColor: Colors.green,
                      hideDirectionUp: true),
                ),
                Text(
                  'Subscribers',
                  textAlign: TextAlign.center,
                  softWrap: true,
                  style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                    fontSize: 15,
                  ),
                ),
              ],
            ),
            SizedBox(height: height / 30),
            Container(
              height: height / 3,
              width: width / 1.2,
              child: BarChartWidget(
                points: pricePoints,
              ),
            ),
            SizedBox(height: height / 50),
            TextButton(
              onPressed: () => showAssetDetailsPopup(
                context,
                'ATLANTIS 1',
                items: <String>[],
              ),
              child: Container(
                width: width / 2.5,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'View Asset Details',
                      style: TextStyle(
                          decoration: TextDecoration.underline,
                          fontFamily: fontsemibold,
                          fontSize: 13,
                          color: notifier.getbluewhitecolor),
                    ),
                    Icon(
                      Icons.arrow_forward,
                      color: notifier.getbluewhitecolor,
                      size: 18,
                    ),
                  ],
                ),
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Button(
              'Payout Proceeds',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {},
            ),
            SizedBox(height: height / 70),
            ButtonOutlined(
              'Liquidate Asset',
              notifier.getwihitecolor,
              Colors.red,
              borderColor: Colors.red,
              onTap: () {},
            ),
            SizedBox(
              height: height / 10,
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

  Widget pill(
    String name, {
    required Color backColor,
    required Color foreColor,
    double? fontSize: 12,
    bool hideDirectionUp = false,
  }) {
    return Padding(
      padding: const EdgeInsets.all(3.0),
      child: Container(
        decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(10.0)),
            color: backColor),
        child: Padding(
          padding: const EdgeInsets.all(5.0),
          child: Wrap(
            alignment: WrapAlignment.center,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              Row(
                children: [
                  Text(
                    name,
                    textAlign: TextAlign.center,
                    softWrap: true,
                    style: TextStyle(
                        color: foreColor,
                        fontFamily: fontbody,
                        fontSize: fontSize),
                  ),
                  if (!hideDirectionUp) ...[
                    Icon(
                      Icons.arrow_upward,
                      color: foreColor,
                      size: 18,
                    ),
                  ],
                ],
              ),
              SizedBox(
                width: width / 70,
              ),
            ],
          ),
        ),
      ),
    );
  }
}

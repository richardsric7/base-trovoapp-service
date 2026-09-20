import 'package:flutter/material.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TotalSales extends StatefulWidget {
  const TotalSales({Key? key}) : super(key: key);

  @override
  State<TotalSales> createState() => _TotalSalesState();
}

class _TotalSalesState extends State<TotalSales> with TickerProviderStateMixin {
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
              'Total Sales',
              notifier.getbluewhitecolor,
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
                      SizedBox(height: height / 70),
                      Text(
                        'Total Sales',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(height: height / 70),
                      Text(
                        '467,000 cNGN',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(height: height / 50),
                      Text(
                        '4,670 ATLANTIS 1',
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(height: height / 50),
                    ],
                  ),
                ),
              ),
            ),
            SizedBox(height: height / 70),
            Column(
              children: [
                assetTile('20', 'Kenny', '10 minutes ago', '2000'),
                assetTile('40', 'Onoja', '10 days ago', '3000'),
                assetTile('80', 'Ric1', '15 days ago', '6000'),
                SizedBox(height: height / 20),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget assetTile(
    String amount,
    String name,
    String timeago,
    String trovAmount,
  ) {
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(15.0)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
          title: Row(
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '$amount ATLANTIS 1',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  SizedBox(height: height / 70),
                  Row(
                    children: [
                      Text(
                        'Purchased by ',
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      Text(
                        name,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: height / 70),
                  Text(
                    timeago,
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ],
          ),
          trailing: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Text(
              '$trovAmount cNGN',
              style: TextStyle(
                fontSize: 15,
                fontFamily: fontsemibold,
                color: notifier.getbluewhitecolor,
              ),
            ),
          ),
        ),
      ),
    );
  }
}

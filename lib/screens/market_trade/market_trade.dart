import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/graph/graph.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/bar_chart.dart';
import 'package:trovo_wallet/widgets/price_points.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class MarketTrade extends StatefulWidget {
  const MarketTrade({Key? key}) : super(key: key);

  @override
  State<MarketTrade> createState() => _MarketTradeState();
}

class _MarketTradeState extends State<MarketTrade>
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
              LanguageEn.markettrade,
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            Wrap(
              spacing: 5,
              runSpacing: 5,
              children: [
                for (var i = 0; i <= 5; i++) ...[
                  chartCard(
                      Image.asset(
                        'assets/images/trovo.png',
                        height: height / 50,
                      ),
                      'TROV/NGN',
                      '27,763.32',
                      '8.46%',
                      true),
                  chartCard(
                    Image.asset(
                      'assets/images/trovo.png',
                      height: height / 50,
                    ),
                    'BTC/USDT',
                    '27,763.32',
                    '8.46%',
                    false,
                  ),
                ]
              ],
            ),
            SizedBox(height: height / 50),
          ],
        ),
      ),
    );
  }

  Widget chartCard(
    Image assetLogo,
    String currency,
    String amount,
    String percentage,
    bool directionUp,
  ) {
    return Container(
      width: width / 2.3,
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
              page: MarketTradeInfoViewPageConfig,
            );
          },
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(
                height: height / 70,
              ),
              Row(
                children: [
                  assetLogo,
                  SizedBox(
                    width: width / 70,
                  ),
                  Text(
                    currency,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Text(
                amount,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '${directionUp ? '+' : '-'}${percentage}',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontbody,
                      color: directionUp ? notifier.getgreencolor : Colors.red,
                    ),
                  ),
                  Image.asset(
                    directionUp
                        ? 'assets/images/graph_up.png'
                        : 'assets/images/graph_down.png',
                    height: directionUp ? height / 50 : height / 60,
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
    );
  }
}

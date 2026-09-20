import 'dart:math';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/widgets/utilities.dart';

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
  List<Map> marketPairs = <Map>[
    {
      'pair': 'TROV/CNGN',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'TROV/USDC',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'TROV/USDT',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'TROV/ETH',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'ETH/CNGN',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'ETH/USDC',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'ETH/USDT',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'TROV/CNGN',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'TROV/USDC',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'TROV/USDT',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'TROV/ETH',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'ETH/CNGN',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'ETH/USDC',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'ETH/USDT',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'TROV/USDT',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'TROV/USDC',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
    {
      'pair': 'TROV/CNGN',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': true,
    },
    {
      'pair': 'TROV/ETH',
      'price': formatNumberShort(Random().nextInt(10000).toDouble()),
      'isGreen': false,
    },
  ];

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
            // ignore: dead_code
            if (false) ...[
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                'DEX Trade',
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              Container(
                height: height / 1.2,
                child: Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Image.asset(
                        "assets/images/transfer.png",
                        height: height / 2.5,
                      ),
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
                                  "welcometoassettokenization2".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontSize: 15,
                                    fontFamily: fontbody,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                                SizedBox(height: height / 70),
                                Text(
                                  "welcometoassettokenization3".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontSize: 13,
                                    height: 1.4,
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
                      SizedBox(height: height / 20),
                      Text(
                        'Coming Soon ...',
                        overflow: TextOverflow.ellipsis,
                        style: TextStyle(
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      SizedBox(height: height / 90),
                    ],
                  ),
                ),
              ),
            ] else ...[
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                'Favorites',
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SmallButton(
                "marketpairs".tr(),
                notifier.getbluewhitecolor,
                notifier.getwihitecolor,
                onTap: () {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: MarketPairsViewPageConfig,
                  );
                },
              ),
              SizedBox(height: height / 50),
              Wrap(
                spacing: 5,
                runSpacing: 5,
                children: [
                  for (var i = 0; i < marketPairs.length; i++) ...[
                    chartCard(
                      Image.asset(
                        'assets/images/trovo.png',
                        height: height / 50,
                      ),
                      marketPairs[i]['pair'],
                      marketPairs[i]['price'],
                      '8.46%',
                      marketPairs[i]['isGreen'],
                    ),
                  ],
                ],
              ),
            ],
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
              SizedBox(height: height / 70),
              Row(
                children: [
                  assetLogo,
                  SizedBox(width: width / 70),
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
              SizedBox(height: height / 70),
              Text(
                amount,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
              SizedBox(height: height / 50),
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
              SizedBox(height: height / 50),
            ],
          ),
        ),
      ),
    );
  }
}

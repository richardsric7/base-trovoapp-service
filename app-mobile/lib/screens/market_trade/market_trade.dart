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
import 'package:trovo_app/network/market_requests.dart';

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
  List<MarketPair>? pairs;
  Set<String> favorites = {};
  bool loading = true;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  Future<void> load() async {
    setState(() => loading = true);
    final favs = await MarketFavorites.load();
    List<MarketPair>? list;
    try {
      list = await MarketApi(appState).pairs();
    } catch (_) {
      list = null;
    }
    if (!mounted) return;
    setState(() {
      favorites = favs;
      pairs = list;
      loading = false;
    });
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    load();
  }

  void openPair(MarketPair pair) {
    appState.viewData ??= {};
    appState.viewData![MarketTradeInfoViewPageConfig.key] = pair.toMap();
    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: MarketTradeInfoViewPageConfig,
    );
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    final all = pairs ?? [];
    final starred =
        all.where((p) => favorites.contains(MarketFavorites.id(p))).toList();
    final shown = starred.isNotEmpty ? starred : all;
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: RefreshIndicator(
        onRefresh: load,
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                starred.isNotEmpty ? 'Favorites' : 'DEX Trade',
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SmallButton(
                "marketpairs".tr(),
                notifier.getbluewhitecolor,
                notifier.getwihitecolor,
                onTap: () async {
                  appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: MarketPairsViewPageConfig,
                  );
                },
              ),
              SizedBox(height: height / 50),
              if (loading && pairs == null)
                Padding(
                  padding: EdgeInsets.only(top: height / 5),
                  child: const CircularProgressIndicator(),
                )
              else if (pairs == null)
                message(
                  'The market could not be loaded. Pull down to try again.',
                )
              else if (all.isEmpty)
                message(
                  'No tokens are open for trading yet. Pull down to check again.',
                )
              else
                Wrap(
                  spacing: 5,
                  runSpacing: 5,
                  children: [
                    for (final pair in shown) ...[
                      chartCard(
                        pair,
                        tokenLogo(pair.base.imageUrl),
                        pair.symbol,
                        '${formatMarketString(pair.lastPrice)} ${pair.counter.code}',
                        '${pair.changePercent24h.abs().toStringAsFixed(2)}%',
                        pair.isUp,
                      ),
                    ],
                  ],
                ),
              SizedBox(height: height / 50),
            ],
          ),
        ),
      ),
    );
  }

  Widget message(String text) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 20, vertical: height / 6),
      child: Text(
        text,
        textAlign: TextAlign.center,
        style: TextStyle(
          fontSize: 14,
          fontFamily: fontbody,
          color: notifier.getbluewhitecolor,
        ),
      ),
    );
  }

  Widget tokenLogo(String url) {
    final fallback = Image.asset('assets/images/trovo.png', height: height / 50);
    if (url.isEmpty) return fallback;
    return Image.network(
      url,
      height: height / 50,
      errorBuilder: (_, __, ___) => fallback,
    );
  }

  Widget chartCard(
    MarketPair pair,
    Widget assetLogo,
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
          onPressed: () => openPair(pair),
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
                    '${directionUp ? '+' : '-'}$percentage',
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

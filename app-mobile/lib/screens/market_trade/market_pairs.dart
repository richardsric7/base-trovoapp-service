import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/network/market_requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/widgets/utilities.dart';

import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class MarketPairs extends StatefulWidget {
  const MarketPairs({Key? key}) : super(key: key);

  @override
  State<MarketPairs> createState() => _MarketPairsState();
}

class _MarketPairsState extends State<MarketPairs>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  List<MarketPair>? pairs;
  Set<String> favorites = {};
  bool loading = true;
  String token = allTokens;
  String search = '';

  static const allTokens = 'All tokens';

  List<String> get options => [
        allTokens,
        ...{for (final p in pairs ?? <MarketPair>[]) p.base.code},
      ];

  List<DropdownMenuItem<String>> get getOptions {
    List<DropdownMenuItem<String>> myOptions = [];
    for (final value in options) {
      myOptions.add(
        DropdownMenuItem(
          child: Text(value, overflow: TextOverflow.ellipsis),
          value: value,
        ),
      );
    }
    return myOptions;
  }

  List<MarketPair> get shown {
    final q = search.trim().toLowerCase();
    return (pairs ?? []).where((p) {
      if (token != allTokens && p.base.code != token) return false;
      if (q.isEmpty) return true;
      return p.symbol.toLowerCase().contains(q) ||
          p.base.name.toLowerCase().contains(q) ||
          p.counter.name.toLowerCase().contains(q);
    }).toList();
  }

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
                "marketpairs".tr(),
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              Row(
                children: [
                  Container(
                    width: width / 2,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 10.0),
                      child: dropdown(
                        (value) => setState(() => token = '$value'),
                        getOptions,
                        token,
                        allTokens,
                        context,
                        null,
                      ),
                    ),
                  ),
                  CustomTextFormField.textField(
                    'Search Pairs',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    45.sp,
                    150.sp,
                    onChanged: (value) => setState(() => search = '$value'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              if (loading && pairs == null)
                Padding(
                  padding: EdgeInsets.only(top: height / 5),
                  child: const CircularProgressIndicator(),
                )
              else if (pairs == null)
                message('The market could not be loaded. Pull down to try again.')
              else if (shown.isEmpty)
                message(pairs!.isEmpty
                    ? 'No tokens are open for trading yet.'
                    : 'No pairs match your search.')
              else
                table(),
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

  Widget table() {
    return Column(
      children: [
        Table(children: getTableRows(shown)),
      ],
    );
  }

  List<TableRow> getTableRows(List<MarketPair> marketPairs) {
    List<TableRow> tableRows = [];

    tableRows.add(
      TableRow(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            child: Text(
              'Favorites',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getsplashgrey,
                fontSize: 13,
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            child: Text(
              'Pairs',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getsplashgrey,
                fontSize: 13,
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            child: Text(
              'Price',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getsplashgrey,
                fontSize: 13,
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
            child: Text(
              'Chart',
              textAlign: TextAlign.center,
              style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getsplashgrey,
                fontSize: 13,
              ),
            ),
          ),
        ],
      ),
    );

    for (var i = 0; i < marketPairs.length; i++) {
      tableRows.add(
        TableRow(
          children: [
            Transform.scale(
              scale: 0.9.sp,
              child: Checkbox(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.all(Radius.circular(5.sp)),
                ),
                activeColor: notifier.isDark
                    ? notifier.getbluecolor50
                    : notifier.getbluecolor90,
                side: BorderSide(
                  color: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor90,
                ),
                value: favorites.contains(MarketFavorites.id(marketPairs[i])),
                onChanged: (bool? value) async {
                  final favs = await MarketFavorites.toggle(marketPairs[i]);
                  if (mounted) setState(() => favorites = favs);
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 15),
              child: Text(
                marketPairs[i].symbol,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 15),
              child: Text(
                '${formatMarketString(marketPairs[i].lastPrice)} ${marketPairs[i].counter.code}',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
            TextButton(
              onPressed: () => openPair(marketPairs[i]),
              child: Image.asset(
                'assets/images/gotochart.png',
                height: height / 50,
              ),
            ),
          ],
        ),
      );
    }

    return tableRows;
  }
}

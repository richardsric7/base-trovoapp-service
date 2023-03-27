import 'package:flutter/material.dart';
import 'package:flutter_html/custom_render.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/utils/enstring.dart';
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
  List<Map> marketPairs = <Map>[
    {
      'pair': 'TROV/CNGN',
      'price': '450 CNGN',
      'isChecked': true,
    },
    {
      'pair': 'TROV/USDC',
      'price': '0.5 USDC',
      'isChecked': false,
    },
    {
      'pair': 'TROV/USDT',
      'price': '0.49 USDT',
      'isChecked': true,
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
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              LanguageEn.marketpairs,
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            table(),
            SizedBox(height: height / 50),
          ],
        ),
      ),
    );
  }

  Widget table() {
    return SingleChildScrollView(
      child: Column(
        children: [
          Table(
            children: getTableRows(
              'TROV/CNGN',
              '450 CNGN',
              true,
            ),
          ),
        ],
      ),
    );
  }

  List<TableRow> getTableRows(String pairs, String price, bool isChecked) {
    List<TableRow> tableRows = [];
    bool checked = isChecked;

    tableRows.add(TableRow(children: [
      Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
        child: Text(
          'Favorites',
          textAlign: TextAlign.center,
          style: TextStyle(
              fontFamily: fontsemibold,
              color: notifier.getsplashgrey,
              fontSize: 13),
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
              fontSize: 13),
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
              fontSize: 13),
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
              fontSize: 13),
        ),
      ),
    ]));

    for (var i = 0; i < marketPairs.length; i++) {
      tableRows.add(
        TableRow(
          children: [
            Transform.scale(
              scale: 0.9.sp,
              child: Checkbox(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.all(
                    Radius.circular(5.sp),
                  ),
                ),
                activeColor: notifier.isDark
                    ? notifier.getbluecolor50
                    : notifier.getbluecolor90,
                side: BorderSide(
                  color: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor90,
                ),
                value: marketPairs[i]['isChecked'],
                onChanged: (bool? value) {
                  setState(() {
                    marketPairs[i]['isChecked'] = value!;
                  });
                },
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 15),
              child: Text(
                marketPairs[i]['pair'],
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
                marketPairs[i]['price'],
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
            TextButton(
              onPressed: () {},
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

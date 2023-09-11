import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  List<String> options = [
    'TROV',
    'USDC',
    'USDT',
    'XBN',
  ];

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
    {
      'pair': 'TROV/XBN',
      'price': '1000 XBN',
      'isChecked': false,
    },
    {
      'pair': 'XBN/CNGN',
      'price': '0.5 CNGN',
      'isChecked': false,
    },
    {
      'pair': 'XBN/USDC',
      'price': '0.0003 USDC',
      'isChecked': true,
    },
    {
      'pair': 'XBN/USDT',
      'price': '0.0003 USDT',
      'isChecked': false,
    },
    {
      'pair': 'TROV/CNGN',
      'price': '200 CNGN',
      'isChecked': false,
    },
    {
      'pair': 'TROV/USDC',
      'price': '0.24 USDC',
      'isChecked': true,
    },
    {
      'pair': 'TROV/USDT',
      'price': '1.5 USDT',
      'isChecked': false,
    },
    {
      'pair': 'TROV/XBN',
      'price': '1000 XBN',
      'isChecked': true,
    },
    {
      'pair': 'XBN/CNGN',
      'price': '0.5 CNGN',
      'isChecked': false,
    },
    {
      'pair': 'XBN/USDC',
      'price': '0.49 USDC',
      'isChecked': false,
    },
    {
      'pair': 'XBN/USDT',
      'price': '0.5 USDT',
      'isChecked': true,
    },
    {
      'pair': 'TROV/USDT',
      'price': '0.4 USDT',
      'isChecked': false,
    },
    {
      'pair': 'TROV/USDC',
      'price': '0.82 USDT',
      'isChecked': false,
    },
    {
      'pair': 'TROV/CNGN',
      'price': '30 CNGN',
      'isChecked': true,
    },
    {
      'pair': 'TROV/XBN',
      'price': '1.4 XBN',
      'isChecked': false,
    },
  ];

  List<DropdownMenuItem<String>> get getOptions {
    List<DropdownMenuItem<String>> myOptions = [];
    options.forEach((value) {
      myOptions.add(DropdownMenuItem(
          child: Text(
            value,
            overflow: TextOverflow.ellipsis,
          ),
          value: value));
    });
    return myOptions;
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
                      (value) {},
                      getOptions,
                      null,
                      'Insurance',
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
                  // controller: referrerController,
                  // validator: validateReferrer,
                  onSaved: (value) {},
                ),
              ],
            ),
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
    // bool checked = isChecked;

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

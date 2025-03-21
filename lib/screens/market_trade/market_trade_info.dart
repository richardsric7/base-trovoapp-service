import 'dart:convert';
import 'dart:math';

import 'package:candlesticks/candlesticks.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:http/http.dart' as http;

class MarketTradeInfo extends StatefulWidget {
  const MarketTradeInfo({Key? key}) : super(key: key);

  @override
  State<MarketTradeInfo> createState() => _MarketTradeInfoState();
}

class _MarketTradeInfoState extends State<MarketTradeInfo>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  List<Candle> candles = [];
  bool themeIsDark = false;
  late TabController tabController;

  List<String> options = [
    'Hour',
    'Day',
    'Week',
    'Month',
    'Year',
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

  final candle = Candle(
    date: DateTime.now(),
    open: 1780.36,
    high: 1873.93,
    low: 1755.34,
    close: 1848.56,
    volume: 0,
  );

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  Future<List<Candle>> fetchCandles() async {
    final uri = Uri.parse(
        "https://api.binance.com/api/v3/klines?symbol=BTCUSDT&interval=1h");
    final res = await http.get(uri);
    return (jsonDecode(res.body) as List<dynamic>)
        .map((e) => Candle.fromJson(e))
        .toList()
        .reversed
        .toList();
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
    tabController = TabController(length: 7, vsync: this);
    fetchCandles().then((value) {
      setState(() {
        candles = value;
      });
    });
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
              'TROV/NGN',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            TabBar(
              isScrollable: true,
              controller: tabController,
              labelColor: notifier.getbluewhitecolor,
              indicatorColor: notifier.getbluewhitecolor,
              labelStyle: TextStyle(
                fontSize: 12.sp,
                fontFamily: fontsemibold,
              ),
              tabs: [
                Tab(
                  height: 20,
                  text: 'Buy',
                ),
                Tab(
                  height: 20,
                  text: 'Sell',
                ),
                Tab(
                  height: 20,
                  text: 'Chart',
                ),
                Tab(
                  height: 20,
                  text: 'Order Book',
                ),
                Tab(
                  height: 20,
                  text: 'Last Trades',
                ),
                Tab(
                  height: 20,
                  text: 'Trades',
                ),
                Tab(
                  height: 20,
                  text: 'Orders',
                ),
              ],
            ),
            SizedBox(height: height / 70),
            Container(
              height: height / 1.18,
              child: TabBarView(
                controller: tabController,
                children: [
                  buyAndSell(true),
                  buyAndSell(false),
                  chart(),
                  orderBook(),
                  lastTrades(),
                  trades(),
                  orders(),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget orderBook() {
    return SingleChildScrollView(
      child: Column(
        children: [
          Table(
            children: getOrderBookTableRows(
              Colors.red[400]!,
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Divider(),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Row(
              children: [
                Text(
                  'Spread',
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                SizedBox(
                  width: width / 5,
                ),
                Text(
                  '956,800.00 NGN',
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: Colors.red[400],
                  ),
                ),
              ],
            ),
          ),
          Divider(),
          SizedBox(
            height: height / 70,
          ),
          Table(
            children: getOrderBookTableRows(notifier.getgreencolor,
                showHeaderRow: false),
          ),
        ],
      ),
    );
  }

  Widget buyAndSell(bool isBuy) {
    return Column(
      children: [
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Price Per TROV',
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: CustomTextFormField.textField(
                '400',
                notifier.getbluecolor,
                null,
                notifier.getgrey,
                null,
                notifier.getblck,
                notifier.getgrey,
                70.sp,
                300.sp,
                // controller: referrerController,
                // validator: validateReferrer,
                onSaved: (value) {},
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Quantity of TROV',
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: CustomTextFormField.textField(
                '',
                notifier.getbluecolor,
                null,
                notifier.getgrey,
                null,
                notifier.getblck,
                notifier.getgrey,
                70.sp,
                300.sp,
                // controller: referrerController,
                // validator: validateReferrer,
                onSaved: (value) {},
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
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
                  SizedBox(
                    height: height / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Trading Fee',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Text(
                          '50%',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    height: height / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Total',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Text(
                          '36,180,000.00 NGN',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                ],
              ),
            ),
          ),
        ),
        SizedBox(
          height: height / 30,
        ),
        Button(
          '${isBuy ? 'Buy' : 'Sell'} TROV',
          isBuy ? notifier.getbluecolor : Colors.red[400],
          wihitecolor,
          onTap: () {
            // Navigator.of(context).pop();
          },
        ),
      ],
    );
  }

  Widget chart() {
    return Column(
      children: [
        Container(
          color: notifier.getaddsubwalletgrey,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10.0, vertical: 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Volume 467k',
                  style: TextStyle(
                    fontSize: 12.sp,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                Text(
                  'High 25,000.00',
                  style: TextStyle(
                    fontSize: 12.sp,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                Text(
                  'Low 20,250.00',
                  style: TextStyle(
                    fontSize: 12.sp,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ],
            ),
          ),
        ),
        SizedBox(height: height / 70),
        Container(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'May 26',
                  style: TextStyle(
                    fontSize: 12.sp,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
                Container(
                  width: width / 4,
                  child: dropdown(
                    (value) {},
                    getOptions,
                    null,
                    'Hour',
                    context,
                    null,
                  ),
                ),
              ],
            ),
          ),
        ),
        SizedBox(height: height / 50),
        Container(
          height: height / 1.4,
          child: Candlesticks(
            candles: candles,
          ),
        ),
      ],
    );
  }

  List<TableRow> getOrderBookTableRows(Color color,
      {bool showHeaderRow = true}) {
    List<TableRow> tableRows = [];

    if (showHeaderRow) {
      tableRows.add(TableRow(children: [
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
          child: Text(
            'Price NGN',
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
            'Amount TROV',
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
            'Total NGN',
            textAlign: TextAlign.center,
            style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getsplashgrey,
                fontSize: 13),
          ),
        ),
      ]));
    }

    for (var i = 0; i <= 10; i++) {
      tableRows.add(
        TableRow(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                formatHistoryNumber(Random().nextDouble() * 256, 99000000000),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: color,
                  fontSize: 13,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                Random().nextInt(1000000).toDouble().toStringAsFixed(2),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                Random().nextInt(1000000).toDouble().toStringAsFixed(2),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
          ],
        ),
      );
    }
    return tableRows;
  }

  Widget lastTrades() {
    return SingleChildScrollView(
      child: Column(
        children: [
          Table(
            children: getLastTradesRows(
              Colors.red[400]!,
            ),
          ),
        ],
      ),
    );
  }

  List<TableRow> getLastTradesRows(Color color, {bool showHeaderRow = true}) {
    List<TableRow> tableRows = [];

    if (showHeaderRow) {
      tableRows.add(TableRow(children: [
        Padding(
          padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
          child: Text(
            'Time',
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
            'Price NGN',
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
            'Amount TROV',
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
            'Total TROV',
            textAlign: TextAlign.center,
            style: TextStyle(
                fontFamily: fontsemibold,
                color: notifier.getsplashgrey,
                fontSize: 13),
          ),
        ),
      ]));
    }

    for (var i = 0; i <= 10; i++) {
      tableRows.add(
        TableRow(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                DateFormat.Hms().format(
                    DateTime.now().add(Duration(hours: Random().nextInt(100)))),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: color,
                  fontSize: 13,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                Random().nextInt(1000000).toDouble().toStringAsFixed(2),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                Random().nextInt(1000000).toDouble().toStringAsFixed(2),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                Random().nextInt(1000000).toDouble().toStringAsFixed(2),
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                  fontSize: 13,
                ),
              ),
            ),
          ],
        ),
      );
    }
    return tableRows;
  }

  Widget trades() {
    return Column(
      children: [
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Mon, Feb 20',
                style: TextStyle(
                  fontSize: 13,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.02962',
          'TROV X 21,004.22 NGN',
          '622.43 NGN',
          '14:08',
          true,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.02962',
          'TROV X 21,004.22 NGN',
          '622.43 NGN',
          '14:08',
          false,
        ),
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Tue, Feb 21',
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.02962',
          'TROV X 21,004.22 NGN',
          '622.43 NGN',
          '14:08',
          true,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.02962',
          'TROV X 21,004.22 NGN',
          '622.43 NGN',
          '14:08',
          false,
        ),
        SizedBox(
          height: height / 30,
        ),
      ],
    );
  }

  Widget orders() {
    return Column(
      children: [
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Thu, Feb 23',
                style: TextStyle(
                  fontSize: 13,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.02962',
          'TROV X 12.45 NGN',
          '17,300 NGN',
          '14:08',
          true,
          showMore: true,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.00063',
          'TROV X 10.08 NGN',
          '16,000 NGN',
          '14:08',
          false,
          showMore: true,
        ),
        SizedBox(
          height: height / 50,
        ),
        Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                'Tue, Feb 21',
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontsemibold,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
          ],
        ),
        SizedBox(
          height: height / 50,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.02962',
          'TROV X 12.45 NGN',
          '17,300 NGN',
          '14:08',
          true,
          showMore: true,
        ),
        tradeItemCard(
          'TROV/NGN',
          '0.00063',
          'TROV X 10.08 NGN',
          '16,000 NGN',
          '14:08',
          false,
          showMore: true,
        ),
        SizedBox(
          height: height / 30,
        ),
      ],
    );
  }

  Widget tradeItemCard(
    String asset,
    String rateAmountText,
    String rateAmountText2,
    String totalAmount,
    String time,
    bool isBuy, {
    bool showMore = false,
  }) {
    return Padding(
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
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      asset,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Row(
                      children: [
                        pill(
                          isBuy ? 'Buy' : 'Sell',
                          backColor:
                              isBuy ? notifier.getgreencolor : Colors.red[400],
                          foreColor: wihitecolor,
                        ),
                        if (showMore)
                          IconButton(
                            onPressed: () {},
                            icon: Icon(
                              Icons.more_vert,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                      ],
                    ),
                  ],
                ),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      rateAmountText,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: isBuy ? notifier.getgreencolor : Colors.red[400],
                      ),
                    ),
                    SizedBox(
                      width: 5,
                    ),
                    Text(
                      rateAmountText2,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
              ),
              SizedBox(
                height: height / 70,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      totalAmount,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    Text(
                      time,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
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

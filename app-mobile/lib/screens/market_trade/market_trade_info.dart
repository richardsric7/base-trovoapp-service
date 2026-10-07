import 'dart:async';

import 'package:candlesticks/candlesticks.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
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
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/market_requests.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/utils/local_auth.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';

import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

// MarketTradeInfo is one pair of the DEX market: placing buy and sell
// offers, the price chart, the order book, the last trades on the market
// and the wallet's own trades and open orders. The pair comes from
// appState.viewData[MarketTradeInfoViewPageConfig.key].
class MarketTradeInfo extends StatefulWidget {
  const MarketTradeInfo({Key? key}) : super(key: key);

  @override
  State<MarketTradeInfo> createState() => _MarketTradeInfoState();
}

class _MarketTradeInfoState extends State<MarketTradeInfo>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late TabController tabController;
  late MarketApi api;
  MarketPair? pair;

  List<Candle> candles = [];
  bool candlesLoading = true;
  String resolution = 'Hour';
  MarketOrderBook? book;
  List<MarketTrade>? marketTrades;
  List<MarketOffer>? offers;
  Timer? refresher;

  final buyPrice = TextEditingController();
  final buyQuantity = TextEditingController();
  final sellPrice = TextEditingController();
  final sellQuantity = TextEditingController();

  List<DropdownMenuItem<String>> get getOptions => [
        for (final value in marketResolutions.keys)
          DropdownMenuItem(
            child: Text(value, overflow: TextOverflow.ellipsis),
            value: value,
          ),
      ];

  String get base => pair?.base.code ?? '';
  String get counter => pair?.counter.code ?? '';
  bool get signedIn =>
      appState.secretKeys.isNotEmpty && appState.userInfo?.wallets != null;

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
    tabController = TabController(length: 7, vsync: this);
    appState = Provider.of<DataProvider>(context, listen: false);
    api = MarketApi(appState);
    final data = appState.viewData?[MarketTradeInfoViewPageConfig.key];
    if (data is Map) {
      pair = MarketPair(data);
      final p = pair!;
      final last = double.tryParse(p.lastPrice) ?? 0;
      final ask = double.tryParse(p.bestAsk) ?? 0;
      final bid = double.tryParse(p.bestBid) ?? 0;
      buyPrice.text = _plain(ask > 0 ? ask : last);
      sellPrice.text = _plain(bid > 0 ? bid : last);
      for (final c in [buyPrice, buyQuantity, sellPrice, sellQuantity]) {
        c.addListener(() => setState(() {}));
      }
      loadCandles();
      refresh();
      refresher = Timer.periodic(const Duration(seconds: 20), (_) => refresh());
    }
  }

  @override
  void dispose() {
    refresher?.cancel();
    tabController.dispose();
    for (final c in [buyPrice, buyQuantity, sellPrice, sellQuantity]) {
      c.dispose();
    }
    super.dispose();
  }

  static String _plain(double v) =>
      v <= 0 ? '' : double.parse(v.toStringAsPrecision(8)).toString();

  Future<void> loadCandles() async {
    final p = pair;
    if (p == null) return;
    setState(() => candlesLoading = true);
    List<Candle>? list;
    try {
      list = await api.candles(p, marketResolutions[resolution]!);
    } catch (_) {}
    if (!mounted) return;
    setState(() {
      candles = list ?? [];
      candlesLoading = false;
    });
  }

  Future<T?> _quiet<T>(Future<T?> Function() f) async {
    try {
      return await f();
    } catch (_) {
      return null;
    }
  }

  // refresh reloads the pair's figures, the order book, the market's
  // trades and the wallet's offers.
  Future<void> refresh() async {
    final p = pair;
    if (p == null) return;
    final pairs = _quiet(api.pairs);
    final b = _quiet(() => api.orderBook(p));
    final t = _quiet(() => api.trades(p));
    final mine = signedIn ? _quiet(api.myOffers) : Future.value(null);
    final freshPairs = await pairs;
    final freshBook = await b;
    final freshTrades = await t;
    final myOffers = await mine;
    if (!mounted) return;
    setState(() {
      final match = freshPairs?.where(
          (x) => MarketFavorites.id(x) == MarketFavorites.id(p));
      if (match != null && match.isNotEmpty) pair = match.first;
      book = freshBook ?? book;
      marketTrades = freshTrades ?? marketTrades;
      if (myOffers != null) {
        offers = myOffers.where((o) => o.forPair(p)).toList();
      }
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
              pair?.symbol ?? 'DEX Trade',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            if (pair == null)
              message('Choose a pair from the market first.')
            else ...[
              SizedBox(height: height / 50),
              TabBar(
                isScrollable: true,
                controller: tabController,
                labelColor: notifier.getbluewhitecolor,
                indicatorColor: notifier.getbluewhitecolor,
                labelStyle:
                    TextStyle(fontSize: 12.sp, fontFamily: fontsemibold),
                tabs: const [
                  Tab(height: 20, text: 'Buy'),
                  Tab(height: 20, text: 'Sell'),
                  Tab(height: 20, text: 'Chart'),
                  Tab(height: 20, text: 'Order Book'),
                  Tab(height: 20, text: 'Last Trades'),
                  Tab(height: 20, text: 'Trades'),
                  Tab(height: 20, text: 'Orders'),
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
          ],
        ),
      ),
    );
  }

  Widget message(String text) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 20, vertical: height / 10),
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

  Widget label(String text) {
    return Row(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Text(
            text,
            style: TextStyle(
              fontSize: 15,
              fontFamily: fontsemibold,
              color: notifier.getbluewhitecolor,
            ),
          ),
        ),
      ],
    );
  }

  Widget numberField(TextEditingController controller, String hint) {
    return Row(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: CustomTextFormField.textField(
            hint,
            notifier.getbluecolor,
            null,
            notifier.getgrey,
            null,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            controller: controller,
            keyboardtype: const TextInputType.numberWithOptions(decimal: true),
            inputFormatters: [
              FilteringTextInputFormatter.allow(RegExp(r'[0-9.]')),
            ],
          ),
        ),
      ],
    );
  }

  double _value(TextEditingController c) =>
      double.tryParse(c.text.replaceAll(',', '').trim()) ?? 0;

  // _balance is what the primary wallet holds of a token.
  double _balance(MarketToken t) {
    if (!signedIn) return 0;
    final assets = appState.primaryWallet.claimedAssets ?? [];
    for (final a in assets) {
      if ((a.contractAddress ?? '').toLowerCase() ==
          t.contractAddress.toLowerCase()) {
        return a.amount ?? 0;
      }
    }
    return 0;
  }

  Widget buyAndSell(bool isBuy) {
    final p = pair!;
    final priceC = isBuy ? buyPrice : sellPrice;
    final qtyC = isBuy ? buyQuantity : sellQuantity;
    final total = _value(priceC) * _value(qtyC);
    final available = isBuy
        ? '${formatMarketNumber(_balance(p.counter))} $counter'
        : '${formatMarketNumber(_balance(p.base))} $base';
    return SingleChildScrollView(
      child: Column(
        children: [
          SizedBox(height: height / 50),
          label('Price per $base ($counter)'),
          SizedBox(height: height / 50),
          numberField(priceC, '0.00'),
          SizedBox(height: height / 50),
          label('Quantity of $base'),
          SizedBox(height: height / 50),
          numberField(qtyC, '0'),
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
              child: Column(
                children: [
                  SizedBox(height: height / 70),
                  summaryRow('Available', available),
                  SizedBox(height: height / 70),
                  summaryRow('Trading Fee', 'None'),
                  SizedBox(height: height / 70),
                  summaryRow(
                    isBuy ? 'You pay' : 'You receive',
                    '${formatMarketNumber(total)} $counter',
                    bold: true,
                  ),
                  SizedBox(height: height / 70),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 8.0),
                    child: Text(
                      'Your ${isBuy ? counter : base} is held on the market until the order fills or you cancel it. A small network fee may apply.',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ),
          ),
          SizedBox(height: height / 30),
          Button(
            '${isBuy ? 'Buy' : 'Sell'} $base',
            isBuy ? notifier.getbluecolor : Colors.red[400],
            wihitecolor,
            onTap: () => place(isBuy),
          ),
          SizedBox(height: height / 30),
        ],
      ),
    );
  }

  Widget summaryRow(String name, String value, {bool bold = false}) {
    final style = TextStyle(
      fontSize: 15,
      fontFamily: bold ? fontsemibold : fontbody,
      color: notifier.getbluewhitecolor,
    );
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 8.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(name, style: style),
          Flexible(
            child: Text(value, textAlign: TextAlign.end, style: style),
          ),
        ],
      ),
    );
  }

  String? _problem(bool isBuy) {
    final p = pair!;
    final price = _value(isBuy ? buyPrice : sellPrice);
    final qty = _value(isBuy ? buyQuantity : sellQuantity);
    if (!signedIn) return 'Sign in to your wallet to trade.';
    if (price <= 0) return 'Enter the price per $base.';
    if (qty <= 0) return 'Enter how much $base to ${isBuy ? 'buy' : 'sell'}.';
    if (isBuy && price * qty > _balance(p.counter)) {
      return 'Not enough $counter in your wallet for this order.';
    }
    if (!isBuy && qty > _balance(p.base)) {
      return 'Not enough $base in your wallet for this order.';
    }
    return null;
  }

  static String _tx(Map r) =>
      r['data'] is Map ? '${r['data']['transaction'] ?? ''}' : '';

  static List<String> _messages(Map r) => r['data'] is Map
      ? ((r['data']['messages'] ?? []) as List).map((m) => '$m').toList()
      : [];

  // place builds the offer on app-backend, shows it for confirmation, then
  // signs and submits it (or sends it for approval on a shared wallet).
  Future<void> place(bool isBuy) async {
    final problem = _problem(isBuy);
    if (problem != null) {
      popup(context, title: 'Check your order', message: problem);
      return;
    }
    final p = pair!;
    final price =
        (isBuy ? buyPrice : sellPrice).text.replaceAll(',', '').trim();
    final qty =
        (isBuy ? buyQuantity : sellQuantity).text.replaceAll(',', '').trim();
    final Wallet wallet = appState.primaryWallet;
    showLoader(context);
    try {
      final built =
          await api.placeOffer(p, buy: isBuy, price: price, quantity: qty);
      hideLoader(context);
      final tx = _tx(built);
      if (built['statusCode'] != 200 || tx.isEmpty) {
        popup(context,
            title: 'Order not placed',
            message: MarketApi.message(built, 'Please try again.'));
        return;
      }
      final total =
          (double.tryParse(price) ?? 0) * (double.tryParse(qty) ?? 0);
      final ok = await confirm(
        '${isBuy ? 'Buy' : 'Sell'} ${formatMarketString(qty)} $base',
        [
          ['Price', '${formatMarketString(price)} $counter'],
          [
            isBuy ? 'You pay' : 'You receive',
            '${formatMarketNumber(total)} $counter'
          ],
          ['Trading fee', 'None'],
          ['Wallet', wallet.alias ?? wallet.address ?? ''],
        ],
        _messages(built),
      );
      if (!ok) return;
      final shared = wallet.isSharedWalletAndCanInitiate;
      showLoader(context);
      final sent = await api.placeOffer(
        p,
        buy: isBuy,
        price: price,
        quantity: qty,
        transaction: tx,
        signature: shared
            ? ''
            : TrovoWalletSDK().signBase64Txn(appState.secretKeys[0], tx, ''),
        commit: shared,
      );
      hideLoader(context);
      if (sent['statusCode'] == 200) {
        (isBuy ? buyQuantity : sellQuantity).clear();
        popup(
          context,
          title: shared ? 'Sent for approval' : 'Order placed',
          bodyColor: notifier.getgreencolor,
          message: shared
              ? 'Your order goes on the market once the wallet\'s approvers approve it.'
              : 'Your order is on its way to the market. It shows under Orders in a few seconds.',
        );
        tabController.animateTo(6);
        refresh();
      } else {
        popup(context,
            title: 'Order not placed',
            message: MarketApi.message(sent, 'Please try again.'));
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: 'Error', message: e.toString());
    }
  }

  // confirm shows a summary and asks for biometrics when the phone has them.
  Future<bool> confirm(
      String title, List<List<String>> rows, List<String> messages) async {
    final ok = await showModalBottomSheet<bool>(
          context: context,
          isScrollControlled: true,
          backgroundColor: notifier.getwihitecolor,
          shape: const RoundedRectangleBorder(
            borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
          ),
          builder: (context) => SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  const SizedBox(height: 16),
                  for (final r in rows) ...[
                    summaryRow(r[0], r[1]),
                    const SizedBox(height: 8),
                  ],
                  for (final m in messages)
                    Padding(
                      padding: const EdgeInsets.only(top: 6),
                      child: Text(
                        m,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  const SizedBox(height: 20),
                  Button(
                    'Confirm',
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: () => Navigator.of(context).pop(true),
                  ),
                  Center(
                    child: TextButton(
                      onPressed: () => Navigator.of(context).pop(false),
                      child: Text(
                        'Cancel',
                        style: TextStyle(color: notifier.getbluewhitecolor),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ) ??
        false;
    if (!ok) return false;
    try {
      final authenticator = Authenticator();
      if (await authenticator.canCheckBiometrics()) {
        return await authenticator.authenticateMe();
      }
    } catch (_) {}
    return true;
  }

  Widget chart() {
    final p = pair!;
    return Column(
      children: [
        Container(
          color: notifier.getaddsubwalletgrey,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10.0, vertical: 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                stat('Volume ${formatMarketString(p.baseVolume24h)}'),
                stat('High ${formatMarketString(p.high24h)}'),
                stat('Low ${formatMarketString(p.low24h)}'),
              ],
            ),
          ),
        ),
        SizedBox(height: height / 70),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 10.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '${formatMarketString(p.lastPrice)} $counter  '
                '${p.isUp ? '+' : '-'}${p.changePercent24h.abs().toStringAsFixed(2)}%',
                style: TextStyle(
                  fontSize: 12.sp,
                  fontFamily: fontsemibold,
                  color: p.isUp ? notifier.getgreencolor : Colors.red[400],
                ),
              ),
              Container(
                width: width / 3,
                child: dropdown(
                  (value) {
                    setState(() => resolution = '$value');
                    loadCandles();
                  },
                  getOptions,
                  resolution,
                  null,
                  context,
                  null,
                ),
              ),
            ],
          ),
        ),
        SizedBox(height: height / 50),
        Container(
          height: height / 1.4,
          child: candlesLoading && candles.isEmpty
              ? const Center(child: CircularProgressIndicator())
              : candles.length < 2
                  ? message(
                      'Not enough trades yet to draw a chart for $base/$counter.')
                  : Candlesticks(candles: candles),
        ),
      ],
    );
  }

  Widget stat(String text) {
    return Flexible(
      child: Text(
        text,
        overflow: TextOverflow.ellipsis,
        style: TextStyle(
          fontSize: 12.sp,
          fontFamily: fontbody,
          color: notifier.getbluewhitecolor,
        ),
      ),
    );
  }

  Widget orderBook() {
    final b = book;
    if (b == null) return const Center(child: CircularProgressIndicator());
    if (b.asks.isEmpty && b.bids.isEmpty) {
      return message('No orders on the book for $base/$counter yet.');
    }
    final spread = b.spread;
    return RefreshIndicator(
      onRefresh: refresh,
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: Column(
          children: [
            // asks highest first, so the best ask sits next to the spread
            Table(
              children: getOrderBookTableRows(
                b.asks.reversed.toList(),
                Colors.red[400]!,
              ),
            ),
            SizedBox(height: height / 70),
            const Divider(),
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
                  SizedBox(width: width / 5),
                  Text(
                    spread == null
                        ? '-'
                        : '${formatMarketNumber(spread)} $counter',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: Colors.red[400],
                    ),
                  ),
                ],
              ),
            ),
            const Divider(),
            SizedBox(height: height / 70),
            Table(
              children: getOrderBookTableRows(
                b.bids,
                notifier.getgreencolor,
                showHeaderRow: false,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget headerCell(String text) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
      child: Text(
        text,
        textAlign: TextAlign.center,
        style: TextStyle(
          fontFamily: fontsemibold,
          color: notifier.getsplashgrey,
          fontSize: 13,
        ),
      ),
    );
  }

  Widget cell(String text, {Color? color}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
      child: Text(
        text,
        textAlign: TextAlign.center,
        style: TextStyle(
          fontFamily: fontbody,
          color: color ?? notifier.getbluewhitecolor,
          fontSize: 13,
        ),
      ),
    );
  }

  List<TableRow> getOrderBookTableRows(
    List<MarketLevel> levels,
    Color color, {
    bool showHeaderRow = true,
  }) {
    return [
      if (showHeaderRow)
        TableRow(
          children: [
            headerCell('Price $counter'),
            headerCell('Amount $base'),
            headerCell('Total $counter'),
          ],
        ),
      for (final l in levels)
        TableRow(
          children: [
            cell(formatMarketNumber(l.price), color: color),
            cell(formatMarketNumber(l.amount)),
            cell(formatMarketNumber(l.total)),
          ],
        ),
    ];
  }

  Widget lastTrades() {
    final t = marketTrades;
    if (t == null) return const Center(child: CircularProgressIndicator());
    if (t.isEmpty) return message('No trades on $base/$counter yet.');
    return RefreshIndicator(
      onRefresh: refresh,
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: Table(children: getLastTradesRows(t)),
      ),
    );
  }

  List<TableRow> getLastTradesRows(List<MarketTrade> list) {
    return [
      TableRow(
        children: [
          headerCell('Time'),
          headerCell('Price $counter'),
          headerCell('Amount $base'),
          headerCell('Total $counter'),
        ],
      ),
      for (var i = 0; i < list.length; i++)
        TableRow(
          children: [
            cell(DateFormat.Hms().format(list[i].at.toLocal())),
            // red when the price fell from the trade before it
            cell(
              formatMarketNumber(list[i].price),
              color: i + 1 < list.length && list[i].price < list[i + 1].price
                  ? Colors.red[400]
                  : notifier.getgreencolor,
            ),
            cell(formatMarketNumber(list[i].baseAmount)),
            cell(formatMarketNumber(list[i].counterAmount)),
          ],
        ),
    ];
  }

  Widget signInFirst() =>
      message('Sign in to your wallet to see your orders.');

  // trades are the wallet's offers on this pair that have traded.
  Widget trades() {
    if (!signedIn) return signInFirst();
    final o = offers;
    if (o == null) return const Center(child: CircularProgressIndicator());
    final traded = o.where((x) => x.fills > 0).toList();
    if (traded.isEmpty) {
      return message('You have no trades on $base/$counter yet.');
    }
    return RefreshIndicator(
      onRefresh: refresh,
      child: ListView(
        children: [
          SizedBox(height: height / 50),
          for (final x in traded)
            tradeItemCard(
              pair!.symbol,
              formatMarketNumber(x.averagePrice),
              '$counter x ${formatMarketNumber(x.filledBase)} $base',
              '${formatMarketNumber(x.filledCounter)} $counter',
              x.status,
              x.isBuy,
            ),
          SizedBox(height: height / 30),
        ],
      ),
    );
  }

  // orders are the wallet's offers on this pair still on the market or on
  // their way to it.
  Widget orders() {
    if (!signedIn) return signInFirst();
    final o = offers;
    if (o == null) return const Center(child: CircularProgressIndicator());
    final live = o
        .where(
            (x) => const ['Open', 'Part filled', 'Pending'].contains(x.status))
        .toList();
    if (live.isEmpty) {
      return message('You have no open orders on $base/$counter.');
    }
    return RefreshIndicator(
      onRefresh: refresh,
      child: ListView(
        children: [
          SizedBox(height: height / 50),
          for (final x in live)
            tradeItemCard(
              pair!.symbol,
              formatMarketNumber(x.pricePerUnit),
              '$counter x ${formatMarketNumber(x.quantity)} $base',
              '${formatMarketNumber(x.pricePerUnit * x.quantity)} $counter',
              x.status,
              x.isBuy,
              onCancel: x.cancellable ? () => cancel(x) : null,
            ),
          SizedBox(height: height / 30),
        ],
      ),
    );
  }

  Future<void> cancel(MarketOffer offer) async {
    final Wallet wallet = appState.primaryWallet;
    showLoader(context);
    try {
      final built = await api.cancelOffer(offer.id);
      hideLoader(context);
      final tx = _tx(built);
      if (built['statusCode'] != 200 || tx.isEmpty) {
        popup(context,
            title: 'Not cancelled',
            message: MarketApi.message(built, 'Please try again.'));
        return;
      }
      final ok = await confirm(
        'Cancel this order?',
        [
          [
            'Order',
            '${offer.isBuy ? 'Buy' : 'Sell'} ${formatMarketNumber(offer.quantity)} $base'
          ],
          ['Price', '${formatMarketNumber(offer.pricePerUnit)} $counter'],
          [
            'Returned to wallet',
            '${formatMarketNumber(offer.remaining)} ${offer.isBuy ? counter : base}'
          ],
        ],
        _messages(built),
      );
      if (!ok) return;
      final shared = wallet.isSharedWalletAndCanInitiate;
      showLoader(context);
      final sent = await api.cancelOffer(
        offer.id,
        transaction: tx,
        signature: shared
            ? ''
            : TrovoWalletSDK().signBase64Txn(appState.secretKeys[0], tx, ''),
        commit: shared,
      );
      hideLoader(context);
      if (sent['statusCode'] == 200) {
        popup(
          context,
          title: shared ? 'Sent for approval' : 'Order cancelled',
          bodyColor: notifier.getgreencolor,
          message: shared
              ? 'The order is cancelled once the wallet\'s approvers approve it.'
              : 'What was left of the order goes back to your wallet.',
        );
        refresh();
      } else {
        popup(context,
            title: 'Not cancelled',
            message: MarketApi.message(sent, 'Please try again.'));
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: 'Error', message: e.toString());
    }
  }

  Widget tradeItemCard(
    String asset,
    String rateAmountText,
    String rateAmountText2,
    String totalAmount,
    String status,
    bool isBuy, {
    VoidCallback? onCancel,
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
        child: Column(
          children: [
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    asset,
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
                            isBuy ? notifier.getgreencolor : Colors.red[400]!,
                        foreColor: wihitecolor,
                      ),
                      if (onCancel != null)
                        PopupMenuButton<String>(
                          icon: Icon(
                            Icons.more_vert,
                            color: notifier.getbluewhitecolor,
                          ),
                          onSelected: (_) => onCancel(),
                          itemBuilder: (_) => const [
                            PopupMenuItem(
                              value: 'cancel',
                              child: Text('Cancel order'),
                            ),
                          ],
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
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: isBuy ? notifier.getgreencolor : Colors.red[400],
                    ),
                  ),
                  const SizedBox(width: 5),
                  Flexible(
                    child: Text(
                      rateAmountText2,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: height / 70),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    totalAmount,
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  Text(
                    status,
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(height: height / 50),
          ],
        ),
      ),
    );
  }
}

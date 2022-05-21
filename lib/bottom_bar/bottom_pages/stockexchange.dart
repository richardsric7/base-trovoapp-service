import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stock_exchange_tabs/buy.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';

import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'stock_exchange_tabs/sell.dart';

class StockExchange extends StatefulWidget {
  const StockExchange({Key? key}) : super(key: key);

  @override
  State<StockExchange> createState() => _StockExchangeState();
}

class _StockExchangeState extends State<StockExchange>
    with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return ScreenUtilInit(
      builder: () => DefaultTabController(
        length: 2,
        child: Scaffold(
          backgroundColor: notifier.getwihitecolor,
          appBar: AppBar(
            centerTitle: true,
            leading: GestureDetector(
              onTap: () {
                Get.back();
              },
              child: Image.asset("assets/images/back.png", scale: 6),
            ),
            title: Text(
              LanguageEn.exchange,
              style:
                  TextStyle(color: notifier.getblck, fontFamily: 'Gilroy_Bold'),
            ),
            backgroundColor: notifier.getfavorites,
            elevation: 0,
            bottom: PreferredSize(
              preferredSize: Size.fromHeight(80.sp),
              child: Padding(
                padding: const EdgeInsets.only(left: 13, right: 13),
                child: Card(elevation: 0,
                  color: notifier.getwihitecolor,
                  child: Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: Container(
                      color: notifier.getfavorites,
                      padding: EdgeInsets.all(8.sp),
                      child: TabBar(
                          padding: EdgeInsets.zero,
                          unselectedLabelColor:
                             notifier.getgrey,
                          indicator: BoxDecoration(
                              borderRadius: BorderRadius.circular(5.sp),
                              color: notifier.getbluecolor),
                          tabs: [
                            Tab(
                              child: Align(
                                alignment: Alignment.center,
                                child: Text(
                                  LanguageEn.buy,
                                ),
                              ),
                            ),
                            Tab(
                              child: Align(
                                alignment: Alignment.center,
                                child: Text(
                                  LanguageEn.sell,
                              ),
                            )),
                          ]),
                    ),
                  ),
                ),
              ),
            ),
          ),
          body: const TabBarView(

              children: [Buy(), Sell()]),
        ),
      ),
    );
  }
}

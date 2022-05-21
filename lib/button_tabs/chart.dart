import 'package:buttons_tabbar/buttons_tabbar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stock_exchange_tabs/selectcrypto.dart';
import 'package:gocrypto/graph_tabs/fiveyear.dart';
import 'package:gocrypto/graph_tabs/oned.dart';
import 'package:gocrypto/graph_tabs/onemonth.dart';
import 'package:gocrypto/graph_tabs/oneweek.dart';
import 'package:gocrypto/graph_tabs/oneyear.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../utils/medeiaqury/medeiaqury.dart';

class Chart extends StatefulWidget {
  const Chart({Key? key}) : super(key: key);

  @override
  State<Chart> createState() => _ChartState();
}

class _ChartState extends State<Chart> {
  late ColorNotifier notifier;
  int selectedindex = -1;
  int touchedIndex = -1;
  List chartduration = ["1D", "1W", "1Y", "5Y"];

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
    return ScreenUtilInit(
      builder: () => Scaffold(
        floatingActionButton: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(width: width / 13),
            GestureDetector(
              onTap: () {
                Get.to(() => const SelectCrypto());
              },
              child: button(LanguageEn.sell, const Color(0xffF05150),
                  notifier.getwihitecolor, Colors.transparent),
            ),
            SizedBox(width: width / 30),
            GestureDetector(
              onTap: () {
                Get.to(() => const SelectCrypto());
              },
              child: button(LanguageEn.buy, const Color(0xff2873FF),
                  notifier.getwihitecolor, Colors.transparent),
            ),
          ],
        ),
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            notifier.getwihitecolor, LanguageEn.ltcusd, notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Row(
                children: [
                  SizedBox(width: width / 13),
                  ltc(),
                ],
              ),
              SizedBox(height: height / 40),
              Container(
                color: Colors.transparent,
                height: width / 1.5,
                width: 320.w,
                child: DefaultTabController(
                  length: 5,
                  child: Column(
                    children: <Widget>[
                      ButtonsTabBar(
                        backgroundColor:
                            const Color.fromARGB(106, 149, 180, 253),
                        unselectedBackgroundColor: const Color(0xffeff6ff),
                        unselectedLabelStyle:
                            const TextStyle(color: Colors.black),
                        labelStyle: const TextStyle(
                            color: Colors.green, fontFamily: 'Gilroy_Bold'),
                        tabs: [
                          Tab(
                            child: Container(
                              height: height / 19,
                              width: width / 8,
                              decoration: const BoxDecoration(
                                // border: Border.all(color: const Color(0xff8f94b0)),
                                borderRadius: BorderRadius.all(
                                  Radius.circular(5),
                                ),
                              ),
                              child: Center(
                                child: Text(
                                  "1D",
                                  style: TextStyle(
                                      fontFamily: 'Gilroy_Bold',
                                      fontSize: 14.sp),
                                ),
                              ),
                            ),
                          ),
                          Tab(
                            child: Container(
                              height: height / 19,
                              width: width / 8,
                              decoration: const BoxDecoration(
                                borderRadius: BorderRadius.all(
                                  Radius.circular(5),
                                ),
                              ),
                              child: Center(
                                child: Text(
                                  "1W",
                                  style: TextStyle(
                                      fontFamily: 'Gilroy_Bold',
                                      fontSize: 14.sp),
                                ),
                              ),
                            ),
                          ),
                          Tab(
                            child: Container(
                              height: height / 19,
                              width: width / 8,
                              decoration: const BoxDecoration(
                                borderRadius: BorderRadius.all(
                                  Radius.circular(5),
                                ),
                              ),
                              child: Center(
                                child: Text(
                                  "1M",
                                  style: TextStyle(
                                      fontFamily: 'Gilroy_Bold',
                                      fontSize: 14.sp),
                                ),
                              ),
                            ),
                          ),
                          Tab(
                            child: Container(
                              height: height / 19,
                              width: width / 8,
                              decoration: const BoxDecoration(
                                borderRadius: BorderRadius.all(
                                  Radius.circular(5),
                                ),
                              ),
                              child: Center(
                                child: SingleChildScrollView(
                                  scrollDirection: Axis.horizontal,
                                  child: Text(
                                    "1Y",
                                    style: TextStyle(
                                        fontFamily: 'Gilroy_Bold',
                                        fontSize: 14.sp),
                                  ),
                                ),
                              ),
                            ),
                          ),
                          Tab(
                            child: Container(
                              height: height / 19,
                              width: width / 8,
                              decoration: const BoxDecoration(
                                borderRadius: BorderRadius.all(
                                  Radius.circular(5),
                                ),
                              ),
                              child: Center(
                                child: Text(
                                  "5Y",
                                  style: TextStyle(
                                      fontFamily: 'Gilroy_Bold',
                                      fontSize: 14.sp),
                                ),
                              ),
                            ),
                          ),
                        ],
                      ),
                      const Expanded(
                        child: TabBarView(
                          children: <Widget>[
                            Oneday(),
                            Oneweek(),
                            Onemonth(),
                            OneYear(),
                            FiveYear(),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 20),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.description,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 18.sp),
                  ),
                ],
              ),
              Container(
                color: Colors.transparent,
                height: height / 4.2,
                child: SingleChildScrollView(
                  child: Column(
                    children: [
                      Padding(
                        padding: EdgeInsets.only(
                            left: width / 20, right: width / 30),
                        child: Text(
                          LanguageEn.litecoinisa,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 14.sp,
                              fontFamily: 'Gilroy_Medium'),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.only(
                            left: width / 20, right: width / 30),
                        child: Text(
                          LanguageEn.litecoinisa,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 14.sp,
                              fontFamily: 'Gilroy_Medium'),
                        ),
                      ),
                      Padding(
                        padding: EdgeInsets.only(
                            left: width / 20, right: width / 30),
                        child: Text(
                          LanguageEn.litecoinisa,
                          style: TextStyle(
                              color: notifier.getgrey,
                              fontSize: 14.sp,
                              fontFamily: 'Gilroy_Medium'),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget button(buttontext, colorbutton, buttontextcolor, bordercolor) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(15),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          LayoutBuilder(builder: (context, constraints) {
            return Container(
              height: height / 15,
              width: width / 2.3,
              decoration: BoxDecoration(
                border: Border.all(color: bordercolor),
                color: colorbutton,
                borderRadius: BorderRadius.circular(15),
              ),
              child: Center(
                child: Text(
                  buttontext,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      fontFamily: 'Gilroy_Medium',
                      fontSize: 15.sp,
                      color: buttontextcolor),
                ),
              ),
            );
          }),
        ],
      ),
    );
  }

  Widget ltc() {
    return Row(
      children: [
        Image.asset("assets/images/kotok.png", height: height / 17),
        SizedBox(width: width / 30),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              LanguageEn.litecoinltc,
              style: TextStyle(
                  color: notifier.getgrey,
                  fontSize: 13.sp,
                  fontFamily: 'Gilroy_Medium'),
            ),
            Row(
              children: [
                Text(
                  "\$34,970.98",
                  style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 18.sp,
                      fontFamily: 'Gilroy_Bold'),
                ),
                SizedBox(width: width / 3.8),
                Container(
                  height: height / 30,
                  width: width / 6,
                  decoration: BoxDecoration(
                    color: const Color(0xff22C36B),
                    borderRadius: BorderRadius.all(
                      Radius.circular(30.sp),
                    ),
                  ),
                  child: Center(
                    child: Text(
                      "+10%",
                      style: TextStyle(
                          color: notifier.getwihitecolor, fontSize: 13.sp),
                    ),
                  ),
                )
              ],
            ),
          ],
        ),
      ],
    );
  }
}

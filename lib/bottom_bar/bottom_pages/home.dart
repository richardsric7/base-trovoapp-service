import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_stoc/custtomstoc.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Stock%20Gains%20card/stockgainscard.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_slock_list/custtom_slock_list.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/selectstocks.dart';
import 'package:trovo_wallet/button_tabs/chart.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';

import '../../utils/medeiaqury/medeiaqury.dart';
import 'stock_exchange_tabs/notification.dart';

class Home extends StatefulWidget {
  const Home({Key? key}) : super(key: key);

  @override
  State<Home> createState() => _HomeState();
}

class _HomeState extends State<Home> {
  late ColorNotifier notifier;

  int touchedIndex = -1;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(70.sp),
          // here the desired height
          child: AppBar(
            actions: [
              GestureDetector(
                onTap: () {
                  Get.to(() => const SelectStocks());
                },
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Image.asset("assets/images/searsh.png"),
                ),
              ),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Notifi());
                },
                child: Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Image.asset("assets/images/bell.png"),
                ),
              ),
            ],
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            title: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: height / 70),
                Text(
                  LanguageEn.gocrypto,
                  style: TextStyle(
                      fontSize: 20.sp,
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Bold'),
                ),
                Text(
                  LanguageEn.buildingtrustinthecrypto,
                  style: TextStyle(
                      fontSize: 13.5.sp,
                      color: notifier.getgrey,
                      fontFamily: 'Gilroy-Regular'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget mypoyfoliyo(image, title, stockprice, updown) {
    return Container(
      margin: EdgeInsets.only(left: width / 20),
      height: height / 4.5,
      width: width / 1.8,
      decoration: BoxDecoration(
        color: notifier.getfavorites,
        borderRadius: BorderRadius.all(
          Radius.circular(20.sp),
        ),
      ),
      child: Column(
        children: [
          Row(
            children: [
              Padding(
                padding: EdgeInsets.only(left: width / 20, top: height / 60),
                child: Image.asset(image, height: height / 15),
              ),
              SizedBox(width: width / 20),
              // SizedBox(height: height / 50),
              Column(
                children: [
                  SizedBox(height: height / 50),
                  Text(
                    title,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 18.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
            ],
          ),
          SizedBox(height: height / 15),
          Row(
            children: [
              SizedBox(width: width / 25),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    LanguageEn.portofolio,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontFamily: 'Gilroy_Medium',
                        fontSize: 13.sp),
                  ),
                  Text(
                    stockprice,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 21.sp),
                  ),
                ],
              ),
              const Spacer(),
              Text(
                updown,
                style:
                    TextStyle(color: const Color(0xff22C36B), fontSize: 13.sp),
              ),
              SizedBox(width: width / 25),
            ],
          )
        ],
      ),
    );
  }

  LineChartData mainData(color) {
    List<Color> gradientColors = [color];
    return LineChartData(
      gridData: FlGridData(
        show: false,
        drawVerticalLine: true,
        getDrawingHorizontalLine: (value) {
          return FlLine(strokeWidth: 1, color: Colors.yellow);
        },
        getDrawingVerticalLine: (value) {
          return FlLine(
            strokeWidth: 0,
          );
        },
      ),
      titlesData: FlTitlesData(
        show: false,
        leftTitles: SideTitles(
          showTitles: false,
          reservedSize: 20,
          margin: 8,
        ),
      ),
      borderData: FlBorderData(
          show: true,
          border: Border.all(color: notifier.getbluecolor, width: 0)),
      minX: 0,
      maxX: 8,
      minY: 0,
      maxY: 5,
      lineBarsData: [
        LineChartBarData(
          spots: [
            const FlSpot(0, 2.5),
            const FlSpot(1, 2),
            const FlSpot(2, 4),
            const FlSpot(3, 3.1),
            const FlSpot(4, 4),
            const FlSpot(5, 2),
            const FlSpot(6, 4),
            const FlSpot(7, 3.1),
            const FlSpot(8, 2),
            const FlSpot(9, 1.5),
            const FlSpot(10, 3),
          ],
          isCurved: true,
          colors: gradientColors,
          barWidth: 2,
          isStrokeCapRound: true,
          dotData: FlDotData(
            show: false,
          ),
          belowBarData: BarAreaData(
            show: true,
            colors:
                gradientColors.map((color) => color.withOpacity(0.1)).toList(),
          ),
        ),
      ],
    );
  }
}

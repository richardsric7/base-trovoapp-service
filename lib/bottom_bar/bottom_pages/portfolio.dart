import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/button_tabs/chart.dart';
import 'package:provider/provider.dart';
import 'package:syncfusion_flutter_charts/charts.dart';

import '../../Custom_BlocObserver/custtom_slock_list/custtom_slock_list.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class Portfolio extends StatefulWidget {
  const Portfolio({Key? key}) : super(key: key);

  @override
  State<Portfolio> createState() => _PortfolioState();
}

class _PortfolioState extends State<Portfolio> {
  late ColorNotifier notifier;

  // late List<GDPData> _chartData;
  // late TooltipBehavior _tooltipBehavior;

  @override
  void initState() {
    // _chartData = getChartData();
    // _tooltipBehavior = TooltipBehavior(enable: true);
    super.initState();
  }

  List<ChartData> chartData = [
    const ChartData('IND', 24),
    const ChartData('AUS', 20),
    const ChartData('USA', 27),
    const ChartData('DEU', 57),
    const ChartData('ITA', 30),
    const ChartData('UK', 41),
  ];

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "Portfolio", notifier.getblck,
            height: height / 15),
        body: Column(
          children: [
            SizedBox(height: height / 40),
            Row(
              children: [
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Chart());
                  },
                  child: updownrate("assets/images/Today_Gains.png", "\$2,209",
                      "Today Gains"),
                ),
                const Spacer(),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Chart());
                  },
                  child: updownrate("assets/images/Overall Loss.png", "\$5,440",
                      "Overall Loss"),
                ),
                SizedBox(width: width / 15),
              ],
            ),
            SizedBox(height: height / 25),
            Row(
              children: [
                SizedBox(width: width / 20),
                Text(
                  "Portfolio Balance",
                  style: TextStyle(
                      fontSize: 14.sp,
                      color: notifier.getgrey,
                      fontFamily: 'Gilroy_Medium'),
                ),
                const Spacer(),
                Image.asset("assets/images/ball.png", height: height / 20),
                SizedBox(width: width / 15),
              ],
            ),
            SizedBox(height: height / 200),
            Row(
              children: [
                SizedBox(width: width / 20),
                Text(
                  "\$97,326.46",
                  style: TextStyle(
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Bold',
                      fontSize: 27.sp),
                ),
              ],
            ),
            SizedBox(height: height / 200),
            Row(
              children: [
                SizedBox(width: width / 25),
                const Icon(
                  Icons.arrow_drop_up_outlined,
                  color: Color(0xff19C09A),
                ),
                Text(
                  "65.63 (76,23%)",
                  style: TextStyle(
                      color: const Color(0xff19C09A),
                      fontSize: 12.sp,
                      fontFamily: 'Gilroy_Bold'),
                ),
                SizedBox(width: width / 100),
                Text(
                  "Today",
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 12.sp,
                      fontFamily: 'Gilroy_Bold'),
                ),
              ],
            ),
            SizedBox(height: height / 50),
            Expanded(
              child: ListView.builder(
                itemCount: 1,
                itemBuilder: (context, index) {
                  return Column(
                    children: [
                      columnchart(),
                      SizedBox(height: height / 80),
                      Row(
                        children: [
                          SizedBox(width: width / 20),
                          Text(
                            "List Stocks",
                            style: TextStyle(
                                color: notifier.getblck,
                                fontSize: 18.sp,
                                fontFamily: 'Gilroy_Bold'),
                          ),
                        ],
                      ),
                      SizedBox(height: height / 30),
                      GestureDetector(
                        onTap: () {
                          Get.to(() => const Chart());
                        },
                        child: const CusttomButton(
                            "assets/images/kotok.png",
                            "Kotak Bank",
                            "Kotak PVL.",
                            "\$326,23",
                            "2,87%",
                            Color(0xffF65556)),
                      ),
                      SizedBox(height: height / 30),
                      GestureDetector(
                        onTap: () {
                          Get.to(() => const Chart());
                        },
                        child: const CusttomButton(
                            "assets/images/hdfc.png",
                            "HDFC bank",
                            "HDFC Inc.",
                            "\$252,12",
                            "10.03%",
                            Color(0xff22C36B)),
                      ),
                      SizedBox(height: height / 30),
                      GestureDetector(
                        onTap: () {
                          Get.to(() => const Chart());
                        },
                        child: const CusttomButton(
                            "assets/images/Ambuja_logo.png",
                            "Ambuja Cement",
                            "Ambuja",
                            "\$360,23",
                            "10.03%",
                            Color(0xff22C36B)),
                      ),
                      SizedBox(height: height / 30),
                      GestureDetector(
                        onTap: () {
                          Get.to(() => const Chart());
                        },
                        child: const CusttomButton(
                            "assets/images/airtel.jpg",
                            "Airtel bharti",
                            "Airtel",
                            "\$326,23",
                            "2,87%",
                            Color(0xffF65556)),
                      ),
                      SizedBox(height: height / 30),
                      GestureDetector(
                        onTap: () {
                          Get.to(() => const Chart());
                        },
                        child: const CusttomButton(
                            "assets/images/icici.png",
                            "icici lombard gen.",
                            "icici",
                            "\$360,23",
                            "10,03%",
                            Color(0xffF65556)),
                      ),
                      SizedBox(height: height / 30),
                    ],
                  );
                },
              ),
            )
          ],
        ),
      ),
    );
  }

  Widget columnchart() {
    return SfCartesianChart(
        primaryXAxis: CategoryAxis(interval: 1),
        annotations: const <CartesianChartAnnotation>[
          CartesianChartAnnotation(
              coordinateUnit: CoordinateUnit.percentage,
              region: AnnotationRegion.plotArea,
              widget: Text('',
                  style: TextStyle(
                    fontSize: 14,
                  )),
              x: '50%',
              y: '50%')
        ],
        series: <ChartSeries<ChartData, String>>[
          ColumnSeries<ChartData, String>(
              color: notifier.getbluecolor,
              dataSource: chartData,
              xValueMapper: (ChartData data, _) => data.x,
              yValueMapper: (ChartData data, _) => data.y)
        ]);
  }

  Widget updownrate(image, rate, days) {
    return Container(
      color: notifier.getwihitecolor,
      height: height / 12,
      width: width / 2.2,
      child: Card(
        color: notifier.getwihitecolor,
        elevation: 1,
        child: Row(
          children: [
            SizedBox(width: width / 20),
            Image.asset(
              image,
              height: height / 17,
            ),
            SizedBox(width: width / 30),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: height / 100),
                Text(
                  rate,
                  style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 17.sp,
                    fontFamily: 'Gilroy_Bold',
                  ),
                ),
                Text(
                  days,
                  style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 13.sp,
                    fontFamily: 'Gilroy_Medium',
                  ),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }
}

class ChartData {
  const ChartData(this.x, this.y);

  final String x;
  final int y;
}

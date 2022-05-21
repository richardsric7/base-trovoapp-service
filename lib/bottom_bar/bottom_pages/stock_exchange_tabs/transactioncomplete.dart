import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/portfolio.dart';
import 'package:gocrypto/bottom_bar/bottombar.dart';
import 'package:provider/provider.dart';

import '../../../Custom_BlocObserver/graph/graph.dart';
import '../../../utils/medeiaqury/medeiaqury.dart';

class TransactionComplete extends StatefulWidget {
  const TransactionComplete({Key? key}) : super(key: key);

  @override
  State<TransactionComplete> createState() => _TransactionCompleteState();
}

class _TransactionCompleteState extends State<TransactionComplete> {
  late ColorNotifier notifier;
  int touchedIndex = -1;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 7),
              Center(
                child: Image.asset("assets/images/TransactionComplete.png",
                    height: height / 5),
              ),
              SizedBox(height: height / 50),
              Text(
                "Transaction Complete",
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 21.sp,
                    fontFamily: 'Gilroy_Bold'),
              ),
              SizedBox(height: height / 100),
              Text(
                "Your transaction has  been completed.\nYou purchased \$127,00 of Apple.",
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: 'Gilroy-Regular'),
              ),
              Padding(
                padding: EdgeInsets.only(left: width / 15, right: width / 3.5),
                child: Container(
                  color: Colors.transparent,
                  height: height / 2.8,
                  width: double.infinity,
                  child: LineChart(
                    mainData(notifier.getbluecolor, 0.1),
                  ),
                ),
              ),
              InkWell(
                onHover: (value) {
                  setState(() {
                    touchedIndex = -1;
                  });
                },
                hoverColor: Colors.transparent,
                onTap: () {},
                child: Row(
                  children: const <Widget>[],
                ),
              ),
              GestureDetector(
                  onTap: () {
                    Get.to(const Portfolio());
                  },
                  child: Button("View Portfolio", notifier.getbluecolor,
                      notifier.getwihitecolor)),
              GestureDetector(
                  onTap: () {
                    Get.to(const BottomHome());
                  },
                  child: Button("Go to Home", notifier.getwihitecolor,
                      notifier.getbluecolor)),
            ],
          ),
        ),
      ),
    );
  }
}

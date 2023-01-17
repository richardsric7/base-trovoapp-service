import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/selectcrypto.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/selectstocks.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';

import '../../../utils/medeiaqury/medeiaqury.dart';

class Buy extends StatefulWidget {
  const Buy({Key? key}) : super(key: key);

  @override
  State<Buy> createState() => _BuyState();
}

class _BuyState extends State<Buy> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getfavorites,
        body: Column(
          children: [
            GestureDetector(
              onTap: () {
                Get.to(
                  const SelectStocks(),
                );
              },
              child: Padding(
                padding: const EdgeInsets.only(left: 13, right: 13),
                child: exchangestock(),
              ),
            ),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.only(left: 13, right: 13),
              child: exchangefree(),
            ),
            SizedBox(height: height / 25),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  LanguageEn.clickherefor,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 11.sp,
                      fontFamily: 'Gilroy_Medium'),
                ),
                Text(
                  LanguageEn.termsandcondition,
                  style: TextStyle(
                      color: notifier.getbluecolor,
                      fontSize: 11.sp,
                      fontFamily: 'Gilroy_Medium'),
                ),
              ],
            ),
            Text(
              LanguageEn.forthistranjection,
              style: TextStyle(
                  color: notifier.getgrey,
                  fontSize: 11.sp,
                  fontFamily: 'Gilroy_Medium'),
            ),
            SizedBox(height: height / 17),
            GestureDetector(
              onTap: () {
                Get.to(
                  const SelectCrypto(),
                  // const SelectStocks(),
                );
              },
              child: Button(LanguageEn.exchangenow, notifier.getbluecolor,
                  notifier.getwihitecolor),
            ),
          ],
        ),
      ),
    );
  }

  Widget exchangefree() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        Card(
          elevation: 0,
          color: notifier.getwihitecolor,
          // color: notifier.getfavorites,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.all(
              Radius.circular(3.sp),
            ),
          ),
          child: Container(
            color: Colors.transparent,
            height: height / 10,
            width: width / 1.65,
            child: Row(
              children: [
                SizedBox(width: width / 30),
                Image.asset("assets/images/exchangenow.png",
                    height: height / 17),
                SizedBox(width: width / 40),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    SizedBox(height: height / 35),
                    Text(
                      LanguageEn.exchangefee,
                      style: TextStyle(
                          color: notifier.getgrey,
                          fontSize: 12.sp,
                          fontFamily: 'Gilroy_Medium'),
                    ),
                    Text(
                      "0.08%",
                      style: TextStyle(
                          color: notifier.getblck,
                          fontSize: 13.sp,
                          fontFamily: 'Gilroy_Medium'),
                    )
                  ],
                )
              ],
            ),
          ),
        ),
        Card(
          elevation: 0,
          color: notifier.getwihitecolor,
          // color: notifier.getfavorites,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.all(
              Radius.circular(3.sp),
            ),
          ),
          child: Container(
            color: Colors.transparent,
            height: height / 10,
            width: width / 3.6,
            child: Center(
              child: Text(
                "\$32",
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 18.sp,
                    fontFamily: 'Gilroy_Bold'),
              ),
            ),
          ),
        )
      ],
    );
  }

  Widget exchangestock() {
    return Center(
      child: Card(
        elevation: 0,
        color: notifier.getwihitecolor,
        // color: notifier.getfavorites,
        child: Container(
          color: Colors.transparent,
          height: height / 3.4,
          width: width / 1,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: height / 50),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.yousend,
                    style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 14.5.sp,
                      fontFamily: 'Gilroy_Medium',
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    "0.14689",
                    style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 15.sp,
                      fontFamily: 'Gilroy_Bold',
                    ),
                  ),
                  SizedBox(width: width / 2.6),
                  Image.asset("assets/images/airtel.jpg", height: height / 30),
                  SizedBox(width: width / 50),
                  Text(
                    LanguageEn.eth,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getgrey,
                  )
                ],
              ),
              SizedBox(height: height / 30),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Expanded(
                    child: Divider(
                      indent: 20.0,
                      endIndent: 10.0,
                      thickness: 1,
                    ),
                  ),
                  Image.asset("assets/images/data.png", height: height / 40),
                  const Expanded(
                    child: Divider(
                      indent: 10.0,
                      endIndent: 20.0,
                      thickness: 1,
                    ),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.youreceive,
                    style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 14.5.sp,
                      fontFamily: 'Gilroy_Medium',
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    "22878.12",
                    style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 15.sp,
                      fontFamily: 'Gilroy_Bold',
                    ),
                  ),
                  SizedBox(width: width / 2.7),
                  Image.asset("assets/images/Dai.png", height: height / 30),
                  SizedBox(width: width / 50),
                  Text(
                    LanguageEn.dai,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 13.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                  Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getgrey,
                  )
                ],
              ),
              SizedBox(height: height / 30),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.circle, color: notifier.getbluecolor, size: 10.sp),
                  SizedBox(width: width / 50),
                  Text(
                    "1 ETH = 12.489 DAI",
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontFamily: 'Gilroy_Medium',
                        fontSize: 12.sp),
                  )
                ],
              )
            ],
          ),
        ),
      ),
    );
  }
}

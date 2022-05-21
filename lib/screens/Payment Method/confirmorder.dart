import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stock_exchange_tabs/transactioncomplete.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../card_type/custtomcsrdtype.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ConfirmOrder extends StatefulWidget {
  const ConfirmOrder({Key? key}) : super(key: key);

  @override
  State<ConfirmOrder> createState() => _ConfirmOrderState();
}

class _ConfirmOrderState extends State<ConfirmOrder> {
  late ColorNotifier notifier;
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
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            notifier.getwihitecolor, "Confirm Order", notifier.getblck,
            height: height / 15),

      ),
    );
  }

  Widget exchangestock(image, title, subtitle, price) {
    return Center(
      child: Container(
        height: height / 12,
        width: width / 1.1,
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15)),
          color: notifier.getconcirmstockbuycolor,
        ),
        child: Row(
          children: [
            SizedBox(width: width / 12),
            Image.asset(image, height: height / 20),
            SizedBox(width: width / 15),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: height / 50),
                Text(
                  title,
                  style: TextStyle(
                      fontSize: 15.sp,
                      fontFamily: 'Gilroy_Bold',
                      color: notifier.getblck),
                ),
                SizedBox(width: width / 100),
                Text(
                  subtitle,
                  style: TextStyle(
                      fontSize: 12.sp,
                      fontFamily: 'Gilroy-Regular',
                      color: notifier.getgrey),
                ),
              ],
            ),
            const Spacer(),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: height / 35),
                Row(
                  children: [
                    Text(
                      price,
                      style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 17.sp,
                        fontFamily: 'Gilroy_Bold',
                      ),
                    ),
                    SizedBox(width: width / 15),
                  ],
                ),
              ],
            )
          ],
        ),
      ),
    );
  }

  Widget paymentdetails(image, name, subname, color) {
    return Row(
      children: [
        SizedBox(width: width / 17),
        Image.asset(image, height: height / 30),
        SizedBox(width: width / 40),
        Text(
          name,
          style: TextStyle(
              fontFamily: 'Gilroy_Medium',
              fontSize: 15.sp,
              color: notifier.getblck),
        ),
        const Spacer(),
        Text(
          subname,
          style: TextStyle(
              fontFamily: 'Gilroy_Bold', fontSize: 15.sp, color: color),
        ),
        SizedBox(width: width / 17),
      ],
    );
  }
}

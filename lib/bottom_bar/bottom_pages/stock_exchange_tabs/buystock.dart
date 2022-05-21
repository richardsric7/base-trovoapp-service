import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/Payment%20Method/paymentmethod.dart';
import 'package:provider/provider.dart';

import '../../../utils/medeiaqury/medeiaqury.dart';

class BuyStock extends StatefulWidget {
  const BuyStock({Key? key}) : super(key: key);

  @override
  State<BuyStock> createState() => _BuyStockState();
}

class _BuyStockState extends State<BuyStock> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          notifier.getwihitecolor,
          "Buy Stock",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    "\$",
                    style: TextStyle(
                        fontSize: 40.sp,
                        color: notifier.getblck,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                  SizedBox(width: width / 70),
                  pricetextfild()
                ],
              ),
              SizedBox(height: height / 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  maxprice("\$1"),
                  maxprice("\$10"),
                  maxprice("\$50"),
                  maxprice("\$100"),
                ],
              ),
              SizedBox(height: height / 20),
              GestureDetector(
                  onTap: () {
                    Get.to(const PaymentMethod());
                  },
                  child: Button(
                      "Buy", notifier.getbluecolor, notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }

  Widget maxprice(txt) {
    return Container(
      color: Colors.transparent,
      height: height / 15,
      width: width / 5.5,
      child: Card(
        elevation: 1,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12.0),
        ),
        child: Center(
          child: Text(
            txt,
            style: TextStyle(
                fontSize: 25.sp,
                color: notifier.getgrey,
                fontFamily: 'Gilroy_Bold'),
          ),
        ),
      ),
    );
  }

  Widget pricetextfild() {
    return Container(
      color: Colors.transparent,
      width: width / 5,
      child: TextField(
        keyboardType: TextInputType.number,
        style: TextStyle(
            fontSize: 30.sp,
            color: notifier.getblck,
            fontFamily: 'Gilroy_Bold'),
        decoration: InputDecoration(
            border: InputBorder.none,
            hintText: "0",
            hintStyle: TextStyle(color: notifier.getblck, fontSize: 40.sp)),
      ),
    );
  }
}

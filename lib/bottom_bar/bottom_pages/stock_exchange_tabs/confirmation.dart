import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/success.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class Confirmation extends StatefulWidget {
  const Confirmation({Key? key}) : super(key: key);

  @override
  State<Confirmation> createState() => _ConfirmationState();
}

class _ConfirmationState extends State<Confirmation> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(context, notifier.getwihitecolor,
            LanguageEn.confirmation, notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 40),
              confirmation(),
              SizedBox(height: height / 35),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Icon(
                    Icons.lock,
                    color: notifier.getgrey,
                    size: 16.sp,
                  ),
                  SizedBox(width: width / 100),
                  Text(
                    LanguageEn.processedbycryptoline,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 12.sp,
                        fontFamily: 'Gilroy_Medium'),
                  )
                ],
              ),
              Row(
                children: [
                  SizedBox(width: width / 9.5),
                  Text(
                    LanguageEn.securedbythe,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontSize: 12.sp,
                        fontFamily: 'Gilroy_Medium'),
                  ),
                  Text(
                    LanguageEn.privacypolicy,
                    style: TextStyle(
                        color: notifier.getbluecolor,
                        fontSize: 12.sp,
                        fontFamily: 'Gilroy_Medium'),
                  ),
                ],
              ),
              SizedBox(height: height / 3.8),
              GestureDetector(
                onTap: () {
                  Get.to(() => const Success());
                },
                child: Button(LanguageEn.confirmorder, notifier.getbluecolor,
                    notifier.getwihitecolor),
              )
            ],
          ),
        ),
      ),
    );
  }

  Widget confirmation() {
    return Center(
      child: Card(
        color: notifier.getfavorites,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.all(
            Radius.circular(10.sp),
          ),
        ),
        child: Container(
          color: Colors.transparent,
          height: height / 2.7,
          width: width / 1.1,
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Column(
                    children: [
                      Image.asset("assets/images/airtel.jpg",
                          height: height / 15),
                      SizedBox(height: height / 50),
                      Text(
                        LanguageEn.eth,
                        style: TextStyle(
                            fontFamily: 'Gilroy_Bold', fontSize: 13.sp),
                      )
                    ],
                  ),
                  SizedBox(width: width / 15),
                  Column(
                    children: [
                      SizedBox(height: height / 40),
                      Text(
                        LanguageEn.to,
                        style: TextStyle(
                            fontSize: 15.sp,
                            fontFamily: 'Gilroy_Medium',
                            color: notifier.getgrey),
                      )
                    ],
                  ),
                  SizedBox(width: width / 15),
                  Column(
                    children: [
                      Image.asset("assets/images/Dai.png", height: height / 15),
                      SizedBox(height: height / 50),
                      Text(
                        LanguageEn.dai,
                        style: TextStyle(
                            fontFamily: 'Gilroy_Bold', fontSize: 13.sp),
                      ),
                    ],
                  ),
                ],
              ),
              Container(
                color: Colors.transparent,
                height: height / 50,
                width: width,
                child: Padding(
                  padding: EdgeInsets.only(
                      left: width / 30, right: width / 30, top: height / 30),
                  child: Divider(
                    color: notifier.getgrey,
                  ),
                ),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.amount,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontFamily: 'Gilroy_Medium',
                        fontSize: 13.sp),
                  ),
                  const Spacer(),
                  Text(
                    "22878.12 DAI",
                    style: TextStyle(
                        color: notifier.getblck,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 13.sp),
                  ),
                  SizedBox(width: width / 20),
                ],
              ),
              SizedBox(height: height / 50),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.free,
                    style: TextStyle(
                        color: notifier.getgrey,
                        fontFamily: 'Gilroy_Medium',
                        fontSize: 13.sp),
                  ),
                  const Spacer(),
                  Text(
                    "0",
                    style: TextStyle(
                        color: notifier.getblck,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 13.sp),
                  ),
                  SizedBox(width: width / 20),
                ],
              ),
              Container(
                color: Colors.transparent,
                height: height / 50,
                width: width,
                child: Padding(
                  padding: EdgeInsets.only(
                      left: width / 30, right: width / 30, top: height / 30),
                  child: Divider(
                    color: notifier.getgrey,
                  ),
                ),
              ),
              SizedBox(height: height / 25),
              Row(
                children: [
                  SizedBox(width: width / 20),
                  Text(
                    LanguageEn.total,
                    style: TextStyle(
                        color: notifier.getblck,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 13.sp),
                  ),
                  const Spacer(),
                  Text(
                    "22878.12 DAI",
                    style: TextStyle(
                        color: notifier.getbluecolor,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 13.sp),
                  ),
                  SizedBox(width: width / 20),
                ],
              )
            ],
          ),
        ),
      ),
    );
  }
}

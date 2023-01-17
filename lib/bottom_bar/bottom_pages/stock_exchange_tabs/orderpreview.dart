import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/success.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class OrderPreview extends StatefulWidget {
  const OrderPreview({Key? key}) : super(key: key);

  @override
  State<OrderPreview> createState() => _OrderPreviewState();
}

class _OrderPreviewState extends State<OrderPreview> {
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
            LanguageEn.orderpreview, notifier.getblck,
            height: height / 15),
        body: Column(
          children: [
            SizedBox(height: height / 20),
            Center(
              child: Text(
                LanguageEn.balance,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 14.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
            ),
            SizedBox(height: height / 80),
            Text(
              "\$9,876.09",
              style: TextStyle(
                  color: notifier.getblck,
                  fontSize: 27.sp,
                  fontFamily: 'Gilroy_Bold'),
            ),
            SizedBox(height: height / 45),
            balancecard(),
            SizedBox(height: height / 40),
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
            SizedBox(height: height / 9.5),
            GestureDetector(
                onTap: () {
                  Get.to(() => const Success());
                },
                child: Button(LanguageEn.confirmorder, notifier.getbluecolor,
                    notifier.getwihitecolor))
          ],
        ),
      ),
    );
  }

  Widget balancecard() {
    return Card(
      color: notifier.getfavorites,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.all(
          Radius.circular(10.sp),
        ),
      ),
      child: Container(
        color: Colors.transparent,
        height: height / 2.43,
        width: width / 1.1,
        child: Column(
          children: [
            SizedBox(height: height / 50),
            Row(
              children: [
                SizedBox(width: width / 20),
                Image.asset("assets/images/RBL.jpg", height: height / 22),
                SizedBox(width: width / 30),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      LanguageEn.bitcoin,
                      style: TextStyle(
                          color: notifier.getblck,
                          fontFamily: 'Gilroy_Bold',
                          fontSize: 16.sp),
                    ),
                    Text(
                      LanguageEn.btc,
                      style: TextStyle(
                          color: notifier.getgrey,
                          fontFamily: 'Gilroy_Medium',
                          fontSize: 13.sp),
                    ),
                  ],
                ),
                SizedBox(width: width / 3),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text(
                      "\$567.06",
                      style: TextStyle(
                          color: notifier.getblck,
                          fontFamily: 'Gilroy_Bold',
                          fontSize: 17.sp),
                    ),
                    Text(
                      "1.786876 BTC",
                      style: TextStyle(
                          color: notifier.getgrey,
                          fontFamily: 'Gilroy_Medium',
                          fontSize: 13.sp),
                    ),
                  ],
                )
              ],
            ),
            Container(
              color: Colors.transparent,
              height: height / 15,
              width: width / 1.2,
              child: Divider(
                color: notifier.getgrey,
              ),
            ),
            Row(
              children: [
                SizedBox(width: width / 20),
                Text(
                  LanguageEn.paymentmethod,
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getgrey,
                      fontFamily: 'Gilroy_Medium'),
                ),
                const Spacer(),
                Image.asset("assets/images/VisaLogo.png", height: height / 50),
                SizedBox(width: width / 20),
              ],
            ),
            Container(
              color: Colors.transparent,
              height: height / 15,
              width: width / 1.2,
              child: Divider(
                color: notifier.getgrey,
              ),
            ),
            Row(
              children: [
                SizedBox(width: width / 20),
                Text(
                  LanguageEn.purchaseprice,
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getgrey,
                      fontFamily: 'Gilroy_Medium'),
                ),
                const Spacer(),
                Text(
                  "\$168,00",
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Medium'),
                ),
                SizedBox(width: width / 20),
              ],
            ),
            SizedBox(height: height / 100),
            Row(
              children: [
                SizedBox(width: width / 20),
                Text(
                  LanguageEn.free,
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getgrey,
                      fontFamily: 'Gilroy_Medium'),
                ),
                const Spacer(),
                Text(
                  "\$30,06",
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Medium'),
                ),
                SizedBox(width: width / 20),
              ],
            ),
            Container(
              color: Colors.transparent,
              height: height / 15,
              width: width / 1.2,
              child: Divider(
                color: notifier.getgrey,
              ),
            ),
            Row(
              children: [
                SizedBox(width: width / 20),
                Text(
                  LanguageEn.total,
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Medium'),
                ),
                const Spacer(),
                Text(
                  "\$298,06",
                  style: TextStyle(
                      fontSize: 15.sp,
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Bold'),
                ),
                SizedBox(width: width / 20),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

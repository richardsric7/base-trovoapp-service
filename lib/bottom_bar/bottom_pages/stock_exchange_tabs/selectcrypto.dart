import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/custtom_slock_list/custtom_slock_list.dart';
import 'package:gocrypto/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stock_exchange_tabs/confirmation.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:gocrypto/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class SelectCrypto extends StatefulWidget {
  const SelectCrypto({Key? key}) : super(key: key);

  @override
  State<SelectCrypto> createState() => _SelectCryptoState();
}

class _SelectCryptoState extends State<SelectCrypto> {
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
            notifier.getwihitecolor, LanguageEn.selectcrypto, notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.only(left: 10, right: 10),
            child: Column(
              children: [
                // SizedBox(height: height / 50),
                Center(
                  child: Customtextfild.textField(
                      "Search  company, stocks...",
                      notifier.getbluecolor,
                      Icons.search_rounded,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      notifier.getgrey,
                      50.sp,
                      310.sp),
                ),

                SizedBox(height: height / 20),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/airtel.jpg",
                      LanguageEn.ethereum,
                      LanguageEn.eth,
                      "\$127,00",
                      "10,03%",
                      const Color(0xff22C36B)),
                ),
                SizedBox(height: height / 30),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/Ambuja_logo.png",
                      LanguageEn.binance,
                      LanguageEn.bnb,
                      "\$297,64",
                      "2,87%",
                      const Color(0xffF65556)),
                ),
                SizedBox(height: height / 30),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/kotok.png",
                      LanguageEn.litecoin,
                      LanguageEn.ltc,
                      "\$326,23",
                      "2,87%",
                      const Color(0xff22C36B)),
                ),
                SizedBox(height: height / 30),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/icici.png",
                      LanguageEn.xrp,
                      LanguageEn.xrp,
                      "\$326,23",
                      "2,87%",
                      const Color(0xffF65556)),
                ),
                SizedBox(height: height / 30),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/Ambuja_logo.png",
                      LanguageEn.binance,
                      LanguageEn.bnb,
                      "\$297,64",
                      "2,87%",
                      const Color(0xffF65556)),
                ),
                SizedBox(height: height / 30),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/icici.png",
                      LanguageEn.xrp,
                      LanguageEn.xrp,
                      "\$326,23",
                      "2,87%",
                      const Color(0xffF65556)),
                ),
                SizedBox(height: height / 30),
                GestureDetector(
                  onTap: () {
                    Get.to(() => const Confirmation());
                  },
                  child: CusttomButton(
                      "assets/images/Ambuja_logo.png",
                      LanguageEn.binance,
                      LanguageEn.bnb,
                      "\$326,23",
                      "2,87%",
                      const Color(0xff22C36B)),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

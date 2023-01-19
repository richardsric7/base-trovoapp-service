import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_bar.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';

class Success extends StatefulWidget {
  const Success({Key? key}) : super(key: key);

  @override
  State<Success> createState() => _SuccessState();
}

class _SuccessState extends State<Success> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 8),
              Center(
                child: Image.asset("assets/images/vrfcomplate.png",
                    height: height / 3),
              ),
              SizedBox(height: height / 20),
              Text(
                LanguageEn.success,
                style: TextStyle(
                    fontFamily: 'Gilroy_Bold',
                    fontSize: 24.sp,
                    color: notifier.getblck),
              ),
              SizedBox(height: height / 20),
              Container(
                height: height / 15,
                width: width / 2.8,
                decoration: BoxDecoration(
                  color: notifier.getfavorites,
                  borderRadius: BorderRadius.all(
                    Radius.circular(13.sp),
                  ),
                ),
                child: Center(
                  child: Text(
                    "22878.12 DAI",
                    style: TextStyle(
                        color: notifier.getbluecolor,
                        fontFamily: 'Gilroy_Bold',
                        fontSize: 14.sp),
                  ),
                ),
              ),
              SizedBox(height: height / 40),
              Text(
                LanguageEn.hasbeenexchange,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 12.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
              SizedBox(height: height / 5.2),
              GestureDetector(
                  onTap: () {
                    Get.to(() => const BottomHome());
                  },
                  child: Button(LanguageEn.done, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }
}

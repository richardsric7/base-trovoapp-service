import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/button_tabs/chart.dart';
import 'package:provider/provider.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class Ordertabs extends StatefulWidget {
  const Ordertabs({Key? key}) : super(key: key);

  @override
  State<Ordertabs> createState() => _OrdertabsState();
}

class _OrdertabsState extends State<Ordertabs> {
  late ColorNotifier notifier;

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
              SizedBox(height: height / 9),
              Center(
                child: Image.asset("assets/images/noorder.png",
                    height: height / 5),
              ),
              SizedBox(height: height / 50),
              Text(
                "No Order Yet",
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 21.sp,
                    fontFamily: 'Gilroy_Bold'),
              ),
              SizedBox(height: height / 100),
              Text(
                "There is no recent stock you order, let’s\nmake your first investment!",
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: 'Gilroy-Regular'),
              ),
              SizedBox(height: height / 30),
              GestureDetector(
                  onTap: () {
                    Get.to(const Chart());
                  },
                  child: buton(notifier.getbluecolor, "View Stocks",
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }

  Widget buton(colorbutton, buttontext, buttontextcolor) {
    return Center(
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(15),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            LayoutBuilder(builder: (context, constraints) {
              return Container(
                height: height / 15,
                width: width / 1.9,
                decoration: BoxDecoration(
                  color: colorbutton!,
                  borderRadius: BorderRadius.circular(15.sp),
                ),
                child: Center(
                  child: Text(
                    buttontext!,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                        fontFamily: 'Gilroy_Medium',
                        fontSize: 15.sp,
                        color: buttontextcolor),
                  ),
                ),
              );
            }),
          ],
        ),
      ),
    );
  }
}

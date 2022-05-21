import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';
import '../notifire_clor.dart';

// ignore: must_be_immutable
class CusttomStoc extends StatefulWidget {
  final String? image;
  final String? title;
  final String? subtitle;
  final String? stockprice;
  final String? updown;
  final Color? iconcolor;
  final String? graphimage;

  const CusttomStoc(this.image, this.title, this.subtitle, this.stockprice,
      this.updown, this.iconcolor, this.graphimage,
      {Key? key})
      : super(key: key);

  @override
  State<CusttomStoc> createState() => _CusttomStocState();
}

class _CusttomStocState extends State<CusttomStoc> {
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
      builder: () => Container(
        margin: EdgeInsets.only(left: width / 20),
        height: height / 2.5,
        width: width / 2.3,
        decoration: BoxDecoration(
          color: notifier.getfavorites,
          borderRadius: BorderRadius.all(
            Radius.circular(20.sp),
          ),
        ),
        child: Column(
          children: [
            Row(
              children: [
                Padding(
                  padding: EdgeInsets.only(left: width / 20, top: height / 60),
                  child: Image.asset(widget.image!, height: height / 20),
                ),
                SizedBox(width: width / 20),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    SizedBox(height: height / 50),
                    Text(
                      widget.title!,
                      style: TextStyle(
                          color: notifier.getblck,
                          fontSize: 12.sp,
                          fontFamily: 'Gilroy_Bold'),
                    ),
                    Text(
                      widget.subtitle!,
                      style: TextStyle(
                          color: notifier.getgrey,
                          fontSize: 11.sp,
                          fontFamily: 'Gilroy-Regular'),
                    ),
                  ],
                )
              ],
            ),
            SizedBox(height: height / 30),
            Center(
              child: Image.asset(widget.graphimage!, height: height / 13),
            ),
            SizedBox(height: height / 50),
            Row(
              children: [
                SizedBox(width: width / 25),
                Text(
                  widget.stockprice!,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontFamily: 'Gilroy_Bold',
                      fontSize: 15.sp),
                ),
                const Spacer(),
                Text(
                  widget.updown!,
                  style: TextStyle(
                      color: const Color(0xff22C36B), fontSize: 13.sp),
                ),
                SizedBox(width: width / 25),
              ],
            )
          ],
        ),
      ),
    );
  }
}

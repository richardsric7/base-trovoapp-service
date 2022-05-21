import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';
import '../notifire_clor.dart';

class CusttomButton extends StatefulWidget {
  final String? image;
  final String? title;
  final String? subtitle;
  final String? price;
  final String? updown;
  final Color? graphcolor;

  const CusttomButton(this.image, this.title, this.subtitle, this.price,
      this.updown, this.graphcolor,
      {Key? key})
      : super(key: key);

  @override
  State<CusttomButton> createState() => _CusttomButtonState();
}

class _CusttomButtonState extends State<CusttomButton> {
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
      builder: () => Center(
        child: Container(
          color: Colors.transparent,
          height: height / 17,
          child: Row(
            children: [
              SizedBox(width: width / 15),
              Image.asset(widget.image!, height: height / 20),
              SizedBox(width: width / 30),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    widget.title!,
                    style: TextStyle(
                        fontSize: 11.sp,
                        fontFamily: 'Gilroy_Bold',
                        color: notifier.getblck),
                  ),
                  Text(
                    widget.subtitle!,
                    style: TextStyle(
                        fontSize: 12.sp,
                        fontFamily: 'Gilroy-Regular',
                        color: notifier.getgrey),
                  ),
                ],
              ),
              // SizedBox(width: width,),
              const Spacer(),
              Image.asset("assets/images/Watchlist_chart.png",
                  color: widget.graphcolor, height: height / 30),
              const Spacer(),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    widget.price!,
                    style: TextStyle(
                      color: notifier.getblck,
                      fontSize: 15.sp,
                      fontFamily: 'Gilroy_Bold',
                    ),
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Icon(
                        Icons.arrow_drop_up_outlined,
                        color: const Color(0xff19C09A),
                        size: 20.sp,
                      ),
                      Text(
                        widget.updown!,
                        style: TextStyle(
                            color: const Color(0xff19C09A), fontSize: 12.sp),
                      ),
                      SizedBox(width: width / 15),
                    ],
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

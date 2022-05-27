// ignore_for_file: must_be_immutable

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';
import '../fonts.dart';
import '../notifire_clor.dart';

class Button extends StatefulWidget {
  final String? buttontext;
  final Color? colorbutton;
  final Color? buttontextcolor;

  const Button(this.buttontext, this.colorbutton, this.buttontextcolor,
      {Key? key})
      : super(key: key);

  @override
  State<Button> createState() => _ButtonState();
}

class _ButtonState extends State<Button> {
  get borderRadius => BorderRadius.circular(15);

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
          decoration: BoxDecoration(
            borderRadius: borderRadius,
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              LayoutBuilder(builder: (context, constraints) {
                return Container(
                  height: height / 15,
                  width: width / 1.1,
                  decoration: BoxDecoration(
                    color: widget.colorbutton!,
                    borderRadius: borderRadius,
                  ),
                  child: Center(
                    child: Text(
                      widget.buttontext!,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                          fontFamily: fontbody,
                          fontSize: 15.sp,
                          color: widget.buttontextcolor),
                    ),
                  ),
                );
              }),
            ],
          ),
        ),
      ),
    );
  }
}

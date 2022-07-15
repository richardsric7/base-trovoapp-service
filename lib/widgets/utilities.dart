import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/widgets/popups.dart';

void showSnackBar(String rel, BuildContext context) {
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      backgroundColor: notifier.getbluecolor,
      content: Text(
        '$rel copied successfully',
        style: TextStyle(
          color: notifier.getwihitecolor,
          fontSize: 12.sp,
          fontWeight: FontWeight.w500,
          fontFamily: fontbody,
        ),
      ),
      action: SnackBarAction(
        label: 'DISMISS',
        textColor: Colors.white,
        onPressed: () => {
          ScaffoldMessenger.of(context).clearSnackBars(),
        },
      ),
    ),
  );
}

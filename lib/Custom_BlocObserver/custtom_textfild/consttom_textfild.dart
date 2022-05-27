import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gocrypto/Custom_BlocObserver/fonts.dart';

class Customtextfild {
  static Widget textField(labletext, focuscolor, preicon, lablecolor, iconcolor,
      textcolor, bordercolor, h, w) {
    return ScreenUtilInit(
      builder: () => Container(
        color: Colors.transparent,
        height: h,
        width: w,
        child: TextField(
          style: TextStyle(color: textcolor, fontFamily: fontbody),
          cursorColor: lablecolor,
          onChanged: (value) {},
          // obscureText: hidePassword, //show/hide password
          decoration: InputDecoration(
            label: Text(labletext),
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            prefixIcon: Icon(preicon, color: iconcolor),
            labelStyle: TextStyle(color: lablecolor),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            enabledBorder: OutlineInputBorder(
              borderSide: BorderSide(color: bordercolor, width: 1),
              borderRadius: BorderRadius.circular(15.sp),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: BorderSide(color: focuscolor, width: 1),
              borderRadius: BorderRadius.circular(15.sp),
            ),
          ),
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gocrypto/Custom_BlocObserver/fonts.dart';

class Custompasswordtextfild {
  static Widget textField(
      labletext, focuscolor, preicon, lablecolor, iconcolor, textcolor) {
    bool hidePassword = true;
    return ScreenUtilInit(
      builder: () => Container(
        color: Colors.transparent,
        height: 45.h,
        width: 300.w,
        child: TextField(
          style: TextStyle(color: textcolor, fontFamily: fontbody),
          onChanged: (value) {},
          obscureText: hidePassword, //show/hide password
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
              borderSide: const BorderSide(color: Colors.grey, width: 1.0),
              borderRadius: BorderRadius.circular(15.sp),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: BorderSide(color: focuscolor, width: 1.0),
              borderRadius: BorderRadius.circular(15.sp),
            ),
          ),
        ),
      ),
    );
  }
}

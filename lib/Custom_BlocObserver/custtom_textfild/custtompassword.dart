import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:stellar_flutter_sdk/stellar_flutter_sdk.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';

class Custompasswordtextfild {
  static Widget textField(
      labletext, focuscolor, preicon, lablecolor, iconcolor, textcolor) {
    bool hidePassword = true;
    return ScreenUtilInit(
      builder: (context, child) => Container(
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

class CustomPasswordFormField {
  static Widget textField(
      labletext, focuscolor, preicon, lablecolor, iconcolor, textcolor, h, w,
      {onChanged, validator, onSaved}) {
    bool hidePassword = true;
    return ScreenUtilInit(
      builder: (context, child) => Container(
        color: Colors.transparent,
        height: h,
        width: w,
        child: TextFormField(
          style: TextStyle(color: textcolor, fontFamily: fontbody),
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
          onChanged: onChanged,
          validator: validator,
          onSaved: onSaved,
        ),
      ),
    );
  }
}

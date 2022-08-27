import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';

class Customtextfild {
  static Widget textField(labletext, focuscolor, preicon, lablecolor, iconcolor,
      textcolor, bordercolor, h, w) {
    return ScreenUtilInit(
      builder: (context, child) => Container(
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

class CustomTextFormField {
  static Widget textField(
    labletext,
    focuscolor,
    preicon,
    lablecolor,
    iconcolor,
    textcolor,
    bordercolor,
    h,
    w, {
    initialValue,
    onChanged,
    maxLength,
    validator,
    onSaved,
    keyboardtype,
    helperText,
    inputFormatters,
    controller,
    buildCounter,
    key,
  }) {
    return ScreenUtilInit(
      builder: (context, child) => Container(
        color: Colors.transparent,
        height: h,
        width: w,
        child: TextFormField(
          key: key,
          maxLength: maxLength,
          style: TextStyle(color: textcolor, fontFamily: fontbody),
          initialValue: initialValue,
          cursorColor: lablecolor,
          onChanged: onChanged,
          decoration: InputDecoration(
            counterStyle: TextStyle(
              fontFamily: fontbody,
              color: textcolor,
            ),
            errorStyle: TextStyle(
              fontFamily: fontbody,
            ),
            helperText: helperText,
            helperStyle: TextStyle(
              fontSize: 12,
              fontFamily: fontbody,
            ),
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
          inputFormatters: inputFormatters,
          keyboardType: keyboardtype,
          validator: validator,
          controller: controller,
          onSaved: onSaved,
          buildCounter: buildCounter,
        ),
      ),
    );
  }
}

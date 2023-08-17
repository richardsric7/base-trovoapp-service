import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';

class Customtextfild {
  static Widget textField(labletext, focuscolor, preicon, lablecolor, iconcolor,
      textcolor, bordercolor, h, w) {
    return Container(
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
    readOnly = false,
    onTap,
    key,
  }) {
    return Container(
      color: Colors.transparent,
      height: h,
      width: w,
      child: TextFormField(
        key: key,
        maxLength: maxLength,
        readOnly: readOnly,
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
          prefixIcon: preicon == null ? null : Icon(preicon, color: iconcolor),
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
        onTap: onTap,
        buildCounter: buildCounter,
      ),
    );
  }

  static Widget textFieldWithoutIcon(
    labletext,
    focuscolor,
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
    readOnly = false,
    key,
  }) {
    return Container(
      color: Colors.transparent,
      height: h,
      width: w,
      child: TextFormField(
        key: key,
        maxLength: maxLength,
        readOnly: readOnly,
        style: TextStyle(
            color: textcolor,
            overflow: TextOverflow.visible,
            fontFamily: fontbody),
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
            borderRadius: BorderRadius.circular(10.sp),
          ),
          labelStyle: TextStyle(color: lablecolor),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(10.sp),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: bordercolor, width: 1),
            borderRadius: BorderRadius.circular(10.sp),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: focuscolor, width: 1),
            borderRadius: BorderRadius.circular(10.sp),
          ),
        ),
        inputFormatters: inputFormatters,
        keyboardType: keyboardtype,
        validator: validator,
        controller: controller,
        onSaved: onSaved,
        buildCounter: buildCounter,
      ),
    );
  }
}

Widget multilineInput(
  labletext,
  focuscolor,
  lablecolor,
  textcolor,
  bordercolor,
  h,
  w, {
  onChanged,
  maxLength,
  minLines,
  maxLines,
  validator,
  onSaved,
  keyboardtype,
  focusNode,
}) {
  return Container(
    height: h,
    width: w,
    child: TextFormField(
      focusNode: focusNode,
      maxLength: maxLength,
      minLines: minLines,
      maxLines: maxLines,
      style: TextStyle(color: textcolor, fontFamily: fontbody),
      cursorColor: lablecolor,
      onChanged: onChanged,
      decoration: InputDecoration(
        hintText: labletext,
        hintStyle: TextStyle(color: lablecolor),
        disabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(15),
        ),
        labelStyle: TextStyle(color: lablecolor),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(15),
        ),
        enabledBorder: OutlineInputBorder(
          borderSide: BorderSide(color: bordercolor, width: 1),
          borderRadius: BorderRadius.circular(15),
        ),
        focusedBorder: OutlineInputBorder(
          borderSide: BorderSide(color: focuscolor, width: 1),
          borderRadius: BorderRadius.circular(15),
        ),
      ),
      keyboardType: keyboardtype,
      validator: validator,
      onSaved: onSaved,
    ),
  );
}

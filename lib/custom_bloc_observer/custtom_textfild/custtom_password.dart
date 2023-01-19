import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';

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

class CustomPasswordFormField extends StatefulWidget {
  String? labelText;
  Color? focusColor;
  IconData? preIcon;
  Color? labelColor;
  Color? iconColor;
  Color? textColor;
  double? height;
  double? width;
  int? maxLength;
  TextEditingController? controller;
  final void Function(String?)? onChanged;
  final void Function(String?)? onSubmitted;
  final String? Function(String?)? validator;
  final void Function(String?)? onSaved;
  TextInputAction? textInputAction;
  FocusNode? focusNode;

  CustomPasswordFormField(
    this.labelText,
    this.focusColor,
    this.preIcon,
    this.labelColor,
    this.iconColor,
    this.textColor,
    this.height,
    this.width, {
    Key? key,
    this.maxLength,
    this.onSubmitted,
    this.onChanged,
    this.validator,
    this.onSaved,
    this.focusNode,
    this.controller,
    this.textInputAction,
  }) : super(key: key);

  @override
  State<CustomPasswordFormField> createState() =>
      _CustomPasswordFormFieldState();
}

class _CustomPasswordFormFieldState extends State<CustomPasswordFormField> {
  bool hidePassword = true;

  @override
  Widget build(BuildContext context) {
    return ScreenUtilInit(
      builder: (context, child) => Container(
        color: Colors.transparent,
        height: widget.height,
        width: widget.width,
        child: TextFormField(
          focusNode: widget.focusNode,
          maxLength: widget.maxLength,
          controller: widget.controller,
          style: TextStyle(color: widget.textColor, fontFamily: fontbody),
          obscureText: hidePassword, //show/hide password
          textInputAction: widget.textInputAction,
          decoration: InputDecoration(
            counterStyle: TextStyle(
              fontFamily: fontbody,
              color: widget.textColor,
            ),
            label: Text(widget.labelText!),
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            prefixIcon: Icon(widget.preIcon, color: widget.iconColor),
            suffixIcon: IconButton(
                onPressed: () {
                  setState(() {
                    hidePassword = !hidePassword;
                  });
                },
                icon: Icon(
                  getSuffixIcon(),
                  color: widget.textColor,
                  size: height / 50,
                )),
            labelStyle: TextStyle(color: widget.labelColor),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            enabledBorder: OutlineInputBorder(
              borderSide: const BorderSide(color: Colors.grey, width: 1.0),
              borderRadius: BorderRadius.circular(15.sp),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: BorderSide(color: widget.focusColor!, width: 1.0),
              borderRadius: BorderRadius.circular(15.sp),
            ),
          ),
          onChanged: widget.onChanged,
          onFieldSubmitted: widget.onSubmitted,
          validator: widget.validator,
          onSaved: widget.onSaved,
        ),
      ),
    );
  }

  IconData getSuffixIcon() =>
      hidePassword ? CupertinoIcons.eye : CupertinoIcons.eye_slash;
}

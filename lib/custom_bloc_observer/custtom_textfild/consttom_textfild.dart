import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

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
            borderRadius: BorderRadius.circular(15),
          ),
          prefixIcon: Icon(preicon, color: iconcolor),
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
    hintText,
    autoFormatNumber = false,
    inputFormatters,
    controller,
    buildCounter,
    readOnly = false,
    onTap,
    key,
  }) {
    if (autoFormatNumber && controller == null) {
      throw Exception(
          "Please supply controller in order to enable number auto formatting!");
    }

    return Container(
      color: Colors.transparent,
      height: double.parse(h.toString()),
      width: double.parse(w.toString()),
      child: Focus(
        onFocusChange: (hasFocus) {
          if (autoFormatNumber) {
            if (!hasFocus) {
              var newVal = controller.text;
              // check if there are multiple dots on the text
              var splitText = newVal.split('.');
              if (splitText.length > 2) {
                // remove all dots except the first one.
                newVal = '${splitText[0]}.${splitText[1]}';
              }

              if (newVal == '.') {
                newVal = '';
                controller.text = newVal;
              }

              if (newVal.isNotEmpty) {
                controller!.text = truncateNumber(double.parse(
                    newVal.toString().replaceAll(',', '').replaceAll('-', '')));

                if (!newVal.endsWith('.')) {
                  controller.selection =
                      TextSelection.collapsed(offset: controller.selection.end);
                }
                newVal = newVal.replaceAll(',', '');
              }
              print('not has focus');
            } else {
              print('has focus');
            }
          }
        },
        child: TextFormField(
          key: key,
          maxLength: maxLength,
          readOnly: readOnly,
          style: TextStyle(color: textcolor, fontFamily: fontbody),
          initialValue: initialValue,
          cursorColor: lablecolor,
          onChanged: (newVal) {
            if (onChanged != null) {
              onChanged(newVal);
            }
          },
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
            hintText: hintText,
            label: labletext != null ? Text(labletext) : null,
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15),
            ),
            prefixIcon:
                preicon == null ? null : Icon(preicon, color: iconcolor),
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
          inputFormatters: inputFormatters,
          keyboardType: keyboardtype,
          validator: (value) {
            if (validator != null) {
              var newVal = value;
              if (autoFormatNumber) {
                newVal = value.toString().replaceAll(',', '');
              }

              return validator(newVal);
            }
            return null;
          },
          controller: controller,
          onSaved: (value) {
            if (onSaved != null) {
              var newVal = value;
              if (autoFormatNumber) {
                newVal = value.toString().replaceAll(',', '');
              }

              onSaved(newVal);
            }
          },
          onTap: onTap,
          buildCounter: buildCounter,
        ),
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
    fontsize,
    inputFormatters,
    controller,
    buildCounter,
    readOnly = false,
    key,
  }) {
    return Container(
      color: Colors.transparent,
      width: w,
      constraints: BoxConstraints(minHeight: h),
      child: TextFormField(
        key: key,
        maxLength: maxLength,
        readOnly: readOnly,
        style: TextStyle(
          color: textcolor,
          overflow: TextOverflow.visible,
          fontFamily: fontbody,
          fontSize: fontsize,
        ),
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
            borderRadius: BorderRadius.circular(10),
          ),
          labelStyle: TextStyle(color: lablecolor),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(10),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: bordercolor, width: 1),
            borderRadius: BorderRadius.circular(10),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: focuscolor, width: 1),
            borderRadius: BorderRadius.circular(10),
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
  initialValue,
}) {
  return Container(
    height: h,
    width: w,
    child: TextFormField(
      focusNode: focusNode,
      maxLength: maxLength,
      minLines: minLines,
      maxLines: maxLines,
      initialValue: initialValue,
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

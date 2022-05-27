import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/fonts.dart';

class CustomAppBar extends PreferredSize {
  final double height;
  final String txt;
  final Color color;
  final Color titlecolor;

  CustomAppBar(
    this.color,
    this.txt,
    this.titlecolor, {
    Key? key,
    required this.height,
  }) : super(
          key: key,
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: color,
            leading: GestureDetector(
              onTap: () {
                Get.back();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            title: Text(
              txt,
              style: TextStyle(color: titlecolor, fontFamily: fontsemibold),
            ),
          ),
          preferredSize: Size.fromHeight(height),
        );
}

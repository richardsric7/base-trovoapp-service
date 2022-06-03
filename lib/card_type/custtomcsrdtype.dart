import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

import '../utils/medeiaqury/medeiaqury.dart';

Widget paymenttype(
    image, cardname, disc, cardnamecolor, disccolor, imagecolor, imagehight) {
  return ScreenUtilInit(
    builder: (context, child) => Padding(
      padding: EdgeInsets.only(left: width / 50, right: width / 50),
      child: Container(
        height: height / 12,
        width: width,
        decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15)),
            color: Colors.transparent,
            border: Border.all(color: const Color(0xff9EA3AE))),
        child: Row(
          children: [
            SizedBox(width: width / 25),
            Image.asset(image, height: imagehight, color: imagecolor),
            SizedBox(width: width / 40),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: height / 50),
                Text(
                  cardname,
                  style: TextStyle(
                      color: cardnamecolor,
                      fontFamily: 'Gilroy_Bold',
                      fontSize: 14.sp),
                ),
                Text(
                  disc,
                  style: TextStyle(
                      color: disccolor,
                      fontFamily: 'Gilroy_Medium',
                      fontSize: 12.sp),
                )
              ],
            ),
          ],
        ),
      ),
    ),
  );
}

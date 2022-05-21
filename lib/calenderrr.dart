import 'package:calendar_calendar/calendar_calendar.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';

Widget calender(
    bgcolor, bgcolortwo, textcolor, bgcolorthree, stcolor, sttwocolor) {
  return ScreenUtilInit(
    builder: () => Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: <Widget>[
          Container(
            color: Colors.transparent,
            height: 305.h,
            child: Calendar(
              weekendOpacityEnable: true,
              previous: Container(
                decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(500),
                    boxShadow: [
                      BoxShadow(
                          color: Colors.grey[300]!,
                          spreadRadius: 1.5,
                          blurRadius: 5,
                          offset: const Offset(2.0, 0.0))
                    ]),
                child: CircleAvatar(
                  radius: 14,
                  backgroundColor: bgcolor,
                  child: const Icon(
                    Icons.arrow_back_ios,
                    size: 16,
                    color: Color(0xff3B82F6),
                  ),
                ),
              ),
              next: Container(
                decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(500),
                    boxShadow: [
                      BoxShadow(
                          color: Colors.grey[300]!,
                          spreadRadius: 1.5,
                          blurRadius: 5,
                          offset: const Offset(2.0, 0.0))
                    ]),
                child: CircleAvatar(
                  radius: 14,
                  backgroundColor: bgcolorthree,
                  child: const Icon(
                    Icons.arrow_forward_ios,
                    size: 16,
                    color: Color(0xff3B82F6),
                  ),
                ),
              ),
              space: 20,
              onSelected: print,
              backgroundColor: bgcolortwo,
              activeColor: const Color(0xff3B82F6),
              textStyleDays:
                  TextStyle(fontWeight: FontWeight.normal, color: stcolor),
              textStyleWeekDay:
                  TextStyle(fontWeight: FontWeight.bold, color: sttwocolor),
              titleStyle:
                  const TextStyle(fontWeight: FontWeight.bold, fontSize: 20),
              selectedStyle:
                  TextStyle(fontWeight: FontWeight.bold, color: textcolor),
            ),
          )
        ],
      ),
    ),
  );
}

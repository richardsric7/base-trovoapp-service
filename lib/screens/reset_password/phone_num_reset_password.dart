import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Auth/vericication.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class PhoneNumResetPassword extends StatefulWidget {
  const PhoneNumResetPassword({Key? key}) : super(key: key);

  @override
  State<PhoneNumResetPassword> createState() => _PhoneNumResetPasswordState();
}

class _PhoneNumResetPasswordState extends State<PhoneNumResetPassword> {
  late ColorNotifier notifier;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  String dropdownvalue = '+91';

  // List of items in our dropdown menu
  var items = [
    '+61',
    '+91',
    '+92',
    '+152',
    '+139',
  ];

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        LanguageEn.setupsecondstep,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 24.sp,
                            fontFamily: 'Gilroy_Bold'),
                      ),
                      SizedBox(height: height / 40),
                      Text(
                        LanguageEn.enteryourphonenumber,
                        style: TextStyle(
                            fontSize: 13.5.sp, color: notifier.getgrey),
                      ),
                      SizedBox(height: height / 50),
                      verification()
                      // Row(
                      //   children: [
                      //     DropdownButtonHideUnderline(
                      //       child: ButtonTheme(
                      //         child: DropdownButton<String>(
                      //           hint: Image.asset("assets/images/flagfour.png",
                      //               height: height / 25),
                      //           value: _selectedindex,
                      //           onChanged: (newValue) {
                      //             setState(() {
                      //               _selectedindex = newValue;
                      //             });
                      //           },
                      //           items: _myjson.map((Map map) {
                      //             return DropdownMenuItem<String>(
                      //               value: map["id"].toString(),
                      //               child: Row(
                      //                 children: <Widget>[
                      //                   Image.asset(
                      //                     map["image"].toString(),
                      //                     width: 35.w,
                      //                   ),
                      //                 ],
                      //               ),
                      //             );
                      //           }).toList(),
                      //         ),
                      //       ),
                      //     ),
                      //     SizedBox(width: width / 20),
                      //     Container(color: Colors.transparent,
                      //       width: width / 1.6,
                      //       child: Customtextfild.textField(
                      //           "Mobile No",
                      //           notifier.getbluecolor,
                      //           Icons.add,
                      //           notifier.getgrey,
                      //           notifier.getblck,
                      //           notifier.getblck,
                      //           notifier.getgrey,45.sp,300.sp),
                      //     ),
                      //   ],
                      // ),
                    ],
                  ),
                ],
              ),
              SizedBox(height: height / 11),
              GestureDetector(
                onTap: () {
                  Get.to(
                    const Veryfication(),
                  );
                },
                child: Button(LanguageEn.continuee, notifier.getbluecolor,
                    notifier.getwihitecolor),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget verification() {
    return Container(
      height: height / 15,
      width: width / 1.2,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.all(
          Radius.circular(15.sp),
        ),
        border: Border.all(
          color: notifier.getbluecolor,
        ),
      ),
      child: Row(
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 10),
            child: DropdownButton(
              dropdownColor: notifier.getwihitecolor,
              underline: const SizedBox(),
              // Initial Value
              value: dropdownvalue,

              // Down Arrow Icon
              icon: Icon(
                Icons.arrow_drop_down,
                color: notifier.getblck,
              ),

              // Array list of items
              items: items.map((String items) {
                return DropdownMenuItem(
                  value: items,
                  child: Text(
                    items,
                    style: TextStyle(color: notifier.getblck),
                  ),
                );
              }).toList(),
              // After selecting the desired option,it will
              // change button value to selected value
              onChanged: (String? newValue) {
                setState(() {
                  dropdownvalue = newValue!;
                });
              },
            ),
          ),
          SizedBox(width: width / 100),
          Container(
            height: height / 22,
            width: width / 250,
            color: notifier.getgrey,
          ),
          SizedBox(width: width / 22),
          Container(
            color: Colors.transparent,
            height: 25.h,
            width: 200.w,
            child: TextField(
              style: TextStyle(color: notifier.getblck),
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(border: InputBorder.none),
            ),
          )
        ],
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stockexchange.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';
import '../notifire_clor.dart';

class CusttomStockGainscard extends StatefulWidget {
  const CusttomStockGainscard({Key? key}) : super(key: key);

  @override
  State<CusttomStockGainscard> createState() => _CusttomStockGainscardState();
}

class _CusttomStockGainscardState extends State<CusttomStockGainscard> {
  get borderRadius => BorderRadius.circular(15);

  bool showAvg = false;

  int touchedIndex = -1;
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

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    // String dropdownvalue = 'This Week';
    // var items = [
    //   'This Week',
    //   'This Month',
    //   'This year',
    //   'This Hours',
    // ];
    return ScreenUtilInit(
      builder: () => Column(
        children: [
          Stack(
            children: [
              Center(
                child: Container(
                  color: Colors.transparent,
                  height: height / 5.5,
                  width: width / 1.1,
                  child: Card(
                    color: const Color(0xff2873FF),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(15.sp),
                    ),
                    child: Column(
                      children: [
                        SizedBox(height: height / 50),
                        Padding(
                          padding: EdgeInsets.only(
                              left: width / 25, right: width / 25),
                          child: Row(
                            // mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                LanguageEn.myPortfolio,
                                style: TextStyle(
                                    color: notifier.getwihitecolor,
                                    fontSize: 15.sp,
                                    fontFamily: 'Gilroy_Medium'),
                              ),
                              const Spacer(),
                              Image.asset("assets/images/statistics.png",
                                  height: height / 25),
                              SizedBox(width: width / 100),
                              // Stack(
                              //   children: [
                              //     Opacity(
                              //       opacity: 0.15,
                              //       child: Container(
                              //         height: height / 20,
                              //         width: width / 3.7,
                              //         decoration: BoxDecoration(
                              //           color: notifier.getwihitecolor,
                              //           borderRadius: BorderRadius.all(
                              //             Radius.circular(10.sp),
                              //           ),
                              //         ),
                              //       ),
                              //     ),
                              //     // Row(
                              //     //   children: [
                              //     //     Padding(
                              //     //       padding: EdgeInsets.only(left: width / 26),
                              //     //       child: DropdownButton(
                              //     //         dropdownColor: notifier.getbluecolor,
                              //     //         underline: const SizedBox(),
                              //     //         value: dropdownvalue,
                              //     //         icon: Icon(
                              //     //             Icons.keyboard_arrow_down_rounded,
                              //     //             color: notifier.getwihitecolor,
                              //     //             size: 17.sp),
                              //     //         items: items.map((String items) {
                              //     //           return DropdownMenuItem(
                              //     //             value: items,
                              //     //             child: Text(
                              //     //               items,
                              //     //               style: TextStyle(
                              //     //                   fontFamily: 'Gilroy_Medium',
                              //     //                   color: notifier.getwihitecolor,
                              //     //                   fontSize: 11.sp),
                              //     //             ),
                              //     //           );
                              //     //         }).toList(),
                              //     //         onChanged: (String? newValue) {
                              //     //           setState(() {
                              //     //             dropdownvalue = newValue!;
                              //     //           });
                              //     //         },
                              //     //       ),
                              //     //     ),
                              //     //   ],
                              //     // ),
                              //   ],
                              // )
                            ],
                          ),
                        ),
                        // Padding(
                        //   padding: EdgeInsets.only(left: width / 25),
                        //   child: Row(
                        //     children: [
                        //       Text(
                        //         "\$24,320+",
                        //         style: TextStyle(
                        //             color: notifier.getwihitecolor,
                        //             fontSize: 20.sp,
                        //             fontFamily: 'Gilroy_Bold'),
                        //       ),
                        //     ],
                        //   ),
                        // ),
                        // Expanded(
                        //   child: Container(
                        //     child: LineChart(
                        //       mainData(notifier.getwihitecolor, 0.5),
                        //     ),
                        //   ),
                        // ),
                        // InkWell(
                        //   onHover: (value) {
                        //     setState(() {
                        //       touchedIndex = -1;
                        //     });
                        //   },
                        //   hoverColor: Colors.transparent,
                        //   onTap: () {},
                        //   child: Row(
                        //     children: const <Widget>[],
                        //   ),
                        // ),
                        SizedBox(height: height / 100),
                        Row(
                          children: [
                            SizedBox(width: width / 27),
                            Text(
                              "\$34,010.00",
                              style: TextStyle(
                                  color: notifier.getwihitecolor,
                                  fontSize: 20.sp,
                                  fontFamily: 'Gilroy_Bold'),
                            ),
                            const Spacer(),
                            Text(
                              "+2,5%",
                              style: TextStyle(
                                  fontSize: 13.sp,
                                  fontFamily: 'Gilroy_Medium',
                                  color: notifier.getwihitecolor),
                            ),
                            SizedBox(width: width / 23),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              Column(
                children: [
                  SizedBox(height: height / 7.3),
                  Row(
                    children: [
                      SizedBox(width: width / 10),
                      GestureDetector(
                          onTap: () {
                            Get.to(const StockExchange());
                          },
                          child: cardbutton(
                              "assets/images/deposit.png", LanguageEn.deposit)),
                      SizedBox(width: width / 50),
                      GestureDetector(
                          onTap: () {
                            Get.to(const StockExchange());
                          },
                          child: cardbutton("assets/images/withdraw.png",
                              LanguageEn.withdraw)),
                    ],
                  ),
                ],
              )
            ],
          ),
        ],
      ),
    );
  }

  Widget cardbutton(image, name) {
    return Card(
      shape: RoundedRectangleBorder(
        side: const BorderSide(color: Colors.white70, width: 1),
        borderRadius: BorderRadius.circular(12.sp),
      ),
      elevation: 1,
      child: Container(
        color: Colors.transparent,
        height: height / 18,
        width: width / 2.7,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Image.asset(image, height: height / 35, color: Colors.black),
            SizedBox(width: width / 70),
            Text(
              name,
              style: TextStyle(
                  color: Colors.black,
                  fontSize: 14.sp,
                  fontFamily: 'Gilroy_Medium'),
            )
          ],
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/stock_exchange_tabs/notification.dart';
import 'package:trovo_wallet/button_tabs/chart.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';

import '../../../Custom_BlocObserver/custtom_slock_list/custtom_slock_list.dart';
import '../../../Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import '../../../utils/medeiaqury/medeiaqury.dart';

class SearchView extends StatefulWidget {
  const SearchView({Key? key}) : super(key: key);

  @override
  State<SearchView> createState() => _SearchViewState();
}

class _SearchViewState extends State<SearchView> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(70.sp),
          // here the desired height
          child: AppBar(
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            title: Text(
              LanguageEn.search,
              style: TextStyle(
                  fontSize: 20.sp,
                  color: notifier.getblck,
                  fontFamily: 'Gilroy_Bold'),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.only(left: 10, right: 10),
            child: Column(
              children: [
                SizedBox(height: height / 80),
                Center(
                  child: Customtextfild.textField(
                      "Search  assets, nfts...",
                      notifier.getbluecolor,
                      Icons.search_rounded,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      notifier.getgrey,
                      50.sp,
                      310.sp),
                ),
                SizedBox(height: height / 20),
                Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Center(
                      child: Text(
                        'This is search view...',
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.bold,
                          fontFamily: fontsemibold,
                          color: notifier.getblck,
                        ),
                      ),
                    )
                  ],
                )
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget marketrate(txt, rate, updown) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          txt,
          style: TextStyle(
              color: notifier.getgrey,
              fontFamily: 'Gilroy_Medium',
              fontSize: 13.sp),
        ),
        SizedBox(height: height / 200),
        Text(
          rate,
          style: TextStyle(
              color: notifier.getblck,
              fontFamily: 'Gilroy_Bold',
              fontSize: 14.sp),
        ),
        SizedBox(height: height / 200),
        Text(
          updown,
          style: TextStyle(
              color: const Color(0xff22C36B),
              fontFamily: 'Gilroy_Bold',
              fontSize: 14.sp),
        )
      ],
    );
  }
}

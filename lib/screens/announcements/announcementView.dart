import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:intl/intl.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/announcement.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';

class AnnouncementView extends StatelessWidget {
  AnnouncementView({Key? key}) : super(key: key);
  late DataProvider appState;
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    Announcement announcement =
        appState.viewData![AnnouncementViewPageConfig.key];

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
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Container(
                width: width / 1.2,
                child: Text(
                  announcement.title!,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontFamily: fontsemibold,
                      fontSize: 17.sp),
                ),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(
                    vertical: 15.0, horizontal: 25.0),
                child: Text(
                  announcement.message!,
                  textAlign: TextAlign.justify,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontFamily: fontbody,
                      fontSize: 15.sp),
                ),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(
                    vertical: 15.0, horizontal: 25.0),
                child: Text(
                  DateFormat('MMMM dd, yyyy hh:mm a').format(
                    announcement.createdAt!,
                  ),
                  // textAlign: TextAlign.justify,
                  style: TextStyle(
                      color: notifier.getblck,
                      fontFamily: fontbody,
                      fontSize: 15.sp),
                ),
              ),
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }
}

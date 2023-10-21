import 'package:flutter/material.dart';
import 'package:flutter_html/flutter_html.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/announcement.dart';
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
          child: CustomAppBar(
            context,
            notifier.getwihitecolor,
            announcement.title!,
            notifier.getbluewhitecolor,
            height: height / 15,
          ).getBar(),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 30.0),
                child: Html(data: announcement.message!, style: {
                  "*": Style(
                      color: notifier.getbluewhitecolor,
                      fontSize: FontSize.large,
                      lineHeight: LineHeight.number(1.2),
                      wordSpacing: 1.2,
                      textAlign: TextAlign.justify),
                  "h1, h2, h3, h4": Style(
                    fontFamily: fontsemibold,
                    fontSize: FontSize.large,
                  ),
                }),
              ),
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }
}

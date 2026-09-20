import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/models/announcement.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/storage/store.dart';
import 'package:trovo_app/widgets/utilities.dart';

import '../../../utils/medeiaqury/medeiaqury.dart';

class AnnouncementsView extends StatefulWidget {
  const AnnouncementsView({Key? key}) : super(key: key);

  @override
  State<AnnouncementsView> createState() => _AnnouncementsViewState();
}

class _AnnouncementsViewState extends State<AnnouncementsView> {
  late ColorNotifier notifier;
  late DataProvider appState;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);

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
            centerTitle: true,
            title: Text(
              "announcements".tr(),
              style: TextStyle(
                fontSize: 20.sp,
                color: notifier.getblck,
                fontFamily: 'Gilroy_Bold',
              ),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 80),
              FutureBuilder<dynamic>(
                future: StoreData().storeGetData('announcements'),
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return Container(
                      height: height / 1.5,
                      width: width,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          CircularProgressIndicator(
                            backgroundColor: notifier.getbluecolor,
                            valueColor: new AlwaysStoppedAnimation<Color>(
                              notifier.getgreencolor,
                            ),
                            strokeWidth: 3.0,
                          ),
                        ],
                      ),
                    );
                  } else if (snapshot.connectionState == ConnectionState.done) {
                    if (snapshot.hasError) {
                      return Container(
                        height: height / 1.5,
                        child: Padding(
                          padding: const EdgeInsets.all(8.0),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                "somethingwentwrong".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getblck,
                                  fontFamily: fontbody,
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    } else if (snapshot.hasData) {
                      var announcements = Announcement().deserializeJsonList(
                        snapshot.data,
                      );
                      return Column(
                        children: [
                          for (var announcement in announcements) ...[
                            GestureDetector(
                              onTap: () {
                                appState.viewData![AnnouncementViewPageConfig
                                        .key] =
                                    announcement;
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: AnnouncementViewPageConfig,
                                );

                                announcement.isViewed = true;
                                StoreData().storeInsertData(
                                  'lastNotificationViewDate',
                                  DateTime.now().toIso8601String(),
                                );
                                StoreData().storeInsertData(
                                  'announcements',
                                  Announcement().toJSONEncodableList(
                                    announcements,
                                  ),
                                );
                              },
                              child: Card(
                                elevation: notifier.isDark ? 0 : 5,
                                shadowColor: Colors.black,
                                color: notifier.gettilewihitecolor,
                                margin: EdgeInsets.symmetric(
                                  vertical: 10,
                                  horizontal: 20,
                                ),
                                shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(15.0),
                                ),
                                child: Padding(
                                  padding: const EdgeInsets.symmetric(
                                    vertical: 8.0,
                                  ),
                                  child: ListTile(
                                    title: Row(
                                      children: [
                                        Icon(
                                          Icons.circle_sharp,
                                          color: announcement.isViewed
                                              ? Colors.grey
                                              : Colors.green,
                                          size: 15,
                                        ),
                                        SizedBox(width: 10),
                                        Column(
                                          crossAxisAlignment:
                                              CrossAxisAlignment.start,
                                          children: [
                                            Container(
                                              width: width / 1.4,
                                              child: Text(
                                                announcement.title!,
                                                style: TextStyle(
                                                  fontSize: 12,
                                                  fontFamily: fontsemibold,
                                                  color: notifier.getblck,
                                                ),
                                              ),
                                            ),
                                            SizedBox(height: height / 90),
                                            Padding(
                                              padding:
                                                  const EdgeInsets.fromLTRB(
                                                    0,
                                                    3.0,
                                                    0,
                                                    0,
                                                  ),
                                              child: Container(
                                                width: width / 1.4,
                                                child: Text(
                                                  truncate(
                                                    announcement.message!,
                                                    length: 60,
                                                  ),
                                                  style: TextStyle(
                                                    fontSize: 9,
                                                    fontFamily: fontbody,
                                                    color: notifier.getblck,
                                                  ),
                                                ),
                                              ),
                                            ),
                                            SizedBox(height: height / 90),
                                            Padding(
                                              padding:
                                                  const EdgeInsets.fromLTRB(
                                                    0,
                                                    3.0,
                                                    0,
                                                    0,
                                                  ),
                                              child: Container(
                                                width: width / 1.4,
                                                child: Text(
                                                  DateFormat(
                                                    'MMMM dd, yyyy hh:mm a',
                                                  ).format(
                                                    announcement.createdAt!,
                                                  ),
                                                  style: TextStyle(
                                                    fontSize: 9,
                                                    fontFamily: fontbody,
                                                    color: notifier.getblck,
                                                  ),
                                                ),
                                              ),
                                            ),
                                          ],
                                        ),
                                      ],
                                    ),
                                  ),
                                ),
                              ),
                            ),
                          ],
                          SizedBox(height: height / 10),
                        ],
                      );
                    }
                  }

                  return Container(
                    height: height / 1.5,
                    width: width,
                    child: Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            'No announcement here',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.bold,
                              fontFamily: fontsemibold,
                              color: notifier.getblck,
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                },
              ),
            ],
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
            fontSize: 13.sp,
          ),
        ),
        SizedBox(height: height / 200),
        Text(
          rate,
          style: TextStyle(
            color: notifier.getblck,
            fontFamily: 'Gilroy_Bold',
            fontSize: 14.sp,
          ),
        ),
        SizedBox(height: height / 200),
        Text(
          updown,
          style: TextStyle(
            color: const Color(0xff22C36B),
            fontFamily: 'Gilroy_Bold',
            fontSize: 14.sp,
          ),
        ),
      ],
    );
  }
}

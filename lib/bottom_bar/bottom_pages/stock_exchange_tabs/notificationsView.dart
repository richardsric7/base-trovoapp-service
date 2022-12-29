import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import '../../../utils/medeiaqury/medeiaqury.dart';

class NotificationsView extends StatefulWidget {
  const NotificationsView({Key? key}) : super(key: key);

  @override
  State<NotificationsView> createState() => _NotificationsViewState();
}

class _NotificationsViewState extends State<NotificationsView> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Future<Map> notificationsList;

  @override
  void initState() {
    super.initState();
    notificationsList = fetchNotifications();
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
              'Notifications',
              style: TextStyle(
                  fontSize: 20.sp,
                  color: notifier.getblck,
                  fontFamily: 'Gilroy_Bold'),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 80),
              FutureBuilder<Map>(
                future: notificationsList,
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
                              marketrate('Notification', '200', 'updown'),
                              Text(
                                LanguageEn.somethingwentwrong,
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                    fontSize: 16,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              ElevatedButton(
                                onPressed: () {
                                  setState(() {
                                    notificationsList = fetchNotifications();
                                  });
                                },
                                style: ButtonStyle(
                                  backgroundColor:
                                      MaterialStateProperty.all<Color>(
                                          notifier.getbluecolor!),
                                ),
                                child: Text(
                                  LanguageEn.retry,
                                  style: TextStyle(
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    } else if (snapshot.hasData) {
                      return marketrate('Notification', '200', 'updown');
                    }
                  }

                  return Text(
                    'You have no notifications yet',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
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

  Future<Map> fetchNotifications() async {
    try {
      print('fetching announcements');
      var uri = '/v1/announcements';

      Map responseData = await makeUnSecuredGetRequest(
        Uri.encodeFull(uri),
      );

      print('response: ${responseData}');

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class MessageSupport extends StatefulWidget {
  const MessageSupport({Key? key}) : super(key: key);

  @override
  State<MessageSupport> createState() => _MessageSupportState();
}

class _MessageSupportState extends State<MessageSupport> {
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
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(context, notifier.getwihitecolor,
            "Message Support", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Container(
                color: Colors.transparent,
                height: height / 1.25,
                child: SingleChildScrollView(
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 50,
                      ),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        child: chat(
                            "You can also turn chat features off in the Messages deactivation web portal"),
                      ),
                      SizedBox(height: height / 15),
                      Padding(
                        padding: const EdgeInsets.only(right: 20),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.end,
                          children: [
                            Container(
                              height: height / 9,
                              decoration: BoxDecoration(
                                color: notifier.getbluecolor,
                                borderRadius: BorderRadius.all(
                                  Radius.circular(15.sp),
                                ),
                              ),
                              child: Padding(
                                padding:
                                    const EdgeInsets.symmetric(horizontal: 10),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    SizedBox(height: height / 80),
                                    Text(
                                      "Why are my card purchases\ntiming out?",
                                      style: TextStyle(
                                          fontFamily: 'Gilroy_Medium',
                                          fontSize: 13.sp,
                                          color: notifier.getwihitecolor),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                            SizedBox(width: width / 50),
                            Image.asset(
                              "assets/images/profile.png",
                              height: height / 15,
                            ),
                          ],
                        ),
                      ),
                      SizedBox(height: height / 15),
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20),
                        child: chat(
                            "If your carrier supports chat features but your device isn’t automatically set up for chat features, you may get a notification to “Do more with Messages.” If you get this notification:"),
                      ),
                    ],
                  ),
                ),
              ),
              SizedBox(height: height / 80),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  Container(
                    color: Colors.transparent,
                    height: height / 17,
                    width: width / 1.3,
                    child: TextField(
                      style: TextStyle(color: notifier.getblck),
                      textAlign: TextAlign.left,
                      decoration: InputDecoration(
                          hintText: "Write a reply...",
                          hintStyle: TextStyle(
                              color: notifier.getgrey, fontSize: 15.sp),
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(30.sp),
                            borderSide: const BorderSide(
                              width: 0,
                              style: BorderStyle.none,
                            ),
                          ),
                          filled: true,
                          fillColor: notifier.getwihitecolor),
                    ),
                  ),
                  Image.asset("assets/images/GIF.png", height: height / 40),
                  Image.asset("assets/images/Emoticon.png",
                      height: height / 40),
                  Image.asset("assets/images/Attach.png", height: height / 40),
                  SizedBox(width: width / 100),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget chat(txt) {
    return Row(
      children: [
        Image.asset(
          "assets/images/profile.png",
          height: height / 15,
        ),
        SizedBox(width: width / 50),
        Container(
          width: width / 1.5,
          decoration: BoxDecoration(
            color: notifier.getwihitecolor,
            borderRadius: const BorderRadius.all(
              Radius.circular(15),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                SizedBox(height: height / 100),
                Text(
                  txt,
                  style: TextStyle(
                      fontFamily: 'Gilroy_Medium',
                      fontSize: 14.sp,
                      color: notifier.getblck),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }
}

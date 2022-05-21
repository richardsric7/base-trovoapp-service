import 'package:accordion/accordion.dart';
import 'package:accordion/controllers.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class FAQ extends StatefulWidget {
  const FAQ({Key? key}) : super(key: key);

  @override
  State<FAQ> createState() => _FAQState();
}

class _FAQState extends State<FAQ> {
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

  final _loremIpsum =
      "Open the gocrypto app to get started and follow the\nsteps. gocrypto doesn’t charge a fee to create or\nmaintain your gocrypto account.";
  final _contentStyle = const TextStyle(
      color: Color(0xff999999), fontSize: 14, fontWeight: FontWeight.normal);

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(notifier.getwihitecolor, "FAQ", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 30),
              Center(
                child: Customtextfild.textField(
                    "Search topics or questions",
                    notifier.getbluecolor,
                    Icons.search_rounded,
                    notifier.getgrey,
                    notifier.getblck,
                    notifier.getblck,
                    notifier.getgrey,
                    45.sp,
                    300.sp),
              ),
              SizedBox(height: height / 17),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  GestureDetector(
                      onTap: () {
                        Get.to(() => const Pages(""));
                      },
                      child: faqtype(
                          "assets/images/Getting.png", "Getting\nStarted")),
                  GestureDetector(
                      onTap: () {
                        Get.to(() => const Pages(""));
                      },
                      child: faqtype(
                          "assets/images/Prices.png", "How to\nInvest")),
                  GestureDetector(
                      onTap: () {
                        Get.to(() => const Pages(""));
                      },
                      child: faqtype(
                          "assets/images/Card-payment.png", "Payment\nMethod")),
                ],
              ),
              SizedBox(height: height / 10),
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Text(
                    "Top Questions",
                    style: TextStyle(
                        color: notifier.getblck,
                        fontSize: 16.sp,
                        fontFamily: 'Gilroy_Bold'),
                  ),
                ],
              ),
              SizedBox(height: height / 50),
              Accordion(
                disableScrolling: true,
                flipRightIconIfOpen: true,
                contentVerticalPadding: 0,
                scrollIntoViewOfItems: ScrollIntoViewOfItems.fast,
                contentBorderColor: Colors.transparent,
                maxOpenSections: 1,
                headerBackgroundColorOpened: Colors.black54,
                headerPadding:
                    const EdgeInsets.symmetric(vertical: 7, horizontal: 15),
                children: [
                  AccordionSection(
                    sectionClosingHapticFeedback: SectionHapticFeedback.light,
                    headerBackgroundColor: notifier.getwihitecolor,
                    header: Text(
                      'How to create a gocrypto account?',
                      style: TextStyle(
                          color: notifier.getblck,
                          fontSize: 15,
                          fontWeight: FontWeight.bold),
                    ),
                    content: Text(_loremIpsum, style: _contentStyle),
                    contentHorizontalPadding: 20,
                    contentBorderWidth: 1,
                  ),
                  AccordionSection(
                    flipRightIconIfOpen: true,
                    headerBackgroundColor: notifier.getwihitecolor,
                    header: Text(
                      'How to add a payment method?',
                      style: TextStyle(
                          color: notifier.getblck,
                          fontSize: 15,
                          fontWeight: FontWeight.bold),
                    ),
                    content: Text(_loremIpsum, style: _contentStyle),
                    contentHorizontalPadding: 20,
                    contentBorderWidth: 1,
                  ),
                  AccordionSection(
                    flipRightIconIfOpen: true,
                    headerBackgroundColor: notifier.getwihitecolor,
                    header: Text(
                      'How to buy stocks?',
                      style: TextStyle(
                          color: notifier.getblck,
                          fontSize: 15,
                          fontWeight: FontWeight.bold),
                    ),
                    content: Text(_loremIpsum, style: _contentStyle),
                    contentHorizontalPadding: 20,
                    contentBorderWidth: 1,
                  ),
                ],
              ),
              // type("How to buy stocks?", Icons.add),
            ],
          ),
        ),
      ),
    );
  }

  Widget type(txt, icon) {
    return Row(
      children: [
        SizedBox(width: width / 15),
        Text(
          txt,
          style: TextStyle(
              color: notifier.getblck,
              fontFamily: 'Gilroy_Bold',
              fontSize: 14.sp),
        ),
        const Spacer(),
        Icon(icon, color: notifier.getbluecolor, size: 20.sp),
        SizedBox(width: width / 15),
      ],
    );
  }

  Widget faqtype(image, txt) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Image.asset(image, height: height / 25),
        SizedBox(height: height / 100),
        Text(
          txt,
          style: TextStyle(color: notifier.getblck, fontSize: 13.sp),
        )
      ],
    );
  }
}

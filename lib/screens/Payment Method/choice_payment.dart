import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/screens/Payment%20Method/confirmorder.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../card_type/custtomcsrdtype.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ChoicePayment extends StatefulWidget {
  const ChoicePayment({Key? key}) : super(key: key);

  @override
  State<ChoicePayment> createState() => _ChoicePaymentState();
}

class _ChoicePaymentState extends State<ChoicePayment> {
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
        appBar: CustomAppBar(context, notifier.getwihitecolor, "Payment Method",
            notifier.getblck,
            height: height / 15),
      ),
    );
  }
}

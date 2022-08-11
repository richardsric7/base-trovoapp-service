import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class PaymentMethod extends StatefulWidget {
  const PaymentMethod({Key? key}) : super(key: key);

  @override
  State<PaymentMethod> createState() => _PaymentMethodState();
}

class _PaymentMethodState extends State<PaymentMethod> {
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

  Widget banktype(txt) {
    return Center(
      child: Card(
        color: notifier.getfavorites,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.all(
            Radius.circular(10.sp),
          ),
        ),
        child: Container(
          color: Colors.transparent,
          height: height / 13,
          width: width / 1.1,
          child: Row(
            children: [
              SizedBox(width: width / 20),
              Text(
                txt,
                style: TextStyle(
                    color: notifier.getbluecolor,
                    fontFamily: 'Gilroy_Bold',
                    fontSize: 13.sp),
              ),
              const Spacer(),
              Icon(
                Icons.arrow_forward_ios,
                size: 14.sp,
                color: notifier.getgrey,
              ),
              SizedBox(width: width / 20),
            ],
          ),
        ),
      ),
    );
  }

  Widget selectpaymenttype(image) {
    return Container(
      height: height / 20,
      width: width / 4.5,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.all(
          Radius.circular(13.sp),
        ),
        border: Border.all(color: notifier.getgrey),
      ),
      child: Center(
        child: Image.asset(image, height: height / 55),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/button/custtom_button.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/screens/Auth/verifyyouridentity.dart';
import 'package:gocrypto/utils/enstring.dart';
import 'package:gocrypto/utils/medeiaqury/medeiaqury.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class Emailpassword extends StatefulWidget {
  const Emailpassword({Key? key}) : super(key: key);

  @override
  State<Emailpassword> createState() => _EmailpasswordState();
}

class _EmailpasswordState extends State<Emailpassword> {
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

  String dropdownvalue = '+91';

  // List of items in our dropdown menu
  var items = [
    '+61',
    '+91',
    '+92',
    '+152',
    '+139',
  ];

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 20),
              Center(
                  child: Image.asset("assets/images/lock.png",
                      height: height / 3.9)),
              SizedBox(height: height / 25),
              Text(
                LanguageEn.resetpass,
                style: TextStyle(
                    color: notifier.getblck,
                    fontSize: 22.sp,
                    fontFamily: 'Gilroy_Bold'),
              ),
              SizedBox(height: height / 100),
              Text(
                LanguageEn.enteranemailadress,
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getgrey,
                    fontSize: 15.sp,
                    fontFamily: 'Gilroy_Medium'),
              ),
              SizedBox(height: height / 30),
              verification(),
              SizedBox(height: height / 4.7),
              GestureDetector(
                  onTap: () {
                    Get.to(const VerifyYourIdentity());
                  },
                  child: Button(LanguageEn.continuee, notifier.getbluecolor,
                      notifier.getwihitecolor))
            ],
          ),
        ),
      ),
    );
  }

  Widget verification() {
    return Container(
      height: height / 15,
      width: width / 1.2,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.all(
          Radius.circular(15.sp),
        ),
        border: Border.all(
          color: notifier.getbluecolor,
        ),
      ),
      child: Row(
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 10),
            child: DropdownButton(
              dropdownColor: notifier.getwihitecolor,
              underline: const SizedBox(),
              // Initial Value
              value: dropdownvalue,

              // Down Arrow Icon
              icon: Icon(
                Icons.arrow_drop_down,
                color: notifier.getblck,
              ),

              // Array list of items
              items: items.map((String items) {
                return DropdownMenuItem(
                  value: items,
                  child: Text(
                    items,
                    style: TextStyle(color: notifier.getblck),
                  ),
                );
              }).toList(),
              // After selecting the desired option,it will
              // change button value to selected value
              onChanged: (String? newValue) {
                setState(() {
                  dropdownvalue = newValue!;
                });
              },
            ),
          ),
          SizedBox(width: width / 100),
          Container(
            height: height / 22,
            width: width / 250,
            color: notifier.getgrey,
          ),
          SizedBox(width: width / 22),
          Container(
            color: Colors.transparent,
            height: 25.h,
            width: 200.w,
            child: TextField(
              style: TextStyle(color: notifier.getblck),
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(border: InputBorder.none),
            ),
          )
        ],
      ),
    );
  }
}

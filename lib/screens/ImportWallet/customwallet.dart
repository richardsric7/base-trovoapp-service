import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/termsOfService.dart';

class CustomWallet extends StatefulWidget {
  const CustomWallet({Key? key}) : super(key: key);

  @override
  State<CustomWallet> createState() => _CustomWalletState();
}

class _CustomWalletState extends State<CustomWallet> {
  late ColorNotifier notifier;
  bool hasAgreed = false;
  bool showTermsError = false;
  bool usePassPhrase = false;
  late FocusNode passPhraseFocusNode;
  late FocusNode secretKeyFocusNode;
  final _formKey = GlobalKey<FormState>();
  String? email;
  String? username;
  String? passPhrase;
  String? secretKey;
  String? password;

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
    passPhraseFocusNode = FocusNode();
    secretKeyFocusNode = FocusNode();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: height / 10),
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          LanguageEn.createcustomwallet,
                          style: TextStyle(
                              color: notifier.getblck,
                              fontSize: 26.sp,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(height: height / 10),
                        // Email address
                        CustomTextFormField.textField(
                          LanguageEn.customwalletname,
                          notifier.getbluecolor,
                          Icons.email,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          validator: (value) {
                            var trimmedVal = value!.trim().replaceAll(' ', '');
                            if (trimmedVal.isEmpty) {
                              return LanguageEn.accountaliasoremailempty;
                            }

                            if (trimmedVal.length < 3) {
                              return LanguageEn.accountaliasoremailinvalid;
                            }
                          },
                          // onSaved: storeUsernameOrEmail,
                          keyboardtype: TextInputType.emailAddress,
                        ),
                        SizedBox(height: height / 40),
                        CustomTextFormField.textField(
                          LanguageEn.customwalletaddress,
                          notifier.getbluecolor,
                          Icons.email,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          validator: (value) {
                            var trimmedVal = value!.trim().replaceAll(' ', '');
                            if (trimmedVal.isEmpty) {
                              return LanguageEn.accountaliasoremailempty;
                            }

                            if (trimmedVal.length < 3) {
                              return LanguageEn.accountaliasoremailinvalid;
                            }
                          },
                          // onSaved: storeUsernameOrEmail,
                          keyboardtype: TextInputType.emailAddress,
                        ),
                      ],
                    ),
                  )
                ],
              ),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                notifier.getwihitecolor,
                onTap: () => {},
              ),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

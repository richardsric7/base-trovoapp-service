import 'dart:convert';
import 'dart:typed_data';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:image_picker/image_picker.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetTokenInformation extends StatefulWidget {
  const AssetTokenInformation({Key? key}) : super(key: key);

  @override
  State<AssetTokenInformation> createState() => _AssetTokenInformation();
}

class _AssetTokenInformation extends State<AssetTokenInformation>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  String? imageThumbnail;
  bool addAdditionalKyc = false;

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getCurrencies {
    List<DropdownMenuItem<String>> currencies = [];
    appState.fiatRate.forEach((key, value) {
      currencies.add(DropdownMenuItem(
          child: Text(
            key,
            overflow: TextOverflow.ellipsis,
          ),
          value: key));
    });
    return currencies;
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
    appState = Provider.of<DataProvider>(context, listen: true);
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SingleChildScrollView(
        child: Column(
          children: [
            CustomAppBar(
              context,
              notifier.getwihitecolor,
              'Asset Token Information',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Enter Asset Code',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Upload Asset Logo',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Upload Asset Logo',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            Row(
              children: [
                GestureDetector(
                  onTap: () {
                    getImage();
                  },
                  child: Column(
                    children: [
                      SizedBox(
                        height: height / 50,
                      ),
                      Padding(
                        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                        child: Container(
                          decoration: BoxDecoration(
                            border: Border.all(
                                color: notifier.getbluewhitecolor, width: 1),
                            borderRadius:
                                const BorderRadius.all(Radius.circular(15.0)),
                            color: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                          ),
                          child: Column(
                            children: [
                              SizedBox(
                                  width: width / 1.2,
                                  height: height / 6,
                                  child: Center(
                                    child: Wrap(
                                      alignment: WrapAlignment.center,
                                      children: [
                                        Text(
                                          'Browse file to select image.\nMax image size: 10MB',
                                          textAlign: TextAlign.center,
                                          style: TextStyle(
                                            color: notifier.getbluewhitecolor,
                                            fontFamily: fontsemibold,
                                            fontSize: 12.sp,
                                          ),
                                        ),
                                      ],
                                    ),
                                  )),
                              const SizedBox(height: 2),
                            ],
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            if (imageThumbnail != null) ...[
              GestureDetector(
                onTap: () {
                  getImage();
                },
                child: Padding(
                  padding: const EdgeInsets.all(15.0),
                  child: Image.memory(
                    base64Decode(imageThumbnail!),
                    width: width / 1.3,
                    height: height / 6,
                  ),
                ),
              ),
            ],
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Number of Token to be Issued',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Number to be issued',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Number of Token to be sold',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Number to be sold',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Total token held by Manager',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Total held by Manager',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Select Wallet to hold assets not up for sale',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'hint',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Or create a new wallet for this',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Asset Sales & and Pricing',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Price Per Token',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Price per token',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    '(Calculated from Asset value)',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Select Funding Method',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'e-Naira',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Text(
                          'Sales start from',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    CustomTextFormField.textField(
                      'start',
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      50.sp,
                      width / 2.7,
                      // controller: referrerController,
                      // validator: validateReferrer,
                      onSaved: (value) {},
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                        signed: true,
                      ),
                    ),
                  ],
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Text(
                          'Sales end on',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    CustomTextFormField.textField(
                      'ends',
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      50.sp,
                      width / 2.5,
                      // controller: referrerController,
                      // validator: validateReferrer,
                      onSaved: (value) {},
                      keyboardtype: TextInputType.numberWithOptions(
                        decimal: true,
                        signed: true,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Asset Wallets',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Minting Wallet',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Minting wallet',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Market Making Wallet',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Market Making Wallet',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Cap on Purchase',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
                Spacer(),
                Transform.scale(
                  scale: 0.7,
                  child: CupertinoSwitch(
                    activeColor: notifier.getgreencolor,
                    value: notifier.getIsDark,
                    onChanged: (val) async {
                      final prefs = await SharedPreferences.getInstance();
                      setState(() {
                        notifier.setIsDark = val;
                        prefs.setBool("setIsDark", val);
                      });
                    },
                  ),
                ),
              ],
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Container(
                    width: width / 1.2,
                    child: Text(
                      'Cap is the maximum amount of token a person is able to buy within the cap period, after which they can buy any amount they want or as allowed by regulation in their jurisdiction',
                      style: TextStyle(
                        fontSize: 9,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Cap Quantity (Quantity to be bought)',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Quantity',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Cap Duration',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'Days',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Proceed Payout',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Proceed Payout Cycle',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'Monthly',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Proceed Payout Currency',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'Naira',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Primary Buyer Requirement',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Container(
                    width: width / 1.2,
                    child: Text(
                      'Select Countries that are exempted from participation if any',
                      style: TextStyle(
                        fontSize: 9,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'Algeria',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                  color: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 20.0, vertical: 15.0),
                      child: Column(
                        children: [
                          Container(
                              width: width / 1.8,
                              child: Wrap(
                                alignment: WrapAlignment.center,
                                children: [
                                  if (true) ...[
                                    userItem('Algeria', () {},
                                        foreColor: wihitecolor,
                                        backColor: notifier.getbluebackcolor),
                                    userItem('China', () {},
                                        foreColor: wihitecolor,
                                        backColor: notifier.getbluebackcolor),
                                    userItem('India', () {},
                                        foreColor: wihitecolor,
                                        backColor: notifier.getbluebackcolor),
                                    userItem('Pakistan', () {},
                                        foreColor: wihitecolor,
                                        backColor: notifier.getbluebackcolor)
                                  ] else ...[
                                    Text(
                                      'Name of exempted countries appear here',
                                      textAlign: TextAlign.center,
                                      style: TextStyle(
                                          color: notifier.getbluewhitecolor,
                                          fontFamily: fontbody,
                                          fontSize: 15.sp),
                                    ),
                                  ]
                                ],
                              )),
                          SizedBox(height: 2),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
            SizedBox(height: height / 50),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Are there additional KYC Requirement?',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Column(
              children: [
                ListTile(
                  title: Text(
                    'Yes',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  leading: Radio(
                    value: addAdditionalKyc,
                    groupValue: true,
                    activeColor: notifier.getbluewhitecolor,
                    fillColor:
                        MaterialStateProperty.all(notifier.getbluewhitecolor),
                    onChanged: (value) {
                      setState(() {
                        addAdditionalKyc = true;
                        print('addAdditionalKyc: $value');
                      });
                    },
                  ),
                ),
                ListTile(
                  title: Text(
                    'No',
                    style: TextStyle(
                      fontSize: 15,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                  leading: Radio(
                    value: addAdditionalKyc,
                    groupValue: false,
                    fillColor:
                        MaterialStateProperty.all(notifier.getbluewhitecolor),
                    activeColor: notifier.getbluewhitecolor,
                    onChanged: (value) {
                      setState(() {
                        addAdditionalKyc = false;
                        print('addAdditionalKyc: $value');
                      });
                    },
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        'List the Requirements',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 70,
                ),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: multilineInput(
                        '',
                        notifier.getbluecolor,
                        notifier.getgrey,
                        notifier.getblck,
                        notifier.getgrey,
                        100.sp,
                        300.sp,
                        validator: (value) {},
                        onSaved: (value) {},
                        minLines: 3,
                        maxLines: null,
                        keyboardtype: TextInputType.multiline,
                      ),
                    ),
                  ],
                ),
              ],
            ),
            accreditInvestors(),
            SizedBox(
              height: height / 30,
            ),
            Button(
              'Save',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {},
            ),
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
  }

  Widget accreditInvestors() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(5.sp),
              ),
            ),
            activeColor: notifier.isDark
                ? notifier.getbluecolor50
                : notifier.getbluecolor90,
            side: BorderSide(
              color: notifier.isDark
                  ? notifier.getbluecolor50
                  : notifier.getbluecolor90,
            ),
            value: true,
            onChanged: (bool? value) {
              setState(() {});
            },
          ),
        ),
        Container(
          width: width / 1.2,
          child: Text(
            'Investor must be accredited',
            overflow: TextOverflow.visible,
            style: TextStyle(
                fontSize: 15,
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody),
          ),
        ),
      ],
    );
  }

  Widget multilineInput(
    labletext,
    focuscolor,
    lablecolor,
    textcolor,
    bordercolor,
    h,
    w, {
    onChanged,
    maxLength,
    minLines,
    maxLines,
    validator,
    onSaved,
    keyboardtype,
    focusNode,
  }) {
    return Container(
      height: h,
      width: w,
      child: TextFormField(
        focusNode: focusNode,
        maxLength: maxLength,
        minLines: minLines,
        maxLines: maxLines,
        style: TextStyle(color: textcolor, fontFamily: fontbody),
        cursorColor: lablecolor,
        onChanged: onChanged,
        decoration: InputDecoration(
          hintText: labletext,
          hintStyle: TextStyle(color: lablecolor),
          disabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          labelStyle: TextStyle(color: lablecolor),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: bordercolor, width: 1),
            borderRadius: BorderRadius.circular(15),
          ),
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: focuscolor, width: 1),
            borderRadius: BorderRadius.circular(15),
          ),
        ),
        keyboardType: keyboardtype,
        validator: validator,
        onSaved: onSaved,
      ),
    );
  }

  Future<void> getImage() async {
    var image = await ImagePicker().pickImage(source: ImageSource.gallery);
    if (image != null) {
      var imageBase64Uncompressed = await getBase64Image(image);
      await StoreData().storeInsertData('image', imageBase64Uncompressed);
      setState(() {
        imageThumbnail = imageBase64Uncompressed;
      });
    }
  }

  Future<dynamic> getBase64Image(XFile image) async {
    //
    List<int> imageBytes = await image.readAsBytes();
    String imageB64 = base64Encode(imageBytes);
    return imageB64;
  }

  Uint8List getBase64Decode(String image) {
    Uint8List imageString = base64Decode(image);
    return imageString;
  }
}

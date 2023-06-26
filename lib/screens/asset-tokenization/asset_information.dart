import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetInformation extends StatefulWidget {
  const AssetInformation({Key? key}) : super(key: key);

  @override
  State<AssetInformation> createState() => _AssetInformation();
}

class _AssetInformation extends State<AssetInformation>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  List<String> options = [
    'Insurance',
    'Alarm System',
    'Surveillance System',
    'Physical Security',
    'Inspection and Management',
    'Other',
  ];

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getOptions {
    List<DropdownMenuItem<String>> myOptions = [];
    options.forEach((value) {
      myOptions.add(DropdownMenuItem(
          child: Text(
            value,
            overflow: TextOverflow.ellipsis,
          ),
          value: value));
    });
    return myOptions;
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
              'Asset Information',
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
                    'Asset Description',
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
                    validator: (value) {
                      if (value.isEmpty) {
                        return "enterpassphraseempty".tr();
                      }
                    },
                    onSaved: (value) {},
                    minLines: 3,
                    maxLines: null,
                    keyboardtype: TextInputType.multiline,
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
                    'Enter Asset Physical Address',
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
                  child: CustomTextFormField.textField(
                    'Asset Physical Address',
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
                    'Enter Asset Google Map Coordinates',
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
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                CustomTextFormField.textField(
                  'Latitute',
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
                CustomTextFormField.textField(
                  'Longitude',
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
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Ownership',
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
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                CheckItem(
                  'Direct Ownership',
                  () {},
                  borderColor: notifier.getbluewhitecolor,
                  foreColor: notifier.getwihitecolor,
                  backColor: notifier.getbluewhitecolor,
                ),
                CheckItem(
                  'Third Party',
                  () {},
                  borderColor: notifier.getbluewhitecolor,
                  foreColor: notifier.getbluewhitecolor,
                  backColor: notifier.getwihitecolor,
                )
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Container(
                    width: width / 1.17,
                    child: Text(
                      'Select what best describes the third party owner',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                CheckItem(
                  'Individual',
                  () {},
                  borderColor: notifier.getbluewhitecolor,
                  foreColor: notifier.getwihitecolor,
                  backColor: notifier.getbluewhitecolor,
                ),
                CheckItem(
                  'Organization',
                  () {},
                  borderColor: notifier.getbluewhitecolor,
                  foreColor: notifier.getbluewhitecolor,
                  backColor: notifier.getwihitecolor,
                )
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
                    'Name of Organization',
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
                    'Name of Organization',
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
                    'Address of Organization',
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
                    'Address of Organization',
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
                    'Asset Custodian',
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
                  child: Text(
                    'Please provide the following information about the asset custodian ',
                    style: TextStyle(
                      fontSize: 9,
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
                    'Who is the Asset Custodian?',
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
                    'Asset Custodian',
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
                    'What is their address?',
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
                    'Address',
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
                    'Asset Manager',
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
                  child: Text(
                    'Please provide the following information about the asset manager ',
                    style: TextStyle(
                      fontSize: 9,
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
                    'Who is the Asset Manager?',
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
                    'Asset Manager',
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
                    'What is their address?',
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
                    'Address',
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
                    'Asset Value',
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
                  child: Text(
                    'Please provide the following information about the asset value ',
                    style: TextStyle(
                      fontSize: 9,
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
                    'What is the current value of the Asset?',
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
                    'Current Value of Asset',
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
                    'Percentage of Asset to be Tokenized? (1%-100%)',
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
                    'Percentage to be Tokenized',
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
                    'Value of tokenized asset',
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
                    'Value of tokenized asset',
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
                    '(asset value x asset percentage)',
                    style: TextStyle(
                      fontSize: 9,
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
                    'What is the real asset protection in place?',
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
                  child: Text(
                    'Select one or more protection options',
                    style: TextStyle(
                      fontSize: 9,
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
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getOptions,
                null,
                'Insurance',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: Row(
                children: [
                  userItem(
                    'Insurance',
                    () {},
                    foreColor: notifier.getwihitecolor,
                    backColor: notifier.getbluewhitecolor,
                  ),
                  userItem(
                    'Alarm System',
                    () {},
                    foreColor: notifier.getwihitecolor,
                    backColor: notifier.getbluewhitecolor,
                  )
                ],
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
                    'Insurance Company Name',
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
                    'Company Name',
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
                    'Insurance Policy Number',
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
                    'Insurance Policy Number',
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
                    'Insurance Policy Holder',
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
                    'Insurance Policy Holder',
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
                    'Percentage Value of Insurance',
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
                    'Percentage Value of Insurance',
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
            confirmLiensAndEncumbrance(),
            SizedBox(
              height: height / 30,
            ),
            Button(
              'Save',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                Navigator.of(context).pop();
              },
            ),
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
  }

  Widget confirmLiensAndEncumbrance() {
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
            'I confirm that this asset is completely free of all liens and encumbrance',
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
}

Widget CheckItem(
  String name,
  void Function()? onClick, {
  required Color backColor,
  required Color foreColor,
  required Color borderColor,
  double? fontSize: 15,
}) {
  return Padding(
    padding: const EdgeInsets.all(3.0),
    child: Container(
      decoration: BoxDecoration(
          border: Border.all(color: borderColor, width: 1),
          borderRadius: const BorderRadius.all(Radius.circular(10.0)),
          color: backColor),
      child: Padding(
        padding: const EdgeInsets.all(15.0),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                  color: foreColor, fontFamily: fontbody, fontSize: fontSize),
            ),
            SizedBox(
              width: width / 70,
            ),
          ],
        ),
      ),
    ),
  );
}

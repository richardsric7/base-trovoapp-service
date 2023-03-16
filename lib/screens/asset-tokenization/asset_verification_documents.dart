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
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetVerificationDocuments extends StatefulWidget {
  const AssetVerificationDocuments({Key? key}) : super(key: key);

  @override
  State<AssetVerificationDocuments> createState() =>
      _AssetVerificationDocuments();
}

class _AssetVerificationDocuments extends State<AssetVerificationDocuments>
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
              'Asset Verification Documents',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(
              height: height / 30,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Proof of Asset Existence',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 70,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          'Purchase receipt',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Text(
                          'Atlantis  Invoice PDF',
                          style: TextStyle(
                            fontSize: 12,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Icon(
                          CupertinoIcons.delete,
                          size: 20,
                        ),
                      ],
                    ),
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Proof of Asset Existence',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Proof of Ownership',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Proof of Ownership',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Asset Status Verification',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Asset Status Verification',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Asset Custodian Agreement',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Asset Custodian Agreement',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Proof of Asset Manager',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Proof of Asset Manager',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Asset Protection Document',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Asset Protection Document',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            myContainer(
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.start,
                    children: [
                      Text(
                        'Asset Valuation Certificate ',
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.27,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        TextButton(
                          onPressed: () => showDocumentUploadPopup(
                            context,
                            'Asset Protection Document',
                            onDone: () {},
                            dropdownItems: getCurrencies,
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                'Add Proof',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              SizedBox(
                                width: width / 20,
                              ),
                              Icon(
                                CupertinoIcons.add_circled_solid,
                                size: 20,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
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

  Widget myContainer(Widget body) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 15.0),
              child: Column(
                children: [
                  body,
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

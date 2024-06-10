import 'dart:convert';
import 'dart:developer';
import 'dart:typed_data';
import 'package:country_picker/country_picker.dart';
import 'package:easy_localization/easy_localization.dart';
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
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
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
  final _formKey = GlobalKey<FormState>();
  String? assetLogo;
  String fundingMethod = 'e-Naira';
  bool hasAdditionalKYCRequirements = false;
  String proceedPayoutCurrency = 'e-Naira';
  late int numberOfTokenToBeSold;
  late int numberOfTokenToBeIssued;
  late int totalTokenHeldByManager;
  late int pricePerToken;
  late String assetCode;
  late DateTime? startDate;
  late DateTime? endDate;
  late int capQuantity;
  late int capDurationInDays;
  late String proceedCycle;
  late List<String> exemptedCountries;
  late String additionalKYCRequirements;
  late bool investorAccreditationRequired;
  late bool capOnPurchase;
  late dynamic data = {};

  List<String> fundingOptions = [
    'e-Naira',
    'TROV',
    'XBN',
  ];
  List<String> payoutCycleOptions = [
    'Monthly',
    'Quarterly',
    'Yearly',
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

  List<DropdownMenuItem<String>> get getFundingOptions {
    List<DropdownMenuItem<String>> options = [];
    fundingOptions.forEach((item) {
      options.add(DropdownMenuItem(
          child: Text(
            item,
            overflow: TextOverflow.ellipsis,
          ),
          value: item));
    });
    return options;
  }

  List<DropdownMenuItem<String>> get getStandardWallets {
    List<DropdownMenuItem<String>> wallets = [];
    appState.userInfo!.getStandardWallets.forEach((wallet) {
      wallets.add(DropdownMenuItem(
          child: Text(
            wallet.alias!,
            overflow: TextOverflow.ellipsis,
          ),
          value: wallet.publicKey));
    });
    return wallets;
  }

  List<DropdownMenuItem<String>> get getPayoutCycles {
    List<DropdownMenuItem<String>> cycles = [];
    payoutCycleOptions.forEach((item) {
      cycles.add(DropdownMenuItem(
          child: Text(
            item,
            overflow: TextOverflow.ellipsis,
          ),
          value: item));
    });
    return cycles;
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    inspect(appState.viewData);
    data = appState.viewData;

    numberOfTokenToBeSold = data['numberOfTokenToBeSold'];
    numberOfTokenToBeIssued = data['numberOfTokenToBeIssued'];
    totalTokenHeldByManager = data['totalTokenHeldByManager'];
    pricePerToken = data['pricePerToken'];
    assetCode = data['assetCode'];
    startDate = DateTime.parse(data['salesStart'].toString());
    endDate = DateTime.parse(data['salesEnd'].toString());
    capOnPurchase = data['capOnPurchase'] == 1;
    capQuantity = data['capQuantity'];
    capDurationInDays = data['capDurationInDays'];
    proceedCycle = data['proceedCycle'];
    assetLogo = data['assetLogo'];
    exemptedCountries = data['exemptedCountries'].toString().isEmpty
        ? ['Pakistan']
        : data['exemptedCountries'].toString().split(',');
    // fundingMethod = data[''];
    hasAdditionalKYCRequirements = data['hasAdditionalKYCRequirements'] == 1;
    proceedPayoutCurrency = data['proceedPayoutCurrency'];
    additionalKYCRequirements = data['additionalKYCRequirements'];
    investorAccreditationRequired = data['investorAccreditationRequired'] == 1;
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                "assettokeninfo".tr(),
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
                      "enterassetcode".tr(),
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
                      "enterassetcode".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          assetCode = value!;
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
                    child: Text(
                      "uploadassetlogo".tr(),
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
                                            "browsefiles".tr(),
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
              if (assetLogo != null && assetLogo!.isNotEmpty) ...[
                GestureDetector(
                  onTap: () {
                    getImage();
                  },
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Image.memory(
                      base64Decode(assetLogo!),
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
                      "nooftokentoissue".tr(),
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
                      "notobeissued".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          numberOfTokenToBeIssued = value!;
                        });
                      },
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
                      "notobesold".tr(),
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
                      "notobesold".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          numberOfTokenToBeSold = value!;
                        });
                      },
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
                      "totaltobeheldbymanager".tr(),
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
                      "totalheldbymanager".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          totalTokenHeldByManager = value!;
                        });
                      },
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
                      "selectwallettoholdassetnotforsale".tr(),
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
                  getStandardWallets,
                  null,
                  appState.userInfo!.getStandardWallets.length > 0
                      ? appState.userInfo!.getStandardWallets.first.alias
                      : '',
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
                      "orcreateanewwallet".tr(),
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
                      "assetsalesandpricing".tr(),
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
                      "pricepertoken".tr(),
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
                      "pricepertoken".tr(),
                      notifier.getbluecolor,
                      null,
                      notifier.getgrey,
                      null,
                      notifier.getblck,
                      notifier.getgrey,
                      70.sp,
                      300.sp,
                      validator: (value) {
                        if (value.isEmpty) {
                          return "fieldcannotbeempty".tr();
                        }
                        return null;
                      },
                      onSaved: (value) {
                        setState(() {
                          pricePerToken = value!;
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
                    child: Text(
                      "calculatedassetvalue".tr(),
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
                      "selectfundingmethod".tr(),
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
                  (value) {
                    setState(() {
                      fundingMethod = value.toString();
                    });
                  },
                  getFundingOptions,
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
                            "salesstartsfrom".tr(),
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
                      ButtonOutlined(
                        startDate != null
                            ? DateFormat('MMMM dd, yyyy').format(startDate!)
                            : "starts".tr(),
                        notifier.getwihitecolor,
                        notifier.getgrey,
                        borderColor: notifier.getgrey,
                        width: width / 2.7,
                        height: 50.sp,
                        onTap: () {
                          showDatePicker(
                            context: context,
                            initialDate: DateTime.now(),
                            firstDate:
                                DateTime.fromMicrosecondsSinceEpoch(1000),
                            lastDate: DateTime.now().add(Duration(days: 730)),
                          ).then((value) => setState(() {
                                startDate = value;
                              }));
                        },
                      ),
                    ],
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            "salesendon".tr(),
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
                      ButtonOutlined(
                        endDate != null
                            ? DateFormat('MMMM dd, yyyy').format(endDate!)
                            : "ends".tr(),
                        notifier.getwihitecolor,
                        notifier.getgrey,
                        borderColor: notifier.getgrey,
                        width: width / 2.7,
                        height: 50.sp,
                        onTap: () {
                          showDatePicker(
                            context: context,
                            initialDate: DateTime.now(),
                            firstDate:
                                DateTime.fromMicrosecondsSinceEpoch(1000),
                            lastDate: DateTime.now().add(Duration(days: 730)),
                          ).then((value) => {
                                setState(() {
                                  endDate = value;
                                })
                              });
                        },
                      ),
                    ],
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         'Asset Wallets',
              //         style: TextStyle(
              //           fontSize: 15,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         'Minting Wallet',
              //         style: TextStyle(
              //           fontSize: 12,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: CustomTextFormField.textField(
              //         'Minting wallet',
              //         notifier.getbluecolor,
              //         null,
              //         notifier.getgrey,
              //         null,
              //         notifier.getblck,
              //         notifier.getgrey,
              //         70.sp,
              //         300.sp,
              //         // controller: referrerController,
              //         // validator: validateReferrer,
              //         onSaved: (value) {},
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: Text(
              //         'Market Making Wallet',
              //         style: TextStyle(
              //           fontSize: 12,
              //           fontFamily: fontsemibold,
              //           color: notifier.getbluewhitecolor,
              //         ),
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              // Row(
              //   children: [
              //     Padding(
              //       padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //       child: CustomTextFormField.textField(
              //         'Market Making Wallet',
              //         notifier.getbluecolor,
              //         null,
              //         notifier.getgrey,
              //         null,
              //         notifier.getblck,
              //         notifier.getgrey,
              //         70.sp,
              //         300.sp,
              //         // controller: referrerController,
              //         // validator: validateReferrer,
              //         onSaved: (value) {},
              //       ),
              //     ),
              //   ],
              // ),
              // SizedBox(
              //   height: height / 50,
              // ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "caponpurchase".tr(),
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
                      value: capOnPurchase,
                      onChanged: (val) async {
                        setState(() {
                          capOnPurchase = !capOnPurchase;
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
                        "capexplained".tr(),
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
              if (capOnPurchase) ...[
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20.0),
                      child: Text(
                        "capquantity".tr(),
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
                        "quantity".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        70.sp,
                        300.sp,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            capQuantity = value!;
                          });
                        },
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
                        "capduration".tr(),
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
                        "noofdays".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        70.sp,
                        300.sp,
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            capQuantity = value!;
                          });
                        },
                        keyboardtype: TextInputType.multiline,
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 30,
                ),
              ],
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      "proceedpayout".tr(),
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
                      "proceedpayoutcycle".tr(),
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
                  (value) {
                    setState(() {
                      proceedCycle = value.toString();
                    });
                  },
                  getPayoutCycles,
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
                      "proceedpayoutcurrency".tr(),
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
                  (value) {
                    setState(() {
                      proceedPayoutCurrency = value.toString();
                    });
                  },
                  getFundingOptions,
                  null,
                  'e-Naira',
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
                      "primarybuyerrequirement".tr(),
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
                        "selectexemptedcountries".tr(),
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
                padding: const EdgeInsets.symmetric(horizontal: 15.0),
                child: Container(
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(10.0),
                    ),
                    color: notifier.isDark
                        ? notifier.getbluecolor90
                        : notifier.getaddsubwalletgrey,
                    child: TextButton(
                      onPressed: () {
                        showCountryPicker(
                          context: context,
                          onSelect: (Country country) {
                            setState(() {
                              exemptedCountries.add(country.name);
                            });
                          },
                        );
                      },
                      style: ButtonStyle(
                          elevation: MaterialStateProperty.all<double>(0)),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            exemptedCountries.length > 0
                                ? exemptedCountries.first
                                : 'Pakistan',
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          Icon(Icons.keyboard_arrow_down_rounded),
                        ],
                      ),
                    ),
                  ),
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
                                    if (exemptedCountries.isNotEmpty) ...[
                                      for (var country
                                          in exemptedCountries) ...[
                                        userItem(
                                          country,
                                          () {
                                            setState(() {
                                              exemptedCountries.remove(country);
                                            });
                                          },
                                          foreColor: wihitecolor,
                                          backColor: notifier.getbluebackcolor,
                                        ),
                                      ]
                                    ] else ...[
                                      Text(
                                        "nameofexemptedcountriesappearhere"
                                            .tr(),
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
                      "arethereadditionkycreq".tr(),
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
                      'yes'.tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    leading: Radio(
                      value: hasAdditionalKYCRequirements,
                      groupValue: true,
                      activeColor: notifier.getbluewhitecolor,
                      fillColor:
                          MaterialStateProperty.all(notifier.getbluewhitecolor),
                      onChanged: (value) {
                        setState(() {
                          hasAdditionalKYCRequirements = true;
                        });
                      },
                    ),
                  ),
                  ListTile(
                    title: Text(
                      'no'.tr(),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    leading: Radio(
                      value: hasAdditionalKYCRequirements,
                      groupValue: false,
                      fillColor:
                          MaterialStateProperty.all(notifier.getbluewhitecolor),
                      activeColor: notifier.getbluewhitecolor,
                      onChanged: (value) {
                        setState(() {
                          hasAdditionalKYCRequirements = false;
                          print('addAdditionalKyc: $value');
                        });
                      },
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  if (hasAdditionalKYCRequirements) ...[
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 20.0),
                          child: Text(
                            "listrequirements".tr(),
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
                            additionalKYCRequirements,
                            notifier.getbluecolor,
                            notifier.getgrey,
                            notifier.getblck,
                            notifier.getgrey,
                            100.sp,
                            300.sp,
                            validator: (value) {
                              if (value.isEmpty) {
                                return "fieldcannotbeempty".tr();
                              }
                              return null;
                            },
                            onSaved: (value) {
                              setState(() {
                                additionalKYCRequirements = value!;
                              });
                            },
                            keyboardtype: TextInputType.multiline,
                            minLines: 3,
                            maxLines: null,
                          ),
                        ),
                      ],
                    ),
                  ],
                ],
              ),
              accreditInvestors(),
              SizedBox(
                height: height / 30,
              ),
              Button(
                "saveandcontinuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  var form = _formKey.currentState;
                  if (form!.validate()) {
                    submitForm();
                  }
                },
              ),
              SizedBox(
                height: height / 10,
              ),
            ],
          ),
        ),
      ),
    );
  }

  void submitForm() async {
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credential
      var mintingWallet = appState.userInfo!.getMintingWallets[0];

      data['numberOfTokenToBeSold'] = numberOfTokenToBeSold;
      data['numberOfTokenToBeIssued'] = numberOfTokenToBeIssued;
      data['totalTokenHeldByManager'] = totalTokenHeldByManager;
      data['pricePerToken'] = pricePerToken;
      data['assetCode'] = assetCode;
      data['startDate'] = startDate!.toUtc().toString();
      data['endDate'] = endDate!.toUtc().toString();
      data['capOnPurchase'] = capOnPurchase ? 1 : 0;
      data['capQuantity'] = capQuantity;
      data['capDurationInDays'] = capDurationInDays;
      data['proceedCycle'] = proceedCycle;
      data['assetLogo'] = assetLogo;
      data['exemptedCountries'] = exemptedCountries.join(',');
      data['hasAdditionalKYCRequirements'] =
          hasAdditionalKYCRequirements ? 1 : 0;
      data['proceedPayoutCurrency'] = proceedPayoutCurrency;
      data['additionalKYCRequirements'] = additionalKYCRequirements;
      data['investorAccreditationRequired'] =
          investorAccreditationRequired ? 1 : 0;

      inspect(data);
      String requestBody = jsonEncode(data);
      print('requestBody =======> $requestBody');
      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: mintingWallet.publicKey!,
      );

      hideLoader(context);

      print('responseData ${responseData['data']}');
      inspect(responseData['data']);

      if (responseData['statusCode'] == 200) {
        Navigator.of(context).pop();
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
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
            value: investorAccreditationRequired,
            onChanged: (bool? value) {
              setState(() {
                investorAccreditationRequired = !investorAccreditationRequired;
              });
            },
          ),
        ),
        Container(
          width: width / 1.2,
          child: Text(
            "investormustbeaccredited".tr(),
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
        assetLogo = imageBase64Uncompressed;
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

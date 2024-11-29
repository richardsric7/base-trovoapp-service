import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:image_picker/image_picker.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TokenizationFeePayment extends StatefulWidget {
  const TokenizationFeePayment({Key? key}) : super(key: key);

  @override
  State<TokenizationFeePayment> createState() => _TokenizationFeePayment();
}

class _TokenizationFeePayment extends State<TokenizationFeePayment>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  final formKey = GlobalKey<FormState>();
  late TokenizedAsset tokenizedAsset;
  late double fiatFeeCap;
  late double tokenFee;
  late double fiatFee;
  PlatformFile? recieptFile;
  String errorMsg = '';
  String transactionId = '';
  bool hasMadePayment = false;
  String preferredPaymentMethod = 'CNGN';
  List<String> paymentMethods = ['Fiat', 'CNGN'];

  List<DropdownMenuItem<String>> get getPaymentMethods {
    List<DropdownMenuItem<String>> cycles = [];
    paymentMethods.forEach((item) {
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
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    tokenizedAsset = TokenizedAsset().deserializeJson(appState.viewData!);
    var fiatPercentage = appState.tokenizationData["tokenizationFees"]
        [tokenizedAsset.tokenizationFeeId]['feeFiatPercentage'];
    var assetPercentage = appState.tokenizationData["tokenizationFees"]
        [tokenizedAsset.tokenizationFeeId]['feeAssetPercentage'];
    fiatFeeCap = double.parse(appState.tokenizationData["tokenizationFees"]
            [tokenizedAsset.tokenizationFeeId]['feeFiatCap']
        .toString());
    tokenFee = tokenizedAsset.numberOfTokenToBeIssued! * assetPercentage;
    fiatFee = tokenizedAsset.assetCurrentValue! * fiatPercentage;
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
        appBar: CustomAppBar(context, notifier.getwihitecolor, "payment".tr(),
                notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "selectpaymentmethod".tr(),
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
              ),
              SizedBox(height: height / 90),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    setState(() {
                      preferredPaymentMethod = value.toString();
                    });
                  },
                  getPaymentMethods,
                  preferredPaymentMethod.isEmpty
                      ? null
                      : preferredPaymentMethod,
                  'Select payment method',
                  context,
                  null,
                ),
              ),
              SizedBox(
                height: height / 30,
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "payto".tr(args: [
                    "${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} ${preferredPaymentMethod == 'CNGN' ? preferredPaymentMethod : tokenizedAsset.proceedPayoutCurrency}"
                  ]),
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w400,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              if (preferredPaymentMethod == 'CNGN') ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(10.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(
                              vertical: 10.0, horizontal: 15),
                          child: Text(
                            "walletaddress".tr(),
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                        Row(
                          children: [
                            Padding(
                              padding: const EdgeInsets.all(10.0),
                              child: Container(
                                width: width / 1.4,
                                child: Text(
                                  tokenizedAsset.walletToHoldAssetsNotForSale!
                                      .toLowerCase(),
                                  style: TextStyle(
                                      fontSize: 15,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody),
                                ),
                              ),
                            ),
                            IconButton(
                              onPressed: () {
                                Clipboard.setData(
                                  ClipboardData(
                                    text: tokenizedAsset
                                        .walletToHoldAssetsNotForSale!,
                                  ),
                                );
                                showSnackBar("publickey".tr(), context);
                              },
                              icon: Icon(Icons.copy,
                                  size: 20, color: notifier.getbluewhitecolor),
                            ),
                          ],
                        ),
                        SizedBox(height: height / 90),
                      ],
                    ),
                  ),
                ),
              ] else ...[
                Padding(
                  padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius:
                          const BorderRadius.all(Radius.circular(10.0)),
                      color: notifier.isDark
                          ? darktilewhitecolor
                          : notifier.getaddsubwalletgrey,
                    ),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(
                          vertical: 10.0, horizontal: 15),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            "bankaccountdetails".tr(),
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                            ),
                          ),
                          SizedBox(height: height / 50),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                "accountname".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Text(
                                "Trovo Tokenizer",
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                            ],
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                "accountnumber".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    "0123908438",
                                    style: TextStyle(
                                        fontSize: 15,
                                        fontWeight: FontWeight.w700,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold),
                                  ),
                                  IconButton(
                                    onPressed: () {
                                      Clipboard.setData(
                                        ClipboardData(
                                          text: tokenizedAsset
                                              .walletToHoldAssetsNotForSale!,
                                        ),
                                      );
                                      showSnackBar(
                                          "accountnumber".tr(), context);
                                    },
                                    icon: Icon(Icons.copy,
                                        size: 20,
                                        color: notifier.getbluewhitecolor),
                                  ),
                                ],
                              ),
                            ],
                          ),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                "bank".tr(),
                                style: TextStyle(
                                    fontSize: 15,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody),
                              ),
                              Text(
                                "First Bank",
                                style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w700,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                            ],
                          ),
                          SizedBox(height: height / 90),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
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
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Column(
                      children: [
                        Text(
                          "${"important".tr()}!",
                          style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w700,
                              color: Colors.red,
                              fontFamily: fontbody),
                        ),
                        SizedBox(
                          height: height / 70,
                        ),
                        Text(
                          "yourtokenizationfeeis".tr(),
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody),
                        ),
                        SizedBox(
                          height: height / 90,
                        ),
                        Text(
                          "\$${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} + ${formatNumber(tokenFee)} ${tokenizedAsset.assetCode}",
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(
                          height: height / 70,
                        ),
                        Text(
                          "pleasepayasap".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody),
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              SizedBox(height: height / 50),
              Container(
                width: width / 1.08,
                child: Row(
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
                        value: hasMadePayment,
                        onChanged: (bool? value) {
                          setState(() {
                            hasMadePayment = value!;
                          });
                        },
                      ),
                    ),
                    Container(
                      width: width / 1.3,
                      child: Text(
                        "ihavemadepayment".tr(),
                        overflow: TextOverflow.visible,
                        style: TextStyle(
                            fontSize: 15,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ],
                ),
              ),
              if (hasMadePayment && transactionId.isEmpty) ...[
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20),
                      child: Text(
                        "uploadreceipt".tr(),
                        textAlign: TextAlign.start,
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w400,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold,
                        ),
                      ),
                    ),
                  ],
                ),
                GestureDetector(
                  onTap: () {
                    getImage();
                  },
                  child: Padding(
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
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Padding(
                            padding: const EdgeInsets.all(15.0),
                            child: Column(
                              children: [
                                SizedBox(height: height / 70),
                                Icon(
                                  Icons.file_present_rounded,
                                  color: notifier.getbluewhitecolor,
                                  size: 35,
                                ),
                                SizedBox(height: height / 70),
                                if (recieptFile != null) ...[
                                  Text(
                                    "fileuploaded".tr(args: [
                                      truncate(recieptFile!.name.toString(),
                                          length: 25)
                                    ]),
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                        fontSize: 15,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody),
                                  ),
                                  SizedBox(height: 10),
                                  Row(
                                    children: [
                                      OutlinedButton(
                                        onPressed: () async {
                                          getImage();
                                        },
                                        style: ButtonStyle(
                                          side: MaterialStateProperty.all(
                                            BorderSide(
                                                color:
                                                    notifier.getbluewhitecolor,
                                                width:
                                                    2), // Example of BorderSide
                                          ),
                                          shape: MaterialStateProperty.all<
                                              RoundedRectangleBorder>(
                                            const RoundedRectangleBorder(
                                              borderRadius: BorderRadius.all(
                                                Radius.circular(10),
                                              ),
                                            ),
                                          ),
                                        ),
                                        child: Row(
                                          children: [
                                            Icon(
                                              Icons.file_upload_outlined,
                                            ),
                                            SizedBox(width: 5),
                                            Text(
                                              "replacefile".tr(),
                                              style: TextStyle(
                                                fontSize: 12,
                                                fontFamily: fontsemibold,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      SizedBox(width: 10),
                                      OutlinedButton(
                                        onPressed: () async {
                                          setState(() {
                                            recieptFile = null;
                                          });
                                        },
                                        style: ButtonStyle(
                                          side: MaterialStateProperty.all(
                                            BorderSide(
                                                color: Colors.red,
                                                width:
                                                    2), // Example of BorderSide
                                          ),
                                          shape: MaterialStateProperty.all<
                                              RoundedRectangleBorder>(
                                            const RoundedRectangleBorder(
                                              borderRadius: BorderRadius.all(
                                                Radius.circular(10),
                                              ),
                                            ),
                                          ),
                                        ),
                                        child: Row(
                                          children: [
                                            Icon(
                                              Icons.cancel_outlined,
                                              color: Colors.red,
                                            ),
                                            SizedBox(width: 5),
                                            Text(
                                              "removefile".tr(),
                                              style: TextStyle(
                                                fontFamily: fontsemibold,
                                                fontSize: 12,
                                                color: Colors.red,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                    ],
                                  ),
                                ] else ...[
                                  Text(
                                    "browseimageorpdf".tr(),
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                        fontSize: 15,
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody),
                                  ),
                                ],
                                SizedBox(height: height / 70),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              ],
              if (hasMadePayment &&
                  (recieptFile == null && transactionId.isEmpty)) ...[
                Row(
                  children: [
                    Expanded(
                      child: Container(
                        margin: const EdgeInsets.only(left: 27.0, right: 27.0),
                        child: Divider(
                          color: notifier.getgrey,
                          height: 50,
                        ),
                      ),
                    ),
                    Text(
                      "oR".tr(),
                      style: TextStyle(color: notifier.getgrey),
                    ),
                    Expanded(
                      child: Container(
                        margin: const EdgeInsets.only(left: 27.0, right: 27.0),
                        child: Divider(
                          color: notifier.getgrey,
                          height: 50,
                        ),
                      ),
                    ),
                  ],
                ),
              ],
              if (hasMadePayment && recieptFile == null) ...[
                SizedBox(height: height / 50),
                Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 20),
                      child: Text(
                        "entertransactionid".tr(),
                        textAlign: TextAlign.start,
                        style: TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w400,
                          color: notifier.getbluewhitecolor,
                          fontFamily: fontsemibold,
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
                        "transactionid".tr(),
                        notifier.getbluecolor,
                        null,
                        notifier.getgrey,
                        null,
                        notifier.getblck,
                        notifier.getgrey,
                        70.sp,
                        300.sp,
                        onChanged: (value) {
                          setState(() {
                            transactionId = value;
                          });
                        },
                        validator: (value) {
                          if (value.isEmpty) {
                            return "fieldcannotbeempty".tr();
                          }
                          return null;
                        },
                        onSaved: (value) {
                          setState(() {
                            transactionId = value!;
                          });
                        },
                      ),
                    ),
                  ],
                ),
              ],
              SizedBox(
                height: height / 20,
              ),
              Button(
                "confirmpayment".tr(),
                hasMadePayment &&
                        (transactionId.isNotEmpty || recieptFile != null)
                    ? notifier.getbluecolor
                    : notifier.getbluecolor80,
                wihitecolor,
                onTap: () async {
                  if (hasMadePayment &&
                      (transactionId.isNotEmpty || recieptFile != null)) {
                    showLoader(context);
                    await Future.delayed(Duration(seconds: 1));
                    hideLoader(context);
                    appState.viewData![SuccessViewPageConfig.key] = {
                      'title': '',
                      'buttonText': 'Go to Home',
                      // 'useOnDone': true,
                      // 'onDone': () {
                      //   appState.currentAction = PageAction(
                      //     state: PageState.addPage,
                      //     page: TokenizationFeePaymentViewPageConfig,
                      //   );
                      // },
                      'message':
                          'Your Proof of Payment has been submitted successfully and is awaiting confirmation. Your asset tokenization application will be processed once payment has been confirmed.',
                    };
                    appState.currentAction = PageAction(
                        state: PageState.addPage, page: SuccessViewPageConfig);
                  }
                },
              ),
              SizedBox(
                height: height / 20,
              ),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> getImage() async {
    recieptFile = await getFile();
    if (recieptFile != null && recieptFile!.size > 900000) {
      errorMsg = "filesizeerror".tr();
      recieptFile = null;
    }
    setState(() {});
  }

  Future<dynamic> getBase64Image(XFile image) async {
    //
    List<int> imageBytes = await image.readAsBytes();
    String imageB64 = base64Encode(imageBytes);
    return imageB64;
  }
}

import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/models/tokenizedAsset.dart';
import 'package:trovo_wallet/network/requests.dart';
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
  String transactionReference = '';
  bool hasMadePayment = false;
  String preferredPaymentMethod = 'FIAT';
  List<String> paymentMethods = [];

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
    for (var i = 0;
        i < appState.tokenizationData['feePaymentMethods'].length;
        i++) {
      paymentMethods
          .add(appState.tokenizationData['feePaymentMethods'][i]['id']);
    }
    initializeData();
  }

  void initializeData() {
    tokenizedAsset = TokenizedAsset().deserializeJson(appState.viewData!);
    print('deserialized tokenized asset');
    inspect(tokenizedAsset);
    var fiatPercentage = appState.tokenizationData["tokenizationFees"]
        [tokenizedAsset.tokenizationFeeId]['feeFiatPercentage'];
    var assetPercentage = appState.tokenizationData["tokenizationFees"]
        [tokenizedAsset.tokenizationFeeId]['feeAssetPercentage'];
    fiatFeeCap = double.parse(appState.tokenizationData["tokenizationFees"]
            [tokenizedAsset.tokenizationFeeId]['feeFiatCap']
        .toString());
    tokenFee = tokenizedAsset.numberOfTokenToBeIssued! * assetPercentage;
    fiatFee = tokenizedAsset.assetCurrentValue! * fiatPercentage;
    hasMadePayment =
        tokenizedAsset.proofOfPaymentDocuments?.isNotEmpty ?? false;
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
                    "${formatNumber(fiatFee > fiatFeeCap ? fiatFeeCap : fiatFee)} ${preferredPaymentMethod == 'STABLE COIN' ? preferredPaymentMethod : tokenizedAsset.proceedPayoutCurrency}"
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
              if (preferredPaymentMethod == 'STABLE COIN') ...[
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
              for (var item in tokenizedAsset.proofOfPaymentDocuments!) ...[
                uploadedDoc(item.documentUrl!, item.id.toString(),
                    item.transactionReference),
              ],
              if (hasMadePayment) ...[
                TextButton(
                  onPressed: () {
                    uploadTokenizationFeePopup(context,
                        onSubmit: (file, transactionReference) async {
                      transactionReference = transactionReference;
                      await uploadFile(file, transactionReference);
                    });
                  },
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.file_upload_outlined,
                        size: 22,
                        color: notifier.getgreencolor,
                      ),
                      SizedBox(width: 3),
                      Text(
                        'uploadproof'.tr(),
                        style: TextStyle(
                          decoration: TextDecoration.underline,
                          fontSize: 16,
                          fontFamily: fontbody,
                          color: notifier.getgreencolor,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
              SizedBox(
                height: height / 20,
              ),
              Button(
                "confirmpayment".tr(),
                tokenizedAsset.proofOfPaymentDocuments?.isNotEmpty ?? false
                    ? notifier.getbluecolor
                    : notifier.getbluecolor80,
                wihitecolor,
                onTap: () async {
                  if (tokenizedAsset.proofOfPaymentDocuments?.isNotEmpty ??
                      false) {
                    confirmPayments();
                  }
                },
              ),
              SizedBox(height: 10),
              ButtonOutlined(
                "Go to homepage".tr(),
                wihitecolor,
                tokenizedAsset.proofOfPaymentDocuments?.isNotEmpty ?? false
                    ? notifier.getbluecolor
                    : notifier.getbluecolor80,
                onTap: () async {
                  appState.currentAction = PageAction(
                      state: PageState.replaceAll, page: BottomHomePageConfig);
                },
              ),
              SizedBox(height: height / 20),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  void confirmPayments() async {
    try {
      showLoader(context);
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;
      // String requestBody = jsonEncode(appState.viewData);
      // print('requestBody  =======> $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/tokenization/fee/${tokenizedAsset.id}',
        body: "",
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: mintingWalletPublicKey,
      );

      hideLoader(context);

      print('responseData token information  ${responseData['data']}');
      inspect(responseData);

      if (responseData['statusCode'] == 200) {
        Navigator.of(context).pop();
        appState.viewData![SuccessViewPageConfig.key] = {
          'title': '',
          'buttonText': 'Go to Home',
          'message':
              'Your Proof of Payment has been submitted successfully and is awaiting confirmation. Your asset tokenization application will be processed once payment has been confirmed.',
        };
        appState.currentAction =
            PageAction(state: PageState.addPage, page: SuccessViewPageConfig);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Widget uploadedDoc(String url, String docId, String? ref) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        TextButton(
          onPressed: () {
            if (url.isNotEmpty && url.endsWith('.pdf')) {
              appState.pdfUrl = url;
              appState.currentAction =
                  PageAction(state: PageState.addPage, page: PdfViewPageConfig);

              return;
            }

            appState.goToWebView(url);
          },
          child: Column(
            children: [
              Text(
                truncateString(url.toString()),
                textAlign: TextAlign.center,
                style: TextStyle(
                    decoration: TextDecoration.underline,
                    fontSize: 15,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody),
              ),
              if (ref != null && ref.isNotEmpty) ...[
                SizedBox(
                  width: 240,
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Ref:',
                        textAlign: TextAlign.start,
                        style: TextStyle(
                            fontSize: 15,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold),
                      ),
                      SizedBox(width: 2),
                      SizedBox(
                        width: 200,
                        child: Text(
                          ref,
                          textAlign: TextAlign.start,
                          style: TextStyle(
                              fontSize: 15,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ],
          ),
        ),
        TextButton(
          onPressed: () async {
            await deleteFile(docId);
            setState(() {});
          },
          child: Icon(
            CupertinoIcons.trash,
            color: Colors.red,
            size: 20,
          ),
        ),
      ],
    );
    // return Padding(
    //   padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
    //   child: Container(
    //     decoration: BoxDecoration(
    //       border: Border.all(color: notifier.getbluewhitecolor, width: 1),
    //       borderRadius: const BorderRadius.all(Radius.circular(15.0)),
    //       color: notifier.isDark
    //           ? darktilewhitecolor
    //           : notifier.getaddsubwalletgrey,
    //     ),
    //     child: Row(
    //       mainAxisAlignment: MainAxisAlignment.center,
    //       children: [
    //         Padding(
    //           padding: const EdgeInsets.all(5.0),
    //           child: Column(
    //             children: [
    //               SizedBox(height: height / 70),
    //               Icon(
    //                 Icons.file_present_rounded,
    //                 color: notifier.getbluewhitecolor,
    //                 size: 35,
    //               ),
    //               SizedBox(height: height / 70),

    //               SizedBox(height: 10),
    //               Row(
    //                 children: [
    //                   OutlinedButton(
    //                     onPressed: () async {
    //                       if (url.isNotEmpty && url.endsWith('.pdf')) {
    //                         appState.pdfUrl = url;
    //                         appState.currentAction = PageAction(
    //                             state: PageState.addPage,
    //                             page: PdfViewPageConfig);

    //                         return;
    //                       }

    //                       appState.goToWebView(url);
    //                     },
    //                     style: ButtonStyle(
    //                       side: MaterialStateProperty.all(
    //                         BorderSide(
    //                             color: notifier.getbluewhitecolor,
    //                             width: 2), // Example of BorderSide
    //                       ),
    //                       shape:
    //                           MaterialStateProperty.all<RoundedRectangleBorder>(
    //                         const RoundedRectangleBorder(
    //                           borderRadius: BorderRadius.all(
    //                             Radius.circular(10),
    //                           ),
    //                         ),
    //                       ),
    //                     ),
    //                     child: Row(
    //                       children: [
    //                         Icon(
    //                           Icons.file_present,
    //                         ),
    //                         SizedBox(width: 5),
    //                         Text(
    //                           "viewfile".tr(),
    //                           style: TextStyle(
    //                             fontSize: 12,
    //                             fontFamily: fontsemibold,
    //                           ),
    //                         ),
    //                       ],
    //                     ),
    //                   ),
    //                   SizedBox(width: 10),
    //                   OutlinedButton(
    //                     onPressed: () async {
    //                       await deleteFile(docId);
    //                       setState(() {});
    //                     },
    //                     style: ButtonStyle(
    //                       side: MaterialStateProperty.all(
    //                         BorderSide(
    //                             color: Colors.red,
    //                             width: 2), // Example of BorderSide
    //                       ),
    //                       shape:
    //                           MaterialStateProperty.all<RoundedRectangleBorder>(
    //                         const RoundedRectangleBorder(
    //                           borderRadius: BorderRadius.all(
    //                             Radius.circular(10),
    //                           ),
    //                         ),
    //                       ),
    //                     ),
    //                     child: Row(
    //                       children: [
    //                         Icon(
    //                           Icons.cancel_outlined,
    //                           color: Colors.red,
    //                         ),
    //                         SizedBox(width: 5),
    //                         Text(
    //                           "removefile".tr(),
    //                           style: TextStyle(
    //                             fontFamily: fontsemibold,
    //                             fontSize: 12,
    //                             color: Colors.red,
    //                           ),
    //                         ),
    //                       ],
    //                     ),
    //                   ),
    //                 ],
    //               ),
    //               SizedBox(height: height / 70),
    //             ],
    //           ),
    //         ),
    //       ],
    //     ),
    //   ),
    // );
  }

  Future<void> uploadFile(
    PlatformFile file,
    String? transactionReference,
  ) async {
    try {
      showLoader(context);
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;

      Map responseData = await makePutRequestForFeeRecieptUpload(
        uri: '/v1/tokenization/fee/${tokenizedAsset.id}',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: mintingWalletPublicKey,
        file: file,
        transactionReference: transactionReference ?? '',
        tokenizationFeePaymentMethodID: preferredPaymentMethod,
      );

      // print('responseData $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(
        context,
        title: "error".tr(),
        message: "Sorry, something went wrong. Please try again.",
      );
    }
  }

  Future<void> deleteFile(String documentId) async {
    try {
      showLoader(context);
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;
      Map requestBody = {};
      Map responseData = await makeDeleteRequest(
        uri: '/v1/tokenization/fee/$documentId',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: mintingWalletPublicKey,
        body: jsonEncode(requestBody),
      );

      print("response ============> ${responseData}");
      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
      }
      hideLoader(context);
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(
        context,
        title: "error".tr(),
        message: "Sorry, something went wrong. Please try again.",
      );
    }
  }

  Future<void> refreshCurrentTokenizationInfo() async {
    print('refreshing tokenization info');
    try {
      var uri = '/v1/tokenization/detail/${appState.viewData!['id']}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      inspect(responseData);
      // print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        print('success');
        appState.viewData = responseData['data'];
        initializeData();
        setState(() {});
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  // Future<void> getImage() async {
  //   recieptFile = await getFile();
  //   if (recieptFile != null && recieptFile!.size > 900000) {
  //     errorMsg = "filesizeerror".tr();
  //     recieptFile = null;
  //   } else {
  //     await uploadFile(recieptFile!);
  //   }
  //   setState(() {});
  // }
}

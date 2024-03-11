import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:easy_localization/easy_localization.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:get/get_connect.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
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

  String selectedProofOfExistenceOption = '';
  String proofOfExistenceFileName = '';

  String selectedProofOfOwnershipOption = '';
  String proofOfOwnershipFileName = '';

  String selectedAssetStatusVerificationOption = '';
  String assetStatusVerificationFileName = '';

  String selectedAssetCustodianAgreementOption = '';
  String assetCustodianAgreementFileName = '';

  String selectedProofOfAssetManagerOption = '';
  String proofOfAssetManagerFileName = '';

  String selectedAssetProtectionDocumentOption = '';
  String assetProtectionDocumentFileName = '';

  String selectedAssetValuationCertificateOption = '';
  String assetValuationCertificateFileName = '';

  String selectedProofOfAdditionalCostOutsideValuationOption = '';
  String additionalCostOutsideValuationFileName = '';

  String selectedProofOfAssetConditionOption = '';
  String proofOfAssetConditionFileName = '';

  String selectedThirdPartyTokenizationAgreementOption = '';
  String thirdPartyTokenizationAgreementFileName = '';

  String selectedThirdPartyAssetOwnerBusinessRegOption = '';
  String thirdPartyAssetOwnerBusinessRegFileName = '';

  String selectedThirdPartyAssetOwnerProofOfAddressOption = '';
  String thirdPartyAssetOwnerProofOfAddressFileName = '';

  String selectedSecApprovalOption = '';
  String secRegFileName = '';

  String selectedProofOfComplianceOption = '';
  String proofOfComplianceFileName = '';

  String selectedProofOfEnvComplianceOption = '';
  String proofOfEnvComplianceFileName = '';

  String selectedEnvImpactAssessmentReportOption = '';
  String envImpactAssessmentReportFileName = '';

  String selectedProofofLegalCounselOption = '';
  String proofOfLegalCounselFileName = '';

  String selectedLegalAdvisorsContactOption = '';
  String legalAdvisorsContactFileName = '';

  String selectedProofofMortgagesorLiensOption = '';
  String proofofMortgagesorLiensFileName = '';

  String selectedProofofOutstandingLoansOption = '';
  String proofofOutstandingLoansFileName = '';

  String selectedProofofLegalDisputesOnAssetOption = '';
  String proofofLegalDisputesOnAssetFileName = '';

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
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  uploadFile(file, 1, 'ProofOfAssetExistence');
                  setState(() {
                    selectedProofOfExistenceOption = selectedOption;
                    proofOfExistenceFileName = uploadedFileName;
                  });
                },
                label: 'Proof of Asset Existence',
                selectedOption: selectedProofOfExistenceOption,
                uploadedFileName: proofOfExistenceFileName,
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofOfOwnershipOption = selectedOption;
                    proofOfOwnershipFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofOfOwnershipOption,
                uploadedFileName: proofOfOwnershipFileName,
                label: 'Proof of Ownership',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedAssetStatusVerificationOption = selectedOption;
                    assetStatusVerificationFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedAssetStatusVerificationOption,
                uploadedFileName: assetStatusVerificationFileName,
                label: 'Asset Status Verification',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedAssetCustodianAgreementOption = selectedOption;
                    assetCustodianAgreementFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedAssetCustodianAgreementOption,
                uploadedFileName: assetCustodianAgreementFileName,
                label: 'Asset Custodian Agreement',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofOfAssetManagerOption = selectedOption;
                    proofOfAssetManagerFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofOfAssetManagerOption,
                uploadedFileName: proofOfAssetManagerFileName,
                label: 'Proof of Asset Manager',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedAssetProtectionDocumentOption = selectedOption;
                    assetProtectionDocumentFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedAssetProtectionDocumentOption,
                uploadedFileName: assetProtectionDocumentFileName,
                label: 'Asset Protection Document',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedAssetValuationCertificateOption = selectedOption;
                    assetValuationCertificateFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedAssetValuationCertificateOption,
                uploadedFileName: assetValuationCertificateFileName,
                label: 'Asset Valuation Certificate',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofOfAdditionalCostOutsideValuationOption =
                        selectedOption;
                    additionalCostOutsideValuationFileName = uploadedFileName;
                  });
                },
                selectedOption:
                    selectedProofOfAdditionalCostOutsideValuationOption,
                uploadedFileName: additionalCostOutsideValuationFileName,
                label: 'Proof of Additional Cost Outside Valuation',
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofOfAssetConditionOption = selectedOption;
                    proofOfAssetConditionFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofOfAssetConditionOption,
                uploadedFileName: proofOfAssetConditionFileName,
                label: "Proof of Asset's Condition",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedThirdPartyTokenizationAgreementOption =
                        selectedOption;
                    thirdPartyTokenizationAgreementFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedThirdPartyTokenizationAgreementOption,
                uploadedFileName: thirdPartyTokenizationAgreementFileName,
                label: "Third Party Tokenization Agreement",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedThirdPartyAssetOwnerBusinessRegOption =
                        selectedOption;
                    thirdPartyAssetOwnerBusinessRegFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedThirdPartyAssetOwnerBusinessRegOption,
                uploadedFileName: thirdPartyAssetOwnerBusinessRegFileName,
                label: "Third Party Asset Owner’s Business Registration",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedThirdPartyAssetOwnerProofOfAddressOption =
                        selectedOption;
                    thirdPartyAssetOwnerProofOfAddressFileName =
                        uploadedFileName;
                  });
                },
                selectedOption:
                    selectedThirdPartyAssetOwnerProofOfAddressOption,
                uploadedFileName: thirdPartyAssetOwnerProofOfAddressFileName,
                label: "Third Party Asset Owner’s Proof of Address",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedSecApprovalOption = selectedOption;
                    secRegFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedSecApprovalOption,
                uploadedFileName: secRegFileName,
                label: "SEC Registration/Tokenization Approval",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofOfComplianceOption = selectedOption;
                    proofOfComplianceFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofOfComplianceOption,
                uploadedFileName: proofOfComplianceFileName,
                label: "Proof of Compliance with Local Laws and Regulations",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofOfEnvComplianceOption = selectedOption;
                    proofOfEnvComplianceFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofOfEnvComplianceOption,
                uploadedFileName: proofOfEnvComplianceFileName,
                label: "Proof of Compliance with Environmental Standards",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedEnvImpactAssessmentReportOption = selectedOption;
                    envImpactAssessmentReportFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedEnvImpactAssessmentReportOption,
                uploadedFileName: envImpactAssessmentReportFileName,
                label: "Environmental Impact Assessment Report",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofofLegalCounselOption = selectedOption;
                    proofOfLegalCounselFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofofLegalCounselOption,
                uploadedFileName: proofOfLegalCounselFileName,
                label: "Proof of Legal/Financial Counsel",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofofLegalCounselOption = selectedOption;
                    proofOfLegalCounselFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofofLegalCounselOption,
                uploadedFileName: proofOfLegalCounselFileName,
                label: "Legal/Financial Advisors Contact",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofofMortgagesorLiensOption = selectedOption;
                    proofofMortgagesorLiensFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofofMortgagesorLiensOption,
                uploadedFileName: proofofMortgagesorLiensFileName,
                label: "Proof of Existing Mortgages or Liens on Asset",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofofOutstandingLoansOption = selectedOption;
                    proofofOutstandingLoansFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofofOutstandingLoansOption,
                uploadedFileName: proofofOutstandingLoansFileName,
                label: "Proof of Outstanding Loans on Asset",
                documentOptions: getCurrencies),
            myContainer(
                onDone: (String selectedOption, String uploadedFileName,
                    PlatformFile file) {
                  setState(() {
                    selectedProofofLegalDisputesOnAssetOption = selectedOption;
                    proofofLegalDisputesOnAssetFileName = uploadedFileName;
                  });
                },
                selectedOption: selectedProofofLegalDisputesOnAssetOption,
                uploadedFileName: proofofLegalDisputesOnAssetFileName,
                label: "Proof of  Legal Disputes or Encumbrances on Asset",
                documentOptions: getCurrencies),
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

  Widget myContainer({
    required String label,
    required List<DropdownMenuItem<String>> documentOptions,
    required String selectedOption,
    required String uploadedFileName,
    required void Function(
            String selectedOption, String uploadedFileName, PlatformFile file)
        onDone,
  }) {
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
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.start,
                        children: [
                          Container(
                            width: width / 1.3,
                            child: Text(
                              label,
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
                      Container(
                        width: width / 1.27,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              selectedOption,
                              style: TextStyle(
                                fontSize: 12,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                            Text(
                              uploadedFileName,
                              style: TextStyle(
                                decoration: TextDecoration.underline,
                                fontSize: 12,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                            if (uploadedFileName.isNotEmpty) ...[
                              Icon(
                                CupertinoIcons.delete,
                                size: 20,
                              ),
                            ]
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
                                onDone: onDone,
                                dropdownItems: documentOptions,
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
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  void uploadFile(
      PlatformFile file, int documentType, String documentTitle) async {
    try {
      print('requestBody =======> $file');
      showLoader(context);
      var mintingWallet = appState.userInfo!.getMintingWallets[0];

      Map responseData = await makePutRequestForMultipartDocumentUpload(
        uri: '/v1/tokenization/document',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: mintingWallet.publicKey!,
        file: file,
        tokenizedAssetId: appState.viewData!['id'],
        documentTitle: documentTitle,
        documentType: documentType,
      );

      print('imageUpload response=========>$responseData');

      if (responseData['statusCode'] == 200) {
        String imageUrl = responseData['data'].toString().replaceAll('"', '');
        print('imageUpload response=========>$imageUrl');
        hideLoader(context);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}

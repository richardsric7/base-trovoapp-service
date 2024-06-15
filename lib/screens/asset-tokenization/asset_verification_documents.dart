import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/widgets/loader.dart';
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
  Map<String, List<String>> documentOptions = {
    "proofOfExistenceFiles": <String>[
      "Purchase Receipt",
      "Proof of Address",
      "Other",
    ],
    "proofOfOwnershipFiles": <String>[
      "Title Deed",
      "Bill of sale",
      "Signed transfer of ownership",
      "Certificate of ownership",
      "Other",
    ],
    "assetStatusVerificationFiles": <String>[
      "Inspection reports",
      "Maintenance/repair reports",
      "Photos",
      "Other",
    ],
    "assetCustodianAgreementFiles": <String>[
      "Asset Custodian Agreement",
      "Other",
    ],
    "proofOfAssetManagerFiles": <String>[
      "Asset Management Agreement",
      "Other",
    ],
    "assetProtectionDocumentFiles": <String>[
      "Insurance Policy Document",
      "Bill of sale",
      "Premium payment receipts",
      "Inspection/maintenance report",
      "Other",
    ],
    "assetValuationCertificateFiles": <String>[
      "Asset valaution report",
      "Asset valaution certificate",
      "Other",
    ],
  };

  Map<String, int> documentTypeAndCodes = {
    "proofOfAssetExistence": 1,
    "proofOfAssetOwnership": 2,
    "proofOfAssetStatusVerification": 3,
    "assetCustodianAgreement": 4,
    "proofOfAssetManager": 5,
    "assetProtectionDocument": 6,
    "assetValuationCertificate": 7,
    "assetOwnerGovernmentID": 8,
    "proofOfAssetCondition": 9,
    "thirdPartyTokenizationAgreement": 10,
    "thirdPartyAssetOwnerBusinessRegistration": 11,
    "thirdPartyAssetOwnerProofOfAddress": 12,
    "secApproval": 13,
    "proofOfCompliance": 14,
    "proofOfEnvCompliance": 15,
    "envImpactAssessmentReport": 16,
    "proofOfLegalCounsel": 17,
    "legalAdvisorsContract": 18,
    "proofofMortgagesorLiens": 19,
    "proofofOutstandingLoans": 20,
    "proofofLegalDisputesOnAsset": 21,
  };

  List<DropdownMenuItem<String>> getDocumentOptions(String rel) {
    List<DropdownMenuItem<String>> documentOption = [];
    if (documentOptions[rel] != null) {
      documentOptions[rel]!.forEach((item) {
        documentOption.add(DropdownMenuItem(
            child: Text(
              item,
              overflow: TextOverflow.ellipsis,
            ),
            value: item));
      });
    }
    return documentOption;
  }

  String selectedProofOfExistenceOption = '';
  Map<String, dynamic> proofOfExistenceFiles = {};

  String selectedProofOfOwnershipOption = '';
  Map<String, dynamic> proofOfOwnershipFiles = {};

  String selectedAssetStatusVerificationOption = '';
  Map<String, dynamic> assetStatusVerificationFiles = {};

  String selectedAssetCustodianAgreementOption = '';
  Map<String, dynamic> assetCustodianAgreementFiles = {};

  String selectedProofOfAssetManagerOption = '';
  Map<String, dynamic> proofOfAssetManagerFiles = {};

  String selectedAssetProtectionDocumentOption = '';
  Map<String, dynamic> assetProtectionDocumentFiles = {};

  String selectedAssetValuationCertificateOption = '';
  Map<String, dynamic> assetValuationCertificateFiles = {};

  String selectedProofOfAdditionalCostOutsideValuationOption = '';
  Map<String, dynamic> additionalCostOutsideValuationFiles = {};

  String selectedProofOfAssetConditionOption = '';
  Map<String, dynamic> proofOfAssetConditionFiles = {};

  String selectedThirdPartyTokenizationAgreementOption = '';
  Map<String, dynamic> thirdPartyTokenizationAgreementFiles = {};

  String selectedThirdPartyAssetOwnerBusinessRegOption = '';
  Map<String, dynamic> thirdPartyAssetOwnerBusinessRegFiles = {};

  String selectedThirdPartyAssetOwnerProofOfAddressOption = '';
  Map<String, dynamic> thirdPartyAssetOwnerProofOfAddressFiles = {};

  String selectedSecApprovalOption = '';
  Map<String, dynamic> secRegFiles = {};

  String selectedProofOfComplianceOption = '';
  Map<String, dynamic> proofOfComplianceFiles = {};

  String selectedProofOfEnvComplianceOption = '';
  Map<String, dynamic> proofOfEnvComplianceFiles = {};

  String selectedEnvImpactAssessmentReportOption = '';
  Map<String, dynamic> envImpactAssessmentReportFiles = {};

  String selectedProofofLegalCounselOption = '';
  Map<String, dynamic> proofOfLegalCounselFiles = {};

  String selectedLegalAdvisorsContactOption = '';
  Map<String, dynamic> legalAdvisorsContactFiles = {};

  String selectedProofofMortgagesorLiensOption = '';
  Map<String, dynamic> proofofMortgagesorLiensFiles = {};

  String selectedProofofOutstandingLoansOption = '';
  Map<String, dynamic> proofofOutstandingLoansFiles = {};

  String selectedProofofLegalDisputesOnAssetOption = '';
  Map<String, dynamic> proofofLegalDisputesOnAssetFiles = {};

  late dynamic documents = {};

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
    appState = Provider.of<DataProvider>(context, listen: false);
    initializeData();
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
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfAssetExistence']!,
                    selectedOption.isEmpty
                        ? "Proof Of Asset Existence"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfExistenceOption = selectedOption;
                        proofOfExistenceFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                label: 'Proof of Asset Existence',
                selectedOption: selectedProofOfExistenceOption,
                uploadedFiles: proofOfExistenceFiles,
                documentOptions: getDocumentOptions('proofOfExistenceFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfAssetOwnership']!,
                    selectedOption.isEmpty
                        ? "Proof Of Asset Ownership"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfOwnershipOption = selectedOption;
                        proofOfOwnershipFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofOfOwnershipOption,
                uploadedFiles: proofOfOwnershipFiles,
                label: 'Proof of Ownership',
                documentOptions: getDocumentOptions('proofOfOwnershipFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfAssetStatusVerification']!,
                    selectedOption.isEmpty
                        ? "Proof Of Asset Status Verification"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedAssetStatusVerificationOption = selectedOption;
                        assetStatusVerificationFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedAssetStatusVerificationOption,
                uploadedFiles: assetStatusVerificationFiles,
                label: 'Asset Status Verification',
                documentOptions:
                    getDocumentOptions('assetStatusVerificationFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['assetCustodianAgreement']!,
                    selectedOption.isEmpty
                        ? "Asset Custodian Agreement"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedAssetCustodianAgreementOption = selectedOption;
                        assetCustodianAgreementFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedAssetCustodianAgreementOption,
                uploadedFiles: assetCustodianAgreementFiles,
                label: 'Asset Custodian Agreement',
                documentOptions:
                    getDocumentOptions('assetCustodianAgreementFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfAssetManager']!,
                    selectedOption.isEmpty
                        ? "Proof Of Asset Manager"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfAssetManagerOption = selectedOption;
                        proofOfAssetManagerFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofOfAssetManagerOption,
                uploadedFiles: proofOfAssetManagerFiles,
                label: 'Proof of Asset Manager',
                documentOptions:
                    getDocumentOptions('proofOfAssetManagerFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['assetProtectionDocument']!,
                    selectedOption.isEmpty
                        ? "Asset Protection Document"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedAssetProtectionDocumentOption = selectedOption;
                        assetProtectionDocumentFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedAssetProtectionDocumentOption,
                uploadedFiles: assetProtectionDocumentFiles,
                label: 'Asset Protection Document',
                documentOptions:
                    getDocumentOptions('assetProtectionDocumentFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['assetValuationCertificate']!,
                    selectedOption.isEmpty
                        ? "Asset Valuation Certificate"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedAssetValuationCertificateOption =
                            selectedOption;
                        assetValuationCertificateFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedAssetValuationCertificateOption,
                uploadedFiles: assetValuationCertificateFiles,
                label: 'Asset Valuation Certificate',
                documentOptions:
                    getDocumentOptions('assetValuationCertificateFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['assetOwnerGovernmentID']!,
                    selectedOption.isEmpty
                        ? "Asset Owner Government ID"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfAdditionalCostOutsideValuationOption =
                            selectedOption;
                        additionalCostOutsideValuationFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption:
                    selectedProofOfAdditionalCostOutsideValuationOption,
                uploadedFiles: additionalCostOutsideValuationFiles,
                label: 'Proof of Additional Cost Outside Valuation',
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfAssetCondition']!,
                    selectedOption.isEmpty
                        ? "Proof Of Asset Condition"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfAssetConditionOption = selectedOption;
                        proofOfAssetConditionFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofOfAssetConditionOption,
                uploadedFiles: proofOfAssetConditionFiles,
                label: "Proof of Asset's Condition",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['thirdPartyTokenizationAgreement']!,
                    selectedOption.isEmpty
                        ? "Third Party Tokenization Agreement"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedThirdPartyTokenizationAgreementOption =
                            selectedOption;
                        thirdPartyTokenizationAgreementFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedThirdPartyTokenizationAgreementOption,
                uploadedFiles: thirdPartyTokenizationAgreementFiles,
                label: "Third Party Tokenization Agreement",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes[
                        'thirdPartyAssetOwnerBusinessRegistration']!,
                    selectedOption.isEmpty
                        ? "Third Party Asset Owner Business Registration"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedThirdPartyAssetOwnerBusinessRegOption =
                            selectedOption;
                        thirdPartyAssetOwnerBusinessRegFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedThirdPartyAssetOwnerBusinessRegOption,
                uploadedFiles: thirdPartyAssetOwnerBusinessRegFiles,
                label: "Third Party Asset Owner’s Business Registration",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['thirdPartyAssetOwnerProofOfAddress']!,
                    selectedOption.isEmpty
                        ? "Third Party Asset Owner Proof Of Address"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedThirdPartyAssetOwnerProofOfAddressOption =
                            selectedOption;
                        thirdPartyAssetOwnerProofOfAddressFiles[
                            selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption:
                    selectedThirdPartyAssetOwnerProofOfAddressOption,
                uploadedFiles: thirdPartyAssetOwnerProofOfAddressFiles,
                label: "Third Party Asset Owner’s Proof of Address",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['secApproval']!,
                    selectedOption.isEmpty ? "SEC Approval" : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedSecApprovalOption = selectedOption;
                        secRegFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedSecApprovalOption,
                uploadedFiles: secRegFiles,
                label: "SEC Registration/Tokenization Approval",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfCompliance']!,
                    selectedOption.isEmpty
                        ? "Proof Of Compliance"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfComplianceOption = selectedOption;
                        proofOfComplianceFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofOfComplianceOption,
                uploadedFiles: proofOfComplianceFiles,
                label: "Proof of Compliance with Local Laws and Regulations",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfEnvCompliance']!,
                    selectedOption.isEmpty
                        ? "Proof Of Env. Compliance"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofOfEnvComplianceOption = selectedOption;
                        proofOfEnvComplianceFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofOfEnvComplianceOption,
                uploadedFiles: proofOfEnvComplianceFiles,
                label: "Proof of Compliance with Environmental Standards",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['envImpactAssessmentReport']!,
                    selectedOption.isEmpty
                        ? "Env. Impact Assessment Report"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedEnvImpactAssessmentReportOption =
                            selectedOption;
                        envImpactAssessmentReportFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedEnvImpactAssessmentReportOption,
                uploadedFiles: envImpactAssessmentReportFiles,
                label: "Environmental Impact Assessment Report",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofOfLegalCounsel']!,
                    selectedOption.isEmpty
                        ? "Proof Of Legal Counsel"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofofLegalCounselOption = selectedOption;
                        proofOfLegalCounselFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofofLegalCounselOption,
                uploadedFiles: proofOfLegalCounselFiles,
                label: "Proof of Legal/Financial Counsel",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['legalAdvisorsContract']!,
                    selectedOption.isEmpty
                        ? "Legal Advisors Contract"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofofLegalCounselOption = selectedOption;
                        proofOfLegalCounselFiles[selectedOption] = documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofofLegalCounselOption,
                uploadedFiles: proofOfLegalCounselFiles,
                label: "Legal/Financial Advisors Contact",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofofMortgagesorLiens']!,
                    selectedOption.isEmpty
                        ? "Proof of Mortgages or Liens"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofofMortgagesorLiensOption = selectedOption;
                        proofofMortgagesorLiensFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofofMortgagesorLiensOption,
                uploadedFiles: proofofMortgagesorLiensFiles,
                label: "Proof of Existing Mortgages or Liens on Asset",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofofOutstandingLoans']!,
                    selectedOption.isEmpty
                        ? "Proof of Outstanding Loans"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofofOutstandingLoansOption = selectedOption;
                        proofofOutstandingLoansFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofofOutstandingLoansOption,
                uploadedFiles: proofofOutstandingLoansFiles,
                label: "Proof of Outstanding Loans on Asset",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(
                    file,
                    documentTypeAndCodes['proofofLegalDisputesOnAsset']!,
                    selectedOption.isEmpty
                        ? "Proof of Legal Disputes On Asset"
                        : selectedOption,
                    onDone: (documentUrl) {
                      setState(() {
                        selectedProofofLegalDisputesOnAssetOption =
                            selectedOption;
                        proofofLegalDisputesOnAssetFiles[selectedOption] =
                            documentUrl;
                      });
                    },
                  );
                },
                selectedOption: selectedProofofLegalDisputesOnAssetOption,
                uploadedFiles: proofofLegalDisputesOnAssetFiles,
                label: "Proof of  Legal Disputes or Encumbrances on Asset",
                documentOptions: getDocumentOptions('String rel')),
            SizedBox(
              height: height / 30,
            ),
            Button(
              'Continue',
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

  Widget proofDocumentItem({
    required String label,
    required List<DropdownMenuItem<String>> documentOptions,
    required String selectedOption,
    required Map<String, dynamic> uploadedFiles,
    required void Function(String selectedOption, PlatformFile file) onDone,
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
                      for (var item in uploadedFiles.keys) ...[
                        Container(
                          width: width / 1.27,
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              if (item.isNotEmpty) ...[
                                Text(
                                  truncate(item, length: 18),
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontbody,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                              TextButton(
                                onPressed: () {
                                  var fileUrl = uploadedFiles[item]
                                          ['documentUrl']
                                      .toString();
                                  if (fileUrl.isNotEmpty &&
                                      fileUrl.endsWith('.pdf')) {
                                    appState.viewData!['pdfUrl'] =
                                        uploadedFiles[item];
                                    appState.currentAction = PageAction(
                                        state: PageState.addPage,
                                        page: PdfViewPageConfig);

                                    return;
                                  }

                                  appState.goToWebView(fileUrl);
                                },
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Text(
                                      truncatePublicKey(
                                          uploadedFiles[item]['documentUrl']!),
                                      style: TextStyle(
                                        decoration: TextDecoration.underline,
                                        fontSize: 12,
                                        fontFamily: fontbody,
                                        color: notifier.getbluewhitecolor,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              IconButton(
                                icon: Icon(
                                  CupertinoIcons.delete,
                                  size: 20,
                                ),
                                onPressed: (() async {
                                  // setState(() async {
                                  deleteFile(
                                      uploadedFiles[item]['id'].toString());
                                  // uploadedFiles.remove(item);
                                  // await fetchTokenizationData();
                                  // });
                                }),
                              )
                            ],
                          ),
                        ),
                      ],
                      Container(
                        width: width / 1.27,
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            TextButton(
                              onPressed: () {
                                // showDocumentUploadPopup(
                                //   context,
                                //   label,
                                //   onDone: onDone,
                                //   dropdownItems: documentOptions,
                                // );
                                fetchCurrentTokenizationInfo();
                              },
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

  void initializeData() {
    inspect(appState.viewData);
    documents = appState.viewData!['AssetTokenizationDocuments'];
    inspect(documents);
    for (var item in documents) {
      switch (item['documentType']) {
        case 1:
          proofOfExistenceFiles[item['documentTitle']] = item;
          break;
        case 2:
          proofOfOwnershipFiles[item['documentTitle']] = item;
          break;
        case 3:
          assetStatusVerificationFiles[item['documentTitle']] = item;
          break;
        case 4:
          assetCustodianAgreementFiles[item['documentTitle']] = item;
          break;
        case 5:
          proofOfAssetManagerFiles[item['documentTitle']] = item;
          break;
        case 6:
          assetProtectionDocumentFiles[item['documentTitle']] = item;
          break;
        case 7:
          assetValuationCertificateFiles[item['documentTitle']] = item;
          break;
        case 8:
          additionalCostOutsideValuationFiles[item['documentTitle']] = item;
          break;
        case 9:
          proofOfAssetConditionFiles[item['documentTitle']] = item;
          break;
        case 10:
          thirdPartyTokenizationAgreementFiles[item['documentTitle']] = item;
          break;
        case 11:
          thirdPartyAssetOwnerBusinessRegFiles[item['documentTitle']] = item;
          break;
        case 12:
          thirdPartyAssetOwnerProofOfAddressFiles[item['documentTitle']] = item;
          break;
        case 13:
          secRegFiles[item['documentTitle']] = item;
          break;
        case 14:
          proofOfComplianceFiles[item['documentTitle']] = item;
          break;
        case 15:
          proofOfEnvComplianceFiles[item['documentTitle']] = item;
          break;
        case 16:
          envImpactAssessmentReportFiles[item['documentTitle']] = item;
          break;
        case 17:
          proofOfLegalCounselFiles[item['documentTitle']] = item;
          break;
        case 18:
          legalAdvisorsContactFiles[item['documentTitle']] = item;
          break;
        case 19:
          proofofMortgagesorLiensFiles[item['documentTitle']] = item;
          break;
        case 20:
          proofofOutstandingLoansFiles[item['documentTitle']] = item;
          break;
        case 21:
          proofofLegalDisputesOnAssetFiles[item['documentTitle']] = item;
          break;
        default:
      }
    }
  }

  Future<Map> fetchCurrentTokenizationInfo() async {
    try {
      var uri = '/v1/tokenization/detail/${appState.viewData!['id']}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      print('===============> response ${responseData}');
      if (responseData['statusCode'] == 200) {
        print('success');
        return responseData['data'];
      } else {
        print('success');
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  Future<void> uploadFile(
    PlatformFile file,
    int documentType,
    String documentTitle, {
    required void Function(String documentUrl) onDone,
  }) async {
    try {
      // print('requestBody =======> $file');
      showLoader(context);
      var mintingWalletPublicKey = appState.activeTokenizationWalletPublicKey!;

      Map responseData = await makePutRequestForMultipartDocumentUpload(
        uri: '/v1/tokenization/document',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: mintingWalletPublicKey,
        file: file,
        tokenizedAssetId: appState.viewData!['id'],
        documentTitle: documentTitle.toLowerCase().replaceAll(' ', '-'),
        documentType: documentType,
      );

      // print('imageUpload response=========>$responseData');

      if (responseData['statusCode'] == 200) {
        String imageUrl = responseData['data'].toString().replaceAll('"', '');
        // print('imageUpload response=========>$imageUrl');
        onDone(imageUrl);
        // appState.viewData = await fetchTokenizationData();
        hideLoader(context);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
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
      print('object ${mintingWalletPublicKey}');
      Map responseData = await makeDeleteRequest(
        uri: '/v1/tokenization/document/$documentId',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: mintingWalletPublicKey,
        body: jsonEncode(requestBody),
      );

      print("response ============> ${responseData}");
      if (responseData['statusCode'] == 200) {
        print("response ============> ${responseData['data']}");
        hideLoader(context);
      } else {
        popup(context,
            title: "error".tr(), message: responseData['data']['message']);
        hideLoader(context);
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
}

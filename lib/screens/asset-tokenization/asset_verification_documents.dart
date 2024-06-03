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
  Map<String, String> proofOfExistenceFiles = {};

  String selectedProofOfOwnershipOption = '';
  Map<String, String> proofOfOwnershipFiles = {};

  String selectedAssetStatusVerificationOption = '';
  Map<String, String> assetStatusVerificationFiles = {};

  String selectedAssetCustodianAgreementOption = '';
  Map<String, String> assetCustodianAgreementFiles = {};

  String selectedProofOfAssetManagerOption = '';
  Map<String, String> proofOfAssetManagerFiles = {};

  String selectedAssetProtectionDocumentOption = '';
  Map<String, String> assetProtectionDocumentFiles = {};

  String selectedAssetValuationCertificateOption = '';
  Map<String, String> assetValuationCertificateFiles = {};

  String selectedProofOfAdditionalCostOutsideValuationOption = '';
  Map<String, String> additionalCostOutsideValuationFiles = {};

  String selectedProofOfAssetConditionOption = '';
  Map<String, String> proofOfAssetConditionFiles = {};

  String selectedThirdPartyTokenizationAgreementOption = '';
  Map<String, String> thirdPartyTokenizationAgreementFiles = {};

  String selectedThirdPartyAssetOwnerBusinessRegOption = '';
  Map<String, String> thirdPartyAssetOwnerBusinessRegFiles = {};

  String selectedThirdPartyAssetOwnerProofOfAddressOption = '';
  Map<String, String> thirdPartyAssetOwnerProofOfAddressFiles = {};

  String selectedSecApprovalOption = '';
  Map<String, String> secRegFiles = {};

  String selectedProofOfComplianceOption = '';
  Map<String, String> proofOfComplianceFiles = {};

  String selectedProofOfEnvComplianceOption = '';
  Map<String, String> proofOfEnvComplianceFiles = {};

  String selectedEnvImpactAssessmentReportOption = '';
  Map<String, String> envImpactAssessmentReportFiles = {};

  String selectedProofofLegalCounselOption = '';
  Map<String, String> proofOfLegalCounselFiles = {};

  String selectedLegalAdvisorsContactOption = '';
  Map<String, String> legalAdvisorsContactFiles = {};

  String selectedProofofMortgagesorLiensOption = '';
  Map<String, String> proofofMortgagesorLiensFiles = {};

  String selectedProofofOutstandingLoansOption = '';
  Map<String, String> proofofOutstandingLoansFiles = {};

  String selectedProofofLegalDisputesOnAssetOption = '';
  Map<String, String> proofofLegalDisputesOnAssetFiles = {};

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
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  uploadFile(file, 1, 'ProofOfAssetExistence');
                  setState(() {
                    selectedProofOfExistenceOption = selectedOption;
                    proofOfExistenceFiles[selectedOption] = file.name;
                  });
                },
                label: 'Proof of Asset Existence',
                selectedOption: selectedProofOfExistenceOption,
                uploadedFiles: proofOfExistenceFiles,
                documentOptions: getDocumentOptions('proofOfExistenceFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofOfOwnershipOption = selectedOption;
                    proofOfOwnershipFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofOfOwnershipOption,
                uploadedFiles: proofOfOwnershipFiles,
                label: 'Proof of Ownership',
                documentOptions: getDocumentOptions('proofOfOwnershipFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedAssetStatusVerificationOption = selectedOption;
                    assetStatusVerificationFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedAssetStatusVerificationOption,
                uploadedFiles: assetStatusVerificationFiles,
                label: 'Asset Status Verification',
                documentOptions:
                    getDocumentOptions('assetStatusVerificationFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedAssetCustodianAgreementOption = selectedOption;
                    assetCustodianAgreementFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedAssetCustodianAgreementOption,
                uploadedFiles: assetCustodianAgreementFiles,
                label: 'Asset Custodian Agreement',
                documentOptions:
                    getDocumentOptions('assetCustodianAgreementFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofOfAssetManagerOption = selectedOption;
                    proofOfAssetManagerFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofOfAssetManagerOption,
                uploadedFiles: proofOfAssetManagerFiles,
                label: 'Proof of Asset Manager',
                documentOptions:
                    getDocumentOptions('proofOfAssetManagerFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedAssetProtectionDocumentOption = selectedOption;
                    assetProtectionDocumentFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedAssetProtectionDocumentOption,
                uploadedFiles: assetProtectionDocumentFiles,
                label: 'Asset Protection Document',
                documentOptions:
                    getDocumentOptions('assetProtectionDocumentFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedAssetValuationCertificateOption = selectedOption;
                    assetValuationCertificateFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedAssetValuationCertificateOption,
                uploadedFiles: assetValuationCertificateFiles,
                label: 'Asset Valuation Certificate',
                documentOptions:
                    getDocumentOptions('assetValuationCertificateFiles')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofOfAdditionalCostOutsideValuationOption =
                        selectedOption;
                    additionalCostOutsideValuationFiles[selectedOption] =
                        file.name;
                  });
                },
                selectedOption:
                    selectedProofOfAdditionalCostOutsideValuationOption,
                uploadedFiles: additionalCostOutsideValuationFiles,
                label: 'Proof of Additional Cost Outside Valuation',
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofOfAssetConditionOption = selectedOption;
                    proofOfAssetConditionFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofOfAssetConditionOption,
                uploadedFiles: proofOfAssetConditionFiles,
                label: "Proof of Asset's Condition",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedThirdPartyTokenizationAgreementOption =
                        selectedOption;
                    thirdPartyTokenizationAgreementFiles[selectedOption] =
                        file.name;
                  });
                },
                selectedOption: selectedThirdPartyTokenizationAgreementOption,
                uploadedFiles: thirdPartyTokenizationAgreementFiles,
                label: "Third Party Tokenization Agreement",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedThirdPartyAssetOwnerBusinessRegOption =
                        selectedOption;
                    thirdPartyAssetOwnerBusinessRegFiles[selectedOption] =
                        file.name;
                  });
                },
                selectedOption: selectedThirdPartyAssetOwnerBusinessRegOption,
                uploadedFiles: thirdPartyAssetOwnerBusinessRegFiles,
                label: "Third Party Asset Owner’s Business Registration",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedThirdPartyAssetOwnerProofOfAddressOption =
                        selectedOption;
                    thirdPartyAssetOwnerProofOfAddressFiles[selectedOption] =
                        file.name;
                  });
                },
                selectedOption:
                    selectedThirdPartyAssetOwnerProofOfAddressOption,
                uploadedFiles: thirdPartyAssetOwnerProofOfAddressFiles,
                label: "Third Party Asset Owner’s Proof of Address",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedSecApprovalOption = selectedOption;
                    secRegFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedSecApprovalOption,
                uploadedFiles: secRegFiles,
                label: "SEC Registration/Tokenization Approval",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofOfComplianceOption = selectedOption;
                    proofOfComplianceFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofOfComplianceOption,
                uploadedFiles: proofOfComplianceFiles,
                label: "Proof of Compliance with Local Laws and Regulations",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofOfEnvComplianceOption = selectedOption;
                    proofOfEnvComplianceFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofOfEnvComplianceOption,
                uploadedFiles: proofOfEnvComplianceFiles,
                label: "Proof of Compliance with Environmental Standards",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedEnvImpactAssessmentReportOption = selectedOption;
                    envImpactAssessmentReportFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedEnvImpactAssessmentReportOption,
                uploadedFiles: envImpactAssessmentReportFiles,
                label: "Environmental Impact Assessment Report",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofofLegalCounselOption = selectedOption;
                    proofOfLegalCounselFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofofLegalCounselOption,
                uploadedFiles: proofOfLegalCounselFiles,
                label: "Proof of Legal/Financial Counsel",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofofLegalCounselOption = selectedOption;
                    proofOfLegalCounselFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofofLegalCounselOption,
                uploadedFiles: proofOfLegalCounselFiles,
                label: "Legal/Financial Advisors Contact",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofofMortgagesorLiensOption = selectedOption;
                    proofofMortgagesorLiensFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofofMortgagesorLiensOption,
                uploadedFiles: proofofMortgagesorLiensFiles,
                label: "Proof of Existing Mortgages or Liens on Asset",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofofOutstandingLoansOption = selectedOption;
                    proofofOutstandingLoansFiles[selectedOption] = file.name;
                  });
                },
                selectedOption: selectedProofofOutstandingLoansOption,
                uploadedFiles: proofofOutstandingLoansFiles,
                label: "Proof of Outstanding Loans on Asset",
                documentOptions: getDocumentOptions('String rel')),
            proofDocumentItem(
                onDone: (String selectedOption, PlatformFile file) {
                  setState(() {
                    selectedProofofLegalDisputesOnAssetOption = selectedOption;
                    proofofLegalDisputesOnAssetFiles[selectedOption] =
                        file.name;
                  });
                },
                selectedOption: selectedProofofLegalDisputesOnAssetOption,
                uploadedFiles: proofofLegalDisputesOnAssetFiles,
                label: "Proof of  Legal Disputes or Encumbrances on Asset",
                documentOptions: getDocumentOptions('String rel')),
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

  Widget proofDocumentItem({
    required String label,
    required List<DropdownMenuItem<String>> documentOptions,
    required String selectedOption,
    required Map<String, String> uploadedFiles,
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
                                  item,
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontbody,
                                    color: notifier.getbluewhitecolor,
                                  ),
                                ),
                              ],
                              Text(
                                truncate(uploadedFiles[item]!, length: 15),
                                style: TextStyle(
                                  decoration: TextDecoration.underline,
                                  fontSize: 12,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                              IconButton(
                                icon: Icon(
                                  CupertinoIcons.delete,
                                  size: 20,
                                ),
                                onPressed: (() {
                                  setState(() {
                                    uploadedFiles.remove(item);
                                  });
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
                              onPressed: () => showDocumentUploadPopup(
                                context,
                                label,
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

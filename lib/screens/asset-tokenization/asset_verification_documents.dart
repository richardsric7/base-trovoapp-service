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

  Map<String, Map<String, dynamic>> get initDocTypeAndCodes => {
        "proofOfAssetExistence": {
          'name': 'Proof of Asset Existence',
          'documentType': '1',
          'options': <String>[
            "Purchase Receipt",
            "Proof of Address",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfAssetOwnership": {
          'name': 'Proof Of Asset Ownership',
          'documentType': '2',
          'options': <String>[
            "Title Deed",
            "Bill of sale",
            "Signed transfer of ownership",
            "Certificate of ownership",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfAssetStatusVerification": {
          'name': 'Proof Of Asset Status Verification',
          'documentType': '3',
          'options': <String>[
            "Inspection reports",
            "Maintenance/repair reports",
            "Photos",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "assetCustodianAgreement": {
          'name': 'Asset Custodian Agreement',
          'documentType': '4',
          'options': <String>[
            "Asset Custodian Agreement",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfAssetManager": {
          'name': 'Proof of Asset Manager',
          'documentType': '5',
          'options': <String>[
            "Asset Management Agreement",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "assetProtectionDocument": {
          'name': 'Asset Protection Document',
          'documentType': '6',
          'options': <String>[
            "Insurance Policy Document",
            "Bill of sale",
            "Premium payment receipts",
            "Inspection/maintenance report",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "assetValuationCertificate": {
          'name': 'Asset Valuation Certificate',
          'documentType': '7',
          'options': <String>[
            "Asset valaution report",
            "Asset valaution certificate",
            "Other",
          ],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "additionalCostOutsideValuation": {
          'name': 'Proof of Additional Cost Outside Valuation',
          'documentType': '8',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfAssetCondition": {
          'name': 'Proof of Asset\'s Condition',
          'documentType': '9',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "thirdPartyTokenizationAgreement": {
          'name': 'Third Party Tokenization Agreement',
          'documentType': '10',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "thirdPartyAssetOwnerBusinessRegistration": {
          'name': 'Third Party Asset Owner Business Registration',
          'documentType': '11',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "thirdPartyAssetOwnerProofOfAddress": {
          'name': 'Third Party Asset Owner Proof Of Address',
          'documentType': '12',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "secApproval": {
          'name': 'SEC Registration/Tokenization Approval',
          'documentType': '13',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfCompliance": {
          'name': 'Proof of Compliance with Local Laws and Regulations',
          'documentType': '14',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfEnvCompliance": {
          'name': 'Proof of Compliance with Environmental Standards',
          'documentType': '15',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "envImpactAssessmentReport": {
          'name': 'Environmental Impact Assessment Report',
          'documentType': '16',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofOfLegalCounsel": {
          'name': 'Proof of Legal/Financial Counsel',
          'documentType': '17',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "legalAdvisorsContact": {
          'name': 'Legal/Financial Advisors Contact',
          'documentType': '18',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofofMortgagesorLiens": {
          'name': 'Proof of Mortgages or Liens',
          'documentType': '19',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofofOutstandingLoans": {
          'name': 'Proof of Outstanding Loans',
          'documentType': '20',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
        "proofofLegalDisputesOnAsset": {
          'name': 'Proof of  Legal Disputes or Encumbrances on Asset',
          'documentType': '21',
          'options': <String>[],
          'files': <String, dynamic>{},
          'selectedFileOption': '',
        },
      };

  Map<String, Map<String, dynamic>> documentTypeAndCodes = {};

  List<DropdownMenuItem<String>> getDocumentOptions(
      List<String> documentOptions) {
    List<DropdownMenuItem<String>> documentOption = [];
    if (documentOptions.length > 0) {
      documentOptions.forEach((item) {
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

  late dynamic documents = {};

  List<DropdownMenuItem<String>> get getDocumentsList {
    List<DropdownMenuItem<String>> documentOptions = [];

    documentTypeAndCodes.forEach((key, value) {
      documentOptions.add(
        DropdownMenuItem(
          child: Text(
            value['name'].toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: key,
        ),
      );
    });

    return documentOptions;
  }

  Map<String, Widget> selectedDocuments = {};

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
    notifier = Provider.of<ColorNotifier>(context, listen: false);
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
            if (documentTypeAndCodes.isNotEmpty) ...[
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 10.0),
                child: dropdown(
                  (value) {
                    var val = documentTypeAndCodes[value] as dynamic;
                    setState(() {
                      var files = val['files'];
                      selectedDocuments[val['documentType']] =
                          proofDocumentItem(
                        onDone: (String selectedOption, PlatformFile file) {
                          uploadFile(
                            file,
                            val['documentType'],
                            selectedOption.isEmpty
                                ? val['name']
                                : selectedOption,
                          );
                        },
                        label: val['name'],
                        selectedOption: val['selectedFileOption'],
                        uploadedFiles: files,
                        documentOptions: getDocumentOptions(val['options']),
                      );
                    });
                  },
                  getDocumentsList,
                  null,
                  'Select document',
                  context,
                  null,
                ),
              ),
            ],
            for (var item in selectedDocuments.entries) ...[
              item.value,
            ],
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
                                    appState.pdfUrl = fileUrl;
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
                                  deleteFile(
                                      uploadedFiles[item]['id'].toString());
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
                                showDocumentUploadPopup(
                                  context,
                                  label,
                                  onDone: onDone,
                                  dropdownItems: documentOptions,
                                );
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
    documentTypeAndCodes = initDocTypeAndCodes;
    selectedDocuments.clear();
    // inspect(appState.viewData!['AssetTokenizationDocuments']);
    print(selectedDocuments);
    documents = appState.viewData!['AssetTokenizationDocuments'] ?? [];
    for (var item in documents) {
      switch (item['documentType']) {
        case '1':
          var val = documentTypeAndCodes['proofOfAssetExistence'];
          val?['files'][item['documentTitle'].toString()] = item;
          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '2':
          var val = documentTypeAndCodes['proofOfAssetOwnership'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '3':
          var val = documentTypeAndCodes['proofOfAssetStatusVerification'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '4':
          var val = documentTypeAndCodes['assetCustodianAgreement'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '5':
          var val = documentTypeAndCodes['proofOfAssetManager'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '6':
          var val = documentTypeAndCodes['assetProtectionDocument'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '7':
          var val = documentTypeAndCodes['assetValuationCertificate'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '8':
          var val = documentTypeAndCodes['additionalCostOutsideValuation'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '9':
          var val = documentTypeAndCodes['proofOfAssetCondition'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '10':
          var val = documentTypeAndCodes['thirdPartyTokenizationAgreement'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '11':
          var val =
              documentTypeAndCodes['thirdPartyAssetOwnerBusinessRegistration'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '12':
          var val = documentTypeAndCodes['thirdPartyAssetOwnerProofOfAddress'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '13':
          var val = documentTypeAndCodes['secApproval'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '14':
          var val = documentTypeAndCodes['proofOfCompliance'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '15':
          var val = documentTypeAndCodes['proofOfEnvCompliance'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '16':
          var val = documentTypeAndCodes['envImpactAssessmentReport'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '17':
          var val = documentTypeAndCodes['proofOfLegalCounsel'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '18':
          var val = documentTypeAndCodes['legalAdvisorsContact'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '19':
          var val = documentTypeAndCodes['proofofMortgagesorLiens'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '20':
          var val = documentTypeAndCodes['proofofOutstandingLoans'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        case '21':
          var val = documentTypeAndCodes['proofofLegalDisputesOnAsset'];
          val?['files'][item['documentTitle'].toString()] = item;

          selectedDocuments[item['documentType']] = proofDocumentItem(
            onDone: (String selectedOption, PlatformFile file) {
              uploadFile(
                file,
                val['documentType'],
                selectedOption.isEmpty ? val['name'] : selectedOption,
              );
            },
            label: val!['name'],
            selectedOption: val['selectedFileOption'],
            uploadedFiles: val['files'],
            documentOptions: getDocumentOptions(
              val['options'],
            ),
          );
          break;
        default:
      }
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

  Future<void> uploadFile(
    PlatformFile file,
    String documentType,
    String documentTitle,
  ) async {
    try {
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

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
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
      Map responseData = await makeDeleteRequest(
        uri: '/v1/tokenization/document/$documentId',
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
}

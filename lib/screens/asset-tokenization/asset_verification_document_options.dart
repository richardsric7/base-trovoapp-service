import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetVerificationDocumentOptionsView extends StatefulWidget {
  const AssetVerificationDocumentOptionsView({Key? key}) : super(key: key);

  @override
  State<AssetVerificationDocumentOptionsView> createState() =>
      _AssetVerificationDocumentOptionsView();
}

class _AssetVerificationDocumentOptionsView
    extends State<AssetVerificationDocumentOptionsView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  int documentUploadCount = 1;
  int totalDocuments = 5;

  Map<String, Map<String, dynamic>> documentTypeAndCodes = {};

  List<DropdownMenuItem<String>> getDocumentOptions(
    List<String> documentOptions,
  ) {
    List<DropdownMenuItem<String>> documentOption = [];
    if (documentOptions.length > 0) {
      documentOptions.forEach((item) {
        documentOption.add(
          DropdownMenuItem(
            child: Text(item, overflow: TextOverflow.ellipsis),
            value: item,
          ),
        );
      });
    }
    return documentOption;
  }

  late dynamic documents = {};

  List<DropdownMenuItem<String>> getDocumentsList(bool isSelected) {
    List<DropdownMenuItem<String>> documentOptions = [];

    documentTypeAndCodes.forEach((key, value) {
      if (value['required'] == 1) {
        documentOptions.add(
          DropdownMenuItem(
            child: Row(
              children: [
                Container(
                  constraints: BoxConstraints(maxWidth: 250),
                  child: Padding(
                    padding: EdgeInsets.symmetric(vertical: isSelected ? 0 : 6),
                    child: Text(
                      value['name'].toString(),
                      overflow: isSelected
                          ? TextOverflow.ellipsis
                          : TextOverflow.visible,
                    ),
                  ),
                ),
                if (selectedDocuments[value['documentType']] != null) ...[
                  SizedBox(width: 3),
                  Icon(Icons.check, size: 18, color: notifier.getbluecolor),
                ],
              ],
            ),
            value: key,
          ),
        );
      }
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
              appState.viewData?['title'],
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: Row(
                children: [
                  SizedBox(
                    width: 350,
                    child: Text(
                      "Select a file to upload. Files marked with * are required to proceed with your application",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
            ),

            // if (documentTypeAndCodes.isNotEmpty) ...[
            //   Padding(
            //     padding: const EdgeInsets.symmetric(horizontal: 10.0),
            //     child: dropdown(
            //       (value) {
            //         var val = documentTypeAndCodes[value] as dynamic;
            //         setState(() {
            //           var files = val['files'];
            //           selectedDocuments[val['documentType']] =
            //               proofDocumentItem(
            //                 onDone: (String selectedOption, PlatformFile file) {
            //                   uploadFile(
            //                     file,
            //                     val['documentType'],
            //                     selectedOption.isEmpty
            //                         ? val['name']
            //                         : selectedOption,
            //                   );
            //                 },
            //                 label: val['name'],
            //                 selectedOption: val['selectedFileOption'],
            //                 uploadedFiles: files,
            //                 documentOptions: getDocumentOptions(val['options']),
            //               );
            //         });
            //       },
            //       getDocumentsList(false),
            //       null,
            //       'Select document',
            //       context,
            //       (context) {
            //         return getDocumentsList(true);
            //       },
            //     ),
            //   ),
            // ],
            // for (var item in selectedDocuments.entries) ...[item.value],
            SizedBox(height: 20),
            proofItemCard(
              title: "Asset Protection Documents",
              onTap: () {
                uploadFileBottomSheet(title: 'Asset Protection Documents');
              },
            ),
            SizedBox(height: height / 30),
            Button(
              'Continue',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                Navigator.of(context).pop();
              },
            ),
            SizedBox(height: height / 10),
          ],
        ),
      ),
    );
  }

  void uploadFileBottomSheet({required String title}) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: notifier.getwihitecolor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
          top: Radius.circular(16),
        ), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.75,
          minChildSize: 0.25,
          expand: false,
          builder: (context, scrollController) {
            return SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  SizedBox(height: 10),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 14.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text(
                          title,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 18,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(Icons.cancel_outlined),
                          style: ButtonStyle(
                            padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            minimumSize: WidgetStatePropertyAll(Size.zero),
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: 20),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 10.0),
                    child: Row(
                      children: [
                        SizedBox(
                          width: 350,
                          child: Text(
                            "Select a file to upload.",
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                      ],
                    ),
                  ),
                  Row(
                    children: [
                      GestureDetector(
                        onTap: () {
                          getFile();
                        },
                        child: Column(
                          children: [
                            SizedBox(height: height / 50),
                            Padding(
                              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
                              child: Container(
                                width: 320,
                                padding: EdgeInsets.all(20),
                                decoration: BoxDecoration(
                                  border: Border.all(
                                    color: notifier.getbluewhitecolor,
                                    width: 1,
                                  ),
                                  borderRadius: const BorderRadius.all(
                                    Radius.circular(15.0),
                                  ),
                                  color: notifier.getwihitecolor,
                                ),
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  spacing: 10,
                                  children: [
                                    Icon(
                                      Icons.file_upload_outlined,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                    Text(
                                      "Upload file here",
                                      textAlign: TextAlign.center,
                                      style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontsemibold,
                                        fontSize: 12.sp,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 60),
                ],
              ),
            );
          },
        );
      },
    );
  }

  Widget proofItemCard({
    required String title,
    required void Function() onTap,
    bool isRequired = true,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 10),
        child: Container(
          width: double.infinity,
          padding: EdgeInsets.symmetric(horizontal: 5, vertical: 20),
          decoration: BoxDecoration(
            border: Border.all(color: notifier.getsplashgrey, width: 1),
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: wihitecolor,
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  if (isRequired) ...[
                    Text(
                      "*",
                      style: TextStyle(
                        fontSize: 16,
                        fontFamily: fontsemibold,
                        color: Colors.red,
                      ),
                    ),
                  ],
                  Text(
                    "Proof of Asset Address",
                    style: TextStyle(
                      fontSize: 16,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
            ],
          ),
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
              padding: const EdgeInsets.symmetric(
                horizontal: 20.0,
                vertical: 15.0,
              ),
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
                      SizedBox(height: height / 70),
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
                                  var fileUrl =
                                      uploadedFiles[item]['documentUrl']
                                          .toString();
                                  if (fileUrl.isNotEmpty &&
                                      fileUrl.endsWith('.pdf')) {
                                    appState.pdfUrl = fileUrl;
                                    appState.currentAction = PageAction(
                                      state: PageState.addPage,
                                      page: PdfViewPageConfig,
                                    );

                                    return;
                                  }

                                  appState.goToWebView(fileUrl);
                                },
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Text(
                                      truncatePublicKey(
                                        uploadedFiles[item]['documentUrl']!,
                                      ),
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
                                  color: Colors.red,
                                ),
                                onPressed: (() async {
                                  deleteFile(
                                    uploadedFiles[item]['id'].toString(),
                                  );
                                }),
                              ),
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
                                  SizedBox(width: width / 20),
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
    selectedDocuments.clear();
    var assetAlreadyExists = appState.viewData!['assetAlreadyExists'];
    var assetOwnership = appState.viewData!['ownershipType'];
    var requiredProofOfContributedValue =
        assetAlreadyExists == 1 &&
            appState.viewData!['assetOwnerRetainedOrContributedValue'] > 0
        ? 1
        : 0;
    documentTypeAndCodes = {
      '1': {
        'name': 'Proof of Asset Existence',
        'documentType': '1',
        'options': <String>["Purchase Receipt", "Proof of Address", "Other"],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      '2': {
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
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "3": {
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
        'required': 1,
      },
      "6": {
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
        'required': 1,
      },
      "7": {
        'name': 'Asset Valuation Certificate',
        'documentType': '7',
        'options': <String>[
          "Asset valaution report",
          "Asset valaution certificate",
          "Other",
        ],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "8": {
        'name': 'Proof of Additional Cost Outside Valuation',
        'documentType': '8',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "9": {
        'name': 'Proof of Asset\'s Condition',
        'documentType': '9',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "10": {
        'name': 'Third Party Tokenization Agreement',
        'documentType': '10',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetOwnership == 'THIRD-PARTY' ? 1 : 0,
      },
      "11": {
        'name': 'Third Party Asset Owner Business Registration',
        'documentType': '11',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetOwnership == 'THIRD-PARTY' ? 1 : 0,
      },
      "12": {
        'name': 'Third Party Asset Owner Proof Of Address',
        'documentType': '12',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetOwnership == 'THIRD-PARTY' ? 1 : 0,
      },
      "14": {
        'name': 'Proof of Compliance with Local Laws and Regulations',
        'documentType': '14',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "15": {
        'name': 'Proof of Compliance with Environmental Standards',
        'documentType': '15',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "16": {
        'name': 'Environmental Impact Assessment Report',
        'documentType': '16',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "17": {
        'name': 'Proof of Legal/Financial Counsel',
        'documentType': '17',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "18": {
        'name': 'Legal/Financial Advisors Contact',
        'documentType': '18',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "19": {
        'name': 'Proof of Mortgages or Liens',
        'documentType': '19',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 0,
      },
      "20": {
        'name': 'Proof of Outstanding Loans',
        'documentType': '20',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 0,
      },
      "21": {
        'name': 'Proof of Legal Disputes on Asset',
        'documentType': '21',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 0,
      },
      "22": {
        'name': 'Approved Project Budget',
        'documentType': '22',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 0 ? 1 : 0,
      },
      "23": {
        'name': 'Proof of Contribution from Sponsor',
        'documentType': '23',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': requiredProofOfContributedValue,
      },
      "24": {
        'name': 'Title Deeds or Certificates',
        'documentType': '24',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "25": {
        'name': 'Original Ownership Agreements',
        'documentType': '25',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "26": {
        'name':
            'Government Issued ID  of Original Asset Owner(s) or Authorized Pepresentatives',
        'documentType': '26',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "27": {
        'name': 'Infrastructure Inspection Reports',
        'documentType': '27',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "28": {
        'name': 'Maintenance Records',
        'documentType': '28',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "29": {
        'name': 'Asset Financial Performance Report',
        'documentType': '29',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': assetAlreadyExists == 1 ? 1 : 0,
      },
      "30": {
        'name': 'Statutory Licenses and Permits/Compliance Certificates',
        'documentType': '30',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "31": {
        'name': 'Project Proposal Concept Note/Business Case',
        'documentType': '31',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "32": {
        'name': 'Market Analysis',
        'documentType': '32',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "33": {
        'name': 'Feasibility Study',
        'documentType': '33',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "34": {
        'name': 'Business Plan',
        'documentType': '34',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "35": {
        'name': 'Project Financial Model',
        'documentType': '35',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "36": {
        'name':
            'Legal Documentation for Development, Construction, and Operation',
        'documentType': '36',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "37": {
        'name': 'Contracts',
        'documentType': '37',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "38": {
        'name': 'Permits',
        'documentType': '38',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "39": {
        'name': 'licenses',
        'documentType': '39',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "40": {
        'name':
            'Construction Plans and Specifications for Development, Construction, and Operation',
        'documentType': '40',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "41": {
        'name': 'Detailed Plans',
        'documentType': '41',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "42": {
        'name': 'Drawings',
        'documentType': '42',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "43": {
        'name': 'Specifications',
        'documentType': '43',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "44": {
        'name': 'Site Survey Reports',
        'documentType': '44',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "45": {
        'name': 'Risk Management Plan',
        'documentType': '45',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "46": {
        'name': 'Operational Plans and Procedures',
        'documentType': '46',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "47": {
        'name': 'Project Timeline',
        'documentType': '47',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "48": {
        'name': 'Project Team and Partnerships',
        'documentType': '48',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "49": {
        'name': 'Project Developers',
        'documentType': '49',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "50": {
        'name': 'Consultants',
        'documentType': '50',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "51": {
        'name': 'Other Key Stakeholders Involved in the Project',
        'documentType': '51',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "52": {
        'name': 'Regulatory and Compliance Documentation',
        'documentType': '52',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "53": {
        'name': 'Insurance Policies',
        'documentType': '53',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "54": {
        'name': 'Third-Party Reports and Due Diligence',
        'documentType': '54',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "55": {
        'name': 'Proof of Additional Cost Incurred Outside Valuation',
        'documentType': '55',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "56": {
        'name': 'Comprehensive Insurance',
        'documentType': '56',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "57": {
        'name': 'Contractual Protections',
        'documentType': '57',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "58": {
        'name': 'Revenue Guarantees',
        'documentType': '58',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "59": {
        'name': 'Performance Bond',
        'documentType': '59',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "60": {
        'name': 'Service Level Agreement',
        'documentType': '60',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "61": {
        'name': 'Risk Sharing Mechanisms',
        'documentType': '61',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "62": {
        'name': 'Public Private Partnerships (PPPs) Agreement',
        'documentType': '62',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "63": {
        'name': 'Hedging Instruments',
        'documentType': '63',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "64": {
        'name': 'Completion Guarantees',
        'documentType': '64',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "65": {
        'name': 'Governance and Oversight',
        'documentType': '65',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "66": {
        'name': 'Independent Monitoring Engagement Agreements',
        'documentType': '66',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "67": {
        'name': 'Environmental, Social, and Governance (ESG) Safeguards',
        'documentType': '67',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "68": {
        'name': 'Sustainability Certifications',
        'documentType': '68',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "69": {
        'name': 'Community Engagement Plans',
        'documentType': '69',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "70": {
        'name': 'Security Measures Documentation (optional)',
        'documentType': '70',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "71": {
        'name':
            'Proof of Unpaid Taxes, Utility Bills, Fees, or Other Property-related Expenses.',
        'documentType': '71',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "72": {
        'name':
            'Proof of Hidden Liabilities or Obligations that have not been Disclosed.',
        'documentType': '72',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "73": {
        'name':
            'Proof of Unresolved Structural, Maintenance, or Safety Issues.',
        'documentType': '73',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
      "74": {
        'name': 'Contractors',
        'documentType': '74',
        'options': <String>[],
        'files': <String, dynamic>{},
        'selectedFileOption': '',
        'required': 1,
      },
    };

    documents = appState.viewData!['AssetTokenizationDocuments'] ?? [];
    for (var item in documents) {
      var val = documentTypeAndCodes[item['documentType']];
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
        documentOptions: getDocumentOptions(val['options']),
      );
    }
  }

  Future<void> refreshCurrentTokenizationInfo() async {
    try {
      var uri = '/v1/tokenization/detail/${appState.viewData!['id']}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );
      if (responseData['statusCode'] == 200) {
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
      Map responseData = await makePutRequestForMultipartDocumentUpload(
        uri: '/v1/tokenization/document',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: appState.primaryWallet.signer!,
        file: file,
        tokenizedAssetId: appState.viewData!['id'],
        documentTitle: documentTitle.toLowerCase().replaceAll(' ', '-'),
        documentType: documentType,
      );
      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
        hideLoader(context);
      } else {
        popup(
          context,
          title: "error".tr(),
          message:
              responseData['data']['message'] ?? responseData['data']['error'],
        );
        hideLoader(context);
      }
    } catch (e) {
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
      Map requestBody = {};
      Map responseData = await makeDeleteRequest(
        uri: '/v1/tokenization/document/$documentId',
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: appState.primaryWallet.signer!,
        body: jsonEncode(requestBody),
      );

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
      hideLoader(context);
    } catch (e) {
      hideLoader(context);
      popup(
        context,
        title: "error".tr(),
        message: "Sorry, something went wrong. Please try again.",
      );
    }
  }
}

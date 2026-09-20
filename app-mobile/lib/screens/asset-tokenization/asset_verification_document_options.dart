import 'dart:convert';
import 'dart:developer';

import 'package:dotted_border/dotted_border.dart';
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
import 'package:uuid/uuid.dart';

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
    inspect(appState.tokenizationData);
    inspect(appState.viewData);
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
                    width: 330.sp,
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
            SizedBox(height: 20),
            for (var item in appState.viewData!['files']) ...[
              if (appState.viewData!["fundingStructure"] == 0 &&
                  item["requiredForDebtAndHybrid"] == true) ...[
                SizedBox.shrink(),
              ] else ...[
                Builder(
                  builder: (context) {
                    var documents = <Map>[];
                    for (var file
                        in appState.viewData!['AssetTokenizationDocuments']) {
                      if (file["documentType"] == item["documentType"]) {
                        documents.add(file);
                      }
                    }
                    return proofItemCard(item: item, documents: documents);
                  },
                ),
                SizedBox(height: 10.sp),
              ],
            ],
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

  void uploadFileBottomSheet({
    required String title,
    required String documentType,
  }) {
    String errorMsg = '';

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      constraints: BoxConstraints(maxHeight: 400.h),
      backgroundColor: notifier.getwihitecolor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
          top: Radius.circular(16),
        ), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return StatefulBuilder(
          builder: (BuildContext context, StateSetter setModalState) {
            return DraggableScrollableSheet(
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
                          crossAxisAlignment: CrossAxisAlignment.center,
                          children: [
                            SizedBox(
                              width: 280.sp,
                              child: Text(
                                appState.viewData?['title'],
                                overflow: TextOverflow.ellipsis,
                                style: TextStyle(
                                  fontSize: 18,
                                  fontFamily: fontsemibold,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ),
                            TextButton(
                              onPressed: () => Navigator.of(context).pop(),
                              child: Icon(Icons.cancel_outlined, size: 20),
                              style: ButtonStyle(
                                padding: WidgetStatePropertyAll(
                                  EdgeInsets.all(7),
                                ),
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
                                "Upload ${title}",
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
                      if (errorMsg.isNotEmpty) ...[
                        SizedBox(height: 5.h),
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 10.0),
                          child: Row(
                            children: [
                              SizedBox(
                                width: 350,
                                child: Text(
                                  errorMsg,
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontFamily: fontbody,
                                    color: Colors.red,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                      GestureDetector(
                        onTap: () async {
                          var file = await getFile();
                          if (file != null && file.size > 900000) {
                            errorMsg = "filesizeerror".tr();
                            file = null;
                            setModalState(() {});
                            return;
                          }

                          var uuid = const Uuid();
                          String shortId = uuid.v4().split('-').first;

                          uploadFile(
                            file!,
                            documentType,
                            "${title.split('(').first.trim()}-$shortId.${file.path?.split('.').last}"
                                .toLowerCase(),
                          );
                        },
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Column(
                              crossAxisAlignment: CrossAxisAlignment.center,
                              children: [
                                SizedBox(height: height / 50),
                                Padding(
                                  padding: const EdgeInsets.fromLTRB(
                                    20,
                                    10,
                                    20,
                                    3,
                                  ),
                                  child: DottedBorder(
                                    ignoring: false,
                                    options: RectDottedBorderOptions(
                                      dashPattern: [10, 4],
                                      color: notifier.getbluecolor80,
                                    ),
                                    child: Container(
                                      width: 320,
                                      padding: EdgeInsets.all(20),
                                      child: Row(
                                        mainAxisAlignment:
                                            MainAxisAlignment.center,
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
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),
                      SizedBox(height: 20.sp),
                    ],
                  ),
                );
              },
            );
          },
        );
      },
    );
  }

  Widget proofItemCard({required dynamic item, required List<Map> documents}) {
    String title = item['name'];
    String documentType = item['documentType'];
    bool isRequired = item['required'];
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: 10),
      child: Container(
        width: double.infinity,
        padding: EdgeInsets.symmetric(horizontal: 10, vertical: 20),
        decoration: BoxDecoration(
          border: Border.all(color: notifier.getsplashgrey, width: 1),
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: wihitecolor,
        ),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  spacing: 5,
                  children: [
                    SizedBox(
                      width: 300,
                      child: Text.rich(
                        TextSpan(
                          text: title,
                          style: TextStyle(
                            fontSize: 14,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                          children: [
                            if (isRequired) ...[
                              TextSpan(
                                text: ' *',
                                style: TextStyle(
                                  fontWeight: FontWeight.bold,
                                  color: Colors.red,
                                  fontSize: 16,
                                ),
                              ),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
              ],
            ),
            for (var file in documents) ...[
              SizedBox(height: 10),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  TextButton(
                    style: ButtonStyle(
                      padding: WidgetStatePropertyAll(EdgeInsets.all(0)),
                      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      minimumSize: WidgetStatePropertyAll(Size.zero),
                    ),
                    onPressed: () {
                      var fileUrl = file['documentUrl'].toString();
                      if (fileUrl.isNotEmpty && fileUrl.endsWith('.pdf')) {
                        appState.pdfUrl = fileUrl;
                        appState.currentAction = PageAction(
                          state: PageState.addPage,
                          page: PdfViewPageConfig,
                        );

                        return;
                      }

                      appState.initialUrl = fileUrl;
                      appState.setPage(page: AppImageViewerPageConfig);
                    },
                    child: Row(
                      children: [
                        Row(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          spacing: 3,
                          children: [
                            Icon(
                              Icons.file_present_outlined,
                              color: notifier.getbluewhitecolor,
                            ),
                            Container(
                              constraints: BoxConstraints(maxWidth: 250.sp),
                              child: Text(
                                file['documentTitle'],
                                style: TextStyle(
                                  decoration: TextDecoration.underline,
                                  fontSize: 12.sp,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                  TextButton(
                    onPressed: () {
                      confirmVerificationDocumentDeletePopup(
                        context,
                        fileName: file['documentTitle'],
                        onConfirmationSuccess: () {
                          deleteFile(file['id'].toString());
                        },
                      );
                    },
                    style: ButtonStyle(
                      padding: WidgetStatePropertyAll(EdgeInsets.all(0)),
                      tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      minimumSize: WidgetStatePropertyAll(Size.zero),
                    ),
                    child: Icon(
                      CupertinoIcons.trash,
                      size: 15,
                      color: Colors.red,
                    ),
                  ),
                ],
              ),
            ],
            SizedBox(height: 10.sp),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: TextButton(
                onPressed: () {
                  uploadFileBottomSheet(
                    title: title,
                    documentType: documentType,
                  );
                },
                style: ButtonStyle(
                  overlayColor: WidgetStateProperty.all<Color>(
                    notifier.getsplashgrey,
                  ),
                  fixedSize: WidgetStateProperty.all(Size(150.sp, 20.sp)),
                  side: WidgetStateProperty.all(
                    BorderSide(
                      color: notifier.getbluecolor90,
                      width: 1,
                      style: BorderStyle.solid,
                    ),
                  ),
                  shape: WidgetStateProperty.all<RoundedRectangleBorder>(
                    const RoundedRectangleBorder(
                      borderRadius: BorderRadius.all(Radius.circular(10)),
                    ),
                  ),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.add_circle_outlined,
                      color: notifier.getbluewhitecolor,
                      size: 18,
                    ),
                    SizedBox(width: 3),
                    Text(
                      "Add File",
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> refreshCurrentTokenizationInfo() async {
    try {
      var uri = '/v1/tokenization/detail/${appState.viewData!['id']}';

      Map responseData = await makeGetRequest(
        uri: Uri.encodeFull(uri),
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        address: appState.primaryWallet.signer!,
      );
      inspect(responseData['data']);
      if (responseData['statusCode'] == 200) {
        appState.viewData = {
          ...responseData['data'] as Map,
          'title': appState.viewData!['title'],
          'files': appState.viewData!['files'],
        };
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
        address: appState.primaryWallet.signer!,
        file: file,
        tokenizedAssetId: appState.viewData!['id'],
        documentTitle: documentTitle,
        documentType: documentType,
      );
      inspect(responseData['data']);
      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo();
        hideLoader(context);
        Navigator.of(context).pop();
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
        address: appState.primaryWallet.signer!,
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

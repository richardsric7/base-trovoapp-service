import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
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
import 'package:trovo_app/router/ui_pages.dart';
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
  Map documents = {};
  late Future<dynamic> formJsonFuture;

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
    super.initState();
    getdarkmodepreviousstate();
    formJsonFuture = fetchFormJson(appState.viewData?['assetType']);
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

            FutureBuilder<dynamic>(
              future: formJsonFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return SizedBox(
                    height: height / 2,
                    child: Center(
                      child: CircularProgressIndicator(
                        backgroundColor: notifier.getbluecolor,
                        valueColor: new AlwaysStoppedAnimation<Color>(
                          notifier.getgreencolor,
                        ),
                        strokeWidth: 3.0,
                      ),
                    ),
                  );
                } else if (snapshot.connectionState == ConnectionState.done) {
                  if (snapshot.hasError) {
                    return Padding(
                      padding: const EdgeInsets.all(8.0),
                      child: SizedBox(
                        height: height / 6,
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              "somethingwentwrong".tr(),
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontSize: 16,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontbody,
                              ),
                            ),
                            ElevatedButton(
                              onPressed: () {
                                setState(() {
                                  formJsonFuture = fetchFormJson(
                                    appState.viewData?['assetType'],
                                  );
                                });
                              },
                              style: ButtonStyle(
                                backgroundColor: WidgetStateProperty.all<Color>(
                                  notifier.getbluecolor!,
                                ),
                                foregroundColor: WidgetStateProperty.all<Color>(
                                  notifier.getwihitecolor,
                                ),
                              ),
                              child: Text(
                                "retry".tr(),
                                style: TextStyle(fontFamily: fontsemibold),
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  } else if (snapshot.hasData) {
                    var documents = snapshot.data! as Map<dynamic, dynamic>;
                    return Column(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 10.0),
                          child: Row(
                            children: [
                              Text(
                                "Upload the files required in each folder",
                                style: TextStyle(
                                  fontSize: 14,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ],
                          ),
                        ),
                        SizedBox(height: 20),
                        if (documents.isEmpty) ...[
                          Padding(
                            padding: const EdgeInsets.all(8.0),
                            child: SizedBox(
                              height: height / 6,
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    "No document to upload here. Please contact administrator",
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                      fontSize: 16,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ] else ...[
                          for (var item in documents.entries) ...[
                            Builder(
                              builder: (context) {
                                int uploadedFileTypeCount = 0;
                                int uploadedRequiredFileTypeCount = 0;
                                int requiredFilesCount = 0;
                                // get the number of documents that has been
                                // uploaded so far in order to calculate
                                // the number that is remaining
                                for (var file in item.value['files']) {
                                  if (file['required'] == true) {
                                    requiredFilesCount += 1;
                                  }

                                  for (var uploadedFile
                                      in appState
                                          .viewData!['AssetTokenizationDocuments']) {
                                    if (file['documentType'] ==
                                        uploadedFile['documentType']) {
                                      uploadedFileTypeCount += 1;
                                      if (file['required'] == true) {
                                        uploadedRequiredFileTypeCount += 1;
                                      }
                                      break;
                                    }
                                  }
                                }
                                return documentCardItem(
                                  title: item.value['sectionName'],
                                  onTap: () {
                                    appState.viewData!['title'] =
                                        item.value['sectionName'];
                                    appState.viewData!['files'] =
                                        item.value['files'];
                                    appState.setPage(
                                      page:
                                          AssetVerificationDocumentOptionsViewPageConfig,
                                    );
                                  },
                                  documentUploadCount: uploadedFileTypeCount,
                                  uploadedRequiredFileTypeCount:
                                      uploadedRequiredFileTypeCount,
                                  requiredFilesCount: requiredFilesCount,
                                  totalDocuments: item.value['files'].length,
                                );
                              },
                            ),
                            SizedBox(height: 10),
                          ],
                        ],
                      ],
                    );
                  }
                }
                return Text(
                  '',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.bold,
                    fontFamily: fontsemibold,
                  ),
                );
              },
            ),
            SizedBox(height: height / 30),
            Align(
              alignment: Alignment.bottomCenter,
              child: Button(
                'Continue',
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  Navigator.of(context).pop();
                },
              ),
            ),
            SizedBox(height: height / 10),
          ],
        ),
      ),
    );
  }

  Widget documentCardItem({
    required String title,
    required void Function() onTap,
    required int documentUploadCount,
    required int uploadedRequiredFileTypeCount,
    required int requiredFilesCount,
    required int totalDocuments,
  }) {
    double documentUploadProgress = 0;

    documentUploadProgress = documentUploadCount == totalDocuments
        ? 1
        : documentUploadCount / totalDocuments; // Convert to 0-1 range
    return GestureDetector(
      onTap: onTap,
      child: Padding(
        padding: EdgeInsets.symmetric(horizontal: 10),
        child: Container(
          width: double.infinity,
          padding: EdgeInsets.symmetric(horizontal: 5, vertical: 10),
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(15.0)),
            color: notifier.isDark
                ? darktilewhitecolor
                : notifier.getaddsubwalletgrey,
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Icon(Icons.folder, color: notifier.getbluewhitecolor, size: 35),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  SizedBox(
                    width: 280.sp,
                    child: Text(
                      title,
                      style: TextStyle(
                        fontSize: 14.sp,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(height: 10.sp),
                  Container(
                    height: 5.sp,
                    width: 280.sp,
                    child: LinearProgressIndicator(
                      value: documentUploadProgress, // Show progress (0 to 1)
                      minHeight: 10,
                      borderRadius: BorderRadius.circular(10),
                      backgroundColor: Colors.grey[300],
                      color: uploadedRequiredFileTypeCount == requiredFilesCount
                          ? Colors.green
                          : Colors.blue[400],
                    ),
                  ),
                  SizedBox(height: 5),
                  Text(
                    "$documentUploadCount/$totalDocuments files uploaded ($uploadedRequiredFileTypeCount/$requiredFilesCount required)",
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getgrey,
                    ),
                  ),
                  SizedBox(height: 5),
                ],
              ),
              Icon(
                Icons.arrow_forward_ios_rounded,
                color: notifier.getbluewhitecolor,
                size: 18,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<dynamic> fetchFormJson(String formId) async {
    var uri = '/v1/forms/$formId';

    Map responseData = await makeGetRequest(
      uri: Uri.encodeFull(uri),
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.signer!,
    );
    // inspect(responseData['data']);
    if (responseData['statusCode'] == 200) {
      try {
        var formStr = responseData['data']['formString'];
        var first = jsonDecode(formStr);
        var jsonObj = first is String ? jsonDecode(first) : first;
        inspect(jsonObj);
        return jsonObj['documents'];
      } catch (e) {
        print(e);
        // inspect(e);
      }
    }

    return Future.error('Error fetching form data');
  }
}

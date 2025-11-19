import 'dart:convert';
import 'dart:developer';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class TokenizeAsset extends StatefulWidget {
  const TokenizeAsset({Key? key}) : super(key: key);

  @override
  State<TokenizeAsset> createState() => _TokenizeAssetState();
}

class _TokenizeAssetState extends State<TokenizeAsset>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  var formData = {};
  late Future<dynamic> formJsonFuture;
  late RefreshController _refreshController;

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
    notifier = Provider.of<ColorNotifier>(context, listen: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    _refreshController = RefreshController(initialRefresh: false);
    getdarkmodepreviousstate();
    formJsonFuture = fetchFormJson(appState.viewData?['assetType']);
  }

  void refreshData() async {
    try {
      await refreshCurrentTokenizationInfo(appState);
      formJsonFuture = fetchFormJson(appState.viewData?['assetType']);
      _refreshController.refreshCompleted();
      appState.updateListeners();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  @override
  Widget build(BuildContext context) {
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: SmartRefresher(
        enablePullDown: true,
        controller: _refreshController,
        onRefresh: refreshData,
        child: SingleChildScrollView(
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                "assetdetails".tr(),
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Text(
                  "fillouttocreatetoken".tr(),
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontbody,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
              SizedBox(height: height / 50),
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
                                  backgroundColor:
                                      WidgetStateProperty.all<Color>(
                                        notifier.getbluecolor!,
                                      ),
                                  foregroundColor:
                                      WidgetStateProperty.all<Color>(
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
                      appState.viewData!['documents'] = documents;

                      int uploadedRequiredFileTypeCount = 0;
                      int requiredFilesCount = 0;

                      for (var item in documents.entries) {
                        for (var file in item.value['files']) {
                          if (file['required'] == true) {
                            requiredFilesCount += 1;
                          }

                          for (var uploadedFile
                              in appState
                                  .viewData!['AssetTokenizationDocuments']) {
                            if (file['documentType'] ==
                                uploadedFile['documentType']) {
                              if (file['required'] == true) {
                                uploadedRequiredFileTypeCount += 1;
                              }
                              break;
                            }
                          }
                        }
                      }

                      return Column(
                        children: [
                          detailItem(
                            "assetinformation".tr(),
                            "providebasicinfo".tr(),
                            "1",
                            appState.viewData!['assetName'] != null &&
                                    appState.viewData!['assetName'].length > 0
                                ? "continuee".tr()
                                : "start".tr(),
                            onTap: () {
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: getFormPage(),
                              );
                            },
                          ),
                          SizedBox(height: height / 50),
                          detailItem(
                            "assetverificationdocs".tr(),
                            "provideverificationdocs".tr(),
                            "2",
                            appState.viewData!['AssetTokenizationDocuments'] !=
                                        null &&
                                    appState
                                            .viewData!['AssetTokenizationDocuments']
                                            .length >
                                        0
                                ? "continuee".tr()
                                : "start".tr(),
                            onTap: () {
                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: AssetVerificationDocumentsViewPageConfig,
                              );
                            },
                          ),
                          SizedBox(height: height / 50),
                          detailItem(
                            "assettokeninfo".tr(),
                            "providetokeninfo".tr(),
                            "3",
                            appState.viewData!['assetCode'] != null &&
                                    appState.viewData!['assetCode'].length > 0
                                ? "continuee".tr()
                                : "start".tr(),
                            isDisabled:
                                appState.viewData!['assetName'] == null ||
                                appState.viewData!['assetName'].length == 0,
                            onTap: () {
                              if (appState.viewData!['assetName'] != null &&
                                  appState.viewData!['assetName'].length > 0) {
                                appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: AssetTokenInformationViewPageConfig,
                                );
                              }
                            },
                          ),
                          SizedBox(height: height / 30),
                          Button(
                            "completetokenization".tr(),
                            notifier.getbluecolor,
                            wihitecolor,
                            onTap: () async {
                              var assetCode =
                                  appState.viewData!['assetCode'] ?? '';
                              var assetName =
                                  appState.viewData!['assetName'] ?? '';

                              if (assetCode.length == 0 ||
                                  assetName.length == 0) {
                                popup(
                                  context,
                                  title: "formincomplete".tr(),
                                  message: "pleasefillouttokenizationform".tr(
                                    args: ['3'],
                                  ),
                                );
                                return;
                              }

                              if (uploadedRequiredFileTypeCount <
                                  requiredFilesCount) {
                                popup(
                                  context,
                                  title: "formincomplete".tr(),
                                  message:
                                      "Please upload all required documents to proceed.",
                                );
                                return;
                              }

                              appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: ConfirmTokenizationDetailsViewPageConfig,
                              );
                            },
                          ),
                          SizedBox(height: 10),
                          Button(
                            "deletetokenization".tr(),
                            Colors.red,
                            wihitecolor,
                            onTap: () {
                              confirmTokenizationDeletePopup(
                                context,
                                onConfirmationSuccess: deleteTokenization,
                              );
                            },
                          ),
                          SizedBox(height: height / 10),
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
            ],
          ),
        ),
      ),
    );
  }

  PageConfiguration getFormPage() {
    var data = appState.tokenizationData;
    for (var i = 0; i < data['assetTypes'].length; i++) {
      if (data['assetTypes'][i]['id'].toString() ==
          appState.viewData!['assetType'].toString()) {
        appState.viewData!['assetFormName'] =
            data['assetTypes'][i]['assetType'];

        break;
      }
    }

    switch (appState.viewData!['assetType'].toString()) {
      case '1123':
      case '1124':
      case '1125':
      case '1114':
      case '1115':
      case '1182':
      case '1180':
      case '1121':
      case '1122':
      case '1174':
        return EquityMutualFundsAssetInformationViewPageConfig;
      default:
        if (appState.viewData!['assetAlreadyExists'] == 1) {
          return AssetInformationViewPageConfig;
        }

        return UpcomingAssetInformationViewPageConfig;
    }
  }

  Widget detailItem(
    String title,
    String description,
    String number,
    String status, {
    required void Function() onTap,
    bool isDisabled = false,
  }) {
    return TextButton(
      onPressed: onTap,
      style: TextButton.styleFrom(
        padding: EdgeInsets.zero,
        minimumSize: Size(50, 30),
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
        alignment: Alignment.centerLeft,
      ),
      child: Stack(
        alignment: AlignmentDirectional.centerStart,
        children: [
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
            child: Container(
              decoration: BoxDecoration(
                border: Border.all(
                  color: isDisabled
                      ? notifier.getsplashgrey
                      : notifier.getbluewhitecolor,
                  width: 1.5,
                ),
                borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                color: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
              ),
              child: Padding(
                padding: const EdgeInsets.fromLTRB(25.0, 15.0, 5.0, 15.0),
                child: Row(
                  children: [
                    Container(
                      width: width / 1.27,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          SizedBox(width: width / 50),
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Container(
                                width: width / 1.87,
                                child: Text(
                                  title,
                                  style: TextStyle(
                                    fontSize: 15,
                                    fontWeight: FontWeight.w400,
                                    color: isDisabled
                                        ? notifier.getsplashgrey
                                        : notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold,
                                  ),
                                ),
                              ),
                              Text(
                                '${status} >>>',
                                style: TextStyle(
                                  fontStyle: FontStyle.italic,
                                  fontSize: 12,
                                  fontWeight: FontWeight.w400,
                                  color: isDisabled
                                      ? notifier.getsplashgrey
                                      : notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ],
                          ),
                          SizedBox(height: 20),
                          Text(
                            description,
                            overflow: TextOverflow.visible,
                            style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w400,
                              color: isDisabled
                                  ? notifier.getsplashgrey
                                  : notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                          SizedBox(height: 20),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 5.0),
            child: Container(
              decoration: BoxDecoration(
                border: Border.all(
                  color: isDisabled
                      ? notifier.getsplashgrey
                      : notifier.getbluewhitecolor,
                  width: 1.5,
                ),
                shape: BoxShape.circle,
                color: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
              ),
              child: Padding(
                padding: const EdgeInsets.all(10.0),
                child: Text(
                  number,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w400,
                    color: isDisabled
                        ? notifier.getsplashgrey
                        : notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<dynamic> fetchFormJson(String formId) async {
    var uri = '/v1/forms';

    switch (appState.viewData!['assetType'].toString()) {
      case '1123':
      case '1124':
      case '1125':
      case '1114':
      case '1115':
      case '1182':
      case '1180':
      case '1121':
      case '1122':
      case '1174':
        uri = '$uri/$formId';
      default:
        if (appState.viewData!['fundingStructure'] == 0) {
          uri = '$uri/1'; // equity verification documents
        } else {
          uri = '$uri/2'; // debt verification documents
        }
    }

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

  void deleteTokenization() async {
    try {
      showLoader(context);

      Map responseData = await makeDeleteRequest(
        uri: '/v1/tokenization/${appState.viewData!['id']}',
        body: "",
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      hideLoader(context);
      if (responseData['statusCode'] == 200) {
        appState.currentAction = PageAction(
          state: PageState.replaceAll,
          page: BottomHomePageConfig,
        );
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }
}

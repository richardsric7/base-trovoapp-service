import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
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
  late TabController tabController;

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
    tabController = TabController(length: 2, vsync: this);
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
            detailItem(
              "assetinformation".tr(),
              "providebasicinfo".tr(),
              "1",
              appState.viewData!['assetDescription'] != null &&
                      appState.viewData!['assetDescription'].length > 0
                  ? "continuee".tr()
                  : "start".tr(),
              onTap: () {
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: appState.viewData!['assetAlreadyExists'] == 1
                      ? AssetInformationViewPageConfig
                      : UpcomingAssetInformationViewPageConfig,
                );
              },
            ),
            SizedBox(height: height / 50),
            detailItem(
              "assetverificationdocs".tr(),
              "provideverificationdocs".tr(),
              "2",
              appState.viewData!['AssetTokenizationDocuments'] != null &&
                      appState.viewData!['AssetTokenizationDocuments'].length >
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
                  appState.viewData!['assetDescription'] == null ||
                  appState.viewData!['assetDescription'].length == 0,
              onTap: () {
                if (appState.viewData!['assetDescription'] != null &&
                    appState.viewData!['assetDescription'].length > 0) {
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
                var assetCode = appState.viewData!['assetCode'] ?? '';
                var documents =
                    appState.viewData!['AssetTokenizationDocuments'] ?? [];
                var assetDescription =
                    appState.viewData!['assetDescription'] ?? '';

                if (assetCode.length == 0 ||
                    documents.length == 0 ||
                    assetDescription.length == 0) {
                  popup(
                    context,
                    title: "formincomplete".tr(),
                    message: "pleasefillouttokenizationform".tr(args: ['3']),
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
        ),
      ),
    );
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

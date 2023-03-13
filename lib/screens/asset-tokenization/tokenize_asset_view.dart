import 'package:flutter/material.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
              'Create Asset Token',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(height: height / 50),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: Text(
                LanguageEn.fillouttocreatetoken,
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Asset Category',
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
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Select the asset category from the list below ',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'hint',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Country',
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
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Select Country where the asset is located',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 70,
            ),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencies,
                null,
                'hint',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Asset Details',
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
              height: height / 50,
            ),
            detailItem(
              'Asset Information',
              'Provide the basic information about the asset you want to tokenize',
              onTap: () {
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: AssetInformationViewPageConfig,
                );
              },
            ),
            SizedBox(
              height: height / 50,
            ),
            detailItem(
              'Asset Verification Documents',
              'Provide the basic information about the asset you want to tokenize',
              onTap: () {
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: AssetInformationViewPageConfig,
                );
              },
            ),
            SizedBox(
              height: height / 50,
            ),
            detailItem(
              'Asset Token Information',
              'Provide information about the asset token you want to create',
              onTap: () {
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: AssetTokenInformationViewPageConfig,
                );
              },
            ),
            SizedBox(
              height: height / 30,
            ),
            Button(
              'Complete Tokenization',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {},
            ),
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
  }

  Widget detailItem(String title, String description,
      {required void Function() onTap}) {
    return Column(
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
          child: Container(
            decoration: BoxDecoration(
              border: Border.all(color: notifier.getbluewhitecolor, width: 1.5),
              borderRadius: const BorderRadius.all(Radius.circular(15.0)),
              color: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 10.0, vertical: 15.0),
              child: Row(
                children: [
                  Container(
                    width: width / 1.2,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        SizedBox(
                          width: width / 50,
                        ),
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Text(
                              title,
                              style: TextStyle(
                                fontSize: 15,
                                fontWeight: FontWeight.w400,
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                              ),
                            ),
                            TextButton(
                              onPressed: onTap,
                              style: TextButton.styleFrom(
                                  padding: EdgeInsets.zero,
                                  minimumSize: Size(50, 30),
                                  tapTargetSize:
                                      MaterialTapTargetSize.shrinkWrap,
                                  alignment: Alignment.centerLeft),
                              child: Text(
                                'Start',
                                style: TextStyle(
                                  fontStyle: FontStyle.italic,
                                  fontSize: 15,
                                  fontWeight: FontWeight.w400,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                          ],
                        ),
                        SizedBox(
                          height: 20,
                        ),
                        Text(
                          description,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.w400,
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                          ),
                        ),
                        SizedBox(
                          height: 20,
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}

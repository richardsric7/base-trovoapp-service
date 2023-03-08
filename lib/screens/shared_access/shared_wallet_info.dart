import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';

class SharedWalletInfo extends StatefulWidget {
  const SharedWalletInfo({Key? key}) : super(key: key);

  @override
  State<SharedWalletInfo> createState() => _SharedWalletInfoState();
}

class _SharedWalletInfoState extends State<SharedWalletInfo> {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  late Future<Map> responseData;

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
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet =
        appState.userInfo!.getWallet(appState.viewData!['walletPublicKey']);

    responseData = fetchWalletBalance(
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: wallet.publicKey!);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      appBar: CustomAppBar(
        context,
        notifier.getwihitecolor,
        'Shared Access',
        notifier.getbluewhitecolor,
        height: height / 15,
      ),
      body: Column(
        children: [
          SizedBox(
            height: height / 20,
          ),
          Text(
            wallet.alias!,
            style: TextStyle(
              fontSize: 20,
              fontFamily: fontsemibold,
              color: notifier.getbluewhitecolor,
            ),
          ),
          SizedBox(
            height: height / 50,
          ),
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
            child: Container(
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                color: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
              ),
              child: Padding(
                padding: const EdgeInsets.symmetric(
                    horizontal: 20.0, vertical: 15.0),
                child: Column(
                  children: [
                    Text(
                      'Description',
                      style: TextStyle(
                        fontSize: 17,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 90,
                    ),
                    Text(
                      wallet.description!,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 16,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      'Owner',
                      style: TextStyle(
                        fontSize: 17,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 90,
                    ),
                    Text(
                      wallet.owner!,
                      style: TextStyle(
                        fontSize: 16,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      'Permissions',
                      style: TextStyle(
                        fontSize: 17,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 90,
                    ),
                    Container(
                      width: width / 1.3,
                      child: Wrap(alignment: WrapAlignment.center, children: [
                        Text(
                          'You have ',
                          style: TextStyle(
                            fontSize: 16,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Text(
                          wallet.accesses![0].toString().toLowerCase(),
                          style: TextStyle(
                            fontSize: 16,
                            fontFamily: fontsemibold,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        if (wallet.accesses!.length > 1) ...[
                          SizedBox(
                            width: width / 90,
                          ),
                          Text(
                            'and',
                            style: TextStyle(
                              fontSize: 16,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            width: width / 90,
                          ),
                          Text(
                            wallet.accesses![1].toString().toLowerCase(),
                            style: TextStyle(
                              fontSize: 16,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                          SizedBox(
                            width: width / 90,
                          ),
                        ],
                        Text(
                          'access on this wallet',
                          style: TextStyle(
                            fontSize: 16,
                            fontFamily: fontbody,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                      ]),
                    ),
                    SizedBox(
                      height: height / 90,
                    ),
                    SizedBox(height: 2),
                  ],
                ),
              ),
            ),
          ),
          SizedBox(
            height: height / 20,
          ),
          FutureBuilder<Map>(
            future: responseData,
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return Center(
                  child: CircularProgressIndicator(
                    backgroundColor: notifier.getbluecolor,
                    valueColor: new AlwaysStoppedAnimation<Color>(
                      notifier.getgreencolor,
                    ),
                    strokeWidth: 3.0,
                  ),
                );
              } else if (snapshot.connectionState == ConnectionState.done) {
                if (snapshot.hasError) {
                  return Padding(
                    padding: const EdgeInsets.all(8.0),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          LanguageEn.somethingwentwrong,
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              fontSize: 16,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody),
                        ),
                        ElevatedButton(
                          onPressed: () {
                            setState(() {
                              responseData = fetchWalletBalance(
                                  signer: appState.activeWallet!.signer!,
                                  secretKey: appState.secretKeys[0],
                                  publicKey: wallet.publicKey!);
                            });
                          },
                          style: ButtonStyle(
                            backgroundColor: MaterialStateProperty.all<Color>(
                                notifier.getbluecolor!),
                          ),
                          child: Text(
                            LanguageEn.retry,
                            style: TextStyle(
                              fontFamily: fontsemibold,
                            ),
                          ),
                        ),
                      ],
                    ),
                  );
                } else if (snapshot.hasData) {
                  return Column(
                    children: [
                      Button(
                        'View wallet',
                        notifier.getbluecolor,
                        wihitecolor,
                        onTap: () {
                          appState.viewData = {
                            'rel': 'sharedWalletView',
                            'walletPublicKey': wallet.publicKey,
                          };

                          appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: WalletDetailsViewPageConfig);
                        },
                      ),
                      SizedBox(
                        height: height / 50,
                      ),
                      ButtonOutlined(
                        'View transaction history',
                        notifier.getbluecolor80,
                        wihitecolor,
                        onTap: () {
                          appState.viewData = {
                            'rel': 'sharedWalletView',
                            'walletPublicKey': wallet.publicKey,
                          };

                          appState.currentAction = PageAction(
                              state: PageState.addPage,
                              page: PaymentHistoryViewPageConfig);

                          appState.setFilterQuery = "";

                          appState.getHistory(context, wallet.publicKey!);
                        },
                      ),
                      if (wallet.isInitiator) ...[
                        SizedBox(height: height / 50),
                        ButtonOutlined(
                          'Modify shared access',
                          notifier.getwihitecolor,
                          notifier.getbluewhitecolor,
                          onTap: () {
                            appState.viewData = {
                              'walletPublicKey': wallet.publicKey,
                            };
                            appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: UpdateSharedAccessViewPageConfig);
                          },
                        ),
                      ]
                    ],
                  );
                } else {
                  return const Text('Empty data');
                }
              } else {
                return Text('State: ${snapshot.connectionState}');
              }
            },
          )
        ],
      ),
    );
  }

  Future<Map> fetchWalletBalance(
      {required String signer,
      required String secretKey,
      required String publicKey}) async {
    try {
      Map responseData = await makeGetRequest(
        uri: '/v1/shared-access/wallet-balances',
        signer: signer,
        secretKey: secretKey, // the primary wallet secret key
        publicKey: publicKey,
      );

      print('response: ${responseData}');

      if (responseData['statusCode'] == 200) {
        return responseData['data'];
      } else {
        return Future.error('Error! Something went wrong.');
      }
    } catch (e) {
      return Future.error('Error! ${e}');
    }
  }

  @override
  void dispose() {
    appState.viewData!['rel'] = '';
    appState.returnView = null;
    super.dispose();
  }
}

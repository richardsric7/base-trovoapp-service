import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Permission.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
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
  var viewData;
  bool isInitiator = false;
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
    viewData = appState.viewData![SharedWalletInfoViewPageConfig.key];
    print(viewData);

    for (var i = 0; i < viewData['permissions'].length; i++) {
      if (viewData['permissions'][i] == 'INITIATOR') isInitiator = true;
    }

    responseData = fetchWalletBalance(
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0],
        publicKey: viewData['walletPublicKey']);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
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
              viewData['walletAlias'],
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
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Padding(
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
                            viewData['walletDescription'],
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
                            viewData['owner'],
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
                            child: Wrap(
                                alignment: WrapAlignment.center,
                                children: [
                                  Text(
                                    'You have ',
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontFamily: fontbody,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  Text(
                                    viewData['permissions'][0]
                                        .toString()
                                        .toLowerCase(),
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontFamily: fontsemibold,
                                      color: notifier.getbluewhitecolor,
                                    ),
                                  ),
                                  if (viewData['permissions'].length > 1) ...[
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
                                      viewData['permissions'][1]
                                          .toString()
                                          .toLowerCase(),
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
                  ],
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
                                    publicKey: viewData['walletPublicKey']);
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
                            appState.viewData![SharedWalletDetailsViewPageConfig
                                .key] = viewData;

                            appState.viewData![SharedWalletDetailsViewPageConfig
                                    .key]['claimed'] =
                                snapshot.data!['assetBalances']['claimed'];

                            appState.viewData![SharedWalletDetailsViewPageConfig
                                    .key]['unclaimed'] =
                                snapshot.data!['assetBalances']['unclaimed'];

                            print('================${appState.viewData}');
                            appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: SharedWalletDetailsViewPageConfig);
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
                            appState.viewData![
                                PaymentHistoryViewPageConfig.key] = viewData;

                            appState.viewData![PaymentHistoryViewPageConfig.key]
                                    ['claimed'] =
                                snapshot.data!['assetBalances']['claimed'];

                            appState.viewData![PaymentHistoryViewPageConfig.key]
                                    ['unclaimed'] =
                                snapshot.data!['assetBalances']['unclaimed'];

                            appState.currentAction = PageAction(
                                state: PageState.addPage,
                                page: PaymentHistoryViewPageConfig);
                            appState.setFilterQuery = "";
                            appState.getHistory(
                                context, viewData['walletPublicKey']);
                          },
                        ),
                        if (isInitiator) ...[
                          SizedBox(height: height / 50),
                          ButtonOutlined(
                            'Modify shared access',
                            notifier.getwihitecolor,
                            notifier.getbluewhitecolor,
                            onTap: () {
                              var viewers = <Permission>[];
                              var approvers = <Permission>[];
                              var initiators = <Permission>[];

                              for (var i = 0;
                                  i <
                                      viewData['walletSettings']['permissions']
                                          .length;
                                  i++) {
                                if (viewData['walletSettings']['permissions'][i]
                                        ['permission'] ==
                                    'VIEW-ONLY') {
                                  viewers.add(
                                    Permission(
                                      targetUsername: viewData['walletSettings']
                                          ['permissions'][i]['targetUsername'],
                                      fullName: viewData['walletSettings']
                                          ['permissions'][i]['fullName'],
                                      permission: viewData['walletSettings']
                                          ['permissions'][i]['permission'],
                                    ),
                                  );
                                }

                                if (viewData['walletSettings']['permissions'][i]
                                        ['permission'] ==
                                    'APPROVER') {
                                  approvers.add(Permission(
                                      targetUsername: viewData['walletSettings']
                                          ['permissions'][i]['targetUsername'],
                                      fullName: viewData['walletSettings']
                                          ['permissions'][i]['fullName'],
                                      permission: viewData['walletSettings']
                                          ['permissions'][i]['permission']));
                                }

                                if (viewData['walletSettings']['permissions'][i]
                                        ['permission'] ==
                                    'INITIATOR') {
                                  initiators.add(Permission(
                                      targetUsername: viewData['walletSettings']
                                          ['permissions'][i]['targetUsername'],
                                      fullName: viewData['walletSettings']
                                          ['permissions'][i]['fullName'],
                                      permission: viewData['walletSettings']
                                          ['permissions'][i]['permission']));
                                }
                              }

                              appState.viewData![
                                  UpdateSharedAccessViewPageConfig.key] = {};

                              appState.viewData![
                                      UpdateSharedAccessViewPageConfig.key]
                                  ['walletAlias'] = viewData['walletAlias'];

                              appState.viewData![
                                          UpdateSharedAccessViewPageConfig.key]
                                      ['walletPublicKey'] =
                                  viewData['walletPublicKey'];

                              appState.viewData![
                                      UpdateSharedAccessViewPageConfig.key]
                                  ['viewers'] = viewers;

                              appState.viewData![
                                      UpdateSharedAccessViewPageConfig.key]
                                  ['isPrimaryWallet'] = 0;

                              appState.viewData![
                                      UpdateSharedAccessViewPageConfig.key]
                                  ['approvers'] = approvers;

                              appState.viewData![
                                      UpdateSharedAccessViewPageConfig.key]
                                  ['initiators'] = initiators;

                              appState.viewData![
                                          UpdateSharedAccessViewPageConfig.key]
                                      ['numberOfApprovalsNeeded'] =
                                  viewData!['walletSettings']
                                      ['numberOfApprovalsNeeded'];

                              appState.currentAction = PageAction(
                                  state: PageState.addPage,
                                  page: UpdateSharedAccessViewPageConfig);

                              print('================${appState.viewData}');
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
}

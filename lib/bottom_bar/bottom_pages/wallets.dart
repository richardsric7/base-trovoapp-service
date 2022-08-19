import 'dart:convert';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/constants.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../Custom_BlocObserver/notifire_clor.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class Wallets extends StatefulWidget {
  const Wallets({Key? key}) : super(key: key);

  @override
  State<Wallets> createState() => _WalletsState();
}

enum WalletAction { import, createNew }

enum WalletView { listWallets, addSubWallet, confirmAddSubWallet }

class _WalletsState extends State<Wallets> with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  WalletAction? action = WalletAction.createNew;
  WalletView walletView = WalletView.listWallets;
  final Authenticator _authenticator = Authenticator();
  bool isTileView = false;
  bool isImport = true;
  String? tag;
  String? description;
  String? secretKey;
  String password = '';
  late Account primaryWalletKeyPair;
  late Account newSubWalletKeyPair;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  List<Wallet>? wallets;
  Wallet? mainWallet;
  var claimedAssets;
  var unclaimedAssets;
  final _formKey = GlobalKey<FormState>();
  final _formKey2 = GlobalKey<FormState>();
  late RefreshController _refreshController;
  var actionIcon = Icons.add_circle_outline_sharp;
  var actionText = LanguageEn.addsubwallet;

  @override
  void initState() {
    super.initState();
    _refreshController = RefreshController(initialRefresh: false);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    userInfo = appState.userInfo!;
    assetBalances = appState.assetBalances;
    wallets =
        userInfo.wallets!.where((wallet) => wallet.primaryWallet == 0).toList();
    mainWallet =
        userInfo.wallets!.firstWhere((wallet) => wallet.primaryWallet == 1);
    nfts = appState.nfts;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
          resizeToAvoidBottomInset: false,
          backgroundColor: notifier.getwihitecolor,
          appBar: AppBar(
            centerTitle: true,
            title: Text(
              LanguageEn.wallets,
              style: TextStyle(
                  color: notifier.getblck,
                  fontWeight: FontWeight.bold,
                  fontFamily: fontsemibold),
            ),
            backgroundColor: notifier.getfavorites,
            elevation: 0,
          ),
          body: SmartRefresher(
            enablePullDown: true,
            controller: _refreshController,
            onRefresh: refreshData,
            child: ListView(
              children: [
                if (walletView == WalletView.listWallets) ...[
                  Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: Container(
                        color: notifier.getfavorites,
                        padding: EdgeInsets.all(8.sp),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.end,
                          children: [
                            GestureDetector(
                              onTap: () => setState(() {
                                isTileView = !isTileView;
                              }),
                              child: Padding(
                                padding: const EdgeInsets.all(8.0),
                                child: SvgPicture.asset(isTileView
                                    ? "assets/images/listview.svg"
                                    : "assets/images/tileview.svg"),
                              ),
                            ),
                          ],
                        )),
                  )
                ] else ...[
                  SizedBox(
                    height: height / 40,
                  )
                ],
                GestureDetector(
                  onTap: () {
                    appState.setActiveWallet = mainWallet;
                    appState.currentAction = PageAction(
                        state: PageState.addPage,
                        page: WalletDetailsViewPageConfig);
                  },
                  child: walletListItem(
                      mainWallet!.alias!.capitalizeFirst,
                      '4,014 USD',
                      notifier.getbluecolor,
                      assetBalances[mainWallet!.publicKey]['claimed'][0]),
                ),
                SizedBox(
                  height: height / 50,
                ),
                GestureDetector(
                  onTap: () => setState(() {
                    if (walletView == WalletView.listWallets) {
                      actionIcon = Icons.cancel_outlined;
                      actionText = LanguageEn.cancel;
                      walletView = WalletView.addSubWallet;
                    } else if (walletView == WalletView.addSubWallet) {
                      actionIcon = Icons.add_circle_outline_sharp;
                      actionText = LanguageEn.addsubwallet;
                      walletView = WalletView.listWallets;
                    } else if (walletView == WalletView.confirmAddSubWallet) {
                      actionIcon = Icons.cancel_outlined;
                      actionText = LanguageEn.cancel;
                      walletView = WalletView.addSubWallet;
                    }
                  }),
                  child: Column(
                    children: [
                      Icon(
                        actionIcon,
                        color: notifier.getbluecolor,
                        size: 35,
                      ),
                      SizedBox(
                        height: 5,
                      ),
                      Text(
                        actionText,
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w500,
                          color: notifier.getbluecolor,
                          fontFamily: fontbody,
                        ),
                      ),
                    ],
                  ),
                ),
                SizedBox(
                  height: height / 50,
                ),
                if (walletView == WalletView.addSubWallet) ...[
                  addSubwallet()
                ] else if (walletView == WalletView.confirmAddSubWallet) ...[
                  confirmAddSubwallet()
                ] else ...[
                  isTileView ? gridView() : walletListView(),
                ],
                Padding(
                    padding: EdgeInsets.only(
                        bottom: MediaQuery.of(context).viewInsets.bottom)),
              ],
            ),
          )),
    );
  }

  Widget gridView() {
    return Container(
      height: height / 2,
      child: GridView.count(
        primary: true,
        padding: const EdgeInsets.fromLTRB(15, 20, 15, 70),
        crossAxisCount: 2,
        mainAxisSpacing: 20,
        crossAxisSpacing: 20,
        childAspectRatio: 1.05,
        children: [
          for (var i = 0; i < wallets!.length; i++) ...[
            GestureDetector(
              onTap: () {
                appState.setActiveWallet = wallets![i];
                appState.currentAction = PageAction(
                    state: PageState.addPage,
                    page: WalletDetailsViewPageConfig);
              },
              child: walletTile(
                  wallets![i].alias!.capitalizeFirst!,
                  '4,014 USD',
                  notifier.getbluecolor80,
                  assetBalances[wallets![i].publicKey]['claimed'][0]),
            ),
          ]
        ],
      ),
    );
  }

  Widget walletTile(walletName, usdBal, color, asset) {
    var balance =
        "${formatHistoryNumber(double.parse(asset['amount']))} ${asset["assetCode"].toString().isEmpty ? 'XBN' : asset["assetCode"]}";
    return Card(
      elevation: 5,
      shadowColor: Colors.black,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      color: color,
      child: Stack(
        children: [
          Center(
            child: Padding(
              padding:
                  const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
              child: Image.asset(
                'assets/images/trovo_white.png',
                fit: BoxFit.cover,
                height: 100,
                width: 100,
              ),
            ),
          ),
          Center(
            child: Column(
              children: [
                SizedBox(
                  height: height / 50,
                ),
                Text(
                  walletName,
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getwihitecolor,
                  ),
                ),
                Padding(
                  padding:
                      const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
                  child: Text(
                    appState.hideBalances ? hideBalanceText : balance,
                    style: TextStyle(
                      fontSize: 18,
                      fontFamily: fontbody,
                      color: notifier.getwihitecolor,
                    ),
                  ),
                ),
                Padding(
                  padding:
                      const EdgeInsets.symmetric(vertical: 8, horizontal: 10),
                  child: Text(
                    appState.hideBalances ? hideBalanceText : usdBal,
                    style: TextStyle(
                      fontSize: 13,
                      fontFamily: fontbody,
                      color: notifier.getwihitecolor,
                    ),
                  ),
                ),
                Icon(
                  CupertinoIcons.eye_slash,
                  color: notifier.getwihitecolor,
                ),
              ],
            ),
          ),
        ],
      ), //SizedBox
    );
  }

  Widget walletListView() {
    return Column(
      children: [
        for (var i = 0; i < wallets!.length; i++) ...[
          GestureDetector(
            onTap: () {
              appState.activeWallet = wallets![i];
              appState.currentAction = PageAction(
                  state: PageState.addPage, page: WalletDetailsViewPageConfig);
            },
            child: walletListItem(
                wallets![i].alias!.capitalizeFirst!,
                '4,014 USD',
                notifier.getbluecolor80,
                assetBalances[wallets![i].publicKey]['claimed'][0]),
          ),
          SizedBox(
            height: height / 50,
          ),
        ],
        SizedBox(
          height: height / 15,
        ),
      ],
    );
  }

  Widget walletListItem(walletName, balanceUsd, color, asset) {
    var balance =
        "${formatHistoryNumber(double.parse(asset['amount']))} ${asset["assetCode"].toString().isEmpty ? 'XBN' : asset["assetCode"]}";
    return Container(
      height: height / 6.6,
      margin: EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        borderRadius: const BorderRadius.all(Radius.circular(20.0)),
        color: color,
        // color: colors[i - 1],
      ),
      child: Stack(children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.end,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
              child: Image.asset(
                'assets/images/trovo_white.png',
                fit: BoxFit.cover,
                height: 100,
                width: 100,
              ),
            ),
          ],
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0, vertical: 20.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Row(
                children: [
                  Text(
                    walletName,
                    style: TextStyle(
                      fontSize: 16,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  Text(
                    appState.hideBalances ? hideBalanceText : balance,
                    style: TextStyle(
                      fontSize: 20,
                      color: notifier.getwihitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                  SizedBox(
                    width: 20,
                  ),
                  Icon(
                    CupertinoIcons.eye_slash,
                    size: 25,
                    color: notifier.getwihitecolor,
                  ),
                ],
              ),
              SizedBox(height: height / 80),
              Text(
                appState.hideBalances ? hideBalanceText : balanceUsd,
                style: TextStyle(
                  fontWeight: FontWeight.w300,
                  fontSize: 13,
                  color: notifier.getwihitecolor,
                  fontFamily: fontbody,
                ),
              ),
            ],
          ),
        ),
      ]),
    );
  }

  Widget addSubwallet() {
    return Form(
      key: _formKey2,
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Card(
              shadowColor: Colors.black,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(15.0),
              ),
              color: notifier.getaddsubwalletgrey,
              child: Center(
                child: Column(
                  children: [
                    SizedBox(
                      height: height / 50,
                    ),
                    Container(
                      width: width / 1.4,
                      child: Text(
                        LanguageEn.abouttocreatesubwallet,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluecolor,
                        ),
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Text(
                      LanguageEn.chooseamethod,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        // fontWeight: FontWeight.bold,
                        fontFamily: fontsemibold,
                        color: notifier.getbluecolor,
                      ),
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      children: [
                        SizedBox(
                          width: width / 10,
                        ),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.import,
                            groupValue: action,
                            activeColor: notifier.getbluecolor,
                            onChanged: (value) => {
                              setState(
                                () {
                                  action = value;
                                },
                              )
                            },
                          ),
                        ),
                        Text(
                          LanguageEn.importexistingwallet,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.bold,
                            fontFamily: fontbody,
                            color: notifier.getbluecolor,
                          ),
                        ),
                      ],
                    ),
                    Row(
                      children: [
                        SizedBox(
                          width: width / 10,
                        ),
                        Transform.scale(
                          scale: 1.5,
                          child: Radio<WalletAction>(
                            value: WalletAction.createNew,
                            activeColor: notifier.getbluecolor,
                            groupValue: action,
                            onChanged: (value) => {
                              setState(
                                () {
                                  action = value;
                                },
                              )
                            },
                          ),
                        ),
                        Text(
                          LanguageEn.createnewwallet,
                          style: TextStyle(
                            fontSize: 15,
                            fontWeight: FontWeight.bold,
                            fontFamily: fontbody,
                            color: notifier.getbluecolor,
                          ),
                        ),
                      ],
                    ),
                    SizedBox(
                      height: height / 50,
                    ),
                  ],
                ),
              ),
            ),
          ),
          SizedBox(
            height: height / 20,
          ),
          // Tag name
          CustomTextFormField.textField(
            LanguageEn.tag,
            notifier.getbluecolor,
            Icons.tag,
            notifier.getgrey,
            notifier.getbluecolor,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            initialValue: tag,
            onChanged: (value) {
              setState(() {
                tag = value.trim().replaceAll(' ', '');
              });
            },
            onSaved: (value) {
              print('tag: $value');
              tag = value.trim().replaceAll(' ', '');
            },
            keyboardtype: TextInputType.text,
            maxLength: 6,
            validator: validateTag,
            helperText: tag == null || tag!.isEmpty
                ? ''
                : "${appState.userInfo!.username}_$tag",
          ),
          SizedBox(height: height / 50),
          CustomTextFormField.textField(
            LanguageEn.description,
            notifier.getbluecolor,
            Icons.description,
            notifier.getgrey,
            notifier.getbluecolor,
            notifier.getblck,
            notifier.getgrey,
            70.sp,
            300.sp,
            initialValue: description,
            onSaved: (value) {
              print('description: $value');
              description = value;
            },
            keyboardtype: TextInputType.text,
            maxLength: 100,
            validator: validateDescription,
          ),
          if (action == WalletAction.import) ...[
            SizedBox(height: height / 50),
            // Secret Key
            CustomPasswordFormField(
              LanguageEn.secretkey,
              notifier.getbluecolor,
              Icons.lock,
              notifier.getgrey,
              notifier.getbluecolor,
              notifier.getblck,
              70.sp,
              300.sp,
              validator: (value) {
                var trimmedVal = value!.trim().replaceAll(' ', '');
                if (trimmedVal.isEmpty) {
                  return LanguageEn.entersecretkeyempty;
                }

                if (trimmedVal.length < 56) {
                  return LanguageEn.secretkeyinvalid;
                }

                try {
                  TrovoWalletSDK().parseSecretKey(value);
                } catch (e) {
                  return 'Secret Key is invalid';
                }

                return null;
              },
              onSaved: (value) {
                print('email: $value');
                secretKey = value!.trim().replaceAll(' ', '');
              },
              maxLength: 56,
            ),
          ],
          SizedBox(height: height / 30),
          Button(
            LanguageEn.continuee,
            notifier.getbluecolor,
            notifier.getwihitecolor,
            onTap: () => submitForm(),
          ),
          SizedBox(height: height / 20),
        ],
      ),
    );
  }

  Widget confirmAddSubwallet() {
    return Column(children: [
      Padding(
        padding: const EdgeInsets.symmetric(horizontal: 20.0),
        child: Card(
          shadowColor: Colors.black,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15.0),
          ),
          color: notifier.getaddsubwalletgrey,
          child: Center(
            child: Form(
              key: _formKey,
              child: Column(
                children: [
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.4,
                    child: Text(
                      LanguageEn.requesttocreatesubwallet,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.bold,
                        fontFamily: fontbody,
                        color: notifier.getbluecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.tag,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  Text(
                    "${userInfo.username!}_$tag",
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.description,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  Text(
                    description!,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.method,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  Text(
                    action == WalletAction.import
                        ? LanguageEn.importsubwallet
                        : LanguageEn.createnewsubwallet,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Text(
                    LanguageEn.publickey,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontbody,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 15.0),
                    child: Text(
                      newSubWalletKeyPair.publicKey,
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w500,
                        fontFamily: fontbody,
                        color: notifier.getbluecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 30,
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
      SizedBox(height: height / 50),
      // Secret Key
      CustomPasswordFormField(
        LanguageEn.password,
        notifier.getbluecolor,
        Icons.lock,
        notifier.getgrey,
        notifier.getbluecolor,
        notifier.getblck,
        70.sp,
        300.sp,
        validator: validatePassword,
        onChanged: (value) {
          setState(() {
            password = value!.trim().replaceAll(' ', '');
          });
        },
        onSaved: (value) {
          print('email: $value');
          secretKey = value!.trim().replaceAll(' ', '');
        },
      ),
      SizedBox(height: height / 30),

      if (appState.biometricEnabled && password.isEmpty) ...[
        Button(
          LanguageEn.authorizewithbiometrics,
          notifier.getbluecolor,
          notifier.getwihitecolor,
          onTap: toggleSwitch,
        ),
      ] else ...[
        Button(
          LanguageEn.authorize,
          notifier.getbluecolor,
          notifier.getwihitecolor,
          onTap: handleAuthorization,
        ),
      ],
      SizedBox(height: height / 20),
    ]);
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  submitForm() async {
    print('submitting...');
    var form = _formKey2.currentState;
    if (!form!.validate()) {
      return;
    }
    form.save();
    generateKeyPairs();
    setState(() {
      walletView = WalletView.confirmAddSubWallet;
      actionIcon = Icons.arrow_circle_left_outlined;
      actionText = LanguageEn.back;
      password = '';
      secretKey = '';
    });
  }

  generateKeyPairs() {
    primaryWalletKeyPair =
        TrovoWalletSDK().parseSecretKey(appState.secretKeys[0]);
    setState(() {
      if (action == WalletAction.import) {
        try {
          // parse supplied secret to get the keypair
          newSubWalletKeyPair = TrovoWalletSDK().parseSecretKey(secretKey);
        } catch (e) {
          popup(context, title: 'Error!', message: 'Secret Key is invalid');
        }
      } else {
        // generate keypair for the new subwallet
        newSubWalletKeyPair = TrovoWalletSDK().createAccount();
      }
    });
  }

  String? validateTag(String? value) {
    print('validating tag...');
    if (value!.isEmpty) return 'Enter wallet tag';

    String pattern = r'^[a-zA-Z0-9]*$';
    RegExp regex = new RegExp(pattern);

    if (!regex.hasMatch(value.trim().replaceAll(' ', ''))) {
      return 'Invalid tag name';
    }

    return null;
  }

  String? validateDescription(String? value) {
    print('validating description...');
    if (value!.isEmpty) return 'Enter wallet description';

    return null;
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        sendDataToServer();
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  void handleAuthorization() {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      sendDataToServer();
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }

  void sendDataToServer() async {
    showLoader(context);

    try {
      // make initial request to the server using the
      // following credentials
      Map map = {
        "publickey": newSubWalletKeyPair.publicKey,
        "walletTag": tag,
        "WalletDescription": description,
      };
      String requestBody = jsonEncode(map);

      print(requestBody);

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        var messageLength = responseData['data']['messages'].length;
        var messageShown = 0;

        print('new dialog $messageLength');
        print('messagecount $messageShown');
        await postProcessData(
            messageShown, messageLength, responseData['data']);
        // print('sending full data to server.........');
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  postProcessData(messageShown, messageLength, data) {
    print('messageShown: $messageShown messageLength $messageLength');
    // we would like to display all messages returned from the initial
    // request to server using a popup. In order to achieve that we
    // employ the use of a little recursion here. Please recursive
    // functions can turn into a nightmare fast so be carefull here.
    if (messageShown <= messageLength - 1) {
      showResponseMessage(
          context,
          data['messages'][messageShown],
          () => {
                print('postProcessData: $messageShown'),
                postProcessData(messageShown, messageLength, data),
              });

      messageShown++;
      return;
    }

    sendFullDataToServer(data);
    return;
  }

  void sendFullDataToServer(responseBody) async {
    try {
      showLoader(context);
      // get primary signature
      var primarySignature = TrovoWalletSDK().signBase64Txn(
        primaryWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      // get secondary signature
      var subWalletSignature = TrovoWalletSDK().signBase64Txn(
        newSubWalletKeyPair.secretKey,
        responseBody['transaction'],
        responseBody['networkPassPhrase'],
      );

      print('this is primary sign: $primarySignature');
      print('this is subwallet sign: $subWalletSignature');

      Map map = {
        "publickey": responseBody['publicKey'],
        "walletTag": responseBody['walletTag'],
        "WalletDescription": responseBody['walletDescription'],
        "transaction": responseBody['transaction'],
        "primarySignature": primarySignature,
        "subWalletSignature": subWalletSignature,
        "transactionId": responseBody['transactionId'],
        "networkPassPhrase": responseBody['networkPassPhrase'],
        "channelAccount": responseBody['channelAccount'],
        "channelAccountSignature": responseBody['channelAccountSignature'],
        "subWalletMustSign": responseBody['subWalletMustSign'],
      };

      String requestBody = jsonEncode(map);

      print('this is request body: $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/subwallet',
        body: requestBody,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        // add the secret key of this new subwallet to
        // the existing list of secrets
        appState.secretKeys.add(newSubWalletKeyPair.secretKey);
        // store back the list of secret keys but this time it
        // contains the secret key of the newly created subwallet
        await StoreData().storeInsertData('secretKey', appState.secretKeys);
        await updateUserInfo();
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
    }

    hideLoader(context);
  }

  void resetForm() {
    tag = '';
    description = '';
    secretKey = '';
    actionIcon = Icons.add_circle_outline_sharp;
    actionText = LanguageEn.addsubwallet;
    walletView = WalletView.listWallets;
  }

  Future<void> updateUserInfo() async {
    var keyPair =
        TrovoWalletSDK().parseSecretKey(primaryWalletKeyPair.secretKey);
    Map responseData = await makeGetRequest(
        uri: '/v1/users/${userInfo.username!.trim().replaceAll(' ', '')}',
        signer: keyPair.publicKey,
        publicKey: keyPair.publicKey,
        secretKey: keyPair.secretKey);

    print('response: ${responseData}');

    if (responseData['statusCode'] == 200) {
      await storeUserInfo(responseData['data']);
    }
  }

  Future<void> storeUserInfo(userInfoMap) async {
    print('userInfoMap: ${userInfoMap['userData']}');
    var userInfo = userInfoMap['userData'] ?? {};
    var assetBalances = userInfoMap['assetBalances'] ?? {};
    var nfts = userInfoMap['nfts'] ?? {};
    var thirdPartyWalletAccess = userInfoMap['thirdPartyWalletAccess'] ?? [];
    var defaultAssets = userInfoMap['defaultAssets'] ?? [];

    await StoreData().storeInsertData('userInfo', userInfo);
    await StoreData().storeInsertData('assetBalances', assetBalances);
    await StoreData().storeInsertData('nftBalances', nfts);
    await StoreData()
        .storeInsertData('thirdPartyWalletAccess', thirdPartyWalletAccess);
    await StoreData().storeInsertData('defaultAssets', defaultAssets);

    // save useInfo to appstate
    appState.setUser = UserInfo().deserializeJson(userInfo);
    appState.setNFTs = nfts;
    appState.setassetBalances = assetBalances;
    appState.activeWallet = appState.userInfo!.wallets!.firstWhere(
        (wallet) => wallet.publicKey == newSubWalletKeyPair.publicKey);
    appState.activeWallet!.secretKey = newSubWalletKeyPair.secretKey;
    appState.currentAction =
        PageAction(state: PageState.addPage, page: CongratulationsPageConfig);
    resetForm();
    print('stored new user data.................');
  }

  void refreshData() async {
    try {
      await appState.refreshData();
      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/services/push_fcm_service.dart';
import 'package:trovo_wallet/storage/cache.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../../Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import '../../Custom_BlocObserver/button/custtom_button.dart';
import '../../Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import '../../Custom_BlocObserver/custtom_textfild/custtompassword.dart';
import '../../Models/User.dart';
import '../../network/requests.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../storage/store.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import '../../widgets/loader.dart';
import '../../widgets/popups.dart';
import '../../widgets/termsOfService.dart';

class ImportWallet extends StatefulWidget {
  const ImportWallet({Key? key}) : super(key: key);

  @override
  State<ImportWallet> createState() => _ImportWalletState();
}

class _ImportWalletState extends State<ImportWallet> {
  late ColorNotifier notifier;
  bool hasAgreed = false;
  bool showTermsError = false;
  bool usePassPhrase = false;
  late FocusNode passPhraseFocusNode;
  late FocusNode secretKeyFocusNode;
  final _formKey = GlobalKey<FormState>();
  String? username;
  String? passPhrase;
  String? secretKey;
  String? password;
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
    super.initState();
    getdarkmodepreviousstate();
    passPhraseFocusNode = FocusNode();
    secretKeyFocusNode = FocusNode();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
            context, notifier.getwihitecolor, "", notifier.getblck,
            height: height / 15),
        body: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              SizedBox(height: height / 20),
              Row(
                children: [
                  SizedBox(width: width / 15),
                  Form(
                    key: _formKey,
                    child: Column(
                      // crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Text(
                              LanguageEn.import,
                              style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 26.sp,
                                  fontFamily: fontsemibold),
                            ),
                            SizedBox(
                              width: width / 50,
                            ),
                            Text(
                              LanguageEn.wallet,
                              style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  fontSize: 26.sp,
                                  fontFamily: fontsemibold),
                            ),
                          ],
                        ),
                        SizedBox(height: height / 15),
                        // Email address
                        CustomTextFormField.textField(
                          LanguageEn.usernameoremail,
                          notifier.getbluecolor,
                          Icons.email,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          300.sp,
                          validator: (value) {
                            var trimmedVal = value!.trim().replaceAll(' ', '');
                            if (trimmedVal.isEmpty) {
                              return LanguageEn.usernameoremailempty;
                            }

                            if (trimmedVal.length < 3) {
                              return LanguageEn.usernameoremailinvalid;
                            }
                          },
                          onSaved: storeUsernameOrEmail,
                          keyboardtype: TextInputType.emailAddress,
                        ),
                        Row(
                          children: [
                            Container(
                              width: width / 1.2,
                              child: checkUsePassphrase(),
                            ),
                          ],
                        ),
                        if (usePassPhrase) ...[
                          // Pass phrase/Mnemonic
                          passPhraseInput(
                            LanguageEn.passphrase,
                            notifier.getbluecolor,
                            notifier.getgrey,
                            notifier.getblck,
                            notifier.getgrey,
                            100.sp,
                            300.sp,
                            validator: (value) {
                              if (value.isEmpty) {
                                return LanguageEn.enterpassphraseempty;
                              }
                            },
                            onSaved: (value) {
                              passPhrase = value;
                            },
                            minLines: 3,
                            maxLines: null,
                            keyboardtype: TextInputType.multiline,
                            focusNode: passPhraseFocusNode,
                          ),
                        ] else ...[
                          // Secret Key
                          CustomPasswordFormField(
                            LanguageEn.secretkey,
                            notifier.getbluecolor,
                            Icons.lock,
                            notifier.getgrey,
                            notifier.getprefixicon,
                            notifier.getblck,
                            70.sp,
                            300.sp,
                            validator: (value) {
                              var trimmedVal =
                                  value!.trim().replaceAll(' ', '');
                              if (trimmedVal.isEmpty) {
                                return LanguageEn.entersecretkeyempty;
                              }

                              if (trimmedVal.length < 56) {
                                return LanguageEn.secretkeyinvalid;
                              }
                            },
                            onSaved: (value) {
                              print('email: $value');
                              secretKey = value!.trim().replaceAll(' ', '');
                            },
                            maxLength: 56,
                            focusNode: secretKeyFocusNode,
                          )
                        ],
                        SizedBox(height: height / 40),
                        CustomPasswordFormField(
                          LanguageEn.password,
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70.sp,
                          300.sp,
                          onChanged: (value) {
                            setState(() {
                              password = value!.trim().replaceAll(' ', '');
                            });
                          },
                          validator: validatePassword,
                        ),
                        SizedBox(height: height / 80),
                        CustomPasswordFormField(
                          LanguageEn.confirmPassword,
                          notifier.getbluecolor,
                          Icons.lock,
                          notifier.getgrey,
                          notifier.getprefixicon,
                          notifier.getblck,
                          70.sp,
                          300.sp,
                          validator: validateConfirmPassword,
                        ),
                        // Terms of Service
                        TermsOfService(
                          value: hasAgreed,
                          showError: showTermsError,
                          onChanged: (bool? value) {
                            setState(() {
                              hasAgreed = value!;
                            });
                          },
                        ),
                      ],
                    ),
                  )
                ],
              ),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () => saveForm(),
              ),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                    bottom: MediaQuery.of(context).viewInsets.bottom),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget checkUsePassphrase() {
    return Row(
      children: [
        Transform.scale(
          scale: 1.sp,
          child: Checkbox(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.all(
                Radius.circular(5.sp),
              ),
            ),
            activeColor: notifier.isDark
                ? notifier.getbluecolor50
                : notifier.getbluecolor90,
            side: BorderSide(
              color: notifier.isDark
                  ? notifier.getbluecolor50
                  : notifier.getbluecolor90,
            ),
            value: usePassPhrase,
            onChanged: (bool? value) {
              setState(() {
                usePassPhrase = value!;
                if (usePassPhrase) {
                  passPhraseFocusNode.requestFocus();
                } else {
                  secretKeyFocusNode.requestFocus();
                }
              });
            },
          ),
        ),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Text(
                  LanguageEn.enterpassphrase,
                  style: TextStyle(
                      fontSize: height / 55,
                      color: notifier.getblck,
                      fontFamily: fontbody),
                ),
              ],
            ),
          ],
        )
      ],
    );
  }

  Widget passPhraseInput(
    labletext,
    focuscolor,
    lablecolor,
    textcolor,
    bordercolor,
    h,
    w, {
    onChanged,
    maxLength,
    minLines,
    maxLines,
    validator,
    onSaved,
    keyboardtype,
    focusNode,
  }) {
    return ScreenUtilInit(
      builder: (context, child) => Container(
        color: Colors.transparent,
        height: h,
        width: w,
        child: TextFormField(
          focusNode: focusNode,
          maxLength: maxLength,
          minLines: minLines,
          maxLines: maxLines,
          style: TextStyle(color: textcolor, fontFamily: fontbody),
          cursorColor: lablecolor,
          onChanged: onChanged,
          decoration: InputDecoration(
            hintText: labletext,
            hintStyle: TextStyle(color: lablecolor),
            disabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            labelStyle: TextStyle(color: lablecolor),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(15.sp),
            ),
            enabledBorder: OutlineInputBorder(
              borderSide: BorderSide(color: bordercolor, width: 1),
              borderRadius: BorderRadius.circular(15.sp),
            ),
            focusedBorder: OutlineInputBorder(
              borderSide: BorderSide(color: focuscolor, width: 1),
              borderRadius: BorderRadius.circular(15.sp),
            ),
          ),
          keyboardType: keyboardtype,
          validator: validator,
          onSaved: onSaved,
        ),
      ),
    );
  }

  String? storeUsernameOrEmail(String? value) {
    var currValue = value!.trim().replaceAll(' ', '');
    if (currValue.isEmpty) {
      return LanguageEn.emailvalidateempty;
    }

    setState(() {
      username = currValue;
    });
    return null;
  }

  saveForm() async {
    try {
      print('saving form...');
      final form = _formKey.currentState;
      if (!form!.validate()) {
        checkTerms();
        return;
      }

      // check that terms and conditions has been accepted
      if (!checkTerms()) return;

      showLoader(context);
      form.save();

      Account? creds = usePassPhrase
          ? await getCredsFromPassPhrase()
          : parseKey(secretKey!)!;

      if (creds != null) {
        // store these credentials and the password before making request
        // to get userinfo from the server. This way we can use these stored data
        // to create new user account if the provided user account does not exist
        appState.setTempPassword = password;
        appState.setTempPublicKey = creds.publicKey;
        appState.setTempSecretKey = creds.secretKey;
        String? token = await StoreData().storeGetData('token');

        if (token == null) {
          token = await FCM().getPushNotificationToken();
        }

        Map responseData = await makeGetRequest(
            uri: '/v1/users/${username}?type=import&pnt=$token',
            signer: creds.publicKey,
            publicKey: creds.publicKey,
            secretKey: creds.secretKey);
        print('response: ${responseData}');

        if (responseData['statusCode'] == 200) {
          fetchNotifications(appState);
          getFiatRates(creds.publicKey, creds.secretKey, creds.publicKey,
              username, appState);
          storeUserInfo(responseData['data']);
          appState.currentAction =
              PageAction(state: PageState.addPage, page: FingerprintPageConfig);
        } else if (responseData['statusCode'] == 404) {
          accountNotFoundPopup(context);
        } else {
          // must be some sort of server error
          // let's throw it
          popup(context,
              title: LanguageEn.error,
              message: responseData['data']['message']);
        }
      }
      hideLoader(context);
    } catch (e) {
      hideLoader(context);
      print('object');
      print(e);
    }
  }

  storeUserInfo(userInfoMap) async {
    print('this is userinfo map: ${userInfoMap}');
    var userInfo = userInfoMap['userData'] ?? {};
    var assetBalances = userInfoMap['assetBalances'] ?? {};
    var nfts = userInfoMap['nfts'] ?? {};
    var walletsSharedWithUser = userInfoMap['walletsSharedWithUser'] ?? [];
    var defaultAssets = userInfoMap['defaultAssets'] ?? [];

    // delete all user data already stored on the app
    await StoreData().storeDeleteData();

    await StoreData().storeInsertData('userInfo', userInfo);
    await StoreData().storeInsertData('assetBalances', assetBalances);
    await StoreData().storeInsertData('nfts', nfts);
    await StoreData()
        .storeInsertData('walletsSharedWithUser', walletsSharedWithUser);
    await StoreData().storeInsertData('defaultAssets', defaultAssets);
    await StoreData().storeInsertData('isFirstTime', false);
    await StoreData().storeInsertData('password', appState.tempPassword);
    await StoreData().storeInsertData('publicKey', appState.tempPublicKey);
    await StoreData()
        .storeInsertData('secretKey', <String>[appState.tempSecretKey]);

    // save useInfo to appstate
    appState.setUser = UserInfo().deserializeJson(userInfo);
    appState.setNFTs = nfts;
    appState.setSharedWallets = walletsSharedWithUser;
    appState.assetBalances = assetBalances;

    // save secrets to appstate
    appState.setSecretKeys = await StoreData().storeGetData('secretKey');
    appState.setPassword = appState.tempPassword;
    appState.activeWallet = appState.userInfo!.wallets!
        .firstWhere((wallet) => wallet.primaryWallet == 1);
  }

  String? validatePassword(value) {
    print('password: $value');
    if (value.isEmpty) {
      //return "Enter a password";
      return LanguageEn.passwordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      //return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    return null;
  }

  String? validateConfirmPassword(value) {
    print('confirm password: ${value.trim().replaceAll(' ', '')} & $password');
    if (value.isEmpty) {
      // return "Confirm your password";
      return LanguageEn.confirmpasswordemptyerror;
    }

    if (value.trim().replaceAll(' ', '').length < 6) {
      // return 'Use 6 characters or more for your password';
      return LanguageEn.hinterrorpassword;
    }

    if (password != value.trim().replaceAll(' ', '')) {
      //  return 'Those passwords didn\’t match. Try again.';
      return LanguageEn.passwordmismatcherror;
    }

    return null;
  }

  bool checkTerms() {
    if (!hasAgreed) {
      // show error message if the user has not
      // agreed to terms and conditions
      setState(() {
        showTermsError = true;
      });
      return false;
    } else {
      setState(() {
        showTermsError = false;
      });
      return true;
    }
  }

  Future<Account?> getCredsFromPassPhrase() async {
    try {
      var trimmedPassprase = passPhrase!.trimLeft().trimRight();
      print('trimmed pass phrase ${trimmedPassprase}');
      Account account = await TrovoWalletSDK()
          .retrieveCredentialsFromPassPhrase(trimmedPassprase);
      print(account);
      return account;
    } catch (e) {
      print('from getSecretKey');
      print(e);
      // must be some sort of server error
      // let's throw it
      popup(context,
          title: LanguageEn.error, message: LanguageEn.invalidcredentials);
      return null;
    }
  }

  Account? parseKey(String secretKey) {
    try {
      Account account =
          TrovoWalletSDK().parseSecretKey(secretKey.toUpperCase());
      print(account);
      return account;
    } catch (e) {
      print('from parseKey');
      print(e);
      // must be some sort of server error
      // let's throw it
      popup(context,
          title: LanguageEn.error, message: LanguageEn.invalidcredentials);
      return null;
    }
  }

  @override
  void dispose() {
    passPhraseFocusNode.dispose();
    secretKeyFocusNode.dispose();
    super.dispose();
  }
}

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/custtompassword.dart';
import 'package:trovo_wallet/models/Wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class AuthorizeActionView extends StatefulWidget {
  const AuthorizeActionView({Key? key}) : super(key: key);

  @override
  State<AuthorizeActionView> createState() => _AuthorizeActionViewState();
}

class _AuthorizeActionViewState extends State<AuthorizeActionView> {
  late DataProvider appState;
  late ColorNotifier notifier;
  String password = '';
  final formKey = GlobalKey<FormState>();
  final Authenticator _authenticator = Authenticator();
  Wallet? activeWallet;
  var viewData;
  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    viewData = appState.viewData![AuthorizeActionViewPageConfig.key];
    activeWallet = appState.activeWallet;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(height: height / 10),
              Text(
                '${LanguageEn.authorize} ${viewData!['serviceShortName']}',
                textAlign: TextAlign.center,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 27.sp),
              ),
              SizedBox(height: height / 50),
              Padding(
                padding: const EdgeInsets.all(15.0),
                child: Text(
                  viewData!['description'].toString().capitalizeFirst!,
                  textAlign: TextAlign.center,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontSize: 19.sp,
                      fontFamily: fontbody),
                ),
              ),
              SizedBox(
                height: height / 10,
              ),
              Form(
                key: formKey,
                child: CustomPasswordFormField(
                  LanguageEn.password,
                  notifier.getbluewhitecolor,
                  Icons.lock,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  70.sp,
                  300.sp,
                  validator: validatePassword,
                  onChanged: (value) {
                    setState(() {
                      password = value!.trim().replaceAll(' ', '');
                    });
                  },
                ),
              ),
              SizedBox(
                height: height / 20,
              ),
              if (appState.biometricEnabled && password.isEmpty) ...[
                Button(
                  LanguageEn.authorizewithbiometrics,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: toggleSwitch,
                ),
              ] else ...[
                Button(
                  LanguageEn.authorize,
                  notifier.getbluecolor,
                  wihitecolor,
                  onTap: handleAuthorization,
                ),
              ],
              SizedBox(height: height / 50.5),
              ButtonOutlined(
                LanguageEn.cancel,
                notifier.getwihitecolor,
                notifier.getbluewhitecolor,
                onTap: () {
                  Navigator.of(context).pop();
                },
              ),
              Padding(
                  padding: EdgeInsets.only(
                      bottom: MediaQuery.of(context).viewInsets.bottom)),
            ],
          ),
        ),
      ),
    );
  }

  void handleAuthorization() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      authorize();
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        authorize();
        // aparently we need the code below to make the
        // screen updata to show loader
        // after authorizing with biometrics
        setState(() {});
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }

  authorize() async {
    print('sending to server.... $viewData');

    try {
      showLoader(context);
      // sign transaction

      Map responseData = await makePostRequest(
        uri:
            '/v1/users/servicelinks/authorize/approval/${viewData['targetUser']}?ownerUsername=${viewData['ownerUsername']}&authId=${viewData['authId']}',
        body: "",
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
      );

      print('response: $responseData');
      if (responseData['statusCode'] == 200) {
        appState.viewData![TransactionSuccessViewPageConfig.key] =
            responseData['data'];
        showSuccessAlert(context,
            onTap: () => appState.currentAction = PageAction(
                state: PageState.replaceAll, page: BottomHomePageConfig));
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
}

import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class CompleteAccountRecovery extends StatefulWidget {
  const CompleteAccountRecovery({Key? key}) : super(key: key);

  @override
  State<CompleteAccountRecovery> createState() => _CompleteAccountRecovery();
}

class _CompleteAccountRecovery extends State<CompleteAccountRecovery>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  String password = '';
  final formKey = GlobalKey<FormState>();
  var viewData;
  bool hasBackedUp = false;

  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    viewData = appState.viewData![CompleteAccountRecoveryViewPageConfig.key];

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
                context, notifier.getwihitecolor, "", notifier.getblck,
                height: height / 15)
            .getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    LanguageEn.account,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 26.sp,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(
                    width: width / 50,
                  ),
                  Text(
                    LanguageEn.recovery,
                    style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 26.sp,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              SizedBox(
                height: height / 20,
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
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Column(
                      children: [
                        SizedBox(
                          height: height / 50,
                        ),
                        Text(
                          LanguageEn.haveyoubackedup,
                          overflow: TextOverflow.visible,
                          style: TextStyle(
                              fontSize: 15,
                              fontWeight: FontWeight.w700,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold),
                        ),
                        SizedBox(
                          height: height / 50,
                        ),
                      ],
                    ),
                  ),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  Transform.scale(
                    scale: 1.sp,
                    child: Checkbox(
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.all(
                          Radius.circular(5.sp),
                        ),
                      ),
                      activeColor: notifier.getbluecolor,
                      side: BorderSide(color: notifier.getbluewhitecolor),
                      value: hasBackedUp,
                      onChanged: (value) => setState(() {
                        hasBackedUp = value!;
                      }),
                    ),
                  ),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Container(
                        padding: const EdgeInsets.all(16.0),
                        width: width / 1.2,
                        child: Text(
                          LanguageEn.ihavebackedupmywallet,
                          style: TextStyle(
                              fontSize: height / 55,
                              color: notifier.getgrey,
                              fontFamily: fontbody),
                        ),
                      ),
                    ],
                  )
                ],
              ),
              SizedBox(
                height: height / 20,
              ),
              Button(
                LanguageEn.completeaccountrecovery,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  if (hasBackedUp) {
                    completeAccountRecovery();
                  } else {
                    popup(context,
                        title: LanguageEn.important,
                        message: LanguageEn.ensurebackedup);
                  }
                },
              ),
              SizedBox(
                height: height / 20,
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
    sendDataToServer();
  }

  void completeAccountRecovery() {
    var messageLength = viewData['messages'].length;
    var messageShown = 0;
    postProcessData(messageShown, messageLength, viewData);
  }

  sendDataToServer() async {
    print('sending to server....');

    try {
      showLoader(context);
      viewData['commit'] = 1;
      String requestBody = jsonEncode(viewData);
      print('this is request body $requestBody');

      Map responseData = await makePostRequest(
        uri: '/v1/users/account/recover',
        body: requestBody,
        signer: appState.tempPublicKey,
        secretKey: appState.tempSecretKey, // the primary wallet secret key
        publicKey: appState.tempPublicKey,
      );

      print('response: $responseData');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        appState.currentAction = PageAction(
            state: PageState.addPage,
            page: AccountRecoverySuccessViewPageConfig);
      } else {
        hideLoader(context);
        popup(context,
            title: LanguageEn.error, message: responseData['data']['error']);
      }
    } catch (e) {
      print(e);
      popup(context, title: LanguageEn.error, message: e.toString());
      hideLoader(context);
    }
  }
}

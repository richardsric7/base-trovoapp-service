import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../../custom_bloc_observer/button/custtom_button.dart';
import '../../custom_bloc_observer/fonts.dart';
import '../../custom_bloc_observer/notifire_clor.dart';
import '../../router/page_actions.dart';
import '../../router/ui_pages.dart';
import '../../storage/state.dart';
import '../../utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class WelcomeToSharedAccess extends StatelessWidget {
  WelcomeToSharedAccess({Key? key}) : super(key: key);
  late DataProvider appState;
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          LanguageEn.sharedaccess,
          notifier.getblck,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(
                    vertical: 15.0, horizontal: 25.0),
                child: RichText(
                  text: TextSpan(
                    text:
                        'Welcome to Shared Access. Here you can give others various access rights to your wallet.\n\n',
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontbody,
                      color: notifier.getbluewhitecolor,
                    ),
                    children: [
                      TextSpan(
                        text: 'Viewer Access ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text:
                            'enables other users to view your wallet balance, receive payment into your wallet and view your wallet history.\n\n',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: 'Initiator Access ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text:
                            'enables other users in addition to viewer access, to initiate a transaction from your wallet and pass it to the appropriate approvers to approve.\n\n',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: 'Approver Access ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text:
                            'enables other users in addition to viewer access, to become approvers on your wallet, this means that whenever a transaction is initiated from your wallet by the initiators, it must be approved by the required number of approvers out of the added approvers for the transaction to successfully go through.\n\n',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text: 'Note: ',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextSpan(
                        text:
                            'You must be careful when giving approver access because once you give others approver access on any of your wallets, the wallet seizes to be your sole wallet, it now becomes a jointly owned wallet that must get the approval of all the required approvers for any transaction to take place on it successfully.\n\nIf you happen not to be an initiator on the wallet with approver access, even though the wallet is originally your wallet, you will no longer be able to initiate transactions from the wallet.\n\nAlso if you happen not to be an Approver on the wallet with approver access enabled, you cannot approve any transaction on the wallet as well, you can only view the wallet going forward.\n\nIndividuals can use Approver Access, however, it is best suited for organizations, businesses, associations, and any other use case where more than 1 person is required to operate a wallet.',
                        style: TextStyle(
                          fontSize: 17,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                  textAlign: TextAlign.justify,
                ),
              ),
              SizedBox(height: height / 20),
              Button(
                LanguageEn.continuee,
                notifier.getbluecolor,
                wihitecolor,
                onTap: () async {
                  await StoreData()
                      .storeInsertData('introducedSharedAccess', true);
                  appState.setIntroducedSharedAccess = true;
                  appState.currentAction = PageAction(
                      state: PageState.replace,
                      page: SharedAccessViewPageConfig);
                },
              ),
              SizedBox(height: height / 15),
            ],
          ),
        ),
      ),
    );
  }
}

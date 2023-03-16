import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
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

class ProceedsPayOut extends StatefulWidget {
  const ProceedsPayOut({Key? key}) : super(key: key);

  @override
  State<ProceedsPayOut> createState() => _ProceedsPayOut();
}

class _ProceedsPayOut extends State<ProceedsPayOut>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
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
              'Payout Proceeds',
              notifier.getbluewhitecolor,
              height: height / 15,
            ).getBar(),
            SizedBox(
              height: height / 30,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Enter Amount',
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
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: CustomTextFormField.textField(
                    'Amount',
                    notifier.getbluecolor,
                    null,
                    notifier.getgrey,
                    null,
                    notifier.getblck,
                    notifier.getgrey,
                    70.sp,
                    300.sp,
                    // controller: referrerController,
                    // validator: validateReferrer,
                    onSaved: (value) {},
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Description of Payout',
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
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: multilineInput(
                    '',
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
                    onSaved: (value) {},
                    minLines: 3,
                    maxLines: null,
                    keyboardtype: TextInputType.multiline,
                  ),
                ),
              ],
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Container(
                    width: width / 1.2,
                    child: Text(
                      'Please request for approval from the following users in order to proceed',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
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
              child: Card(
                shadowColor: Colors.black,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(15.0),
                ),
                color: notifier.isDark
                    ? notifier.getbluecolor90
                    : notifier.getaddsubwalletgrey,
                child: Center(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 20.0),
                    child: Wrap(
                      children: [
                        userItem(
                          'Onoja',
                          null,
                          foreColor: notifier.getwihitecolor,
                          backColor: notifier.getbluewhitecolor,
                        ),
                        userItem(
                          'Muche',
                          null,
                          foreColor: notifier.getwihitecolor,
                          backColor: notifier.getbluewhitecolor,
                        ),
                        userItem(
                          'Nancy',
                          null,
                          foreColor: notifier.getwihitecolor,
                          backColor: notifier.getbluewhitecolor,
                        ),
                        userItem(
                          'Obiaruku',
                          null,
                          foreColor: notifier.getwihitecolor,
                          backColor: notifier.getbluewhitecolor,
                        ),
                      ],
                    ),
                  ),
                ),
              ),
            ),
            SizedBox(
              height: height / 20,
            ),
            Button(
              'Request Approval',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                appState.viewData![SuccessViewPageConfig.key] = {
                  'title': '',
                  'message':
                      'Your Payout proceed Requests for [Atlantis Asset] have been submitted and are awaiting approval. You’ll be notified when all requests have been approved.',
                };
                appState.currentAction = PageAction(
                    state: PageState.replace, page: SuccessViewPageConfig);
              },
            ),
            SizedBox(
              height: height / 10,
            ),
          ],
        ),
      ),
    );
  }
}

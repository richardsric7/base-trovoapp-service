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

class BuyTokens extends StatefulWidget {
  const BuyTokens({Key? key}) : super(key: key);

  @override
  State<BuyTokens> createState() => _BuyTokens();
}

class _BuyTokens extends State<BuyTokens> with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  List<String> currencies = [
    'e-Naira',
    'TROV',
    'XBN',
  ];

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  List<DropdownMenuItem<String>> get getCurrencyOptions {
    List<DropdownMenuItem<String>> options = [];
    currencies.forEach((item) {
      options.add(DropdownMenuItem(
          child: Text(
            item,
            overflow: TextOverflow.ellipsis,
          ),
          value: item));
    });
    return options;
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
              'Buy Atlantis 1',
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
                    'Select Currency',
                    style: TextStyle(
                      fontSize: 12,
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
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {},
                getCurrencyOptions,
                null,
                'Select currency',
                context,
                null,
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: Text(
                    'Quantity of Atlantis 1',
                    style: TextStyle(
                      fontSize: 12,
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
                    'Quantity',
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
                    'Amount',
                    style: TextStyle(
                      fontSize: 12,
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
              height: height / 30,
            ),
            Button(
              'Pay',
              notifier.getbluecolor,
              wihitecolor,
              onTap: () {
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: ConfirmBuyViewPageConfig,
                );
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

  Widget confirmLiensAndEncumbrance() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
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
            value: true,
            onChanged: (bool? value) {
              setState(() {});
            },
          ),
        ),
        Container(
          width: width / 1.2,
          child: Text(
            'I confirm that this asset is completely free of all liens and encumbrance',
            overflow: TextOverflow.visible,
            style: TextStyle(
                fontSize: 15,
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody),
          ),
        ),
      ],
    );
  }
}

Widget CheckItem(
  String name,
  void Function()? onClick, {
  required Color backColor,
  required Color foreColor,
  required Color borderColor,
  double? fontSize: 15,
}) {
  return Padding(
    padding: const EdgeInsets.all(3.0),
    child: Container(
      decoration: BoxDecoration(
          border: Border.all(color: borderColor, width: 1),
          borderRadius: const BorderRadius.all(Radius.circular(10.0)),
          color: backColor),
      child: Padding(
        padding: const EdgeInsets.all(15.0),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                  color: foreColor, fontFamily: fontbody, fontSize: fontSize),
            ),
            SizedBox(
              width: width / 70,
            ),
          ],
        ),
      ),
    ),
  );
}

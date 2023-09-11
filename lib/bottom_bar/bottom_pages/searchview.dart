import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import '../../../custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import '../../../utils/medeiaqury/medeiaqury.dart';

class SearchView extends StatefulWidget {
  const SearchView({Key? key}) : super(key: key);

  @override
  State<SearchView> createState() => _SearchViewState();
}

class _SearchViewState extends State<SearchView> {
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: PreferredSize(
          preferredSize: Size.fromHeight(70),
          // here the desired height
          child: AppBar(
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            title: Text(
              "search".tr(),
              style: TextStyle(
                  fontSize: 20,
                  color: notifier.getblck,
                  fontFamily: 'Gilroy_Bold'),
            ),
          ),
        ),
        body: SingleChildScrollView(
          child: Padding(
            padding: const EdgeInsets.only(left: 10, right: 10),
            child: Column(
              children: [
                SizedBox(height: height / 80),
                Center(
                  child: Customtextfild.textField(
                      "Search  assets, nfts...",
                      notifier.getbluecolor,
                      Icons.search_rounded,
                      notifier.getgrey,
                      notifier.getblck,
                      notifier.getgrey,
                      notifier.getgrey,
                      50,
                      310),
                ),
                SizedBox(height: height / 20),
                Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Center(
                      child: Text(
                        "thisissearchview".tr(),
                        style: TextStyle(
                          fontSize: 13,
                          fontWeight: FontWeight.bold,
                          fontFamily: fontsemibold,
                          color: notifier.getblck,
                        ),
                      ),
                    )
                  ],
                )
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget marketrate(txt, rate, updown) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          txt,
          style: TextStyle(
              color: notifier.getgrey,
              fontFamily: 'Gilroy_Medium',
              fontSize: 13),
        ),
        SizedBox(height: height / 200),
        Text(
          rate,
          style: TextStyle(
              color: notifier.getblck, fontFamily: 'Gilroy_Bold', fontSize: 14),
        ),
        SizedBox(height: height / 200),
        Text(
          updown,
          style: TextStyle(
              color: const Color(0xff22C36B),
              fontFamily: 'Gilroy_Bold',
              fontSize: 14),
        )
      ],
    );
  }
}

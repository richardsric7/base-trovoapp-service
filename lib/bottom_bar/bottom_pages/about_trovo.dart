import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/storage/state.dart';

class AboutTrovoView extends StatelessWidget {
  const AboutTrovoView({super.key});

  @override
  Widget build(BuildContext context) {
    var notifier = Provider.of<ColorNotifier>(context, listen: true);
    var appState = Provider.of<DataProvider>(context, listen: true);

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      appBar: CustomAppBar(
        context,
        notifier.getwihitecolor,
        '',
        notifier.getblck,
        height: 50,
      ).getBar(),
      body: Padding(
        padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.center,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Image.asset('assets/images/trovo_app.png', width: 250),
                SizedBox(height: 30.0),
                Text(
                  'V ${appState.appVersion}',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.normal,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: 10.0),
                SizedBox(
                  width: 300,
                  child: Text(
                    'The Trovo App is a  decentralized, non-custodial platform for asset tokenization. It is a primary gateway for tokenized securities and features multi-signature wallets, multi-sub-wallets, and a range of advanced features for the secure storage and exchange of digital assets.',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.normal,
                      fontFamily: fontbody,
                    ),
                  ),
                ),
                SizedBox(height: 20.0),
                Text(
                  '© 2025 Trovotech Limited',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.normal,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: 10.0),
                Text(
                  'All rights reserved',
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.normal,
                    fontFamily: fontbody,
                  ),
                ),
                SizedBox(height: 2),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

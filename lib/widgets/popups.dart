import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import '../Custom_BlocObserver/notifire_clor.dart';
import '../utils/enstring.dart';

late ColorNotifier notifier;

popup(context,
    {required String title,
    Color bodyColor = Colors.red,
    required String message}) async {
  notifier = Provider.of<ColorNotifier>(context, listen: false);
  height = MediaQuery.of(context).size.height;
  width = MediaQuery.of(context).size.width;
  return showDialog(
      context: context,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return AlertDialog(
            scrollable: true,
            backgroundColor: Colors.transparent,
            insetPadding: const EdgeInsets.all(20),
            content: Container(
              decoration: BoxDecoration(
                color: notifier.getwihitecolor,
                borderRadius: BorderRadius.all(
                  Radius.circular(23),
                ),
              ),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                mainAxisSize: MainAxisSize.min,
                children: [
                  Padding(
                    padding: const EdgeInsets.all(20.0),
                    child: Center(
                      child: Text(
                        title,
                        style: TextStyle(
                            color: notifier.getblck,
                            fontSize: 18,
                            fontWeight: FontWeight.bold,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  Container(
                    constraints: BoxConstraints(
                      maxHeight: height / 5,
                    ),
                    // height: height / 5,
                    child: SingleChildScrollView(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Padding(
                            padding: const EdgeInsets.symmetric(
                                vertical: 10.0, horizontal: 5.0),
                            child: Text(
                              message,
                              style: TextStyle(
                                fontSize: 17,
                                fontWeight: FontWeight.w300,
                                color: bodyColor,
                              ),
                              textAlign: TextAlign.center,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.all(10.0),
                    child: ElevatedButton(
                      onPressed: () =>
                          Navigator.of(context).pop(), // dismiss dialog,
                      style: ButtonStyle(
                        fixedSize: MaterialStateProperty.all(
                          Size(width / 1.5, height / 20),
                        ),
                        backgroundColor: MaterialStateProperty.all<Color>(
                            notifier.getbluecolor),
                        shape:
                            MaterialStateProperty.all<RoundedRectangleBorder>(
                          const RoundedRectangleBorder(
                            borderRadius: BorderRadius.all(
                              Radius.circular(10),
                            ),
                          ),
                        ),
                      ),
                      child: Text(
                        LanguageEn.continuee,
                        style: TextStyle(
                            color: notifier.getwihitecolor,
                            fontFamily: fontbody),
                      ),
                    ),
                  ),
                  SizedBox(height: height / 50),
                ],
              ),
            ));
      });
}

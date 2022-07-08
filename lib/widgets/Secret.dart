import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

import '../Custom_BlocObserver/fonts.dart';

class Secret extends StatefulWidget {
  late String alias;
  late String secret;
  late String publicKey;
  // late bool? showSecret; // this will be controlled by external code
  Secret(
    this.alias,
    this.secret,
    this.publicKey, {
    Key? key,
  }) : super(key: key);

  @override
  State<Secret> createState() => _SecretState();
}

class _SecretState extends State<Secret> {
  bool show = false; // internal state for this widget
  String hiddenText = '**********';

  @override
  Widget build(BuildContext context) {
    return Card(
      shadowColor: notifier.getblck,
      color: notifier.getwihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      elevation: 5,
      child: ListTile(
        title: Padding(
          padding: const EdgeInsets.symmetric(vertical: 20.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                'Alias: ',
                style: TextStyle(
                  fontFamily: fontsemibold,
                  fontWeight: FontWeight.bold,
                  color: notifier.getblck,
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    widget.alias,
                    style: TextStyle(
                      fontFamily: fontbody,
                      color: notifier.getblck,
                    ),
                  ),
                  IconButton(
                    onPressed: () => {
                      Clipboard.setData(
                        ClipboardData(text: widget.alias),
                      ),
                      showSnackBar('Alias', context),
                    },
                    icon: Icon(Icons.copy),
                    color: notifier.getblck,
                  ),
                ],
              ),
              Divider(),
              Text(
                'Public Key: ',
                style: TextStyle(
                  fontFamily: fontsemibold,
                  fontWeight: FontWeight.bold,
                  color: notifier.getblck,
                ),
              ),
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    width: 200,
                    child: Text(
                      widget.publicKey,
                      style: TextStyle(
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                  IconButton(
                    onPressed: () => {
                      Clipboard.setData(
                        ClipboardData(text: widget.publicKey),
                      ),
                      showSnackBar('Public Key', context),
                    },
                    icon: Icon(Icons.copy),
                    color: notifier.getblck,
                  ),
                ],
              ),
              Divider(),
              Text(
                'Secret Key: ',
                style: TextStyle(
                  fontFamily: fontsemibold,
                  fontWeight: FontWeight.bold,
                  color: notifier.getblck,
                ),
              ),
              SizedBox(height: height / 50),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    width: 200,
                    child: Text(
                      getSecretText(),
                      style: TextStyle(
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                  IconButton(
                    onPressed: () {
                      setState(() {
                        show = !show;
                      });
                    },
                    icon: Icon(
                        show ? CupertinoIcons.eye_slash : CupertinoIcons.eye),
                    color: notifier.getblck,
                  ),
                  IconButton(
                    onPressed: () => {
                      Clipboard.setData(
                        ClipboardData(text: widget.secret),
                      ),
                      showSnackBar('Secret', context),
                    },
                    icon: Icon(Icons.copy),
                    color: notifier.getblck,
                  ),
                ],
              ),
              ElevatedButton(
                onPressed: () => {
                  Clipboard.setData(
                    ClipboardData(
                        text:
                            'Alias:  ${widget.alias}\n\nPublic Key:  ${widget.publicKey}\n\nSecretKey:  ${widget.secret}'),
                  ),
                  showSnackBar('Wallet Details', context),
                },
                style: ButtonStyle(
                  backgroundColor:
                      MaterialStateProperty.all<Color>(notifier.getbluecolor!),
                ),
                child: Text(
                  'Copy All',
                  style: TextStyle(
                    fontFamily: fontsemibold,
                  ),
                ),
              )
            ],
          ),
        ),
        // subtitle: Padding(
        //   padding: const EdgeInsets.symmetric(vertical: 8.0),
        //   child: Text(
        //     getSecretText(),
        //     style: TextStyle(
        //       fontFamily: fontbody,
        //       color: notifier.getblck,
        //     ),
        //   ),
        // ),
        // trailing: Row(
        //   mainAxisSize: MainAxisSize.min,
        //   children: [
        //     IconButton(
        //       onPressed: () => {
        //         Clipboard.setData(
        //           ClipboardData(text: widget.secret),
        //         ),
        //         showSnackBar('Secret', context),
        //       },
        //       icon: Icon(Icons.copy),
        //       color: notifier.getblck,
        //     ),
        //     IconButton(
        //       onPressed: () {
        //         setState(() {
        //           show = !show;
        //         });
        //       },
        //       icon: Icon(show ? CupertinoIcons.eye_slash : CupertinoIcons.eye),
        //       color: notifier.getblck,
        //     ),
        //   ],
        // ),
      ),
    );
  }

  String getSecretText() => show ? widget.secret : hiddenText;
}

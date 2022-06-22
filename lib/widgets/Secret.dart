import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

import '../Custom_BlocObserver/fonts.dart';

class Secret extends StatefulWidget {
  late String alias;
  late String secret;
  // late bool? showSecret; // this will be controlled by external code
  Secret(
    this.alias,
    this.secret, {
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
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: Text(
            widget.alias,
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getblck,
            ),
          ),
        ),
        subtitle: Padding(
          padding: const EdgeInsets.symmetric(vertical: 8.0),
          child: Text(
            getSecretText(),
            style: TextStyle(
              fontFamily: fontbody,
              color: notifier.getblck,
            ),
          ),
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
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
            IconButton(
              onPressed: () {
                setState(() {
                  show = !show;
                });
              },
              icon: Icon(show ? CupertinoIcons.eye_slash : CupertinoIcons.eye),
              color: notifier.getblck,
            ),
          ],
        ),
      ),
    );
  }

  String getSecretText() {
    // print('widget: ${widget.showSecret} local: ${show}');
    // if ((widget.showSecret != null && widget.showSecret!) || show) {
    //   return widget.secret;
    // }
    // return hiddenText;
    return show ? widget.secret : hiddenText;
  }
}

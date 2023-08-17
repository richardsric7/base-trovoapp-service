import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

import '../custom_bloc_observer/fonts.dart';

class Secret extends StatefulWidget {
  late String alias;
  late String secret;
  late String publicKey;
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
  late ColorNotifier notifier;

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: false);

    return Card(
      shadowColor: darkblck,
      color: notifier.gettilewihitecolor,
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
                  Container(
                    width: 200,
                    child: Text(
                      widget.alias,
                      style: TextStyle(
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
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
      ),
    );
  }

  String getSecretText() => show ? widget.secret : hiddenText;
}

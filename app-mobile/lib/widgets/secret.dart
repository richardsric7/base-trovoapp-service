import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:google_sign_in/google_sign_in.dart';
import 'package:googleapis/drive/v3.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/functions/google_drive_client.dart';
import 'package:trovo_app/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_app/widgets/utilities.dart';

import '../custom_bloc_observer/fonts.dart';

class Secret extends StatefulWidget {
  late final String alias;
  late final String secret;
  late final String address;
  Secret(this.alias, this.secret, this.address, {Key? key}) : super(key: key);

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
                      Clipboard.setData(ClipboardData(text: widget.alias)),
                      showSnackBar('Alias', context),
                    },
                    icon: Icon(Icons.copy),
                    color: notifier.getblck,
                  ),
                ],
              ),
              Divider(),
              Text(
                'Address: ',
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
                      widget.address,
                      style: TextStyle(
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                  IconButton(
                    onPressed: () => {
                      Clipboard.setData(ClipboardData(text: widget.address)),
                      showSnackBar('Address', context),
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
                      show ? CupertinoIcons.eye_slash : CupertinoIcons.eye,
                    ),
                    color: notifier.getblck,
                  ),
                  IconButton(
                    onPressed: () => {
                      Clipboard.setData(ClipboardData(text: widget.secret)),
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
                          'Alias:  ${widget.alias}\n\nAddress:  ${widget.address}\n\nSecretKey:  ${widget.secret}',
                    ),
                  ),
                  showSnackBar('Wallet Details', context),
                },
                style: ButtonStyle(
                  backgroundColor: WidgetStateProperty.all<Color>(
                    notifier.getbluecolor!,
                  ),
                  foregroundColor: WidgetStateProperty.all<Color>(
                    notifier.getwihitecolor,
                  ),
                ),
                child: Text(
                  'Copy All',
                  style: TextStyle(
                    fontFamily: fontsemibold,
                    color: wihitecolor,
                  ),
                ),
              ),
              ElevatedButton(
                onPressed: () async {
                  final googleSignIn = GoogleSignIn.standard(
                    scopes: [DriveApi.driveFileScope],
                  );
                  var account = await googleSignIn.signIn();
                  final GoogleSignInAuthentication? auth =
                      await account?.authentication;
                  final String accessToken = auth!.accessToken!;
                  var client = await GoogleDriveClient.create(
                    googleSignIn.currentUser!,
                    accessToken,
                  );

                  var fileContent = await client.downloadFile();

                  if (fileContent == null ||
                      !fileContent.contains(
                        "${widget.alias}|${widget.secret}|${widget.address}",
                      )) {
                    client.uploadFile(
                      '${fileContent ?? ''}\n${widget.alias}|${widget.secret}|${widget.address}',
                    );
                  } else {}

                  googleSignIn.signOut();
                  showSnackBarForInfo("backupsuccess".tr(), context);
                },
                style: ButtonStyle(
                  backgroundColor: WidgetStateProperty.all<Color>(
                    notifier.getbluecolor80!,
                  ),
                  foregroundColor: WidgetStateProperty.all<Color>(
                    notifier.getwihitecolor,
                  ),
                ),
                child: Text(
                  "backupongoogledrive".tr(),
                  style: TextStyle(
                    fontFamily: fontsemibold,
                    color: wihitecolor,
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String getSecretText() => show ? widget.secret : hiddenText;
}

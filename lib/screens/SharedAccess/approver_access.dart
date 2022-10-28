import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';

import '../../Models/Wallet.dart';

class ApproverAccess extends StatefulWidget {
  const ApproverAccess({Key? key}) : super(key: key);

  @override
  State<ApproverAccess> createState() => _ApproverAccessState();
}

class _ApproverAccessState extends State<ApproverAccess> {
  late ColorNotifier notifier;
  late DataProvider appState;
  TextEditingController approversController = TextEditingController();
  TextEditingController initiatorsController = TextEditingController();
  List<Wallet>? wallets;
  Wallet? activeWallet;
  dynamic selectedWallet = '';

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
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    wallets = appState.userInfo!.wallets!;
    activeWallet = appState.activeWallet;
    selectedWallet = activeWallet!.publicKey;
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          'Approver Access',
          notifier.getbluewhitecolor,
          height: height / 15,
        ),
        body: Column(
          children: [
            Center(
              child: Text('Approvers'),
            ),
          ],
        ),
      ),
    );
  }
}

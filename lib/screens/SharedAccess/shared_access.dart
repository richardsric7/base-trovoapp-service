import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/utils/local_auth.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class SharedAccess extends StatefulWidget {
  const SharedAccess({Key? key}) : super(key: key);

  @override
  State<SharedAccess> createState() => _SharedAccessState();
}

class _SharedAccessState extends State<SharedAccess>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late TabController _tabController;
  TextEditingController viewersController = TextEditingController();
  TextEditingController approversController = TextEditingController();
  TextEditingController initiatorsController = TextEditingController();
  final formKey = GlobalKey<FormState>();
  List<Wallet>? wallets;
  Wallet? activeWallet;
  dynamic selectedWallet = '';
  dynamic selectedAccessType = 'Viewer';
  List<String> accessTypes = ['Viewer', 'Approver'];
  List<String> filter = ['All', 'Viewer', 'Initiator', 'Approver'];
  List<String> accessMode = [
    'Access granted by me',
    'Access granted to me'
  ]; // 'mode' for want for a better name
  dynamic selectedAccessMode = 'Access granted to me';
  bool addApprovers = false;

  final Authenticator _authenticator = Authenticator();
  var password = '';
  var viewers = <String>[];
  var initiators = <String>[]; // holds usernames of initiators
  var approvers = <String>[]; // holds usernames of approvers
  String noOfApprovalsNeeded = '';

  List<DropdownMenuItem<String>> get walletDropdownItems {
    return wallets!
        .map<DropdownMenuItem<String>>((wallet) => DropdownMenuItem(
            child: Text(
              wallet.alias!,
              overflow: TextOverflow.ellipsis,
            ),
            value: wallet.publicKey))
        .toList();
  }

  List<DropdownMenuItem<String>> get accessTypeDropdownItems {
    return accessTypes
        .map<DropdownMenuItem<String>>((accessType) => DropdownMenuItem(
            child: Text(
              accessType,
              overflow: TextOverflow.ellipsis,
            ),
            value: accessType))
        .toList();
  }

  List<DropdownMenuItem<String>> get sortDropdownItems {
    return filter
        .map<DropdownMenuItem<String>>((item) => DropdownMenuItem(
            child: Text(
              item,
              overflow: TextOverflow.ellipsis,
            ),
            value: item))
        .toList();
  }

  List<DropdownMenuItem<String>> get accessModeDropdownItems {
    return accessMode
        .map<DropdownMenuItem<String>>((item) => DropdownMenuItem(
            child: Text(
              item,
              overflow: TextOverflow.ellipsis,
            ),
            value: item))
        .toList();
  }

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
    _tabController = TabController(length: 3, vsync: this);
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
          LanguageEn.sharedaccess,
          notifier.getbluewhitecolor,
          height: height / 15,
        ),
        body: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
              child: TabBar(
                controller: _tabController,
                labelColor: notifier.getbluewhitecolor,
                indicatorColor: notifier.getbluewhitecolor,
                labelStyle: TextStyle(
                  fontSize: 15.sp,
                  fontWeight: FontWeight.w600,
                  fontFamily: fontsemibold,
                ),
                tabs: [
                  Tab(
                    height: 50,
                    text: LanguageEn.grantaccess,
                  ),
                  Tab(
                    height: 50,
                    text: LanguageEn.pendingapprovals,
                  ),
                  Tab(
                    height: 50,
                    text: LanguageEn.accesslist,
                  ),
                ],
              ),
            ),
            Container(
              height: height / 1.22,
              child: TabBarView(controller: _tabController, children: [
                SingleChildScrollView(
                  child: grantAccess(),
                ),
                SingleChildScrollView(
                  child: pendingApprovals(),
                ),
                SingleChildScrollView(
                  child: accessList(),
                ),
              ]),
            ),
          ],
        ),
      ),
    );
  }

  Widget pendingApprovals() {
    return Container(
      height: height / 3,
      child: Center(
        child: Text(
          LanguageEn.nopendingapprovals,
          overflow: TextOverflow.visible,
          textAlign: TextAlign.center,
          style: TextStyle(
            fontSize: 15,
            fontFamily: fontsemibold,
            color: notifier.getbluewhitecolor,
          ),
        ),
      ),
    );
  }

  Widget accessList() {
    return Column(
      children: [
        SizedBox(height: height / 30),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 15.0),
          child: Row(
            children: [
              Text(
                'Mode',
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                    fontSize: 15.sp),
              ),
              SizedBox(
                width: width / 10,
              ),
              Expanded(
                child: DropdownButtonFormField(
                  isExpanded: true,
                  dropdownColor: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                  decoration: InputDecoration(
                    contentPadding:
                        EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                    enabledBorder: OutlineInputBorder(
                      borderSide: BorderSide.none,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    border: OutlineInputBorder(
                      borderSide: BorderSide.none,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    filled: true,
                    fillColor: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  value: selectedAccessMode,
                  icon: Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                  elevation: 0,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 15.sp,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500),
                  onChanged: (newValue) {
                    setState(() {
                      selectedAccessMode = newValue!;
                    });
                  },
                  items: accessModeDropdownItems,
                ),
              ),
            ],
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 15.0),
          child: Row(
            children: [
              Text(
                LanguageEn.filterby,
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontbody,
                    fontSize: 15.sp),
              ),
              SizedBox(
                width: width / 10,
              ),
              Expanded(
                child: DropdownButtonFormField(
                  isExpanded: true,
                  dropdownColor: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                  decoration: InputDecoration(
                    contentPadding:
                        EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                    enabledBorder: OutlineInputBorder(
                      borderSide: BorderSide.none,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    border: OutlineInputBorder(
                      borderSide: BorderSide.none,
                      borderRadius: BorderRadius.circular(20),
                    ),
                    filled: true,
                    fillColor: notifier.isDark
                        ? darktilewhitecolor
                        : notifier.getaddsubwalletgrey,
                  ),
                  value: selectedAccessType,
                  icon: Icon(
                    Icons.keyboard_arrow_down_rounded,
                    color: notifier.getbluewhitecolor,
                  ),
                  elevation: 0,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: 15.sp,
                      fontFamily: fontsemibold,
                      fontWeight: FontWeight.w500),
                  onChanged: (newValue) {
                    setState(() {
                      selectedAccessType = newValue!;
                      appState.activeWallet = wallets!
                          .firstWhere((wallet) => wallet.publicKey == newValue);
                    });
                  },
                  items: sortDropdownItems,
                ),
              ),
            ],
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
        for (var i = 0; i < appState.sharedWallets.length; i++) ...[
          tiles(
            walletOwner: appState.sharedWallets[i]['owner']!,
            walletAlias: appState.sharedWallets[i]['walletAlias']!,
            accessType: appState.sharedWallets[i]['permission']!,
          ),
        ],
        SizedBox(
          height: height / 10,
        ),
        Padding(
            padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom)),
      ],
    );
  }

  Widget tiles(
      {required String walletOwner,
      required String walletAlias,
      required String accessType}) {
    return Card(
      elevation: notifier.isDark ? 0 : 5,
      shadowColor: Colors.black,
      color: notifier.gettilewihitecolor,
      margin: EdgeInsets.symmetric(vertical: 10, horizontal: 20),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(15.0),
      ),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 8.0),
        child: ListTile(
          title: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    walletOwner,
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontsemibold,
                      color: notifier.getblck,
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      accessType,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      walletAlias,
                      style: TextStyle(
                        fontSize: 9,
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                ],
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  ElevatedButton(
                    onPressed: () => {},
                    style: ButtonStyle(
                      backgroundColor: MaterialStateProperty.all<Color>(
                          notifier.getbluecolor!),
                    ),
                    child: Text(
                      'Approve Revoke Access',
                      style: TextStyle(
                        fontFamily: fontsemibold,
                        fontSize: 10.sp,
                      ),
                    ),
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      '3 out of 4 initiated',
                      style: TextStyle(
                        fontSize: 9,
                        fontFamily: fontbody,
                        color: notifier.getblck,
                      ),
                    ),
                  ),
                ],
              )
            ],
          ),
        ),
      ),
    );
  }

  Widget grantAccess() {
    return Column(
      children: [
        SizedBox(height: height / 30),
        chooseWallet(),
        SizedBox(
          height: height / 30,
        ),
        if (viewers.length > 0) ...[
          Padding(
            padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
            child: Container(
              decoration: BoxDecoration(
                borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                color: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 20.0, vertical: 15.0),
                    child: Column(
                      children: [
                        Container(
                            width: width / 1.3,
                            child: Wrap(
                              alignment: WrapAlignment.center,
                              children: [
                                for (var i = 0; i < viewers.length; i++) ...[
                                  userItem(viewers[i], () {
                                    viewers.removeAt(i);
                                  }, notifier.getbluecolor)
                                ],
                              ],
                            )),
                        SizedBox(height: 2),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
        SizedBox(
          height: height / 50,
        ),
        Container(
          width: width / 1.1,
          child: Text(
            LanguageEn.enteraccountsusernameviewers,
            // : LanguageEn.enteraccountsusernameapprovers,
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody,
                fontSize: 15.sp),
          ),
        ),
        SizedBox(
          height: height / 70,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              CustomTextFormField.textFieldWithoutIcon(
                'Viewer',
                notifier.getbluecolor,
                notifier.getgrey,
                notifier.getprefixicon,
                notifier.getblck,
                notifier.getgrey,
                60.sp,
                210.sp,
                controller: viewersController,
              ),
              SizedBox(
                width: width / 20,
              ),
              ElevatedButton(
                onPressed: () => setState(() {
                  if (viewersController.text.isNotEmpty) {
                    viewers.add(viewersController.text);
                    viewersController.text = '';
                  }
                }),
                style: ButtonStyle(
                  backgroundColor:
                      MaterialStateProperty.all<Color>(notifier.getbluecolor!),
                ),
                child: Text(
                  LanguageEn.add,
                  style: TextStyle(
                    fontFamily: fontsemibold,
                  ),
                ),
              ),
            ],
          ),
        ),
        SizedBox(
          height: height / 70,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Transform.scale(
                scale: 1.sp,
                child: Checkbox(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(
                      Radius.circular(5.sp),
                    ),
                  ),
                  activeColor: notifier.getbluecolor,
                  side: BorderSide(color: notifier.getbluewhitecolor),
                  value: addApprovers,
                  onChanged: (value) {
                    setState(() {
                      addApprovers = value ?? false;
                    });
                  },
                ),
              ),
              Container(
                child: Text(
                  LanguageEn.doyouwanttoaddapprovers,
                  overflow: TextOverflow.visible,
                  style: TextStyle(
                    fontSize: 15,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                  ),
                ),
              ),
            ],
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
        if (addApprovers) ...[
          Container(
            width: width / 1.1,
            child: Text(
              LanguageEn.enternoofapprover,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 15.sp),
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                CustomTextFormField.textFieldWithoutIcon(
                  'Total no. of approvers',
                  notifier.getbluecolor,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  notifier.getgrey,
                  70.sp,
                  300.sp,
                  keyboardtype: TextInputType.number,
                  onChanged: (value) =>
                      noOfApprovalsNeeded = value.trim().replaceAll(' ', ''),
                ),
              ],
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Container(
            width: width / 1.1,
            child: Text(
              LanguageEn.enternoofapprovals,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 15.sp),
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                CustomTextFormField.textFieldWithoutIcon(
                  'No. of required approvals',
                  notifier.getbluecolor,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  notifier.getgrey,
                  70.sp,
                  300.sp,
                  keyboardtype: TextInputType.number,
                  onChanged: (value) =>
                      noOfApprovalsNeeded = value.trim().replaceAll(' ', ''),
                ),
              ],
            ),
          ),
          if (approvers.length > 0) ...[
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                  color: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 20.0, vertical: 15.0),
                      child: Column(
                        children: [
                          Container(
                              width: width / 1.3,
                              child: Wrap(
                                alignment: WrapAlignment.center,
                                children: [
                                  for (var i = 0;
                                      i < approvers.length;
                                      i++) ...[
                                    userItem(approvers[i], () {
                                      approvers.removeAt(i);
                                    }, notifier.getgreencolor)
                                  ],
                                ],
                              )),
                          SizedBox(height: 2),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
          SizedBox(
            height: height / 50,
          ),
          Container(
            width: width / 1.1,
            child: Text(
              LanguageEn.enteraccountsusernameapprovers,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 15.sp),
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                CustomTextFormField.textFieldWithoutIcon(
                  'Approver',
                  notifier.getbluecolor,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  notifier.getgrey,
                  60.sp,
                  210.sp,
                  controller: approversController,
                ),
                SizedBox(
                  width: width / 20,
                ),
                ElevatedButton(
                  onPressed: () => setState(() {
                    if (approversController.text.isNotEmpty) {
                      approvers.add(approversController.text);
                      approversController.text = '';
                    }
                  }),
                  style: ButtonStyle(
                    backgroundColor: MaterialStateProperty.all<Color>(
                        notifier.getbluecolor!),
                  ),
                  child: Text(
                    LanguageEn.add,
                    style: TextStyle(
                      fontFamily: fontsemibold,
                    ),
                  ),
                ),
              ],
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          if (initiators.length > 0) ...[
            Padding(
              padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
              child: Container(
                decoration: BoxDecoration(
                  borderRadius: const BorderRadius.all(Radius.circular(15.0)),
                  color: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Padding(
                      padding: const EdgeInsets.symmetric(
                          horizontal: 20.0, vertical: 15.0),
                      child: Column(
                        children: [
                          Container(
                              width: width / 1.3,
                              child: Wrap(
                                alignment: WrapAlignment.center,
                                children: [
                                  for (var i = 0;
                                      i < initiators.length;
                                      i++) ...[
                                    userItem(initiators[i], () {
                                      initiators.removeAt(i);
                                    }, notifier.getbluecolor80)
                                  ],
                                ],
                              )),
                          SizedBox(height: 2),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
          SizedBox(
            height: height / 70,
          ),
          Container(
            width: width / 1.1,
            child: Text(
              LanguageEn.enteraccountsusernameinitiators,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 15.sp),
            ),
          ),
          SizedBox(
            height: height / 70,
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 20.0),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                CustomTextFormField.textFieldWithoutIcon(
                  'Initiator',
                  notifier.getbluecolor,
                  notifier.getgrey,
                  notifier.getprefixicon,
                  notifier.getblck,
                  notifier.getgrey,
                  60.sp,
                  210.sp,
                  controller: initiatorsController,
                ),
                SizedBox(
                  width: width / 20,
                ),
                ElevatedButton(
                  onPressed: () => setState(() {
                    if (initiatorsController.text.isNotEmpty) {
                      initiators.add(initiatorsController.text);
                      initiatorsController.text = '';
                    }
                  }),
                  style: ButtonStyle(
                    backgroundColor: MaterialStateProperty.all<Color>(
                        notifier.getbluecolor!),
                  ),
                  child: Text(
                    LanguageEn.add,
                    style: TextStyle(
                      fontFamily: fontsemibold,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
        SizedBox(
          height: height / 15,
        ),
        Button(
          LanguageEn.proceed,
          notifier.getbluecolor,
          wihitecolor,
          onTap: () {
            appState.viewData = {
              AddSharedAccessDetailsViewPageConfig.key: {
                'viewers': viewers,
                'approvers': approvers,
                'accessType': selectedAccessType,
                'noOfApprovalsNeeded': noOfApprovalsNeeded,
                'initiators': initiators,
                'addApprovers': addApprovers,
              }
            };
            appState.currentAction = PageAction(
              state: PageState.addPage,
              page: AddSharedAccessDetailsViewPageConfig,
            );
          },
        ),
        SizedBox(
          height: height / 10,
        ),
        Padding(
            padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom)),
      ],
    );
  }

  Widget userItem(String name, void Function() onRemove, Color color) {
    return Padding(
      padding: const EdgeInsets.all(3.0),
      child: Container(
        decoration: BoxDecoration(
            borderRadius: const BorderRadius.all(Radius.circular(10.0)),
            color: color),
        child: Padding(
          padding: const EdgeInsets.all(5.0),
          child: Wrap(
            children: [
              Text(
                name,
                overflow: TextOverflow.ellipsis,
                softWrap: true,
                style: TextStyle(
                    color: wihitecolor, fontFamily: fontbody, fontSize: 15.sp),
              ),
              SizedBox(
                width: width / 70,
              ),
              GestureDetector(
                  onTap: () => setState(() {
                        onRemove();
                      }),
                  child: Icon(
                    Icons.cancel_outlined,
                    color: wihitecolor,
                  ))
            ],
          ),
        ),
      ),
    );
  }

  Widget chooseWallet() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 15.0),
      child: Row(
        children: [
          Text(
            LanguageEn.choosewallet,
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody,
                fontSize: 15.sp),
          ),
          SizedBox(
            width: width / 10,
          ),
          Expanded(
            child: DropdownButtonFormField(
              isExpanded: true,
              dropdownColor: notifier.isDark
                  ? darktilewhitecolor
                  : notifier.getaddsubwalletgrey,
              decoration: InputDecoration(
                contentPadding:
                    EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                enabledBorder: OutlineInputBorder(
                  borderSide: BorderSide.none,
                  borderRadius: BorderRadius.circular(20),
                ),
                border: OutlineInputBorder(
                  borderSide: BorderSide.none,
                  borderRadius: BorderRadius.circular(20),
                ),
                filled: true,
                fillColor: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
              ),
              value: selectedWallet,
              icon: Icon(
                Icons.keyboard_arrow_down_rounded,
                color: notifier.getbluewhitecolor,
              ),
              elevation: 0,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontSize: 15.sp,
                  fontFamily: fontsemibold,
                  fontWeight: FontWeight.w500),
              onChanged: (newValue) {
                setState(() {
                  selectedWallet = newValue!;
                  appState.activeWallet = wallets!
                      .firstWhere((wallet) => wallet.publicKey == newValue);
                });
              },
              items: walletDropdownItems,
            ),
          ),
        ],
      ),
    );
  }

  void handleAuthorization() {
    if (!formKey.currentState!.validate()) {
      return;
    }

    if (password == appState.password!) {
      // sendDataToServer();
    } else {
      popup(context,
          title: LanguageEn.oops, message: LanguageEn.invalidpassword);
    }
  }

  void toggleSwitch() async {
    try {
      bool result = await _authenticator.authenticateMe();
      if (result) {
        // sendDataToServer();
        // aparently we need the code below to make the
        // screen updata to show loader
        // after authorizing with biometrics
        setState(() {});
      }
    } on PlatformException catch (e) {
      if (e.code == auth_error.notEnrolled ||
          e.code == auth_error.notAvailable) {
        biometricsErrorAlert(context);
      }
    }
  }

  String? validatePassword(String? value) {
    if (value!.isEmpty) return 'Enter your password';

    if (value.length < 6) return 'Use 6 characters or more for your password';

    return null;
  }
}

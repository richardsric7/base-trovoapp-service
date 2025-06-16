import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/models/permission.dart';
import 'package:trovo_app/models/wallet.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/router/page_actions.dart';
import 'package:trovo_app/router/ui_pages.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';

class UpdateSharedAccess extends StatefulWidget {
  const UpdateSharedAccess({Key? key}) : super(key: key);

  @override
  State<UpdateSharedAccess> createState() => _UpdateSharedAccessState();
}

class _UpdateSharedAccessState extends State<UpdateSharedAccess>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
  late Wallet wallet;
  var viewers = <Permission>[];
  var initiators = <Permission>[]; // holds usernames of initiators
  var approvers = <Permission>[]; // holds usernames of approvers
  int noOfApprovalsNeeded = 2;
  int noOfApprovers = 3;
  TextEditingController viewersController = TextEditingController();
  TextEditingController approversController = TextEditingController();
  TextEditingController initiatorsController = TextEditingController();
  final approversFormKey = GlobalKey<FormState>();
  String viewerUsernameErrorMessage = "";
  String approverUsernameErrorMessage = "";
  String initiatorUsernameErrorMessage = "";
  bool addApprovers = false;

  late TabController _tabController;

  List<DropdownMenuItem<int>> get getNoOfApproversDropdownItems {
    var items = <DropdownMenuItem<int>>[];
    for (var i = 3; i < 20; i++) {
      items.add(
        DropdownMenuItem(
          child: Text(i.toString(), overflow: TextOverflow.ellipsis),
          value: i,
        ),
      );
    }
    return items;
  }

  List<DropdownMenuItem<int>> get getNoOfApprovalsDropdownItems {
    var items = <DropdownMenuItem<int>>[];
    for (var i = 2; i < noOfApprovers; i++) {
      items.add(
        DropdownMenuItem(
          child: Text(i.toString(), overflow: TextOverflow.ellipsis),
          value: i,
        ),
      );
    }
    return items;
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
    _tabController = TabController(length: 3, vsync: this);
    getdarkmodepreviousstate();
    appState = Provider.of<DataProvider>(context, listen: false);
    wallet = appState.userInfo!.getWallet(
      appState.viewData!['walletPublicKey'],
    );
    noOfApprovalsNeeded = wallet.numberOfApprovalsNeeded!;

    // since you can only pass around objects by reference in dart
    // and since we need to modify permissions without necessarily
    // modifying the original user object until it is sent to the
    // server and committed, we have to clone the permissions object
    // and use it for the necessary modifications without touching
    // the main data
    for (var permission in wallet.permissions!) {
      if (permission.permission == 'VIEW-ONLY') {
        viewers.add(Permission.clone(permission));
      }

      if (permission.permission == 'APPROVER') {
        approvers.add(Permission.clone(permission));
      }

      if (permission.permission == 'INITIATOR') {
        initiators.add(Permission.clone(permission));
      }
    }

    if (approvers.length > 0 && noOfApprovalsNeeded > 0) {
      noOfApprovers = approvers.length;
      addApprovers = true;
    }
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "updatesharedaccess".tr(),
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        body: SingleChildScrollView(
          child: Container(
            width: width,
            child: Column(
              children: [
                if (wallet.isPrimaryWallet || !addApprovers) ...[
                  SingleChildScrollView(
                    child: Column(
                      children: [
                        SizedBox(height: height / 20),
                        Text(
                          "viewaccess".tr(),
                          textAlign: TextAlign.center,
                          style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontsemibold,
                            fontSize: 18.sp,
                          ),
                        ),
                        showViewers(),
                        SizedBox(height: height / 25),
                        Button(
                          "proceed".tr(),
                          notifier.getbluecolor,
                          wihitecolor,
                          onTap: () {
                            if (isValid()) {
                              updateSharedAccess();
                            }
                          },
                        ),
                        Padding(
                          padding: EdgeInsets.only(
                            bottom: MediaQuery.of(context).viewInsets.bottom,
                          ),
                        ),
                      ],
                    ),
                  ),
                ] else ...[
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
                        Tab(height: 50, text: "viewers".tr()),
                        Tab(height: 50, text: "approvers".tr()),
                        Tab(height: 50, text: "initiators".tr()),
                      ],
                    ),
                  ),
                  Container(
                    height: height / 1.45,
                    child: TabBarView(
                      controller: _tabController,
                      children: [
                        SingleChildScrollView(child: showViewers()),
                        SingleChildScrollView(child: showApprovers()),
                        SingleChildScrollView(child: showInitiators()),
                      ],
                    ),
                  ),
                  SizedBox(height: height / 50),
                  Button(
                    "proceed".tr(),
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: () {
                      if (_tabController.index == 2) {
                        if (isValid()) {
                          updateSharedAccess();
                        }
                      } else {
                        _tabController.animateTo(_tabController.index + 1);
                      }
                    },
                  ),
                ],
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget showViewers() {
    return Column(
      children: [
        SizedBox(height: height / 50),
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
                    horizontal: 20.0,
                    vertical: 15.0,
                  ),
                  child: Column(
                    children: [
                      Container(
                        width: width / 1.8,
                        child: Wrap(
                          alignment: WrapAlignment.center,
                          children: [
                            if (viewers.length > 0) ...[
                              for (var i = 0; i < viewers.length; i++) ...[
                                getPermissionItem(viewers, i, 'viewer'),
                              ],
                            ] else ...[
                              Text(
                                "nameofviewersappearhere".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody,
                                  fontSize: 15.sp,
                                ),
                              ),
                            ],
                          ],
                        ),
                      ),
                      SizedBox(height: 2),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        SizedBox(height: height / 30),
        Container(
          width: width / 1.1,
          child: Text(
            "enteraccountsusernameviewers".tr(),
            textAlign: TextAlign.center,
            style: TextStyle(
              color: notifier.getbluewhitecolor,
              fontFamily: fontbody,
              fontSize: 15.sp,
            ),
          ),
        ),
        SizedBox(height: height / 50),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: CustomTextFormField.textFieldWithoutIcon(
            "viewer".tr(),
            notifier.getbluecolor,
            notifier.getgrey,
            notifier.getprefixicon,
            notifier.getblck,
            notifier.getgrey,
            60.sp,
            210.sp,
            controller: viewersController,
          ),
        ),
        // since we cannot use form validators here to check username
        // because of the async process we use this to show error messages
        if (viewerUsernameErrorMessage.isNotEmpty) ...[
          Text(
            viewerUsernameErrorMessage,
            textAlign: TextAlign.center,
            style: TextStyle(
              color: Colors.red,
              fontFamily: fontbody,
              fontSize: 11.sp,
            ),
          ),
        ],
        SizedBox(height: height / 50),
        ElevatedButton(
          onPressed: () async {
            viewerUsernameErrorMessage = '';
            var username = viewersController.text.trim();
            setState(() {});

            if (username.isEmpty) {
              viewerUsernameErrorMessage = "enterusername".tr();
              setState(() {});
              return;
            }

            if (appState.userInfo!.username == username) {
              // viewerUsernameErrorMessage =
              //     'You cannot add yourself as a viewer on this wallet because as the owner of this wallet you already have view access';
              // setState(() {});
              popup(
                context,
                title: "error".tr(),
                message: "cannotaddyourself".tr(),
                bodyColor: Colors.red,
              );
              return;
            }

            if (viewers
                .where((viewer) => viewer.targetUsername == username)
                .isNotEmpty) {
              viewerUsernameErrorMessage = "usernamealreadyadded".tr();
              setState(() {});
              return;
            }

            // if the username is already on the approvers list then there's an error
            if (approvers
                .where((approver) => approver.targetUsername == username)
                .isNotEmpty) {
              popup(
                context,
                title: "alert".tr(),
                message: "usernamealreadyaddedtoapproverslist".tr(),
                bodyColor: notifier.getbluecolor,
              );
              return;
            }

            // if the username is already on the initiators list then there's an error
            if (initiators
                .where((initiator) => initiator.targetUsername == username)
                .isNotEmpty) {
              popup(
                context,
                title: "error".tr(),
                message: "usernamealreadyaddedtoinitiatorslist".tr(),
                bodyColor: notifier.getbluecolor,
              );
              return;
            }

            var userInfo = await checkUsername(username);
            if (userInfo == null) {
              viewerUsernameErrorMessage = "notavalidtrovousername".tr();
              setState(() {});
              return;
            }

            viewersController.text = '';
            viewers.add(
              Permission(
                targetUsername: username,
                fullName: '${userInfo['firstName']} ${userInfo['lastName']}',
                permission: 'VIEW-ONLY',
                permissionState: PermissionState.Added,
              ),
            );
            setState(() {});
          },
          style: ButtonStyle(
            backgroundColor: WidgetStateProperty.all<Color>(
              notifier.getbluecolor!,
            ),
            foregroundColor: WidgetStateProperty.all<Color>(
              notifier.getwihitecolor,
            ),
          ),
          child: Text("add".tr(), style: TextStyle(fontFamily: fontsemibold)),
        ),
        SizedBox(height: height / 50),
        // if wallet is not primary wallet
        // primary wallets can only have view-only shared access
        // the cannot have approver and initiator shared access
        if (!wallet.isInitiator && !wallet.isPrimaryWallet) ...[
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Transform.scale(
                scale: 1.sp,
                child: Checkbox(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(Radius.circular(5.sp)),
                  ),
                  activeColor: notifier.getbluecolor,
                  side: BorderSide(color: notifier.getbluewhitecolor),
                  value: addApprovers,
                  onChanged: (value) {
                    setState(() {
                      if (value == false && approvers.length > 0) {
                        showResponseMessage(
                          context,
                          "allexistingapproverswillberemoved".tr(),
                          () {
                            approvers.forEach((approver) {
                              approver.permissionState =
                                  PermissionState.Revoked;
                            });

                            initiators.forEach((initiator) {
                              initiator.permissionState =
                                  PermissionState.Revoked;
                            });
                            addApprovers = false;
                            setState(() {});
                          },
                        );
                      } else {
                        addApprovers = value ?? false;
                        if (addApprovers) {
                          approvers.forEach((approver) {
                            approver.permissionState = null;
                          });

                          initiators.forEach((initiator) {
                            initiator.permissionState = null;
                          });
                          _tabController.animateTo(1);
                        }
                      }
                    });
                  },
                ),
              ),
              Container(
                child: Text(
                  "doyouwanttoaddapprovers".tr(),
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
        ],
        SizedBox(height: height / 10),
      ],
    );
  }

  Widget showApprovers() {
    return Form(
      key: approversFormKey,
      child: Column(
        children: [
          SizedBox(height: height / 50),
          SizedBox(height: height / 70),
          Container(
            height: height / 10,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Container(
                    height: height / 5,
                    width: width / 3.6,
                    child: Column(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 10.0),
                          child: Text(
                            "approvals".tr(),
                            style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                              fontSize: 15.sp,
                            ),
                          ),
                        ),
                        SizedBox(height: height / 90),
                        Expanded(
                          child: DropdownButtonFormField(
                            isExpanded: true,
                            dropdownColor: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                            decoration: InputDecoration(
                              contentPadding: EdgeInsets.symmetric(
                                vertical: 0,
                                horizontal: 20,
                              ),
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
                            value: noOfApprovalsNeeded,
                            icon: Icon(
                              Icons.keyboard_arrow_down_rounded,
                              color: notifier.getbluewhitecolor,
                            ),
                            elevation: 0,
                            style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontsemibold,
                              fontWeight: FontWeight.w500,
                            ),
                            onChanged: (newValue) {
                              setState(() {
                                noOfApprovalsNeeded = int.parse(
                                  newValue.toString(),
                                );
                              });
                            },
                            items: getNoOfApprovalsDropdownItems,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(width: width / 20),
                  Container(
                    height: height / 5,
                    width: width / 20,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'of'.tr(),
                          style: TextStyle(
                            color: notifier.getbluewhitecolor,
                            fontFamily: fontbody,
                            fontSize: 15.sp,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(width: width / 20),
                  Container(
                    height: height / 5,
                    width: width / 3.6,
                    child: Column(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 10.0),
                          child: Text(
                            "approvers".tr(),
                            style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                              fontSize: 15.sp,
                            ),
                          ),
                        ),
                        SizedBox(height: height / 90),
                        Expanded(
                          child: DropdownButtonFormField(
                            isExpanded: true,
                            dropdownColor: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                            decoration: InputDecoration(
                              contentPadding: EdgeInsets.symmetric(
                                vertical: 0,
                                horizontal: 20,
                              ),
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
                            value: noOfApprovers,
                            icon: Icon(
                              Icons.keyboard_arrow_down_rounded,
                              color: notifier.getbluewhitecolor,
                            ),
                            elevation: 0,
                            style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontSize: 15.sp,
                              fontFamily: fontsemibold,
                              fontWeight: FontWeight.w500,
                            ),
                            onChanged: (newValue) {
                              setState(() {
                                var newValueInt = int.parse(
                                  newValue.toString(),
                                );
                                if (noOfApprovalsNeeded > newValueInt) {
                                  noOfApprovalsNeeded = newValueInt - 1;
                                }
                                noOfApprovers = newValueInt;
                              });
                            },
                            items: getNoOfApproversDropdownItems,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 10.0),
            child: Text(
              "noapprovalsrequiredoutofno".tr(
                args: [
                  noOfApprovalsNeeded.toString(),
                  noOfApprovers.toString(),
                ],
              ),
              textAlign: TextAlign.center,
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody,
                fontSize: 13.sp,
              ),
            ),
          ),
          SizedBox(height: height / 30),
          Container(
            child: Text(
              "enteraccountsusernameapprovers".tr(),
              textAlign: TextAlign.center,
              style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody,
                fontSize: 15.sp,
              ),
            ),
          ),
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
                      horizontal: 20.0,
                      vertical: 15.0,
                    ),
                    child: Column(
                      children: [
                        Container(
                          width: width / 1.8,
                          child: Wrap(
                            alignment: WrapAlignment.center,
                            children: [
                              if (approvers.length > 0) ...[
                                for (var i = 0; i < approvers.length; i++) ...[
                                  getPermissionItem(approvers, i, 'approver'),
                                ],
                              ] else ...[
                                Text(
                                  "noapproversyet".tr(),
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody,
                                    fontSize: 15.sp,
                                  ),
                                ),
                              ],
                            ],
                          ),
                        ),
                        SizedBox(height: 2),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
          SizedBox(height: height / 50),
          CustomTextFormField.textFieldWithoutIcon(
            "approver".tr(),
            notifier.getbluecolor,
            notifier.getgrey,
            notifier.getprefixicon,
            notifier.getblck,
            notifier.getgrey,
            60.sp,
            210.sp,
            controller: approversController,
          ),
          // since we cannot use form validators here because of the async process
          // to check username we use this to show error messages
          if (approverUsernameErrorMessage.isNotEmpty) ...[
            Text(
              approverUsernameErrorMessage,
              textAlign: TextAlign.center,
              style: TextStyle(
                color: Colors.red,
                fontFamily: fontbody,
                fontSize: 11.sp,
              ),
            ),
          ],
          SizedBox(height: height / 50),
          ElevatedButton(
            onPressed: () async {
              approverUsernameErrorMessage = '';
              var username = approversController.text.trim();
              setState(() {});

              if (username.isEmpty) {
                approverUsernameErrorMessage = "enterusername".tr();
                setState(() {});
                return;
              }

              if (approvers
                  .where((approver) => approver.targetUsername == username)
                  .isNotEmpty) {
                approverUsernameErrorMessage = "usernamealreadyadded".tr();
                setState(() {});
                return;
              }

              if (approvers
                      .where(
                        (approver) =>
                            approver.permissionState != PermissionState.Revoked,
                      )
                      .length ==
                  noOfApprovers) {
                popup(
                  context,
                  title: "error".tr(),
                  message: "usernamecannotbemorethannoapprovers".tr(),
                );
                return;
              }

              // if the username is already on the viewers list then there's
              // no need to check again that the username is valid so we add it to
              // to the approvers list
              var tempViewer = viewers.where(
                (viewer) =>
                    viewer.targetUsername == username &&
                    viewer.permissionState != PermissionState.Revoked,
              );
              if (tempViewer.isNotEmpty) {
                showResponseMessage(
                  context,
                  "willrevokeviewonlyaccess".tr(),
                  () {
                    approvers.add(
                      Permission(
                        targetUsername: username,
                        fullName: '${tempViewer.first.fullName}',
                        permission: 'APPROVER',
                        permissionState: PermissionState.Added,
                      ),
                    );
                    approversController.text = '';
                    if (tempViewer.first.permissionState ==
                        PermissionState.Added) {
                      viewers.removeWhere(
                        (permission) => permission.targetUsername == username,
                      );
                    } else {
                      tempViewer.first.permissionState =
                          PermissionState.Revoked;
                    }
                    setState(() {});
                  },
                );
                return;
              }

              // if the username is already on the approvers list then there's
              // no need to check again that the username is valid so we add it to
              // to the approvers list
              var tempInitiator = initiators.where(
                (viewer) => viewer.targetUsername == username,
              );
              if (tempInitiator.isNotEmpty) {
                approvers.add(
                  Permission(
                    targetUsername: username,
                    fullName: '${tempInitiator.first.fullName}',
                    permission: 'APPROVER',
                    permissionState: PermissionState.Added,
                  ),
                );
                approversController.text = '';
                setState(() {});
                return;
              }

              var userInfo = await checkUsername(username);
              if (userInfo == null) {
                approverUsernameErrorMessage = "notavalidtrovousername".tr();
                setState(() {});
                return;
              }

              approvers.add(
                Permission(
                  targetUsername: username,
                  fullName: '${userInfo['firstName']} ${userInfo['lastName']}',
                  permission: 'APPROVER',
                  permissionState: PermissionState.Added,
                ),
              );
              approversController.text = '';
              setState(() {});
            },
            style: ButtonStyle(
              backgroundColor: WidgetStateProperty.all<Color>(
                notifier.getbluecolor!,
              ),
              foregroundColor: WidgetStateProperty.all<Color>(
                notifier.getwihitecolor,
              ),
            ),
            child: Text("add".tr(), style: TextStyle(fontFamily: fontsemibold)),
          ),
          SizedBox(height: height / 20),
          Padding(
            padding: EdgeInsets.only(
              bottom: MediaQuery.of(context).viewInsets.bottom,
            ),
          ),
        ],
      ),
    );
  }

  Widget showInitiators() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        SizedBox(height: height / 50),
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
                    horizontal: 20.0,
                    vertical: 15.0,
                  ),
                  child: Column(
                    children: [
                      Container(
                        width: width / 1.8,
                        child: Wrap(
                          alignment: WrapAlignment.center,
                          children: [
                            if (initiators.length > 0) ...[
                              for (var i = 0; i < initiators.length; i++) ...[
                                getPermissionItem(initiators, i, 'initiator'),
                              ],
                            ] else ...[
                              Text(
                                "noinitatorsyet".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody,
                                  fontSize: 15.sp,
                                ),
                              ),
                            ],
                          ],
                        ),
                      ),
                      SizedBox(height: 2),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
        SizedBox(height: height / 50),
        Container(
          child: Text(
            "enteraccountsusernameinitiators".tr(),
            textAlign: TextAlign.center,
            style: TextStyle(
              color: notifier.getbluewhitecolor,
              fontFamily: fontbody,
              fontSize: 15.sp,
            ),
          ),
        ),
        SizedBox(height: height / 50),
        CustomTextFormField.textFieldWithoutIcon(
          "initiator".tr(),
          notifier.getbluecolor,
          notifier.getgrey,
          notifier.getprefixicon,
          notifier.getblck,
          notifier.getgrey,
          60.sp,
          210.sp,
          controller: initiatorsController,
        ),
        // since we cannot use form validators here because of the async process
        // to check username we use this to show error messages
        if (initiatorUsernameErrorMessage.isNotEmpty) ...[
          Text(
            initiatorUsernameErrorMessage,
            textAlign: TextAlign.center,
            style: TextStyle(
              color: Colors.red,
              fontFamily: fontbody,
              fontSize: 11.sp,
            ),
          ),
        ],
        SizedBox(height: height / 50),
        ElevatedButton(
          onPressed: () async {
            initiatorUsernameErrorMessage = '';
            var username = initiatorsController.text.trim();
            setState(() {});

            if (username.isEmpty) {
              initiatorUsernameErrorMessage = "enterusername".tr();
              setState(() {});
              return;
            }

            if (initiators
                .where((initiator) => initiator.targetUsername == username)
                .isNotEmpty) {
              initiatorUsernameErrorMessage = "usernamealreadyadded".tr();
              setState(() {});
              return;
            }

            // if the username is already on the viewers list then there's
            // no need to check again that the username is valid so we add it to
            // to the approvers list
            var tempViewer = viewers.where(
              (viewer) =>
                  viewer.targetUsername == username &&
                  viewer.permissionState != PermissionState.Revoked,
            );
            if (tempViewer.isNotEmpty) {
              showResponseMessage(context, "willrevokeviewonlyaccess".tr(), () {
                initiators.add(
                  Permission(
                    targetUsername: username,
                    fullName: '${tempViewer.first.fullName}',
                    permission: 'INITIATOR',
                    permissionState: PermissionState.Added,
                  ),
                );
                initiatorsController.text = '';
                if (tempViewer.first.permissionState == PermissionState.Added) {
                  viewers.removeWhere(
                    (permission) => permission.targetUsername == username,
                  );
                } else {
                  tempViewer.first.permissionState = PermissionState.Revoked;
                }

                setState(() {});
              });
              return;
            }

            // if the username is already on the approvers list then there's
            // no need to check again that the username is valid so we add it to
            // to the approvers list
            var tempApprover = approvers.where(
              (approver) => approver.targetUsername == username,
            );
            if (tempApprover.isNotEmpty) {
              initiators.add(
                Permission(
                  targetUsername: username,
                  fullName: '${tempApprover.first.fullName}',
                  permission: 'INITIATOR',
                  permissionState: PermissionState.Added,
                ),
              );
              initiatorsController.text = '';
              setState(() {});
              return;
            }

            var userInfo = await checkUsername(username);
            if (userInfo == null) {
              initiatorUsernameErrorMessage = "notavalidtrovousername".tr();
              setState(() {});
              return;
            }
            initiators.add(
              Permission(
                targetUsername: username,
                fullName: '${userInfo['firstName']} ${userInfo['lastName']}',
                permission: 'INITIATOR',
                permissionState: PermissionState.Added,
              ),
            );
            initiatorsController.text = '';
            setState(() {});
          },
          style: ButtonStyle(
            backgroundColor: WidgetStateProperty.all<Color>(
              notifier.getbluecolor!,
            ),
            foregroundColor: WidgetStateProperty.all<Color>(
              notifier.getwihitecolor,
            ),
          ),
          child: Text("add".tr(), style: TextStyle(fontFamily: fontsemibold)),
        ),
        SizedBox(height: height / 20),
        Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom,
          ),
        ),
      ],
    );
  }

  void updateSharedAccess() {
    appState.viewData = {
      'walletPublicKey': wallet.publicKey,
      'walletAlias': wallet.alias,
      'viewers': viewers,
      'addApprovers': addApprovers,
      'approvers': approvers,
      'noOfApprovers': addApprovers ? noOfApprovers : 0,
      'noOfApprovalsNeeded': addApprovers ? noOfApprovalsNeeded : 0,
      'initiators': initiators,
    };

    appState.currentAction = PageAction(
      state: PageState.addPage,
      page: UpdateSharedAccessDetailsViewPageConfig,
    );
    // popup(context, title: 'e pass', message: 'This one don pass');
  }

  // we need to check that the username entered here is a valid
  // username of an active trovo account
  Future<Map?> checkUsername(String username) async {
    try {
      showLoader(context);
      Map responseData = await makeGetRequest(
        uri: '/v1/users/$username',
        signer: appState.activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.activeWallet!.publicKey!,
      );

      print('response: ${responseData}');
      hideLoader(context);

      if (responseData['statusCode'] == 200) {
        return responseData['data']['userData'];
      }

      return null;
    } catch (e) {
      popup(context, title: "error".tr(), message: e.toString());
      hideLoader(context);
      return null;
    }
  }

  Widget getPermissionItem(List<Permission> permList, int index, String rel) {
    switch (permList[index].permissionState) {
      case PermissionState.Revoked:
        return userItem(
          '${permList[index].targetUsername} [${permList[index].fullName}]',
          () {
            permList[index].permissionState = null;
            setState(() {});
          },
          backColor: Colors.red,
          foreColor: wihitecolor,
          restoreMode: true,
        );
      case PermissionState.Added:
        return userItem(
          '${permList[index].targetUsername} [${permList[index].fullName}]',
          () {
            permList.removeAt(index);
            setState(() {});
          },
          backColor: notifier.getgreencolor,
          foreColor: wihitecolor,
        );

      default:
        return userItem(
          '${permList[index].targetUsername} [${permList[index].fullName}]',
          () {
            showResponseMessage(
              context,
              "abouttorevokeaccess".tr(args: [rel]),
              () {
                permList[index].permissionState = PermissionState.Revoked;
                setState(() {});
              },
            );
          },
          backColor: notifier.getbluebackcolor,
          foreColor: wihitecolor,
        );
    }
  }

  bool isValid() {
    var approversList = approvers.where(
      (permission) => permission.permissionState != PermissionState.Revoked,
    );
    var initiatorsList = initiators.where(
      (permission) => permission.permissionState != PermissionState.Revoked,
    );

    if (addApprovers) {
      if (approversList.isEmpty) {
        popup(context, title: "error".tr(), message: "addapprovers".tr());
        _tabController.animateTo(1);
        return false;
      }

      if (approversList.length < noOfApprovers ||
          approversList.length > noOfApprovers) {
        popup(
          context,
          title: "error".tr(),
          message: "approverscannotbelessnoofapprover".tr(),
        );
        _tabController.animateTo(1);
        return false;
      }

      if (initiatorsList.isEmpty) {
        popup(
          context,
          title: "error".tr(),
          message: "musthaveatleastoneinitiator".tr(),
        );
        _tabController.animateTo(2);
        return false;
      }

      if (approversList.isEmpty) {
        popup(
          context,
          title: "error".tr(),
          message: "cannothaveinitiatorswithoutapprovers".tr(),
        );
        _tabController.animateTo(1);
        return false;
      }
    }

    return true;
  }

  @override
  void dispose() {
    super.dispose();
    print('disposing...');
    appState.viewData![UpdateSharedAccessViewPageConfig.key] = null;
    print('disposed');
  }
}

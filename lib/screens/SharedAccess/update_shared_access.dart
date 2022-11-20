import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/button/custtom_button.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Permission.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

class UpdateSharedAccess extends StatefulWidget {
  const UpdateSharedAccess({Key? key}) : super(key: key);

  @override
  State<UpdateSharedAccess> createState() => _UpdateSharedAccessState();
}

class _UpdateSharedAccessState extends State<UpdateSharedAccess>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late DataProvider appState;
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
  var viewData;
  bool addApprovers = false;

  late TabController _tabController;

  List<DropdownMenuItem<int>> get getNoOfApproversDropdownItems {
    var items = <DropdownMenuItem<int>>[];
    for (var i = 1; i < 20; i++) {
      items.add(DropdownMenuItem(
          child: Text(
            i.toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: i));
    }
    return items;
  }

  List<DropdownMenuItem<int>> get getNoOfApprovalsDropdownItems {
    var items = <DropdownMenuItem<int>>[];
    for (var i = 1; i < noOfApprovers; i++) {
      items.add(DropdownMenuItem(
          child: Text(
            i.toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: i));
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
    viewData = appState.viewData![UpdateSharedAccessViewPageConfig.key];
    for (var i = 0; i < viewData['viewers'].length; i++) {
      viewers.add(Permission(
        targetUsername: viewData['viewers'][i].targetUsername,
        fullName: viewData['viewers'][i].fullName,
        permission: viewData['viewers'][i].permission,
      ));
    }

    for (var i = 0; i < viewData['approvers'].length; i++) {
      approvers.add(Permission(
        targetUsername: viewData['approvers'][i].targetUsername,
        fullName: viewData['approvers'][i].fullName,
        permission: viewData['approvers'][i].permission,
      ));
    }

    for (var i = 0; i < viewData['initiators'].length; i++) {
      initiators.add(Permission(
        targetUsername: viewData['initiators'][i].targetUsername,
        fullName: viewData['initiators'][i].fullName,
        permission: viewData['initiators'][i].permission,
      ));
    }

    if (approvers.length > 0 && viewData['numberOfApprovalsNeeded'] > 0) {
      noOfApprovers = approvers.length;
      noOfApprovalsNeeded = viewData['numberOfApprovalsNeeded'];
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
          'Update Shared Access',
          notifier.getbluewhitecolor,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Container(
            width: width,
            child: Column(
              children: [
                if (viewData['isPrimaryWallet'] == 1 || !addApprovers) ...[
                  SingleChildScrollView(
                    child: Column(
                      children: [
                        SizedBox(
                          height: height / 20,
                        ),
                        Text(
                          'View Access',
                          textAlign: TextAlign.center,
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontsemibold,
                              fontSize: 18.sp),
                        ),
                        showViewers(),
                        SizedBox(height: height / 25),
                        Button(
                          'Proceed',
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
                              bottom: MediaQuery.of(context).viewInsets.bottom),
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
                        Tab(
                          height: 50,
                          text: 'Viewers',
                        ),
                        Tab(
                          height: 50,
                          text: 'Approvers',
                        ),
                        Tab(
                          height: 50,
                          text: 'Initiators',
                        ),
                      ],
                    ),
                  ),
                  Container(
                    height: height / 1.45,
                    child: TabBarView(
                      controller: _tabController,
                      children: [
                        SingleChildScrollView(
                          child: showViewers(),
                        ),
                        SingleChildScrollView(
                          child: showApprovers(),
                        ),
                        SingleChildScrollView(
                          child: showInitiators(),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: height / 50),
                  Button(
                    'Proceed',
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
        SizedBox(
          height: height / 50,
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
                      horizontal: 20.0, vertical: 15.0),
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
                                  'Name of viewers appear here',
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                      fontSize: 15.sp),
                                ),
                              ]
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
        SizedBox(
          height: height / 30,
        ),
        Container(
          width: width / 1.1,
          child: Text(
            LanguageEn.enteraccountsusernameviewers,
            textAlign: TextAlign.center,
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody,
                fontSize: 15.sp),
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20.0),
          child: CustomTextFormField.textFieldWithoutIcon(
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
        ),
        // since we cannot use form validators here to check username
        // because of the async process we use this to show error messages
        if (viewerUsernameErrorMessage.isNotEmpty) ...[
          Text(
            viewerUsernameErrorMessage,
            textAlign: TextAlign.center,
            style: TextStyle(
                color: Colors.red, fontFamily: fontbody, fontSize: 11.sp),
          ),
        ],
        SizedBox(
          height: height / 50,
        ),
        ElevatedButton(
          onPressed: () async {
            viewerUsernameErrorMessage = '';
            var username = viewersController.text.trim();
            setState(() {});

            if (username.isEmpty) {
              viewerUsernameErrorMessage = 'Please enter a username';
              setState(() {});
              return;
            }

            if (appState.userInfo!.username == username) {
              // viewerUsernameErrorMessage =
              //     'You cannot add yourself as a viewer on this wallet because as the owner of this wallet you already have view access';
              // setState(() {});
              popup(
                context,
                title: 'Error!',
                message:
                    'You cannot add yourself as a viewer on this wallet because as the owner of this wallet you already have view access.',
                bodyColor: Colors.red,
              );
              return;
            }

            if (viewers
                .where((viewer) => viewer.targetUsername == username)
                .isNotEmpty) {
              viewerUsernameErrorMessage = 'Username already added';
              setState(() {});
              return;
            }

            // if the username is already on the approvers list then there's an error
            if (approvers
                .where((approver) => approver.targetUsername == username)
                .isNotEmpty) {
              popup(
                context,
                title: 'Alert',
                message:
                    'This user is already added to approver access which gives them implicit view access. Please remove them from approver access if you want to grant them view-only access.',
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
                title: 'Error',
                message:
                    'This user is already added to initiator access which gives them implicit view access. Please remove them from initiator access if you want to grant them view-only access.',
                bodyColor: notifier.getbluecolor,
              );
              return;
            }

            var userInfo = await checkUsername(username);
            if (userInfo == null) {
              viewerUsernameErrorMessage = 'This is not a valid Trovo username';
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
        SizedBox(
          height: height / 50,
        ),
        // if wallet is not primary wallet
        // primary wallets can only have view-only shared access
        // the cannot have approver and initiator shared access
        if (viewData['isPrimaryWallet'] == 0) ...[
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
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
                        if (value == false && approvers.length > 0) {
                          showResponseMessage(context,
                              'All existing approvers and initiator access will be revoked. Do you want to proceed?',
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
                          });
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
                    }),
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
        ],
        SizedBox(
          height: height / 10,
        ),
      ],
    );
  }

  Widget showApprovers() {
    return Form(
      key: approversFormKey,
      child: Column(
        children: [
          SizedBox(
            height: height / 50,
          ),
          SizedBox(
            height: height / 70,
          ),
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
                            'Approvals',
                            style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                                fontSize: 15.sp),
                          ),
                        ),
                        SizedBox(
                          height: height / 90,
                        ),
                        Expanded(
                          child: DropdownButtonFormField(
                            isExpanded: true,
                            dropdownColor: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                            decoration: InputDecoration(
                              contentPadding: EdgeInsets.symmetric(
                                  vertical: 0, horizontal: 20),
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
                                fontWeight: FontWeight.w500),
                            onChanged: (newValue) {
                              setState(() {
                                noOfApprovalsNeeded =
                                    int.parse(newValue.toString());
                              });
                            },
                            items: getNoOfApprovalsDropdownItems,
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    width: width / 20,
                  ),
                  Container(
                    height: height / 5,
                    width: width / 20,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Text(
                          'of',
                          style: TextStyle(
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                              fontSize: 15.sp),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(
                    width: width / 20,
                  ),
                  Container(
                    height: height / 5,
                    width: width / 3.6,
                    child: Column(
                      children: [
                        Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 10.0),
                          child: Text(
                            'Approvers',
                            style: TextStyle(
                                color: notifier.getbluewhitecolor,
                                fontFamily: fontsemibold,
                                fontSize: 15.sp),
                          ),
                        ),
                        SizedBox(
                          height: height / 90,
                        ),
                        Expanded(
                          child: DropdownButtonFormField(
                            isExpanded: true,
                            dropdownColor: notifier.isDark
                                ? darktilewhitecolor
                                : notifier.getaddsubwalletgrey,
                            decoration: InputDecoration(
                              contentPadding: EdgeInsets.symmetric(
                                  vertical: 0, horizontal: 20),
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
                                fontWeight: FontWeight.w500),
                            onChanged: (newValue) {
                              setState(() {
                                var newValueInt =
                                    int.parse(newValue.toString());
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
              '${noOfApprovalsNeeded} approvals required out of ${noOfApprovers} approvers',
              textAlign: TextAlign.center,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 13.sp),
            ),
          ),
          SizedBox(
            height: height / 30,
          ),
          Container(
            child: Text(
              LanguageEn.enteraccountsusernameapprovers,
              textAlign: TextAlign.center,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 15.sp),
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
                        horizontal: 20.0, vertical: 15.0),
                    child: Column(
                      children: [
                        Container(
                            width: width / 1.8,
                            child: Wrap(
                              alignment: WrapAlignment.center,
                              children: [
                                if (approvers.length > 0) ...[
                                  for (var i = 0;
                                      i < approvers.length;
                                      i++) ...[
                                    getPermissionItem(approvers, i, 'approver'),
                                  ],
                                ] else ...[
                                  Text(
                                    'You have added no approvers yet',
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                        color: notifier.getbluewhitecolor,
                                        fontFamily: fontbody,
                                        fontSize: 15.sp),
                                  ),
                                ]
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
          SizedBox(
            height: height / 50,
          ),
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
          // since we cannot use form validators here because of the async process
          // to check username we use this to show error messages
          if (approverUsernameErrorMessage.isNotEmpty) ...[
            Text(
              approverUsernameErrorMessage,
              textAlign: TextAlign.center,
              style: TextStyle(
                  color: Colors.red, fontFamily: fontbody, fontSize: 11.sp),
            ),
          ],
          SizedBox(
            height: height / 50,
          ),
          ElevatedButton(
            onPressed: () async {
              approverUsernameErrorMessage = '';
              var username = approversController.text.trim();
              setState(() {});

              if (username.isEmpty) {
                approverUsernameErrorMessage = 'Please enter a username';
                setState(() {});
                return;
              }

              if (approvers
                  .where((approver) => approver.targetUsername == username)
                  .isNotEmpty) {
                approverUsernameErrorMessage = 'Username already added';
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
                popup(context,
                    title: 'Error!',
                    message:
                        'Number of usernames cannot be more than the number of approvers you selected');
                return;
              }

              // if the username is already on the viewers list then there's
              // no need to check again that the username is valid so we add it to
              // to the approvers list
              var tempViewer = viewers.where((viewer) =>
                  viewer.targetUsername == username &&
                  viewer.permissionState != PermissionState.Revoked);
              if (tempViewer.isNotEmpty) {
                showResponseMessage(context,
                    'This user\'s view-only access will be revoked since they will have implicit view access as an approver',
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
                        (permission) => permission.targetUsername == username);
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
              var tempInitiator = initiators
                  .where((viewer) => viewer.targetUsername == username);
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
                approverUsernameErrorMessage =
                    'This is not a valid Trovo username';
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
          SizedBox(
            height: height / 20,
          ),
          Padding(
              padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom)),
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
                      horizontal: 20.0, vertical: 15.0),
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
                                  'You have added no initiators yet',
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                      fontSize: 15.sp),
                                ),
                              ]
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
        SizedBox(
          height: height / 50,
        ),
        Container(
          child: Text(
            LanguageEn.enteraccountsusernameinitiators,
            textAlign: TextAlign.center,
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontbody,
                fontSize: 15.sp),
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
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
        // since we cannot use form validators here because of the async process
        // to check username we use this to show error messages
        if (initiatorUsernameErrorMessage.isNotEmpty) ...[
          Text(
            initiatorUsernameErrorMessage,
            textAlign: TextAlign.center,
            style: TextStyle(
                color: Colors.red, fontFamily: fontbody, fontSize: 11.sp),
          ),
        ],
        SizedBox(
          height: height / 50,
        ),
        ElevatedButton(
          onPressed: () async {
            initiatorUsernameErrorMessage = '';
            var username = initiatorsController.text.trim();
            setState(() {});

            if (username.isEmpty) {
              initiatorUsernameErrorMessage = 'Please enter a username';
              setState(() {});
              return;
            }

            if (initiators
                .where((initiator) => initiator.targetUsername == username)
                .isNotEmpty) {
              initiatorUsernameErrorMessage = 'Username already added';
              setState(() {});
              return;
            }

            // if the username is already on the viewers list then there's
            // no need to check again that the username is valid so we add it to
            // to the approvers list
            var tempViewer = viewers.where((viewer) =>
                viewer.targetUsername == username &&
                viewer.permissionState != PermissionState.Revoked);
            if (tempViewer.isNotEmpty) {
              showResponseMessage(context,
                  'This user\'s view-only access will be revoked since they will have implicit view access as an initiator',
                  () {
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
                      (permission) => permission.targetUsername == username);
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
            var tempApprover = approvers
                .where((approver) => approver.targetUsername == username);
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
              initiatorUsernameErrorMessage =
                  'This is not a valid Trovo username';
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
        SizedBox(
          height: height / 20,
        ),
        Padding(
            padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom)),
      ],
    );
  }

  void updateSharedAccess() {
    appState.viewData![UpdateSharedAccessDetailsViewPageConfig.key] = {
      'walletPublicKey': viewData['walletPublicKey'],
      'walletAlias': viewData['walletAlias'],
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
      popup(context, title: LanguageEn.error, message: e.toString());
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
          Colors.red,
          restoreMode: true,
        );
      case PermissionState.Added:
        return userItem(
          '${permList[index].targetUsername} [${permList[index].fullName}]',
          () {
            permList.removeAt(index);
            setState(() {});
          },
          notifier.getgreencolor,
        );

      default:
        return userItem(
          '${permList[index].targetUsername} [${permList[index].fullName}]',
          () {
            showResponseMessage(context,
                'You are about to revoke this user\'s $rel access to this wallet. Do you want to proceed?',
                () {
              permList[index].permissionState = PermissionState.Revoked;
              setState(() {});
            });
          },
          notifier.getbluecolor,
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
        popup(context, title: 'Error!', message: 'Please add approvers.');
        _tabController.animateTo(1);
        return false;
      }

      if (approversList.length < noOfApprovers) {
        popup(context,
            title: 'Error!',
            message:
                'Number of approver usernames cannot be less than the number of approvers you selected. Please add more approvers.');
        _tabController.animateTo(1);
        return false;
      }

      if (initiatorsList.isEmpty) {
        popup(
          context,
          title: 'Error!',
          message:
              'You must have at least one user with initiator access to this wallet.',
        );
        _tabController.animateTo(2);
        return false;
      }

      if (approversList.isEmpty) {
        popup(context,
            title: 'Error!',
            message:
                'You cannot have initiators without having approvers. Please add approvers.');
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

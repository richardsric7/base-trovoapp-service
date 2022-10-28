import 'dart:math';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
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
import 'package:trovo_wallet/widgets/utilities.dart';
import '../../utils/medeiaqury/medeiaqury.dart';
import 'package:local_auth/error_codes.dart' as auth_error;

class SharedAccess extends StatefulWidget {
  const SharedAccess({Key? key}) : super(key: key);

  @override
  State<SharedAccess> createState() => _SharedAccessState();
}

enum GrantSharedAccessView { viewers, approvers, initiators }

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
  int currentStep = 0;

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
        appBar: PreferredSize(
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: notifier.getwihitecolor,
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            title: Text(
              LanguageEn.sharedaccess,
              style: TextStyle(
                  color: notifier.getbluewhitecolor, fontFamily: fontsemibold),
            ),
          ),
          preferredSize: Size.fromHeight(height / 15),
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
                    text: LanguageEn.accesslist,
                  ),
                  Tab(
                    height: 50,
                    text: LanguageEn.pendingapprovals,
                  ),
                  Tab(
                    height: 50,
                    text: LanguageEn.grantaccess,
                  ),
                ],
              ),
            ),
            Container(
              height: height / 1.22,
              child: TabBarView(controller: _tabController, children: [
                SingleChildScrollView(
                  child: accessList(),
                ),
                SingleChildScrollView(
                  child: pendingApprovals(),
                ),
                SingleChildScrollView(
                  child: grantAccess(),
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
    return Container(
      height: height / 1.22,
      child: Scaffold(
        floatingActionButton: FloatingActionButton(
            onPressed: () {
              selectAccessTypePopup(context, appState);
            },
            backgroundColor: notifier.getbluecolor,
            child: Icon(
              Icons.add,
              size: 30.sp,
            )),
        body: SingleChildScrollView(
          child: Column(
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
                            appState.activeWallet = wallets!.firstWhere(
                                (wallet) => wallet.publicKey == newValue);
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
          ),
        ),
      ),
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

  List<Step> getAllSteps() {
    return <Step>[
      Step(
        state: currentStep > 0 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 0,
        title: Text("Choose wallet",
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
                fontSize: 15.sp)),
        content: Column(
          children: [
            chooseWallet(),
            SizedBox(
              height: height / 30,
            )
          ],
        ),
      ),
      Step(
        state: currentStep > 1 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 1,
        title: Text("Grant viewer access",
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
                fontSize: 15.sp)),
        content: Column(
          children: [
            showViewers(),
          ],
        ),
      ),
      Step(
        state: currentStep > 2 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 2,
        title: Text("Grant approver access",
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
                fontSize: 15.sp)),
        content: Column(
          children: [
            showApprovers(),
          ],
        ),
      ),
      Step(
        state: currentStep > 3 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 3,
        title: Text("Grant initiator access",
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
                fontSize: 15.sp)),
        content: Column(
          children: [
            showInitiators(),
          ],
        ),
      ),
    ];
  }

  List<Step> getOneStep() {
    return <Step>[
      Step(
        state: currentStep > 0 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 0,
        title: Text("Choose wallet",
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
                fontSize: 15.sp)),
        content: Column(
          children: [
            chooseWallet(),
            SizedBox(
              height: height / 30,
            )
          ],
        ),
      ),
      Step(
        state: currentStep > 1 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 1,
        title: Text("Grant viewer access",
            style: TextStyle(
                color: notifier.getbluewhitecolor,
                fontFamily: fontsemibold,
                fontSize: 15.sp)),
        content: Column(
          children: [
            showViewers(),
          ],
        ),
      ),
    ];
  }

  Widget grantAccess() {
    return Column(
      children: [
        // SizedBox(height: height / 30),
        // chooseWallet(),
        if (addApprovers) ...[
          Container(
              height: height / 1.219,
              width: width,
              // color: Colors.black,
              child: Stepper(
                key: Key(Random.secure().nextDouble().toString()),
                type: StepperType.vertical,
                currentStep: currentStep,

                // onStepCancel: () => currentStep == 0
                //     ? null
                //     : setState(() {
                //         currentStep -= 1;
                //       }),
                // onStepContinue: () {
                //   bool isLastStep = (currentStep == getAllSteps().length - 1);
                //   if (isLastStep) {
                //     //Do something with this information
                //   } else {
                //     setState(() {
                //       currentStep += 1;
                //     });
                //   }
                // },
                controlsBuilder: (context, _) {
                  return Column(
                    children: [
                      SizedBox(
                        height: height / 50,
                      ),
                      ElevatedButton(
                        onPressed: () {
                          bool isLastStep =
                              (currentStep == getAllSteps().length - 1);
                          if (isLastStep) {
                            submitSharedAccessForm();
                          } else {
                            setState(() {
                              currentStep += 1;
                            });
                          }
                          setState(() {});
                        },
                        style: ButtonStyle(
                          backgroundColor: MaterialStateProperty.all<Color>(
                              notifier.getbluecolor!),
                        ),
                        child: Padding(
                          padding: const EdgeInsets.all(20.0),
                          child: Text(
                            addApprovers
                                ? 'Add approver access'
                                : 'Grant viewer access',
                            style: TextStyle(
                              fontFamily: fontsemibold,
                              fontSize: 14.sp,
                            ),
                          ),
                        ),
                      ),
                    ],
                  );
                },
                onStepTapped: (step) => setState(() {
                  currentStep = step;
                }),
                steps: getAllSteps(),
              )),
        ] else ...[
          Container(
              height: height / 1.219,
              width: width,
              // color: Colors.black,
              child: Stepper(
                type: StepperType.vertical,
                currentStep: currentStep,
                // onStepCancel: () => currentStep == 0
                //     ? null
                //     : setState(() {
                //         currentStep -= 1;
                //       }),
                // onStepContinue: () {
                //   bool isLastStep = (currentStep == getAllSteps().length - 1);
                //   if (isLastStep) {
                //     //Do something with this information
                //   } else {
                //     setState(() {
                //       currentStep += 1;
                //     });
                //   }
                // },
                onStepTapped: (step) => setState(() {
                  currentStep = step;
                }),
                controlsBuilder: (context, _) {
                  return Column(
                    children: [
                      SizedBox(
                        height: height / 50,
                      ),
                      ElevatedButton(
                        onPressed: () {
                          bool isLastStep =
                              (currentStep == getOneStep().length - 1);
                          if (isLastStep) {
                            submitSharedAccessForm();
                          } else {
                            setState(() {
                              currentStep += 1;
                            });
                          }
                          setState(() {});
                        },
                        style: ButtonStyle(
                          backgroundColor: MaterialStateProperty.all<Color>(
                              notifier.getbluecolor!),
                        ),
                        child: Padding(
                          padding: const EdgeInsets.all(20.0),
                          child: Text(
                            addApprovers
                                ? 'Add approver access'
                                : 'Grant viewer access',
                            style: TextStyle(
                              fontFamily: fontsemibold,
                              fontSize: 14.sp,
                            ),
                          ),
                        ),
                      ),
                    ],
                  );
                },
                steps: getOneStep(),
              )),
        ],
        Padding(
            padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom)),
      ],
    );
  }

  void submitSharedAccessForm() {
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
  }

  Widget showViewers() {
    return Column(
      children: [
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
                                  userItem(viewers[i], () {
                                    viewers.removeAt(i);
                                  }, notifier.getbluecolor)
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
        SizedBox(
          height: height / 50,
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
        // SizedBox(
        //   height: height / 10,
        // ),
        // Button(
        //   addApprovers ? 'Grant approver access' : 'Grant viewer access',
        //   notifier.getbluecolor,
        //   wihitecolor,
        //   onTap: () {
        //     appState.viewData = {
        //       AddSharedAccessDetailsViewPageConfig.key: {
        //         'viewers': viewers,
        //         'approvers': approvers,
        //         'accessType': selectedAccessType,
        //         'noOfApprovalsNeeded': noOfApprovalsNeeded,
        //         'initiators': initiators,
        //         'addApprovers': addApprovers,
        //       }
        //     };

        //     if (addApprovers) {
        //       appState.grantSharedAccessView.view =
        //           GrantSharedAccessView.approvers;
        //     } else {
        //       appState.currentAction = PageAction(
        //         state: PageState.addPage,
        //         page: AddSharedAccessDetailsViewPageConfig,
        //       );
        //     }
        //     setState(() {});
        //   },
        // ),
      ],
    );
  }

  Widget showApprovers() {
    return Column(
      children: [
        Container(
          child: Text(
            'Grant approver access',
            overflow: TextOverflow.visible,
            style: TextStyle(
              fontSize: 17,
              fontFamily: fontsemibold,
              color: notifier.getbluewhitecolor,
            ),
          ),
        ),
        SizedBox(
          height: height / 20,
        ),
        Container(
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
        CustomTextFormField.textFieldWithoutIcon(
          'e.g, 10',
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
        SizedBox(
          height: height / 70,
        ),
        Container(
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
        CustomTextFormField.textFieldWithoutIcon(
          'e.g, 5',
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
        SizedBox(
          height: height / 50,
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
                              if (approvers.length > 0) ...[
                                for (var i = 0; i < approvers.length; i++) ...[
                                  userItem(approvers[i], () {
                                    approvers.removeAt(i);
                                  }, notifier.getbluecolor)
                                ],
                              ] else ...[
                                Text(
                                  'Name of approvers appear here',
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
        SizedBox(
          height: height / 50,
        ),
        ElevatedButton(
          onPressed: () => setState(() {
            if (approversController.text.isNotEmpty) {
              approvers.add(approversController.text);
              approversController.text = '';
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
        SizedBox(
          height: height / 20,
        ),
      ],
    );
  }

  Widget showInitiators() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        // Container(
        //   child: Text(
        //     'Grant initiator access',
        //     overflow: TextOverflow.visible,
        //     style: TextStyle(
        //       fontSize: 17,
        //       fontFamily: fontsemibold,
        //       color: notifier.getbluewhitecolor,
        //     ),
        //   ),
        // ),
        SizedBox(
          height: height / 30,
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
                                  userItem(viewers[i], () {
                                    viewers.removeAt(i);
                                  }, notifier.getbluecolor)
                                ],
                              ] else ...[
                                Text(
                                  'Name of initiators appear here',
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
        SizedBox(
          height: height / 50,
        ),
        ElevatedButton(
          onPressed: () => setState(() {
            if (initiatorsController.text.isNotEmpty) {
              initiators.add(initiatorsController.text);
              initiatorsController.text = '';
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
        SizedBox(
          height: height / 10,
        ),
        // Button(
        //   LanguageEn.proceed,
        //   notifier.getbluecolor,
        //   wihitecolor,
        //   onTap: () {
        //     appState.viewData = {
        //       AddSharedAccessDetailsViewPageConfig.key: {
        //         'viewers': viewers,
        //         'approvers': approvers,
        //         'accessType': selectedAccessType,
        //         'noOfApprovalsNeeded': noOfApprovalsNeeded,
        //         'initiators': initiators,
        //         'addApprovers': addApprovers,
        //       }
        //     };
        //     appState.currentAction = PageAction(
        //       state: PageState.addPage,
        //       page: AddSharedAccessDetailsViewPageConfig,
        //     );
        //   },
        // ),
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
    return Container(
      height: height / 10,
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 15.0),
        child: Column(
          children: [
            // Text(
            //   LanguageEn.choosewallet,
            //   style: TextStyle(
            //       color: notifier.getbluewhitecolor,
            //       fontFamily: fontbody,
            //       fontSize: 15.sp),
            // ),
            SizedBox(
              height: height / 50,
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

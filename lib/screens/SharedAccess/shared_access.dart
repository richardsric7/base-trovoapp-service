import 'dart:math';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_wallet/Custom_BlocObserver/colors.dart';
import 'package:trovo_wallet/Custom_BlocObserver/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/Permission.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

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
  final approversFormKey = GlobalKey<FormState>();

  List<Wallet>? wallets;
  Wallet? activeWallet;
  dynamic selectedWallet = '';
  dynamic selectedFilter = 'All';
  List<String> accessTypes = ['Viewer', 'Approver'];
  List<String> filter = ['All', 'Viewer', 'Initiator', 'Approver'];
  List<String> accessMode = [
    'Access granted to me',
    'Access granted by me',
  ]; // 'mode' for want for a better name
  dynamic selectedAccessMode = 'Access granted to me';
  bool addApprovers = false;
  int currentStep = 0;
  String viewerUsernameErrorMessage = "";
  String approverUsernameErrorMessage = "";
  String initiatorUsernameErrorMessage = "";
  var allKey = Key(Random.secure().nextDouble().toString());

  var password = '';
  var viewers = <String>[];
  var initiators = <String>[]; // holds usernames of initiators
  var approvers = <String>[]; // holds usernames of approvers
  var userFullnames = {};
  int noOfApprovalsNeeded = 2;
  int noOfApprovers = 3;

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
    // if (noOfApprovers > 0) {
    for (var i = 1; i < noOfApprovers; i++) {
      items.add(DropdownMenuItem(
          child: Text(
            i.toString(),
            overflow: TextOverflow.ellipsis,
          ),
          value: i));
    }
    // } else {
    //   items.add(DropdownMenuItem(
    //       child: Text(
    //         '0',
    //         overflow: TextOverflow.ellipsis,
    //       ),
    //       value: 0));
    // }
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
        floatingActionButton: FloatingActionButton(
            onPressed: () async {
              shareAccessInfoPopup(context);
            },
            backgroundColor: notifier.getbluecolor,
            child: Icon(
              Icons.question_mark,
              size: 30.sp,
            )),
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
                        value: selectedFilter,
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
                            selectedFilter = newValue!;
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
              if (selectedAccessMode == 'Access granted to me') ...[
                getAccessGrantedToMe(appState),
              ] else ...[
                Builder(builder: (context) {
                  // we rename 'Viewer' to 'VIEW-ONLY' because that's what is
                  // returned from the server.
                  String filter = selectedFilter == 'Viewer'
                      ? 'VIEW-ONLY'
                      : selectedFilter.toString().toUpperCase();
                  // filter the wallets according to the access type that
                  // the user selected
                  var filteredWallets = <Wallet>[];
                  wallets!.forEach((wallet) {
                    wallet.permissions!.forEach((permissionObj) {
                      if (filter == 'ALL' ||
                          permissionObj.permission == filter) {
                        filteredWallets.add(wallet);
                      }
                    });
                  });
                  return Column(
                    children: [
                      if (filteredWallets.length > 0) ...[
                        for (var walletIndex = 0;
                            walletIndex < filteredWallets.length;
                            walletIndex++) ...{
                          if (filteredWallets[walletIndex].permissions !=
                                  null &&
                              filteredWallets[walletIndex].permissions!.length >
                                  0) ...[
                            accessGrantedByMe(
                                walletAlias:
                                    filteredWallets[walletIndex].alias!,
                                permissions:
                                    filteredWallets[walletIndex].permissions),
                          ]
                        }
                      ] else ...[
                        Center(
                          heightFactor: 15.sp,
                          child: Text(
                            'Nothing to show here',
                            style: TextStyle(
                              fontSize: 17,
                              fontFamily: fontbody,
                              color: notifier.getbluecolor,
                            ),
                          ),
                        ),
                      ]
                    ],
                  );
                })
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

  Widget getAccessGrantedToMe(DataProvider appState) {
    // create a map to hold each wallet info
    var wallets = {};
    // store the wallet aliases as keys here so that we can use it to easily get the values back
    // from the wallets map since we cannot create widgets by looping through the map using map.forEach((k,v))
    var walletKeys = [];
    for (var i = 0; i < appState.sharedWallets.length; i++) {
      // filter shared access by permission
      if (appState.sharedWallets[i]['permission'] ==
              selectedFilter.toString().toUpperCase() ||
          selectedFilter == 'All') {
        // if the wallet is not already added to the map then add it
        if (wallets[appState.sharedWallets[i]['walletAlias']] == null) {
          walletKeys.add(appState.sharedWallets[i]['walletAlias']);
          wallets[appState.sharedWallets[i]['walletAlias']] = {
            'walletAlias': appState.sharedWallets[i]['walletAlias'],
            'permissions': <String>[appState.sharedWallets[i]['permission']],
            'owner': appState.sharedWallets[i]['owner'],
            'walletPublicKey': appState.sharedWallets[i]['walletPublicKey'],
            'walletDescription': appState.sharedWallets[i]['walletDescription'],
          };
        } else {
          // if we got here then the wallet is already on the map so we add
          // this permission to the list of permissions granted to the user on the
          // wallet
          wallets[appState.sharedWallets[i]['walletAlias']]['permissions']
              .add(appState.sharedWallets[i]['permission']);
        }
      }
    }

    return Column(
      children: [
        if (walletKeys.length > 0) ...[
          // loop through the walletKeys array and use each key to get the values
          // stored in the wallets map above.
          for (var i = 0; i < walletKeys.length; i++) ...[
            GestureDetector(
              onTap: () {
                // add the shared access data to viewData so we can pass it to
                // shared access details view when user taps on it
                appState.viewData![SharedWalletInfoViewPageConfig.key] =
                    wallets[walletKeys[i]];
                appState.currentAction = PageAction(
                  state: PageState.addPage,
                  page: SharedWalletInfoViewPageConfig,
                );
              },
              child: accessGrantedToMe(
                walletOwner: wallets[walletKeys[i]]['owner'],
                walletAlias: walletKeys[i],
                publicKey: wallets[walletKeys[i]]['walletPublicKey'],
                permissions: wallets[walletKeys[i]]['permissions'],
                walletDescription: wallets[walletKeys[i]]['walletDescription'],
              ),
            )
          ],
        ] else ...[
          Center(
            heightFactor: 15.sp,
            child: Text(
              'Nothing to show here',
              style: TextStyle(
                fontSize: 17,
                fontFamily: fontbody,
                color: notifier.getbluecolor,
              ),
            ),
          ),
        ]
      ],
    );
  }

  Widget accessGrantedToMe(
      {required String walletOwner,
      required String walletAlias,
      required String publicKey,
      required List<String> permissions,
      required walletDescription}) {
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
          title: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    walletAlias,
                    style: TextStyle(
                      fontSize: 17,
                      fontFamily: fontsemibold,
                      color: notifier.getbluecolor,
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      'Owner:',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: width / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      walletOwner,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluecolor,
                      ),
                    ),
                  ),
                ],
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      'Permissions:',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: width / 70,
                  ),
                  for (var i = 0; i < permissions.length; i++) ...[
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        permissions[i].toLowerCase(),
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontbody,
                          color: notifier.getbluecolor,
                        ),
                      ),
                    ),
                    if (i < permissions.length - 1) ...[
                      Padding(
                        padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                        child: Text(
                          ',',
                          style: TextStyle(
                            fontSize: 15,
                            fontFamily: fontbody,
                            color: notifier.getbluecolor,
                          ),
                        ),
                      ),
                    ],
                    SizedBox(
                      width: width / 70,
                    ),
                  ]
                ],
              ),
              SizedBox(
                height: height / 90,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget accessGrantedByMe(
      {required String walletAlias, required List<Permission>? permissions}) {
    // sort the list first
    List<Permission> approversList = [];
    List<Permission> initiatorsList = [];
    List<Permission> viewersList = [];
    if (permissions != null && permissions.length > 0) {
      for (var permIndex = 0; permIndex < permissions.length; permIndex++) {
        if (permissions[permIndex].permission == 'VIEW-ONLY') {
          viewersList.add(permissions[permIndex]);
        }
        if (permissions[permIndex].permission == 'APPROVER') {
          approversList.add(permissions[permIndex]);
        }
        if (permissions[permIndex].permission == 'INITIATOR') {
          initiatorsList.add(permissions[permIndex]);
        }
      }
    }
    return Column(
      children: [
        Card(
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
              title: Column(
                children: [
                  Text(
                    walletAlias,
                    style: TextStyle(
                      fontSize: 19,
                      fontFamily: fontsemibold,
                      color: notifier.getbluecolor,
                    ),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  // filter by viewer
                  if (selectedFilter == 'All' ||
                      selectedFilter == 'Viewer') ...[
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                          child: Text(
                            'Viewer access',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluecolor,
                            ),
                          ),
                        ),
                      ],
                    ),
                    if (viewersList.length > 0) ...[
                      for (var permIndex = 0;
                          permIndex < viewersList.length;
                          permIndex++) ...[
                        Row(children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              '${viewersList[permIndex].targetUsername}',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluecolor,
                              ),
                            ),
                          ),
                        ])
                      ]
                    ] else ...[
                      Row(
                        children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Container(
                              width: width / 1.5,
                              child: Text(
                                'You have granted no view-only access on this wallet',
                                overflow: TextOverflow.visible,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontbody,
                                  color: notifier.getbluecolor,
                                ),
                              ),
                            ),
                          ),
                        ],
                      )
                    ],
                  ],
                  // filter by approver
                  if (selectedFilter == 'All' ||
                      selectedFilter == 'Approver') ...[
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                          child: Text(
                            'Approver access',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluecolor,
                            ),
                          ),
                        ),
                      ],
                    ),
                    if (approversList.length > 0) ...[
                      for (var permIndex = 0;
                          permIndex < approversList.length;
                          permIndex++) ...[
                        Row(children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              '${approversList[permIndex].targetUsername}',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluecolor,
                              ),
                            ),
                          ),
                        ])
                      ]
                    ] else ...[
                      Row(
                        children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Container(
                              width: width / 1.5,
                              child: Text(
                                'You have granted no approver access on this wallet',
                                overflow: TextOverflow.visible,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontbody,
                                  color: notifier.getbluecolor,
                                ),
                              ),
                            ),
                          ),
                        ],
                      )
                    ],
                  ],
                  // filter by initiator
                  if (selectedFilter == 'All' ||
                      selectedFilter == 'Initiator') ...[
                    SizedBox(
                      height: height / 50,
                    ),
                    Row(
                      children: [
                        Padding(
                          padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                          child: Text(
                            'Initiator access',
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluecolor,
                            ),
                          ),
                        ),
                      ],
                    ),
                    if (initiatorsList.length > 0) ...[
                      for (var permIndex = 0;
                          permIndex < initiatorsList.length;
                          permIndex++) ...[
                        Row(children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              '${initiatorsList[permIndex].targetUsername}',
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluecolor,
                              ),
                            ),
                          ),
                        ])
                      ]
                    ] else ...[
                      Row(
                        children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Container(
                              width: width / 1.5,
                              child: Text(
                                'You have granted no initiator access on this wallet',
                                overflow: TextOverflow.visible,
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontbody,
                                  color: notifier.getbluecolor,
                                ),
                              ),
                            ),
                          ),
                        ],
                      )
                    ],
                    SizedBox(
                      height: height / 50,
                    ),
                  ]
                ],
              ),
            ),
          ),
        ),
        SizedBox(
          height: height / 50,
        ),
      ],
    );
  }

  List<Step> getSteps() {
    return <Step>[
      Step(
        state: currentStep > 0 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 0,
        title: Text("Add viewer access",
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
        state: currentStep > 1 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 1,
        title: Text("Add approver access",
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
        state: currentStep > 2 ? StepState.complete : StepState.indexed,
        isActive: currentStep >= 2,
        title: Text("Add initiator access",
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

  Widget grantAccess() {
    return Column(
      children: [
        if (addApprovers && activeWallet!.primaryWallet == 0) ...[
          Container(
              height: height / 1.219,
              width: width,
              child: Stepper(
                key: allKey,
                type: StepperType.vertical,
                currentStep: currentStep,
                controlsBuilder: (context, _) {
                  return Column(
                    children: [],
                  );
                },
                onStepTapped: (step) => setState(() {
                  currentStep = step;
                }),
                steps: getSteps(),
              )),
        ] else ...[
          SizedBox(
            height: height / 30,
          ),
          Text("Add viewer access",
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontsemibold,
                  fontSize: 17.sp)),
          SizedBox(
            height: height / 30,
          ),
          Container(
              // height: height / 1.219,
              width: width,
              // color: Colors.black,
              child: showViewers()),
        ],
      ],
    );
  }

  void submitSharedAccessForm() {
    appState.viewData = {
      AddSharedAccessDetailsViewPageConfig.key: {
        'viewers': viewers,
        'addApprovers': addApprovers,
        'approvers': approvers,
        'noOfApprovers': addApprovers ? noOfApprovers : 0,
        'noOfApprovalsNeeded': addApprovers ? noOfApprovalsNeeded : 0,
        'initiators': initiators,
        'userFullnames': userFullnames,
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
        chooseWallet(),
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
                                  userItem(
                                      '${viewers[i]} [${userFullnames[viewers[i]]}]',
                                      () {
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

            if (viewers.contains(username)) {
              viewerUsernameErrorMessage = 'Username already added';
              setState(() {});
              return;
            }

            // if the username is already on the approvers list then there's an error
            if (approvers.contains(username)) {
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
            if (initiators.contains(username)) {
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

            userFullnames[username] =
                '${userInfo['firstName']} ${userInfo['lastName']}';
            viewers.add(username);
            viewersController.text = '';
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
          height: height / 70,
        ),
        // if wallet is not primary wallet
        // primary wallets can only have view-only shared access
        // the cannot have approver and initiator shared access
        if (activeWallet!.primaryWallet == 0) ...[
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
                        addApprovers = value ?? false;
                      });

                      setState(() {
                        // move straight to the next step
                        currentStep = 1;
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
          height: height / 50,
        ),
        ElevatedButton(
          onPressed: () {
            if (!addApprovers && viewers.isEmpty) {
              popup(context,
                  title: 'Error!',
                  message:
                      'Please enter the username of those you want to grant access to this wallet');
              return;
            }

            if (addApprovers) {
              bool isLastStep = (currentStep == getSteps().length - 1);
              if (isLastStep) {
                submitSharedAccessForm();
              } else {
                setState(() {
                  currentStep += 1;
                });
              }
            } else {
              submitSharedAccessForm();
            }
          },
          style: ButtonStyle(
            backgroundColor:
                MaterialStateProperty.all<Color>(notifier.getbluecolor!),
            shape: MaterialStateProperty.all<RoundedRectangleBorder>(
              const RoundedRectangleBorder(
                borderRadius: BorderRadius.all(
                  Radius.circular(15),
                ),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Text(
              addApprovers ? 'Add approver access' : 'Proceed',
              style: TextStyle(
                fontFamily: fontsemibold,
                fontSize: 14.sp,
              ),
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

  Widget showApprovers() {
    return Form(
      key: approversFormKey,
      child: Column(
        children: [
          SizedBox(
            height: height / 70,
          ),
          Container(
            height: height / 10,
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
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
                                    userItem(
                                        '${approvers[i]} [${userFullnames[approvers[i]]}]',
                                        () {
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

              if (approvers.contains(username)) {
                approverUsernameErrorMessage = 'Username already added';
                setState(() {});
                return;
              }

              if (approvers.length == noOfApprovers) {
                popup(context,
                    title: 'Error!',
                    message:
                        'Number of usernames cannot be more than the number of approvers you selected');
                return;
              }

              // if the username is already on the viewers list then there's
              // no need to check again that the username is valid so we add it to
              // to the approvers list
              if (viewers.contains(username)) {
                showResponseMessage(context,
                    'This user will be removed from the view-only access as they will have implicit view access as an approver',
                    () {
                  viewers.removeWhere((userItem) => userItem == username);
                  approvers.add(username);
                  approversController.text = '';
                });
                setState(() {});
                return;
              }

              // if the username is already on the approvers list then there's
              // no need to check again that the username is valid so we add it to
              // to the approvers list
              if (initiators.contains(username)) {
                approvers.add(username);
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
              userFullnames[username] =
                  '${userInfo['firstName']} ${userInfo['lastName']}';
              approvers.add(username);
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
          SizedBox(
            height: height / 50,
          ),
          ElevatedButton(
            onPressed: () {
              if (!approversFormKey.currentState!.validate()) {
                return;
              }

              if (approvers.isEmpty) {
                popup(context,
                    title: 'Error!',
                    message:
                        'Please enter the username of those you want to grant approver access to this wallet');
                return;
              }

              if (approvers.length < noOfApprovers) {
                popup(context,
                    title: 'Error!',
                    message:
                        'Number of usernames cannot be less than the number of approvers you selected');
                return;
              }

              bool isLastStep = (currentStep == getSteps().length - 1);
              if (isLastStep) {
                submitSharedAccessForm();
              } else {
                setState(() {
                  currentStep += 1;
                });
              }
            },
            style: ButtonStyle(
              backgroundColor:
                  MaterialStateProperty.all<Color>(notifier.getbluecolor!),
              shape: MaterialStateProperty.all<RoundedRectangleBorder>(
                const RoundedRectangleBorder(
                  borderRadius: BorderRadius.all(
                    Radius.circular(15),
                  ),
                ),
              ),
            ),
            child: Padding(
              padding: const EdgeInsets.all(20.0),
              child: Text(
                'Add initiator access',
                style: TextStyle(
                  fontFamily: fontsemibold,
                  fontSize: 14.sp,
                ),
              ),
            ),
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
                                  userItem(
                                      '${initiators[i]} [${userFullnames[initiators[i]]}]',
                                      () {
                                    initiators.removeAt(i);
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

            if (initiators.contains(username)) {
              initiatorUsernameErrorMessage = 'Username already added';
              setState(() {});
              return;
            }

            // if the username is already on the viewers list then there's
            // no need to check again that the username is valid so we add it to
            // to the approvers list
            if (viewers.contains(username)) {
              // initiators.add(username);
              // initiatorsController.text = '';
              showResponseMessage(context,
                  'This user will be removed from the view-only access as they will have implicit view access as an initiator',
                  () {
                viewers.removeWhere((username) => username == username);
                initiators.add(username);
                initiatorsController.text = '';
              });
              setState(() {});
              return;
            }

            // if the username is already on the approvers list then there's
            // no need to check again that the username is valid so we add it to
            // to the approvers list
            if (approvers.contains(username)) {
              initiators.add(username);
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
            userFullnames[username] =
                '${userInfo['firstName']} ${userInfo['lastName']}';
            initiators.add(username);
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
        SizedBox(
          height: height / 50,
        ),
        ElevatedButton(
          onPressed: () {
            if (initiators.isEmpty) {
              popup(context,
                  title: 'Error!',
                  message:
                      'Please enter the username of those you want to grant initiator access to this wallet');
              return;
            }

            if (approvers.length < noOfApprovers) {
              popup(context,
                  title: 'Error!',
                  message:
                      'Number of approver usernames cannot be less than the number of approvers you selected. Please go back and add more approvers.');
              return;
            }

            if (approvers.isEmpty) {
              popup(context,
                  title: 'Error!',
                  message:
                      'You cannot have initiators without having approvers. Please add approvers.');
              return;
            }

            bool isLastStep = (currentStep == getSteps().length - 1);
            if (isLastStep) {
              submitSharedAccessForm();
            } else {
              setState(() {
                currentStep += 1;
              });
            }
          },
          style: ButtonStyle(
            backgroundColor:
                MaterialStateProperty.all<Color>(notifier.getbluecolor!),
            shape: MaterialStateProperty.all<RoundedRectangleBorder>(
              const RoundedRectangleBorder(
                borderRadius: BorderRadius.all(
                  Radius.circular(15),
                ),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.all(20.0),
            child: Text(
              'Proceed',
              style: TextStyle(
                fontFamily: fontsemibold,
                fontSize: 14.sp,
              ),
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
            alignment: WrapAlignment.center,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              Text(
                name,
                textAlign: TextAlign.center,
                softWrap: true,
                style: TextStyle(
                    color: wihitecolor, fontFamily: fontbody, fontSize: 12.sp),
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
            Text(
              LanguageEn.choosewallet,
              style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontFamily: fontbody,
                  fontSize: 15.sp),
            ),
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
                    addApprovers = false;
                    initiators = [];
                    approvers = [];
                    viewers = [];
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

  // we need to check that the username entered here is a valid
  // username of an active trovo account
  Future<Map?> checkUsername(String username) async {
    try {
      showLoader(context);
      Map responseData = await makeGetRequest(
        uri: '/v1/users/$username',
        signer: activeWallet!.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: activeWallet!.publicKey!,
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
}

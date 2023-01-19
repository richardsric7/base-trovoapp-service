import 'dart:math';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:get/get_utils/src/extensions/string_extensions.dart';
import 'package:intl/intl.dart';
import 'package:loadmore/loadmore.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/permission.dart';
import 'package:trovo_wallet/models/wallet.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/router/page_actions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import 'package:trovo_wallet/widgets/utilities.dart';
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
  TextEditingController viewersController = TextEditingController();
  TextEditingController approversController = TextEditingController();
  TextEditingController initiatorsController = TextEditingController();
  final approversFormKey = GlobalKey<FormState>();
  late RefreshController _refreshController;
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
  ApprovalsListFilterType filterType =
      ApprovalsListFilterType.TransactionStatus;
  var filterTypesMap = {
    ApprovalsListFilterType.TransactionStatus: "Transaction status",
    ApprovalsListFilterType.TransactionType: "Transaction type",
    ApprovalsListFilterType.DateRange: "Date range",
    ApprovalsListFilterType.TransactionId: "Transaction ID",
    ApprovalsListFilterType.Initiator: "Initiator",
    ApprovalsListFilterType.Description: "Description",
    ApprovalsListFilterType.WalletPublicKey: "Wallet public key",
    ApprovalsListFilterType.WalletAlias: "Wallet alias",
  };
  var transactionTypes = <String>[
    'ALL',
    'SWAP',
    'PAYMENT',
    'MODIFY SHARED ACCESS',
    'DISABLE SHARED ACCESS',
    'OPT IN ASSET',
    'OPT OUT ASSET',
    'ACCEPT PENDING ASSET',
    'REJECT PENDING ASSET',
    'MAKE MARKET OFFER',
    'DELETE MARKET OFFER',
    'MODIFY MARKET OFFER',
    'MINT TOKEN',
    'BURN TOKEN',
    'BULK PAYMENT',
  ];

  var transactionStatus = <String>[
    'ALL',
    'PENDING',
    'COMPLETED',
    'REJECTED',
  ];

  List<DropdownMenuItem<String>> get transactionTypeDropdownItems {
    return transactionTypes
        .map<DropdownMenuItem<String>>((item) => DropdownMenuItem(
            child: Text(
              item.capitalizeFirst!,
              overflow: TextOverflow.ellipsis,
            ),
            value: item))
        .toList();
  }

  List<DropdownMenuItem<ApprovalsListFilterType>> get filterTypeDropdownItems {
    List<DropdownMenuItem<ApprovalsListFilterType>> items = [];
    filterTypesMap.forEach((key, value) {
      items.add(
        DropdownMenuItem(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(
                  value,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
            ),
            value: key),
      );
    });

    return items;
  }

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
    for (var i = 3; i < 20; i++) {
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
    for (var i = 2; i < noOfApprovers; i++) {
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
    getdarkmodepreviousstate();
    _refreshController = RefreshController(initialRefresh: false);
    appState = Provider.of<DataProvider>(context, listen: false);
    appState.totalRecords = 0;
    appState.filterTransactionStatus = 'Pending';
    appState.filterQuery = "&transactionStatus=PENDING";
    appState.approvals = appState.fetchApprovals(
      limit: appState.limit.toString(),
      query: appState.filterQuery,
    );
    appState.sharedAccesstabController = TabController(length: 3, vsync: this);
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
        body: SmartRefresher(
          enablePullDown: true,
          controller: _refreshController,
          onRefresh: refreshData,
          child: Column(
            children: [
              Padding(
                padding: const EdgeInsets.fromLTRB(20, 12.0, 20, 10.0),
                child: TabBar(
                  controller: appState.sharedAccesstabController,
                  labelColor: notifier.getbluewhitecolor,
                  indicatorColor: notifier.getbluewhitecolor,
                  labelStyle: TextStyle(
                    fontSize: 15.sp,
                    fontFamily: fontsemibold,
                  ),
                  tabs: [
                    Tab(
                      height: 50,
                      child: Icon(Icons.people_alt_outlined),
                    ),
                    Tab(
                      height: 50,
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(CupertinoIcons.square_list),
                          FutureBuilder<Map>(
                            future: appState.approvals,
                            builder: (context, snapshot) {
                              if (snapshot.connectionState ==
                                      ConnectionState.done &&
                                  snapshot.hasData) {
                                var noOfTransactionsToSign =
                                    appState.filterQuery.contains('PENDING')
                                        ? snapshot.data!['totalRecords'] ?? 0
                                        : 0;
                                if (noOfTransactionsToSign > 0) {
                                  return Text(
                                    '($noOfTransactionsToSign)',
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.bold,
                                      fontFamily: fontsemibold,
                                    ),
                                  );
                                }
                              }

                              return Text(
                                '',
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 15,
                                  fontWeight: FontWeight.bold,
                                  fontFamily: fontsemibold,
                                ),
                              );
                            },
                          ),
                        ],
                      ),
                    ),
                    Tab(
                      height: 50,
                      child: Icon(Icons.group_add_outlined),
                    ),
                  ],
                ),
              ),
              Container(
                height: height / 1.22,
                child: TabBarView(
                    controller: appState.sharedAccesstabController,
                    children: [
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
      ),
    );
  }

  Widget pendingApprovals() {
    return Container(
      height: height / 1.22,
      child: SingleChildScrollView(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            SizedBox(height: height / 30),
            Container(
              width: width,
              child: Row(
                children: [
                  SizedBox(
                    width: width / 50,
                  ),
                  Expanded(
                    flex: 2,
                    child: dropdown(
                      (newValue) async {
                        filterType = newValue as ApprovalsListFilterType;
                        showPopup(filterType);
                        setState(() {});
                      },
                      filterTypeDropdownItems,
                      ApprovalsListFilterType.TransactionStatus,
                      null,
                    ),
                  ),
                  Expanded(
                    flex: 2,
                    child: getContent(filterType),
                  ),
                  SizedBox(
                    width: width / 50,
                  ),
                ],
              ),
            ),
            SizedBox(
              height: height / 50,
            ),
            FutureBuilder<Map>(
              future: appState.approvals,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return Container(
                    height: height / 1.5,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        CircularProgressIndicator(
                          backgroundColor: notifier.getbluecolor,
                          valueColor: new AlwaysStoppedAnimation<Color>(
                            notifier.getgreencolor,
                          ),
                          strokeWidth: 3.0,
                        ),
                      ],
                    ),
                  );
                } else if (snapshot.connectionState == ConnectionState.done) {
                  if (snapshot.hasError) {
                    return Container(
                      height: height / 1.5,
                      child: Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Text(
                              LanguageEn.somethingwentwrong,
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody),
                            ),
                            ElevatedButton(
                              onPressed: () {
                                setState(() {
                                  appState.getApprovals();
                                });
                              },
                              style: ButtonStyle(
                                backgroundColor:
                                    MaterialStateProperty.all<Color>(
                                        notifier.getbluecolor!),
                              ),
                              child: Text(
                                LanguageEn.retry,
                                style: TextStyle(
                                  fontFamily: fontsemibold,
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  } else if (snapshot.hasData) {
                    var records = snapshot.data!['records'];
                    appState.totalRecords = snapshot.data!['totalRecords'];

                    WidgetsBinding.instance.addPostFrameCallback((_) async {
                      if (appState.filterQuery.contains('PENDING') &&
                          appState.totalRecords == 0) {
                        appState.excludeUserApproved = 0;
                        appState.filterTransactionStatus = 'All';
                        appState.filterQuery = '';
                        await appState.getApprovals();
                      }
                    });

                    if (records.length > 0) {
                      return LoadMore(
                        isFinish: records.length == appState.totalRecords,
                        onLoadMore: () async {
                          appState.limit += 2;
                          await appState.getApprovals();
                          return records.length <= appState.totalRecords!;
                        },
                        textBuilder: (LoadMoreStatus status) {
                          String text;
                          switch (status) {
                            case LoadMoreStatus.fail:
                              text = "Tap to load more";
                              break;
                            case LoadMoreStatus.idle:
                              text = "Tap to load more";
                              break;
                            default:
                              text = "";
                          }
                          return text;
                        },
                        child: Column(
                          children: [
                            for (var i = 0; i < records.length; i++) ...[
                              GestureDetector(
                                onTap: () {
                                  appState.viewData![
                                          ApprovalDetailsViewPageConfig.key] =
                                      records[i];

                                  appState.currentAction = PageAction(
                                    state: PageState.addPage,
                                    page: ApprovalDetailsViewPageConfig,
                                  );
                                },
                                child: approvalItem(
                                  walletAlias: records[i]['alias'],
                                  initiator: records[i]['initiator'],
                                  transactionType: records[i]
                                      ['transactionType'],
                                  approvalsNeeded: records[i]
                                      ['approvalsNeeded'],
                                  approvalsGotten: records[i]
                                      ['approvalsGotten'],
                                  transactionStatus: records[i]
                                      ['transactionStatus'],
                                  createdAt: DateTime.tryParse(
                                      records[i]['createdAt'])!,
                                ),
                              ),
                            ],
                          ],
                        ),
                      );
                    }

                    return Container(
                      height: height / 1.9,
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            'No results here.',
                            overflow: TextOverflow.visible,
                            textAlign: TextAlign.center,
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ],
                      ),
                    );
                  } else {
                    return Center(
                      child: Text(
                        'Error fetching data. Please try again',
                        overflow: TextOverflow.visible,
                        textAlign: TextAlign.center,
                        style: TextStyle(
                          fontSize: 15,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    );
                  }
                } else {
                  return Text('State: ${snapshot.connectionState}');
                }
              },
            ),
            SizedBox(
              height: height / 15,
            ),
          ],
        ),
      ),
    );
  }

  Widget approvalItem(
      {required String walletAlias,
      required String initiator,
      required String transactionType,
      required int approvalsNeeded,
      required int approvalsGotten,
      required String transactionStatus,
      required DateTime createdAt}) {
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
                      fontSize: 20,
                      fontFamily: fontsemibold,
                      color: notifier.getbluewhitecolor,
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 70,
              ),
              Container(
                width: width / 1.2,
                child: Row(
                  children: [
                    Padding(
                      padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                      child: Text(
                        'Transaction type:',
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
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Container(
                      width: width / 1.5,
                      child: Text(
                        transactionType.capitalizeFirst!,
                        overflow: TextOverflow.visible,
                        style: TextStyle(
                          fontSize: 14,
                          fontFamily: fontbody,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 90,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      'Initiated by:',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: width / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      initiator,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 90,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      'Approval status:',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: width / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      '$approvalsGotten/$approvalsNeeded',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 90,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      'Transaction status:',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: width / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      transactionStatus.capitalizeFirst!,
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
              SizedBox(
                height: height / 90,
              ),
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      'Initiated:',
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                  SizedBox(
                    width: width / 70,
                  ),
                  Padding(
                    padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                    child: Text(
                      DateFormat('MMMM dd, yyyy hh:mm a').format(createdAt),
                      style: TextStyle(
                        fontSize: 15,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget accessList() {
    return Container(
      height: height / 1.22,
      child: SingleChildScrollView(
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
            if (selectedAccessMode == 'Access granted to me') ...[
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
              getAccessGrantedToMe(appState),
            ] else ...[
              Builder(builder: (context) {
                // filter the wallets to get the one that granted only viewer
                // access to others. Once I grant others approver and initiator access
                // the wallet no longer belongs to me.
                var filteredWallets = <Wallet>[];
                wallets!.forEach((wallet) {
                  if (wallet.permissions!.isNotEmpty &&
                      wallet.permissions!
                          .where((permission) =>
                              permission.permission == 'INITIATOR' ||
                              permission.permission == 'APPROVER')
                          .isEmpty) {
                    filteredWallets.add(wallet);
                  }
                });
                return Column(
                  children: [
                    if (filteredWallets.length > 0) ...[
                      for (var walletIndex = 0;
                          walletIndex < filteredWallets.length;
                          walletIndex++) ...{
                        if (filteredWallets[walletIndex].permissions != null &&
                            filteredWallets[walletIndex].permissions!.length >
                                0) ...[
                          accessGrantedByMe(
                              numberOfApprovalsNeeded:
                                  filteredWallets[walletIndex]
                                      .numberOfApprovalsNeeded!,
                              isPrimaryWallet:
                                  filteredWallets[walletIndex].primaryWallet!,
                              walletPublicKey:
                                  filteredWallets[walletIndex].publicKey!,
                              walletAlias: filteredWallets[walletIndex].alias!,
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
                            color: notifier.getbluewhitecolor,
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
    );
  }

  Widget getAccessGrantedToMe(DataProvider appState) {
    // we rename 'Viewer' to 'VIEW-ONLY' because that's what is
    // returned from the server.
    String filter = selectedFilter == 'Viewer'
        ? 'VIEW-ONLY'
        : selectedFilter.toString().toUpperCase();
    // create a map to hold each wallet info
    var wallets = {};
    // store the wallet aliases as keys here so that we can use it to easily get the values back
    // from the wallets map since we cannot create widgets by looping through the map using map.forEach((k,v))
    var walletKeys = [];
    for (var i = 0; i < appState.sharedWallets.length; i++) {
      // filter shared access by permission
      if (appState.sharedWallets[i]['permission'] == filter ||
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

        if (appState.sharedWallets[i]['walletSettings'] != null) {
          wallets[appState.sharedWallets[i]['walletAlias']]['walletSettings'] =
              appState.sharedWallets[i]['walletSettings'];
        }

        print(
            '=========== ${wallets[appState.sharedWallets[i]['walletAlias']]}');
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
                color: notifier.getbluewhitecolor,
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
                      color: notifier.getbluewhitecolor,
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
                        color: notifier.getbluewhitecolor,
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
                        color: notifier.getbluewhitecolor,
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
                        color: notifier.getbluewhitecolor,
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
                          color: notifier.getbluewhitecolor,
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
                            color: notifier.getbluewhitecolor,
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
      {required String walletAlias,
      required String walletPublicKey,
      required int isPrimaryWallet,
      required int numberOfApprovalsNeeded,
      required List<Permission>? permissions}) {
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
        GestureDetector(
          onTap: () {
            appState.viewData![UpdateSharedAccessViewPageConfig.key] = {
              'walletPublicKey': walletPublicKey,
              'walletAlias': walletAlias,
              'isPrimaryWallet': isPrimaryWallet,
              'viewers': viewersList,
              'approvers': approversList,
              'initiators': initiatorsList,
              'numberOfApprovalsNeeded': numberOfApprovalsNeeded,
            };
            print('=======viewData: ${appState.viewData}');
            appState.currentAction = PageAction(
                state: PageState.addPage,
                page: UpdateSharedAccessViewPageConfig);
          },
          child: Card(
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
                        color: notifier.getbluewhitecolor,
                      ),
                    ),

                    // filter by viewer
                    if (selectedFilter == 'All' ||
                        selectedFilter == 'Viewer') ...[
                      if (viewersList.length > 0) ...[
                        SizedBox(
                          height: height / 50,
                        ),
                        Row(
                          children: [
                            Padding(
                              padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                              child: Text(
                                'Viewer access',
                                style: TextStyle(
                                  fontSize: 15,
                                  fontFamily: fontsemibold,
                                  color: notifier.getbluewhitecolor,
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
                                getNames(viewersList),
                                style: TextStyle(
                                  fontSize: 13,
                                  fontFamily: fontbody,
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ),
                            SizedBox(width: width / 90),
                          ],
                        )
                      ]
                    ],
                    // filter by approver
                    if (selectedFilter == 'All' ||
                        selectedFilter == 'Approver') ...[
                      if (approversList.length > 0) ...[
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
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ),
                          ],
                        ),
                        Row(children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              getNames(approversList),
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ),
                        ])
                      ]
                    ],
                    // filter by initiator
                    if (selectedFilter == 'All' ||
                        selectedFilter == 'Initiator') ...[
                      if (initiatorsList.length > 0) ...[
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
                                  color: notifier.getbluewhitecolor,
                                ),
                              ),
                            ),
                          ],
                        ),
                        Row(children: [
                          Padding(
                            padding: const EdgeInsets.fromLTRB(0, 3.0, 0, 0),
                            child: Text(
                              getNames(initiatorsList),
                              style: TextStyle(
                                fontSize: 13,
                                fontFamily: fontbody,
                                color: notifier.getbluewhitecolor,
                              ),
                            ),
                          ),
                        ])
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
                                    setState(() {});
                                  },
                                      foreColor: wihitecolor,
                                      backColor: notifier.getbluebackcolor)
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
                    'This user is already added to approver access which gives them view access. Please remove them from approver access if you want to grant them view-only access.',
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
                    'This user is already added to initiator access which gives them view access. Please remove them from initiator access if you want to grant them view-only access.',
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
                                if (noOfApprovalsNeeded >= newValueInt) {
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
                                      setState(() {});
                                    },
                                        foreColor: wihitecolor,
                                        backColor: notifier.getbluebackcolor)
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
                    'This user will be removed from the view-only access as they will have view access as an approver/initiator',
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
                                    setState(() {});
                                  },
                                      foreColor: wihitecolor,
                                      backColor: notifier.getbluebackcolor)
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
                  'This user will be removed from the view-only access as they will have view access as an initiator',
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

  void refreshData() async {
    try {
      appState.getApprovals();
      await appState.refreshData();

      _refreshController.refreshCompleted();
    } catch (e) {
      _refreshController.refreshFailed();
    }
  }

  Widget dropdown(void Function(Object?) onChanged,
      List<DropdownMenuItem<Object>> items, Object? value, String? hint) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: DropdownButtonFormField(
        isDense: true,
        isExpanded: true,
        hint: Container(
          // width: 150, //and here
          child: hint != null
              ? Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      hint,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                        color: notifier.getbluewhitecolor,
                        fontSize: 15,
                        fontFamily: fontsemibold,
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ],
                )
              : null,
        ),
        dropdownColor:
            notifier.isDark ? darktilewhitecolor : notifier.getaddsubwalletgrey,
        decoration: InputDecoration(
          contentPadding: EdgeInsets.symmetric(vertical: 0, horizontal: 10),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide.none,
            borderRadius: BorderRadius.circular(10),
          ),
          border: OutlineInputBorder(
            borderSide: BorderSide.none,
            borderRadius: BorderRadius.circular(10),
          ),
          filled: true,
          fillColor: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        value: value,
        icon: Icon(
          Icons.keyboard_arrow_down_rounded,
          color: notifier.getbluewhitecolor,
        ),
        elevation: 0,
        style: TextStyle(
          color: notifier.getbluewhitecolor,
          fontSize: 15,
          fontFamily: fontsemibold,
          fontWeight: FontWeight.w500,
        ),
        onChanged: onChanged,
        items: items,
      ),
    );
  }

  void showPopup(ApprovalsListFilterType filterType) {
    switch (filterType) {
      case ApprovalsListFilterType.TransactionStatus:
        approvalListTransactionTypePopup(
          context,
          transactionStatus,
          'Select transaction status',
          (status) {
            appState.setFilterTransactionStatus = status.capitalizeFirst;
            appState.setFilterQuery =
                status == 'ALL' ? '' : "&transactionStatus=$status";
            appState.getApprovals();
            Navigator.of(context).pop(); // dismiss dialog,
          },
          rel: ApprovalsListFilterType.TransactionStatus,
        );
        break;
      case ApprovalsListFilterType.TransactionType:
        approvalListTransactionTypePopup(
          context,
          transactionTypes,
          'Select transaction type',
          (transactionType) {
            appState.setFilterTransactionType = transactionType.capitalizeFirst;
            appState.setFilterQuery = "&transactionType=$transactionType";
            appState.getApprovals();
            Navigator.of(context).pop(); // dismiss dialog,
          },
        );
        break;
      case ApprovalsListFilterType.WalletAlias:
        approvalTextFieldPopup(context,
            label: 'Enter wallet alias',
            value: appState.filterWalletAlias,
            placeholder: 'Enter alias', onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterWalletAlias = value;
            appState.setFilterQuery = "&walletAlias=$value";
            await appState.getApprovals(
              onDone: () => {},
            );
          }
        });
        break;
      case ApprovalsListFilterType.Description:
        approvalTextFieldPopup(context,
            label: 'Enter decription',
            value: appState.filterDescription,
            placeholder: 'Enter description', onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterDescription = value;
            appState.setFilterQuery = "&description=$value";
            await appState.getApprovals(
              onDone: () => {},
            );
          }
        });
        break;
      case ApprovalsListFilterType.TransactionId:
        approvalTextFieldPopup(context,
            label: 'Enter Transaction ID',
            value: appState.filterTransactionId,
            placeholder: 'Transaction ID', onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterTransactionId = value;
            appState.setFilterQuery = "&transactionID=$value";
            await appState.getApprovals(
              onDone: () => {},
            );
          }
        });
        break;
      case ApprovalsListFilterType.DateRange:
        customDateRangePopup(context, onDone: () async {
          appState.setFilterQuery =
              "&dateBetween=${DateFormat('yyyy-MM-dd').format(appState.filterStartDate!)}|${DateFormat('yyyy-MM-dd').format(appState.filterEndDate!)}";
          await appState.getApprovals();
        });
        break;
      case ApprovalsListFilterType.WalletPublicKey:
        approvalTextFieldPopup(context,
            label: 'Enter wallet public key',
            value: appState.filterWalletPublicKey,
            placeholder: 'Public Key', onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterWalletPublicKey = value;
            appState.setFilterQuery = "&walletPublicKey=$value";
            await appState.getApprovals(
              onDone: () => {},
            );
          }
        });
        break;
      case ApprovalsListFilterType.Initiator:
        approvalTextFieldPopup(context,
            label: 'Enter initiator username',
            value: appState.filterInitiatorUsername,
            placeholder: 'Username', onDone: (value) async {
          if (value != null && value.toString().isNotEmpty) {
            appState.setFilterInitiatorUsername = value;
            appState.setFilterQuery = "&initiator=$value";
            await appState.getApprovals(
              onDone: () => {},
            );
          }
        });
        break;
      default:
        break;
    }
  }

  Widget content({required void Function() onPressed, required String label}) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 5.0),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(10.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: TextButton(
          onPressed: onPressed,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Container(
                constraints: BoxConstraints(
                  maxWidth: width / 2.9,
                ),
                child: Text(
                  label,
                  textAlign: TextAlign.start,
                  overflow: TextOverflow.ellipsis,
                  style: TextStyle(
                      color: notifier.getbluewhitecolor,
                      fontSize: appState.filterUsername != null ? 12 : 15,
                      fontFamily: fontsemibold),
                ),
              ),
              Icon(
                Icons.keyboard_arrow_down_rounded,
                color: notifier.getbluewhitecolor,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget getContent(ApprovalsListFilterType type) {
    switch (type) {
      case ApprovalsListFilterType.TransactionStatus:
        return content(
          onPressed: () {
            approvalListTransactionTypePopup(
              context,
              transactionStatus,
              'Select transaction status',
              (status) {
                print(status);
                appState.setFilterTransactionStatus = status.capitalizeFirst;
                appState.setFilterQuery =
                    status == 'ALL' ? '' : "&transactionStatus=$status";
                appState.getApprovals();
                Navigator.of(context).pop(); // dismiss dialog,
              },
              rel: ApprovalsListFilterType.TransactionStatus,
            );
          },
          label: appState.filterTransactionStatus.isEmpty
              ? "Choose status"
              : appState.filterTransactionStatus,
        );
      case ApprovalsListFilterType.WalletAlias:
        return content(
          onPressed: () {
            approvalTextFieldPopup(context,
                label: 'Enter wallet alias',
                value: appState.filterWalletAlias,
                placeholder: 'Enter alias', onDone: (value) async {
              if (value != null && value.toString().isNotEmpty) {
                appState.setFilterWalletAlias = value;
                appState.setFilterQuery = "&walletAlias=$value";
                await appState.getApprovals();
              }
            });
          },
          label: appState.filterWalletAlias.isEmpty
              ? "Enter alias"
              : appState.filterWalletAlias,
        );
      case ApprovalsListFilterType.Description:
        return content(
          onPressed: () {
            approvalTextFieldPopup(context,
                label: 'Enter decription',
                value: appState.filterDescription,
                placeholder: 'Enter description', onDone: (value) async {
              if (value != null && value.toString().isNotEmpty) {
                appState.setFilterDescription = value;
                appState.setFilterQuery = "&description=$value";
                await appState.getApprovals();
              }
            });
          },
          label: appState.filterDescription.isEmpty
              ? "Enter description"
              : truncate(appState.filterDescription, length: 30),
        );

      case ApprovalsListFilterType.WalletPublicKey:
        return content(
          onPressed: () {
            approvalTextFieldPopup(context,
                label: 'Enter wallet public key',
                value: appState.filterWalletPublicKey,
                placeholder: 'Public Key', onDone: (value) async {
              if (value != null && value.toString().isNotEmpty) {
                appState.setFilterWalletPublicKey = value;
                appState.setFilterQuery = "&walletPublicKey=$value";
                await appState.getApprovals();
              }
            });
          },
          label: getTruncatedPublicKey(appState.filterWalletPublicKey),
        );
      case ApprovalsListFilterType.TransactionId:
        return content(
          onPressed: () {
            approvalTextFieldPopup(context,
                label: 'Enter Transaction ID',
                value: appState.filterTransactionId,
                placeholder: 'Transaction ID', onDone: (value) async {
              if (value != null && value.toString().isNotEmpty) {
                appState.setFilterTransactionId = value;
                appState.setFilterQuery = "&transactionID=$value";
                await appState.getApprovals();
              }
            });
          },
          label: appState.filterTransactionId.isEmpty
              ? 'Enter ID'
              : appState.filterTransactionId,
        );
      case ApprovalsListFilterType.Initiator:
        return content(
          onPressed: () {
            approvalTextFieldPopup(context,
                label: 'Enter initiator username',
                value: appState.filterInitiatorUsername,
                placeholder: 'Username', onDone: (value) async {
              if (value != null && value.toString().isNotEmpty) {
                appState.setFilterInitiatorUsername = value;
                appState.setFilterQuery = "&initiator=$value";
                await appState.getApprovals();
              }
            });
          },
          label: appState.filterInitiatorUsername.isEmpty
              ? 'Enter username'
              : appState.filterInitiatorUsername,
        );
      case ApprovalsListFilterType.DateRange:
        return content(
          onPressed: () {
            customDateRangePopup(context, onDone: () async {
              appState.setFilterQuery =
                  "&dateBetween=${DateFormat('yyyy-MM-dd').format(appState.filterStartDate!)}|${DateFormat('yyyy-MM-dd').format(appState.filterEndDate!)}";
              await appState.getApprovals();
            });
          },
          label: getDateRangeValue(),
        );
      default: // ApprovalsListFilterType.TransactionType
        return content(
          onPressed: () {
            approvalListTransactionTypePopup(
              context,
              transactionTypes,
              'Select transaction type',
              (transactionType) {
                appState.setFilterTransactionType =
                    transactionType.capitalizeFirst;
                appState.setFilterQuery = "&transactionType=$transactionType";
                appState.getApprovals();
                Navigator.of(context).pop(); // dismiss dialog,
              },
            );
          },
          label: appState.filterTransactionType,
        );
    }
  }

  getDateRangeValue() {
    if (appState.filterStartDate != null && appState.filterEndDate != null) {
      return "${DateFormat('dd/MM/yy').format(appState.filterStartDate!)} - ${DateFormat('dd/MM/yy').format(appState.filterEndDate!)} ";
    }

    return 'Enter range';
  }

  getTruncatedPublicKey(String publicKey) {
    if (publicKey.isEmpty) return "Enter public key";
    if (publicKey.length <= 7) return publicKey;
    return truncate(publicKey, length: 7) +
        publicKey.substring(publicKey.length - 7);
  }

  @override
  void dispose() {
    appState.excludeUserApproved = 1;
    appState.totalRecords = 0;
    appState.filterQuery = '';
    super.dispose();
  }
}

enum ApprovalsListFilterType {
  TransactionType,
  TransactionStatus,
  TransactionId,
  DateRange,
  Initiator,
  Description,
  WalletPublicKey,
  WalletAlias,
}

String getNames(listOfNames) {
  var names = <String>[];
  for (var permIndex = 0; permIndex < listOfNames.length; permIndex++) {
    if (permIndex < 3) {
      names.add(listOfNames[permIndex].targetUsername);
    }
  }

  if (listOfNames.length > 3) {
    names.add('...tap to view all');
  }

  return names.join(', ');
}

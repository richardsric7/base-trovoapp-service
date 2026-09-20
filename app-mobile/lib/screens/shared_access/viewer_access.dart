// import 'package:flutter/material.dart';
// import 'package:flutter_screenutil/flutter_screenutil.dart';
// import 'package:provider/provider.dart';
// import 'package:shared_preferences/shared_preferences.dart';
// import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
// import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
// import 'package:trovo_app/custom_bloc_observer/colors.dart';
// import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
// import 'package:trovo_app/custom_bloc_observer/fonts.dart';
// import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
// import 'package:trovo_app/storage/state.dart';
// import 'package:trovo_app/utils/enstring.dart';
// import 'package:trovo_app/utils/medeiaqury/medeiaqury.dart';
// import 'package:trovo_app/models/wallet.dart';

// class ViewerAccess extends StatefulWidget {
//   const ViewerAccess({Key? key}) : super(key: key);

//   @override
//   State<ViewerAccess> createState() => _ViewerAccessState();
// }

// class _ViewerAccessState extends State<ViewerAccess> {
//   late ColorNotifier notifier;
//   late DataProvider appState;
//   TextEditingController approversController = TextEditingController();
//   TextEditingController viewersController = TextEditingController();
//   List<Wallet>? wallets;
//   Wallet? activeWallet;
//   dynamic selectedWallet = '';
//   var viewers = <String>[];

//   getdarkmodepreviousstate() async {
//     final prefs = await SharedPreferences.getInstance();
//     bool? previusstate = prefs.getBool("setIsDark");
//     if (previusstate == null) {
//       notifier.setIsDark = false;
//     } else {
//       notifier.setIsDark = previusstate;
//     }
//   }

//   @override
//   void initState() {
//     super.initState();
//     getdarkmodepreviousstate();
//   }

//   @override
//   Widget build(BuildContext context) {
//     notifier = Provider.of<ColorNotifier>(context, listen: true);
//     appState = Provider.of<DataProvider>(context, listen: true);
//     height = MediaQuery.of(context).size.height;
//     width = MediaQuery.of(context).size.width;
//     wallets = appState.userInfo!.wallets!;
//     activeWallet = appState.activeWallet;
//     selectedWallet = activeWallet!.address;
//     return ScreenUtilInit(
//       builder: (context, child) => Scaffold(
//         resizeToAvoidBottomInset: false,
//         backgroundColor: notifier.getwihitecolor,
//         appBar: CustomAppBar(
//           context,
//           notifier.getwihitecolor,
//           'Grant viewer access',
//           notifier.getbluewhitecolor,
//           height: height / 15,
//         ),
//         body: Column(
//           children: [
//             SizedBox(
//               height: height / 30,
//             ),
//             showViewers(),
//           ],
//         ),
//       ),
//     );
//   }

//   Widget showViewers() {
//     return Column(
//       children: [
//         Padding(
//           padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
//           child: Container(
//             decoration: BoxDecoration(
//               borderRadius: const BorderRadius.all(Radius.circular(15.0)),
//               color: notifier.isDark
//                   ? darktilewhitecolor
//                   : notifier.getaddsubwalletgrey,
//             ),
//             child: Row(
//               mainAxisAlignment: MainAxisAlignment.center,
//               children: [
//                 Padding(
//                   padding: const EdgeInsets.symmetric(
//                       horizontal: 20.0, vertical: 15.0),
//                   child: Column(
//                     children: [
//                       Container(
//                           width: width / 1.3,
//                           child: Wrap(
//                             alignment: WrapAlignment.center,
//                             children: [
//                               if (viewers.length > 0) ...[
//                                 for (var i = 0; i < viewers.length; i++) ...[
//                                   userItem(viewers[i], () {
//                                     viewers.removeAt(i);
//                                   }, notifier.getbluecolor)
//                                 ],
//                               ] else ...[
//                                 Text(
//                                   'Name of viewers appear here',
//                                   // : "enteraccountsusernameapprovers".tr(),
//                                   style: TextStyle(
//                                       color: notifier.getbluewhitecolor,
//                                       fontFamily: fontbody,
//                                       fontSize: 15.sp),
//                                 ),
//                               ]
//                             ],
//                           )),
//                       SizedBox(height: 2),
//                     ],
//                   ),
//                 ),
//               ],
//             ),
//           ),
//         ),
//         SizedBox(
//           height: height / 30,
//         ),
//         Container(
//           width: width / 1.1,
//           child: Text(
//             "enteraccountsusernameviewers".tr() + 'kenmaddy_bantu',
//             textAlign: TextAlign.center,
//             style: TextStyle(
//                 color: notifier.getbluewhitecolor,
//                 fontFamily: fontbody,
//                 fontSize: 15.sp),
//           ),
//         ),
//         SizedBox(
//           height: height / 50,
//         ),
//         Padding(
//           padding: const EdgeInsets.symmetric(horizontal: 20.0),
//           child: CustomTextFormField.textFieldWithoutIcon(
//             'Viewer',
//             notifier.getbluecolor,
//             notifier.getgrey,
//             notifier.getprefixicon,
//             notifier.getblck,
//             notifier.getgrey,
//             60.sp,
//             210.sp,
//             controller: viewersController,
//           ),
//         ),
//         SizedBox(
//           height: height / 50,
//         ),
//         ElevatedButton(
//           onPressed: () => setState(() {
//             if (viewersController.text.isNotEmpty) {
//               viewers.add(viewersController.text);
//               viewersController.text = '';
//             }
//           }),
//           style: ButtonStyle(
//             backgroundColor:
//                 MaterialStateProperty.all<Color>(notifier.getbluecolor!),
//           ),
//           child: Text(
//             "add".tr(),
//             style: TextStyle(
//               fontFamily: fontsemibold,
//             ),
//           ),
//         ),
//         SizedBox(
//           height: height / 70,
//         ),
//         // Padding(
//         //   padding: const EdgeInsets.symmetric(horizontal: 20.0),
//         //   child: Row(
//         //     crossAxisAlignment: CrossAxisAlignment.center,
//         //     children: [
//         //       Transform.scale(
//         //         scale: 1.sp,
//         //         child: Checkbox(
//         //           shape: RoundedRectangleBorder(
//         //             borderRadius: BorderRadius.all(
//         //               Radius.circular(5.sp),
//         //             ),
//         //           ),
//         //           activeColor: notifier.getbluecolor,
//         //           side: BorderSide(color: notifier.getbluewhitecolor),
//         //           value: addApprovers,
//         //           onChanged: (value) {
//         //             setState(() {
//         //               addApprovers = value ?? false;
//         //             });
//         //           },
//         //         ),
//         //       ),
//         //       Container(
//         //         child: Text(
//         //           "doyouwanttoaddapprovers".tr(),
//         //           overflow: TextOverflow.visible,
//         //           style: TextStyle(
//         //             fontSize: 15,
//         //             fontFamily: fontsemibold,
//         //             color: notifier.getbluewhitecolor,
//         //           ),
//         //         ),
//         //       ),
//         //     ],
//         //   ),
//         // ),
//         SizedBox(
//           height: height / 10,
//         ),
//         Button(
//           'Proceed to grant approver access',
//           notifier.getbluecolor,
//           wihitecolor,
//           onTap: () {
//             // appState.viewData = {
//             //   AddSharedAccessDetailsViewPageConfig.key: {
//             //     'viewers': viewers,
//             //     'approvers': approvers,
//             //     'accessType': selectedAccessType,
//             //     'noOfApprovalsNeeded': noOfApprovalsNeeded,
//             //     'initiators': initiators,
//             //     'addApprovers': addApprovers,
//             //   }
//             // };
//             // appState.currentAction = PageAction(
//             //   state: PageState.addPage,
//             //   page: AddSharedAccessDetailsViewPageConfig,
//             // );
//             setState(() {});
//           },
//         ),
//         SizedBox(height: height / 50),
//         ButtonOutlined(
//           'Proceed without granting approver access',
//           notifier.getwihitecolor,
//           notifier.getbluewhitecolor,
//           onTap: () {
//             // appState.viewData = {
//             //   AddSharedAccessDetailsViewPageConfig.key: {
//             //     'viewers': viewers,
//             //     'approvers': approvers,
//             //     'accessType': selectedAccessType,
//             //     'noOfApprovalsNeeded': noOfApprovalsNeeded,
//             //     'initiators': initiators,
//             //     'addApprovers': addApprovers,
//             //   }
//             // };
//             // appState.currentAction = PageAction(
//             //   state: PageState.addPage,
//             //   page: AddSharedAccessDetailsViewPageConfig,
//             // );
//             // appState.grantSharedAccessView.view = GrantSharedAccessView.viewers;
//             // setState(() {});
//           },
//         ),
//       ],
//     );
//   }

//   Widget userItem(String name, void Function() onRemove, Color color) {
//     return Padding(
//       padding: const EdgeInsets.all(3.0),
//       child: Container(
//         decoration: BoxDecoration(
//             borderRadius: const BorderRadius.all(Radius.circular(10.0)),
//             color: color),
//         child: Padding(
//           padding: const EdgeInsets.all(5.0),
//           child: Wrap(
//             children: [
//               Text(
//                 name,
//                 overflow: TextOverflow.ellipsis,
//                 softWrap: true,
//                 style: TextStyle(
//                     color: wihitecolor, fontFamily: fontbody, fontSize: 15.sp),
//               ),
//               SizedBox(
//                 width: width / 70,
//               ),
//               GestureDetector(
//                   onTap: () => setState(() {
//                         onRemove();
//                       }),
//                   child: Icon(
//                     Icons.cancel_outlined,
//                     color: wihitecolor,
//                   ))
//             ],
//           ),
//         ),
//       ),
//     );
//   }
// }

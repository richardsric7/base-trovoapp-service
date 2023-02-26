import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:image_cropper/image_cropper.dart';
import 'package:image_picker/image_picker.dart';
import 'package:trovo_wallet/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/functions/trovo-sdk.dart';
import 'package:trovo_wallet/network/requests.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/loader.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ProfileDetails extends StatefulWidget {
  const ProfileDetails({Key? key}) : super(key: key);

  @override
  State<ProfileDetails> createState() => _ProfileDetailsState();
}

class _ProfileDetailsState extends State<ProfileDetails> {
  late ColorNotifier notifier;
  late DataProvider appState;

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
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          LanguageEn.myprofile,
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 20,
              ),
              GestureDetector(
                onTap: () {
                  imageSourceDialog(
                    context,
                    onCamera: () {
                      getImage(ImageSource.camera);
                    },
                    onGallery: () {
                      getImage(ImageSource.gallery);
                    },
                  );
                },
                child: Center(
                  child: GestureDetector(
                    onTap: () {
                      imageSourceDialog(
                        context,
                        onCamera: () {
                          getImage(ImageSource.camera);
                        },
                        onGallery: () {
                          getImage(ImageSource.gallery);
                        },
                      );
                    },
                    child: CircleAvatar(
                        radius: width / 10,
                        backgroundColor: notifier.getbluecolor70,
                        child: GestureDetector(
                          onTap: () {
                            imageSourceDialog(
                              context,
                              onCamera: () {
                                getImage(ImageSource.camera);
                              },
                              onGallery: () {
                                getImage(ImageSource.gallery);
                              },
                            );
                          },
                          child: ClipRRect(
                            borderRadius: BorderRadius.circular(100.0),
                            child: Image.network(
                              appState.userInfo!.imageThumbnailURL!,
                              width: width / 5.3,
                              // height: width / 10,
                              fit: BoxFit.fill,
                              errorBuilder: (context, error, stackTrace) {
                                return Image.asset(
                                  'assets/images/trovo.png',
                                  width: width / 9,
                                );
                              },
                            ),
                          ),
                        )),
                  ),
                ),
              ),
              SizedBox(
                height: height / 80,
              ),
              Text(
                '${appState.userInfo!.firstName} ${appState.userInfo!.lastName} ${appState.userInfo!.isCorporate ? '(Corporate)' : ''}',
                style: TextStyle(
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold,
                    fontSize: 16.sp),
              ),
              Text(
                '@${appState.userInfo!.username}',
                style: TextStyle(
                    color: notifier.getgrey,
                    fontFamily: fontsemibold,
                    fontSize: 13.sp),
              ),
              Text(
                '${LanguageEn.referralid}: ${appState.userInfo!.username}',
                style: TextStyle(
                    color: notifier.getgrey,
                    fontFamily: fontsemibold,
                    fontSize: 13.sp),
              ),
              // SizedBox(height: height / 20),
              // Padding(
              //   padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //   child: Row(
              //     children: [
              //       Text(
              //         LanguageEn.bio,
              //         style: TextStyle(
              //             color: notifier.getbluewhitecolor,
              //             fontFamily: fontsemibold,
              //             fontSize: 16.sp),
              //       ),
              //     ],
              //   ),
              // ),
              bioInfo(),
              // SizedBox(height: height / 20),
              // Padding(
              //   padding: const EdgeInsets.symmetric(horizontal: 20.0),
              //   child: Row(
              //     children: [
              //       Text(
              //         LanguageEn.socials,
              //         style: TextStyle(
              //             color: notifier.getbluewhitecolor,
              //             fontFamily: fontsemibold,
              //             fontSize: 16.sp),
              //       ),
              //     ],
              //   ),
              // ),
              // socials(),
              SizedBox(height: height / 50),
              Text(
                'Referral Info (Downlines)',
                style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: notifier.getbluewhitecolor,
                    fontFamily: fontsemibold),
              ),
              SizedBox(height: height / 50),
              Container(
                color: Colors.white,
                padding: EdgeInsets.symmetric(horizontal: 20.0),
                child: Table(
                  border: TableBorder.all(
                      width: 1.5,
                      color: notifier.getbluewhitecolor,
                      borderRadius: BorderRadius.circular(15)),
                  children: getTableRows(),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }

  List<TableRow> getTableRows() {
    List<TableRow> tableRows = [];
    var index = 1;

    tableRows.add(TableRow(children: [
      Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
        child: Text(
          'Referral Level',
          style: TextStyle(
              fontFamily: fontsemibold, color: notifier.getbluewhitecolor),
        ),
      ),
      Padding(
        padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 10),
        child: Text(
          'No. Referred',
          style: TextStyle(
              fontFamily: fontsemibold, color: notifier.getbluewhitecolor),
        ),
      ),
    ]));

    for (MapEntry entry in appState.userInfo!.referralInfo!.downlines.entries) {
      tableRows.add(
        TableRow(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                'Level $index',
                style: TextStyle(
                    fontFamily: fontbody, color: notifier.getbluewhitecolor),
              ),
            ),
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 5, horizontal: 10),
              child: Text(
                entry.value.toString(),
                style: TextStyle(
                    fontFamily: fontbody, color: notifier.getbluewhitecolor),
              ),
            ),
          ],
        ),
      );
      index++;
    }
    return tableRows;
  }

  Future<void> getImage(ImageSource source) async {
    var image = await ImagePicker().pickImage(source: source);
    if (image != null) {
      var croppedImage = await ImageCropper().cropImage(
          sourcePath: image.path,
          cropStyle: CropStyle.circle,
          aspectRatio: CropAspectRatio(ratioX: 1, ratioY: 1),
          compressQuality: 100,
          maxHeight: 800,
          maxWidth: 800,
          compressFormat: ImageCompressFormat.jpg,
          uiSettings: [
            AndroidUiSettings(
              toolbarColor: notifier.getbluecolor80,
              toolbarTitle: 'Crop Image',
            ),
            IOSUiSettings(
              title: 'Crop Image',
            ),
          ]);
      if (croppedImage != null) {
        await uploadImage(croppedImage);
      }
    }
  }

  Future<void> uploadImage(croppedImage) async {
    print('-----------------------${croppedImage.path}');
    try {
      showLoader(context);
      // make initial request to the server using the
      // following credentials
      var primaryWalletKeyPair =
          TrovoWalletSDK().parseSecretKey(appState.secretKeys[0]);

      Map responseData = await makePutRequestForMultipartFile(
        uri: '/v1/users/upload-picture',
        multipartFilePath: croppedImage.path,
        signer: primaryWalletKeyPair.publicKey,
        secretKey: primaryWalletKeyPair.secretKey,
        publicKey: primaryWalletKeyPair.publicKey,
      );
      print('----------this is responseData: $responseData');
      if (responseData['statusCode'] == 200) {
        String imageUrl = responseData['data'].toString().replaceAll('"', '');
        appState.userInfo!.imageThumbnailURL = imageUrl;
        appState.updateListeners();
        print(
            '---------------------appState.userInfo!.imageThumbnailURL: ${appState.userInfo!.imageThumbnailURL}');
        await StoreData()
            .storeInsertData('userInfo', appState.userInfo!.toJSONEncodable());
        setState(() {});
        hideLoader(context);
      } else {
        popup(context,
            title: LanguageEn.error, message: responseData['data']['message']);
        hideLoader(context);
      }
    } catch (e) {
      print(e);
      hideLoader(context);
      popup(context, title: LanguageEn.error, message: e.toString());
    }
  }

  Widget bioInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  Text(
                    LanguageEn.emailadress,
                    style: TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluewhitecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  Text(
                    appState.userInfo!.email!,
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluewhitecolor,
                      fontFamily: fontbody,
                    ),
                  ),
                  SizedBox(
                    height: height / 25,
                  ),
                  // Phone number
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                LanguageEn.phonenumber,
                                style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: height / 90,
                              ),
                              Row(
                                children: [
                                  Text(
                                    appState.userInfo!.mobile!,
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluewhitecolor,
                                      fontFamily: fontbody,
                                    ),
                                  ),
                                  SizedBox(
                                    width: 10,
                                  ),
                                  // GestureDetector(
                                  //   onTap: () {},
                                  //   child: Text(
                                  //     LanguageEn.edit,
                                  //     style: TextStyle(
                                  //       fontSize: 13,
                                  //       fontWeight: FontWeight.w400,
                                  //       color: notifier.getbluewhitecolor,
                                  //       fontFamily: fontbody,
                                  //     ),
                                  //   ),
                                  // ),
                                ],
                              )
                            ],
                          ),
                        ),
                        // Container(
                        //   child: Column(
                        //     children: [
                        //       Text(
                        //         LanguageEn.unverified,
                        //         style: TextStyle(
                        //             fontSize: 11,
                        //             fontWeight: FontWeight.w600,
                        //             color: notifier.getbluewhitecolor,
                        //             fontFamily: fontsemibold),
                        //       ),
                        //       SizedBox(
                        //         height: 5,
                        //       ),
                        //       GestureDetector(
                        //         onTap: () {},
                        //         child: Text(
                        //           LanguageEn.taptoverify,
                        //           style: TextStyle(
                        //             fontSize: 10,
                        //             fontWeight: FontWeight.w400,
                        //             color: notifier.getbluewhitecolor,
                        //             fontFamily: fontbody,
                        //           ),
                        //         ),
                        //       )
                        //     ],
                        //   ),
                        // )
                      ],
                    ),
                  ),
                  // Advanced KYC
                  // SizedBox(
                  //   height: height / 25,
                  // ),
                  // Container(
                  //   width: width / 1.29,
                  //   child: Row(
                  //     mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  //     crossAxisAlignment: CrossAxisAlignment.center,
                  //     children: [
                  //       Container(
                  //         child: Text(
                  //           LanguageEn.advancedkyc,
                  //           style: TextStyle(
                  //               fontSize: 16,
                  //               fontWeight: FontWeight.w600,
                  //               color: notifier.getbluewhitecolor,
                  //               fontFamily: fontsemibold),
                  //         ),
                  //       ),
                  //       Container(
                  //         child: Column(
                  //           children: [
                  //             Text(
                  //               LanguageEn.unverified,
                  //               style: TextStyle(
                  //                   fontSize: 11,
                  //                   fontWeight: FontWeight.w600,
                  //                   color: notifier.getbluewhitecolor,
                  //                   fontFamily: fontsemibold),
                  //             ),
                  //             SizedBox(
                  //               height: 5,
                  //             ),
                  //             GestureDetector(
                  //               onTap: () {},
                  //               child: Text(
                  //                 LanguageEn.taptostart,
                  //                 style: TextStyle(
                  //                   fontSize: 10,
                  //                   fontWeight: FontWeight.w400,
                  //                   color: notifier.getbluewhitecolor,
                  //                   fontFamily: fontbody,
                  //                 ),
                  //               ),
                  //             )
                  //           ],
                  //         ),
                  //       )
                  //     ],
                  //   ),
                  // ),
                  SizedBox(height: 2),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget socials() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.isDark
              ? darktilewhitecolor
              : notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  // Twitter
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Image.asset(
                                "assets/images/twitter.png",
                                height: height / 30,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        ),
                        Container(
                          child: Column(
                            children: [
                              Text(
                                LanguageEn.unverified,
                                style: TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.w600,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: 5,
                              ),
                              GestureDetector(
                                onTap: () {},
                                child: Text(
                                  LanguageEn.taptoverify,
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.w400,
                                    color: notifier.getbluewhitecolor,
                                    fontFamily: fontbody,
                                  ),
                                ),
                              )
                            ],
                          ),
                        )
                      ],
                    ),
                  ),
                  // Instagram
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Image.asset(
                            "assets/images/instagram.png",
                            height: height / 30,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        Container(
                          child: Column(
                            children: [
                              Image.asset(
                                "assets/images/tick.png",
                                height: height / 30,
                                color: notifier.getbluewhitecolor,
                              ),
                            ],
                          ),
                        )
                      ],
                    ),
                  ),
                  // Instagram
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Image.asset(
                            "assets/images/gemcave.png",
                            height: height / 30,
                            color: notifier.getbluewhitecolor,
                          ),
                        ),
                        GestureDetector(
                          onTap: () {},
                          child: Text(
                            LanguageEn.taptoconnect,
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluewhitecolor,
                              fontFamily: fontbody,
                            ),
                          ),
                        )
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
    );
  }
}

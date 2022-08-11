import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_share/flutter_share.dart';
import 'package:image_picker/image_picker.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/widgets/popups.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class ProfileDetails extends StatefulWidget {
  const ProfileDetails({Key? key}) : super(key: key);

  @override
  State<ProfileDetails> createState() => _ProfileDetailsState();
}

class _ProfileDetailsState extends State<ProfileDetails> {
  late ColorNotifier notifier;
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
                  child: Image.asset("assets/images/avatar.png",
                      height: height / 10),
                ),
              ),
              TextButton(
                onPressed: () {
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
                child: Text(
                  LanguageEn.changepicture,
                  style: TextStyle(
                      color: notifier.getgrey,
                      fontFamily: fontsemibold,
                      fontSize: 13.sp),
                ),
              ),
              Text(
                'Obi Enechi',
                style: TextStyle(
                    color: notifier.getbluecolor,
                    fontFamily: 'Gilroy_Bold',
                    fontSize: 16.sp),
              ),
              SizedBox(height: height / 70),
              Text(
                '@efizee',
                style: TextStyle(
                    color: notifier.getgrey,
                    fontFamily: fontsemibold,
                    fontSize: 13.sp),
              ),
              SizedBox(height: height / 70),
              Text(
                '${LanguageEn.referralid}: efizee',
                style: TextStyle(
                    color: notifier.getgrey,
                    fontFamily: fontsemibold,
                    fontSize: 13.sp),
              ),
              SizedBox(height: height / 20),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Row(
                  children: [
                    Text(
                      LanguageEn.bio,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontFamily: 'Gilroy_Bold',
                          fontSize: 16.sp),
                    ),
                  ],
                ),
              ),
              bioInfo(),
              SizedBox(height: height / 20),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: Row(
                  children: [
                    Text(
                      LanguageEn.socials,
                      style: TextStyle(
                          color: notifier.getbluecolor,
                          fontFamily: 'Gilroy_Bold',
                          fontSize: 16.sp),
                    ),
                  ],
                ),
              ),
              socials(),
              SizedBox(height: height / 20),
            ],
          ),
        ),
      ),
    );
  }

  Future<void> getImage(ImageSource source) async {
    var image = await ImagePicker().pickImage(source: source);
    if (image != null) {
      print('${image.mimeType}, ${image.path}');
      image.saveTo('images/profile-pic.${image.name.split('.')[1]}');
    }
  }

  Widget bioInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
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
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(
                    height: height / 90,
                  ),
                  Text(
                    'thundeyy@trovo.io',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluecolor,
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
                                'Phone Number',
                                style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: notifier.getbluecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: height / 90,
                              ),
                              Row(
                                children: [
                                  Text(
                                    '+2347062685682',
                                    textAlign: TextAlign.center,
                                    style: TextStyle(
                                      fontSize: 15,
                                      fontWeight: FontWeight.w400,
                                      color: notifier.getbluecolor,
                                      fontFamily: fontbody,
                                    ),
                                  ),
                                  SizedBox(
                                    width: 10,
                                  ),
                                  GestureDetector(
                                    onTap: () {},
                                    child: Text(
                                      LanguageEn.edit,
                                      style: TextStyle(
                                        fontSize: 13,
                                        fontWeight: FontWeight.w400,
                                        color: notifier.getbluecolor,
                                        fontFamily: fontbody,
                                      ),
                                    ),
                                  ),
                                ],
                              )
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
                                    color: notifier.getbluecolor,
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
                                    color: notifier.getbluecolor,
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
                  // Advanced KYC
                  SizedBox(
                    height: height / 25,
                  ),
                  Container(
                    width: width / 1.29,
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      crossAxisAlignment: CrossAxisAlignment.center,
                      children: [
                        Container(
                          child: Text(
                            LanguageEn.advancedkyc,
                            style: TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.w600,
                                color: notifier.getbluecolor,
                                fontFamily: fontsemibold),
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
                                    color: notifier.getbluecolor,
                                    fontFamily: fontsemibold),
                              ),
                              SizedBox(
                                height: 5,
                              ),
                              GestureDetector(
                                onTap: () {},
                                child: Text(
                                  LanguageEn.taptostart,
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.w400,
                                    color: notifier.getbluecolor,
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
          color: notifier.getaddsubwalletgrey,
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
                                color: notifier.getbluecolor,
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
                                    color: notifier.getbluecolor,
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
                                    color: notifier.getbluecolor,
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
                            color: notifier.getbluecolor,
                          ),
                        ),
                        Container(
                          child: Column(
                            children: [
                              Image.asset(
                                "assets/images/tick.png",
                                height: height / 30,
                                color: notifier.getbluecolor,
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
                            color: notifier.getbluecolor,
                          ),
                        ),
                        GestureDetector(
                          onTap: () {},
                          child: Text(
                            LanguageEn.taptoconnect,
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w400,
                              color: notifier.getbluecolor,
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

  Future<void> share() async {
    await FlutterShare.share(
        title: 'Example share',
        text: 'Example share text',
        linkUrl: 'https://flutter.dev/',
        chooserTitle: 'Example Chooser Title');
  }

  Widget invitefriend(colorbutton, buttontext, buttontextcolor) {
    return Center(
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(15),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            LayoutBuilder(builder: (context, constraints) {
              return Container(
                height: height / 10,
                width: width / 1.1,
                decoration: BoxDecoration(
                  color: colorbutton!,
                  borderRadius: BorderRadius.circular(15),
                ),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                  children: [
                    Image.asset("assets/images/referrals.png",
                        height: height / 30),
                    Container(
                      width: width / 1.7,
                      child: Text(
                        buttontext!,
                        textAlign: TextAlign.start,
                        style: TextStyle(
                            fontFamily: fontbody,
                            fontSize: 13.sp,
                            color: buttontextcolor),
                      ),
                    ),
                    Icon(
                      Icons.arrow_forward_ios,
                      size: 12.sp,
                      color: notifier.getwihitecolor,
                    )
                  ],
                ),
              );
            }),
          ],
        ),
      ),
    );
  }

  Widget iteamlist(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 25),
          Image.asset(
            image,
            height: height / 30,
            color: notifier.getbluecolor,
          ),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
          ),
          const Spacer(),
          SizedBox(width: width / 100),
          Icon(Icons.arrow_forward_ios, color: notifier.getgrey, size: 17.sp),
          SizedBox(width: width / 15),
        ],
      ),
    );
  }

  Widget logout(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 25),
          Image.asset(
            image,
            height: height / 30,
            color: notifier.getbluecolor,
          ),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
          ),
        ],
      ),
    );
  }

  Widget darkmode(image, txt, name) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          SizedBox(width: width / 25),
          Image.asset(
            image,
            height: height / 30,
            color: notifier.getbluecolor,
          ),
          SizedBox(width: width / 40),
          Text(
            name,
            style: TextStyle(
                color: notifier.getblck,
                fontSize: 15.sp,
                fontFamily: 'Gilroy_Medium'),
          ),
          const Spacer(),
          SizedBox(width: width / 100),
          Transform.scale(
            scale: 0.7,
            child: CupertinoSwitch(
              activeColor: notifier.getbluecolor,
              value: notifier.getIsDark,
              onChanged: (val) async {
                final prefs = await SharedPreferences.getInstance();
                setState(() {
                  notifier.setIsDark = val;
                  prefs.setBool("setIsDark", val);
                });
              },
            ),
          ),
          SizedBox(width: width / 15),
        ],
      ),
    );
  }
}

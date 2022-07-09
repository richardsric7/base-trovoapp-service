import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:flutter_svg/svg.dart';
import 'package:get/get.dart';
import 'package:trovo_wallet/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:trovo_wallet/Custom_BlocObserver/fonts.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/User.dart';
import 'package:trovo_wallet/Models/Wallet.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/router/PageActions.dart';
import 'package:trovo_wallet/router/ui_pages.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/enstring.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class AssetDetails extends StatefulWidget {
  const AssetDetails({Key? key}) : super(key: key);

  @override
  State<AssetDetails> createState() => _AssetDetailsState();
}

class _AssetDetailsState extends State<AssetDetails>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  late TabController _tabController;
  late DataProvider appState;
  late UserInfo userInfo;
  var assetBalances;
  var nfts;
  Wallet? activeWallet;
  var claimedAssets;
  var unclaimedAssets;
  int tabLength = 2;
  int touchedIndex = -1;

  @override
  void initState() {
    // TODO: implement initState
    super.initState();
    _tabController = TabController(length: tabLength, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    assetBalances = appState.assetBalances;
    nfts = appState.nfts;
    activeWallet = appState.activeWallet;
    claimedAssets = assetBalances[activeWallet!.publicKey]['claimed'];
    unclaimedAssets = assetBalances[activeWallet!.publicKey]['unclaimed'];
    if (unclaimedAssets.length > 0) {
      setState(() {
        tabLength = 3;
        _tabController = TabController(length: tabLength, vsync: this);
      });
    }
    return ScreenUtilInit(
      builder: (context, child) => Scaffold(
        resizeToAvoidBottomInset: false,
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "",
          notifier.getblck,
          height: height / 15,
        ),
        body: SingleChildScrollView(
          child: Column(
            children: [
              SizedBox(
                height: height / 50,
              ),
              Row(
                children: [
                  SizedBox(
                    width: 20,
                  ),
                  Text(
                    'TROV',
                    style: TextStyle(
                        fontSize: 22,
                        fontWeight: FontWeight.bold,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                ],
              ),
              walletSlides(),
              SizedBox(
                height: height / 30,
              ),
              assetInfo(),
              SizedBox(
                height: height / 50,
              ),
              actionButtons(),
            ],
          ),
        ),
      ),
    );
  }

  Widget actionButtons() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceAround,
      children: [
        ElevatedButton(
          onPressed: () {},
          child: Container(
            child: Padding(
              padding: const EdgeInsets.symmetric(
                vertical: 8.0,
              ),
              child: Column(
                children: [
                  SvgPicture.asset(
                    "assets/images/send.svg",
                    width: width / 6,
                  ),
                  Text(
                    'Send',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
        ElevatedButton(
          onPressed: () {},
          child: Container(
            width: width / 3.9,
            height: height / 10,
            child: Padding(
              padding: const EdgeInsets.symmetric(
                vertical: 8.0,
              ),
              child: Column(
                children: [
                  SvgPicture.asset(
                    "assets/images/recieve.svg",
                    width: width / 6,
                  ),
                  Text(
                    'Recieve',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget walletSlides() {
    var colors = [notifier.getbluecolor, Colors.red, Colors.green];
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: colors[0],
          // color: colors[i - 1],
        ),
        child: Stack(children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 35.0, horizontal: 20),
                child: Image.asset('assets/images/trovo_white.png'),
              ),
            ],
          ),
          Padding(
            padding:
                const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  activeWallet!.alias!.capitalizeFirst!,
                  style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: notifier.getwihitecolor,
                      fontFamily: fontsemibold),
                ),
                SizedBox(
                  height: height / 50,
                ),
                Row(
                  children: [
                    Text(
                      LanguageEn.totalbalance,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getwihitecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ],
                ),
                SizedBox(
                  height: height / 98.0,
                ),
                Text(
                  '2,082,898 NGN',
                  style: TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: notifier.getwihitecolor,
                    fontFamily: fontsemibold,
                  ),
                ),
                SizedBox(height: 2),
                Text(
                  '4,014 USD',
                  style: TextStyle(
                    fontWeight: FontWeight.w300,
                    fontSize: 13,
                    color: notifier.getwihitecolor,
                    fontFamily: fontbody,
                  ),
                ),
              ],
            ),
          ),
        ]),
      ),
    );
  }

  Widget assetInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        height: height / 2.5,
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.all(Radius.circular(15.0)),
          color: notifier.getaddsubwalletgrey,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Padding(
              padding:
                  const EdgeInsets.symmetric(horizontal: 20.0, vertical: 35.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    'TROV Token',
                    style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: notifier.getbluecolor,
                        fontFamily: fontsemibold),
                  ),
                  SizedBox(
                    height: height / 50,
                  ),
                  Container(
                    width: width / 1.3,
                    child: Text(
                      'TROV token (TROV) is the utility token that powers the Trovotech ecosystem. TROV token is used to access discounts, voting rights, airdrops, NFTs and other community incentives. ',
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 15,
                        fontWeight: FontWeight.w400,
                        color: notifier.getbluecolor,
                        fontFamily: fontbody,
                      ),
                    ),
                  ),
                  SizedBox(
                    height: height / 50.0,
                  ),
                  Text(
                    'www.trovotech.io',
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w400,
                      color: notifier.getbluecolor,
                      fontFamily: fontbody,
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

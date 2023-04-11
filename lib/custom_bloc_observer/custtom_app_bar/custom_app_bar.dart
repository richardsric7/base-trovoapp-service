import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/storage/state.dart';

class CustomAppBarWithoutBanner extends PreferredSize {
  final BuildContext context;
  final double height;
  final String txt;
  final Color color;
  final Color titlecolor;

  CustomAppBarWithoutBanner(
    this.context,
    this.color,
    this.txt,
    this.titlecolor, {
    Key? key,
    required this.height,
  }) : super(
          key: key,
          child: AppBar(
            centerTitle: true,
            elevation: 0,
            backgroundColor: color,
            leading: GestureDetector(
              onTap: () {
                Navigator.of(context).pop();
              },
              child: Image.asset("assets/images/back.png", scale: 5),
            ),
            title: Text(
              txt,
              style: TextStyle(color: titlecolor, fontFamily: fontsemibold),
            ),
          ),
          preferredSize: Size.fromHeight(height),
        );
}

class CustomAppBar {
  final BuildContext context;
  final double height;
  final String txt;
  final Color color;
  final Color titlecolor;
  late DataProvider appState;

  CustomAppBar(
    this.context,
    this.color,
    this.txt,
    this.titlecolor, {
    required this.height,
  });

  PreferredSize getBar() {
    appState = Provider.of<DataProvider>(context, listen: true);
    return PreferredSize(
      child: AppBar(
        centerTitle: true,
        elevation: 0,
        backgroundColor: color,
        leading: GestureDetector(
          onTap: () {
            Navigator.of(context).pop();
          },
          child: Image.asset("assets/images/back.png", scale: 5),
        ),
        title: Text(
          txt,
          style: TextStyle(color: titlecolor, fontFamily: fontsemibold),
        ),
        actions: [
          if (appState.walletMode == "Testnet") ...[
            Visibility(
              visible: true,
              child: Container(
                color: Color(0xFFAA453E),
                width: 18,
                child: RotatedBox(
                  quarterTurns: 1,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 2),
                        child: Text(
                          'Testnet',
                          style: TextStyle(
                            fontFamily: fontsemibold,
                            color: wihitecolor,
                            fontSize: 11,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ]
        ],
      ),
      preferredSize: Size.fromHeight(height),
    );
  }
}

class CustomAppBarWithoutLeading {
  final BuildContext context;
  final double height;
  final String? txt;
  final Color? titlecolor;
  final Color color;
  final GlobalKey<ScaffoldState>? scaffoldKey; // Create a key
  final bool showMenu;
  late DataProvider appState;

  CustomAppBarWithoutLeading(
    this.context,
    this.color, {
    Key? key,
    this.scaffoldKey,
    this.showMenu = false,
    this.txt,
    this.titlecolor,
    required this.height,
  });

  PreferredSize getBar() {
    appState = Provider.of<DataProvider>(context, listen: true);
    return PreferredSize(
      child: AppBar(
        centerTitle: true,
        elevation: 0,
        backgroundColor: color,
        leading: showMenu
            ? IconButton(
                onPressed: () {
                  scaffoldKey?.currentState!.openDrawer();
                },
                icon: Icon(
                  Icons.menu,
                  size: 35,
                  color: titlecolor,
                ),
              )
            : null,
        title: Text(
          txt ?? '',
          style: TextStyle(color: titlecolor, fontFamily: fontsemibold),
        ),
        actions: [
          if (appState.walletMode == "Testnet") ...[
            Visibility(
              visible: true,
              child: Container(
                color: Color(0xFFAA453E),
                width: 18,
                child: RotatedBox(
                  quarterTurns: 1,
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 2),
                        child: Text(
                          'Testnet',
                          style: TextStyle(
                            fontFamily: fontsemibold,
                            color: wihitecolor,
                            fontSize: 11,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ]
        ],
      ),
      preferredSize: Size.fromHeight(height),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:trovo_wallet/Custom_BlocObserver/notifire_clor.dart';
import 'package:trovo_wallet/Models/BottomTabPage.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/home.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/settings.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/payment_history.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'bottom_pages/swap_assets.dart';
import 'bottom_pages/wallets.dart';

class BottomHome extends StatefulWidget {
  const BottomHome({Key? key}) : super(key: key);

  @override
  _BottomHomeState createState() => _BottomHomeState();
}

class _BottomHomeState extends State<BottomHome> {
  int _selectedIndex = 0;

  late ColorNotifier notifire;
  late DataProvider appState;
  bool isTapped = false;
  final _controller = PageController();

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    // set the page controller to appState
    appState.bottomTabPageController = _controller;
  }

  @override
  Widget build(BuildContext context) {
    notifire = Provider.of<ColorNotifier>(context, listen: true);
    appState = Provider.of<DataProvider>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return WillPopScope(
      onWillPop: () {
        Navigator.pop(context);
        return Future.value(false);
      },
      child: Scaffold(
        resizeToAvoidBottomInset: false,
        bottomNavigationBar: BottomNavigationBar(
          backgroundColor: notifire.getwihitecolor,
          unselectedItemColor: notifire.getgrey.withOpacity(.80),
          selectedLabelStyle: const TextStyle(fontFamily: 'Gilroy_Medium'),
          type: BottomNavigationBarType.fixed,
          selectedItemColor: notifire.getbluecolor,
          unselectedLabelStyle: const TextStyle(fontFamily: 'Gilroy_Medium'),
          currentIndex: _selectedIndex,
          showSelectedLabels: true,
          showUnselectedLabels: true,
          items: [
            BottomNavigationBarItem(
                backgroundColor: notifire.getwihitecolor,
                icon: AnimatedContainer(
                  duration: Duration(milliseconds: 2000),
                  curve: Curves.fastOutSlowIn,
                  child: Image.asset(
                    "assets/images/home.png",
                    color: _selectedIndex == ButtomTabPage.Dashboard.index
                        ? notifire.isDark
                            ? notifire.getbluecolor60
                            : notifire.getbluecolor
                        : notifire.getblck,
                    height: _selectedIndex == ButtomTabPage.Dashboard.index
                        ? height / 29
                        : height / 35,
                    fit: BoxFit.contain,
                  ),
                ),
                label: ''),
            BottomNavigationBarItem(
                backgroundColor: notifire.getwihitecolor,
                icon: AnimatedContainer(
                  duration: Duration(milliseconds: 2000),
                  curve: Curves.fastOutSlowIn,
                  child: Image.asset(
                    "assets/images/wallets.png",
                    color: _selectedIndex == ButtomTabPage.Wallets.index
                        ? notifire.isDark
                            ? notifire.getbluecolor60
                            : notifire.getbluecolor
                        : notifire.getblck,
                    height: _selectedIndex == ButtomTabPage.Wallets.index
                        ? height / 29
                        : height / 35,
                    fit: BoxFit.fitHeight,
                  ),
                ),
                label: ''),
            BottomNavigationBarItem(
                backgroundColor: notifire.getwihitecolor,
                icon: AnimatedContainer(
                  duration: Duration(milliseconds: 2000),
                  curve: Curves.fastOutSlowIn,
                  child: Image.asset("assets/images/history.png",
                      color: _selectedIndex ==
                              ButtomTabPage.TransactionHistory.index
                          ? notifire.isDark
                              ? notifire.getbluecolor60
                              : notifire.getbluecolor
                          : notifire.getblck,
                      height: _selectedIndex ==
                              ButtomTabPage.TransactionHistory.index
                          ? height / 29
                          : height / 35),
                ),
                label: ''),
            BottomNavigationBarItem(
                backgroundColor: notifire.getwihitecolor,
                icon: AnimatedContainer(
                  duration: Duration(milliseconds: 2000),
                  curve: Curves.fastOutSlowIn,
                  child: Image.asset("assets/images/swap.png",
                      color: _selectedIndex == ButtomTabPage.Swap.index
                          ? notifire.isDark
                              ? notifire.getbluecolor60
                              : notifire.getbluecolor
                          : notifire.getblck,
                      height: _selectedIndex == ButtomTabPage.Swap.index
                          ? height / 29
                          : height / 35),
                ),
                label: ''),
            BottomNavigationBarItem(
              backgroundColor: notifire.getwihitecolor,
              icon: AnimatedContainer(
                duration: Duration(milliseconds: 2000),
                curve: Curves.fastOutSlowIn,
                child: Image.asset(
                  "assets/images/settings.png",
                  color: _selectedIndex == ButtomTabPage.Settings.index
                      ? notifire.isDark
                          ? notifire.getbluecolor60
                          : notifire.getbluecolor
                      : notifire.getblck,
                  height: _selectedIndex == ButtomTabPage.Settings.index
                      ? height / 27
                      : height / 33,
                ),
              ),
              label: '',
            ),
          ],
          // onTap: (index) {
          //   changeTabMethod(index);
          //   isTapped = true;
          // },
          onTap: _onItemTapped,
        ),
        // body: Stack(
        //   children: [
        //     _buildOffstageNavigator(0),
        //     _buildOffstageNavigator(1),
        //     _buildOffstageNavigator(2),
        //     _buildOffstageNavigator(3),
        //     _buildOffstageNavigator(4),
        //   ],
        // ),
        body: PageView(
          controller: _controller,
          onPageChanged: (index) {
            changeTabMethod(index);
          },
          children: _pages,
        ),
      ),
    );
  }

  changeTabMethod(index) {
    setState(() {
      if (_selectedIndex != ButtomTabPage.TransactionHistory.index &&
          index == ButtomTabPage.TransactionHistory.index) {
        appState.getHistory();
      }
      _selectedIndex = index;
      appState.currentBottomTabIndex = _selectedIndex;
    });
  }

  void _onItemTapped(int index) {
    changeTabMethod(index);
    _controller.animateToPage(index,
        duration: const Duration(milliseconds: 500), curve: Curves.ease);
  }

  final List<Widget> _pages = [
    Home(),
    Wallets(),
    PaymentHistory(),
    SwapAssets(),
    Settings(),
  ];

  // Map<String, WidgetBuilder> _routeBuilders(BuildContext context, int index) {
  //   return {
  //     '/': (context) {
  //       return [
  //         Home(onButtonPressed: changeTabMethod),
  //         Wallets(),
  //         PaymentHistory(),
  //         SwapAssets(),
  //         Settings(),
  //       ].elementAt(index);
  //     },
  //   };
  // }
  // Widget _buildOffstageNavigator(int index) {
  //   var routeBuilders = _routeBuilders(context, index);

  //   return Offstage(
  //     offstage: _selectedIndex != index,
  //     child: Navigator(
  //       onGenerateRoute: (routeSettings) {
  //         return MaterialPageRoute(
  //           builder: (context) => routeBuilders[routeSettings.name]!(context),
  //         );
  //       },
  //     ),
  //   );
  // }
}

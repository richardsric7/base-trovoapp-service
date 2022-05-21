import 'package:flutter/material.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/home.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/profile.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stock_exchange_tabs/selectstocks.dart';
import 'package:gocrypto/bottom_bar/bottom_pages/stockexchange.dart';
import 'package:provider/provider.dart';

import '../utils/medeiaqury/medeiaqury.dart';

class BottomHome extends StatefulWidget {
  const BottomHome({Key? key}) : super(key: key);

  @override
  _BottomHomeState createState() => _BottomHomeState();
}

class _BottomHomeState extends State<BottomHome> {
  int _selectedIndex = 0;

  late ColorNotifier notifire;

  @override
  Widget build(BuildContext context) {
    notifire = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return WillPopScope(
      onWillPop: () {
        Navigator.pop(context);
        return Future.value(false);
      },
      child: Scaffold(
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
                icon: Image.asset("assets/images/Home-Filled.png",
                    color: _selectedIndex == 0
                        ? notifire.getbluecolor
                        : notifire.getgrey,
                    height: height / 35),
                label: ''),
            BottomNavigationBarItem(
                backgroundColor: notifire.getwihitecolor,
                icon: Image.asset("assets/images/Portfolio-outline.png",
                    color: _selectedIndex == 1
                        ? notifire.getbluecolor
                        : notifire.getgrey,
                    height: height / 35),
                label: ''),
            BottomNavigationBarItem(
                backgroundColor: notifire.getwihitecolor,
                icon: Image.asset("assets/images/exchange.png",
                    color: _selectedIndex == 2
                        ? notifire.getbluecolor
                        : notifire.getgrey,
                    height: height / 35),
                label: ''),
            BottomNavigationBarItem(
              backgroundColor: notifire.getwihitecolor,
              icon: Image.asset("assets/images/User-outline.png",
                  color: _selectedIndex == 3
                      ? notifire.getbluecolor
                      : notifire.getgrey,
                  height: height / 35),
              label: '',
            ),
          ],
          onTap: (index) {
            setState(() {
              _selectedIndex = index;
            });
          },
        ),
        body: Stack(
          children: [
            _buildOffstageNavigator(0),
            _buildOffstageNavigator(1),
            _buildOffstageNavigator(2),
            _buildOffstageNavigator(3),
          ],
        ),
      ),
    );
  }

  Map<String, WidgetBuilder> _routeBuilders(BuildContext context, int index) {
    return {
      '/': (context) {
        return [
          const Home(),
          const SelectStocks(),
          const StockExchange(),
          const Profile(),
        ].elementAt(index);
      },
    };
  }

  Widget _buildOffstageNavigator(int index) {
    var routeBuilders = _routeBuilders(context, index);

    return Offstage(
      offstage: _selectedIndex != index,
      child: Navigator(
        onGenerateRoute: (routeSettings) {
          return MaterialPageRoute(
            builder: (context) => routeBuilders[routeSettings.name]!(context),
          );
        },
      ),
    );
  }
}

import 'package:flutter/cupertino.dart';

import 'colors.dart';

class ColorNotifier with ChangeNotifier {
  bool isDark = false;

  set setIsDark(v) {
    isDark = v;
    notifyListeners();
  }

  get getIsDark => isDark;

  get getwihitecolor => isDark ? darkwihitecolor : wihitecolor;

  get getbluecolor => isDark ? blue : trovoblue;
  get getbluecolor90 => isDark ? darkblue90 : trovoblue90;
  get getbluecolor80 => isDark ? darkblue80 : trovoblue80;
  get getbluecolor70 => isDark ? darkblue70 : trovoblue70;
  get getbluecolor60 => isDark ? darkblue60 : trovoblue60;
  get getbluecolor50 => isDark ? darkblue50 : trovoblue50;

  // get getbluecolor => isDark ? blue : darkblue;

  get getgrey => isDark ? grey : darkgrey;

  get getdarkgrey => darkgrey;

  get getsplashgrey => splashgrey;

  get getgreencolor => green;

  get getblck => isDark ? blck : darkblck;

  get getpinauth => isDark ? pinauth : darkpinauth;

  get getInvestmentbluecolor => isDark ? investmentcolor : darkInvestmentcolor;

  get getconcirmstockbuycolor =>
      isDark ? darkconcirmstockbuycolor : concirmstockbuycolor;

  get getprefixicon => isDark ? prefixicon : darkprefixicon;

  get getidentyfiymethod => isDark ? darkidentyfiymethod : identyfiymethod;

  get getfavorites => isDark ? darkfavorites : favorites;

  get getaddsubwalletgrey => addsubwalletgrey;
}

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

  get gettilewihitecolor => isDark ? darktilewhitecolor : wihitecolor;

  get getbluewhitecolor => isDark ? wihitecolor : trovoblue;

  get getbluebackcolor => isDark ? trovoblue90 : trovoblue;

  get getbluecolor => isDark ? darkblue90 : trovoblue;
  get getbluecolor90 => isDark ? darkblue90 : trovoblue90;
  get getbluecolor80 => isDark ? darkblue80 : trovoblue80;
  get getbluecolor70 => isDark ? darkblue70 : trovoblue70;
  get getbluecolor60 => isDark ? darkblue60 : trovoblue60;
  get getbluecolor50 => isDark ? darkblue50 : trovoblue50;
  get getbottombarblue => isDark ? bottombarblue : bottombarblue;

  WalletTileColor get getstructuredbluecolor => WalletTileColor(
      backColor: isDark ? darkblue90 : trovoblue, foreColor: wihitecolor);
  WalletTileColor get getstructuredbluecolor90 => WalletTileColor(
      backColor: isDark ? darkblue90 : trovoblue90, foreColor: wihitecolor);
  WalletTileColor get getstructuredbluecolor80 => WalletTileColor(
      backColor: isDark ? darkblue80 : trovoblue80, foreColor: wihitecolor);
  WalletTileColor get getstructuredbluecolor70 => WalletTileColor(
      backColor: isDark ? darkblue70 : trovoblue70,
      foreColor: getbluewhitecolor);
  WalletTileColor get getstructuredbluecolor60 => WalletTileColor(
      backColor: isDark ? darkblue60 : trovoblue60,
      foreColor: getbluewhitecolor);
  WalletTileColor get getstructuredbluecolor50 => WalletTileColor(
      backColor: isDark ? darkblue50 : trovoblue50,
      foreColor: getbluewhitecolor);

  WalletTileColor get getpinkcolor =>
      WalletTileColor(backColor: colorPink, foreColor: forecolorblue);
  WalletTileColor get getpinkcolor90 =>
      WalletTileColor(backColor: colorPink90, foreColor: forecolorblue);
  WalletTileColor get getpinkcolor80 =>
      WalletTileColor(backColor: colorPink80, foreColor: forecolorblue);
  WalletTileColor get getpinkcolor70 =>
      WalletTileColor(backColor: colorPink70, foreColor: forecolorblue);
  WalletTileColor get getpinkcolor60 =>
      WalletTileColor(backColor: colorPink60, foreColor: forecolorblue);
  WalletTileColor get getpinkcolor50 =>
      WalletTileColor(backColor: colorPink50, foreColor: forecolorblue);

  WalletTileColor get getstructuredgreencolor =>
      WalletTileColor(backColor: colorGreen, foreColor: forecolorblue);
  WalletTileColor get getstructuredgreencolor90 =>
      WalletTileColor(backColor: colorGreen90, foreColor: forecolorblue);
  WalletTileColor get getstructuredgreencolor80 =>
      WalletTileColor(backColor: colorGreen80, foreColor: forecolorblue);
  WalletTileColor get getstructuredgreencolor70 =>
      WalletTileColor(backColor: colorGreen70, foreColor: forecolorblue);
  WalletTileColor get getstructuredgreencolor60 =>
      WalletTileColor(backColor: colorGreen60, foreColor: forecolorblue);
  WalletTileColor get getstructuredgreencolor50 =>
      WalletTileColor(backColor: colorGreen50, foreColor: forecolorblue);

  WalletTileColor get getorangecolor =>
      WalletTileColor(backColor: colorOrange, foreColor: forecolorblue);
  WalletTileColor get getorangecolor90 =>
      WalletTileColor(backColor: colorOrange90, foreColor: forecolorblue);
  WalletTileColor get getorangecolor80 =>
      WalletTileColor(backColor: colorOrange80, foreColor: forecolorblue);
  WalletTileColor get getorangecolor70 =>
      WalletTileColor(backColor: colorOrange70, foreColor: forecolorblue);
  WalletTileColor get getorangecolor60 =>
      WalletTileColor(backColor: colorOrange60, foreColor: forecolorblue);
  WalletTileColor get getorangecolor50 =>
      WalletTileColor(backColor: colorOrange50, foreColor: forecolorblue);

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
  get getpillbg => pillbg;
  get getplatinumcolor => platinum;
  get getdiamondcolor => diamond;
  get getgoldcolor => gold;
}

class WalletTileColor {
  Color backColor;
  Color foreColor;

  WalletTileColor({required this.backColor, required this.foreColor});
}

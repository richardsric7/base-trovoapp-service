import 'package:flutter/material.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallets.dart';

class WalletsListViewData {
  WalletView view;
  IconData actionIcon;
  String actionText;
  WalletsListViewData({
    required this.view,
    required this.actionIcon,
    required this.actionText,
  });
}

import 'package:flutter/material.dart';

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

enum WalletAction { import, createNew }

enum WalletView { listWallets, addSubWallet, confirmAddSubWallet }

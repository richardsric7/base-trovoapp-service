import 'package:flutter/material.dart';
import 'package:trovo_wallet/bottom_bar/bottom_pages/wallets.dart';

import '../screens/SharedAccess/shared_access.dart';

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

class GrantSharedAccessViewData {
  GrantSharedAccessView view;
  IconData actionIcon;
  String actionText;
  GrantSharedAccessViewData({
    required this.view,
    required this.actionIcon,
    required this.actionText,
  });
}

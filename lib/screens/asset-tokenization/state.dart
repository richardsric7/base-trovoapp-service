import 'package:flutter/material.dart';

class AssetTokenizationViewsState with ChangeNotifier {
  // view data is where all the data that a particular view needs
  // to do its work is. So when you want to pass any data from one view to
  // another, assign it to viewData and then get it back when you get
  // to the view. viewData is of type Map<String, dynamic>? where the string key
  // is the ViewPageConfig.key and the value is the data you want to pass to the
  // view. The value is of dynamic type so you can pass any data type you want.
  Map<String, dynamic>? viewData = {};
  dynamic tokenizationData = {};
}

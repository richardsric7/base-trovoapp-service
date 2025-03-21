import 'package:flutter/material.dart';
import 'package:trovo_app/router/ui_pages.dart';

enum PageState { none, addPage, addAll, addWidget, pop, replace, replaceAll }

class PageAction {
  PageState state;
  PageConfiguration? page;
  List<PageConfiguration>? pages;
  Widget? widget;

  PageAction({
    this.state = PageState.none,
    this.page = null,
    this.pages = null,
    this.widget = null,
  });
}

import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:gocrypto/Custom_BlocObserver/Custtom_app_bar/custtomappbar.dart';
import 'package:gocrypto/Custom_BlocObserver/notifire_clor.dart';
import 'package:gocrypto/bottom_bar/oder_tabs/historytabs.dart';
import 'package:gocrypto/bottom_bar/oder_tabs/ordertab.dart';
import 'package:provider/provider.dart';

import '../../utils/medeiaqury/medeiaqury.dart';

class Order extends StatefulWidget {
  const Order({Key? key}) : super(key: key);

  @override
  State<Order> createState() => _OrderState();
}

class _OrderState extends State<Order> with SingleTickerProviderStateMixin {
  late ColorNotifier notifier;
  TabController? _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: list.length, vsync: this);
  }

  List<Widget> list = [
    Row(
      children: const [
        Tab(
          icon: Icon(Icons.stacked_bar_chart, color: Colors.black),
        ),
        SizedBox(width: 4),
        Text("Order", style: TextStyle(color: Colors.black))
      ],
    ),
    Row(
      children: const [
        Tab(
          icon: Icon(Icons.local_activity, color: Colors.white),
        ),
        SizedBox(width: 4),
        Text("History", style: TextStyle(color: Colors.white))
      ],
    ),
  ];

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    return ScreenUtilInit(
      builder: () => Scaffold(
        backgroundColor: notifier.getwihitecolor,
        appBar: CustomAppBar(notifier.getwihitecolor, "Order", notifier.getblck,
            height: height / 15),
     
      ),
    );
  }
}

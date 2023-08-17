import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/widgets/price_points.dart';

class BarChartWidget extends StatefulWidget {
  const BarChartWidget({Key? key, required this.points}) : super(key: key);

  final List<PricePoint> points;

  @override
  State<BarChartWidget> createState() =>
      _BarChartWidgetState(points: this.points);
}

class _BarChartWidgetState extends State<BarChartWidget> {
  final List<PricePoint> points;
  late ColorNotifier notifier;

  _BarChartWidgetState({required this.points});

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    super.initState();
    getdarkmodepreviousstate();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return AspectRatio(
      aspectRatio: 1,
      child: BarChart(
        BarChartData(
          barGroups: _chartGroups(),
          borderData: FlBorderData(
              border: Border(
                  bottom: BorderSide(color: notifier.getbluewhitecolor),
                  left: BorderSide(color: notifier.getbluewhitecolor))),
          gridData: FlGridData(show: false),
          titlesData: FlTitlesData(
            bottomTitles: AxisTitles(sideTitles: _bottomTitles),
            leftTitles: AxisTitles(sideTitles: _leftTitles),
            topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
            rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
          ),
        ),
      ),
    );
  }

  List<BarChartGroupData> _chartGroups() {
    return points
        .map((point) => BarChartGroupData(x: point.x.toInt(), barRods: [
              BarChartRodData(toY: point.y, color: Colors.blueAccent),
              BarChartRodData(toY: point.y, color: Colors.green)
            ]))
        .toList();
  }

  SideTitles get _bottomTitles => SideTitles(
        showTitles: true,
        getTitlesWidget: (value, meta) {
          String text = '';
          switch (value.toInt()) {
            case 0:
              text = 'Sun';
              break;
            case 2:
              text = 'Mon';
              break;
            case 4:
              text = 'Tue';
              break;
            case 6:
              text = 'Wed';
              break;
            case 8:
              text = 'Thur';
              break;
            case 10:
              text = 'Fri';
              break;
            case 12:
              text = 'Sat';
              break;
          }

          return Text(
            text,
            style: TextStyle(
              fontSize: 13,
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          );
        },
      );
  SideTitles get _leftTitles => SideTitles(
        showTitles: true,
        getTitlesWidget: (value, meta) {
          String text = '';
          if (value > 0.1) text = '0.1';
          if (value > 0.2) text = '0.2';
          if (value > 0.3) text = '0.3';
          if (value > 0.4) text = '0.4';
          if (value > 0.5) text = '0.5';
          if (value > 0.6) text = '0.6';
          if (value > 0.7) text = '0.7';
          if (value > 0.8) text = '0.8';
          if (value > 0.9) text = '0.9';

          return Text(
            text,
            style: TextStyle(
              fontSize: 13,
              fontFamily: fontbody,
              color: notifier.getbluewhitecolor,
            ),
          );
        },
      );
}

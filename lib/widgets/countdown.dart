import 'dart:async';

import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';

class Countdown extends StatefulWidget {
  final DateTime startDate;
  const Countdown({super.key, required this.startDate});

  @override
  State<Countdown> createState() => _CountdownState();
}

class _CountdownState extends State<Countdown> {
  late Timer _timer;
  Widget countdown = Row();
  late ColorNotifier notifier;
  @override
  initState() {
    var diff = widget.startDate.difference(DateTime.now()).inDays;
    if (diff > 0) {
      countdown = Row(
        children: [
          for (var i = 0; i < diff.toString().length; i++) ...[
            Card(
              margin: EdgeInsets.zero,
              shadowColor: Colors.black,
              shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(5.0),
                  side: BorderSide(
                    color: notifier.getbluewhitecolor,
                    width: 1,
                  )),
              color: notifier.isDark
                  ? notifier.getbluecolor90
                  : notifier.getaddsubwalletgrey,
              child: Padding(
                padding: const EdgeInsets.all(5.0),
                child: Text(
                  countdown.toString()[i],
                  textAlign: TextAlign.center,
                  style: TextStyle(
                    fontSize: 12,
                    fontFamily: fontsemibold,
                    color: notifier.getbluewhitecolor,
                    overflow: TextOverflow.visible,
                  ),
                ),
              ),
            ),
            SizedBox(width: 2),
          ],
          SizedBox(width: 4),
          Text(
            'days left',
            style: TextStyle(
              fontSize: 12,
              fontFamily: fontbody,
              color: notifier.getblck,
            ),
          ),
        ],
      );
    } else {
      _timer = Timer.periodic(Duration(seconds: 1), (timer) {
        print('dfjslkd');
        setState(() {
          DateTime now = DateTime.now();
          var endTime = DateTime(
              now.year, now.month, now.day + 1, 0, 0, 0); // 12 AM tomorrow
          var remainingTime = endTime.difference(now);
          if (remainingTime.isNegative) remainingTime = Duration.zero;
          var time = formatTime(remainingTime).split(':');
          countdown = Row(
            children: [
              for (var i = 0; i < time.length; i++) ...[
                Card(
                  margin: EdgeInsets.zero,
                  shadowColor: Colors.black,
                  shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(5.0),
                      side: BorderSide(
                        color: notifier.getbluewhitecolor,
                        width: 1,
                      )),
                  color: notifier.isDark
                      ? notifier.getbluecolor90
                      : notifier.getaddsubwalletgrey,
                  child: Padding(
                    padding: const EdgeInsets.all(5.0),
                    child: Text(
                      time[i],
                      textAlign: TextAlign.center,
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontsemibold,
                        color: notifier.getbluewhitecolor,
                        overflow: TextOverflow.visible,
                      ),
                    ),
                  ),
                ),
                if (i + 1 < time.length) ...[
                  SizedBox(width: 2),
                  Text(
                    ':',
                    style: TextStyle(
                      fontSize: 12,
                      fontFamily: fontbody,
                      color: notifier.getblck,
                    ),
                  ),
                ],
                SizedBox(width: 2),
              ],
              SizedBox(width: 4),
              Text(
                'hours left',
                style: TextStyle(
                  fontSize: 12,
                  fontFamily: fontbody,
                  color: notifier.getblck,
                ),
              ),
            ],
          );
        });
      });
    }
    super.initState();
  }

  @override
  void dispose() {
    _timer.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    return countdown;
  }

  String formatTime(Duration duration) {
    int hours = duration.inHours;
    int minutes = (duration.inMinutes % 60);
    int seconds = (duration.inSeconds % 60);
    return '$hours:${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}';
  }
}

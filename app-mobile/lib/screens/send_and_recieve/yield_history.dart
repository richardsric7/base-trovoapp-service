import 'package:flutter/material.dart';
import 'package:trovo_app/screens/send_and_recieve/dividend_history.dart';

// YieldHistoryView is the interest payouts of a debt asset: the same
// payouts as dividends, paid by the payout engine.
class YieldHistoryView extends StatelessWidget {
  const YieldHistoryView({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) =>
      const DividendHistoryView(title: 'Interest History');
}

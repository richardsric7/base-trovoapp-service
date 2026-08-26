import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/storage/state.dart';

Future<void> bootstrapTestApp({
  required WidgetTester tester,
  required Widget child,
  required DataProvider dataProvider,
  required ColorNotifier colorNotifier,
  Locale locale = const Locale('en', 'US'),
}) async {
  TestWidgetsFlutterBinding.ensureInitialized();

  await tester.pumpWidget(
    MultiProvider(
      providers: [
        ChangeNotifierProvider<DataProvider>.value(value: dataProvider),
        ChangeNotifierProvider<ColorNotifier>.value(value: colorNotifier),
      ],
      child: MaterialApp(
        locale: locale,
        supportedLocales: const [Locale('en', 'US')],
        home: child,
      ),
    ),
  );

  await tester.pumpAndSettle();
}

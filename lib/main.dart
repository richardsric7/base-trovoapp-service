import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/Custom_BlocObserver/provider.dart';
import 'package:trovo_wallet/storage/store.dart';
import 'Custom_BlocObserver/notifire_clor.dart';
import 'config/app_settings.config.dart';
import 'screens/Splash_Screen/splashscreen.dart';

void main() async {
  await GetStorage.init();
  await Firebase.initializeApp();
  BlocOverrides.runZoned(
    () => runApp(const App()),
  );
}

class App extends StatefulWidget {
  const App({Key? key}) : super(key: key);

  @override
  State<App> createState() => _AppState();
}

class _AppState extends State<App> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => ColorNotifier()),
        ChangeNotifierProvider(create: (_) => DataProvider())
      ],
      child: const GetMaterialApp(
        debugShowCheckedModeBanner: false,
        home: SpashScreen(),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:provider/provider.dart';

import 'Custom_BlocObserver/notifire_clor.dart';
import 'screens/Splash_Screen/splashscreen.dart';

void main() async {
  await GetStorage.init();
  BlocOverrides.runZoned(
    () => runApp(const App()),
  );
}

class App extends StatelessWidget {
  const App({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => ColorNotifier()),
      ],
      child: const GetMaterialApp(
        debugShowCheckedModeBanner: false,
        home: SpashScreen(),
      ),
    );
  }
}

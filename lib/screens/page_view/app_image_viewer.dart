import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:webview_flutter/webview_flutter.dart';

class AppImageViewer extends StatefulWidget {
  const AppImageViewer({super.key});

  @override
  State<AppImageViewer> createState() => _AppImageViewerState();
}

class _AppImageViewerState extends State<AppImageViewer> {
  late DataProvider appState;
  bool isLoading = false;
  late ColorNotifier notifier;
  late WebViewController controller;

  @override
  void initState() {
    super.initState();
    appState = Provider.of<DataProvider>(context, listen: false);
    notifier = Provider.of<ColorNotifier>(context, listen: false);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: notifier.getwihitecolor,
      appBar: AppBar(
        backgroundColor: notifier.getwihitecolor,
        elevation: 0,
        toolbarHeight: 10,
        leading: Container(),
      ),
      body: Stack(
        children: [
          Center(
            child: Image.network(
              appState.initialUrl,
              loadingBuilder: (context, child, loadingProgress) {
                if (loadingProgress == null) return child;
                return const Center(child: CircularProgressIndicator());
              },
              errorBuilder: (context, error, stackTrace) =>
                  const Icon(Icons.error, size: 48),
            ),
          ),
        ],
      ),
    );
  }
}

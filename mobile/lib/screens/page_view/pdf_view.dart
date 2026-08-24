import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:flutter_pdfview/flutter_pdfview.dart';
import 'dart:io';
import 'package:http/http.dart' as http;
import 'package:path_provider/path_provider.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/utils/medeiaqury/medeiaqury.dart';

class PdfViewer extends StatefulWidget {
  const PdfViewer({super.key});

  @override
  State<PdfViewer> createState() => _PdfViewerState();
}

class _PdfViewerState extends State<PdfViewer> {
  String urlPDFPath = "";
  bool exists = true;
  int _totalPages = 0;
  int _currentPage = 0;
  bool pdfReady = false;
  late PDFViewController _pdfViewController;
  bool loaded = false;
  late DataProvider appState;
  late ColorNotifier notifier;

  Future<File> getFileFromUrl(String url, {name}) async {
    var fileName = 'testonline';
    if (name != null) {
      fileName = name;
    }
    try {
      var data = await http.get(Uri.parse(url));
      var bytes = data.bodyBytes;
      var dir = await getApplicationDocumentsDirectory();
      File file = File("${dir.path}/" + fileName + ".pdf");
      File urlFile = await file.writeAsBytes(bytes);
      return urlFile;
    } catch (e) {
      throw Exception("Error opening url file");
    }
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    getFileFromUrl(appState.pdfUrl!)
        .then(
          (value) => {
            setState(() {
              urlPDFPath = value.path;
              loaded = true;
              exists = true;
            }),
          },
        )
        .onError(
          (error, stackTrace) => {
            setState(() {
              exists = false;
            }),
          },
        );
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    if (loaded) {
      return Scaffold(
        appBar: CustomAppBar(
          context,
          notifier.getwihitecolor,
          "pdfviewer".tr(),
          notifier.getbluewhitecolor,
          height: height / 15,
        ).getBar(),
        body: PDFView(
          filePath: urlPDFPath,
          autoSpacing: true,
          enableSwipe: true,
          pageSnap: true,
          swipeHorizontal: true,
          nightMode: false,
          onError: (e) {
            //Show some error message or UI
          },
          onRender: (_pages) {
            setState(() {
              _totalPages = _pages ?? 0;
              pdfReady = true;
            });
          },
          onViewCreated: (PDFViewController vc) {
            setState(() {
              _pdfViewController = vc;
            });
          },
          onPageChanged: (int? page, int? total) {
            setState(() {
              _currentPage = page ?? 1;
            });
          },
          onPageError: (page, e) {},
        ),
        floatingActionButton: Row(
          mainAxisAlignment: MainAxisAlignment.end,
          children: <Widget>[
            IconButton(
              icon: Icon(Icons.chevron_left),
              iconSize: 50,
              color: Colors.black,
              onPressed: () {
                setState(() {
                  if (_currentPage > 0) {
                    _currentPage--;
                    _pdfViewController.setPage(_currentPage);
                  }
                });
              },
            ),
            Text(
              "${_currentPage + 1}/$_totalPages",
              style: TextStyle(color: Colors.black, fontSize: 20),
            ),
            IconButton(
              icon: Icon(Icons.chevron_right),
              iconSize: 50,
              color: Colors.black,
              onPressed: () {
                setState(() {
                  if (_currentPage < _totalPages - 1) {
                    _currentPage++;
                    _pdfViewController.setPage(_currentPage);
                  }
                });
              },
            ),
          ],
        ),
      );
    } else {
      if (exists) {
        //Replace with your loading UI
        return Scaffold(
          body: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                "pdfviewer".tr(),
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              Text(
                "loading".tr(),
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ],
          ),
        );
      } else {
        //Replace Error UI
        return Scaffold(
          body: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                "pdfviewer".tr(),
                notifier.getbluewhitecolor,
                height: height / 15,
              ).getBar(),
              Text(
                "pdfnotavailable".tr(),
                style: TextStyle(
                  fontSize: 15,
                  fontFamily: fontbody,
                  color: notifier.getbluewhitecolor,
                ),
              ),
            ],
          ),
        );
      }
    }
  }

  @override
  void dispose() {
    appState.pdfUrl = null;
    super.dispose();
  }
}

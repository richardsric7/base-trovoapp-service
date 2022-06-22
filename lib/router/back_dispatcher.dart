// /*
//  * Copyright (c) 2021 Razeware LLC
//  *
//  * Permission is hereby granted, free of charge, to any person obtaining a copy
//  * of this software and associated documentation files (the "Software"), to deal
//  * in the Software without restriction, including without limitation the rights
//  * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//  * copies of the Software, and to permit persons to whom the Software is
//  * furnished to do so, subject to the following conditions:
//  *
//  * The above copyright notice and this permission notice shall be included in
//  * all copies or substantial portions of the Software.
//  *
//  * Notwithstanding the foregoing, you may not use, copy, modify, merge, publish,
//  * distribute, sublicense, create a derivative work, and/or sell copies of the
//  * Software in any work that is designed, intended, or marketed for pedagogical or
//  * instructional purposes related to programming, coding, application development,
//  * or information technology.  Permission for such use, copying, modification,
//  * merger, publication, distribution, sublicensing, creation of derivative works,
//  * or sale is expressly withheld.
//  *
//  * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//  * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//  * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
//  * THE SOFTWARE.
//  */
// import 'package:flutter/material.dart';
// import 'router_delegate.dart';

// class FakeBantuPayBackButtonDispatcher extends RootBackButtonDispatcher {
//   final FakeBantuPayRouterDelegate _routerDelegate;

//   FakeBantuPayBackButtonDispatcher(this._routerDelegate) : super();
//   bool dialogOpen = false;

//   @override
//   Future<bool> didPopRoute() async {
//     if (_routerDelegate.pages.length <= 1) {
//       if (dialogOpen) {
//         Navigator.of(
//           _routerDelegate.navigatorKey.currentContext!,
//           rootNavigator: true,
//         ).pop(true);
//         dialogOpen = false;
//         return true;
//       }
//       dialogOpen = true;
//       return await _confirmAppExit() ?? true;
//     } else {
//       return _routerDelegate.popRoute();
//     }
//   }

//   Future<bool?> _confirmAppExit() {
//     return showDialog<bool>(
//         context: _routerDelegate.navigatorKey.currentContext!,
//         builder: (context) {
//           return AlertDialog(
//             shape: const RoundedRectangleBorder(
//                 borderRadius: BorderRadius.all(Radius.circular(25.0))),
//             actions: [
//               TextButton(
//                 child: Text(
//                   'Yes',
//                   style: TextStyle(fontSize: 16.0, color: Colors.orange[800]),
//                 ),
//                 onPressed: () => Navigator.of(
//                   context,
//                   rootNavigator: true,
//                 ).pop(false),
//               ),
//               TextButton(
//                   child: Text(
//                     'No',
//                     style: TextStyle(fontSize: 16.0, color: Colors.orange[800]),
//                   ),
//                   onPressed: () {
//                     Navigator.of(
//                       context,
//                       rootNavigator: true,
//                     ).pop(true);
//                     dialogOpen = false;
//                   }),
//             ],
//             content: Container(
//               decoration: const BoxDecoration(
//                 color: Colors.white,
//                 borderRadius: BorderRadius.all(
//                   Radius.circular(30),
//                 ),
//               ),
//               child: Text(
//                 'Are you sure you want to close this application?',
//                 style: TextStyle(
//                   fontSize: 18.0,
//                   fontWeight: FontWeight.w400,
//                 ),
//               ),
//             ),
//           );
//         });
//   }
// }

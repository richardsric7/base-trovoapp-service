import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/popups.dart';

// Network fee preference: the stablecoin the user's wallets pay network
// fees in (GET/PUT /v1/users/settings/gas-fee-asset(s) on app-backend).
// Every send is paid for by the wallet itself - in the chosen stablecoin
// when the wallet holds enough of it, otherwise in ETH.
Future<void> showNetworkFeeSetting(BuildContext context) async {
  final appState = Provider.of<DataProvider>(context, listen: false);
  final wallet = appState.primaryWallet;
  final signer = wallet.signer!;
  final address = wallet.address!;
  final secretKey = appState.secretKeys[0];

  final res = await makeGetRequest(
    uri: '/v1/users/settings/gas-fee-assets',
    signer: signer,
    address: address,
    secretKey: secretKey,
  );
  if (!context.mounted) return;
  if (res['statusCode'] != 200) {
    popup(
      context,
      title: "error".tr(),
      message: res['data']?['message'] ?? "somethingwentwrong".tr(),
    );
    return;
  }
  final assets = (res['data']['assets'] as List? ?? [])
      .map((a) => a['assetCode'] as String)
      .toList();
  String selected = res['data']['gasFeeAsset'] ?? '';

  await showModalBottomSheet(
    context: context,
    builder: (sheetContext) {
      final notifier = Provider.of<ColorNotifier>(sheetContext);
      return StatefulBuilder(
        builder: (sheetContext, setSheetState) {
          Widget option(String value, String label) => RadioListTile<String>(
            value: value,
            title: Text(
              label,
              style: TextStyle(
                color: notifier.getblck,
                fontFamily: fontsemibold,
                fontSize: 14,
              ),
            ),
          );
          return SafeArea(
            child: Padding(
              padding: const EdgeInsets.all(20),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    "networkfees".tr(),
                    style: TextStyle(
                      color: notifier.getblck,
                      fontFamily: fontsemibold,
                      fontSize: 16,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    "networkfeesdescription".tr(),
                    style: TextStyle(color: notifier.getblck, fontSize: 13),
                  ),
                  const SizedBox(height: 8),
                  RadioGroup<String>(
                    groupValue: selected,
                    onChanged: (v) => setSheetState(() => selected = v ?? ''),
                    child: Column(
                      children: [
                        option('', "payfeesineth".tr()),
                        ...assets.map((code) => option(code, code)),
                      ],
                    ),
                  ),
                  const SizedBox(height: 12),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: () async {
                        final saved = await makePutRequest(
                          uri: '/v1/users/settings/gas-fee-asset',
                          body: jsonEncode({'assetCode': selected}),
                          signer: signer,
                          address: address,
                          secretKey: secretKey,
                        );
                        if (!sheetContext.mounted) return;
                        Navigator.of(sheetContext).pop();
                        if (saved['statusCode'] != 200 && context.mounted) {
                          popup(
                            context,
                            title: "error".tr(),
                            message:
                                saved['data']?['message'] ??
                                "somethingwentwrong".tr(),
                          );
                        }
                      },
                      child: Text("save".tr()),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      );
    },
  );
}

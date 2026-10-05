import 'dart:convert';

import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:trovo_app/functions/trovo-sdk.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/storage/state.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';

// Account recovery changes each covered wallet in its own operation (turning
// recovery on or off, cancelling a recovery): the backend returns one
// transaction per wallet ("transactions", with their aliases in "wallets"),
// and the app signs each with the user's key.

// signRecoveryTransactions adds the user's signature of every transaction in
// data (one per wallet) as "transactionSignatures".
void signRecoveryTransactions(Map data, String secretKey) {
  final List txs = (data['transactions'] as List?) ?? [];
  data['transactionSignatures'] = txs
      .map(
        (tx) => TrovoWalletSDK().signBase64Txn(
          secretKey,
          tx.toString(),
          data['networkPassPhrase'] ?? '',
        ),
      )
      .toList();
}

// recoveryMessages are the messages to show before signing, starting with
// how many wallets the user signs for.
List recoveryMessages(Map data) {
  final messages = List.from((data['messages'] as List?) ?? []);
  final wallets = (data['wallets'] as List?) ?? [];
  if (wallets.length > 1) {
    messages.insert(0, "recoverysignperwallet".tr(args: [wallets.join(', ')]));
  }
  return messages;
}

// showMessagesThen shows each message in turn, then runs done.
void showMessagesThen(
  BuildContext context,
  List messages,
  void Function() done, [
  int index = 0,
]) {
  if (index >= messages.length) {
    done();
    return;
  }
  showResponseMessage(
    context,
    messages[index],
    () => showMessagesThen(context, messages, done, index + 1),
  );
}

// cancelAccountRecovery cancels a recovery of the signed-in user's account:
// it fetches the operations, shows what they do, has the user's key sign
// them and submits them.
Future<void> cancelAccountRecovery(
  BuildContext context,
  DataProvider appState, {
  required String signer,
  required String address,
  required String secretKey,
}) async {
  const uri = '/v1/users/account/recovery/cancel';
  try {
    showLoader(context);
    Map prepared = await makePostRequest(
      uri: uri,
      body: jsonEncode({}),
      signer: signer,
      secretKey: secretKey,
      address: address,
    );
    hideLoader(context);
    if (prepared['statusCode'] == 200) {
      popup(
        context,
        title: "success".tr(),
        message: "cancelaccountrecoverysuccess".tr(),
      );
      return;
    }
    if (prepared['statusCode'] == 404) {
      popup(
        context,
        title: "accountrecovery".tr(),
        message: "norecoveryinprogress".tr(),
      );
      return;
    }
    if (prepared['statusCode'] != 202) {
      popup(
        context,
        title: "error".tr(),
        message: prepared['data']['message'] ?? '',
      );
      return;
    }
    Map data = prepared['data'];
    showMessagesThen(context, recoveryMessages(data), () async {
      try {
        showLoader(context);
        signRecoveryTransactions(data, secretKey);
        Map submitted = await makePostRequest(
          uri: uri,
          body: jsonEncode(data),
          signer: signer,
          secretKey: secretKey,
          address: address,
        );
        hideLoader(context);
        if (submitted['statusCode'] == 200) {
          popup(
            context,
            title: "success".tr(),
            message: "cancelaccountrecoverysuccess".tr(),
          );
        } else {
          popup(
            context,
            title: "error".tr(),
            message: submitted['data']['message'] ?? '',
          );
        }
      } catch (e) {
        hideLoader(context);
        popup(context, title: "error".tr(), message: e.toString());
      }
    });
  } catch (e) {
    hideLoader(context);
    popup(context, title: "error".tr(), message: e.toString());
  }
}

import 'package:trovo_wallet/functions/trovo-sdk.dart';

import '../network/requests.dart';

class AuthFunctions {
  createAccount({account, userInfo}) async {
    try {
      var account = TrovoWalletSDK().createAccount();

      Map responseData = await makePostRequest(
          uri: '/v1/users',
          signer: account.publicKey,
          publicKey: account.publicKey,
          secretKey: account.secretKey,
          body: userInfo);
      print('$responseData');
    } catch (e) {
      print(e);
    }
  }
}

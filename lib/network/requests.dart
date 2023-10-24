import 'dart:async';
import 'dart:convert';
import 'dart:io';
import 'package:trovo_wallet/functions/helpers.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import '../functions/trovo-sdk.dart';

//final String trovoBaseUrl = 'https://api-alpha.dev.bantupay.org'; // alpha
//final String trovoBaseUrl = 'https://api-beta.dev.bantupay.org'; // beta
//final String trovoBaseUrl = 'https://api-prod.bantupay.org'; // pre launch production
//final String trovoBaseUrl = 'https://api.bantupay.org'; // Production

//final String walletApiBaseUrl =
//'https://api-wallet-alpha.dev.bantupay.org'; // alpha
//final String walletApiBaseUrl = 'https://api-wallet-beta.dev.bantupay.org'; // beta
//final String walletApiBaseUrl = 'https://api-wallet.bantupay.org'; // production

String getTrovoBaseURL() {
  // String trovoBaseURL;
  // if (GlobalConfiguration().getString("network") == "Development") {
  //   print('development...');
  //   trovoBaseURL = 'https://api.trovotechnologies.com';
  // } else {
  //   print('not development...');
  //   trovoBaseURL = 'https://api.trovotechnologies.com';
  // }
  // return trovoBaseURL;
  return 'https://api.trovotechnologies.com';
}

String getTrovoTestnetBaseURL() {
  // String trovoBaseURL;
  // if (GlobalConfiguration().getString("network") == "Development") {
  //   print('development...');
  //   trovoBaseURL = 'https://api.trovotechnologies.com';
  // } else {
  //   print('not development...');
  //   trovoBaseURL = 'https://api.trovotechnologies.com';
  // }
  // return trovoBaseURL;
  return 'https://apidev.trovotechnologies.com';
}

String getTrovoWalletApiBaseURL() {
  // String trovoWalletBaseURL;
  // if (GlobalConfiguration().getString("network") == "Development") {
  //   print('development...');
  //   trovoWalletBaseURL = 'https://api.trovotechnologies.com';
  // } else {
  //   print('not development...');
  //   trovoWalletBaseURL = 'https://api.trovotechnologies.com';
  // }
  // return trovoWalletBaseURL;
  return 'https://api.trovotechnologies.com';
}

Future<Map> makePostRequest({
  required String uri,
  required String body,
  required String signer,
  required String secretKey,
  required String publicKey,
}) async {
  Map<String, String> headers = await getRequestHeader(
    uri: uri,
    signer: signer,
    publicKey: publicKey,
    secretKey: secretKey,
  );

  //print('frist body: $body, pubkey: $publicKey, url: $baseUrlTest$uri');

  try {
    http.Response response = await http
        .post(Uri.parse(getTrovoTestnetBaseURL() + uri),
            body: body, headers: headers)
        .timeout(Duration(seconds: 60));
    // print("The statucode is: ${response.statusCode}");
    // print("The Response Body is: ${response.body}");
    return {
      'statusCode': response.statusCode,
      'data': json.decode(response.body)
    };
  } on SocketException catch (e) {
    print("The Catch Error on makePostRequest() Is: $e");
    // print('No Internet connection 😑');
    // return {'statusCode': 505, 'data': 'No Internet connection'};
    Map errorResponse = {
      "data": "$e",
      "error": "SocketException",
      "message": "No Internet connection"
    };
    return {'statusCode': 505, 'data': errorResponse};
  } on HttpException catch (e) {
    print("The Catch Error on makePostRequest() Is: $e");
    // print("Couldn't find the post 😱");
    // return {'statusCode': 505, 'data': "Couldn't find the post. Try again"};
    Map errorResponse = {
      "data": "$e",
      "error": "HttpException",
      "message": "Couldn't find the post"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on FormatException catch (e) {
    print("The Catch Error on makePostRequest() Is: $e");
    // print("Bad response format 👎");
    // return {'statusCode': 505, 'data': 'Bad response format'};

    Map errorResponse = {
      "data": "$e",
      "error": "FormatException",
      "message": "Bad response format"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on TimeoutException catch (e) {
    print("The Catch Error on makePostRequest() Is: $e");
    print("Request Time Out");
    // return {'statusCode': 505, 'data': 'Request Time Out'};
    Map errorResponse = {
      "data": "$e",
      "error": "TimeoutException",
      "message": "Request Time Out"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on Exception catch (e) {
    print("The Catch Error on makePostRequest() Is: $e");
    // return {'statusCode': 505, 'data': 'Request failed. Try again'};
    Map errorResponse = {
      "data": "$e",
      "error": "UnknownException",
      "message": "Sorry, something went wrong. Please try again"
    };

    return {'statusCode': 505, 'data': errorResponse};
  }
}

Future<Map> makeGetRequest({
  required String uri,
  required String signer,
  required String publicKey,
  required String secretKey,
}) async {
  Map<String, String> headers = await getRequestHeader(
    uri: uri,
    signer: signer,
    secretKey: secretKey,
    publicKey: publicKey,
  );

  try {
    http.Response response = await http
        .get(Uri.parse(getTrovoTestnetBaseURL() + uri), headers: headers)
        .timeout(Duration(seconds: 60));
    //  print("The statucode is: ${response.statusCode}");
    //  print("The Response Body is: ${response.body}");

    return {
      'statusCode': response.statusCode,
      'data': json.decode(response.body)
    };
  } on SocketException catch (e) {
    print("The Catch Error on makeGetRequest() Is: $e");
    // print('No Internet connection 😑');
    // return {'statusCode': 505, 'data': 'No Internet connection'};
    Map errorResponse = {
      "data": "$e",
      "error": "SocketException",
      "message": "No Internet connection"
    };
    return {'statusCode': 505, 'data': errorResponse};
  } on HttpException catch (e) {
    print("The Catch Error on makeGetRequest() Is: $e");
    // print("Couldn't find the post 😱");
    // return {'statusCode': 505, 'data': "Couldn't find the post. Try again"};
    Map errorResponse = {
      "data": "$e",
      "error": "HttpException",
      "message": "Couldn't find the post"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on FormatException catch (e) {
    print("The Catch Error on makeGetRequest() Is: $e");
    // print("Bad response format 👎");
    // return {'statusCode': 505, 'data': 'Bad response format'};

    Map errorResponse = {
      "data": "$e",
      "error": "FormatException",
      "message": "Bad response format"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on TimeoutException catch (e) {
    print("The Catch Error on makeGetRequest() Is: $e");
    print("Request Time Out");
    // return {'statusCode': 505, 'data': 'Request Time Out'};
    Map errorResponse = {
      "data": "$e",
      "error": "TimeoutException",
      "message": "Request Time Out"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on Exception catch (e) {
    print("The Catch Error on makeGetRequest() Is: $e");
    // return {'statusCode': 505, 'data': 'Request failed. Try again'};
    Map errorResponse = {
      "data": "$e",
      "error": "UnknownException",
      "message": "Unknown error. Try again"
    };

    return {'statusCode': 505, 'data': errorResponse};
  }
}

Future<Map> makePutRequest({
  required String uri,
  required String signer,
  required String body,
  required String secretKey,
  required String publicKey,
}) async {
  Map<String, String> headers = await getRequestHeader(
    uri: uri,
    signer: signer,
    secretKey: secretKey,
    publicKey: publicKey,
  );

  //print('frist body: $body, pubkey: $publicKey, url: $baseUrlTest$uri');

  try {
    http.Response response = await http
        .put(Uri.parse(getTrovoTestnetBaseURL() + uri),
            body: body, headers: headers)
        .timeout(Duration(seconds: 60));
    // print("The statucode is: ${response.statusCode}");
    // print("The Response Body is: ${response.body}");

    return {
      'statusCode': response.statusCode,
      'data': json.decode(response.body)
    };
  } on SocketException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // print('No Internet connection 😑');
    // return {'statusCode': 505, 'data': 'No Internet connection'};
    Map errorResponse = {
      "data": "$e",
      "error": "SocketException",
      "message": "No Internet connection"
    };
    return {'statusCode': 505, 'data': errorResponse};
  } on HttpException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // print("Couldn't find the post 😱");
    // return {'statusCode': 505, 'data': "Couldn't find the post. Try again"};
    Map errorResponse = {
      "data": "$e",
      "error": "HttpException",
      "message": "Couldn't find the post"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on FormatException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // print("Bad response format 👎");
    // return {'statusCode': 505, 'data': 'Bad response format'};

    Map errorResponse = {
      "data": "$e",
      "error": "FormatException",
      "message": "Bad response format"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on TimeoutException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    print("Request Time Out");
    // return {'statusCode': 505, 'data': 'Request Time Out'};
    Map errorResponse = {
      "data": "$e",
      "error": "TimeoutException",
      "message": "Request Time Out"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on Exception catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // return {'statusCode': 505, 'data': 'Request failed. Try again'};
    Map errorResponse = {
      "data": "$e",
      "error": "UnknownException",
      "message": "Unknown error. Try again"
    };

    return {'statusCode': 505, 'data': errorResponse};
  }
}

Future<Map> makeUnSecuredGetRequest(String path) async {
  try {
    http.Response response = await http
        .get(Uri.parse(getTrovoTestnetBaseURL() + path))
        .timeout(Duration(seconds: 60));
    //  print("The statucode is: ${response.statusCode}");
    //  print("The Response Body is: ${response.body}");

    return {
      'statusCode': response.statusCode,
      'data': json.decode(response.body)
    };
  } on SocketException catch (e) {
    print("The Catch Error on makeUnSecuredGetRequest() Is: $e");
    // print('No Internet connection 😑');
    // return {'statusCode': 505, 'data': 'No Internet connection'};
    Map errorResponse = {
      "data": "$e",
      "error": "SocketException",
      "message": "No Internet connection"
    };
    return {'statusCode': 505, 'data': errorResponse};
  } on HttpException catch (e) {
    print("The Catch Error on createBantuUser() Is: $e");
    // print("Couldn't find the post 😱");
    // return {'statusCode': 505, 'data': "Couldn't find the post. Try again"};
    Map errorResponse = {
      "data": "$e",
      "error": "HttpException",
      "message": "Couldn't find the post"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on FormatException catch (e) {
    print("The Catch Error on makeUnSecuredGetRequest() Is: $e");
    // print("Bad response format 👎");
    // return {'statusCode': 505, 'data': 'Bad response format'};

    Map errorResponse = {
      "data": "$e",
      "error": "FormatException",
      "message": "Bad response format"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on TimeoutException catch (e) {
    print("The Catch Error on makeUnSecuredGetRequest() Is: $e");
    print("Request Time Out");
    // return {'statusCode': 505, 'data': 'Request Time Out'};
    Map errorResponse = {
      "data": "$e",
      "error": "TimeoutException",
      "message": "Request Time Out"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on Exception catch (e) {
    print("The Catch Error on makeUnSecuredGetRequest() Is: $e");
    // return {'statusCode': 505, 'data': 'Request failed. Try again'};
    Map errorResponse = {
      "data": "$e",
      "error": "UnknownException",
      "message": "Unknown error. Try again"
    };

    return {'statusCode': 505, 'data': errorResponse};
  }
}

Future<Map> makePutRequestForMultipartFile({
  required String uri,
  required String signer,
  required String multipartFilePath,
  required String secretKey,
  required String publicKey,
}) async {
  Map<String, String> headers = await getRequestHeader(
    uri: uri,
    signer: signer,
    secretKey: secretKey,
    publicKey: publicKey,
  );

  //print('frist body: $body, pubkey: $publicKey, url: $baseUrlTest$uri');

  try {
    var request = await http.MultipartRequest(
        'PUT', Uri.parse(getTrovoTestnetBaseURL() + uri));
    request.headers.addAll(headers);
    request.files.add(await http.MultipartFile.fromPath(
        'profilePicture', multipartFilePath,
        contentType: MediaType('image', 'jpeg')));
    var response = await request.send();
    var responseString = await response.stream.bytesToString();
    print("The statucode is: ${response.statusCode}");
    print("The Response Body is: ${responseString}");

    return {
      'statusCode': response.statusCode,
      'data': responseString,
    };
  } on SocketException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // print('No Internet connection 😑');
    // return {'statusCode': 505, 'data': 'No Internet connection'};
    Map errorResponse = {
      "data": "$e",
      "error": "SocketException",
      "message": "No Internet connection"
    };
    return {'statusCode': 505, 'data': errorResponse};
  } on HttpException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // print("Couldn't find the post 😱");
    // return {'statusCode': 505, 'data': "Couldn't find the post. Try again"};
    Map errorResponse = {
      "data": "$e",
      "error": "HttpException",
      "message": "Couldn't find the post"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on FormatException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // print("Bad response format 👎");
    // return {'statusCode': 505, 'data': 'Bad response format'};

    Map errorResponse = {
      "data": "$e",
      "error": "FormatException",
      "message": "Bad response format"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on TimeoutException catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    print("Request Time Out");
    // return {'statusCode': 505, 'data': 'Request Time Out'};
    Map errorResponse = {
      "data": "$e",
      "error": "TimeoutException",
      "message": "Request Time Out"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on Exception catch (e) {
    print("The Catch Error on makePutRequest() Is: $e");
    // return {'statusCode': 505, 'data': 'Request failed. Try again'};
    Map errorResponse = {
      "data": "$e",
      "error": "UnknownException",
      "message": "Unknown error. Try again"
    };

    return {'statusCode': 505, 'data': errorResponse};
  }
}

Future<Map> makeDeleteRequest({
  required String uri,
  required String body,
  required String signer,
  required String secretKey,
  required String publicKey,
}) async {
  Map<String, String> headers = await getRequestHeader(
    uri: uri,
    signer: signer,
    publicKey: publicKey,
    secretKey: secretKey,
  );

  //print('frist body: $body, pubkey: $publicKey, url: $baseUrlTest$uri');

  try {
    http.Response response = await http
        .delete(Uri.parse(getTrovoTestnetBaseURL() + uri),
            body: body, headers: headers)
        .timeout(Duration(seconds: 60));
    // print("The statucode is: ${response.statusCode}");
    // print("The Response Body is: ${response.body}");
    return {
      'statusCode': response.statusCode,
      'data': json.decode(response.body)
    };
  } on SocketException catch (e) {
    print("The Catch Error on makeDeleteRequest() Is: $e");
    // print('No Internet connection 😑');
    // return {'statusCode': 505, 'data': 'No Internet connection'};
    Map errorResponse = {
      "data": "$e",
      "error": "SocketException",
      "message": "No Internet connection"
    };
    return {'statusCode': 505, 'data': errorResponse};
  } on HttpException catch (e) {
    print("The Catch Error on makeDeleteRequest() Is: $e");
    // print("Couldn't find the post 😱");
    // return {'statusCode': 505, 'data': "Couldn't find the post. Try again"};
    Map errorResponse = {
      "data": "$e",
      "error": "HttpException",
      "message": "Couldn't find the post"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on FormatException catch (e) {
    print("The Catch Error on makeDeleteRequest() Is: $e");
    // print("Bad response format 👎");
    // return {'statusCode': 505, 'data': 'Bad response format'};

    Map errorResponse = {
      "data": "$e",
      "error": "FormatException",
      "message": "Bad response format"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on TimeoutException catch (e) {
    print("The Catch Error on makeDeleteRequest() Is: $e");
    print("Request Time Out");
    // return {'statusCode': 505, 'data': 'Request Time Out'};
    Map errorResponse = {
      "data": "$e",
      "error": "TimeoutException",
      "message": "Request Time Out"
    };

    return {'statusCode': 505, 'data': errorResponse};
  } on Exception catch (e) {
    print("The Catch Error on makeDeleteRequest() Is: $e");
    // return {'statusCode': 505, 'data': 'Request failed. Try again'};
    Map errorResponse = {
      "data": "$e",
      "error": "UnknownException",
      "message": "Unknown error. Try again"
    };

    return {'statusCode': 505, 'data': errorResponse};
  }
}

getRequestHeader({uri, signer, publicKey, secretKey}) async {
  var deviceID = await getDeviceDetails();
  var appVersion = await getAppVersion();
  var ms = (new DateTime.now().toUtc()).millisecondsSinceEpoch;
  var serverTs = (ms / 1000).round().toString();
  var toSign = uri + signer + serverTs;
  print('toSign: $toSign');
  var signHTTP =
      TrovoWalletSDK().signHTTP(toSign: toSign, secretKey: secretKey);

  print('pubkey: $publicKey uri: $uri');

  Map<String, String> headers = {
    "X-TW-SIGNATURE": signHTTP,
    "X-TW-PUBLIC-KEY": publicKey,
    "X-TW-SIGNER": signer,
    "X-TW-DEVICE-ID": deviceID,
    "X-TW-APP-VERSION": appVersion,
    "X-TW-TIMESTAMP": serverTs
  };

  print('The printed header is $headers');

  return headers;
}

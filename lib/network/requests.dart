import 'dart:async';
import 'dart:convert';
import 'dart:developer';
import 'dart:io';
import 'package:file_picker/file_picker.dart';
import 'package:mime/mime.dart';
import 'package:trovo_wallet/functions/helpers.dart';
import 'package:http/http.dart' as http;
import 'package:http_parser/http_parser.dart';
import 'package:trovo_wallet/storage/store.dart';
import '../functions/trovo-sdk.dart';

Future<String> getTrovoAppBaseURL() async {
  String trovoBaseURL;
  if (await StoreData().storeGetData('walletMode') == "Testnet") {
    trovoBaseURL = 'https://apidev.trovotechnologies.com';
    print('testnet... $trovoBaseURL');
  } else {
    trovoBaseURL = 'https://api.trovotechnologies.com';
    print('mainnet... $trovoBaseURL');
  }
  return trovoBaseURL;
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
        .post(Uri.parse(await getTrovoAppBaseURL() + uri),
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
        .get(Uri.parse(await getTrovoAppBaseURL() + uri), headers: headers)
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
        .put(Uri.parse(await getTrovoAppBaseURL() + uri),
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
        .get(Uri.parse(await getTrovoAppBaseURL() + path))
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

Future<Map> makeUnSecuredPostRequest({
  required String uri,
  required String body,
  Map<String, String>? headers,
}) async {
  try {
    http.Response response = await http
        .post(
          Uri.parse(await getTrovoAppBaseURL() + uri),
          body: body,
          headers: headers,
        )
        .timeout(Duration(seconds: 60));
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
        'PUT', Uri.parse(await getTrovoAppBaseURL() + uri));
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

Future<Map> makePutRequestForMultipartDocumentUpload({
  required String uri,
  required String signer,
  required String secretKey,
  required String publicKey,
  required PlatformFile file,
  required String tokenizedAssetId,
  required String documentType,
  required String documentTitle,
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
        'PUT', Uri.parse(await getTrovoAppBaseURL() + uri));
    Map<String, String> map = {
      "tokenizedAssetID": tokenizedAssetId,
      "documentType": documentType.toString(),
      "documentTitle": documentTitle,
    };
    // print('mappppppppppp $map');
    request.headers.addAll(headers);
    request.fields.addAll(map);
    final mimeType = lookupMimeType(file.path!);
    final contentType = mimeType != null ? MediaType.parse(mimeType) : null;
    request.files.add(await http.MultipartFile.fromPath(
      'documentFile',
      file.path!,
      contentType: contentType,
    ));
    var response = await request.send();
    var responseString = await response.stream.bytesToString();
    // print("The statucode is: ${response.statusCode}");
    // print("The Response Body is: ${responseString}");

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

Future<Map> makePutRequestForFeeRecieptUpload({
  required String uri,
  required String signer,
  required String secretKey,
  required String publicKey,
  required PlatformFile file,
  required String tokenizationFeePaymentMethodID,
  required String transactionReference,
}) async {
  Map<String, String> headers = await getRequestHeader(
    uri: uri,
    signer: signer,
    secretKey: secretKey,
    publicKey: publicKey,
  );

  try {
    var request = await http.MultipartRequest(
        'PUT', Uri.parse(await getTrovoAppBaseURL() + uri));
    Map<String, String> map = {
      "tokenizationFeePaymentMethodID": tokenizationFeePaymentMethodID,
      "transactionReference": transactionReference.toString(),
    };
    print('mappppppppppp $map');
    request.headers.addAll(headers);
    request.fields.addAll(map);
    final mimeType = lookupMimeType(file.path!);
    final contentType = mimeType != null ? MediaType.parse(mimeType) : null;
    request.files.add(await http.MultipartFile.fromPath(
      'documentFile',
      file.path!,
      contentType: contentType,
    ));
    inspect(request);
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
        .delete(Uri.parse(await getTrovoAppBaseURL() + uri),
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

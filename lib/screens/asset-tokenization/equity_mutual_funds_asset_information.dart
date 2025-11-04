import 'dart:async';
import 'dart:convert';
import 'dart:developer';
import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter_screenutil/flutter_screenutil.dart';
import 'package:trovo_app/custom_bloc_observer/button/custtom_button.dart';
import 'package:trovo_app/custom_bloc_observer/colors.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_app_bar/custom_app_bar.dart';
import 'package:trovo_app/custom_bloc_observer/custtom_textfild/consttom_textfild.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';
import 'package:trovo_app/custom_bloc_observer/notifire_clor.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:trovo_app/network/requests.dart';
import 'package:trovo_app/widgets/loader.dart';
import 'package:trovo_app/widgets/popups.dart';
import 'package:trovo_app/widgets/utilities.dart';
import '../../storage/state.dart';
import '../../utils/medeiaqury/medeiaqury.dart';

class EquityMutualFundsAssetInformationView extends StatefulWidget {
  const EquityMutualFundsAssetInformationView({Key? key}) : super(key: key);

  @override
  State<EquityMutualFundsAssetInformationView> createState() =>
      _EquityMutualFundsAssetInformationView();
}

class _EquityMutualFundsAssetInformationView
    extends State<EquityMutualFundsAssetInformationView>
    with TickerProviderStateMixin {
  late ColorNotifier notifier;
  final _formKey = GlobalKey<FormState>();
  late DataProvider appState;
  late Future<dynamic> formJsonFuture;
  bool formHasError = false;
  late dynamic data = {};
  var formData = {};

  final valueOfAssetController = TextEditingController();
  final miscCostOfAssetController = TextEditingController();
  final percentageFromPromotersController = TextEditingController();
  final assetOwnerRetainedOrContributedValueController =
      TextEditingController();
  final percentValueOfInsuranceController = TextEditingController();

  getdarkmodepreviousstate() async {
    final prefs = await SharedPreferences.getInstance();
    bool? previusstate = prefs.getBool("setIsDark");
    if (previusstate == null) {
      notifier.setIsDark = false;
    } else {
      notifier.setIsDark = previusstate;
    }
  }

  @override
  void initState() {
    appState = Provider.of<DataProvider>(context, listen: false);
    data = appState.viewData;
    inspect(data);
    super.initState();
    getdarkmodepreviousstate();
    formJsonFuture = fetchFormJson(appState.viewData?['assetType']);
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: false,
      backgroundColor: notifier.getwihitecolor,
      body: Form(
        key: _formKey,
        child: SingleChildScrollView(
          child: Column(
            children: [
              CustomAppBar(
                context,
                notifier.getwihitecolor,
                'Asset Information',
                notifier.getbluewhitecolor,
                height: height / 15,
                // fontSize: 13,
              ).getBar(),
              notifyAdditionalInfo(),
              SizedBox(height: height / 50),
              FutureBuilder<dynamic>(
                future: formJsonFuture,
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return SizedBox(
                      height: height / 2,
                      child: Center(
                        child: CircularProgressIndicator(
                          backgroundColor: notifier.getbluecolor,
                          valueColor: new AlwaysStoppedAnimation<Color>(
                            notifier.getgreencolor,
                          ),
                          strokeWidth: 3.0,
                        ),
                      ),
                    );
                  } else if (snapshot.connectionState == ConnectionState.done) {
                    if (snapshot.hasError) {
                      return Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: SizedBox(
                          height: height / 6,
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(
                                "somethingwentwrong".tr(),
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                  fontSize: 16,
                                  color: notifier.getbluewhitecolor,
                                  fontFamily: fontbody,
                                ),
                              ),
                              ElevatedButton(
                                onPressed: () {
                                  setState(() {
                                    formJsonFuture = fetchFormJson(
                                      appState.viewData?['assetType'],
                                    );
                                  });
                                },
                                style: ButtonStyle(
                                  backgroundColor:
                                      WidgetStateProperty.all<Color>(
                                        notifier.getbluecolor!,
                                      ),
                                  foregroundColor:
                                      WidgetStateProperty.all<Color>(
                                        notifier.getwihitecolor,
                                      ),
                                ),
                                child: Text(
                                  "retry".tr(),
                                  style: TextStyle(fontFamily: fontsemibold),
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    } else if (snapshot.hasData) {
                      formData = snapshot.data!;
                      inspect(formData);
                      return Column(
                        children: [
                          for (var section in formData.entries) ...[
                            Padding(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 20.0,
                              ),
                              child: Row(
                                mainAxisAlignment:
                                    MainAxisAlignment.spaceBetween,
                                children: [
                                  SizedBox(
                                    width: 295,
                                    child: Text.rich(
                                      TextSpan(
                                        text: section.value['sectionName'],
                                        style: TextStyle(
                                          fontSize: 18,
                                          fontFamily: fontsemibold,
                                          color: notifier.getbluewhitecolor,
                                        ),
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            SizedBox(height: height / 70),
                            for (var item in section.value['body'].entries) ...[
                              if (checkConditions(item)) ...[
                                Padding(
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 20.0,
                                  ),
                                  child: Row(
                                    mainAxisAlignment:
                                        MainAxisAlignment.spaceBetween,
                                    children: [
                                      SizedBox(
                                        width: 295,
                                        child: Text.rich(
                                          TextSpan(
                                            text: item.value['label'],
                                            style: TextStyle(
                                              fontSize: 12,
                                              fontFamily: fontsemibold,
                                              color: notifier.getbluewhitecolor,
                                            ),
                                            children: [
                                              if (item.value['required'] ==
                                                  true) ...[
                                                TextSpan(
                                                  text: ' *',
                                                  style: TextStyle(
                                                    fontWeight: FontWeight.bold,
                                                    color: Colors.red,
                                                  ),
                                                ),
                                              ],
                                            ],
                                          ),
                                        ),
                                      ),
                                      if (item.value['widgetType'] ==
                                          'list') ...[
                                        TextButton(
                                          onPressed: () {
                                            addMilestone(
                                              label:
                                                  'Add ${item.value['label']}',
                                              placeholder:
                                                  item.value['placeholderText'],
                                              onDone: (value) {
                                                setState(() {
                                                  var val =
                                                      item.value['value'] ==
                                                              null ||
                                                          item.value['value']
                                                              .toString()
                                                              .isEmpty
                                                      ? []
                                                      : item.value['value']
                                                            .toString()
                                                            .split(',');
                                                  val.add(value);
                                                  formData[section
                                                      .key]['body'][item
                                                      .key]['value'] = val.join(
                                                    ',',
                                                  );
                                                });
                                              },
                                            );
                                          },
                                          style: ButtonStyle(
                                            padding: WidgetStatePropertyAll(
                                              EdgeInsets.all(7),
                                            ),
                                            tapTargetSize: MaterialTapTargetSize
                                                .shrinkWrap,
                                            minimumSize: WidgetStatePropertyAll(
                                              Size.zero,
                                            ),
                                          ),
                                          child: Icon(
                                            Icons.add_circle,
                                            size: 20,
                                          ),
                                        ),
                                      ],
                                    ],
                                  ),
                                ),
                                SizedBox(height: height / 70),
                                getFormElement(item),
                              ],
                            ],
                          ],
                        ],
                      );
                    }
                  }
                  return Text(
                    '',
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      fontFamily: fontsemibold,
                    ),
                  );
                },
              ),

              Button(
                "saveandcontinuee".tr(),
                notifier.getbluecolor,
                wihitecolor,
                onTap: () {
                  var form = _formKey.currentState;
                  if (form!.validate()) {
                    form.save();
                    submitForm();
                  }
                },
              ),
              SizedBox(height: height / 10),
              Padding(
                padding: EdgeInsets.only(
                  bottom: MediaQuery.of(context).viewInsets.bottom,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Future<dynamic> fetchFormJson(String formId) async {
    var uri = '/v1/forms/$formId';

    Map responseData = await makeGetRequest(
      uri: Uri.encodeFull(uri),
      signer: appState.primaryWallet.signer!,
      secretKey: appState.secretKeys[0], // the primary wallet secret key
      publicKey: appState.primaryWallet.signer!,
    );
    // inspect(responseData['data']);
    if (responseData['statusCode'] == 200) {
      try {
        var formStr = responseData['data']['formString'];
        var first = jsonDecode(formStr);
        var jsonObj = first is String ? jsonDecode(first) : first;
        return jsonObj;
      } catch (e) {
        print(e);
        // inspect(e);
      }
    }

    return Future.error('Error fetching form data');
  }

  bool checkConditions(MapEntry<dynamic, dynamic> item) {
    bool value = true;
    var condition = item.value['when'];
    if (condition != null) {
      var conditionValue = condition?['value'].toString().toLowerCase();
      var inputValue =
          formData[condition['section']
                  .toString()]['body'][condition['targetName']]?['value']
              .toString()
              .toLowerCase();
      print("Input value: $inputValue  ========> $conditionValue");

      if (condition?['type'] == 'hasvalue') {
        if (inputValue == null || inputValue.isEmpty) {
          print('========== we are here 1');
          return false;
        }
      } else if (inputValue != conditionValue) {
        print('========== we are here 2');
        return false;
      }
    }
    return value;
  }

  Timer? _timer;

  Widget getFormElement(MapEntry<dynamic, dynamic> item) {
    if (item.key == 'fundInstrumentType') {
      print(
        'widget key ${item.key} ========> widget value ${item.value['value']} ====> ${data[item.key]}',
      );
    }
    switch (item.value['widgetType']) {
      case 'list':
        if (item.value['value'] == null) {
          item.value['value'] = data[item.key].toString();
        }

        var listItems =
            item.value['value'] == null ||
                item.value['value'].toString().isEmpty
            ? []
            : item.value['value'].toString().split(',');
        return Column(
          children: [
            if (item.value['required'] == true && listItems.isEmpty) ...[
              Row(
                children: [
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 20.0),
                    child: Text(
                      'This field is required',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: Colors.red,
                      ),
                    ),
                  ),
                ],
              ),
            ] else ...[
              for (var it in listItems) ...[
                SizedBox(height: height / 70),
                listItem(
                  item: it,
                  onDelete: (val) {
                    setState(() {
                      listItems.removeWhere((i) => i == val);
                      item.value['value'] = listItems.join(',');
                    });
                  },
                ),
              ],
            ],
            SizedBox(height: height / 50),
          ],
        );
      case 'datetime':
        if (item.value['value'] == null) {
          var parsedDate = DateTime.parse(data[item.key].toString());
          item.value['value'] = parsedDate.year == DateTime(0001).year
              ? null
              : parsedDate;
        }

        return Column(
          children: [
            ButtonOutlined(
              item.value['value'] != null
                  ? DateFormat(
                      'MMMM dd, yyyy',
                    ).format(DateTime.parse(item.value['value'].toString()))
                  : "Select date",
              notifier.getwihitecolor,
              notifier.getgrey,
              borderColor: notifier.getgrey,
              width: 320,
              height: 50.sp,
              onTap: () {
                showDatePicker(
                  context: context,
                  initialDate: DateTime.now(),
                  firstDate: DateTime.fromMicrosecondsSinceEpoch(1000),
                  lastDate: DateTime.now().add(Duration(days: 730)),
                ).then(
                  (value) => {
                    setState(() {
                      item.value['value'] = value;
                    }),
                  },
                );
              },
            ),
            SizedBox(height: height / 50),
          ],
        );
      case 'dropdown':
        List<DropdownMenuItem<String>> options = [];
        var dropdownOptions = item.value['options'];
        for (var i = 0; i < dropdownOptions.length; i++) {
          options.add(
            DropdownMenuItem(
              child: Text(dropdownOptions[i], overflow: TextOverflow.ellipsis),
              value: dropdownOptions[i],
            ),
          );
        }

        if (item.value['value'] == null) {
          item.value['value'] = data[item.key].toString().nullIfEmpty();
        }

        return Column(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10.0),
              child: dropdown(
                (value) {
                  setState(() {
                    item.value['value'] = value;
                  });
                },
                options,
                item.value['value'],
                item.value['placeholderText'],
                context,
                null,
                validator: (value) {
                  if (item.value['required'] != true) return null;

                  if (item.value['value'] == null) {
                    return "Please select an item";
                  }
                  return null;
                },
              ),
            ),
            SizedBox(height: height / 50),
          ],
        );
      case 'longtext':
        if (item.value['value'] == null) {
          item.value['value'] = data[item.key]?.toString() ?? '';
        }

        return Column(
          children: [
            Row(
              children: [
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 20.0),
                  child: multilineInput(
                    item.value['placeholderText'],
                    notifier.getbluecolor,
                    notifier.getgrey,
                    notifier.getblck,
                    notifier.getgrey,
                    100.sp,
                    width / 1.12,
                    initialValue: item.value['value'] ?? '',
                    validator: (value) {
                      if (item.value['required'] != true) return null;

                      if (value.isEmpty) {
                        return "fieldcannotbeempty".tr();
                      }
                      return null;
                    },
                    onChanged: (value) {
                      _timer?.cancel();
                      _timer = Timer(Duration(milliseconds: 200), () {
                        // setState(() {
                        item.value['value'] = value;
                        // });
                      });
                    },
                    onSaved: (value) {
                      setState(() {
                        item.value['value'] = value;
                      });
                    },
                    minLines: 3,
                    maxLines: null,
                    keyboardtype: TextInputType.multiline,
                  ),
                ),
              ],
            ),
            SizedBox(height: height / 50),
          ],
        );
      default:
        if (item.value['type'] == 'double') {
          if (item.value['value'] == null) {
            item.value['value'] = double.parse(data[item.key].toString());
          }

          if (item.value['controller'] == null) {
            item.value['controller'] = TextEditingController(
              text: item.value['value'] == 0
                  ? ''
                  : formatNumberForInput(item.value['value']),
            );
          }

          return Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: CustomTextFormField.textField(
                  null,
                  notifier.getbluecolor,
                  null,
                  notifier.getgrey,
                  null,
                  notifier.getblck,
                  notifier.getgrey,
                  hintText: item.value['placeholderText'],
                  85,
                  300.sp,
                  validator: (value) {
                    if (item.value['required'] != true) return null;

                    if (value.isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                  onChanged: (value) {
                    _timer?.cancel();
                    _timer = Timer(Duration(milliseconds: 200), () {
                      // setState(() {
                      item.value['value'] = double.parse(value);
                      // });
                    });
                  },
                  onSaved: (value) {
                    setState(() {
                      item.value['value'] = double.parse(value);
                    });
                  },
                  autoFormatNumber: true,
                  controller: item.value['controller'],
                  keyboardtype: TextInputType.numberWithOptions(decimal: true),
                ),
              ),
            ],
          );
        }

        if (item.value['type'] == 'int') {
          if (item.value['value'] == null) {
            item.value['value'] = int.parse(data[item.key].toString());
          }

          if (item.value['controller'] == null) {
            item.value['controller'] = TextEditingController(
              text: item.value['value'] == 0
                  ? ''
                  : item.value['value'].toString(),
            );
          }

          return Row(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 20.0),
                child: CustomTextFormField.textField(
                  null,
                  notifier.getbluecolor,
                  null,
                  notifier.getgrey,
                  null,
                  notifier.getblck,
                  notifier.getgrey,
                  hintText: item.value['placeholderText'],
                  85.sp,
                  300.sp,
                  validator: (value) {
                    if (item.value['required'] != true) return null;

                    if (value.isEmpty) {
                      return "fieldcannotbeempty".tr();
                    }
                    return null;
                  },
                  onChanged: (value) {
                    _timer?.cancel();
                    _timer = Timer(Duration(milliseconds: 200), () {
                      // setState(() {
                      item.value['value'] = int.tryParse(value) ?? 0;
                      // });
                    });
                  },
                  onSaved: (value) {
                    setState(() {
                      item.value['value'] = int.tryParse(value) ?? 0;
                    });
                  },
                  autoFormatNumber: true,
                  controller: item.value['controller'],
                  keyboardtype: TextInputType.numberWithOptions(decimal: true),
                ),
              ),
            ],
          );
        }

        if (item.value['value'] == null) {
          item.value['value'] = data[item.key] == null
              ? ''
              : data[item.key].toString();
        }

        if (item.key == 'fundInstrumentType') {
          print(
            'widget key ${item.key} ========> widget value ${item.value['value']} ====> ${data[item.key]}',
          );
        }

        return Row(
          children: [
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 20.0),
              child: CustomTextFormField.textField(
                null,
                notifier.getbluecolor,
                null,
                notifier.getgrey,
                null,
                notifier.getblck,
                notifier.getgrey,
                hintText: item.value['placeholderText'],
                85,
                300.sp,
                initialValue: item.value['value'].toString(),
                onChanged: (value) {
                  _timer?.cancel();
                  _timer = Timer(Duration(milliseconds: 200), () {
                    // setState(() {
                    item.value['value'] = value;
                    // });
                  });
                },
                onSaved: (value) {
                  setState(() {
                    item.value['value'] = value;
                  });
                },
                validator: (value) {
                  if (item.value['required'] != true) return null;

                  if (value.isEmpty) {
                    return "fieldcannotbeempty".tr();
                  }
                  return null;
                },
              ),
            ),
          ],
        );
    }
  }

  void submitForm() async {
    try {
      showLoader(context);
      var newData = {...data as Map};

      for (var section in formData.entries) {
        for (var item in section.value['body'].entries) {
          if (item.value['widgetType'] == 'datetime') {
            print('doing date time ===> ${item.value['value']}');
            newData[item.key] = DateFormat(
              "yyyy-MM-ddTHH:mm:ss.SSSSSS'Z'",
            ).format(item.value['value']!.toUtc());
            print(
              'got here ===> ${newData[item.key]} ===> ${DateFormat("yyyy-MM-ddTHH:mm:ss.SSSSSS'Z'").format(item.value['value']!.toUtc())}',
            );
          } else {
            newData[item.key] = item.value['value'];
          }
        }
      }

      // Logger().i(newData);
      // inspect(newData);

      String requestBody = jsonEncode(newData);
      Map responseData = await makePostRequest(
        uri: '/v1/tokenization',
        body: requestBody,
        signer: appState.primaryWallet.signer!,
        secretKey: appState.secretKeys[0], // the primary wallet secret key
        publicKey: appState.primaryWallet.signer!,
      );

      if (responseData['statusCode'] == 200) {
        await refreshCurrentTokenizationInfo(appState);
        Navigator.of(context).pop();
      } else {
        popup(
          context,
          title: "error".tr(),
          message: responseData['data']['message'],
        );
      }
      hideLoader(context);
    } catch (e) {
      hideLoader(context);
      popup(context, title: "error".tr(), message: e.toString());
    }
  }

  Widget notifyAdditionalInfo() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(20, 10, 20, 3),
      child: Container(
        child: Card(
          shadowColor: Colors.black,
          color: notifier.getaddsubwalletgrey,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(5.0),
            side: BorderSide(color: notifier.getaddsubwalletgrey, width: 1),
          ),
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 5.0, vertical: 10),
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.center,
              children: [
                Container(
                  child: Card(
                    shadowColor: Colors.black,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20.0),
                      side: BorderSide(
                        color: notifier.getbluewhitecolor,
                        width: 1,
                      ),
                    ),
                    color: notifier.getaddsubwalletgrey,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 5.0),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Text(
                            "i".tr(),
                            textAlign: TextAlign.start,
                            style: TextStyle(
                              fontSize: 12,
                              fontFamily: fontbody,
                              color: notifier.getbluewhitecolor,
                              overflow: TextOverflow.visible,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
                SizedBox(
                  width: 300,
                  child: Text.rich(
                    TextSpan(
                      text: 'Fields marked in (',
                      style: TextStyle(
                        fontSize: 12,
                        fontFamily: fontbody,
                        color: notifier.getbluewhitecolor,
                      ),
                      children: [
                        TextSpan(
                          text: '*',
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: Colors.red,
                          ),
                        ),
                        TextSpan(text: ') are required'),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget CheckItem(
    String name,
    void Function()? onClick, {
    required Color backColor,
    required Color foreColor,
    required Color borderColor,
    double? fontSize = 15,
  }) {
    return Padding(
      padding: const EdgeInsets.all(3.0),
      child: ElevatedButton(
        onPressed: onClick,
        style: ButtonStyle(
          overlayColor: WidgetStateProperty.all<Color>(notifier.getsplashgrey),
          elevation: WidgetStateProperty.all<double>(0),
          backgroundColor: WidgetStateProperty.all<Color>(backColor),
          foregroundColor: WidgetStateProperty.all<Color>(
            notifier.getwihitecolor,
          ),
          side: WidgetStateProperty.all(
            BorderSide(color: borderColor, width: 1, style: BorderStyle.solid),
          ),
          shape: WidgetStateProperty.all<RoundedRectangleBorder>(
            const RoundedRectangleBorder(
              borderRadius: BorderRadius.all(Radius.circular(10)),
            ),
          ),
        ),
        child: Wrap(
          alignment: WrapAlignment.center,
          crossAxisAlignment: WrapCrossAlignment.center,
          children: [
            Text(
              name,
              textAlign: TextAlign.center,
              softWrap: true,
              style: TextStyle(
                color: foreColor,
                fontFamily: fontbody,
                fontSize: fontSize,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget CheckboxItem({
    required String label,
    required bool value,
    required void Function(bool?) onChanged,
    String? Function(Object?)? validator,
  }) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Transform.scale(
          scale: 1,
          child: FormField(
            builder: (state) {
              return SizedBox(
                width: 24,
                height: 24,
                child: Checkbox(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.all(Radius.circular(5.sp)),
                  ),
                  activeColor: notifier.isDark
                      ? notifier.getbluecolor50
                      : notifier.getbluecolor90,
                  side: BorderSide(
                    color: notifier.isDark
                        ? notifier.getbluecolor50
                        : notifier.getbluecolor90,
                  ),
                  value: value,
                  onChanged: onChanged,
                ),
              );
            },
            validator: validator,
          ),
        ),
        SizedBox(width: 10),
        Container(
          width: width / 1.4,
          child: Text(
            label,
            overflow: TextOverflow.visible,
            style: TextStyle(
              fontSize: 15,
              color: formHasError && !value && validator != null
                  ? Colors.red
                  : notifier.getbluewhitecolor,
              fontFamily: fontbody,
            ),
          ),
        ),
      ],
    );
  }

  void addMilestone({
    required String label,
    required String placeholder,
    required void Function(String val) onDone,
  }) {
    String val = '';
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: notifier.getwihitecolor,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(
          top: Radius.circular(16),
        ), // Rounded top corners
      ),
      builder: (BuildContext context) {
        return DraggableScrollableSheet(
          initialChildSize: 0.75,
          minChildSize: 0.25,
          expand: false,
          builder: (context, scrollController) {
            return SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  SizedBox(height: 10),
                  Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 14.0),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        SizedBox(
                          width: 250.sp,
                          child: Text(
                            label,
                            overflow: TextOverflow.visible,
                            style: TextStyle(
                              fontSize: 15,
                              fontFamily: fontsemibold,
                              color: notifier.getbluewhitecolor,
                            ),
                          ),
                        ),
                        TextButton(
                          onPressed: () => Navigator.of(context).pop(),
                          child: Icon(Icons.cancel_outlined),
                          style: ButtonStyle(
                            padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                            tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                            minimumSize: WidgetStatePropertyAll(Size.zero),
                          ),
                        ),
                      ],
                    ),
                  ),
                  SizedBox(height: 20),
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 20.0),
                        child: CustomTextFormField.textField(
                          placeholder,
                          notifier.getbluecolor,
                          null,
                          notifier.getgrey,
                          null,
                          notifier.getblck,
                          notifier.getgrey,
                          70.sp,
                          width / 1.12,
                          validator: (value) {
                            if (value.isEmpty) {
                              return "fieldcannotbeempty".tr();
                            }
                            return null;
                          },
                          onChanged: (value) {
                            val = value;
                          },
                        ),
                      ),
                    ],
                  ),
                  SizedBox(height: 40),
                  Button(
                    "Done",
                    notifier.getbluecolor,
                    wihitecolor,
                    onTap: () {
                      Navigator.of(context).pop();
                      onDone(val);
                    },
                  ),
                ],
              ),
            );
          },
        );
      },
    );
  }

  Widget listItem({
    required String item,
    required void Function(String item) onDelete,
  }) {
    return Row(
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 20),
          child: Container(
            width: width / 1.12,
            decoration: BoxDecoration(
              borderRadius: const BorderRadius.all(Radius.circular(10)),
              color: notifier.getaddsubwalletgrey,
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 10),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        item,
                        style: TextStyle(
                          fontSize: 12,
                          fontFamily: fontsemibold,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                      TextButton(
                        onPressed: () {
                          onDelete(item);
                        },
                        style: ButtonStyle(
                          padding: WidgetStatePropertyAll(EdgeInsets.all(7)),
                          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
                          minimumSize: WidgetStatePropertyAll(Size.zero),
                        ),
                        child: Icon(
                          CupertinoIcons.trash,
                          size: 15,
                          color: notifier.getbluewhitecolor,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}

enum Comparison {
  equals,
  greater,
  lesser,
  greaterOrEquals,
  lesserOrEquals,
  notEquals,
  isEmpty,
  isNotEmpty,
}

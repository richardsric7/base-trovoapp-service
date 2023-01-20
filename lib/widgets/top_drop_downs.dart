import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_wallet/custom_bloc_observer/colors.dart';
import 'package:trovo_wallet/custom_bloc_observer/fonts.dart';
import 'package:trovo_wallet/custom_bloc_observer/notifire_clor.dart';
import 'package:trovo_wallet/models/asset.dart';
import 'package:trovo_wallet/storage/state.dart';
import 'package:trovo_wallet/utils/medeiaqury/medeiaqury.dart';
import 'package:trovo_wallet/widgets/utilities.dart';

class TopDropdowns extends StatefulWidget {
  void Function(String newValue) onWalletChanged;
  void Function(String newValue)? onAssetChanged;
  List<Asset> claimedAssets;
  String? selectedWallet;
  String? selectedAsset;

  TopDropdowns({
    Key? key,
    required this.onWalletChanged,
    this.onAssetChanged,
    required this.claimedAssets,
    this.selectedAsset,
    required this.selectedWallet,
  }) : super(key: key);

  @override
  State<TopDropdowns> createState() => _TopDropdownsState();
}

class _TopDropdownsState extends State<TopDropdowns> {
  late ColorNotifier notifier;
  late DataProvider appState;
  var assetBalances;
  late String? selectedWallet;
  late String? selectedAsset;

  List<DropdownMenuItem<String>> assetDropdownItems(bool isSelected) {
    List<DropdownMenuItem<String>> menuItems = [];
    for (var asset in widget.claimedAssets) {
      menuItems.add(DropdownMenuItem(
          child: Text(
            isSelected
                ? truncate(
                    getAssetCode(asset.assetCode),
                    length: 3,
                  )
                : getAssetCode(asset.assetCode),
            overflow: TextOverflow.visible,
          ),
          value:
              '${getAssetCode(asset.assetCode)}|${getAssetIssuer(asset.assetIssuer)}'));
    }
    return menuItems;
  }

  List<DropdownMenuItem<String>> walletDropdownItems(bool isSelected) {
    var walletsList = <DropdownMenuItem<String>>[];
    appState.allWallets.forEach((key, value) {
      walletsList.add(
        DropdownMenuItem(
          child: Row(
            children: [
              Container(
                constraints: isSelected
                    ? BoxConstraints(maxWidth: width / 4)
                    : BoxConstraints(maxWidth: width / 2.5),
                child: Text(
                  value['alias'],
                  overflow:
                      isSelected ? TextOverflow.ellipsis : TextOverflow.visible,
                ),
              ),
              if (value['sharedAccessEnabled'] == 1) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.people_outline,
                  size: 17,
                  color: notifier.getbluecolor,
                )
              ],
              if (!isSelected && key == selectedWallet) ...[
                SizedBox(
                  width: 2,
                ),
                Icon(
                  Icons.check,
                  size: 18,
                  color: notifier.getbluecolor,
                )
              ],
            ],
          ),
          value: key,
        ),
      );
    });

    return walletsList;
  }

  @override
  Widget build(BuildContext context) {
    notifier = Provider.of<ColorNotifier>(context, listen: true);
    height = MediaQuery.of(context).size.height;
    width = MediaQuery.of(context).size.width;
    appState = Provider.of<DataProvider>(context, listen: true);
    selectedAsset = widget.selectedAsset;
    selectedWallet = widget.selectedWallet;

    return Container(
      width: width / 1.2,
      child: Row(
        children: [
          Expanded(
              flex: 5,
              child: dropdown(
                (newValue) {
                  widget.onWalletChanged(newValue.toString());
                },
                walletDropdownItems(false),
                selectedWallet.toString().isEmpty ? null : selectedWallet,
                null,
                context,
                (context) {
                  return walletDropdownItems(true);
                },
              )),
          if (widget.claimedAssets != null) ...[
            SizedBox(
              width: width / 40,
            ),
            Expanded(
              flex: 3,
              child: DropdownButtonFormField(
                dropdownColor: notifier.isDark
                    ? darktilewhitecolor
                    : notifier.getaddsubwalletgrey,
                decoration: InputDecoration(
                  contentPadding:
                      EdgeInsets.symmetric(vertical: 0, horizontal: 20),
                  enabledBorder: OutlineInputBorder(
                    borderSide: BorderSide.none,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  border: OutlineInputBorder(
                    borderSide: BorderSide.none,
                    borderRadius: BorderRadius.circular(10),
                  ),
                  filled: true,
                  fillColor: notifier.isDark
                      ? darktilewhitecolor
                      : notifier.getaddsubwalletgrey,
                ),
                value: selectedAsset,
                icon: Icon(
                  Icons.keyboard_arrow_down_rounded,
                  color: notifier.getbluewhitecolor,
                ),
                style: TextStyle(
                  color: notifier.getbluewhitecolor,
                  fontSize: 15,
                  fontFamily: fontsemibold,
                ),
                onChanged: (newValue) {
                  widget.onAssetChanged!(newValue.toString());
                },
                borderRadius: BorderRadius.all(
                  Radius.circular(15),
                ),
                items: assetDropdownItems(false),
                selectedItemBuilder: (context) {
                  return assetDropdownItems(true);
                },
              ),
            ),
            SizedBox(
              width: width / 20,
            ),
          ]
        ],
      ),
    );
  }
}

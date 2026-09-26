import 'package:easy_localization/easy_localization.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:trovo_app/custom_bloc_observer/p2p_theme.dart';
import 'package:trovo_app/storage/state.dart';

// P2PWalletPayoutSelect is the BUY-offer wallet dropdown (Plan: "an
// ordered list of their wallets showing the wallet alias only"). Sourced
// from the same wallet list already loaded into DataProvider for the
// wallet switcher, rather than a second P2P-scoped endpoint.
class P2PWalletPayoutSelect extends StatelessWidget {
  final String? value;
  final ValueChanged<String?> onChanged;
  const P2PWalletPayoutSelect({super.key, required this.value, required this.onChanged});

  @override
  Widget build(BuildContext context) {
    final appState = Provider.of<DataProvider>(context, listen: false);
    final wallets = appState.userInfo?.wallets ?? [];
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text('p2preceiveassetinto'.tr(), style: const TextStyle(fontWeight: FontWeight.w600)),
        const SizedBox(height: P2PTheme.space1),
        Container(
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(P2PTheme.space2),
          ),
          child: DropdownButtonFormField<String>(
            value: value,
            isExpanded: true,
            decoration: InputDecoration(
              hintText: 'p2pselectwallet'.tr(),
              filled: true,
              fillColor: Colors.white,
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(P2PTheme.space2),
                borderSide: BorderSide.none,
              ),
            ),
            items: wallets
                .map((w) => DropdownMenuItem(value: w.address, child: Text(w.alias ?? '', overflow: TextOverflow.ellipsis)))
                .toList(),
            onChanged: onChanged,
          ),
        ),
      ],
    );
  }
}

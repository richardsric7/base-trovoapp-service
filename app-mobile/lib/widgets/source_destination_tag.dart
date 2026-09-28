import 'package:flutter/material.dart';
import 'package:trovo_app/custom_bloc_observer/fonts.dart';

// A small rounded pill labelling one side (source or destination) of a
// payment-history row: which asset/network the funds left with, and which
// they arrived as. For a plain payment the two sides are the same, so
// SourceDestinationTags collapses to a single combined tag; a row where
// they differ is a swap, shown as two distinct From/To tags.
class SourceDestinationTag extends StatelessWidget {
  final String label;
  final Color backgroundColor;
  final Color textColor;

  const SourceDestinationTag({
    super.key,
    required this.label,
    required this.backgroundColor,
    required this.textColor,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: BorderRadius.circular(20),
      ),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: textColor,
          fontFamily: fontbody,
        ),
      ),
    );
  }
}

class SourceDestinationTags extends StatelessWidget {
  final bool isSwap;
  final String sourceAssetCode;
  final String sourceNetwork;
  final String destinationAssetCode;
  final String destinationNetwork;

  const SourceDestinationTags({
    super.key,
    required this.isSwap,
    required this.sourceAssetCode,
    required this.sourceNetwork,
    required this.destinationAssetCode,
    required this.destinationNetwork,
  });

  @override
  Widget build(BuildContext context) {
    if (!isSwap) {
      return Wrap(
        spacing: 6,
        runSpacing: 6,
        children: [
          SourceDestinationTag(
            label: '$destinationAssetCode · $destinationNetwork',
            backgroundColor: const Color(0xFFE3EBFA),
            textColor: const Color(0xFF20427C),
          ),
        ],
      );
    }

    return Wrap(
      spacing: 6,
      runSpacing: 6,
      children: [
        SourceDestinationTag(
          label: 'From $sourceAssetCode · $sourceNetwork',
          backgroundColor: const Color(0xFFFCE9CF),
          textColor: const Color(0xFF8A5A1F),
        ),
        SourceDestinationTag(
          label: 'To $destinationAssetCode · $destinationNetwork',
          backgroundColor: const Color(0xFFDCF4E6),
          textColor: const Color(0xFF1E7A46),
        ),
      ],
    );
  }
}

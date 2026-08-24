import 'package:flutter/services.dart';

class WordLengthLimitingTextInputFormatter extends TextInputFormatter {
  final int maxWords;

  WordLengthLimitingTextInputFormatter(this.maxWords);

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    // Trim leading/trailing spaces
    final trimmedText = newValue.text.trim();

    if (trimmedText.isEmpty) {
      return newValue;
    }

    // Split by whitespace (handles multiple spaces)
    final words = trimmedText.split(RegExp(r'\s+'));

    if (words.length <= maxWords) {
      return newValue;
    }

    // If word limit exceeded, revert to old value
    return oldValue;
  }
}

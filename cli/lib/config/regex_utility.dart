part of wdtk_config;

/// Regular expression for parsing parameters
final RegExp _paramExpr = RegExp("[^{}]+(?=})");

/// Callback for parsing a group
typedef ParseGroupFn = String Function(String);

/// A utility class for parsing parameters
class ParameterParser {
  final String inString;

  ParameterParser({required this.inString});

  /// Parse a string
  String parse(ParseGroupFn parseCallback) {
    String result = "";
    int previousStart = 0;

    final parameters = _paramExpr.allMatches(inString);
    if (parameters.isEmpty) {
      return inString;
    }

    for (var parameter in parameters) {
      // We are only capturing the inner part of the ${} so gotta remove the starting '${' and the trailing '}'
      if (previousStart != parameter.start - 2) {
        // Don't add empty
        result += inString.substring(previousStart, parameter.start - 2);
      }
      previousStart = parameter.end + 1;

      result += parseCallback(parameter.group(0)!);
    }

    // Add remaining to the end
    result += inString.substring(previousStart);

    return result;
  }
}

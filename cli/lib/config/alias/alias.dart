part of wdtk_config;


/// Base class for aliases
abstract class Alias {
  /// Get computed value of an alias
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false});

  /// Compute the value of the alias
  void _compute(WDTKConfig config);

  /// Executes an alias from the given argument string
  static String _executeOnString(Alias toExecute, String argsString) {
    return _executeOnList(toExecute, argsString.split(","));
  }

  /// Execute an alias on a list of arguments
  static String _executeOnList(Alias toExecute, List<String> argsList) {
    Map<String, String> arguments = {};
    bool forward = false;

    for (int i = 0; i < argsList.length; i++) {
      if (argsList[i].trim() == "...") {
        // Forward remaining arguments
        if (i != argsList.length - 1) {
          Logger.warning(
              "Forwarding operator should be last, any subsequent entries will be ignored");
        }

        forward = true;
        break;
      }

      var entry = argsList[i].split(":");
      arguments[entry[0].trim()] = entry[1].trim();
    }

    return toExecute.getComputedValue(
        args: arguments, forwardRemaining: forward);
  }
}

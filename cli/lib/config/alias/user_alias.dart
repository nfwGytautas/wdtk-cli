part of wdtk_config;

/// Class for holding and processing a user alias
class UserAlias implements Alias {
  final String name;
  final String value;

  // Computation
  String _computeValue = "";

  UserAlias({required this.name, required this.value});

  /// Get compute value of an alias
  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    return ParameterParser(inString: _computeValue).parse((parameter) {
      var name = parameter.substring(1);

      if (args != null) {
        String? value = args[name];
        if (value != null) {
          return value;
        }
      }

      if (forwardRemaining) {
        return "\${#$name}";
      }

      Logger.error("Alias '${this.name}' is missing an argument '$name'");
      return "@@ERROR@@";
    });
  }

  /// Compute the value of the alias
  @override
  void _compute(WDTKConfig config) {
    if (_computeValue.isNotEmpty) {
      // Already computed
      return;
    }
    Logger.verbose("[Alias] Computing '$name'");

    _computeValue = ParameterParser(inString: value).parse((parameter) {
      var param = parameter.split(",");

      var name = param[0];

      if (name[0] == '#') {
        // This is a function argument, leave it
        return "\${#${name.substring(1)}}";
      } else {
        // This is a reference to another alias
        Alias? alias = config.getAlias(name);

        if (alias == null) {
          Logger.error("Failed to find dependant alias $name");
          return "@@ERROR@@";
        }

        alias._compute(config);

        // Arguments
        return Alias._executeOnList(alias, param.sublist(1));
      }
    });

    Logger.verbose("[Alias] Computed '$name' to '$_computeValue'");
  }
}

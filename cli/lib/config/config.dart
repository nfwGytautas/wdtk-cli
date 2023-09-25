part of wdtk_config;

/// Class containing information read from 'wdtk.yaml' file
class WDTKConfig {
  late final String package;
  late final String name;

  late final Map<String, Alias>? aliases;
  late final Map<String, Deployment> deployments;
  late final FrontendConfig? frontend;
  late final Map<String, Service> services;

  Deployment? _selectedDeployment;
  Service? _currentService;
  String? _currentFrontend;

  /// Get all services with the specified type
  List<Service> getServicesOfType(String type) {
    List<Service> result = List.empty(growable: true);
    for (var service in services.values) {
      if (service.source.getType() == type) {
        result.add(service);
      }
    }
    return result;
  }

  /// Compute the in string and replace all aliases where possible
  String getStringValue(String input) {
    return ParameterParser(inString: input).parse((parameter) {
      return _aliasValue(parameter);
    });
  }

  /// Get an alias value from the specified alias string
  String getAliasValue(String aliasString) {
    if (!aliasString.startsWith("\${")) {
      Logger.error("Not an alias string $aliasString");
      return "@@ERROR@@";
    }

    // Remove the brace ${ and remaining }
    var cleanString = aliasString.substring(2, aliasString.length - 1);
    return _aliasValue(cleanString);
  }

  /// Get alias value
  String _aliasValue(String args) {
    var entries = args.split(",");

    Alias? alias = getAlias(entries[0]);

    if (alias == null) {
      Logger.error("Alias '${entries[0]}' is not defined");
      return "@@ERROR@@";
    }

    return Alias._executeOnList(alias, entries.sublist(1));
  }

  /// Get an alias with the specified name
  Alias? getAlias(String name) {
    if (name.startsWith("::")) {
      // Deployment alias
      return _selectedDeployment!.getAlias(name.substring(2));
    }

    if (aliases == null) {
      return null;
    }

    if (!aliases!.containsKey(name)) {
      return null;
    }

    return aliases![name];
  }

  /// Select the current deployment
  void selectDeployment(String name) {
    // TODO: Refactor to be able to use without selecting
    if (!deployments.containsKey(name)) {
      print("Unknown deployment $name");
      return;
    }

    Logger.verbose("Selecting deployment $name");

    aliases!["__DEPLOYMENT__"] = DeploymentAlias(deploymentName: name);
    aliases!["__SERVICE__"] = ServiceAlias(config: this);
    aliases!["__DEPLOYMENT_DIR__"] = DeploymentDirAlias(config: this);

    _selectedDeployment = deployments[name];
    _computeAliases();
  }

  /// Select a service for processing
  void selectService(String name) {
    // TODO: Refactor to be able to use without selecting
    if (!services.containsKey(name)) {
      print("Unknown service $name");
      return;
    }

    Logger.verbose("Selecting service $name");

    _currentService = services[name];
    _currentFrontend = null;
  }

  /// Select a frontend entry for processing
  void selectFrontend(String name) {
    if (frontend == null) {
      // Do nothing
      return;
    }

    Logger.verbose("Selecting frontend $name");

    _currentService = null;
    _currentFrontend = name;
  }

  /// Computes all aliases
  void _computeAliases() {
    if (aliases == null) {
      return;
    }

    // Compute globals
    for (var alias in aliases!.values) {
      alias._compute(this);
    }

    // Compute deployment aliases
    _selectedDeployment!._computeAliases(this);
  }

  /// Load the configuration file
  static WDTKConfig? load({String? path}) {
    File file = File(path ?? 'wdtk.yaml');
    if (file.existsSync() == true) {
      WDTKConfig result = WDTKConfig();
      result._load(loadYaml(file.readAsStringSync()));
      return result;
    }

    return null;
  }

  /// Load data from a map into the config
  void _load(Map data) {
    // Simple
    package = data["package"];
    name = data["name"];

    // Aliases
    aliases = <String, Alias>{};
    if (data.containsKey("aliases")) {
      for (var entry in data["aliases"].entries) {
        aliases![entry.key] =
            UserAlias(name: entry.key, value: entry.value.toString());
      }
    }

    // Deployments
    deployments = <String, Deployment>{};
    for (var entry in data["deployments"]) {
      var dep = Deployment(entry);
      deployments[dep.name] = dep;
    }

    // Frontend
    if (data.containsKey("frontend")) {
      frontend = FrontendConfig(data["frontend"]);
    }

    // Services
    services = {};
    for (var entry in data["services"]) {
      var ser = Service(entry);
      services[ser.name] = ser;
    }

    _addDefaultAliases();
  }

  /// Adds predefined default aliases for wdtk
  void _addDefaultAliases() {
    aliases!["__HOME__"] = HomeAlias();
    aliases!["__PACKAGE__"] = PackageAlias(packageName: name);
    aliases!["__PACKAGE_ROOT__"] = PackageRootAlias();
    aliases!["__USERNAME__"] = UsernameAlias();
  }
}

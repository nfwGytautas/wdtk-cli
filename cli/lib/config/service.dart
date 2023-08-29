part of wdtk_config;

/// WDTK predefined options for services
class ServiceOptions {
  late final String? gateway;

  ServiceOptions(Map data) {
    gateway = data["gateway"] ? data["gateway"].toString() : "false";
  }
}

/// Class for containing information about a service
class Service {
  late final String name;
  late final ServiceSource source;
  late final ServiceOptions? options;
  late final Map? config;

  Service(Map data) {
    name = data["name"];
    source = ServiceSource.createSource(data["source"]);

    if (data.containsKey("options")) {
      options = ServiceOptions(data["options"]);
    } else {
      options = null;
    }

    config = data["config"];
  }

  /// Returns the path to the service (from root)
  String getPath() {
    return Path.join(source.getPath(), "$name/");
  }

  /// Returns the path to the service output directory (from root)
  String getOutputDir() {
    return Path.absolute(".wdtk/bin/services/$name/");
  }

  /// Get the path to the config file of the service
  String getConfigFile(String deployment) {
    return Path.join(".wdtk/generated/configs/$deployment/$name.json");
  }
}

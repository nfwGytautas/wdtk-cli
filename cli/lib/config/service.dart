part of wdtk_config;

/// WDTK predefined options for services
class ServiceOptions {
  late final bool? gateway;

  ServiceOptions(Map data) {
    if (data.containsKey("gateway")) {
      gateway = data["gateway"];
    }
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
    }

    if (data.containsKey("config")) {
      config = data["config"];
    }
  }
}

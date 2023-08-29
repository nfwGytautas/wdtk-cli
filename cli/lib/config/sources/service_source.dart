part of wdtk_config;

/// A class for holding service types
class ServiceType {
  static final String git = "git";
  static final String binary = "bin";
  static final String local = "src";
}

/// Class containing service source configuration
abstract class ServiceSource {
  /// Get the type of the source
  String getType();

  /// Get the path to the source
  String getPath();

  /// Create source from YAML map
  static ServiceSource createSource(Map data) {
    final type = data["type"];

    if (type == ServiceType.git) {
      return GitSource(data);
    }

    if (type == ServiceType.local) {
      return CompiledSource(data);
    }

    print("Unknown source type $type");
    throw ArgumentError("Unknown source type $type");
  }
}

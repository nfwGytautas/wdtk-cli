part of wdtk_config;

/// Class containing service source configuration
abstract class ServiceSource {
  /// Get the type of the source
  String getType();

  /// Create source from YAML map
  static ServiceSource createSource(Map data) {
    final type = data["type"];

    if (type == "git") {
      return GitSource(data);
    }

    if (type == "src") {
      return CompiledSource(data);
    }

    print("Unknown source type $type");
    throw ArgumentError("Unknown source type $type");
  }
}

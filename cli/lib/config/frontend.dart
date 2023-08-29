part of wdtk_config;

/// Possible toolchains
class Toolchains {
  static final String flutter = "flutter";
}

/// Class for holding information about a single frontend platform entry
class PlatformEntry {
  late final String type;
  late final String toolchain;

  PlatformEntry(Map data) {
    type = data["type"];
    toolchain = data["toolchain"];
  }

  /// Returns the path to the service output directory (from root)
  String getOutputDir() {
    return Path.absolute(".wdtk/bin/frontend/$type/");
  }

  /// Get the path to the config file of the service
  String getConfigFile(String deployment) {
    return Path.join(".wdtk/generated/configs/$deployment/$type.json");
  }
}

/// Class for holding information about frontends
class FrontendConfig {
  late final List<PlatformEntry> platforms;

  FrontendConfig(Map data) {
    platforms = List.empty(growable: true);
    for (var platform in data["platforms"]) {
      platforms.add(PlatformEntry(platform));
    }
  }
}

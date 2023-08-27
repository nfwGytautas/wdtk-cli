part of wdtk_config;

/// Class for holding information about a single frontend platform entry
class PlatformEntry {
  late final String type;
  late final String toolchain;

  PlatformEntry(Map data) {
    type = data["type"];
    toolchain = data["toolchain"];
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

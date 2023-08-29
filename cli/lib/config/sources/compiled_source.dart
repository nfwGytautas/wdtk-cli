part of wdtk_config;

/// Class for holding compile language constants
class CompileLanguage {
  static final String go = "go";
  static final String flutter = "flutter";
}

/// A implementation of a service source that requires compilation
class CompiledSource implements ServiceSource {
  late final String language;

  CompiledSource(Map data) {
    language = data["language"];
  }

  @override
  String getType() {
    return ServiceType.local;
  }

  @override
  String getPath() {
    return "services/";
  }
}

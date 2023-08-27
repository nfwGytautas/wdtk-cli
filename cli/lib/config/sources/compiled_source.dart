part of wdtk_config;

/// A implementation of a service source that requires compilation
class CompiledSource implements ServiceSource {
  late final String language;

  CompiledSource(Map data) {
    language = data["language"];
  }

  @override
  String getType() {
    return "compiled";
  }
}

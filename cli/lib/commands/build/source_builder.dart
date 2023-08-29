part of wdtk_commands;

/// Base class for source builders
abstract class SourceBuilder {
  Future<BuildResult> buildService(Service service) {
    throw UnimplementedError("Frontend build not implemented");
  }

  Future<BuildResult> buildFrontend(PlatformEntry platform) {
    throw UnimplementedError("Frontend build not implemented");
  }

  /// Create a builder for the specified source
  static SourceBuilder? fromSource(ServiceSource source) {
    if (source is! CompiledSource) {
      return null;
    }

    return fromString(source.language);
  }

  /// Create a builder for the specified language
  static SourceBuilder? fromString(String language) {
    if (language == CompileLanguage.go) {
      return GoBuilder();
    }

    if (language == "flutter") {
      return FlutterBuilder();
    }

    return null;
  }
}

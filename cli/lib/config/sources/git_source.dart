part of wdtk_config;

/// A implementation of service source that is received from git repository
class GitSource implements ServiceSource {
  late final String remote;
  late final String language;

  GitSource(Map data) {
    remote = data["remote"];
    language = data["language"];
  }

  @override
  String getType() {
    return "git";
  }
}

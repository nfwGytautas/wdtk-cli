part of wdtk_config;

/// A implementation of service source that is received from git repository
class GitSource extends CompiledSource {
  late final String remote;
  late final String branch;

  GitSource(Map data) : super(data) {
    remote = data["remote"];
    branch = data["branch"] ?? "master";
  }

  @override
  String getType() {
    return ServiceType.git;
  }

  @override
  String getPath() {
    final parts = remote.split("/");
    final remoteName = parts.sublist(0, 3).join("/");
    return Path.join(".wdtk/remotes/", remoteName, branch);
  }
}

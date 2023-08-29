part of wdtk_commands;

/// Builder for go services
class GoBuilder extends SourceBuilder {
  @override
  Future<BuildResult> buildService(Service service) async {
    if (!await _goTidy(service)) {
      return BuildResult(service: service.name, success: false);
    }

    if (!await _goGet(service)) {
      return BuildResult(service: service.name, success: false);
    }

    if (!await _goBuild(service)) {
      return BuildResult(service: service.name, success: false);
    }

    return BuildResult(service: service.name, success: true);
  }

  // Run go mod tidy
  Future<bool> _goTidy(Service service) async {
    // go tidy
    var result = await Process.run("go", ["mod", "tidy"],
        workingDirectory: service.getPath());
    if (result.exitCode != 0) {
      return false;
    }

    return true;
  }

  /// Run go get ./
  Future<bool> _goGet(Service service) async {
    // go get
    var result = await Process.run("go", ["get", "./"],
        workingDirectory: service.getPath());
    if (result.exitCode != 0) {
      return false;
    }

    return true;
  }

  /// Run go build
  Future<bool> _goBuild(Service service) async {
    // go build
    var result = await Process.run("go",
        ["build", "-o", Path.join(service.getOutputDir(), service.name), "."],
        workingDirectory: service.getPath());
    if (result.exitCode != 0) {
      Logger.error(result.stderr);
      return false;
    }

    return true;
  }
}

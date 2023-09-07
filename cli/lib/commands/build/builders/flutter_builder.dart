part of wdtk_commands;

/// Builder for flutter frontend
class FlutterBuilder extends SourceBuilder {
  @override
  Future<BuildResult> buildFrontend(PlatformEntry platform) async {
    Logger.verbose("Building frontend ${platform.type} with flutter");
    var result = await Process.run("flutter", ["build", platform.type],
        workingDirectory: "frontend/_flutter/");

    if (result.exitCode != 0) {
      return BuildResult(service: platform.type, success: false);
    }

    final destination = Directory(".wdtk/bin/frontend/${platform.type}/");
    Utility.copyDirectory(
        Directory("frontend/_flutter/build/${platform.type}/"), destination);

    return BuildResult(service: platform.type, success: true);
  }
}

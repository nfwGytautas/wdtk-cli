part of wdtk_commands;

/// An action used to create a flutter project
class CreateFlutterProject implements ScaffoldAction {
  @override
  String get name => "Create flutter project";

  @override
  Future<ActionResult> execute(WDTKConfig config) async {
    if (config.frontend == null) {
      return ActionResult.nothingToDo;
    }

    String platforms = "";
    for (final frontendEntry in config.frontend!.platforms) {
      if (frontendEntry.toolchain == Toolchains.flutter) {
        platforms += "${frontendEntry.type},";
      }
    }

    if (platforms.isEmpty) {
      return ActionResult.nothingToDo;
    }

    // Remote trailing ","
    platforms = platforms.substring(0, platforms.length - 1);
    Logger.verbose("Creating flutter for platforms: '$platforms'");

    var result = await Process.run(
        "flutter", ["create", "_flutter", "--platforms=$platforms"],
        workingDirectory: "frontend/");

    if (result.stdout != "") {
      Logger.verbose(result.stdout, indent: Indent(indentSize: 2));
    }

    if (result.stderr != "") {
      Logger.verbose(result.stderr, indent: Indent(indentSize: 2));
    }

    if (result.exitCode == 0) {
      return ActionResult.success;
    }

    Logger.error("Failed to create flutter project");
    return ActionResult.error;
  }
}

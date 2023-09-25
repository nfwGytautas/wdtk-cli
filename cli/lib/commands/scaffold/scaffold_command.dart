part of wdtk_commands;

/// Scaffold command, used to create services, frontends, run checks, etc.
class ScaffoldCommand extends CliCommand {
  @override
  final name = "scaffold";

  @override
  final description =
      "Pull remote services, create flutter projects, update configuration files, etc.";

  ScaffoldCommand();

  @override
  void run() async {
    super.run();

    if (config == null) {
      // WDTK not initialized
      Logger.error(
          "No wdtk.yaml found, run wdtk init -n 'name', to create a project");
      return;
    }

    await Scaffold.run(config!);
  }
}

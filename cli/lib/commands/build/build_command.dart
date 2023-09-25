part of wdtk_commands;

/// Build command, used to build/copy/etc. all services
class BuildCommand extends CliCommand {
  @override
  final name = "build";

  @override
  final description = "Build all services";

  BuildCommand();

  @override
  void run() async {
    super.run();

    if (config == null) {
      // WDTK not initialized
      Logger.error(
          "No wdtk.yaml found, run wdtk init -n 'name', to create a project");
      return;
    }

    await Build.run(config!);
  }
}

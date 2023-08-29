part of wdtk_commands;

/// Scaffold command, used to create services, frontends, run checks, etc.
class ScaffoldCommand extends CliCommand {
  @override
  final name = "scaffold";

  @override
  final description = "Initialize the basic folder structure for wdtk";

  final List<ScaffoldAction> _actions = [
    CreateLocalServices(),
    PullGitServices(),
    GenerateConfigs(),
    WriteGoWork(),
    CreateFlutterProject()
  ];

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

    List<Future<ActionResult>> futures = List.empty(growable: true);

    for (var action in _actions) {
      futures.add(action.execute(config!));
    }

    final results = await Future.wait(futures);

    Logger.verbose("Scaffold summary");
    for (int i = 0; i < results.length; i++) {
      Logger.verbose("${_actions[i].name} : ${results[i]}", indent: Indent());
    }
  }
}

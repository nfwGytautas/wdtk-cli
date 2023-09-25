part of wdtk_commands;

/// An implementation of a scaffold action
class Scaffold {
  static final List<ScaffoldAction> _actions = [
    CreateLocalServices(),
    GenerateConfigs(),
    WriteGoWork(),
  ];

  static Future<bool> run(
    WDTKConfig config, {
    bool? pullGit = true,
    bool? generateFlutter = true,
  }) async {
    List<Future<ActionResult>> futures = List.empty(growable: true);
    bool result = true;

    var actionsToExecute = [..._actions];

    if (pullGit!) {
      actionsToExecute.add(PullGitServices());
    }

    if (generateFlutter!) {
      actionsToExecute.add(CreateFlutterProject());
    }

    for (var action in _actions) {
      futures.add(action.execute(config));
    }

    final results = await Future.wait(futures);

    Logger.verbose("Scaffold summary");
    for (int i = 0; i < results.length; i++) {
      Logger.verbose("${_actions[i].name} : ${results[i]}", indent: Indent());

      if (results[i] == ActionResult.error) {
        result = false;
      }
    }

    return result;
  }
}

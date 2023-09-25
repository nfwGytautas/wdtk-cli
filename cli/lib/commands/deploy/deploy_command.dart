part of wdtk_commands;

/// Deploy command used to deploy services
class DeployCommand extends CliCommand {
  @override
  final name = "deploy";

  @override
  final description = "Deploy services with a given deployment";

  DeployCommand() {
    argParser.addOption("deployment",
        abbr: "d", help: "Deployment to use", mandatory: true);
  }

  @override
  void run() async {
    super.run();

    if (config == null) {
      // WDTK not initialized
      Logger.error(
          "No wdtk.yaml found, run wdtk init -n 'name', to create a project");
      return;
    }

    if (!config!.deployments.containsKey(argResults!["deployment"])) {
      Logger.error("Unknown deployment ${argResults!["deployment"]}");
      return;
    }

    await Deploy.run(config!, argResults!["deployment"]);
  }
}

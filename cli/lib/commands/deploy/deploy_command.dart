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

    final deployment = config!.deployments[argResults!["deployment"]]!;

    config!.selectDeployment(deployment.name);

    List<Future<DeployResult>> futures = List.empty(growable: true);
    for (var service in config!.services.values) {
      config!.selectService(service.name);
      final deploymentEntry = deployment.getServiceDeployment(service.name);
      final deployer = Deployer.fromIp(deploymentEntry.ip!);

      if (deployer == null) {
        Logger.error(
            "Invalid deploy location ${deploymentEntry.ip} for ${service.name}");
        return;
      }

      final args = DeploySettings(
          name: service.name,
          configFile: service.getConfigFile(deployment.name),
          inputPath: service.getOutputDir(),
          outDirectory: config!.getStringValue(deploymentEntry.deploymentDir!));

      futures.add(deployer.deploy(args));
    }

    if (config!.frontend != null) {
      for (var entry in config!.frontend!.platforms) {
        config!.selectFrontend(entry.type);
        final deploymentEntry = deployment.getServiceDeployment(entry.type);
        final deployer = Deployer.fromIp(deploymentEntry.ip!);

        if (deployer == null) {
          Logger.error(
              "Invalid deploy location ${deploymentEntry.ip} for ${entry.type}");
          return;
        }

        final args = DeploySettings(
            name: entry.type,
            configFile: entry.getConfigFile(deployment.name),
            inputPath: entry.getOutputDir(),
            outDirectory:
                config!.getStringValue(deploymentEntry.deploymentDir!),
            configFileOverride: "assets/"
            );

        futures.add(deployer.deploy(args));
      }
    }

    final results = await Future.wait(futures);
    Logger.verbose("Deploy summary");
    for (final result in results) {
      Logger.verbose("${result.service} : ${result.success}", indent: Indent());
    }
  }
}

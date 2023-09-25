part of wdtk_commands;

class Deploy {
  static Future<bool> run(WDTKConfig config, String deploymentName, {bool? deployFrontend = true}) async {
    List<Future<DeployResult>> futures = List.empty(growable: true);

    final deployment = config.deployments[deploymentName]!;
    config.selectDeployment(deployment.name);

    for (var service in config.services.values) {
      config.selectService(service.name);
      final deploymentEntry = deployment.getServiceDeployment(service.name);
      final remoteTarget = deploymentEntry.sshUser == null
          ? null
          : "${deploymentEntry.sshUser}@${deploymentEntry.ip}";
      final deployer = Deployer.fromIp(deploymentEntry.ip!);

      if (deployer == null) {
        Logger.error(
            "Invalid deploy location ${deploymentEntry.ip} for ${service.name}");
        return false;
      }

      final args = DeploySettings(
          name: service.name,
          configFile: service.getConfigFile(deployment.name),
          inputPath: service.getOutputDir(),
          outDirectory: config.getStringValue(deploymentEntry.deploymentDir!),
          remoteTarget: remoteTarget);

      futures.add(deployer.deploy(args));
    }

    if (deployFrontend!) {
      if (config.frontend != null) {
        for (var entry in config.frontend!.platforms) {
          config.selectFrontend(entry.type);
          final deploymentEntry = deployment.getServiceDeployment(entry.type);
          final remoteTarget = deploymentEntry.sshUser == null
              ? null
              : "${deploymentEntry.sshUser}@${deploymentEntry.ip}";
          final deployer = Deployer.fromIp(deploymentEntry.ip!);

          if (deployer == null) {
            Logger.error(
                "Invalid deploy location ${deploymentEntry.ip} for ${entry.type}");
            return false;
          }

          final args = DeploySettings(
              name: entry.type,
              configFile: entry.getConfigFile(deployment.name),
              inputPath: entry.getOutputDir(),
              outDirectory: config.getStringValue(deploymentEntry.deploymentDir!),
              configFileOverride: "assets/",
              remoteTarget: remoteTarget);

          futures.add(deployer.deploy(args));
        }
      }
    }

    final results = await Future.wait(futures);
    bool runResult = true;
    Logger.verbose("Deploy summary");
    for (final result in results) {
      Logger.verbose("${result.service} : ${result.success}", indent: Indent());

      if (!result.success) {
        runResult = false;
      }
    }

    return runResult;
  }
}

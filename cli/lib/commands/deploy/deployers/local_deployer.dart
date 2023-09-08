part of wdtk_commands;

/// Deployer for deploying locally e.g. to the host pc
class LocalDeployer implements Deployer {
  @override
  Future<DeployResult> deploy(DeploySettings args) async {
    Logger.verbose("Local deploy for: ${args.name} to ${args.outDirectory}");

    try {
      await Directory(args.outDirectory).create(recursive: true);

      Utility.copyDirectory(
          Directory(args.inputPath), Directory(args.outDirectory));

      if (args.configFileOverride == null) {
        await File(args.configFile)
            .copy(Path.join(args.outDirectory, "WdtkConfig.json"));
      } else {
        Logger.verbose(
            "Config file deploy override to '${args.configFileOverride!}'");

        await Directory(Path.join(args.outDirectory, args.configFileOverride!))
            .create(recursive: true);

        File(args.configFile).copySync(Path.join(
            args.outDirectory, args.configFileOverride!, "WdtkConfig.json"));
      }
    } catch (ex) {
      Logger.error("Failed to deploy ${ex.toString()}");
      return DeployResult(service: args.name, success: false);
    }

    return DeployResult(service: args.name, success: true);
  }
}

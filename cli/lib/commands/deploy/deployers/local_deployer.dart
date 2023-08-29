part of wdtk_commands;

/// Deployer for deploying locally e.g. to the host pc
class LocalDeployer implements Deployer {
  @override
  Future<DeployResult> deploy(DeploySettings args) async {
    Logger.verbose("Local deploy for: ${args.name} to ${args.outDirectory}");

    await Directory(args.outDirectory).create(recursive: true);

    await Utility.copyDirectory(
        Directory(args.inputPath), Directory(args.outDirectory));
    await File(args.configFile).copy(Path.join(args.outDirectory, "WdtkConfig.json"));

    return DeployResult(service: args.name, success: true);
  }
}

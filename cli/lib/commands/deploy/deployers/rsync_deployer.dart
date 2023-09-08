part of wdtk_commands;

/// Deployer for rsync deploy
class RsyncDeployer implements Deployer {
  @override
  Future<DeployResult> deploy(DeploySettings args) async {
    Logger.verbose("Rsync deploy for: ${args.name} to ${args.outDirectory}");

    if (args.remoteTarget == null) {
      Logger.error("Remote target isn't defined");
      return DeployResult(service: args.name, success: false);
    }

    final configFileOut = Path.join(
        args.outDirectory, args.configFileOverride ?? "", "WdtkConfig.json");

    try {
      var mkdirCommand = await Process.run(
          "ssh", [args.remoteTarget!, 'mkdir -p ${args.outDirectory}'],
          runInShell: true);

      _printIfNotEmpty(mkdirCommand.stdout);
      _printIfNotEmpty(mkdirCommand.stderr);

      if (mkdirCommand.exitCode != 0) {
        return DeployResult(service: args.name, success: false);
      }

      var executableCopyCommand = await Process.run("rsync",
          ["-r", args.inputPath, "${args.remoteTarget!}:${args.outDirectory}"]);


      _printIfNotEmpty(executableCopyCommand.stdout);
      _printIfNotEmpty(executableCopyCommand.stderr);

      if (executableCopyCommand.exitCode != 0) {
        return DeployResult(service: args.name, success: false);
      }

      var configCopyCommand = await Process.run("rsync", [
        "-r",
        args.configFile,
        "${args.remoteTarget!}:$configFileOut"
      ]);

      _printIfNotEmpty(configCopyCommand.stdout);
      _printIfNotEmpty(configCopyCommand.stderr);

      if (configCopyCommand.exitCode != 0) {
        return DeployResult(service: args.name, success: false);
      }
    } catch (ex) {
      Logger.error("Failed to deploy '${ex.toString()}'");
      return DeployResult(service: args.name, success: false);
    }

    return DeployResult(service: args.name, success: true);
  }

  void _printIfNotEmpty(dynamic input) {
    if (input.isNotEmpty) {
      print(input);
    }
  }
}

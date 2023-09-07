part of wdtk_commands;

/// Settings for deploy command
class DeploySettings {
  final String name;
  final String configFile;
  final String inputPath;
  final String outDirectory;

  // Optional override for putting config files somewhere other than root outDirectory
  final String? configFileOverride;

  DeploySettings(
      {required this.name,
      required this.configFile,
      required this.inputPath,
      required this.outDirectory,
      this.configFileOverride});
}

/// Base class for possible deployment strategies
abstract class Deployer {
  Future<DeployResult> deploy(DeploySettings args);

  /// Create a deployer for the specified ip
  static Deployer? fromIp(String ip) {
    // Local
    if (ip == "127.0.0.1") {
      return LocalDeployer();
    }

    // Remote
    // TODO: Implement

    return null;
  }
}

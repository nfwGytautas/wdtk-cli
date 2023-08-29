part of wdtk_commands;

/// Settings for deploy command
class DeploySettings {
  final String name;
  final String configFile;
  final String inputPath;
  final String outDirectory;

  DeploySettings(
      {required this.name, required this.configFile, required this.inputPath, required this.outDirectory});
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

part of wdtk_commands;

/// A base class for CLI commands
abstract class CliCommand extends Command {
  WDTKConfig? config;

  @override
  void run() {
    // Setup logging
    final settings = LoggerSettings(verbose: globalResults!["verbose"]);
    Logger.setSettings(settings);

    // Load config if possible
    config = WDTKConfig.load();
  }
}

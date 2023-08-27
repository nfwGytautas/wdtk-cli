part of wdtk_cli_logging;

class LoggerSettings {
  final bool verbose;

  LoggerSettings({required this.verbose});

  /// Returns default settings for the logger
  static LoggerSettings getDefault() {
    return LoggerSettings(verbose: false);
  }
}

/// Logger class
class Logger {
  static LoggerSettings _settings = LoggerSettings.getDefault();

  /// Set the settings of the logger
  static void setSettings(LoggerSettings settings) {
    _settings = settings;
  }

  /// Message only shown when verbose mode is active
  static void verbose(String message) {
    if (_settings.verbose) {
      print("[Verbose] $message");
    }
  }

  /// Print a warning message
  static void warning(String message) {
    print("[Warning] $message");
  }

  /// Print a error message
  static void error(String message) {
    print("[Error] $message");
  }
}

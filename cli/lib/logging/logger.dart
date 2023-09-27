part of wdtk_cli_logging;

/// Possible logger settings
class LoggerSettings {
  final bool verbose;

  LoggerSettings({required this.verbose});

  /// Returns default settings for the logger
  static LoggerSettings getDefault() {
    return LoggerSettings(verbose: false);
  }
}

/// Indents for logging
class Indent {
  static final String _indentSize = "   ";

  String _indent = "";

  @override
  String toString() {
    return _indent;
  }

  Indent({int? indentSize}) : _indent = _indentSize * (indentSize ?? 1);
}

/// Logger class
class Logger {
  static LoggerSettings _settings = LoggerSettings.getDefault();

  /// Set the settings of the logger
  static void setSettings(LoggerSettings settings) {
    _settings = settings;
  }

  /// Message only shown when verbose mode is active
  static void verbose(String message, {Indent? indent}) {
    if (_settings.verbose) {
      print("[Verbose] ${indent ?? ""}$message");
    }
  }

  /// Info message
  static void info(String message) {
    print(message);
  }

  /// Print a warning message
  static void warning(String message) {
    print("\x1B[33m$message\x1B[0m");
  }

  /// Print a error message
  static void error(String message) {
    print("\x1B[31m$message\x1B[0m");
  }
}

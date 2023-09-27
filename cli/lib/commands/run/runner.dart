part of wdtk_commands;

/// A simple service runner
class ServiceRunner {
  final Service service;
  final bool? watchChanges;

  Process? _process = null;
  bool _localMutex = false;

  ServiceRunner({required this.service, this.watchChanges}) {
    if (watchChanges ?? false) {
      _watch();
    }
  }

  /// Run the method
  void run() async {
    while (_localMutex) {}
    _localMutex = true;

    if (_process != null) {
      // Kill first
      Logger.info("Killing ${service.name}");

      _process!.kill();
    }

    Logger.info("Starting ${service.name}");
    _process = await Process.start("./${service.name}", [],
        workingDirectory: "dev/${service.name}/", runInShell: true);

    // Logging information
    _process!.stdout.transform(utf8.decoder).forEach((element) {
      _handleOutput(element);
    });

    _process!.stderr.transform(utf8.decoder).forEach((element) {
      _handleOutput(element, error: true);
    });

    _localMutex = false;
  }

  /// Stop service
  Future<void> stop() async {
    while (_localMutex) {}
    _localMutex = true;

    if (_process != null) {
      // Kill first
      Logger.info("Killing ${service.name}");

      _process!.kill();
    }

    _localMutex = false;
  }

  /// Watch for changes
  void _watch() async {
    var watcher = DirectoryWatcher(
        Path.join(Directory.current.path, "services", service.name));
    watcher.events.listen(_onModified);
  }

  /// Called when a service is changed
  void _onModified(WatchEvent? event) {
    Logger.info("Service ${service.name} modified");
    run();
  }

  /// Handle output from stdin or stderr
  void _handleOutput(String output, {bool error = false}) {
    for (final line in output.split("\n")) {
      if (line.isEmpty || line == "\n") {
        continue;
      }

      final message = "[${service.name.padLeft(20, " ")}] $line";

      if (error) {
        Logger.error(message);
      } else {
        Logger.info(message);
      }
    }
  }
}

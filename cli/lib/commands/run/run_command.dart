part of wdtk_commands;

/// Run command used to quickly execute all services locally in one place
class RunCommand extends CliCommand {
  @override
  final name = "run";

  @override
  final description = "Run the services [Local deployment only]";
  RunCommand();

  final List<ServiceRunner> _localServices = List.empty(growable: true);
  final List<ServiceRunner> _remoteServices = List.empty(growable: true);
  bool _mutex = false;

  @override
  void run() async {
    super.run();

    if (config == null) {
      // WDTK not initialized
      Logger.error(
          "No wdtk.yaml found, run wdtk init -n 'name', to create a project");
      return;
    }

    try {
      await _prepareRun();

      _runLocalServices();
      _runRemoteServices();

      _watchConsole();
      _watchConfig();
    } catch (ex) {
      Logger.error("Got a WDTK exception while running. Stopping all services");
      for (final service in _localServices) {
        service.stop();
      }

      for (final service in _remoteServices) {
        service.stop();
      }
    }
  }

  /// Load lists from wdtk config
  Future<void> _loadLists(WDTKConfig config) async {
    for (final service in _localServices) {
      await service.stop();
    }

    for (final service in _remoteServices) {
      await service.stop();
    }

    _localServices.clear();
    _remoteServices.clear();

    for (final service in config.services.values) {
      if (service.source.getType() == ServiceType.local) {
        _localServices.add(ServiceRunner(service: service, watchChanges: true));
      } else {
        _remoteServices.add(ServiceRunner(service: service));
      }
    }
  }

  /// Create a watcher for wdtk.yaml
  void _watchConfig() {
    // Start watching wdtk.yaml
    var configWatcher =
        FileWatcher(Path.join(Directory.current.path, "wdtk.yaml"));
    configWatcher.events.listen(_onConfigChange);
  }

  /// Watch console for 'R' press to reload
  void _watchConsole() async {
    stdin.lineMode = false;
    stdin.echoMode = false;

    await for (var inputList in stdin) {
      final input = Utf8Decoder().convert(inputList);
      print(input);
      if (input == "r") {
        for (int i = 0; i < stdout.terminalLines; i++) {
          stdout.writeln();
        }

        final message = "RELOAD";
        final spacing =
            " " * ((stdout.terminalColumns / 2).floor() - message.length);

        stdout.writeln("=" * stdout.terminalColumns);
        stdout.writeln("$spacing$message");
        stdout.writeln("=" * stdout.terminalColumns);

        _onConfigChange(null);
      }
    }
  }

  /// Called when a config is changed
  void _onConfigChange(WatchEvent? event) async {
    while (_mutex) {}

    _mutex = true;

    print("Config changed running scaffold, build and deploy");

    await _prepareRun();

    _runLocalServices();
    _runRemoteServices();

    _mutex = false;
  }

  /// Run remote services
  void _runRemoteServices() {
    // Run all services
    for (final service in _remoteServices) {
      service.run();
    }
  }

  /// Run local services
  void _runLocalServices() async {
    // Run all services
    for (final service in _localServices) {
      service.run();
    }
  }

  /// Scaffold, Build, Deploy cycle
  Future<void> _prepareRun() async {
    WDTKConfig? newConfig = WDTKConfig.load();
    if (newConfig == null) {
      Logger.error("Failed to load new config fix errors");
      return;
    }

    await _loadLists(newConfig);

    // Run scaffold
    bool scaffold = await Scaffold.run(newConfig);
    if (!scaffold) {
      Logger.error(
          "Scaffold command failed, check that wdtk.yaml config is without errors, no services reloaded");
      return;
    }

    // Build services
    bool build = await Build.run(newConfig, buildFrontend: false);
    if (!build) {
      Logger.error(
          "Build command failed, correct errors, no services reloaded");
      return;
    }

    // Deploy to dev
    bool deploy = await Deploy.run(newConfig, "dev", deployFrontend: false);
    if (!deploy) {
      Logger.error(
          "Deploy command failed, correct errors, no services reloaded");
      return;
    }
  }
}

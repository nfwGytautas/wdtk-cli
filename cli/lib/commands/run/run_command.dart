part of wdtk_commands;

/// Run command used to quickly execute all services locally in one place
class RunCommand extends CliCommand {
  @override
  final name = "run";

  @override
  final description = "Run the services [Local deployment only]";
  RunCommand();

  @override
  void run() async {
    super.run();

    if (config == null) {
      // WDTK not initialized
      Logger.error(
          "No wdtk.yaml found, run wdtk init -n 'name', to create a project");
      return;
    }

    List<Future<Process>> futures = List.empty(growable: true);

    for (var service in config!.services.values) {
      futures.add(_executeService(service));
    }

    var processes = await Future.wait(futures);

    for (var process in processes) {
      await process.exitCode;
    }
  }

  /// Execute a single service
  Future<Process> _executeService(Service service) async {
    Logger.info("Starting ${service.name}");

    var p = await Process.start("./${service.name}", [],
        workingDirectory: "dev/${service.name}/", runInShell: true);


    // Logging information
    p.stdout.transform(utf8.decoder).forEach((element) {
      for (final line in element.split("\n")) {
        if (line.isEmpty || line == "\n") {
          continue;
        }

        print("[${service.name.padLeft(20, " ")}] $line");
      }
    });

    return p;
  }
}

part of wdtk_commands;

/// Build command, used to build/copy/etc. all services
class BuildCommand extends CliCommand {
  @override
  final name = "build";

  @override
  final description = "Build all services";

  BuildCommand();

  @override
  void run() async {
    super.run();

    if (config == null) {
      // WDTK not initialized
      Logger.error(
          "No wdtk.yaml found, run wdtk init -n 'name', to create a project");
      return;
    }

    List<Future<BuildResult>> futures = List.empty(growable: true);
    for (var service in config!.services.values) {
      futures.add(_buildService(service));
    }

    if (config!.frontend != null) {
      for (var entry in config!.frontend!.platforms) {
        futures.add(_buildFrontend(entry));
      }
    }

    final results = await Future.wait(futures);
    Logger.verbose("Build summary");
    for (final result in results) {
      Logger.verbose("${result.service} : ${result.success}", indent: Indent());
    }
  }

  /// Build a service and return a result
  Future<BuildResult> _buildService(Service service) async {
    if (service.source.getType() == ServiceType.binary) {
      // For binary just copy
      // TODO: Implement
      return BuildResult(service: service.name, success: true);
    }

    // Git and local need building
    final builder = SourceBuilder.fromSource(service.source);
    if (builder == null) {
      Logger.error("Unsupported language for ${service.name}");
      return BuildResult(service: service.name, success: false);
    }

    return builder.buildService(service);
  }

  /// Build a frontend and return a result
  Future<BuildResult> _buildFrontend(PlatformEntry platform) async {
    // Git and local need building
    final builder = SourceBuilder.fromString(platform.toolchain);
    if (builder == null) {
      Logger.error("Unsupported language for ${platform.type}");
      return BuildResult(service: platform.type, success: false);
    }

    return builder.buildFrontend(platform);
  }
}

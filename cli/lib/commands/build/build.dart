part of wdtk_commands;

/// A build command do
class Build {
  static Future<bool> run(WDTKConfig config, {bool? buildFrontend = true}) async {
    List<Future<BuildResult>> futures = List.empty(growable: true);
    bool runResult = true;

    for (var service in config.services.values) {
      futures.add(_buildService(service));
    }

    if (buildFrontend!) {
      if (config.frontend != null) {
        for (var entry in config.frontend!.platforms) {
          futures.add(_buildFrontend(entry));
        }
      }
    }

    final results = await Future.wait(futures);
    Logger.verbose("Build summary");
    for (final result in results) {
      Logger.verbose("${result.service} : ${result.success}", indent: Indent());

      if (!result.success) {
        runResult = false;
      }
    }

    return runResult;
  }

  /// Build a service and return a result
  static Future<BuildResult> _buildService(Service service) async {
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
  static Future<BuildResult> _buildFrontend(PlatformEntry platform) async {
    // Git and local need building
    final builder = SourceBuilder.fromString(platform.toolchain);
    if (builder == null) {
      Logger.error("Unsupported language for ${platform.type}");
      return BuildResult(service: platform.type, success: false);
    }

    return builder.buildFrontend(platform);
  }
}

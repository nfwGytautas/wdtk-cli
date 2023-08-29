part of wdtk_config;

/// Alias for processing __HOME__
class HomeAlias implements Alias {
  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    final homeDir =
        Platform.environment['HOME'] ?? Platform.environment['USERPROFILE'];

    return homeDir!;
  }

  @override
  void _compute(WDTKConfig config) {
    // Nothing to do
  }
}

/// Alias for processing __PACKAGE__
class PackageAlias implements Alias {
  final String packageName;

  PackageAlias({required this.packageName});

  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    return packageName;
  }

  @override
  void _compute(WDTKConfig config) {
    // Nothing to do
  }
}

/// Alias for processing __PACKAGE_ROOT__
class PackageRootAlias implements Alias {
  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    return Directory.current.path;
  }

  @override
  void _compute(WDTKConfig config) {
    // Nothing to do
  }
}

/// Alias for processing __DEPLOYMENT__
class DeploymentAlias implements Alias {
  final String deploymentName;

  DeploymentAlias({required this.deploymentName});

  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    return deploymentName;
  }

  @override
  void _compute(WDTKConfig config) {
    // Nothing to do
  }
}

/// Alias for processing __SERVICE__
class ServiceAlias implements Alias {
  final WDTKConfig config;

  ServiceAlias({required this.config});

  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    if (config._currentService == null) {
      if (config._currentFrontend == null) {
        Logger.warning("No service currently being processed");
        return "@@NULL@@";
      }

      return config._currentFrontend!;
    }

    return config._currentService!.name;
  }

  @override
  void _compute(WDTKConfig config) {
    // Nothing to do
  }
}

/// Alias for processing __DEPLOYMENT_DIR__
class DeploymentDirAlias implements Alias {
  final WDTKConfig config;

  DeploymentDirAlias({required this.config});

  @override
  String getComputedValue(
      {Map<String, String>? args, bool forwardRemaining = false}) {
    if (args == null || !args.containsKey("service")) {
      Logger.error("'__DEPLOYMENT_DIR__' alias requires 'service' argument");
      return "@@ERROR@@";
    }

    final service = args["service"]!;

    DeploymentEntry entry =
        config._selectedDeployment!.getServiceDeployment(service);
    if (entry.deploymentDir == null) {
      return "@@ERROR@@";
    }

    final deploymentDir =
        entry.deploymentDir!.replaceAll(RegExp(r"\${__SERVICE__}"), service);
    return config.getStringValue(deploymentDir);
  }

  @override
  void _compute(WDTKConfig config) {
    // Nothing to do
  }
}

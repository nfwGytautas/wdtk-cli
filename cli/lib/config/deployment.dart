part of wdtk_config;

/// Class for a single deployment entry
class DeploymentEntry {
  late final String? ip;
  late final String? apiKey;
  late final String? deploymentDir;
  late final String? port;
  late final String? sshUser;

  DeploymentEntry.empty();

  DeploymentEntry(Map data) {
    ip = data["ip"];
    apiKey = data["apiKey"];
    deploymentDir = data["dir"];
    port = data["port"]?.toString();
    sshUser = data["ssh-user"]?.toString();
  }

  /// Create a new entry where empty fields from entry are filled with ones from defaults
  static DeploymentEntry fill(DeploymentEntry defaults, DeploymentEntry entry) {
    DeploymentEntry result = DeploymentEntry.empty();

    result.ip = entry.ip ?? defaults.ip;
    result.apiKey = entry.apiKey ?? defaults.apiKey;
    result.deploymentDir = entry.deploymentDir ?? defaults.deploymentDir;
    result.port = entry.port ?? defaults.port;
    result.sshUser = entry.sshUser ?? defaults.sshUser;

    return result;
  }
}

/// Class for holding information about a deployment configuration
class Deployment {
  late final String name;

  late final DeploymentEntry? defaults;
  late final Map<String, Alias>? aliases;
  late final Map<String, DeploymentEntry>? settings;

  Deployment(Map data) {
    name = data["name"];

    if (data.containsKey("defaults")) {
      defaults = DeploymentEntry(data["defaults"]);
    }

    if (data.containsKey("aliases")) {
      aliases = <String, Alias>{};
      for (var entry in data["aliases"].entries) {
        aliases![entry.key] = UserAlias(name: entry.key, value: entry.value.toString());
      }
    }

    if (data.containsKey("settings")) {
      settings = <String, DeploymentEntry>{};
      for (var entry in data["settings"].entries) {
        settings![entry.key] = DeploymentEntry(entry.value);
      }
    }

    if (defaults == null && settings == null) {
      Logger.error(
          "If deployment doesn't specify settings for each service it must define defaults");
    }
  }

  /// Get deployment alias
  Alias? getAlias(String name) {
    if (aliases == null) {
      return null;
    }

    if (!aliases!.containsKey(name)) {
      return null;
    }

    return aliases![name];
  }

  /// Returns the deployment entry for the specified service
  DeploymentEntry getServiceDeployment(String service) {
    if (settings == null) {
      return defaults!;
    }

    if (!settings!.containsKey(service)) {
      return defaults!;
    }

    DeploymentEntry serviceSettings = settings![service]!;
    return DeploymentEntry.fill(defaults!, serviceSettings);
  }

  /// Compute deployment alias values
  void _computeAliases(WDTKConfig config) {
    if (aliases == null) {
      return;
    }

    // Compute globals
    for (var alias in aliases!.values) {
      alias._compute(config);
    }
  }
}

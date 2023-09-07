part of wdtk_commands;

/// An action used to generate service and frontend configuration files
class GenerateConfigs implements ScaffoldAction {
  final encoder = JsonEncoder.withIndent(' ' * 4);

  @override
  String get name => "Create config files";

  @override
  Future<ActionResult> execute(WDTKConfig config) async {
    ActionResult result = ActionResult.nothingToDo;

    Service? gatewayService;

    // Find gateway service
    for (var service in config.services.values) {
      if (service.options == null) {
        continue;
      }

      if (service.options!.gateway == "true") {
        if (gatewayService != null) {
          Logger.error(
              "Multiple gateway services current: ${service.name}, previous: ${gatewayService.name}");
          return ActionResult.error;
        }

        gatewayService = service;
      }
    }

    if (gatewayService == null) {
      Logger.error("No gateway service specified");
      return ActionResult.error;
    }

    // Create deployment configs
    for (var deployment in config.deployments.values) {
      await _createConfigs(config, gatewayService, deployment);
    }

    return result;
  }

  /// Create configs
  Future<void> _createConfigs(
      WDTKConfig config, Service gatewayService, Deployment deployment) async {
    config.selectDeployment(deployment.name);

    final gatewayDeployment =
        deployment.getServiceDeployment(gatewayService.name);

    var locatorTable = [];

    final gatewayIp = config
        .getStringValue("${gatewayDeployment.ip}:${gatewayDeployment.port}");

    for (var service in config.services.values) {
      // Skip gateway
      if (service.name == gatewayService.name) {
        continue;
      }

      config.selectService(service.name);

      final serviceDeployment = deployment.getServiceDeployment(service.name);

      final outPath = Path.join(
          ".wdtk/generated/configs/", deployment.name, "${service.name}.json");

      var file = await File(outPath).create(recursive: true);

      // Standard settings
      final runIp = config
          .getStringValue("${serviceDeployment.ip}:${serviceDeployment.port}");
      var configMap = <String,dynamic>{
        "runAddress": runIp,
        "gatewayIp": gatewayIp,
        "apiKey": config.getStringValue(serviceDeployment.apiKey!)
      };

      // User config map
      if (service.config != null) {
        service.config!.forEach((key, value) {
          if (value is String) {
            configMap[key] = config.getStringValue(value.toString());
          } else {
            configMap[key] = value;
          }
        });
      }

      Logger.verbose("Writing $configMap", indent: Indent());

      locatorTable.add({"service": service.name, "ip": runIp});

      file.writeAsString(encoder.convert(configMap));
    }

    if (config.frontend != null) {
      for (var entry in config.frontend!.platforms) {
        final outPath = Path.join(
            ".wdtk/generated/configs/", deployment.name, "${entry.type}.json");

        var file = await File(outPath).create(recursive: true);
        file.writeAsString(encoder.convert({"gatewayIp": gatewayIp}));
      }
    }

    config.selectService(gatewayService.name);
    var gatewayConfig = <String, dynamic>{
      "apiKey": config.getStringValue(gatewayDeployment.apiKey!),
      "runAddress": config
          .getStringValue("${gatewayDeployment.ip}:${gatewayDeployment.port}"),
    };

    gatewayConfig["locatorTable"] = locatorTable;

    Logger.verbose("Writing $gatewayConfig", indent: Indent());

    final outPath = Path.join(".wdtk/generated/configs/", deployment.name,
        "${gatewayService.name}.json");
    var file = await File(outPath).create(recursive: true);
    file.writeAsString(encoder.convert(gatewayConfig));
  }
}

part of wdtk_commands;

/// An action used to create services that are specified in wdtk.yaml but don't exist in the directory yet
class CreateLocalServices implements ScaffoldAction {
  @override
  String get name => "Create local templates";

  @override
  Future<ActionResult> execute(WDTKConfig config) async {
    ActionResult result = ActionResult.nothingToDo;

    for (var service in config.getServicesOfType(ServiceType.local)) {
      final exists = await Directory("services/${service.name}/").exists();
      if (exists) {
        // Do nothing
        continue;
      }

      Logger.verbose("Creating $service");

      final language = (service.source as CompiledSource).language;

      final template = LanguageTemplate.fromString(language);
      if (template == null) {
        Logger.error("Unsupported language $language for ${service.name}");
        continue;
      }

      final args = LanguageTemplateArgs(
          rootPath: "services/${service.name}/", domain: config.package, service: service.name);

      template.write(args);

      result = ActionResult.success;
    }

    return result;
  }
}

part of wdtk_commands;

/// An action used to update go.work file in root directory
class WriteGoWork implements ScaffoldAction {
  @override
  String get name => "Write go.work";

  @override
  Future<ActionResult> execute(WDTKConfig config) async {
    ActionResult result = ActionResult.nothingToDo;

    var file = File("go.work").openWrite();

    file.write("go 1.20\n");

    for (var service in config.getServicesOfType(ServiceType.local)) {
      if ((service.source as CompiledSource).language == CompileLanguage.go) {
        result = ActionResult.success;

        file.write("use ${service.getPath()}\n");
      }
    }

    return result;
  }
}

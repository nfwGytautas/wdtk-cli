part of wdtk_templates;

/// Golang template
class GoLanguageTemplate implements LanguageTemplate {
  @override
  Future<void> write(LanguageTemplateArgs args) async {
    await Directory(args.rootPath).create(recursive: true);

    await Future.wait([_writeMain(args), _writeGoMod(args)]);
  }

  /// Write main.go
  Future<void> _writeMain(LanguageTemplateArgs args) async {
    final mainGoTemplate = """
package main

func main() {
	println("Running service")
}
""";

    File(Path.join(args.rootPath, "main.go")).writeAsString(mainGoTemplate);
  }

  /// Write go.mod
  Future<void> _writeGoMod(LanguageTemplateArgs args) async {
    final goModTemplate = """
module ${Path.join(args.domain, "services", args.service)}

go 1.20
""";

    File(Path.join(args.rootPath, "go.mod")).writeAsString(goModTemplate);
  }
}

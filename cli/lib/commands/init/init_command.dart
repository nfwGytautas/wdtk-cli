part of wdtk_commands;

/// Init command, used to setup the folder structure and basic templates for wdtk
class InitCommand extends CliCommand {
  @override
  final name = "init";

  @override
  final description = "Initialize the basic folder structure for wdtk";

  InitCommand() {
    argParser.addOption("name",
        abbr: "n", help: "Name of the project", mandatory: true);
  }

  @override
  void run() async {
    super.run();

    if (config != null) {
      // Already created
      Logger.info("Valid 'wdtk.yaml' file already exists");
      return;
    }

    // Check that the directory is empty
    final directoryEmpty =
        await Directory(Directory.current.path).list().isEmpty;
    if (!directoryEmpty) {
      Logger.error("The directory isn't empty");
      return;
    }

    await Future.wait([
      _writeConfigFile(),
      _createDirectoryStructure(),
    ]);
  }

  /// Write a template wdtk.yaml file
  Future<void> _writeConfigFile() async {
    final name = argResults!["name"];

    final contents = templateWdtkYaml("$name.com/$name/", name);
    await File("wdtk.yaml").writeAsString(contents);
  }

  /// Create the directory structure
  Future<void> _createDirectoryStructure() async {
    final name = argResults!["name"];

    Directory("services").create();
    Directory("frontend").create();

    Directory(".wdtk/logs/").create(recursive: true);
    Directory(".wdtk/generated/").create(recursive: true);
    Directory(".wdtk/bin/services/").create(recursive: true);
    Directory(".wdtk/bin/frontends/").create(recursive: true);
    Directory(".wdtk/remotes/").create(recursive: true);

    (await File(".gitignore").create()).writeAsString(templateGitIgnore());
    (await File("README.md").create()).writeAsString(templateRootReadme(name));
    File("CHANGELOG.md").create();

    (await File("services/README.md").create(recursive: true))
        .writeAsString(templateServiceReadme());
    (await File("frontend/README.md").create(recursive: true))
        .writeAsString(templateFrontendReadme());
  }
}

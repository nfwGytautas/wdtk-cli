import 'dart:io';

import 'package:args/command_runner.dart';
import 'package:wdtk_cli/commands/commands.dart';

void main(List<String> arguments) {
  // TODO: Automatically get it from pubspec.yaml
  final version = "0.0.0";

  var runner = CommandRunner("wdtk", "Webdev-Toolkit (v$version)");

  runner.addCommand(InitCommand());
  runner.addCommand(ScaffoldCommand());
  runner.addCommand(BuildCommand());
  runner.addCommand(DeployCommand());

  runner.argParser.addFlag('verbose',
      negatable: false, abbr: "v", help: "Enable verbose logging");

  runner.run(arguments).catchError((error) {
    if (error is UsageException) {
      print(error);
      exit(64);
    }

    throw error;
  });

  return;
}

library wdtk_commands;

import 'dart:io';
import 'dart:convert';

import 'package:path/path.dart' as Path;
import 'package:args/command_runner.dart';
import 'package:wdtk_cli/config/wdtk_config.dart';
import 'package:wdtk_cli/logging/logging.dart';
import 'package:wdtk_cli/templates/templates.dart';

part 'cli_command.dart';


part 'build/build_command.dart';
part 'build/build_result.dart';
part 'build/source_builder.dart';

part 'build/builders/go_builder.dart';
part 'build/builders/flutter_builder.dart';

part 'deploy/deploy_command.dart';
part 'deploy/deploy_result.dart';
part 'deploy/deployer.dart';

part 'deploy/deployers/local_deployer.dart';

part 'init/init_command.dart';

part 'run/run_command.dart';

part 'scaffold/scaffold_command.dart';
part 'scaffold/scaffold_action.dart';

part 'scaffold/actions/create_flutter_project.dart';
part 'scaffold/actions/create_local_services.dart';
part 'scaffold/actions/generate_configs.dart';
part 'scaffold/actions/pull_git_services.dart';
part 'scaffold/actions/write_go_work.dart';

part 'utility/utility.dart';

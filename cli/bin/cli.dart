import 'package:wdtk_cli/config/wdtk_config.dart';
import 'package:wdtk_cli/logging/logging.dart';

void main(List<String> arguments) {
  Logger.setSettings(LoggerSettings(verbose: true));

  WDTKConfig? config = WDTKConfig.load();

  if (config == null) {
    print("Failed to load config");
    return;
  }

  config.selectDeployment("dev");
  print(config
      .getAlias("::databaseString")!
      .getComputedValue(args: {"database": "auth"}));

  print(config.getAliasValue("\${::databaseString, database: auth}"));

  print(config.getStringValue(
      "\${__HOME__}/\${__PACKAGE__}/\${__DEPLOYMENT__}/\${__SERVICE__}"));

  print(config.getStringValue("Non alias string"));

  print(config.getStringValue("\${__DEPLOYMENT_DIR__, service: Web}"));
}

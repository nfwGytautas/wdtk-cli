import 'dart:io';

import 'package:test/test.dart';
import 'package:wdtk_cli/config/wdtk_config.dart';

void main() {
  WDTKConfig? config = WDTKConfig.load(path: "${Directory.current.path}/test/config/wdtk.yaml");
  assert(config != null, "Failed to load config");

  config!.selectDeployment("dev");

  test("getAliasValue", () {
    expect(config.getAliasValue("\${apiKey}"), "API_KEY_GOES_HERE");
  });

  test("getComputedValue", () {
    final alias = config.getAlias("::databaseString");

    expect(alias, isNot(null));

    final computedValue = alias!.getComputedValue(args: {"database": "auth"});

    expect(computedValue,
        "user:password@tcp(127.0.0.1:3306)/auth?charset=utf8mb4&parseTime=True&loc=Local");
  });

  test("getStringValue [no service]", () {
    final value = config.getStringValue(
        "\${__HOME__}/\${__PACKAGE__}/\${__DEPLOYMENT__}/\${__SERVICE__}");

    expect(value, "/Users/gytautaskazlauskas/test/dev/@@NULL@@");
  });

  test("getStringValue [with service]", () {
    final homeDir =
        Platform.environment['HOME'] ?? Platform.environment['USERPROFILE'];
    config.selectService("Gateway");
    final value = config.getStringValue(
        "\${__HOME__}/\${__PACKAGE__}/\${__DEPLOYMENT__}/\${__SERVICE__}");

    expect(value, "$homeDir/test/dev/Gateway");
  });

  test("getStringValue [no aliases]", () {
    final value = config.getStringValue("No aliases");
    expect(value, "No aliases");
  });

  // Internal aliases
  test("__HOME__", () {
    final homeDir =
        Platform.environment['HOME'] ?? Platform.environment['USERPROFILE'];
    final value = config.getStringValue("\${__HOME__}");
    expect(value, homeDir);
  });

  test("__PACKAGE__", () {
    final value = config.getStringValue("\${__PACKAGE__}");
    expect(value, "test");
  });

  test("__DEPLOYMENT__", () {
    final value = config.getStringValue("\${__DEPLOYMENT__}");
    expect(value, "dev");
  });

  test("__SERVICE__", () {
    config.selectService("Gateway");
    final value = config.getStringValue("\${__SERVICE__}");
    expect(value, "Gateway");
  });

  test("__DEPLOYMENT_DIR__", () {
    final homeDir =
        Platform.environment['HOME'] ?? Platform.environment['USERPROFILE'];
    config.selectService("Gateway");
    final value = config.getStringValue("\${__DEPLOYMENT_DIR__, service: Web}");
    expect(value, "$homeDir/test/dev/Web/");
  });
}

part of wdtk_commands;

/// Build result class to hold the results for a build
class BuildResult {
  final String service;
  final bool success;

  BuildResult({required this.service, required this.success});
}

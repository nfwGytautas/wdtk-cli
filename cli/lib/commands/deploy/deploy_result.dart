part of wdtk_commands;

/// Class for holding deploy results
class DeployResult {
  final String service;
  final bool success;

  DeployResult({required this.service, required this.success});
}

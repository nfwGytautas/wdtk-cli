part of wdtk_commands;

/// Possible results of a scaffold action
enum ActionResult {
  success,
  error,
  nothingToDo,
}

/// Base class for scaffold actions
abstract class ScaffoldAction {
  /// Name of the action
  String get name;

  /// Execute the action
  Future<ActionResult> execute(WDTKConfig config);
}

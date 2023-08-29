part of wdtk_commands;

/// Holding class for remote info
class RemoteEntry {
  final String remote;
  final String branch;

  RemoteEntry({required this.remote, required this.branch});

  @override
  bool operator ==(other) {
    if (other is RemoteEntry) {
      return other.remote == remote && other.branch == branch;
    }

    return false;
  }

  @override
  int get hashCode => Object.hash(remote, branch);

  @override
  String toString() {
    return "Remote: $remote, Branch: $branch";
  }
}

/// An action used to pull remote services from git
class PullGitServices implements ScaffoldAction {
  @override
  String get name => "Pull Git services";

  @override
  Future<ActionResult> execute(WDTKConfig config) async {
    Set<RemoteEntry> entries = {};

    // Get remote list
    for (var service in config.getServicesOfType(ServiceType.git)) {
      final source = (service.source as GitSource);

      final parts = source.remote.split("/");
      final remoteName = parts.sublist(0, 3).join("/");

      final entry = RemoteEntry(remote: remoteName, branch: source.branch);
      entries.add(entry);
    }

    if (entries.isEmpty) {
      return ActionResult.nothingToDo;
    }

    // Clone/Pull
    List<Future<void>> futures = List.empty(growable: true);
    Logger.verbose("Remotes to pull:");
    for (var element in entries) {
      Logger.verbose(element.toString(), indent: Indent());

      futures.add(_cloneOrPull(element));
    }

    await Future.wait(futures);

    return ActionResult.success;
  }

  /// Clone/Pull a git service
  Future<void> _cloneOrPull(RemoteEntry entry) async {
    final path = Path.join(".wdtk/remotes/", entry.remote, entry.branch);

    Logger.verbose(path);

    if (await Directory(path).exists()) {
      // Already exists, pull
      Logger.verbose("Pulling", indent: Indent());
      var result = await Process.run("git", ["pull"], workingDirectory: path);

      if (result.exitCode != 0) {
        Logger.error("Failed to pull ${entry.remote}");
      }

      if (result.stdout != "") {
        Logger.verbose(result.stdout, indent: Indent(indentSize: 2));
      }

      if (result.stderr != "") {
        Logger.verbose(result.stderr, indent: Indent(indentSize: 2));
      }

      return;
    }

    // Clone
    Logger.verbose("Cloning", indent: Indent());
    await Directory(path).create(recursive: true);

    var result = await Process.run(
        "git", ["clone", "-b", entry.branch, "https://${entry.remote}", "."],
        workingDirectory: path);

    if (result.exitCode != 0) {
      Logger.error("Failed to clone ${entry.remote}");
    }

    if (result.stdout != "") {
      Logger.verbose(result.stdout, indent: Indent(indentSize: 2));
    }

    if (result.stderr != "") {
      Logger.verbose(result.stderr, indent: Indent(indentSize: 2));
    }
  }
}

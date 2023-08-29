part of wdtk_commands;

/// Utility class
class Utility {
  /// Copy directory
  static Future<void> copyDirectory(Directory source, Directory destination) async {
    destination.create(recursive: true);

    source
        .list(recursive: false)
        .forEach((var entity) async {
      if (entity is Directory) {
        var newDirectory = Directory(
            Path.join(destination.absolute.path, Path.basename(entity.path)));
        newDirectory.createSync(recursive: true);

        await copyDirectory(entity.absolute, newDirectory);
      } else if (entity is File) {
        entity
            .copySync(Path.join(destination.path, Path.basename(entity.path)));
      }
    });
  }
}

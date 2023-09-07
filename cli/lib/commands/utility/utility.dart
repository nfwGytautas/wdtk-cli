part of wdtk_commands;

/// Utility class
class Utility {
  /// Copy directory
  static void copyDirectory(Directory source, Directory destination) {
    destination.createSync(recursive: true);

    source
        .listSync(recursive: false)
        .forEach((var entity) {
      if (entity is Directory) {
        var newDirectory = Directory(
            Path.join(destination.absolute.path, Path.basename(entity.path)));
        newDirectory.createSync(recursive: true);

        copyDirectory(entity.absolute, newDirectory);
      } else if (entity is File) {
        entity
            .copySync(Path.join(destination.path, Path.basename(entity.path)));
      }
    });
  }
}

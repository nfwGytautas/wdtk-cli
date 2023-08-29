part of wdtk_templates;

/// Arguments for writing a language template
class LanguageTemplateArgs {
  final String rootPath;
  final String domain;
  final String service;

  LanguageTemplateArgs({required this.rootPath, required this.domain, required this.service});
}

/// Base class for language templates
abstract class LanguageTemplate {
  /// Write a template to the specified directory
  Future<void> write(LanguageTemplateArgs args);

  /// Create language template from string
  static LanguageTemplate? fromString(String language) {
    if (language == "go") {
      return GoLanguageTemplate();
    }

    return null;
  }
}

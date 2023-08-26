package language

import "errors"

type LanguageTemplateData struct {
	Root        string
	Directory   string
	ServiceName string
}

func Template(data LanguageTemplateData, lang string) error {
	if lang == "go" {
		return writeGolangTemplate(data)
	}

	return errors.New("unsupported language " + lang)
}

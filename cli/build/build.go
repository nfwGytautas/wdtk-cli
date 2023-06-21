package build

import "errors"

type BuildData struct {
	SourceDir   string
	OutDir      string
	ServiceName string
}

func Build(data BuildData, lang string) error {
	if lang == "go" {
		return buildGo(data)
	}

	return errors.New("unsupported language " + lang)
}

package build

import "errors"

type ServiceBuildData struct {
	SourceDir   string
	OutDir      string
	ServiceName string
}

type FrontendBuildData struct {
	Type   string
	OutDir string
}

func BuildService(data ServiceBuildData, lang string) error {
	if lang == "go" {
		return buildGo(data)
	}

	return errors.New("unsupported language " + lang)
}

func BuildFrontend(data FrontendBuildData, toolchain string) error {
	if toolchain == "flutter" {
		return buildFlutter(data)
	}

	return errors.New("unsupported toolchain " + toolchain)
}

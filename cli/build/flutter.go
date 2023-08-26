package build

import (
	"errors"
	"os"

	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/util"
)

func buildFlutter(data FrontendBuildData) error {
	if data.Type == "web" {
		return flutterBuildWeb(data)
	}

	return errors.New("unsupported type " + data.Type + " for flutter toolchain")
}

func flutterBuildWeb(data FrontendBuildData) error {
	println(util.SPACING_2 + "Running 'flutter build web'")
	err := util.ExecuteCommand(util.Command{
		Command:        "flutter",
		Args:           []string{"build", "web"},
		Directory:      "frontend/_flutter/",
		PrintToConsole: true,
	})

	if err != nil {
		return err
	}

	// Copy to .wdtk/bin/frontends/
	if !file.Exists(".wdtk/bin/frontends/web") {
		err := os.Mkdir(".wdtk/bin/frontends/web", os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Copy directory
	err = file.CopyDirectory("frontend/_flutter/build/web/", ".wdtk/bin/frontends/web/")
	if err != nil {
		return err
	}

	return nil
}

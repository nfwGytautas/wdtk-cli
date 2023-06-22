package build

import "github.com/nfwGytautas/webdev-tk/cli/util"

func buildGo(data ServiceBuildData) error {
	var err error

	err = modTidy(data)
	if err != nil {
		return err
	}

	err = goGet(data)
	if err != nil {
		return err
	}

	err = goBuild(data)
	if err != nil {
		return err
	}

	return nil
}

func modTidy(data ServiceBuildData) error {
	println(util.SPACING_2 + "Running 'go mod tidy'")
	return util.ExecuteCommand(util.Command{
		Command:        "go",
		Args:           []string{"mod", "tidy"},
		Directory:      data.SourceDir,
		PrintToConsole: true,
	})
}

func goGet(data ServiceBuildData) error {
	println(util.SPACING_2 + "Running 'go get ./'")
	return util.ExecuteCommand(util.Command{
		Command:        "go",
		Args:           []string{"get", "./"},
		Directory:      data.SourceDir,
		PrintToConsole: true,
	})
}

func goBuild(data ServiceBuildData) error {
	println(util.SPACING_2 + "Running 'go build -o " + data.OutDir + data.ServiceName + " .'")
	return util.ExecuteCommand(util.Command{
		Command:        "go",
		Args:           []string{"build", "-o", data.OutDir + data.ServiceName, "."},
		Directory:      data.SourceDir,
		PrintToConsole: true,
	})
}

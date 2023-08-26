package deploy

import (
	"os"

	"github.com/nfwGytautas/gdev/file"
)

func DeployLocal(data DeployData) error {
	// Local deployment is just a copy

	// Check if dir exists
	if !file.Exists(data.OutputDir) {
		// Create
		err := os.MkdirAll(data.OutputDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	// Copy service config
	err := file.CopyFile(data.ConfigFile, data.OutputDir+"/assets/"+data.ConfigFileName)
	if err != nil {
		return err
	}

	// Copy executable
	if data.ServiceName[len(data.ServiceName)-1:] == "/" {
		// Directory
		err = file.CopyDirectory(data.InputDir+data.ServiceName, data.OutputDir+"/")
		if err != nil {
			return err
		}
	} else {
		// File
		err = file.CopyFile(data.InputDir+data.ServiceName, data.OutputDir+"/"+data.ServiceName)
		if err != nil {
			return err
		}
	}

	return nil
}

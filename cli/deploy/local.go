package deploy

import (
	"fmt"
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
	serviceConfigPath := fmt.Sprintf("deploy/generated/%s_ServiceConfig_%s.json", data.ServiceName, data.DeploymentName)
	err := file.CopyFile(serviceConfigPath, data.OutputDir+"/ServiceConfig.json")
	if err != nil {
		return err
	}

	// Copy executable
	err = file.CopyFile("deploy/bin/"+data.ServiceName, data.OutputDir+"/"+data.ServiceName)
	if err != nil {
		return err
	}

	return nil
}

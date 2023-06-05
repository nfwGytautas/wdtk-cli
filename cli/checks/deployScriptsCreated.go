package checks

import (
	"fmt"
	"os"
	"strings"

	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/templates"
	"github.com/nfwGytautas/webdev-tk/cli/types"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

// Checks if all services are created that have been specified
func DeployScriptsExist(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("📦  Creating deployment scripts")

	for _, service := range cfg.Services {
		// Doesn't exist create
		stats.NumCreatedDeployScripts++

		err := createUnixBuildScript(service)
		if err != nil {
			return err
		}

		for _, deployment := range cfg.Deployments {
			filled, err := cfg.GetFilledDeployment(service, deployment.Name)
			if err != nil {
				return err
			}

			err = createDeploymentScript(service, filled)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================

func createUnixBuildScript(service types.ServiceDescriptionConfig) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	data := templates.UNIXDeployData{
		ServiceName: service.Name,
		RootDir:     currentDir,
		ServiceLang: service.Language,
	}

	outFile := fmt.Sprintf("deploy/unix/%s_BUILD_UNIX.sh", service.Name)

	err = file.WriteTemplate(outFile, templates.UnixHeaderDeployTemplate, data)
	if err != nil {
		return err
	}

	if service.Language == "go" {
		goBuildData := templates.GoBuildData{
			ServiceName: service.Name,
			SourceDir:   currentDir + "/services/" + service.Name + "/",
			OutDir:      currentDir + "/deploy/bin/unix/",
		}

		err = file.AppendTemplate(outFile, templates.GoBuildDeployTemplate, goBuildData)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDeploymentScript(service types.ServiceDescriptionConfig, deployment types.DeploymentConfig) error {
	// TODO: Remote deploy

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	data := templates.UNIXDeployData{
		ServiceName: service.Name,
		RootDir:     currentDir,
	}

	outFile := fmt.Sprintf("deploy/unix/%s_DEPLOY_%s.sh", service.Name, deployment.Name)

	err = file.WriteTemplate(outFile, templates.UnixHeaderDeployTemplate, data)
	if err != nil {
		return err
	}

	rootDeploymentDirectory := strings.Replace(*deployment.DeployDir, "%serviceName", service.Name, -1)

	err = os.MkdirAll(rootDeploymentDirectory, os.ModePerm)
	if err != nil {
		return err
	}

	deploymentData := templates.DeployData{
		Deployment: deployment.Name,
		InFile:     fmt.Sprintf("../bin/unix/%s", service.Name),
		OutDir:     rootDeploymentDirectory,
	}

	err = file.AppendTemplate(outFile, templates.LocalDeployTemplate, deploymentData)
	if err != nil {
		return err
	}

	return nil
}

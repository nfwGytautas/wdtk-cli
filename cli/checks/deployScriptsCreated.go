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

	balancerLang := ""
	serviceLang := ""

	if service.Source.Balancer != nil {
		balancerLang = service.Source.Balancer.Language
	}

	if service.Source.Service != nil {
		serviceLang = service.Source.Service.Language
	}

	data := templates.UNIXDeployData{
		ServiceName:  service.Name,
		RootDir:      currentDir,
		BalancerLang: balancerLang,
		ServiceLang:  serviceLang,
	}

	outFile := fmt.Sprintf("deploy/unix/%s_BUILD_UNIX.sh", service.Name)

	err = file.WriteTemplate(outFile, templates.UnixHeaderDeployTemplate, data)
	if err != nil {
		return err
	}

	if balancerLang == "go" {
		goBuildData := templates.GoBuildData{
			BuildName:   "balancer",
			ServiceName: service.Name,
			SourceDir:   currentDir + "/services/" + service.Name + "/balancer/",
			OutDir:      currentDir + "/deploy/bin/unix/",
		}

		err = file.AppendTemplate(outFile, templates.GoBuildDeployTemplate, goBuildData)
		if err != nil {
			return err
		}
	}

	if serviceLang == "go" {
		goBuildData := templates.GoBuildData{
			BuildName:   "service",
			ServiceName: service.Name,
			SourceDir:   currentDir + "/services/" + service.Name + "/service/",
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

	if service.Source.Balancer != nil {
		deploymentData := templates.DeployData{
			InFile: fmt.Sprintf("../bin/unix/%s_balancer", service.Name),
			OutDir: rootDeploymentDirectory,
		}

		err = file.AppendTemplate(outFile, templates.LocalDeployTemplate, deploymentData)
		if err != nil {
			return err
		}
	}

	if service.Source.Service != nil {
		deploymentData := templates.DeployData{
			InFile: fmt.Sprintf("../bin/unix/%s_service", service.Name),
			OutDir: rootDeploymentDirectory,
		}

		err = file.AppendTemplate(outFile, templates.LocalDeployTemplate, deploymentData)
		if err != nil {
			return err
		}
	}

	return nil
}

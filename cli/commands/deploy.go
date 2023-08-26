package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nfwGytautas/gdev/array"
	"github.com/nfwGytautas/webdev-tk/cli/deploy"
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/nfwGytautas/webdev-tk/cli/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

func DeployCommand() *cli.Command {
	return &cli.Command{
		Flags:     []cli.Flag{},
		Name:      "deploy",
		Usage:     "Deploy services",
		ArgsUsage: "[target] [all|services...]",
		Action:    runDeploy,
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

func runDeploy(ctx *cli.Context) error {
	if ctx.NArg() < 2 {
		println("❌  Deploy command expects a target and either 'all', a service name or a list of services that you want to deploy")
		return nil
	}

	// Read wdtk.yml
	cfg := types.WDTKConfig{}
	err := cfg.Read()
	if err != nil {
		return err
	}

	println("✈️   Deploying...")

	deployment := getDeployment(&cfg, ctx.Args().Get(0))
	if deployment == nil {
		return nil
	}

	// Create deploy log file
	logFile := fmt.Sprintf(".wdtk/logs/%s.deploy.log", time.Now().Format("2006-01-02 15:04:05"))
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	servicesToDeploy := ctx.Args().Slice()[1:]
	numDeployed := 0
	numFailed := 0

	deployAll := ctx.Args().Get(1) == "all"

	for _, service := range cfg.Services {
		shouldDeploy := deployAll || array.IsElementInArray(servicesToDeploy, service.Name)

		filledDeployment, err := cfg.GetFilledServiceDeployment(service, ctx.Args().Get(0))
		if err != nil {
			return err
		}

		if shouldDeploy {
			name := service.Name
			println(util.SPACING_1 + "- [ Service] " + name)

			rootDeploymentDirectory := strings.Replace(*deployment.DeployDir, "%serviceName", name, -1)
			serviceConfigPath := fmt.Sprintf(".wdtk/generated/%s_ServiceConfig_%s.json", service.Name, deployment.Name)

			data := deploy.DeployData{
				InputDir:       ".wdtk/bin/services/",
				ConfigFile:     serviceConfigPath,
				ConfigFileName: "ServiceConfig.json",
				OutputDir:      rootDeploymentDirectory,
				ServiceName:    service.Name,
				DeploymentName: deployment.Name,
				Directory:      true,
			}

			err := deploySingle(data, filledDeployment)
			if err != nil {
				log.Println(err)
				numFailed++
			} else {
				numDeployed++
			}
		}
	}

	if cfg.Frontend != nil {
		for _, frontend := range cfg.Frontend.Platforms {
			shouldDeploy := deployAll || array.IsElementInArray(servicesToDeploy, frontend.Type)

			filledDeployment, err := cfg.GetFilledFrontendDeployment(frontend, ctx.Args().Get(0))
			if err != nil {
				return err
			}

			if shouldDeploy {
				name := cases.Title(language.English, cases.Compact).String(frontend.Type)
				println(util.SPACING_1 + "- [Frontend] " + name)

				rootDeploymentDirectory := strings.Replace(*deployment.DeployDir, "%serviceName", name, -1)
				serviceConfigPath := fmt.Sprintf(".wdtk/generated/%s_FrontendConfig_%s.json", frontend.Type, deployment.Name)

				data := deploy.DeployData{
					InputDir:       ".wdtk/bin/frontends/",
					ConfigFile:     serviceConfigPath,
					ConfigFileName: "FrontendConfig.json",
					OutputDir:      rootDeploymentDirectory,
					ServiceName:    frontend.Type + "/",
					DeploymentName: deployment.Name,
					Directory:      true,
				}

				err := deploySingle(data, filledDeployment)
				if err != nil {
					log.Println(err)
					numFailed++
				} else {
					numDeployed++
				}
			}
		}
	}

	println(fmt.Sprintf("--- %d deployed, %d failed ---", numDeployed, numFailed))

	if numFailed != 0 {
		return errors.New("one or more deployments failed")
	}

	return nil
}

func getDeployment(cfg *types.WDTKConfig, target string) *types.DeploymentConfig {
	failMessage := ""

	for _, deployment := range cfg.Deployments {
		if deployment.Name == target {
			return &deployment
		}

		failMessage += deployment.Name + ","
	}

	println("❌  Unknown target '" + target + "' valid options: [" + failMessage[:len(failMessage)-1] + "]")

	return nil
}

func deploySingle(data deploy.DeployData, deployment types.DeploymentConfig) error {
	// Check if this is a local deployment or remote
	if *deployment.IP == "127.0.0.1" || *deployment.IP == "localhost" {
		return deploy.DeployLocal(data)
	}

	return deploy.DeployRemote(data)
}

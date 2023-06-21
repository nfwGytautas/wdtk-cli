package scaffold

import (
	"fmt"
	"os"

	"github.com/nfwGytautas/gdev/array"
	"github.com/nfwGytautas/gdev/file"
	"github.com/nfwGytautas/webdev-tk/cli/templates/language"
	"github.com/nfwGytautas/webdev-tk/cli/types"
	"github.com/nfwGytautas/webdev-tk/cli/util"
)

// PUBLIC FUNCTIONS
// ========================================================================
func CreateLocalServices(cfg types.WDTKConfig, stats *types.ServiceCheckStats) error {
	println("🗂️   Creating missing local services")

	serviceNames := []string{}

	for _, service := range cfg.GetServicesOfType(types.SERVICE_TYPE_LOCAL) {
		serviceNames = append(serviceNames, service.Name)
		path := fmt.Sprintf("services/%s", service.Name)

		if !file.Exists(path) {
			fmt.Printf(util.SPACING_1+"- Creating %s\n", service.Name)

			// Doesn't exist, create
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}

			data := language.LanguageTemplateData{
				Directory:   path,
				ServiceName: service.Name,
				Root:        cfg.Package,
			}

			err = language.Template(data, *service.Source.Language)
			if err != nil {
				return err
			}

			stats.NumCreatedServices++
		}
	}

	dirs, err := file.GetDirectories("services/")
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if !array.IsElementInArray(serviceNames, dir) {
			stats.UnusedServices = append(stats.UnusedServices, dir)
		}
	}

	return nil
}

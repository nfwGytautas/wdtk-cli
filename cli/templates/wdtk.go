package templates

// PUBLIC TYPES
// ========================================================================

// Data for balancer template
type WDTKTemplateData struct {
	CurrentDir  string
	ProjectName string
}

// Template for balancer
const WDTKTemplate = `
# Generic project settings
package: {{.ProjectName}}.com/{{.ProjectName}}/
name: {{.ProjectName}}

# List of valid deployment targets, every service needs to have all of these defined in their deployment tag
deployments:
  - name: dev
    # You can define defaults for a target here
    ip: 127.0.0.1
    dir: {{.CurrentDir}}/dev/%serviceName
	apiKey: API_KEY_GOES_HERE

# Services array must define a service with the name 'Authentication' and name 'Gateway'
services:
  # wdtk_service is a reserved keyword, which means that the service is going to be taken from wdtk-services repository
  - name: Gateway
    type: wdtk
    deployment:
      - name: dev
        port: 8080

  - name: Authentication
    type: wdtk
    deployment:
      - name: dev
        port: 8081
        # Config key can be used for additional configuration options these will be stored inside the generated service config files
        config:
          connectionString: "user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"

  # Describe services here
  - name: ExampleService
    type: service
    language: go
    deployment:
      - name: dev
        port: 8090
`

// Template for README.md in frontend directory
const FrontendReadME = `
# Frontend
Directory for all supported frontends
`

// Template for README.md in root directory
const RootReadME = `
# Project
Project description here

## WDTK
A project utilizing WebDev Toolkit https://github.com/nfwGytautas/webdev-tk
`

// Template for .gitignore in root directory
const GitIgnore = `
# deploy related directories
deploy/
`

// Template for update go mods template
const UnixUpdateGoMods = `
`

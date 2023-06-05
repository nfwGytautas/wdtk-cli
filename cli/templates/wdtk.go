package templates

// PUBLIC TYPES
// ========================================================================

// Data for balancer template
type WDTKTemplateData struct {
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
    dir: ~/{{.ProjectName}}/dev/%serviceName

# Gateway settings
apiGateway:
  # Describe gateway deployments
  deployment:
    - name: dev
      port: 8080

# Authentication service
authentication:
  deployment:
    - name: dev
      connectionString: "user:password@tcp(127.0.0.1:3306)/database?charset=utf8mb4&parseTime=True&loc=Local"

# Services
services:
  - name: ExampleService
    source:
      service:
        language: go
      balancer: # <- Can also specify null to indicate that the service has no balancer
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

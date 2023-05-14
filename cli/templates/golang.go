package templates

// PUBLIC TYPES
// ========================================================================

// Data for go.mod template
type GoModFileData struct {
	Root        string
	ServiceName string
	GoVersion   string
	Suffix      string
}

// Template for go.mod file
const GoModTemplate = `
module {{.Root}}services/{{.ServiceName}}/{{.Suffix}}

go {{.GoVersion}}
`

package types

// PUBLIC TYPES
// ========================================================================

// Struct for holding service check stats
type ServiceCheckStats struct {
	NumCreatedServices      int
	NumModifiedServices     int
	NumCreatedDeployScripts int
	UnusedServices          []string
}

// Scaffold command action
type ScaffoldAction func(cfg WDTKConfig, stats *ServiceCheckStats) error

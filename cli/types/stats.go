package types

// PUBLIC TYPES
// ========================================================================

// Struct for holding service check stats
type ServiceCheckStats struct {
	NumCreatedServices  int
	NumModifiedServices int
	UnusedServices      []string
}

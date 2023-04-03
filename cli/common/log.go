package common

import "log"

// PUBLIC TYPES
// ========================================================================

const (
	LOG_LEVEL_NONE  = iota
	LOG_LEVEL_ERROR = iota
	LOG_LEVEL_INFO  = iota
	LOG_LEVEL_TRACE = iota
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_ALL   = iota
)

/*
Log level
*/
var LogLevel = LOG_LEVEL_TRACE

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Utility log function for debug messages
*/
func LogDebug(fmt string, v ...any) {
	if LogLevel >= LOG_LEVEL_DEBUG {
		log.Printf(fmt, v...)
	}
}

/*
Utility log function for trace messages
*/
func LogTrace(fmt string, v ...any) {
	if LogLevel >= LOG_LEVEL_TRACE {
		log.Printf(fmt, v...)
	}
}

/*
Utility log function for info messages
*/
func LogInfo(fmt string, v ...any) {
	if LogLevel >= LOG_LEVEL_INFO {
		log.Printf(fmt, v...)
	}
}

/*
Utility log function for error messages
*/
func LogError(fmt string, v ...any) {
	if LogLevel >= LOG_LEVEL_ERROR {
		log.Printf(fmt, v...)
	}
}

/*
Utility log function for panics
*/
func LogPanic(fmt string, v ...any) {
	log.Panicf(fmt, v...)
}

// PRIVATE FUNCTIONS
// ========================================================================

package common

import (
	"runtime"
	"time"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Panics if err is not nil
*/
func PanicOnError(err error, message string) {
	if err != nil {
		LogError(err.Error())
		LogPanic(message)
	}
}

/*
Time the execution time of a function
*/
func TimeFn(name string) func() {
	start := time.Now()
	return func() {
		LogDebug("%s took %v\n", name, time.Since(start))
	}
}

/*
Time the function that called this method
*/
func TimeCurrentFn() func() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return TimeFn(frame.Function)
}

// PRIVATE FUNCTIONS
// ========================================================================

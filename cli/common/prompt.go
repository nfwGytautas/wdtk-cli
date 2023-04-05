package common

import (
	"bufio"
	"fmt"
	"os"
)

// PUBLIC TYPES
// ========================================================================

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
A prompt to the user asking for confirmation. Returns true if confirmed, false otherwise or error on failure to open console for input
*/
func YNPrompt(message string, defaultValue bool) (bool, error) {
	if defaultValue {
		fmt.Printf("%s [Y/n]", message)
	} else {
		fmt.Printf("%s [y/N]", message)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')

	if err != nil {
		return defaultValue, err
	}

	value := string([]byte(input)[0])

	if value == "" {
		return defaultValue, nil
	} else {
		return value == "y", nil
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

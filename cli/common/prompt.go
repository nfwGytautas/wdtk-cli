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
	for {
		if defaultValue {
			fmt.Printf("%s [Y/n] ", message)
		} else {
			fmt.Printf("%s [y/N] ", message)
		}

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')

		if err != nil {
			return defaultValue, err
		}

		value := string([]byte(input)[0])
		if value == "y" || value == "n" {
			return value == "y", nil
		} else if value == "\n" {
			return defaultValue, nil
		}
	}
}

// PRIVATE FUNCTIONS
// ========================================================================

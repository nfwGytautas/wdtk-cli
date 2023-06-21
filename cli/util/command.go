package util

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/nfwGytautas/gdev/file"
)

type Command struct {
	Command        string
	Args           []string
	Directory      string
	PrintToConsole bool
	LogFile        string
}

func ExecuteCommand(command Command) error {
	var outb, errb bytes.Buffer

	cmd := exec.Command(command.Command, command.Args...)
	cmd.Dir = command.Directory

	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()

	if command.LogFile != "" {
		file.Append(command.LogFile, outb.String())
		file.Append(command.LogFile, errb.String())

		if err != nil {
			file.Append(command.LogFile, err.Error())
		}
	}

	if command.PrintToConsole {
		if outb.String() != "" {
			fmt.Println(outb.String())
		}

		if errb.String() != "" {
			fmt.Println(errb.String())
		}

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return err
}

package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

var ErrEmptyCommandName = errors.New("undefined command name")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
//
//nolint:unparam
func RunCmd(cmd []string, env Environment) int {
	if len(cmd) == 0 {
		panic(ErrEmptyCommandName)
	}

	for key, value := range env {
		if value.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				fmt.Println(err)
			}
		}
		err := os.Setenv(key, value.Value)
		if err != nil {
			fmt.Println(err)
		}
	}

	//nolint:gosec
	c := exec.Command(cmd[0], cmd[1:4]...)

	c.Stdout = os.Stdout

	if err := c.Run(); err != nil {
		panic(err)
	}

	return c.ProcessState.ExitCode()
}

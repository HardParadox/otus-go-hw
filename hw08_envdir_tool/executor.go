package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
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

	c := exec.Command(cmd[0], cmd[1:4]...)

	c.Stdout = os.Stdout

	if err := c.Run(); err != nil {
		panic(err)
	}

	return c.ProcessState.ExitCode()
}

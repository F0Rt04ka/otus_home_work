package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {

	for envName, envVal := range env {
		if envVal.NeedRemove {
			err := os.Unsetenv(envName)
			if err != nil {
				return 1
			}
			continue
		}

		err := os.Setenv(envName, envVal.Value)
		if err != nil {
			return 1
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode()
		}
	}

	return
}

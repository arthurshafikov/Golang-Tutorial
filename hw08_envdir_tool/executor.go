package main

import (
	"errors"
	"os"
	"os/exec"
)

func setEnvVariables(env Environment) {
	for key, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(key)
			continue
		}

		os.Setenv(key, envValue.Value)
	}
}

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	setEnvVariables(env)
	name := cmd[0]
	args := cmd[1:]
	command := exec.Command(name, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Run()
	if err := command.Run(); err != nil {
		var e *exec.ExitError
		if errors.Is(err, e) {
			//nolint:errorlint
			return err.(*exec.ExitError).ExitCode()
		}
	}
	return 0
}

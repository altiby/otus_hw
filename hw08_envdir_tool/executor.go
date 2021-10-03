package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	systemEnv := os.Environ()
	program := cmd[0]
	args := cmd[1:]
	command := exec.Command(program, args...)
	command.Stdout = os.Stdout
	command.Env = make([]string, 0, len(systemEnv)+len(env))
	command.Env = os.Environ()
	for envKey, envValue := range env {
		if envValue.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=", envKey))
		} else {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", envKey, envValue.Value))
		}
	}
	err := command.Run()
	if err != nil {
		exitError := exec.ExitError{}
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
	}
	return -1
}

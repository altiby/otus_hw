package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	systemEnv := os.Environ()
	command := exec.Command(cmd[0], cmd[1:]...)
	command.Stdout = os.Stdout
	command.Env = make([]string, 0, len(systemEnv)+len(env))
	command.Env = append(os.Environ())
	for envKey, envValue := range env {
		if envValue.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=", envKey))
		} else {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", envKey, envValue.Value))
		}
	}
	err := command.Run()
	if err != nil {
		return err.(*exec.ExitError).ExitCode()
	}
	return 0
}

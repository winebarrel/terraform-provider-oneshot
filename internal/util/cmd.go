package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

func ExecCmd(shell string, command string, extraEnvs ...string) (string, string, error) {
	envs, args, err := shellwords.ParseWithEnvs(shell)

	if err != nil {
		return "", "", err
	}

	args = append(args, command)
	envs = append(envs, extraEnvs...)

	var cmd *exec.Cmd

	if len(args) > 1 {
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		cmd = exec.Command(args[0])
	}

	cmd.Env = append(os.Environ(), envs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		return "", "", fmt.Errorf("Failed to execute command: %w\n[STDOUT] %s\n[STDERR] %s\n", err, stdout.String(), stderr.String())
	}

	return stdout.String(), stderr.String(), nil
}

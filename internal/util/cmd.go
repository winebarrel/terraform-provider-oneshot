package util

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/mattn/go-shellwords"
)

type Cmd struct {
	Shell  string
	Stdout string
	Stderr string
}

func NewCmd(shell string) *Cmd {
	cmd := &Cmd{
		Shell: shell,
	}

	return cmd
}

func NewCmdWithLog(shell string, stdout string, stderr string) *Cmd {
	cmd := &Cmd{
		Shell:  shell,
		Stdout: stdout,
		Stderr: stderr,
	}

	return cmd
}

func (c *Cmd) Run(command string, extraEnvs ...string) (string, string, error) {
	envs, args, err := shellwords.ParseWithEnvs(c.Shell)

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

	if c.Stdout != "" {
		f, err := os.OpenFile(c.Stdout, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)

		if err != nil {
			return "", "", err
		}

		defer f.Close()
		cmd.Stdout = io.MultiWriter(&stdout, f)
	} else {
		cmd.Stdout = &stdout
	}

	if c.Stderr != "" {
		f, err := os.OpenFile(c.Stderr, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0600)

		if err != nil {
			return "", "", err
		}

		defer f.Close()
		cmd.Stderr = io.MultiWriter(&stderr, f)
	} else {
		cmd.Stderr = &stderr
	}

	err = cmd.Run()

	if err != nil {
		return "", "", fmt.Errorf("Failed to execute command: %w\n[STDOUT] %s\n[STDERR] %s\n", err, stdout.String(), stderr.String())
	}

	return stdout.String(), stderr.String(), nil
}

package docker

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
)


var (
	logger *zap.SugaredLogger
)

func init() {
	var err error
	devLogger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger = devLogger.Sugar()
}


func Command(command string, args ...string) *containerCmd {
	return &containerCmd{
		command:  command,
		args:     args,
	}
}

// containerCmd implements exec.Cmd for docker containers
type containerCmd struct {
	nameOrID string // the container name or ID
	command  string
	args     []string
	env      []string
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
}

func (c *containerCmd) Run() error {
	var args []string

	// set env
	for _, env := range c.env {
		args = append(args, "-e", env)
	}
	args = append(
		args,
		// finally, with the caller args
		c.args...,
	)

	cmd := exec.Command("docker", args...)
	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}
	if c.stderr != nil {
		cmd.Stderr = c.stderr
	}
	if c.stdout != nil {
		cmd.Stdout = c.stdout
	}
	return cmd.Run()
}

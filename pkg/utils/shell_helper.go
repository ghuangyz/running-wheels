package utils

import (
	"bytes"
	"github.com/ghuangyz/running-wheels/pkg/types"
	"io"
	"os/exec"
)

func RunCommandGroup(
	cg []*types.Command,
	output *bytes.Buffer,
	stderr *bytes.Buffer,
	usePipe bool,
) []error {
	var errors []error
	if len(cg) == 1 || !usePipe {
		for i := 0; i < len(cg); i++ {
			cmd := exec.Command(cg[i].Name, cg[i].Arguments...)
			cmd.Stdout = output
			cmd.Stderr = stderr
			err := cmd.Run()
			if err != nil {
				errors = append(errors, err)
				break
			}
		}
	} else {
		return runCommandsWithPipe(cg, output, stderr)
	}
	return errors
}

func runCommandsWithPipe(
	cg []*types.Command,
	output *bytes.Buffer,
	stderr *bytes.Buffer,
) []error {
	var errors []error
	count := len(cg)
	commands := make([]*exec.Cmd, count)
	writers := make([]*io.PipeWriter, count)
	for i, command := range cg {
		commands[i] = exec.Command(command.Name, command.Arguments...)
		if i > 0 {
			r, w := io.Pipe()
			writers[i-1] = w
			commands[i-1].Stdout = w
			commands[i].Stdin = r
			commands[i].Stderr = stderr
		}

		if i == count-1 {
			commands[i].Stdout = output
		}
	}

	for i := 0; i < count; i++ {
		err := commands[i].Start()
		if err != nil {
			errors = append(errors, err)
		}
	}

	for i := 0; i < count; i++ {
		err := commands[i].Wait()
		if err != nil {
			errors = append(errors, err)
		}
		if writers[i] != nil {
			writers[i].Close()
		}
	}

	return errors
}

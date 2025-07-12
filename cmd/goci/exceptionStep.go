package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type exceptionStep struct {
	name    string
	exe     string
	args    []string
	message string
	proj    string
}

func (s exceptionStep) execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	if out.Len() > 0 {
		return "", &stepErr{
			step:  s.name,
			msg:   fmt.Sprintf("invalid format:%s", out.String()),
			cause: nil,
		}
	}

	return s.message, nil
}

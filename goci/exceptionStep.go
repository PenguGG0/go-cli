package main

import (
	"bytes"
	"os/exec"
)

type exceptionStep struct {
	name    string
	exe     string
	message string
	proj    string
	args    []string
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
			msg:   "invalid format:" + out.String(),
			cause: nil,
		}
	}

	return s.message, nil
}

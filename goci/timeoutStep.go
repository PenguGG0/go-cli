package main

import (
	"context"
	"errors"
	"os/exec"
	"time"
)

var command = exec.CommandContext

type timeoutStep struct {
	name    string
	exe     string
	message string
	proj    string
	args    []string
	timeout time.Duration
}

func (s timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := command(ctx, s.exe, s.args...)

	cmd.Dir = s.proj

	if err := cmd.Run(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return "", &stepErr{
				step:  s.name,
				msg:   "failed time out",
				cause: context.DeadlineExceeded,
			}
		}

		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}

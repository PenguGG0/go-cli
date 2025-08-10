package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
	ErrSignal     = errors.New("received signal")
)

type stepErr struct {
	cause error
	step  string
	msg   string
}

func (s *stepErr) Error() string {
	return fmt.Sprintf("Step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *stepErr) Is(target error) bool {
	var t *stepErr
	if !errors.As(target, &t) {
		return false
	}

	return s.step == t.step
}

func (s *stepErr) Unwrap() error {
	return s.cause
}

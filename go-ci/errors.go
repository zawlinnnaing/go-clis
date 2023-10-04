package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation  = errors.New("validation failed")
	ErrSignal      = errors.New("signal received")
	ErrEmptyDir    = errors.New("empty directory")
	ErrInvalidFile = errors.New("invalid file")
)

type StepError struct {
	step  string
	msg   string
	cause error
}

func (s *StepError) Error() string {
	return fmt.Sprintf("step: %q: %s: Cause: %v", s.step, s.msg, s.cause)
}

func (s *StepError) Is(target error) bool {
	t, ok := target.(*StepError)
	if !ok {
		return s.cause == target
	}

	return t.step == s.step
}

func (s *StepError) UnWrap() error {
	return s.cause
}

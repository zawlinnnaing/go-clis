package main

import (
	"errors"
	"fmt"
)

var (
	ErrValidation = errors.New("validation failed")
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
		return false
	}

	return t.step == s.step
}

func (s *StepError) UnWrap() error {
	return s.cause
}

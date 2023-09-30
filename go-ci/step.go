package main

import (
	"fmt"
	"os/exec"
)

type step struct {
	name    string
	exe     string
	args    []string
	message string
	project string
}

func (s *step) execute() (string, error) {
	if s.project == "" {
		return "", fmt.Errorf("project directory is required: %w", ErrValidation)
	}
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.project
	if err := cmd.Run(); err != nil {
		return "", &StepError{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	return s.message, nil
}

func NewStep(name, exe, message, project string, args []string) *step {
	return &step{
		name:    name,
		exe:     exe,
		args:    args,
		message: message,
		project: project,
	}
}

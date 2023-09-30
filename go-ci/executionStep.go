package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

type executionStep struct {
	step
}

func (execStep *executionStep) execute() (string, error) {
	cmd := exec.Command(execStep.exe, execStep.args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = execStep.project
	if err := cmd.Run(); err != nil {
		return "", &StepError{
			step:  execStep.name,
			msg:   fmt.Sprintf("failed to execute step: %s", execStep.name),
			cause: err,
		}
	}
	if out.Len() > 0 {
		return "", &StepError{
			step:  execStep.name,
			msg:   fmt.Sprintf("invalid output: %s", out.String()),
			cause: nil,
		}
	}
	return execStep.message, nil
}

func NewExecutionStep(name, exe, message, project string, args []string) *executionStep {
	return &executionStep{
		step: *NewStep(name, exe, message, project, args),
	}
}

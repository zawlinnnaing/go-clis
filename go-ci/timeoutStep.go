package main

import (
	"context"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

var command = exec.CommandContext

func (timeoutStep *timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutStep.timeout)
	defer cancel()
	cmd := command(ctx, timeoutStep.exe, timeoutStep.args...)
	cmd.Dir = timeoutStep.project
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", &StepError{
				msg:   "timeout exceeded",
				step:  timeoutStep.name,
				cause: context.DeadlineExceeded,
			}
		}
		return "", &StepError{
			msg:   "failed to execute",
			step:  timeoutStep.name,
			cause: err,
		}
	}
	return timeoutStep.message, nil
}

func NewTimeoutStep(name, exe, message, project string, args []string, timeout time.Duration) *timeoutStep {
	timeoutStep := &timeoutStep{
		step: *NewStep(name, exe, message, project, args),
	}
	if timeout == 0 {
		timeoutStep.timeout = 30 * time.Second
	} else {
		timeoutStep.timeout = timeout
	}
	return timeoutStep
}

package cmd

import (
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}

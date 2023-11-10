package pomodoro_test

import (
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {

	}
}

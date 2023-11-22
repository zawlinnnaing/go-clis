//go:build !inmemory

package pomodoro_test

import (
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/repository"
	"os"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	tempF, err := os.CreateTemp("", "pomodoro-")
	if err != nil {
		t.Fatal(err)
	}
	if err = tempF.Close(); err != nil {
		t.Fatal(err)
	}
	dbRepo, err := repository.
		NewSQLite3Repo(tempF.Name())
	if err != nil {
		t.Fatal(err)
	}

	return dbRepo, func() {
		os.Remove(tempF.Name())
	}
}

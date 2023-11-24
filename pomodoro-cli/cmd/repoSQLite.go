//go:build !inmemory
// +build !inmemory

package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/repository"
)

func getRepo() (pomodoro.Repository, error) {
	repo, err := repository.NewSQLite3Repo(viper.GetString("db"))
	fmt.Println("Using database file ===>", viper.GetString("db"))
	if err != nil {
		return nil, err
	}
	return repo, nil
}

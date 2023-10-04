package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type executer interface {
	execute() (string, error)
}

func main() {
	project := flag.String("p", "", "Project directory")
	file := flag.String("f", "", "Config file path")
	flag.Parse()

	_, err := parseFile(*file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := run(*project, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, out io.Writer) error {
	pipeline := []executer{
		NewStep("go build", "go", "Go build: success", project, []string{"build", "."}),
		NewStep("go test", "go", "Go test: success", project, []string{"test", "-v", "."}),
		NewExecutionStep("go format", "gofmt", "Go format: success", project, []string{"-l"}),
		NewTimeoutStep("git push", "git", "Git push: success", project, []string{"push", "origin", "master"}, 10*time.Second),
	}
	sigChan := make(chan os.Signal, 1)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for _, step := range pipeline {
			output, err := step.execute()
			if err != nil {
				errCh <- err
				return
			}
			_, err = fmt.Fprintln(out, output)
			if err != nil {
				errCh <- err
				return
			}
		}
		close(doneCh)
	}()
	for {
		select {
		case received := <-sigChan:
			signal.Stop(sigChan)
			return fmt.Errorf("%s: Exiting: %w", received, ErrSignal)
		case err := <-errCh:
			return err
		case <-doneCh:
			return nil
		}
	}
}

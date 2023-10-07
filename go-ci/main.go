package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

type executer interface {
	execute() (string, error)
}

func main() {
	project := flag.String("p", "", "Project directory")
	file := flag.String("f", "", "Config file path")
	flag.Parse()

	if err := run(*project, *file, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, file string, out io.Writer) error {
	if project == "" && file == "" {
		return fmt.Errorf("must provide either project or file")
	}
	var pipeline []executer
	if file != "" {
		steps, err := parseFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		pipeline = steps
	} else {
		pipeline = createDefaultPipeline(project)
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

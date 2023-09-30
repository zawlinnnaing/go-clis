package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type executer interface {
	execute() (string, error)
}

func main() {
	project := flag.String("p", "", "Project directory")
	flag.Parse()

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
	}
	for _, step := range pipeline {
		output, err := step.execute()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, output)
		if err != nil {
			return err
		}
	}
	return nil
}

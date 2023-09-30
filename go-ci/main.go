package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	project := flag.String("p", "", "Project directory")
	flag.Parse()

	if err := run(*project, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("project directory is required: %w", ErrValidation)
	}
	args := []string{"build", "."}

	cmd := exec.Command("go", args...)
	cmd.Dir = project

	if err := cmd.Run(); err != nil {
		return &StepError{step: "go build", msg: "go build failed", cause: err}
	}

	_, err := fmt.Fprintln(out, "go build success")

	return err
}

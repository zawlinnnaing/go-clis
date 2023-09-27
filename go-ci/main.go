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
		return fmt.Errorf("project directory is required")
	}
	// Need to add error to end of the list so that go doesn't create executable
	args := []string{"build", ".", "error"}

	cmd := exec.Command("go", args...)
	cmd.Dir = project

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	_, err := fmt.Fprintln(out, "go build success")

	return err
}

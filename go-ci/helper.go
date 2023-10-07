package main

import "time"

func createDefaultPipeline(project string) []executer {
	return []executer{
		NewStep("go build", "go", "Go build: success", project, []string{"build", "."}),
		NewStep("go test", "go", "Go test: success", project, []string{"test", "-v", "."}),
		NewExecutionStep("go format", "gofmt", "Go format: success", project, []string{"-l"}),
		NewTimeoutStep("git push", "git", "Git push: success", project, []string{"push", "origin", "master"}, 10*time.Second),
	}
}

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	_, err := exec.LookPath("git")
	if err != nil {
		t.Skip("Git not installed. Skipping test...")
	}
	testCases := []struct {
		proj     string
		output   string
		expErr   error
		name     string
		setUpGit bool
	}{
		{
			proj: "./testdata/tool",
			name: "success", expErr: nil,
			output:   "Go build: success\nGo test: success\nGo format: success\nGit push: success\n",
			setUpGit: true,
		},
		{
			proj:     "./testdata/toolErr",
			name:     "fail",
			expErr:   &StepError{step: "go build"},
			output:   "",
			setUpGit: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setUpGit {
				cleanUp := setUpGit(t, testCase.proj)
				defer cleanUp()
			}
			out := bytes.Buffer{}
			err := run(testCase.proj, &out)
			if testCase.expErr != nil {
				if err == nil {
					t.Errorf("Expected error; received none")
					return
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error %v; received %v", testCase.expErr, err)
					return
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if out.String() != testCase.output {
				t.Errorf("Expected out: %s, received %s", testCase.output, out.String())
			}
		})
	}
}

func setUpGit(t *testing.T, project string) func() {
	t.Helper()
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	tempDir, err := ioutil.TempDir("", "go-ci-test")
	if err != nil {
		t.Fatal(err)
	}
	projectPath, err := filepath.Abs(project)
	if err != nil {
		t.Fatal(err)
	}
	remoteURI := fmt.Sprintf("file://%s", tempDir)
	gitCommandList := []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projectPath, nil},
		{[]string{"remote", "add", "origin", remoteURI}, projectPath, nil},
		{[]string{"add", "."}, projectPath, nil},
		{[]string{"commit", "-m", "test"}, projectPath, []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
		}},
	}
	for _, gitCommand := range gitCommandList {
		command := exec.Command(gitExec, gitCommand.args...)
		command.Dir = gitCommand.dir
		if gitCommand.env != nil {
			command.Env = append(os.Environ(), gitCommand.env...)
		}
		if err := command.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projectPath, ".git"))
	}
}

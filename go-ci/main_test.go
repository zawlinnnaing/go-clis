package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		proj     string
		file     string
		output   string
		expErr   error
		name     string
		setUpGit bool
		mockCmd  func(ctx context.Context, name string, args ...string) *exec.Cmd
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
		{
			proj:     "./testdata/tool",
			name:     "successMock",
			expErr:   nil,
			output:   "Go build: success\nGo test: success\nGo format: success\nGit push: success\n",
			mockCmd:  mockCommandContext,
			setUpGit: false,
		},
		{
			proj:     "./testdata/tool",
			name:     "failTimeout",
			expErr:   context.DeadlineExceeded,
			output:   "",
			setUpGit: false,
			mockCmd:  mockCommandTimeout,
		},
		{
			name:     "successFile",
			file:     "./testdata/valid.yaml",
			expErr:   nil,
			proj:     "",
			output:   "go build success\ngo test success\ngo format success\ngit push success\n",
			setUpGit: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.setUpGit {
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git not installed. Skipping this test...")
				}
				cleanUp := setUpGit(t, testCase.proj)
				defer cleanUp()
			}
			if testCase.mockCmd != nil {
				command = testCase.mockCmd
			}
			out := bytes.Buffer{}
			err := run(testCase.proj, testCase.file, &out)
			if testCase.expErr != nil {
				if err == nil {
					t.Errorf("Expected error; received none")
					return
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error: %q, received: %q", testCase.expErr, err)
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

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("GO_HELPER_TIMEOUT") == "1" {
		time.Sleep(15 * time.Second)
	}
	if os.Args[2] == "git" {
		fmt.Fprintln(os.Stdout, "Everything up-to-date")
		os.Exit(0)
	}
	os.Exit(1)
}

func TestRunKill(t *testing.T) {
	testCases := []struct {
		name    string
		project string
		expErr  error
		sig     syscall.Signal
	}{
		{
			name:    "SIGINT",
			project: "./testdata/tool",
			expErr:  ErrSignal,
			sig:     syscall.SIGINT,
		},
		{
			name:    "SIGTERM",
			project: "./testdata/tool",
			expErr:  ErrSignal,
			sig:     syscall.SIGTERM,
		},
		{
			name:    "SIGQUIT",
			project: "./testdata/tool",
			expErr:  nil,
			sig:     syscall.SIGQUIT,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			command = mockCommandTimeout
			errCh := make(chan error)
			ignoreSig := make(chan os.Signal, 1)
			expSig := make(chan os.Signal, 1)

			signal.Notify(ignoreSig, syscall.SIGQUIT)
			defer signal.Stop(ignoreSig)

			signal.Notify(expSig, testCase.sig)
			defer signal.Stop(expSig)

			go func() {
				errCh <- run(testCase.project, "", ioutil.Discard)
			}()
			go func() {
				time.Sleep(2 * time.Second)
				syscall.Kill(os.Getpid(), testCase.sig)
			}()

			select {
			case err := <-errCh:
				if err == nil {
					t.Errorf("Expected error, received none")
					return
				}
				if !errors.Is(err, testCase.expErr) {
					t.Errorf("Expected error: %q, received error: %q", testCase.expErr, err)
					return
				}
				select {
				case rec := <-expSig:
					if rec != testCase.sig {
						t.Errorf("Expected signal: %v, received signal: %v", testCase.sig, rec)
						return
					}
				default:
					t.Error("Signal not received")
				}
			case <-ignoreSig:
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

func mockCommandContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess"}
	cs = append(cs, exe)
	cs = append(cs, args...)
	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func mockCommandTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cmd := mockCommandContext(ctx, exe, args...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

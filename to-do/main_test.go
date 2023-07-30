package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo-app"
	fileName = ".todo-test.json"
)

func TestMain(m *testing.M) {
	os.Setenv("TODO_FILENAME", fileName)
	fmt.Println("Building tool")
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %v, error: %v", build, err)
		os.Exit(1)
	}

	fmt.Println("Running test")
	result := m.Run()

	fmt.Println("Cleaning up")
	os.Remove(binName)
	os.Remove(fileName)
	os.Unsetenv("TODO_FILENAME")

	os.Exit(result)
}

func TestCLI(t *testing.T) {
	task := "New Task to add"
	dir, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddTaskFromArguments", func(t *testing.T) {
		err := exec.Command(cmdPath, "-add", task).Run()
		if err != nil {
			fmt.Printf("Error: %v", err)
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"

	t.Run("AddTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdStdIn, task2)
		cmdStdIn.Close()
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		out, err := exec.Command(cmdPath, "-list").CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("1: %s\n2: %s\n", task, task2)
		if string(out) != expected {
			t.Errorf("Expected %v, Received %v", expected, string(out))
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		err := exec.Command(cmdPath, "-complete", "1").Run()
		if err != nil {
			t.Fatal(err)
		}
		out, err := exec.Command(cmdPath, "-list").CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("X 1: %s\n2: %s\n", task, task2)
		if string(out) != expected {
			t.Errorf("Expected %v, received %v", expected, string(out))
		}
	})
}

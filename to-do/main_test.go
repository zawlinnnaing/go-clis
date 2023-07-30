package main_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo-app"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
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

	os.Exit(result)
}

func TestCLI(t *testing.T) {
	task := "New Task to add"
	dir, err := os.Getwd()

	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddTask", func(t *testing.T) {
		err := exec.Command(cmdPath, "-task", task).Run()
		if err != nil {
			fmt.Printf("Error: %v", err)
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		out, err := exec.Command(cmdPath, "-list").CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("1: %s\n", task)
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
		expected := fmt.Sprintf("X 1: %s\n", task)
		if string(out) != expected {
			t.Errorf("Expected %v, received %v", expected, string(out))
		}
	})
}

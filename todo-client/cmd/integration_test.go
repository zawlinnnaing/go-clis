//go:build integration
// +build integration

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func randomTaskName(t *testing.T) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var p strings.Builder
	for i := 0; i < 32; i++ {
		p.WriteByte(chars[r.Intn(len(chars))])
	}
	return p.String()
}

func TestIntegration(t *testing.T) {
	apiRoot := "http://127.0.0.1:8080"
	if os.Getenv("TODO_API_ROOT") != "" {
		apiRoot = os.Getenv("TODO_API_ROOT")
	}
	today := time.Now().Format("Jan/02")

	taskName := randomTaskName(t)
	taskId := ""

	t.Run("AddTask", func(t *testing.T) {
		args := []string{taskName}
		expOut := fmt.Sprintf("Added task %q to the list.\n", taskName)
		var out bytes.Buffer
		if err := addAction(apiRoot, args, &out); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if out.String() != expOut {
			t.Errorf("Expected output: %q, received: %q", expOut, out.String())
		}
	})
	t.Run("ListTasks", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, apiRoot); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		outList := ""
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), taskName) {
				outList = scanner.Text()
				break
			}
		}
		if outList == "" {
			t.Errorf("Task %s not included in the list.", taskName)
		}
		taskCompletedStatus := strings.Fields(outList)[0]

		if taskCompletedStatus != "-" {
			t.Errorf("Expected status: %q, got %q", "-", taskCompletedStatus)
		}
		taskId = strings.Fields(outList)[1]
	})
	viewRes := t.Run("ViewTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := viewAction(apiRoot, taskId, &out); err != nil {
			t.Fatal(err)
		}
		viewOut := strings.Split(out.String(), "\n")
		if !strings.Contains(viewOut[0], taskName) {
			t.Fatalf("Expected task: %s, recevied: %s", taskName, viewOut[0])
		}
		if !strings.Contains(viewOut[1], today) {
			t.Errorf("Expected created date: %s, recevied: %s", today, viewOut[1])
		}
		if !strings.Contains(viewOut[2], "No") {
			t.Errorf("Expected completed status: %s, recevied: %s", "No", viewOut[2])
		}
	})
	if !viewRes {
		t.Fatal("View task failed, aborting integration tests.")
	}
	t.Run("CompleteTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := completeAction(apiRoot, taskId, &out); err != nil {
			t.Fatal(err)
		}
		expOut := fmt.Sprintf("Item number %s marked as completed.\n", taskId)
		if out.String() != expOut {
			t.Fatalf("Expected output: %s, received: %s", expOut, out.String())
		}
	})
	t.Run("ListCompletedTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, apiRoot); err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		outList := ""
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), taskName) {
				outList = scanner.Text()
				break
			}
		}
		if outList == "" {
			t.Errorf("Task %s not included in the list.", taskName)
		}
		taskCompletedStatus := strings.Fields(outList)[0]

		if taskCompletedStatus != "X" {
			t.Errorf("Expected status: %q, got %q", "-", taskCompletedStatus)
		}
		taskId = strings.Fields(outList)[1]
	})
	t.Run("DeleteTask", func(t *testing.T) {
		deletedId, err := deleteAction(apiRoot, taskId)
		if err != nil {
			t.Fatal(err)
		}
		deletedIdStr := strconv.Itoa(deletedId)
		if deletedIdStr != taskId {
			t.Fatalf("Expected delete id: %q, received: %q", taskId, deletedIdStr)
		}
	})
	t.Run("ListDeleteTask", func(t *testing.T) {
		var out bytes.Buffer
		if err := listAction(&out, apiRoot); err != nil {
			t.Fatal(err)
		}
		scanner := bufio.NewScanner(&out)
		for scanner.Scan() {
			if strings.Contains(taskName, scanner.Text()) {
				t.Errorf("Expected task not to be found. %s", scanner.Text())
				break
			}
		}
	})
}

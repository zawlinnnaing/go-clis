package todo_test

import (
	"os"
	"testing"

	"github.com/zawlinnnaing/go-clis/to-do/todo"
)

func TestAddItem(t *testing.T) {
	list := todo.TaskList{}

	taskName := "New Task"
	list.Add(taskName)

	if list[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, list[0].Task)
	}
}

func TestComplete(t *testing.T) {
	list := todo.TaskList{}

	taskName := "New task"

	list.Add(taskName)

	list.Complete(0)

	if !list[0].Done {
		t.Errorf("Expected task to be done, received %v", list[0].Done)
	}
}

func TestDelete(t *testing.T) {
	list := todo.TaskList{}

	list.Add("Task")

	list.Delete(0)

	if len(list) != 0 {
		t.Errorf("Expected list to be %v, received %v", 0, len(list))
	}

}

func TestSaveLoad(t *testing.T) {
	list1 := todo.TaskList{}
	list2 := todo.TaskList{}

	list1.Add("New Task")

	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Failed to create temp file")
	}

	defer os.Remove(file.Name())

	if err := list1.Save(file.Name()); err != nil {
		t.Fatalf("Failed to save TaskList to file %v", file.Name())
	}

	if err := list2.Load(file.Name()); err != nil {
		t.Fatalf("Failed to load TaskList from file %v", file.Name())
	}

	if list1[0].Task != list2[0].Task {
		t.Errorf("Expected TaskLists to be equal. Received from TaskList 1: %v, and from 2: %v", list1[0].Task, list2[0].Task)
	}

}

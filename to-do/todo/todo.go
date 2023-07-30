package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type TaskList []item

func (list *TaskList) Add(task string) {
	t := item{Task: task, Done: false, CreatedAt: time.Now()}
	*list = append(*list, t)
}

func (list *TaskList) validIndex(index int) bool {
	return index >= 0 && index < len(*list)
}

func (list *TaskList) Complete(index int) error {
	deList := *list
	if !list.validIndex(index) {
		return fmt.Errorf("item %v does not exist", index)
	}

	item := &deList[index]
	item.Done = true
	item.CompletedAt = time.Now()

	return nil
}

func (list *TaskList) Delete(index int) error {
	if !list.validIndex(index) {
		return fmt.Errorf("item %v does not exist", index)
	}

	*list = append((*list)[:index], (*list)[index+1:]...)

	return nil
}

func (list *TaskList) Save(filename string) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

func (list *TaskList) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, list)
}

func (list *TaskList) String() string {
	formatted := ""

	for index, task := range *list {
		prefix := ""
		if task.Done {
			prefix = "X "
		}

		formatted += fmt.Sprintf("%v%v: %v\n", prefix, index+1, task.Task)
	}

	return formatted
}

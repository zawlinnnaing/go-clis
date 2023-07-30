package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zawlinnnaing/go-clis/to-do/todo"
)

const todoFileName = ".todo.json"

func main() {

	task := flag.String("task", "", "Task to be added in the ToDo list")
	showList := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Parse()

	list := &todo.TaskList{}

	if err := list.Load(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *showList:
		for _, item := range *list {
			if item.Done {
				continue
			}
			fmt.Println(item.Task)
		}
	case *complete > 0:
		if err := list.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := list.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		list.Add(*task)
		if err := list.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

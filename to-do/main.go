package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zawlinnnaing/go-clis/to-do/todo"
)

const todoFileName = ".todo.json"

func main() {
	list := &todo.TaskList{}

	if err := list.Load(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case len(os.Args) == 1:
		for _, item := range *list {
			fmt.Println(item.Task)
		}
	default:
		item := strings.Join(os.Args[1:], " ")
		list.Add(item)
		if err := list.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

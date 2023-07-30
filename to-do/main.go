package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/zawlinnnaing/go-clis/to-do/todo"
)

var todoFileName = ".todo.json"

func main() {

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2023\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	add := flag.Bool("add", false, "Task to be added in the ToDo list")
	showList := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	list := &todo.TaskList{}

	if err := list.Load(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *showList:
		fmt.Print(list)
	case *complete > 0:
		if err := list.Complete(*complete - 1); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := list.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, task := range tasks {
			list.Add(task)
		}
		if err := list.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}

func getTask(reader io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		return []string{strings.Join(args, " ")}, nil
	}

	scanner := bufio.NewScanner(reader)
	strs := []string{}
	for scanner.Scan() {
		strs = append(strs, scanner.Text())
	}
	if len(strs) == 0 {
		return []string{}, fmt.Errorf("task cannot be blank")
	}

	return strs, nil
}

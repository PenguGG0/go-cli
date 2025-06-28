// To-Do List CLI
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/PenguGG0/go-cli/internal/todo"
)

func main() {
	var todoFileName = ".todo.json"

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	task := flag.String("task", "", "Task to be included in the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by Pengu_GG\n", filepath.Base(os.Args[0]))
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2025\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage information:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	l := todo.List{}

	// read from the todoFile
	if err := l.Get(todoFileName); err != nil {
		log.Fatalln(err)
	}

	switch {
	// if -list flag is set, list the to-do items
	case *list:
		fmt.Print(l.String())
	// if -complete flag is set, complete the item according to the given value and save
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			log.Fatalln(err)
		}
		if err := l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}
	// if it's followed by a task string, add the task to l and save
	case *task != "":
		l.Add(*task)
		if err := l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}
	// if there's no argument, return an error
	default:
		log.Fatalln("Invalid option")
	}
}

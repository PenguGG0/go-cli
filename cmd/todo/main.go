// To-Do List CLI
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PenguGG0/go-cli/internal/todo"
)

func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	s := bufio.NewScanner(r)
	s.Scan()
	if err := s.Err(); err != nil {
		return "", err
	}
	if len(s.Text()) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}
	return s.Text(), nil
}

func main() {
	var todoFileName = ".todo.json"

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Item to be deleted")

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
	// if -complete flag is set, complete the specified item by number
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			log.Fatalln(err)
		}
		if err := l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}
	// if -del flag is set, delete the specified item by number
	case *del > 0:
		if err := l.Delete(*del); err != nil {
			log.Fatalln(err)
		}
		if err := l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}
	// if it's followed by a task string, add the task to l and save
	case *add:
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			log.Fatalln(err)
		}

		l.Add(t)
		if err = l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}
	// if there's no argument, return an error
	default:
		log.Fatalln("Invalid option")
	}
}

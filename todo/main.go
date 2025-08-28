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
	"time"

	todo "github.com/PenguGG0/go-cli/todo/internal"
)

func getTask(r io.Reader, args ...string) ([]string, error) {
	tasks := make([]string, 0)

	// Deal with command line args first
	if len(args) > 0 {
		tasks = append(tasks, strings.Join(args, " "))

		return tasks, nil
	}

	// If there's no args, read from io.Reader
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		tasks = append(tasks, line)
	}
	if err := scanner.Err(); err != nil {
		return tasks, fmt.Errorf("error scanning input: %w", err)
	}

	return tasks, nil
}

func main() {
	todoFileName := "./tmp/.todo.json"

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	add := flag.Bool("add", false, "Add task to the ToDo list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("del", 0, "Item to be deleted")
	verbose := flag.Bool("verbose", false, "Enable verbose output, showing information like date/time")
	onlyShowPending := flag.Bool("pending", false, "Only display uncompleted items")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed by Pengu_GG\n", filepath.Base(os.Args[0]))
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2025\n")
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nAdd new task: todo -add [task]\n")
	}

	flag.Parse()

	l := todo.List{}

	// read from the todoFile
	if err := l.Get(todoFileName); err != nil {
		log.Fatalln(err)
	}

	// if -verbose flag is set, show information like data/time
	if *verbose {
		fmt.Printf("%v\n", time.Now().Format("2006-01-02/15:04:05"))
	}

	switch {
	// if -list flag is set, list the to-do items
	case *list:
		fmt.Print(l.String(*onlyShowPending))

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
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			log.Fatalln(err)
		}
		for _, t := range tasks {
			l.Add(t)
		}
		if err = l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}

	// if there's no argument, return an error
	default:
		log.Fatalln("Invalid option")
	}
}

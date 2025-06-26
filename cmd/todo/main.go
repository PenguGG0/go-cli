// To-Do List CLI
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PenguGG0/go-cli/internal/todo"
)

const todoFileName = ".todo.json"

func main() {
	l := todo.List{}

	if err := l.Get(todoFileName); err != nil {
		log.Fatalln(err)
	}

	switch {
	// For no extra arguments, list the to-do items
	case len(os.Args) == 1:
		for _, item := range l {
			fmt.Println(item.Task)
		}
	// Concatenate all provided arguments with a space
	// Add it to the list and save the new list
	default:
		item := strings.Join(os.Args[1:], " ")
		l.Add(item)
		if err := l.Save(todoFileName); err != nil {
			log.Fatalln(err)
		}
	}

}

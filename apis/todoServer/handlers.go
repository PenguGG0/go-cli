package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/PenguGG0/go-cli/internal/todo"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("incalid data")
)

func replyTextContent(w http.ResponseWriter, r *http.Request, status int, content string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write([]byte(content))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	content := "There's an API here"
	replyTextContent(w, r, http.StatusOK, content)
}

func todoRouter(todoFile string, l sync.Locker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list := &todo.List{}

		l.Lock()
		defer l.Unlock()

		if err := list.Get(todoFile); err != nil {
			replyError(w, r, http.StatusInternalServerError, err.Error())
			return
		}
	}
}

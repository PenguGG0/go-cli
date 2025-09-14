package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	tempTodoFile, err := os.CreateTemp("", "todotest")
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(newMux(tempTodoFile.Name()))

	for i := 1; i < 3; i++ {
		var body bytes.Buffer

		taskName := fmt.Sprintf("Task number %d.", i)
		item := struct {
			Task string `json:"task"`
		}{
			Task: taskName,
		}

		if err := json.NewDecoder(&body).Decode(&item); err != nil {
			t.Fatal(err)
		}

		r, err := http.Post(ts.URL+"/todo", "application/json", &body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		if r.StatusCode != http.StatusCreated {
			t.Fatalf("Failed to add initial items: Status: %d", r.StatusCode)
		}
	}

	return ts.URL, func() {
		ts.Close()
		os.Remove(tempTodoFile.Name())
	}
}

func TestGet(t *testing.T) {
	testcases := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{
			name:       "GetRoot",
			path:       "/",
			expCode:    http.StatusOK,
			expContent: "There's an API here",
		},
		{
			name:    "NotFound",
			path:    "/todo/500",
			expCode: http.StatusNotFound,
		},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)

			r, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}
			defer r.Body.Close()

			if r.StatusCode != tc.expCode {
				t.Fatalf("Got %q, expected %q.",
					http.StatusText(r.StatusCode),
					http.StatusText(tc.expCode))
			}

			switch {
			case strings.Contains(r.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(r.Body); err != nil {
					t.Error(err)
				}

				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Got %q, expected %q.", string(body), tc.expContent)
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", r.Header.Get("Content-Type"))
			}
		})
	}
}

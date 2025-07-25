package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName      = "todo"
	testFileName = ".todoTest.json"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("TODO_FILENAME", testFileName); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("test: Building tool...")
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		log.Fatalln("Cannot build tool todo:", err)
	}

	fmt.Println("test: Running tests...")
	result := m.Run()

	fmt.Println("test: Cleaning up...")
	_ = os.Remove(binName)
	_ = os.Remove(testFileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	testTask1 := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", testTask1)

		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	testTask2 := "test task number 2"
	testTask3 := "test task number 3"

	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdIn, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}

		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if _, err = io.WriteString(cmdStdIn, testTask2+"\n"+testTask3); err != nil {
			t.Fatal(err)
		}

		if err = cmdStdIn.Close(); err != nil {
			t.Fatal(err)
		}

		if err = cmd.Run(); err != nil {
			t.Fatalf("%v: %v", err, stderr.String())
		}
	})

	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-del", "2")

		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")

		if err = cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		expected := fmt.Sprintf("Done 1: %s\n     2: %s\n", testTask1, testTask3)
		if expected != string(out) {
			t.Errorf("Expected %v, got %v instead\n", expected, string(out))
		}
	})
}

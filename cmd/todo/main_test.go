package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
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
	_ = os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	testTask := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-task", testTask)

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

		expected := fmt.Sprintf("     1: %s\n", testTask)
		if expected != string(out) {
			t.Errorf("Expected %v, got %v instead\n", expected, string(out))
		}
	})
}

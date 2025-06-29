package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	resultFile = "./tmp/test1.md.html"
	goldenFile = "./testdata/test1.md.html"
)

func TestRun(t *testing.T) {
	if err := run(inputFile, true); err != nil {
		t.Fatal(err)
	}

	result, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("golden:\n%s\n", expected)
		t.Logf("result:\n%s\n", result)
		t.Error("Result content does not match golden file")
	}
	err = os.Remove(resultFile)
	if err != nil {
		t.Fatal(err)
	}
}

package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
)

func TestRun(t *testing.T) {
	outName, err := run(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result, err := os.ReadFile(outName)
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
	err = os.Remove(outName)
	if err != nil {
		t.Fatal(err)
	}
}

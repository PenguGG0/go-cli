package main

import (
	"os"
	"testing"
)

func TestValid(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      []string
		minSize  int64
		expected bool
	}{
		{"FilterNoExtension", "testdata/dir.log", []string{""}, 0, true},
		{"FilterExtensionMatch", "testdata/dir.log", []string{".log"}, 0, true},
		{"FilterExtensionNoMatch", "testdata/dir.log", []string{".sh"}, 0, false},
		{"FilterSizeMatch", "testdata/dir.log", []string{".log"}, 10, true},
		{"FilterSizeNoMatch", "testdata/dir.log", []string{".log"}, 20, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			if got := valid(tc.file, tc.ext, tc.minSize, info); got != tc.expected {
				t.Errorf("Got %v, expected %v", got, tc.expected)
			}
		})
	}
}

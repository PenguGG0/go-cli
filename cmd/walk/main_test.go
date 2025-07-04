package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createTempDir(t *testing.T, files map[string]int) (string, func()) {
	// Mark as a test helper
	t.Helper()

	tempDir, err := os.MkdirTemp(os.TempDir(), "walk_test_*")
	if err != nil {
		return "", nil
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fileName := fmt.Sprintf("file%d%s", j, k)
			filePath := filepath.Join(tempDir, fileName)
			if err = os.WriteFile(filePath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	cleanup := func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			t.Fatal(err)
		}
	}

	return tempDir, cleanup
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		app      application
		expected string
	}{
		{
			name:     "NoFilter",
			root:     "testdata",
			app:      application{ext: "", minSize: 0, list: true},
			expected: "testdata\\dir.log\ntestdata\\dir2\\script.sh\n",
		},
		{
			name:     "FilterExtensionMatch",
			root:     "testdata",
			app:      application{ext: ".log", minSize: 0, list: true},
			expected: "testdata\\dir.log\n",
		},
		{
			name:     "FilterExtensionSizeMatch",
			root:     "testdata",
			app:      application{ext: ".log", minSize: 10, list: true},
			expected: "testdata\\dir.log\n",
		},
		{
			name:     "FilterSizeNoMatch",
			root:     "testdata",
			app:      application{ext: ".log", minSize: 20, list: true},
			expected: "",
		},
		{
			name:     "FilterExtensionNoMatch",
			root:     "testdata",
			app:      application{ext: ".gz", minSize: 0, list: true},
			expected: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)

			err := run(tc.root, out, tc.app)
			if err != nil {
				t.Fatal(err)
			}

			gotOut := out.String()
			if gotOut != tc.expected {
				t.Errorf("run() got %v, expected %v", gotOut, tc.expected)
			}
		})
	}
}

func TestRunDel(t *testing.T) {
	testCases := []struct {
		name             string
		extDelete        string
		extNoDelete      string
		nDelete          int
		nNoDelete        int
		expected         string
		expectedLogLines int // nDelete+1
	}{
		{
			name:             "DeleteExtensionNoMatch",
			extDelete:        ".log",
			extNoDelete:      ".gz",
			nDelete:          0,
			nNoDelete:        10,
			expected:         "",
			expectedLogLines: 1,
		},
		{
			name:             "DeleteExtensionMatch",
			extDelete:        ".log",
			extNoDelete:      "",
			nDelete:          10,
			nNoDelete:        0,
			expected:         "",
			expectedLogLines: 11,
		},
		{
			name:             "DeleteExtensionMixed",
			extDelete:        ".log",
			extNoDelete:      ".gz",
			nDelete:          5,
			nNoDelete:        5,
			expected:         "",
			expectedLogLines: 6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)

			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.extDelete:   tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			app := application{ext: tc.extDelete, del: true, wLog: &logBuffer}
			if err := run(tempDir, &buffer, app); err != nil {
				t.Fatal(err)
			}

			if tc.expected != buffer.String() {
				t.Errorf("run() with '-delete' got %v, expected %v", buffer.String(), tc.expected)
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("run() with '-delete' got %v files left, expected %v", len(filesLeft), tc.nNoDelete)
			}

			logLines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(logLines) != tc.expectedLogLines {
				t.Errorf("run() with '-delete' got %v log lines, expected %v", len(logLines), tc.expectedLogLines)
			}
		})
	}
}

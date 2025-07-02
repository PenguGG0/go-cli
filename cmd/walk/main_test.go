package main

import (
	"bytes"
	"testing"
)

func Test_run(t *testing.T) {
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

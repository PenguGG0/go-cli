package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		proj   string
		outStr string
		expErr error
	}{
		{
			name:   "allSuccess",
			proj:   "./testdata/tool",
			outStr: "Go Build: SUCCESS\nGo Test: SUCCESS\nGofmt: SUCCESS\n",
			expErr: nil,
		},
		{
			name:   "buildFail",
			proj:   "./testdata/toolBuildErr",
			outStr: "",
			expErr: &stepErr{step: "go build"},
		},
		{
			name:   "fmtFail",
			proj:   "./testdata/toolFmtErr",
			outStr: "",
			expErr: &stepErr{step: "go fmt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			err := run(tc.proj, &out)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Got error: nil, expected: %q", tc.expErr)
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Got error: %q, expected: %q", err, tc.expErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != tc.outStr {
				t.Errorf("Got output: %q, expected: %q", out.String(), tc.outStr)
			}
		})
	}
}

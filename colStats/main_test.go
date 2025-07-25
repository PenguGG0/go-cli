package main

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name   string
		col    int
		op     string
		exp    string
		files  []string
		expErr error
	}{
		{
			name:   "RunAvg1File",
			col:    3,
			op:     "avg",
			exp:    "227.6\n",
			files:  []string{"./testdata/example.csv"},
			expErr: nil,
		},
		{
			name:   "RunAvg1File",
			col:    3,
			op:     "avg",
			exp:    "233.84\n",
			files:  []string{"./testdata/example.csv", "./testdata/example2.csv"},
			expErr: nil,
		},
		{
			name:   "RunFailRead",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata/example.csv", "./testdata/fakefile.csv"},
			expErr: os.ErrNotExist,
		},
		{
			name:   "RunFailColumn",
			col:    0,
			op:     "avg",
			exp:    "",
			files:  []string{"./testdata/example.csv"},
			expErr: ErrInvalidColumn,
		},
		{
			name:   "RunFailNoFiles",
			col:    2,
			op:     "avg",
			exp:    "",
			files:  []string{},
			expErr: ErrNoFiles,
		},
		{
			name:   "RunFailOperation",
			col:    2,
			op:     "invalid",
			exp:    "",
			files:  []string{"./testdata/example.csv"},
			expErr: ErrInvalidOperation,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var res bytes.Buffer

			err := run(tc.files, tc.op, tc.col, &res)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Got no error, expected error %v", tc.expErr)
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Got error %v, expected error %v", err, tc.expErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if res.String() != tc.exp {
				t.Errorf("Got %v, expected %v", res.String(), tc.exp)
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	fileNames, err := filepath.Glob("./testdata/benchmark/*.csv")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = run(fileNames, "avg", 2, io.Discard); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkMin(b *testing.B) {
	size := 10_000_000
	testData := make([]float64, 0, size)
	for range size {
		testData = append(testData, rand.Float64())
	}

	b.ResetTimer()
	_ = dataMin(testData)
}

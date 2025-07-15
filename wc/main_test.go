package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name       string
		r          io.Reader
		countLines bool
		countBytes bool
		wantCount  int
		wantByte   int
		wantErr    error
	}{
		{
			name:       "four words",
			r:          bytes.NewBufferString("word1 word2 word3 word4\n"),
			countLines: false,
			countBytes: false,
			wantCount:  4,
			wantByte:   0,
			wantErr:    nil,
		},
		{
			name:       "three lines",
			r:          bytes.NewBufferString("word1 word2 word3\nline2\nline3 word4"),
			countLines: true,
			countBytes: false,
			wantCount:  3,
			wantByte:   0,
			wantErr:    nil,
		},
		{
			name:       "three lines and 38 bytes",
			r:          bytes.NewBufferString("word1 word2 word3\nline2\nline3 word4"),
			countLines: true,
			countBytes: true,
			wantCount:  3,
			wantByte:   35,
			wantErr:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotCount, gotByte, err := count(test.r, test.countLines, test.countBytes)
			if !errors.Is(err, test.wantErr) {
				t.Errorf("count() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotCount != test.wantCount {
				t.Errorf("count() gotCount = %v, wantCount %v", gotCount, test.wantCount)
			}
			if gotByte != test.wantByte {
				t.Errorf("count() gotByte = %v, wantByte %v", gotByte, test.wantByte)
			}
		})
	}
}

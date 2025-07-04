package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func valid(path string, ext string, minSize int64, fileInfo fs.FileInfo) bool {
	// The file can't be a directory
	// The size of file must be larger than min size
	// If ext is not set, assume the extension matches ext
	// The extension of file must match ext
	if !fileInfo.IsDir() &&
		fileInfo.Size() >= minSize &&
		(ext == "" || ext == filepath.Ext(path)) {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	if err != nil {
		return err
	}

	return nil
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}

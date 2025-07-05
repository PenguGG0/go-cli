package main

import (
	"compress/gzip"
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

func archiveFile(destDir, root, sourcePath string) error {
	// Check if destDir is a directory
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// Get relative directory path of the file
	relDir, err := filepath.Rel(root, filepath.Dir(sourcePath))
	if err != nil {
		return err
	}

	// Get file name of destination compressed file
	targetPath := filepath.Join(destDir, relDir, filepath.Base(sourcePath)+".gz")

	// Open the destination file, create it if it doesn't exist
	if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}
	dest, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err = dest.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Open the source file
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer func() {
		if err = source.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Compress the source to the dest using gzip.Writer
	zw := gzip.NewWriter(dest)
	zw.Name = filepath.Base(sourcePath)
	if _, err = io.Copy(zw, source); err != nil {
		return err
	}
	defer func() {
		if err = zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	return nil
}

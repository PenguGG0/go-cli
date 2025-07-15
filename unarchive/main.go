package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func unarchive(destDir, archive, sourcePath string) error {
	// Check if destDir is a directory
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// Get relative directory path of the file
	relDir, err := filepath.Rel(archive, filepath.Dir(sourcePath))
	if err != nil {
		return err
	}

	// Get file name of destination compressed file
	// If file doesn't end with ".gz", skip it
	targetName, found := strings.CutSuffix(filepath.Base(sourcePath), ".gz")
	if found == false {
		return nil
	}
	targetPath := filepath.Join(destDir, relDir, targetName)

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

	// Extract the source file to target file
	zr, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	zr.Name = filepath.Base(targetPath)
	if _, err = io.Copy(dest, zr); err != nil {
		return err
	}

	defer func() {
		if err = zr.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	log.Println(targetPath)

	return nil
}

func main() {
	log.SetPrefix("RECOVER FILE: ")

	root := flag.String("root", ".", "Root directory to start")
	archive := flag.String("archive", "", "Archive directory")
	flag.Parse()

	if *archive == "" {
		log.Fatalln("Archive directory is empty")
	}

	err := filepath.WalkDir(*archive, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err = unarchive(*root, *archive, path); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalln(err)
	}
}

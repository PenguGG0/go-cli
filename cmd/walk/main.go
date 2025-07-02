package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type application struct {
	ext     string // extension to filter out
	minSize int64  // min file size
	list    bool   // list files or not
}

func run(root string, out io.Writer, app application) error {
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If the file is not valid, skip the rest of the function
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}
		if !valid(path, app.ext, app.minSize, fileInfo) {
			return nil
		}

		if !app.list {
			fmt.Println("--- list is false ---")
		}

		return listFile(path, out)
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	app := application{
		ext:     *ext,
		minSize: *size,
		list:    *list,
	}

	if err := run(*root, os.Stdout, app); err != nil {
		log.Fatalln(err)
	}
}

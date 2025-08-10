package main

import (
	"flag"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type application struct {
	wLog       io.Writer
	archive    string
	extensions []string
	minSize    int64
	list       bool
	del        bool
}

func run(root string, out io.Writer, app application) error {
	delLogger := log.New(app.wLog, "DELETED FILE: ", log.LstdFlags)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If the file is not valid, skip the rest of the function
		fileInfo, err := d.Info()
		if err != nil {
			return err
		}

		if !valid(path, app.extensions, app.minSize, fileInfo) {
			return nil
		}

		// If list was explicitly set, just list file and skip the actions after
		if app.list {
			return listFile(path, out)
		}

		// Archive files and continue if successful
		if app.archive != "" {
			if err = archiveFile(app.archive, root, path); err != nil {
				return err
			}
		}

		if app.del {
			return delFile(path, delLogger)
		}

		// List the file by default
		return listFile(path, out)
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	root := flag.String("root", ".", "Root directory to start")
	logFileName := flag.String("log", "", "Log deletes to this file")
	list := flag.Bool("list", false, "List files only")
	archive := flag.String("archive", "", "Archive directory")
	del := flag.Bool("del", false, "Delete files")
	ext := flag.String("ext", "", "File extension to filter out\n"+
		"This flag allows multiple values separated by ','\n"+
		"e.g., '-ext .txt,.exe'")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	extensions := strings.Split(*ext, ",")

	var (
		logFile *os.File
		err     error
	)

	if *logFileName != "" {
		logFile, err = os.OpenFile(*logFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		logFile = os.Stdout
	}

	defer func() {
		if logFile != nil && logFile != os.Stdout {
			if err = logFile.Close(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	app := application{
		extensions: extensions,
		minSize:    *size,
		list:       *list,
		del:        *del,
		wLog:       logFile,
		archive:    *archive,
	}

	if err = run(*root, os.Stdout, app); err != nil {
		log.Fatalln(err)
	}
}

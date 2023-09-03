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

type config struct {
	ext      string
	size     int64
	list     bool
	del      bool
	delWrite io.Writer
}

func main() {
	root := flag.String("root", ".", "Root directory to start")

	list := flag.Bool("list", false, "List files only")

	del := flag.Bool("del", false, "Delete files")

	ext := flag.String("ext", "", "File extension to filter out.")

	size := flag.Int64("size", 0, "Minimum file size")

	logFileName := flag.String("log", "", "Log deletes to this file")

	flag.Parse()

	var (
		logFile = os.Stdout
		err     error
	)

	if *logFileName != "" {
		logFile, err = os.OpenFile(*logFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer logFile.Close()
	}

	conf := config{
		ext:      *ext,
		list:     *list,
		size:     *size,
		del:      *del,
		delWrite: logFile,
	}

	if err := run(*root, os.Stdout, conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filterOut(path, cfg, info) {
			return nil
		}

		if cfg.list {
			return listFile(path, out)
		}

		if cfg.del {
			logger := log.New(cfg.delWrite, "Deleted file:", log.LstdFlags)
			return delFile(path, logger)
		}

		return listFile(path, out)
	})
}

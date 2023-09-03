package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path string, conf config, fileInfo fs.FileInfo) bool {
	if fileInfo.IsDir() || fileInfo.Size() < conf.size {
		return true
	}

	if conf.ext != "" && filepath.Ext(path) != conf.ext {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, logger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	logger.Println(path)
	return nil
}

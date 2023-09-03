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

func archiveFile(destDir string, root string, path string) error {
	fileInfo, err := os.Stat(destDir)

	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetFile := filepath.Join(destDir, relDir, fileName)

	if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetFile, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	gzipWriter := gzip.NewWriter(out)
	gzipWriter.Name = filepath.Base(path)

	if _, err := io.Copy(gzipWriter, in); err != nil {
		return err
	}

	if err := gzipWriter.Close(); err != nil {
		return err
	}

	return out.Close()
}

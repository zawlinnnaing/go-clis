package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	root := flag.String("root", "", "Root directory for archived files")

	target := flag.String("target", "", "Target directory to place unarchive files")

	flag.Parse()

	err := run(*root, *target, os.Stdout)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(rootDir string, targetDir string, out io.Writer) error {
	return filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		err = unarchiveFile(rootDir, targetDir, path)
		if err != nil {
			return err
		}

		fmt.Fprint(out, path)
		return nil
	})
}

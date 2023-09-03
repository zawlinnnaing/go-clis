package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func unarchiveFile(rootDir string, targetDir string, sourceZipFilePath string) error {
	targetInfo, err := os.Stat(targetDir)
	if err != nil {
		return err
	}

	if !targetInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", targetDir)
	}

	relDir, err := filepath.Rel(rootDir, filepath.Dir(sourceZipFilePath))
	if err != nil {
		return err
	}

	targetFilePath := filepath.Join(targetDir, relDir, strings.ReplaceAll(filepath.Base(sourceZipFilePath), ".gzip", ""))

	err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
	if err != nil {
		return err
	}

	targetFile, err := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	sourceZipFile, err := os.Open(sourceZipFilePath)
	if err != nil {
		return err
	}
	defer sourceZipFile.Close()

	gzipReader, err := gzip.NewReader(sourceZipFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	content, err := io.ReadAll(gzipReader)
	if err != nil {
		return err
	}

	_, err = targetFile.Write(content)
	if err != nil {
		return err
	}

	return targetFile.Close()
}

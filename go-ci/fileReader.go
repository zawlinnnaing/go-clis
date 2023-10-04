package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var stepTypes = []string{"step", "execution", "timeout"}

type YamlStep struct {
	Name    string
	Project string
	Step    string `yaml:"step"`
	Exe     string
	Args    string
}

func readFile(filePath string) ([]executer, error) {
	if filePath == "" {
		return nil, ErrEmptyDir
	}
	ext := filepath.Ext(filePath)
	if ext != ".yaml" && ext != ".yml" {
		return nil, ErrInvalidFile
	}

	yamlFile, err := os.ReadFile(filePath)
	fmt.Println("🚀 ~ file: fileReader.go:30 ~ funcreadFile ~ yamlFile:", yamlFile)
	if err != nil {
		return nil, err
	}
	var yamlSteps []YamlStep
	err = yaml.Unmarshal(yamlFile, yamlSteps)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

var stepTypes = []string{"step", "execution", "timeout"}

type YamlStep struct {
	Name    string `yaml:"name"`
	Project string `yaml:"project"`
	Step    string `yaml:"type"`
	Exe     string `yaml:"exe"`
	Args    string `yaml:"args"`
	Message string `yaml:"message"`
}

func readFile(filePath string) (*[]YamlStep, error) {
	if filePath == "" {
		return nil, ErrEmptyDir
	}
	ext := filepath.Ext(filePath)
	if ext != ".yaml" && ext != ".yml" {
		return nil, ErrInvalidFile
	}

	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var yamlSteps []YamlStep = []YamlStep{}
	err = yaml.Unmarshal(yamlFile, &yamlSteps)
	if err != nil {
		return nil, err
	}

	return &yamlSteps, nil
}

func validateSteps(yamlSteps *[]YamlStep) error {
	errCh := make(chan error, len(*yamlSteps))
	var wg sync.WaitGroup
	for idx, step := range *yamlSteps {
		wg.Add(1)
		go func(step YamlStep, idx int) {
			defer wg.Done()
			if !slices.Contains[[]string](stepTypes, step.Step) {
				errCh <- fmt.Errorf("%w: Invalid step type at %d with step %s. Supported steps: %v", ErrInvalidStep, idx, step.Step, stepTypes)
				return
			}
			_, err := os.Stat(step.Project)
			if errors.Is(err, os.ErrNotExist) {
				errCh <- fmt.Errorf("%w: Project not exists", ErrInvalidStep)
				return
			}
			if err != nil {
				errCh <- err
			}
		}(step, idx)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}

func parseFile(filePath string) ([]executer, error) {
	yamlSteps, err := readFile(filePath)
	if err != nil {
		return nil, err
	}
	err = validateSteps(yamlSteps)
	if err != nil {
		return nil, err
	}
	executers := []executer{}
	for _, step := range *yamlSteps {
		var executor executer
		args := strings.Split(step.Args, " ")
		switch step.Step {
		case "step":
			executor = NewStep(step.Name, step.Exe, step.Message, step.Project, args)
		case "timeout":
			executor = NewTimeoutStep(step.Name, step.Exe, step.Message, step.Project, args, 30*time.Second)
		case "execution":
			executor = NewExecutionStep(step.Name, step.Exe, step.Message, step.Project, args)
		default:
			return nil, errors.New("invalid step type. Should not reached here")
		}
		executers = append(executers, executor)
	}

	return executers, nil
}

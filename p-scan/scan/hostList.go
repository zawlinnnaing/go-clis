package scan

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

var (
	ErrHostExists    = errors.New("host already exists in the list")
	ErrHostNotExists = errors.New("host does not exist in the list")
)

type HostsList struct {
	Hosts []string
}

func (list *HostsList) search(host string) (bool, int) {
	sort.Strings(list.Hosts)
	i := sort.SearchStrings(list.Hosts, host)
	if i < len(list.Hosts) && list.Hosts[i] == host {
		return true, i
	}
	return false, -1
}

func (list *HostsList) Add(host string) error {
	if found, _ := list.search(host); found {
		return fmt.Errorf("%w: %s", ErrHostExists, host)
	}
	list.Hosts = append(list.Hosts, host)
	return nil
}

func (list *HostsList) Remove(host string) error {
	found, index := list.search(host)
	if !found {
		return fmt.Errorf("%w: %s", ErrHostNotExists, host)
	}
	list.Hosts = append(list.Hosts[:index], list.Hosts[index+1:]...)
	return nil
}

func (list *HostsList) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		list.Hosts = append(list.Hosts, scanner.Text())
	}

	return nil
}

func (list *HostsList) Save(file string) error {
	output := ""
	for _, host := range list.Hosts {
		output += fmt.Sprintln(host)
	}

	return os.WriteFile(file, []byte(output), 0644)
}

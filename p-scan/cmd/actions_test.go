package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/zawlinnnaing/go-clis/p-scan/scan"
)

func setUp(t *testing.T, hosts []string, initializeList bool) (string, func()) {
	tempFile, err := os.CreateTemp("", "p-scan_")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	if initializeList {
		hostsList := scan.HostsList{}
		for _, host := range hosts {
			if err := hostsList.Add(host); err != nil {
				t.Fatal(err)
			}
		}
		if err := hostsList.Save(tempFile.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tempFile.Name(), func() {
		os.Remove(tempFile.Name())
	}
}

func TestActions(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}
	testCases := []struct {
		name           string
		hosts          []string
		initializeList bool
		action         func(io.Writer, string, []string) error
		expectedOut    string
	}{
		{
			name:           "AddAction",
			hosts:          hosts,
			initializeList: false,
			action:         addAction,
			expectedOut:    "Added host: host1\nAdded host: host2\nAdded host: host3\n",
		},
		{
			name:           "ListAction",
			hosts:          hosts,
			initializeList: true,
			expectedOut:    "host1\nhost2\nhost3\n",
			action:         listAction,
		},
		{
			name:           "DeleteAction",
			hosts:          hosts[:2],
			expectedOut:    "Removed host: host1\nRemoved host: host2\n",
			action:         deleteAction,
			initializeList: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tempFile, cleanUp := setUp(t, hosts, testCase.initializeList)
			defer cleanUp()
			var out bytes.Buffer
			if err := testCase.action(&out, tempFile, testCase.hosts); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if out.String() != testCase.expectedOut {
				t.Errorf("Expected: %s, received %s", testCase.expectedOut, out.String())
			}
		})
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}
	tempFile, cleanup := setUp(t, hosts, false)
	defer cleanup()
	deleteHost := "host2"
	hostsEnd := []string{
		"host1",
		"host3",
	}
	var out bytes.Buffer
	expectedOut := ""
	for _, host := range hosts {
		expectedOut += fmt.Sprintf("Added host: %s\n", host)
	}
	expectedOut += strings.Join(hosts, "\n")
	expectedOut += fmt.Sprintln()

	expectedOut += fmt.Sprintf("Removed host: %s\n", deleteHost)

	expectedOut += strings.Join(hostsEnd, "\n")
	expectedOut += fmt.Sprintln()

	if err := addAction(&out, tempFile, hosts); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if err := listAction(&out, tempFile, hosts); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if err := deleteAction(&out, tempFile, []string{deleteHost}); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if err := listAction(&out, tempFile, hostsEnd); err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if out.String() != expectedOut {
		t.Errorf("Expected: %s, received: %s", expectedOut, out.String())
	}
}

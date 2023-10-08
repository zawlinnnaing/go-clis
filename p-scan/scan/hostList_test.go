package scan_test

import (
	"errors"
	"os"
	"testing"

	"github.com/zawlinnnaing/go-clis/p-scan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{
			name:      "AddNew",
			host:      "host2",
			expectLen: 2,
			expectErr: nil,
		},
		{
			name:      "AddExists",
			host:      "host1",
			expectLen: 1,
			expectErr: scan.ErrHostExists,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			hostsList := &scan.HostsList{}
			if err := hostsList.Add("host1"); err != nil {
				t.Fatal(err)
			}
			err := hostsList.Add(testCase.host)
			if testCase.expectErr != nil {
				if err == nil {
					t.Errorf("Expect error: %v, received none", testCase.expectErr)
					return
				}
				if !errors.Is(err, testCase.expectErr) {
					t.Errorf("Expected error: %v, received error %v", testCase.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(hostsList.Hosts) != testCase.expectLen {
				t.Errorf("Expected length: %d, received length: %d", testCase.expectLen, len(hostsList.Hosts))
				return
			}
			if hostsList.Hosts[1] != testCase.host {
				t.Errorf("Expect index 1 to be %s, received %s", testCase.host, hostsList.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{
			name:      "RemoveExists",
			host:      "host1",
			expectLen: 1,
			expectErr: nil,
		},
		{
			name:      "RemoveNotExists",
			host:      "host3",
			expectLen: 2,
			expectErr: scan.ErrHostNotExists,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			hostsList := scan.HostsList{
				Hosts: []string{"host1", "host2"},
			}
			err := hostsList.Remove(testCase.host)
			if testCase.expectErr != nil {
				if err == nil {
					t.Errorf("Expected error: %v, received none", testCase.expectErr)
					return
				}
				if !errors.Is(err, testCase.expectErr) {
					t.Errorf("Expected error: %v, received error %v", testCase.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(hostsList.Hosts) != testCase.expectLen {
				t.Errorf("Expected length: %d, received length: %d", testCase.expectLen, len(hostsList.Hosts))
				return
			}
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hostsList1 := &scan.HostsList{}
	hostsList2 := &scan.HostsList{}

	if err := hostsList1.Add("host1"); err != nil {
		t.Fatalf("Error adding value: %s", err)
	}
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())
	if err := hostsList1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving file: %s", err)
	}
	if err := hostsList2.Load(tf.Name()); err != nil {
		t.Fatalf("Error loading file: %s", err)
	}
	if hostsList1.Hosts[0] != hostsList2.Hosts[0] {
		t.Errorf("Host %q should match %q", hostsList1.Hosts[0], hostsList2.Hosts[0])
	}
}

func TestLoadFileNotExists(t *testing.T) {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	if err := os.Remove(tempFile.Name()); err != nil {
		t.Fatalf("Error removing temp file: %s", err)
	}
	hostsList := scan.HostsList{}
	if err := hostsList.Load(tempFile.Name()); err != nil {
		t.Errorf("Expected no error, received error: %s", err)
	}
}

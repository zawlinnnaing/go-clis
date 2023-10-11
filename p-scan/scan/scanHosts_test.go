package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/zawlinnnaing/go-clis/p-scan/scan"
)

func TestStateString(t *testing.T) {
	portState := scan.PortState{}

	if portState.Open.String() != "closed" {
		t.Errorf("Expected %v, received %v", "closed", portState.Open.String())
	}

	portState.Open = true

	if portState.Open.String() != "open" {
		t.Errorf("Expected %v, received: %v", "open", portState.Open.String())
	}
}

func TestHostFound(t *testing.T) {
	host := "localhost"
	hostsList := scan.HostsList{}
	hostsList.Add(host)

	testCases := []struct {
		name        string
		expectState string
	}{
		{
			name:        "OpenPort",
			expectState: "open",
		},
		{
			name:        "ClosePort",
			expectState: "closed",
		},
	}

	availablePorts := []int{}
	for _, testCase := range testCases {
		conn, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		_, portStr, err := net.SplitHostPort(conn.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		availablePorts = append(availablePorts, port)
		if testCase.expectState == "closed" {
			conn.Close()
		}
	}
	res := scan.Run(&hostsList, availablePorts)
	if len(res) != 1 {
		t.Errorf("Expected length 1, received %d", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Expected host: %s, received host: %s", host, res[0].Host)
	}

	if res[0].NotFound {
		t.Error("Expected host to be found")
	}

	if len(res[0].PortStates) != 2 {
		t.Errorf("Expected 2 port states, received %d", len(res[0].PortStates))
	}
	for i, testCase := range testCases {
		result := &res[0]
		if result.PortStates[i].Port != availablePorts[i] {
			t.Errorf("Expected port to be %d, received %d", availablePorts[i], result.PortStates[i].Port)
		}
		if result.PortStates[i].Open.String() != testCase.expectState {
			t.Errorf("Expected port state to be %s, received %s", testCase.expectState, result.PortStates[i].Open.String())
		}
	}
}

func TestHostNotFound(t *testing.T) {
	host := "389.389.389.389"
	hostsList := scan.HostsList{}
	hostsList.Add(host)

	results := scan.Run(&hostsList, []int{})

	if len(results) != 1 {
		t.Fatalf("Expected length to be: %d, received %d", 1, len(results))
	}

	if results[0].Host != host {
		t.Errorf("Expected host to be %s, received %s", host, results[0].Host)
	}

	if !results[0].NotFound {
		t.Errorf("Expected not found, received %v", results[0].NotFound)
	}

	if len(results[0].PortStates) != 0 {
		t.Errorf("Expected port states length to be 0, received %d", len(results[0].PortStates))
	}
}

package scan_test

import (
	"fmt"
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
		{
			name:        "UdpPort",
			expectState: "open",
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
	res := scan.Run(&hostsList, availablePorts, "tcp")
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

	results := scan.Run(&hostsList, []int{}, "tcp")

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

func TestUDP(t *testing.T) {
	testCases := []struct {
		name string
		open bool
	}{
		{
			name: "OpenPort",
			open: true,
		},
	}
	host := "localhost"
	hostsList := scan.HostsList{}
	hostsList.Add(host)
	availablePorts := []int{}
	for _, testCase := range testCases {

		conn, err := net.ListenPacket("udp", fmt.Sprintf("%s:0", host))
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()
		_, portStr, err := net.SplitHostPort(conn.LocalAddr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}
		availablePorts = append(availablePorts, port)
		if !testCase.open {

			conn.Close()
		}
	}
	res := scan.Run(&hostsList, availablePorts, "udp")
	if len(res) != 1 {
		t.Fatal("Expected length to be 1")
	}
	result := res[0]
	if len(result.PortStates) != len(availablePorts) {
		t.Errorf("Expected ports to be: %d, received %d", len(availablePorts), len(result.PortStates))
	}
	for i, testCase := range testCases {
		if result.PortStates[i].Port != availablePorts[i] {
			t.Errorf("Expected port: %v, received: %v", availablePorts[i], result.PortStates[i].Port)
		}
		if bool(result.PortStates[i].Open) != testCase.open {
			t.Errorf("Expected port open: %v, received %v", testCase.open, result.PortStates[i].Open)
		}
	}
}

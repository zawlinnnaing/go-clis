package scan

import (
	"fmt"
	"net"
	"time"
)

type state bool

type PortState struct {
	Port int
	Open state
}

func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

type Result struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

func scanPort(host string, port int, networkType string) PortState {
	portState := PortState{
		Port: port,
	}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	scanConnect, err := net.DialTimeout(networkType, address, 1*time.Second)
	if err != nil {
		return portState
	}
	scanConnect.Close()
	portState.Open = true
	return portState
}

func Run(hostsList *HostsList, ports []int, networkType string) []Result {
	results := []Result{}
	for _, host := range hostsList.Hosts {
		result := Result{
			Host: host,
		}
		if _, err := net.LookupHost(host); err != nil {
			result.NotFound = true
			results = append(results, result)
			continue
		}
		for _, port := range ports {
			portState := scanPort(host, port, networkType)
			result.PortStates = append(result.PortStates, portState)
		}
		results = append(results, result)
	}

	return results
}

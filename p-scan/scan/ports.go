package scan

import (
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidPortsFormat = errors.New("invalid ports format")
)

func parsePortsRange(portRangeStr string) (int, int, error) {
	portRanges := strings.Split(portRangeStr, "-")
	if len(portRanges) != 2 {
		return 0, 0, errors.New("invalid port range")
	}
	startPort, err := strconv.Atoi(portRanges[0])
	if err != nil {
		return 0, 0, err
	}
	endPort, err := strconv.Atoi(portRanges[1])
	if err != nil {
		return 0, 0, err
	}
	return startPort, endPort, nil
}

func getPorts(startPort int, endPort int) []int {
	portsLen := (endPort - startPort) + 1
	ranges := make([]int, portsLen)
	for i := 0; i < portsLen; i++ {
		ranges[i] = startPort + i
	}
	return ranges
}

func ParsePorts(ports string) ([]int, error) {
	var portsArr []int
	if strings.ContainsAny(ports, ",") {
		portsStr := strings.Split(ports, ",")
		for _, portStr := range portsStr {
			port, err := strconv.Atoi(portStr)
			if err != nil {
				return nil, err
			}
			portsArr = append(portsArr, port)
		}
		return portsArr, nil
	}
	if strings.ContainsAny(ports, "-") {
		startPort, endPort, err := parsePortsRange(ports)
		if err != nil {
			return nil, err
		}
		portsArr = getPorts(startPort, endPort)
		return portsArr, nil
	}
	return nil, ErrInvalidPortsFormat
}

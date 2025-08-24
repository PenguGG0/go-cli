package scan

import (
	"net"
	"strconv"
	"time"
)

type state bool

func (s state) String() string {
	if s {
		return "open"
	}

	return "closed"
}

type PortState struct {
	Port int
	Open state
}

// Results represent the scan results for a single host
type Results struct {
	Host       string
	PortStates []PortState
}

// scanPort scans a single TCP port
func scanPort(host string, port int, timeout int) PortState {
	// Assume the port is close
	p := PortState{
		Port: port,
		Open: false,
	}

	// Join the host and port to a address, and try to connect the address
	address := net.JoinHostPort(host, strconv.Itoa(port))
	scanConn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	// When the connection succeeds, set the value od p.Open to true
	if err == nil {
		p.Open = true
		scanConn.Close()
	}

	return p
}

func Run(hl *HostsList, ports []int, timeout int) []Results {
	res := make([]Results, 0, len(hl.Hosts))

	// Scan every host in the list
	for _, h := range hl.Hosts {
		r := Results{
			Host:       h,
			PortStates: []PortState{},
		}

		// Valid the address, skip scanning ports if it's unvalid
		if _, err := net.LookupHost(h); err != nil {
			res = append(res, r)
			continue
		}

		// Scan ports of the host
		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p, timeout))
		}

		res = append(res, r)
	}

	return res
}

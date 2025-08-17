package scan_test

import (
	"net"
	"strconv"
	"testing"

	"github.com/PenguGG0/go-cli/pScan/scan"
)

func TestStateString(t *testing.T) {
	ps := scan.PortState{}
	if ps.Open.String() != "closed" {
		t.Errorf("Got %q, expected %q\n", ps.Open.String(), "closed")
	}

	ps.Open = true
	if ps.Open.String() != "open" {
		t.Errorf("Got %q, expected %q\n", ps.Open.String(), "open")
	}
}

func TestRunHostFound(t *testing.T) {
	testCases := []struct {
		name        string
		expectState string
	}{
		{"OpenPort", "open"},
		{"ClosedPort", "closed"},
	}

	host := "localhost"
	hl := &scan.HostsList{}
	hl.Add(host)

	// Initialize ports values for the testcases
	ports := []int{}
	for _, tc := range testCases {
		// Use port 0 means to get a temporary port assigned dynamically
		ln, err := net.Listen("tcp", net.JoinHostPort(host, "0"))
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()

		// Get the port value
		_, portStr, err := net.SplitHostPort(ln.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			t.Fatal(err)
		}

		// Add the port value to our ports slice
		// Close it if testcase wants a closed port
		ports = append(ports, port)
		if tc.name == "ClosedPort" {
			ln.Close()
		}
	}

	res := scan.Run(hl, ports)

	if len(res) != 1 {
		t.Fatalf("Got %q, expected 1\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Got %q, expected %q\n", res[0].Host, host)
	}
	if len(res[0].PortStates) == 0 {
		t.Errorf("Got no host, expected host: %q\n", host)
	}
	if len(res[0].PortStates) != 2 {
		t.Errorf("Got %q host, expected 2 host", len(res[0].PortStates))
	}

	for i, tc := range testCases {
		if res[0].PortStates[i].Port != ports[i] {
			t.Errorf("Got port[%q] %q, expected port[%q] %q\n", i, res[0].PortStates[i].Port, i, ports[i])
		}
		if res[0].PortStates[i].Open.String() != tc.expectState {
			t.Errorf("Got port %q state %q, expected %q\n", ports[i], res[0].PortStates[i].Open.String(), tc.expectState)
		}
	}
}

func TestHostNotFound(t *testing.T) {
	// An unvalided host
	host := "389.389.389.389"
	hl := &scan.HostsList{}
	hl.Add(host)

	res := scan.Run(hl, []int{})

	if len(res) != 1 {
		t.Fatalf("Got %q, expected 1\n", len(res))
	}
	if res[0].Host != host {
		t.Errorf("Got %q, expected %q\n", res[0].Host, host)
	}
	if len(res[0].PortStates) != 0 {
		t.Errorf("Got %q port state, expected 0 port state", len(res[0].PortStates))
	}
}

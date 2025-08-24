package cmd

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"testing"

	"github.com/PenguGG0/go-cli/pScan/scan"
)

var (
	timeout  int  = 1
	showOpen bool = false
)

func setup(t *testing.T, hosts []string, initList bool) (string, func()) {
	// Create temp file
	tf, err := os.CreateTemp("", "pScan")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()

	// Initialize lis if needed
	if initList {
		hl := &scan.HostsList{}

		for _, h := range hosts {
			hl.Add(h)
		}

		if err := hl.Save(tf.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tf.Name(), func() {
		os.Remove(tf.Name())
	}
}

func TestHostActions(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	testCases := []struct {
		name       string
		args       []string
		expectOut  string
		initList   bool
		actionFunc func(io.Writer, string, []string) error
	}{
		{
			name:       "AddAction",
			args:       hosts,
			expectOut:  "Added host: host1\nAdded host: host2\nAdded host: host3\n",
			initList:   false,
			actionFunc: addAction,
		},
		{
			name:       "ListAction",
			expectOut:  "host1\nhost2\nhost3\n",
			initList:   true,
			actionFunc: listAction,
		},
		{
			name:       "DeleteAction",
			args:       []string{"host1", "host2"},
			expectOut:  "Deleted host: host1\nDeleted host: host2\n",
			initList:   true,
			actionFunc: deleteAction,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tf, cleanup := setup(t, hosts, tc.initList)
			defer cleanup()

			var out bytes.Buffer

			if err := tc.actionFunc(&out, tf, tc.args); err != nil {
				t.Fatalf("Got %q, expected no error.\n", err)
			}
			if out.String() != tc.expectOut {
				t.Errorf("Got %q, expected %q\n", out.String(), tc.expectOut)
			}
		})
	}
}

func TestScanAction(t *testing.T) {
	hosts := []string{
		"localhost",
		"unknownhostoutthere",
	}

	tf, cleanup := setup(t, hosts, true)
	defer cleanup()

	// Init 2 ports, 1 open, 1 closed
	ports := []int{}
	for i := 0; i < 2; i++ {
		// Use port 0 means to get a temporary port assigned dynamically
		ln, err := net.Listen("tcp", net.JoinHostPort("localhost", "0"))
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
		if i == 1 {
			ln.Close()
		}
	}

	expectedOut := fmt.Sprintln("localhost:")
	expectedOut += fmt.Sprintf("\t%d: open\n", ports[0])
	expectedOut += fmt.Sprintf("\t%d: closed\n", ports[1])
	expectedOut += fmt.Sprintln()
	expectedOut += fmt.Sprintln("unknownhostoutthere: Host not found")
	expectedOut += fmt.Sprintln()

	var out bytes.Buffer
	if err := scanAction(&out, tf, ports, showOpen, timeout); err != nil {
		t.Fatalf("Expected no error, got %q\n", err)
	}
	if out.String() != expectedOut {
		t.Errorf("Expected output %q, got %q\n", expectedOut, out.String())
	}
}

func TestIntegration(t *testing.T) {
	hosts := []string{
		"host1",
		"host2",
		"host3",
	}

	tf, cleanup := setup(t, hosts, false)
	defer cleanup()

	delHost := "host2"

	hostsEnd := []string{"host1", "host3"}

	var out bytes.Buffer

	expectOut := ""

	// Add hosts to the list
	if err := addAction(&out, tf, hosts); err != nil {
		t.Fatalf("Got %q, expected no error.\n", err)
	}
	for _, h := range hosts {
		expectOut += fmt.Sprintf("Added host: %s\n", h)
	}

	// List hosts
	if err := listAction(&out, tf, hosts); err != nil {
		t.Fatalf("Got %q, expected no error.\n", err)
	}
	for _, h := range hosts {
		expectOut += fmt.Sprintf("%s\n", h)
	}

	// Delete host2
	if err := deleteAction(&out, tf, []string{delHost}); err != nil {
		t.Fatalf("Got %q, expected no error.\n", err)
	}
	expectOut += fmt.Sprintf("Deleted host: %s\n", delHost)

	// List hosts after delete
	if err := listAction(&out, tf, hosts); err != nil {
		t.Fatalf("Got %q, expected no error.\n", err)
	}
	for _, h := range hostsEnd {
		expectOut += fmt.Sprintf("%s\n", h)
	}

	// Scan hosts
	if err := scanAction(&out, tf, nil, showOpen, timeout); err != nil {
		t.Fatalf("Got %q, expected no error.\n", err)
	}
	for _, v := range hostsEnd {
		expectOut += fmt.Sprintf("%s: Host not found\n\n", v)
	}

	if out.String() != expectOut {
		t.Errorf("Got %q, expected %q\n", out.String(), expectOut)
	}
}

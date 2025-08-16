package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/PenguGG0/go-cli/pScan/scan"
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

	if out.String() != expectOut {
		t.Errorf("Got %q, expected %q\n", out.String(), expectOut)
	}
}

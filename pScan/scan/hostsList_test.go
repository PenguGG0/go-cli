package scan_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/PenguGG0/go-cli/pScan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 1, scan.ErrExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}

			err := hl.Add(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatal("Got nil, expected error")
				}
				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Got error %q, expected error %q\n", err, tc.expectErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Got %q, expected no error", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Got %d list length, expected %d list length\n", len(hl.Hosts), tc.expectLen)
			}

			if hl.Hosts[1] != tc.host {
				t.Errorf("Got host name %q, expected %q\n", hl.Hosts[1], tc.host)
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
		{"RemoveExisting", "host1", 1, nil},
		{"RemoveNotFound", "host3", 1, scan.ErrNotExists},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}

			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}

			err := hl.Remove(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatal("Got nil, expected error")
				}
				if !errors.Is(err, tc.expectErr) {
					t.Errorf("Got error %q, expected error %q\n", err, tc.expectErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Got %q, expected no error", err)
			}

			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("Got %d list length, expected %d list length\n", len(hl.Hosts), tc.expectLen)
			}

			// if hl.Hosts[0] != tc.host {
			// 	t.Errorf("Got host name %q, expected %q\n", hl.Hosts[1], tc.host)
			// }
		})
	}
}

func TestSaveLoad(t *testing.T) {
	hl1 := scan.HostsList{}
	hl2 := scan.HostsList{}
	hostName := "host1"
	hl1.Add(hostName)
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())
	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}
	if err := hl2.Load(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}
	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("Host %q should match %q host.", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tempFile, err := os.CreateTemp(os.TempDir(), "pScan_*")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	filePath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	
	if err = os.Remove(filePath); err != nil {
		t.Fatalf("Error deleting temp file: %s", err)
	}

	hl := &scan.HostsList{}
	if err = hl.Load(filePath); err != nil {
		t.Errorf("Expected no error, got %q instead\n", err)
	}
}

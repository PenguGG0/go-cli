package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func setupGit(t *testing.T, proj string) func() {
	t.Helper()

	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "goci_test")
	if err != nil {
		t.Fatal(err)
	}

	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatal(err)
	}

	remoteURI := fmt.Sprintf("file://%s", tempDir)

	gitCmdList := []struct {
		args []string
		dir  string
		env  []string
	}{
		{[]string{"init", "--bare"}, tempDir, nil},
		{[]string{"init"}, projPath, nil},
		{[]string{"remote", "add", "origin", remoteURI}, projPath, nil},
		{[]string{"add", "."}, projPath, nil},
		{[]string{"commit", "-m", "test"}, projPath, []string{
			"GIT_COMMITTER_NAME=test",
			"GIT_COMMITTER_EMAIL=test@example.com",
			"GIT_AUTHOR_NAME=test",
			"GIT_AUTHOR_EMAIL=test@example.com",
		}},
	}

	for _, g := range gitCmdList {
		gitCmd := exec.Command(gitExec, g.args...)
		gitCmd.Dir = g.dir

		if g.env != nil {
			gitCmd.Env = append(os.Environ(), g.env...)
		}

		if err = gitCmd.Run(); err != nil {
			t.Fatal(err)
		}
	}

	return func() {
		if err = os.RemoveAll(tempDir); err != nil {
			return
		}
		if err = os.RemoveAll(filepath.Join(projPath, ".git")); err != nil {
			return
		}
	}
}

func TestHelperContext(t *testing.T) {
	if os.Getenv("GO_HELPER_Context") != "1" {
		return
	}

	if os.Args[2] == "git" {
		_, err := fmt.Fprintln(os.Stdout, "Everything up-to-date")
		if err != nil {
			return
		}
		os.Exit(0)
	}

	os.Exit(1)
}

func TestHelperTimeout(t *testing.T) {
	if os.Getenv("GO_HELPER_TIMEOUT") != "1" {
		return
	}

	time.Sleep(3 * time.Second)
}

func mockCmdContext(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperContext"}
	cs = append(cs, exe)
	cs = append(cs, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = append(cmd.Env, "GO_HELPER_Context=1")
	return cmd
}

func mockCmdTimeout(ctx context.Context, exe string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperTimeout"}
	cs = append(cs, exe)
	cs = append(cs, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = append(cmd.Env, "GO_HELPER_TIMEOUT=1")
	return cmd
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		proj     string
		outStr   string
		expErr   error
		setupGit bool
		mockCmd  func(ctx context.Context, name string, args ...string) *exec.Cmd
	}{
		{
			name: "allSuccess",
			proj: "./testdata/tool",
			outStr: "Go Build: SUCCESS\n" +
				"Go Test: SUCCESS\n" +
				"Gofmt: SUCCESS\n" +
				"Git Push: SUCCESS\n",
			expErr:   nil,
			setupGit: true,
			mockCmd:  nil,
		},
		{
			name: "allSuccessMock",
			proj: "./testdata/tool",
			outStr: "Go Build: SUCCESS\n" +
				"Go Test: SUCCESS\n" +
				"Gofmt: SUCCESS\n" +
				"Git Push: SUCCESS\n",
			expErr:   nil,
			setupGit: false,
			mockCmd:  mockCmdContext,
		},
		{
			name:     "buildFail",
			proj:     "./testdata/toolBuildErr",
			outStr:   "",
			expErr:   &stepErr{step: "go build"},
			setupGit: false,
			mockCmd:  mockCmdContext,
		},
		{
			name:     "fmtFail",
			proj:     "./testdata/toolFmtErr",
			outStr:   "",
			expErr:   &stepErr{step: "go fmt"},
			setupGit: false,
			mockCmd:  nil,
		},
		{
			name:     "timeoutFail",
			proj:     "./testdata/tool",
			outStr:   "",
			expErr:   context.DeadlineExceeded,
			setupGit: false,
			mockCmd:  mockCmdTimeout,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupGit {
				_, err := exec.LookPath("git")
				if err != nil {
					t.Skip("Git is not installed, skipping test")
				}
				cleanup := setupGit(t, tc.proj)
				defer cleanup()
			}

			if tc.mockCmd != nil {
				command = tc.mockCmd
			}

			var out bytes.Buffer

			err := run(tc.proj, &out)
			if tc.expErr != nil {
				if err == nil {
					t.Errorf("Got error: nil, expected: %q", tc.expErr)
					return
				}
				if !errors.Is(err, tc.expErr) {
					t.Errorf("Got error: %q, expected: %q", err, tc.expErr)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %q", err)
			}

			if out.String() != tc.outStr {
				t.Errorf("Got output: %q, expected: %q", out.String(), tc.outStr)
			}
		})
	}
}

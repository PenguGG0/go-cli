/*
Copyright © 2025 The Pragmatic Programmers, LLC
Copyrights apply to this source code.
Check LICENSE for details.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/PenguGG0/go-cli/pScan/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a port scan on the hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile := viper.GetString("hosts-file")

		portsStr, err := cmd.Flags().GetString("ports")
		if err != nil {
			return err
		}

		ports, err := parsePortsString(portsStr)
		if err != nil {
			return err
		}

		showOpen, err := cmd.Flags().GetBool("show-open")
		if err != nil {
			return err
		}

		timeout, err := cmd.Flags().GetInt("timeout")
		if err != nil {
			return err
		}

		return scanAction(os.Stdout, hostsFile, ports, showOpen, timeout)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// By default, this flag sets the ports to be scanned as 22, 80 and 443
	// Users can change them using "--ports" or "-p"
	scanCmd.Flags().StringP("ports", "p", "22,80,443", "ports to scan (e.g., 22,80,443 or 1-1024)")

	scanCmd.Flags().BoolP("show-open", "s", false, "only show open ports")

	scanCmd.Flags().IntP("timeout", "t", 1, "timeout for the scan (s)")
}

func parsePortsString(portsStr string) ([]int, error) {
	ports := []int{}
	parts := strings.SplitSeq(portsStr, ",")

	for part := range parts {
		if strings.Contains(part, "-") {
			// range ports (e.g., 1-1024)
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid ports range format: %s", part)
			}

			startPort, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid start port in range %s: %w", part, err)
			}

			endPort, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid end port in range %s: %w", part, err)
			}

			if startPort < 1 || endPort > 65535 || startPort > endPort {
				return nil, fmt.Errorf("invalid port range %s: ports must be between 1 and 65535, start <= end", part)
			}

			for p := startPort; p <= endPort; p++ {
				ports = append(ports, p)
			}
		} else {
			// single port (e.g., 8080 or 443)
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port number %s: %w", part, err)
			}

			if p < 1 || p > 65535 {
				return nil, fmt.Errorf("invalid port number %s: port must be between 1 and 65535", part)
			}

			ports = append(ports, p)
		}
	}

	return ports, nil
}

func printResults(out io.Writer, results []scan.Results, showOpen bool) error {
	message := ""
	for _, r := range results {
		message += fmt.Sprintf("%s:", r.Host)

		if len(r.PortStates) == 0 {
			message += " Host not active\n\n"
			continue
		}

		message += "\n"
		for _, p := range r.PortStates {
			// skip the closed ports when showOpen is true
			if showOpen && !bool(p.Open) {
				continue
			}
			message += fmt.Sprintf("\t%d: %s\n", p.Port, p.Open.String())
		}
		message += "\n"
	}

	_, err := fmt.Fprint(out, message)
	return err
}

func scanAction(out io.Writer, hostsFile string, ports []int, showOpen bool, timeout int) error {
	hl := &scan.HostsList{}

	if err := hl.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(hl, ports, timeout)

	return printResults(out, results, showOpen)
}

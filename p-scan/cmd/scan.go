/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zawlinnnaing/go-clis/p-scan/scan"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Run a scan on specified ports on the host",
	RunE: func(cmd *cobra.Command, args []string) error {
		portsFlag, err := cmd.Flags().GetString("ports")
		if err != nil {
			return err
		}
		hostsFile := viper.GetString("hosts-file")
		ports, err := scan.ParsePorts(portsFlag)
		if err != nil {
			return err
		}
		networkFlag, err := cmd.Flags().GetString("network")
		if err != nil {
			return err
		}
		return scanAction(os.Stdout, hostsFile, ports, networkFlag)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringP("ports", "p", "22,80,443", "Ports to scan. Must be a list of ports or ports range. (eg, 22,80,443 or 22-443)")
	scanCmd.Flags().StringP("network", "n", "tcp", "Network type: tcp or udp. (default: tcp)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scanAction(writer io.Writer, hostsFile string, ports []int, networkType string) error {
	hostsList := scan.HostsList{}

	if err := hostsList.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(&hostsList, ports, networkType)

	return printResults(writer, results)
}

func printResults(writer io.Writer, results []scan.Result) error {
	message := ""

	for _, result := range results {
		message += fmt.Sprintf("%s\n", result.Host)
		if result.NotFound {
			message += "Host not found\n"
			continue
		}
		for _, portState := range result.PortStates {
			message += fmt.Sprintf("\t%d: %s\n", portState.Port, portState.Open)
		}
		message += fmt.Sprintln()
	}
	_, err := fmt.Fprint(writer, message)
	return err
}

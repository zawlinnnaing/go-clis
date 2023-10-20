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
		hostsFile := viper.GetString("hosts-file")
		ports, err := cmd.Flags().GetIntSlice("ports")
		if err != nil {
			return err
		}
		return scanAction(os.Stdout, hostsFile, ports)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().IntSliceP("ports", "p", []int{22, 80, 443}, "Ports to scan")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scanAction(writer io.Writer, hostsFile string, ports []int) error {
	hostsList := scan.HostsList{}

	if err := hostsList.Load(hostsFile); err != nil {
		return err
	}

	results := scan.Run(&hostsList, ports)

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

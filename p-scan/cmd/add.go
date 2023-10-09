/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/zawlinnnaing/go-clis/p-scan/scan"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <host1>...<hostN>",
	Short:        "Add new host(s) to list",
	Aliases:      []string{"a"},
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostFilePath, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return addAction(os.Stdout, hostFilePath, args)
	},
}

func addAction(writer io.Writer, hostsFile string, args []string) error {
	hostsList := scan.HostsList{}
	if err := hostsList.Load(hostsFile); err != nil {
		return err
	}
	for _, host := range args {
		if err := hostsList.Add(host); err != nil {
			return err
		}
		fmt.Fprintln(writer, "Added host: ", host)
	}
	return hostsList.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List hosts in hosts list",
	Aliases: []string{"l"},
	RunE: func(cmd *cobra.Command, args []string) error {
		hostFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return listAction(os.Stdout, hostFile, args)
	},
}

func listAction(writer io.Writer, hostFile string, args []string) error {
	hostsList := &scan.HostsList{}

	if err := hostsList.Load(hostFile); err != nil {
		return err
	}

	for _, host := range hostsList.Hosts {
		if _, err := fmt.Fprintln(writer, host); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	hostsCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

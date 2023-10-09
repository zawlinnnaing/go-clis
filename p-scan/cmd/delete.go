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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <host1>...<hostN>",
	Short:        "Delete host(s) from list",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		hostsFile, err := cmd.Flags().GetString("hosts-file")
		if err != nil {
			return err
		}
		return deleteAction(os.Stdout, hostsFile, args)
	},
}

func deleteAction(writer io.Writer, hostsFile string, args []string) error {
	hostsList := scan.HostsList{}
	if err := hostsList.Load(hostsFile); err != nil {
		return err
	}
	for _, host := range args {
		if err := hostsList.Remove(host); err != nil {
			return err
		}
		fmt.Fprintln(writer, "Removed host: ", host)
	}
	return hostsList.Save(hostsFile)
}

func init() {
	hostsCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

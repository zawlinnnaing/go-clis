/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:          "add <task>",
	Short:        "Add a new task to the list",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		return addAction(apiRoot, args, os.Stdout)
	},
}

func addAction(apiRoot string, args []string, writer io.Writer) error {
	task := strings.Join(args, " ")
	if err := addItem(apiRoot, task); err != nil {
		return err
	}
	return printAdd(writer, task)
}

func printAdd(writer io.Writer, task string) error {
	_, err := fmt.Fprintf(writer, "Added task %q to the list.\n", task)
	return err
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

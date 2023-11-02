/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:          "complete <id>",
	Short:        "Marks an item as completed",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		return completeAction(apiRoot, args[0], os.Stdout)
	},
}

func completeAction(apiRoot string, idArg string, writer io.Writer) error {
	id, err := strconv.Atoi(idArg)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotNumber, err)
	}
	err = completeItem(apiRoot, id)
	if err != nil {
		return err
	}

	return printComplete(writer, id)
}

func printComplete(writer io.Writer, id int) error {
	_, err := fmt.Fprintf(writer, "Item number %d marked as completed.\n", id)
	return err
}

func init() {
	rootCmd.AddCommand(completeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

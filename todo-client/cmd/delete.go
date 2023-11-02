/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:          "delete <id>",
	Short:        "Delete an item from list",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		err := deleteAction(apiRoot, args[0])
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(os.Stdout, "Item number %s has been deleted.", args[0])
		return err
	},
}

func deleteAction(apiRoot string, idArg string) error {
	id, err := strconv.Atoi(idArg)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNotNumber, err)
	}
	return deleteItem(apiRoot, id)
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

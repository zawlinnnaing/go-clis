/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:          "view <id>",
	Short:        "View details about a single item",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiRoot := viper.GetString("api-root")
		return viewAction(apiRoot, args[0], os.Stdout)

	},
}

func viewAction(apiRoot string, idArg string, writer io.Writer) error {
	id, err := strconv.Atoi(idArg)
	if err != nil {
		return err
	}
	item, err := getOne(apiRoot, id)
	if err != nil {
		return err
	}
	return printOne(writer, item)
}

func printOne(writer io.Writer, item Item) error {
	tabWriter := tabwriter.NewWriter(writer, 14, 2, 0, ' ', 0)
	fmt.Fprintf(tabWriter, "Task:\t%s\n", item.Task)
	fmt.Fprintf(tabWriter, "CreatedAt:\t%s\n", item.CreatedAt.Format(timeFormat))
	if item.Done {
		fmt.Fprintf(tabWriter, "Completed:\t%s\n", "Yes")
		fmt.Fprintf(tabWriter, "CompletedAt:\t%s\n", item.CompletedAt.Format(timeFormat))
	} else {
		fmt.Fprintf(tabWriter, "Completed:\t%s\n", "No")
	}
	return tabWriter.Flush()
}

func init() {
	rootCmd.AddCommand(viewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// viewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// viewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

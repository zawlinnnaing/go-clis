/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion for your command",
	Long: `To load your completions run source <(pScan completion)
	To load completions automatically on login, add this line to you .bashrc file:
	$ ~/.bashrc source <(pScan completion)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completeAction(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// completionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// completionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func completeAction(writer io.Writer) error {
	return rootCmd.GenZshCompletion(writer)
}

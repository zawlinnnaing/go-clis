/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate documentation for command",
	RunE: func(cmd *cobra.Command, args []string) error {
		directory, err := cmd.Flags().GetString("dir")
		if err != nil {
			return err
		}
		if directory == "" {
			if directory, err = os.MkdirTemp("", "p-scan"); err != nil {
				return err
			}
		}
		return docsAction(directory, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// docsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// docsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	docsCmd.Flags().StringP("dir", "d", "", "Destination directory for generated docs")
}

func docsAction(dir string, writer io.Writer) error {
	if err := doc.GenMarkdownTree(rootCmd, dir); err != nil {
		return err
	}
	_, err := fmt.Fprintf(writer, "Documentation successfully generated at %s", dir)
	return err
}

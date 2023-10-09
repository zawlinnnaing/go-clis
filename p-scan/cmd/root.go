/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "p-scan",
	Short: "Fast TCP port scanner",
	Long: `p-scan -- short for Port Scanner -- executes TCP port scan on a list of hosts.
	p-scan allows you to add, list, and delete hosts from the list.
	p-scan executes a port scan on specified TCP ports.
	You can customize the target ports using a command line flag.
	REMEMBER: Please only execute this program against the system you own.
	`,
	Version: "0.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("hosts-file", "f", "p-scan.hosts", "p-scan hosts")
	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	rootCmd.SetVersionTemplate(versionTemplate)
}

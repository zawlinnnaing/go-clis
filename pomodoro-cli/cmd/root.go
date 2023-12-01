/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/app"
	"github.com/zawlinnnaing/go-clis/pomodoro-cli/pomodoro"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pomodoro-cli",
	Short: "Interactive Pomodoro Timer",
	Long:  `An application for focusing on task with short and long breaks`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := getRepo()
		if err != nil {
			return err
		}
		config := pomodoro.NewIntervalConfig(
			repo,
			viper.GetDuration("pomo"),
			viper.GetDuration("short"),
			viper.GetDuration("long"),
			viper.GetDuration("until"),
		)
		return rootAction(os.Stdout, config)
	},
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pomodoro-cli.yaml)")
	rootCmd.Flags().StringP("db", "d", "./pomodoro.db", "Database file")

	rootCmd.Flags().DurationP("until", "u", 0, "Total duration to run until (Time.Now + duration)")
	rootCmd.Flags().DurationP("pomo", "p", 25*time.Minute, "Pomodoro duration")
	rootCmd.Flags().DurationP("short", "s", 5*time.Minute, "Short break duration")
	rootCmd.Flags().DurationP("long", "l", 15*time.Minute, "Long break duration")

	viper.BindPFlag("db", rootCmd.Flags().Lookup("db"))
	viper.BindPFlag("pomo", rootCmd.Flags().Lookup("pomo"))
	viper.BindPFlag("short", rootCmd.Flags().Lookup("short"))
	viper.BindPFlag("long", rootCmd.Flags().Lookup("long"))
	viper.BindPFlag("until", rootCmd.Flags().Lookup("until"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pomodoro-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pomodoro-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func rootAction(writer io.Writer, config *pomodoro.IntervalConfig) error {
	a, err := app.New(config)
	if err != nil {
		return err
	}
	return a.Run()
}

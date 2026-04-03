package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:     "skipper",
	Version: "v1.0:beta",
	Short:   "skipper is a cli tool for managing ssh connections",
	Long: `skipper is a cli tool for managing ssh connections, It allows you to select your preferred ssh connection and execute commands on it. eg. skipper

Usage: skipper <command> [flags]

Flags:
  -c, --config string   path to ssh config file, defaults to ~/.ssh/config
  -h, --help            help for skipper
  -v, --version         print version information`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "~/.ssh/config", "path to ssh config file, defaults to ~/.ssh/config")
	rootCmd.Flags().BoolP("version", "v", false, "print version information")
}

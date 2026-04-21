package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <alias> <user@host[:port]>",
	Short: "Add a host entry to the SSH config",
	Long:  "Add a host entry to the SSH config under the given alias. The target must be in the form user@host[:port].",
	Example: `  skipper add devone user@10.0.0.8:9000
  skipper add bastion admin@10.0.0.5`,
	Args: cobra.ExactArgs(2),
	RunE: runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	path, err := resolveConfigPath(configPath)
	if err != nil {
		return err
	}

	host, created, err := addHost(path, args[0], args[1])
	if err != nil {
		return err
	}

	if created {
		fmt.Printf("added host %q for %s\n", host.Alias, hostTarget(host))
	} else {
		fmt.Printf("host %q already exists for %s, no change\n", host.Alias, hostTarget(host))
	}
	return nil
}

func init() {
	rootCmd.AddCommand(addCmd)
}

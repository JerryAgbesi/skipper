package cmd

import (
	"fmt"

	"github.com/jerryagbesi/skipper/internal/sshconfig"
	"github.com/jerryagbesi/skipper/internal/ui/addform"

	"github.com/spf13/cobra"
)

var addFormRunner = addform.Run

var addCmd = &cobra.Command{
	Use:   "add [<alias> <user@host[:port]>]",
	Short: "Add a host entry to the SSH config",
	Long: `Add a host entry to the SSH config.

Run with no arguments to launch an interactive form prompting for alias,
user, host name, and port. Alternatively pass an alias and a target in
the form user@host[:port] for a non-interactive add.`,
	Example: `  skipper add
  skipper add devone user@10.0.0.8:9000
  skipper add bastion admin@10.0.0.5`,
	Args: cobra.MaximumNArgs(2),
	RunE: runAdd,
}

func runAdd(_ *cobra.Command, args []string) error {
	path, err := resolveConfigPath(configPath)
	if err != nil {
		return err
	}

	switch len(args) {
	case 2:
		return addFromArgs(path, args[0], args[1])
	case 1:
		return fmt.Errorf("expected either no arguments or both <alias> and <user@host[:port]>")
	default:
		return addInteractive(path)
	}
}

func addFromArgs(path, alias, target string) error {
	host, created, err := addHost(path, alias, target)
	if err != nil {
		return err
	}
	printAddResult(host, created)
	return nil
}

func addInteractive(path string) error {
	result, err := addFormRunner(addform.Input{})
	if err != nil {
		return err
	}
	if result.Cancelled {
		return nil
	}

	host := sshconfig.Host{
		Alias:    result.Alias,
		Hostname: result.Hostname,
		User:     result.User,
		Port:     result.Port,
	}

	added, created, err := sshconfig.AddHost(path, host)
	if err != nil {
		return err
	}
	printAddResult(added, created)
	return nil
}

func printAddResult(host *sshconfig.Host, created bool) {
	if created {
		fmt.Printf("added host %q for %s\n", host.Alias, hostTarget(host))
		return
	}
	fmt.Printf("host %q already exists for %s, no change\n", host.Alias, hostTarget(host))
}

func init() {
	rootCmd.AddCommand(addCmd)
}

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jerryagbesi/skipper/internal/connect"
	"github.com/jerryagbesi/skipper/internal/sshconfig"
	"github.com/jerryagbesi/skipper/internal/ui"

	"github.com/spf13/cobra"
)

var (
	configPath string
	addAlias   string
	findQuery  string
)

var version = "dev"

var rootCmd = &cobra.Command{
	Use:           "skipper <command> [flags]",
	Version:       version,
	Short:         "skipper is a cli tool for managing ssh connections",
	Example:       "skipper --version",
	RunE:          runRoot,
	SilenceErrors: true,
	Long:          `skipper is a cli tool for managing ssh connections, It allows you to select your preferred ssh host alias, connect to it, and execute commands.`,
}

func runRoot(cmd *cobra.Command, args []string) error {
	path, err := resolveConfigPath(configPath)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("add") {
		host, err := addHost(path, addAlias, args)
		if err != nil {
			return err
		}

		fmt.Printf("added host %q for %s\n", host.Alias, hostTarget(host))
		return nil
	}

	hosts, err := sshconfig.ParseHosts(path)
	if err != nil {
		return err
	}

	if len(hosts) == 0 {
		return fmt.Errorf("no hosts found in config file")
	}

	options, hosts, err := prepareHostSelection(cmd, hosts)
	if err != nil {
		return err
	}

	result, err := ui.Run(hosts, options)
	if err != nil {
		return err
	}

	if result.Cancelled {
		return nil
	}

	return connect.Connect(result.Host, exec.Command)
}

func resolveConfigPath(path string) (string, error) {
	if path != "" {
		return path, nil
	}

	return sshconfig.DefaultConfigPath()
}

func prepareHostSelection(cmd *cobra.Command, hosts []sshconfig.Host) (ui.RunOptions, []sshconfig.Host, error) {
	options := ui.RunOptions{}
	if !cmd.Flags().Changed("find") {
		return options, hosts, nil
	}

	options.StartFiltering = findQuery == ""
	if findQuery == "" {
		return options, hosts, nil
	}

	filtered := filterHosts(hosts, findQuery)
	if len(filtered) == 0 {
		return ui.RunOptions{}, nil, fmt.Errorf("no hosts found matching %q", findQuery)
	}

	return options, filtered, nil
}

func addHost(path, alias string, args []string) (*sshconfig.Host, error) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return nil, fmt.Errorf("--add requires an alias")
	}

	if len(args) != 1 {
		return nil, fmt.Errorf("--add requires exactly one target in the format user@host[:port]")
	}

	host, err := connect.ParseTarget(args[0])
	if err != nil {
		return nil, err
	}

	host.Alias = alias
	return sshconfig.AddHost(path, *host)
}

func filterHosts(hosts []sshconfig.Host, query string) []sshconfig.Host {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return hosts
	}

	filtered := make([]sshconfig.Host, 0, len(hosts))
	for _, host := range hosts {
		if hostMatchesQuery(host, query) {
			filtered = append(filtered, host)
		}
	}

	return filtered
}

func hostMatchesQuery(host sshconfig.Host, query string) bool {
	fields := []string{host.Alias, host.Hostname, host.User, host.IdentityFile}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}

	return host.Port != 0 && strings.Contains(fmt.Sprintf("%d", host.Port), query)
}

func hostTarget(host *sshconfig.Host) string {
	target := host.User + "@" + host.Hostname
	if host.Port > 0 {
		return fmt.Sprintf("%s:%d", target, host.Port)
	}

	return target
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to ssh config file, defaults to ~/.ssh/config")
	rootCmd.Flags().StringVarP(&addAlias, "add", "a", "", "add a host alias using `<alias>` user@host[:port]")
	rootCmd.Flags().StringVarP(&findQuery, "find", "f", "", "start in find mode or pre-filter hosts by a search term")
	rootCmd.Flags().Lookup("find").NoOptDefVal = ""
	rootCmd.Flags().BoolP("version", "v", false, "print version information")
}

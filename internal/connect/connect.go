package connect

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/jerryagbesi/skipper/internal/sshconfig"
)

// Commander allows for mocking exec.Command.
type Commander func(name string, args ...string) *exec.Cmd

func Connect(host *sshconfig.Host, commander Commander) error {
	target := host.Alias

	var cmd *exec.Cmd

	switch {
	case target != "":
		cmd = commander("ssh", host.Alias)
	case host.Hostname != "":
		target = host.Hostname
		cmd = commander("ssh", host.User+"@"+host.Hostname, "-p", strconv.Itoa(host.Port))
	default:
		return fmt.Errorf("no target specified")
	}

	fmt.Printf("Connecting to %s\n", target)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

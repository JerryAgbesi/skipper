package connect

import (
	"os/exec"
	"testing"

	"github.com/jerryagbesi/skipper/internal/sshconfig"
)

func TestConnect_WithAlias_BuildsCorrectCommand(t *testing.T) {
	host := &sshconfig.Host{
		Alias:    "bastion",
		Hostname: "10.0.0.1",
	}

	var capturedName string
	var capturedArgs []string

	fakeCommander := func(name string, args ...string) *exec.Cmd {
		capturedName = name
		capturedArgs = args
		return exec.Command("true")
	}

	Connect(host, fakeCommander)

	if capturedName != "ssh" {
		t.Errorf("expected ssh, got %s", capturedName)
	}
	if capturedArgs[0] != "bastion" {
		t.Errorf("expected bastion, got %s", capturedArgs[0])
	}
}

func TestConnect_NoAlias_BuildsCorrectCommand(t *testing.T) {
	host := &sshconfig.Host{
		Alias:    "",
		Hostname: "10.0.0.1",
		User:     "ubuntu",
		Port:     2222,
	}

	var capturedArgs []string

	fakeCommander := func(_ string, args ...string) *exec.Cmd {
		capturedArgs = args
		return exec.Command("true")
	}

	Connect(host, fakeCommander)

	if capturedArgs[0] != "ubuntu@10.0.0.1" {
		t.Errorf("expected ubuntu@10.0.0.1, got %s", capturedArgs[0])
	}
	if capturedArgs[2] != "2222" {
		t.Errorf("expected port 2222, got %s", capturedArgs[2])
	}
}

func TestConnect_EmptyHost_ReturnsError(t *testing.T) {
	host := &sshconfig.Host{}

	fakeCommander := func(_ string, _ ...string) *exec.Cmd {
		t.Error("commander should not be called for empty host")
		return nil
	}

	err := Connect(host, fakeCommander)
	if err == nil {
		t.Error("expected error for empty host, got nil")
	}
}
